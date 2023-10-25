//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	// From ossrs.
	"github.com/ossrs/go-oryx-lib/errors"
	ohttp "github.com/ossrs/go-oryx-lib/http"
	"github.com/ossrs/go-oryx-lib/logger"

	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
)

var callbackWorker *CallbackWorker

type CallbackWorker struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// The ephemeral callback config.
	ephemeralConfig CallbackConfig
	// Whether update the config immediately.
	updateConfig chan bool

	lock sync.Mutex
}

func NewCallbackWorker() *CallbackWorker {
	return &CallbackWorker{
		updateConfig: make(chan bool, 1),
	}
}

func (v *CallbackWorker) Handle(ctx context.Context, handler *http.ServeMux) error {
	ep := "/terraform/v1/mgmt/hooks/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				All   *bool   `json:"all"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			var conf CallbackConfig
			if err := conf.Load(ctx); err != nil {
				return errors.Wrapf(err, "load")
			}

			req, err := rdb.HGet(ctx, SRS_HOOKS, "req").Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v req", SRS_HOOKS)
			}

			res, err := rdb.HGet(ctx, SRS_HOOKS, "res").Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v res", SRS_HOOKS)
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Request  string `json:"req"`
				Response string `json:"res"`
				*CallbackConfig
			}{
				Request:        req,
				Response:       res,
				CallbackConfig: &conf,
			})
			logger.Tf(ctx, "hooks apply ok, %v, token=%vB", conf.String(), len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/mgmt/hooks/apply"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var conf CallbackConfig
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				*CallbackConfig
			}{
				Token:          &token,
				CallbackConfig: &conf,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if err := rdb.HSet(ctx, SRS_HOOKS, "target", conf.Target).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v target %v", SRS_HOOKS, conf.Target)
			}
			if err := rdb.HSet(ctx, SRS_HOOKS, "opaque", conf.Opaque).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v opaque %v", SRS_HOOKS, conf.Opaque)
			}
			if err := rdb.HSet(ctx, SRS_HOOKS, "all", fmt.Sprintf("%v", conf.All)).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v all %v", SRS_HOOKS, conf.All)
			}

			// Notify the callback worker to update the config.
			select {
			case v.updateConfig <- true:
			case <-ctx.Done():
			default:
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "hooks apply ok, %v, token=%vB", conf.String(), len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/mgmt/hooks/example"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			q := r.URL.Query()

			var fail bool
			if value := q.Get("fail"); value == "true" || value == "1" {
				fail = true
			}

			var action, opaque string
			var streamObj SrsStream
			if err := ParseBody(ctx, r.Body, &struct {
				Action *string `json:"action"`
				Opaque *string `json:"opaque"`
				*SrsStream
			}{
				Action:    &action,
				Opaque:    &opaque,
				SrsStream: &streamObj,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			if fail {
				return errors.Errorf("fail as required, action=%v, opaque=%v, %v",
					action, opaque, streamObj.String())
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "hooks example ok, action=%v, opaque=%v, %v", action, opaque, streamObj.String())
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	return nil
}

func (v *CallbackWorker) Close() error {
	if v.cancel != nil {
		v.cancel()
	}
	v.wg.Wait()
	return nil
}

func (v *CallbackWorker) Start(ctx context.Context) error {
	wg := &v.wg

	ctx, cancel := context.WithCancel(ctx)
	v.cancel = cancel

	ctx = logger.WithContext(ctx)
	logger.Tf(ctx, "callback start a worker")

	wg.Add(1)
	go func() {
		defer wg.Done()

		for ctx.Err() == nil {
			var conf CallbackConfig
			if err := conf.Load(ctx); err != nil {
				logger.Wf(ctx, "load config %v err %+v", conf, err)

				select {
				case <-ctx.Done():
				case <-time.After(10 * time.Second):
				}
				continue
			}

			func() {
				v.lock.Lock()
				defer v.lock.Unlock()
				v.ephemeralConfig = conf
			}()

			select {
			case <-ctx.Done():
			case <-time.After(time.Second * 3):
			case <-v.updateConfig:
			}
		}
	}()

	return nil
}

func (v *CallbackWorker) OnMessage(ctx context.Context, action SrsAction, streamObj *SrsStream) error {
	if action != SrsActionOnPublish && action != SrsActionOnUnpublish {
		return nil
	}

	var conf CallbackConfig
	func() {
		v.lock.Lock()
		defer v.lock.Unlock()
		conf = v.ephemeralConfig
	}()

	if !conf.All || conf.Target == "" {
		return nil
	}

	req := &struct {
		Action string `json:"action"`
		Opaque string `json:"opaque"`
		*SrsStream
	}{
		Action:    string(action),
		Opaque:    conf.Opaque,
		SrsStream: streamObj,
	}

	pfn4 := func(b, b2 []byte, code int) error {
		if code != 0 {
			return errors.Errorf("response code %v", code)
		}

		logger.Tf(ctx, "callback ok, post %v with %s, response %v", conf.String(), string(b), string(b2))
		return nil
	}

	pfn3 := func(b, b2 []byte) error {
		if code, err := strconv.ParseInt(string(b2), 10, 64); err == nil {
			return pfn4(b, b2, int(code))
		}

		var code int
		if err := json.Unmarshal(b2, &struct {
			Code *int `json:"code"`
		}{
			Code: &code,
		}); err != nil {
			return errors.Wrapf(err, "unmarshal response")
		}
		return pfn4(b, b2, code)
	}

	pfn2 := func(b []byte) error {
		res, err := http.Post(conf.Target, "application/json", bytes.NewReader(b))
		if err != nil {
			return errors.Wrapf(err, "http post")
		}
		if res.StatusCode != http.StatusOK {
			return errors.Errorf("response status %v", res.StatusCode)
		}
		defer res.Body.Close()

		b2, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrapf(err, "read body")
		}

		if err := rdb.HSet(ctx, SRS_HOOKS, "res", string(b2)).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hset %v res %v", SRS_HOOKS, string(b2))
		}

		if err := pfn3(b, b2); err != nil {
			return errors.Wrapf(err, "res body %v", string(b2))
		}

		return nil
	}

	pfn := func() error {
		b, err := json.Marshal(req)
		if err != nil {
			return errors.Wrapf(err, "marshal req")
		}

		if err := rdb.HSet(ctx, SRS_HOOKS, "req", string(b)).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hset %v req %v", SRS_HOOKS, string(b))
		}

		if err := pfn2(b); err != nil {
			return errors.Wrapf(err, "post with %s", string(b))
		}

		return nil
	}

	if err := pfn(); err != nil {
		return errors.Wrapf(err, "callback with conf %v, req %v", conf.String(), req)
	}
	return nil
}

type CallbackConfig struct {
	// The callback target.
	Target string `json:"target"`
	// The opaque string, for example, the token.
	Opaque string `json:"opaque"`
	// Whether to callback all streams.
	All bool `json:"all"`
}

func (v CallbackConfig) String() string {
	return fmt.Sprintf("target=%v, opaque=%v, all=%v", v.Target, v.Opaque, v.All)
}

func (v *CallbackConfig) Load(ctx context.Context) (err error) {
	if v.Target, err = rdb.HGet(ctx, SRS_HOOKS, "target").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v target", SRS_HOOKS)
	}

	if v.Opaque, err = rdb.HGet(ctx, SRS_HOOKS, "opaque").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v opaque", SRS_HOOKS)
	}

	if all, err := rdb.HGet(ctx, SRS_HOOKS, "all").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v all", SRS_HOOKS)
	} else if all == "true" {
		v.All = true
	}

	return nil
}
