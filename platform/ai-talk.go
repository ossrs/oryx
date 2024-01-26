package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	errors_std "errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/ossrs/go-oryx-lib/errors"
	ohttp "github.com/ossrs/go-oryx-lib/http"
	"github.com/ossrs/go-oryx-lib/logger"
	"github.com/sashabaranov/go-openai"
	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
)

var talkServer *TalkServer
var workDir string

type ASRResult struct {
	Text     string
	Duration time.Duration
}

type ASRService interface {
	RequestASR(ctx context.Context, filepath, language, prompt string) (*ASRResult, error)
}

type TTSService interface {
	RequestTTS(ctx context.Context, buildFilepath func(ext string) string, text string) error
}

type openaiASRService struct {
	conf openai.ClientConfig
	// The callback before start ASR request.
	onBeforeRequest func()
}

func NewOpenAIASRService(conf openai.ClientConfig, opts ...func(service *openaiASRService)) ASRService {
	v := &openaiASRService{conf: conf}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

func (v *openaiASRService) RequestASR(ctx context.Context, inputFile, language, prompt string) (*ASRResult, error) {
	outputFile := fmt.Sprintf("%v.m4a", inputFile)
	defer os.Remove(outputFile)

	// Transcode input audio in opus or aac, to aac in m4a format.
	if true {
		err := exec.CommandContext(ctx, "ffmpeg",
			"-i", inputFile,
			"-vn", "-c:a", "aac", "-ac", "1", "-ar", "16000", "-ab", "30k",
			outputFile,
		).Run()

		if err != nil {
			return nil, errors.Errorf("Error converting the file")
		}
		logger.Tf(ctx, "Convert audio %v to %v ok", inputFile, outputFile)
	}

	duration, _, err := ffprobeAudio(ctx, outputFile)
	if err != nil {
		return nil, errors.Wrapf(err, "ffprobe")
	}

	if v.onBeforeRequest != nil {
		v.onBeforeRequest()
	}

	// Request ASR.
	client := openai.NewClientWithConfig(v.conf)
	resp, err := client.CreateTranscription(
		ctx,
		openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: outputFile,
			Format:   openai.AudioResponseFormatJSON,
			Language: language,
			Prompt:   prompt,
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "asr")
	}

	return &ASRResult{Text: resp.Text, Duration: time.Duration(duration * float64(time.Second))}, nil
}

func ffprobeAudio(ctx context.Context, filename string) (duration float64, bitrate int, err error) {
	args := []string{
		"-show_error", "-show_private_data", "-v", "quiet", "-find_stream_info", "-print_format", "json",
		"-show_format",
	}
	args = append(args, "-i", filename)

	stdout, err := exec.CommandContext(ctx, "ffprobe", args...).Output()
	if err != nil {
		err = errors.Wrapf(err, "probe %v", filename)
		return
	}

	type VLiveFileFormat struct {
		Starttime string `json:"start_time"`
		Duration  string `json:"duration"`
		Bitrate   string `json:"bit_rate"`
		Streams   int32  `json:"nb_streams"`
		Score     int32  `json:"probe_score"`
		HasVideo  bool   `json:"has_video"`
		HasAudio  bool   `json:"has_audio"`
	}

	format := struct {
		Format VLiveFileFormat `json:"format"`
	}{}
	if err = json.Unmarshal([]byte(stdout), &format); err != nil {
		err = errors.Wrapf(err, "parse format %v", stdout)
		return
	}

	var fv float64
	if fv, err = strconv.ParseFloat(format.Format.Duration, 64); err != nil {
		err = errors.Wrapf(err, "parse duration %v", format.Format.Duration)
		return
	} else {
		duration = fv
	}

	var iv int64
	if iv, err = strconv.ParseInt(format.Format.Bitrate, 10, 64); err != nil {
		err = errors.Wrapf(err, "parse bitrate %v", format.Format.Bitrate)
		return
	} else {
		bitrate = int(iv)
	}

	logger.Tf(ctx, "FFprobe input=%v, duration=%v, bitrate=%v", filename, duration, bitrate)
	return
}

type openaiChatService struct {
	// The AI configuration.
	conf openai.ClientConfig
	// The callback for the first response.
	onFirstResponse func(ctx context.Context, text string)
}

func (v *openaiChatService) RequestChat(ctx context.Context, rid string, stage *Stage, robot *Robot) error {
	if stage.previousUser != "" && stage.previousAssitant != "" {
		stage.histories = append(stage.histories, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: stage.previousUser,
		}, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: stage.previousAssitant,
		})

		for len(stage.histories) > robot.chatWindow*2 {
			stage.histories = stage.histories[1:]
		}
	}

	stage.previousUser = stage.previousAsrText
	stage.previousAssitant = ""

	system := robot.prompt
	system += fmt.Sprintf(" Keep your reply neat, limiting the reply to %v words.", robot.replyLimit)
	logger.Tf(ctx, "AI system prompt: %v", system)
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: system},
	}

	messages = append(messages, stage.histories...)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: stage.previousAsrText,
	})

	model := robot.chatModel
	maxTokens := 1024
	temperature := float32(0.9)
	logger.Tf(ctx, "robot=%v(%v), OPENAI_PROXY: %v, AIT_CHAT_MODEL: %v, AIT_MAX_TOKENS: %v, AIT_TEMPERATURE: %v, window=%v, histories=%v",
		robot.uuid, robot.label, v.conf.BaseURL, model, maxTokens, temperature, robot.chatWindow, len(stage.histories))

	client := openai.NewClientWithConfig(v.conf)
	gptChatStream, err := client.CreateChatCompletionStream(
		ctx, openai.ChatCompletionRequest{
			Model:       model,
			Messages:    messages,
			Stream:      true,
			Temperature: temperature,
			MaxTokens:   maxTokens,
		},
	)
	if err != nil {
		return errors.Wrapf(err, "create chat")
	}

	// Never wait for any response.
	go func() {
		defer gptChatStream.Close()
		if err := v.handle(ctx, stage, robot, rid, gptChatStream); err != nil {
			logger.Ef(ctx, "Handle stream failed, err %+v", err)
		}
	}()

	return nil
}

func (v *openaiChatService) handle(ctx context.Context, stage *Stage, robot *Robot, rid string, gptChatStream *openai.ChatCompletionStream) error {
	stage.generating = true
	defer func() {
		stage.generating = false
	}()

	filterAIResponse := func(response *openai.ChatCompletionStreamResponse, err error) (bool, string, error) {
		finished := errors_std.Is(err, io.EOF)
		if err != nil && !finished {
			return finished, "", errors.Wrapf(err, "recv chat")
		}

		if len(response.Choices) == 0 {
			return finished, "", nil
		}

		choice := response.Choices[0]
		dc := choice.Delta.Content
		if dc == "" {
			return finished, "", nil
		}

		filteredStencese := strings.ReplaceAll(dc, "\n\n", "\n")
		filteredStencese = strings.ReplaceAll(filteredStencese, "\n", " ")
		logger.Tf(ctx, "AI response: text=%v", dc)

		return finished, filteredStencese, nil
	}

	gotNewSentence := func(sentence string, firstSentense bool) bool {
		newSentence := false

		// Ignore empty.
		if sentence == "" {
			return newSentence
		}

		// Any ASCII character to split sentence.
		if strings.ContainsAny(sentence, ",.?!\n") {
			newSentence = true
		}

		// Any Chinese character to split sentence.
		if strings.ContainsRune(sentence, '。') ||
			strings.ContainsRune(sentence, '？') ||
			strings.ContainsRune(sentence, '！') ||
			strings.ContainsRune(sentence, '，') {
			newSentence = true
		}

		isEnglish := func(s string) bool {
			for _, r := range s {
				if r > unicode.MaxASCII {
					return false
				}
			}
			return true
		}

		// Determine whether new sentence by length.
		if isEnglish(sentence) {
			maxWords, minWords := 30, 3
			if !firstSentense {
				maxWords, minWords = 50, 5
			}

			if nn := strings.Count(sentence, " "); nn >= maxWords {
				newSentence = true
			} else if nn < minWords {
				newSentence = false
			}
		} else {
			maxWords, minWords := 50, 3
			if !firstSentense {
				maxWords, minWords = 100, 5
			}

			if nn := utf8.RuneCount([]byte(sentence)); nn >= maxWords {
				newSentence = true
			} else if nn < minWords {
				newSentence = false
			}
		}

		return newSentence
	}

	var previous *AnswerSegment
	commitAISentence := func(sentence string, firstSentense bool) {
		if firstSentense {
			if robot.prefix != "" {
				sentence = fmt.Sprintf("%v %v", robot.prefix, sentence)
			}
			if v.onFirstResponse != nil {
				v.onFirstResponse(ctx, sentence)
			}
		}

		current := NewAnswerSegment(func(segment *AnswerSegment) {
			segment.rid = rid
			segment.text = sentence
			segment.first = firstSentense
		})
		stage.ttsWorker.SubmitSegment(ctx, stage, current, previous)
		previous = current
		return
	}

	var sentence string
	var finished bool
	firstSentense := true
	for !finished && ctx.Err() == nil {
		response, err := gptChatStream.Recv()
		if tvFinished, tvSentense, err := filterAIResponse(&response, err); err != nil {
			return errors.Wrapf(err, "filter")
		} else {
			finished, sentence = tvFinished, sentence+tvSentense
		}

		newSentence := gotNewSentence(sentence, firstSentense)
		if !finished && !newSentence {
			continue
		}

		// Use the sentence for prompt and logging.
		stage.previousAssitant += sentence + " "
		// We utilize user ASR and AI responses as prompts for the subsequent ASR, given that this is
		// a chat-based scenario where the user converses with the AI, and the following audio should pertain to both user and AI text.
		stage.previousAsrText += " " + sentence
		// Commit the sentense to TTS worker and callbacks.
		commitAISentence(sentence, firstSentense)
		// Reset the sentence, because we have committed it.
		sentence, firstSentense = "", false
	}

	return nil
}

type openaiTTSService struct {
	conf openai.ClientConfig
}

func NewOpenAITTSService(conf openai.ClientConfig) TTSService {
	return &openaiTTSService{conf: conf}
}

func (v *openaiTTSService) RequestTTS(ctx context.Context, buildFilepath func(ext string) string, text string) error {
	ttsFile := buildFilepath("aac")

	client := openai.NewClientWithConfig(v.conf)
	resp, err := client.CreateSpeech(ctx, openai.CreateSpeechRequest{
		Model:          openai.TTSModel1,
		Input:          text,
		Voice:          openai.VoiceNova,
		ResponseFormat: openai.SpeechResponseFormatAac,
	})
	if err != nil {
		return errors.Wrapf(err, "create speech")
	}
	defer resp.Close()

	out, err := os.Create(ttsFile)
	if err != nil {
		return errors.Errorf("Unable to create the file %v for writing", ttsFile)
	}
	defer out.Close()

	if _, err = io.Copy(out, resp); err != nil {
		return errors.Errorf("Error writing the file")
	}

	return nil
}

// The Robot is a robot that user can talk with.
type Robot struct {
	// The robot uuid.
	uuid string
	// The robot label.
	label string
	// The robot prompt.
	prompt string
	// The robot ASR language.
	asrLanguage string
	// The prefix for TTS for the first sentence if too short.
	prefix string
	// The welcome voice url.
	voice string
	// Reply words limit.
	replyLimit int
	// AI Chat model.
	chatModel string
	// AI Chat message window.
	chatWindow int
}

func (v Robot) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("uuid:%v,label:%v,asr:%v", v.uuid, v.label, v.asrLanguage))
	if v.prefix != "" {
		sb.WriteString(fmt.Sprintf(",prefix:%v", v.prefix))
	}
	sb.WriteString(fmt.Sprintf(",voice=%v,limit=%v,model=%v,window=%v,prompt:%v",
		v.voice, v.replyLimit, v.chatModel, v.chatWindow, v.prompt))
	return sb.String()
}

// The StageMessage is a message from user or AI.
type StageMessage struct {
	// The message UUID.
	MessageUUID string `json:"mid"`
	// The request UUID.
	RequestUUID string `json:"rid"`
	// The message role, text or audio.
	Role string `json:"role"`

	// For role text.
	// The message content.
	Message string `json:"msg"`

	// For role audio.
	// The audio segment uuid.
	SegmentUUID string `json:"asid"`
	// The audio tts file for audio message.
	audioFile string

	// Whether flushed to client.
	flushed bool
}

func (v *StageMessage) Close() error {
	if v.audioFile != "" {
		if _, err := os.Stat(v.audioFile); err == nil {
			_ = os.Remove(v.audioFile)
		}
	}
	return nil
}

// The StagePopout is a popout from a stage.
type StagePopout struct {
	// StagePopout UUID.
	spid string
	// The logging context, to write all logs in one context for a sage.
	loggingCtx context.Context
	// The messages from user ASR and AI responses.
	messages []*StageMessage

	// The owner room, never changes.
	room *SrsLiveRoom
	// The owner stage, never changes.
	stage *Stage
}

func NewStagePopout(opts ...func(*StagePopout)) *StagePopout {
	v := &StagePopout{
		// Create new UUID.
		spid: uuid.NewString(),
	}

	for _, opt := range opts {
		opt(v)
	}
	return v
}

func (v *StagePopout) Close() error {
	for _, message := range v.messages {
		_ = message.Close()
	}
	return nil
}

func (v *StagePopout) addUserTextMessage(rid, msg string) {
	v.messages = append(v.messages, &StageMessage{
		MessageUUID: uuid.NewString(), RequestUUID: rid, Role: "user", Message: msg,
	})
}

func (v *StagePopout) addRobotAudioMessage(ctx context.Context, rid, msg, segmentUUID, ttsFile string) error {
	// Build a new copy file of ttsFile.
	ttsExt := path.Ext(ttsFile)
	copyFile := fmt.Sprintf("%v-copy-%v%v", ttsFile[:len(ttsFile)-len(ttsExt)], v.spid, ttsExt)

	// Copy the ttsFile to copyFile.
	if err := func() error {
		src, err := os.Open(ttsFile)
		if err != nil {
			return errors.Errorf("open %v for reading", ttsFile)
		}
		defer src.Close()

		dst, err := os.OpenFile(copyFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return errors.Errorf("open %v for writing", copyFile)
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return errors.Errorf("copy file content")
		}

		logger.Tf(ctx, "AITalk: Copy %v to %v ok, room=%v, sid=%v, spid=%v",
			ttsFile, copyFile, v.room.UUID, v.stage.sid, v.spid)
		return nil
	}(); err != nil {
		// TODO: FIXME: Cleanup message if error.
		return errors.Wrapf(err, "copy %v to %v", ttsFile, copyFile)
	}

	message := &StageMessage{
		MessageUUID: uuid.NewString(), RequestUUID: rid, SegmentUUID: segmentUUID,
		Role: "robot", Message: msg, audioFile: copyFile,
	}

	// Always close message if timeout.
	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(30 * time.Second):
			message.Close()
		}
	}()

	v.messages = append(v.messages, message)

	return nil
}

func (v *StagePopout) flushMessages() []*StageMessage {
	messages := []*StageMessage{}
	for _, message := range v.messages {
		if !message.flushed {
			messages = append(messages, message)
			message.flushed = true
		}
	}
	return messages
}

func (v *StagePopout) queryAudioFile(asid string) string {
	for _, message := range v.messages {
		if message.SegmentUUID == asid {
			return message.audioFile
		}
	}
	return ""
}

func (v *StagePopout) removeMessage(asid string) error {
	for i, message := range v.messages {
		if message.SegmentUUID == asid {
			v.messages = append(v.messages[:i], v.messages[i+1:]...)
			return message.Close()
		}
	}
	return nil
}

// The Stage is a stage of conversation, when user click start with a scenario,
// we will create a stage object.
type Stage struct {
	// Stage UUID
	sid string
	// Last update of stage.
	update time.Time
	// The TTS worker for this stage.
	ttsWorker *TTSWorker
	// The logging context, to write all logs in one context for a sage.
	loggingCtx context.Context
	// Previous ASR text, to use as prompt for next ASR.
	previousAsrText string
	// Previous chat text, to use as prompt for next chat.
	previousUser, previousAssitant string
	// The chat history, to use as prompt for next chat.
	histories []openai.ChatCompletionMessage
	// Whether the stage is generating more sentences.
	generating bool

	// For time cost statistic.
	lastSentence time.Time
	// The time for last upload audio.
	lastUploadAudio time.Time
	// The time for last extract audio for ASR.
	lastExtractAudio time.Time
	// The time for last request ASR result.
	lastRequestASR time.Time
	// The last request ASR text.
	lastRequestAsrText string
	// The ASR duration of audio file.
	lastAsrDuration time.Duration
	// The time for last request Chat result, the first segment.
	lastRequestChat time.Time
	// The last response text of robot.
	lastRobotFirstText string
	// The time for last request TTS result, the first segment.
	lastRequestTTS time.Time
	// The time for last download the TTS result, the first segment.
	lastDownloadAudio time.Time

	// The robot created for this stage.
	robot *Robot
	// The AI configuration.
	aiConfig openai.ClientConfig
	// The room it belongs to.
	room *SrsLiveRoom
	// All the popouts binding to this stage.
	popouts []*StagePopout
}

func NewStage(opts ...func(*Stage)) *Stage {
	v := &Stage{
		// Create new UUID.
		sid: uuid.NewString(),
		// Update time.
		update: time.Now(),
		// The TTS worker.
		ttsWorker: NewTTSWorker(),
	}

	for _, opt := range opts {
		opt(v)
	}
	return v
}

func (v *Stage) Close() error {
	for _, popout := range v.popouts {
		popout.Close()
	}
	return v.ttsWorker.Close()
}

func (v *Stage) Expired() bool {
	return time.Since(v.update) > 600*time.Second
}

func (v *Stage) KeepAlive() {
	v.update = time.Now()
}

func (v *Stage) total() float64 {
	if v.lastDownloadAudio.After(v.lastSentence) {
		return float64(v.lastDownloadAudio.Sub(v.lastSentence)) / float64(time.Second)
	}
	return 0
}

func (v *Stage) upload() float64 {
	if v.lastUploadAudio.After(v.lastSentence) {
		return float64(v.lastUploadAudio.Sub(v.lastSentence)) / float64(time.Second)
	}
	return 0
}

func (v *Stage) exta() float64 {
	if v.lastExtractAudio.After(v.lastUploadAudio) {
		return float64(v.lastExtractAudio.Sub(v.lastUploadAudio)) / float64(time.Second)
	}
	return 0
}

func (v *Stage) asr() float64 {
	if v.lastRequestASR.After(v.lastExtractAudio) {
		return float64(v.lastRequestASR.Sub(v.lastExtractAudio)) / float64(time.Second)
	}
	return 0
}

func (v *Stage) chat() float64 {
	if v.lastRequestChat.After(v.lastRequestASR) {
		return float64(v.lastRequestChat.Sub(v.lastRequestASR)) / float64(time.Second)
	}
	return 0
}

func (v *Stage) tts() float64 {
	if v.lastRequestTTS.After(v.lastRequestChat) {
		return float64(v.lastRequestTTS.Sub(v.lastRequestChat)) / float64(time.Second)
	}
	return 0
}

func (v *Stage) download() float64 {
	if v.lastDownloadAudio.After(v.lastRequestTTS) {
		return float64(v.lastDownloadAudio.Sub(v.lastRequestTTS)) / float64(time.Second)
	}
	return 0
}

func (v *Stage) addPopout(popout *StagePopout) {
	v.popouts = append(v.popouts, popout)
}

func (v *Stage) queryPopout(spid string) *StagePopout {
	for _, popout := range v.popouts {
		if popout.spid == spid {
			return popout
		}
	}
	return nil
}

// The AnswerSegment is a segment of answer, which is a sentence.
type AnswerSegment struct {
	// Request UUID.
	rid string
	// Answer segment UUID.
	asid string
	// The text of this answer segment.
	text string
	// The TTS file path.
	ttsFile string
	// Whether TTS is done, ready to play.
	ready bool
	// Whether TTS is error, failed.
	err error
	// Whether dummy segment, to identify the request is alive.
	dummy bool
	// Signal to remove the TTS file immediately.
	removeSignal chan bool
	// Whether we have logged this segment.
	logged bool
	// Whether the segment is the first response.
	first bool
}

func NewAnswerSegment(opts ...func(segment *AnswerSegment)) *AnswerSegment {
	v := &AnswerSegment{
		// Request UUID.
		rid: uuid.NewString(),
		// Audio Segment UUID.
		asid: uuid.NewString(),
		// Signal to remove the TTS file.
		removeSignal: make(chan bool, 1),
	}

	for _, opt := range opts {
		opt(v)
	}
	return v
}

// The TalkServer is the AI talk server, manage stages.
type TalkServer struct {
	// All stages created by user.
	stages []*Stage
	// The lock to protect fields.
	lock sync.Mutex
}

func NewTalkServer() *TalkServer {
	return &TalkServer{
		stages: []*Stage{},
	}
}

func (v *TalkServer) Close() error {
	return nil
}

func (v *TalkServer) AddStage(stage *Stage) {
	v.lock.Lock()
	defer v.lock.Unlock()

	v.stages = append(v.stages, stage)
}

func (v *TalkServer) RemoveStage(stage *Stage) {
	v.lock.Lock()
	defer v.lock.Unlock()

	for i, s := range v.stages {
		if s.sid == stage.sid {
			v.stages = append(v.stages[:i], v.stages[i+1:]...)
			return
		}
	}
}

func (v *TalkServer) CountStage() int {
	v.lock.Lock()
	defer v.lock.Unlock()

	return len(v.stages)
}

func (v *TalkServer) QueryStage(rid string) *Stage {
	v.lock.Lock()
	defer v.lock.Unlock()

	for _, s := range v.stages {
		if s.sid == rid {
			return s
		}
	}

	return nil
}

// The TTSWorker is a worker to convert answers from text to audio.
type TTSWorker struct {
	segments []*AnswerSegment
	lock     sync.Mutex
	wg       sync.WaitGroup
}

func NewTTSWorker() *TTSWorker {
	return &TTSWorker{
		segments: []*AnswerSegment{},
	}
}

func (v *TTSWorker) Close() error {
	v.wg.Wait()
	return nil
}

func (v *TTSWorker) QuerySegment(rid, asid string) *AnswerSegment {
	v.lock.Lock()
	defer v.lock.Unlock()

	for _, s := range v.segments {
		if s.rid == rid && s.asid == asid {
			return s
		}
	}

	return nil
}

// TODO: FIMXE: Remove it.
func (v *TTSWorker) QueryAnyReadySegment(ctx context.Context, stage *Stage, rid string) *AnswerSegment {
	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
		case <-time.After(100 * time.Millisecond):
		}

		// When there is no segments, and AI is generating the sentence, we need to wait. For example,
		// if the first sentence is very short, maybe we got it quickly, but the second sentence is very
		// long so that the AI need more time to generate it.
		var s *AnswerSegment
		for ctx.Err() == nil && s == nil && stage.generating {
			if s = v.query(rid); s == nil {
				select {
				case <-ctx.Done():
				case <-time.After(100 * time.Millisecond):
				}
			}
		}

		// Try to fetch one again, because maybe there is new segment.
		s = v.query(rid)

		// All segments are consumed, we return nil.
		if s == nil {
			return nil
		}

		// Wait for dummy segment to be removed.
		if s.dummy {
			continue
		}

		// When segment is finished(ready or error), we return it.
		if s.ready || s.err != nil {
			return s
		}
	}

	return nil
}

func (v *TTSWorker) query(rid string) *AnswerSegment {
	v.lock.Lock()
	defer v.lock.Unlock()

	for _, s := range v.segments {
		if s.rid == rid {
			return s
		}
	}

	return nil
}

func (v *TTSWorker) RemoveSegment(asid string) {
	v.lock.Lock()
	defer v.lock.Unlock()

	for i, s := range v.segments {
		if s.asid == asid {
			v.segments = append(v.segments[:i], v.segments[i+1:]...)
			return
		}
	}
}

func (v *TTSWorker) SubmitSegment(ctx context.Context, stage *Stage, segment, previous *AnswerSegment) {
	// Append the sentence to queue.
	func() {
		v.lock.Lock()
		defer v.lock.Unlock()

		v.segments = append(v.segments, segment)
	}()

	// Ignore the dummy sentence.
	if segment.dummy {
		return
	}

	// Now that we have a real sentence, we should remove the dummy sentence.
	if dummy := v.query(segment.rid); dummy != nil && dummy.dummy {
		v.RemoveSegment(dummy.asid)
	}

	// Start a goroutine to do TTS task.
	v.wg.Add(1)
	go func() {
		defer v.wg.Done()

		ttsService := NewOpenAITTSService(stage.aiConfig)
		if err := ttsService.RequestTTS(ctx, func(ext string) string {
			segment.ttsFile = path.Join(workDir,
				fmt.Sprintf("assistant-%v-sentence-%v-tts.%v", segment.rid, segment.asid, ext),
			)
			return segment.ttsFile
		}, segment.text); err != nil {
			segment.err = err
		} else {
			segment.ready = true
			if segment.first {
				stage.lastRequestTTS = time.Now()
			}
			logger.Tf(ctx, "File saved to %v, %v", segment.ttsFile, segment.text)
		}

		// TODO: FIXME: Keep the same order.
		// Because we submit all TTS tasks to OpenAI simultaneously, to keep the answer segments in the
		// same order to subscribers, so we need to wait for previous tasks to finish.
		if previous != nil && !previous.ready {
			for ctx.Err() == nil && previous != nil && !previous.ready {
				select {
				case <-ctx.Done():
				case <-time.After(300 * time.Millisecond):
				}
			}

			// Maybe previous segment is pushing to queue, so wait for a while.
			select {
			case <-ctx.Done():
			case <-time.After(300 * time.Millisecond):
			}
		}

		// Notify all popouts the segment is ready.
		for _, popout := range stage.popouts {
			if err := popout.addRobotAudioMessage(ctx, segment.rid, segment.text, segment.asid, segment.ttsFile); err != nil {
				segment.err = errors.Wrapf(err, "copy to popout")
			}
		}

		// Start a goroutine to remove the sentence.
		v.wg.Add(1)
		go func() {
			defer v.wg.Done()

			select {
			case <-ctx.Done():
			case <-time.After(300 * time.Second):
			case <-segment.removeSignal:
			}

			logger.Tf(ctx, "Remove %v %v", segment.asid, segment.ttsFile)

			stage.ttsWorker.RemoveSegment(segment.asid)

			if segment.ttsFile != "" && os.Getenv("AIT_KEEP_FILES") != "true" {
				if _, err := os.Stat(segment.ttsFile); err == nil {
					os.Remove(segment.ttsFile)
				}
			}
		}()
	}()
}

func handleAITalkService(ctx context.Context, handler *http.ServeMux) error {
	workDir = path.Join(conf.Pwd, "containers/data/ai-talk")
	logger.Tf(ctx, "AI-Talk work dir: %v", workDir)

	ep := "/terraform/v1/ai-talk/stage/start"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var rid string
			if err := ParseBody(ctx, r.Body, &struct {
				Token    *string `json:"token"`
				RoomUUID *string `json:"room"`
			}{
				Token: &token, RoomUUID: &rid,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			var room SrsLiveRoom
			if r0, err := rdb.HGet(ctx, SRS_LIVE_ROOM, rid).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_LIVE_ROOM, rid)
			} else if r0 == "" {
				return errors.Errorf("live room %v not exists", rid)
			} else if err = json.Unmarshal([]byte(r0), &room); err != nil {
				return errors.Wrapf(err, "unmarshal %v %v", rid, r0)
			}

			// Allow to reuse exists stage.
			stage := talkServer.QueryStage(room.StageUUID)
			if stage != nil {
				logger.Tf(ctx, "Stage: Reuse stage sid=%v, all=%v", stage.sid, talkServer.CountStage())
			} else {
				ctx = logger.WithContext(ctx)
				stage = NewStage(func(stage *Stage) {
					stage.loggingCtx = ctx

					// Create robot for the stage, which attach to a special room.
					stage.robot = &Robot{
						uuid: uuid.NewString(), label: "Default", voice: "hello-english.aac",
						prompt: room.AIChatPrompt, asrLanguage: room.AIASRLanguage, replyLimit: room.AIChatMaxWords,
						chatModel: room.AIChatModel, chatWindow: room.AIChatMaxWindow,
					}
					if room.AIASRLanguage == "zh" {
						stage.robot.voice = "hello-chinese.aac"
					}

					// Initialize the AI services.
					stage.aiConfig = openai.DefaultConfig(room.AISecretKey)
					stage.aiConfig.BaseURL = room.AIBaseURL

					// Bind stage to room.
					room.StageUUID = stage.sid
					stage.room = &room
				})

				if b, err := json.Marshal(room); err != nil {
					return errors.Wrapf(err, "marshal room")
				} else if err := rdb.HSet(ctx, SRS_LIVE_ROOM, room.UUID, string(b)).Err(); err != nil {
					return errors.Wrapf(err, "hset %v %v %v", SRS_LIVE_ROOM, room.UUID, string(b))
				}

				talkServer.AddStage(stage)
				logger.Tf(ctx, "Stage: Create new stage sid=%v, all=%v", stage.sid, talkServer.CountStage())

				go func() {
					defer stage.Close()

					for ctx.Err() == nil {
						select {
						case <-ctx.Done():
						case <-time.After(3 * time.Second):
							if stage.Expired() {
								logger.Tf(ctx, "Stage: Remove %v for expired, update=%v",
									stage.sid, stage.update.Format(time.RFC3339))
								talkServer.RemoveStage(stage)
								return
							}
						}
					}
				}()
			}

			type StageRobotResult struct {
				UUID  string `json:"uuid"`
				Label string `json:"label"`
				Voice string `json:"voice"`
			}
			type StageResult struct {
				StageID string           `json:"sid"`
				Robot   StageRobotResult `json:"robot"`
			}
			r0 := &StageResult{
				StageID: stage.sid,
				Robot: StageRobotResult{
					UUID:  stage.robot.uuid,
					Label: stage.robot.label,
					Voice: stage.robot.voice,
				},
			}

			ohttp.WriteData(ctx, w, r, r0)
			logger.Tf(ctx, "srs ai-talk create stage ok")
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai-talk/stage/upload"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var sid, robotUUID, audioBase64Data string
			var userMayInput float64
			if err := ParseBody(ctx, r.Body, &struct {
				Token        *string  `json:"token"`
				StageUUID    *string  `json:"sid"`
				RobotUUID    *string  `json:"robot"`
				UserMayInput *float64 `json:"umi"`
				AudioData    *string  `json:"audio"`
			}{
				Token: &token, StageUUID: &sid, RobotUUID: &robotUUID,
				UserMayInput: &userMayInput, AudioData: &audioBase64Data,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if sid == "" {
				return errors.Errorf("empty sid")
			}
			if robotUUID == "" {
				return errors.Errorf("empty robot")
			}

			stage := talkServer.QueryStage(sid)
			if stage == nil {
				return errors.Errorf("invalid sid %v", sid)
			}

			// Keep alive the stage.
			stage.KeepAlive()
			// For statistic, this is the last sentence from user, we start to count the cost from here.
			stage.lastSentence = time.Now()
			// Switch to the context of stage.
			ctx = stage.loggingCtx

			robot := stage.robot
			if robot == nil {
				return errors.Errorf("invalid robot %v", robotUUID)
			}

			// The rid is the request id, which identify this request, generally a question.
			rid := uuid.NewString()
			inputFile := path.Join(workDir, fmt.Sprintf("assistant-%v-input.audio", rid))
			logger.Tf(ctx, "Stage: Got question sid=%v, umi=%v, robot=%v(%v), rid=%v, input=%v",
				sid, userMayInput, robot.uuid, robot.label, rid, inputFile)
			defer os.Remove(inputFile)
			if err := func() error {
				data, err := base64.StdEncoding.DecodeString(audioBase64Data)
				if err != nil {
					return errors.Errorf("decode base64 from %v", audioBase64Data)
				}

				out, err := os.Create(inputFile)
				if err != nil {
					return errors.Errorf("Unable to create the file for writing")
				}
				defer out.Close()

				nn, err := io.Copy(out, bytes.NewReader([]byte(data)))
				if err != nil {
					return errors.Errorf("Error writing the file")
				}
				logger.Tf(ctx, "File saved to %v, size: %v", inputFile, nn)
				return nil
			}(); err != nil {
				return errors.Wrapf(err, "copy %v", inputFile)
			}
			stage.lastUploadAudio = time.Now()

			// Do ASR, convert to text.
			var asrText string
			asrService := NewOpenAIASRService(stage.aiConfig, func(*openaiASRService) {
				stage.lastExtractAudio = time.Now()
			})
			if resp, err := asrService.RequestASR(ctx, inputFile, robot.asrLanguage, stage.previousAsrText); err != nil {
				return errors.Wrapf(err, "transcription")
			} else {
				asrText = strings.TrimSpace(resp.Text)
				stage.previousAsrText = asrText
				stage.lastRequestASR = time.Now()
				stage.lastAsrDuration = resp.Duration
				stage.lastRequestAsrText = asrText
			}
			logger.Tf(ctx, "ASR ok, robot=%v(%v), lang=%v, prompt=<%v>, resp is <%v>",
				robot.uuid, robot.label, robot.asrLanguage, stage.previousAsrText, asrText)

			// Important trace log.
			logger.Tf(ctx, "You: %v", asrText)

			// Detect empty input and filter badcase.
			if asrText == "" {
				return errors.Errorf("empty asr")
			}
			if robot.asrLanguage == "zh" {
				if strings.Contains(asrText, "请不吝点赞") ||
					strings.Contains(asrText, "支持明镜与点点栏目") ||
					strings.Contains(asrText, "谢谢观看") ||
					strings.Contains(asrText, "請不吝點贊") ||
					strings.Contains(asrText, "支持明鏡與點點欄目") {
					return errors.Errorf("badcase: %v", asrText)
				}
				if strings.Contains(asrText, "字幕由") && strings.Contains(asrText, "社群提供") {
					return errors.Errorf("badcase: %v", asrText)
				}
			} else if robot.asrLanguage == "en" {
				if strings.ToLower(asrText) == "you" ||
					strings.Count(asrText, ".") == len(asrText) {
					return errors.Errorf("badcase: %v", asrText)
				}
			}

			// Keep alive the stage.
			stage.KeepAlive()

			// Insert a dummy sentence to identify the request is alive.
			stage.ttsWorker.SubmitSegment(ctx, stage, NewAnswerSegment(func(segment *AnswerSegment) {
				segment.rid = rid
				segment.dummy = true
			}), nil)

			// Do chat, get the response in stream.
			chatService := &openaiChatService{
				conf: stage.aiConfig,
				onFirstResponse: func(ctx context.Context, text string) {
					stage.lastRequestChat = time.Now()
					stage.lastRobotFirstText = text
				},
			}
			if err := chatService.RequestChat(ctx, rid, stage, robot); err != nil {
				return errors.Wrapf(err, "chat")
			}

			// Notify all popouts about the ASR text.
			for _, popout := range stage.popouts {
				popout.addUserTextMessage(rid, asrText)
			}

			// Response the request UUID and pulling the response.
			ohttp.WriteData(ctx, w, r, struct {
				RequestUUID string `json:"rid"`
				ASR         string `json:"asr"`
			}{
				RequestUUID: rid,
				ASR:         asrText,
			})
			logger.Tf(ctx, "srs ai-talk stage upload ok, rid=%v, asr=%v", rid, asrText)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai-talk/stage/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var sid, rid string
			if err := ParseBody(ctx, r.Body, &struct {
				Token       *string `json:"token"`
				StageUUID   *string `json:"sid"`
				RequestUUID *string `json:"rid"`
			}{
				Token: &token, StageUUID: &sid, RequestUUID: &rid,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if sid == "" {
				return errors.Errorf("empty sid")
			}
			if rid == "" {
				return errors.Errorf("empty rid")
			}

			stage := talkServer.QueryStage(sid)
			if stage == nil {
				return errors.Errorf("invalid sid %v", sid)
			}

			// Keep alive the stage.
			stage.KeepAlive()
			// Switch to the context of stage.
			ctx = stage.loggingCtx
			logger.Tf(ctx, "Stage: Query sid=%v, rid=%v", sid, rid)

			segment := stage.ttsWorker.QueryAnyReadySegment(ctx, stage, rid)
			if segment == nil {
				logger.Tf(ctx, "TTS: No segment for sid=%v, rid=%v", sid, rid)
				ohttp.WriteData(ctx, w, r, struct {
					AnswerSegmentUUID string `json:"asid"`
				}{})
				return nil
			}

			ohttp.WriteData(ctx, w, r, struct {
				// Whether is processing.
				Processing bool `json:"processing"`
				// The UUID for this answer segment.
				AnswerSegmentUUID string `json:"asid"`
				// The TTS text.
				TTS string `json:"tts"`
			}{
				// Whether is processing.
				Processing: segment.dummy || (!segment.ready && segment.err == nil),
				// The UUID for this answer segment.
				AnswerSegmentUUID: segment.asid,
				// The TTS text.
				TTS: segment.text,
			})

			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai-talk/stage/remove"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var sid, rid, asid string
			if err := ParseBody(ctx, r.Body, &struct {
				Token            *string `json:"token"`
				StageUUID        *string `json:"sid"`
				RequestUUID      *string `json:"rid"`
				AudioSegmentUUID *string `json:"asid"`
			}{
				Token: &token, StageUUID: &sid, RequestUUID: &rid, AudioSegmentUUID: &asid,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if sid == "" {
				return errors.Errorf("empty sid")
			}
			if rid == "" {
				return errors.Errorf("empty rid")
			}
			if asid == "" {
				return errors.Errorf("empty asid")
			}

			stage := talkServer.QueryStage(sid)
			if stage == nil {
				return errors.Errorf("invalid sid %v", sid)
			}

			// Keep alive the stage.
			stage.KeepAlive()
			// Switch to the context of stage.
			ctx = stage.loggingCtx
			logger.Tf(ctx, "Stage: Remove sid=%v, rid=%v, asid=%v", sid, rid, asid)

			// Notify to remove the segment.
			segment := stage.ttsWorker.QuerySegment(rid, asid)
			if segment == nil {
				return errors.Errorf("no segment for %v %v", rid, asid)
			}

			// Remove it.
			stage.ttsWorker.RemoveSegment(asid)

			select {
			case <-ctx.Done():
			case segment.removeSignal <- true:
			}

			ohttp.WriteData(ctx, w, r, nil)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai-talk/stage/examples/"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			filename := r.URL.Path[len("/terraform/v1/ai-talk/stage/examples/"):]
			if !strings.Contains(filename, ".") {
				filename = fmt.Sprintf("%v.aac", filename)
			}

			// If there is an optional stage id, we will use the logging context of stage.
			q := r.URL.Query()
			if sid := q.Get("sid"); sid != "" {
				if stage := talkServer.QueryStage(sid); stage != nil {
					ctx = stage.loggingCtx
				}
			}

			ext := strings.Trim(path.Ext(filename), ".")
			contentType := fmt.Sprintf("audio/%v", ext)
			logger.Tf(ctx, "Serve example file=%v, ext=%v, contentType=%v", filename, ext, contentType)

			w.Header().Set("Content-Type", contentType)
			http.ServeFile(w, r, path.Join(workDir, filename))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai-talk/popout/token"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var rid string
			if err := ParseBody(ctx, r.Body, &struct {
				Token    *string `json:"token"`
				RoomUUID *string `json:"room"`
			}{
				Token: &token, RoomUUID: &rid,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			var room SrsLiveRoom
			if r0, err := rdb.HGet(ctx, SRS_LIVE_ROOM, rid).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_LIVE_ROOM, rid)
			} else if r0 == "" {
				return errors.Errorf("live room %v not exists", rid)
			} else if err = json.Unmarshal([]byte(r0), &room); err != nil {
				return errors.Wrapf(err, "unmarshal %v %v", rid, r0)
			}

			if room.PopoutToken == "" {
				room.PopoutToken = uuid.NewString()

				if b, err := json.Marshal(room); err != nil {
					return errors.Wrapf(err, "marshal room")
				} else if err := rdb.HSet(ctx, SRS_LIVE_ROOM, room.UUID, string(b)).Err(); err != nil {
					return errors.Wrapf(err, "hset %v %v %v", SRS_LIVE_ROOM, room.UUID, string(b))
				}
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Token string `json:"token"`
			}{
				Token: room.PopoutToken,
			})
			logger.Tf(ctx, "srs ai-talk create popout token ok")
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai-talk/popout/verify"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var rid string
			if err := ParseBody(ctx, r.Body, &struct {
				Token    *string `json:"token"`
				RoomUUID *string `json:"room"`
			}{
				Token: &token, RoomUUID: &rid,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			var room SrsLiveRoom
			if r0, err := rdb.HGet(ctx, SRS_LIVE_ROOM, rid).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_LIVE_ROOM, rid)
			} else if r0 == "" {
				return errors.Errorf("live room %v not exists", rid)
			} else if err = json.Unmarshal([]byte(r0), &room); err != nil {
				return errors.Wrapf(err, "unmarshal %v %v", rid, r0)
			}

			if room.PopoutToken != token {
				return errors.Errorf("invalid token")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			ohttp.WriteData(ctx, w, r, &struct {
				Token string `json:"token"`
			}{
				Token: apiSecret,
			})
			logger.Tf(ctx, "srs ai-talk verify popout token ok")
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai-talk/subscribe/start"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var rid string
			if err := ParseBody(ctx, r.Body, &struct {
				Token    *string `json:"token"`
				RoomUUID *string `json:"room"`
			}{
				Token: &token, RoomUUID: &rid,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if rid == "" {
				return errors.Errorf("empty rid")
			}

			var room SrsLiveRoom
			if r0, err := rdb.HGet(ctx, SRS_LIVE_ROOM, rid).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_LIVE_ROOM, rid)
			} else if r0 == "" {
				return errors.Errorf("live room %v not exists", rid)
			} else if err = json.Unmarshal([]byte(r0), &room); err != nil {
				return errors.Wrapf(err, "unmarshal %v %v", rid, r0)
			}

			// Use the latest stage as popout's source stage, note that it may change.
			stage := talkServer.QueryStage(room.StageUUID)
			if stage == nil {
				return errors.Errorf("no stage in room %v", rid)
			}

			// TODO: FIXME: Cleanup popouts for a room.
			ctx = logger.WithContext(ctx)
			popout := NewStagePopout(func(popout *StagePopout) {
				popout.loggingCtx = ctx
				popout.room = &room

				// Bind the popout to the stage.
				popout.stage = stage
				stage.addPopout(popout)
			})

			// Keep alive the stage.
			stage.KeepAlive()

			type StageRobotResult struct {
				UUID  string `json:"uuid"`
				Label string `json:"label"`
				Voice string `json:"voice"`
			}
			type StageResult struct {
				StageID       string           `json:"sid"`
				StagePopoutID string           `json:"spid"`
				Robot         StageRobotResult `json:"robot"`
			}
			r0 := &StageResult{
				StageID:       stage.sid,
				StagePopoutID: popout.spid,
				Robot: StageRobotResult{
					UUID:  stage.robot.uuid,
					Label: stage.robot.label,
					Voice: stage.robot.voice,
				},
			}

			ohttp.WriteData(ctx, w, r, &r0)
			logger.Tf(ctx, "Stage: create popout ok, room=%v, sid=%v, spid=%v",
				room.UUID, stage.sid, popout.spid)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai-talk/subscribe/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var sid, spid string
			if err := ParseBody(ctx, r.Body, &struct {
				Token         *string `json:"token"`
				StageUUID     *string `json:"sid"`
				StagePopoutID *string `json:"spid"`
			}{
				Token: &token, StageUUID: &sid, StagePopoutID: &spid,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if sid == "" {
				return errors.Errorf("empty sid")
			}
			if spid == "" {
				return errors.Errorf("empty spid")
			}

			stage := talkServer.QueryStage(sid)
			if stage == nil {
				return errors.Errorf("invalid sid %v", sid)
			}

			// Keep alive the stage.
			stage.KeepAlive()
			// Switch to the context of stage.
			ctx = stage.loggingCtx
			logger.Tf(ctx, "Stage: Query sid=%v, rid=%v", sid, stage.room.UUID)

			popout := stage.queryPopout(spid)
			if popout == nil {
				return errors.Errorf("invalid spid %v of sid %v", spid, sid)
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Messages []*StageMessage `json:"msgs"`
			}{
				Messages: popout.flushMessages(),
			})
			logger.Tf(ctx, "srs ai-talk query popout stage ok")
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai-talk/subscribe/tts"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			q := r.URL.Query()

			sid := q.Get("sid")
			if sid == "" {
				return errors.Errorf("empty sid")
			}

			spid := q.Get("spid")
			if spid == "" {
				return errors.Errorf("empty spid")
			}

			asid := q.Get("asid")
			if asid == "" {
				return errors.Errorf("empty asid")
			}

			stage := talkServer.QueryStage(sid)
			if stage == nil {
				return errors.Errorf("invalid sid %v", sid)
			}

			// Keep alive the stage.
			stage.KeepAlive()
			// Switch to the context of stage.
			ctx = stage.loggingCtx
			logger.Tf(ctx, "Stage: Download sid=%v, spid=%v, asid=%v", sid, spid, asid)

			popout := stage.queryPopout(spid)
			if popout == nil {
				return errors.Errorf("invalid spid %v of sid %v", spid, sid)
			}

			audioFile := popout.queryAudioFile(asid)

			// Read the ttsFile and response it as opus audio.
			if strings.HasSuffix(audioFile, ".wav") {
				w.Header().Set("Content-Type", "audio/wav")
			} else {
				w.Header().Set("Content-Type", "audio/aac")
			}
			http.ServeFile(w, r, audioFile)

			logger.Tf(ctx, "srs ai-talk play tts popout stage ok")
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai-talk/subscribe/remove"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var sid, spid, asid string
			if err := ParseBody(ctx, r.Body, &struct {
				Token            *string `json:"token"`
				StageUUID        *string `json:"sid"`
				StagePopoutID    *string `json:"spid"`
				AudioSegmentUUID *string `json:"asid"`
			}{
				Token: &token, StageUUID: &sid, StagePopoutID: &spid, AudioSegmentUUID: &asid,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if sid == "" {
				return errors.Errorf("empty sid")
			}
			if spid == "" {
				return errors.Errorf("empty spid")
			}
			if asid == "" {
				return errors.Errorf("empty asid")
			}

			stage := talkServer.QueryStage(sid)
			if stage == nil {
				return errors.Errorf("invalid sid %v", sid)
			}

			// Keep alive the stage.
			stage.KeepAlive()
			// Switch to the context of stage.
			ctx = stage.loggingCtx
			logger.Tf(ctx, "Stage: Query sid=%v, rid=%v", sid, stage.room.UUID)

			popout := stage.queryPopout(spid)
			if popout == nil {
				return errors.Errorf("invalid spid %v of sid %v", spid, sid)
			}

			if err := popout.removeMessage(asid); err != nil {
				return errors.Wrapf(err, "remove message asid=%v of sid=%v, spid=%v", asid, sid, spid)
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "srs ai-talk remove popout stage file ok")
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	return nil
}
