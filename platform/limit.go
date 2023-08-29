//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ossrs/go-oryx-lib/errors"
	"strconv"
	"sync"
	"time"
)

var limitWorker *LimitWorker

type LimitWorker struct {
	// The limit of total duration in seconds for each month.
	TotalDuration int64
	// The date to start to apply the limit.
	StartDate time.Time
	// The lastest date to update the limit.
	LastUpdateDate time.Time
	// The consumed duration in seconds for current month.
	ConsumedDuration int64

	lock sync.Mutex
	wg   sync.WaitGroup
}

func NewLimitWorker() *LimitWorker {
	return &LimitWorker{}
}

func (v *LimitWorker) Close() error {
	v.wg.Wait()
	return nil
}

func (v *LimitWorker) Start(ctx context.Context) error {
	if err := v.Parse(ctx); err != nil {
		return errors.Wrapf(err, "parse limit")
	}

	return nil
}

func (v *LimitWorker) Parse(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	r0, err := rdb.HGetAll(ctx, SRS_LIMIT).Result()
	if err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hgetall %v", SRS_BEIAN)
	}

	for k, value := range r0 {
		if k == "totalDuration" {
			if v.TotalDuration, err = strconv.ParseInt(value, 10, 64); err != nil {
				return errors.Wrapf(err, "parse limit %v", value)
			}
		} else if k == "startDate" {
			if v.StartDate, err = time.Parse(time.RFC3339, value); err != nil {
				return errors.Wrapf(err, "parse limit %v", value)
			}
		} else if k == "lastUpdateDate" {
			if v.LastUpdateDate, err = time.Parse(time.RFC3339, value); err != nil {
				return errors.Wrapf(err, "parse limit %v", value)
			}
		} else if k == "consumedDuration" {
			if v.ConsumedDuration, err = strconv.ParseInt(value, 10, 64); err != nil {
				return errors.Wrapf(err, "parse limit %v", value)
			}
		}
	}
	return nil
}

func (v *LimitWorker) Flush(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	v.LastUpdateDate = time.Now()
	if v.StartDate.IsZero() {
		v.StartDate = time.Now()
	}

	if err := rdb.HSet(ctx, SRS_LIMIT, "totalDuration", fmt.Sprintf("%v", v.TotalDuration)).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_LIMIT, "totalDuration", v.TotalDuration)
	}
	if err := rdb.HSet(ctx, SRS_LIMIT, "startDate", v.StartDate.Format(time.RFC3339)).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_LIMIT, "startDate", v.StartDate)
	}
	if err := rdb.HSet(ctx, SRS_LIMIT, "lastUpdateDate", v.LastUpdateDate.Format(time.RFC3339)).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_LIMIT, "lastUpdateDate", v.LastUpdateDate)
	}
	if err := rdb.HSet(ctx, SRS_LIMIT, "consumedDuration", fmt.Sprintf("%v", v.ConsumedDuration)).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_LIMIT, "consumedDuration", v.ConsumedDuration)
	}
	return nil
}
