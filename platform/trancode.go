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
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"

	// From ossrs.
	"github.com/ossrs/go-oryx-lib/errors"
	ohttp "github.com/ossrs/go-oryx-lib/http"
	"github.com/ossrs/go-oryx-lib/logger"

	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var transcodeWorker *TranscodeWorker

type TranscodeWorker struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// The global transcode task, only support one transcode task.
	task *TranscodeTask
}

func NewTranscodeWorker() *TranscodeWorker {
	v := &TranscodeWorker{}
	v.task = NewTranscodeTask()
	v.task.transcodeWorker = v
	return v
}

func (v *TranscodeWorker) Handle(ctx context.Context, handler *http.ServeMux) error {
	ep := "/terraform/v1/ffmpeg/transcode/query"
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

			var config TranscodeConfig
			if b, err := rdb.HGet(ctx, SRS_TRANSCODE_CONFIG, "global").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v global", SRS_TRANSCODE_CONFIG)
			} else if len(b) > 0 {
				if err := json.Unmarshal([]byte(b), &config); err != nil {
					return errors.Wrapf(err, "unmarshal %v", b)
				}
			}

			ohttp.WriteData(ctx, w, r, &config)
			logger.Tf(ctx, "transcode query ok, %v, token=%vB", config, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ffmpeg/transcode/apply"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var config TranscodeConfig
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				*TranscodeConfig
			}{
				Token:           &token,
				TranscodeConfig: &config,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if b, err := json.Marshal(config); err != nil {
				return errors.Wrapf(err, "marshal conf %v", config)
			} else if err := rdb.HSet(ctx, SRS_TRANSCODE_CONFIG, "global", string(b)).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v global %v", SRS_TRANSCODE_CONFIG, string(b))
			}

			if err := v.task.Restart(ctx); err != nil {
				return errors.Wrapf(err, "restart task %v", config.String())
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "transcode apply ok, %v, token=%vB", config, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ffmpeg/transcode/task"
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

			var config TranscodeConfig
			if b, err := rdb.HGet(ctx, SRS_TRANSCODE_CONFIG, "global").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v global", SRS_TRANSCODE_CONFIG)
			} else if len(b) > 0 {
				if err := json.Unmarshal([]byte(b), &config); err != nil {
					return errors.Wrapf(err, "unmarshal %v", b)
				}
			}

			pid, input, output, frame, update := v.task.queryFrame()

			res := struct {
				// The task uuid.
				UUID string `json:"uuid"`
				// Whether task is enabled.
				Enabled bool `json:"enabled"`
				// The input stream URL.
				InputStream string `json:"input"`
				// The output stream URL.
				OutputStream string `json:"output"`
				// The FFmpeg log.
				Frame struct {
					// The FFmpeg log lines.
					Log string `json:"log"`
					// The last update time.
					Update string `json:"update"`
				} `json:"frame"`
			}{}
			res.Enabled = config.All
			res.UUID = v.task.UUID
			if pid > 0 {
				res.InputStream = input
				res.OutputStream = output
				res.Frame.Log = frame
				res.Frame.Update = update
			}

			ohttp.WriteData(ctx, w, r, &res)
			logger.Tf(ctx, "transcode task ok, %v, pid=%v, input=%v, output=%v, frame=%v, update=%v, token=%vB",
				config, pid, input, output, frame, update, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	return nil
}

func (v *TranscodeWorker) Close() error {
	if v.cancel != nil {
		v.cancel()
	}
	v.wg.Wait()
	return nil
}

func (v *TranscodeWorker) Start(ctx context.Context) error {
	wg := &v.wg

	ctx, cancel := context.WithCancel(ctx)
	v.cancel = cancel

	ctx = logger.WithContext(ctx)
	logger.Tf(ctx, "transcode start a worker")

	// Load tasks from redis and force to kill all.
	if objs, err := rdb.HGetAll(ctx, SRS_TRANSCODE_TASK).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hgetall %v", SRS_TRANSCODE_TASK)
	} else if len(objs) > 0 {
		for uuid, obj := range objs {
			logger.Tf(ctx, "Load task %v object %v", uuid, obj)

			var task TranscodeTask
			if err = json.Unmarshal([]byte(obj), &task); err != nil {
				return errors.Wrapf(err, "unmarshal %v %v", uuid, obj)
			}

			if task.PID > 0 {
				task.cleanup(ctx)
			}
		}

		if err = rdb.Del(ctx, SRS_TRANSCODE_TASK).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "del %v", SRS_TRANSCODE_TASK)
		}
	}

	// Start global transcode task.
	wg.Add(1)
	go func() {
		defer wg.Done()

		task := v.task
		for ctx.Err() == nil {
			var duration time.Duration
			if err := task.Run(ctx); err != nil {
				logger.Wf(ctx, "run task %v err %+v", task.String(), err)
				duration = 10 * time.Second
			} else {
				duration = 3 * time.Second
			}

			select {
			case <-ctx.Done():
			case <-time.After(duration):
			}
		}
	}()

	return nil
}

type TranscodeConfig struct {
	// Whether transcode all streams.
	All bool `json:"all"`
	// The video codec name.
	VideoCodec string `json:"vcodec"`
	// The audio codec name.
	AudioCodec string `json:"acodec"`
	// The video bitrate in kbps.
	VideoBitrate int `json:"vbitrate"`
	// The audio bitrate in kbps.
	AudioBitrate int `json:"abitrate"`
	// The video profile, for example, baseline.
	VideoProfile string `json:"vprofile"`
	// The video preset, for example, veryfast.
	VideoPreset string `json:"vpreset"`
	// The audio channels.
	AudioChannels int `json:"achannels"`
	// The RTMP server url, for example, rtmp://localhost/live
	Server string `json:"server"`
	// The RTMP stream and secret, for example, livestream
	Secret string `json:"secret"`
}

func (v TranscodeConfig) String() string {
	return fmt.Sprintf("all=%v, vcodec=%v, acodec=%v, vbitrate=%v, abitrate=%v, achannels=%v, vprofile=%v, vpreset=%v, server=%v, secret=%v",
		v.All, v.VideoCodec, v.AudioCodec, v.VideoBitrate, v.AudioBitrate, v.AudioChannels, v.VideoProfile,
		v.VideoPreset, v.Server, v.Secret,
	)
}

type TranscodeTask struct {
	// The ID for task.
	UUID string `json:"uuid"`

	// The input url.
	Input string `json:"input"`
	// The input stream URL.
	inputStreamURL string
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

	// The configure for transcode task.
	config TranscodeConfig
	// The transcode worker.
	transcodeWorker *TranscodeWorker

	// To protect the fields.
	lock sync.Mutex
}

func NewTranscodeTask() *TranscodeTask {
	return &TranscodeTask{UUID: uuid.NewString()}
}

func (v *TranscodeTask) String() string {
	return fmt.Sprintf("uuid=%v, pid=%v, config is %v",
		v.UUID, v.PID, v.config.String(),
	)
}

func (v *TranscodeTask) Restart(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if v.cancel != nil {
		v.cancel()
	}

	return nil
}

func (v *TranscodeTask) Run(ctx context.Context) error {
	ctx = logger.WithContext(ctx)
	logger.Tf(ctx, "transcode run task %v", v.String())

	isSameStream := func(a, b string) bool {
		ua, err := url.Parse(a)
		if err != nil {
			return false
		}

		ub, err := url.Parse(b)
		if err != nil {
			return false
		}

		if path.Clean(ua.Path) == path.Clean(ub.Path) {
			return true
		}

		return false
	}

	// TODO: FIXME: Should select stream again when stream republished.
	selectActiveStream := func() (*SrsStream, error) {
		streams, err := rdb.HGetAll(ctx, SRS_STREAM_ACTIVE).Result()
		if err != nil {
			return nil, errors.Wrapf(err, "hgetall %v", SRS_STREAM_ACTIVE)
		}

		var best *SrsStream
		for _, value := range streams {
			var stream SrsStream
			if err := json.Unmarshal([]byte(value), &stream); err != nil {
				return nil, errors.Wrapf(err, "unmarshal %v", value)
			}

			// Ignore the transcode stream itself.
			if isSameStream(fmt.Sprintf("%v/%v", v.config.Server, v.config.Secret),
				fmt.Sprintf("rtmp://%v/%v/%v", stream.Vhost, stream.App, stream.Stream)) {
				continue
			}

			if best == nil {
				best = &stream
				continue
			}

			bestUpdate, err := time.Parse(time.RFC3339, best.Update)
			if err != nil {
				return nil, errors.Wrapf(err, "parse %v", best.Update)
			}

			streamUpdate, err := time.Parse(time.RFC3339, stream.Update)
			if err != nil {
				return nil, errors.Wrapf(err, "parse %v", stream.Update)
			}

			if bestUpdate.Before(streamUpdate) {
				best = &stream
			}
		}

		// Ignore if no active stream.
		if best == nil {
			return nil, nil
		}

		logger.Tf(ctx, "transcode use best=%v as input", best.StreamURL())
		return best, nil
	}

	pfn := func(ctx context.Context) error {
		if b, err := rdb.HGet(ctx, SRS_TRANSCODE_CONFIG, "global").Result(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hget %v global", SRS_TRANSCODE_CONFIG)
		} else if len(b) > 0 {
			if err := json.Unmarshal([]byte(b), &v.config); err != nil {
				return errors.Wrapf(err, "unmarshal %v", b)
			}
		}

		// Ignore if not enabled.
		if !v.config.All {
			return nil
		}

		if err := v.saveTask(ctx); err != nil {
			return errors.Wrapf(err, "save task")
		}

		// Use a active stream as input.
		input, err := selectActiveStream()
		if err != nil {
			return errors.Wrapf(err, "select input")
		}

		if input == nil {
			return nil
		}

		// Start transcode task.
		if err := v.doTranscode(ctx, input); err != nil {
			return errors.Wrapf(err, "do transcode")
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

func (v *TranscodeTask) doTranscode(ctx context.Context, input *SrsStream) error {
	// Create context for current task.
	parentCtx := ctx
	ctx, cancel := context.WithCancel(ctx)
	v.cancel = cancel

	// Build input URL.
	host := "localhost"
	inputURL := fmt.Sprintf("rtmp://%v/%v/%v", host, input.App, input.Stream)

	// Build output URL.
	outputServer := strings.ReplaceAll(v.config.Server, "localhost", host)
	if !strings.HasSuffix(outputServer, "/") && !strings.HasPrefix(v.config.Secret, "/") {
		outputServer += "/"
	}
	outputURL := fmt.Sprintf("%v%v", outputServer, v.config.Secret)

	// Start FFmpeg process.
	args := []string{
		"-stream_loop", "-1", "-i", inputURL,
		"-vcodec", v.config.VideoCodec,
		"-profile:v", v.config.VideoProfile,
		"-preset:v", v.config.VideoPreset,
		"-tune", "zerolatency", // Low latency mode.
		"-b:v", fmt.Sprintf("%vk", v.config.VideoBitrate),
		"-r", "25", "-g", "50", // Set gop to 2s.
		"-bf", "0", // Disable B frame for WebRTC.
		"-acodec", v.config.AudioCodec,
		"-b:a", fmt.Sprintf("%vk", v.config.AudioBitrate),
	}
	if v.config.AudioChannels > 0 {
		args = append(args, "-ac", fmt.Sprintf("%v", v.config.AudioChannels))
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
	v.Input, v.inputStreamURL, v.Output = inputURL, input.StreamURL(), outputURL
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
	logger.Tf(ctx, "transcode start, stream=%v, pid=%v", input.StreamURL(), v.PID)

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
				logger.Tf(ctx, "transcode FFmpeg not update for a while, restart it")
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
	logger.Tf(ctx, "transcode done, stream=%v, pid=%v, err=%v",
		input.StreamURL(), v.PID, err,
	)
	return nil
}

func (v *TranscodeTask) updateFrame(frame string) {
	v.lock.Lock()
	defer v.lock.Unlock()

	v.frame = strings.TrimSpace(frame)
	v.update = time.Now()
}

func (v *TranscodeTask) queryFrame() (int32, string, string, string, string) {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.PID, v.inputStreamURL, v.Output, v.frame, v.update.Format(time.RFC3339)
}

func (v *TranscodeTask) saveTask(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if b, err := json.Marshal(v); err != nil {
		return errors.Wrapf(err, "marshal %v", v.String())
	} else if err = rdb.HSet(ctx, SRS_TRANSCODE_TASK, v.UUID, string(b)).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_TRANSCODE_TASK, v.UUID, string(b))
	}

	return nil
}

func (v *TranscodeTask) cleanup(ctx context.Context) error {
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
