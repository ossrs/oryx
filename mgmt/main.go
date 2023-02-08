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
	"os"
	"os/exec"
	"os/signal"
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

	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		return errors.Wrapf(err, "load .env")
	}
	setEnvDefault("REDIS_PORT", "6379")
	setEnvDefault("MGMT_LISTEN", "2022")
	logger.Tf(ctx, "load .env as MGMT_PASSWORD=%vB, SRS_PLATFORM_SECRET=%vB, CLOUD=%v, REGION=%v, SOURCE=%v, "+
		"NODE_ENV=%v, LOCAL_RELEASE=%v, SRS_DOCKER=%v, USE_DOCKER=%v, SRS_UTEST=%v, REDIS_PASSWORD=%vB, REDIS_PORT=%v, "+
		"PUBLIC_URL=%v, BUILD_PATH=%v, REACT_APP_LOCALE=%v, MGMT_LISTEN=%v",
		len(os.Getenv("MGMT_PASSWORD")), len(os.Getenv("SRS_PLATFORM_SECRET")), os.Getenv("CLOUD"),
		os.Getenv("REGION"), os.Getenv("SOURCE"), os.Getenv("NODE_ENV"), os.Getenv("LOCAL_RELEASE"),
		os.Getenv("SRS_DOCKER"), os.Getenv("USE_DOCKER"), os.Getenv("SRS_UTEST"),
		len(os.Getenv("REDIS_PASSWORD")), os.Getenv("REDIS_PORT"), os.Getenv("PUBLIC_URL"),
		os.Getenv("BUILD_PATH"), os.Getenv("REACT_APP_LOCALE"), os.Getenv("MGMT_LISTEN"),
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

	// Setup the OS for redis, which should never depends on redis.
	if err := initOS(ctx); err != nil {
		return errors.Wrapf(err, "init os")
	}

	// Always restart the redis container.
	redisManager := NewDockerRedisManager()

	stopRedisTime := time.Now()
	redisStopCtx, _ := context.WithTimeout(ctx, 15*time.Second)
	if err := redisManager.Stop(redisStopCtx, 15*time.Second); err != nil {
		logger.Tf(ctx, "ignore stop redis err %v", err)
	}

	startRedisTime := time.Now()
	if err := redisManager.Start(ctx); err != nil {
		return errors.Wrapf(err, "start redis")
	}
	logger.Tf(ctx, "restart redis container, stop=%vms, start=%vms",
		startRedisTime.Sub(stopRedisTime), time.Now().Sub(startRedisTime))

	// Wait for redis to be ready.
	redisStartCtx, _ := context.WithTimeout(ctx, 30*time.Second)
	if err := redisManager.Ready(redisStartCtx); err != nil {
		return errors.Wrapf(err, "wait redis ready")
	}
	logger.Tf(ctx, "redis is running")

	// Initialize global rdb, the redis client.
	if err := InitRdb(); err != nil {
		return errors.Wrapf(err, "init rdb")
	}
	logger.Tf(ctx, "init rdb(redis client) ok")

	// We must initialize the mgmt after redis is ready.
	if err := initMmgt(ctx); err != nil {
		return errors.Wrapf(err, "init mgmt")
	}
	logger.Tf(ctx, "initialize mgmt region=%v, registry=%v, version=%v", conf.Region, conf.Registry, version)

	// Note that we always restart the platform container.
	if err := removeContainer(ctx, platformDockerName); err != nil {
		logger.Tf(ctx, "ignore restart platform err %v", err)
	}

	// Create backend service manager.
	service := NewDockerBackendService(redisManager, NewDockerPlatformManager())
	if err := service.Start(ctx); err != nil {
		return errors.Wrapf(err, "start service manager")
	}
	defer service.Close()

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
	// Initialize global config.
	conf = NewConfig()

	// Initialize pwd.
	if conf.Pwd, err = os.Getwd(); err != nil {
		return errors.Wrapf(err, "get pwd")
	}

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

	// Discover the platform, not the GOOS, for report and statistic only.
	if conf.Platform, err = discoverPlatform(ctx, conf.Cloud); err != nil {
		return errors.Wrapf(err, "discover platform")
	}

	// Create directories for data, allow user to link it.
	for _, dir := range []string{
		"containers/data/dvr", "containers/data/record", "containers/data/vod",
		"containers/data/upload", "containers/data/vlive",
	} {
		if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
			if err = os.MkdirAll(dir, os.ModeDir|os.FileMode(0755)); err != nil {
				return errors.Wrapf(err, "create dir %v", dir)
			}
		}
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
			if name, ipv4, err := discoverPrivateIPv4(ctx); err != nil {
				logger.Wf(ctx, "ignore ipv4 discover err %v", err)
			} else if name != "" && ipv4 != nil {
				conf.ipv4 = ipv4
				conf.Iface = name
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
func initMmgt(ctx context.Context) error {
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

	// Create api secret if not exists, see setupApiSecret
	if token, err := rdb.HGet(ctx, SRS_PLATFORM_SECRET, "token").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v token", SRS_PLATFORM_SECRET)
	} else if token == "" {
		token = fmt.Sprintf("srs-v1-%v", strings.ReplaceAll(uuid.NewString(), "-", ""))
		if err = rdb.HSet(ctx, SRS_PLATFORM_SECRET, "token", token).Err(); err != nil {
			return errors.Wrapf(err, "hset %v token %v", SRS_PLATFORM_SECRET, token)
		}

		update := time.Now().Format(time.RFC3339)
		if err = rdb.HSet(ctx, SRS_PLATFORM_SECRET, "update", update).Err(); err != nil {
			return errors.Wrapf(err, "hset %v update %v", SRS_PLATFORM_SECRET, update)
		}
		logger.Tf(ctx, "Platform api secret update, token=%vB, update=%v", len(token), update)
	}

	// Load the cloud first, because it never changed.
	if cloud, err := rdb.HGet(ctx, SRS_TENCENT_LH, "cloud").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v cloud", SRS_TENCENT_LH)
	} else if cloud == "" || conf.Cloud != cloud {
		if err = rdb.HSet(ctx, SRS_TENCENT_LH, "cloud", conf.Cloud).Err(); err != nil {
			return errors.Wrapf(err, "hset %v cloud %v", SRS_TENCENT_LH, conf.Cloud)
		}
		logger.Tf(ctx, "Update cloud=%v", conf.Cloud)
	}

	// Load the region first, because it never changed.
	if region, err := rdb.HGet(ctx, SRS_TENCENT_LH, "region").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v region", SRS_TENCENT_LH)
	} else if region == "" || conf.Region != region {
		if err = rdb.HSet(ctx, SRS_TENCENT_LH, "region", conf.Region).Err(); err != nil {
			return errors.Wrapf(err, "hset %v region %v", SRS_TENCENT_LH, conf.Region)
		}
		logger.Tf(ctx, "Update region=%v", conf.Region)
	}

	// Always update the source, because it might change.
	if source, err := rdb.HGet(ctx, SRS_TENCENT_LH, "source").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v source", SRS_TENCENT_LH)
	} else if source == "" || conf.Source != source {
		if err = rdb.HSet(ctx, SRS_TENCENT_LH, "source", conf.Source).Err(); err != nil {
			return errors.Wrapf(err, "hset %v source %v", SRS_TENCENT_LH, conf.Source)
		}
		logger.Tf(ctx, "Update source=%v", conf.Source)
	}

	// Always update the registry, because it might change.
	if registry, err := rdb.HGet(ctx, SRS_TENCENT_LH, "registry").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v registry", SRS_TENCENT_LH)
	} else if registry == "" || conf.Registry != registry {
		if err = rdb.HSet(ctx, SRS_TENCENT_LH, "registry", conf.Registry).Err(); err != nil {
			return errors.Wrapf(err, "hset %v registry %v", SRS_TENCENT_LH, conf.Registry)
		}
		logger.Tf(ctx, "Update registry=%v", conf.Registry)
	}

	// Load the platform first, because it never changed.
	if platform, err := rdb.HGet(ctx, SRS_TENCENT_LH, "platform").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v platform", SRS_TENCENT_LH)
	} else if platform == "" || conf.Platform != platform {
		if err = rdb.HSet(ctx, SRS_TENCENT_LH, "platform", conf.Platform).Err(); err != nil {
			return errors.Wrapf(err, "hset %v platform %v", SRS_TENCENT_LH, conf.Platform)
		}
		logger.Tf(ctx, "Update platform=%v", conf.Platform)
	}

	// Refresh the env file.
	if envs, err := godotenv.Read(".env"); err != nil {
		return errors.Wrapf(err, "load envs")
	} else {
		envs["CLOUD"] = conf.Cloud
		envs["REGION"] = conf.Region
		envs["SOURCE"] = conf.Source
		if os.Getenv("MGMT_PASSWORD") != "" {
			envs["MGMT_PASSWORD"] = os.Getenv("MGMT_PASSWORD")
		}

		if err := godotenv.Write(envs, ".env"); err != nil {
			return errors.Wrapf(err, "write .env")
		}
		logger.Tf(ctx, "Refresh .env ok")
	}

	// Query the api secret from redis, cache it to env.
	if os.Getenv("SRS_PLATFORM_SECRET") == "" {
		if token, err := rdb.HGet(ctx, SRS_PLATFORM_SECRET, "token").Result(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hget %v token", SRS_PLATFORM_SECRET)
		} else {
			os.Setenv("SRS_PLATFORM_SECRET", token)
			logger.Tf(ctx, "Update api secret to %vB", len(token))
		}
	}

	// Disable srs-dev, only enable srs-server.
	if srsDevEnabled, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, srsDevDockerName).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_CONTAINER_DISABLED, srsDevDockerName)
	} else if srsEnabled, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, srsDockerName).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_CONTAINER_DISABLED, srsDockerName)
	} else if srsDevEnabled != "true" && srsEnabled == "true" {
		r0 := rdb.HSet(ctx, SRS_CONTAINER_DISABLED, srsDevDockerName, "true").Err()
		r1 := rdb.HSet(ctx, SRS_CONTAINER_DISABLED, srsDockerName, "false").Err()
		r2 := removeContainer(ctx, srsDevDockerName)
		logger.Wf(ctx, "Disable srs-dev r0=%v, r2=%v, only enable srs-server r1=%v", r0, r2, r1)
	}

	// Cleanup the previous unused images.
	newImage := fmt.Sprintf("%v/ossrs/srs-cloud:mgmt-%v", conf.Registry, version)
	if previousImage, err := rdb.HGet(ctx, SRS_DOCKER_IMAGES, mgmtDockerName).Result(); err != nil && err != redis.Nil {
		return err
	} else {
		if err = rdb.HSet(ctx, SRS_DOCKER_IMAGES, mgmtDockerName, newImage).Err(); err != nil {
			return err
		}
		if previousImage != "" && previousImage != newImage {
			r0 := exec.CommandContext(ctx, "docker", "rmi", previousImage).Run()
			logger.Tf(ctx, "remove previous mgmt image %v, err %v", previousImage, r0)
		}
	}

	return nil
}
