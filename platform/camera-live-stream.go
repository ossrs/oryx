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
	"net/http"
	"os"
	"os/exec"
	"path"
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

var cameraWorker *CameraWorker

type CameraWorker struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// The tasks we have started to IP camera streams,, key is platform in string, value is *CameraTask.
	tasks sync.Map
}

func NewCameraWorker() *CameraWorker {
	return &CameraWorker{}
}

func (v *CameraWorker) GetTask(platform string) *CameraTask {
	if task, loaded := v.tasks.Load(platform); loaded {
		return task.(*CameraTask)
	}
	return nil
}

func (v *CameraWorker) Handle(ctx context.Context, handler *http.ServeMux) error {
	ep := "/terraform/v1/ffmpeg/camera/secret"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token, action string
			var userConf CameraConfigure
			if err := ParseBody(ctx, r.Body, &struct {
				Token  *string `json:"token"`
				Action *string `json:"action"`
				*CameraConfigure
			}{
				Token: &token, Action: &action, CameraConfigure: &userConf,
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
				if len(userConf.Streams) == 0 {
					return errors.New("no files")
				}
			}

			if action == "update" {
				var targetConf CameraConfigure
				if config, err := rdb.HGet(ctx, SRS_CAMERA_CONFIG, userConf.Platform).Result(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hget %v %v", SRS_CAMERA_CONFIG, userConf.Platform)
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
					} else if err = rdb.HSet(ctx, SRS_CAMERA_CONFIG, userConf.Platform, string(newB)).Err(); err != nil && err != redis.Nil {
						return errors.Wrapf(err, "hset %v %v %v", SRS_CAMERA_CONFIG, userConf.Platform, string(newB))
					}
				}

				// Restart the IP camera if exists.
				if task := cameraWorker.GetTask(userConf.Platform); task != nil {
					if err := task.Restart(ctx); err != nil {
						return errors.Wrapf(err, "restart task %v", userConf.String())
					}
				}

				ohttp.WriteData(ctx, w, r, nil)
				logger.Tf(ctx, "Camera: update secret ok, token=%vB", len(token))
				return nil
			} else {
				confObjs := make(map[string]*CameraConfigure)
				if configs, err := rdb.HGetAll(ctx, SRS_CAMERA_CONFIG).Result(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hgetall %v", SRS_CAMERA_CONFIG)
				} else {
					for k, v := range configs {
						var obj CameraConfigure
						if err = json.Unmarshal([]byte(v), &obj); err != nil {
							return errors.Wrapf(err, "unmarshal %v %v", k, v)
						}
						confObjs[k] = &obj
					}
				}

				ohttp.WriteData(ctx, w, r, confObjs)
				logger.Tf(ctx, "Camera: query configures ok, token=%vB", len(token))
				return nil
			}
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ffmpeg/camera/streams"
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
			if configs, err := rdb.HGetAll(ctx, SRS_CAMERA_CONFIG).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hgetall %v", SRS_CAMERA_CONFIG)
			} else if len(configs) > 0 {
				for k, v := range configs {
					var config CameraConfigure
					if err = json.Unmarshal([]byte(v), &config); err != nil {
						return errors.Wrapf(err, "unmarshal %v %v", k, v)
					}

					var pid int32
					var inputUUID, frame, update string
					if task := cameraWorker.GetTask(config.Platform); task != nil {
						pid, inputUUID, frame, update = task.queryFrame()
					}

					elem := map[string]interface{}{
						"platform":   config.Platform,
						"enabled":    config.Enabled,
						"custom":     config.Customed,
						"label":      config.Label,
						"files":      config.Streams,
						"extraAudio": config.ExtraAudio,
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
			logger.Tf(ctx, "Camera: Query streams ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ffmpeg/camera/stream-url"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
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

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			// Parse URL to object.
			u, err := RebuildStreamURL(qUrl)
			if err != nil {
				return errors.Wrapf(err, "parse %v", qUrl)
			}

			// Check url if valid rtmp or rtsp or http-flv or https-flv or hls live url
			if u.Scheme != "rtmp" && u.Scheme != "rtsp" && u.Scheme != "http" && u.Scheme != "https" {
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

			logger.Tf(ctx, "Camera: Update stream url ok, url=%v, uuid=%v", qUrl, targetUUID)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ffmpeg/camera/source"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			type CameraTempFile struct {
				// The file name.
				Name string `json:"name"`
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
			var streams []*CameraTempFile
			if err := ParseBody(ctx, r.Body, &struct {
				Token    *string            `json:"token"`
				Platform *string            `json:"platform"`
				Streams  *[]*CameraTempFile `json:"files"`
			}{
				Token: &token, Platform: &platform, Streams: &streams,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if len(streams) == 0 {
				return errors.New("no files")
			}

			// Check files.
			for _, stream := range streams {
				if stream.Target == "" {
					return errors.New("no target")
				}
				if stream.Type != FFprobeSourceTypeStream {
					return errors.Errorf("invalid target %v type %v", stream.Target, stream.Type)
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
			var parsedStreams []*FFprobeSource

			// Parse file information and move file from camera stream URL.
			for _, stream := range streams {
				// Probe file information.
				toCtx, toCancelFunc := context.WithTimeout(ctx, 15*time.Second)
				defer toCancelFunc()

				args := []string{
					"-show_error", "-show_private_data", "-v", "quiet", "-find_stream_info", "-print_format", "json",
					"-show_format", "-show_streams",
				}
				// For RTSP stream source, always use TCP transport.
				if strings.HasPrefix(stream.Target, "rtsp://") {
					args = append(args, "-rtsp_transport", "tcp")
				}
				// Rebuild the stream url, because it may contain special characters.
				if strings.Contains(stream.Target, "://") {
					if u, err := RebuildStreamURL(stream.Target); err != nil {
						return errors.Wrapf(err, "rebuild %v", stream.Target)
					} else {
						args = append(args, "-i", u.String())
					}
				} else {
					args = append(args, "-i", stream.Target)
				}

				stdout, err := exec.CommandContext(toCtx, "ffprobe", args...).Output()
				if err != nil {
					return errors.Wrapf(err, "probe %v with ffprobe %v", stream.Target, args)
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
					if limits, err := rdb.HGet(ctx, SRS_SYS_LIMITS, "camera").Int64(); err != nil && err != redis.Nil {
						return errors.Wrapf(err, "hget %v camera", SRS_SYS_LIMITS)
					} else {
						if limits == 0 {
							limits = SrsSysLimitsCamera // in Kbps.
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

				parsedStream := &FFprobeSource{
					Name: stream.Name, Size: uint64(stream.Size), UUID: stream.UUID,
					Target: stream.Target,
					Type:   stream.Type,
					Format: &format.Format, Video: matchVideo, Audio: matchAudio,
				}

				parsedStreams = append(parsedStreams, parsedStream)
				logger.Tf(ctx, "Camera: process stream %v", parsedStream.String())
			}

			// Update redis object.
			confObj := CameraConfigure{Platform: platform}
			if conf, err := rdb.HGet(ctx, SRS_CAMERA_CONFIG, platform).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_CAMERA_CONFIG, platform)
			} else if conf != "" {
				if err = json.Unmarshal([]byte(conf), &confObj); err != nil {
					return errors.Wrapf(err, "parse %v", conf)
				}
			}

			confObj.Streams = parsedStreams

			if b, err := json.Marshal(&confObj); err != nil {
				return errors.Wrapf(err, "marshal %v", confObj.String())
			} else if err = rdb.HSet(ctx, SRS_CAMERA_CONFIG, platform, string(b)).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v %v %v", SRS_CAMERA_CONFIG, platform, string(b))
			}

			// Restart the IP camera if exists.
			if task := cameraWorker.GetTask(platform); task != nil {
				if err := task.Restart(ctx); err != nil {
					return errors.Wrapf(err, "restart task %v", platform)
				}
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Platform string           `json:"platform"`
				Files    []*FFprobeSource `json:"files"`
			}{
				Platform: platform, Files: parsedStreams,
			})
			logger.Tf(ctx, "Camera:: Update ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	return nil
}

func (v *CameraWorker) Close() error {
	if v.cancel != nil {
		v.cancel()
	}
	v.wg.Wait()
	return nil
}

func (v *CameraWorker) Start(ctx context.Context) error {
	wg := &v.wg

	ctx, cancel := context.WithCancel(ctx)
	v.cancel = cancel

	ctx = logger.WithContext(ctx)
	logger.Tf(ctx, "Camera: start a worker")

	// Load tasks from redis and force to kill all.
	if objs, err := rdb.HGetAll(ctx, SRS_CAMERA_TASK).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hgetall %v", SRS_CAMERA_TASK)
	} else if len(objs) > 0 {
		for uuid, obj := range objs {
			logger.Tf(ctx, "Load task %v object %v", uuid, obj)

			var task CameraTask
			if err = json.Unmarshal([]byte(obj), &task); err != nil {
				return errors.Wrapf(err, "unmarshal %v %v", uuid, obj)
			}

			if task.PID > 0 {
				task.cleanup(ctx)
			}
		}

		if err = rdb.Del(ctx, SRS_CAMERA_TASK).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "del %v", SRS_CAMERA_TASK)
		}
	}

	// Load all configurations from redis.
	loadTasks := func() error {
		configItems, err := rdb.HGetAll(ctx, SRS_CAMERA_CONFIG).Result()
		if err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hgetall %v", SRS_CAMERA_CONFIG)
		}
		if len(configItems) == 0 {
			return nil
		}

		for platform, configItem := range configItems {
			var config CameraConfigure
			if err = json.Unmarshal([]byte(configItem), &config); err != nil {
				return errors.Wrapf(err, "unmarshal %v %v", platform, configItem)
			}

			var task *CameraTask
			if tv, loaded := v.tasks.LoadOrStore(config.Platform, &CameraTask{
				UUID:     uuid.NewString(),
				Platform: config.Platform,
				config:   &config,
			}); loaded {
				// Ignore if exists.
				continue
			} else {
				task = tv.(*CameraTask)
				logger.Tf(ctx, "Camera: create platform=%v task is %v", platform, task.String())
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
		logger.Tf(ctx, "Camera: Start to run tasks")

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

// CameraConfigure is the configure for IP camera.
type CameraConfigure struct {
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
	// The extra audio stream strategy.
	ExtraAudio string `json:"extraAudio"`

	// The input files for IP camera.
	Streams []*FFprobeSource `json:"files"`
}

func (v CameraConfigure) String() string {
	return fmt.Sprintf("platform=%v, server=%v, secret=%v, enabled=%v, customed=%v, label=%v, files=%v, extraAudio=%v",
		v.Platform, v.Server, v.Secret, v.Enabled, v.Customed, v.Label, v.Streams, v.ExtraAudio,
	)
}

func (v *CameraConfigure) Update(u *CameraConfigure) error {
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
	v.Streams = append([]*FFprobeSource{}, u.Streams...)
	v.ExtraAudio = u.ExtraAudio
	return nil
}

// CameraTask is a task for FFmpeg to IP camera stream, with a configure.
type CameraTask struct {
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
	update time.Time

	// The context for current task.
	cancel context.CancelFunc

	// The configure for IP camera task.
	config *CameraConfigure
	// The IP camera worker.
	cameraWorker *CameraWorker

	// To protect the fields.
	lock sync.Mutex
}

func (v *CameraTask) String() string {
	return fmt.Sprintf("uuid=%v, platform=%v, input=%v, output=%v, pid=%v, frame=%vB, config is %v",
		v.UUID, v.Platform, v.Input, v.Output, v.PID, len(v.frame), v.config.String(),
	)
}

func (v *CameraTask) saveTask(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if b, err := json.Marshal(v); err != nil {
		return errors.Wrapf(err, "marshal %v", v.String())
	} else if err = rdb.HSet(ctx, SRS_CAMERA_TASK, v.UUID, string(b)).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_CAMERA_TASK, v.UUID, string(b))
	}

	return nil
}

func (v *CameraTask) cleanup(ctx context.Context) error {
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

func (v *CameraTask) Restart(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if v.cancel != nil {
		v.cancel()
	}

	// Reload config from redis.
	if b, err := rdb.HGet(ctx, SRS_CAMERA_CONFIG, v.Platform).Result(); err != nil {
		return errors.Wrapf(err, "hget %v %v", SRS_CAMERA_CONFIG, v.Platform)
	} else if err = json.Unmarshal([]byte(b), v.config); err != nil {
		return errors.Wrapf(err, "unmarshal %v", b)
	}

	return nil
}

func (v *CameraTask) updateFrame(frame string) {
	v.lock.Lock()
	defer v.lock.Unlock()

	v.frame = frame
	v.update = time.Now()
}

func (v *CameraTask) queryFrame() (int32, string, string, string) {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.PID, v.inputUUID, v.frame, v.update.Format(time.RFC3339)
}

func (v *CameraTask) Initialize(ctx context.Context, w *CameraWorker) error {
	v.cameraWorker = w
	logger.Tf(ctx, "Camera: Initialize uuid=%v, platform=%v", v.UUID, v.Platform)

	if err := v.saveTask(ctx); err != nil {
		return errors.Wrapf(err, "save task")
	}

	return nil
}

func (v *CameraTask) Run(ctx context.Context) error {
	ctx = logger.WithContext(ctx)
	logger.Tf(ctx, "Camera: Run task %v", v.String())

	selectInputFile := func() *FFprobeSource {
		v.lock.Lock()
		defer v.lock.Unlock()

		if len(v.config.Streams) == 0 {
			return nil
		}

		file := v.config.Streams[0]
		logger.Tf(ctx, "Camera: Use file=%v as input for platform=%v", file.UUID, v.Platform)
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

		// Start IP camera task.
		if err := v.doCameraStreaming(ctx, input); err != nil {
			return errors.Wrapf(err, "do IP camera")
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

func (v *CameraTask) doCameraStreaming(ctx context.Context, input *FFprobeSource) error {
	// Create context for current task.
	parentCtx := ctx
	ctx, cancel := context.WithCancel(ctx)
	v.cancel = cancel

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
	if input.Type != FFprobeSourceTypeStream {
		args = append(args, "-stream_loop", "-1")
	}
	args = append(args, "-re",
		"-fflags", "nobuffer", // Reduce the latency introduced by optional buffering.
	)
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
		}
	} else {
		args = append(args, "-i", input.Target)
	}
	// Whether insert extra audio stream.
	if v.config.ExtraAudio == "silent" {
		args = append(args,
			"-f", "lavfi", "-i", "anullsrc=channel_layout=stereo:sample_rate=44100", // Silent audio stream.
			"-map", "0:v", "-map", "1:a", // Ignore the original audio stream.
			"-c:a", "aac", "-ac", "2", "-ar", "44100", "-b:a", "20k", // Encode audio stream.
			"-c:v", "copy", // Copy video stream.
		)
	} else {
		args = append(args, "-c", "copy")
	}
	args = append(args, "-f", "flv", outputURL)
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
	logger.Tf(ctx, "Camera: Start, platform=%v, input=%v, pid=%v", v.Platform, input.Target, v.PID)

	if err := v.saveTask(ctx); err != nil {
		return errors.Wrapf(err, "save task %v", v.String())
	}

	// Monitor FFmpeg update, restart if not update for a while.
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(3 * time.Second):
			}

			if v.update.Add(10 * time.Second).Before(time.Now()) {
				logger.Tf(ctx, "Camera: FFmpeg not update for a while, restart it")
				cancel()
				return
			}
		}
	}()

	// Read stderr to update status and output of FFmpeg.
	pollingCtx, pollingCancel := context.WithCancel(ctx)
	go func() {
		defer pollingCancel()
		buf := make([]byte, 4096)
		for ctx.Err() == nil {
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
	}()

	// Process terminated, or user cancel the process.
	select {
	case <-parentCtx.Done():
	case <-ctx.Done():
	case <-pollingCtx.Done():
	}

	err = cmd.Wait()
	logger.Tf(ctx, "Camera: Cycle done, platform=%v, input=%v, pid=%v, err=%v",
		v.Platform, input.Target, v.PID, err,
	)

	return nil
}
