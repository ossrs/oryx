package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
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

func TestSystem_WithStream_PublishRtmpKickoff(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsLongTimeout)*time.Millisecond)
	defer cancel()

	if *noMediaTest {
		return
	}

	var r0, r1, r2, r3, r4, r5 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1, r2, r3, r4, r5); err != nil {
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

	var wg sync.WaitGroup
	defer wg.Wait()

	// Start FFmpeg to publish stream.
	streamID := fmt.Sprintf("stream-%v-%v", os.Getpid(), rand.Int())
	streamURL := fmt.Sprintf("%v/live/%v?secret=%v", *endpointRTMP, streamID, pubSecret)
	ffmpeg := NewFFmpeg(func(v *ffmpegClient) {
		v.args = []string{
			"-re", "-stream_loop", "-1", "-i", *srsInputFile, "-c", "copy",
			"-f", "flv", streamURL,
		}
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r1 = ffmpeg.Run(ctx, cancel)
	}()

	// Wait for stream to be ready.
	select {
	case <-ctx.Done():
	case <-time.After(5 * time.Second):
	}

	// Query all streams.
	type StreamQueryResult struct {
		Vhost  string `json:"vhost"`
		App    string `json:"app"`
		Stream string `json:"stream"`
	}
	var streams []StreamQueryResult
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/streams/query", nil, &struct {
		Streams *[]StreamQueryResult `json:"streams"`
	}{
		Streams: &streams,
	}); err != nil {
		r0 = err
		return
	}

	// Find the target stream.
	var stream *StreamQueryResult
	for _, s := range streams {
		if s.Stream == streamID {
			stream = &s
			break
		}
	}
	if stream == nil {
		r0 = errors.Errorf("stream %v not found", streamID)
		return
	}

	// Kickoff the stream.
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/streams/kickoff", stream, nil); err != nil {
		r0 = err
		return
	}

	streams = nil
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/streams/query", nil, &struct {
		Streams *[]StreamQueryResult `json:"streams"`
	}{
		Streams: &streams,
	}); err != nil {
		r0 = err
		return
	}

	for _, s := range streams {
		if s.Stream == streamID {
			r0 = errors.Errorf("stream %v should be kicked off", streamID)
			return
		}
	}

	time.Sleep(3 * time.Second)
	logger.Tf(ctx, "kickoff ok")
	cancel()
}
