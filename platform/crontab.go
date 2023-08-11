//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: MIT
//
package main

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
)

var crontabWorker *CrontabWorker

type CrontabWorker struct {
	wg sync.WaitGroup
}

func NewCrontabWorker() *CrontabWorker {
	return &CrontabWorker{}
}

func (v *CrontabWorker) Close() error {
	v.wg.Wait()
	return nil
}

func (v *CrontabWorker) Start(ctx context.Context) error {
	v.wg.Add(1)
	go func() {
		defer v.wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(24*3600) * time.Second):
			}

			logger.Tf(ctx, "crontab: start to query latest version")
			if versions, err := queryLatestVersion(ctx); err != nil {
				logger.Wf(ctx, "crontab: ignore err %v", err)
			} else {
				logger.Tf(ctx, "crontab: query version ok, result is %v", versions.String())
			}
		}
	}()

	v.wg.Add(1)
	go func() {
		defer v.wg.Done()

		for {
			logger.Tf(ctx, "crontab: start to refresh ssl cert")
			if err := refreshSSLCert(ctx); err != nil {
				logger.Wf(ctx, "crontab: ignore err %v", err)
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(24*3600) * time.Second):
			}
		}
	}()

	return nil
}

func refreshSSLCert(ctx context.Context) error {
	provider, err := rdb.Get(ctx, SRS_HTTPS).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if provider != "lets" {
		logger.Tf(ctx, "crontab: ignore ssl provider %v", provider)
		return nil
	}

	domain, err := rdb.Get(ctx, SRS_HTTPS_DOMAIN).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if domain == "" {
		logger.Tf(ctx, "crontab: ignore ssl domain empty")
		return nil
	}

	if err := renewLetsEncrypt(ctx, domain); err != nil {
		return err
	} else {
		logger.Tf(ctx, "crontab: renew ssl cert ok")
	}

	if err := nginxGenerateConfig(ctx); err != nil {
		return errors.Wrapf(err, "nginx config and reload")
	}

	logger.Tf(ctx, "crontab: refresh ssl cert ok")
	return nil
}
