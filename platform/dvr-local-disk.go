//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
package main

import (
	"context"
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

	"github.com/ossrs/go-oryx-lib/errors"
	ohttp "github.com/ossrs/go-oryx-lib/http"
	"github.com/ossrs/go-oryx-lib/logger"

	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type RecordPostProcess string

const (
	RecordPostProcessCpFile RecordPostProcess = "post-cp-file"
)

var recordWorker *RecordWorker

type RecordWorker struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Got message from SRS, a new TS segment file is generated.
	msgs chan *SrsOnHlsObject
	// The streams we're recording, key is m3u8 URL in string, value is m3u8 object *RecordM3u8Stream.
	streams sync.Map
}

func NewRecordWorker() *RecordWorker {
	return &RecordWorker{
		msgs: make(chan *SrsOnHlsObject, 1024),
	}
}

func (v *RecordWorker) Handle(ctx context.Context, handler *http.ServeMux) error {
	ep := "/terraform/v1/hooks/record/query"
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

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if all, err := rdb.HGet(ctx, SRS_RECORD_PATTERNS, "all").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v all", SRS_RECORD_PATTERNS)
			} else if globs, err := rdb.HGet(ctx, SRS_RECORD_PATTERNS, "globs").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v globs", SRS_RECORD_PATTERNS)
			} else if processCpDir, err := rdb.HGet(ctx, SRS_RECORD_PATTERNS, string(RecordPostProcessCpFile)).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", string(RecordPostProcessCpFile))
			} else {
				globFilters := []string{}
				if globs != "" {
					if err := json.Unmarshal([]byte(globs), &globFilters); err != nil {
						return errors.Wrapf(err, "parse %v", globs)
					}
				}

				type RecordQueryResult struct {
					// Whether enable record.
					All bool `json:"all"`
					// The home directory of record.
					Home string `json:"home"`
					// The glob filters for record.
					Globs []string `json:"globs"`
					// The post process to copy file to dir for record.
					ProcessCpDir string `json:"processCpDir"`
				}

				ohttp.WriteData(ctx, w, r, &RecordQueryResult{
					All: all == "true", Home: "/data/record", Globs: globFilters,
					ProcessCpDir: processCpDir,
				})
			}

			logger.Tf(ctx, "record query ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/hooks/record/apply"
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

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if err := rdb.HSet(ctx, SRS_RECORD_PATTERNS, "all", fmt.Sprintf("%v", all)).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v all %v", SRS_RECORD_PATTERNS, all)
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "record apply ok, all=%v, token=%vB", all, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/hooks/record/globs"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var globs []string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string   `json:"token"`
				Globs *[]string `json:"globs"`
			}{
				Token: &token, Globs: &globs,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			filteredGlobs := []string{}
			for _, glob := range globs {
				if glob != "" {
					filteredGlobs = append(filteredGlobs, glob)
				}
			}

			if b, err := json.Marshal(filteredGlobs); err != nil {
				return errors.Wrapf(err, "marshal %v", filteredGlobs)
			} else if err := rdb.HSet(ctx, SRS_RECORD_PATTERNS, "globs", string(b)).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v globs %v", SRS_RECORD_PATTERNS, string(b))
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "record update globs ok, glob=%v, token=%vB", filteredGlobs, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/hooks/record/post-processing"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var postProcess, PostCpDir string
			if err := ParseBody(ctx, r.Body, &struct {
				Token       *string `json:"token"`
				PostProcess *string `json:"postProcess"`
				PostCpDir   *string `json:"postCpDir"`
			}{
				Token: &token, PostProcess: &postProcess, PostCpDir: &PostCpDir,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if RecordPostProcess(postProcess) != RecordPostProcessCpFile {
				return errors.Errorf("invalid post process %v", postProcess)
			}
			if PostCpDir != "" {
				if _, err := os.Stat(PostCpDir); err != nil {
					return errors.Wrapf(err, "stat dir %v", PostCpDir)
				}
			}

			if err := rdb.HSet(ctx, SRS_RECORD_PATTERNS, string(RecordPostProcessCpFile), PostCpDir).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v %v %v", SRS_RECORD_PATTERNS, RecordPostProcessCpFile, PostCpDir)
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "record update post processing ok, postProcess=%v, postCpDir=%v, token=%vB",
				postProcess, PostCpDir, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/hooks/record/remove"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token, uuid string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				UUID  *string `json:"uuid"`
			}{
				Token: &token, UUID: &uuid,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if uuid == "" {
				return errors.New("no uuid")
			}

			var metadata M3u8VoDArtifact
			if M3u8VoDMetadata, err := rdb.HGet(ctx, SRS_RECORD_M3U8_ARTIFACT, uuid).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_RECORD_M3U8_ARTIFACT, uuid)
			} else if M3u8VoDMetadata == "" {
				return errors.Errorf("no record for uuid=%v", uuid)
			} else if err = json.Unmarshal([]byte(M3u8VoDMetadata), &metadata); err != nil {
				return errors.Wrapf(err, "parse %v", M3u8VoDMetadata)
			}

			// Remove all ts files.
			for _, file := range metadata.Files {
				if _, err := os.Stat(file.Key); err == nil {
					os.Remove(file.Key)
				}
			}

			// Remove m3u8 file.
			m3u8File := path.Join("record", uuid, "index.m3u8")
			if _, err := os.Stat(m3u8File); err == nil {
				os.Remove(m3u8File)
			}

			// Remove mp4 file.
			mp4File := path.Join("record", uuid, "index.mp4")
			if _, err := os.Stat(mp4File); err == nil {
				os.Remove(mp4File)
			}

			// Remove ts directory.
			m3u8Directory := path.Join("record", uuid)
			if _, err := os.Stat(m3u8Directory); err == nil {
				os.RemoveAll(m3u8Directory)
			}

			// Remove HLS from list.
			if err := rdb.HDel(ctx, SRS_RECORD_M3U8_ARTIFACT, uuid).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hdel %v %v", SRS_RECORD_M3U8_ARTIFACT, uuid)
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "record remove ok, uuid=%v, token=%vB", uuid, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/hooks/record/end"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token, uuid string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				UUID  *string `json:"uuid"`
			}{
				Token: &token, UUID: &uuid,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if uuid == "" {
				return errors.New("no uuid")
			}

			task := recordWorker.QueryTask(uuid)
			if task == nil {
				return errors.Errorf("no record task for uuid=%v", uuid)
			}

			// Make the task to expire, to end it ASAP.
			task.Expired = true

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "record end ok, uuid=%v, token=%vB", uuid, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/hooks/record/files"
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

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			keys, cursor, err := rdb.HScan(ctx, SRS_RECORD_M3U8_ARTIFACT, 0, "*", 100).Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hscan %v 0 * 100", SRS_RECORD_M3U8_ARTIFACT)
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
				})
			}

			ohttp.WriteData(ctx, w, r, files)
			logger.Tf(ctx, "record files ok, cursor=%v, token=%vB", cursor, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	m3u8Handler := func(w http.ResponseWriter, r *http.Request) error {
		// Format is :uuid.m3u8 or :uuid/index.m3u8
		filename := r.URL.Path[len("/terraform/v1/hooks/record/hls/"):]
		// Format is :uuid.m3u8
		filename = strings.ReplaceAll(filename, "/index.m3u8", ".m3u8")
		uuid := filename[:len(filename)-len(path.Ext(filename))]
		if len(uuid) == 0 {
			return errors.Errorf("invalid uuid %v from %v of %v", uuid, filename, r.URL.Path)
		}

		var metadata M3u8VoDArtifact
		if m3u8Metadata, err := rdb.HGet(ctx, SRS_RECORD_M3U8_ARTIFACT, uuid).Result(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hget %v %v", SRS_RECORD_M3U8_ARTIFACT, uuid)
		} else if m3u8Metadata == "" {
			return errors.Errorf("no m3u8 of uuid=%v", uuid)
		} else if err = json.Unmarshal([]byte(m3u8Metadata), &metadata); err != nil {
			return errors.Wrapf(err, "parse %v", m3u8Metadata)
		}

		prefix := "/terraform/v1/hooks/record/hls/"
		contentType, m3u8Body, duration, err := buildVodM3u8ForLocal(ctx, metadata.Files, true, prefix)
		if err != nil {
			return errors.Wrapf(err, "build vod m3u8 of %v with prefix=%v", metadata.String(), prefix)
		}

		w.Header().Set("Content-Type", contentType)
		w.Write([]byte(m3u8Body))
		logger.Tf(ctx, "record generate m3u8 ok, uuid=%v, duration=%v", uuid, duration)
		return nil
	}

	tsHandler := func(w http.ResponseWriter, r *http.Request) error {
		// Format is :dir/:m3u8/:uuid.ts
		filename := r.URL.Path[len("/terraform/v1/hooks/record/hls/"):]
		fileDir, fileBase := path.Dir(filename), path.Base(filename)
		uuid := fileBase[:len(fileBase)-len(path.Ext(fileBase))]
		dir, m3u8 := path.Dir(fileDir), path.Base(fileDir)
		if len(uuid) == 0 {
			return errors.Errorf("invalid uuid %v from %v of %v", uuid, fileBase, r.URL.Path)
		}
		if len(dir) == 0 {
			return errors.Errorf("invalid dir %v from %v of %v", dir, fileDir, r.URL.Path)
		}
		if len(m3u8) == 0 {
			return errors.Errorf("invalid m3u8 %v from %v of %v", m3u8, fileDir, r.URL.Path)
		}

		tsFilePath := path.Join(dir, m3u8, fmt.Sprintf("%v.ts", uuid))
		if _, err := os.Stat(tsFilePath); err != nil {
			return errors.Wrapf(err, "no ts file %v", tsFilePath)
		}

		if tsFile, err := os.Open(tsFilePath); err != nil {
			return errors.Wrapf(err, "open file %v", tsFilePath)
		} else {
			defer tsFile.Close()
			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
			io.Copy(w, tsFile)
		}

		logger.Tf(ctx, "record server ts file ok, uuid=%v, ts=%v", uuid, tsFilePath)
		return nil
	}

	mp4Handler := func(w http.ResponseWriter, r *http.Request) error {
		// Format is :uuid/index.mp4
		filename := r.URL.Path[len("/terraform/v1/hooks/record/hls/"):]
		uuid := path.Dir(filename)
		if len(uuid) == 0 {
			return errors.Errorf("invalid uuid %v from %v of %v", uuid, filename, r.URL.Path)
		}

		var metadata M3u8VoDArtifact
		if m3u8Metadata, err := rdb.HGet(ctx, SRS_RECORD_M3U8_ARTIFACT, uuid).Result(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hget %v %v", SRS_RECORD_M3U8_ARTIFACT, uuid)
		} else if m3u8Metadata == "" {
			return errors.Errorf("no m3u8 of uuid=%v", uuid)
		} else if err = json.Unmarshal([]byte(m3u8Metadata), &metadata); err != nil {
			return errors.Wrapf(err, "parse %v", m3u8Metadata)
		}

		mp4FilePath := path.Join("record", uuid, "index.mp4")
		stats, err := os.Stat(mp4FilePath)
		if err != nil {
			return errors.Wrapf(err, "no mp4 file %v", mp4FilePath)
		}

		mp4File, err := os.Open(mp4FilePath)
		if err != nil {
			return errors.Wrapf(err, "open file %v", mp4FilePath)
		}
		defer mp4File.Close()

		// No range request.
		rangeHeader := r.Header.Get("Range")
		if rangeHeader == "" {
			w.Header().Set("Content-Type", "video/mp4")
			io.Copy(w, mp4File)
			logger.T(ctx, "record serve full mp4=%v", mp4FilePath)
			return nil
		}

		// Support range request.
		var start, end int64
		fmt.Sscanf(rangeHeader, "bytes=%u-%u", &start, &end)
		if end == 0 {
			end = stats.Size() - 1
		}

		if _, err := mp4File.Seek(start, io.SeekStart); err != nil {
			return errors.Wrapf(err, "seek to %v of %v", start, mp4FilePath)
		}

		w.Header().Set("ccept-Ranges", "bytes")
		w.Header().Set("Content-Length", fmt.Sprintf("%v", end+1-start))
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %v-%v/%v", start, end, stats.Size()))

		w.WriteHeader(http.StatusPartialContent)
		w.Header().Set("Content-Type", "video/mp4")
		io.CopyN(w, mp4File, end+1-start)

		logger.Tf(ctx, "record serve partial ok, uuid=%v, mp4=%v", uuid, mp4FilePath)
		return nil
	}

	ep = "/terraform/v1/hooks/record/hls/"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			if strings.HasSuffix(r.URL.Path, ".m3u8") {
				return m3u8Handler(w, r)
			} else if strings.HasSuffix(r.URL.Path, ".ts") {
				return tsHandler(w, r)
			} else if strings.HasSuffix(r.URL.Path, ".mp4") {
				return mp4Handler(w, r)
			}

			return errors.Errorf("invalid handler for %v", r.URL.Path)
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	return nil
}

func (v *RecordWorker) OnHlsTsMessage(ctx context.Context, msg *SrsOnHlsMessage) error {
	// Copy the ts file to temporary cache dir.
	tsid := uuid.NewString()
	tsfile := path.Join("record", fmt.Sprintf("%v.ts", tsid))

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

func (v *RecordWorker) Close() error {
	if v.cancel != nil {
		v.cancel()
	}
	v.wg.Wait()
	return nil
}

func (v *RecordWorker) QueryTask(uuid string) *RecordM3u8Stream {
	var target *RecordM3u8Stream
	v.streams.Range(func(key, value interface{}) bool {
		if task := value.(*RecordM3u8Stream); task.UUID == uuid {
			target = task
			return false
		}
		return true
	})
	return target
}

func (v *RecordWorker) Start(ctx context.Context) error {
	wg := &v.wg

	ctx, cancel := context.WithCancel(ctx)
	v.cancel = cancel

	ctx = logger.WithContext(ctx)
	logger.Tf(ctx, "Record: start a worker")

	// Load all objects from redis.
	if objs, err := rdb.HGetAll(ctx, SRS_RECORD_M3U8_WORKING).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hgetall %v", SRS_RECORD_M3U8_WORKING)
	} else if len(objs) > 0 {
		for m3u8URL, value := range objs {
			logger.Tf(ctx, "Load %v object %v", m3u8URL, value)

			var m3u8LocalObj RecordM3u8Stream
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
		logger.Tf(ctx, "Record: Got message %v", msg.String())

		// Filter the stream by glob filters.
		var globFilters []string
		if globs, err := rdb.HGet(ctx, SRS_RECORD_PATTERNS, "globs").Result(); err != nil && err != redis.Nil {
			return errors.Wrapf(err, "hget %v globs", SRS_RECORD_PATTERNS)
		} else if globs != "" {
			if err := json.Unmarshal([]byte(globs), &globFilters); err != nil {
				return errors.Wrapf(err, "parse %v", globs)
			}
		}

		// If glob filters are empty, ignore it, and record all streams.
		if len(globFilters) > 0 {
			var globMatched bool
			streamURL := fmt.Sprintf("/%v/%v", msg.Msg.App, msg.Msg.Stream)
			for _, globFilter := range globFilters {
				if ok, err := path.Match(globFilter, streamURL); err != nil {
					return errors.Wrapf(err, "match %v", globFilter)
				} else if ok {
					logger.Tf(ctx, "match stream %v by glob filter %v in %v", streamURL, globFilter, globFilters)
					globMatched = true
				}
			}

			if !globMatched {
				logger.Wf(ctx, "ignore stream %v by glob filters %v", streamURL, globFilters)
				return nil
			}
		}

		// Load stream local object.
		var m3u8LocalObj *RecordM3u8Stream
		var freshObject bool
		if obj, loaded := v.streams.LoadOrStore(msg.Msg.M3u8URL, &RecordM3u8Stream{
			M3u8URL: msg.Msg.M3u8URL, UUID: uuid.NewString(), recordWorker: v,
		}); true {
			m3u8LocalObj, freshObject = obj.(*RecordM3u8Stream), !loaded
		}

		// Initialize the fresh object.
		if freshObject {
			if err := m3u8LocalObj.Initialize(ctx, v); err != nil {
				return errors.Wrapf(err, "init %v", m3u8LocalObj.String())
			}
		}

		// Append new ts file to object.
		m3u8LocalObj.addMessage(ctx, msg)

		// Always save the object to redis, for reloading it when restart.
		if err := m3u8LocalObj.saveObject(ctx); err != nil {
			return errors.Wrapf(err, "save %v", m3u8LocalObj.String())
		}

		// Serve object if fresh one.
		if freshObject {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := m3u8LocalObj.Run(ctx); err != nil {
					logger.Wf(ctx, "serve m3u8 %v err %+v", m3u8LocalObj.String(), err)
				}
			}()
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

	return nil
}

// RecordM3u8Stream is the current active local object for a HLS stream.
// When recording done, it will generate a M3u8VoDArtifact, which is a HLS VoD object.
type RecordM3u8Stream struct {
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
	// Whether task is set to expire by user.
	Expired bool `json:"expired"`

	// The ts files of this m3u8.
	Messages []*SrsOnHlsObject `json:"msgs"`

	// The worker which owns this object.
	recordWorker *RecordWorker
	// The artifact we're working for.
	artifact *M3u8VoDArtifact
	// To protect the fields.
	lock sync.Mutex
}

func (v RecordM3u8Stream) String() string {
	return fmt.Sprintf("url=%v, uuid=%v, done=%v, update=%v, messages=%v, expired=%v",
		v.M3u8URL, v.UUID, v.Done, v.Update, len(v.Messages), v.Expired,
	)
}

func (v *RecordM3u8Stream) deleteObject(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if err := rdb.HDel(ctx, SRS_RECORD_M3U8_WORKING, v.M3u8URL).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hdel %v %v", SRS_RECORD_M3U8_WORKING, v.M3u8URL)
	}

	return nil
}

func (v *RecordM3u8Stream) saveObject(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if b, err := json.Marshal(v); err != nil {
		return errors.Wrapf(err, "marshal object")
	} else if err = rdb.HSet(ctx, SRS_RECORD_M3U8_WORKING, v.M3u8URL, string(b)).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_RECORD_M3U8_WORKING, v.M3u8URL, string(b))
	}
	return nil
}

func (v *RecordM3u8Stream) saveArtifact(ctx context.Context, artifact *M3u8VoDArtifact) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if b, err := json.Marshal(artifact); err != nil {
		return errors.Wrapf(err, "marshal %v", artifact.String())
	} else if err = rdb.HSet(ctx, SRS_RECORD_M3U8_ARTIFACT, v.UUID, string(b)).Err(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_RECORD_M3U8_ARTIFACT, v.UUID, string(b))
	}
	return nil
}

func (v *RecordM3u8Stream) updateArtifact(ctx context.Context, artifact *M3u8VoDArtifact, msg *SrsOnHlsObject) {
	v.lock.Lock()
	defer v.lock.Unlock()

	artifact.Vhost = msg.Msg.Vhost
	artifact.App = msg.Msg.App
	artifact.Stream = msg.Msg.Stream

	artifact.Files = append(artifact.Files, msg.TsFile)
	artifact.NN = len(artifact.Files)

	artifact.Update = time.Now().Format(time.RFC3339)
}

func (v *RecordM3u8Stream) finishArtifact(ctx context.Context, artifact *M3u8VoDArtifact) {
	v.lock.Lock()
	defer v.lock.Unlock()

	artifact.Processing = false
	artifact.Update = time.Now().Format(time.RFC3339)
}

func (v *RecordM3u8Stream) addMessage(ctx context.Context, msg *SrsOnHlsObject) {
	v.lock.Lock()
	defer v.lock.Unlock()

	v.Messages = append(v.Messages, msg)
	v.NN = len(v.Messages)
	v.Update = time.Now().Format(time.RFC3339)
}

func (v *RecordM3u8Stream) copyMessages() []*SrsOnHlsObject {
	v.lock.Lock()
	defer v.lock.Unlock()

	return append([]*SrsOnHlsObject{}, v.Messages...)
}

func (v *RecordM3u8Stream) removeMessage(msg *SrsOnHlsObject) {
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
}

func (v *RecordM3u8Stream) expired(ctx context.Context) bool {
	v.lock.Lock()
	defer v.lock.Unlock()

	if v.Expired {
		return true
	}

	update, err := time.Parse(time.RFC3339, v.Update)
	if err != nil {
		return true
	}

	var enabled bool
	if all, err := rdb.HGet(ctx, SRS_RECORD_PATTERNS, "all").Result(); err == nil {
		enabled = all == "true"
	}

	duration := 30 * time.Second
	if enabled && os.Getenv("NODE_ENV") != "development" {
		duration = 300 * time.Second
	}

	if update.Add(duration).Before(time.Now()) {
		return true
	}

	return false
}

// Initialize to load artifact. There is no simultaneously access, so no certFileLock is needed.
func (v *RecordM3u8Stream) Initialize(ctx context.Context, r *RecordWorker) error {
	v.recordWorker = r
	logger.Tf(ctx, "record initialize url=%v, uuid=%v", v.M3u8URL, v.UUID)

	// Try to load artifact from redis. The final artifact is VoD HLS object.
	if value, err := rdb.HGet(ctx, SRS_RECORD_M3U8_ARTIFACT, v.UUID).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_RECORD_M3U8_ARTIFACT, v.UUID)
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
		if err := v.saveArtifact(ctx, v.artifact); err != nil {
			return errors.Wrapf(err, "save artifact %v", v.artifact.String())
		}
	}

	return nil
}

// Run to serve the current recording object.
func (v *RecordM3u8Stream) Run(ctx context.Context) error {
	parentCtx := logger.WithContext(ctx)
	ctx, cancel := context.WithCancel(parentCtx)
	logger.Tf(ctx, "record run task %v", v.String())

	if true {
		message, err := v.callbackBegin(ctx)
		if err != nil {
			logger.Wf(ctx, "ignore task %v callback start err %+v", v.String(), err)
		}

		defer func() {
			// Keep in mind that we should not use the ctx for this record task, as it will be canceled once the
			// task is completed. Instead, we need to use the parent ctx, which is the function parameter representing
			// the server context that remains active until the server is shut down.
			ctx := parentCtx

			if err := v.callbackEnd(ctx, message); err != nil {
				logger.Wf(ctx, "ignore task %v callback end err %+v", v.String(), err)
			}
		}()
	}

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
		if !v.expired(ctx) {
			return nil
		}

		// Try to finish the object.
		if err := v.finishM3u8(ctx); err != nil {
			return errors.Wrapf(err, "finish m3u8")
		}

		// Do post processing.
		if err := v.postProcessing(ctx); err != nil {
			return errors.Wrapf(err, "post processing")
		}

		// Now HLS is done
		logger.Tf(ctx, "Record is done, hls is %v, artifact is %v", v.String(), v.artifact.String())
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

func (v *RecordM3u8Stream) serveMessage(ctx context.Context, msg *SrsOnHlsObject) error {
	// We always remove the msg from current object.
	defer v.removeMessage(msg)

	// Ignore file if not exists.
	if _, err := os.Stat(msg.TsFile.File); err != nil {
		return err
	}

	tsDir := path.Join("record", v.UUID)
	key := path.Join(tsDir, fmt.Sprintf("%v.ts", msg.TsFile.TsID))
	msg.TsFile.Key = key

	if err := os.MkdirAll(tsDir, 0755); err != nil {
		return errors.Wrapf(err, "mkdir %v", tsDir)
	}

	if err := os.Rename(msg.TsFile.File, key); err != nil {
		return errors.Wrapf(err, "rename %v to %v", msg.TsFile.File, key)
	}

	// Update the metadata for m3u8.
	v.updateArtifact(ctx, v.artifact, msg)
	if err := v.saveArtifact(ctx, v.artifact); err != nil {
		return errors.Wrapf(err, "save artifact %v", v.artifact.String())
	}

	logger.Tf(ctx, "record consume msg %v", msg.String())
	return nil
}

func (v *RecordM3u8Stream) finishM3u8(ctx context.Context) error {
	contentType, m3u8Body, duration, err := buildVodM3u8ForLocal(ctx, v.artifact.Files, false, "")
	if err != nil {
		return errors.Wrapf(err, "build vod")
	}

	hls := path.Join("record", v.UUID, "index.m3u8")
	if f, err := os.OpenFile(hls, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return errors.Wrapf(err, "open file %v", hls)
	} else {
		defer f.Close()
		if _, err = f.Write([]byte(m3u8Body)); err != nil {
			return errors.Wrapf(err, "write hls %v to %v", m3u8Body, hls)
		}
	}
	logger.Tf(ctx, "record to %v ok, type=%v, duration=%v", hls, contentType, duration)

	mp4 := path.Join("record", v.UUID, "index.mp4")
	if b, err := exec.CommandContext(ctx, "ffmpeg", "-i", hls, "-c", "copy", "-y", mp4).Output(); err != nil {
		return errors.Wrapf(err, "covert to mp4 %v err %v", mp4, string(b))
	}
	logger.Tf(ctx, "record to %v ok", mp4)

	// Remove object from worker.
	v.recordWorker.streams.Delete(v.M3u8URL)

	// Update artifact after finally.
	v.finishArtifact(ctx, v.artifact)
	r0 := v.saveArtifact(ctx, v.artifact)
	r1 := v.deleteObject(ctx)
	logger.Tf(ctx, "record cleanup ok, r0=%v, r1=%v", r0, r1)

	// Do final cleanup, because new messages might arrive while converting to mp4, which takes a long time.
	files := v.copyMessages()
	for _, file := range files {
		r2 := os.Remove(file.TsFile.File)
		logger.Tf(ctx, "drop %v r2=%v", file.String(), r2)
	}

	return nil
}

func (v *RecordM3u8Stream) postProcessing(ctx context.Context) error {
	processCpDir, err := rdb.HGet(ctx, SRS_RECORD_PATTERNS, string(RecordPostProcessCpFile)).Result()
	if err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", string(RecordPostProcessCpFile))
	}

	if processCpDir == "" {
		return nil
	}

	artifactPath := path.Join("record", v.UUID, "index.mp4")
	targetPath := path.Join(processCpDir, fmt.Sprintf("%v.mp4", v.artifact.UUID))
	if err = exec.CommandContext(ctx, "cp", "-f", artifactPath, targetPath).Run(); err != nil {
		return errors.Wrapf(err, "cp %v to %v", artifactPath, targetPath)
	}
	logger.Tf(ctx, "record post process, cp %v to %v ok", artifactPath, targetPath)

	return nil
}

func (v *RecordM3u8Stream) callbackBegin(ctx context.Context) (*SrsOnHlsObject, error) {
	messages := v.copyMessages()
	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages")
	}

	message := messages[0]
	if err := callbackWorker.OnRecordMessage(ctx, SrsActionOnRecordBegin, v.UUID, message.Msg, nil); err != nil {
		return message, errors.Wrapf(err, "on record end %v", message)
	}

	return message, nil
}

func (v *RecordM3u8Stream) callbackEnd(ctx context.Context, message *SrsOnHlsObject) error {
	if message == nil {
		return nil
	}

	if err := callbackWorker.OnRecordMessage(ctx, SrsActionOnRecordEnd, v.UUID, message.Msg, v.artifact); err != nil {
		return errors.Wrapf(err, "on record end %v", message)
	}

	return nil
}
