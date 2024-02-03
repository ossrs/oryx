package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
)

func TestScenario_WithStream_PublishCameraStreamUrl(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	if *noMediaTest {
		return
	}

	var r0, r1, r2, r3, r4, r5 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1, r2, r3, r4, r5); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	var pubSecret string
	if err := NewApi().WithAuth(ctx, "/terraform/v1/hooks/srs/secret/query", nil, &struct {
		Publish *string `json:"publish"`
	}{
		Publish: &pubSecret,
	}); err != nil {
		r0 = err
		return
	}

	var wg sync.WaitGroup
	defer wg.Wait()
	// Start FFmpeg to publish stream.
	streamID := fmt.Sprintf("stream-%v-%v", os.Getpid(), rand.Int())
	streamURL := fmt.Sprintf("%v/live/%v?secret=%v", *endpointRTMP, streamID, pubSecret)
	ffmpeg := NewFFmpeg(func(v *ffmpegClient) {
		v.args = []string{
			"-re", "-stream_loop", "-1", "-i", *srsInputFile, "-c", "copy",
			"-f", "flv", streamURL,
		}
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r1 = ffmpeg.Run(ctx, cancel)
	}()

	defer cancel()
	select {
	case <-ctx.Done():
		return
	case <-ffmpeg.ReadyCtx().Done():
	}

	// Use the publish stream url as the camera input.
	res := struct {
		Name   string `json:"name"`
		Target string `json:"target"`
		UUID   string `json:"uuid"`
		Size   int64  `json:"size"`
		Type   string `json:"type"`
	}{}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/camera/stream-url", &struct {
		StreamURL string `json:"url"`
	}{
		StreamURL: streamURL,
	}, &res); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg camera streamUrl failed")
		return
	}

	// Use the publish stream url as the camera source.
	res.Size = 0
	res.Type = "stream"
	codec := struct {
		UUID  string `json:"uuid"`
		Audio struct {
			CodecName  string `json:"codec_name"`
			Channels   int    `json:"channels"`
			SampleRate string `json:"sample_rate"`
		} `json:"audio"`
		Video struct {
			CodecName string `json:"codec_name"`
			Profile   string `json:"profile"`
			Width     int    `json:"width"`
			Height    int    `json:"height"`
		} `json:"video"`
	}{}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/camera/source", &struct {
		Platform string        `json:"platform"`
		Files    []interface{} `json:"files"`
	}{
		Platform: "bilibili",
		Files:    []interface{}{res},
	}, &struct {
		Files []interface{} `json:"files"`
	}{
		Files: []interface{}{&codec},
	}); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg camera source failed")
		return
	}

	if err := func() error {
		if codec.UUID != res.UUID {
			return errors.Errorf("invalid codec uuid=%v, %v", codec.UUID, res.UUID)
		}
		if codec.Audio.CodecName != "aac" || codec.Audio.Channels != 2 || codec.Audio.SampleRate != "44100" {
			return errors.Errorf("invalid codec audio=%v", codec.Audio)
		}
		if codec.Video.CodecName != "h264" || codec.Video.Profile != "High" || codec.Video.Width != 768 || codec.Video.Height != 320 {
			return errors.Errorf("invalid codec video=%v", codec.Video)
		}
		return nil
	}(); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg camera source failed")
		return
	}

	// Start camera streaming.
	type VLiveConfig struct {
		Platform string      `json:"platform"`
		Server   string      `json:"server"`
		Secret   string      `json:"secret"`
		Enabled  bool        `json:"enabled"`
		Custom   bool        `json:"custom"`
		Label    string      `json:"label"`
		Files    interface{} `json:"files"`
		Action   string      `json:"action"`
	}
	conf := make(map[string]*VLiveConfig)
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/camera/secret", nil, &conf); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg camera secret failed")
		return
	}

	bilibili, ok := conf["bilibili"]
	if !ok || bilibili == nil {
		r0 = errors.Errorf("invalid bilibili secret")
		return
	}

	// Restore the state of enabled.
	backup := *bilibili
	defer func() {
		backup.Action = "update"
		logger.Tf(ctx, "restore config %v", backup)

		if backup.Server == "" {
			backup.Server = bilibili.Server
		}
		if backup.Secret == "" {
			backup.Secret = bilibili.Secret
		}

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/camera/secret", backup, nil)
	}()

	publishStreamID := fmt.Sprintf("publish-stream-%v-%v", os.Getpid(), rand.Int())
	bilibili.Secret = fmt.Sprintf("%v?secret=%v", publishStreamID, pubSecret)
	bilibili.Server = "rtmp://localhost/live/"
	bilibili.Enabled = true
	bilibili.Action = "update"
	bilibili.Custom = true
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/camera/secret", &bilibili, nil); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg camera secret failed")
		return
	}

	// Start FFprobe to detect and verify stream.
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", publishStreamID)
		v.streamURL = fmt.Sprintf("%v/live/%v.flv", *endpointHTTP, publishStreamID)
		v.duration, v.timeout = duration, time.Duration(*srsFFprobeTimeout)*time.Millisecond
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r2 = ffprobe.Run(ctx, cancel)
	}()

	// Fast quit for probe done.
	select {
	case <-ctx.Done():
	case <-ffprobe.ProbeDoneCtx().Done():
		cancel()
	}

	str, m := ffprobe.Result()
	if len(m.Streams) != 2 {
		r3 = errors.Errorf("invalid streams=%v, %v, %v", len(m.Streams), m.String(), str)
	}

	if ts := 90; m.Format.ProbeScore < ts {
		r4 = errors.Errorf("low score=%v < %v, %v, %v", m.Format.ProbeScore, ts, m.String(), str)
	}
	if dv := m.Duration(); dv < duration/2 {
		r5 = errors.Errorf("short duration=%v < %v, %v, %v", dv, duration, m.String(), str)
	}
}

func TestScenario_WithStream_PublishCameraVideoOnly(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	if *noMediaTest {
		return
	}

	var r0, r1, r2, r3, r4, r5 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1, r2, r3, r4, r5); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	var pubSecret string
	if err := NewApi().WithAuth(ctx, "/terraform/v1/hooks/srs/secret/query", nil, &struct {
		Publish *string `json:"publish"`
	}{
		Publish: &pubSecret,
	}); err != nil {
		r0 = err
		return
	}

	var wg sync.WaitGroup
	defer wg.Wait()
	// Start FFmpeg to publish stream, in video only mode, drop audio.
	streamID := fmt.Sprintf("stream-%v-%v", os.Getpid(), rand.Int())
	streamURL := fmt.Sprintf("%v/live/%v?secret=%v", *endpointRTMP, streamID, pubSecret)
	ffmpeg := NewFFmpeg(func(v *ffmpegClient) {
		v.args = []string{
			"-re", "-stream_loop", "-1", "-i", *srsInputFile, "-c:v", "copy", "-an",
			"-f", "flv", streamURL,
		}
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r1 = ffmpeg.Run(ctx, cancel)
	}()

	defer cancel()
	select {
	case <-ctx.Done():
		return
	case <-ffmpeg.ReadyCtx().Done():
	}

	// Use the publish stream url as the camera input.
	res := struct {
		Name   string `json:"name"`
		Target string `json:"target"`
		UUID   string `json:"uuid"`
		Size   int64  `json:"size"`
		Type   string `json:"type"`
	}{}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/camera/stream-url", &struct {
		StreamURL string `json:"url"`
	}{
		StreamURL: streamURL,
	}, &res); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg camera streamUrl failed")
		return
	}

	// Use the publish stream url as the camera source.
	res.Size = 0
	res.Type = "stream"
	codec := struct {
		UUID  string `json:"uuid"`
		Video struct {
			CodecName string `json:"codec_name"`
			Profile   string `json:"profile"`
			Width     int    `json:"width"`
			Height    int    `json:"height"`
		} `json:"video"`
	}{}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/camera/source", &struct {
		Platform string        `json:"platform"`
		Files    []interface{} `json:"files"`
	}{
		Platform: "bilibili",
		Files:    []interface{}{res},
	}, &struct {
		Files []interface{} `json:"files"`
	}{
		Files: []interface{}{&codec},
	}); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg camera source failed")
		return
	}

	if err := func() error {
		if codec.UUID != res.UUID {
			return errors.Errorf("invalid codec uuid=%v, %v", codec.UUID, res.UUID)
		}
		if codec.Video.CodecName != "h264" || codec.Video.Profile != "High" || codec.Video.Width != 768 || codec.Video.Height != 320 {
			return errors.Errorf("invalid codec video=%v", codec.Video)
		}
		return nil
	}(); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg camera source failed")
		return
	}

	// Start camera streaming.
	type VLiveConfig struct {
		Platform   string      `json:"platform"`
		Server     string      `json:"server"`
		Secret     string      `json:"secret"`
		Enabled    bool        `json:"enabled"`
		Custom     bool        `json:"custom"`
		Label      string      `json:"label"`
		Files      interface{} `json:"files"`
		ExtraAudio string      `json:"extraAudio"`
		Action     string      `json:"action"`
	}
	conf := make(map[string]*VLiveConfig)
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/camera/secret", nil, &conf); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg camera secret failed")
		return
	}

	bilibili, ok := conf["bilibili"]
	if !ok || bilibili == nil {
		r0 = errors.Errorf("invalid bilibili secret")
		return
	}

	// Restore the state of enabled.
	backup := *bilibili
	defer func() {
		backup.Action = "update"
		logger.Tf(ctx, "restore config %v", backup)

		if backup.Server == "" {
			backup.Server = bilibili.Server
		}
		if backup.Secret == "" {
			backup.Secret = bilibili.Secret
		}

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/camera/secret", backup, nil)
	}()

	// Create extra silent audio stream.
	bilibili.ExtraAudio = "silent"

	publishStreamID := fmt.Sprintf("publish-stream-%v-%v", os.Getpid(), rand.Int())
	bilibili.Secret = fmt.Sprintf("%v?secret=%v", publishStreamID, pubSecret)
	bilibili.Server = "rtmp://localhost/live/"
	bilibili.Enabled = true
	bilibili.Action = "update"
	bilibili.Custom = true
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/camera/secret", &bilibili, nil); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg camera secret failed")
		return
	}

	// Start FFprobe to detect and verify stream.
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", publishStreamID)
		v.streamURL = fmt.Sprintf("%v/live/%v.flv", *endpointHTTP, publishStreamID)
		v.duration, v.timeout = duration, time.Duration(*srsFFprobeTimeout)*time.Millisecond
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r2 = ffprobe.Run(ctx, cancel)
	}()

	// Fast quit for probe done.
	select {
	case <-ctx.Done():
	case <-ffprobe.ProbeDoneCtx().Done():
		cancel()
	}

	str, m := ffprobe.Result()

	// Note that for camera streaming, should automatically add audio stream,
	// even if the source stream is video only.
	if len(m.Streams) != 2 {
		r3 = errors.Errorf("invalid streams=%v, %v, %v", len(m.Streams), m.String(), str)
	}

	if ts := 90; m.Format.ProbeScore < ts {
		r4 = errors.Errorf("low score=%v < %v, %v, %v", m.Format.ProbeScore, ts, m.String(), str)
	}
	if dv := m.Duration(); dv < duration/2 {
		r5 = errors.Errorf("short duration=%v < %v, %v, %v", dv, duration, m.String(), str)
	}
}

func TestScenario_WithStream_PublishCameraDuration(t *testing.T) {
	// Note that we must use long timeout, because we will interrupt FFmpeg and cause the HLS discontinuity.
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	// Control the FFmpeg to exit after N s, which cause FFmpeg restart and HLS discontinuity.
	ffmpegConfig := "max-stream-duration=15s&abnormal-fast-speed=1.5"

	if *noMediaTest {
		return
	}

	var r0, r1, r2, r3, r4, r5 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1, r2, r3, r4, r5); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	var pubSecret string
	if err := NewApi().WithAuth(ctx, "/terraform/v1/hooks/srs/secret/query", nil, &struct {
		Publish *string `json:"publish"`
	}{
		Publish: &pubSecret,
	}); err != nil {
		r0 = err
		return
	}

	var wg sync.WaitGroup
	defer wg.Wait()
	// Start FFmpeg to publish stream.
	streamID := fmt.Sprintf("stream-%v-%v", os.Getpid(), rand.Int())
	streamURL := fmt.Sprintf("%v/live/%v?secret=%v", *endpointRTMP, streamID, pubSecret)
	ffmpeg := NewFFmpeg(func(v *ffmpegClient) {
		v.args = []string{
			"-re", "-stream_loop", "-1", "-i", *srsInputFile, "-c", "copy",
			"-f", "flv", streamURL,
		}
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r1 = ffmpeg.Run(ctx, cancel)
	}()

	defer cancel()
	select {
	case <-ctx.Done():
		return
	case <-ffmpeg.ReadyCtx().Done():
	}

	// Use the publish stream url as the camera input.
	res := struct {
		Name   string `json:"name"`
		Target string `json:"target"`
		UUID   string `json:"uuid"`
		Size   int64  `json:"size"`
		Type   string `json:"type"`
	}{}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/camera/stream-url", &struct {
		StreamURL string `json:"url"`
	}{
		StreamURL: fmt.Sprintf("%v&%v", streamURL, ffmpegConfig),
	}, &res); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg camera streamUrl failed")
		return
	}

	// Use the publish stream url as the camera source.
	res.Size = 0
	res.Type = "stream"
	codec := struct {
		UUID  string `json:"uuid"`
		Audio struct {
			CodecName  string `json:"codec_name"`
			Channels   int    `json:"channels"`
			SampleRate string `json:"sample_rate"`
		} `json:"audio"`
		Video struct {
			CodecName string `json:"codec_name"`
			Profile   string `json:"profile"`
			Width     int    `json:"width"`
			Height    int    `json:"height"`
		} `json:"video"`
	}{}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/camera/source", &struct {
		Platform string        `json:"platform"`
		Files    []interface{} `json:"files"`
	}{
		Platform: "bilibili",
		Files:    []interface{}{res},
	}, &struct {
		Files []interface{} `json:"files"`
	}{
		Files: []interface{}{&codec},
	}); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg camera source failed")
		return
	}

	if err := func() error {
		if codec.UUID != res.UUID {
			return errors.Errorf("invalid codec uuid=%v, %v", codec.UUID, res.UUID)
		}
		if codec.Audio.CodecName != "aac" || codec.Audio.Channels != 2 || codec.Audio.SampleRate != "44100" {
			return errors.Errorf("invalid codec audio=%v", codec.Audio)
		}
		if codec.Video.CodecName != "h264" || codec.Video.Profile != "High" || codec.Video.Width != 768 || codec.Video.Height != 320 {
			return errors.Errorf("invalid codec video=%v", codec.Video)
		}
		return nil
	}(); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg camera source failed")
		return
	}

	// Start camera streaming.
	type VLiveConfig struct {
		Platform string      `json:"platform"`
		Server   string      `json:"server"`
		Secret   string      `json:"secret"`
		Enabled  bool        `json:"enabled"`
		Custom   bool        `json:"custom"`
		Label    string      `json:"label"`
		Files    interface{} `json:"files"`
		Action   string      `json:"action"`
	}
	conf := make(map[string]*VLiveConfig)
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/camera/secret", nil, &conf); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg camera secret failed")
		return
	}

	bilibili, ok := conf["bilibili"]
	if !ok || bilibili == nil {
		r0 = errors.Errorf("invalid bilibili secret")
		return
	}

	// Restore the state of enabled.
	backup := *bilibili
	defer func() {
		backup.Action = "update"
		logger.Tf(ctx, "restore config %v", backup)

		if backup.Server == "" {
			backup.Server = bilibili.Server
		}
		if backup.Secret == "" {
			backup.Secret = bilibili.Secret
		}

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/camera/secret", backup, nil)
	}()

	publishStreamID := fmt.Sprintf("publish-stream-%v-%v", os.Getpid(), rand.Int())
	bilibili.Secret = fmt.Sprintf("%v?secret=%v", publishStreamID, pubSecret)
	bilibili.Server = "rtmp://localhost/live/"
	bilibili.Enabled = true
	bilibili.Action = "update"
	bilibili.Custom = true
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/camera/secret", &bilibili, nil); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg camera secret failed")
		return
	}

	// Because every N s, the FFmpeg will restart, so we need to check multiple times.
	for ffmpegID := 0; ffmpegID < 3; ffmpegID++ {
		logger.Tf(ctx, "Start the %v FFmpeg", ffmpegID)
		ffprobeCtx, ffprobeCancel := context.WithCancel(ctx)

		// Cancel the FFmpeg publisher.
		go func(ctx context.Context, ffmpegID int) {
			<-ffprobeCtx.Done()
			if ffmpegID >= 2 {
				cancel()
			}
		}(ffprobeCtx, ffmpegID)

		// Start FFprobe to detect and verify stream.
		duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
		ffprobe := NewFFprobe(func(v *ffprobeClient) {
			v.dvrFile = fmt.Sprintf("srs-ffprobe-%v-FFmpeg%v.flv", publishStreamID, ffmpegID)
			v.streamURL = fmt.Sprintf("%v/live/%v.flv", *endpointHTTP, publishStreamID)
			v.duration, v.timeout = duration, time.Duration(*srsFFprobeTimeout)*time.Millisecond
		})
		wg.Add(1)
		go func() {
			defer wg.Done()
			r2 = ffprobe.Run(ctx, ffprobeCancel)
		}()

		// Fast quit for probe done.
		select {
		case <-ctx.Done():
		case <-ffprobeCtx.Done():
		case <-ffprobe.ProbeDoneCtx().Done():
		}
		ffprobeCancel()

		str, m := ffprobe.Result()
		if len(m.Streams) != 2 {
			r3 = errors.Errorf("ffmpeg=%v, invalid streams=%v, %v, %v",
				ffmpegID, len(m.Streams), m.String(), str)
		}

		if ts := 90; m.Format.ProbeScore < ts {
			r4 = errors.Errorf("ffmpeg=%v, low score=%v < %v, %v, %v",
				ffmpegID, m.Format.ProbeScore, ts, m.String(), str)
		}
		if dv := m.Duration(); dv < duration/3 {
			r5 = errors.Errorf("ffmpeg=%v, short duration=%v < %v, %v, %v",
				ffmpegID, dv, duration, m.String(), str)
		}
		if r2 != nil || r3 != nil || r4 != nil || r5 != nil {
			return
		}

		// Wait for the next FFmpeg restart.
		select {
		case <-ctx.Done():
		case <-time.After(3 * time.Second):
		}
	}
}
