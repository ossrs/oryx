package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
)

func TestMedia_WithStream_PublishRtmpPlayFlv_SecretQuery(t *testing.T) {
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
	if dv := m.Duration(); dv < duration/5 {
		r5 = errors.Errorf("short duration=%v < %v, %v, %v", dv, duration, m.String(), str)
	}
}

func TestMedia_WithStream_PublishRtmpPlayFlv_SecretStream(t *testing.T) {
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
	streamID := fmt.Sprintf("stream-%v-%v-%v", pubSecret, os.Getpid(), rand.Int())
	streamURL := fmt.Sprintf("%v/live/%v", *endpointRTMP, streamID)
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
	if dv := m.Duration(); dv < duration/5 {
		r5 = errors.Errorf("short duration=%v < %v, %v, %v", dv, duration, m.String(), str)
	}
}

func TestMedia_WithStream_PublishRtmpPlayHls_SecretQuery(t *testing.T) {
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

	// Start FFprobe to detect and verify stream.
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", streamID)
		v.streamURL = fmt.Sprintf("%v/live/%v.m3u8", *endpointHTTP, streamID)
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

	// Note that HLS score is low, so we only check duration. Note that only check half of duration, because we
	// might get only some pieces of segments.
	if dv := m.Duration(); dv < duration/5 {
		r4 = errors.Errorf("short duration=%v < %v, %v, %v", dv, duration/2, m.String(), str)
	}
}

func TestMedia_WithStream_PublishSrtPlayFlv_SecretQuery(t *testing.T) {
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
	streamURL := fmt.Sprintf("%v?streamid=#!::r=live/%v?secret=%v,m=publish", *endpointSRT, streamID, pubSecret)
	ffmpeg := NewFFmpeg(func(v *ffmpegClient) {
		v.args = []string{
			"-re", "-stream_loop", "-1", "-i", *srsInputFile, "-c", "copy",
			"-f", "mpegts", streamURL,
		}
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r1 = ffmpeg.Run(ctx, cancel)
	}()

	// Start to probe SRT stream after published N seconds, to wait stream ready.
	select {
	case <-ctx.Done():
	case <-ffmpeg.ReadyCtx().Done():
	}
	select {
	case <-ctx.Done():
	case <-time.After(time.Duration(*srsFFprobeTimeout) / 10 * time.Millisecond):
	}

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
	if dv := m.Duration(); dv < duration/5 {
		r5 = errors.Errorf("short duration=%v < %v, %v, %v", dv, duration, m.String(), str)
	}
}

func TestMedia_WithStream_PublishRtmpPlayHls_NoHlsCtx(t *testing.T) {
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

	type Data struct {
		NoHlsCtx bool `json:"noHlsCtx"`
	}

	if true {
		initData := Data{}
		if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/hphls/query", nil, &initData); err != nil {
			r0 = err
			return
		}
		defer func() {
			// TODO: FIXME: Remove it after fix the bug.
			time.Sleep(10 * time.Second)

			// The ctx has already been cancelled by test case, which will cause the request failed.
			ctx := context.Background()
			if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/hphls/update", &initData, nil); err != nil {
				logger.Tf(ctx, "restore hphls config failed %+v", err)
			}
		}()
	}

	noHlsCtx := Data{NoHlsCtx: true}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/hphls/update", &noHlsCtx, nil); err != nil {
		r0 = err
		return
	}

	verifyData := Data{}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/hphls/query", nil, &verifyData); err != nil {
		r0 = err
		return
	} else if verifyData.NoHlsCtx != true {
		r0 = errors.Errorf("invalid response %+v", verifyData)
	}

	// TODO: FIXME: Remove it after fix the bug.
	time.Sleep(10 * time.Second)

	///////////////////////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////////////////////
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

	// Don't cancel the context, for we need to verify the HLS stream.
	_, noCancel := context.WithCancel(context.Background())

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
		r1 = ffmpeg.Run(ctx, noCancel)
	}()

	// Start FFprobe to detect and verify stream.
	var hlsStreamURL string
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", streamID)
		v.streamURL = fmt.Sprintf("%v/live/%v.m3u8", *endpointHTTP, streamID)
		v.duration, v.timeout = duration, time.Duration(*srsFFprobeTimeout)*time.Millisecond
		hlsStreamURL = v.streamURL
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r2 = ffprobe.Run(ctx, noCancel)
	}()

	// Don't quit for probe done.
	select {
	case <-ctx.Done():
	case <-ffprobe.ProbeDoneCtx().Done():
	}

	str, m := ffprobe.Result()
	if len(m.Streams) != 2 {
		r3 = errors.Errorf("invalid streams=%v, %v, %v", len(m.Streams), m.String(), str)
	}

	// Note that HLS score is low, so we only check duration. Note that only check half of duration, because we
	// might get only some pieces of segments.
	if dv := m.Duration(); dv < duration/5 {
		r4 = errors.Errorf("short duration=%v < %v, %v, %v", dv, duration/2, m.String(), str)
	}

	// Check the HLS playlist, should with hls context.
	var body string
	if err := NewApi().Request(ctx, hlsStreamURL, nil, false, &body); err != nil {
		r5 = err
		return
	}

	// #EXTM3U
	// #EXT-X-VERSION:3
	// #EXT-X-MEDIA-SEQUENCE:0
	// #EXT-X-TARGETDURATION:15
	// #EXT-X-DISCONTINUITY
	// #EXTINF:10.008, no desc
	// stream-15318-7260362267190950336-0.ts
	// #EXTINF:11.989, no desc
	// stream-15318-7260362267190950336-1.ts
	// #EXTINF:11.994, no desc
	// stream-15318-7260362267190950336-2.ts
	if strings.Contains(body, ".m3u8?hls_ctx=") {
		r5 = errors.Errorf("invalid hls playlist=%v", body)
	}
	if !strings.Contains(body, "#EXTINF:") {
		r5 = errors.Errorf("invalid hls playlist=%v", body)
	}

	cancel()
}

func TestMedia_WithStream_PublishRtmpPlayHls_WithHlsCtx(t *testing.T) {
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

	type Data struct {
		NoHlsCtx bool `json:"noHlsCtx"`
	}

	if true {
		initData := Data{}
		if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/hphls/query", nil, &initData); err != nil {
			r0 = err
			return
		}
		defer func() {
			// TODO: FIXME: Remove it after fix the bug.
			time.Sleep(10 * time.Second)

			// The ctx has already been cancelled by test case, which will cause the request failed.
			ctx := context.Background()
			if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/hphls/update", &initData, nil); err != nil {
				logger.Tf(ctx, "restore hphls config failed %+v", err)
			}
		}()
	}

	noHlsCtx := Data{NoHlsCtx: false}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/hphls/update", &noHlsCtx, nil); err != nil {
		r0 = err
		return
	}

	verifyData := Data{}
	if err := NewApi().WithAuth(ctx, "/terraform/v1/mgmt/hphls/query", nil, &verifyData); err != nil {
		r0 = err
		return
	} else if verifyData.NoHlsCtx != false {
		r0 = errors.Errorf("invalid response %+v", verifyData)
	}

	// TODO: FIXME: Remove it after fix the bug.
	time.Sleep(10 * time.Second)

	///////////////////////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////////////////////
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

	// Don't cancel the context, for we need to verify the HLS stream.
	_, noCancel := context.WithCancel(context.Background())

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
		r1 = ffmpeg.Run(ctx, noCancel)
	}()

	// Start FFprobe to detect and verify stream.
	var hlsStreamURL string
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", streamID)
		v.streamURL = fmt.Sprintf("%v/live/%v.m3u8", *endpointHTTP, streamID)
		v.duration, v.timeout = duration, time.Duration(*srsFFprobeTimeout)*time.Millisecond
		hlsStreamURL = v.streamURL
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r2 = ffprobe.Run(ctx, noCancel)
	}()

	// Don't quit for probe done.
	select {
	case <-ctx.Done():
	case <-ffprobe.ProbeDoneCtx().Done():
	}

	str, m := ffprobe.Result()
	if len(m.Streams) != 2 {
		r3 = errors.Errorf("invalid streams=%v, %v, %v", len(m.Streams), m.String(), str)
	}

	// Note that HLS score is low, so we only check duration. Note that only check half of duration, because we
	// might get only some pieces of segments.
	if dv := m.Duration(); dv < duration/5 {
		r4 = errors.Errorf("short duration=%v < %v, %v, %v", dv, duration/2, m.String(), str)
	}

	// Check the HLS playlist, should with hls context.
	var body string
	if err := NewApi().Request(ctx, hlsStreamURL, nil, false, &body); err != nil {
		r5 = err
		return
	}

	// #EXTM3U
	// #EXT-X-STREAM-INF:BANDWIDTH=1,AVERAGE-BANDWIDTH=1
	// /live/stream-22525-1463247945540465917.m3u8?hls_ctx=84q1332s
	if !strings.Contains(body, ".m3u8?hls_ctx=") {
		r5 = errors.Errorf("invalid hls playlist=%v", body)
	}
	if strings.Contains(body, "#EXTINF:") {
		r5 = errors.Errorf("invalid hls playlist=%v", body)
	}

	cancel()
}
