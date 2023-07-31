package main

import (
	"context"
	"testing"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
)

func TestApi_Empty(t *testing.T) {
	ctx := logger.WithContext(context.Background())
	logger.Tf(ctx, "test done")
}

func TestApi_Ready(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	if err := waitForServiceReady(ctx); err != nil {
		t.Errorf("Fail for err %+v", err)
	} else {
		logger.Tf(ctx, "test done")
	}
}

func TestApi_QueryPublishSecret(t *testing.T) {
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
	if err := apiRequest(ctx, "/terraform/v1/hooks/srs/secret/query", nil, &struct {
		Publish *string `json:"publish"`
	}{
		Publish: &publishSecret,
	}); err != nil {
		r0 = err
	} else if publishSecret == "" {
		r0 = errors.New("empty publish secret")
	}
}

func TestApi_LoginByPassword(t *testing.T) {
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
	if err := apiRequest(ctx, "/terraform/v1/mgmt/login", &struct {
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

func TestApi_BootstrapQueryEnvs(t *testing.T) {
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
		Secret     bool `json:"secret"`
		HTTPS      bool `json:"https"`
		MgmtDocker bool `json:"mgmtDocker"`
	}{}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/envs", nil, &res); err != nil {
		r0 = err
	} else if !res.Secret || !res.HTTPS || !res.MgmtDocker {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestApi_BootstrapQueryInit(t *testing.T) {
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
	if err := apiRequest(ctx, "/terraform/v1/mgmt/init", nil, &res); err != nil {
		r0 = err
	} else if !res.Init {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestApi_BootstrapQueryCheck(t *testing.T) {
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
	if err := apiRequest(ctx, "/terraform/v1/mgmt/check", nil, &res); err != nil {
		r0 = err
	} else if res.Upgrading {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestApi_BootstrapQueryVersions(t *testing.T) {
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
	if err := apiRequest(ctx, "/terraform/v1/mgmt/versions", nil, &res); err != nil {
		r0 = err
	} else if res.Version == "" {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestApi_SetupWebsiteFooter(t *testing.T) {
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

	req := struct {
		Beian string `json:"beian"`
		Text  string `json:"text"`
	}{
		Beian: "icp",
		Text:  "TestFooter",
	}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/beian/update", &req, nil); err != nil {
		r0 = err
	}

	res := struct {
		ICP string `json:"icp"`
	}{}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/beian/query", nil, &res); err != nil {
		r0 = err
	} else if res.ICP != "TestFooter" {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestApi_SetupWebsiteTitle(t *testing.T) {
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

	req := struct {
		Beian string `json:"beian"`
		Text  string `json:"text"`
	}{
		Beian: "title",
		Text:  "TestTitle",
	}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/beian/update", &req, nil); err != nil {
		r0 = err
	}

	res := struct {
		Title string `json:"title"`
	}{}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/beian/query", nil, &res); err != nil {
		r0 = err
	} else if res.Title != "TestTitle" {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestApi_SetupPublishSecret(t *testing.T) {
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

	req := struct {
		Secret string `json:"secret"`
	}{
		Secret: "TestPublish",
	}
	if err := apiRequest(ctx, "/terraform/v1/hooks/srs/secret/update", &req, nil); err != nil {
		r0 = err
	}

	res := struct {
		Publish string `json:"publish"`
	}{}
	if err := apiRequest(ctx, "/terraform/v1/hooks/srs/secret/query", nil, &res); err != nil {
		r0 = err
	} else if res.Publish != "TestPublish" {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestApi_TutorialsQueryBilibili(t *testing.T) {
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

	req := struct {
		BVID string `json:"bvid"`
	}{
		BVID: "BV1844y1L7dL",
	}
	res := struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
	}{}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/bilibili", &req, &res); err != nil {
		r0 = err
	} else if res.Title == "" || res.Desc == "" {
		r0 = errors.Errorf("invalid response %v", res)
	}
}
