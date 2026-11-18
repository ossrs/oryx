// Copyright (c) 2022-2026 Winlin
//
// SPDX-License-Identifier: MIT
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
)

type Versions struct {
	Version string `json:"version"`
	Stable  string `json:"stable"`
	Latest  string `json:"latest"`
}

func (v Versions) String() string {
	return fmt.Sprintf("version=%v, latest=%v, stable=%v", v.Version, v.Latest, v.Stable)
}

// TODO: FIXME: Should be merged to mgmt.
type Config struct {
	IsDarwin bool
	// The working directory, at xxx/oryx/platform.
	Pwd string

	Cloud    string
	Region   string
	Source   string
	Registry string

	ipv4  net.IP
	Iface string

	Versions Versions
}

func NewConfig() *Config {
	return &Config{
		ipv4:     net.IPv4zero,
		IsDarwin: runtime.GOOS == "darwin",
		Versions: Versions{
			Version: "v0.0.0",
			Latest:  "v0.0.0",
			Stable:  "v0.0.0",
		},
	}
}

func (v *Config) IPv4() string {
	return v.ipv4.String()
}

func (v *Config) String() string {
	return fmt.Sprintf("darwin=%v, cloud=%v, region=%v, source=%v, registry=%v, iface=%v, ipv4=%v, pwd=%v, "+
		"mgmtPwd=%v, version=%v, latest=%v, stable=%v",
		v.IsDarwin, v.Cloud, v.Region, v.Source, v.Registry, v.Iface, v.IPv4(), v.Pwd, v.Pwd, v.Versions.Version,
		v.Versions.Latest, v.Versions.Stable,
	)
}

func discoverRegion(ctx context.Context) (cloud, region string, err error) {
	if envCloud() == "BT" {
		return "BT", "ap-beijing", nil
	}

	if envCloud() == "BIN" {
		return "BIN", "ap-beijing", nil
	}

	if envCloud() == "AAPANEL" {
		return "AAPANEL", "ap-singapore", nil
	}

	if envCloud() == "DOCKER" {
		return "DOCKER", "ap-beijing", nil
	}

	if envCloud() != "" && envRegion() != "" {
		return envCloud(), envRegion(), nil
	}

	if conf.IsDarwin {
		return "DEV", "ap-beijing", nil
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
	case "DEV", "BT", "BIN", "DOCKER":
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
	case "BIN":
		return "bin", nil
	case "AAPANEL":
		return "aapanel", nil
	case "DOCKER":
		return "docker", nil
	}

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

// TODO: FIXME: Remove it.
const srsDockerName = "srs-server"

// TODO: FIXME: Remove it.
const srsDevDockerName = "srs-dev"

const (
	// Host information, such as the cloud, region and source.
	SRS_TENCENT_LH           = "SRS_TENCENT_LH"
	SRS_HP_HLS               = "SRS_HP_HLS"
	SRS_LL_HLS               = "SRS_LL_HLS"
	SRS_TENCENT_CAM          = "SRS_TENCENT_CAM"
	SRS_TENCENT_COS          = "SRS_TENCENT_COS"
	SRS_TENCENT_VOD          = "SRS_TENCENT_VOD"
	SRS_RECORD_PATTERNS      = "SRS_RECORD_PATTERNS"
	SRS_RECORD_M3U8_WORKING  = "SRS_RECORD_M3U8_WORKING"
	SRS_RECORD_M3U8_ARTIFACT = "SRS_RECORD_M3U8_ARTIFACT"
	SRS_DVR_PATTERNS         = "SRS_DVR_PATTERNS"
	SRS_DVR_M3U8_WORKING     = "SRS_DVR_M3U8_WORKING"
	SRS_DVR_M3U8_ARTIFACT    = "SRS_DVR_M3U8_ARTIFACT"
	SRS_VOD_PATTERNS         = "SRS_VOD_PATTERNS"
	SRS_VOD_M3U8_WORKING     = "SRS_VOD_M3U8_WORKING"
	SRS_VOD_M3U8_ARTIFACT    = "SRS_VOD_M3U8_ARTIFACT"
	// The COS token and target of each VoD upload.
	SRS_VOD_COS_TOKEN      = "SRS_VOD_COS_TOKEN"
	SRS_FORWARD_CONFIG     = "SRS_FORWARD_CONFIG"
	SRS_FORWARD_TASK       = "SRS_FORWARD_TASK"
	SRS_VLIVE_CONFIG       = "SRS_VLIVE_CONFIG"
	SRS_VLIVE_TASK         = "SRS_VLIVE_TASK"
	SRS_CAMERA_CONFIG      = "SRS_CAMERA_CONFIG"
	SRS_CAMERA_TASK        = "SRS_CAMERA_TASK"
	SRS_TRANSCODE_CONFIG   = "SRS_TRANSCODE_CONFIG"
	SRS_TRANSCODE_TASK     = "SRS_TRANSCODE_TASK"
	SRS_TRANSCRIPT_CONFIG  = "SRS_TRANSCRIPT_CONFIG"
	SRS_TRANSCRIPT_TASK    = "SRS_TRANSCRIPT_TASK"
	SRS_OCR_CONFIG         = "SRS_OCR_CONFIG"
	SRS_OCR_TASK           = "SRS_OCR_TASK"
	SRS_STREAM_ACTIVE      = "SRS_STREAM_ACTIVE"
	SRS_STREAM_SRT_ACTIVE  = "SRS_STREAM_SRT_ACTIVE"
	SRS_STREAM_RTC_ACTIVE  = "SRS_STREAM_RTC_ACTIVE"
	SRS_STAT_COUNTER       = "SRS_STAT_COUNTER"
	SRS_CONTAINER_DISABLED = "SRS_CONTAINER_DISABLED"
	SRS_LIVE_ROOM          = "SRS_LIVE_ROOM"
	SRS_DUBBING_PROJECTS   = "SRS_DUBBING_PROJECTS"
	SRS_DUBBING_TASKS      = "SRS_DUBBING_TASKS"
	SRS_AUTH_SECRET        = "SRS_AUTH_SECRET"
	SRS_SECRET_PUBLISH     = "SRS_SECRET_PUBLISH"
	SRS_LOCALE             = "SRS_LOCALE"
	SRS_FIRST_BOOT         = "SRS_FIRST_BOOT"
	SRS_UPGRADING          = "SRS_UPGRADING"
	SRS_UPGRADE_WINDOW     = "SRS_UPGRADE_WINDOW"
	SRS_PLATFORM_SECRET    = "SRS_PLATFORM_SECRET"
	SRS_CACHE_BILIBILI     = "SRS_CACHE_BILIBILI"
	SRS_BEIAN              = "SRS_BEIAN"
	SRS_HTTPS              = "SRS_HTTPS"
	SRS_HTTPS_DOMAIN       = "SRS_HTTPS_DOMAIN"
	SRS_HOOKS              = "SRS_HOOKS"
	SRS_SYS_LIMITS         = "SRS_SYS_LIMITS"
	SRS_SYS_OPENAI         = "SRS_SYS_OPENAI"
)

func GenerateRoomPublishKey(roomStreamName string) string {
	return fmt.Sprintf("room-pub-%v", roomStreamName)
}

// The default, in kbps.
const SrsSysLimitsVLive = 5 * 1000

// The default, in kbps.
const SrsSysLimitsCamera = 5 * 1000

const (
	TENCENT_CLOUD_CAM_ENDPOINT = "cam.tencentcloudapi.com"
	TENCENT_CLOUD_VOD_ENDPOINT = "vod.tencentcloudapi.com"
)

type FFprobeSourceType string

const FFprobeSourceTypeUpload FFprobeSourceType = "upload"
const FFprobeSourceTypeFile FFprobeSourceType = "file"
const FFprobeSourceTypeYTDL FFprobeSourceType = "ytdl"
const FFprobeSourceTypeStream FFprobeSourceType = "stream"

var dirUploadPath = path.Join(".", "upload")
var dirVLivePath = path.Join(".", "vlive")
var dirDubbingPath = path.Join(".", "dub")

// Oryx may use the files in this directory.
const serverDataDirectory = "/data"

var serverAllowVideoFiles []string = []string{".mp4", ".flv", ".ts", ".mkv", ".mov"}

var serverAllowAudioFiles []string = []string{".mp3", ".aac", ".m4a"}

func envApiSecret() string {
	return os.Getenv("SRS_PLATFORM_SECRET")
}

func envNodeEnv() string {
	return os.Getenv("NODE_ENV")
}

func envMgmtPassword() string {
	return os.Getenv("MGMT_PASSWORD")
}

func envSelfSignedCertificate() string {
	return os.Getenv("AUTO_SELF_SIGNED_CERTIFICATE")
}

func envPlatformListen() string {
	return os.Getenv("PLATFORM_LISTEN")
}

func envCloud() string {
	return os.Getenv("CLOUD")
}

func envNameLookup() string {
	return os.Getenv("NAME_LOOKUP")
}

func envPlatformDocker() string {
	return os.Getenv("PLATFORM_DOCKER")
}

func envCandidate() string {
	return os.Getenv("CANDIDATE")
}

func envMgmtListen() string {
	return os.Getenv("MGMT_LISTEN")
}

func envRegion() string {
	return os.Getenv("REGION")
}

func envSource() string {
	return os.Getenv("SOURCE")
}

func envHttpListen() string {
	return os.Getenv("HTTPS_LISTEN")
}

func envSrtListen() string {
	return os.Getenv("SRT_PORT")
}

func envRtcListen() string {
	return os.Getenv("RTC_PORT")
}

func envLocalRelease() string {
	return os.Getenv("LOCAL_RELEASE")
}

func envKeepFiles() string {
	return os.Getenv("AIT_KEEP_FILES")
}

func envRedisPassword() string {
	return os.Getenv("REDIS_PASSWORD")
}

func envRedisPort() string {
	return os.Getenv("REDIS_PORT")
}

func envRedisHost() string {
	return os.Getenv("REDIS_HOST")
}

func envRedisDatabase() string {
	return os.Getenv("REDIS_DATABASE")
}

func envRtmpPort() string {
	return os.Getenv("RTMP_PORT")
}

func envPublicUrl() string {
	return os.Getenv("PUBLIC_URL")
}

func envBuildPath() string {
	return os.Getenv("BUILD_PATH")
}

func envReactAppLocale() string {
	return os.Getenv("REACT_APP_LOCALE")
}

func envHttpPort() string {
	return os.Getenv("HTTP_PORT")
}

func envRegistry() string {
	return os.Getenv("REGISTRY")
}

func envPath() string {
	return os.Getenv("PATH")
}

func envForwardLimit() string {
	return os.Getenv("SRS_FORWARD_LIMIT")
}

func envVLiveLimit() string {
	return os.Getenv("SRS_VLIVE_LIMIT")
}

func envCameraLimit() string {
	return os.Getenv("SRS_CAMERA_LIMIT")
}

func envGoPprof() string {
	return os.Getenv("GO_PPROF")
}

func envYtdlProxy() string {
	return os.Getenv("YTDL_PROXY")
}

var rdb *redis.Client

func InitRdb() error {
	redisDatabase, err := strconv.Atoi(envRedisDatabase())
	if err != nil {
		return errors.Wrapf(err, "invalid REDIS_DATABASE %v", envRedisDatabase())
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", envRedisHost(), envRedisPort()),
		Password: envRedisPassword(),
		DB:       redisDatabase,
	})
	return nil
}

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

func refreshIPv4(ctx context.Context) error {
	discoverPrivateIPv4 := func(ctx context.Context) (string, net.IP, error) {
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

			// The IP address might change.
			time.Sleep(time.Duration(60) * time.Second)
		}
	}()

	select {
	case <-ctx.Done():
	case <-ipv4Ctx.Done():
	}

	return nil
}

func setEnvDefault(key, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}

func srsGenerateConfig(ctx context.Context) error {
	hlsConf := []string{
		"",
		"hls {",
		"    enabled on;",
	}
	if hlsLowLatency, err := rdb.HGet(ctx, SRS_LL_HLS, "hlsLowLatency").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v hls", SRS_LL_HLS)
	} else {
		if hlsLowLatency != "true" {
			hlsConf = append(hlsConf, []string{
				"    hls_fragment 10;",
				"    hls_window 60;",
			}...)
		} else {
			hlsConf = append(hlsConf, []string{
				"    hls_fragment 2;",
				"    hls_window 16;",
			}...)
		}
	}
	hlsConf = append(hlsConf, []string{
		"    hls_aof_ratio 2.1;",
		"    hls_path ./objs/nginx/html;",
		"    hls_m3u8_file [app]/[stream].m3u8;",
		"    hls_ts_file [app]/[stream]-[seq]-[timestamp].ts;",
		"    hls_wait_keyframe on;",
		"    hls_dispose 15;",
	}...)
	if noHlsCtx, err := rdb.HGet(ctx, SRS_HP_HLS, "noHlsCtx").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v hls", SRS_HP_HLS)
	} else if noHlsCtx == "true" {
		hlsConf = append(hlsConf, []string{
			"    hls_ctx off;",
			"    hls_ts_ctx off;",
		}...)
	}
	hlsConf = append(hlsConf, "}")

	if true {
		confLines := []string{
			"# !!! Important: This file is produced and maintained by the Oryx, please never modify it.",
		}
		confLines = append(confLines, "", "")

		confData := strings.Join(confLines, "\n")
		fileName := path.Join(conf.Pwd, "containers/data/config/srs.server.conf")
		if f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
			return errors.Wrapf(err, "open file %v", fileName)
		} else {
			defer f.Close()
			if _, err = f.Write([]byte(confData)); err != nil {
				return errors.Wrapf(err, "write file %v with %v", fileName, confData)
			}
		}
	}
	if true {
		confLines := []string{
			"# !!! Important: This file is produced and maintained by the Oryx, please never modify it.",
		}
		confLines = append(confLines, hlsConf...)
		confLines = append(confLines, "", "")

		confData := strings.Join(confLines, "\n")
		fileName := path.Join(conf.Pwd, "containers/data/config/srs.vhost.conf")
		if f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
			return errors.Wrapf(err, "open file %v", fileName)
		} else {
			defer f.Close()
			if _, err = f.Write([]byte(confData)); err != nil {
				return errors.Wrapf(err, "write file %v with %v", fileName, confData)
			}
		}
	}

	// The ID of the SRS reload transaction.
	fetchReload := func(ctx context.Context) (string, error) {
		// TODO: FIXME: Remove it after SRS merged https://github.com/ossrs/srs/pull/3768
		return time.Now().String(), nil

		/*api := "http://127.0.0.1:1985/api/v1/raw?rpc=reload-fetch"
		req, err := http.NewRequestWithContext(ctx, "GET", api, nil)
		if err != nil {
			return "", errors.Wrapf(err, "reload fetch srs %v", api)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", errors.Wrapf(err, "reload srs %v", api)
		}
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", errors.Wrapf(err, "read srs %v", api)
		}

		if res.StatusCode != http.StatusOK {
			return "", errors.Errorf("reload srs %v, code=%v, body=%v", api, res.StatusCode, string(b))
		}

		resObj := struct {
			Code int `json:"code"`
			Data struct {
				Err   int    `json:"err"`
				Msg   string `json:"msg"`
				State int    `json:"state"`
				RID   string `json:"rid"`
			} `json:"data"`
		}{}
		if err := json.Unmarshal(b, &resObj); err != nil {
			return "", errors.Wrapf(err, "unmarshal reload id %v", string(b))
		}

		if resObj.Code != 0 || resObj.Data.Err != 0 {
			return "", errors.Errorf("reload srs code=%v, err=%v, invalid %v",
				resObj.Code, resObj.Data.Err, string(b),
			)
		}
		logger.Tf(ctx, "reload fetch srs ok")

		return resObj.Data.RID, nil*/
	}
	reloadID, err := fetchReload(ctx)
	if err != nil {
		return errors.Wrapf(err, "fetch reload id")
	}

	if true {
		api := "http://127.0.0.1:1985/api/v1/raw?rpc=reload"
		res, err := http.DefaultClient.Get(api)
		if err != nil {
			return errors.Wrapf(err, "reload srs %v", api)
		}
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrapf(err, "read srs %v", api)
		}

		if res.StatusCode != http.StatusOK {
			return errors.Errorf("reload srs %v, code=%v, body=%v", api, res.StatusCode, string(b))
		}
		logger.Tf(ctx, "reload submit srs ok")
	}

	for i := 0; i < 10; i++ {
		if newReloadID, err := fetchReload(ctx); err != nil {
			return errors.Wrapf(err, "fetch reload id")
		} else if newReloadID != reloadID {
			logger.Tf(ctx, "reload fetch srs ok")
			return nil
		}

		select {
		case <-ctx.Done():
		case <-time.After(time.Second):
		}
	}

	return errors.New("reload srs timeout")
}

func nginxGenerateConfig(ctx context.Context) error {
	sslConf := []string{}
	if ssl, err := rdb.Get(ctx, SRS_HTTPS).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "get %v", SRS_HTTPS)
	} else if ssl == "ssl" || ssl == "lets" {
		sslConf = []string{
			"",
			"# For SSL/TLS config.",
			"listen       443 ssl;",
			"listen       [::]:443 ssl;",
			"ssl_certificate /data/config/nginx.crt;",
			"ssl_certificate_key /data/config/nginx.key;",
			"ssl_protocols TLSv1.1 TLSv1.2 TLSv1.3;",
			`add_header Strict-Transport-Security "max-age=0";`,
			"ssl_session_cache shared:SSL:10m;",
			"ssl_session_timeout 10m;",
			"",
		}
	}

	uploadLimit := []string{
		"",
		"# Limit for upload file size",
		"client_max_body_size 100g;",
	}

	if true {
		confLines := []string{
			"# !!! Important: This file is produced and maintained by the Oryx, please never modify it.",
		}
		confLines = append(confLines, "", "")

		confData := strings.Join(confLines, "\n")
		fileName := path.Join(conf.Pwd, "containers/data/config/nginx.http.conf")
		if f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
			return errors.Wrapf(err, "open file %v", fileName)
		} else {
			defer f.Close()
			if _, err = f.Write([]byte(confData)); err != nil {
				return errors.Wrapf(err, "write file %v with %v", fileName, confData)
			}
		}
	}
	if true {
		confLines := []string{
			"# !!! Important: This file is produced and maintained by the Oryx, please never modify it.",
		}
		confLines = append(confLines, uploadLimit...)
		confLines = append(confLines, sslConf...)
		confLines = append(confLines, "", "")

		confData := strings.Join(confLines, "\n")
		fileName := path.Join(conf.Pwd, "containers/data/config/nginx.server.conf")
		if f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
			return errors.Wrapf(err, "open file %v", fileName)
		} else {
			defer f.Close()
			if _, err = f.Write([]byte(confData)); err != nil {
				return errors.Wrapf(err, "write file %v with %v", fileName, confData)
			}
		}
	}

	reloadNginx := func(ctx context.Context) error {
		defer certManager.ReloadCertificate(ctx)

		if conf.IsDarwin {
			logger.T(ctx, "ignore reload nginx on darwin")
			return nil
		}

		fileName := path.Join(conf.Pwd, fmt.Sprintf("containers/data/signals/nginx.reload.%v",
			time.Now().UnixNano()/int64(time.Millisecond),
		))
		if f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
			return errors.Wrapf(err, "open file %v", fileName)
		} else {
			defer f.Close()
			msg := fmt.Sprintf("Oryx reload Nginx at %v\n", time.Now().Format(time.RFC3339))
			if _, err = f.Write([]byte(msg)); err != nil {
				return errors.Wrapf(err, "write file %v", fileName)
			}
		}
		return nil
	}

	if err := reloadNginx(ctx); err != nil {
		return errors.Wrapf(err, "reload nginx")
	}
	logger.Tf(ctx, "NGINX: Refresh nginx conf ok")

	return nil
}

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

func buildVodM3u8ForLocal(
	ctx context.Context, tsFiles []*TsFile, useKey bool, prefix string,
) (
	contentType, m3u8Body string, duration float64, err error,
) {
	if len(tsFiles) == 0 {
		err = errors.Errorf("no files")
		return
	}

	for _, file := range tsFiles {
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
	for index, file := range tsFiles {
		// TODO: FIXME: Identify discontinuity by callback.
		if index < len(tsFiles)-2 {
			next := tsFiles[index+1]
			if file.SeqNo+1 != next.SeqNo {
				m3u8 = append(m3u8, "#EXT-X-DISCONTINUITY")
			}
		}

		m3u8 = append(m3u8, fmt.Sprintf("#EXTINF:%.2f, no desc", file.Duration))

		var tsURL string
		if useKey {
			tsURL = fmt.Sprintf("%v%v", prefix, file.Key)
		} else {
			tsURL = fmt.Sprintf("%v%v.ts", prefix, file.TsID)
		}
		m3u8 = append(m3u8, tsURL)
	}
	m3u8 = append(m3u8, "#EXT-X-ENDLIST")

	contentType = "application/vnd.apple.mpegurl"
	m3u8Body = strings.Join(m3u8, "\n")
	return
}

func buildLiveM3u8ForLocal(
	ctx context.Context, tsFiles []*TsFile, useKey bool, prefix string,
) (
	contentType, m3u8Body string, duration float64, err error,
) {
	if len(tsFiles) == 0 {
		err = errors.Errorf("no files")
		return
	}

	for _, file := range tsFiles {
		duration = math.Max(duration, file.Duration)
	}

	first := tsFiles[0]

	m3u8 := []string{
		"#EXTM3U",
		"#EXT-X-VERSION:3",
		fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%v", first.SeqNo),
		fmt.Sprintf("#EXT-X-TARGETDURATION:%v", math.Ceil(duration)),
	}
	for index, file := range tsFiles {
		// TODO: FIXME: Identify discontinuity by callback.
		if index < len(tsFiles)-2 {
			next := tsFiles[index+1]
			if file.SeqNo+1 != next.SeqNo {
				m3u8 = append(m3u8, "#EXT-X-DISCONTINUITY")
			}
		}

		m3u8 = append(m3u8, fmt.Sprintf("#EXTINF:%.2f, no desc", file.Duration))

		var tsURL string
		if useKey {
			tsURL = fmt.Sprintf("%v%v", prefix, file.Key)
		} else {
			tsURL = fmt.Sprintf("%v%v.ts", prefix, file.TsID)
		}
		m3u8 = append(m3u8, tsURL)
	}

	contentType = "application/vnd.apple.mpegurl"
	m3u8Body = strings.Join(m3u8, "\n")
	return
}

func buildLiveM3u8ForVariantCC(
	ctx context.Context, bitrate int64, lang, stream, subtitles string,
) (contentType, m3u8Body string, err error) {
	m3u8 := []string{
		"#EXTM3U",
		fmt.Sprintf(
			`#EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID="subs",NAME="Subtitle-%v",LANGUAGE="%v",DEFAULT=YES,AUTOSELECT=YES,FORCED=NO,URI="%v"`,
			strings.ToUpper(lang), lang, subtitles,
		),
		fmt.Sprintf(`#EXT-X-STREAM-INF:BANDWIDTH=%v,SUBTITLES="subs"`, bitrate),
		stream,
	}

	contentType = "application/vnd.apple.mpegurl"
	m3u8Body = strings.Join(m3u8, "\n")
	return
}

func slicesContains(arr []string, elem string) bool {
	for _, e := range arr {
		if e == elem {
			return true
		}
	}
	return false
}

type TsFile struct {
	// For the record, the final ts file path, such as
	// record/3ECF0239-708C-42E4-96E1-5AE935C6E6A9/5B7B5C03-8DB4-4ABA-AAF3-CB55902CF177.ts.
	// For DVR and VoD, the COS object key, such as
	// 3ECF0239-708C-42E4-96E1-5AE935C6E6A9/5B7B5C03-8DB4-4ABA-AAF3-CB55902CF177.ts.
	// Not used by the transcript.
	Key string `json:"key,omitempty"`
	// A UUID, such as 5B7B5C03-8DB4-4ABA-AAF3-CB55902CF177.
	TsID string `json:"tsid,omitempty"`
	// The local file, such as record/5B7B5C03-8DB4-4ABA-AAF3-CB55902CF177.ts.
	File string `json:"tsfile,omitempty"`
	// Generated by SRS, such as live/livestream/2015-04-23/01/476584165.ts.
	URL string `json:"url,omitempty"`
	// Generated by SRS.
	SeqNo uint64 `json:"seqno,omitempty"`
	// In seconds, generated by SRS.
	Duration float64 `json:"duration,omitempty"`
	// In bytes.
	Size uint64 `json:"size,omitempty"`
}

func (v *TsFile) String() string {
	return fmt.Sprintf("key=%v, id=%v, url=%v, seqno=%v, duration=%v, size=%v, file=%v",
		v.Key, v.TsID, v.URL, v.SeqNo, v.Duration, v.Size, v.File,
	)
}

// A Dvr, Vod or RecordM3u8Stream may make many VoD files, one artifact each. The VoD file is made
// some seconds after the stream is unpublished, and a republish makes a new artifact.
type M3u8VoDArtifact struct {
	// The number of ts files.
	NN     int    `json:"nn"`
	Update string `json:"update"`
	// Generated by the worker.
	UUID string `json:"uuid"`
	// Generated by SRS, such as live/livestream/live.m3u8.
	M3u8URL string `json:"m3u8_url"`

	Vhost  string `json:"vhost"`
	App    string `json:"app"`
	Stream string `json:"stream"`

	// TODO: FIXME: It's a typo progress.
	// While processing, the local m3u8 is used to preview or download.
	Processing bool      `json:"progress"`
	Done       string    `json:"done"`
	Files      []*TsFile `json:"files"`

	// For DVR only.
	Bucket string `json:"bucket"`
	Region string `json:"region"`

	// For VoD only.
	// From the VoD commit.
	FileID   string `json:"fileId"`
	MediaURL string `json:"mediaUrl"`
	// The remux template and task.
	Definition uint64           `json:"definition"`
	TaskID     string           `json:"taskId"`
	Task       *VodTaskArtifact `json:"taskObj"`
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
	if v.Definition != 0 {
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

// The output file of the VoD remux task.
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

type SrsOnHlsMessage struct {
	// Must be on_hls.
	Action SrsAction `json:"action,omitempty"`
	// Such as ./objs/nginx/html/live/livestream/2015-04-23/01/476584165.ts.
	File string `json:"file,omitempty"`
	// In seconds.
	Duration float64 `json:"duration,omitempty"`
	// Such as live/livestream/live.m3u8.
	M3u8URL string `json:"m3u8_url,omitempty"`
	SeqNo   uint64 `json:"seq_no,omitempty"`

	Vhost  string `json:"vhost,omitempty"`
	App    string `json:"app,omitempty"`
	Stream string `json:"stream,omitempty"`

	// Such as live/livestream/2015-04-23/01/476584165.ts.
	URL string `json:"url,omitempty"`
}

func (v *SrsOnHlsMessage) String() string {
	return fmt.Sprintf("action=%v, file=%v, duration=%v, seqno=%v, m3u8_url=%v, vhost=%v, "+
		"app=%v, stream=%v, url=%v",
		v.Action, v.File, v.Duration, v.SeqNo, v.M3u8URL, v.Vhost, v.App, v.Stream, v.URL,
	)
}

type SrsOnHlsObject struct {
	Msg    *SrsOnHlsMessage `json:"msg"`
	TsFile *TsFile          `json:"tsfile"`
}

func (v *SrsOnHlsObject) String() string {
	return fmt.Sprintf("msg(%v), ts(%v)", v.Msg.String(), v.TsFile.String())
}

type VodTranscodeTemplate struct {
	// A string in the query template API, see https://cloud.tencent.com/document/product/266/33769
	// An integer in the remux task API, see https://cloud.tencent.com/document/product/266/33427
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

type SrsStream struct {
	Vhost  string `json:"vhost,omitempty"`
	App    string `json:"app,omitempty"`
	Stream string `json:"stream,omitempty"`
	Param  string `json:"param,omitempty"`

	Server string `json:"server_id,omitempty"`
	Client string `json:"client_id,omitempty"`

	Update string `json:"update,omitempty"`
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

func ParseBody(ctx context.Context, r io.ReadCloser, v interface{}) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrapf(err, "read body")
	}
	defer r.Close()

	if len(b) == 0 {
		return nil
	}

	if err := json.Unmarshal(b, v); err != nil {
		return errors.Wrapf(err, "json unmarshal %v", string(b))
	}

	return nil
}

func Authenticate(ctx context.Context, apiSecret, token string, header http.Header) error {
	if apiSecret == "" {
		return errors.New("no api secret")
	}

	authorization := header.Get("Authorization")
	if authorization == "" && token == "" {
		return errors.New("no Authorization or token")
	}

	if authorization != "" {
		parseBearerToken := func(authorization string) (string, error) {
			authParts := strings.Split(authorization, " ")
			if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
				return "", errors.New("Invalid Authorization format")
			}

			return authParts[1], nil
		}

		authSecret, err := parseBearerToken(authorization)
		if err != nil {
			return errors.Wrapf(err, "parse bearer token")
		}

		if authSecret != apiSecret {
			return errors.New("invalid bearer token")
		}
		return nil
	}

	// See https://www.npmjs.com/package/jsonwebtoken#errors--codes
	// See https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
	if _, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(apiSecret), nil
	}); err != nil {
		return errors.Wrapf(err, "verify token %v", token)
	}

	return nil
}

func ChooseNotEmpty(strings ...string) string {
	for _, str := range strings {
		if str != "" {
			return str
		}
	}
	return ""
}

func RebuildStreamURL(rawURL string) (*url.URL, error) {
	// Special characters in the username or password may fail the parse, or leave no host.
	if r0, err := url.Parse(rawURL); err == nil {
		isSRT := r0.Scheme == "srt"
		hostIsInvalid := r0.Host == ""
		pathIsInvalid := r0.Path == "" && r0.RequestURI() != "/"
		if isSRT || (!hostIsInvalid && !pathIsInvalid) {
			return r0, err
		}
	}

	// Only compatible with schema://username:password@hostport/path?query#fragment
	schema, u, _ := strings.Cut(rawURL, "://")
	u2, pathQueryFragment, _ := strings.Cut(u, "/")

	hostport, u3 := u2, u2
	if index := strings.LastIndex(u2, "@"); index >= 0 {
		hostport = u2[index+1:]
		u3 = u2[:index]
	}

	username, password, _ := strings.Cut(u3, ":")
	userInfo := url.UserPassword(username, password)

	var sb strings.Builder
	sb.WriteString(schema)
	sb.WriteString("://")
	if ui := userInfo.String(); ui != "" {
		sb.WriteString(ui)
		sb.WriteString("@")
	}
	sb.WriteString(hostport)
	if pathQueryFragment != "" {
		sb.WriteString("/")
		sb.WriteString(pathQueryFragment)
	}

	if r0, err := url.Parse(sb.String()); err == nil {
		return r0, err
	}

	return url.Parse(rawURL)
}

// CORS is always enabled, because HTTP cache is enabled.
func httpAllowCORS(w http.ResponseWriter, r *http.Request) {
	// SRS needs no cookie or credentials, so * is used for the origin, headers, expose headers and methods.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers
	w.Header().Set("Access-Control-Allow-Headers", "*")
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods
	w.Header().Set("Access-Control-Allow-Methods", "*")
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers
	w.Header().Set("Access-Control-Expose-Headers", "*")
	// https://stackoverflow.com/a/24689738/17679565
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func httpCreateProxy(targetURL string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, errors.Wrapf(err, "parse backend %v", targetURL)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ModifyResponse = func(resp *http.Response) error {
		// The platform sets its own Server header.
		resp.Header.Del("Server")

		// The platform sets its own CORS headers.
		resp.Header.Del("Access-Control-Allow-Origin")
		resp.Header.Del("Access-Control-Allow-Headers")
		resp.Header.Del("Access-Control-Allow-Methods")
		resp.Header.Del("Access-Control-Expose-Headers")
		resp.Header.Del("Access-Control-Allow-Credentials")

		// Not used right now.
		resp.Header.Del("Access-Control-Request-Private-Network")

		return nil
	}

	return proxy, nil
}

type whxpResponseModifier struct {
	w http.ResponseWriter
}

func (w *whxpResponseModifier) Header() http.Header {
	return w.w.Header()
}

func (w *whxpResponseModifier) Write(b []byte) (int, error) {
	// TODO: FIXME: Should pass the rtc port to WHIP/WHEP api, because the port maybe not the same length to 8000,
	//  for example, 80, 443, 18000, etc, in such case, the sdp length will change.
	if port := envRtcListen(); port != "8000" {
		scan := bufio.NewScanner(strings.NewReader(string(b)))

		var lines []string
		for scan.Scan() {
			line := scan.Text()
			if strings.Contains(line, "candidate") {
				line = strings.ReplaceAll(line, " 8000 ", fmt.Sprintf(" %v ", port))
			}
			lines = append(lines, line)
		}

		sdp := strings.Join(lines, "\r\n") + "\r\n"

		return w.w.Write([]byte(sdp))
	}
	return w.w.Write(b)
}

func (w *whxpResponseModifier) WriteHeader(statusCode int) {
	w.w.WriteHeader(statusCode)
}

type FFprobeFormat struct {
	// In seconds.
	Starttime string `json:"start_time"`
	// In seconds.
	Duration string `json:"duration"`
	// In bps.
	Bitrate string `json:"bit_rate"`
	// Data streams count too, so 2 streams may be audio and video, video and data, or audio and data.
	Streams int32 `json:"nb_streams"`
	// The confidence of the format detection.
	Score    int32 `json:"probe_score"`
	HasVideo bool  `json:"has_video"`
	HasAudio bool  `json:"has_audio"`
}

func (v *FFprobeFormat) String() string {
	return fmt.Sprintf("starttime=%v, duration=%v, bitrate=%v, streams=%v, score=%v, video=%v, audio=%v",
		v.Starttime, v.Duration, v.Bitrate, v.Streams, v.Score, v.HasVideo, v.HasAudio,
	)
}

type FFprobeVideo struct {
	CodecType string `json:"codec_type"`
	// For example, h264, h265, vp6f, vp8, vp9, av1, or avs3.
	CodecName string `json:"codec_name"`
	// For example, High, Main, Baseline, or Constrained Baseline.
	Profile string `json:"profile"`
	Width   int32  `json:"width"`
	Height  int32  `json:"height"`
	// For example, yuv420p, yuv422p, yuv444p, yuv410p, yuv411p, or yuvj420p.
	PixFormat string `json:"pix_fmt"`
	Level     int32  `json:"level"`
	// In bps.
	Bitrate string `json:"bit_rate"`
	// In seconds.
	Starttime string `json:"start_time"`
	// In seconds.
	Duration string `json:"duration"`
}

func (v *FFprobeVideo) String() string {
	return fmt.Sprintf("codec=%v, profile=%v, width=%v, height=%v, fmt=%v, level=%v, bitrate=%v",
		v.CodecName, v.Profile, v.Width, v.Height, v.PixFormat, v.Level, v.Bitrate,
	)
}

type FFprobeAudio struct {
	CodecType string `json:"codec_type"`
	// For example, aac, mp3, opus, vorbis, or flac.
	CodecName string `json:"codec_name"`
	// For example, AAC LC, AAC HE, AAC HEv2, or AAC LD.
	Profile string `json:"profile"`
	// For example, fltp, s16p, s32p, s64p, or dbl.
	SampleFormat string `json:"sample_fmt"`
	// In Hz.
	SampleRate string `json:"sample_rate"`
	Channels   int32  `json:"channels"`
	// For example, mono, stereo, 5.1, or 7.1.
	ChannelLayout string `json:"channel_layout"`
	// In bps.
	Bitrate string `json:"bit_rate"`
	// In seconds.
	Starttime string `json:"start_time"`
	// In seconds.
	Duration string `json:"duration"`
}

func (v *FFprobeAudio) String() string {
	return fmt.Sprintf("codec=%v, profile=%v, fmt=%v, rate=%v, channels=%v, layout=%v, bitrate=%v",
		v.CodecName, v.Profile, v.SampleFormat, v.SampleRate, v.Channels, v.ChannelLayout, v.Bitrate,
	)
}

type FFprobeSource struct {
	Name string `json:"name"`
	Path string `json:"path"`
	// In bytes.
	Size   uint64            `json:"size"`
	UUID   string            `json:"uuid"`
	Target string            `json:"target"`
	Type   FFprobeSourceType `json:"type"`
	Format *FFprobeFormat    `json:"format"`
	Video  *FFprobeVideo     `json:"video"`
	Audio  *FFprobeAudio     `json:"audio"`
}

func (v *FFprobeSource) String() string {
	return fmt.Sprintf("name=%v, path=%v, size=%v, uuid=%v, target=%v, format=(%v), video=(%v), audio=(%v)",
		v.Name, v.Path, v.Size, v.UUID, v.Target, v.Format, v.Video, v.Audio,
	)
}

type FFmpegHeartbeat struct {
	starttime time.Time
	update    time.Time

	firstReadyTime time.Time
	firstReadyCtx  context.Context

	// Set by the FFmpeg log, such as:
	//		Exiting normally, received signal 2.
	exitingNormally bool
	parsedCount     uint64
	// FFmpeg prints a cycle log every second. Other logs, such as errors, are kept here.
	extraLogs []string
	// The last cycle log whose time changed, with its time and speed.
	line, timestamp, speed string
	failedParsedCount      uint64
	lastFailedParsed       string
	notChangedCount        uint64
	lastNotChanged         string
	failedSpeedCount       uint64
	lastFailedSpeed        string
	veryFastSpeedCount     uint64
	verySlowSpeedCount     uint64

	FrameLogs chan string
	// Canceled when reading the FFmpeg stderr ends.
	PollingCtx   context.Context
	cancelFFmpeg context.CancelFunc

	// Exit after streaming for this duration.
	MaxStreamDuration time.Duration
	AbnormalFastSpeed float64
}

// On a timeout, cancelFFmpeg must stop FFmpeg, or the caller waits for it forever.
func NewFFmpegHeartbeat(cancelFFmpeg context.CancelFunc, opts ...func(*FFmpegHeartbeat)) *FFmpegHeartbeat {
	v := &FFmpegHeartbeat{
		starttime:    time.Now(),
		update:       time.Now(),
		FrameLogs:    make(chan string, 1),
		cancelFFmpeg: cancelFFmpeg,
		// No limit.
		MaxStreamDuration: time.Duration(0),
		AbnormalFastSpeed: FFmpegAbnormalFastSpeed,
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// Works only for a URL input, not a file.
func (v *FFmpegHeartbeat) Parse(u *url.URL) {
	q := u.Query()
	if qv := q.Get("max-stream-duration"); qv != "" {
		v.MaxStreamDuration, _ = time.ParseDuration(qv)
	}
	if qv := q.Get("abnormal-fast-speed"); qv != "" {
		v.AbnormalFastSpeed, _ = strconv.ParseFloat(qv, 64)
	}
}

func (v *FFmpegHeartbeat) Polling(ctx context.Context, stderr io.Reader) {
	pollingReadyCtx, cancelPollingReady := context.WithCancel(ctx)
	logger.Tf(ctx, "FFmpeg: Start to polling heartbeat, start=%v, msd=%v, afs=%v",
		v.starttime, v.MaxStreamDuration, v.AbnormalFastSpeed)

	go func() {
		<-pollingReadyCtx.Done()

		select {
		case <-ctx.Done():
		case <-v.PollingCtx.Done():
		}
		logger.Tf(ctx, "FFmpeg: Quit exit-normally=%v, parsed=%v, failed=%v,<%v>, speed=%v,%v,%v,<%v>, not-change=%v,<%v>, extra logs is %v",
			v.exitingNormally, v.parsedCount, v.failedParsedCount, v.lastFailedParsed, v.failedSpeedCount,
			v.veryFastSpeedCount, v.verySlowSpeedCount, v.lastFailedSpeed, v.notChangedCount, v.lastNotChanged,
			strings.Join(v.extraLogs, " "))
	}()

	go func() {
		<-pollingReadyCtx.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-v.PollingCtx.Done():
				return
			case <-time.After(3 * time.Second):
			}

			if v.update.Add(10 * time.Second).Before(time.Now()) {
				logger.Wf(ctx, "FFmpeg: not update for %v, restart it", time.Since(v.update))
				v.cancelFFmpeg()
				return
			}
		}
	}()

	var firstNormalFrame bool
	var firstReadyCancel context.CancelFunc
	v.firstReadyCtx, firstReadyCancel = context.WithCancel(context.Background())
	handleOutputOfFFmpeg := func(ffmpegLog string) {
		line := strings.TrimSpace(ffmpegLog)
		for strings.Contains(line, "= ") {
			line = strings.ReplaceAll(line, "= ", "=")
		}
		line = strings.ReplaceAll(line, "\n", " ")
		line = strings.ReplaceAll(line, "\r", " ")

		if strings.Contains(line, "Exiting normally") {
			v.exitingNormally = true
		}

		if !strings.Contains(line, "size=") && !strings.Contains(line, "time=") {
			v.extraLogs = append(v.extraLogs, line)
			return
		}
		if strings.Contains(line, "time=N/A") || strings.Contains(line, "speed=N/A") {
			v.extraLogs = append(v.extraLogs, line)
			return
		}

		timestamp, speed, err := ParseFFmpegCycleLog(line)
		if err != nil {
			v.failedParsedCount, v.lastFailedParsed = v.failedParsedCount+1, line
			return
		}

		// FFmpeg is alive only while the time changes.
		if v.timestamp == timestamp {
			v.notChangedCount, v.lastNotChanged = v.notChangedCount+1, line
			return
		}

		// During live streaming, the speed is typically about 1x. It may exceed 1x, mostly during
		// startup, but the mean stays at 1x.
		if speedv, err := strconv.ParseFloat(strings.Trim(speed, "x"), 64); err != nil || speedv <= 0 {
			v.failedSpeedCount, v.lastFailedSpeed = v.failedSpeedCount+1, line
		} else {
			if speedv > v.AbnormalFastSpeed {
				v.veryFastSpeedCount++
			} else {
				v.veryFastSpeedCount = 0
			}

			if speedv < FFmpegAbnormalSlowSpeed {
				v.verySlowSpeedCount++
			} else {
				v.verySlowSpeedCount = 0
			}

			var exitForTimeout bool
			if v.MaxStreamDuration > 0 && time.Now().After(v.starttime.Add(v.MaxStreamDuration)) {
				logger.Tf(ctx, "FFmpeg: exit after duration=%v, start=%v, now=%v",
					v.MaxStreamDuration, v.starttime, time.Now())
				exitForTimeout = true
			}

			if mv := RestartFFmpegCountAbnormalSpeed; v.veryFastSpeedCount > mv || v.verySlowSpeedCount > mv || exitForTimeout {
				logger.Wf(ctx, "FFmpeg: abnormal speed=%v, fast=%v, slow=%v, mv=%v, timeout=%v,%v, restart it",
					speed, v.veryFastSpeedCount, v.verySlowSpeedCount, mv, exitForTimeout, v.MaxStreamDuration)
				v.cancelFFmpeg()
				return
			}
		}

		v.update, v.parsedCount = time.Now(), v.parsedCount+1
		v.line, v.timestamp, v.speed = line, timestamp, speed
		if !firstNormalFrame {
			firstNormalFrame, v.firstReadyTime = true, time.Now()
			firstReadyCancel()

			logger.Tf(ctx, "FFmpeg: First normal frame, parsed=%v, failed=%v,<%v>, speed=%v,%v,%v,<%v>, not-change=%v,<%v>, last=<%v>",
				v.parsedCount, v.failedParsedCount, v.lastFailedParsed, v.failedSpeedCount,
				v.veryFastSpeedCount, v.verySlowSpeedCount, v.lastFailedSpeed, v.notChangedCount, v.lastNotChanged,
				v.line)
		}

		select {
		case <-ctx.Done():
		case v.FrameLogs <- line:
		}
	}

	go func() {
		defer v.cancelFFmpeg()

		var pollingCancel context.CancelFunc
		v.PollingCtx, pollingCancel = context.WithCancel(ctx)
		defer pollingCancel()

		// The other goroutines use v.PollingCtx, so they start only now.
		cancelPollingReady()

		buf := make([]byte, 4096)
		for ctx.Err() == nil {
			nn, err := stderr.Read(buf)
			if nn > 0 {
				handleOutputOfFFmpeg(string(buf[:nn]))
			}

			if err != nil && !v.exitingNormally {
				logger.Wf(ctx, "FFmpeg: read stderr failed: %v", err)
			}
			if nn == 0 && !v.exitingNormally {
				logger.Wf(ctx, "FFmpeg: EOF of stderr")
			}

			if err != nil || nn == 0 {
				return
			}
		}
	}()

	go func() {
		<-pollingReadyCtx.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-v.PollingCtx.Done():
				return
			case <-time.After(15 * time.Second):
				logger.Tf(ctx, "FFmpeg: Running parsed=%v, failed=%v,<%v>, speed=%v,%v,%v,<%v>, not-change=%v,<%v>, last=<%v>",
					v.parsedCount, v.failedParsedCount, v.lastFailedParsed, v.failedSpeedCount,
					v.veryFastSpeedCount, v.verySlowSpeedCount, v.lastFailedSpeed, v.notChangedCount, v.lastNotChanged,
					v.line)
			}
		}
	}()

	// The caller uses v.PollingCtx after Polling returns.
	<-pollingReadyCtx.Done()
}

const FFmpegAbnormalFastSpeed = 1.5

const FFmpegAbnormalSlowSpeed = 0.5

const RestartFFmpegCountAbnormalSpeed = uint64(30)

// The log is mostly like:
//
//	size=18859kB time=00:10:09.38 bitrate=253.5kbits/s speed=1x
//	frame=184 fps=9.7 q=28.0 size=364kB time=00:00:19.41 bitrate=153.7kbits/s dup=0 drop=235 speed=1.03x
func ParseFFmpegCycleLog(line string) (timestamp, speed string, err error) {
	re := regexp.MustCompile(`time=(\S+)( .*)speed=(\S+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 4 {
		err = errors.Errorf("parse %v failed, matches=%v", line, matches)
		return
	}
	timestamp, speed = matches[1], matches[3]
	return
}

type MediaFormat struct {
	Starttime string  `json:"start_time"`
	Duration  float64 `json:"duration"`
	Bitrate   int64   `json:"bit_rate"`
	Streams   int32   `json:"nb_streams"`
	Score     int32   `json:"probe_score"`
	HasVideo  bool    `json:"has_video"`
	HasAudio  bool    `json:"has_audio"`
}

func (v *MediaFormat) FromFFprobeFormat(format *FFprobeFormat) error {
	v.Starttime = format.Starttime
	v.Streams = format.Streams
	v.Score = format.Score
	v.HasVideo = format.HasVideo
	v.HasAudio = format.HasAudio

	if fv, err := strconv.ParseFloat(format.Duration, 64); err != nil {
		return errors.Wrapf(err, "parse duration %v of %v", format.Duration, format)
	} else {
		v.Duration = fv
	}

	if iv, err := strconv.ParseInt(format.Bitrate, 10, 64); err != nil {
		return errors.Wrapf(err, "parse bitrate %v of %v", format.Bitrate, format)
	} else {
		v.Bitrate = iv
	}

	return nil
}

func (v *MediaFormat) String() string {
	return fmt.Sprintf("starttime=%v, duration=%v, bitrate=%v, streams=%v, score=%v, video=%v, audio=%v",
		v.Starttime, v.Duration, v.Bitrate, v.Streams, v.Score, v.HasVideo, v.HasAudio,
	)
}

func FFprobeFileFormat(ctx context.Context, filename string) (format *MediaFormat, video *FFprobeVideo, audio *FFprobeAudio, err error) {
	args := []string{
		"-show_error", "-show_private_data", "-v", "quiet", "-find_stream_info", "-print_format", "json",
		"-show_format", "-show_streams",
	}
	args = append(args, "-i", filename)

	var stdout []byte
	stdout, err = exec.CommandContext(ctx, "ffprobe", args...).Output()
	if err != nil {
		err = errors.Wrapf(err, "probe %v", filename)
		return
	}

	ffprobeFormat := struct {
		Format FFprobeFormat `json:"format"`
	}{}
	if err = json.Unmarshal([]byte(stdout), &ffprobeFormat); err != nil {
		err = errors.Wrapf(err, "parse format %v", stdout)
		return
	}

	videos := struct {
		Streams []FFprobeVideo `json:"streams"`
	}{}
	if err = json.Unmarshal([]byte(stdout), &videos); err != nil {
		err = errors.Wrapf(err, "parse video streams %v", stdout)
		return
	}
	var matchVideo *FFprobeVideo
	for _, video := range videos.Streams {
		if video.CodecType == "video" {
			matchVideo = &video
			ffprobeFormat.Format.HasVideo = true
			break
		}
	}

	audios := struct {
		Streams []FFprobeAudio `json:"streams"`
	}{}
	if err = json.Unmarshal([]byte(stdout), &audios); err != nil {
		err = errors.Wrapf(err, "parse audio streams %v", stdout)
		return
	}
	var matchAudio *FFprobeAudio
	for _, audio := range audios.Streams {
		if audio.CodecType == "audio" {
			matchAudio = &audio
			ffprobeFormat.Format.HasAudio = true
			break
		}
	}

	format = &MediaFormat{}
	if err = format.FromFFprobeFormat(&ffprobeFormat.Format); err != nil {
		err = errors.Wrapf(err, "from ffprobe format %v", ffprobeFormat.Format)
		return
	}

	video = matchVideo
	audio = matchAudio
	return
}
