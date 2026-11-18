// Copyright (c) 2022-2026 Winlin
//
// SPDX-License-Identifier: MIT
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var conf *Config

func init() {
	certManager = NewCertManager()
	conf = NewConfig()

	// The crontab polls to refresh the fast cache, such as the LLHLS config.
	fastCache = NewFastCache()
}

func main() {
	ctx := logger.WithContext(context.Background())
	ctx = logger.WithContext(ctx)

	if err := doMain(ctx); err != nil {
		logger.Tf(ctx, "run err %+v", err)
		return
	}

	logger.Tf(ctx, "run ok")
}

func doMain(ctx context.Context) error {
	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "Print version and quit")
	flag.BoolVar(&showVersion, "version", false, "Print version and quit")
	flag.Parse()

	if showVersion {
		fmt.Println(strings.TrimPrefix(version, "v"))
		os.Exit(0)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		for s := range sc {
			logger.Tf(ctx, "Got signal %v", s)
			cancel()
		}
	}()

	// The main thread may be blocked after the context is cancelled, so force the exit after a timeout.
	go func() {
		<-ctx.Done()
		time.Sleep(30 * time.Second)
		logger.Wf(ctx, "Force to exit by timeout")
		os.Exit(1)
	}()

	// Load the .env, which holds the management password, without relying on Redis.
	if true {
		if pwd, err := os.Getwd(); err != nil {
			return errors.Wrapf(err, "getpwd")
		} else {
			conf.Pwd = pwd
		}

		envFile := path.Join(conf.Pwd, "containers/data/config/.env")
		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Overload(envFile); err != nil {
				return errors.Wrapf(err, "load %v", envFile)
			}
		}
	}

	setEnvDefault("NODE_ENV", "development")
	setEnvDefault("PLATFORM_LISTEN", "2024")
	// The language, en or zh.
	setEnvDefault("REACT_APP_LOCALE", "en")
	setEnvDefault("GO_PPROF", "")

	setEnvDefault("REDIS_DATABASE", "0")
	setEnvDefault("REDIS_HOST", "127.0.0.1")
	setEnvDefault("REDIS_PORT", "6379")
	setEnvDefault("MGMT_LISTEN", "2022")

	setEnvDefault("HTTPS_LISTEN", "2443")
	setEnvDefault("AUTO_SELF_SIGNED_CERTIFICATE", "on")

	setEnvDefault("NAME_LOOKUP", "on")
	setEnvDefault("PLATFORM_DOCKER", "off")

	setEnvDefault("RTMP_PORT", "1935")
	setEnvDefault("HTTP_PORT", "")
	setEnvDefault("SRT_PORT", "10080")
	setEnvDefault("RTC_PORT", "8000")

	setEnvDefault("SRS_FORWARD_LIMIT", "10")
	setEnvDefault("SRS_VLIVE_LIMIT", "10")
	setEnvDefault("SRS_CAMERA_LIMIT", "10")

	logger.Tf(ctx, "load .env as MGMT_PASSWORD=%vB, GO_PPROF=%v, "+
		"SRS_PLATFORM_SECRET=%vB, CLOUD=%v, REGION=%v, SOURCE=%v, SRT_PORT=%v, RTC_PORT=%v, "+
		"NODE_ENV=%v, LOCAL_RELEASE=%v, REDIS_DATABASE=%v, REDIS_HOST=%v, REDIS_PASSWORD=%vB, REDIS_PORT=%v, RTMP_PORT=%v, "+
		"PUBLIC_URL=%v, BUILD_PATH=%v, REACT_APP_LOCALE=%v, PLATFORM_LISTEN=%v, HTTP_PORT=%v, "+
		"REGISTRY=%v, MGMT_LISTEN=%v, HTTPS_LISTEN=%v, AUTO_SELF_SIGNED_CERTIFICATE=%v, "+
		"NAME_LOOKUP=%v, PLATFORM_DOCKER=%v, SRS_FORWARD_LIMIT=%v, SRS_VLIVE_LIMIT=%v, "+
		"SRS_CAMERA_LIMIT=%v, YTDL_PROXY=%v",
		len(envMgmtPassword()), envGoPprof(), len(envApiSecret()), envCloud(),
		envRegion(), envSource(), envSrtListen(), envRtcListen(),
		envNodeEnv(), envLocalRelease(),
		envRedisDatabase(), envRedisHost(), len(envRedisPassword()), envRedisPort(),
		envRtmpPort(), envPublicUrl(),
		envBuildPath(), envReactAppLocale(), envPlatformListen(), envHttpPort(),
		envRegistry(), envMgmtListen(), envHttpListen(),
		envSelfSignedCertificate(), envNameLookup(),
		envPlatformDocker(), envForwardLimit(), envVLiveLimit(),
		envCameraLimit(), envYtdlProxy(),
	)

	if addr := envGoPprof(); addr != "" {
		go func() {
			logger.Tf(ctx, "Start Go pprof at %v", addr)
			http.ListenAndServe(addr, nil)
		}()
	}

	// Must not depend on Redis, which is not ready yet.
	if err := initMgmtOS(ctx); err != nil {
		return errors.Wrapf(err, "init mgmt os")
	}

	if err := InitRdb(); err != nil {
		return errors.Wrapf(err, "init rdb")
	}
	logger.Tf(ctx, "init rdb(redis client) ok")

	if err := initOS(ctx); err != nil {
		return errors.Wrapf(err, "init os")
	}

	if err := initPlatform(ctx); err != nil {
		return errors.Wrapf(err, "init platform")
	}

	if err := initMmgt(ctx); err != nil {
		return errors.Wrapf(err, "init mgmt")
	}
	logger.Tf(ctx, "initialize platform region=%v, registry=%v, version=%v", conf.Region, conf.Registry, version)

	candidateWorker = NewCandidateWorker()
	defer candidateWorker.Close()
	if err := candidateWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start candidate worker")
	}

	callbackWorker = NewCallbackWorker()
	defer callbackWorker.Close()
	if err := callbackWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start callback worker")
	}

	transcriptWorker = NewTranscriptWorker()
	defer transcriptWorker.Close()
	if err := transcriptWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start transcript worker")
	}

	ocrWorker = NewOCRWorker()
	defer ocrWorker.Close()
	if err := ocrWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start OCR worker")
	}

	dubbingServer = NewDubbingServer()
	defer dubbingServer.Close()

	transcodeWorker = NewTranscodeWorker()
	defer transcodeWorker.Close()
	if err := transcodeWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start transcode worker")
	}

	recordWorker = NewRecordWorker()
	defer recordWorker.Close()
	if err := recordWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start record worker")
	}

	dvrWorker = NewDvrWorker()
	defer dvrWorker.Close()
	if err := dvrWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start dvr worker")
	}

	vodWorker = NewVodWorker()
	defer vodWorker.Close()
	if err := vodWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start vod worker")
	}

	forwardWorker = NewForwardWorker()
	defer forwardWorker.Close()
	if err := forwardWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start forward worker")
	}

	vLiveWorker = NewVLiveWorker()
	defer vLiveWorker.Close()
	if err := vLiveWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start vLive worker")
	}

	cameraWorker = NewCameraWorker()
	defer cameraWorker.Close()
	if err := cameraWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start IP camera worker")
	}

	crontabWorker = NewCrontabWorker()
	defer crontabWorker.Close()
	if err := crontabWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start crontab worker")
	}

	httpService := NewHTTPService()
	defer httpService.Close()
	if err := httpService.Run(ctx); err != nil {
		return errors.Wrapf(err, "start http service")
	}

	return nil
}

func initMgmtOS(ctx context.Context) (err error) {
	// On Darwin, add /usr/local/bin to the search PATH for docker.
	// Set the PATH env, not the exec.Cmd.Env.
	if conf.IsDarwin && !strings.Contains(envPath(), "/usr/local/bin") {
		os.Setenv("PATH", fmt.Sprintf("%v:/usr/local/bin", envPath()))
	}

	// Redis is not ready yet, so discover from the env or the network.
	if conf.Cloud, conf.Region, err = discoverRegion(ctx); err != nil {
		return errors.Wrapf(err, "discover region")
	}

	if conf.Source, err = discoverSource(ctx, conf.Cloud, conf.Region); err != nil {
		return errors.Wrapf(err, "discover source")
	}

	if conf.Registry, err = discoverRegistry(ctx, conf.Source); err != nil {
		return errors.Wrapf(err, "discover registry")
	}

	logger.Tf(ctx, "initOS %v", conf.String())
	return
}

func initOS(ctx context.Context) (err error) {
	if token, err := rdb.HGet(ctx, SRS_PLATFORM_SECRET, "token").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v token", SRS_PLATFORM_SECRET)
	} else if token == "" {
		token = envApiSecret()
		if token == "" {
			token = fmt.Sprintf("srs-v2-%v", strings.ReplaceAll(uuid.NewString(), "-", ""))
		}

		if err = rdb.HSet(ctx, SRS_PLATFORM_SECRET, "token", token).Err(); err != nil {
			return errors.Wrapf(err, "hset %v token %v", SRS_PLATFORM_SECRET, token)
		}

		update := time.Now().Format(time.RFC3339)
		if err = rdb.HSet(ctx, SRS_PLATFORM_SECRET, "update", update).Err(); err != nil {
			return errors.Wrapf(err, "hset %v update %v", SRS_PLATFORM_SECRET, update)
		}
		logger.Tf(ctx, "Platform update api secret, token=%vB, at=%v", len(token), update)
	}

	// The platform uses the secret to access the API of mgmt.
	if envApiSecret() == "" {
		if token, err := rdb.HGet(ctx, SRS_PLATFORM_SECRET, "token").Result(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hget %v token", SRS_PLATFORM_SECRET)
		} else {
			os.Setenv("SRS_PLATFORM_SECRET", token)
			logger.Tf(ctx, "Update api secret to %vB", len(token))
		}
	}

	if cloud, err := rdb.HGet(ctx, SRS_TENCENT_LH, "cloud").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v cloud", SRS_TENCENT_LH)
	} else if cloud == "" || conf.Cloud != cloud {
		if err = rdb.HSet(ctx, SRS_TENCENT_LH, "cloud", conf.Cloud).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hset %v cloud %v", SRS_TENCENT_LH, conf.Cloud)
		}
		logger.Tf(ctx, "Update cloud=%v", conf.Cloud)
	}

	if region, err := rdb.HGet(ctx, SRS_TENCENT_LH, "region").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v region", SRS_TENCENT_LH)
	} else if region == "" || conf.Region != region {
		if err = rdb.HSet(ctx, SRS_TENCENT_LH, "region", conf.Region).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hset %v region %v", SRS_TENCENT_LH, conf.Region)
		}
		logger.Tf(ctx, "Update region=%v", conf.Region)
	}

	// The source might change.
	if source, err := rdb.HGet(ctx, SRS_TENCENT_LH, "source").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v source", SRS_TENCENT_LH)
	} else if source == "" || conf.Source != source {
		if err = rdb.HSet(ctx, SRS_TENCENT_LH, "source", conf.Source).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hset %v source %v", SRS_TENCENT_LH, conf.Source)
		}
		logger.Tf(ctx, "Update source=%v", conf.Source)
	}

	if registry, err := rdb.HGet(ctx, SRS_TENCENT_LH, "registry").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v registry", SRS_TENCENT_LH)
	} else if registry != "" {
		conf.Registry = registry
	}

	// The platform is for stat only, not the OS platform.
	if platform, err := discoverPlatform(ctx, conf.Cloud); err != nil {
		return errors.Wrapf(err, "discover platform by cloud=%v", conf.Cloud)
	} else {
		if err = rdb.HSet(ctx, SRS_TENCENT_LH, "platform", platform).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hset %v platform %v", SRS_TENCENT_LH, platform)
		}
		logger.Tf(ctx, "Update platform=%v", platform)
	}

	if registry, err := rdb.HGet(ctx, SRS_TENCENT_LH, "registry").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v registry", SRS_TENCENT_LH)
	} else if registry == "" || conf.Registry != registry {
		if err = rdb.HSet(ctx, SRS_TENCENT_LH, "registry", conf.Registry).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hset %v registry %v", SRS_TENCENT_LH, conf.Registry)
		}
		logger.Tf(ctx, "Update registry=%v", conf.Registry)
	}

	hostPlatform := runtime.GOOS
	if hostPlatform == "darwin" {
		conf.IsDarwin = true
	}

	// Waits for the first private IPv4, then a goroutine keeps updating it.
	if err := refreshIPv4(ctx); err != nil {
		return errors.Wrapf(err, "refresh ipv4")
	}

	logger.Tf(ctx, "initOS %v", conf.String())
	return
}

func initPlatform(ctx context.Context) error {
	// On Darwin, add /usr/local/bin to the search PATH for docker.
	// Set the PATH env, not the exec.Cmd.Env.
	if conf.IsDarwin && !strings.Contains(envPath(), "/usr/local/bin") {
		os.Setenv("PATH", fmt.Sprintf("%v:/usr/local/bin", envPath()))
	}

	// Only missing directories are created, so the user may link one. The user mounts
	// containers/data/srs-s3-bucket to save recordings to cloud storage.
	for _, dir := range []string{
		"containers/data/dvr", "containers/data/record", "containers/data/vod",
		"containers/data/upload", "containers/data/vlive", "containers/data/signals",
		"containers/data/lego", "containers/data/.well-known", "containers/data/config",
		"containers/data/transcript", "containers/data/srs-s3-bucket",
		"containers/data/dubbing", "containers/data/ocr",
	} {
		if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
			if err = os.MkdirAll(dir, os.ModeDir|os.FileMode(0755)); err != nil {
				return errors.Wrapf(err, "create dir %v", dir)
			}
		}
	}

	bootRelease := "v2023-r30"
	if firstRun, err := rdb.HGet(ctx, SRS_FIRST_BOOT, bootRelease).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_FIRST_BOOT, bootRelease)
	} else if firstRun == "" {
		logger.Tf(ctx, "boot setup, v=%v, key=%v", bootRelease, SRS_FIRST_BOOT)

		if err := nginxGenerateConfig(ctx); err != nil {
			return errors.Wrapf(err, "nginx config and reload")
		}

		if err := srsGenerateConfig(ctx); err != nil {
			return errors.Wrapf(err, "srs config and reload")
		}

		if err := rdb.HSet(ctx, SRS_FIRST_BOOT, bootRelease, 1).Err(); err != nil {
			return errors.Wrapf(err, "hset %v %v 1", SRS_FIRST_BOOT, bootRelease)
		}

		logger.Tf(ctx, "boot done, v=%v, key=%v", bootRelease, SRS_FIRST_BOOT)
	} else {
		logger.Tf(ctx, "boot already done, v=%v, key=%v", bootRelease, SRS_FIRST_BOOT)
	}

	go refreshLatestVersion(ctx)

	if srsDevEnabled, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, srsDevDockerName).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_CONTAINER_DISABLED, srsDevDockerName)
	} else if srsEnabled, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, srsDockerName).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_CONTAINER_DISABLED, srsDockerName)
	} else if srsDevEnabled != "true" && srsEnabled == "true" {
		r0 := rdb.HSet(ctx, SRS_CONTAINER_DISABLED, srsDevDockerName, "true").Err()
		r1 := rdb.HSet(ctx, SRS_CONTAINER_DISABLED, srsDockerName, "false").Err()
		logger.Wf(ctx, "Disable srs-dev r0=%v, only enable srs-server r1=%v", r0, r1)
	}

	if srsReleaseDisabled, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, srsDockerName).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_CONTAINER_DISABLED, srsDockerName)
	} else if srsDevDisabled, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, srsDevDockerName).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_CONTAINER_DISABLED, srsDevDockerName)
	} else if srsReleaseDisabled != "true" && srsDevDisabled != "true" {
		r0 := rdb.HSet(ctx, SRS_CONTAINER_DISABLED, srsDevDockerName, true).Err()
		logger.Tf(ctx, "disable srs dev for release enabled, r0=%v", r0)
	}

	if publish, err := rdb.HGet(ctx, SRS_AUTH_SECRET, "pubSecret").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v pubSecret", SRS_AUTH_SECRET)
	} else if publish == "" {
		publish = strings.ReplaceAll(uuid.NewString(), "-", "")
		if err = rdb.HSet(ctx, SRS_AUTH_SECRET, "pubSecret", publish).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hset %v pubSecret %v", SRS_AUTH_SECRET, publish)
		}
		if err = rdb.Set(ctx, SRS_SECRET_PUBLISH, publish, 0).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "set %v %v", SRS_SECRET_PUBLISH, publish)
		}
	}

	// Migrate from previous versions.
	for _, migrate := range []struct {
		PVK string
		CVK string
	}{
		{"SRS_RECORD_M3U8_METADATA", SRS_RECORD_M3U8_ARTIFACT},
		{"SRS_DVR_M3U8_METADATA", SRS_DVR_M3U8_ARTIFACT},
		{"SRS_VOD_M3U8_METADATA", SRS_VOD_M3U8_ARTIFACT},
	} {
		pv, _ := rdb.HLen(ctx, migrate.PVK).Result()
		cv, _ := rdb.HLen(ctx, migrate.CVK).Result()
		if pv > 0 && cv == 0 {
			if vs, err := rdb.HGetAll(ctx, migrate.PVK).Result(); err == nil {
				for k, v := range vs {
					_ = rdb.HSet(ctx, migrate.CVK, k, v)
				}
				logger.Tf(ctx, "migrate %v to %v with %v keys", migrate.PVK, migrate.CVK, len(vs))
			}
		}
	}

	if upgrading, err := rdb.HGet(ctx, SRS_UPGRADING, "upgrading").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v upgrading", SRS_UPGRADING)
	} else if upgrading == "1" {
		if err = rdb.HSet(ctx, SRS_UPGRADING, "upgrading", "0").Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hset %v upgrading 0", SRS_UPGRADING)
		}
	}

	if nid, err := rdb.HGet(ctx, SRS_TENCENT_LH, "node").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v node", SRS_TENCENT_LH)
	} else if nid == "" {
		nid = uuid.NewString()
		if err = rdb.HSet(ctx, SRS_TENCENT_LH, "node", nid).Err(); err != nil {
			return errors.Wrapf(err, "hset %v node %v", SRS_TENCENT_LH, nid)
		}
		logger.Tf(ctx, "Update node nid=%v", nid)
	}

	return nil
}

func initMmgt(ctx context.Context) error {
	dataDir := filepath.Join(conf.Pwd, "containers", "data")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.RemoveAll(dataDir); err != nil {
			return errors.Wrapf(err, "remove data dir %s", dataDir)
		}
	}

	dirs := []string{"redis", "config", "dvr", "record", "vod", "upload", "vlive"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(dataDir, dir), 0755); err != nil {
			return errors.Wrapf(err, "create dir %s", dir)
		}
	}

	envs := map[string]string{}
	envFile := path.Join(conf.Pwd, "containers/data/config/.env")
	if _, err := os.Stat(envFile); err == nil {
		if v, err := godotenv.Read(envFile); err != nil {
			return errors.Wrapf(err, "load envs")
		} else {
			envs = v
		}
	}

	envs["CLOUD"] = conf.Cloud
	envs["REGION"] = conf.Region
	envs["SOURCE"] = conf.Source
	envs["REGISTRY"] = conf.Registry
	if envMgmtPassword() != "" {
		envs["MGMT_PASSWORD"] = envMgmtPassword()
	}

	if err := godotenv.Write(envs, envFile); err != nil {
		return errors.Wrapf(err, "write %v", envFile)
	}
	logger.Tf(ctx, "Refresh %v ok", envFile)

	return nil
}

func refreshLatestVersion(ctx context.Context) error {
	versionsCtx, versionsCancel := context.WithCancel(context.Background())
	go func() {
		ctx := logger.WithContext(ctx)
		for ctx.Err() == nil {
			versions, err := queryLatestVersion(ctx)
			if err == nil && versions != nil && versions.Latest != "" {
				logger.Tf(ctx, "query version ok, result is %v", versions.String())
				conf.Versions = *versions
				versionsCancel()

				// CrontabWorker refreshes the version later.
				break
			}

			select {
			case <-ctx.Done():
			case <-time.After(3 * time.Minute):
			}
		}
	}()

	select {
	case <-ctx.Done():
	case <-versionsCtx.Done():
	}

	return nil
}
