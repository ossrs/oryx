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
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"

	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
)

// BackendService is a manager for backend service, like redis and platform.
type BackendService interface {
	Start(ctx context.Context) error
	Close() error
}

// HttpService is a HTTP server for mgmt and platform.
type HttpService interface {
	Close() error
	Run(ctx context.Context) error
}

// PlatformManager is platform service.
type PlatformManager interface {
	// Start the platform server.
	Start(ctx context.Context) error
}

// RedisManager is Redis based on docker or CLI.
type RedisManager interface {
	// Stop the redis server.
	// Note that we should never use 'docker rm -f redis" or data maybe discard. Instead, we should use command similar
	// to 'docker stop redis' to allow redis to save data to disk.
	Stop(ctx context.Context, timeout time.Duration) error
	// Start the redis server.
	// Now we start redis again, to keep it with the latest configurations and params.
	Start(ctx context.Context) error
	// Ready to wait for redis server ready.
	Ready(ctx context.Context) error
}

// Config is for configuration.
type Config struct {
	IsDarwin bool
	Pwd      string

	Cloud    string
	Region   string
	Source   string
	Registry string

	ipv4  net.IP
	Iface string

	// The platform for SRS cloud, not the GOOS, for report and statistic only.
	Platform string
}

func NewConfig() *Config {
	return &Config{
		ipv4:     net.IPv4zero,
		IsDarwin: runtime.GOOS == "darwin",
	}
}

func (v *Config) IPv4() string {
	return v.ipv4.String()
}

func (v *Config) String() string {
	return fmt.Sprintf("darwin=%v, cloud=%v, region=%v, source=%v, registry=%v, iface=%v, ipv4=%v, pwd=%v, platform=%v",
		v.IsDarwin, v.Cloud, v.Region, v.Source, v.Registry, v.Iface, v.IPv4(), v.Pwd, v.Platform,
	)
}

// conf is a global config object.
var conf *Config

func discoverRegion(ctx context.Context) (cloud, region string, err error) {
	if conf.IsDarwin {
		return "DEV", "ap-beijing", nil
	}

	if os.Getenv("CLOUD") == "BT" {
		return "BT", "ap-beijing", nil
	}

	if os.Getenv("CLOUD") == "AAPANEL" {
		return "AAPANEL", "ap-singapore", nil
	}

	if os.Getenv("CLOUD") != "" && os.Getenv("REGION") != "" {
		return os.Getenv("CLOUD"), os.Getenv("REGION"), nil
	}

	logger.Tf(ctx, "Initialize start to discover region")

	var wg sync.WaitGroup
	defer wg.Wait()

	discoverCtx, discoverCancel := context.WithCancel(ctx)
	result := make(chan Config, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()

		res, err := http.Get("http://metadata.tencentyun.com/latest/meta-data/placement/region")
		if err != nil {
			logger.Tf(ctx, "Ignore tencent region err %v", err)
			return
		}
		defer res.Body.Close()

		b, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Tf(ctx, "Ignore tencent region err %v", err)
			return
		}

		select {
		case <-discoverCtx.Done():
		case result <- Config{Cloud: "TENCENT", Region: string(b)}:
			discoverCancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// See https://docs.digitalocean.com/reference/api/metadata-api/#operation/getRegion
		res, err := http.Get("http://169.254.169.254/metadata/v1/region")
		if err != nil {
			logger.Tf(ctx, "Ignore do region err %v", err)
			return
		}
		defer res.Body.Close()

		b, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Tf(ctx, "Ignore do region err %v", err)
			return
		}

		select {
		case <-discoverCtx.Done():
		case result <- Config{Cloud: "DO", Region: string(b)}:
			discoverCancel()
		}
	}()

	select {
	case <-ctx.Done():
	case r := <-result:
		return r.Cloud, r.Region, nil
	}
	return
}

func discoverSource(ctx context.Context, cloud, region string) (source string, err error) {
	switch cloud {
	case "DEV", "BT":
		return "gitee", nil
	case "DO", "AAPANEL":
		return "github", nil
	}

	for _, r := range []string{
		"ap-guangzhou", "ap-shanghai", "ap-nanjing", "ap-beijing", "ap-chengdu", "ap-chongqing",
	} {
		if strings.HasPrefix(region, r) {
			return "gitee", nil
		}
	}
	return "github", nil
}

func discoverRegistry(ctx context.Context, source string) (registry string, err error) {
	if source == "github" {
		return "docker.io", nil
	}
	return "registry.cn-hangzhou.aliyuncs.com", nil
}

func discoverPlatform(ctx context.Context, cloud string) (platform string, err error) {
	switch cloud {
	case "DEV":
		return "dev", nil
	case "DO":
		return "droplet", nil
	case "BT":
		return "bt", nil
	case "AAPANEL":
		return "aapanel", nil
	}

	// Discover CVM or lighthouse.
	res, err := http.Get("http://metadata.tencentyun.com/latest/meta-data/instance-name")
	if err != nil {
		logger.Tf(ctx, "Ignore tencent platform err %v", err)
		return "dev", nil
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Tf(ctx, "Ignore tencent platform err %v", err)
		return "dev", nil
	}

	if strings.Contains(string(b), "-lhins-") {
		return "lighthouse", nil
	}
	return "cvm", nil
}

func discoverPrivateIPv4(ctx context.Context) (string, net.IP, error) {
	candidates := make(map[string]net.IP)

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", nil, err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", nil, err
		}

		for _, addr := range addrs {
			if addr, ok := addr.(*net.IPNet); ok {
				if addr.IP.To4() != nil && !addr.IP.IsLoopback() {
					candidates[iface.Name] = addr.IP
				}
			}
		}
	}

	var bestMatch string
	for name, _ := range candidates {
		if strings.HasPrefix(name, "en") || strings.HasPrefix(name, "eth") {
			bestMatch = name
			break
		}
	}

	if bestMatch == "" {
		for name, _ := range candidates {
			bestMatch = name
			break
		}
	}

	var privateIPv4 net.IP
	if bestMatch != "" {
		if addr, ok := candidates[bestMatch]; ok {
			privateIPv4 = addr
		}
	}

	logger.Tf(ctx, "Refresh ipv4=%v, bestMatch=%v, candidates=%v", privateIPv4, bestMatch, candidates)
	return bestMatch, privateIPv4, nil
}

// Docker container names.
const redisDockerName = "redis"
const platformDockerName = "platform"
const srsDockerName = "srs-server"

// Note that we only enable srs-server, never enable srs-dev.
const srsDevDockerName = "srs-dev"

// Note that we only use the docker to release binary for different CPU archs.
const mgmtDockerName = "mgmt"

// Redis keys.
const (
	// For LightHouse information, like region or source.
	SRS_TENCENT_LH = "SRS_TENCENT_LH"
	// For container and images.
	SRS_CONTAINER_DISABLED = "SRS_CONTAINER_DISABLED"
	SRS_DOCKER_IMAGES      = "SRS_DOCKER_IMAGES"
	// For system settings.
	SRS_PLATFORM_SECRET = "SRS_PLATFORM_SECRET"
	SRS_UPGRADING       = "SRS_UPGRADING"
	SRS_HTTPS           = "SRS_HTTPS"
)

// rdb is a global redis client object.
var rdb *redis.Client

// InitRdb create and init global rdb, which is a redis client.
func InitRdb() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%v", os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	return nil
}

// For docker state.
type dockerInfo struct {
	Command      string `json:"Command"`
	CreatedAt    string `json:"CreatedAt"`
	ID           string `json:"ID"`
	Image        string `json:"Image"`
	Labels       string `json:"Labels"`
	LocalVolumes string `json:"LocalVolumes"`
	Mounts       string `json:"Mounts"`
	Names        string `json:"Names"`
	Networks     string `json:"Networks"`
	Ports        string `json:"Ports"`
	RunningFor   string `json:"RunningFor"`
	Size         string `json:"Size"`
	State        string `json:"State"`
	Status       string `json:"Status"`
}

func (v *dockerInfo) String() string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("ID=%v, State=%v, Status=%v", v.ID, v.State, v.Status)
}

// queryContainer used to query the state of docker container.
func queryContainer(ctx context.Context, name string) (all, running *dockerInfo) {
	if true {
		cmd := exec.CommandContext(ctx, "docker",
			"ps", "-a", "-f", fmt.Sprintf("name=%v", name), "--format", "'{{json .}}'",
		)
		b, err := cmd.Output()
		s := strings.Trim(strings.TrimSpace(string(b)), "'")
		if err == nil && s != "" {
			all = &dockerInfo{}
			if err = json.Unmarshal([]byte(s), all); err != nil {
				logger.Tf(ctx, "ignore parse %v err %v", s, err)
			}
		}
	}

	if true {
		cmd := exec.CommandContext(ctx, "docker",
			"ps", "-f", fmt.Sprintf("name=%v", name), "--format", "'{{json .}}'",
		)
		b, err := cmd.Output()
		s := strings.Trim(strings.TrimSpace(string(b)), "'")
		if err == nil && s != "" {
			running = &dockerInfo{}
			if err = json.Unmarshal([]byte(s), running); err != nil {
				logger.Tf(ctx, "ignore parse %v err %v", s, err)
			}
		}
	}

	return
}

// removeContainer used to remove the docker container.
func removeContainer(ctx context.Context, name string) error {
	cmd := exec.CommandContext(ctx, "docker", "rm", "-f", name)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "docker rm -f %v", name)
	}
	return nil
}

// setEnvDefault set env key=value if not set.
func setEnvDefault(key, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}

// reloadNginx is used to reload the NGINX server.
func reloadNginx(ctx context.Context) error {
	if conf.IsDarwin {
		return nil
	}

	var nginxServiceExists, nginxPidExists bool
	if _, err := os.Stat("/usr/lib/systemd/system/nginx.service"); err == nil {
		nginxServiceExists = true
	}
	if os.Getenv("NGINX_PID") != "" {
		if _, err := os.Stat(os.Getenv("NGINX_PID")); err == nil {
			nginxPidExists = true
		}
	}
	if !nginxServiceExists && !nginxPidExists {
		return errors.Errorf("Can't reload NGINX, no service or pid %v", os.Getenv("NGINX_PID"))
	}

	// Try to reload by service if exists, try pid if failed.
	if nginxServiceExists {
		if err := exec.CommandContext(ctx, "systemctl", "reload", "nginx.service").Run(); err != nil {
			if !nginxPidExists {
				return errors.Wrapf(err, "reload nginx failed, service=%v, pid=%v", nginxServiceExists, nginxPidExists)
			}
		} else {
			return nil
		}
	}

	if b, err := ioutil.ReadFile(os.Getenv("NGINX_PID")); err != nil {
		return errors.Wrapf(err, "read nginx pid from %v", os.Getenv("NGINX_PID"))
	} else if pid := strings.TrimSpace(string(b)); pid == "" {
		return errors.Errorf("no pid at %v", os.Getenv("NGINX_PID"))
	} else if err = exec.CommandContext(ctx, "kill", "-s", "SIGHUP", pid).Run(); err != nil {
		return errors.Wrapf(err, "reload nginx failed pid=%v", pid)
	}

	return nil
}
