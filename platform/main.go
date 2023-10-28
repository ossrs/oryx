//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
package main

import (
	"context"
	"flag"
	"fmt"
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

	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// conf is a global config object.
var conf *Config

func init() {
	certManager = NewCertManager()
	conf = NewConfig()
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

	// Initialize the management password and load the environment without relying on Redis.
	if true {
		if pwd, err := os.Getwd(); err != nil {
			return errors.Wrapf(err, "getpwd")
		} else {
			conf.Pwd = pwd
		}

		// Note that we only use .env in mgmt.
		envFile := path.Join(conf.Pwd, "containers/data/config/.env")
		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Overload(envFile); err != nil {
				return errors.Wrapf(err, "load %v", envFile)
			}
		}
	}

	// For platform, default to development for Darwin.
	setEnvDefault("NODE_ENV", "development")
	// For platform, HTTP server listen port.
	setEnvDefault("PLATFORM_LISTEN", "2024")
	// Set the default language, en or zh.
	setEnvDefault("REACT_APP_LOCALE", "en")

	// Migrate from mgmt.
	setEnvDefault("REDIS_PORT", "6379")
	setEnvDefault("MGMT_LISTEN", "2022")
	setEnvDefault("PLATFORM_DOCKER", "on")

	// For HTTPS.
	setEnvDefault("HTTPS_LISTEN", "2443")
	setEnvDefault("AUTO_SELF_SIGNED_CERTIFICATE", "off")

	// For feature control.
	setEnvDefault("NAME_LOOKUP", "off")

	logger.Tf(ctx, "load .env as MGMT_PASSWORD=%vB, "+
		"SRS_PLATFORM_SECRET=%vB, CLOUD=%v, REGION=%v, SOURCE=%v, "+
		"NODE_ENV=%v, LOCAL_RELEASE=%v, REDIS_PASSWORD=%vB, REDIS_PORT=%v, "+
		"PUBLIC_URL=%v, BUILD_PATH=%v, REACT_APP_LOCALE=%v, PLATFORM_LISTEN=%v, "+
		"REGISTRY=%v, MGMT_LISTEN=%v, HTTPS_LISTEN=%v, AUTO_SELF_SIGNED_CERTIFICATE=%v, "+
		"NAME_LOOKUP=%v",
		len(os.Getenv("MGMT_PASSWORD")), len(os.Getenv("SRS_PLATFORM_SECRET")), os.Getenv("CLOUD"),
		os.Getenv("REGION"), os.Getenv("SOURCE"), os.Getenv("NODE_ENV"), os.Getenv("LOCAL_RELEASE"),
		len(os.Getenv("REDIS_PASSWORD")), os.Getenv("REDIS_PORT"), os.Getenv("PUBLIC_URL"),
		os.Getenv("BUILD_PATH"), os.Getenv("REACT_APP_LOCALE"), os.Getenv("PLATFORM_LISTEN"),
		os.Getenv("REGISTRY"), os.Getenv("MGMT_LISTEN"), os.Getenv("HTTPS_LISTEN"),
		os.Getenv("AUTO_SELF_SIGNED_CERTIFICATE"), os.Getenv("NAME_LOOKUP"),
	)

	// Setup the base OS for redis, which should never depends on redis.
	if err := initMgmtOS(ctx); err != nil {
		return errors.Wrapf(err, "init mgmt os")
	}

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

	// We must initialize the mgmt after redis is ready.
	if err := initMmgt(ctx); err != nil {
		return errors.Wrapf(err, "init mgmt")
	}
	logger.Tf(ctx, "initialize platform region=%v, registry=%v, version=%v", conf.Region, conf.Registry, version)

	// Create candidate worker for resolving domain to ip.
	candidateWorker = NewCandidateWorker()
	defer candidateWorker.Close()
	if err := candidateWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start candidate worker")
	}

	// Create callback worker.
	callbackWorker = NewCallbackWorker()
	defer callbackWorker.Close()
	if err := callbackWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start callback worker")
	}

	// Create transcode worker for transcoding.
	transcodeWorker = NewTranscodeWorker()
	defer transcodeWorker.Close()
	if err := transcodeWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start transcode worker")
	}

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

	// Create worker for crontab.
	crontabWorker = NewCrontabWorker()
	defer crontabWorker.Close()
	if err := crontabWorker.Start(ctx); err != nil {
		return errors.Wrapf(err, "start crontab worker")
	}

	// Run HTTP service.
	httpService := NewHTTPService()
	defer httpService.Close()
	if err := httpService.Run(ctx); err != nil {
		return errors.Wrapf(err, "start http service")
	}

	return nil
}

// Initialize the source for redis, note that we don't change the env.
func initMgmtOS(ctx context.Context) (err error) {
	// For Darwin, append the search PATH for docker.
	// Note that we should set the PATH env, not the exec.Cmd.Env.
	if conf.IsDarwin && !strings.Contains(os.Getenv("PATH"), "/usr/local/bin") {
		os.Setenv("PATH", fmt.Sprintf("%v:/usr/local/bin", os.Getenv("PATH")))
	}

	// The redis is not available when os startup, so we must directly discover from env or network.
	if conf.Cloud, conf.Region, err = discoverRegion(ctx); err != nil {
		return errors.Wrapf(err, "discover region")
	}

	// Always update the source, because it might change.
	if conf.Source, err = discoverSource(ctx, conf.Cloud, conf.Region); err != nil {
		return errors.Wrapf(err, "discover source")
	}

	// Always update the registry, because it might change.
	if conf.Registry, err = discoverRegistry(ctx, conf.Source); err != nil {
		return errors.Wrapf(err, "discover registry")
	}

	logger.Tf(ctx, "initOS %v", conf.String())
	return
}

// Initialize the source for redis, note that we don't change the env.
func initOS(ctx context.Context) (err error) {
	// Create api secret if not exists, see setupApiSecret
	if token, err := rdb.HGet(ctx, SRS_PLATFORM_SECRET, "token").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v token", SRS_PLATFORM_SECRET)
	} else if token == "" {
		token = os.Getenv("SRS_PLATFORM_SECRET")
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

	// Load the platform from redis, initialized by mgmt.
	if cloud, err := rdb.HGet(ctx, SRS_TENCENT_LH, "cloud").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v cloud", SRS_TENCENT_LH)
	} else if cloud == "" || conf.Cloud != cloud {
		if err = rdb.HSet(ctx, SRS_TENCENT_LH, "cloud", conf.Cloud).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hset %v cloud %v", SRS_TENCENT_LH, conf.Cloud)
		}
		logger.Tf(ctx, "Update cloud=%v", conf.Cloud)
	}

	// Load the region first, because it never changed.
	if region, err := rdb.HGet(ctx, SRS_TENCENT_LH, "region").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v region", SRS_TENCENT_LH)
	} else if region == "" || conf.Region != region {
		if err = rdb.HSet(ctx, SRS_TENCENT_LH, "region", conf.Region).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hset %v region %v", SRS_TENCENT_LH, conf.Region)
		}
		logger.Tf(ctx, "Update region=%v", conf.Region)
	}

	// Always update the source, because it might change.
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

	// Discover and update the platform for stat only, not the OS platform.
	if platform, err := discoverPlatform(ctx, conf.Cloud); err != nil {
		return errors.Wrapf(err, "discover platform by cloud=%v", conf.Cloud)
	} else {
		if err = rdb.HSet(ctx, SRS_TENCENT_LH, "platform", platform).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hset %v platform %v", SRS_TENCENT_LH, platform)
		}
		logger.Tf(ctx, "Update platform=%v", platform)
	}

	// Always update the registry, because it might change.
	if registry, err := rdb.HGet(ctx, SRS_TENCENT_LH, "registry").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v registry", SRS_TENCENT_LH)
	} else if registry == "" || conf.Registry != registry {
		if err = rdb.HSet(ctx, SRS_TENCENT_LH, "registry", conf.Registry).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hset %v registry %v", SRS_TENCENT_LH, conf.Registry)
		}
		logger.Tf(ctx, "Update registry=%v", conf.Registry)
	}

	// Request the host platform OS, whether the OS is Darwin.
	hostPlatform := runtime.GOOS
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

// Initialize the platform before thread run.
func initPlatform(ctx context.Context) error {
	// For Darwin, append the search PATH for docker.
	// Note that we should set the PATH env, not the exec.Cmd.Env.
	// Note that it depends on conf.IsDarwin, so it's unavailable util initOS.
	if conf.IsDarwin && !strings.Contains(os.Getenv("PATH"), "/usr/local/bin") {
		os.Setenv("PATH", fmt.Sprintf("%v:/usr/local/bin", os.Getenv("PATH")))
	}

	// Create directories for data, allow user to link it.
	for _, dir := range []string{
		"containers/data/dvr", "containers/data/record", "containers/data/vod",
		"containers/data/upload", "containers/data/vlive", "containers/data/signals",
		"containers/data/lego", "containers/data/.well-known", "containers/data/config",
	} {
		if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
			if err = os.MkdirAll(dir, os.ModeDir|os.FileMode(0755)); err != nil {
				return errors.Wrapf(err, "create dir %v", dir)
			}
		}
	}

	// Run only once for a special version.
	bootRelease := "v2023-r29"
	if firstRun, err := rdb.HGet(ctx, SRS_FIRST_BOOT, bootRelease).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_FIRST_BOOT, bootRelease)
	} else if firstRun == "" {
		logger.Tf(ctx, "boot setup, v=%v, key=%v", bootRelease, SRS_FIRST_BOOT)

		// Generate the dynamic config for NGINX.
		if err := nginxGenerateConfig(ctx); err != nil {
			return errors.Wrapf(err, "nginx config and reload")
		}

		// Generate the dynamic config for SRS.
		if err := srsGenerateConfig(ctx); err != nil {
			return errors.Wrapf(err, "srs config and reload")
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
	go refreshLatestVersion(ctx)

	// Disable srs-dev, only enable srs-server.
	if srsDevEnabled, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, srsDevDockerName).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_CONTAINER_DISABLED, srsDevDockerName)
	} else if srsEnabled, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, srsDockerName).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_CONTAINER_DISABLED, srsDockerName)
	} else if srsDevEnabled != "true" && srsEnabled == "true" {
		r0 := rdb.HSet(ctx, SRS_CONTAINER_DISABLED, srsDevDockerName, "true").Err()
		r1 := rdb.HSet(ctx, SRS_CONTAINER_DISABLED, srsDockerName, "false").Err()
		logger.Wf(ctx, "Disable srs-dev r0=%v, only enable srs-server r1=%v", r0, r1)
	}

	// For SRS, if release enabled, disable dev automatically.
	if srsReleaseDisabled, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, srsDockerName).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_CONTAINER_DISABLED, srsDockerName)
	} else if srsDevDisabled, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, srsDevDockerName).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_CONTAINER_DISABLED, srsDevDockerName)
	} else if srsReleaseDisabled != "true" && srsDevDisabled != "true" {
		r0 := rdb.HSet(ctx, SRS_CONTAINER_DISABLED, srsDevDockerName, true).Err()
		logger.Tf(ctx, "disable srs dev for release enabled, r0=%v", r0)
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

	// Cancel upgrading.
	if upgrading, err := rdb.HGet(ctx, SRS_UPGRADING, "upgrading").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v upgrading", SRS_UPGRADING)
	} else if upgrading == "1" {
		if err = rdb.HSet(ctx, SRS_UPGRADING, "upgrading", "0").Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hset %v upgrading 0", SRS_UPGRADING)
		}
	}

	// Initialize the node id.
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

// Initialize the platform before thread run.
func initMmgt(ctx context.Context) error {
	// Always create the data dir and sub dirs.
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

	// Refresh the env file.
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
	if os.Getenv("MGMT_PASSWORD") != "" {
		envs["MGMT_PASSWORD"] = os.Getenv("MGMT_PASSWORD")
	}

	if err := godotenv.Write(envs, envFile); err != nil {
		return errors.Wrapf(err, "write %v", envFile)
	}
	logger.Tf(ctx, "Refresh %v ok", envFile)

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
