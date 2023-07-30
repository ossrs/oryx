package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ossrs/go-oryx-lib/logger"
)

func TestApi_QuerySecret(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	logger.Tf(ctx, "TestApi_QuerySecret with %v", options())

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done with err %+v", err)
		}
	}(ctx)
	defer cancel()

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
