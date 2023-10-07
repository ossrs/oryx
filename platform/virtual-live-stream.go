//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
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

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
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
				if !slicesContains(allowedPlatforms, userConf.Platform) {
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
				logger.Tf(ctx, "vLive update secret ok, token=%vB", len(token))
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
				logger.Tf(ctx, "vLive query configures ok, token=%vB", len(token))
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

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			res := make([]map[string]interface{}, 0)
			if configs, err := rdb.HGetAll(ctx, SRS_VLIVE_CONFIG).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hgetall %v", SRS_VLIVE_CONFIG)
			} else if len(configs) > 0 {
				for k, v := range configs {
					var conf VLiveConfigure
					if err = json.Unmarshal([]byte(v), &conf); err != nil {
						return errors.Wrapf(err, "unmarshal %v %v", k, v)
					}

					var pid int32
					var inputUUID, frame, update string
					if task := vLiveWorker.GetTask(conf.Platform); task != nil {
						pid, inputUUID, frame, update = task.queryFrame()
					}

					elem := map[string]interface{}{
						"platform": conf.Platform,
						"enabled":  conf.Enabled,
						"custom":   conf.Customed,
						"label":    conf.Label,
						"files":    conf.Files,
					}

					if pid > 0 {
						elem["source"] = inputUUID
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
			logger.Tf(ctx, "Query vLive streams ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ffmpeg/vlive/streamUrl/"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			q := r.URL.Query()
			qUrl := q.Get("url")
			u, err := url.Parse(qUrl)
			if err != nil {
				return errors.Wrapf(err, "parse %v", qUrl)
			}
			// check url if valid rtmp or http-flv or https-flv or hls live url
			if u.Scheme != "rtmp" && u.Scheme != "http" && u.Scheme != "https" {
				return errors.Errorf("invalid url scheme %v", u.Scheme)
			}
			if u.Scheme == "http" || u.Scheme == "https" {
				if u.Path == "" {
					return errors.Errorf("url path %v empty", u.Path)
				}
				if !strings.HasSuffix(u.Path, ".flv") && !strings.HasSuffix(u.Path, ".m3u8") {
					return errors.Errorf("invalid url path suffix %v", u.Path)
				}
			}
			targetUUID := uuid.NewString()
			ohttp.WriteData(ctx, w, r, &struct {
				Name string `json:"name"`
				UUID string `json:"uuid"`
				Target string `json:"target"`
			}{
				 Name: path.Base(u.Path),
				 UUID: targetUUID,
				 Target: qUrl,
			})
			logger.Tf(ctx, "vLive stream url ok, url=%v, uuid=%v", qUrl, targetUUID)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ffmpeg/vlive/server/"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			q := r.URL.Query()
			qFile := filepath.Clean(q.Get("file"))
			fileAbsPath, err := filepath.Abs(qFile)
			if err != nil {
				return errors.Wrapf(err, "abs %v", qFile)
			}

			if !strings.HasPrefix(fileAbsPath, serverDataDirectory) && !strings.HasPrefix(qFile, dirUploadPath) {
				return errors.Errorf("invalid file %v, should in %v", fileAbsPath, serverDataDirectory)
			}

			var validExtension bool
			for _, ext := range serverAllowVideoFiles {
				if strings.HasSuffix(fileAbsPath, ext) {
					validExtension = true
					break
				}
			}
			if !validExtension {
				return errors.Errorf("invalid file extension %v, should be %v", fileAbsPath, serverAllowVideoFiles)
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
			logger.Tf(ctx, "Got vlive local file target=%v, size=%v", targetFileName, info.Size())
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
			logger.Tf(ctx, "create %v for %v", targetFileName, filename)

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
				logger.Tf(ctx, "start part for %v", targetFileName)

				partStarttime := time.Now()
				if nn, err := io.Copy(targetFile, part); err != nil {
					return errors.Wrapf(err, "copy %v to %v", targetFile, filename)
				} else {
					written += nn
					logger.Tf(ctx, "finish part for %v, nn=%v, writen=%v, cost=%v",
						targetFileName, nn, written, time.Now().Sub(partStarttime),
					)
				}
			}

			ohttp.WriteData(ctx, w, r, &struct {
				UUID   string `json:"uuid"`
				Target string `json:"target"`
			}{
				UUID: targetUUID, Target: targetFileName,
			})
			logger.Tf(ctx, "Got vlive target=%v, size=%v, cost=%v", targetFileName, written, time.Now().Sub(starttime))
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
				Name   string `json:"name"`
				Size   int64  `json:"size"`
				UUID   string `json:"uuid"`
				Target string `json:"target"`
				Type  string `json:"type"`
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

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if len(files) == 0 {
				return errors.New("no files")
			}

			// Always cleanup the files in upload.
			var tempFiles []string
			for _, f := range files {
				if f.Type != SRS_SOURCE_TYPE_STREAM {
					tempFiles = append(tempFiles, f.Target)
				}
			}
			defer func() {
				for _, tempFile := range tempFiles {
					if _, err := os.Stat(tempFile); err == nil {
						os.Remove(tempFile)
						logger.Tf(ctx, "vLive cleanup %v", tempFile)
					}
				}
			}()

			// Check files.
			for _, f := range files {
				if f.Target == "" {
					return errors.New("no target")
				}
				if f.Type != SRS_SOURCE_TYPE_STREAM {
					if _, err := os.Stat(f.Target); err != nil {
						return errors.Wrapf(err, "no file %v", f.Target)
					}
					if !strings.HasPrefix(f.Target, dirUploadPath) {
						return errors.Errorf("invalid target %v", f.Target)
					}
				}
			}

			// Check platform.
			allowedPlatforms := []string{"wx", "bilibili", "kuaishou"}
			if platform == "" {
				return errors.New("no platform")
			}
			if !slicesContains(allowedPlatforms, platform) {
				return errors.Errorf("invalid platform %v", platform)
			}

			// Parsed source files.
			var parsedFiles []*VLiveSourceFile

			// Parse file information and move file from dirUploadPath to dirVLivePath.
			for _, file := range files {
				// Probe file information.
				stdout, err := exec.CommandContext(ctx, "ffprobe",
					"-show_error", "-show_private_data", "-v", "quiet", "-find_stream_info", "-print_format", "json",
					"-show_format", "-show_streams", file.Target,
				).Output()
				if err != nil {
					return errors.Wrapf(err, "probe %v", file.Target)
				}

				format := struct {
					Format VLiveFileFormat `json:"format"`
				}{}
				if err = json.Unmarshal([]byte(stdout), &format); err != nil {
					return errors.Wrapf(err, "parse format %v", stdout)
				}

				videos := struct {
					Streams []VLiveFileVideo `json:"streams"`
				}{}
				if err = json.Unmarshal([]byte(stdout), &videos); err != nil {
					return errors.Wrapf(err, "parse video streams %v", stdout)
				}
				var matchVideo *VLiveFileVideo
				for _, video := range videos.Streams {
					if video.CodecType == "video" {
						matchVideo = &video
						format.Format.HasVideo = true
						break
					}
				}

				audios := struct {
					Streams []VLiveFileAudio `json:"streams"`
				}{}
				if err = json.Unmarshal([]byte(stdout), &audios); err != nil {
					return errors.Wrapf(err, "parse audio streams %v", stdout)
				}
				var matchAudio *VLiveFileAudio
				for _, audio := range audios.Streams {
					if audio.CodecType == "audio" {
						matchAudio = &audio
						format.Format.HasAudio = true
						break
					}
				}

				parsedFile := &VLiveSourceFile{
					Name: file.Name, Size: uint64(file.Size), UUID: file.UUID,
					Target: file.Target,
					Type: file.Type,
					Format: &format.Format, Video: matchVideo, Audio: matchAudio,
				}
				if file.Type != SRS_SOURCE_TYPE_STREAM {
					parsedFile.Target = path.Join(dirVLivePath, fmt.Sprintf("%v%v", file.UUID, path.Ext(file.Target)))
					if err = os.Rename(file.Target, parsedFile.Target); err != nil {
						return errors.Wrapf(err, "rename %v to %v", file.Target, parsedFile.Target)
					}
				}

				parsedFiles = append(parsedFiles, parsedFile)
				logger.Tf(ctx, "vLive process file %v", parsedFile.String())
			}

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
				if f.Type != SRS_SOURCE_TYPE_STREAM {
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

			ohttp.WriteData(ctx, w, r, &struct {
				Platform string             `json:"platform"`
				Files    []*VLiveSourceFile `json:"files"`
			}{
				Platform: platform, Files: parsedFiles,
			})
			logger.Tf(ctx, "Update vLive ok, token=%vB", len(token))
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
	logger.Tf(ctx, "vLive start a worker")

	// Load tasks from redis and force to kill all.
	if objs, err := rdb.HGetAll(ctx, SRS_VLIVE_TASK).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hgetall %v", SRS_VLIVE_TASK)
	} else if len(objs) > 0 {
		for uuid, obj := range objs {
			logger.Tf(ctx, "Load task %v object %v", uuid, obj)

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
		configs, err := rdb.HGetAll(ctx, SRS_VLIVE_CONFIG).Result()
		if err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hgetall %v", SRS_VLIVE_CONFIG)
		}
		if len(configs) == 0 {
			return nil
		}

		for platform, config := range configs {
			var conf VLiveConfigure
			if err = json.Unmarshal([]byte(config), &conf); err != nil {
				return errors.Wrapf(err, "unmarshal %v %v", platform, config)
			}

			var task *VLiveTask
			if tv, loaded := v.tasks.LoadOrStore(conf.Platform, &VLiveTask{
				UUID:     uuid.NewString(),
				Platform: conf.Platform,
				config:   &conf,
			}); loaded {
				// Ignore if exists.
				continue
			} else {
				task = tv.(*VLiveTask)
				logger.Tf(ctx, "vLive create platform=%v task is %v", platform, task.String())
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
		logger.Tf(ctx, "vLive start to run tasks")

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

type VLiveFileFormat struct {
	Duration string `json:"duration"`
	Bitrate  string `json:"bit_rate"`
	Streams  int32  `json:"nb_streams"`
	Score    int32  `json:"probe_score"`
	HasVideo bool   `json:"has_video"`
	HasAudio bool   `json:"has_audio"`
}

func (v *VLiveFileFormat) String() string {
	return fmt.Sprintf("duration=%v, bitrate=%v, streams=%v, score=%v, video=%v, audio=%v",
		v.Duration, v.Bitrate, v.Streams, v.Score, v.HasVideo, v.HasAudio,
	)
}

type VLiveFileVideo struct {
	CodecType string `json:"codec_type"`
	CodecName string `json:"codec_name"`
	Profile   string `json:"profile"`
	Width     int32  `json:"width"`
	Height    int32  `json:"height"`
	PixFormat string `json:"pix_fmt"`
	Level     int32  `json:"level"`
	Bitrate   string `json:"bit_rate"`
}

func (v *VLiveFileVideo) String() string {
	return fmt.Sprintf("codec=%v, profile=%v, width=%v, height=%v, fmt=%v, level=%v, bitrate=%v",
		v.CodecName, v.Profile, v.Width, v.Height, v.PixFormat, v.Level, v.Bitrate,
	)
}

type VLiveFileAudio struct {
	CodecType     string `json:"codec_type"`
	CodecName     string `json:"codec_name"`
	Profile       string `json:"profile"`
	SampleFormat  string `json:"sample_fmt"`
	SampleRate    string `json:"sample_rate"`
	Channels      int32  `json:"channels"`
	ChannelLayout string `json:"channel_layout"`
	Bitrate       string `json:"bit_rate"`
}

func (v *VLiveFileAudio) String() string {
	return fmt.Sprintf("codec=%v, profile=%v, fmt=%v, rate=%v, channels=%v, layout=%v, bitrate=%v",
		v.CodecName, v.Profile, v.SampleFormat, v.SampleRate, v.Channels, v.ChannelLayout, v.Bitrate,
	)
}

type VLiveSourceFile struct {
	Name   string           `json:"name"`
	Size   uint64           `json:"size"`
	UUID   string           `json:"uuid"`
	Target string           `json:"target"`
	Type   string 		    `json:"type"`
	Format *VLiveFileFormat `json:"format"`
	Video  *VLiveFileVideo  `json:"video"`
	Audio  *VLiveFileAudio  `json:"audio"`
}

func (v *VLiveSourceFile) String() string {
	return fmt.Sprintf("name=%v, size=%v, uuid=%v, target=%v, format=(%v), video=(%v), audio=(%v)",
		v.Name, v.Size, v.UUID, v.Target, v.Format, v.Video, v.Audio,
	)
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
	Files []*VLiveSourceFile `json:"files"`
}

func (v *VLiveConfigure) String() string {
	return fmt.Sprintf("platform=%v, server=%v, secret=%v, enabled=%v, customed=%v, label=%v, files=%v",
		v.Platform, v.Server, v.Secret, v.Enabled, v.Customed, v.Label, v.Files,
	)
}

func (v *VLiveConfigure) Update(u *VLiveConfigure) error {
	if u.Platform != "" {
		v.Platform = u.Platform
	}
	if u.Server != "" {
		v.Server = u.Server
	}
	if u.Secret != "" {
		v.Secret = u.Secret
	}
	if u.Label != "" {
		v.Label = u.Label
	}
	v.Enabled = u.Enabled
	v.Customed = u.Customed
	v.Files = append([]*VLiveSourceFile{}, u.Files...)
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
	update string

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
	v.update = time.Now().Format(time.RFC3339)
}

func (v *VLiveTask) queryFrame() (int32, string, string, string) {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.PID, v.inputUUID, v.frame, v.update
}

func (v *VLiveTask) Initialize(ctx context.Context, w *VLiveWorker) error {
	v.vLiveWorker = w
	logger.Tf(ctx, "vLive initialize uuid=%v, platform=%v", v.UUID, v.Platform)

	if err := v.saveTask(ctx); err != nil {
		return errors.Wrapf(err, "save task")
	}

	return nil
}

func (v *VLiveTask) Run(ctx context.Context) error {
	ctx = logger.WithContext(ctx)
	logger.Tf(ctx, "vLive run task %v", v.String())

	selectInputFile := func() *VLiveSourceFile {
		v.lock.Lock()
		defer v.lock.Unlock()

		if len(v.config.Files) == 0 {
			return nil
		}

		file := v.config.Files[0]
		logger.Tf(ctx, "vLive use file=%v as input for platform=%v", file.UUID, v.Platform)
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
		if err := v.doVLive(ctx, input); err != nil {
			return errors.Wrapf(err, "do vLive")
		}

		return nil
	}

	for ctx.Err() == nil {
		if err := pfn(ctx); err != nil {
			logger.Wf(ctx, "ignore %v err %+v", v.String(), err)

			select {
			case <-ctx.Done():
			case <-time.After(10 * time.Second):
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

func (v *VLiveTask) doVLive(ctx context.Context, input *VLiveSourceFile) error {
	// Create context for current task.
	parentCtx := ctx
	ctx, v.cancel = context.WithCancel(ctx)

	// Build input URL.
	host := "localhost"

	// Build output URL.
	outputServer := strings.ReplaceAll(v.config.Server, "localhost", host)
	if !strings.HasSuffix(outputServer, "/") && !strings.HasPrefix(v.config.Secret, "/") {
		outputServer += "/"
	}
	outputURL := fmt.Sprintf("%v%v", outputServer, v.config.Secret)

	// Start FFmpeg process.
	args := []string{}
	if input.Type != SRS_SOURCE_TYPE_STREAM {
		args = append(args, "-stream_loop", "-1")
	}
	args = append(args, "-re", "-i", input.Target, "-c", "copy", "-f", "flv", outputURL)
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
	logger.Tf(ctx, "vLive start, platform=%v, input=%v, pid=%v", v.Platform, input.Target, v.PID)

	if err := v.saveTask(ctx); err != nil {
		return errors.Wrapf(err, "save task %v", v.String())
	}

	buf := make([]byte, 4096)
	for {
		nn, err := stderr.Read(buf)
		if err != nil || nn == 0 {
			break
		}

		line := string(buf[:nn])
		for strings.Contains(line, "= ") {
			line = strings.ReplaceAll(line, "= ", "=")
		}
		v.updateFrame(line)
	}

	err = cmd.Wait()
	logger.Tf(ctx, "vLive done, platform=%v, input=%v, pid=%v, err=%v",
		v.Platform, input.Target, v.PID, err,
	)

	return nil
}
