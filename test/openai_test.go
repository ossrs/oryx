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
	var target *TranscriptSegment
	segments := querySegments("/terraform/v1/ai/transcript/fix-queue")
	for _, s := range segments {
		if s.TsID == segment.TsID {
			target = &s
			break
		}
	}
	// Maybe already pushed in to overlay queue.
	if target == nil {
		segments = querySegments("/terraform/v1/ai/transcript/overlay-queue")
		for _, s := range segments {
			if s.TsID == segment.TsID {
				target = &s
				break
			}
		}
	}

	// Check target segment.
	if target == nil {
		r0 = errors.Errorf("invalid segments %v", segments)
		return
	}
	if !target.UserClearASR {
		r0 = errors.Errorf("invalid segment %v", target)
	}
	return
}

func TestOpenAI_LiveRoomTextAssistant_RoomToken(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	// Ignore the test case if api secret key not set.
	apiKey, baseUrl := OpenAIConfig()
	if apiKey == "" {
		return
	}

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
		// Live room UUID.
		UUID string `json:"uuid"`
		// Live room title.
		Title string `json:"title"`
		// Live room secret.
		Secret string `json:"secret"`
		// Create time.
		CreatedAt string `json:"created_at"`
		// The room level authentication token, for example, popout application with this token to verify
		// the room, to prevent leaking of the bearer token.
		RoomToken string `json:"roomToken"`

		// Whether enable the AI assistant.
		Assistant bool `json:"assistant"`
		// Whether enable the AI ASR.
		AIASREnabled bool `json:"aiAsrEnabled"`
		// Whether enable the AI processing.
		AIChatEnabled bool `json:"aiChatEnabled"`
		// Whether enable the AI TTS.
		AITTSEnabled bool `json:"aiTtsEnabled"`

		// The AI provider.
		AIProvider string `json:"aiProvider"`
		// The AI secret key.
		AISecretKey string `json:"aiSecretKey"`
		// The AI base URL.
		AIBaseURL string `json:"aiBaseURL"`
		// The AI model name.
		AIChatModel string `json:"aiChatModel"`
	}
	var conf LiveRoomCreateResult
	if err := NewApi().WithAuth(ctx, "/terraform/v1/live/room/create", &struct {
		Title string `json:"title"`
	}{
		Title: roomTitle,
	}, &conf); err != nil {
		r0 = errors.Wrapf(err, "create room title=%v", roomTitle)
		return
	}

	backup := conf
	defer func() {
		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/live/room/remove", &backup, nil)
	}()

	// Enable the AI Assistant of live room.
	conf.Assistant, conf.AISecretKey, conf.AIBaseURL = true, apiKey, baseUrl
	conf.AIProvider, conf.AIChatModel = "openai", "gpt-3.5-turbo"
	if err := NewApi().WithAuth(ctx, "/terraform/v1/live/room/update", &conf, nil); err != nil {
		r0 = errors.Wrapf(err, "enable assistant for room id=%v", conf.UUID)
		return
	}

	// Create a stage.
	type StageResult struct {
		StageID   string `json:"sid"`
		RoomToken string `json:"roomToken"`
		UserID    string `json:"userId"`
		// AI Configurations.
		AIASREnabled bool `json:"aiAsrEnabled"`
	}
	var stage StageResult
	if err := NewApi().NoAuth(ctx, "/terraform/v1/ai-talk/stage/start", &struct {
		RoomUUID  string `json:"room"`
		RoomToken string `json:"roomToken"`
	}{
		RoomUUID: conf.UUID, RoomToken: conf.RoomToken,
	}, &stage); err != nil {
		r0 = errors.Wrapf(err, "create stage for room id=%v", conf.UUID)
		return
	}

	if stage.StageID == "" || stage.UserID == "" || stage.RoomToken == "" {
		r0 = errors.Errorf("invalid stage %v", stage)
		return
	}

	// Start a conversation.
	type ConversationResult struct {
		RequestUUID string `json:"rid"`
	}
	var conversation ConversationResult
	if err := NewApi().NoAuth(ctx, "/terraform/v1/ai-talk/stage/conversation", &struct {
		StageUUID string `json:"sid"`
		RoomUUID  string `json:"room"`
		RoomToken string `json:"roomToken"`
	}{
		StageUUID: stage.StageID, RoomUUID: conf.UUID, RoomToken: conf.RoomToken,
	}, &conversation); err != nil {
		r0 = errors.Wrapf(err, "start conversation for stage id=%v", stage.StageID)
		return
	}

	if conversation.RequestUUID == "" {
		r0 = errors.Errorf("invalid conversation %v", conversation)
		return
	}

	// Start a subscriber to pull messages from server. Start the subscriber before send any messages,
	// otherwise, you may miss some messages.
	type SubscribeResult struct {
		StageID      string `json:"sid"`
		SubscriberID string `json:"spid"`
		Voice        string `json:"voice"`
	}
	var subscribe SubscribeResult
	if err := NewApi().NoAuth(ctx, "/terraform/v1/ai-talk/subscribe/start", &struct {
		RoomUUID  string `json:"room"`
		RoomToken string `json:"roomToken"`
	}{
		RoomUUID: conf.UUID, RoomToken: conf.RoomToken,
	}, &subscribe); err != nil {
		r0 = errors.Wrapf(err, "subscribe stage id=%v", stage.StageID)
		return
	}

	if subscribe.SubscriberID == "" {
		r0 = errors.Errorf("invalid subscribe %v", subscribe)
		return
	}

	// Send a text message to AI.
	type TextMessageResult struct {
		RequestUUID string `json:"rid"`
		ASR         string `json:"asr"`
	}
	var txtMessage TextMessageResult
	if err := NewApi().NoAuth(ctx, "/terraform/v1/ai-talk/stage/upload", &struct {
		RoomUUID    string `json:"room"`
		RoomToken   string `json:"roomToken"`
		StageUUID   string `json:"sid"`
		UserID      string `json:"userId"`
		RequestUUID string `json:"rid"`
		TextMessage string `json:"text"`
	}{
		RoomUUID: conf.UUID, RoomToken: conf.RoomToken, StageUUID: stage.StageID,
		UserID: stage.UserID, RequestUUID: conversation.RequestUUID, TextMessage: "hello",
	}, &txtMessage); err != nil {
		r0 = errors.Wrapf(err, "send message for stage id=%v", stage.StageID)
		return
	}

	if txtMessage.ASR != "hello" {
		r0 = errors.Errorf("invalid message %v", txtMessage)
		return
	}

	// Actually, we can also directly pull the result of AI. But to conver as more API as possible, we wait
	// for stage to be ready.
	for {
		type StageStatusResult struct {
			Finished bool `json:"finished"`
		}
		var status StageStatusResult
		if err := NewApi().NoAuth(ctx, "/terraform/v1/ai-talk/stage/query", &struct {
			RoomUUID    string `json:"room"`
			RoomToken   string `json:"roomToken"`
			StageUUID   string `json:"sid"`
			RequestUUID string `json:"rid"`
		}{
			RoomUUID: conf.UUID, RoomToken: conf.RoomToken, StageUUID: stage.StageID,
			RequestUUID: conversation.RequestUUID,
		}, &status); err != nil {
			r0 = errors.Wrapf(err, "query stage id=%v", stage.StageID)
			return
		}

		if status.Finished {
			break
		}

		select {
		case <-ctx.Done():
			r0 = errors.Errorf("timeout for stage id=%v", stage.StageID)
			return
		case <-time.After(5 * time.Second):
		}
	}

	// Should got messages for subscriber, because stage is finished.
	type StageMessage struct {
		// The message UUID.
		MessageUUID string `json:"mid"`
		// The request UUID.
		RequestUUID string `json:"rid"`
		// The message role, text or audio.
		Role string `json:"role"`
		// The username who send this message.
		Username string `json:"username,omitempty"`

		// For role text.
		// The message content.
		Message string `json:"msg"`

		// For role audio.
		// The audio segment uuid.
		SegmentUUID string `json:"asid"`
		// Whether has audio file.
		HasAudioFile bool `json:"hasAudio"`
	}
	type MessageResult struct {
		// Finished messages.
		Messages []*StageMessage `json:"msgs"`
		// Is there any pending messages.
		Pending bool `json:"pending"`
	}
	var message MessageResult
	if err := NewApi().NoAuth(ctx, "/terraform/v1/ai-talk/subscribe/query", &struct {
		RoomUUID     string `json:"room"`
		RoomToken    string `json:"roomToken"`
		StageUUID    string `json:"sid"`
		SubscriberID string `json:"spid"`
		UserID       string `json:"userId"`
	}{
		RoomUUID: conf.UUID, RoomToken: conf.RoomToken, StageUUID: stage.StageID,
		SubscriberID: subscribe.SubscriberID, UserID: stage.UserID,
	}, &message); err != nil {
		r0 = errors.Wrapf(err, "query subscribe stage id=%v", stage.StageID)
		return
	}

	if message.Pending {
		r0 = errors.Errorf("invalid message %v", message)
		return
	}
	// At least two or more messages, one is user message, one is assistant message.
	if len(message.Messages) <= 1 {
		r0 = errors.Errorf("invalid message %v", message)
		return
	}

	// Remove audio segment file of all messages.
	for _, msg := range message.Messages {
		// Ignore user messages.
		if msg.Role == "user" {
			continue
		}

		// Verify the TTS result.
		if !msg.HasAudioFile {
			r0 = errors.Errorf("invalid message %v", msg)
			return
		}

		var body string
		fileURL := fmt.Sprintf(
			"%v/terraform/v1/ai-talk/subscribe/tts?sid=%v&spid=%v&asid=%v&room=%v&roomToken=%v",
			*endpointHTTP, stage.StageID, subscribe.SubscriberID, msg.SegmentUUID, conf.UUID, conf.RoomToken)
		if err := NewApi().Request(ctx, fileURL, nil, false, &body); err != nil {
			r0 = err
			return
		}
		if body == "" {
			r0 = errors.Errorf("invalid message %v", msg)
			return
		}

		// Cleanup message.
		if err := NewApi().NoAuth(ctx, "/terraform/v1/ai-talk/subscribe/remove", &struct {
			RoomUUID         string `json:"room"`
			RoomToken        string `json:"roomToken"`
			StageUUID        string `json:"sid"`
			SubscriberID     string `json:"spid"`
			AudioSegmentUUID string `json:"asid"`
		}{
			RoomUUID: conf.UUID, RoomToken: conf.RoomToken, StageUUID: stage.StageID,
			SubscriberID: subscribe.SubscriberID, AudioSegmentUUID: msg.SegmentUUID,
		}, nil); err != nil {
			r0 = errors.Wrapf(err, "remove message %v", message)
			return
		}
	}

	// Update user information.
	type StageUserResult struct {
		// The user UUID.
		UserID string `json:"userId"`
		// The username.
		Username string `json:"username,omitempty"`

		// The language the host use.
		Language string `json:"language,omitempty"`
		// The AI welcome voice url, binding to the language of user.
		Voice string `json:"voice,omitempty"`
	}
	var stageUser StageUserResult
	if err := NewApi().NoAuth(ctx, "/terraform/v1/ai-talk/user/update", &struct {
		RoomUUID     string `json:"room"`
		RoomToken    string `json:"roomToken"`
		StageUUID    string `json:"sid"`
		UserID       string `json:"userId"`
		UserLanguage string `json:"lang"`
	}{
		RoomUUID: conf.UUID, RoomToken: conf.RoomToken, StageUUID: stage.StageID,
		UserID: stage.UserID, UserLanguage: "zh",
	}, &stageUser); err != nil {
		r0 = errors.Wrapf(err, "query user for stage id=%v", stage.StageID)
		return
	}

	if stageUser.Language != "zh" {
		r0 = errors.Errorf("invalid stage user %v", stageUser)
		return
	}
}

func TestOpenAI_LiveRoomTextAssistant_BearerToken(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	// Ignore the test case if api secret key not set.
	apiKey, baseUrl := OpenAIConfig()
	if apiKey == "" {
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
		// Live room title.
		Title string `json:"title"`
		// Live room secret.
		Secret string `json:"secret"`
		// Create time.
		CreatedAt string `json:"created_at"`
		// The room level authentication token, for example, popout application with this token to verify
		// the room, to prevent leaking of the bearer token.
		RoomToken string `json:"roomToken"`

		// Whether enable the AI assistant.
		Assistant bool `json:"assistant"`
		// Whether enable the AI ASR.
		AIASREnabled bool `json:"aiAsrEnabled"`
		// Whether enable the AI processing.
		AIChatEnabled bool `json:"aiChatEnabled"`
		// Whether enable the AI TTS.
		AITTSEnabled bool `json:"aiTtsEnabled"`

		// The AI provider.
		AIProvider string `json:"aiProvider"`
		// The AI secret key.
		AISecretKey string `json:"aiSecretKey"`
		// The AI base URL.
		AIBaseURL string `json:"aiBaseURL"`
		// The AI model name.
		AIChatModel string `json:"aiChatModel"`
	}
	var conf LiveRoomCreateResult
	if err := NewApi().WithAuth(ctx, "/terraform/v1/live/room/create", &struct {
		Title string `json:"title"`
	}{
		Title: roomTitle,
	}, &conf); err != nil {
		r0 = errors.Wrapf(err, "create room title=%v", roomTitle)
		return
	}

	backup := conf
	defer func() {
		// The ctx has already been cancelled by test case, which will cause the request failed.
		ctx := context.Background()
		NewApi().WithAuth(ctx, "/terraform/v1/live/room/remove", &backup, nil)
	}()

	// Enable the AI Assistant of live room.
	conf.Assistant, conf.AIProvider, conf.AISecretKey, conf.AIBaseURL, conf.AIChatModel = true, "openai", apiKey, baseUrl, "gpt-3.5-turbo"
	if err := NewApi().WithAuth(ctx, "/terraform/v1/live/room/update", &conf, nil); err != nil {
		r0 = errors.Wrapf(err, "enable assistant for room id=%v", conf.UUID)
		return
	}

	// Create a stage.
	type StageResult struct {
		StageID   string `json:"sid"`
		RoomToken string `json:"roomToken"`
		UserID    string `json:"userId"`
		// AI Configurations.
		AIASREnabled bool `json:"aiAsrEnabled"`
	}
	var stage StageResult
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai-talk/stage/start", &struct {
		RoomUUID string `json:"room"`
	}{
		RoomUUID: conf.UUID,
	}, &stage); err != nil {
		r0 = errors.Wrapf(err, "create stage for room id=%v", conf.UUID)
		return
	}

	if stage.StageID == "" || stage.UserID == "" || stage.RoomToken == "" {
		r0 = errors.Errorf("invalid stage %v", stage)
		return
	}

	// Start a conversation.
	type ConversationResult struct {
		RequestUUID string `json:"rid"`
	}
	var conversation ConversationResult
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai-talk/stage/conversation", &struct {
		StageUUID string `json:"sid"`
		RoomUUID  string `json:"room"`
	}{
		StageUUID: stage.StageID, RoomUUID: conf.UUID,
	}, &conversation); err != nil {
		r0 = errors.Wrapf(err, "start conversation for stage id=%v", stage.StageID)
		return
	}

	if conversation.RequestUUID == "" {
		r0 = errors.Errorf("invalid conversation %v", conversation)
		return
	}

	// Start a subscriber to pull messages from server. Start the subscriber before send any messages,
	// otherwise, you may miss some messages.
	type SubscribeResult struct {
		StageID      string `json:"sid"`
		SubscriberID string `json:"spid"`
		Voice        string `json:"voice"`
	}
	var subscribe SubscribeResult
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai-talk/subscribe/start", &struct {
		RoomUUID string `json:"room"`
	}{
		RoomUUID: conf.UUID,
	}, &subscribe); err != nil {
		r0 = errors.Wrapf(err, "subscribe stage id=%v", stage.StageID)
		return
	}

	if subscribe.SubscriberID == "" {
		r0 = errors.Errorf("invalid subscribe %v", subscribe)
		return
	}

	// Send a text message to AI.
	type TextMessageResult struct {
		RequestUUID string `json:"rid"`
		ASR         string `json:"asr"`
	}
	var txtMessage TextMessageResult
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai-talk/stage/upload", &struct {
		RoomUUID    string `json:"room"`
		StageUUID   string `json:"sid"`
		UserID      string `json:"userId"`
		RequestUUID string `json:"rid"`
		TextMessage string `json:"text"`
	}{
		RoomUUID: conf.UUID, StageUUID: stage.StageID,
		UserID: stage.UserID, RequestUUID: conversation.RequestUUID, TextMessage: "hello",
	}, &txtMessage); err != nil {
		r0 = errors.Wrapf(err, "send message for stage id=%v", stage.StageID)
		return
	}

	if txtMessage.ASR != "hello" {
		r0 = errors.Errorf("invalid message %v", txtMessage)
		return
	}

	// Actually, we can also directly pull the result of AI. But to conver as more API as possible, we wait
	// for stage to be ready.
	for {
		type StageStatusResult struct {
			Finished bool `json:"finished"`
		}
		var status StageStatusResult
		if err := NewApi().WithAuth(ctx, "/terraform/v1/ai-talk/stage/query", &struct {
			RoomUUID    string `json:"room"`
			StageUUID   string `json:"sid"`
			RequestUUID string `json:"rid"`
		}{
			RoomUUID: conf.UUID, StageUUID: stage.StageID,
			RequestUUID: conversation.RequestUUID,
		}, &status); err != nil {
			r0 = errors.Wrapf(err, "query stage id=%v", stage.StageID)
			return
		}

		if status.Finished {
			break
		}

		select {
		case <-ctx.Done():
			r0 = errors.Errorf("timeout for stage id=%v", stage.StageID)
			return
		case <-time.After(5 * time.Second):
		}
	}

	// Should got messages for subscriber, because stage is finished.
	type StageMessage struct {
		// The message UUID.
		MessageUUID string `json:"mid"`
		// The request UUID.
		RequestUUID string `json:"rid"`
		// The message role, text or audio.
		Role string `json:"role"`
		// The username who send this message.
		Username string `json:"username,omitempty"`

		// For role text.
		// The message content.
		Message string `json:"msg"`

		// For role audio.
		// The audio segment uuid.
		SegmentUUID string `json:"asid"`
		// Whether has audio file.
		HasAudioFile bool `json:"hasAudio"`
	}
	type MessageResult struct {
		// Finished messages.
		Messages []*StageMessage `json:"msgs"`
		// Is there any pending messages.
		Pending bool `json:"pending"`
	}
	var message MessageResult
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai-talk/subscribe/query", &struct {
		RoomUUID     string `json:"room"`
		StageUUID    string `json:"sid"`
		SubscriberID string `json:"spid"`
		UserID       string `json:"userId"`
	}{
		RoomUUID: conf.UUID, StageUUID: stage.StageID,
		SubscriberID: subscribe.SubscriberID, UserID: stage.UserID,
	}, &message); err != nil {
		r0 = errors.Wrapf(err, "query subscribe stage id=%v", stage.StageID)
		return
	}

	if message.Pending {
		r0 = errors.Errorf("invalid message %v", message)
		return
	}
	// At least two or more messages, one is user message, one is assistant message.
	if len(message.Messages) <= 1 {
		r0 = errors.Errorf("invalid message %v", message)
		return
	}

	// Remove audio segment file of all messages.
	for _, msg := range message.Messages {
		// Ignore user messages.
		if msg.Role == "user" {
			continue
		}

		// Verify the TTS result.
		if !msg.HasAudioFile {
			r0 = errors.Errorf("invalid message %v", msg)
			return
		}

		var body string
		fileURL := fmt.Sprintf(
			"%v/terraform/v1/ai-talk/subscribe/tts?sid=%v&spid=%v&asid=%v&room=%v&roomToken=%v",
			*endpointHTTP, stage.StageID, subscribe.SubscriberID, msg.SegmentUUID, conf.UUID, conf.RoomToken)
		if err := NewApi().Request(ctx, fileURL, nil, false, &body); err != nil {
			r0 = err
			return
		}
		if body == "" {
			r0 = errors.Errorf("invalid message %v", msg)
			return
		}

		// Cleanup segment.
		if err := NewApi().WithAuth(ctx, "/terraform/v1/ai-talk/subscribe/remove", &struct {
			RoomUUID         string `json:"room"`
			StageUUID        string `json:"sid"`
			SubscriberID     string `json:"spid"`
			AudioSegmentUUID string `json:"asid"`
		}{
			RoomUUID: conf.UUID, StageUUID: stage.StageID,
			SubscriberID: subscribe.SubscriberID, AudioSegmentUUID: msg.SegmentUUID,
		}, nil); err != nil {
			r0 = errors.Wrapf(err, "remove message %v", message)
			return
		}
	}

	// Update user information.
	type StageUserResult struct {
		// The user UUID.
		UserID string `json:"userId"`
		// The username.
		Username string `json:"username,omitempty"`

		// The language the host use.
		Language string `json:"language,omitempty"`
		// The AI welcome voice url, binding to the language of user.
		Voice string `json:"voice,omitempty"`
	}
	var stageUser StageUserResult
	if err := NewApi().WithAuth(ctx, "/terraform/v1/ai-talk/user/update", &struct {
		RoomUUID     string `json:"room"`
		RoomToken    string `json:"roomToken"`
		StageUUID    string `json:"sid"`
		UserID       string `json:"userId"`
		UserLanguage string `json:"lang"`
	}{
		RoomUUID: conf.UUID, RoomToken: conf.RoomToken, StageUUID: stage.StageID,
		UserID: stage.UserID, UserLanguage: "zh",
	}, &stageUser); err != nil {
		r0 = errors.Wrapf(err, "query user for stage id=%v", stage.StageID)
		return
	}

	if stageUser.Language != "zh" {
		r0 = errors.Errorf("invalid stage user %v", stageUser)
		return
	}
}
