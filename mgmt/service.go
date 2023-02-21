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
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	ohttp "github.com/ossrs/go-oryx-lib/http"
	"github.com/ossrs/go-oryx-lib/logger"

	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
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
	return fmt.Sprintf("darwin=%v, cloud=%v, region=%v, source=%v, registry=%v, iface=%v, ipv4=%v, pwd=%v",
		v.IsDarwin, v.Cloud, v.Region, v.Source, v.Registry, v.Iface, v.IPv4(), v.Pwd,
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

func NewDockerRedisManager() RedisManager {
	return &dockerRedisManager{}
}

type dockerRedisManager struct {
}

func (v *dockerRedisManager) Stop(ctx context.Context, timeout time.Duration) error {
	args := []string{"stop", "-t", fmt.Sprintf("%v", int64(timeout/time.Second)), redisDockerName}
	cmd := exec.CommandContext(ctx, "docker", args...)

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "docker %v", strings.Join(args, " "))
	}
	return nil
}

func (v *dockerRedisManager) Start(ctx context.Context) error {
	all, _ := queryContainer(ctx, redisDockerName)
	if all != nil && all.ID != "" {
		err := removeContainer(ctx, redisDockerName)
		logger.Tf(ctx, "docker remove name=%v, id=%v, err=%v", redisDockerName, all.ID, err)
	}

	args := []string{
		"run", "-d", "--restart=always", "--privileged", fmt.Sprintf("--name=%v", redisDockerName),
		"--env", fmt.Sprintf("SRS_REGION=%v", conf.Region),
		"--env", fmt.Sprintf("SRS_SOURCE=%v", conf.Source),
		"--log-driver=json-file", "--log-opt=max-size=1g", "--log-opt=max-file=3",
		"-v", fmt.Sprintf("%v/containers/data/redis:/data", conf.Pwd),
		"-v", fmt.Sprintf("%v/containers/conf/redis.conf:/etc/redis/redis.conf", conf.Pwd),
		"-p", fmt.Sprintf("%v:%v/tcp", os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PORT")),
	}
	if !conf.IsDarwin {
		args = append(args, "--network=srs-cloud")
	}
	args = append(args,
		fmt.Sprintf("%v/ossrs/redis", conf.Registry),
		"redis-server",
		"/etc/redis/redis.conf",
	)
	if os.Getenv("REDIS_PASSWORD") != "" {
		args = append(args, "--requirepass", os.Getenv("REDIS_PASSWORD"))
	}
	args = append(args, "--port", os.Getenv("REDIS_PORT"))

	cmd := exec.CommandContext(ctx, "docker", args...)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "docker %v", strings.Join(args, " "))
	}

	logger.Tf(ctx, "docker %v", strings.Join(args, " "))
	return nil
}

func (v *dockerRedisManager) Ready(ctx context.Context) error {
	for {
		all, running := queryContainer(ctx, redisDockerName)
		if all != nil && all.ID != "" && running != nil && running.ID != "" {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(300 * time.Millisecond):
		}
	}
}

func NewDockerPlatformManager() PlatformManager {
	return &dockerPlatformManager{}
}

type dockerPlatformManager struct {
}

func (v *dockerPlatformManager) Start(ctx context.Context) error {
	all, _ := queryContainer(ctx, platformDockerName)
	if all != nil && all.ID != "" {
		err := removeContainer(ctx, platformDockerName)
		logger.Tf(ctx, "docker remove name=%v, id=%v, err=%v", platformDockerName, all.ID, err)
	}

	args := []string{
		"run", "-d", "--restart=always", "--privileged", fmt.Sprintf("--name=%v", platformDockerName),
		"--env", fmt.Sprintf("SRS_REGION=%v", conf.Region),
		"--env", fmt.Sprintf("SRS_SOURCE=%v", conf.Source),
		"--log-driver=json-file", "--log-opt=max-size=1g", "--log-opt=max-file=3",
		// Note that we should mount .env to mgmt, because in platform only use this path.
		"-v", fmt.Sprintf("%v/.env:/usr/local/srs-cloud/mgmt/.env", conf.Pwd),
		// We mount the containers to mgmt in platform container, which links to platform.
		"-v", fmt.Sprintf("%v/containers:/usr/local/srs-cloud/mgmt/containers", conf.Pwd),
		"--env", fmt.Sprintf("SRS_DOCKER=%v", os.Getenv("SRS_DOCKER")),
		"--env", fmt.Sprintf("USE_DOCKER=%v", os.Getenv("USE_DOCKER")),
		// If use docker, should always use production to connect to redis.
		"--env", "NODE_ENV=production",
		"-p", "2024:2024/tcp",
		"--add-host", fmt.Sprintf("mgmt.srs.local:%v", conf.IPv4()),
	}
	if !conf.IsDarwin {
		args = append(args, "--network=srs-cloud")
	}
	if conf.IsDarwin {
		args = append(args, "--env", "REDIS_HOST=host.docker.internal")
	} else {
		args = append(args, "--env", "REDIS_HOST=redis")
	}
	args = append(args,
		fmt.Sprintf("%v/ossrs/srs-cloud:platform-%v", conf.Registry, version),
	)

	cmd := exec.CommandContext(ctx, "docker", args...)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "docker %v", strings.Join(args, " "))
	}

	logger.Tf(ctx, "docker %v", strings.Join(args, " "))
	return nil
}

func NewDockerHTTPService() HttpService {
	return &dockerHTTPService{}
}

type dockerHTTPService struct {
	server *http.Server
}

func (v *dockerHTTPService) Close() error {
	server := v.server
	v.server = nil

	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Tf(ctx, "ignore HTTP server shutdown err %v", err)
		}
	}

	return nil
}

func (v *dockerHTTPService) Run(ctx context.Context) error {
	addr := os.Getenv("MGMT_LISTEN")
	if !strings.HasPrefix(addr, ":") {
		addr = fmt.Sprintf(":%v", addr)
	}
	logger.Tf(ctx, "HTTP listen at %v", addr)

	handler := http.NewServeMux()
	if err := handleDockerHTTPService(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle service")
	}

	server := &http.Server{Addr: addr, Handler: handler}
	v.server = server

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		logger.Tf(ctx, "shutting down HTTP server...")
		v.Close()
	}()

	if err := server.ListenAndServe(); err != nil && ctx.Err() != context.Canceled {
		return errors.Wrapf(err, "listen %v", addr)
	}
	logger.Tf(ctx, "HTTP server is done")

	return nil
}

type dockerServerRequest struct {
	Action string        `json:"action"`
	Token  string        `json:"token"`
	Args   []interface{} `json:"args"`
}

func (v *dockerServerRequest) String() string {
	return fmt.Sprintf("action=%v, args=%v", v.Action, v.Args)
}

func (v *dockerServerRequest) ArgsAsString() []string {
	args := []string{}
	for _, arg := range v.Args {
		if s, ok := arg.(string); ok {
			args = append(args, s)
		}
	}
	return args
}

func (v *dockerServerRequest) ArgsAsMap() []map[string]interface{} {
	args := []map[string]interface{}{}
	for _, arg := range v.Args {
		if m, ok := arg.(map[string]interface{}); ok {
			args = append(args, m)
		}
	}
	return args
}

func (v *dockerServerRequest) ArgsAsSlices() [][]interface{} {
	args := [][]interface{}{}
	for _, arg := range v.Args {
		if m, ok := arg.([]interface{}); ok {
			args = append(args, m)
		}
	}
	return args
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

	handlers := make(map[string]func(ctx context.Context, w http.ResponseWriter, r *http.Request, sr *dockerServerRequest) error)

	// Current work directory.
	handlers["cwd"] = func(ctx context.Context, w http.ResponseWriter, r *http.Request, sr *dockerServerRequest) error {
		ohttp.WriteData(ctx, w, r, &struct {
			Cwd string `json:"cwd"`
		}{
			conf.Pwd,
		})
		logger.Tf(ctx, "execApi req=%v, pwd=%v", sr, conf.Pwd)
		return nil
	}

	// Current host platform name.
	handlers["hostPlatform"] = func(ctx context.Context, w http.ResponseWriter, r *http.Request, sr *dockerServerRequest) error {
		// The platform must be the os name, such as darwin or linux, equals to nodejs process.platform, for Go it
		// should be runtime.GOOS, not the conf.Platform which is for statistic only.
		ohttp.WriteData(ctx, w, r, &struct {
			Platform string `json:"platform"`
		}{
			Platform: runtime.GOOS,
		})
		logger.Tf(ctx, "execApi req=%v,platform=%v", sr, runtime.GOOS)
		return nil
	}

	// Fetch the container.
	handlers["fetchContainer"] = func(ctx context.Context, w http.ResponseWriter, r *http.Request, sr *dockerServerRequest) error {
		name := sr.ArgsAsString()[0]
		if name == "" {
			return errors.New("no name")
		}

		all, running := queryContainer(ctx, name)
		ohttp.WriteData(ctx, w, r, &struct {
			All     *dockerInfo `json:"all"`
			Running *dockerInfo `json:"running"`
		}{
			All: all, Running: running,
		})

		logger.Tf(ctx, "execApi req=%v, all=%v, running=%v", sr, all, running)
		return nil
	}

	// Remove the container.
	handlers["removeContainer"] = func(ctx context.Context, w http.ResponseWriter, r *http.Request, sr *dockerServerRequest) error {
		name := sr.ArgsAsString()[0]
		if name == "" {
			return errors.New("no name")
		}

		err := removeContainer(ctx, name)
		if err != nil {
			return errors.Wrapf(err, "remove container %v", name)
		}

    ohttp.WriteData(ctx, w, r, nil)
		logger.Tf(ctx, "execApi req=%v", sr)
		return nil
	}

	// Query all the containers, ignore if not exists.
	handlers["queryContainers"] = func(ctx context.Context, w http.ResponseWriter, r *http.Request, sr *dockerServerRequest) error {
		names := sr.ArgsAsString()
		if len(names) == 0 {
			return errors.Errorf("no name")
		}

		// Convert names to docker names.
		var dockerNames []string
		for _, name := range names {
			if name == "srs" {
				dockerNames = append(dockerNames, srsDockerName)
			} else if name == "srsDev" {
				return errors.New("srs dev is not supported")
			} else {
				dockerNames = append(dockerNames, name)
			}
		}

		containers := []interface{}{}
		for _, name := range dockerNames {
			disabled, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, name).Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_CONTAINER_DISABLED, name)
			}

			r0 := &struct {
				Name      string `json:"name"`
				Enabled   bool   `json:"enabled"`
				Container struct {
					ID     string `json:"ID"`
					State  string `json:"State"`
					Status string `json:"Status"`
				} `json:"container"`
			}{
				Name: name, Enabled: disabled != "true",
			}
			all, _ := queryContainer(ctx, name)
			if all != nil {
				r0.Container.ID = all.ID
				r0.Container.State = all.State
				r0.Container.Status = all.Status
			}

			containers = append(containers, r0)
		}

		ohttp.WriteData(ctx, w, r, &struct {
			Containers interface{} `json:"containers"`
		}{
			Containers: containers,
		})
		logger.Tf(ctx, "execApi req=%v, containers=%v", sr, len(containers))
		return nil
	}

	// Reload NGINX.
	handlers["reloadNginx"] = func(ctx context.Context, w http.ResponseWriter, r *http.Request, sr *dockerServerRequest) error {
		if err := reloadNginx(ctx); err != nil {
			return errors.Wrapf(err, "reload nginx")
		}
		logger.Tf(ctx, "NGINX: Refresh dynamic.conf ok")

		ohttp.WriteData(ctx, w, r, nil)
		logger.Tf(ctx, "execApi req=%v", sr)
		return nil
	}

	// Remove the specified container.
	handlers["rmContainer"] = func(ctx context.Context, w http.ResponseWriter, r *http.Request, sr *dockerServerRequest) error {
		name := sr.ArgsAsString()[0]
		if name == "" {
			return errors.Errorf("no name")
		}

		if err := removeContainer(ctx, name); err != nil {
			return errors.Wrapf(err, "remove container name=%v", name)
		}

		ohttp.WriteData(ctx, w, r, nil)
		logger.Tf(ctx, "execApi req=%v", sr)
		return nil
	}

	// Current ipv4 internal address.
	handlers["ipv4"] = func(ctx context.Context, w http.ResponseWriter, r *http.Request, sr *dockerServerRequest) error {
		ohttp.WriteData(ctx, w, r, &struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		}{
			Name: conf.Iface, Address: conf.IPv4(),
		})
		logger.Tf(ctx, "execApi req=%v, iface=%v, ipv4=%v", sr, conf.Iface, conf.IPv4())
		return nil
	}

	// Start the container with args.
	handlers["startContainer"] = func(ctx context.Context, w http.ResponseWriter, r *http.Request, sr *dockerServerRequest) error {
		name, dockerArgs := sr.ArgsAsString()[0], sr.ArgsAsSlices()[0]
		if name == "" {
			return errors.New("no name")
		}
		if dockerArgs == nil {
			return errors.New("no args")
		}

		var args []string
		for _, a := range dockerArgs {
			if s, ok := a.(string); ok {
				args = append(args, s)
			}
		}

		if err := removeContainer(ctx, name); err != nil {
			logger.Tf(ctx, "ignore remove docker name=%v err %v", name, err)
		}
		if err := exec.CommandContext(ctx, "docker", args...).Run(); err != nil {
			return errors.Wrapf(err, "docker %v", strings.Join(args, " "))
		}

		ohttp.WriteData(ctx, w, r, nil)
		logger.Tf(ctx, "execApi req=%v", sr)
		return nil
	}

	// Reload the env from .env
	handlers["reloadEnv"] = func(ctx context.Context, w http.ResponseWriter, r *http.Request, sr *dockerServerRequest) error {
		if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
			return errors.Wrapf(err, "load .env")
		}

		ohttp.WriteData(ctx, w, r, nil)
		logger.Tf(ctx, "execApi req=%v", sr)
		return nil
	}

	// Start upgrade.
	handlers["execUpgrade"] = func(ctx context.Context, w http.ResponseWriter, r *http.Request, sr *dockerServerRequest) error {
		target := sr.ArgsAsString()[0]
		if target == "" {
			return errors.New("no target")
		}

		cmd := exec.CommandContext(ctx, "bash", "upgrade", target)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return errors.Wrapf(err, "pipe stdout")
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return errors.Wrapf(err, "pipe stderr")
		}

		logger.Tf(ctx, "start upgrade to %v", target)
		if err = cmd.Start(); err != nil {
			return errors.Wrapf(err, "start upgrade")
		}

		rs := io.MultiReader(stdout, stderr)
		buf := make([]byte, 4096)
		for {
			if nn, err := rs.Read(buf); err != nil || nn == 0 {
				if err != io.EOF {
					logger.Tf(ctx, "read nn=%v, err %v", nn, err)
				}
				break
			} else if s := buf[:nn]; true {
				logger.Tf(ctx, "%v", string(s))
			}
		}

		if err = cmd.Wait(); err != nil {
			logger.Tf(ctx, "wait err %v", err)
		}

		ohttp.WriteData(ctx, w, r, nil)
		logger.Tf(ctx, "execApi req=%v target=%v", sr, target)
		return nil
	}

	ep = "/terraform/v1/host/exec"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return err
			}

			req := &dockerServerRequest{}
			if err = json.Unmarshal(b, req); err != nil {
				return errors.Wrapf(err, "parse %v", string(b))
			}
			logger.Tf(ctx, "exec action=%v, token=%v, args=%v", req.Action, req.Token, req.Args)

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			// Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
			// See https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
			if _, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
				return []byte(apiSecret), nil
			}); err != nil {
				return errors.Wrapf(err, "verify token %v", req.Token)
			}

			if fn, ok := handlers[req.Action]; !ok {
				return errors.Errorf("no handler for action=%v, args=%v", req.Action, req.Args)
			} else {
				if err := fn(ctx, w, r, req); err != nil {
					return errors.Wrapf(err, "handle action=%v, args=%v", req.Action, req.Args)
				}
			}

			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

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

	proxy2024, err := createProxy("http://127.0.0.1:2024")
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

	fileServer := http.FileServer(http.Dir(path.Join(conf.Pwd, "containers/www")))

	ep = "/"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
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
			fileServer.ServeHTTP(w, r)
			return
		}

		// For registered modules, by /terraform/v1/tencent/
		// Note that this module has been migrated to platform.
		// For registered modules, by /terraform/v1/ffmpeg/
		// Note that this module has been migrated to platform.
		// For platform apis, by /terraform/v1/mgmt/
		// TODO: FIXME: Proxy all mgmt APIs to platform.
		// For registered modules, by /terraform/v1/hooks/
		// Note that this module has been migrated to platform.
		if strings.HasPrefix(r.URL.Path, "/terraform/") {
			logger.Tf(ctx, "Proxy %v to backend 2024", r.URL.Path)
			proxy2024.ServeHTTP(w, r)
			return
		}

		// The UI proxy to platform UI, system mgmt UI.
		if strings.HasPrefix(r.URL.Path, "/mgmt") {
			logger.Tf(ctx, "Proxy %v to backend 2024", r.URL.Path)
			proxy2024.ServeHTTP(w, r)
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

		w.Write([]byte("Hello world!"))
	})

	return nil
}

func NewDockerBackendService(r RedisManager, p PlatformManager) BackendService {
	return &dockerBackendService{
		redisManager: r, platformManager: p,
	}
}

type dockerBackendService struct {
	wg sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc

	redisManager    RedisManager
	platformManager PlatformManager
}

func (v *dockerBackendService) Close() error {
	if v.cancel != nil {
		v.cancel()
	}

	logger.Tf(v.ctx, "backend service quiting...")
	v.wg.Wait()
	logger.Tf(v.ctx, "backend service done")

	return nil
}

func (v *dockerBackendService) Start(ctx context.Context) error {
	ctx = logger.WithContext(ctx)
	v.ctx, v.cancel = context.WithCancel(ctx)

	// Manage the redis service by docker.
	v.wg.Add(1)
	go func() {
		defer v.wg.Done()

		for ctx.Err() == nil {
			duration := 10 * time.Second
			if err := func() error {
				all, running := queryContainer(ctx, redisDockerName)
				if all != nil && all.ID != "" && running != nil && running.ID != "" {
					logger.Tf(ctx, "query ID=%v, State=%v, Status=%v, running=%v", all.ID, all.State, all.Status, running.ID)
					return nil
				}

				if disabled, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, redisDockerName).Result(); err != nil && err != redis.Nil {
					return err
				} else if disabled == "true" {
					logger.Tf(ctx, "container %v disabled", redisDockerName)
					return nil
				}

				logger.Tf(ctx, "restart redis container")
				if err := v.redisManager.Start(ctx); err != nil {
					return err
				}

				return nil
			}(); err != nil {
				duration = 30 * time.Second
				logger.Tf(ctx, "ignore err %v", err)
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(duration):
			}
		}
	}()

	// Manage the platform service by docker.
	v.wg.Add(1)
	go func() {
		defer v.wg.Done()

		for ctx.Err() == nil {
			duration := 10 * time.Second
			if err := func() error {
				all, running := queryContainer(ctx, platformDockerName)
				if all != nil && all.ID != "" && running != nil && running.ID != "" {
					logger.Tf(ctx, "query ID=%v, State=%v, Status=%v, running=%v", all.ID, all.State, all.Status, running.ID)
					return nil
				}

				if disabled, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, platformDockerName).Result(); err != nil && err != redis.Nil {
					return err
				} else if disabled == "true" {
					logger.Tf(ctx, "container %v disabled", platformDockerName)
					return nil
				}

				logger.Tf(ctx, "restart platform container")
				if err := v.platformManager.Start(ctx); err != nil {
					return err
				}

				// Cleanup the previous unused images.
				newImage := fmt.Sprintf("%v/ossrs/srs-cloud:platform-%v", conf.Registry, version)
				if previousImage, err := rdb.HGet(ctx, SRS_DOCKER_IMAGES, platformDockerName).Result(); err != nil && err != redis.Nil {
					return err
				} else {
					if err = rdb.HSet(ctx, SRS_DOCKER_IMAGES, platformDockerName, newImage).Err(); err != nil {
						return err
					}
					if previousImage != "" && previousImage != newImage {
						r0 := exec.CommandContext(ctx, "docker", "rmi", previousImage).Run()
						logger.Tf(ctx, "remove previous platform image %v, err %v", previousImage, r0)
					}
				}

				return nil
			}(); err != nil {
				duration = 30 * time.Second
				logger.Tf(ctx, "ignore err %v", err)
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(duration):
			}
		}
	}()

	return nil
}
