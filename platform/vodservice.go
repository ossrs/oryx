//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: MIT
//
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	ohttp "github.com/ossrs/go-oryx-lib/http"
	"github.com/ossrs/go-oryx-lib/logger"

	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var vodWorker *VodWorker

type VodWorker struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Tencent Cloud credentials.
	secretId  string
	secretKey string
	vodAppID  uint64
	vodClient *vod.Client

	// Got message from SRS, a new TS segment file is generated.
	msgs chan *SrsOnHlsObject
	// The streams we're voding, key is m3u8 URL in string, value is m3u8 object *VodM3u8Stream.
	streams sync.Map
}

func NewVodWorker() *VodWorker {
	return &VodWorker{
		msgs: make(chan *SrsOnHlsObject, 1024),
	}
}

func (v *VodWorker) ready() bool {
	return v.secretId != "" && v.secretKey != "" && v.vodAppID != 0 && v.vodClient != nil
}

func (v *VodWorker) Handle(ctx context.Context, handler *http.ServeMux) error {
	ep := "/terraform/v1/hooks/vod/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var token string
			if err := json.Unmarshal(b, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "json unmarshal %v", string(b))
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			// Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
			// See https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
			if _, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(apiSecret), nil
			}); err != nil {
				return errors.Wrapf(err, "verify token %v", token)
			}

			all, err := rdb.HGet(ctx, SRS_VOD_PATTERNS, "all").Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v all", SRS_VOD_PATTERNS)
			}

			appId, _ := rdb.HGet(ctx, SRS_TENCENT_CAM, "appId").Result()
			secretId, _ := rdb.HGet(ctx, SRS_TENCENT_CAM, "secretId").Result()
			secretKey, _ := rdb.HGet(ctx, SRS_TENCENT_CAM, "secretKey").Result()

			// VoD service status.
			service, _ := rdb.HGet(ctx, SRS_TENCENT_VOD, "service").Result()
			storage, _ := rdb.HGet(ctx, SRS_TENCENT_VOD, "storage").Result()

			ohttp.WriteData(ctx, w, r, &struct {
				All     bool `json:"all"`
				Secret  bool `json:"secret"`
				Service bool `json:"service"`
				Storage bool `json:"storage"`
			}{
				All:    all == "true",
				Secret: appId != "" && secretId != "" && secretKey != "",
				// For previous platform, the service is set to string ok, so we also create an application.
				Service: service != "" && service != "ok",
				Storage: storage != "",
			})

			logger.Tf(ctx, "vod query ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/hooks/vod/apply"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var token string
			var all bool
			if err := json.Unmarshal(b, &struct {
				Token *string `json:"token"`
				All   *bool   `json:"all"`
			}{
				Token: &token, All: &all,
			}); err != nil {
				return errors.Wrapf(err, "json unmarshal %v", string(b))
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			// Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
			// See https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
			if _, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(apiSecret), nil
			}); err != nil {
				return errors.Wrapf(err, "verify token %v", token)
			}

			if all, err := rdb.HSet(ctx, SRS_VOD_PATTERNS, "all", fmt.Sprintf("%v", all)).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v all %v", SRS_VOD_PATTERNS, all)
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "vod apply ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/hooks/vod/files"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var token string
			if err := json.Unmarshal(b, &struct {
				Token *string `json:"token"`
			}{
				Token: &token,
			}); err != nil {
				return errors.Wrapf(err, "json unmarshal %v", string(b))
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			// Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
			// See https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
			if _, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(apiSecret), nil
			}); err != nil {
				return errors.Wrapf(err, "verify token %v", token)
			}

			keys, cursor, err := rdb.HScan(ctx, SRS_VOD_M3U8_ARTIFACT, 0, "*", 100).Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hscan %v 0 * 100", SRS_VOD_M3U8_ARTIFACT)
			}

			// Unmarshal artifacts from redis.
			artifacts, pendingArtifacts := []*M3u8VoDArtifact{}, []*M3u8VoDArtifact{}
			for i := 0; i < len(keys); i += 2 {
				var metadata M3u8VoDArtifact
				if err := json.Unmarshal([]byte(keys[i+1]), &metadata); err != nil {
					return errors.Wrapf(err, "json parse %v", keys[i+1])
				}

				if metadata.Definition != 0 && metadata.TaskID != "" && metadata.Task == nil {
					pendingArtifacts = append(pendingArtifacts, &metadata)
				}

				artifacts = append(artifacts, &metadata)
			}

			// TODO: FIXME: Query in worker thread.
			// Query media information for artifacts.
			// See https://cloud.tencent.com/document/product/266/31763
			if len(pendingArtifacts) > 0 {
				var fileIDs []string
				fileIDKVs := make(map[string]*M3u8VoDArtifact)
				for _, artifact := range pendingArtifacts {
					fileIDs = append(fileIDs, artifact.FileID)
					fileIDKVs[artifact.FileID] = artifact
				}

				request := vod.NewDescribeMediaInfosRequest()
				request.FileIds = common.StringPtrs(fileIDs)
				request.SubAppId = common.Uint64Ptr(v.vodAppID)
				request.Filters = common.StringPtrs([]string{"transcodeInfo"})

				var mediaInfoSet []*vod.MediaInfo
				if response, err := v.vodClient.DescribeMediaInfosWithContext(ctx, request); err != nil {
					return errors.Wrapf(err, "describe media info")
				} else if len(response.Response.MediaInfoSet) > 0 {
					mediaInfoSet = response.Response.MediaInfoSet
				}

				for _, mis := range mediaInfoSet {
					if mis.FileId == nil || mis == nil || mis.TranscodeInfo == nil || len(mis.TranscodeInfo.TranscodeSet) == 0 {
						continue
					}

					artifact, ok := fileIDKVs[*mis.FileId]
					if !ok {
						continue
					}

					var matched *vod.MediaTranscodeItem
					for _, ts := range mis.TranscodeInfo.TranscodeSet {
						if ts.Definition == nil || uint64(*ts.Definition) != artifact.Definition {
							continue
						}

						matched = ts
						break
					}
					if matched == nil {
						continue
					}

					artifact.Task = &VodTaskArtifact{
						URL:      *matched.Url,
						Bitrate:  *matched.Bitrate,
						Height:   int32(*matched.Height),
						Width:    int32(*matched.Width),
						Size:     *matched.Size,
						Duration: *matched.Duration,
						MD5:      *matched.Md5,
					}

					if b, err := json.Marshal(artifact); err != nil {
						return errors.Wrapf(err, "marshal %v", artifact.String())
					} else if err = rdb.HSet(ctx, SRS_VOD_M3U8_ARTIFACT, artifact.UUID, string(b)).Err(); err != nil && err != redis.Nil {
						return errors.Wrapf(err, "hset %v %v %v", SRS_VOD_M3U8_ARTIFACT, artifact.UUID, string(b))
					}
					logger.Tf(ctx, "vod update task %v", artifact.String())
				}
			}

			// Build response from artifacts.
			files := []map[string]interface{}{}
			for _, metadata := range artifacts {
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
					// For VoD only.
					"file":  metadata.FileID,
					"media": metadata.MediaURL,
					"task":  metadata.Task,
				})
			}

			ohttp.WriteData(ctx, w, r, files)
			logger.Tf(ctx, "vod files ok, cursor=%v, token=%vB", cursor, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/hooks/vod/hls/"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			// Format is :uuid.m3u8 or :uuid/index.m3u8
			filename := r.URL.Path[len("/terraform/v1/hooks/vod/hls/"):]
			// Format is :uuid.m3u8
			filename = strings.ReplaceAll(filename, "/index.m3u8", ".m3u8")
			uuid := filename[:len(filename)-len(path.Ext(filename))]
			if len(uuid) == 0 {
				return errors.Errorf("invalid uuid %v from %v of %v", uuid, filename, r.URL.Path)
			}

			var metadata M3u8VoDArtifact
			if m3u8Metadata, err := rdb.HGet(ctx, SRS_VOD_M3U8_ARTIFACT, uuid).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_VOD_M3U8_ARTIFACT, uuid)
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
			logger.Tf(ctx, "vod generate m3u8 ok, uuid=%v, duration=%v", uuid, duration)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	return nil
}

func (v *VodWorker) OnHlsTsMessage(ctx context.Context, msg *SrsOnHlsMessage) error {
	// Ignore for Tencent Cloud credentials not ready.
	if !v.ready() {
		return nil
	}

	// Copy the ts file to temporary cache dir.
	tsid := uuid.NewString()
	tsfile := path.Join("vod", fmt.Sprintf("%v.ts", tsid))

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

func (v *VodWorker) Close() error {
	if v.cancel != nil {
		v.cancel()
	}
	v.wg.Wait()
	return nil
}

func (v *VodWorker) Start(ctx context.Context) error {
	wg := &v.wg

	ctx, cancel := context.WithCancel(ctx)
	v.cancel = cancel

	ctx = logger.WithContext(ctx)
	logger.Tf(ctx, "Vod: start a worker")

	// When system startup, we initialize the credentials, for worker to use it.
	if err := v.updateCredential(ctx); err != nil {
		return errors.Wrapf(err, "update credential")
	}

	// Load all objects from redis.
	if objs, err := rdb.HGetAll(ctx, SRS_VOD_M3U8_WORKING).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hgetall %v", SRS_VOD_M3U8_WORKING)
	} else if len(objs) > 0 {
		for m3u8URL, value := range objs {
			logger.Tf(ctx, "Load %v object %v", m3u8URL, value)

			var m3u8LocalObj VodM3u8Stream
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
		logger.Tf(ctx, "Vod: Got message %v", msg.String())

		// Load stream local object.
		var m3u8LocalObj *VodM3u8Stream
		var freshObject bool
		if obj, loaded := v.streams.LoadOrStore(msg.Msg.M3u8URL, &VodM3u8Stream{
			M3u8URL: msg.Msg.M3u8URL, UUID: uuid.NewString(), vodWorker: v,
		}); true {
			m3u8LocalObj, freshObject = obj.(*VodM3u8Stream), !loaded
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

func (v *VodWorker) updateCredential(ctx context.Context) error {
	previous := &struct {
		SecretId, SecretKey string
		VodAppID            uint64
	}{
		SecretId: v.secretId, SecretKey: v.secretKey, VodAppID: v.vodAppID,
	}

	// The credential might not be ready, so we ignore error.
	if secretId, err := rdb.HGet(ctx, SRS_TENCENT_CAM, "secretId").Result(); err == nil {
		v.secretId = secretId
	}

	if secretKey, err := rdb.HGet(ctx, SRS_TENCENT_CAM, "secretKey").Result(); err == nil {
		v.secretKey = secretKey
	}

	if service, err := rdb.HGet(ctx, SRS_TENCENT_VOD, "service").Result(); err == nil && service != "ok" {
		if tv, err := strconv.ParseInt(service, 10, 64); err != nil {
			return errors.Wrapf(err, "parse vod appid %v", service)
		} else {
			v.vodAppID = uint64(tv)
		}
	}

	changed := v.secretId != previous.SecretId || v.secretKey != previous.SecretKey || v.vodAppID != previous.VodAppID
	credentialOK := v.secretId != "" && v.secretKey != "" && v.vodAppID > 0
	if (v.vodClient == nil || changed) && credentialOK {
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = TENCENT_CLOUD_VOD_ENDPOINT

		client, err := vod.NewClient(common.NewCredential(v.secretId, v.secretKey), conf.Region, cpf)
		if err != nil {
			return errors.Wrapf(err, "create vod sdk client, region=%v", conf.Region)
		}
		v.vodClient = client
		logger.Tf(ctx, "create vod client ok, appid=%v", v.vodAppID)
	}
	return nil
}

// VodCosToken is the token of COS for VoD.
type VodCosToken struct {
	M3u8URL string               `json:"m3u8_url"`
	UUID    string               `json:"uuid"`
	Bucket  string               `json:"bucket"`
	Region  string               `json:"region"`
	Key     string               `json:"key"`
	Session string               `json:"session"`
	Cert    *vod.TempCertificate `json:"cert"`
	Update  string               `json:"update"`
}

func (v *VodCosToken) String() string {
	var token string
	if v.Cert.Token != nil {
		token = *v.Cert.Token
	}
	return fmt.Sprintf("url=%v, uuid=%v, bucket=%v, region=%v, key=%v, session=%v, cert/token=%vB, update=%v",
		v.M3u8URL, v.UUID, v.Bucket, v.Region, v.Key, v.Session, len(token), v.Update,
	)
}

// VodM3u8Stream is the current active local object for a HLS stream.
// When voding done, it will generate a M3u8VoDArtifact, which is a HLS VoD object.
type VodM3u8Stream struct {
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
	vodWorker *VodWorker
	// The artifact we're working for.
	artifact *M3u8VoDArtifact
	// To protect the fields.
	lock sync.Mutex
}

func (v VodM3u8Stream) String() string {
	return fmt.Sprintf("url=%v, uuid=%v, done=%v, update=%v, messages=%v",
		v.M3u8URL, v.UUID, v.Done, v.Update, len(v.Messages),
	)
}

func (v *VodM3u8Stream) deleteObject(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if err := rdb.HDel(ctx, SRS_VOD_M3U8_WORKING, v.M3u8URL).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hdel %v %v", SRS_VOD_M3U8_WORKING, v.M3u8URL)
	}

	return nil
}

func (v *VodM3u8Stream) saveObject(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if b, err := json.Marshal(v); err != nil {
		return errors.Wrapf(err, "marshal object")
	} else if err = rdb.HSet(ctx, SRS_VOD_M3U8_WORKING, v.M3u8URL, string(b)).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_VOD_M3U8_WORKING, v.M3u8URL, string(b))
	}
	return nil
}

func (v *VodM3u8Stream) saveArtifact(ctx context.Context, artifact *M3u8VoDArtifact) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if b, err := json.Marshal(artifact); err != nil {
		return errors.Wrapf(err, "marshal %v", artifact.String())
	} else if err = rdb.HSet(ctx, SRS_VOD_M3U8_ARTIFACT, v.UUID, string(b)).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_VOD_M3U8_ARTIFACT, v.UUID, string(b))
	}
	return nil
}

func (v *VodM3u8Stream) updateArtifact(ctx context.Context, artifact *M3u8VoDArtifact, msg *SrsOnHlsObject) {
	v.lock.Lock()
	defer v.lock.Unlock()

	artifact.Vhost = msg.Msg.Vhost
	artifact.App = msg.Msg.App
	artifact.Stream = msg.Msg.Stream

	artifact.Files = append(artifact.Files, msg.TsFile)
	artifact.NN = len(artifact.Files)

	artifact.Update = time.Now().Format(time.RFC3339)
}

func (v *VodM3u8Stream) finishArtifact(ctx context.Context, artifact *M3u8VoDArtifact) {
	v.lock.Lock()
	defer v.lock.Unlock()

	artifact.Processing = false
	artifact.Update = time.Now().Format(time.RFC3339)
}

func (v *VodM3u8Stream) addMessage(ctx context.Context, msg *SrsOnHlsObject) {
	v.lock.Lock()
	defer v.lock.Unlock()

	v.Messages = append(v.Messages, msg)
	v.NN = len(v.Messages)
	v.Update = time.Now().Format(time.RFC3339)
}

func (v *VodM3u8Stream) copyMessages() []*SrsOnHlsObject {
	v.lock.Lock()
	defer v.lock.Unlock()

	return append([]*SrsOnHlsObject{}, v.Messages...)
}

func (v *VodM3u8Stream) vodCommit(ctx context.Context, artifact *M3u8VoDArtifact, fileID, mediaURL string) {
	v.lock.Lock()
	defer v.lock.Unlock()

	artifact.FileID = fileID
	artifact.MediaURL = mediaURL
}

func (v *VodM3u8Stream) vodRemux(ctx context.Context, artifact *M3u8VoDArtifact, definition uint64, taskID string) {
	v.lock.Lock()
	defer v.lock.Unlock()

	artifact.Definition = definition
	artifact.TaskID = taskID
}

func (v *VodM3u8Stream) removeMessage(ctx context.Context, msg *SrsOnHlsObject) {
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

func (v *VodM3u8Stream) expired() bool {
	v.lock.Lock()
	defer v.lock.Unlock()

	update, err := time.Parse(time.RFC3339, v.Update)
	if err != nil {
		return true
	}

	duration := 30 * time.Second
	if os.Getenv("NODE_ENV") != "development" {
		duration = 300 * time.Second
	}

	if update.Add(duration).Before(time.Now()) {
		return true
	}

	return false
}

// Initialize to load artifact. There is no simultaneously access, so no lock is needed.
func (v *VodM3u8Stream) Initialize(ctx context.Context, r *VodWorker) error {
	v.vodWorker = r
	logger.Tf(ctx, "vod initialize url=%v, uuid=%v", v.M3u8URL, v.UUID)

	// Try to load artifact from redis. The final artifact is VoD HLS object.
	if value, err := rdb.HGet(ctx, SRS_VOD_M3U8_ARTIFACT, v.UUID).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_VOD_M3U8_ARTIFACT, v.UUID)
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
			Update:     time.Now().Format(time.RFC3339),
			UUID:       v.UUID,
			M3u8URL:    v.M3u8URL,
			Processing: true,
		}
		if err := v.saveArtifact(ctx, v.artifact); err != nil {
			return errors.Wrapf(err, "save artifact %v", v.artifact.String())
		}
	}

	return nil
}

// Run to serve the current voding object.
func (v *VodM3u8Stream) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(logger.WithContext(ctx))
	logger.Tf(ctx, "vod run task %v", v.String())

	var cosClient *cos.Client
	var cosToken *VodCosToken

	pfn := func() error {
		// Refresh cos client.
		if tc, tt, err := v.refreshCosClient(ctx, cosClient, cosToken); err != nil {
			return errors.Wrapf(err, "refresh cos client")
		} else {
			cosClient, cosToken = tc, tt
		}

		// Process message and remove it.
		msgs := v.copyMessages()
		for _, msg := range msgs {
			if err := v.serveMessage(ctx, msg, cosClient, cosToken); err != nil {
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
		if err := v.finishM3u8(ctx, cosClient, cosToken); err != nil {
			return errors.Wrapf(err, "finish m3u8")
		}

		// Now HLS is done
		logger.Tf(ctx, "Vod is done, hls is %v, artifact is %v", v.String(), v.artifact.String())
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

func (v *VodM3u8Stream) refreshCosClient(ctx context.Context, oldClient *cos.Client, oldToken *VodCosToken) (*cos.Client, *VodCosToken, error) {
	var update string
	if oldToken != nil {
		update = oldToken.Update
	}

	refreshCosToken := func(ctx context.Context, cosToken *VodCosToken) (*VodCosToken, error) {
		if cosToken == nil {
			if token, err := rdb.HGet(ctx, SRS_VOD_COS_TOKEN, v.UUID).Result(); err != nil && err != redis.Nil {
				return nil, errors.Wrapf(err, "hget %v %v", SRS_VOD_COS_TOKEN, v.UUID)
			} else if token != "" {
				cosToken = &VodCosToken{}
				if err = json.Unmarshal([]byte(token), cosToken); err != nil {
					return nil, errors.Wrapf(err, "unmarshal %v", token)
				}
			}
		}

		// If not expired, reuse the session.
		if cosToken != nil && cosToken.Update != "" {
			duration := 1800 * time.Second
			if os.Getenv("NODE_ENV") == "development" {
				duration = 30 * time.Second
			}

			if update, err := time.Parse(time.RFC3339, cosToken.Update); err == nil {
				if update.Add(duration).After(time.Now()) {
					return cosToken, nil
				}
				logger.Tf(ctx, "VoD session expired, update=%v", update.Format(time.RFC3339))
			}
		}

		// See https://cloud.tencent.com/document/product/266/31767
		if true {
			type ExtendedApplyUploadRequest struct {
				*vod.ApplyUploadRequest
				// To use the same session, to extend the expire time.
				VodSessionKey string `json:"VodSessionKey,omitempty"`
			}

			request := &ExtendedApplyUploadRequest{}
			request.ApplyUploadRequest = vod.NewApplyUploadRequest()
			request.MediaType = common.StringPtr("m3u8")
			request.SubAppId = common.Uint64Ptr(v.vodWorker.vodAppID)
			request.SetContext(ctx)

			// To reuse the same session.
			if cosToken != nil {
				request.VodSessionKey = cosToken.Session
			}

			response := vod.NewApplyUploadResponse()
			if err := v.vodWorker.vodClient.Send(request, response); err != nil {
				return nil, errors.Wrapf(err, "apply upload")
			} else if res := response.Response; res == nil || res.TempCertificate == nil {
				return nil, errors.Errorf("empty cert for %v", response.ToJsonString())
			} else {
				cosToken = &VodCosToken{
					M3u8URL: v.M3u8URL,
					UUID:    v.UUID,
					Bucket:  *response.Response.StorageBucket,
					Region:  *response.Response.StorageRegion,
					Key:     *response.Response.MediaStoragePath,
					Session: *response.Response.VodSessionKey,
					Cert:    response.Response.TempCertificate,
					Update:  time.Now().Format(time.RFC3339),
				}
				logger.Tf(ctx, "vod refresh token ok, update=%v", cosToken.Update)
			}
		}

		// Save token to redis.
		if b, err := json.Marshal(cosToken); err != nil {
			return nil, errors.Wrapf(err, "marshal %v", cosToken.String())
		} else if err := rdb.HSet(ctx, SRS_VOD_COS_TOKEN, v.UUID, string(b)).Err(); err != nil && err != redis.Nil {
			return nil, errors.Wrapf(err, "hset %v %v %v", SRS_VOD_COS_TOKEN, v.UUID, string(b))
		}

		return cosToken, nil
	}

	// Refresh the COS token to upload file.
	cosToken, err := refreshCosToken(ctx, oldToken)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "refresh cos token")
	}

	// No changed, reuse previous one.
	if update == cosToken.Update && oldClient != nil {
		return oldClient, cosToken, nil
	}

	// Create new COS client.
	location := fmt.Sprintf("%v.cos.%v.myqcloud.com", cosToken.Bucket, cosToken.Region)
	u, err := url.Parse(fmt.Sprintf("https://%v", location))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "parse %v", fmt.Sprintf("https://%v", location))
	}

	// Create the bucket with private ACL.
	// See https://cloud.tencent.com/document/product/436/65639
	expire := time.Duration(*cosToken.Cert.ExpiredTime-uint64(time.Now().Unix())) * time.Second
	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID: *cosToken.Cert.SecretId, SecretKey: *cosToken.Cert.SecretKey,
			SessionToken: *cosToken.Cert.Token, Expire: expire,
		},
	})
	logger.Tf(ctx, "vod create cos client ok, bucket=%v, location=%v", cosToken.Bucket, location)

	return client, cosToken, nil
}

func (v *VodM3u8Stream) serveMessage(ctx context.Context, msg *SrsOnHlsObject, cosClient *cos.Client, cosToken *VodCosToken) error {
	// We always remove the msg from current object.
	defer v.removeMessage(ctx, msg)

	// Ignore file if credential is not ready.
	if !v.vodWorker.ready() {
		return nil
	}

	// Ignore file if not exists.
	if _, err := os.Stat(msg.TsFile.File); err != nil {
		return err
	}

	// Upload file to COS of VoD.
	tsDir := path.Dir(cosToken.Key)
	key := fmt.Sprintf("%v/%v.ts", tsDir, msg.TsFile.TsID)
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
	if _, err = cosClient.Object.Put(ctx, key, f, &opt); err != nil {
		return errors.Wrapf(err, "cos put object %v", key)
	}

	// Update the metadata for m3u8.
	v.updateArtifact(ctx, v.artifact, msg)
	if err := v.saveArtifact(ctx, v.artifact); err != nil {
		return errors.Wrapf(err, "save artifact %v", v.artifact.String())
	}

	logger.Tf(ctx, "vod consume msg %v", msg.String())
	return nil
}

func (v *VodM3u8Stream) finishM3u8(ctx context.Context, cosClient *cos.Client, cosToken *VodCosToken) error {
	contentType, m3u8Body, duration, err := buildVodM3u8(ctx, v.artifact, false, "", false, "")
	if err != nil {
		return errors.Wrapf(err, "build vod")
	}

	// Upload to COS bucket.
	// See https://cloud.tencent.com/document/product/436/64980
	opt := cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType:   contentType,
			ContentLength: int64(len(m3u8Body)),
		},
	}
	if _, err = cosClient.Object.Put(ctx, cosToken.Key, strings.NewReader(m3u8Body), &opt); err != nil {
		return errors.Wrapf(err, "cos put object %v", cosToken.Key)
	}
	logger.Tf(ctx, "vod to %v, duration=%v ok", cosToken.Key, duration)

	// Commit the upload for cloud VOD.
	// See https://cloud.tencent.com/document/product/266/31766
	if true {
		request := vod.NewCommitUploadRequest()

		request.VodSessionKey = common.StringPtr(cosToken.Session)
		request.SubAppId = common.Uint64Ptr(v.vodWorker.vodAppID)

		if response, err := v.vodWorker.vodClient.CommitUploadWithContext(ctx, request); err != nil {
			return errors.Wrapf(err, "vod commit, key=%v", cosToken.Key)
		} else {
			v.vodCommit(ctx, v.artifact, *response.Response.FileId, *response.Response.MediaUrl)
		}
	}

	// Start a remux task to covert HLS to MP4.
	// See https://cloud.tencent.com/document/product/266/33427
	var definition int64
	if remux, err := rdb.HGet(ctx, SRS_TENCENT_VOD, "remux").Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v remux", SRS_TENCENT_VOD)
	} else if remux != "" {
		transcode := &VodTranscodeTemplate{}
		if err = json.Unmarshal([]byte(remux), transcode); err != nil {
			return errors.Wrapf(err, "unmarshal %v", remux)
		}

		definition, err = strconv.ParseInt(transcode.Definition, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "invalid remux %v", transcode.String())
		}
	}

	var taskID string
	if definition > 0 {
		request := vod.NewProcessMediaRequest()

		request.FileId = common.StringPtr(v.artifact.FileID)
		request.SubAppId = common.Uint64Ptr(v.vodWorker.vodAppID)
		request.MediaProcessTask = &vod.MediaProcessTaskInput{
			TranscodeTaskSet: []*vod.TranscodeTaskInput{
				&vod.TranscodeTaskInput{
					Definition: common.Uint64Ptr(uint64(definition)),
				},
			},
		}

		if response, err := v.vodWorker.vodClient.ProcessMediaWithContext(ctx, request); err != nil {
			return errors.Wrapf(err, "vod remux")
		} else {
			taskID = *response.Response.TaskId
		}
	}

	if definition > 0 || taskID != "" {
		v.vodRemux(ctx, v.artifact, uint64(definition), taskID)
		logger.Tf(ctx, "vod remux ok, definition=%v, taskID=%v", definition, taskID)
	}

	// Remove object from worker.
	v.vodWorker.streams.Delete(v.M3u8URL)

	// Update artifact after finally.
	v.finishArtifact(ctx, v.artifact)
	r0 := v.saveArtifact(ctx, v.artifact)
	r1 := v.deleteObject(ctx)
	logger.Tf(ctx, "vod cleanup ok, r0=%v, r1=%v", r0, r1)

	// Do final cleanup, because new messages might arrive while converting to mp4, which takes a long time.
	files := v.copyMessages()
	for _, file := range files {
		r2 := os.Remove(file.TsFile.File)
		logger.Tf(ctx, "drop %v r2=%v", file.String(), r2)
	}

	return nil
}
