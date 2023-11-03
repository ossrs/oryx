package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
)

func TestScenario_WithStream_PublishVliveStreamUrl(t *testing.T) {
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

	// Use the publish stream url as the vlive input.
	res := struct {
		Name   string `json:"name"`
		Target string `json:"target"`
		UUID   string `json:"uuid"`
		Size   int64  `json:"size"`
		Type   string `json:"type"`
	}{}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/vlive/streamUrl?url="+streamURL, nil, &res); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg vlive streamUrl failed")
		return
	}

	// Use the publish stream url as the vlive source.
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
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/vlive/source", &struct {
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
		r0 = errors.Wrapf(err, "request ffmpeg vlive source failed")
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
		r0 = errors.Wrapf(err, "request ffmpeg vlive source failed")
		return
	}

	// Start virtual live streaming.
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
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/vlive/secret", nil, &conf); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg vlive secret failed")
		return
	}

	bilibili, ok := conf["bilibili"]
	if !ok || bilibili == nil {
		r0 = errors.Errorf("invalid bilibili secret")
		return
	}
	bilibili.Action = "update"
	bilibili.Custom = true

	// Restore the state of enabled.
	backup := *bilibili
	defer func() {
		logger.Tf(ctx, "restore config %v", backup)

		if backup.Server == "" {
			backup.Server = bilibili.Server
		}
		if backup.Secret == "" {
			backup.Secret = bilibili.Secret
		}

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/vlive/secret", backup, nil)
	}()

	publishStreamID := fmt.Sprintf("publish-stream-%v-%v", os.Getpid(), rand.Int())
	bilibili.Secret = fmt.Sprintf("%v?secret=%v", publishStreamID, pubSecret)
	bilibili.Server = "rtmp://localhost/live/"
	bilibili.Enabled = true
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/vlive/secret", &bilibili, nil); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg vlive secret failed")
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

func TestScenario_WithStream_PublishVLivePlayFlv(t *testing.T) {
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

	// Copy virtual live source file to /data/upload and platform/containers/data/upload
	destDirs := []string{
		"/data/upload/",
		"platform/containers/data/upload",
		"../platform/containers/data/upload",
	}
	if err := copyToDest(ctx, *srsInputFile, destDirs...); err != nil {
		r0 = errors.Wrapf(err, "copy %v to %v", *srsInputFile, destDirs)
		return
	}

	// Get first matched source file.
	sourceFile := getExistsFile(ctx, filepath.Base(*srsInputFile), destDirs...)
	if sourceFile == "" {
		r0 = errors.Errorf("no source file found")
		return
	}

	// If not absolution path, always use short path in upload.
	if !strings.HasPrefix(sourceFile, "/data/upload/") {
		sourceFile = "upload/" + path.Base(sourceFile)
	}

	// Use the file as uploaded file.
	res := struct {
		Name   string `json:"name"`
		Size   int64  `json:"size"`
		Target string `json:"target"`
		UUID   string `json:"uuid"`
	}{}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/vlive/server?file="+sourceFile, nil, &res); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg vlive server failed")
		return
	}

	// Use the file as source file.
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
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/vlive/source", &struct {
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
		r0 = errors.Wrapf(err, "request ffmpeg vlive source failed")
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
		r0 = errors.Wrapf(err, "request ffmpeg vlive source failed")
		return
	}

	// Start virtual live streaming.
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
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/vlive/secret", nil, &conf); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg vlive secret failed")
		return
	}

	bilibili, ok := conf["bilibili"]
	if !ok || bilibili == nil {
		r0 = errors.Errorf("invalid bilibili secret")
		return
	}
	bilibili.Action = "update"

	// Restore the state of enabled.
	backup := *bilibili
	defer func() {
		logger.Tf(ctx, "restore config %v", backup)

		if backup.Server == "" {
			backup.Server = bilibili.Server
		}
		if backup.Secret == "" {
			backup.Secret = bilibili.Secret
		}

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/vlive/secret", backup, nil)
	}()

	streamID := fmt.Sprintf("stream-%v-%v", os.Getpid(), rand.Int())
	bilibili.Secret = fmt.Sprintf("%v?secret=%v", streamID, pubSecret)
	bilibili.Server = "rtmp://localhost/live/"
	bilibili.Enabled = true
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/vlive/secret", &bilibili, nil); err != nil {
		r0 = errors.Wrapf(err, "request ffmpeg vlive secret failed")
		return
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	// Start FFprobe to detect and verify stream.
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", streamID)
		v.streamURL = fmt.Sprintf("%v/live/%v.flv", *endpointHTTP, streamID)
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

func TestScenario_WithStream_PublishRtmpRecordMp4(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsLongTimeout)*time.Millisecond)
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

	// Query the old config.
	backup := make(map[string]interface{})
	if err := NewApi().WithAuth(ctx, "/terraform/v1/hooks/record/query", nil, &backup); err != nil {
		r0 = errors.Wrapf(err, "request record query failed")
		return
	}
	defer func() {
		logger.Tf(ctx, "restore config %v", backup)

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/hooks/record/apply", backup, nil)
	}()

	// Enable the record worker.
	if err := NewApi().WithAuth(ctx, "/terraform/v1/hooks/record/apply", &struct {
		All bool `json:"all"`
	}{true}, nil); err != nil {
		r0 = errors.Wrapf(err, "request record apply failed")
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

	// Wait for record to save file.
	select {
	case <-ctx.Done():
	case <-time.After(25 * time.Second):
	}

	// Stop record worker.
	if err := NewApi().WithAuth(ctx, "/terraform/v1/hooks/record/apply", &struct {
		All bool `json:"all"`
	}{false}, nil); err != nil {
		r0 = errors.Wrapf(err, "request record apply failed")
		return
	}
	logger.Tf(ctx, "stop record worker done")

	// Query the record file.
	type RecordFile struct {
		Stream   string  `json:"stream"`
		UUID     string  `json:"uuid"`
		Duration float64 `json:"duration"`
		Progress bool    `json:"progress"`
	}
	var recordFile *RecordFile
	defer func() {
		if recordFile == nil || recordFile.UUID == "" {
			return
		}
		logger.Tf(ctx, "remove record file %v", recordFile)

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/hooks/record/remove", &struct {
			UUID string `json:"uuid"`
		}{recordFile.UUID}, nil)
	}()
	defer cancel()

	for i := 0; i < 60; i++ {
		files := []RecordFile{}
		if err := NewApi().WithAuth(ctx, "/terraform/v1/hooks/record/files", nil, &files); err != nil {
			r0 = errors.Wrapf(err, "request record files failed")
			return
		}

		for _, file := range files {
			if file.Stream == streamID {
				recordFile = &file
				break
			}
		}

		if recordFile == nil || recordFile.Progress {
			select {
			case <-ctx.Done():
				r0 = errors.Wrapf(ctx.Err(), "record file not found")
				return
			case <-time.After(1 * time.Second):
				continue
			}
		}
		break
	}

	if recordFile == nil {
		r0 = errors.Errorf("record file not found")
		return
	}
	if recordFile.Progress {
		r0 = errors.Errorf("record file is progress, %v", recordFile)
		return
	}
	if recordFile.Duration < 10 {
		r0 = errors.Errorf("record file duration too short, %v", recordFile)
		return
	}

	time.Sleep(3 * time.Second)
	logger.Tf(ctx, "record ok, file is %v", recordFile)
	cancel()
}

func TestScenario_WithStream_PublishRtmpForwardPlatform(t *testing.T) {
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

	// Query the old config.
	type ForwardConfig struct {
		Action   string `json:"action"`
		Enabled  bool   `json:"enabled"`
		Custom   bool   `json:"custom"`
		Label    string `json:"label"`
		Platform string `json:"platform"`
		Secret   string `json:"secret"`
		Server   string `json:"server"`
	}
	conf := make(map[string]*ForwardConfig)
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/forward/secret", nil, &conf); err != nil {
		r0 = errors.Wrapf(err, "request forward query failed")
		return
	}

	forwardStreamID := fmt.Sprintf("forward-stream-%v-%v", os.Getpid(), rand.Int())
	bilibili, ok := conf["bilibili"]
	if !ok || bilibili == nil {
		bilibili = &ForwardConfig{
			Enabled:  false,
			Custom:   true,
			Label:    "Test",
			Platform: "bilibili",
			Secret:   fmt.Sprintf("%v?secret=%v", forwardStreamID, pubSecret),
			Server:   "rtmp://localhost/live/",
		}
		conf["bilibili"] = bilibili
	}
	bilibili.Action = "update"

	// Restore the state of forward.
	backup := *bilibili
	defer func() {
		logger.Tf(ctx, "restore config %v", backup)

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/forward/secret", backup, nil)
	}()

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

	// Enable the forward worker.
	select {
	case <-ctx.Done():
	case <-time.After(3 * time.Second):
	}

	bilibili.Secret = fmt.Sprintf("%v?secret=%v", forwardStreamID, pubSecret)
	bilibili.Server = "rtmp://localhost/live/"
	bilibili.Enabled = true
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/forward/secret", bilibili, nil); err != nil {
		r0 = errors.Wrapf(err, "request forward apply failed")
		return
	}

	// Start FFprobe to detect and verify stream.
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", forwardStreamID)
		v.streamURL = fmt.Sprintf("%v/live/%v.flv", *endpointHTTP, forwardStreamID)
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

func TestScenario_WithStream_PublishRtmpTranscodeDefault(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	if *noMediaTest {
		return
	}

	var r0, r1, r2, r3, r4, r5, r6 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1, r2, r3, r4, r5, r6); err != nil {
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

	// Query the old config.
	type TranscodeConfig struct {
		All           bool   `json:"all"`
		VideoCodec    string `json:"vcodec"`
		VideoProfile  string `json:"vprofile"`
		VideoPreset   string `json:"vpreset"`
		VideoBitrate  int    `json:"vbitrate"`
		AudioCodec    string `json:"acodec"`
		AudioBitrate  int    `json:"abitrate"`
		AudioChannels int    `json:"achannels"`
		Server        string `json:"server"`
		Secret        string `json:"secret"`
	}
	var conf TranscodeConfig
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/transcode/query", nil, &conf); err != nil {
		r0 = errors.Wrapf(err, "request transcode query failed")
		return
	}

	// Restore the state of transcode.
	backup := conf
	defer func() {
		logger.Tf(ctx, "restore config %v", backup)

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/transcode/apply", backup, nil)
	}()

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

	// Enable the transcode worker.
	select {
	case <-ctx.Done():
	case <-time.After(3 * time.Second):
	}

	transcodeStreamID := fmt.Sprintf("transcoded-stream-%v-%v", os.Getpid(), rand.Int())
	conf.All = true
	conf.Server = "rtmp://localhost/live/"
	conf.Secret = fmt.Sprintf("%v?secret=%v", transcodeStreamID, pubSecret)
	conf.VideoCodec = "libx264"
	conf.VideoBitrate = 200
	conf.VideoProfile = "baseline"
	conf.VideoPreset = "ultrafast"
	conf.AudioCodec = "aac"
	conf.AudioBitrate = 16
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/transcode/apply", &conf, nil); err != nil {
		r0 = errors.Wrapf(err, "request transcode apply failed")
		return
	}

	// Check the transcode task.
	select {
	case <-ctx.Done():
	case <-time.After(1 * time.Second):
	}

	task := struct {
		UUID    string `json:"uuid"`
		Enabled bool   `json:"enabled"`
		Input   string `json:"input"`
		Output  string `json:"output"`
		Frame   struct {
			Log    string `json:"log"`
			Update string `json:"update"`
		} `json:"frame"`
	}{}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/transcode/task", nil, &task); err != nil {
		r0 = errors.Wrapf(err, "request transcode query failed")
		return
	}
	if !task.Enabled || task.Input == "" || task.Output == "" || task.Frame.Log == "" || task.Frame.Update == "" ||
		task.UUID == "" {
		r0 = errors.Errorf("invalid task=%v", task)
		return
	}

	// Start FFprobe to detect and verify stream.
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", transcodeStreamID)
		v.streamURL = fmt.Sprintf("%v/live/%v.flv", *endpointHTTP, transcodeStreamID)
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

	if m.Audio().Channels != 2 {
		r6 = errors.Errorf("invalid audio channels=%v, %v, %v", m.Audio().Channels, m.String(), str)
	}
}

func TestScenario_WithStream_PublishRtmpTranscodeFollow(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	if *noMediaTest {
		return
	}

	var r0, r1, r2, r3, r4, r5, r6 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1, r2, r3, r4, r5, r6); err != nil {
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

	// Query the old config.
	type TranscodeConfig struct {
		All           bool   `json:"all"`
		VideoCodec    string `json:"vcodec"`
		VideoProfile  string `json:"vprofile"`
		VideoPreset   string `json:"vpreset"`
		VideoBitrate  int    `json:"vbitrate"`
		AudioCodec    string `json:"acodec"`
		AudioBitrate  int    `json:"abitrate"`
		AudioChannels int    `json:"achannels"`
		Server        string `json:"server"`
		Secret        string `json:"secret"`
	}
	var conf TranscodeConfig
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/transcode/query", nil, &conf); err != nil {
		r0 = errors.Wrapf(err, "request transcode query failed")
		return
	}

	// Restore the state of transcode.
	backup := conf
	defer func() {
		logger.Tf(ctx, "restore config %v", backup)

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/transcode/apply", backup, nil)
	}()

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

	// Enable the transcode worker.
	select {
	case <-ctx.Done():
	case <-time.After(3 * time.Second):
	}

	transcodeStreamID := fmt.Sprintf("transcoded-stream-%v-%v", os.Getpid(), rand.Int())
	conf.All = true
	conf.Server = "rtmp://localhost/live/"
	conf.Secret = fmt.Sprintf("%v?secret=%v", transcodeStreamID, pubSecret)
	conf.VideoCodec = "libx264"
	conf.VideoBitrate = 200
	conf.VideoProfile = "baseline"
	conf.VideoPreset = "ultrafast"
	conf.AudioCodec = "aac"
	conf.AudioBitrate = 16
	conf.AudioChannels = 0
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/transcode/apply", &conf, nil); err != nil {
		r0 = errors.Wrapf(err, "request transcode apply failed")
		return
	}

	// Start FFprobe to detect and verify stream.
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", transcodeStreamID)
		v.streamURL = fmt.Sprintf("%v/live/%v.flv", *endpointHTTP, transcodeStreamID)
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

	if m.Audio().Channels != 2 {
		r6 = errors.Errorf("invalid audio channels=%v, %v, %v", m.Audio().Channels, m.String(), str)
	}
}

func TestScenario_WithStream_PublishRtmpTranscodeMono(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	if *noMediaTest {
		return
	}

	var r0, r1, r2, r3, r4, r5, r6 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1, r2, r3, r4, r5, r6); err != nil {
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

	// Query the old config.
	type TranscodeConfig struct {
		All           bool   `json:"all"`
		VideoCodec    string `json:"vcodec"`
		VideoProfile  string `json:"vprofile"`
		VideoPreset   string `json:"vpreset"`
		VideoBitrate  int    `json:"vbitrate"`
		AudioCodec    string `json:"acodec"`
		AudioBitrate  int    `json:"abitrate"`
		AudioChannels int    `json:"achannels"`
		Server        string `json:"server"`
		Secret        string `json:"secret"`
	}
	var conf TranscodeConfig
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/transcode/query", nil, &conf); err != nil {
		r0 = errors.Wrapf(err, "request transcode query failed")
		return
	}

	// Restore the state of transcode.
	backup := conf
	defer func() {
		logger.Tf(ctx, "restore config %v", backup)

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/transcode/apply", backup, nil)
	}()

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

	// Enable the transcode worker.
	select {
	case <-ctx.Done():
	case <-time.After(3 * time.Second):
	}

	transcodeStreamID := fmt.Sprintf("transcoded-stream-%v-%v", os.Getpid(), rand.Int())
	conf.All = true
	conf.Server = "rtmp://localhost/live/"
	conf.Secret = fmt.Sprintf("%v?secret=%v", transcodeStreamID, pubSecret)
	conf.VideoCodec = "libx264"
	conf.VideoBitrate = 200
	conf.VideoProfile = "baseline"
	conf.VideoPreset = "ultrafast"
	conf.AudioCodec = "aac"
	conf.AudioBitrate = 16
	conf.AudioChannels = 1
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/transcode/apply", &conf, nil); err != nil {
		r0 = errors.Wrapf(err, "request transcode apply failed")
		return
	}

	// Start FFprobe to detect and verify stream.
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", transcodeStreamID)
		v.streamURL = fmt.Sprintf("%v/live/%v.flv", *endpointHTTP, transcodeStreamID)
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

	if m.Audio().Channels != 1 {
		r6 = errors.Errorf("invalid audio channels=%v, %v, %v", m.Audio().Channels, m.String(), str)
	}
}

func TestScenario_WithStream_PublishRtmpTranscodeStereo(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	if *noMediaTest {
		return
	}

	var r0, r1, r2, r3, r4, r5, r6 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1, r2, r3, r4, r5, r6); err != nil {
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

	// Query the old config.
	type TranscodeConfig struct {
		All           bool   `json:"all"`
		VideoCodec    string `json:"vcodec"`
		VideoProfile  string `json:"vprofile"`
		VideoPreset   string `json:"vpreset"`
		VideoBitrate  int    `json:"vbitrate"`
		AudioCodec    string `json:"acodec"`
		AudioBitrate  int    `json:"abitrate"`
		AudioChannels int    `json:"achannels"`
		Server        string `json:"server"`
		Secret        string `json:"secret"`
	}
	var conf TranscodeConfig
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/transcode/query", nil, &conf); err != nil {
		r0 = errors.Wrapf(err, "request transcode query failed")
		return
	}

	// Restore the state of transcode.
	backup := conf
	defer func() {
		logger.Tf(ctx, "restore config %v", backup)

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/transcode/apply", backup, nil)
	}()

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

	// Enable the transcode worker.
	select {
	case <-ctx.Done():
	case <-time.After(3 * time.Second):
	}

	transcodeStreamID := fmt.Sprintf("transcoded-stream-%v-%v", os.Getpid(), rand.Int())
	conf.All = true
	conf.Server = "rtmp://localhost/live/"
	conf.Secret = fmt.Sprintf("%v?secret=%v", transcodeStreamID, pubSecret)
	conf.VideoCodec = "libx264"
	conf.VideoBitrate = 200
	conf.VideoProfile = "baseline"
	conf.VideoPreset = "ultrafast"
	conf.AudioCodec = "aac"
	conf.AudioBitrate = 16
	conf.AudioChannels = 2
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ffmpeg/transcode/apply", &conf, nil); err != nil {
		r0 = errors.Wrapf(err, "request transcode apply failed")
		return
	}

	// Start FFprobe to detect and verify stream.
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", transcodeStreamID)
		v.streamURL = fmt.Sprintf("%v/live/%v.flv", *endpointHTTP, transcodeStreamID)
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

	if m.Audio().Channels != 2 {
		r6 = errors.Errorf("invalid audio channels=%v, %v, %v", m.Audio().Channels, m.String(), str)
	}
}

func TestScenario_WithStream_CallbackOnPublishSuccess(t *testing.T) {
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

	type CallbackConfig struct {
		All    bool   `json:"all"`
		Opaque string `json:"opaque"`
		Target string `json:"target"`
	}
	var conf CallbackConfig
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/hooks/query", nil, &conf); err != nil {
		r0 = errors.Wrapf(err, "request hooks apply failed")
		return
	}

	// Restore the state of transcode.
	backup := conf
	defer func() {
		logger.Tf(ctx, "restore config %v", backup)

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/mgmt/hooks/apply", backup, nil)
	}()

	// Enable the callback worker.
	conf.All = true
	conf.Target = fmt.Sprintf("%v/terraform/v1/mgmt/hooks/example?fail=false", *endpoint)
	conf.Opaque = fmt.Sprintf("opaque-%v", rand.Int())
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/hooks/apply", &conf, nil); err != nil {
		r0 = errors.Wrapf(err, "request hooks apply failed")
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

	// Start FFprobe to detect and verify stream.
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", streamID)
		v.streamURL = fmt.Sprintf("%v/live/%v.flv", *endpointHTTP, streamID)
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
	if dv := m.Duration(); dv < duration/3 {
		r5 = errors.Errorf("short duration=%v < %v, %v, %v", dv, duration, m.String(), str)
	}
}

func TestScenario_WithStream_CallbackOnPublishFailed(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	if *noMediaTest {
		return
	}

	var r0, r1, r2, r3 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1, r2, r3); err != nil {
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

	type CallbackConfig struct {
		All    bool   `json:"all"`
		Opaque string `json:"opaque"`
		Target string `json:"target"`
	}
	var conf CallbackConfig
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/hooks/query", nil, &conf); err != nil {
		r0 = errors.Wrapf(err, "request hooks apply failed")
		return
	}

	// Restore the state of transcode.
	backup := conf
	defer func() {
		logger.Tf(ctx, "restore config %v", backup)

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/mgmt/hooks/apply", backup, nil)
	}()

	// Enable the callback worker.
	conf.All = true
	conf.Target = fmt.Sprintf("%v/terraform/v1/mgmt/hooks/example?fail=true", *endpoint)
	conf.Opaque = fmt.Sprintf("opaque-%v", rand.Int())
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/hooks/apply", &conf, nil); err != nil {
		r0 = errors.Wrapf(err, "request hooks apply failed")
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

	// Start FFprobe to detect and verify stream.
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", streamID)
		v.streamURL = fmt.Sprintf("%v/live/%v.flv", *endpointHTTP, streamID)
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

	// Should have no stream for callback failed.
	str, m := ffprobe.Result()
	if len(m.Streams) != 0 {
		r3 = errors.Errorf("invalid streams=%v, %v, %v", len(m.Streams), m.String(), str)
	}
}
