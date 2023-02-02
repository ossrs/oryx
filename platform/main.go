//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: MIT
//
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	ctx := logger.WithContext(context.Background())

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

	// For platform, we need to load .env from pwd or mgmt. Right now, we still don't know about the .env file path,
	// util we load the conf.Pwd from mgmt by execApi in initOS.
	if true {
		pwd, err := os.Getwd()
		if err != nil {
			return errors.Wrapf(err, "getpwd")
		}

		envFile := path.Join(pwd, "../mgmt/.env")
		if _, err := os.Stat(".env"); err == nil {
			if err := godotenv.Load(".env"); err != nil {
				return errors.Wrapf(err, "load .env")
			}
		} else if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Load(envFile); err != nil {
				return errors.Wrapf(err, "load %v", envFile)
			}
		} else {
			return errors.Errorf("no .env or %v", envFile)
		}
	}

	// For platform, default to development for Darwin.
	setEnvDefault("NODE_ENV", "development")
	// For platform, HTTP server listen port.
	setEnvDefault("PLATFORM_LISTEN", "2024")

	setEnvDefault("REDIS_PORT", "6379")
	logger.Tf(ctx, "load .env as MGMT_PASSWORD=%vB, SRS_PLATFORM_SECRET=%vB, CLOUD=%v, REGION=%v, SOURCE=%v, "+
		"NODE_ENV=%v, LOCAL_RELEASE=%v, SRS_DOCKER=%v, USE_DOCKER=%v, SRS_UTEST=%v, REDIS_PASSWORD=%vB, REDIS_PORT=%v, "+
		"PUBLIC_URL=%v, BUILD_PATH=%v, REACT_APP_LOCALE=%v, PLATFORM_LISTEN=%v",
		len(os.Getenv("MGMT_PASSWORD")), len(os.Getenv("SRS_PLATFORM_SECRET")), os.Getenv("CLOUD"),
		os.Getenv("REGION"), os.Getenv("SOURCE"), os.Getenv("NODE_ENV"), os.Getenv("LOCAL_RELEASE"),
		os.Getenv("SRS_DOCKER"), os.Getenv("USE_DOCKER"), os.Getenv("SRS_UTEST"),
		len(os.Getenv("REDIS_PASSWORD")), os.Getenv("REDIS_PORT"), os.Getenv("PUBLIC_URL"),
		os.Getenv("BUILD_PATH"), os.Getenv("REACT_APP_LOCALE"), os.Getenv("PLATFORM_LISTEN"),
	)

	// Install signals.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		for s := range sc {
			logger.Tf(ctx, "Got signal %v", s)
			cancel()
		}
	}()

	// Initialize global rdb, the redis client.
	if err := InitRdb(); err != nil {
		return errors.Wrapf(err, "init rdb")
	}
	logger.Tf(ctx, "init rdb(redis client) ok")

	// For platform, we should initOS after redis.
	// Setup the OS for redis, which should never depends on redis.
	if err := initOS(ctx); err != nil {
		return errors.Wrapf(err, "init os")
	}

	// We must initialize the platform after redis is ready.
	if err := initPlatform(ctx); err != nil {
		return errors.Wrapf(err, "init platform")
	}
	logger.Tf(ctx, "initialize platform region=%v, registry=%v, version=%v", conf.Region, conf.Registry, version)

	// Note that we always restart the SRS container.
	if err := execApi(ctx, "removeContainer", []string{srsDockerName}, nil); err != nil {
		logger.Tf(ctx, "ignore restart SRS err %v", err)
	}

	// Create backend service manager.
	service := NewDockerBackendService(NewDockerSrsManager())
	if err := service.Start(ctx); err != nil {
		return errors.Wrapf(err, "start service manager")
	}
	defer service.Close()

	// Create worker for RECORD, covert live stream to local file.
	recordWorker = NewRecordWorker()
	defer recordWorker.Close()
	if err := recordWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start record worker")
	}

	// Create worker for DVR, covert live stream to local file.
	dvrWorker = NewDvrWorker()
	defer dvrWorker.Close()
	if err := dvrWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start dvr worker")
	}

	// Create worker for VoD, covert live stream to local file.
	vodWorker = NewVodWorker()
	defer vodWorker.Close()
	if err := vodWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start vod worker")
	}

	// Create worker for forwarding.
	forwardWorker = NewForwardWorker()
	defer forwardWorker.Close()
	if err := forwardWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start forward worker")
	}

	// Create worker for vLive.
	vLiveWorker = NewVLiveWorker()
	defer vLiveWorker.Close()
	if err := vLiveWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start vLive worker")
	}

	// Run HTTP service.
	httpService := NewDockerHTTPService()
	defer httpService.Close()
	if err := httpService.Run(ctx); err != nil {
		return errors.Wrapf(err, "start http service")
	}

	return nil
}

// Initialize the source for redis, note that we don't change the env.
func initOS(ctx context.Context) (err error) {
	// For platform, we must use the secret to access API of mgmt.
	// Query the api secret from redis, cache it to env.
	if os.Getenv("SRS_PLATFORM_SECRET") == "" {
		if token, err := rdb.HGet(ctx, SRS_PLATFORM_SECRET, "token").Result(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hget %v token", SRS_PLATFORM_SECRET)
		} else {
			os.Setenv("SRS_PLATFORM_SECRET", token)
			logger.Tf(ctx, "Update api secret to %vB", len(token))
		}
	}

	// Initialize global config.
	conf = NewConfig()

	// Initialize pwd.
	if err = execApi(ctx, "cwd", nil, &struct {
		Cwd *string `json:"cwd"`
	}{
		Cwd: &conf.Pwd,
	}); err != nil {
		return errors.Wrapf(err, "get pwd")
	}

	// Load the platform from redis, initialized by mgmt.
	if cloud, err := rdb.HGet(ctx, SRS_TENCENT_LH, "cloud").Result(); err != nil {
		return errors.Wrapf(err, "hget %v cloud", SRS_TENCENT_LH)
	} else {
		conf.Cloud = cloud
	}

	if region, err := rdb.HGet(ctx, SRS_TENCENT_LH, "region").Result(); err != nil {
		return errors.Wrapf(err, "hget %v region", SRS_TENCENT_LH)
	} else {
		conf.Region = region
	}

	if source, err := rdb.HGet(ctx, SRS_TENCENT_LH, "source").Result(); err != nil {
		return errors.Wrapf(err, "hget %v source", SRS_TENCENT_LH)
	} else {
		conf.Source = source
	}

	if registry, err := rdb.HGet(ctx, SRS_TENCENT_LH, "registry").Result(); err != nil {
		return errors.Wrapf(err, "hget %v registry", SRS_TENCENT_LH)
	} else {
		conf.Registry = registry
	}

	// Request the host platform OS, whether the OS is Darwin.
	var hostPlatform string
	if err = execApi(ctx, "hostPlatform", nil, &struct {
		Platform *string `json:"platform"`
	}{
		Platform: &hostPlatform,
	}); err != nil {
		return errors.Wrapf(err, "get platform")
	}
	// Because platform might run in docker, so we overwrite it by query from mgmt.
	if hostPlatform == "darwin" {
		conf.IsDarwin = true
	}

	// Start a goroutine to update ipv4 of config.
	if err := refreshIPv4(ctx); err != nil {
		return errors.Wrapf(err, "refresh ipv4")
	}

	logger.Tf(ctx, "initOS %v", conf.String())
	return
}

// Refresh the ipv4 address.
func refreshIPv4(ctx context.Context) error {
	ipv4Ctx, ipv4Cancel := context.WithCancel(context.Background())
	go func() {
		ctx := logger.WithContext(ctx)
		for ctx.Err() == nil {
			var address string
			if err := execApi(ctx, "ipv4", nil, &struct {
				Name    *string `json:"name"`
				Address *string `json:"address"`
			}{
				Name: &conf.Iface, Address: &address,
			}); err == nil && address != "" {
				logger.Tf(ctx, "query ipv4 ok, result is %v %v", conf.Iface, address)
				conf.ipv4 = net.ParseIP(address)
				ipv4Cancel()
			}

			duration := time.Duration(24*3600) * time.Second
			if os.Getenv("NODE_ENV") == "development" {
				duration = time.Duration(30) * time.Second
			}
			time.Sleep(duration)
		}
	}()

	select {
	case <-ctx.Done():
	case <-ipv4Ctx.Done():
	}

	return nil
}

// Initialize the platform before thread run.
func initPlatform(ctx context.Context) error {
	// For Darwin, append the search PATH for docker.
	// Note that we should set the PATH env, not the exec.Cmd.Env.
	// Note that it depends on conf.IsDarwin, so it's unavailable util initOS.
	if conf.IsDarwin && !strings.Contains(os.Getenv("PATH"), "/usr/local/bin") {
		os.Setenv("PATH", fmt.Sprintf("%v:/usr/local/bin", os.Getenv("PATH")))
	}

	// Run only once for a special version.
	bootRelease := "v23"
	if firstRun, err := rdb.HGet(ctx, SRS_FIRST_BOOT, bootRelease).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_FIRST_BOOT, bootRelease)
	} else if firstRun == "" {
		logger.Tf(ctx, "boot setup, v=%v, key=%v", bootRelease, SRS_FIRST_BOOT)

		// Generate the dynamic config for NGINX.
		if err := execApi(ctx, "nginxGenerateConfig", nil, nil); err != nil {
			return errors.Wrapf(err, "execApi nginxGenerateConfig")
		}

		// Remove containers for IP might change, and use network srs-cloud.
		names := []string{srsDockerName, srsDevDockerName}
		for _, name := range names {
			if err := execApi(ctx, "rmContainer", []string{name}, nil); err != nil {
				return errors.Wrapf(err, "execApi rmContainer %v", []string{name})
			}
		}

		// Remove the unused containers.
		names = []string{
			"prometheus", "node-exporter", "srs-hooks", "ffmpeg", "tencent-cloud",
		}
		for _, name := range names {
			if err := execApi(ctx, "rmContainer", []string{name}, nil); err != nil {
				return errors.Wrapf(err, "execApi rmContainer %v", []string{name})
			}
		}

		// Run once, record in redis.
		if err := rdb.HSet(ctx, SRS_FIRST_BOOT, bootRelease, 1).Err(); err != nil {
			return errors.Wrapf(err, "hset %v %v 1", SRS_FIRST_BOOT, bootRelease)
		}

		logger.Tf(ctx, "boot done, v=%v, key=%v", bootRelease, SRS_FIRST_BOOT)
	} else {
		logger.Tf(ctx, "boot already done, v=%v, key=%v", bootRelease, SRS_FIRST_BOOT)
	}

	// For development, request the releases from itself which proxy to the releases service.
	if err := refreshLatestVersion(ctx); err != nil {
		return errors.Wrapf(err, "refresh latest version")
	}

	// For SRS, if release enabled, disable dev automatically.
	if srsReleaseDisabled, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, srsDockerName).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_CONTAINER_DISABLED, srsDockerName)
	} else if srsDevDisabled, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, srsDevDockerName).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_CONTAINER_DISABLED, srsDevDockerName)
	} else if srsReleaseDisabled != "true" && srsDevDisabled != "true" {
		r0 := rdb.HSet(ctx, SRS_CONTAINER_DISABLED, srsDevDockerName, true).Err()
		r1 := execApi(ctx, "rmContainer", []string{srsDevDockerName}, nil)
		logger.Tf(ctx, "disable srs dev for release enabled, r0=%v, r1=%v", r0, r1)
	}

	// Setup the publish secret for first run.
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

	return nil
}

// Refresh the latest version.
func refreshLatestVersion(ctx context.Context) error {
	versionsCtx, versionsCancel := context.WithCancel(context.Background())
	go func() {
		ctx := logger.WithContext(ctx)
		for ctx.Err() == nil {
			versions, err := queryLatestVersion(ctx)
			if err == nil && versions.Latest != "" {
				logger.Tf(ctx, "query version ok, result is %v", versions.String())
				conf.Versions = *versions
				versionsCancel()
			}

			duration := time.Duration(24*3600) * time.Second
			if os.Getenv("NODE_ENV") == "development" {
				duration = time.Duration(30) * time.Second
			}
			time.Sleep(duration)
		}
	}()

	select {
	case <-ctx.Done():
	case <-versionsCtx.Done():
	}

	return nil
}
