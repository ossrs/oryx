package main

import (
	"context"
	"testing"
	"time"

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
