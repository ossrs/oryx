//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
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

// Versions is latest and stable version from SRS Stack API.
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
	// Current working directory, at xxx/srs-stack/platform.
	Pwd string

	Cloud    string
	Region   string
	Source   string
	Registry string

	// Discover by iface.
	ipv4  net.IP
	Iface string

	// The latest and stable version from SRS Stack API.
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
	if os.Getenv("CLOUD") == "BT" {
		return "BT", "ap-beijing", nil
	}

	if os.Getenv("CLOUD") == "BIN" {
		return "BIN", "ap-beijing", nil
	}

	if os.Getenv("CLOUD") == "AAPANEL" {
		return "AAPANEL", "ap-singapore", nil
	}

	if os.Getenv("CLOUD") == "DOCKER" {
		return "DOCKER", "ap-beijing", nil
	}

	if os.Getenv("CLOUD") != "" && os.Getenv("REGION") != "" {
		return os.Getenv("CLOUD"), os.Getenv("REGION"), nil
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

// Docker container names.
// TODO: FIXME: Remove it.
const srsDockerName = "srs-server"

// Note that we only enable srs-server, never enable srs-dev.
// TODO: FIXME: Remove it.
const srsDevDockerName = "srs-dev"

// Redis keys.
const (
	// For LightHouse information, like region or source.
	SRS_TENCENT_LH = "SRS_TENCENT_LH"
	// For SRS stream status.
	SRS_HP_HLS = "SRS_HP_HLS"
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
	SRS_VLIVE_TASK   = "SRS_VLIVE_TASK"
	// For transcoding.
	SRS_TRANSCODE_CONFIG = "SRS_TRANSCODE_CONFIG"
	SRS_TRANSCODE_TASK   = "SRS_TRANSCODE_TASK"
	// For SRS stream status.
	SRS_STREAM_ACTIVE     = "SRS_STREAM_ACTIVE"
	SRS_STREAM_SRT_ACTIVE = "SRS_STREAM_SRT_ACTIVE"
	SRS_STREAM_RTC_ACTIVE = "SRS_STREAM_RTC_ACTIVE"
	// For feature statistics.
	SRS_STAT_COUNTER = "SRS_STAT_COUNTER"
	// For container and images.
	SRS_CONTAINER_DISABLED = "SRS_CONTAINER_DISABLED"
	// For system settings.
	SRS_LOCALE          = "SRS_LOCALE"
	SRS_SECRET_PUBLISH  = "SRS_SECRET_PUBLISH"
	SRS_AUTH_SECRET     = "SRS_AUTH_SECRET"
	SRS_FIRST_BOOT      = "SRS_FIRST_BOOT"
	SRS_UPGRADING       = "SRS_UPGRADING"
	SRS_UPGRADE_WINDOW  = "SRS_UPGRADE_WINDOW"
	SRS_PLATFORM_SECRET = "SRS_PLATFORM_SECRET"
	SRS_CACHE_BILIBILI  = "SRS_CACHE_BILIBILI"
	SRS_BEIAN           = "SRS_BEIAN"
	SRS_HTTPS           = "SRS_HTTPS"
	SRS_HTTPS_DOMAIN    = "SRS_HTTPS_DOMAIN"
	SRS_HOOKS           = "SRS_HOOKS"
)

// Tencent cloud consts.
const (
	TENCENT_CLOUD_CAM_ENDPOINT = "cam.tencentcloudapi.com"
	TENCENT_CLOUD_VOD_ENDPOINT = "vod.tencentcloudapi.com"
)

// SrsVLiveSourceType defines the source type of vLive.
type SrsVLiveSourceType string

const SrsVLiveSourceTypeUpload SrsVLiveSourceType = "upload"
const SrsVLiveSourceTypeFile SrsVLiveSourceType = "file"
const SrsVLiveSourceTypeStream SrsVLiveSourceType = "stream"

// For vLive upload directory.
var dirUploadPath = path.Join(".", "upload")
var dirVLivePath = path.Join(".", "vlive")

// For SRS Stack to use the files.
const serverDataDirectory = "/data"

// The files allowed to use by SRS Stack.
var serverAllowVideoFiles []string = []string{".mp4", ".flv", ".ts"}

// rdb is a global redis client object.
var rdb *redis.Client

// InitRdb create and init global rdb, which is a redis client.
func InitRdb() error {
	addr := "localhost"
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

// Refresh the ipv4 address.
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

			// The ip address might change, so we should always resolve it.
			time.Sleep(time.Duration(30) * time.Second)
		}
	}()

	select {
	case <-ctx.Done():
	case <-ipv4Ctx.Done():
	}

	return nil
}

// setEnvDefault set env key=value if not set.
func setEnvDefault(key, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}

// srsGenerateConfig is to build SRS configuration and reload SRS.
func srsGenerateConfig(ctx context.Context) error {
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Build the High Performance HLS config.
	hlsConf := []string{
		"",
		"hls {",
		"    enabled on;",
		"    hls_fragment 10;",
		"    hls_window 60;",
		"    hls_path ./objs/nginx/html;",
		"    hls_m3u8_file [app]/[stream].m3u8;",
		"    hls_ts_file [app]/[stream]-[seq]-[timestamp].ts;",
		"    hls_wait_keyframe on;",
		"    hls_dispose 10;",
	}
	if noHlsCtx, err := rdb.HGet(ctx, SRS_HP_HLS, "noHlsCtx").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v hls", SRS_HP_HLS)
	} else if noHlsCtx == "true" {
		hlsConf = append(hlsConf, []string{
			"    hls_ctx off;",
			"    hls_ts_ctx off;",
		}...)
	}
	hlsConf = append(hlsConf, "}")

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Build the config for SRS.
	if true {
		confLines := []string{
			"# !!! Important: This file is produced and maintained by the SRS Stack, please never modify it.",
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
			"# !!! Important: This file is produced and maintained by the SRS Stack, please never modify it.",
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

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Fetch the reload result, the ID which represents the reload transaction.
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

	// Reload SRS to apply the new config.
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

	// Check reload result.
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

// nginxGenerateConfig is to build NGINX configuration and reload NGINX.
func nginxGenerateConfig(ctx context.Context) error {
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Build the SSL/TLS config.
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

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Build the default root.
	// Note that it's been removed, see SRS_HTTP_PROXY.

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Build the upload limit for uploader(vLive).
	uploadLimit := []string{
		"",
		"# Limit for upload file size",
		"client_max_body_size 100g;",
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Build the config for NGINX.
	if true {
		confLines := []string{
			"# !!! Important: This file is produced and maintained by the SRS Stack, please never modify it.",
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
			"# !!! Important: This file is produced and maintained by the SRS Stack, please never modify it.",
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

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Reload NGINX to apply the new config.
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
			msg := fmt.Sprintf("SRS Stack reload Nginx at %v\n", time.Now().Format(time.RFC3339))
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
	// Note that for Transcript, the key is not used.
	Key string `json:"key,omitempty"`
	// The TS local ID, a uuid string, such as 5B7B5C03-8DB4-4ABA-AAF3-CB55902CF177
	TsID string `json:"tsid,omitempty"`
	// The TS local file, format is record/:uuid.ts, such as record/5B7B5C03-8DB4-4ABA-AAF3-CB55902CF177.ts
	File string `json:"tsfile,omitempty"`
	// The TS url, generated by SRS, such as live/livestream/2015-04-23/01/476584165.ts
	URL string `json:"url,omitempty"`
	// The seqno of TS, generated by SRS, such as 100
	SeqNo uint64 `json:"seqno,omitempty"`
	// The duration of TS in seconds, generated by SRS, such as 9.36
	Duration float64 `json:"duration,omitempty"`
	// The size of TS file in bytes, such as 1934897
	Size uint64 `json:"size,omitempty"`
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
	Definition uint64 `json:"definition"`
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
	Action SrsAction `json:"action,omitempty"`
	// The ts file path, such as ./objs/nginx/html/live/livestream/2015-04-23/01/476584165.ts
	File string `json:"file,omitempty"`
	// The duration of ts file, in seconds, such as 9.36
	Duration float64 `json:"duration,omitempty"`
	// The url of m3u8, such as live/livestream/live.m3u8
	M3u8URL string `json:"m3u8_url,omitempty"`
	// The sequence number of ts, such as 100
	SeqNo uint64 `json:"seq_no,omitempty"`

	// The vhost of stream, generated by SRS, such as video.test.com
	Vhost string `json:"vhost,omitempty"`
	// The app of stream, generated by SRS, such as live
	App string `json:"app,omitempty"`
	// The name of stream, generated by SRS, such as livestream
	Stream string `json:"stream,omitempty"`

	// The TS url, generated by SRS, such as live/livestream/2015-04-23/01/476584165.ts
	URL string `json:"url,omitempty"`
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

// ParseBody read the body from r, and unmarshal JSON to v.
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

// Authenticate check by Bearer or token.
// If use bearer secret, there is the header Authorization: Bearer {apiSecret}.
// If use token, there is a JWT token which is signed by apiSecret.
func Authenticate(ctx context.Context, apiSecret, token string, header http.Header) error {
	// Check system api secret.
	if apiSecret == "" {
		return errors.New("no api secret")
	}

	// Should use bearer secret or token.
	authorization := header.Get("Authorization")
	if authorization == "" && token == "" {
		return errors.New("no Authorization or token")
	}

	// Verify bearer secret first.
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

	// Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
	// See https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
	if _, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(apiSecret), nil
	}); err != nil {
		return errors.Wrapf(err, "verify token %v", token)
	}

	return nil
}

// httpAllowCORS allow CORS for HTTP request.
// Note that we always enable CROS because we enable HTTP cache.
func httpAllowCORS(w http.ResponseWriter, r *http.Request) {
	// SRS does not need cookie or credentials, so we disable CORS credentials, and use * for CORS origin,
	// headers, expose headers and methods.
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

// httpCreateProxy create a reverse proxy for target URL.
func httpCreateProxy(targetURL string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, errors.Wrapf(err, "parse backend %v", targetURL)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ModifyResponse = func(resp *http.Response) error {
		// We will set the server field.
		resp.Header.Del("Server")

		// We will set the CORS headers.
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
	if port := os.Getenv("RTC_PORT"); port != "8000" {
		// Read line by line, replace " 8000 " with " {port} " if contains "candidate".
		scan := bufio.NewScanner(strings.NewReader(string(b)))

		var lines []string
		for scan.Scan() {
			line := scan.Text()
			if strings.Contains(line, "candidate") {
				line = strings.ReplaceAll(line, " 8000 ", fmt.Sprintf(" %v ", port))
			}
			lines = append(lines, line)
		}

		// Join lines with "\r\n"
		sdp := strings.Join(lines, "\r\n") + "\r\n"

		return w.w.Write([]byte(sdp))
	}
	return w.w.Write(b)
}

func (w *whxpResponseModifier) WriteHeader(statusCode int) {
	w.w.WriteHeader(statusCode)
}
