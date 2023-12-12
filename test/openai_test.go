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

func TestOpenAI_TranscriptCheckConnection(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	// Ignore the test case if api secret key not set.
	apiKey, baseUrl := OpenAIConfig()
	if apiKey == "" {
		return
	}

	if r0 = NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/check", &struct {
		SecretKey string `json:"secretKey"`
		BaseUrl   string `json:"baseUrl"`
	}{
		SecretKey: apiKey, BaseUrl: baseUrl,
	}, nil); r0 != nil {
		return
	}
}

func TestOpenAI_TranscriptApplyQuery(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	// Ignore the test case if api secret key not set.
	apiKey, baseUrl := OpenAIConfig()
	if apiKey == "" {
		return
	}

	type TranscriptConfig struct {
		All       bool   `json:"all"`
		SecretKey string `json:"secretKey"`
		BaseURL   string `json:"baseURL"`
		Language  string `json:"lang"`
	}
	var conf TranscriptConfig
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/query", nil, &struct {
		Config *TranscriptConfig `json:"config"`
	}{
		Config: &conf,
	}); err != nil {
		r0 = errors.Wrapf(err, "request query failed")
		return
	}

	// Restore the state of transcode.
	backup := conf
	defer func() {
		logger.Tf(ctx, "restore config %v", backup)

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/apply", backup, nil)
	}()

	// Enable transcript.
	conf.All, conf.SecretKey, conf.BaseURL, conf.Language = true, apiKey, baseUrl, "en"
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/apply", conf, nil); err != nil {
		r0 = errors.Wrapf(err, "request apply failed")
		return
	}

	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/query", nil, &struct {
		Config *TranscriptConfig `json:"config"`
	}{
		Config: &conf,
	}); err != nil {
		r0 = errors.Wrapf(err, "request query failed")
		return
	} else if !conf.All || conf.SecretKey != apiKey || conf.BaseURL != baseUrl || conf.Language != "en" {
		r0 = errors.Errorf("invalid config %v", conf)
		return
	}

	// Disable transcript.
	conf.All, conf.SecretKey, conf.BaseURL, conf.Language = false, apiKey, baseUrl, "en"
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/apply", conf, nil); err != nil {
		r0 = errors.Wrapf(err, "request apply failed")
		return
	}

	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/query", nil, &struct {
		Config *TranscriptConfig `json:"config"`
	}{
		Config: &conf,
	}); err != nil {
		r0 = errors.Wrapf(err, "request query failed")
		return
	} else if conf.All || conf.SecretKey != apiKey || conf.BaseURL != baseUrl || conf.Language != "en" {
		r0 = errors.Errorf("invalid config %v", conf)
		return
	}
}

func TestOpenAI_WithStream_TranscriptASR(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsLongTimeout)*time.Millisecond)
	defer cancel()

	if *noMediaTest {
		return
	}

	var r0, r1 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1); err != nil {
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

	// Ignore the test case if api secret key not set.
	apiKey, baseUrl := OpenAIConfig()
	if apiKey == "" {
		return
	}

	type TranscriptConfig struct {
		All       bool   `json:"all"`
		SecretKey string `json:"secretKey"`
		BaseURL   string `json:"baseURL"`
		Language  string `json:"lang"`
	}
	var conf TranscriptConfig

	type TranscriptTask struct {
		UUID string `json:"uuid"`
	}
	var task TranscriptTask

	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/query", nil, &struct {
		Config *TranscriptConfig `json:"config"`
		Task   *TranscriptTask   `json:"task"`
	}{
		Config: &conf, Task: &task,
	}); err != nil {
		r0 = errors.Wrapf(err, "request query failed")
		return
	}

	// Restore the state of transcode.
	backup := conf
	defer func() {
		logger.Tf(ctx, "restore config %v", backup)

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/apply", backup, nil)
	}()

	// Enable transcript.
	conf.All, conf.SecretKey, conf.BaseURL, conf.Language = true, apiKey, baseUrl, "en"
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/apply", conf, nil); err != nil {
		r0 = errors.Wrapf(err, "request apply failed")
		return
	}

	// Always disable and cleanup transcript.
	defer func(ctx context.Context) {
		// Disable transcript.
		conf.All, conf.SecretKey, conf.BaseURL, conf.Language = false, apiKey, baseUrl, "en"
		if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/apply", conf, nil); err != nil {
			r0 = errors.Wrapf(err, "request apply failed")
			return
		}

		// Reset and cleanup transcript.
		if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/reset", task, nil); err != nil {
			r0 = errors.Wrapf(err, "request reset failed")
			return
		}
	}(ctx)

	// Context for publish stream.
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	// Start publish stream, about 10s.
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
	go func(ctx context.Context) {
		defer wg.Done()
		r1 = ffmpeg.Run(ctx, cancel)
	}(ctx)

	// Wait for record to save file.
	// There should have some transcript files.
	type TranscriptSegment struct {
		ASR      string  `json:"asr"`
		Duration float64 `json:"duration"`
		Size     int64   `json:"size"`
	}
	var segments []TranscriptSegment

	for i := 0; i < 10; i++ {
		if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/overlay-queue", nil, &struct {
			Count    int                  `json:"count"`
			Segments *[]TranscriptSegment `json:"segments"`
		}{
			Segments: &segments,
		}); err != nil {
			r0 = errors.Wrapf(err, "request query failed")
			return
		}

		if len(segments) > 0 {
			break
		}

		select {
		case <-ctx.Done():
		case <-time.After(5 * time.Second):
		}
	}

	// Cancel ffmpeg publisher.
	defer cancel()

	// Check result.
	if len(segments) < 1 {
		r0 = errors.Errorf("invalid segments %v", segments)
		return
	}

	if segment := segments[0]; segment.ASR == "" || segment.Duration <= 0 || segment.Size <= 0 {
		r0 = errors.Errorf("invalid segment %v", segment)
		return
	}
}

func TestOpenAI_WithStream_Transcript_ClearSubtitle(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsLongTimeout)*time.Millisecond)
	defer cancel()

	if *noMediaTest {
		return
	}

	var r0, r1 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1); err != nil {
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

	// Ignore the test case if api secret key not set.
	apiKey, baseUrl := OpenAIConfig()
	if apiKey == "" {
		return
	}

	type TranscriptConfig struct {
		All       bool   `json:"all"`
		SecretKey string `json:"secretKey"`
		BaseURL   string `json:"baseURL"`
		Language  string `json:"lang"`
	}
	var conf TranscriptConfig

	type TranscriptTask struct {
		UUID string `json:"uuid"`
	}
	var task TranscriptTask

	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/query", nil, &struct {
		Config *TranscriptConfig `json:"config"`
		Task   *TranscriptTask   `json:"task"`
	}{
		Config: &conf, Task: &task,
	}); err != nil {
		r0 = errors.Wrapf(err, "request query failed")
		return
	}

	// Restore the state of transcode.
	backup := conf
	defer func() {
		logger.Tf(ctx, "restore config %v", backup)

		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/apply", backup, nil)
	}()

	// Enable transcript.
	conf.All, conf.SecretKey, conf.BaseURL, conf.Language = true, apiKey, baseUrl, "en"
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/apply", conf, nil); err != nil {
		r0 = errors.Wrapf(err, "request apply failed")
		return
	}

	// Always disable and cleanup transcript.
	defer func(ctx context.Context) {
		// Disable transcript.
		conf.All, conf.SecretKey, conf.BaseURL, conf.Language = false, apiKey, baseUrl, "en"
		if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/apply", conf, nil); err != nil {
			r0 = errors.Wrapf(err, "request apply failed")
			return
		}

		// Reset and cleanup transcript.
		if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/reset", task, nil); err != nil {
			r0 = errors.Wrapf(err, "request reset failed")
			return
		}
	}(ctx)

	// Context for publish stream.
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	// Start publish stream, about 10s.
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

	ffmpegCtx, ffmpegCancel := context.WithCancel(ctx)
	wg.Add(1)
	go func() {
		defer wg.Done()
		r1 = ffmpeg.Run(ffmpegCtx, cancel)
	}()
	defer ffmpegCancel()

	// Wait for record to save file.
	// There should have some transcript files.
	type TranscriptSegment struct {
		TsID     string  `json:"tsid"`
		ASR      string  `json:"asr"`
		Duration float64 `json:"duration"`
		Size     int64   `json:"size"`
		// User clear the ASR subtitle.
		UserClearASR bool `json:"uca"`
	}

	querySegments := func(api string) []TranscriptSegment {
		var segments []TranscriptSegment
		for i := 0; i < 10; i++ {
			if err := NewApi().WithAuth(ctx, api, nil, &struct {
				Count    int                  `json:"count"`
				Segments *[]TranscriptSegment `json:"segments"`
			}{
				Segments: &segments,
			}); err != nil {
				r0 = errors.Wrapf(err, "request query %v failed", api)
				return nil
			}

			if len(segments) > 0 {
				break
			}

			select {
			case <-ctx.Done():
			case <-time.After(5 * time.Second):
			}
		}
		return segments
	}

	if segments := querySegments("/terraform/v1/ai/transcript/overlay-queue"); len(segments) < 1 {
		r0 = errors.Errorf("invalid segments %v", segments)
		return
	} else if segment := segments[0]; segment.ASR == "" || segment.Duration <= 0 || segment.Size <= 0 {
		r0 = errors.Errorf("invalid segment %v", segment)
		return
	}

	// Cancel ffmpeg publisher.
	ffmpegCancel()

	// Cancel the task.
	defer cancel()

	// Query the fix queue, should have at least one segment.
	var segment TranscriptSegment
	if segments := querySegments("/terraform/v1/ai/transcript/fix-queue"); len(segments) < 1 {
		r0 = errors.Errorf("invalid segments %v", segments)
		return
	} else if segment = segments[0]; segment.TsID == "" || segment.ASR == "" || segment.Duration <= 0 || segment.Size <= 0 {
		r0 = errors.Errorf("invalid segment %v", segment)
		return
	}

	// Clear the subtitle.
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai/transcript/clear-subtitle", &struct {
		UUID string `json:"uuid"`
		TSID string `json:"tsid"`
	}{
		UUID: task.UUID, TSID: segment.TsID,
	}, nil); err != nil {
		r0 = errors.Wrapf(err, "request clear subtitle %v failed", segment)
		return
	}

	// Check the segment again, should be cleared.
	segments := querySegments("/terraform/v1/ai/transcript/fix-queue")
	for _, s := range segments {
		if s.TsID == segment.TsID {
			if !s.UserClearASR {
				r0 = errors.Errorf("invalid segment %v", s)
			}
			return
		}
	}

	r0 = errors.Errorf("invalid segments %v", segments)
	return
}
