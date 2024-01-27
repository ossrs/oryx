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
	"regexp"
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
	outputFile := fmt.Sprintf("%v.mp4", inputFile)
	defer os.Remove(outputFile)

	// Transcode input audio in opus or aac, to aac in m4a/mp4 format.
	// If need to encode to aac, use:
	//		"-c:a", "aac", "-ac", "1", "-ar", "16000", "-ab", "30k",
	if err := exec.CommandContext(ctx, "ffmpeg",
		"-i", inputFile,
		"-vn", "-c:a", "copy",
		outputFile,
	).Run(); err != nil {
		return nil, errors.Errorf("Error converting the file")
	}
	logger.Tf(ctx, "Convert audio %v to %v ok", inputFile, outputFile)

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
			// Note that must use verbose JSON, to get the duration of file.
			Format:   openai.AudioResponseFormatVerboseJSON,
			Language: language,
			Prompt:   prompt,
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "asr")
	}

	return &ASRResult{Text: resp.Text, Duration: time.Duration(resp.Duration * float64(time.Second))}, nil
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

func (v *openaiChatService) RequestChat(ctx context.Context, sreq *StageRequest, stage *Stage, robot *Robot) error {
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
		if err := v.handle(ctx, stage, robot, sreq, gptChatStream); err != nil {
			logger.Ef(ctx, "Handle stream failed, err %+v", err)
		}
	}()

	return nil
}

func (v *openaiChatService) handle(ctx context.Context, stage *Stage, robot *Robot, sreq *StageRequest, gptChatStream *openai.ChatCompletionStream) error {
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

		return finished, filteredStencese, nil
	}

	gotNewSentence := func(sentence, lastWords string, firstSentense bool) bool {
		newSentence := false

		isEnglish := func(s string) bool {
			for _, r := range s {
				if r > unicode.MaxASCII {
					return false
				}
			}
			return true
		}

		// Ignore empty.
		if sentence == "" {
			return newSentence
		}

		// Any ASCII character to split sentence.
		if strings.ContainsAny(lastWords, ",.?!\n") {
			newSentence = true
		}

		// Any Chinese character to split sentence.
		if strings.ContainsRune(lastWords, '。') ||
			strings.ContainsRune(lastWords, '？') ||
			strings.ContainsRune(lastWords, '！') ||
			strings.ContainsRune(lastWords, '，') {
			newSentence = true
		}

		// Badcase, for number such as 1.3, or 1,300,000.
		var badcase bool
		if match, _ := regexp.MatchString(`\d+(\.|,)\d*$`, sentence); match {
			badcase, newSentence = true, false
		}

		// Determine whether new sentence by length.
		if isEnglish(sentence) {
			maxWords, minWords := 30, 3
			if !firstSentense || badcase {
				maxWords, minWords = 50, 5
			}

			if nn := strings.Count(sentence, " "); nn >= maxWords {
				newSentence = true
			} else if nn < minWords {
				newSentence = false
			}
		} else {
			maxWords, minWords := 50, 3
			if !firstSentense || badcase {
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

	commitAISentence := func(sentence string, firstSentense bool) {
		filteredSentence := sentence
		if strings.TrimSpace(sentence) == "" {
			return
		}

		if firstSentense {
			if robot.prefix != "" {
				filteredSentence = fmt.Sprintf("%v %v", robot.prefix, filteredSentence)
			}
			if v.onFirstResponse != nil {
				v.onFirstResponse(ctx, filteredSentence)
			}
		}

		segment := NewAnswerSegment(func(segment *AnswerSegment) {
			segment.request = sreq
			segment.text = filteredSentence
			segment.first = firstSentense
		})
		stage.ttsWorker.SubmitSegment(ctx, stage, sreq, segment)

		logger.Tf(ctx, "TTS: Commit segment rid=%v, asid=%v, first=%v, sentence is %v",
			sreq.rid, segment.asid, firstSentense, filteredSentence)
		return
	}

	var sentence, lastWords string
	isFinished, firstSentense := false, true
	for !isFinished && ctx.Err() == nil {
		response, err := gptChatStream.Recv()
		if finished, words, err := filterAIResponse(&response, err); err != nil {
			return errors.Wrapf(err, "filter")
		} else {
			isFinished, sentence, lastWords = finished, sentence+words, words
		}
		//logger.Tf(ctx, "AI response: text=%v plus %v", lastWords, sentence)

		newSentence := gotNewSentence(sentence, lastWords, firstSentense)
		if !isFinished && !newSentence {
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

// The StageRequest is a request from user, submited to ASR and AI, generated answer segments,
// finally delivered as stage messages to subscribers.
type StageRequest struct {
	// The request UUID.
	rid string
	// The upload audio file, the input file.
	inputFile string
	// The ASR text, converted from audio.
	asrText string
	// Whether the request is finished.
	finished bool

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

	// The answer segments and pieces, answer in text by AI, and in audio by TTS.
	segments []*AnswerSegment
	// The owner stage.
	stage *Stage
}

func (v *StageRequest) onSegmentReady(segment *AnswerSegment) {
	if segment.first {
		v.lastRequestTTS = time.Now()
	}

	for _, s := range v.segments {
		if !s.ready && s.err == nil {
			return
		}
	}
	v.finished = true
}

func (v *StageRequest) Close() error {
	return v.FastDispose()
}

// Fast cleanup the files of request, after converted to text by ASR service.
func (v *StageRequest) FastDispose() error {
	if v.inputFile != "" {
		if _, err := os.Stat(v.inputFile); err == nil {
			os.Remove(v.inputFile)
		}
	}

	return nil
}

func (v *StageRequest) asrAudioToText(ctx context.Context, aiConfig openai.ClientConfig, asrLanguage, previousAsrText string) error {
	var asrText string
	var asrDuration time.Duration

	asrService := NewOpenAIASRService(aiConfig, func(*openaiASRService) {
		v.lastExtractAudio = time.Now()
	})

	if resp, err := asrService.RequestASR(ctx, v.inputFile, asrLanguage, previousAsrText); err != nil {
		return errors.Wrapf(err, "transcription")
	} else {
		asrText, asrDuration = strings.TrimSpace(resp.Text), resp.Duration
	}

	// Detect empty input and filter badcase.
	if asrText == "" {
		return errors.Errorf("empty asr")
	}
	if asrLanguage == "zh" {
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
	} else if asrLanguage == "en" {
		if strings.ToLower(asrText) == "you" ||
			strings.Count(asrText, ".") == len(asrText) {
			return errors.Errorf("badcase: %v", asrText)
		}
	}

	v.asrText = asrText
	v.lastRequestASR = time.Now()
	v.lastAsrDuration = asrDuration
	v.lastRequestAsrText = asrText

	return nil
}

func (v *StageRequest) receiveInputFile(ctx context.Context, audioBase64Data string) error {
	data, err := base64.StdEncoding.DecodeString(audioBase64Data)
	if err != nil {
		return errors.Errorf("decode base64 from %v", audioBase64Data)
	}

	out, err := os.Create(v.inputFile)
	if err != nil {
		return errors.Errorf("Unable to create the file for writing")
	}
	defer out.Close()

	nn, err := io.Copy(out, bytes.NewReader([]byte(data)))
	if err != nil {
		return errors.Errorf("Error writing the file")
	}
	logger.Tf(ctx, "File saved to %v, size: %v", v.inputFile, nn)

	v.lastUploadAudio = time.Now()
	return nil
}

func (v *StageRequest) total() float64 {
	if v.lastDownloadAudio.After(v.lastSentence) {
		return float64(v.lastDownloadAudio.Sub(v.lastSentence)) / float64(time.Second)
	}
	return 0
}

func (v *StageRequest) upload() float64 {
	if v.lastUploadAudio.After(v.lastSentence) {
		return float64(v.lastUploadAudio.Sub(v.lastSentence)) / float64(time.Second)
	}
	return 0
}

func (v *StageRequest) exta() float64 {
	if v.lastExtractAudio.After(v.lastUploadAudio) {
		return float64(v.lastExtractAudio.Sub(v.lastUploadAudio)) / float64(time.Second)
	}
	return 0
}

func (v *StageRequest) asr() float64 {
	if v.lastRequestASR.After(v.lastExtractAudio) {
		return float64(v.lastRequestASR.Sub(v.lastExtractAudio)) / float64(time.Second)
	}
	return 0
}

func (v *StageRequest) chat() float64 {
	if v.lastRequestChat.After(v.lastRequestASR) {
		return float64(v.lastRequestChat.Sub(v.lastRequestASR)) / float64(time.Second)
	}
	return 0
}

func (v *StageRequest) tts() float64 {
	if v.lastRequestTTS.After(v.lastRequestChat) {
		return float64(v.lastRequestTTS.Sub(v.lastRequestChat)) / float64(time.Second)
	}
	return 0
}

func (v *StageRequest) download() float64 {
	if v.lastDownloadAudio.After(v.lastRequestTTS) {
		return float64(v.lastDownloadAudio.Sub(v.lastRequestTTS)) / float64(time.Second)
	}
	return 0
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

	// Whether ready to flush to client.
	finished bool
	// The error object of message.
	err error

	// The owner subscriber.
	subscriber *StageSubscriber
	// The source segment.
	segment *AnswerSegment
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

// The StageSubscriber is a subscriber from a stage.
type StageSubscriber struct {
	// StageSubscriber UUID.
	spid string
	// Last update of stage.
	update time.Time
	// The logging context, to write all logs in one context for a sage.
	loggingCtx context.Context
	// The messages from user ASR and AI responses.
	messages []*StageMessage

	// The owner room, never changes.
	room *SrsLiveRoom
	// The owner stage, never changes.
	stage *Stage
}

func NewStageSubscriber(opts ...func(*StageSubscriber)) *StageSubscriber {
	v := &StageSubscriber{
		// Create new UUID.
		spid: uuid.NewString(),
	}

	for _, opt := range opts {
		opt(v)
	}
	return v
}

func (v *StageSubscriber) Close() error {
	for _, message := range v.messages {
		_ = message.Close()
	}
	return nil
}

func (v *StageSubscriber) Expired() bool {
	return time.Since(v.update) > 300*time.Second
}

func (v *StageSubscriber) KeepAlive() {
	v.update = time.Now()
}

func (v *StageSubscriber) addUserTextMessage(rid, msg string) {
	v.messages = append(v.messages, &StageMessage{
		finished: true, MessageUUID: uuid.NewString(), subscriber: v,
		RequestUUID: rid, Role: "user", Message: msg,
	})
}

// Create a robot empty message, to keep the order of messages.
func (v *StageSubscriber) createRobotEmptyMessage() *StageMessage {
	message := &StageMessage{
		finished: false, MessageUUID: uuid.NewString(), subscriber: v,
	}
	v.messages = append(v.messages, message)
	return message
}

func (v *StageSubscriber) completeRobotAudioMessage(ctx context.Context, sreq *StageRequest, segment *AnswerSegment, message *StageMessage) {
	// Build a new copy file of ttsFile.
	ttsExt := path.Ext(segment.ttsFile)
	copyFile := fmt.Sprintf("%v-copy-%v%v", segment.ttsFile[:len(segment.ttsFile)-len(ttsExt)], v.spid, ttsExt)

	// Copy the ttsFile to copyFile.
	if err := func() error {
		// If segment is error, ignore.
		if segment.err != nil {
			return nil
		}

		src, err := os.Open(segment.ttsFile)
		if err != nil {
			return errors.Errorf("open %v for reading", segment.ttsFile)
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
			segment.ttsFile, copyFile, v.room.UUID, v.stage.sid, v.spid)
		return nil
	}(); err != nil {
		message.err = errors.Wrapf(err, "copy %v to %v", segment.ttsFile, copyFile)
	}

	message.finished, message.segment = true, segment
	message.RequestUUID, message.SegmentUUID = sreq.rid, segment.asid
	message.Role, message.Message, message.audioFile = "robot", segment.text, copyFile

	// Always close message if timeout.
	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(30 * time.Second):
			message.Close()
		}
	}()
}

// TODO: Cleanup flushed messages.
func (v *StageSubscriber) flushMessages() (messages []*StageMessage, pending bool) {
	for _, message := range v.messages {
		// TODO: Dispose and cleanup flushed messages.
		// Ignore if already flushed.
		if message.flushed {
			continue
		}

		// If message not ready, break to keep the order.
		if !message.finished {
			pending = true
			break
		}

		// Flush the message once finished.
		message.flushed = true

		// Ignore if error.
		err := message.err
		if err == nil && message.segment != nil {
			err = message.segment.err
		}
		if err != nil {
			// TODO: Handle the error.
			continue
		}

		// Got a good piece of message for subscriber.
		messages = append(messages, message)
	}
	return
}

func (v *StageSubscriber) queryAudioFile(asid string) *StageMessage {
	for _, message := range v.messages {
		if message.SegmentUUID == asid {
			return message
		}
	}
	return nil
}

func (v *StageSubscriber) removeMessage(asid string) error {
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

	// The robot created for this stage.
	robot *Robot
	// The AI configuration.
	aiConfig openai.ClientConfig
	// The room it belongs to.
	room *SrsLiveRoom
	// All the requests from user.
	requests []*StageRequest
	// All the subscribers binding to this stage.
	subscribers []*StageSubscriber
	// Cache the room level token for popout.
	roomToken string
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
	for _, subscriber := range v.subscribers {
		subscriber.Close()
	}
	for _, request := range v.requests {
		request.Close()
	}
	return v.ttsWorker.Close()
}

func (v *Stage) Expired() bool {
	return time.Since(v.update) > 600*time.Second
}

func (v *Stage) KeepAlive() {
	v.update = time.Now()
}

func (v *Stage) addRequest(request *StageRequest) {
	v.requests = append(v.requests, request)
}

func (v *Stage) queryRequest(rid string) *StageRequest {
	for _, request := range v.requests {
		if request.rid == rid {
			return request
		}
	}
	return nil
}

func (v *Stage) addSubscriber(subscriber *StageSubscriber) {
	v.subscribers = append(v.subscribers, subscriber)
}

func (v *Stage) querySubscriber(spid string) *StageSubscriber {
	for _, subscriber := range v.subscribers {
		if subscriber.spid == spid {
			return subscriber
		}
	}
	return nil
}

func (v *Stage) removeSubscriber(subscriber *StageSubscriber) {
	for i, s := range v.subscribers {
		if s.spid == subscriber.spid {
			v.subscribers = append(v.subscribers[:i], v.subscribers[i+1:]...)
			return
		}
	}
}

// The AnswerSegment is a segment of answer, which is a sentence.
type AnswerSegment struct {
	// Request UUID.
	request *StageRequest
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
	// Signal to remove the TTS file immediately.
	removeSignal chan bool
	// Whether we have logged this segment.
	logged bool
	// Whether the segment is the first response.
	first bool
}

func NewAnswerSegment(opts ...func(segment *AnswerSegment)) *AnswerSegment {
	v := &AnswerSegment{
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
	// TODO: FIXME: Remove this because they are stored in the stage request.
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

func (v *TTSWorker) SubmitSegment(ctx context.Context, stage *Stage, sreq *StageRequest, segment *AnswerSegment) {
	var messages []*StageMessage

	func() {
		v.lock.Lock()
		defer v.lock.Unlock()

		// Append the sentence to queue.
		v.segments = append(v.segments, segment)

		// TODO: Should not use the lock of tts worker.
		// Add segment to stage request.
		sreq.segments = append(sreq.segments, segment)

		// TODO: Should not use the lock of tts worker.
		// Add message to subscriber, to keep the same order as segments.
		for _, subscriber := range stage.subscribers {
			message := subscriber.createRobotEmptyMessage()
			messages = append(messages, message)
		}
	}()

	// Start a goroutine to do TTS task.
	v.wg.Add(1)
	go func() {
		defer v.wg.Done()

		ttsService := NewOpenAITTSService(stage.aiConfig)
		if err := ttsService.RequestTTS(ctx, func(ext string) string {
			segment.ttsFile = path.Join(workDir,
				fmt.Sprintf("assistant-%v-sentence-%v-tts.%v", sreq.rid, segment.asid, ext),
			)
			return segment.ttsFile
		}, segment.text); err != nil {
			segment.err = err
		} else {
			segment.ready = true
			sreq.onSegmentReady(segment)
			logger.Tf(ctx, "File saved to %v, %v", segment.ttsFile, segment.text)
		}

		// Update all messages.
		for _, m := range messages {
			m.subscriber.completeRobotAudioMessage(ctx, sreq, segment, m)
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
			var roomUUID string
			if err := ParseBody(ctx, r.Body, &struct {
				Token    *string `json:"token"`
				RoomUUID *string `json:"room"`
			}{
				Token: &token, RoomUUID: &roomUUID,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			var room SrsLiveRoom
			if r0, err := rdb.HGet(ctx, SRS_LIVE_ROOM, roomUUID).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_LIVE_ROOM, roomUUID)
			} else if r0 == "" {
				return errors.Errorf("live room %v not exists", roomUUID)
			} else if err = json.Unmarshal([]byte(r0), &room); err != nil {
				return errors.Wrapf(err, "unmarshal %v %v", roomUUID, r0)
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
					// Cache the room token to stage.
					stage.roomToken = room.PopoutToken

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
				StageID   string           `json:"sid"`
				RoomToken string           `json:"roomToken"`
				Robot     StageRobotResult `json:"robot"`
			}
			r0 := &StageResult{
				StageID:   stage.sid,
				RoomToken: stage.roomToken,
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

	ep = "/terraform/v1/ai-talk/stage/conversation"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var sid, robotUUID string
			if err := ParseBody(ctx, r.Body, &struct {
				Token     *string `json:"token"`
				StageUUID *string `json:"sid"`
				RobotUUID *string `json:"robot"`
			}{
				Token: &token, StageUUID: &sid, RobotUUID: &robotUUID,
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
			// Switch to the context of stage.
			ctx = stage.loggingCtx

			robot := stage.robot
			if robot == nil {
				return errors.Errorf("invalid robot %v", robotUUID)
			}

			// The rid is the request id, which identify this request, generally a question.
			sreq := &StageRequest{rid: uuid.NewString(), stage: stage}
			sreq.lastSentence = time.Now()
			stage.addRequest(sreq)

			// Keep alive the stage.
			stage.KeepAlive()

			// Response the request UUID and pulling the response.
			ohttp.WriteData(ctx, w, r, struct {
				RequestUUID string `json:"rid"`
			}{
				RequestUUID: sreq.rid,
			})
			logger.Tf(ctx, "srs ai-talk stage create conversation ok, rid=%v", sreq.rid)
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
			var sid, rid, robotUUID, audioBase64Data string
			var userMayInput float64
			if err := ParseBody(ctx, r.Body, &struct {
				Token        *string  `json:"token"`
				StageUUID    *string  `json:"sid"`
				RobotUUID    *string  `json:"robot"`
				RequestUUID  *string  `json:"rid"`
				UserMayInput *float64 `json:"umi"`
				AudioData    *string  `json:"audio"`
			}{
				Token: &token, StageUUID: &sid, RobotUUID: &robotUUID, RequestUUID: &rid,
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
			if rid == "" {
				return errors.Errorf("empty rid")
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
			// Switch to the context of stage.
			ctx = stage.loggingCtx

			robot := stage.robot
			if robot == nil {
				return errors.Errorf("invalid robot %v", robotUUID)
			}

			// Query the request.
			sreq := stage.queryRequest(rid)
			if sreq == nil {
				return errors.Errorf("invalid sid=%v, rid=%v", sid, rid)
			}

			// The rid is the request id, which identify this request, generally a question.
			defer sreq.FastDispose()

			sreq.inputFile = path.Join(workDir, fmt.Sprintf("assistant-%v-input.audio", sreq.rid))
			logger.Tf(ctx, "Stage: Got question sid=%v, umi=%v, robot=%v(%v), rid=%v, input=%v",
				sid, userMayInput, robot.uuid, robot.label, sreq.rid, sreq.inputFile)

			// Save audio input to file.
			if err := sreq.receiveInputFile(ctx, audioBase64Data); err != nil {
				return errors.Wrapf(err, "save %vB audio to file %v", len(audioBase64Data), sreq.inputFile)
			}

			// Do ASR, convert to text.
			if err := sreq.asrAudioToText(ctx, stage.aiConfig, robot.asrLanguage, stage.previousAsrText); err != nil {
				return errors.Wrapf(err, "asr lang=%v, previous=%v", robot.asrLanguage, stage.previousAsrText)
			}
			logger.Tf(ctx, "ASR ok, robot=%v(%v), lang=%v, prompt=<%v>, resp is <%v>",
				robot.uuid, robot.label, robot.asrLanguage, stage.previousAsrText, sreq.asrText)

			// Important trace log.
			stage.previousAsrText = sreq.asrText
			logger.Tf(ctx, "You: %v", sreq.asrText)

			// Keep alive the stage.
			stage.KeepAlive()

			// Do chat, get the response in stream.
			chatService := &openaiChatService{
				conf: stage.aiConfig,
				onFirstResponse: func(ctx context.Context, text string) {
					sreq.lastRequestChat = time.Now()
					sreq.lastRobotFirstText = text
				},
			}
			if err := chatService.RequestChat(ctx, sreq, stage, robot); err != nil {
				return errors.Wrapf(err, "chat")
			}

			// Notify all subscribers about the ASR text.
			for _, subscriber := range stage.subscribers {
				subscriber.addUserTextMessage(sreq.rid, sreq.asrText)
			}

			// Response the request UUID and pulling the response.
			ohttp.WriteData(ctx, w, r, struct {
				RequestUUID string `json:"rid"`
				ASR         string `json:"asr"`
			}{
				RequestUUID: sreq.rid,
				ASR:         sreq.asrText,
			})
			logger.Tf(ctx, "srs ai-talk stage upload ok, rid=%v, asr=%v", sreq.rid, sreq.asrText)
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

			// Query the request.
			sreq := stage.queryRequest(rid)
			if sreq == nil {
				return errors.Errorf("invalid sid=%v, rid=%v", sid, rid)
			}

			ohttp.WriteData(ctx, w, r, struct {
				Finished bool `json:"finished"`
			}{
				Finished: !sreq.finished,
			})

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

	ep = "/terraform/v1/ai-talk/stage/verify"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var roomToken string
			var roomUUID string
			if err := ParseBody(ctx, r.Body, &struct {
				RoomUUID  *string `json:"room"`
				RoomToken *string `json:"roomToken"`
			}{
				RoomToken: &roomToken, RoomUUID: &roomUUID,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			var room SrsLiveRoom
			if r0, err := rdb.HGet(ctx, SRS_LIVE_ROOM, roomUUID).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_LIVE_ROOM, roomUUID)
			} else if r0 == "" {
				return errors.Errorf("live room %v not exists", roomUUID)
			} else if err = json.Unmarshal([]byte(r0), &room); err != nil {
				return errors.Wrapf(err, "unmarshal %v %v", roomUUID, r0)
			}

			if room.PopoutToken != roomToken {
				return errors.Errorf("invalid room token")
			}

			// TODO: To improve security level, we should not response the bearer token, instead we can
			//   support authentication with room token.
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
			var roomUUID string
			if err := ParseBody(ctx, r.Body, &struct {
				Token    *string `json:"token"`
				RoomUUID *string `json:"room"`
			}{
				Token: &token, RoomUUID: &roomUUID,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if roomUUID == "" {
				return errors.Errorf("empty room id")
			}

			var room SrsLiveRoom
			if r0, err := rdb.HGet(ctx, SRS_LIVE_ROOM, roomUUID).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_LIVE_ROOM, roomUUID)
			} else if r0 == "" {
				return errors.Errorf("live room %v not exists", roomUUID)
			} else if err = json.Unmarshal([]byte(r0), &room); err != nil {
				return errors.Wrapf(err, "unmarshal %v %v", roomUUID, r0)
			}

			// Use the latest stage as subscriber's source stage, note that it may change.
			stage := talkServer.QueryStage(room.StageUUID)
			if stage == nil {
				return errors.Errorf("no stage in room %v", roomUUID)
			}

			// TODO: FIXME: Cleanup subscribers for a room.
			ctx = logger.WithContext(ctx)
			subscriber := NewStageSubscriber(func(subscriber *StageSubscriber) {
				subscriber.loggingCtx = ctx
				subscriber.room = &room

				// Bind the subscriber to the stage.
				subscriber.stage = stage
				stage.addSubscriber(subscriber)
			})

			go func() {
				defer subscriber.Close()

				for ctx.Err() == nil {
					select {
					case <-ctx.Done():
					case <-time.After(3 * time.Second):
						if subscriber.Expired() {
							logger.Tf(ctx, "Stage: Remove spid=%v from stage sid=%v for expired, update=%v",
								subscriber.spid, stage.sid, subscriber.update.Format(time.RFC3339))
							stage.removeSubscriber(subscriber)
							return
						}
					}
				}
			}()

			// Keep alive the stage.
			stage.KeepAlive()
			subscriber.KeepAlive()

			type StageRobotResult struct {
				UUID  string `json:"uuid"`
				Label string `json:"label"`
				Voice string `json:"voice"`
			}
			type StageResult struct {
				StageID      string           `json:"sid"`
				SubscriberID string           `json:"spid"`
				Robot        StageRobotResult `json:"robot"`
			}
			r0 := &StageResult{
				StageID:      stage.sid,
				SubscriberID: subscriber.spid,
				Robot: StageRobotResult{
					UUID:  stage.robot.uuid,
					Label: stage.robot.label,
					Voice: stage.robot.voice,
				},
			}

			ohttp.WriteData(ctx, w, r, &r0)
			logger.Tf(ctx, "Stage: create subscriber ok, room=%v, sid=%v, spid=%v",
				room.UUID, stage.sid, subscriber.spid)
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
				Token        *string `json:"token"`
				StageUUID    *string `json:"sid"`
				SubscriberID *string `json:"spid"`
			}{
				Token: &token, StageUUID: &sid, SubscriberID: &spid,
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

			subscriber := stage.querySubscriber(spid)
			if subscriber == nil {
				return errors.Errorf("invalid spid %v of sid %v", spid, sid)
			}

			// Keep alive the stage.
			stage.KeepAlive()
			subscriber.KeepAlive()
			// Switch to the context of stage.
			ctx = stage.loggingCtx
			// Note that we should disable detail meanless logs for subscribe.
			//logger.Tf(ctx, "Stage: Query sid=%v, room=%v", sid, stage.room.UUID)

			msgs, pending := subscriber.flushMessages()
			ohttp.WriteData(ctx, w, r, &struct {
				// Finished messages.
				Messages []*StageMessage `json:"msgs"`
				// Is there any pending messages.
				Pending bool `json:"pending"`
			}{
				Messages: msgs,
				Pending:  pending,
			})

			// Note that we should disable detail meanless logs for subscribe.
			//logger.Tf(ctx, "srs ai-talk query subscriber stage ok")
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

			roomToken := q.Get("roomToken")
			if asid == "" {
				return errors.Errorf("empty roomToken")
			}

			stage := talkServer.QueryStage(sid)
			if stage == nil {
				return errors.Errorf("invalid sid %v", sid)
			}

			if stage.roomToken != roomToken {
				return errors.Errorf("invalid roomToken")
			}

			subscriber := stage.querySubscriber(spid)
			if subscriber == nil {
				return errors.Errorf("invalid spid %v of sid %v", spid, sid)
			}

			answer := subscriber.queryAudioFile(asid)
			if answer == nil {
				return errors.Errorf("invalid asid %v of sid %v", asid, sid)
			}

			// Keep alive the stage.
			stage.KeepAlive()
			subscriber.KeepAlive()
			// Switch to the context of stage.
			ctx = stage.loggingCtx
			logger.Tf(ctx, "Stage: Download sid=%v, spid=%v, asid=%v", sid, spid, asid)

			// When the first subscriber got the segment, we log the elapsed time.
			if segment := answer.segment; !segment.logged {
				sreq := segment.request
				if segment.first {
					sreq.lastDownloadAudio = time.Now()
					speech := float64(sreq.lastAsrDuration) / float64(time.Second)
					logger.Tf(ctx, "Elapsed cost total=%.1fs, steps=[upload=%.1fs,exta=%.1fs,asr=%.1fs,chat=%.1fs,tts=%.1fs,download=%.1fs], ask=%v, speech=%.1fs, answer=%v",
						sreq.total(), sreq.upload(), sreq.exta(), sreq.asr(), sreq.chat(), sreq.tts(), sreq.download(),
						sreq.lastRequestAsrText, speech, sreq.lastRobotFirstText)
				}

				// Important trace log. Note that browser may request multiple times, so we only log for the first
				// request to reduce logs.
				segment.logged = true
				logger.Tf(ctx, "Bot: %v", segment.text)
			}

			// Read the ttsFile and response it as opus audio.
			if strings.HasSuffix(answer.audioFile, ".wav") {
				w.Header().Set("Content-Type", "audio/wav")
			} else {
				w.Header().Set("Content-Type", "audio/aac")
			}
			http.ServeFile(w, r, answer.audioFile)

			logger.Tf(ctx, "srs ai-talk play tts subscriber stage ok")
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
				SubscriberID     *string `json:"spid"`
				AudioSegmentUUID *string `json:"asid"`
			}{
				Token: &token, StageUUID: &sid, SubscriberID: &spid, AudioSegmentUUID: &asid,
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

			subscriber := stage.querySubscriber(spid)
			if subscriber == nil {
				return errors.Errorf("invalid spid %v of sid %v", spid, sid)
			}

			// Keep alive the stage.
			stage.KeepAlive()
			subscriber.KeepAlive()
			// Switch to the context of stage.
			ctx = stage.loggingCtx
			logger.Tf(ctx, "Stage: Remove segment room=%v, sid=%v, spid=%v, asid=%v",
				stage.room.UUID, sid, spid, asid)

			if err := subscriber.removeMessage(asid); err != nil {
				return errors.Wrapf(err, "remove message asid=%v of sid=%v, spid=%v", asid, sid, spid)
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "srs ai-talk remove subscriber stage file ok")
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	return nil
}
