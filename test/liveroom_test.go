package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
)

func TestMedia_WithStream_LiveRoomCreateQueryRemove(t *testing.T) {
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

	roomTitle := fmt.Sprintf("room-%v-%v", os.Getpid(), rand.Int())
	logger.Tf(ctx, "Test for room title %v", roomTitle)

	type LiveRoomCreateResult struct {
		UUID string `json:"uuid"`
	}
	var roomCreated LiveRoomCreateResult
	if err := NewApi().WithAuth(ctx, "/terraform/v1/live/room/create", &struct {
		Title string `json:"title"`
	}{
		Title: roomTitle,
	}, &roomCreated); err != nil {
		r0 = errors.Wrapf(err, "create room title=%v", roomTitle)
		return
	}

	defer func() {
		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/live/room/remove", &roomCreated, nil)
	}()

	type LiveRoomQueryResult struct {
		// Live room UUID.
		UUID string `json:"uuid"`
		// Live room title.
		Title string `json:"title"`
		// Live room secret.
		Secret string `json:"secret"`
		// Create time.
		CreatedAt string `json:"created_at"`
	}
	var roomQuery LiveRoomQueryResult
	if err := NewApi().WithAuth(ctx, "/terraform/v1/live/room/query", &roomCreated, &roomQuery); err != nil {
		r0 = errors.Wrapf(err, "query room uuid=%v, title=%v", roomCreated.UUID, roomTitle)
		return
	}

	var allRooms []LiveRoomQueryResult
	if err := NewApi().WithAuth(ctx, "/terraform/v1/live/room/list", nil, &struct {
		Rooms *[]LiveRoomQueryResult `json:"rooms"`
	}{
		Rooms: &allRooms,
	}); err != nil {
		r0 = errors.Wrapf(err, "list rooms, uuid=%v, title=%v", roomCreated.UUID, roomTitle)
		return
	}
	var found bool
	for _, v := range allRooms {
		if v.UUID == roomCreated.UUID {
			found = true
			break
		}
	}
	if !found {
		r0 = errors.Errorf("room not found in list, uuid=%v, title=%v", roomCreated.UUID, roomTitle)
		return
	}
}

func TestMedia_WithStream_LiveRoomPublishStream(t *testing.T) {
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

	roomTitle := fmt.Sprintf("room-%v-%v", os.Getpid(), rand.Int())
	logger.Tf(ctx, "Test for room title %v", roomTitle)

	type LiveRoomCreateResult struct {
		// Live room UUID.
		UUID string `json:"uuid"`
		// The stream name, should never use roomUUID because it's secret.
		StreamName string `json:"stream"`
		// Live room title.
		Title string `json:"title"`
		// Live room secret.
		Secret string `json:"secret"`
		// Create time.
		CreatedAt string `json:"created_at"`
	}
	var liveRoom LiveRoomCreateResult
	if err := NewApi().WithAuth(ctx, "/terraform/v1/live/room/create", &struct {
		Title string `json:"title"`
	}{
		Title: roomTitle,
	}, &liveRoom); err != nil {
		r0 = errors.Wrapf(err, "create room title=%v", roomTitle)
		return
	}

	defer func() {
		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/live/room/remove", &liveRoom, nil)
	}()

	var wg sync.WaitGroup
	defer wg.Wait()

	// Start FFmpeg to publish stream.
	streamID := liveRoom.StreamName
	streamURL := fmt.Sprintf("%v/live/%v?secret=%v", *endpointRTMP, streamID, liveRoom.Secret)
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

func TestMedia_WithStream_LiveRoomPublishInvalidStream(t *testing.T) {
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

	roomTitle := fmt.Sprintf("room-%v-%v", os.Getpid(), rand.Int())
	logger.Tf(ctx, "Test for room title %v", roomTitle)

	type LiveRoomCreateResult struct {
		// Live room UUID.
		UUID string `json:"uuid"`
		// The stream name, should never use roomUUID because it's secret.
		StreamName string `json:"stream"`
		// Live room title.
		Title string `json:"title"`
		// Live room secret.
		Secret string `json:"secret"`
		// Create time.
		CreatedAt string `json:"created_at"`
	}
	var liveRoom LiveRoomCreateResult
	if err := NewApi().WithAuth(ctx, "/terraform/v1/live/room/create", &struct {
		Title string `json:"title"`
	}{
		Title: roomTitle,
	}, &liveRoom); err != nil {
		r0 = errors.Wrapf(err, "create room title=%v", roomTitle)
		return
	}

	defer func() {
		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/live/room/remove", &liveRoom, nil)
	}()

	var wg sync.WaitGroup
	defer wg.Wait()

	// Start FFmpeg to publish stream.
	// Use a invalid random stream ID, which should be failed.
	streamID := fmt.Sprintf("stream-%v-%v", os.Getpid(), rand.Int())
	streamURL := fmt.Sprintf("%v/live/%v?secret=%v", *endpointRTMP, streamID, liveRoom.Secret)
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

func TestMedia_WithStream_LiveRoomPublishInvalidSecret(t *testing.T) {
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

	roomTitle := fmt.Sprintf("room-%v-%v", os.Getpid(), rand.Int())
	logger.Tf(ctx, "Test for room title %v", roomTitle)

	type LiveRoomCreateResult struct {
		// Live room UUID.
		UUID string `json:"uuid"`
		// The stream name, should never use roomUUID because it's secret.
		StreamName string `json:"stream"`
		// Live room title.
		Title string `json:"title"`
		// Live room secret.
		Secret string `json:"secret"`
		// Create time.
		CreatedAt string `json:"created_at"`
	}
	var liveRoom LiveRoomCreateResult
	if err := NewApi().WithAuth(ctx, "/terraform/v1/live/room/create", &struct {
		Title string `json:"title"`
	}{
		Title: roomTitle,
	}, &liveRoom); err != nil {
		r0 = errors.Wrapf(err, "create room title=%v", roomTitle)
		return
	}

	defer func() {
		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/live/room/remove", &liveRoom, nil)
	}()

	var wg sync.WaitGroup
	defer wg.Wait()

	// Start FFmpeg to publish stream.
	// Use a invalid random stream ID, which should be failed.
	streamID := liveRoom.StreamName
	streamURL := fmt.Sprintf("%v/live/%v?secret=%v", *endpointRTMP, streamID, uuid.NewString())
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
