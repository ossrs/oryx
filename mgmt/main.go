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
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"

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

	setEnvDefault := func (key, value string) {
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		return errors.Wrapf(err, "load .env")
	}
	setEnvDefault("REDIS_PORT", "6379")
	setEnvDefault("MGMT_LISTEN", "2022")
	setEnvDefault("PLATFORM_DOCKER", "true")
	logger.Tf(ctx, "load .env as MGMT_PASSWORD=%vB, CLOUD=%v, REGION=%v, SOURCE=%v, "+
		"NODE_ENV=%v, LOCAL_RELEASE=%v, SRS_DOCKER=%v, USE_DOCKER=%v, SRS_UTEST=%v, REDIS_PASSWORD=%vB, REDIS_PORT=%v, "+
		"PUBLIC_URL=%v, BUILD_PATH=%v, REACT_APP_LOCALE=%v, MGMT_LISTEN=%v, PLATFORM_DOCKER=%v",
		len(os.Getenv("MGMT_PASSWORD")), os.Getenv("CLOUD"),
		os.Getenv("REGION"), os.Getenv("SOURCE"), os.Getenv("NODE_ENV"), os.Getenv("LOCAL_RELEASE"),
		os.Getenv("SRS_DOCKER"), os.Getenv("USE_DOCKER"), os.Getenv("SRS_UTEST"),
		len(os.Getenv("REDIS_PASSWORD")), os.Getenv("REDIS_PORT"), os.Getenv("PUBLIC_URL"),
		os.Getenv("BUILD_PATH"), os.Getenv("REACT_APP_LOCALE"), os.Getenv("MGMT_LISTEN"),
		os.Getenv("PLATFORM_DOCKER"),
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
	redisManager := NewEmptyRedisManager()

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

	// We must initialize the mgmt after redis is ready.
	if err := initMmgt(ctx, redisManager); err != nil {
		return errors.Wrapf(err, "init mgmt")
	}
	logger.Tf(ctx, "initialize mgmt region=%v, registry=%v, version=%v", conf.Region, conf.Registry, version)

	// Note that we always restart the platform container.
	if err := removeContainer(ctx, platformDockerName); err != nil {
		logger.Tf(ctx, "ignore restart platform err %v", err)
	}

	// Create backend service manager.
	service := NewDockerBackendService(redisManager, NewDockerPlatformManager(redisManager))
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
func initMmgt(ctx context.Context, redisManager RedisManager) error {
	// Refresh the env file.
	if envs, err := godotenv.Read(".env"); err != nil {
		return errors.Wrapf(err, "load envs")
	} else {
		envs["CLOUD"] = conf.Cloud
		envs["REGION"] = conf.Region
		envs["SOURCE"] = conf.Source
		envs["REGISTRY"] = conf.Registry
		if os.Getenv("MGMT_PASSWORD") != "" {
			envs["MGMT_PASSWORD"] = os.Getenv("MGMT_PASSWORD")
		}

		if err := godotenv.Write(envs, ".env"); err != nil {
			return errors.Wrapf(err, "write .env")
		}
		logger.Tf(ctx, "Refresh .env ok")
	}

	return nil
}
