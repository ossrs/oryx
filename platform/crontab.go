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
			if err := certManager.refreshSSLCert(ctx); err != nil {
				logger.Wf(ctx, "crontab: ignore err %v", err)
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(24*3600) * time.Second):
			}
		}
	}()

	v.wg.Add(1)
	go func() {
		defer v.wg.Done()

		for {
			logger.Tf(ctx, "crontab: start to refresh certificate file")
			if err := certManager.reloadCertificateFile(ctx); err != nil {
				logger.Wf(ctx, "crontab: ignore err %v", err)
			}

			select {
			case <-ctx.Done():
				return
			case <-certManager.httpCertificateReload:
			case <-time.After(time.Duration(1*3600) * time.Second):
			}
		}
	}()

	return nil
}
