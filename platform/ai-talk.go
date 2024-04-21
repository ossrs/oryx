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
var aiTalkWorkDir, aiTalkExampleDir string

type ASRResult struct {
	Text     string
	Duration time.Duration
}

type openaiASRService struct {
	conf openai.ClientConfig
	// The callback before start ASR request.
	onBeforeRequest func()
}

func NewOpenAIASRService(conf openai.ClientConfig, opts ...func(service *openaiASRService)) *openaiASRService {
	v := &openaiASRService{conf: conf}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

func (v *openaiASRService) RequestASR(ctx context.Context, inputFile, language, prompt string) (*ASRResult, error) {
	outputFile := fmt.Sprintf("%v.mp4", inputFile)
	defer os.Remove(outputFile)

	// TODO: FIXME: Use client to set the codec and format, skip to copy it in server, because it need about
	//   1s to process in some weak performance VPS.
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

	// Silent detect, see https://platform.openai.com/docs/api-reference/audio/verbose-json-object
	//		no_speech_prob Probability of no speech in the segment. If the value is higher than 1.0 and
	//			the avg_logprob is below -1, consider this segment silent.
	for _, segment := range resp.Segments {
		if segment.NoSpeechProb > 0.8 {
			return nil, errors.Errorf("silent detect, prob=%v, file=%v, text=%v",
				segment.NoSpeechProb, outputFile, segment.Text)
		} else if segment.NoSpeechProb > 0.5 {
			logger.Wf(ctx, "might silent, prob=%v, file=%v, text=%v",
				segment.NoSpeechProb, outputFile, segment.Text)
		}
	}

	return &ASRResult{Text: resp.Text, Duration: time.Duration(resp.Duration * float64(time.Second))}, nil
}

type openaiChatService struct {
	// The AI configuration.
	conf openai.ClientConfig
	// The callback for the first response.
	onFirstResponse func(ctx context.Context, text string)
}

func (v *openaiChatService) RequestChat(ctx context.Context, sreq *StageRequest, stage *Stage, user *StageUser, taskCancel context.CancelFunc) error {
	if stage.previousUser != "" && stage.previousAssitant != "" {
		stage.histories = append(stage.histories, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: stage.previousUser,
		}, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: stage.previousAssitant,
		})

		for len(stage.histories) > stage.chatWindow*2 {
			stage.histories = stage.histories[1:]
		}
	}

	stage.previousUser = user.previousAsrText
	stage.previousAssitant = ""

	system := stage.prompt
	system += fmt.Sprintf(" Keep your reply neat, limiting the reply to %v words.", stage.replyLimit)
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: system},
	}

	messages = append(messages, stage.histories...)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: user.previousAsrText,
	})

	model := stage.chatModel
	maxTokens := 1024
	temperature := float32(0.9)
	logger.Tf(ctx, "AIChat is baseURL=%v, org=%v, model=%v, maxTokens=%v, temperature=%v, window=%v, histories=%v, system is %v",
		v.conf.BaseURL, v.conf.OrgID, model, maxTokens, temperature, stage.chatWindow, len(stage.histories), system)

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

	// Wait for AI got the first sentence response.
	aiFirstResponseCtx, aiFirstResponseCancel := context.WithCancel(ctx)
	defer aiFirstResponseCancel()

	go func() {
		defer gptChatStream.Close()
		if err := v.handle(ctx,
			stage, user, sreq, gptChatStream, aiFirstResponseCancel, taskCancel,
			func(sentence string) {
				stage.previousAssitant += sentence + " "
			},
		); err != nil {
			logger.Ef(ctx, "Handle stream failed, err %+v", err)
		}
	}()

	// Util AI generated the first sentence and commit the audio segment to TTS worker, we are allowed
	// response to client, and client will request the audio segments of this request and wait for segments
	// to be ready. If response before AI generated the first sentence, client will get nothing audio
	// segments and may start a new sentence.
	select {
	case <-ctx.Done():
	case <-aiFirstResponseCtx.Done():
	}

	return nil
}

func (v *openaiChatService) RequestPostProcess(ctx context.Context, sreq *StageRequest, stage *Stage, user *StageUser) error {
	if stage.postPreviousUser != "" && stage.postPreviousAssitant != "" {
		stage.postHistories = append(stage.postHistories, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: stage.postPreviousUser,
		}, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: stage.postPreviousAssitant,
		})

		for len(stage.postHistories) > stage.postChatWindow*2 {
			stage.postHistories = stage.postHistories[1:]
		}
	}

	stage.postPreviousUser = stage.previousAssitant
	stage.postPreviousAssitant = ""

	system := stage.postPrompt
	system += fmt.Sprintf(" Keep your reply neat, limiting the reply to %v words.", stage.postReplyLimit)
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: system},
	}

	messages = append(messages, stage.postHistories...)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: stage.previousAssitant,
	})

	model := stage.postChatModel
	maxTokens := 1024
	temperature := float32(0.9)
	logger.Tf(ctx, "AIPostProcess is baseURL=%v, org=%v, model=%v, maxTokens=%v, temperature=%v, window=%v, histories=%v, system is %v",
		v.conf.BaseURL, v.conf.OrgID, model, maxTokens, temperature, stage.postChatWindow, len(stage.postHistories), system)

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
		return errors.Wrapf(err, "create post-process")
	}

	// Wait for AI got the first sentence response.
	aiFirstResponseCtx, aiFirstResponseCancel := context.WithCancel(ctx)
	defer aiFirstResponseCancel()

	go func() {
		defer gptChatStream.Close()
		if err := v.handle(ctx,
			stage, user, sreq, gptChatStream, aiFirstResponseCancel, nil,
			func(sentence string) {
				stage.postPreviousAssitant += sentence + " "
			},
		); err != nil {
			logger.Ef(ctx, "Handle post-process stream failed, err %+v", err)
		}
	}()

	// Util AI generated the first sentence and commit the audio segment to TTS worker, we are allowed
	// response to client, and client will request the audio segments of this request and wait for segments
	// to be ready. If response before AI generated the first sentence, client will get nothing audio
	// segments and may start a new sentence.
	select {
	case <-ctx.Done():
	case <-aiFirstResponseCtx.Done():
	}

	return nil
}

func (v *openaiChatService) handle(
	ctx context.Context, stage *Stage, user *StageUser, sreq *StageRequest,
	gptChatStream *openai.ChatCompletionStream, aiFirstResponseCancel context.CancelFunc,
	taskCancel context.CancelFunc, onSentence func(string),
) error {
	defer aiFirstResponseCancel()
	if taskCancel != nil {
		defer taskCancel()
	}

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
			if stage.prefix != "" {
				filteredSentence = fmt.Sprintf("%v %v", stage.prefix, filteredSentence)
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

		// We have commit the segment to TTS worker, so we can return the response to client and allow
		// it to query audio segments immediately.
		if firstSentense {
			aiFirstResponseCancel()
		}

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
		if onSentence != nil && sentence != "" {
			onSentence(sentence)
		}
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

func NewOpenAITTSService(conf openai.ClientConfig) *openaiTTSService {
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
	// Whether merged to the next request.
	merged bool
	// Processing errors of this request.
	errs []error

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
		blocks := []string{
			"视频就拍到这里", "視頻就拍到這裡", "社群提供的字幕", "Amara.org社区", "by索兰娅",
			"请不吝点赞", "谢谢观看", "感谢观看", "订阅我的频道", "多謝您收睇", "多谢您收听",
			"多謝您的觀看", "多谢您的观看", "明镜与点点", "明鏡與點點", "视频就分享到这里",
			"关 注 雪 鱼 探 店", "下次见! 再见!",
		}
		for _, block := range blocks {
			if strings.Contains(asrText, block) {
				return errors.Errorf("badcase: %v", asrText)
			}
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
	if v.lastUploadAudio.After(v.lastSentence.Add(100 * time.Millisecond)) {
		return float64(v.lastUploadAudio.Sub(v.lastSentence)) / float64(time.Second)
	}
	return 0
}

func (v *StageRequest) exta() float64 {
	if v.lastExtractAudio.After(v.lastUploadAudio.Add(100 * time.Millisecond)) {
		return float64(v.lastExtractAudio.Sub(v.lastUploadAudio)) / float64(time.Second)
	}
	return 0
}

func (v *StageRequest) asr() float64 {
	if v.lastRequestASR.After(v.lastExtractAudio.Add(100 * time.Millisecond)) {
		return float64(v.lastRequestASR.Sub(v.lastExtractAudio)) / float64(time.Second)
	}
	return 0
}

func (v *StageRequest) chat() float64 {
	if v.lastRequestChat.After(v.lastRequestASR.Add(100 * time.Millisecond)) {
		return float64(v.lastRequestChat.Sub(v.lastRequestASR)) / float64(time.Second)
	}
	return 0
}

func (v *StageRequest) tts() float64 {
	if v.lastRequestTTS.After(v.lastRequestChat.Add(100 * time.Millisecond)) {
		return float64(v.lastRequestTTS.Sub(v.lastRequestChat)) / float64(time.Second)
	}
	return 0
}

func (v *StageRequest) download() float64 {
	if v.lastDownloadAudio.After(v.lastRequestTTS.Add(100 * time.Millisecond)) {
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
	// The username who send this message.
	Username string `json:"username,omitempty"`
	// Whether it's a new sentence.
	NewSentence bool `json:"sentence,omitempty"`

	// TODO: FIXME: Support message category, such as asr, chat, post, etc.
	// For role text.
	// The message content.
	Message string `json:"msg"`

	// For role audio.
	// The audio segment uuid.
	SegmentUUID string `json:"asid"`
	// Whether has audio file.
	HasAudioFile bool `json:"hasAudio"`
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

func (v *StageSubscriber) addUserTextMessage(rid, name, msg string) {
	v.messages = append(v.messages, &StageMessage{
		finished: true, MessageUUID: uuid.NewString(), subscriber: v,
		RequestUUID: rid, Role: "user", Message: msg, Username: name,
	})
}

// Create a robot empty message, to keep the order of messages.
func (v *StageSubscriber) createRobotEmptyMessage() *StageMessage {
	message := &StageMessage{
		finished: false, MessageUUID: uuid.NewString(), subscriber: v,
		Role: "robot",
	}
	v.messages = append(v.messages, message)
	return message
}

func (v *StageSubscriber) completeRobotAudioMessage(ctx context.Context, sreq *StageRequest, segment *AnswerSegment, message *StageMessage) {
	// Build a new copy file of ttsFile.
	var copyFile string
	if !segment.noTTS && segment.ttsFile != "" {
		ttsExt := path.Ext(segment.ttsFile)
		copyFile = fmt.Sprintf("%v-copy-%v%v", segment.ttsFile[:len(segment.ttsFile)-len(ttsExt)], v.spid, ttsExt)
	}

	// Copy the ttsFile to copyFile.
	if err := func() error {
		if copyFile == "" {
			return nil
		}

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

	message.segment = segment
	message.RequestUUID, message.SegmentUUID = sreq.rid, segment.asid
	message.Message, message.audioFile = segment.text, copyFile
	message.Username = v.stage.room.AIName
	message.NewSentence = segment.first

	// User may disable TTS, we only ship the text message to user.
	message.HasAudioFile = !segment.noTTS

	// Now, message is finished.
	message.finished = true

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

// The StageUser is the hosts on the stage, they have different config.
type StageUser struct {
	// The user UUID.
	UserID string `json:"userId"`
	// The username.
	Username string `json:"username,omitempty"`

	// The language the host use.
	Language string `json:"language,omitempty"`
	// The AI welcome voice url, binding to the language of user.
	Voice string `json:"voice,omitempty"`

	// Previous ASR text, to use as prompt for next ASR.
	previousAsrText string

	// Last update of user.
	update time.Time
	// The stage owner.
	stage *Stage
}

func (v *StageUser) Close() error {
	return nil
}

func (v *StageUser) Expired() bool {
	return time.Since(v.update) > 300*time.Second
}

func (v *StageUser) KeepAlive() {
	v.update = time.Now()
}

// TODO: Rename to a better name.
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

	// For post processing.
	// Previous chat text, to use as prompt for next chat.
	postPreviousUser, postPreviousAssitant string
	// The chat history, to use as prompt for next chat.
	postHistories []openai.ChatCompletionMessage
	// Reply words limit.
	postReplyLimit int
	// AI Chat model.
	postChatModel string
	// AI Chat message window.
	postChatWindow int
	// The AI prompt.
	postPrompt string

	// Shared configurations for both chat and post processing.
	// The prefix for TTS for the first sentence if too short.
	prefix string
	// The welcome voice url.
	voice string

	// Previous chat text, to use as prompt for next chat.
	previousUser, previousAssitant string
	// The chat history, to use as prompt for next chat.
	histories []openai.ChatCompletionMessage
	// The AI prompt.
	prompt string
	// Reply words limit.
	replyLimit int
	// AI Chat model.
	chatModel string
	// AI Chat message window.
	chatWindow int

	// The AI ASR language.
	asrLanguage string
	// The AI asr prompt type. user or user-ai.
	asrPrompt string

	// Whether enabled AI services.
	aiASREnabled  bool
	aiChatEnabled bool
	aiPostEnabled bool
	aiTtsEnabled  bool

	// The AI configuration.
	aiConfig openai.ClientConfig
	// The room it belongs to. Note that it's a caching object, update when updating the room. The room object
	// is not the same one, even the uuid is the same. The room is always available when stage is not expired.
	room *SrsLiveRoom
	// All the requests from user.
	requests []*StageRequest
	// All the subscribers binding to this stage.
	subscribers []*StageSubscriber
	// All the users binding to this stage.
	users []*StageUser
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

func (v Stage) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("enabled(asr:%v,chat:%v,post:%v,tts:%v)",
		v.aiASREnabled, v.aiChatEnabled, v.aiPostEnabled, v.aiTtsEnabled))
	sb.WriteString(fmt.Sprintf("sid:%v,asr:%v,asrp:%v",
		v.sid, v.asrLanguage, v.asrPrompt))
	if v.prefix != "" {
		sb.WriteString(fmt.Sprintf(",prefix:%v", v.prefix))
	}
	sb.WriteString(fmt.Sprintf(",voice=%v", v.voice))
	if v.aiChatEnabled {
		sb.WriteString(fmt.Sprintf(",chat(limit=%v,model=%v,window=%v,prompt:%v)",
			v.replyLimit, v.chatModel, v.chatWindow, v.prompt))
	}
	if v.aiPostEnabled {
		sb.WriteString(fmt.Sprintf(",post(limit=%v,model=%v,window=%v,prompt:%v)",
			v.postReplyLimit, v.postChatModel, v.postChatWindow, v.postPrompt))
	}
	return sb.String()
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

func helloVoiceFromLanguage(language string) string {
	if language == "zh" {
		return "hello-chinese.aac"
	}
	return "hello-english.aac"
}

func (v *Stage) UpdateFromRoom(room *SrsLiveRoom) {
	// Whether enabled.
	v.aiASREnabled = room.AIASREnabled
	v.aiChatEnabled = room.AIChatEnabled
	v.aiPostEnabled = room.AIPostEnabled
	v.aiTtsEnabled = room.AITTSEnabled

	// Create robot for the stage, which attach to a special room.
	v.voice = helloVoiceFromLanguage(room.AIASRLanguage)
	v.prompt = room.AIChatPrompt
	v.asrLanguage = room.AIASRLanguage
	v.asrPrompt = room.AIASRPrompt
	v.replyLimit = room.AIChatMaxWords
	v.chatModel = room.AIChatModel
	v.chatWindow = room.AIChatMaxWindow

	// Setup the post processing AI configuration.
	v.postReplyLimit = room.AIPostMaxWords
	v.postChatModel = room.AIPostModel
	v.postChatWindow = room.AIPostMaxWindow
	v.postPrompt = room.AIPostPrompt

	// Initialize the AI services.
	v.aiConfig = openai.DefaultConfig(room.AISecretKey)
	v.aiConfig.OrgID = room.AIOrganization
	v.aiConfig.BaseURL = room.AIBaseURL

	// Bind stage to room.
	room.StageUUID = v.sid
	v.room = room
}

func (v *Stage) Expired() bool {
	return time.Since(v.update) > 600*time.Second
}

func (v *Stage) KeepAlive() {
	v.update = time.Now()
}

func (v *Stage) addUser(user *StageUser) {
	v.users = append(v.users, user)
}

func (v *Stage) removeUser(user *StageUser) {
	for i, u := range v.users {
		if u.UserID == user.UserID {
			v.users = append(v.users[:i], v.users[i+1:]...)
			return
		}
	}
}

func (v *Stage) queryUser(userID string) *StageUser {
	if userID == "" {
		return nil
	}

	for _, user := range v.users {
		if user.UserID == userID {
			return user
		}
	}
	return nil
}

func (v *Stage) addRequest(request *StageRequest) {
	v.requests = append(v.requests, request)
}

func (v *Stage) queryRequest(rid string) *StageRequest {
	if rid == "" {
		return nil
	}

	for _, request := range v.requests {
		if request.rid == rid {
			return request
		}
	}
	return nil
}

func (v *Stage) queryPreviousNotMergedRequests(from *StageRequest) []*StageRequest {
	var requests []*StageRequest
	var matched bool
	for i := len(v.requests) - 1; i >= 0; i-- {
		request := v.requests[i]

		// Ignore requests after the from request.
		if request == from {
			matched = true
		}
		if !matched {
			continue
		}

		// Only return not merged requests.
		if request.merged {
			break
		}

		// Insert request to the head of list.
		requests = append([]*StageRequest{request}, requests...)
	}
	return requests
}

func (v *Stage) addSubscriber(subscriber *StageSubscriber) {
	v.subscribers = append(v.subscribers, subscriber)
}

func (v *Stage) querySubscriber(spid string) *StageSubscriber {
	if spid == "" {
		return nil
	}

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
	// Whether no tts file, as user disabled TTS for example.
	noTTS bool
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

func (v *TalkServer) QueryStageOfRoom(roomUUID string) *Stage {
	v.lock.Lock()
	defer v.lock.Unlock()

	for _, stage := range v.stages {
		if stage.room.UUID == roomUUID {
			return stage
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

		// Always make the segment ready or error, to update the stage to be ready.
		defer sreq.onSegmentReady(segment)

		if stage.aiTtsEnabled {
			ttsService := NewOpenAITTSService(stage.aiConfig)
			if err := ttsService.RequestTTS(ctx, func(ext string) string {
				segment.ttsFile = path.Join(aiTalkWorkDir,
					fmt.Sprintf("assistant-%v-sentence-%v-tts.%v", sreq.rid, segment.asid, ext),
				)
				return segment.ttsFile
			}, segment.text); err != nil {
				segment.err = err
			} else {
				segment.ready, segment.noTTS = true, false
				logger.Tf(ctx, "TTS: Complete rid=%v, asid=%v, file saved to %v, %v",
					sreq.rid, segment.asid, segment.ttsFile, segment.text)
			}
		} else {
			segment.ready, segment.noTTS = true, true
			logger.Tf(ctx, "TTS: Skip rid=%v, asid=%v, %v", sreq.rid, segment.asid, segment.text)
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
	// TODO: FIXME: Should use relative path, never expose absolute path to client.
	aiTalkWorkDir = path.Join(conf.Pwd, "containers/data/ai-talk")
	aiTalkExampleDir = path.Join(conf.Pwd, "containers/conf")
	logger.Tf(ctx, "AI-Talk init workDir=%v, examples=%v", aiTalkWorkDir, aiTalkExampleDir)

	ep := "/terraform/v1/ai-talk/stage/start"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		createNewStage := func(ctx context.Context, room *SrsLiveRoom) (*Stage, error) {
			ctx = logger.WithContext(ctx)

			stage := NewStage(func(stage *Stage) {
				stage.loggingCtx = ctx
				stage.UpdateFromRoom(room)
			})

			// Store the room, as we modify the stage UUID of room.
			if b, err := json.Marshal(room); err != nil {
				return nil, errors.Wrapf(err, "marshal room")
			} else if err := rdb.HSet(ctx, SRS_LIVE_ROOM, room.UUID, string(b)).Err(); err != nil {
				return nil, errors.Wrapf(err, "hset %v %v %v", SRS_LIVE_ROOM, room.UUID, string(b))
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

			return stage, nil
		}

		if err := func() error {
			var token string
			var roomUUID, roomToken string
			if err := ParseBody(ctx, r.Body, &struct {
				Token     *string `json:"token"`
				RoomUUID  *string `json:"room"`
				RoomToken *string `json:"roomToken"`
			}{
				Token:    &token,
				RoomUUID: &roomUUID, RoomToken: &roomToken,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			// Authenticate by bearer token if no room token
			if roomToken == "" {
				apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
				if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
					return errors.Wrapf(err, "authenticate")
				}
			}

			// TODO: FIXME: Should have room object in memory?
			// Get the room to verify user.
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

			// If assistant is disabled in room, fail.
			if !room.Assistant {
				return errors.Errorf("assistant disabled")
			}

			// Authenticate by room token if got one.
			if roomToken != "" && room.RoomToken != roomToken {
				return errors.Errorf("invalid room token %v", roomToken)
			}

			// Allow to reuse exists stage.
			stage := talkServer.QueryStage(room.StageUUID)
			if stage != nil {
				ctx = stage.loggingCtx
				logger.Tf(ctx, "Stage: Reuse stage sid=%v, all=%v", stage.sid, talkServer.CountStage())
			} else {
				var err error
				if stage, err = createNewStage(ctx, &room); err != nil {
					return err
				}
				ctx = stage.loggingCtx
			}

			// Keep alive the stage.
			stage.KeepAlive()

			// Create new user bind to this stage.
			user := &StageUser{
				UserID: uuid.NewString(), stage: stage,
				Language: stage.asrLanguage,
				Voice:    stage.voice,
			}
			user.KeepAlive()
			stage.addUser(user)

			go func() {
				defer user.Close()

				for ctx.Err() == nil {
					select {
					case <-ctx.Done():
					case <-time.After(3 * time.Second):
						if user.Expired() {
							logger.Tf(ctx, "Stage: Remove user=%v from stage sid=%v for expired, update=%v",
								user.UserID, stage.sid, user.update.Format(time.RFC3339))
							stage.removeUser(user)
							return
						}
					}
				}
			}()

			type StageResult struct {
				StageID   string `json:"sid"`
				RoomToken string `json:"roomToken"`
				UserID    string `json:"userId"`
				// AI Configurations.
				AIASREnabled bool `json:"aiAsrEnabled"`
			}
			r0 := &StageResult{
				StageID:   stage.sid,
				RoomToken: stage.room.RoomToken,
				UserID:    user.UserID,
				// AI Configurations.
				AIASREnabled: room.AIASREnabled,
			}

			ohttp.WriteData(ctx, w, r, r0)
			logger.Tf(ctx, "srs ai-talk create stage ok, room=%v, stage=%v, users=%v, subscribers=%v, requests=%v",
				room.UUID, stage.sid, len(stage.users), len(stage.subscribers), len(stage.requests))
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
			var sid, roomUUID, roomToken string
			if err := ParseBody(ctx, r.Body, &struct {
				Token     *string `json:"token"`
				RoomUUID  *string `json:"room"`
				RoomToken *string `json:"roomToken"`
				StageUUID *string `json:"sid"`
			}{
				Token: &token, StageUUID: &sid, RoomUUID: &roomUUID, RoomToken: &roomToken,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			// Authenticate by bearer token if no room token
			if roomToken == "" {
				apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
				if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
					return errors.Wrapf(err, "authenticate")
				}
			}

			if sid == "" {
				return errors.Errorf("empty sid")
			}

			stage := talkServer.QueryStage(sid)
			if stage == nil {
				return errors.Errorf("invalid sid %v", sid)
			}

			// Authenticate by room token if got one.
			if roomToken != "" && stage.room.RoomToken != roomToken {
				return errors.Errorf("invalid room token %v", roomToken)
			}
			if roomUUID != "" && stage.room.UUID != roomUUID {
				return errors.Errorf("invalid room %v", roomUUID)
			}

			// Keep alive the stage.
			stage.KeepAlive()
			// Switch to the context of stage.
			ctx = stage.loggingCtx

			// The rid is the request id, which identify this request, generally a question.
			sreq := &StageRequest{rid: uuid.NewString(), stage: stage}
			sreq.lastSentence = time.Now()
			// TODO: FIMXE: Should cleanup finished requests.
			stage.addRequest(sreq)

			// Keep alive the stage.
			stage.KeepAlive()

			// Response the request UUID and pulling the response.
			ohttp.WriteData(ctx, w, r, struct {
				RequestUUID string `json:"rid"`
			}{
				RequestUUID: sreq.rid,
			})
			logger.Tf(ctx, "ai-talk new conversation, room=%v, sid=%v, rid=%v", roomUUID, sid, sreq.rid)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai-talk/stage/upload"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		var sreq *StageRequest
		if err := func() error {
			var token string
			var sid, rid, userID string
			var roomUUID, roomToken string
			var userMayInput float64
			var audioBase64Data, textMessage string
			var mergeMessages int
			if err := ParseBody(ctx, r.Body, &struct {
				Token        *string  `json:"token"`
				RoomUUID     *string  `json:"room"`
				RoomToken    *string  `json:"roomToken"`
				StageUUID    *string  `json:"sid"`
				UserID       *string  `json:"userId"`
				RequestUUID  *string  `json:"rid"`
				UserMayInput *float64 `json:"umi"`
				AudioData    *string  `json:"audio"`
				TextMessage  *string  `json:"text"`
				// Merge ASR text of conversations, which is small duration audio segment.
				MergeMessages *int `json:"mergeMessages"`
			}{
				Token: &token, StageUUID: &sid, UserID: &userID, RequestUUID: &rid,
				UserMayInput: &userMayInput, TextMessage: &textMessage, AudioData: &audioBase64Data,
				RoomUUID: &roomUUID, RoomToken: &roomToken, MergeMessages: &mergeMessages,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			// Authenticate by bearer token if no room token
			if roomToken == "" {
				apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
				if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
					return errors.Wrapf(err, "authenticate")
				}
			}

			if sid == "" {
				return errors.Errorf("empty sid")
			}
			if rid == "" {
				return errors.Errorf("empty rid")
			}
			if userID == "" {
				return errors.Errorf("empty userId")
			}
			if audioBase64Data == "" && textMessage == "" {
				return errors.Errorf("empty audio and text")
			}

			stage := talkServer.QueryStage(sid)
			if stage == nil {
				return errors.Errorf("invalid sid %v", sid)
			}

			// Authenticate by room token if got one.
			if roomToken != "" && stage.room.RoomToken != roomToken {
				return errors.Errorf("invalid room token %v", roomToken)
			}
			if roomUUID != "" && stage.room.UUID != roomUUID {
				return errors.Errorf("invalid room %v", roomUUID)
			}

			// Keep alive the stage.
			stage.KeepAlive()
			// Switch to the context of stage.
			ctx = stage.loggingCtx

			// Query the request.
			sreq = stage.queryRequest(rid)
			if sreq == nil {
				return errors.Errorf("invalid sid=%v, rid=%v", sid, rid)
			}

			user := stage.queryUser(userID)
			if user == nil {
				return errors.Errorf("invalid user %v of sid %v", userID, sid)
			}
			user.KeepAlive()

			// The rid is the request id, which identify this request, generally a question.
			defer sreq.FastDispose()

			sreq.inputFile = path.Join(aiTalkWorkDir, fmt.Sprintf("assistant-%v-input.audio", sreq.rid))
			logger.Tf(ctx, "Stage: Got question sid=%v, rid=%v, user=%v, umi=%v, input=%v",
				sid, sreq.rid, userID, userMayInput, sreq.inputFile)

			// Whether user input audio.
			if audioBase64Data != "" {
				// Save audio input to file.
				if err := sreq.receiveInputFile(ctx, audioBase64Data); err != nil {
					return errors.Wrapf(err, "save %vB audio to file %v", len(audioBase64Data), sreq.inputFile)
				}

				// Do ASR, convert to text.
				asrLanguage := ChooseNotEmpty(user.Language, stage.asrLanguage)
				if err := sreq.asrAudioToText(ctx, stage.aiConfig, asrLanguage, user.previousAsrText); err != nil {
					return errors.Wrapf(err, "asr lang=%v, previous=%v", asrLanguage, user.previousAsrText)
				}
				logger.Tf(ctx, "ASR ok, sid=%v, rid=%v, user=%v, lang=%v, prompt=<%v>, resp is <%v>",
					sid, sreq.rid, userID, asrLanguage, user.previousAsrText, sreq.asrText)
			} else {
				// Directly update the time for stat.
				sreq.lastUploadAudio = time.Now()
				sreq.lastExtractAudio = time.Now()
				sreq.lastRequestASR = time.Now()
			}

			// Handle user input text.
			if textMessage != "" {
				sreq.asrText = textMessage
				logger.Tf(ctx, "Text ok, sid=%v, rid=%v, user=%v, text=%v",
					sid, sreq.rid, userID, sreq.asrText)
			}

			// Important trace log.
			user.previousAsrText = sreq.asrText
			logger.Tf(ctx, "You: %v", sreq.asrText)

			// Notify all subscribers about the ASR text.
			for _, subscriber := range stage.subscribers {
				subscriber.addUserTextMessage(sreq.rid, user.Username, sreq.asrText)
			}

			// Keep alive the stage.
			stage.KeepAlive()

			// If merge conversation to next one, we do not submit to chat and post processing.
			conversations := stage.queryPreviousNotMergedRequests(sreq)
			mergeToNextConversation := mergeMessages > 0 && len(conversations) < mergeMessages
			if !mergeToNextConversation {
				sreq.merged, user.previousAsrText = true, ""
				// Generate the merged text for chat input.
				for _, conversation := range conversations {
					// If request has error, such as silent or other error, ignore the text.
					if conversation.errs != nil {
						continue
					}
					user.previousAsrText += conversation.asrText
				}
			}

			// Do chat, get the response in stream.
			chatTaskCtx, chatTaskCancel := context.WithCancel(context.Background())
			if !mergeToNextConversation && stage.aiChatEnabled {
				chatService := &openaiChatService{
					conf: stage.aiConfig,
					onFirstResponse: func(ctx context.Context, text string) {
						sreq.lastRequestChat = time.Now()
						sreq.lastRobotFirstText = text
					},
				}
				if err := chatService.RequestChat(ctx, sreq, stage, user, chatTaskCancel); err != nil {
					return errors.Wrapf(err, "chat")
				}
			}

			// Do AI post-process ,get the response in stream.
			// TODO: FIXME: Should use a goroutine to do post-process.
			if !mergeToNextConversation && stage.aiChatEnabled && stage.aiPostEnabled {
				// Wait for chat to be completed.
				select {
				case <-ctx.Done():
				case <-chatTaskCtx.Done():
				}

				// Start post processing task.
				chatService := &openaiChatService{
					conf: stage.aiConfig,
				}
				if err := chatService.RequestPostProcess(ctx, sreq, stage, user); err != nil {
					return errors.Wrapf(err, "post-process")
				}
			}

			// Response the request UUID and pulling the response.
			ohttp.WriteData(ctx, w, r, struct {
				RequestUUID string `json:"rid"`
				ASR         string `json:"asr"`
			}{
				RequestUUID: sreq.rid,
				ASR:         sreq.asrText,
			})
			logger.Tf(ctx, "srs ai-talk stage upload ok, sid=%v, rid=%v, user=%v, asr=%v",
				sid, sreq.rid, userID, sreq.asrText)
			return nil
		}(); err != nil {
			if sreq != nil {
				sreq.errs = append(sreq.errs, err)
			}
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai-talk/stage/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var sid, rid string
			var roomUUID, roomToken string
			if err := ParseBody(ctx, r.Body, &struct {
				Token       *string `json:"token"`
				RoomUUID    *string `json:"room"`
				RoomToken   *string `json:"roomToken"`
				StageUUID   *string `json:"sid"`
				RequestUUID *string `json:"rid"`
			}{
				Token: &token, StageUUID: &sid, RequestUUID: &rid,
				RoomUUID: &roomUUID, RoomToken: &roomToken,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			// Authenticate by bearer token if no room token
			if roomToken == "" {
				apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
				if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
					return errors.Wrapf(err, "authenticate")
				}
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

			// Authenticate by room token if got one.
			if roomToken != "" && stage.room.RoomToken != roomToken {
				return errors.Errorf("invalid room token %v", roomToken)
			}
			if roomUUID != "" && stage.room.UUID != roomUUID {
				return errors.Errorf("invalid room %v", roomUUID)
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
				Finished: sreq.finished,
			})

			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai-talk/stage/hello-voices/"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			filename := r.URL.Path[len("/terraform/v1/ai-talk/stage/hello-voices/"):]
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
			http.ServeFile(w, r, path.Join(aiTalkExampleDir, filename))
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

			// Note that when verifying stage, there may be no stage exists, so we must fetch the room token
			// from redis, should never try to use cached token from stage.
			var room SrsLiveRoom
			if r0, err := rdb.HGet(ctx, SRS_LIVE_ROOM, roomUUID).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_LIVE_ROOM, roomUUID)
			} else if r0 == "" {
				return errors.Errorf("live room %v not exists", roomUUID)
			} else if err = json.Unmarshal([]byte(r0), &room); err != nil {
				return errors.Wrapf(err, "unmarshal %v %v", roomUUID, r0)
			}

			if room.RoomToken != roomToken {
				return errors.Errorf("invalid room token")
			}

			ohttp.WriteData(ctx, w, r, nil)
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
			var roomUUID, roomToken string
			if err := ParseBody(ctx, r.Body, &struct {
				Token     *string `json:"token"`
				RoomUUID  *string `json:"room"`
				RoomToken *string `json:"roomToken"`
			}{
				Token:    &token,
				RoomUUID: &roomUUID, RoomToken: &roomToken,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			// Authenticate by bearer token if no room token
			if roomToken == "" {
				apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
				if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
					return errors.Wrapf(err, "authenticate")
				}
			}

			// When subscribing, there must be a valid stage exists, so we can use cached stage to verify the
			// room token, to make sure the stage started before any subscribers.
			stage := talkServer.QueryStageOfRoom(roomUUID)
			if stage == nil {
				return errors.Errorf("no stage in room %v", roomUUID)
			}

			// Authenticate by room token if got one.
			if roomToken != "" && stage.room.RoomToken != roomToken {
				return errors.Errorf("invalid room token %v", roomToken)
			}

			// TODO: FIXME: Cleanup subscribers for a room.
			ctx = logger.WithContext(ctx)
			subscriber := NewStageSubscriber(func(subscriber *StageSubscriber) {
				subscriber.loggingCtx = ctx
				subscriber.room = stage.room

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

			type SubscribeResult struct {
				StageID      string `json:"sid"`
				SubscriberID string `json:"spid"`
				Voice        string `json:"voice"`
			}
			r0 := &SubscribeResult{
				StageID:      stage.sid,
				SubscriberID: subscriber.spid,
				Voice:        stage.voice,
			}

			ohttp.WriteData(ctx, w, r, &r0)
			logger.Tf(ctx, "Stage: create subscriber ok, room=%v, sid=%v, spid=%v",
				stage.room.UUID, stage.sid, subscriber.spid)
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
			var roomUUID, roomToken string
			// Optional, stream hosts has a userID, no this field for subscribers. Use as heartbeat to
			// make user keep in alive.
			var userID string
			if err := ParseBody(ctx, r.Body, &struct {
				Token        *string `json:"token"`
				RoomUUID     *string `json:"room"`
				RoomToken    *string `json:"roomToken"`
				StageUUID    *string `json:"sid"`
				SubscriberID *string `json:"spid"`
				UserID       *string `json:"userId"`
			}{
				Token: &token, StageUUID: &sid, SubscriberID: &spid,
				RoomUUID: &roomUUID, RoomToken: &roomToken,
				UserID: &userID,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			// Authenticate by bearer token if no room token
			if roomToken == "" {
				apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
				if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
					return errors.Wrapf(err, "authenticate")
				}
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

			// Authenticate by room token if got one.
			if roomToken != "" && stage.room.RoomToken != roomToken {
				return errors.Errorf("invalid room token %v", roomToken)
			}
			if roomUUID != "" && stage.room.UUID != roomUUID {
				return errors.Errorf("invalid room %v", roomUUID)
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

			// For stream hosts, we should keep the user alive.
			if user := stage.queryUser(userID); user != nil {
				user.KeepAlive()
			}

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

	finishAudioSegment := func(segment *AnswerSegment) {
		if segment == nil || segment.logged {
			return
		}

		// Only log the first segment.
		segment.logged = true
		if !segment.first {
			return
		}

		// Time cost logging.
		sreq := segment.request
		sreq.lastDownloadAudio = time.Now()
		speech := float64(sreq.lastAsrDuration) / float64(time.Second)
		logger.Tf(ctx, "Elapsed cost total=%.1fs, steps=[upload=%.1fs,exta=%.1fs,asr=%.1fs,chat=%.1fs,tts=%.1fs,download=%.1fs], ask=%v, speech=%.1fs, answer=%v",
			sreq.total(), sreq.upload(), sreq.exta(), sreq.asr(), sreq.chat(), sreq.tts(), sreq.download(),
			sreq.lastRequestAsrText, speech, sreq.lastRobotFirstText)

		// Important trace log. Note that browser may request multiple times, so we only log for the first
		// request to reduce logs.
		logger.Tf(ctx, "Bot: %v", segment.text)
	}

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

			roomUUID := q.Get("room")
			if roomUUID == "" {
				return errors.Errorf("empty room")
			}

			roomToken := q.Get("roomToken")
			if roomToken == "" {
				return errors.Errorf("empty roomToken")
			}

			stage := talkServer.QueryStage(sid)
			if stage == nil {
				return errors.Errorf("invalid sid %v", sid)
			}

			// Authenticate by room token if got one.
			if roomToken != "" && stage.room.RoomToken != roomToken {
				return errors.Errorf("invalid room token %v", roomToken)
			}
			if roomUUID != "" && stage.room.UUID != roomUUID {
				return errors.Errorf("invalid room %v", roomUUID)
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
			finishAudioSegment(answer.segment)

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
			var roomUUID, roomToken string
			if err := ParseBody(ctx, r.Body, &struct {
				Token            *string `json:"token"`
				RoomUUID         *string `json:"room"`
				RoomToken        *string `json:"roomToken"`
				StageUUID        *string `json:"sid"`
				SubscriberID     *string `json:"spid"`
				AudioSegmentUUID *string `json:"asid"`
			}{
				Token: &token, StageUUID: &sid, SubscriberID: &spid, AudioSegmentUUID: &asid,
				RoomUUID: &roomUUID, RoomToken: &roomToken,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			if roomToken == "" {
				apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
				if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
					return errors.Wrapf(err, "authenticate")
				}
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

			// Authenticate by room token if got one.
			if roomToken != "" && stage.room.RoomToken != roomToken {
				return errors.Errorf("invalid room token %v", roomToken)
			}
			if roomUUID != "" && stage.room.UUID != roomUUID {
				return errors.Errorf("invalid room %v", roomUUID)
			}

			subscriber := stage.querySubscriber(spid)
			if subscriber == nil {
				return errors.Errorf("invalid spid %v of sid %v", spid, sid)
			}

			// If no audio file, we stat the time cost when remove the segment.
			if answer := subscriber.queryAudioFile(asid); answer != nil {
				finishAudioSegment(answer.segment)
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

	ep = "/terraform/v1/ai-talk/user/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var sid, userID string
			var roomUUID, roomToken string
			if err := ParseBody(ctx, r.Body, &struct {
				Token     *string `json:"token"`
				RoomUUID  *string `json:"room"`
				RoomToken *string `json:"roomToken"`
				StageUUID *string `json:"sid"`
				UserID    *string `json:"userId"`
			}{
				Token: &token, StageUUID: &sid, UserID: &userID,
				RoomUUID: &roomUUID, RoomToken: &roomToken,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			if roomToken == "" {
				apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
				if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
					return errors.Wrapf(err, "authenticate")
				}
			}

			if sid == "" {
				return errors.Errorf("empty sid")
			}
			if userID == "" {
				return errors.Errorf("empty userId")
			}

			stage := talkServer.QueryStage(sid)
			if stage == nil {
				return errors.Errorf("invalid sid %v", sid)
			}

			// Authenticate by room token if got one.
			if roomToken != "" && stage.room.RoomToken != roomToken {
				return errors.Errorf("invalid room token %v", roomToken)
			}
			if roomUUID != "" && stage.room.UUID != roomUUID {
				return errors.Errorf("invalid room %v", roomUUID)
			}

			user := stage.queryUser(userID)
			if user == nil {
				return errors.Errorf("invalid user %v of sid %v", userID, sid)
			}

			// Keep alive the stage.
			stage.KeepAlive()
			user.KeepAlive()
			// Switch to the context of stage.
			ctx = stage.loggingCtx
			logger.Tf(ctx, "Stage: Query user room=%v, sid=%v, user=%v",
				stage.room.UUID, sid, userID)

			ohttp.WriteData(ctx, w, r, user)
			logger.Tf(ctx, "srs ai-talk query user ok")
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai-talk/user/update"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var sid, userID, username, userLanguage string
			var roomUUID, roomToken string
			if err := ParseBody(ctx, r.Body, &struct {
				Token        *string `json:"token"`
				RoomUUID     *string `json:"room"`
				RoomToken    *string `json:"roomToken"`
				StageUUID    *string `json:"sid"`
				UserID       *string `json:"userId"`
				Username     *string `json:"name"`
				UserLanguage *string `json:"lang"`
			}{
				Token: &token, StageUUID: &sid, UserID: &userID,
				Username: &username, UserLanguage: &userLanguage,
				RoomUUID: &roomUUID, RoomToken: &roomToken,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			if roomToken == "" {
				apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
				if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
					return errors.Wrapf(err, "authenticate")
				}
			}

			if sid == "" {
				return errors.Errorf("empty sid")
			}
			if userID == "" {
				return errors.Errorf("empty userId")
			}

			stage := talkServer.QueryStage(sid)
			if stage == nil {
				return errors.Errorf("invalid sid %v", sid)
			}

			// Authenticate by room token if got one.
			if roomToken != "" && stage.room.RoomToken != roomToken {
				return errors.Errorf("invalid room token %v", roomToken)
			}
			if roomUUID != "" && stage.room.UUID != roomUUID {
				return errors.Errorf("invalid room %v", roomUUID)
			}

			user := stage.queryUser(userID)
			if user == nil {
				return errors.Errorf("invalid user %v of sid %v", userID, sid)
			}

			// Keep alive the stage.
			stage.KeepAlive()
			user.KeepAlive()
			// Switch to the context of stage.
			ctx = stage.loggingCtx
			logger.Tf(ctx, "Stage: Update user room=%v, sid=%v, user=%v, name=%v, lang=%v",
				stage.room.UUID, sid, userID, username, userLanguage)

			user.Username = username
			user.Language = userLanguage
			user.Voice = helloVoiceFromLanguage(userLanguage)

			ohttp.WriteData(ctx, w, r, user)
			logger.Tf(ctx, "srs ai-talk update user ok, sid=%v, user=%v", sid, userID)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	return nil
}
