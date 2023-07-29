//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: MIT
//
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	ohttp "github.com/ossrs/go-oryx-lib/http"
	"github.com/ossrs/go-oryx-lib/logger"

	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
)

func NewDockerHTTPService() HttpService {
	return &dockerHTTPService{}
}

type dockerHTTPService struct {
	servers []*http.Server
}

func (v *dockerHTTPService) Close() error {
	servers := v.servers
	v.servers = nil

	for _, server := range servers {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Tf(ctx, "ignore HTTP server shutdown err %v", err)
		}
	}

	return nil
}

func (v *dockerHTTPService) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	handler := http.NewServeMux()
	if err := handleDockerHTTPService(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle service")
	}

	var r0 error
	if true {
		addr := os.Getenv("PLATFORM_LISTEN")
		if !strings.HasPrefix(addr, ":") {
			addr = fmt.Sprintf(":%v", addr)
		}
		logger.Tf(ctx, "HTTP listen at %v", addr)

		server := &http.Server{Addr: addr, Handler: handler}
		v.servers = append(v.servers, server)

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			logger.Tf(ctx, "shutting down HTTP server...")
			v.Close()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := server.ListenAndServe(); err != nil && ctx.Err() != context.Canceled {
				r0 = errors.Wrapf(err, "listen %v", addr)
			}
			logger.Tf(ctx, "HTTP server is done")
		}()
	}

	var r1 error
	if true {
		addr := os.Getenv("MGMT_LISTEN")
		if !strings.HasPrefix(addr, ":") {
			addr = fmt.Sprintf(":%v", addr)
		}
		logger.Tf(ctx, "HTTP listen at %v", addr)

		server := &http.Server{Addr: addr, Handler: handler}
		v.servers = append(v.servers, server)

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			logger.Tf(ctx, "shutting down HTTP server...")
			v.Close()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := server.ListenAndServe(); err != nil && ctx.Err() != context.Canceled {
				r1 = errors.Wrapf(err, "listen %v", addr)
			}
			logger.Tf(ctx, "HTTP server is done")
		}()
	}

	for _, r := range []error{r0, r1} {
		if r != nil {
			return r
		}
	}
	return nil
}

func handleDockerHTTPService(ctx context.Context, handler *http.ServeMux) error {
	ohttp.Server = fmt.Sprintf("srs-cloud/%v", version)

	ep := "/terraform/v1/host/versions"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		ohttp.WriteData(ctx, w, r, &struct {
			Version string `json:"version"`
		}{
			Version: strings.TrimPrefix(version, "v"),
		})
	})

	ep = "/terraform/v1/mgmt/versions"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		ohttp.WriteData(ctx, w, r, &struct {
			Version string `json:"version"`
		}{
			Version: strings.TrimPrefix(version, "v"),
		})
	})

	if err := handleDockerHooksService(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle hooks")
	}

	ep = "/terraform/v1/ffmpeg/versions"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		ohttp.WriteData(ctx, w, r, &struct {
			Version string `json:"version"`
		}{
			Version: strings.TrimPrefix(version, "v"),
		})
	})

	if err := forwardWorker.Handle(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle forward")
	}

	if err := vLiveWorker.Handle(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle vLive")
	}

	ep = "/terraform/v1/mgmt/init"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var password string
			if len(b) > 0 {
				if err := json.Unmarshal(b, &struct {
					Password *string `json:"password"`
				}{
					Password: &password,
				}); err != nil {
					return errors.Wrapf(err, "json unmarshal %v", string(b))
				}
			}

			// If no password, query the system init status.
			if password == "" {
				ohttp.WriteData(ctx, w, r, &struct {
					Init bool `json:"init"`
				}{
					Init: os.Getenv("MGMT_PASSWORD") != "",
				})
				return nil
			}

			// If already initialized, never set it again.
			if os.Getenv("MGMT_PASSWORD") != "" {
				return errors.New("already initialized")
			}

			// Initialize the system password, save to env.
			envFile := path.Join(conf.Pwd, "containers/data/config/.env")
			if envs, err := godotenv.Read(envFile); err != nil {
				return errors.Wrapf(err, "load envs from %v", envFile)
			} else {
				envs["MGMT_PASSWORD"] = password
				if err := godotenv.Write(envs, envFile); err != nil {
					return errors.Wrapf(err, "write %v", envFile)
				}
			}
			logger.Tf(ctx, "init mgmt password %vB ok, file=%v", len(password), envFile)

			// Refresh the local token.
			if err := godotenv.Load(envFile); err != nil {
				return errors.Wrapf(err, "load %v", envFile)
			}

			expireAt, createAt, token, err := createToken(ctx, os.Getenv("SRS_PLATFORM_SECRET"))
			if err != nil {
				return errors.Wrapf(err, "build token")
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Token    string `json:"token"`
				CreateAt string `json:"createAt"`
				ExpireAt string `json:"expireAt"`
			}{
				Token: token, CreateAt: createAt.Format(time.RFC3339), ExpireAt: expireAt.Format(time.RFC3339),
			})
			logger.Tf(ctx, "init password ok, create=%v, expire=%v, password=%vB", createAt, expireAt, len(password))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/mgmt/check"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			// Check whether redis is ok.
			if r0, err := rdb.HGet(ctx, SRS_AUTH_SECRET, "pubSecret").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v pubSecret", SRS_AUTH_SECRET)
			} else if r1, err := rdb.HLen(ctx, SRS_FIRST_BOOT).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "get %v", SRS_FIRST_BOOT)
			} else if r2, err := rdb.HLen(ctx, SRS_TENCENT_LH).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "get %v", SRS_TENCENT_LH)
			} else if r0 == "" || r1 <= 0 || r2 <= 0 {
				return errors.New("Redis is not  ready")
			} else {
				logger.Tf(ctx, "system check ok, r0=%v, r1=%v, r2=%v", r0, r1, r2)
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Upgrading bool `json:"upgrading"`
			}{
				Upgrading: false,
			})
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/mgmt/envs"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			r0, err := rdb.HGet(ctx, SRS_PLATFORM_SECRET, "token").Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v token", SRS_PLATFORM_SECRET)
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Secret     bool   `json:"secret"`
				HTTPS      string `json:"https"`
				MgmtDocker bool   `json:"mgmtDocker"`
			}{
				Secret:     r0 != "",
				HTTPS:      os.Getenv("SRS_HTTPS"),
				MgmtDocker: true,
			})
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/mgmt/token"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var token string
			if err := json.Unmarshal(b, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "json unmarshal %v", string(b))
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			// Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
			// See https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
			if _, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(apiSecret), nil
			}); err != nil {
				return errors.Wrapf(err, "verify token %v", token)
			}

			expireAt, createAt, token, err := createToken(ctx, os.Getenv("SRS_PLATFORM_SECRET"))
			if err != nil {
				return errors.Wrapf(err, "build token")
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Token    string `json:"token"`
				CreateAt string `json:"createAt"`
				ExpireAt string `json:"expireAt"`
			}{
				Token: token, CreateAt: createAt.Format(time.RFC3339), ExpireAt: expireAt.Format(time.RFC3339),
			})
			logger.Tf(ctx, "login by token ok, create=%v, expire=%v, token=%vB", createAt, expireAt, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/mgmt/secret/token"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var apiSecret string
			if err := json.Unmarshal(b, &struct {
				ApiSecret *string `json:"apiSecret"`
			}{
				ApiSecret: &apiSecret,
			}); err != nil {
				return errors.Wrapf(err, "json unmarshal %v", string(b))
			}

			if apiSecret != os.Getenv("SRS_PLATFORM_SECRET") {
				return errors.New("apiSecret verify failed")
			}

			expireAt, createAt, token, err := createToken(ctx, apiSecret)
			if err != nil {
				return errors.Wrapf(err, "build token")
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Token    string `json:"token"`
				CreateAt string `json:"createAt"`
				ExpireAt string `json:"expireAt"`
			}{
				Token: token, CreateAt: createAt.Format(time.RFC3339), ExpireAt: expireAt.Format(time.RFC3339),
			})
			logger.Tf(ctx, "create token by apiSecret ok, create=%v, expire=%v, token=%vB", createAt, expireAt, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/mgmt/login"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			if os.Getenv("MGMT_PASSWORD") == "" {
				return errors.New("not init")
			}

			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var password string
			if err := json.Unmarshal(b, &struct {
				Password *string `json:"password"`
			}{
				Password: &password,
			}); err != nil {
				return errors.Wrapf(err, "json unmarshal %v", string(b))
			}

			if password == "" {
				return errors.New("no password")
			}

			if password != os.Getenv("MGMT_PASSWORD") {
				return errors.New("invalid password")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			expireAt, createAt, token, err := createToken(ctx, apiSecret)
			if err != nil {
				return errors.Wrapf(err, "build token")
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Token    string `json:"token"`
				CreateAt string `json:"createAt"`
				ExpireAt string `json:"expireAt"`
			}{
				Token: token, CreateAt: createAt.Format(time.RFC3339), ExpireAt: expireAt.Format(time.RFC3339),
			})
			logger.Tf(ctx, "login by password ok, create=%v, expire=%v, token=%vB", createAt, expireAt, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/mgmt/status"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var token string
			if err := json.Unmarshal(b, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "json unmarshal %v", string(b))
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			// Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
			// See https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
			if _, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(apiSecret), nil
			}); err != nil {
				return errors.Wrapf(err, "verify token %v", token)
			}

			upgrading, err := rdb.HGet(ctx, SRS_UPGRADING, "upgrading").Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v upgrading", SRS_UPGRADING)
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Version   string   `json:"version"`
				Releases  Versions `json:"releases"`
				Upgrading bool     `json:"upgrading"`
				Strategy  string   `json:"strategy"`
			}{
				Version:   conf.Versions.Version,
				Releases:  conf.Versions,
				Upgrading: upgrading == "1",
				Strategy:  "manual",
			})
			logger.Tf(ctx, "status ok, versions=%v, upgrading=%v, token=%vB", conf.Versions.String(), upgrading, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/mgmt/bilibili"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var token, bvid string
			if err := json.Unmarshal(b, &struct {
				Token *string `json:"token"`
				BVID  *string `json:"bvid"`
			}{
				Token: &token, BVID: &bvid,
			}); err != nil {
				return errors.Wrapf(err, "json unmarshal %v", string(b))
			}
			if bvid == "" {
				return errors.New("no bvid")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			// Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
			// See https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
			if _, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(apiSecret), nil
			}); err != nil {
				return errors.Wrapf(err, "verify token %v", token)
			}

			bilibiliObj := struct {
				Update string                 `json:"update"`
				Res    map[string]interface{} `json:"res"`
			}{}
			if bilibili, err := rdb.HGet(ctx, SRS_CACHE_BILIBILI, bvid).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_CACHE_BILIBILI, bvid)
			} else if bilibili != "" {
				if err := json.Unmarshal([]byte(bilibili), &bilibiliObj); err != nil {
					return errors.Wrapf(err, "json unmarshal %v", bilibili)
				}
			}

			var cacheExpired bool
			if bilibiliObj.Update != "" {
				duration := time.Duration(24*3600) * time.Second
				if os.Getenv("NODE_ENV") == "development" {
					duration = time.Duration(300) * time.Second
				}

				updateAt, err := time.Parse(time.RFC3339, bilibiliObj.Update)
				if err != nil {
					cacheExpired = true
				}
				if updateAt.Add(duration).Before(time.Now()) {
					cacheExpired = true
				}
			}

			if bilibiliObj.Res == nil || cacheExpired {
				bilibiliObj.Update = time.Now().Format(time.RFC3339)

				bilibiliURL := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%v", bvid)
				res, err := http.Get(bilibiliURL)
				if err != nil {
					return errors.Wrapf(err, "get %v", bilibiliURL)
				}
				defer res.Body.Close()

				b, err := ioutil.ReadAll(res.Body)
				if err != nil {
					return errors.Wrapf(err, "read %v", bilibiliURL)
				}

				if err := json.Unmarshal(b, &struct {
					Code    int                     `json:"code"`
					Message string                  `json:"message"`
					TTL     int                     `json:"ttl"`
					Data    *map[string]interface{} `json:"data"`
				}{
					Data: &bilibiliObj.Res,
				}); err != nil {
					return errors.Wrapf(err, "json unmarshal %v", string(b))
				}
			}
			if b, err := json.Marshal(bilibiliObj); err != nil {
				return errors.Wrapf(err, "json marshal %v", bilibiliObj)
			} else if err = rdb.HSet(ctx, SRS_CACHE_BILIBILI, bvid, string(b)).Err(); err != nil {
				return errors.Wrapf(err, "update redis for %v", string(b))
			}

			ohttp.WriteData(ctx, w, r, bilibiliObj.Res)
			logger.Tf(ctx, "bilibili cache bvid=%v, update=%v, token=%vB", bvid, bilibiliObj.Update, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/mgmt/beian/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			r0, err := rdb.HGetAll(ctx, SRS_BEIAN).Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hgetall %v", SRS_BEIAN)
			}

			ohttp.WriteData(ctx, w, r, r0)
			logger.Tf(ctx, "beian: query ok")
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/mgmt/secret/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var token string
			if err := json.Unmarshal(b, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "json unmarshal %v", string(b))
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			// Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
			// See https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
			if _, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(apiSecret), nil
			}); err != nil {
				return errors.Wrapf(err, "verify token %v", token)
			}

			ohttp.WriteData(ctx, w, r, apiSecret)
			logger.Tf(ctx, "query apiSecret ok, versions=%v, token=%vB", conf.Versions.String(), len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/mgmt/beian/update"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var token, beian, text string
			if err := json.Unmarshal(b, &struct {
				Token *string `json:"token"`
				Beian *string `json:"beian"`
				Text  *string `json:"text"`
			}{
				Token: &token, Beian: &beian, Text: &text,
			}); err != nil {
				return errors.Wrapf(err, "json unmarshal %v", string(b))
			}
			if beian == "" {
				return errors.New("no beian")
			}
			if text == "" {
				return errors.New("no text")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			// Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
			// See https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
			if _, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(apiSecret), nil
			}); err != nil {
				return errors.Wrapf(err, "verify token %v", token)
			}

			if err := rdb.HSet(ctx, SRS_BEIAN, beian, text).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v %v %v", SRS_BEIAN, beian, text)
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "beian: update ok, beian=%v, text=%v, token=%vB", beian, text, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/mgmt/nginx/hls"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var token string
			var enabled bool
			if err := json.Unmarshal(b, &struct {
				Token   *string `json:"token"`
				Enabled *bool   `json:"enabled"`
			}{
				Token: &token, Enabled: &enabled,
			}); err != nil {
				return errors.Wrapf(err, "json unmarshal %v", string(b))
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			// Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
			// See https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
			if _, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(apiSecret), nil
			}); err != nil {
				return errors.Wrapf(err, "verify token %v", token)
			}

			if err := nginxHlsDelivery(ctx, enabled); err != nil {
				return errors.Wrapf(err, "nginxHlsDelivery %v", enabled)
			}
			if err := nginxGenerateConfig(ctx); err != nil {
				return errors.Wrapf(err, "nginx config and reload")
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "nginx hls ok, enabled=%v, token=%vB", enabled, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	// Serve UI at platform.
	fileRoot := path.Join(conf.Pwd, "../ui/build")
	if os.Getenv("REACT_APP_LOCALE") != "" {
		fileRoot = path.Join(fileRoot, os.Getenv("REACT_APP_LOCALE"))
	} else {
		fileRoot = path.Join(fileRoot, "zh")
	}

	fileServer := http.FileServer(http.Dir(fileRoot))
	logger.Tf(ctx, "File server at %v", fileRoot)

	mgmtHandler := func(w http.ResponseWriter, r *http.Request) {
		// Trim the start prefix.
		r.URL.Path = r.URL.Path[len("/mgmt"):]

		// If home or route page, always use virtual main page to serve it.
		serveAsMainPage := r.URL.Path == "/index.html" || r.URL.Path == "/" || r.URL.Path == ""
		if strings.Contains(r.URL.Path, "/routers-") {
			serveAsMainPage = true
		}
		// Should never use /index.html, which will be redirect to /.
		if serveAsMainPage {
			r.URL.Path = "/"
		}

		// We should never cache the main page for react.
		if !serveAsMainPage {
			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%v", 365*24*3600))
		}

		ohttp.SetHeader(w)
		fileServer.ServeHTTP(w, r)
	}

	ep = "/mgmt"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, mgmtHandler)

	ep = "/mgmt/"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, mgmtHandler)

	// Proxy to other services, migrate from mgmt.
	createProxy := func(target string) (*httputil.ReverseProxy, error) {
		targetObject, err := url.Parse(target)
		if err != nil {
			return nil, errors.Wrapf(err, "parse backend %v", target)
		}
		return httputil.NewSingleHostReverseProxy(targetObject), nil
	}

	proxy2023, err := createProxy("http://127.0.0.1:2023")
	if err != nil {
		return err
	}

	proxy1985, err := createProxy("http://127.0.0.1:1985")
	if err != nil {
		return err
	}

	proxy8080, err := createProxy("http://127.0.0.1:8080")
	if err != nil {
		return err
	}

	mgmtFileServer := http.FileServer(http.Dir(path.Join(conf.Pwd, "containers/www")))

	ep = "/"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		// Set common header.
		ohttp.SetHeader(w)

		// For version management.
		if strings.HasPrefix(r.URL.Path, "/terraform/v1/releases") {
			logger.Tf(ctx, "Proxy %v to backend 2023", r.URL.Path)
			proxy2023.ServeHTTP(w, r)
			return
		}

		// We directly serve the static files, because we overwrite the www for DVR.
		if strings.HasPrefix(r.URL.Path, "/console/") || strings.HasPrefix(r.URL.Path, "/players/") ||
			strings.HasPrefix(r.URL.Path, "/tools/") {
			if r.URL.Path != "/tools/player.html" && r.URL.Path != "/tools/xgplayer.html" {
				w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%v", 365*24*3600))
			}
			mgmtFileServer.ServeHTTP(w, r)
			return
		}

		// Proxy to SRS HTTP streaming, console and player, by /api/, /rtc/, /live/, /console/, /players/
		// See https://github.com/vagusX/koa-proxies
		// TODO: FIXME: Do authentication for api.
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/rtc/") {
			logger.Tf(ctx, "Proxy %v to backend 1985", r.URL.Path)
			proxy1985.ServeHTTP(w, r)
			return
		}

		if strings.HasSuffix(r.URL.Path, ".flv") || strings.HasSuffix(r.URL.Path, ".m3u8") ||
			strings.HasSuffix(r.URL.Path, ".ts") || strings.HasSuffix(r.URL.Path, ".aac") ||
			strings.HasSuffix(r.URL.Path, ".mp3") {
			logger.Tf(ctx, "Proxy %v to backend 8080", r.URL.Path)
			proxy8080.ServeHTTP(w, r)
			return

		}

		http.Redirect(w, r, "/mgmt", http.StatusFound)
	})

	return nil
}
