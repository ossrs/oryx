package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
)

func TestSystem_Empty(t *testing.T) {
	ctx := logger.WithContext(context.Background())
	logger.Tf(ctx, "test done")
}

func TestSystem_Ready(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	if err := waitForServiceReady(ctx); err != nil {
		t.Errorf("Fail for err %+v", err)
	} else {
		logger.Tf(ctx, "test done")
	}
}

func TestSystem_QueryPublishSecret(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	var publishSecret string
	if err := NewApi().WithAuth(ctx, "/terraform/v1/hooks/srs/secret/query", nil, &struct {
		Publish *string `json:"publish"`
	}{
		Publish: &publishSecret,
	}); err != nil {
		r0 = err
	} else if publishSecret == "" {
		r0 = errors.New("empty publish secret")
	}
}

func TestSystem_LoginByPassword(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	var token string
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/login", &struct {
		Password string `json:"password"`
	}{
		Password: *systemPassword,
	}, &struct {
		Token *string `json:"token"`
	}{
		Token: &token,
	}); err != nil {
		r0 = err
	} else if token == "" {
		r0 = errors.New("empty token")
	}
}

func TestSystem_BootstrapQueryEnvs(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	res := struct {
		MgmtDocker bool `json:"mgmtDocker"`
	}{}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/envs", nil, &res); err != nil {
		r0 = err
	} else if !res.MgmtDocker {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestSystem_BootstrapQueryInit(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	res := struct {
		Init bool `json:"init"`
	}{}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/init", nil, &res); err != nil {
		r0 = err
	} else if !res.Init {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestSystem_BootstrapQueryCheck(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	res := struct {
		Upgrading bool `json:"upgrading"`
	}{}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/check", nil, &res); err != nil {
		r0 = err
	} else if res.Upgrading {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestSystem_BootstrapQueryVersions(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	res := struct {
		Version string `json:"version"`
	}{}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/versions", nil, &res); err != nil {
		r0 = err
	} else if res.Version == "" {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestSystem_CandidateByEip(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	var pubSecret string
	if err := NewApi().WithAuth(ctx, "/terraform/v1/hooks/srs/secret/query", nil, &struct {
		Publish *string `json:"publish"`
	}{
		Publish: &pubSecret,
	}); err != nil {
		r0 = err
		return
	}

	// Set the candidate by ?eip=
	offer := strings.ReplaceAll(SrsLarixExampleOffer, "\n", "\r\n")
	streamID := fmt.Sprintf("stream-%v-%v", os.Getpid(), rand.Int())
	eip := "12.34.56.78"
	var answer string
	if err := NewApi().NoAuth(ctx, fmt.Sprintf("/rtc/v1/whip/?app=live&stream=%v&secret=%v&eip=%v", streamID, pubSecret, eip), offer, &answer); err != nil {
		r0 = errors.Wrapf(err, "should ok for rtc publish api")
		return
	} else if line := StringContainsLine(answer, func(l string) bool {
		return strings.Contains(l, eip) && strings.Contains(l, "udp")
	}); line == "" {
		r0 = errors.Errorf("invalid answer %v", answer)
		return
	}
}

func TestSystem_CandidateByHostIp(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	var pubSecret string
	if err := NewApi().WithAuth(ctx, "/terraform/v1/hooks/srs/secret/query", nil, &struct {
		Publish *string `json:"publish"`
	}{
		Publish: &pubSecret,
	}); err != nil {
		r0 = err
		return
	}

	// Set the eip by HTTP header Host.
	offer := strings.ReplaceAll(SrsLarixExampleOffer, "\n", "\r\n")
	streamID := fmt.Sprintf("stream-%v-%v", os.Getpid(), rand.Int())
	eip := "12.34.56.78"
	var answer string
	if err := NewApi(func(v *testApi) {
		v.InjectRequest = func(req *http.Request) {
			req.Header.Set("X-Real-Host", eip)
		}
	}).NoAuth(ctx, fmt.Sprintf("/rtc/v1/whip/?app=live&stream=%v&secret=%v", streamID, pubSecret), offer, &answer); err != nil {
		r0 = errors.Wrapf(err, "should ok for rtc publish api")
		return
	} else if line := StringContainsLine(answer, func(l string) bool {
		return strings.Contains(l, eip) && strings.Contains(l, "udp")
	}); line == "" {
		r0 = errors.Errorf("invalid answer %v", answer)
		return
	}
}
