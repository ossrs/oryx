// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT

//go:build linux

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	ohttp "github.com/ossrs/go-oryx-lib/http"
	"github.com/ossrs/go-oryx-lib/logger"

	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var vLiveWorker *VLiveWorker

type VLiveWorker struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// The tasks we have started to vLive streams,, key is platform in string, value is *VLiveTask.
	tasks sync.Map
}

func NewVLiveWorker() *VLiveWorker {
	return &VLiveWorker{}
}

func (v *VLiveWorker) GetTask(platform string) *VLiveTask {
	if task, loaded := v.tasks.Load(platform); loaded {
		return task.(*VLiveTask)
	}
	return nil
}

func (v *VLiveWorker) Handle(ctx context.Context, handler *http.ServeMux) error {
	ep := "/terraform/v1/ffmpeg/vlive/secret"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token, action string
			var userConf VLiveConfigure
			if err := ParseBody(ctx, r.Body, &struct {
				Token  *string `json:"token"`
				Action *string `json:"action"`
				*VLiveConfigure
			}{
				Token: &token, Action: &action, VLiveConfigure: &userConf,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			allowedActions := []string{"update"}
			allowedPlatforms := []string{"wx", "bilibili", "kuaishou"}
			if action != "" {
				if !slicesContains(allowedActions, action) {
					return errors.Errorf("invalid action=%v", action)
				}

				if userConf.Platform == "" {
					return errors.New("no platform")
				}

				// Platform should be specified platforms, or starts with vlive-.
				if !slicesContains(allowedPlatforms, userConf.Platform) && !strings.Contains(userConf.Platform, "vlive-") {
					return errors.Errorf("invalid platform=%v", userConf.Platform)
				}

				if userConf.Server == "" {
					return errors.New("no server")
				}
				if userConf.Server == "" && userConf.Secret == "" {
					return errors.New("no secret")
				}
				if len(userConf.Files) == 0 {
					return errors.New("no files")
				}
			}

			if action == "update" {
				var targetConf VLiveConfigure
				if config, err := rdb.HGet(ctx, SRS_VLIVE_CONFIG, userConf.Platform).Result(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hget %v %v", SRS_VLIVE_CONFIG, userConf.Platform)
				} else {
					if config != "" {
						if err = json.Unmarshal([]byte(config), &targetConf); err != nil {
							return errors.Wrapf(err, "unmarshal %v", config)
						}
					}
					if err = targetConf.Update(&userConf); err != nil {
						return errors.Wrapf(err, "update %v with %v", targetConf.String(), userConf.String())
					} else if newB, err := json.Marshal(&targetConf); err != nil {
						return errors.Wrapf(err, "marshal %v", targetConf.String())
					} else if err = rdb.HSet(ctx, SRS_VLIVE_CONFIG, userConf.Platform, string(newB)).Err(); err != nil && err != redis.Nil {
						return errors.Wrapf(err, "hset %v %v %v", SRS_VLIVE_CONFIG, userConf.Platform, string(newB))
					}
				}

				// Restart the vLive if exists.
				if task := vLiveWorker.GetTask(userConf.Platform); task != nil {
					if err := task.Restart(ctx); err != nil {
						return errors.Wrapf(err, "restart task %v", userConf.String())
					}
				}

				ohttp.WriteData(ctx, w, r, nil)
				logger.Tf(ctx, "vLive: Update secret ok, token=%vB", len(token))
				return nil
			} else {
				confObjs := make(map[string]*VLiveConfigure)
				if configs, err := rdb.HGetAll(ctx, SRS_VLIVE_CONFIG).Result(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hgetall %v", SRS_VLIVE_CONFIG)
				} else {
					for k, v := range configs {
						var obj VLiveConfigure
						if err = json.Unmarshal([]byte(v), &obj); err != nil {
							return errors.Wrapf(err, "unmarshal %v %v", k, v)
						}
						confObjs[k] = &obj
					}
				}

				ohttp.WriteData(ctx, w, r, confObjs)
				logger.Tf(ctx, "vLive: Query configures ok, token=%vB", len(token))
				return nil
			}
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ffmpeg/vlive/streams"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			res := make([]map[string]interface{}, 0)
			if configs, err := rdb.HGetAll(ctx, SRS_VLIVE_CONFIG).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hgetall %v", SRS_VLIVE_CONFIG)
			} else if len(configs) > 0 {
				for k, v := range configs {
					var config VLiveConfigure
					if err = json.Unmarshal([]byte(v), &config); err != nil {
						return errors.Wrapf(err, "unmarshal %v %v", k, v)
					}

					var pid int32
					var inputUUID, frame, update, starttime, ready string
					if task := vLiveWorker.GetTask(config.Platform); task != nil {
						pid, inputUUID, frame, update, starttime, ready = task.queryFrame()
					}

					elem := map[string]interface{}{
						"platform": config.Platform,
						"enabled":  config.Enabled,
						"custom":   config.Customed,
						"label":    config.Label,
						"files":    config.Files,
					}

					if pid > 0 {
						elem["source"] = inputUUID
						elem["start"] = starttime
						elem["ready"] = ready
						elem["frame"] = map[string]string{
							"log":    frame,
							"update": update,
						}
					}

					res = append(res, elem)
				}
			}

			sort.Slice(res, func(i, j int) bool {
				return res[i]["platform"].(string) < res[j]["platform"].(string)
			})

			ohttp.WriteData(ctx, w, r, res)
			logger.Tf(ctx, "vLive: Query vLive streams ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	streamUrlHandler := func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var qUrl string
			if err := ParseBody(ctx, r.Body, &struct {
				Token     *string `json:"token"`
				StreamURL *string `json:"url"`
			}{
				Token: &token, StreamURL: &qUrl,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			// Parse URL to object.
			u, err := RebuildStreamURL(qUrl)
			if err != nil {
				return errors.Wrapf(err, "parse %v", qUrl)
			}

			// Check url if valid rtmp or rtsp or http-flv or https-flv or hls live url
			if u.Scheme != "rtmp" && u.Scheme != "srt" && u.Scheme != "rtsp" && u.Scheme != "http" && u.Scheme != "https" {
				return errors.Errorf("invalid url scheme %v", u.Scheme)
			}
			if u.Scheme == "http" || u.Scheme == "https" {
				if u.Path == "" {
					return errors.Errorf("url path %v empty", u.Path)
				}
				if !strings.HasSuffix(u.Path, ".flv") && !strings.HasSuffix(u.Path, ".m3u8") && !strings.HasSuffix(u.Path, ".ts") {
					return errors.Errorf("invalid url path suffix %v", u.Path)
				}
			}

			targetUUID := uuid.NewString()
			ohttp.WriteData(ctx, w, r, &struct {
				// The file name.
				Name string `json:"name"`
				// The file UUID.
				UUID string `json:"uuid"`
				// The file target name.
				Target string `json:"target"`
			}{
				Name:   path.Base(u.Path),
				UUID:   targetUUID,
				Target: qUrl,
			})

			logger.Tf(ctx, "vLive: Update stream url ok, url=%v, uuid=%v", qUrl, targetUUID)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	}

	ep = "/terraform/v1/ffmpeg/vlive/streamUrl"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, streamUrlHandler)

	ep = "/terraform/v1/ffmpeg/vlive/stream-url"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, streamUrlHandler)

	ep = "/terraform/v1/ffmpeg/vlive/ytdl"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var qFile string
			if err := ParseBody(ctx, r.Body, &struct {
				Token   *string `json:"token"`
				YtdlURL *string `json:"url"`
			}{
				Token: &token, YtdlURL: &qFile,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if !strings.HasPrefix(qFile, "http") && !strings.HasPrefix(qFile, "https") {
				return errors.Errorf("invalid url %v", qFile)
			}

			// If upload directory is symlink, eval it.
			targetDir := dirUploadPath
			if info, err := os.Lstat(targetDir); err == nil && info.Mode()&os.ModeSymlink != 0 {
				if realPath, err := filepath.EvalSymlinks(targetDir); err != nil {
					return errors.Wrapf(err, "eval symlink %v", targetDir)
				} else {
					targetDir = realPath
				}
			}

			// The prefix for target files.
			targetUUID := uuid.NewString()

			// Cleanup all temporary files created by youtube-dl.
			requestDone, requestDoneCancel := context.WithCancel(context.Background())
			defer requestDoneCancel()
			go func() {
				// If the temporary file still exists for a long time, remove it
				duration := 2 * time.Hour
				if envNodeEnv() == "development" {
					duration = time.Duration(30) * time.Second
				}

				select {
				case <-ctx.Done():
					logger.Tf(ctx, "ytdl: do cleanup immediately when quit")
				case <-requestDone.Done():
					time.Sleep(duration)
				}

				filepath.WalkDir(targetDir, func(p string, info fs.DirEntry, err error) error {
					if err != nil {
						return errors.Wrapf(err, "walk %v", p)
					}

					if !info.IsDir() && strings.HasPrefix(info.Name(), targetUUID) {
						tempFile := path.Join(dirUploadPath, info.Name())
						if _, err := os.Stat(tempFile); err == nil {
							os.Remove(tempFile)
							logger.Wf(ctx, "remove %v, duration=%v", tempFile, duration)
						}
					}

					return nil
				})
			}()

			// Use youtube-dl to download the file.
			ytdlOutput := path.Join(targetDir, targetUUID)
			args := []string{
				"--output", ytdlOutput,
			}
			if proxy := envYtdlProxy(); proxy != "" {
				args = append(args, "--proxy", proxy)
			}
			args = append(args, qFile)
			if err := exec.CommandContext(ctx, "youtube-dl", args...).Run(); err != nil {
				return errors.Wrapf(err, "run youtube-dl %v", args)
			}

			// Find out the downloaded target file.
			var targetFile string
			if err := filepath.WalkDir(targetDir, func(p string, info fs.DirEntry, err error) error {
				if err != nil {
					return errors.Wrapf(err, "walk %v", p)
				}

				if !info.IsDir() && strings.HasPrefix(info.Name(), targetUUID) {
					targetFile = path.Join(dirUploadPath, info.Name())
					return filepath.SkipDir
				}

				return nil
			}); err != nil {
				return errors.Wrapf(err, "walk %v", targetDir)
			}

			if targetFile == "" {
				return errors.Errorf("no target file %v", targetUUID)
			}

			// Get the file information.
			targetFileInfo, err := os.Lstat(targetFile)
			if err != nil {
				return errors.Wrapf(err, "lstat %v", targetFile)
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Name   string `json:"name"`
				UUID   string `json:"uuid"`
				Target string `json:"target"`
				Size   int    `json:"size"`
			}{
				Name:   targetFileInfo.Name(),
				UUID:   targetUUID,
				Target: targetFile,
				Size:   int(targetFileInfo.Size()),
			})
			logger.Tf(ctx, "vLive: Got vlive ytdl file target=%v, size=%v", targetFileInfo.Name(), targetFileInfo.Size())
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ffmpeg/vlive/server"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var qFile string
			if err := ParseBody(ctx, r.Body, &struct {
				Token      *string `json:"token"`
				StreamFile *string `json:"file"`
			}{
				Token: &token, StreamFile: &qFile,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			fileAbsPath, err := filepath.Abs(qFile)
			if err != nil {
				return errors.Wrapf(err, "abs %v", qFile)
			}

			if !strings.HasPrefix(fileAbsPath, serverDataDirectory) && !strings.HasPrefix(qFile, dirUploadPath) {
				return errors.Errorf("invalid file %v, should in %v", fileAbsPath, serverDataDirectory)
			}

			var validExtension bool
			for _, ext := range append(serverAllowVideoFiles, serverAllowAudioFiles...) {
				if strings.HasSuffix(fileAbsPath, ext) {
					validExtension = true
					break
				}
			}
			if !validExtension {
				return errors.Errorf("invalid file extension %v, should be %v",
					fileAbsPath, append(serverAllowVideoFiles, serverAllowAudioFiles...))
			}

			info, err := os.Stat(fileAbsPath)
			if err != nil {
				return errors.Wrapf(err, "stat %v", fileAbsPath)
			}

			targetUUID := uuid.NewString()
			targetFileName := path.Join(dirUploadPath, fmt.Sprintf("%v%v", targetUUID, path.Ext(info.Name())))

			targetFile, err := os.OpenFile(targetFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return errors.Wrapf(err, "open file %v", targetFileName)
			}
			defer targetFile.Close()

			sourceFile, err := os.Open(fileAbsPath)
			if err != nil {
				return errors.Wrapf(err, "open file %v", fileAbsPath)
			}
			defer sourceFile.Close()

			// TODO: FIXME: Cleanup the file if error.
			if _, err = io.Copy(targetFile, sourceFile); err != nil {
				return errors.Wrapf(err, "copy %v to %v", fileAbsPath, targetFileName)
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Name   string `json:"name"`
				UUID   string `json:"uuid"`
				Target string `json:"target"`
				Size   int    `json:"size"`
			}{
				Name:   info.Name(),
				UUID:   targetUUID,
				Target: targetFileName,
				Size:   int(info.Size()),
			})
			logger.Tf(ctx, "vLive: Got vlive local file target=%v, size=%v", targetFileName, info.Size())
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ffmpeg/vlive/upload/"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func(ctx context.Context) error {
			filename := r.URL.Path[len("/terraform/v1/ffmpeg/vlive/upload/"):]

			targetUUID := uuid.NewString()
			targetFileName := path.Join(dirUploadPath, fmt.Sprintf("%v%v", targetUUID, path.Ext(filename)))
			logger.Tf(ctx, "vLive: Create %v for %v", targetFileName, filename)

			var uploadDone bool
			created := time.Now()
			defer func() {
				// If upload is not done, remove the target file
				if !uploadDone {
					os.Remove(targetFileName)
					logger.Wf(ctx, "remove %v, done=%v, created=%v", targetFileName, uploadDone, created)
				}
			}()

			requestDone, requestDoneCancel := context.WithCancel(context.Background())
			defer requestDoneCancel()
			go func() {
				// If the temporary file still exists for a long time, remove it
				duration := 2 * time.Hour
				if envNodeEnv() == "development" {
					duration = time.Duration(30) * time.Second
				}

				select {
				case <-ctx.Done():
					logger.Tf(ctx, "upload: do cleanup immediately when quit")
				case <-requestDone.Done():
					time.Sleep(duration)
				}

				if _, err := os.Stat(targetFileName); err == nil {
					os.Remove(targetFileName)
					logger.Wf(ctx, "remove %v, done=%v, created=%v, duration=%v, elapsed=%v",
						targetFileName, uploadDone, created, duration, time.Now().Sub(created))
				}
			}()

			targetFile, err := os.OpenFile(targetFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return errors.Wrapf(err, "open file %v", targetFileName)
			}
			defer targetFile.Close()

			// See https://github.com/rfielding/uploader/blob/master/uploader.go#L170
			mr, err := r.MultipartReader()
			if err != nil {
				return errors.Wrapf(err, "multi reader")
			}

			starttime := time.Now()
			var written int64
			for {
				part, err := mr.NextPart()
				if err != nil {
					if err == io.EOF {
						break
					}
					return errors.Wrapf(err, "next part")
				}

				if filename != part.FileName() {
					return errors.Errorf("filename=%v mismatch %v", part.FileName(), filename)
				}
				logger.Tf(ctx, "vLive: Start part for %v", targetFileName)

				partStarttime := time.Now()
				if nn, err := io.Copy(targetFile, part); err != nil {
					return errors.Wrapf(err, "copy %v to %v", targetFile, filename)
				} else {
					written += nn
					logger.Tf(ctx, "vLive: Finish part for %v, nn=%v, writen=%v, cost=%v",
						targetFileName, nn, written, time.Now().Sub(partStarttime),
					)
				}
			}

			// After write file success, set the upload done to keep the file.
			uploadDone = true

			ohttp.WriteData(ctx, w, r, &struct {
				// The file UUID.
				UUID string `json:"uuid"`
				// The target file name.
				Target string `json:"target"`
			}{
				UUID: targetUUID, Target: targetFileName,
			})
			logger.Tf(ctx, "vLive: Got vlive target=%v, size=%v, done=%v, cost=%v", targetFileName, written, uploadDone, time.Now().Sub(starttime))
			return nil
		}(logger.WithContext(ctx)); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ffmpeg/vlive/source"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			type VLiveTempFile struct {
				// The file name.
				Name string `json:"name"`
				// The file path.
				Path string `json:"path"`
				// The file size in bytes.
				Size int64 `json:"size"`
				// The UUID for file.
				UUID string `json:"uuid"`
				// The target file name.
				Target string `json:"target"`
				// The source type.
				Type FFprobeSourceType `json:"type"`
			}

			var token, platform string
			var files []*VLiveTempFile
			if err := ParseBody(ctx, r.Body, &struct {
				Token    *string           `json:"token"`
				Platform *string           `json:"platform"`
				Files    *[]*VLiveTempFile `json:"files"`
			}{
				Token: &token, Platform: &platform, Files: &files,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if len(files) == 0 {
				return errors.New("no files")
			}

			// Always cleanup the files in upload.
			var tempFiles []string
			for _, f := range files {
				if f.Type != FFprobeSourceTypeStream {
					tempFiles = append(tempFiles, f.Target)
				}
			}
			defer func() {
				for _, tempFile := range tempFiles {
					if _, err := os.Stat(tempFile); err == nil {
						os.Remove(tempFile)
						logger.Tf(ctx, "vLive: Cleanup %v", tempFile)
					}
				}
			}()

			// Check files.
			for _, f := range files {
				if f.Target == "" {
					return errors.New("no target")
				}
				if f.Type != FFprobeSourceTypeStream {
					if _, err := os.Stat(f.Target); err != nil {
						return errors.Wrapf(err, "no file %v", f.Target)
					}
					if !strings.HasPrefix(f.Target, dirUploadPath) {
						return errors.Errorf("invalid target %v", f.Target)
					}
				}
			}

			// Check platform.
			if platform == "" {
				return errors.New("no platform")
			}
			// For virtual live event only.
			if true {
				allowedPlatforms := []string{"wx", "bilibili", "kuaishou"}

				// Platform should be specified platforms, or starts with vlive-.
				if !slicesContains(allowedPlatforms, platform) && !strings.Contains(platform, "vlive-") {
					return errors.Errorf("invalid platform %v", platform)
				}
			}

			// Parsed source files.
			var parsedFiles []*FFprobeSource

			// Parse file information and move file from dirUploadPath to dirVLivePath.
			for _, file := range files {
				// Probe file information.
				toCtx, toCancelFunc := context.WithTimeout(ctx, 15*time.Second)
				defer toCancelFunc()

				args := []string{
					"-show_error", "-show_private_data", "-v", "quiet", "-find_stream_info", "-print_format", "json",
					"-show_format", "-show_streams",
				}
				// For RTSP stream source, always use TCP transport.
				if strings.HasPrefix(file.Target, "rtsp://") {
					args = append(args, "-rtsp_transport", "tcp")
				}
				// Rebuild the stream url, because it may contain special characters.
				if strings.Contains(file.Target, "://") {
					if u, err := RebuildStreamURL(file.Target); err != nil {
						return errors.Wrapf(err, "rebuild %v", file.Target)
					} else {
						args = append(args, "-i", u.String())
					}
				} else {
					args = append(args, "-i", file.Target)
				}

				// TODO: FIXME: Use FFprobeFileFormat.
				stdout, err := exec.CommandContext(toCtx, "ffprobe", args...).Output()
				if err != nil {
					return errors.Wrapf(err, "probe %v with ffprobe %v", file.Target, args)
				}

				format := struct {
					Format FFprobeFormat `json:"format"`
				}{}
				if err = json.Unmarshal([]byte(stdout), &format); err != nil {
					return errors.Wrapf(err, "parse format %v", stdout)
				}

				// Typically, AWS Lightsail and DigitalOcean Droplets provide 1TB of monthly traffic,
				// permitting a 3Mbps continuous live stream for 7x24 hours. Therefore, it's crucial
				// to restrict the input bitrate to prevent exceeding the traffic limit.
				if format.Format.Bitrate != "" {
					if limits, err := rdb.HGet(ctx, SRS_SYS_LIMITS, "vlive").Int64(); err != nil && err != redis.Nil {
						return errors.Wrapf(err, "hget %v vlive", SRS_SYS_LIMITS)
					} else {
						if limits == 0 {
							limits = SrsSysLimitsVLive // in Kbps.
						}

						if bitrate, err := strconv.ParseInt(format.Format.Bitrate, 10, 64); err != nil {
							return errors.Wrapf(err, "parse bitrate %v", format.Format.Bitrate)
						} else if bitrate > limits*1000 {
							return errors.Errorf("bitrate %vKbps is too large, exceed %vKbps", bitrate/1000, limits)
						}
					}
				}

				videos := struct {
					Streams []FFprobeVideo `json:"streams"`
				}{}
				if err = json.Unmarshal([]byte(stdout), &videos); err != nil {
					return errors.Wrapf(err, "parse video streams %v", stdout)
				}
				var matchVideo *FFprobeVideo
				for _, video := range videos.Streams {
					if video.CodecType == "video" {
						matchVideo = &video
						format.Format.HasVideo = true
						break
					}
				}

				audios := struct {
					Streams []FFprobeAudio `json:"streams"`
				}{}
				if err = json.Unmarshal([]byte(stdout), &audios); err != nil {
					return errors.Wrapf(err, "parse audio streams %v", stdout)
				}
				var matchAudio *FFprobeAudio
				for _, audio := range audios.Streams {
					if audio.CodecType == "audio" {
						matchAudio = &audio
						format.Format.HasAudio = true
						break
					}
				}

				// Only accept common codec for video and audio.
				allowedCodec := []string{"h264", "h265", "aac", "mp3"}
				if matchVideo != nil && !slicesContains(allowedCodec, matchVideo.CodecName) {
					return errors.Errorf("invalid video codec %v, should be %v", matchVideo.CodecName, allowedCodec)
				}
				if matchAudio != nil && !slicesContains(allowedCodec, matchAudio.CodecName) {
					return errors.Errorf("invalid audio codec %v, should be %v", matchAudio.CodecName, allowedCodec)
				}

				parsedFile := &FFprobeSource{
					Name: file.Name, Path: file.Path, Size: uint64(file.Size), UUID: file.UUID,
					Target: file.Target,
					Type:   file.Type,
					Format: &format.Format, Video: matchVideo, Audio: matchAudio,
				}
				if file.Type != FFprobeSourceTypeStream {
					parsedFile.Target = path.Join(dirVLivePath, fmt.Sprintf("%v%v", file.UUID, path.Ext(file.Target)))
					if err = os.Rename(file.Target, parsedFile.Target); err != nil {
						return errors.Wrapf(err, "rename %v to %v", file.Target, parsedFile.Target)
					}
				}

				parsedFiles = append(parsedFiles, parsedFile)
				logger.Tf(ctx, "vLive: Process file %v", parsedFile.String())
			}

			// For virtual live stream only.
			if true {
				// Update redis object.
				confObj := VLiveConfigure{Platform: platform}
				if conf, err := rdb.HGet(ctx, SRS_VLIVE_CONFIG, platform).Result(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hget %v %v", SRS_VLIVE_CONFIG, platform)
				} else if conf != "" {
					if err = json.Unmarshal([]byte(conf), &confObj); err != nil {
						return errors.Wrapf(err, "parse %v", conf)
					}
				}

				// Remove old files.
				for _, f := range confObj.Files {
					if f.Type != FFprobeSourceTypeStream {
						if _, err := os.Stat(f.Target); err == nil {
							os.Remove(f.Target)
						}
					}
				}
				confObj.Files = parsedFiles

				if b, err := json.Marshal(&confObj); err != nil {
					return errors.Wrapf(err, "marshal %v", confObj.String())
				} else if err = rdb.HSet(ctx, SRS_VLIVE_CONFIG, platform, string(b)).Err(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hset %v %v %v", SRS_VLIVE_CONFIG, platform, string(b))
				}

				// Restart the vLive if exists.
				if task := vLiveWorker.GetTask(platform); task != nil {
					if err := task.Restart(ctx); err != nil {
						return errors.Wrapf(err, "restart task %v", platform)
					}
				}
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Platform string           `json:"platform"`
				Files    []*FFprobeSource `json:"files"`
			}{
				Platform: platform, Files: parsedFiles,
			})
			logger.Tf(ctx, "vLive: Update vLive ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	return nil
}

func (v *VLiveWorker) Close() error {
	if v.cancel != nil {
		v.cancel()
	}
	v.wg.Wait()
	return nil
}

func (v *VLiveWorker) Start(ctx context.Context) error {
	wg := &v.wg

	ctx, cancel := context.WithCancel(ctx)
	v.cancel = cancel

	ctx = logger.WithContext(ctx)
	logger.Tf(ctx, "vLive: Start a worker")

	// Load tasks from redis and force to kill all.
	if objs, err := rdb.HGetAll(ctx, SRS_VLIVE_TASK).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hgetall %v", SRS_VLIVE_TASK)
	} else if len(objs) > 0 {
		for uuid, obj := range objs {
			logger.Tf(ctx, "vLive: Load task %v object %v", uuid, obj)

			var task VLiveTask
			if err = json.Unmarshal([]byte(obj), &task); err != nil {
				return errors.Wrapf(err, "unmarshal %v %v", uuid, obj)
			}

			if task.PID > 0 {
				task.cleanup(ctx)
			}
		}

		if err = rdb.Del(ctx, SRS_VLIVE_TASK).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "del %v", SRS_VLIVE_TASK)
		}
	}

	// Load all configurations from redis.
	loadTasks := func() error {
		configItems, err := rdb.HGetAll(ctx, SRS_VLIVE_CONFIG).Result()
		if err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hgetall %v", SRS_VLIVE_CONFIG)
		}
		if len(configItems) == 0 {
			return nil
		}

		for platform, configItem := range configItems {
			var config VLiveConfigure
			if err = json.Unmarshal([]byte(configItem), &config); err != nil {
				return errors.Wrapf(err, "unmarshal %v %v", platform, configItem)
			}

			var task *VLiveTask
			if tv, loaded := v.tasks.LoadOrStore(config.Platform, &VLiveTask{
				UUID:     uuid.NewString(),
				Platform: config.Platform,
				config:   &config,
			}); loaded {
				// Ignore if exists.
				continue
			} else {
				task = tv.(*VLiveTask)
				logger.Tf(ctx, "vLive: Create platform=%v task is %v", platform, task.String())
			}

			// Initialize object.
			if err := task.Initialize(ctx, v); err != nil {
				return errors.Wrapf(err, "init %v", task.String())
			}

			// Store in memory object.
			v.tasks.Store(platform, task)

			wg.Add(1)
			go func() {
				defer wg.Done()

				if err := task.Run(ctx); err != nil {
					logger.Wf(ctx, "run task %v err %+v", task.String(), err)
				}
			}()
		}

		return nil
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		// When startup, we try to wait for client to publish streams.
		select {
		case <-ctx.Done():
		case <-time.After(3 * time.Second):
		}
		logger.Tf(ctx, "vLive: Start to run tasks")

		for ctx.Err() == nil {
			duration := 3 * time.Second
			if err := loadTasks(); err != nil {
				logger.Wf(ctx, "ignore err %+v", err)
				duration = 10 * time.Second
			}

			select {
			case <-ctx.Done():
			case <-time.After(duration):
			}
		}
	}()

	return nil
}

// VLiveConfigure is the configure for vLive.
type VLiveConfigure struct {
	// The platform name, for example, wx
	Platform string `json:"platform"`
	// The RTMP server url, for example, rtmp://localhost/live
	Server string `json:"server"`
	// The RTMP stream and secret, for example, livestream
	Secret string `json:"secret"`
	// Whether enabled.
	Enabled bool `json:"enabled"`
	// Whether custom platform.
	Customed bool `json:"custom"`
	// The label for this configure.
	Label string `json:"label"`

	// The input files for vLive.
	Files []*FFprobeSource `json:"files"`
}

func (v VLiveConfigure) String() string {
	return fmt.Sprintf("platform=%v, server=%v, secret=%v, enabled=%v, customed=%v, label=%v, files=%v",
		v.Platform, v.Server, v.Secret, v.Enabled, v.Customed, v.Label, v.Files,
	)
}

func (v *VLiveConfigure) Update(u *VLiveConfigure) error {
	v.Platform = u.Platform
	v.Server = u.Server
	v.Secret = u.Secret
	v.Label = u.Label
	v.Enabled = u.Enabled
	v.Customed = u.Customed
	v.Files = append([]*FFprobeSource{}, u.Files...)
	return nil
}

// VLiveTask is a task for FFmpeg to vLive stream, with a configure.
type VLiveTask struct {
	// The ID for task.
	UUID string `json:"uuid"`
	// The platform for task.
	Platform string `json:"platform"`

	// The input file path.
	Input string `json:"input"`
	// The input file UUID.
	inputUUID string
	// The output url
	Output string `json:"output"`

	// FFmpeg pid.
	PID int32 `json:"pid"`
	// FFmpeg last frame.
	frame string
	// The last update time.
	update *time.Time
	// The task start time.
	starttime *time.Time
	// The first ready time.
	firstReadyTime *time.Time

	// The context for current task.
	cancel context.CancelFunc

	// The configure for vLive task.
	config *VLiveConfigure
	// The vLive worker.
	vLiveWorker *VLiveWorker

	// To protect the fields.
	lock sync.Mutex
}

func (v *VLiveTask) String() string {
	return fmt.Sprintf("uuid=%v, platform=%v, input=%v, output=%v, pid=%v, frame=%vB, config is %v",
		v.UUID, v.Platform, v.Input, v.Output, v.PID, len(v.frame), v.config.String(),
	)
}

func (v *VLiveTask) saveTask(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if b, err := json.Marshal(v); err != nil {
		return errors.Wrapf(err, "marshal %v", v.String())
	} else if err = rdb.HSet(ctx, SRS_VLIVE_TASK, v.UUID, string(b)).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_VLIVE_TASK, v.UUID, string(b))
	}

	return nil
}

func (v *VLiveTask) cleanup(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if v.PID <= 0 {
		return nil
	}

	logger.Wf(ctx, "kill task pid=%v", v.PID)
	syscall.Kill(int(v.PID), syscall.SIGKILL)

	v.PID = 0
	v.cancel = nil

	return nil
}

func (v *VLiveTask) Restart(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if v.cancel != nil {
		v.cancel()
	}

	// Reload config from redis.
	if b, err := rdb.HGet(ctx, SRS_VLIVE_CONFIG, v.Platform).Result(); err != nil {
		return errors.Wrapf(err, "hget %v %v", SRS_VLIVE_CONFIG, v.Platform)
	} else if err = json.Unmarshal([]byte(b), v.config); err != nil {
		return errors.Wrapf(err, "unmarshal %v", b)
	}

	return nil
}

func (v *VLiveTask) updateFrame(frame string) {
	v.lock.Lock()
	defer v.lock.Unlock()

	v.frame = frame

	var now = time.Now()
	v.update = &now
}

func (v *VLiveTask) queryFrame() (int32, string, string, string, string, string) {
	v.lock.Lock()
	defer v.lock.Unlock()

	ready := ""
	if v.firstReadyTime != nil {
		ready = v.firstReadyTime.Format(time.RFC3339)
	}

	update := ""
	if v.update != nil {
		update = v.update.Format(time.RFC3339)
	}

	starttime := ""
	if v.starttime != nil {
		starttime = v.starttime.Format(time.RFC3339)
	}
	return v.PID, v.inputUUID, v.frame, update, starttime, ready
}

func (v *VLiveTask) Initialize(ctx context.Context, w *VLiveWorker) error {
	v.vLiveWorker = w
	logger.Tf(ctx, "vLive: Initialize uuid=%v, platform=%v", v.UUID, v.Platform)

	if err := v.saveTask(ctx); err != nil {
		return errors.Wrapf(err, "save task")
	}

	return nil
}

func (v *VLiveTask) Run(ctx context.Context) error {
	ctx = logger.WithContext(ctx)
	logger.Tf(ctx, "vLive: Run task %v", v.String())

	selectInputFile := func() *FFprobeSource {
		v.lock.Lock()
		defer v.lock.Unlock()

		if len(v.config.Files) == 0 {
			return nil
		}

		file := v.config.Files[0]
		logger.Tf(ctx, "vLive: Use file=%v as input for platform=%v", file.UUID, v.Platform)
		return file
	}

	pfn := func(ctx context.Context) error {
		// Ignore when not enabled.
		if !v.config.Enabled {
			return nil
		}

		// Use a active stream as input.
		input := selectInputFile()
		if input == nil {
			return nil
		}

		// Start vLive task.
		if err := v.doVirtualLiveStream(ctx, input); err != nil {
			return errors.Wrapf(err, "do vLive")
		}

		return nil
	}

	for ctx.Err() == nil {
		if err := pfn(ctx); err != nil {
			logger.Wf(ctx, "ignore %v err %+v", v.String(), err)

			select {
			case <-ctx.Done():
			case <-time.After(3500 * time.Millisecond):
			}
			continue
		}

		select {
		case <-ctx.Done():
		case <-time.After(300 * time.Millisecond):
		}
	}

	return nil
}

func (v *VLiveTask) doVirtualLiveStream(ctx context.Context, input *FFprobeSource) error {
	// Create context for current task.
	parentCtx := ctx
	ctx, cancel := context.WithCancel(ctx)
	v.cancel = cancel

	// Build input URL.
	host := "localhost"

	// Build output URL.
	outputServer := strings.ReplaceAll(v.config.Server, "localhost", host)
	if !strings.HasSuffix(outputServer, "/") && !strings.HasPrefix(v.config.Secret, "/") && v.config.Secret != "" {
		outputServer += "/"
	}
	outputURL := fmt.Sprintf("%v%v", outputServer, v.config.Secret)

	// Create a heartbeat to poll and manage the status of FFmpeg process.
	heartbeat := NewFFmpegHeartbeat(cancel)
	v.starttime, v.firstReadyTime = &heartbeat.starttime, nil
	defer func() {
		v.starttime = nil
	}()

	// Start FFmpeg process.
	args := []string{}
	if input.Type == FFprobeSourceTypeFile || input.Type == FFprobeSourceTypeUpload || input.Type == FFprobeSourceTypeYTDL {
		args = append(args, "-stream_loop", "-1")
		args = append(args, "-re")
	}
	// For RTSP stream source, always use TCP transport.
	if strings.HasPrefix(input.Target, "rtsp://") {
		args = append(args, "-rtsp_transport", "tcp")
	}
	// Rebuild the stream url, because it may contain special characters.
	if strings.Contains(input.Target, "://") {
		if u, err := RebuildStreamURL(input.Target); err != nil {
			return errors.Wrapf(err, "rebuild %v", input.Target)
		} else {
			args = append(args, "-i", u.String())
			heartbeat.Parse(u)
		}
	} else {
		args = append(args, "-i", input.Target)
	}
	args = append(args, "-c", "copy")
	// If RTMP use flv, if SRT use mpegts, otherwise do not set.
	if strings.HasPrefix(outputURL, "rtmp://") || strings.HasPrefix(outputURL, "rtmps://") {
		args = append(args, "-f", "flv")
	} else if strings.HasPrefix(outputURL, "srt://") {
		args = append(args, "-pes_payload_size", "0", "-f", "mpegts")
	}
	args = append(args, outputURL)
	// Create the command object.
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrapf(err, "pipe process")
	}

	if err := cmd.Start(); err != nil {
		return errors.Wrapf(err, "execute ffmpeg %v", strings.Join(args, " "))
	}

	v.PID = int32(cmd.Process.Pid)
	v.Input, v.inputUUID, v.Output = input.Target, input.UUID, outputURL
	defer func() {
		// If we got a PID, sleep for a while, to avoid too fast restart.
		if v.PID > 0 {
			select {
			case <-ctx.Done():
			case <-time.After(1 * time.Second):
			}
		}

		// When canceled, we should still write to redis, so we must not use ctx(which is cancelled).
		v.cleanup(parentCtx)
		v.saveTask(parentCtx)
	}()
	logger.Tf(ctx, "vLive: Start, platform=%v, input=%v, pid=%v", v.Platform, input.Target, v.PID)

	if err := v.saveTask(ctx); err != nil {
		return errors.Wrapf(err, "save task %v", v.String())
	}

	// Pull the latest log frame.
	heartbeat.Polling(ctx, stderr)
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-heartbeat.firstReadyCtx.Done():
			v.firstReadyTime = &heartbeat.firstReadyTime
		}

		for {
			select {
			case <-ctx.Done():
				return
			case frame := <-heartbeat.FrameLogs:
				v.updateFrame(frame)
			}
		}
	}()

	// Process terminated, or user cancel the process.
	select {
	case <-parentCtx.Done():
	case <-ctx.Done():
	case <-heartbeat.PollingCtx.Done():
	}
	logger.Tf(ctx, "vLive: Cycle stopping, platform=%v, input=%v, pid=%v", v.Platform, input.Target, v.PID)

	err = cmd.Wait()
	logger.Tf(ctx, "vLive: Cycle done, platform=%v, input=%v, pid=%v, err=%v",
		v.Platform, input.Target, v.PID, err,
	)

	return err
}
