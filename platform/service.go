// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"github.com/ossrs/go-oryx-lib/errors"
	ohttp "github.com/ossrs/go-oryx-lib/http"
	"github.com/ossrs/go-oryx-lib/logger"

	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
)

// HttpService is a HTTP server for platform.
type HttpService interface {
	Close() error
	Run(ctx context.Context) error
}

func NewHTTPService() HttpService {
	return &httpService{}
}

type httpService struct {
	servers []*http.Server
}

func (v *httpService) Close() error {
	servers := v.servers
	v.servers = nil

	var wg sync.WaitGroup
	defer wg.Wait()

	for index, server := range servers {
		wg.Add(1)
		go func(index int, server *http.Server) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			logger.Tf(ctx, "service shutting down server #%v/%v: %v", index, len(v.servers), server.Addr)
			err := server.Shutdown(ctx)
			logger.Tf(ctx, "service shutdown ok, server #%v/%v: %v, err=%v", index, len(v.servers), server.Addr, err)
		}(index, server)
	}

	return nil
}

func (v *httpService) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	// For debugging server, listen at 127.0.0.1:22022
	go func() {
		dh := http.NewServeMux()
		handleDebuggingGoroutines(context.Background(), dh)
		server := &http.Server{Addr: "127.0.0.1:22022", Handler: dh}
		server.ListenAndServe()
	}()

	ctx, cancel := context.WithCancel(ctx)

	handler := http.NewServeMux()
	if true {
		serviceHandler := http.NewServeMux()
		if err := handleHTTPService(ctx, serviceHandler); err != nil {
			return errors.Wrapf(err, "handle service")
		}

		handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Set common header.
			ohttp.SetHeader(w)

			// Always allow CORS.
			httpAllowCORS(w, r)

			// Allow OPTIONS for CORS.
			if r.Method == http.MethodOptions {
				w.Write(nil)
				return
			}

			// Handle by service handler.
			serviceHandler.ServeHTTP(w, r)
		})
	}

	var r0 error
	if true {
		addr := envPlatformListen()
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
			logger.Tf(ctx, "shutting down HTTP server, addr=%v", addr)
			v.Close()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			if err := server.ListenAndServe(); err != nil && ctx.Err() != context.Canceled {
				r0 = errors.Wrapf(err, "listen %v", addr)
			}
			logger.Tf(ctx, "HTTP server is done, addr=%v", addr)
		}()
	}

	var r1 error
	if true {
		addr := envMgmtListen()
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
			logger.Tf(ctx, "shutting down HTTP server, addr=%v", addr)
			v.Close()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			if err := server.ListenAndServe(); err != nil && ctx.Err() != context.Canceled {
				r1 = errors.Wrapf(err, "listen %v", addr)
			}
			logger.Tf(ctx, "HTTP server is done, addr=%v", addr)
		}()
	}

	var r2 error
	if true {
		addr := envHttpListen()
		if !strings.HasPrefix(addr, ":") {
			addr = fmt.Sprintf(":%v", addr)
		}
		logger.Tf(ctx, "HTTPS listen at %v", addr)

		server := &http.Server{
			Addr:    addr,
			Handler: handler,
			TLSConfig: &tls.Config{
				GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
					return certManager.httpsCertificate, nil
				},
			},
		}
		v.servers = append(v.servers, server)

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			logger.Tf(ctx, "shutting down HTTPS server, addr=%v", addr)
			v.Close()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			if err := server.ListenAndServeTLS("", ""); err != nil && ctx.Err() != context.Canceled {
				r2 = errors.Wrapf(err, "listen %v", addr)
			}
			logger.Tf(ctx, "HTTPS server is done, addr=%v", addr)
		}()
	}

	wg.Wait()
	for _, r := range []error{r0, r1, r2} {
		if r != nil {
			return r
		}
	}
	return nil
}

func handleHTTPService(ctx context.Context, handler *http.ServeMux) error {
	ohttp.Server = fmt.Sprintf("Oryx/%v", version)

	if err := callbackWorker.Handle(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle callback")
	}

	if err := transcriptWorker.Handle(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle transcript")
	}

	if err := ocrWorker.Handle(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle ocr")
	}

	if err := transcodeWorker.Handle(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle transcode")
	}

	if err := forwardWorker.Handle(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle forward")
	}

	if err := vLiveWorker.Handle(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle vLive")
	}

	if err := cameraWorker.Handle(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle IP camera")
	}

	if err := handleHooksService(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle hooks")
	}

	if err := handleLiveRoomService(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle live room")
	}

	if err := handleDubbingService(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle dubbing")
	}

	if err := handleAITalkService(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle AI talk")
	}

	var ep string

	handleHostVersions(ctx, handler)
	handleMgmtVersions(ctx, handler)
	handleFFmpegVersions(ctx, handler)
	handleMgmtInit(ctx, handler)
	handleMgmtCheck(ctx, handler)
	handleMgmtEnvs(ctx, handler)
	handleMgmtToken(ctx, handler)
	handleMgmtLogin(ctx, handler)
	handleMgmtStatus(ctx, handler)
	handleMgmtBilibili(ctx, handler)
	handleMgmtLimitsQuery(ctx, handler)
	handleMgmtLimitsUpdate(ctx, handler)
	handleMgmtOpenAIQuery(ctx, handler)
	handleMgmtOpenAIUpdate(ctx, handler)
	handleMgmtBeianQuery(ctx, handler)
	handleMgmtSecretQuery(ctx, handler)
	handleMgmtBeianUpdate(ctx, handler)
	handleMgmtNginxHlsUpdate(ctx, handler)
	handleMgmtNginxHlsQuery(ctx, handler)
	handleMgmtHlsLowLatencyUpdate(ctx, handler)
	handleMgmtHlsLowLatencyQuery(ctx, handler)
	handleMgmtAutoSelfSignedCertificate(ctx, handler)
	handleMgmtSsl(ctx, handler)
	handleMgmtLetsEncrypt(ctx, handler)
	handleMgmtCertQuery(ctx, handler)
	handleMgmtStreamsQuery(ctx, handler)
	handleMgmtStreamsKickoff(ctx, handler)
	handleMgmtUI(ctx, handler)

	proxy2023, err := httpCreateProxy("http://" + os.Getenv("SRS_PROXY_HOST") + ":2023")
	if err != nil {
		return err
	}

	proxy1985, err := httpCreateProxy("http://" + os.Getenv("SRS_PROXY_HOST") + ":1985")
	if err != nil {
		return err
	}

	proxyWhxp, err := httpCreateProxy("http://" + os.Getenv("SRS_PROXY_HOST") + ":1985")
	if err != nil {
		return err
	}

	proxy8080, err := httpCreateProxy("http://" + os.Getenv("SRS_PROXY_HOST") + ":8080")
	if err != nil {
		return err
	}

	platformFileServer := http.FileServer(http.Dir(path.Join(conf.Pwd, "containers/www")))
	wellKnownFileServer := http.FileServer(http.Dir(path.Join(conf.Pwd, "containers/data")))
	hlsFileServer := http.FileServer(http.Dir(path.Join(conf.Pwd, "containers/objs/nginx/html")))

	ep = "/"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		// For version management.
		if strings.HasPrefix(r.URL.Path, "/terraform/v1/releases") {
			logger.Tf(ctx, "Proxy %v to backend 2023", r.URL.Path)
			proxy2023.ServeHTTP(w, r)
			return
		}

		// For HTTPS management.
		if strings.HasPrefix(r.URL.Path, "/.well-known/") {
			w.Header().Set("Cache-Control", "no-cache, max-age=0")
			wellKnownFileServer.ServeHTTP(w, r)
			return
		}

		// We directly serve the static files, because we overwrite the www for DVR.
		if strings.HasPrefix(r.URL.Path, "/console/") || strings.HasPrefix(r.URL.Path, "/players/") ||
			strings.HasPrefix(r.URL.Path, "/tools/") {
			if r.URL.Path != "/tools/player.html" && r.URL.Path != "/tools/xgplayer.html" {
				w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%v", 30*24*3600))
			}
			platformFileServer.ServeHTTP(w, r)
			return
		}

		// Proxy to SRS RTC API, by /rtc/ prefix.
		if strings.HasPrefix(r.URL.Path, "/rtc/") {
			q := r.URL.Query()
			if eip := q.Get("eip"); eip != "" {
				logger.Tf(ctx, "Proxy %v to backend 1985, eip=%v, query is %v",
					r.URL.Path, eip, r.URL.RawQuery)
			} else {
				// Allow test to mock and overwrite the host.
				host := r.Header.Get("X-Real-Host")
				if host == "" {
					host = r.Host
				}

				// Resolve the host to ip.
				starttime := time.Now()
				if ip, err := candidateWorker.Resolve(host); err != nil {
					logger.Ef(ctx, "Proxy %v to backend 1985, resolve %v/%v failed, cost=%v, err is %v",
						r.URL.Path, r.Host, host, time.Now().Sub(starttime), err)
					ohttp.WriteError(ctx, w, r, err)
					return
				} else if ip != nil {
					eip = ip.String()
					r.URL.RawQuery += fmt.Sprintf("&eip=%v", eip)
					logger.Tf(ctx, "Proxy %v to backend 1985, host=%v/%v, resolved ip=%v, cost=%v, query is %v",
						r.URL.Path, r.Host, host, eip, time.Now().Sub(starttime), r.URL.RawQuery)
				}
			}

			proxyWhxp.ServeHTTP(&whxpResponseModifier{w}, r)
			return
		}

		// Use versions API as health check API, no auth.
		if r.URL.Path == "/api/v1/versions" {
			logger.Tf(ctx, "Proxy %v to backend 1985", r.URL.Path)
			proxy1985.ServeHTTP(w, r)
			return
		}

		// Proxy to SRS HTTP API, for console, by /api/ prefix.
		if strings.HasPrefix(r.URL.Path, "/api/") {
			token := r.URL.Query().Get("token")
			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				ohttp.WriteError(ctx, w, r, err)
				return
			}

			logger.Tf(ctx, "Proxy %v to backend 1985", r.URL.Path)
			proxy1985.ServeHTTP(w, r)
			return
		}

		// Always directly serve the HLS ts files.
		if fastCache.HLSHighPerformance && strings.HasSuffix(r.URL.Path, ".m3u8") {
			var m3u8ExpireInSeconds int = 10
			if fastCache.HLSLowLatency {
				m3u8ExpireInSeconds = 1 // Note that we use smaller expire time that fragment duration.
			}

			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%v", m3u8ExpireInSeconds))
			hlsFileServer.ServeHTTP(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, ".ts") {
			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%v", 600))
			hlsFileServer.ServeHTTP(w, r)
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

func handleDebuggingGoroutines(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/debug/goroutines"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 1<<16)
		stacklen := runtime.Stack(buf, true)
		fmt.Fprintf(w, "%s", buf[:stacklen])
	})
}

func handleHostVersions(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/host/versions"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		ohttp.WriteData(ctx, w, r, &struct {
			Version string `json:"version"`
		}{
			Version: strings.TrimPrefix(version, "v"),
		})
	})
}

func handleMgmtVersions(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/versions"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		ohttp.WriteData(ctx, w, r, &struct {
			Version string `json:"version"`
		}{
			Version: strings.TrimPrefix(version, "v"),
		})
	})
}

func handleFFmpegVersions(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/ffmpeg/versions"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		ohttp.WriteData(ctx, w, r, &struct {
			Version string `json:"version"`
		}{
			Version: strings.TrimPrefix(version, "v"),
		})
	})
}

func handleMgmtInit(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/init"
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
					Init: envMgmtPassword() != "",
				})
				return nil
			}

			// If already initialized, never set it again.
			if envMgmtPassword() != "" {
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
			if err := godotenv.Overload(envFile); err != nil {
				return errors.Wrapf(err, "load %v", envFile)
			}

			apiSecret := envApiSecret()
			expireAt, createAt, token, err := createToken(ctx, envApiSecret())
			if err != nil {
				return errors.Wrapf(err, "build token")
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Token    string `json:"token"`
				CreateAt string `json:"createAt"`
				ExpireAt string `json:"expireAt"`
				// Allow user to directly use Bearer token.
				Bearer string `json:"bearer"`
			}{
				Token: token, CreateAt: createAt.Format(time.RFC3339), ExpireAt: expireAt.Format(time.RFC3339),
				Bearer: apiSecret,
			})
			logger.Tf(ctx, "init password ok, create=%v, expire=%v, password=%vB", createAt, expireAt, len(password))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtCheck(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/check"
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
}

func handleMgmtEnvs(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/envs"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var locale string
			if err := ParseBody(ctx, r.Body, &struct {
				Locale *string `json:"locale"`
			}{
				Locale: &locale,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			// Filter the locale.
			if locale != "en" && locale != "zh" {
				locale = "un"
			}

			if err := rdb.Set(ctx, SRS_LOCALE, locale, 0).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "set %v %v", SRS_LOCALE, locale)
			}

			var forwardLimit int
			if envForwardLimit() != "" {
				if iv, err := strconv.ParseInt(envForwardLimit(), 10, 64); err != nil {
					return errors.Wrapf(err, "parse env forward limit %v", envForwardLimit())
				} else {
					forwardLimit = int(iv)
				}
			}

			var vLiveLimit int
			if envVLiveLimit() != "" {
				if iv, err := strconv.ParseInt(envVLiveLimit(), 10, 64); err != nil {
					return errors.Wrapf(err, "parse env virtual live limit %v", envVLiveLimit())
				} else {
					vLiveLimit = int(iv)
				}
			}

			var cameraLimit int
			if envCameraLimit() != "" {
				if iv, err := strconv.ParseInt(envCameraLimit(), 10, 64); err != nil {
					return errors.Wrapf(err, "parse env camera limit %v", envCameraLimit())
				} else {
					cameraLimit = int(iv)
				}
			}

			platformDocker := envPlatformDocker() != "off"
			candidate := envCandidate() != ""
			ohttp.WriteData(ctx, w, r, &struct {
				// Whether mgmt run in docker.
				MgmtDocker bool `json:"mgmtDocker"`
				// Whether platform run in docker.
				PlatformDocker bool `json:"platformDocker"`
				// Whether set the env CANDIDATE for WebRTC.
				Candidate bool `json:"candidate"`
				// The exposed RTMP port.
				RTMPPort string `json:"rtmpPort"`
				// The exposed HTTP port.
				HTTPPort string `json:"httpPort"`
				// The exposed SRT port.
				SRTPort string `json:"srtPort"`
				// The exposed RTC port.
				RTCPort string `json:"rtcPort"`
				// The limit of the number of forwarding streams.
				ForwardLimit int `json:"forwardLimit"`
				// The limit of the number of vLive streams.
				VLiveLimit int `json:"vLiveLimit"`
				// The limit of the number of IP camera streams.
				CameraLimit int `json:"cameraLimit"`
			}{
				// Whether in docker.
				MgmtDocker: true,
				// Whether platform in docker.
				PlatformDocker: platformDocker,
				// The candidate IP for WebRTC.
				Candidate: candidate,
				// The export port for RTMP.
				RTMPPort: envRtmpPort(),
				// The export port for HTTP.
				HTTPPort: envHttpPort(),
				// The export port for SRT.
				SRTPort: envSrtListen(),
				// The export port for WebRTC.
				RTCPort: envRtcListen(),
				// The limit of the number of forwarding streams.
				ForwardLimit: forwardLimit,
				// The limit of the number of vLive streams.
				VLiveLimit: vLiveLimit,
				// The limit of the number of IP camera streams.
				CameraLimit: cameraLimit,
			})

			logger.Tf(ctx, "mgmt envs ok, locale=%v, platformDocker=%v, candidate=%v, rtmpPort=%v, httpPort=%v, srtPort=%v, rtcPort=%v, forwardLimit=%v, vLiveLimit=%v, cameraLimit=%v",
				locale, platformDocker, candidate, envRtmpPort(), envHttpPort(),
				envSrtListen(), envRtcListen(), forwardLimit, vLiveLimit, cameraLimit,
			)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtToken(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/token"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			expireAt, createAt, token, err := createToken(ctx, envApiSecret())
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
}

func handleMgmtLogin(ctx context.Context, handler *http.ServeMux) {
	var loginLock sync.Mutex
	ep := "/terraform/v1/mgmt/login"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			if !loginLock.TryLock() {
				return errors.New("login is running, try later")
			}
			defer loginLock.Unlock()

			if envMgmtPassword() == "" {
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

			if password != envMgmtPassword() {
				wait := time.Duration(10) * time.Second
				logger.Wf(ctx, "Invalid password, wait for %v", wait)

				select {
				case <-time.After(wait):
				case <-ctx.Done():
				}

				return errors.Errorf("invalid password, wait %v", wait)
			}

			apiSecret := envApiSecret()
			expireAt, createAt, token, err := createToken(ctx, apiSecret)
			if err != nil {
				return errors.Wrapf(err, "build token")
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Token    string `json:"token"`
				CreateAt string `json:"createAt"`
				ExpireAt string `json:"expireAt"`
				// Allow user to directly use Bearer token.
				Bearer string `json:"bearer"`
			}{
				Token: token, CreateAt: createAt.Format(time.RFC3339), ExpireAt: expireAt.Format(time.RFC3339),
				Bearer: apiSecret,
			})
			logger.Tf(ctx, "login by password ok, create=%v, expire=%v, token=%vB", createAt, expireAt, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtStatus(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/status"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
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
}

func handleMgmtBilibili(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/bilibili"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token, bvid string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				BVID  *string `json:"bvid"`
			}{
				Token: &token, BVID: &bvid,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if bvid == "" {
				return errors.New("no bvid")
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
				if envNodeEnv() == "development" {
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
}

func handleMgmtOpenAIQuery(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/openai/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			aiSecretKey, err := rdb.HGet(ctx, SRS_SYS_OPENAI, "key").Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v key", SRS_SYS_OPENAI)
			}

			aiBaseURL, err := rdb.HGet(ctx, SRS_SYS_OPENAI, "url").Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v url", SRS_SYS_OPENAI)
			}

			aiOrganization, err := rdb.HGet(ctx, SRS_SYS_OPENAI, "org").Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v org", SRS_SYS_OPENAI)
			}

			ohttp.WriteData(ctx, w, r, &struct {
				// The AI secret key.
				AISecretKey string `json:"aiSecretKey"`
				// The AI base url.
				AIBaseURL string `json:"aiBaseURL"`
				// The AI organization.
				AIOrganization string `json:"aiOrganization"`
			}{
				AISecretKey: aiSecretKey, AIBaseURL: aiBaseURL, AIOrganization: aiOrganization,
			})

			logger.Tf(ctx, "settings: query openai ok")
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtOpenAIUpdate(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/openai/update"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var aiSecretKey, aiBaseURL, aiOrganization string
			if err := ParseBody(ctx, r.Body, &struct {
				Token          *string `json:"token"`
				AISecretKey    *string `json:"aiSecretKey"`
				AIBaseURL      *string `json:"aiBaseURL"`
				AIOrganization *string `json:"aiOrganization"`
			}{
				Token: &token, AISecretKey: &aiSecretKey, AIBaseURL: &aiBaseURL,
				AIOrganization: &aiOrganization,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if aiSecretKey == "" {
				return errors.New("no aiSecretKey")
			}
			if aiBaseURL == "" {
				return errors.New("no aiBaseURL")
			}

			if err := rdb.HSet(ctx, SRS_SYS_OPENAI, "key", aiSecretKey).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v key %v", SRS_SYS_OPENAI, aiSecretKey)
			}
			if err := rdb.HSet(ctx, SRS_SYS_OPENAI, "url", aiBaseURL).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v url %v", SRS_SYS_OPENAI, aiBaseURL)
			}
			if err := rdb.HSet(ctx, SRS_SYS_OPENAI, "org", aiOrganization).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v org %v", SRS_SYS_OPENAI, aiOrganization)
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "limits: Update ok, key=%vB, url=%v, org=%v", len(aiSecretKey), aiBaseURL, aiOrganization)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtLimitsQuery(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/limits/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			vLiveLimits, err := rdb.HGet(ctx, SRS_SYS_LIMITS, "vlive").Int64()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v vlive", SRS_SYS_LIMITS)
			} else if vLiveLimits == 0 {
				vLiveLimits = SrsSysLimitsVLive
			}

			ipCameraLimits, err := rdb.HGet(ctx, SRS_SYS_LIMITS, "camera").Int64()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v camera", SRS_SYS_LIMITS)
			} else if ipCameraLimits == 0 {
				ipCameraLimits = SrsSysLimitsCamera
			}

			ohttp.WriteData(ctx, w, r, &struct {
				// The limits for virtual live streaming.
				VLive int64 `json:"vlive"`
				// The limits for IP camera streaming.
				IPCamera int64 `json:"camera"`
			}{
				VLive: vLiveLimits, IPCamera: ipCameraLimits,
			})

			logger.Tf(ctx, "limits: query ok")
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtLimitsUpdate(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/limits/update"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var vlive, camera int64
			if err := ParseBody(ctx, r.Body, &struct {
				Token    *string `json:"token"`
				VLive    *int64  `json:"vlive"`
				IPCamera *int64  `json:"camera"`
			}{
				Token: &token, VLive: &vlive, IPCamera: &camera,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if vlive <= 0 {
				return errors.Errorf("invalid vlive %v", vlive)
			}
			if camera <= 0 {
				return errors.Errorf("invalid vlive %v", vlive)
			}

			if err := rdb.HSet(ctx, SRS_SYS_LIMITS, "vlive", vlive).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v vlive %v", SRS_SYS_LIMITS, vlive)
			}
			if err := rdb.HSet(ctx, SRS_SYS_LIMITS, "camera", camera).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v camera %v", SRS_SYS_LIMITS, camera)
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "limits: Update ok, vlive=%v, camera=%v", vlive, camera)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

// Note that this API is not verified by token.
func handleMgmtBeianQuery(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/beian/query"
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
}

func handleMgmtSecretQuery(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/secret/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			ohttp.WriteData(ctx, w, r, apiSecret)
			logger.Tf(ctx, "query apiSecret ok, versions=%v, token=%vB", conf.Versions.String(), len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtBeianUpdate(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/beian/update"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token, beian, text string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				Beian *string `json:"beian"`
				Text  *string `json:"text"`
			}{
				Token: &token, Beian: &beian, Text: &text,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if beian == "" {
				return errors.New("no beian")
			}
			if text == "" {
				return errors.New("no text")
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
}

func handleMgmtNginxHlsUpdate(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/hphls/update"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var noHlsCtx bool
			if err := ParseBody(ctx, r.Body, &struct {
				Token    *string `json:"token"`
				NoHlsCtx *bool   `json:"noHlsCtx"`
			}{
				Token: &token, NoHlsCtx: &noHlsCtx,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			noHlsCtxValue := fmt.Sprintf("%v", noHlsCtx)
			if err := rdb.HSet(ctx, SRS_HP_HLS, "noHlsCtx", noHlsCtxValue).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v noHlsCtx %v", SRS_HP_HLS, noHlsCtxValue)
			}

			if err := srsGenerateConfig(ctx); err != nil {
				return errors.Wrapf(err, "generate SRS config")
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "nginx hls update ok, enabled=%v, token=%vB", noHlsCtx, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtNginxHlsQuery(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/hphls/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			var enabled bool
			if v, err := rdb.HGet(ctx, SRS_HP_HLS, "noHlsCtx").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_HP_HLS, "noHlsCtx")
			} else {
				enabled = v == "true"
			}

			ohttp.WriteData(ctx, w, r, &struct {
				NoHlsCtx bool `json:"noHlsCtx"`
			}{
				NoHlsCtx: enabled,
			})
			logger.Tf(ctx, "nginx hls query ok, enabled=%v, token=%vB", enabled, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtHlsLowLatencyUpdate(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/hlsll/update"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var hlsLowLatency bool
			if err := ParseBody(ctx, r.Body, &struct {
				Token         *string `json:"token"`
				HlsLowLatency *bool   `json:"hlsLowLatency"`
			}{
				Token: &token, HlsLowLatency: &hlsLowLatency,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			hlsLowLatencyValue := fmt.Sprintf("%v", hlsLowLatency)
			if err := rdb.HSet(ctx, SRS_LL_HLS, "hlsLowLatency", hlsLowLatencyValue).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v hlsLowLatency %v", SRS_LL_HLS, hlsLowLatencyValue)
			}

			if err := srsGenerateConfig(ctx); err != nil {
				return errors.Wrapf(err, "generate SRS config")
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "hls low latency update ok, enabled=%v, token=%vB", hlsLowLatency, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtHlsLowLatencyQuery(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/hlsll/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			var enabled bool
			if v, err := rdb.HGet(ctx, SRS_LL_HLS, "hlsLowLatency").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_LL_HLS, "hlsLowLatency")
			} else {
				enabled = v == "true"
			}

			ohttp.WriteData(ctx, w, r, &struct {
				HlsLowLatency bool `json:"hlsLowLatency"`
			}{
				HlsLowLatency: enabled,
			})
			logger.Tf(ctx, "hls low latency query ok, enabled=%v, token=%vB", enabled, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtAutoSelfSignedCertificate(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/auto-self-signed-certificate"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if err := certManager.createSelfSignCertificate(ctx); err != nil {
				return errors.Wrapf(err, "create self sign certificate")
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "create self-signed cert ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtSsl(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/ssl"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var key, crt string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				Key   *string `json:"key"`
				Crt   *string `json:"crt"`
			}{
				Token: &token, Key: &key, Crt: &crt,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if key = strings.TrimSpace(key); key == "" {
				return errors.New("empty key")
			}
			if crt = strings.TrimSpace(crt); crt == "" {
				return errors.New("empty crt")
			}

			if err := certManager.updateSslFiles(ctx, key+"\n", crt+"\n"); err != nil {
				return errors.Wrapf(err, "updateSslFiles key=%vB, crt=%vB", len(key), len(crt))
			}

			if err := rdb.Set(ctx, SRS_HTTPS, "ssl", 0).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "set %v %v", SRS_HTTPS, "ssl")
			}

			if err := nginxGenerateConfig(ctx); err != nil {
				return errors.Wrapf(err, "nginx config and reload")
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "nginx ssl file ok, key=%vB, crt=%vB, token=%vB", len(key), len(crt), len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtLetsEncrypt(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/letsencrypt"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var domain string
			if err := ParseBody(ctx, r.Body, &struct {
				Token  *string `json:"token"`
				Domain *string `json:"domain"`
			}{
				Token: &token, Domain: &domain,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if domain = strings.TrimSpace(domain); domain == "" {
				return errors.New("empty domain")
			}

			if err := certManager.updateLetsEncrypt(ctx, domain); err != nil {
				return errors.Wrapf(err, "updateSslFiles domain=%v", domain)
			}

			if err := rdb.Set(ctx, SRS_HTTPS, "lets", 0).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "set %v %v", SRS_HTTPS, "lets")
			}
			if err := rdb.Set(ctx, SRS_HTTPS_DOMAIN, domain, 0).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "set %v %v", SRS_HTTPS_DOMAIN, domain)
			}

			if err := nginxGenerateConfig(ctx); err != nil {
				return errors.Wrapf(err, "nginx config and reload")
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "nginx letsencrypt ok, domain=%v, token=%vB", domain, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtCertQuery(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/cert/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			provider, err := rdb.Get(ctx, SRS_HTTPS).Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "get %v", SRS_HTTPS)
			}

			domain, err := rdb.Get(ctx, SRS_HTTPS_DOMAIN).Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "get %v", SRS_HTTPS_DOMAIN)
			}

			var key, crt string
			if provider != "" {
				key, crt, err = certManager.QueryCertificate()
				if err != nil {
					return errors.Wrapf(err, "query cert")
				}
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Provider string `json:"provider"`
				Domain   string `json:"domain"`
				Key      string `json:"key"`
				Crt      string `json:"crt"`
			}{
				Provider: provider, Domain: domain, Key: key, Crt: crt,
			})
			logger.Tf(ctx, "query cert ok, provider=%v, domain=%v, key=%vB, crt=%vB, token=%vB",
				provider, domain, len(key), len(crt), len(token),
			)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtStreamsQuery(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/streams/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			streams, err := rdb.HGetAll(ctx, SRS_STREAM_ACTIVE).Result()
			if err != nil {
				return errors.Wrapf(err, "hgetall %v", SRS_STREAM_ACTIVE)
			}

			var streamObjects []*SrsStream
			for _, value := range streams {
				var stream SrsStream
				if err := json.Unmarshal([]byte(value), &stream); err != nil {
					return errors.Wrapf(err, "unmarshal %v", value)
				}

				streamObjects = append(streamObjects, &stream)
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Streams []*SrsStream `json:"streams"`
			}{
				streamObjects,
			})
			logger.Tf(ctx, "query streams ok, streams=%v, token=%vB", len(streamObjects), len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

// See SRS error code ERROR_RTMP_CLIENT_NOT_FOUND
const ErrorRtmpClientNotFound = 2049

func handleMgmtStreamsKickoff(ctx context.Context, handler *http.ServeMux) {
	ep := "/terraform/v1/mgmt/streams/kickoff"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var vhost, app, stream string
			if err := ParseBody(ctx, r.Body, &struct {
				Token  *string `json:"token"`
				Vhost  *string `json:"vhost"`
				App    *string `json:"app"`
				Stream *string `json:"stream"`
			}{
				Token: &token, Vhost: &vhost, App: &app, Stream: &stream,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if vhost == "" {
				return errors.New("no vhost")
			}
			if app == "" {
				return errors.New("no app")
			}
			if stream == "" {
				return errors.New("no stream")
			}

			streamObject := &SrsStream{Vhost: vhost, App: app, Stream: stream}
			streamURL := streamObject.StreamURL()
			if target, err := rdb.HGet(ctx, SRS_STREAM_ACTIVE, streamURL).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_STREAM_ACTIVE, streamURL)
			} else if target == "" {
				return errors.Errorf("stream not found %v", streamURL)
			} else if err := json.Unmarshal([]byte(target), &streamObject); err != nil {
				return errors.Wrapf(err, "unmarshal %v", target)
			}

			if streamObject.Client == "" {
				return errors.Errorf("no client_id for %v", streamURL)
			}

			// Start request and parse the code.
			requestClient := func(ctx context.Context, clientURL, method string) (int, string, error) {
				req, err := http.NewRequest(method, clientURL, nil)
				if err != nil {
					return 0, "", errors.Wrapf(err, "new request")
				}

				res, err := http.DefaultClient.Do(req.WithContext(ctx))
				if err != nil {
					return 0, "", errors.Wrapf(err, "do request")
				}
				defer res.Body.Close()

				b, err := io.ReadAll(res.Body)
				if err != nil {
					return 0, "", errors.Wrapf(err, "http read body")
				}

				if res.StatusCode != http.StatusOK {
					return 0, "", errors.Errorf("status %v", res.StatusCode)
				}

				var code int
				if err := json.Unmarshal(b, &struct {
					Code *int `json:"code"`
				}{
					Code: &code,
				}); err != nil {
					return 0, "", errors.Wrapf(err, "unmarshal %v", string(b))
				}
				return code, string(b), nil
			}

			// Whether client exists in SRS server.
			var code int
			clientURL := fmt.Sprintf("http://%v:1985/api/v1/clients/%v", os.Getenv("SRS_PROXY_HOST"), streamObject.Client)
			if r0, body, err := requestClient(ctx, clientURL, http.MethodGet); err != nil {
				return errors.Wrapf(err, "http query client %v", clientURL)
			} else if r0 != 0 && r0 != ErrorRtmpClientNotFound {
				return errors.Errorf("invalid code=%v, body=%v", r0, body)
			} else {
				code = r0
			}

			// Kickoff if exists, ignore if not.
			if code == 0 {
				if r0, body, err := requestClient(ctx, clientURL, http.MethodDelete); err != nil {
					return errors.Wrapf(err, "kickoff %v, body %v", clientURL, body)
				} else if r0 != 0 && r0 != ErrorRtmpClientNotFound {
					return errors.Errorf("invalid code=%v, body=%v", r0, body)
				}
			}

			if err := rdb.HDel(ctx, SRS_STREAM_ACTIVE, streamURL).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hdel %v %v", SRS_STREAM_ACTIVE, streamURL)
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "kickoff stream ok, code=%v, token=%vB", code, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})
}

func handleMgmtUI(ctx context.Context, handler *http.ServeMux) {
	// Serve UI at platform.
	fileRoot := path.Join(conf.Pwd, "../ui/build", envReactAppLocale())

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

	ep := "/mgmt"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, mgmtHandler)

	ep = "/mgmt/"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, mgmtHandler)
}
