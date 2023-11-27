package main

import (
	"context"
	"testing"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
)

func TestOpenAI_TranscriptCheckConnection(t *testing.T) {
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

	// Ignore the test case if api secret key not set.
	apiKey, baseUrl := OpenAIConfig()
	if apiKey == "" {
		return
	}

	if r0 = NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/check", &struct {
		SecretKey string `json:"secretKey"`
		BaseUrl   string `json:"baseUrl"`
	}{
		SecretKey: apiKey, BaseUrl: baseUrl,
	}, nil); r0 != nil {
		return
	}
}

func TestOpenAI_TranscriptApplyQuery(t *testing.T) {
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

	// Ignore the test case if api secret key not set.
	apiKey, baseUrl := OpenAIConfig()
	if apiKey == "" {
		return
	}

	type TranscriptConfig struct {
		All       bool   `json:"all"`
		SecretKey string `json:"secretKey"`
		BaseURL   string `json:"baseURL"`
		Language  string `json:"lang"`
	}
	var conf TranscriptConfig
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/query", nil, &struct {
		Config *TranscriptConfig `json:"config"`
	}{
		Config: &conf,
	}); err != nil {
		r0 = errors.Wrapf(err, "request query failed")
		return
	}

	// Restore the state of transcode.
	backup := conf
	defer func() {
		logger.Tf(ctx, "restore config %v", backup)

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/apply", backup, nil)
	}()

	// Enable transcript.
	conf.All, conf.SecretKey, conf.BaseURL, conf.Language = true, apiKey, baseUrl, "en"
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/apply", conf, nil); err != nil {
		r0 = errors.Wrapf(err, "request apply failed")
		return
	}

	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/query", nil, &struct {
		Config *TranscriptConfig `json:"config"`
	}{
		Config: &conf,
	}); err != nil {
		r0 = errors.Wrapf(err, "request query failed")
		return
	} else if !conf.All || conf.SecretKey != apiKey || conf.BaseURL != baseUrl || conf.Language != "en" {
		r0 = errors.Errorf("invalid config %+v", conf)
		return
	}

	// Disable transcript.
	conf.All, conf.SecretKey, conf.BaseURL, conf.Language = false, apiKey, baseUrl, "en"
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/apply", conf, nil); err != nil {
		r0 = errors.Wrapf(err, "request apply failed")
		return
	}

	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/query", nil, &struct {
		Config *TranscriptConfig `json:"config"`
	}{
		Config: &conf,
	}); err != nil {
		r0 = errors.Wrapf(err, "request query failed")
		return
	} else if conf.All || conf.SecretKey != apiKey || conf.BaseURL != baseUrl || conf.Language != "en" {
		r0 = errors.Errorf("invalid config %+v", conf)
		return
	}
}
