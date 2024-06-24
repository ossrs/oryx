// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	ohttp "github.com/ossrs/go-oryx-lib/http"
	"github.com/ossrs/go-oryx-lib/logger"

	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var dvrWorker *DvrWorker

type DvrWorker struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Tencent Cloud credentials.
	secretId   string
	secretKey  string
	bucketName string
	cosClient  *cos.Client

	// Got message from SRS, a new TS segment file is generated.
	msgs chan *SrsOnHlsObject
	// The streams we're dvring, key is m3u8 URL in string, value is m3u8 object *DvrM3u8Stream.
	streams sync.Map
}

func NewDvrWorker() *DvrWorker {
	return &DvrWorker{
		msgs: make(chan *SrsOnHlsObject, 1024),
	}
}

func (v *DvrWorker) ready() bool {
	return v.secretId != "" && v.secretKey != "" && v.bucketName != "" && v.cosClient != nil
}

func (v *DvrWorker) Handle(ctx context.Context, handler *http.ServeMux) error {
	ep := "/terraform/v1/hooks/dvr/query"
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

			all, err := rdb.HGet(ctx, SRS_DVR_PATTERNS, "all").Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v all", SRS_DVR_PATTERNS)
			}

			appId, _ := rdb.HGet(ctx, SRS_TENCENT_CAM, "appId").Result()
			secretId, _ := rdb.HGet(ctx, SRS_TENCENT_CAM, "secretId").Result()
			secretKey, _ := rdb.HGet(ctx, SRS_TENCENT_CAM, "secretKey").Result()

			ohttp.WriteData(ctx, w, r, &struct {
				All    bool `json:"all"`
				Secret bool `json:"secret"`
			}{
				All:    all == "true",
				Secret: appId != "" && secretId != "" && secretKey != "",
			})

			logger.Tf(ctx, "dvr query ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/hooks/dvr/apply"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var all bool
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				All   *bool   `json:"all"`
			}{
				Token: &token, All: &all,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if all, err := rdb.HSet(ctx, SRS_DVR_PATTERNS, "all", fmt.Sprintf("%v", all)).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v all %v", SRS_DVR_PATTERNS, all)
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "dvr query ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/hooks/dvr/files"
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

			keys, cursor, err := rdb.HScan(ctx, SRS_DVR_M3U8_ARTIFACT, 0, "*", 100).Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hscan %v 0 * 100", SRS_DVR_M3U8_ARTIFACT)
			}

			files := []map[string]interface{}{}
			for i := 0; i < len(keys); i += 2 {
				var metadata M3u8VoDArtifact
				if err := json.Unmarshal([]byte(keys[i+1]), &metadata); err != nil {
					return errors.Wrapf(err, "json parse %v", keys[i+1])
				}

				var duration float64
				var size uint64
				for _, file := range metadata.Files {
					duration += file.Duration
					size += file.Size
				}

				files = append(files, map[string]interface{}{
					"uuid":     metadata.UUID,
					"vhost":    metadata.Vhost,
					"app":      metadata.App,
					"stream":   metadata.Stream,
					"progress": metadata.Processing,
					"update":   metadata.Update,
					"nn":       len(metadata.Files),
					"duration": duration,
					"size":     size,
					// For DVR only.
					"bucket": metadata.Bucket,
					"region": metadata.Region,
				})
			}

			ohttp.WriteData(ctx, w, r, files)
			logger.Tf(ctx, "dvr files ok, cursor=%v, token=%vB", cursor, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/hooks/dvr/hls/"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			// Format is :uuid.m3u8 or :uuid/index.m3u8
			filename := r.URL.Path[len("/terraform/v1/hooks/dvr/hls/"):]
			// Format is :uuid.m3u8
			filename = strings.ReplaceAll(filename, "/index.m3u8", ".m3u8")
			uuid := filename[:len(filename)-len(path.Ext(filename))]
			if len(uuid) == 0 {
				return errors.Errorf("invalid uuid %v from %v of %v", uuid, filename, r.URL.Path)
			}

			var metadata M3u8VoDArtifact
			if m3u8Metadata, err := rdb.HGet(ctx, SRS_DVR_M3U8_ARTIFACT, uuid).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_DVR_M3U8_ARTIFACT, uuid)
			} else if m3u8Metadata == "" {
				return errors.Errorf("no m3u8 of uuid=%v", uuid)
			} else if err = json.Unmarshal([]byte(m3u8Metadata), &metadata); err != nil {
				return errors.Wrapf(err, "parse %v", m3u8Metadata)
			}

			contentType, m3u8Body, duration, err := buildVodM3u8(
				ctx, &metadata, true, "", false, "",
			)
			if err != nil {
				return errors.Wrapf(err, "build vod m3u8 of %v", metadata.String())
			}

			w.Header().Set("Content-Type", contentType)
			w.Write([]byte(m3u8Body))
			logger.Tf(ctx, "dvr generate m3u8 ok, uuid=%v, duration=%v", uuid, duration)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	return nil
}

func (v *DvrWorker) OnHlsTsMessage(ctx context.Context, msg *SrsOnHlsMessage) error {
	// Ignore for Tencent Cloud credentials not ready.
	if !v.ready() {
		return nil
	}

	// Copy the ts file to temporary cache dir.
	tsid := uuid.NewString()
	tsfile := path.Join("dvr", fmt.Sprintf("%v.ts", tsid))

	// Always use execFile when params contains user inputs, see https://auth0.com/blog/preventing-command-injection-attacks-in-node-js-apps/
	// Note that should never use fs.copyFileSync(file, tsfile, fs.constants.COPYFILE_FICLONE_FORCE) which fails in macOS.
	if err := exec.CommandContext(ctx, "cp", "-f", msg.File, tsfile).Run(); err != nil {
		return errors.Wrapf(err, "copy file %v to %v", msg.File, tsfile)
	}

	// Get the file size.
	stats, err := os.Stat(msg.File)
	if err != nil {
		return errors.Wrapf(err, "stat file %v", msg.File)
	}

	// Create a local ts file object.
	tsFile := &TsFile{
		TsID:     tsid,
		URL:      msg.URL,
		SeqNo:    msg.SeqNo,
		Duration: msg.Duration,
		Size:     uint64(stats.Size()),
		File:     tsfile,
	}

	// Notify worker asynchronously.
	go func() {
		select {
		case <-ctx.Done():
		case v.msgs <- &SrsOnHlsObject{Msg: msg, TsFile: tsFile}:
		}
	}()
	return nil
}

func (v *DvrWorker) Close() error {
	if v.cancel != nil {
		v.cancel()
	}
	v.wg.Wait()
	return nil
}

func (v *DvrWorker) Start(ctx context.Context) error {
	wg := &v.wg

	ctx, cancel := context.WithCancel(ctx)
	v.cancel = cancel

	ctx = logger.WithContext(ctx)
	logger.Tf(ctx, "Dvr: start a worker")

	// When system startup, we initialize the credentials, for worker to use it.
	if err := v.updateCredential(ctx); err != nil {
		return errors.Wrapf(err, "update credential")
	}

	// Load all objects from redis.
	if objs, err := rdb.HGetAll(ctx, SRS_DVR_M3U8_WORKING).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hgetall %v", SRS_DVR_M3U8_WORKING)
	} else if len(objs) > 0 {
		for m3u8URL, value := range objs {
			logger.Tf(ctx, "Load %v object %v", m3u8URL, value)

			var m3u8LocalObj DvrM3u8Stream
			if err = json.Unmarshal([]byte(value), &m3u8LocalObj); err != nil {
				return errors.Wrapf(err, "load %v", value)
			}

			// Initialize object.
			if err := m3u8LocalObj.Initialize(ctx, v); err != nil {
				return errors.Wrapf(err, "init %v", m3u8LocalObj.String())
			}

			// Save in memory object.
			v.streams.Store(m3u8URL, &m3u8LocalObj)

			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := m3u8LocalObj.Run(ctx); err != nil {
					logger.Wf(ctx, "serve m3u8 %v err %+v", m3u8LocalObj.String(), err)
				}
			}()
		}
	}

	// Create M3u8 object from message.
	buildM3u8Object := func(ctx context.Context, msg *SrsOnHlsObject) error {
		logger.Tf(ctx, "Dvr: Got message %v", msg.String())

		// Load stream local object.
		var m3u8LocalObj *DvrM3u8Stream
		var freshObject bool
		if obj, loaded := v.streams.LoadOrStore(msg.Msg.M3u8URL, &DvrM3u8Stream{
			M3u8URL: msg.Msg.M3u8URL, UUID: uuid.NewString(), dvrWorker: v,
		}); true {
			m3u8LocalObj, freshObject = obj.(*DvrM3u8Stream), !loaded
		}

		// Serve object if fresh one.
		if freshObject {
			// Initialize object.
			if err := m3u8LocalObj.Initialize(ctx, v); err != nil {
				return errors.Wrapf(err, "init %v", m3u8LocalObj.String())
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := m3u8LocalObj.Run(ctx); err != nil {
					logger.Wf(ctx, "serve m3u8 %v err %+v", m3u8LocalObj.String(), err)
				}
			}()
		}

		// Append new ts file to object.
		m3u8LocalObj.addMessage(ctx, msg)

		// Always save the object to redis, for reloading it when restart.
		if err := m3u8LocalObj.saveObject(ctx); err != nil {
			return errors.Wrapf(err, "save %v", m3u8LocalObj.String())
		}

		return nil
	}

	// Process all messages about HLS ts segments.
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-v.msgs:
				if err := buildM3u8Object(ctx, msg); err != nil {
					logger.Wf(ctx, "ignore msg %v ts %v err %+v", msg.Msg.String(), msg.TsFile.String(), err)
				}
			}
		}
	}()

	// Load tencent cloud credentials.
	wg.Add(1)
	go func() {
		defer wg.Done()

		for ctx.Err() == nil {
			var duration time.Duration

			if err := v.updateCredential(ctx); err != nil {
				logger.Wf(ctx, "ignore err %+v", err)
				duration = 30 * time.Second
			} else if !v.ready() {
				duration = 1 * time.Second
			} else {
				duration = 3 * time.Second
			}

			select {
			case <-ctx.Done():
			case <-time.After(duration):
			}
		}
	}()

	return nil
}

func (v *DvrWorker) updateCredential(ctx context.Context) error {
	previous := &struct {
		SecretId, SecretKey, BucketName string
	}{
		SecretId: v.secretId, SecretKey: v.secretKey, BucketName: v.bucketName,
	}

	// The credential might not be ready, so we ignore error.
	if secretId, err := rdb.HGet(ctx, SRS_TENCENT_CAM, "secretId").Result(); err == nil {
		v.secretId = secretId
	}

	if secretKey, err := rdb.HGet(ctx, SRS_TENCENT_CAM, "secretKey").Result(); err == nil {
		v.secretKey = secretKey
	}

	if bucketName, err := rdb.HGet(ctx, SRS_TENCENT_COS, "bucket").Result(); err == nil {
		v.bucketName = bucketName
	}

	changed := v.bucketName != previous.BucketName || v.secretId != previous.SecretId ||
		v.secretKey != previous.SecretKey
	credentialOK := v.secretId != "" && v.secretKey != "" && v.bucketName != ""
	if (v.cosClient == nil || changed) && credentialOK {
		location := fmt.Sprintf("%v.cos.%v.myqcloud.com", v.bucketName, conf.Region)
		u, err := url.Parse(fmt.Sprintf("https://%v", location))
		if err != nil {
			return errors.Wrapf(err, "parse %v", fmt.Sprintf("https://%v", location))
		}

		// Create the bucket with private ACL.
		// See https://cloud.tencent.com/document/product/436/65639
		v.cosClient = cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
			Transport: &cos.AuthorizationTransport{SecretID: v.secretId, SecretKey: v.secretKey},
		})
		logger.Tf(ctx, "create dvr client ok, bucket=%v, location=%v", v.bucketName, location)
	}
	return nil
}

// DvrM3u8Stream is the current active local object for a HLS stream.
// When dvring done, it will generate a M3u8VoDArtifact, which is a HLS VoD object.
type DvrM3u8Stream struct {
	// The url of m3u8, generated by SRS, such as live/livestream/live.m3u8
	M3u8URL string `json:"m3u8_url"`
	// The uuid of M3u8VoDObject, generated by worker, such as 3ECF0239-708C-42E4-96E1-5AE935C6E6A9
	UUID string `json:"uuid"`

	// Number of local files.
	NN int `json:"nn"`
	// The last update time.
	Update string `json:"update"`
	// The done time.
	Done string `json:"done"`

	// The ts files of this m3u8.
	Messages []*SrsOnHlsObject `json:"msgs"`

	// The worker which owns this object.
	dvrWorker *DvrWorker
	// The artifact we're working for.
	artifact *M3u8VoDArtifact
	// To protect the fields.
	lock sync.Mutex
}

func (v DvrM3u8Stream) String() string {
	return fmt.Sprintf("url=%v, uuid=%v, done=%v, update=%v, messages=%v",
		v.M3u8URL, v.UUID, v.Done, v.Update, len(v.Messages),
	)
}

func (v *DvrM3u8Stream) deleteObject(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if err := rdb.HDel(ctx, SRS_DVR_M3U8_WORKING, v.M3u8URL).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hdel %v %v", SRS_DVR_M3U8_WORKING, v.M3u8URL)
	}

	return nil
}

func (v *DvrM3u8Stream) saveObject(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if b, err := json.Marshal(v); err != nil {
		return errors.Wrapf(err, "marshal object")
	} else if err = rdb.HSet(ctx, SRS_DVR_M3U8_WORKING, v.M3u8URL, string(b)).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_DVR_M3U8_WORKING, v.M3u8URL, string(b))
	}
	return nil
}

func (v *DvrM3u8Stream) saveArtifact(ctx context.Context, artifact *M3u8VoDArtifact) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if b, err := json.Marshal(artifact); err != nil {
		return errors.Wrapf(err, "marshal %v", artifact.String())
	} else if err = rdb.HSet(ctx, SRS_DVR_M3U8_ARTIFACT, v.UUID, string(b)).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_DVR_M3U8_ARTIFACT, v.UUID, string(b))
	}
	return nil
}

func (v *DvrM3u8Stream) updateArtifact(ctx context.Context, artifact *M3u8VoDArtifact, msg *SrsOnHlsObject) {
	v.lock.Lock()
	defer v.lock.Unlock()

	artifact.Vhost = msg.Msg.Vhost
	artifact.App = msg.Msg.App
	artifact.Stream = msg.Msg.Stream

	artifact.Files = append(artifact.Files, msg.TsFile)
	artifact.NN = len(artifact.Files)

	artifact.Update = time.Now().Format(time.RFC3339)
}

func (v *DvrM3u8Stream) finishArtifact(ctx context.Context, artifact *M3u8VoDArtifact) {
	v.lock.Lock()
	defer v.lock.Unlock()

	artifact.Processing = false
	artifact.Update = time.Now().Format(time.RFC3339)
}

func (v *DvrM3u8Stream) addMessage(ctx context.Context, msg *SrsOnHlsObject) {
	v.lock.Lock()
	defer v.lock.Unlock()

	v.Messages = append(v.Messages, msg)
	v.NN = len(v.Messages)
	v.Update = time.Now().Format(time.RFC3339)
}

func (v *DvrM3u8Stream) copyMessages() []*SrsOnHlsObject {
	v.lock.Lock()
	defer v.lock.Unlock()

	return append([]*SrsOnHlsObject{}, v.Messages...)
}

func (v *DvrM3u8Stream) removeMessage(ctx context.Context, msg *SrsOnHlsObject) {
	v.lock.Lock()
	defer v.lock.Unlock()

	for index, m := range v.Messages {
		if m == msg {
			v.Messages = append(v.Messages[:index], v.Messages[index+1:]...)
			break
		}
	}

	v.NN = len(v.Messages)
	v.Update = time.Now().Format(time.RFC3339)

	// Remove the tsfile.
	if err := os.Remove(msg.TsFile.File); err != nil {
		logger.Wf(ctx, "ignore remove file %v err %+v", msg.TsFile.File, err)
	}
}

func (v *DvrM3u8Stream) expired() bool {
	v.lock.Lock()
	defer v.lock.Unlock()

	update, err := time.Parse(time.RFC3339, v.Update)
	if err != nil {
		return true
	}

	duration := 30 * time.Second
	if envNodeEnv() != "development" {
		duration = 300 * time.Second
	}

	if update.Add(duration).Before(time.Now()) {
		return true
	}

	return false
}

// Initialize to load artifact. There is no simultaneously access, so no certFileLock is needed.
func (v *DvrM3u8Stream) Initialize(ctx context.Context, r *DvrWorker) error {
	v.dvrWorker = r
	logger.Tf(ctx, "dvr initialize url=%v, uuid=%v", v.M3u8URL, v.UUID)

	// Try to load artifact from redis. The final artifact is VoD HLS object.
	if value, err := rdb.HGet(ctx, SRS_DVR_M3U8_ARTIFACT, v.UUID).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_DVR_M3U8_ARTIFACT, v.UUID)
	} else if value != "" {
		artifact := &M3u8VoDArtifact{}
		if err = json.Unmarshal([]byte(value), artifact); err != nil {
			return errors.Wrapf(err, "unmarshal %v", value)
		} else {
			v.artifact = artifact
		}
	}

	// Create a artifact if new.
	if v.artifact == nil {
		v.artifact = &M3u8VoDArtifact{
			UUID:       v.UUID,
			M3u8URL:    v.M3u8URL,
			Processing: true,
			Update:     time.Now().Format(time.RFC3339),
		}

		// For DVR required parameters.
		v.artifact.Bucket = r.bucketName
		v.artifact.Region = conf.Region

		if err := v.saveArtifact(ctx, v.artifact); err != nil {
			return errors.Wrapf(err, "save artifact %v", v.artifact.String())
		}
	}

	return nil
}

// Run to serve the current dvring object.
func (v *DvrM3u8Stream) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(logger.WithContext(ctx))
	logger.Tf(ctx, "dvr run task %v", v.String())

	pfn := func() error {
		// Process message and remove it.
		msgs := v.copyMessages()
		for _, msg := range msgs {
			if err := v.serveMessage(ctx, msg); err != nil {
				logger.Wf(ctx, "ignore %v err %+v", msg.String(), err)
			}
		}

		// Refresh redis if got messages to serve.
		if len(msgs) > 0 {
			if err := v.saveObject(ctx); err != nil {
				return errors.Wrapf(err, "save object %v", v.String())
			}
		}

		// Ignore if still has messages to process.
		if len(v.Messages) > 0 {
			return nil
		}

		// Check whether expired.
		if !v.expired() {
			return nil
		}

		// Try to finish the object.
		if err := v.finishM3u8(ctx); err != nil {
			return errors.Wrapf(err, "finish m3u8")
		}

		// Now HLS is done
		logger.Tf(ctx, "Dvr is done, hls is %v, artifact is %v", v.String(), v.artifact.String())
		cancel()

		return nil
	}

	for ctx.Err() == nil {
		if err := pfn(); err != nil {
			logger.Wf(ctx, "ignore %v err %+v", v.String(), err)

			select {
			case <-ctx.Done():
			case <-time.After(10 * time.Second):
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

func (v *DvrM3u8Stream) serveMessage(ctx context.Context, msg *SrsOnHlsObject) error {
	// We always remove the msg from current object.
	defer v.removeMessage(ctx, msg)

	// Ignore file if credential is not ready.
	if !v.dvrWorker.ready() {
		return nil
	}

	// Ignore file if not exists.
	if _, err := os.Stat(msg.TsFile.File); err != nil {
		return err
	}

	key := fmt.Sprintf("%v/%v.ts", v.UUID, msg.TsFile.TsID)
	msg.TsFile.Key = key

	f, err := os.Open(msg.TsFile.File)
	if err != nil {
		return errors.Wrapf(err, "open file %v", msg.TsFile.File)
	}

	// Upload to COS bucket.
	// See https://cloud.tencent.com/document/product/436/64980
	opt := cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType:   "video/MP2T",
			ContentLength: int64(msg.TsFile.Size),
		},
	}
	if _, err = v.dvrWorker.cosClient.Object.Put(ctx, key, f, &opt); err != nil {
		return errors.Wrapf(err, "cos put object %v", key)
	}

	// Update the metadata for m3u8.
	v.updateArtifact(ctx, v.artifact, msg)
	if err := v.saveArtifact(ctx, v.artifact); err != nil {
		return errors.Wrapf(err, "save artifact %v", v.artifact.String())
	}

	logger.Tf(ctx, "dvr consume msg %v", msg.String())
	return nil
}

func (v *DvrM3u8Stream) finishM3u8(ctx context.Context) error {
	contentType, m3u8Body, duration, err := buildVodM3u8(ctx, v.artifact, false, "", false, "")
	if err != nil {
		return errors.Wrapf(err, "build vod")
	}

	// Upload the ts file to COS.
	key := fmt.Sprintf("%v/index.m3u8", v.UUID)

	// Upload to COS bucket.
	// See https://cloud.tencent.com/document/product/436/64980
	opt := cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType:   contentType,
			ContentLength: int64(len(m3u8Body)),
		},
	}
	if _, err = v.dvrWorker.cosClient.Object.Put(ctx, key, strings.NewReader(m3u8Body), &opt); err != nil {
		return errors.Wrapf(err, "cos put object %v", key)
	}
	logger.Tf(ctx, "dvr to %v, duration=%v ok", key, duration)

	// Remove object from worker.
	v.dvrWorker.streams.Delete(v.M3u8URL)

	// Update artifact after finally.
	v.finishArtifact(ctx, v.artifact)
	r0 := v.saveArtifact(ctx, v.artifact)
	r1 := v.deleteObject(ctx)
	logger.Tf(ctx, "dvr cleanup ok, r0=%v, r1=%v", r0, r1)

	// Do final cleanup, because new messages might arrive while converting to mp4, which takes a long time.
	files := v.copyMessages()
	for _, file := range files {
		r2 := os.Remove(file.TsFile.File)
		logger.Tf(ctx, "drop %v r2=%v", file.String(), r2)
	}

	return nil
}
