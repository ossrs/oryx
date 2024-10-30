// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"

	// From ossrs.
	"github.com/ossrs/go-oryx-lib/errors"
	ohttp "github.com/ossrs/go-oryx-lib/http"
	"github.com/ossrs/go-oryx-lib/logger"

	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

// The total segments in callback HLS.
const maxCallbackSegments = 9

var ocrWorker *OCRWorker

type OCRWorker struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// The global OCR task, only support one OCR task.
	task *OCRTask

	// Got message from SRS, a new TS segment file is generated.
	tsfiles chan *SrsOnHlsObject
}

func NewOCRWorker() *OCRWorker {
	v := &OCRWorker{
		// TS files.
		tsfiles: make(chan *SrsOnHlsObject, 1024),
	}
	v.task = NewOCRTask()
	v.task.ocrWorker = v
	return v
}

func (v *OCRWorker) Handle(ctx context.Context, handler *http.ServeMux) error {
	ep := "/terraform/v1/ai/ocr/query"
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

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			config := NewOCRConfig()
			if err := config.Load(ctx); err != nil {
				return errors.Wrapf(err, "load config")
			}

			type QueryResponse struct {
				Config *OCRConfig `json:"config"`
				Task   struct {
					UUID string `json:"uuid"`
				} `json:"task"`
			}

			resp := &QueryResponse{
				Config: config,
			}
			resp.Task.UUID = v.task.UUID

			ohttp.WriteData(ctx, w, r, resp)
			logger.Tf(ctx, "ocr query ok, config=<%v>, uuid=%v, token=%vB",
				config, v.task.UUID, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai/ocr/apply"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var uuid string
			var config OCRConfig
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				UUID  *string `json:"uuid"`
				*OCRConfig
			}{
				Token: &token,
				UUID:  &uuid, OCRConfig: &config,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			// Not required yet.
			if uuid != v.task.UUID {
				logger.Wf(ctx, "ocr ignore uuid mismatch, query=%v, task=%v", uuid, v.task.UUID)
			}

			if err := config.Save(ctx); err != nil {
				return errors.Wrapf(err, "save config")
			}

			if err := v.task.restart(ctx); err != nil {
				return errors.Wrapf(err, "restart task %v", config.String())
			}

			type ApplyResponse struct {
				UUID string `json:"uuid"`
			}
			ohttp.WriteData(ctx, w, r, &ApplyResponse{
				UUID: v.task.UUID,
			})
			logger.Tf(ctx, "ocr apply ok, config=<%v>, uuid=%v, token=%vB",
				config, v.task.UUID, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai/ocr/check"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var ocrConfig OCRConfig
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				*OCRConfig
			}{
				Token:     &token,
				OCRConfig: &ocrConfig,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			// Query whisper-1 model detail.
			var config openai.ClientConfig
			config = openai.DefaultConfig(ocrConfig.AISecretKey)
			config.BaseURL = ocrConfig.AIBaseURL

			ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
			defer cancel()

			client := openai.NewClientWithConfig(config)
			model, err := client.GetModel(ctx, "whisper-1")
			if err != nil {
				return errors.Wrapf(err, "query model whisper-1")
			}

			// Start a chat, to check whether the billing is expired.
			resp, err := client.CreateChatCompletion(
				ctx, openai.ChatCompletionRequest{
					Model: openai.GPT4o,
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: "Hello!",
						},
					},
					MaxTokens: 50,
				},
			)
			if err != nil {
				return errors.Wrapf(err, "create chat")
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "ocr check ok, config=<%v>, model=<%v>, msg=<%v>, token=%vB",
				ocrConfig, model.ID, resp.Choices[0].Message.Content, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai/ocr/reset"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var uuid string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				UUID  *string `json:"uuid"`
			}{
				Token: &token,
				UUID:  &uuid,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if uuid != v.task.UUID {
				return errors.Errorf("invalid uuid %v", uuid)
			}

			if err := v.task.reset(ctx); err != nil {
				return errors.Wrapf(err, "restart task %v", uuid)
			}

			type ResetResponse struct {
				UUID string `json:"uuid"`
			}
			ohttp.WriteData(ctx, w, r, &ResetResponse{
				UUID: v.task.UUID,
			})
			logger.Tf(ctx, "ocr reset ok, uuid=%v, new=%v, token=%vB", uuid, v.task.UUID, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai/ocr/live-queue"
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

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			type Segment struct {
				TsID     string  `json:"tsid"`
				SeqNo    uint64  `json:"seqno"`
				URL      string  `json:"url"`
				Duration float64 `json:"duration"`
				Size     uint64  `json:"size"`
			}
			type LiveQueueResponse struct {
				Segments []*Segment `json:"segments"`
				Count    int        `json:"count"`
			}
			res := &LiveQueueResponse{}

			segments := v.task.liveSegments()
			for _, segment := range segments {
				res.Segments = append(res.Segments, []*Segment{&Segment{
					TsID:     segment.TsFile.TsID,
					SeqNo:    segment.TsFile.SeqNo,
					URL:      segment.TsFile.URL,
					Duration: segment.TsFile.Duration,
					Size:     segment.TsFile.Size,
				}}...)
			}

			res.Count = len(res.Segments)

			ohttp.WriteData(ctx, w, r, res)
			logger.Tf(ctx, "ocr query live ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai/ocr/ocr-queue"
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

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			type Segment struct {
				TsID     string  `json:"tsid"`
				SeqNo    uint64  `json:"seqno"`
				URL      string  `json:"url"`
				Duration float64 `json:"duration"`
				Size     uint64  `json:"size"`
				// The source ts file.
				SourceTsID string `json:"stsid"`
				// The cost in ms to extract image.
				ExtractImageCost int32 `json:"eic"`
			}
			type OCRQueueResponse struct {
				Segments []*Segment `json:"segments"`
				Count    int        `json:"count"`
			}
			res := &OCRQueueResponse{}

			segments := v.task.ocrSegments()
			for _, segment := range segments {
				res.Segments = append(res.Segments, []*Segment{&Segment{
					TsID:     segment.ImageFile.TsID,
					SeqNo:    segment.ImageFile.SeqNo,
					URL:      segment.ImageFile.File,
					Duration: segment.ImageFile.Duration,
					Size:     segment.ImageFile.Size,
					// The source ts file.
					SourceTsID: segment.TsFile.TsID,
					// The cost in ms to extract image.
					ExtractImageCost: int32(segment.CostExtractImage.Milliseconds()),
				}}...)
			}

			res.Count = len(res.Segments)

			ohttp.WriteData(ctx, w, r, res)
			logger.Tf(ctx, "ocr query ocr ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai/ocr/callback-queue"
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

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			type Segment struct {
				TsID     string  `json:"tsid"`
				SeqNo    uint64  `json:"seqno"`
				URL      string  `json:"url"`
				Duration float64 `json:"duration"`
				Size     uint64  `json:"size"`
				// The source ts file.
				SourceTsID string `json:"stsid"`
				// The cost in ms to extract image.
				ExtractImageCost int32 `json:"eic"`
				// The OCR text result.
				OCRText string `json:"ocr"`
				// The cost in ms to do OCR.
				OCRCost int32 `json:"ocrc"`
			}
			type OCRQueueResponse struct {
				Segments []*Segment `json:"segments"`
				Count    int        `json:"count"`
			}
			res := &OCRQueueResponse{}

			segments := v.task.callbackSegments()
			for _, segment := range segments {
				res.Segments = append(res.Segments, []*Segment{&Segment{
					TsID:     segment.ImageFile.TsID,
					SeqNo:    segment.ImageFile.SeqNo,
					URL:      segment.ImageFile.File,
					Duration: segment.ImageFile.Duration,
					Size:     segment.ImageFile.Size,
					// The source ts file.
					SourceTsID: segment.TsFile.TsID,
					// The cost in ms to extract image.
					ExtractImageCost: int32(segment.CostExtractImage.Milliseconds()),
					// The OCR text result.
					OCRText: segment.OCRText,
					// The cost in ms to do OCR.
					OCRCost: int32(segment.CostOCR.Milliseconds()),
				}}...)
			}

			res.Count = len(res.Segments)

			ohttp.WriteData(ctx, w, r, res)
			logger.Tf(ctx, "ocr query callback ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai/ocr/cleanup-queue"
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

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			type Segment struct {
				TsID     string  `json:"tsid"`
				SeqNo    uint64  `json:"seqno"`
				URL      string  `json:"url"`
				Duration float64 `json:"duration"`
				Size     uint64  `json:"size"`
				// The source ts file.
				SourceTsID string `json:"stsid"`
				// The cost in ms to extract image.
				ExtractImageCost int32 `json:"eic"`
				// The OCR text result.
				OCRText string `json:"ocr"`
				// The cost in ms to do OCR.
				OCRCost int32 `json:"ocrc"`
				// The cost in ms to do callback.
				CallbackCost int32 `json:"cbc"`
			}
			type OCRQueueResponse struct {
				Segments []*Segment `json:"segments"`
				Count    int        `json:"count"`
			}
			res := &OCRQueueResponse{}

			segments := v.task.cleanupSegments()
			for _, segment := range segments {
				res.Segments = append(res.Segments, []*Segment{&Segment{
					TsID:     segment.ImageFile.TsID,
					SeqNo:    segment.ImageFile.SeqNo,
					URL:      segment.ImageFile.File,
					Duration: segment.ImageFile.Duration,
					Size:     segment.ImageFile.Size,
					// The source ts file.
					SourceTsID: segment.TsFile.TsID,
					// The cost in ms to extract image.
					ExtractImageCost: int32(segment.CostExtractImage.Milliseconds()),
					// The OCR text result.
					OCRText: segment.OCRText,
					// The cost in ms to do OCR.
					OCRCost: int32(segment.CostOCR.Milliseconds()),
					// The cost in msg to do callback.
					CallbackCost: int32(segment.CostCallback.Milliseconds()),
				}}...)
			}

			res.Count = len(res.Segments)

			ohttp.WriteData(ctx, w, r, res)
			logger.Tf(ctx, "ocr query cleanup ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/ai/ocr/image/"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			// Format is /image/:uuid.jpg
			filename := r.URL.Path[len("/terraform/v1/ai/ocr/image/"):]
			// Format is :uuid.jpg
			uuid := filename[:len(filename)-len(path.Ext(filename))]
			if len(uuid) == 0 {
				return errors.Errorf("invalid uuid %v from %v of %v", uuid, filename, r.URL.Path)
			}

			imageFilePath := path.Join("ocr", fmt.Sprintf("%v.jpg", uuid))
			if _, err := os.Stat(imageFilePath); err != nil {
				return errors.Wrapf(err, "no image file %v", imageFilePath)
			}

			if tsFile, err := os.Open(imageFilePath); err != nil {
				return errors.Wrapf(err, "open file %v", imageFilePath)
			} else {
				defer tsFile.Close()
				w.Header().Set("Content-Type", "image/jpeg")
				io.Copy(w, tsFile)
			}

			logger.Tf(ctx, "ocr preview image ok, uuid=%v", uuid)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	return nil
}

func (v *OCRWorker) Enabled() bool {
	return v.task.enabled()
}

func (v *OCRWorker) OnHlsTsMessage(ctx context.Context, msg *SrsOnHlsMessage, data []byte) error {
	// Ignore if not natch the task config.
	if !v.task.match(msg) {
		return nil
	}

	// Copy the ts file to temporary cache dir.
	tsid := fmt.Sprintf("%v-org-%v", msg.SeqNo, uuid.NewString())
	tsfile := path.Join("ocr", fmt.Sprintf("%v.ts", tsid))

	if file, err := os.Create(tsfile); err != nil {
		return errors.Wrapf(err, "create file %v error", tsfile)
	} else {
		defer file.Close()
		io.Copy(file, bytes.NewReader(data))
	}

	// Create a local ts file object.
	tsFile := &TsFile{
		TsID:     tsid,
		URL:      msg.URL,
		SeqNo:    msg.SeqNo,
		Duration: msg.Duration,
		Size:     uint64(len(data)),
		File:     tsfile,
	}

	// Notify worker asynchronously.
	// TODO: FIXME: Should cleanup the temporary file when restart.
	go func() {
		select {
		case <-ctx.Done():
		case v.tsfiles <- &SrsOnHlsObject{Msg: msg, TsFile: tsFile}:
		}
	}()
	return nil
}

func (v *OCRWorker) Close() error {
	if v.cancel != nil {
		v.cancel()
	}

	v.wg.Wait()
	return nil
}

func (v *OCRWorker) Start(ctx context.Context) error {
	wg := &v.wg

	ctx, cancel := context.WithCancel(ctx)
	v.cancel = cancel

	ctx = logger.WithContext(ctx)
	logger.Tf(ctx, "ocr start a worker")

	// Load task from redis and continue to run the task.
	if objs, err := rdb.HGetAll(ctx, SRS_OCR_TASK).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hgetall %v", SRS_OCR_TASK)
	} else if len(objs) != 1 {
		// Only support one task right now.
		if err = rdb.Del(ctx, SRS_OCR_TASK).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "del %v", SRS_OCR_TASK)
		}
	} else {
		for uuid, obj := range objs {
			logger.Tf(ctx, "Load task %v object %v", uuid, obj)

			if err = json.Unmarshal([]byte(obj), v.task); err != nil {
				return errors.Wrapf(err, "unmarshal %v %v", uuid, obj)
			}

			break
		}
	}

	// Start global ocr task.
	wg.Add(1)
	go func() {
		defer wg.Done()

		task := v.task
		for ctx.Err() == nil {
			var duration time.Duration
			if err := task.Run(ctx); err != nil {
				logger.Wf(ctx, "ocr: run task %v err %+v", task.String(), err)
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

	// Consume all ts files by task.
	wg.Add(1)
	go func() {
		defer wg.Done()

		task := v.task
		for ctx.Err() == nil {
			select {
			case <-ctx.Done():
			case msg := <-v.tsfiles:
				if err := task.OnTsSegment(ctx, msg); err != nil {
					logger.Wf(ctx, "ocr: task %v on hls ts message %v err %+v", task.String(), msg.String(), err)
				}
			}
		}
	}()

	// Watch for new stream.
	wg.Add(1)
	go func() {
		defer wg.Done()

		task := v.task
		for ctx.Err() == nil {
			var duration time.Duration
			if err := task.WatchNewStream(ctx); err != nil {
				logger.Wf(ctx, "ocr: task %v watch new stream err %+v", task.String(), err)
				duration = 10 * time.Second
			} else {
				duration = 200 * time.Millisecond
			}

			select {
			case <-ctx.Done():
			case <-time.After(duration):
			}
		}
	}()

	// Drive the live queue to OCR.
	wg.Add(1)
	go func() {
		defer wg.Done()

		task := v.task
		for ctx.Err() == nil {
			var duration time.Duration
			if err := task.DriveLiveQueue(ctx); err != nil {
				logger.Wf(ctx, "ocr: task %v drive live queue err %+v", task.String(), err)
				duration = 10 * time.Second
			} else {
				duration = 200 * time.Millisecond
			}

			select {
			case <-ctx.Done():
			case <-time.After(duration):
			}
		}
	}()

	// Drive the OCR queue to correct queue.
	wg.Add(1)
	go func() {
		defer wg.Done()

		task := v.task
		for ctx.Err() == nil {
			var duration time.Duration
			if err := task.DriveOCRQueue(ctx); err != nil {
				logger.Wf(ctx, "ocr: task %v drive ocr queue err %+v", task.String(), err)
				duration = 10 * time.Second
			} else {
				duration = 200 * time.Millisecond
			}

			select {
			case <-ctx.Done():
			case <-time.After(duration):
			}
		}
	}()

	// Drive the callback queue, notify user's service.
	wg.Add(1)
	go func() {
		defer wg.Done()

		task := v.task
		for ctx.Err() == nil {
			var duration time.Duration
			if err := task.DriveCallbackQueue(ctx); err != nil {
				logger.Wf(ctx, "ocr: task %v drive callback queue err %+v", task.String(), err)
				duration = 10 * time.Second
			} else {
				duration = 200 * time.Millisecond
			}

			select {
			case <-ctx.Done():
			case <-time.After(duration):
			}
		}
	}()

	// Drive the cleanup queue, remove old files.
	wg.Add(1)
	go func() {
		defer wg.Done()

		task := v.task
		for ctx.Err() == nil {
			var duration time.Duration
			if err := task.DriveCleanupQueue(ctx); err != nil {
				logger.Wf(ctx, "ocr: task %v drive cleanup queue err %+v", task.String(), err)
				duration = 10 * time.Second
			} else {
				duration = 200 * time.Millisecond
			}

			select {
			case <-ctx.Done():
			case <-time.After(duration):
			}
		}
	}()

	return nil
}

type OCRConfig struct {
	// Whether ocr all streams.
	All bool `json:"all"`
	// The AI service provider.
	SrsAssistantProvider
	// The AI chat configuration.
	SrsAssistantChat
}

func NewOCRConfig() *OCRConfig {
	v := &OCRConfig{}
	v.All = false
	v.AIChatEnabled = true
	return v
}

func (v OCRConfig) String() string {
	return fmt.Sprintf("all=%v, provider=<%v>, chat=<%v>",
		v.All, v.SrsAssistantProvider.String(), v.SrsAssistantChat.String(),
	)
}

func (v *OCRConfig) Load(ctx context.Context) error {
	if b, err := rdb.HGet(ctx, SRS_OCR_CONFIG, "global").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v global", SRS_OCR_CONFIG)
	} else if len(b) > 0 {
		if err := json.Unmarshal([]byte(b), v); err != nil {
			return errors.Wrapf(err, "unmarshal %v", b)
		}
	}
	return nil
}

func (v *OCRConfig) Save(ctx context.Context) error {
	if b, err := json.Marshal(v); err != nil {
		return errors.Wrapf(err, "marshal conf %v", v)
	} else if err := rdb.HSet(ctx, SRS_OCR_CONFIG, "global", string(b)).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hset %v global %v", SRS_OCR_CONFIG, string(b))
	}
	return nil
}

type OCRSegment struct {
	// The SRS callback message msg.
	Msg *SrsOnHlsMessage `json:"msg,omitempty"`
	// The original source TS file.
	TsFile *TsFile `json:"tsfile,omitempty"`
	// The extracted image file.
	ImageFile *TsFile `json:"image,omitempty"`
	// The ocr result, by AI service.
	OCRText string `json:"ocr,omitempty"`
	// The callback video file.
	CallbackFile *TsFile `json:"callback,omitempty"`

	// The cost to transcode the TS file to image file.
	CostExtractImage time.Duration `json:"eic,omitempty"`
	// The cost to do OCR, converting speech to text.
	CostOCR time.Duration `json:"ocrc,omitempty"`
	// The cost to callback the OCR result.
	CostCallback time.Duration `json:"olc,omitempty"`
}

func (v OCRSegment) String() string {
	var sb strings.Builder
	if v.Msg != nil {
		sb.WriteString(fmt.Sprintf("msg=%v, ", v.Msg.String()))
	}
	if v.TsFile != nil {
		sb.WriteString(fmt.Sprintf("ts=%v, ", v.TsFile.String()))
	}
	if v.ImageFile != nil {
		sb.WriteString(fmt.Sprintf("image=%v, ", v.ImageFile.String()))
		sb.WriteString(fmt.Sprintf("eac=%v, ", v.CostExtractImage))
	}
	if v.OCRText != "" {
		sb.WriteString(fmt.Sprintf("ocr=%v, ", v.OCRText))
		sb.WriteString(fmt.Sprintf("ocrc=%v, ", v.CostOCR))
	}
	if v.CallbackFile != nil {
		sb.WriteString(fmt.Sprintf("callback=%v, ", v.CallbackFile.String()))
		sb.WriteString(fmt.Sprintf("olc=%v, ", v.CostCallback))
	}
	return sb.String()
}

func (v *OCRSegment) Dispose() error {
	// Remove the original ts file.
	if v.TsFile != nil {
		if _, err := os.Stat(v.TsFile.File); err == nil {
			os.Remove(v.TsFile.File)
		}
	}

	// Remove the extracted image file.
	if v.ImageFile != nil {
		if _, err := os.Stat(v.ImageFile.File); err == nil {
			os.Remove(v.ImageFile.File)
		}
	}

	// Remove the callback video file.
	if v.CallbackFile != nil {
		if _, err := os.Stat(v.CallbackFile.File); err == nil {
			os.Remove(v.CallbackFile.File)
		}
	}

	return nil
}

type OCRQueue struct {
	// The ocr segments in the queue.
	Segments []*OCRSegment `json:"segments,omitempty"`

	// To protect the queue.
	lock sync.Mutex
}

func NewOCRQueue() *OCRQueue {
	return &OCRQueue{}
}

func (v *OCRQueue) String() string {
	return fmt.Sprintf("segments=%v", len(v.Segments))
}

func (v *OCRQueue) count() int {
	v.lock.Lock()
	defer v.lock.Unlock()

	return len(v.Segments)
}

func (v *OCRQueue) enqueue(segment *OCRSegment) {
	v.lock.Lock()
	defer v.lock.Unlock()

	v.Segments = append(v.Segments, segment)
}

func (v *OCRQueue) first() *OCRSegment {
	v.lock.Lock()
	defer v.lock.Unlock()

	if len(v.Segments) == 0 {
		return nil
	}

	return v.Segments[0]
}

func (v *OCRQueue) dequeue(segment *OCRSegment) {
	v.lock.Lock()
	defer v.lock.Unlock()

	for i, s := range v.Segments {
		if s == segment {
			v.Segments = append(v.Segments[:i], v.Segments[i+1:]...)
			return
		}
	}
}

func (v *OCRQueue) reset(ctx context.Context) error {
	var segments []*OCRSegment

	func() {
		v.lock.Lock()
		defer v.lock.Unlock()

		segments = v.Segments
		v.Segments = nil
	}()

	for _, segment := range segments {
		segment.Dispose()
	}

	return nil
}

type OCRTask struct {
	// The ID for task.
	UUID string `json:"uuid,omitempty"`

	// The input url.
	Input string `json:"input,omitempty"`
	// The input stream object, select the active stream.
	inputStream *SrsStream

	// The chat history, to use as prompt for next chat.
	histories []openai.ChatCompletionMessage

	// The live queue for the current task. HLS TS segments are copied to the ocr
	// directory, then a segment is created and added to the live queue for the ocr
	// task to process, to convert to image file.
	LiveQueue *OCRQueue `json:"live,omitempty"`
	// The OCR (Optical Character Recognition) queue for the current task. When an image
	// file is generated, the segment is added to the OCR queue, which then requests the
	// AI server to convert the image file into text.
	OCRQueue *OCRQueue `json:"ocr,omitempty"`
	// The callback queue for the current task. It involves drawing OCR (Optical Character
	// Recognition) text onto the video and encoding it into a new video file.
	CallbackQueue *OCRQueue `json:"callback,omitempty"`
	// The cleanup queue for the current task.
	CleanupQueue *OCRQueue `json:"cleanup,omitempty"`

	// The signal to persistence task.
	signalPersistence chan bool
	// The signal to change the active stream for task.
	signalNewStream chan *SrsStream

	// The configure for ocr task.
	config OCRConfig
	// The ocr worker.
	ocrWorker *OCRWorker

	// The context for current task.
	cancel context.CancelFunc

	// To protect the common fields.
	lock sync.Mutex
}

func NewOCRTask() *OCRTask {
	return &OCRTask{
		// Generate a UUID for task.
		UUID: uuid.NewString(),
		// The live queue for current task.
		LiveQueue: NewOCRQueue(),
		// The OCR queue for current task.
		OCRQueue: NewOCRQueue(),
		// The callback queue for current task.
		CallbackQueue: NewOCRQueue(),
		// The cleanup queue for current task.
		CleanupQueue: NewOCRQueue(),
		// Create persistence signal.
		signalPersistence: make(chan bool, 1),
		// Create new stream signal.
		signalNewStream: make(chan *SrsStream, 1),
	}
}

func (v *OCRTask) String() string {
	return fmt.Sprintf("uuid=%v, live=%v, ocr=%v, callback=%v, cleanup=%v, config is %v",
		v.UUID, v.LiveQueue.String(), v.OCRQueue.String(), v.CallbackQueue.String(), v.CleanupQueue.String(),
		v.config.String(),
	)
}

func (v *OCRTask) Run(ctx context.Context) error {
	ctx = logger.WithContext(ctx)
	logger.Tf(ctx, "ocr run task %v", v.String())

	pfn := func(ctx context.Context) error {
		// Load config from redis.
		if err := v.config.Load(ctx); err != nil {
			return errors.Wrapf(err, "load config")
		}

		// Ignore if not enabled.
		if !v.config.All {
			return nil
		}

		// Start ocr task.
		if err := v.doOCR(ctx); err != nil {
			return errors.Wrapf(err, "do ocr")
		}

		return nil
	}

	for ctx.Err() == nil {
		if err := pfn(ctx); err != nil {
			logger.Wf(ctx, "ignore %v err %+v", v.String(), err)

			select {
			case <-ctx.Done():
			case <-time.After(3500 * time.Millisecond):
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

func (v *OCRTask) doOCR(ctx context.Context) error {
	// Create context for current task.
	parentCtx := ctx
	ctx, v.cancel = context.WithCancel(ctx)

	// Main loop, process signals from user or system.
	for ctx.Err() == nil {
		select {
		case <-parentCtx.Done():
		case <-ctx.Done():
		case <-v.signalPersistence:
			if err := v.saveTask(ctx); err != nil {
				return errors.Wrapf(err, "save task %v", v.String())
			}
		case input := <-v.signalNewStream:
			host := "localhost"
			inputURL := fmt.Sprintf("rtmp://%v/%v/%v", host, input.App, input.Stream)
			v.Input, v.inputStream = inputURL, input
		}
	}

	return nil
}

func (v *OCRTask) OnTsSegment(ctx context.Context, msg *SrsOnHlsObject) error {
	// TODO: FIXME: Cleanup the temporary files when task disabled.
	if !v.config.All {
		return nil
	}

	func() {
		// We must not update the queue, when persistence goroutine is working.
		v.lock.Lock()
		v.lock.Unlock()

		v.LiveQueue.enqueue(&OCRSegment{
			Msg:    msg.Msg,
			TsFile: msg.TsFile,
		})
	}()

	// Notify the main loop to persistent current task.
	v.notifyPersistence(ctx)
	return nil
}

// TODO: FIXME: Should reset the stream when republish.
func (v *OCRTask) WatchNewStream(ctx context.Context) error {
	// If not enabled, wait.
	if !v.config.All {
		return nil
	}

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

		if v.inputStream != nil && v.inputStream.StreamURL() != best.StreamURL() {
			logger.Tf(ctx, "ocr use best=%v as input", best.StreamURL())
		}
		return best, nil
	}

	// Use an active stream as input.
	input, err := selectActiveStream()
	if err != nil {
		return errors.Wrapf(err, "select input")
	}

	// Whether stream exists and changed.
	var streamNotChanged bool
	if v.inputStream != nil && input != nil && v.inputStream.StreamURL() == input.StreamURL() {
		streamNotChanged = true
	}

	// No stream, or stream not changed, wait.
	if input == nil || streamNotChanged {
		select {
		case <-ctx.Done():
		case <-time.After(1 * time.Second):
		}
		return nil
	}
	logger.Tf(ctx, "ocr: Got new stream %v", input.String())

	// Notify the main loop to change the active stream.
	select {
	case <-ctx.Done():
	case v.signalNewStream <- input:
		logger.Tf(ctx, "ocr: Notify new stream %v", input.String())
	}

	return nil
}

func (v *OCRTask) DriveLiveQueue(ctx context.Context) error {
	// Ignore if not enabled.
	if !v.config.All {
		return nil
	}

	// Ignore if not enough segments.
	if v.LiveQueue.count() <= 0 {
		return nil
	}

	segment := v.LiveQueue.first()
	starttime := time.Now()

	// Remove segment if file not exists.
	if _, err := os.Stat(segment.TsFile.File); err != nil && os.IsNotExist(err) {
		func() {
			v.lock.Lock()
			defer v.lock.Unlock()
			v.LiveQueue.dequeue(segment)
		}()

		segment.Dispose()
		logger.Tf(ctx, "ocr: remove not exist ts segment %v", segment.String())
		return nil
	}

	// Wait if OCR queue is full.
	if v.OCRQueue.count() >= maxCallbackSegments+1 {
		return nil
	}

	// Transcode to image file, such as jpg.
	imageFile := &TsFile{
		TsID:     fmt.Sprintf("%v-image-%v", segment.TsFile.SeqNo, uuid.NewString()),
		URL:      segment.TsFile.URL,
		SeqNo:    segment.TsFile.SeqNo,
		Duration: segment.TsFile.Duration,
	}
	imageFile.File = path.Join("ocr", fmt.Sprintf("%v.jpg", imageFile.TsID))

	// TODO: FIXME: We should generate a set of images and use the best one.
	args := []string{
		"-i", segment.TsFile.File,
		"-frames:v", "1", "-q:v", "10",
		"-y", imageFile.File,
	}
	if err := exec.CommandContext(ctx, "ffmpeg", args...).Run(); err != nil {
		return errors.Wrapf(err, "transcode %v", args)
	}

	// Update the size of image file.
	stats, err := os.Stat(imageFile.File)
	if err != nil {
		// TODO: FIXME: Cleanup the failed file.
		return errors.Wrapf(err, "stat file %v", imageFile.File)
	}
	imageFile.Size = uint64(stats.Size())

	// Dequeue the segment from live queue and attach to OCR queue.
	func() {
		v.lock.Lock()
		defer v.lock.Unlock()

		v.LiveQueue.dequeue(segment)
		segment.ImageFile = imageFile
		segment.CostExtractImage = time.Since(starttime)
		v.OCRQueue.enqueue(segment)
	}()
	logger.Tf(ctx, "ocr: extract image %v to %v, size=%v, cost=%v",
		segment.TsFile.File, imageFile.File, imageFile.Size, segment.CostExtractImage)

	// Notify the main loop to persistent current task.
	v.notifyPersistence(ctx)
	return nil
}

func (v *OCRTask) DriveOCRQueue(ctx context.Context) error {
	// Ignore if not enabled.
	if !v.config.All {
		return nil
	}

	// Ignore if not enough segments.
	if v.OCRQueue.count() <= 0 {
		return nil
	}

	segment := v.OCRQueue.first()
	starttime := time.Now()

	// Remove segment if file not exists.
	if _, err := os.Stat(segment.ImageFile.File); err != nil && os.IsNotExist(err) {
		func() {
			v.lock.Lock()
			defer v.lock.Unlock()
			v.OCRQueue.dequeue(segment)
		}()
		segment.Dispose()
		logger.Tf(ctx, "ocr: remove not exist image segment %v", segment.String())
		return nil
	}

	// Wait if callback queue is full.
	if v.CallbackQueue.count() >= maxCallbackSegments+1 {
		return nil
	}

	// Read the image file and convert to base64.
	var imageData string
	if data, err := os.ReadFile(segment.ImageFile.File); err != nil {
		return errors.Wrapf(err, "read image from %v", segment.ImageFile.File)
	} else {
		imageData = base64.StdEncoding.EncodeToString(data)
	}

	// Convert the image file to text by AI.
	var config openai.ClientConfig
	config = openai.DefaultConfig(v.config.AISecretKey)
	config.BaseURL = v.config.AIBaseURL
	config.OrgID = v.config.AIOrganization

	prompt := v.config.AIChatPrompt
	system := fmt.Sprintf("Keep your reply neat, limiting the reply to %v words.", v.config.AIChatMaxWords)

	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: system},
	}

	messages = append(messages, v.histories...)
	messages = append(messages, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleUser, Content: prompt,
	})
	messages = append(messages, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleUser, MultiContent: []openai.ChatMessagePart{
			{Type: openai.ChatMessagePartTypeImageURL, ImageURL: &openai.ChatMessageImageURL{
				Detail: openai.ImageURLDetailLow, URL: fmt.Sprintf("data:image/jpeg;base64,%v", imageData),
			}},
		},
	})

	client := openai.NewClientWithConfig(config)
	resp, err := client.CreateChatCompletion(
		ctx, openai.ChatCompletionRequest{
			Model: v.config.AIChatModel, Messages: messages,
		},
	)
	if err != nil {
		return errors.Wrapf(err, "AI process, model=%v, image=%v, messages=%v, system=<%v>, prompt=<%v>",
			v.config.AIChatModel, segment.ImageFile.File, len(messages), system, prompt,
		)
	}

	segment.OCRText = resp.Choices[0].Message.Content
	segment.CostOCR = time.Since(starttime)

	// Build the historical messages.
	if segment.OCRText != "" {
		v.histories = append(v.histories, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		}, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: segment.OCRText,
		})

		for len(v.histories) > v.config.AIChatMaxWindow*2 {
			v.histories = v.histories[1:]
		}
	}

	// Dequeue the segment from OCR queue and attach to correct queue.
	func() {
		v.lock.Lock()
		defer v.lock.Unlock()
		v.OCRQueue.dequeue(segment)
	}()
	func() {
		v.lock.Lock()
		defer v.lock.Unlock()
		v.CallbackQueue.enqueue(segment)
	}()
	logger.Tf(ctx, "ocr: recognize image=%v, model=%v, prompt=%v, text=%v, cost=%v",
		segment.ImageFile.File, v.config.AIChatModel, prompt, segment.OCRText, segment.CostOCR)

	// Notify the main loop to persistent current task.
	v.notifyPersistence(ctx)
	return nil
}

func (v *OCRTask) DriveCallbackQueue(ctx context.Context) error {
	// Ignore if not enabled.
	if !v.config.All {
		return nil
	}

	// Ignore if not enough segments.
	if v.CallbackQueue.count() <= 0 {
		return nil
	}

	// Wait if cleanup queue is full.
	if v.CleanupQueue.count() >= maxCallbackSegments+1 {
		return nil
	}

	segment := v.CallbackQueue.first()
	starttime := time.Now()

	// Do callback to notify user's service.
	if err := callbackWorker.OnOCR(ctx, SrsActionOnOcr, v.UUID, segment.Msg, v.config.AIChatPrompt, segment.OCRText); err != nil {
		logger.Wf(ctx, "ocr: ignore callback %v err %+v", segment.String(), err)
	}

	segment.CostCallback = time.Since(starttime)
	logger.Tf(ctx, "ocr: callback %v, cost=%v", segment.String(), segment.CostCallback)

	// Dequeue the segment from callback queue and attach to cleanup queue.
	func() {
		v.lock.Lock()
		defer v.lock.Unlock()
		v.CallbackQueue.dequeue(segment)
	}()
	func() {
		v.lock.Lock()
		defer v.lock.Unlock()
		v.CleanupQueue.enqueue(segment)
	}()
	logger.Tf(ctx, "ocr: callback image=%v, ocr=%v, cost=%v",
		segment.ImageFile.File, segment.OCRText, segment.CostCallback)

	// Notify the main loop to persistent current task.
	v.notifyPersistence(ctx)
	return nil
}

func (v *OCRTask) DriveCleanupQueue(ctx context.Context) error {
	// Ignore if not enabled.
	if !v.config.All {
		return nil
	}

	// Ignore if not enough segments.
	if v.CallbackQueue.count() <= maxCallbackSegments {
		select {
		case <-ctx.Done():
		case <-time.After(1 * time.Second):
		}
		return nil
	}

	// Cleanup the old segments.
	segment := v.CallbackQueue.first()
	func() {
		v.lock.Lock()
		defer v.lock.Unlock()
		v.CallbackQueue.dequeue(segment)
	}()
	defer segment.Dispose()

	// Notify the main loop to persistent current task.
	v.notifyPersistence(ctx)
	return nil
}

// TODO: FIXME: Should restart task when stream unpublish.
func (v *OCRTask) restart(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if v.cancel != nil {
		v.cancel()
	}

	return nil
}

func (v *OCRTask) reset(ctx context.Context) error {
	if v.config.All {
		// Retry to wait for the task to apply the changed config.
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
			if v.config.All {
				return fmt.Errorf("can not reset when running")
			}
		}
	}

	if err := func() error {
		v.lock.Lock()
		defer v.lock.Unlock()

		// Reset all queues.
		v.LiveQueue.reset(ctx)
		v.OCRQueue.reset(ctx)
		v.CallbackQueue.reset(ctx)
		v.CleanupQueue.reset(ctx)

		// Reset all states.
		v.Input = ""

		// Remove previous task from redis.
		if err := rdb.HDel(ctx, SRS_OCR_TASK, v.UUID).Err(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hdel %v %v", SRS_OCR_TASK, v.UUID)
		}

		// Regenerate new UUID.
		v.UUID = uuid.NewString()

		return nil
	}(); err != nil {
		return errors.Wrapf(err, "reset task")
	}

	// Notify the main loop to persistent current task.
	v.notifyPersistence(ctx)

	// Wait for task to persistence and avoid to reset very fast.
	select {
	case <-ctx.Done():
	case <-time.After(1 * time.Second):
	}
	return nil
}

func (v *OCRTask) enabled() bool {
	v.lock.Lock()
	defer v.lock.Unlock()

	return v.config.All
}

func (v *OCRTask) match(msg *SrsOnHlsMessage) bool {
	v.lock.Lock()
	defer v.lock.Unlock()

	if v.inputStream == nil || !v.config.All {
		return false
	}

	if v.inputStream.App != msg.App || v.inputStream.Stream != msg.Stream {
		return false
	}

	return true
}

func (v *OCRTask) liveSegments() []*OCRSegment {
	v.lock.Lock()
	defer v.lock.Unlock()

	return v.LiveQueue.Segments[:]
}

func (v *OCRTask) ocrSegments() []*OCRSegment {
	v.lock.Lock()
	defer v.lock.Unlock()

	return v.OCRQueue.Segments[:]
}

func (v *OCRTask) callbackSegments() []*OCRSegment {
	v.lock.Lock()
	defer v.lock.Unlock()

	return v.CallbackQueue.Segments[:]
}

func (v *OCRTask) cleanupSegments() []*OCRSegment {
	v.lock.Lock()
	defer v.lock.Unlock()

	return v.CleanupQueue.Segments[:]
}

func (v *OCRTask) notifyPersistence(ctx context.Context) {
	select {
	case <-ctx.Done():
	case v.signalPersistence <- true:
	default:
	}
}

func (v *OCRTask) saveTask(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	starttime := time.Now()

	if b, err := json.Marshal(v); err != nil {
		return errors.Wrapf(err, "marshal %v", v.String())
	} else if err = rdb.HSet(ctx, SRS_OCR_TASK, v.UUID, string(b)).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_OCR_TASK, v.UUID, string(b))
	}

	logger.Tf(ctx, "ocr persistence ok, cost=%v", time.Since(starttime))

	return nil
}
