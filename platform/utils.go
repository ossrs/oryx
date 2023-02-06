//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: MIT
//
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"

	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
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

// SrsManager is SRS based on docker or CLI.
type SrsManager interface {
	// Stop the SRS server.
	// Note that we should never use 'docker rm -f srs" or data maybe discard. Instead, we should use command similar
	// to 'docker stop srs' to allow SRS to save data to disk.
	Stop(ctx context.Context, timeout time.Duration) error
	// Start the SRS server.
	// Now we start SRS again, to keep it with the latest configurations and params.
	Start(ctx context.Context) error
	// Ready to wait for SRS server ready.
	Ready(ctx context.Context) error
}

// Versions is latest and stable version from SRS cloud API.
type Versions struct {
	Version string `json:"version"`
	Stable  string `json:"stable"`
	Latest  string `json:"latest"`
}

func (v Versions) String() string {
	return fmt.Sprintf("version=%v, latest=%v, stable=%v", v.Version, v.Latest, v.Stable)
}

// Config is for configuration.
// TODO: FIXME: Should be merged to mgmt.
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
	// The latest and stable version from SRS cloud API.
	Versions Versions
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
	return fmt.Sprintf("darwin=%v, cloud=%v, region=%v, source=%v, registry=%v, iface=%v, ipv4=%v, pwd=%v, "+
		"platform=%v, version=%v, latest=%v, stable=%v",
		v.IsDarwin, v.Cloud, v.Region, v.Source, v.Registry, v.Iface, v.IPv4(), v.Pwd, v.Platform, v.Versions.Version,
		v.Versions.Latest, v.Versions.Stable,
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
	// For tencent cloud products.
	SRS_TENCENT_CAM = "SRS_TENCENT_CAM"
	SRS_TENCENT_COS = "SRS_TENCENT_COS"
	SRS_TENCENT_VOD = "SRS_TENCENT_VOD"
	// For local record.
	SRS_RECORD_PATTERNS      = "SRS_RECORD_PATTERNS"
	SRS_RECORD_M3U8_WORKING  = "SRS_RECORD_M3U8_WORKING"
	SRS_RECORD_M3U8_ARTIFACT = "SRS_RECORD_M3U8_ARTIFACT"
	// For cloud storage.
	SRS_DVR_PATTERNS      = "SRS_DVR_PATTERNS"
	SRS_DVR_M3U8_WORKING  = "SRS_DVR_M3U8_WORKING"
	SRS_DVR_M3U8_ARTIFACT = "SRS_DVR_M3U8_ARTIFACT"
	// For cloud VoD.
	SRS_VOD_PATTERNS      = "SRS_VOD_PATTERNS"
	SRS_VOD_M3U8_WORKING  = "SRS_VOD_M3U8_WORKING"
	SRS_VOD_M3U8_ARTIFACT = "SRS_VOD_M3U8_ARTIFACT"
	// The cos token and file information for cloud VoD, to upload files.
	SRS_VOD_COS_TOKEN = "SRS_VOD_COS_TOKEN"
	// For stream forwarding by FFmpeg.
	SRS_FORWARD_CONFIG = "SRS_FORWARD_CONFIG"
	SRS_FORWARD_TASK   = "SRS_FORWARD_TASK"
	// For virtual live channel/stream.
	SRS_VLIVE_CONFIG = "SRS_VLIVE_CONFIG"
	SRS_VLIVE_TASK    = "SRS_VLIVE_TASK"
	// For SRS stream status.
	SRS_STREAM_ACTIVE     = "SRS_STREAM_ACTIVE"
	SRS_STREAM_SRT_ACTIVE = "SRS_STREAM_SRT_ACTIVE"
	SRS_STREAM_RTC_ACTIVE = "SRS_STREAM_RTC_ACTIVE"
	// For feature statistics.
	SRS_STAT_COUNTER = "SRS_STAT_COUNTER"
	// For container and images.
	SRS_CONTAINER_DISABLED = "SRS_CONTAINER_DISABLED"
	// For system settings.
	SRS_SECRET_PUBLISH  = "SRS_SECRET_PUBLISH"
	SRS_AUTH_SECRET     = "SRS_AUTH_SECRET"
	SRS_FIRST_BOOT      = "SRS_FIRST_BOOT"
	SRS_UPGRADE_WINDOW  = "SRS_UPGRADE_WINDOW"
	SRS_PLATFORM_SECRET = "SRS_PLATFORM_SECRET"
	SRS_CACHE_BILIBILI  = "SRS_CACHE_BILIBILI"
	SRS_BEIAN           = "SRS_BEIAN"
	SRS_HTTPS           = "SRS_HTTPS"
)

// Tencent cloud consts.
const (
	TENCENT_CLOUD_CAM_ENDPOINT = "cam.tencentcloudapi.com"
	TENCENT_CLOUD_VOD_ENDPOINT = "vod.tencentcloudapi.com"
)

// For vLive upload directory.
var dirUploadPath = path.Join(".", "upload")
var dirVLivePath = path.Join(".", "vlive");

// rdb is a global redis client object.
var rdb *redis.Client

// InitRdb create and init global rdb, which is a redis client.
func InitRdb() error {
	addr := "localhost"
	if os.Getenv("NODE_ENV") != "development" {
		if os.Getenv("REDIS_HOST") != "" {
			addr = os.Getenv("REDIS_HOST")
		} else {
			addr = "mgmt.srs.local"
		}
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", addr, os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	return nil
}

// For platform to build token by jwt.
func createToken(ctx context.Context, apiSecret string) (expireAt, createAt time.Time, token string, err error) {
	createAt, expireAt = time.Now(), time.Now().Add(365*24*time.Hour)

	claims := struct {
		Version string `json:"v"`
		Nonce   string `json:"nonce"`
		jwt.RegisteredClaims
	}{
		Version: "1.0",
		Nonce:   fmt.Sprintf("%x", rand.Uint64()),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(createAt),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(
		[]byte(apiSecret),
	)
	if err != nil {
		return expireAt, createAt, "", errors.Wrapf(err, "jwt sign")
	}

	return expireAt, createAt, token, nil
}

// For platform to execute RPC by HTTP API.
func execApi(ctx context.Context, action string, args interface{}, resObj interface{}) error {
	_, _, token, err := createToken(ctx, os.Getenv("SRS_PLATFORM_SECRET"))
	if err != nil {
		return errors.Wrapf(err, "build token")
	}

	server := "localhost"
	if os.Getenv("NODE_ENV") != "development" {
		server = "mgmt.srs.local"
	}

	body, err := json.Marshal(&struct {
		Token  string      `json:"token"`
		Action string      `json:"action"`
		Args   interface{} `json:"args"`
	}{
		Token: token, Action: action, Args: args,
	})
	if err != nil {
		return errors.Wrapf(err, "build request")
	}

	requestURL := fmt.Sprintf("http://%v:2022/terraform/v1/host/exec", server)
	res, err := http.Post(requestURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return errors.Wrapf(err, "request %v with %v", requestURL, string(body))
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrapf(err, "read response")
	}

	r1 := &struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
	}{
		Data: resObj,
	}
	if err = json.Unmarshal(b, r1); err != nil {
		return errors.Wrapf(err, "unmarshal json from %v", string(b))
	}
	if r1.Code != 0 {
		return errors.Errorf("error response code=%v, %v, action=%v", r1.Code, string(b), action)
	}

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

// setEnvDefault set env key=value if not set.
func setEnvDefault(key, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}

// queryLatestVersion is to query the latest and stable version from SRS cloud API.
func queryLatestVersion(ctx context.Context) (*Versions, error) {
	// Request release api with params.
	params := make(map[string]string)

	// Generate and setup the node id.
	if r0, err := rdb.HGet(ctx, SRS_TENCENT_LH, "node").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v node", SRS_TENCENT_LH)
	} else if r0 != "" {
		params["nid"] = r0
	}

	// Report about local Reocrd.
	if r0, err := rdb.HGet(ctx, SRS_RECORD_PATTERNS, "all").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v all", SRS_RECORD_PATTERNS)
	} else if r0 == "true" {
		params["rkd"] = "1"
	}
	if r0, err := rdb.HLen(ctx, SRS_RECORD_M3U8_ARTIFACT).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_RECORD_M3U8_ARTIFACT)
	} else if r0 > 0 {
		params["rkdn"] = fmt.Sprintf("%v", r0)
	}

	// Report about COS and resource usage.
	if r0, err := rdb.HGet(ctx, SRS_TENCENT_COS, "bucket").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v bucket", SRS_TENCENT_COS)
	} else if r0 == "true" {
		params["cos"] = "1"
	}
	if r0, err := rdb.HLen(ctx, SRS_DVR_M3U8_ARTIFACT).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_DVR_M3U8_ARTIFACT)
	} else if r0 > 0 {
		params["cosn"] = fmt.Sprintf("%v", r0)
	}

	// Report about VoD and resource usage.
	if r0, err := rdb.HGet(ctx, SRS_TENCENT_VOD, "storage").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v storage", SRS_TENCENT_VOD)
	} else if r0 == "true" {
		params["vod"] = "1"
	}
	if r0, err := rdb.HLen(ctx, SRS_DVR_M3U8_ARTIFACT).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_DVR_M3U8_ARTIFACT)
	} else if r0 > 0 {
		params["vodn"] = fmt.Sprintf("%v", r0)
	}

	// Report about FFmpeg forwarding.
	if r0, err := rdb.HLen(ctx, SRS_FORWARD_TASK).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_FORWARD_TASK)
	} else if r0 > 0 {
		params["forward"] = fmt.Sprintf("%v", r0)
	}

	// Report about FFmpeg virtual live from file source.
	if r0, err := rdb.HLen(ctx, SRS_VLIVE_TASK).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_VLIVE_TASK)
	} else if r0 > 0 {
		params["vfile"] = fmt.Sprintf("%v", r0)
	}

	// Report about active streams.
	if r0, err := rdb.HGet(ctx, SRS_STAT_COUNTER, "publish").Int64(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v publish", SRS_STAT_COUNTER)
	} else if r0 > 0 {
		if err = rdb.HSet(ctx, SRS_STAT_COUNTER, "publish", 0).Err(); err != nil && err != redis.Nil {
			return nil, errors.Wrapf(err, "hset %v publish", SRS_STAT_COUNTER)
		}
		params["streams"] = fmt.Sprintf("%v", r0)
	}

	// Report about active players.
	if r0, err := rdb.HGet(ctx, SRS_STAT_COUNTER, "play").Int64(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v play", SRS_STAT_COUNTER)
	} else if r0 > 0 {
		if err = rdb.HSet(ctx, SRS_STAT_COUNTER, "play", 0).Err(); err != nil && err != redis.Nil {
			return nil, errors.Wrapf(err, "hset %v play", SRS_STAT_COUNTER)
		}
		params["players"] = fmt.Sprintf("%v", r0)
	}

	// Report about SRT stream.
	if r0, err := rdb.HLen(ctx, SRS_STREAM_SRT_ACTIVE).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_STREAM_SRT_ACTIVE)
	} else if r0 > 0 {
		params["srt"] = fmt.Sprintf("%v", r0)
	}

	// Report about WebRTC stream.
	if r0, err := rdb.HLen(ctx, SRS_STREAM_RTC_ACTIVE).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_STREAM_RTC_ACTIVE)
	} else if r0 > 0 {
		params["rtc"] = fmt.Sprintf("%v", r0)
	}

	// Report about beian feature.
	if r0, err := rdb.HLen(ctx, SRS_BEIAN).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_BEIAN)
	} else if r0 > 0 {
		params["beian"] = fmt.Sprintf("%v", r0)
	}

	// Report about HTTPS feature.
	if r0, err := rdb.Get(ctx, SRS_HTTPS).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "get %v", SRS_HTTPS)
	} else if r0 != "" {
		params["https"] = r0
	}

	// Report about upgrade window feature.
	if r0, err := rdb.HGet(ctx, SRS_UPGRADE_WINDOW, "update").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v update", SRS_UPGRADE_WINDOW)
	} else if r0 == "true" {
		params["uwin"] = "1"
	}

	// Report whether start as develop environment.
	if os.Getenv("NODE_ENV") == "development" {
		params["dev"] = "1"
	}

	// Report whether enable SRS development version.
	if r0, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, srsDevDockerName).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v %v", SRS_CONTAINER_DISABLED, srsDevDockerName)
	} else if r0 == "false" {
		params["srsd"] = "1"
	}

	// Report about the platform.
	if r0, err := rdb.HGet(ctx, SRS_TENCENT_LH, "platform").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v platform", SRS_TENCENT_LH)
	} else if r0 != "" {
		params["plat"] = r0
	}

	if r0, err := rdb.HGet(ctx, SRS_TENCENT_LH, "cloud").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v cloud", SRS_TENCENT_LH)
	} else if r0 != "" {
		params["cloud"] = r0
	}

	if r0, err := rdb.HGet(ctx, SRS_TENCENT_LH, "region").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v region", SRS_TENCENT_LH)
	} else if r0 != "" {
		params["region"] = r0
	}

	versions := &Versions{}
	if err := execApi(ctx, "refreshVersion", []map[string]string{params}, versions); err != nil {
		return nil, errors.Wrapf(err, "refresh version with %v", params)
	}
	return versions, nil
}

// buildVodM3u8 go generate dynamic m3u8.
func buildVodM3u8(
	ctx context.Context, metadata *M3u8VoDArtifact, absUrl bool, domain string, useKey bool, prefix string,
) (
	contentType, m3u8Body string, duration float64, err error,
) {
	if metadata == nil {
		err = errors.New("no m3u8 metadata")
		return
	}
	if metadata.UUID == "" {
		err = errors.Errorf("no uuid of %v", metadata.String())
		return
	}
	if len(metadata.Files) == 0 {
		err = errors.Errorf("no files of %v", metadata.String())
		return
	}
	if absUrl && metadata.Bucket == "" {
		err = errors.Errorf("no bucket of %v", metadata.String())
		return
	}
	if absUrl && metadata.Region == "" {
		err = errors.Errorf("no region of %v", metadata.String())
		return
	}

	for _, file := range metadata.Files {
		duration += file.Duration
	}

	m3u8 := []string{
		"#EXTM3U",
		"#EXT-X-VERSION:3",
		"#EXT-X-ALLOW-CACHE:YES",
		"#EXT-X-PLAYLIST-TYPE:VOD",
		fmt.Sprintf("#EXT-X-TARGETDURATION:%v", math.Ceil(duration)),
		"#EXT-X-MEDIA-SEQUENCE:0",
	}
	for index, file := range metadata.Files {
		// TODO: FIXME: Identify discontinuity by callback.
		if index < len(metadata.Files)-2 {
			next := metadata.Files[index+1]
			if file.SeqNo+1 != next.SeqNo {
				m3u8 = append(m3u8, "#EXT-X-DISCONTINUITY")
			}
		}

		m3u8 = append(m3u8, fmt.Sprintf("#EXTINF:%.2f, no desc", file.Duration))

		var tsURL string
		if absUrl {
			if domain != "" {
				tsURL = fmt.Sprintf("https://%v/%v", domain, file.Key)
			} else {
				tsURL = fmt.Sprintf("https://%v.cos.%v.myqcloud.com/%v", metadata.Bucket, metadata.Region, file.Key)
			}
		} else {
			if useKey {
				tsURL = fmt.Sprintf("%v%v", prefix, file.Key)
			} else {
				tsURL = fmt.Sprintf("%v%v.ts", prefix, file.TsID)
			}
		}
		m3u8 = append(m3u8, tsURL)
	}
	m3u8 = append(m3u8, "#EXT-X-ENDLIST")

	contentType = "application/vnd.apple.mpegurl"
	m3u8Body = strings.Join(m3u8, "\n")
	return
}

// slicesContains is a function to check whether elem in arr.
func slicesContains(arr []string, elem string) bool {
	for _, e := range arr {
		if e == elem {
			return true
		}
	}
	return false
}

// TsFile is a ts file object.
type TsFile struct {
	// The identify key of TS file, renamed local ts path or COS key, format is record/{m3u8UUID}/{tsID}.ts
	// For example, record/3ECF0239-708C-42E4-96E1-5AE935C6E6A9/5B7B5C03-8DB4-4ABA-AAF3-CB55902CF177.ts
	// For example, 3ECF0239-708C-42E4-96E1-5AE935C6E6A9/5B7B5C03-8DB4-4ABA-AAF3-CB55902CF177.ts
	// Note that for DVR and VoD, the key is key of COS bucket object.
	// Note that for RECORD, the key is the final ts file path.
	Key string `json:"key"`
	// The TS local ID, a uuid string, such as 5B7B5C03-8DB4-4ABA-AAF3-CB55902CF177
	TsID string `json:"tsid"`
	// The TS local file, format is record/:uuid.ts, such as record/5B7B5C03-8DB4-4ABA-AAF3-CB55902CF177.ts
	File string `json:"tsfile"`
	// The TS url, generated by SRS, such as live/livestream/2015-04-23/01/476584165.ts
	URL string `json:"url"`
	// The seqno of TS, generated by SRS, such as 100
	SeqNo uint64 `json:"seqno"`
	// The duration of TS in seconds, generated by SRS, such as 9.36
	Duration float64 `json:"duration"`
	// The size of TS file in bytes, such as 1934897
	Size uint64 `json:"size"`
}

func (v *TsFile) String() string {
	return fmt.Sprintf("key=%v, id=%v, url=%v, seqno=%v, duration=%v, size=%v, file=%v",
		v.Key, v.TsID, v.URL, v.SeqNo, v.Duration, v.Size, v.File,
	)
}

// M3u8VoDArtifact is a HLS VoD object. Because each Dvr/Vod/RecordM3u8Stream might be DVR to many VoD file,
// each is an M3u8VoDArtifact. For example, when user publish live/livestream, there is a Dvr/Vod/RecordM3u8Stream and
// M3u8VoDArtifact, then user unpublish stream and after some seconds a VoD file is generated by M3u8VoDArtifact. Then
// if user republish the stream, there will be a new M3u8VoDArtifact to DVR the stream.
type M3u8VoDArtifact struct {
	// Number of ts files.
	NN int `json:"nn"`
	// The last update time.
	Update string `json:"update"`
	// The uuid of M3u8VoDArtifact, generated by worker, such as 3ECF0239-708C-42E4-96E1-5AE935C6E6A9
	UUID string `json:"uuid"`
	// The url of m3u8, generated by SRS, such as live/livestream/live.m3u8
	M3u8URL string `json:"m3u8_url"`

	// The vhost of stream, generated by SRS, such as video.test.com
	Vhost string `json:"vhost"`
	// The app of stream, generated by SRS, such as live
	App string `json:"app"`
	// The name of stream, generated by SRS, such as livestream
	Stream string `json:"stream"`

	// TODO: FIXME: It's a typo progress.
	// The Record is processing, use local m3u8 address to preview or download.
	Processing bool `json:"progress"`
	// The done time.
	Done string `json:"done"`
	// The ts files of this m3u8.
	Files []*TsFile `json:"files"`

	// For DVR only.
	// The COS bucket name.
	Bucket string `json:"bucket"`
	// The COS bucket region.
	Region string `json:"region"`

	// For VoD only.
	// The file ID generated by VoD commit.
	FileID   string `json:"fileId"`
	MediaURL string `json:"mediaUrl"`
	// The remux task of VoD.
	Definition string `json:"definition"`
	TaskID     string `json:"taskId"`
	// The remux task result.
	Task *VodTaskArtifact `json:"taskObj"`
}

func (v *M3u8VoDArtifact) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("uuid=%v, done=%v, update=%v, processing=%v, files=%v",
		v.UUID, v.Done, v.Update, v.Processing, len(v.Files),
	))
	if v.Bucket != "" {
		sb.WriteString(fmt.Sprintf(", bucket=%v", v.Bucket))
	}
	if v.Region != "" {
		sb.WriteString(fmt.Sprintf(", region=%v", v.Region))
	}
	if v.FileID != "" {
		sb.WriteString(fmt.Sprintf(", fileId=%v", v.FileID))
	}
	if v.MediaURL != "" {
		sb.WriteString(fmt.Sprintf(", mediaUrl=%v", v.MediaURL))
	}
	if v.Definition != "" {
		sb.WriteString(fmt.Sprintf(", definition=%v", v.Definition))
	}
	if v.TaskID != "" {
		sb.WriteString(fmt.Sprintf(", taskId=%v", v.TaskID))
	}
	if v.Task != nil {
		sb.WriteString(fmt.Sprintf(", task=(%v)", v.Task.String()))
	}
	return sb.String()
}

// VodTaskArtifact is the final artifact for remux task.
type VodTaskArtifact struct {
	URL      string  `json:"url"`
	Bitrate  int64   `json:"bitrate"`
	Height   int32   `json:"height"`
	Width    int32   `json:"width"`
	Size     int64   `json:"size"`
	Duration float64 `json:"duration"`
	MD5      string  `json:"md5"`
}

func (v *VodTaskArtifact) String() string {
	return fmt.Sprintf("url=%v", v.URL)
}

// SrsOnHlsMessage is the SRS on_hls callback message.
type SrsOnHlsMessage struct {
	// Must be on_hls
	Action string `json:"action"`
	// The ts file path, such as ./objs/nginx/html/live/livestream/2015-04-23/01/476584165.ts
	File string `json:"file"`
	// The duration of ts file, in seconds, such as 9.36
	Duration float64 `json:"duration"`
	// The url of m3u8, such as live/livestream/live.m3u8
	M3u8URL string `json:"m3u8_url"`
	// The sequence number of ts, such as 100
	SeqNo uint64 `json:"seq_no"`

	// The vhost of stream, generated by SRS, such as video.test.com
	Vhost string `json:"vhost"`
	// The app of stream, generated by SRS, such as live
	App string `json:"app"`
	// The name of stream, generated by SRS, such as livestream
	Stream string `json:"stream"`

	// The TS url, generated by SRS, such as live/livestream/2015-04-23/01/476584165.ts
	URL string `json:"url"`
}

func (v *SrsOnHlsMessage) String() string {
	return fmt.Sprintf("action=%v, file=%v, duration=%v, seqno=%v, m3u8_url=%v, vhost=%v, "+
		"app=%v, stream=%v, url=%v",
		v.Action, v.File, v.Duration, v.SeqNo, v.M3u8URL, v.Vhost, v.App, v.Stream, v.URL,
	)
}

// SrsOnHlsObject contains a SrsOnHlsMessage and a local TsFile.
type SrsOnHlsObject struct {
	Msg    *SrsOnHlsMessage `json:"msg"`
	TsFile *TsFile          `json:"tsfile"`
}

func (v *SrsOnHlsObject) String() string {
	return fmt.Sprintf("msg(%v), ts(%v)", v.Msg.String(), v.TsFile.String())
}

// VodTranscodeTemplate is the transcode template for VoD.
type VodTranscodeTemplate struct {
	// In query template API, it's string. See https://cloud.tencent.com/document/product/266/33769
	// In remux task API, it's integer. See https://cloud.tencent.com/document/product/266/33427
	Definition string `json:"definition"`
	Name       string `json:"name"`
	Comment    string `json:"comment"`
	Container  string `json:"container"`
	Update     string `json:"update"`
}

func (v *VodTranscodeTemplate) String() string {
	return fmt.Sprintf("definition=%v, name=%v, comment=%v, container=%v, update=%v",
		v.Definition, v.Name, v.Comment, v.Container, v.Update,
	)
}

// SrsStream is a stream in SRS.
type SrsStream struct {
	Vhost  string `json:"vhost"`
	App    string `json:"app"`
	Stream string `json:"stream"`
	Param  string `json:"param"`

	Server string `json:"server_id"`
	Client string `json:"client_id"`

	Update string `json:"update"`
}

func (v *SrsStream) String() string {
	return fmt.Sprintf("vhost=%v, app=%v, stream=%v, param=%v, server=%v, client=%v, update=%v",
		v.Vhost, v.App, v.Stream, v.Param, v.Server, v.Client, v.Update,
	)
}

func (v *SrsStream) StreamURL() string {
	streamURL := fmt.Sprintf("%v/%v/%v", v.Vhost, v.App, v.Stream)
	if v.Vhost == "__defaultVhost__" {
		streamURL = fmt.Sprintf("%v/%v", v.App, v.Stream)
	}
	return streamURL
}

func (v *SrsStream) IsSRT() bool {
	return strings.Contains(v.Param, "upstream=srt")
}

func (v *SrsStream) IsRTC() bool {
	return strings.Contains(v.Param, "upstream=rtc")
}
