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
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/ossrs/go-oryx-lib/errors"
	ohttp "github.com/ossrs/go-oryx-lib/http"
	"github.com/ossrs/go-oryx-lib/logger"
	"github.com/sashabaranov/go-openai"
)

var dubbingServer *SrsDubbingServer
var aiDubbingWorkDir = "containers/data/dubbing"

// MergeSegmentBetweenGroups is used for detect the gap for automatically merging.
const MergeSegmentBetweenGroups = 10 * time.Millisecond

// MergeSegmentForSmallWords is used for merging the group if less than this words.
const MergeSegmentForSmallWords = 30

// For production, all should be false.
const (
	alwaysForceRegenerateASRResponse = false
	onlyRegenerateFirstSegment       = false
	alwaysRephraseTranslations       = false
)

func handleDubbingService(ctx context.Context, handler *http.ServeMux) error {
	logger.Tf(ctx, "AI dubbing work dir: %v", aiDubbingWorkDir)

	ep := "/terraform/v1/dubbing/create"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token, title string
			var files []*FFprobeSource
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				// Project title.
				Title *string `json:"title"`
				// File type and path.
				Files *[]*FFprobeSource `json:"files"`
			}{
				Token: &token, Title: &title, Files: &files,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if len(files) != 1 {
				return errors.Errorf("invalid files %v", len(files))
			}

			targetFile := files[0]
			if targetFile == nil {
				return errors.Errorf("invalid file")
			}
			if targetFile.Type != FFprobeSourceTypeFile && targetFile.Type != FFprobeSourceTypeUpload && targetFile.Type != FFprobeSourceTypeYTDL {
				return errors.Errorf("invalid file type %v", targetFile.Type)
			}
			if targetFile.Target == "" {
				return errors.Errorf("invalid file path")
			}

			dubbing := NewSrsDubbingProject(func(dubbing *SrsDubbingProject) {
				dubbing.Title = title
				dubbing.FileType, dubbing.FilePath = targetFile.Type, targetFile.Path
			})

			if err := dubbing.CheckSource(ctx, targetFile.Target); err != nil {
				return errors.Wrapf(err, "check source type=%v, %v", targetFile.Type, targetFile.Target)
			}

			// Create the project home directory.
			if projectDir := path.Join(conf.Pwd, aiDubbingWorkDir, dubbing.UUID); true {
				if err := os.MkdirAll(projectDir, 0755); err != nil {
					return errors.Wrapf(err, "mkdir %v for dubbing project %v", projectDir, dubbing.String())
				}
			}

			if true {
				fileAbsPath, err := filepath.Abs(targetFile.Target)
				if err != nil {
					return errors.Wrapf(err, "abs %v", targetFile.Target)
				}

				info, err := os.Stat(fileAbsPath)
				if err != nil {
					return errors.Wrapf(err, "stat %v", fileAbsPath)
				}

				dubbing.SourceUUID = uuid.NewString()
				dubbing.SourcePath = path.Join(dubbing.UUID,
					fmt.Sprintf("%v%v", dubbing.SourceUUID, path.Ext(info.Name())),
				)

				absSourcePath := path.Join(conf.Pwd, aiDubbingWorkDir, dubbing.SourcePath)
				if err := os.Rename(fileAbsPath, absSourcePath); err != nil {
					return errors.Wrapf(err, "rename %v to %v", fileAbsPath, absSourcePath)
				}
			}

			dubbing.SourceVideo, dubbing.SourceAudio = targetFile.Video, targetFile.Audio
			dubbing.SourceFormat = &MediaFormat{}
			if err := dubbing.SourceFormat.FromFFprobeFormat(targetFile.Format); err != nil {
				return errors.Wrapf(err, "from ffprobe format %v", targetFile.Format)
			}

			if err := dubbing.Save(ctx); err != nil {
				return errors.Wrapf(err, "save dubbing project %v", dubbing.String())
			}

			ohttp.WriteData(ctx, w, r, &dubbing)
			logger.Tf(ctx, "srs dubbing create ok, title=%v, project=%v", title, dubbing.String())
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/dubbing/list"
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

			var projects []*SrsDubbingProject
			if configs, err := rdb.HGetAll(ctx, SRS_DUBBING_PROJECTS).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hgetall %v", SRS_DUBBING_PROJECTS)
			} else {
				for k, v := range configs {
					var obj SrsDubbingProject
					if err = json.Unmarshal([]byte(v), &obj); err != nil {
						return errors.Wrapf(err, "unmarshal %v %v", k, v)
					}
					projects = append(projects, &obj)
				}
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Projects []*SrsDubbingProject `json:"projects"`
			}{
				Projects: projects,
			})
			logger.Tf(ctx, "srs dubbing projects list ok, projects=%v", len(projects))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/dubbing/remove"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token, dubbingUUID string
			if err := ParseBody(ctx, r.Body, &struct {
				Token       *string `json:"token"`
				DubbingUUID *string `json:"uuid"`
			}{
				Token: &token, DubbingUUID: &dubbingUUID,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			dubbing := &SrsDubbingProject{UUID: dubbingUUID}
			if err := dubbing.Load(ctx); err != nil {
				return errors.Wrapf(err, "load dubbing project %v", dubbingUUID)
			}

			if err := rdb.HDel(ctx, SRS_DUBBING_PROJECTS, dubbingUUID).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hdel %v %v", SRS_DUBBING_PROJECTS, dubbingUUID)
			}

			// Remove the project files.
			if dubbing.UUID != "" {
				projectDir := path.Join(conf.Pwd, aiDubbingWorkDir, dubbing.UUID)
				if _, err := os.Stat(projectDir); err == nil {
					if err = os.RemoveAll(projectDir); err != nil {
						logger.Wf(ctx, "ignore dubbing err %v", err)
					}
				}
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "srs remove dubbing ok, uuid=%v", dubbingUUID)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/dubbing/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token, dubbingUUID string
			if err := ParseBody(ctx, r.Body, &struct {
				Token       *string `json:"token"`
				DubbingUUID *string `json:"uuid"`
			}{
				Token: &token, DubbingUUID: &dubbingUUID,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			dubbing := &SrsDubbingProject{UUID: dubbingUUID}
			if err := dubbing.Load(ctx); err != nil {
				return errors.Wrapf(err, "load dubbing project %v", dubbingUUID)
			}

			// If source file does not exists, reset the project.
			absSourcePath := path.Join(conf.Pwd, aiDubbingWorkDir, dubbing.SourcePath)
			if _, err := os.Stat(absSourcePath); err != nil {
				dubbing.SourceUUID, dubbing.SourcePath = "", ""
			}

			ohttp.WriteData(ctx, w, r, &dubbing)
			logger.Tf(ctx, "srs dubbing query ok, uuid=%v, dubbing=%v", dubbingUUID, dubbing.String())
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/dubbing/update"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var dubbing SrsDubbingProject
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				*SrsDubbingProject
			}{
				Token: &token, SrsDubbingProject: &dubbing,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			// TODO: FIXME: Should load dubbing from redis and merge the fields.
			if b, err := json.Marshal(dubbing); err != nil {
				return errors.Wrapf(err, "marshal dubbing")
			} else if err := rdb.HSet(ctx, SRS_DUBBING_PROJECTS, dubbing.UUID, string(b)).Err(); err != nil {
				return errors.Wrapf(err, "hset %v %v %v", SRS_DUBBING_PROJECTS, dubbing.UUID, string(b))
			}

			// Update the task if exists.
			task := dubbingServer.QueryTask(dubbing.TaskUUID)
			if task != nil {
				task.UpdateProject(&dubbing)
			}

			// Limit the changing rate for dubbing.
			select {
			case <-ctx.Done():
			case <-time.After(300 * time.Millisecond):
			}

			ohttp.WriteData(ctx, w, r, &dubbing)
			logger.Tf(ctx, "srs dubbing update ok, dubbing=%v", dubbing.String())
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/dubbing/play"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			q := r.URL.Query()

			token := q.Get("token")
			if token == "" {
				return errors.Errorf("empty token")
			}
			dubbingUUID := q.Get("uuid")
			if dubbingUUID == "" {
				return errors.Errorf("empty uuid")
			}

			// Convert the token in query to header Bearer token.
			r.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, "", r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			dubbing := &SrsDubbingProject{UUID: dubbingUUID}
			if err := dubbing.Load(ctx); err != nil {
				return errors.Wrapf(err, "load dubbing project %v", dubbingUUID)
			}

			// The source file to play.
			filename := path.Join(conf.Pwd, aiDubbingWorkDir, dubbing.SourcePath)
			if dubbing.SourceFormat == nil {
				return errors.Errorf("no format of source %v", filename)
			}

			ext := strings.Trim(path.Ext(filename), ".")
			contentType := fmt.Sprintf("audio/%v", ext)
			if dubbing.SourceFormat.HasVideo {
				contentType = fmt.Sprintf("video/%v", ext)
			}
			logger.Tf(ctx, "Serve example file=%v, ext=%v, contentType=%v", filename, ext, contentType)

			w.Header().Set("Content-Type", contentType)
			http.ServeFile(w, r, filename)

			ohttp.WriteData(ctx, w, r, &dubbing)
			logger.Tf(ctx, "srs dubbing get play src ok, uuid=%v, dubbing=%v", dubbingUUID, dubbing.String())
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/dubbing/export"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var dubbingUUID, taskUUID string
			if err := ParseBody(ctx, r.Body, &struct {
				Token    *string `json:"token"`
				UUID     *string `json:"uuid"`
				TaskUUID *string `json:"task"`
			}{
				Token: &token, UUID: &dubbingUUID, TaskUUID: &taskUUID,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			dubbing := &SrsDubbingProject{UUID: dubbingUUID}
			if err := dubbing.Load(ctx); err != nil {
				return errors.Wrapf(err, "load dubbing project %v", dubbingUUID)
			}

			// TODO: Allow task use different UUID from project.
			if dubbing.TaskUUID != taskUUID {
				return errors.Errorf("invalid task %v, should be %v", taskUUID, dubbing.TaskUUID)
			}

			// User should started the task.
			task := dubbingServer.QueryTask(taskUUID)
			if task == nil {
				return errors.Errorf("task %v not exists", taskUUID)
			}

			// Fail if task not finished.
			if task.status != SrsDubbingTaskStatusDone {
				return errors.Errorf("task %v not finished, status=%v", taskUUID, task.status)
			}

			asrResponse := task.AsrResponse
			if asrResponse == nil {
				return errors.Errorf("task %v not asr response", task.UUID)
			}

			// Fail if group exceed its duration.
			for _, g := range asrResponse.Groups {
				if g.TTSDuration == 0 {
					return errors.Errorf("task %v group %v not tts", task.UUID, g.UUID)
				}
				if g.TTSDuration > g.ASRDuration() {
					return errors.Errorf("task %v group %v tts %v exceed asr %v",
						task.UUID, g.UUID, g.TTSDuration, g.ASRDuration())
				}
				if g.TTS == "" {
					return errors.Errorf("task %v group %v not tts", task.UUID, g.UUID)
				}
				ttsFile := path.Join(conf.Pwd, aiDubbingWorkDir, g.TTS)
				if _, err := os.Stat(ttsFile); err != nil {
					return errors.Wrapf(err, "task %v group %v no tts %v file", task.UUID, g.UUID, g.TTS)
				}
			}

			// Download if file already exists.
			dubbingAudioFile := path.Join(dubbing.UUID, fmt.Sprintf("audio-%v.wav", task.UUID))
			absDubbingAudioFile := path.Join(conf.Pwd, aiDubbingWorkDir, dubbingAudioFile)

			exportFilename := path.Join(dubbing.UUID, fmt.Sprintf("audio-%v.mp4", task.UUID))
			absExportFile := path.Join(conf.Pwd, aiDubbingWorkDir, exportFilename)
			if _, err := os.Stat(absExportFile); err == nil {
				w.Header().Set("Content-Type", "video/mp4")
				http.ServeFile(w, r, absExportFile)
				logger.Tf(ctx, "srs dubbing download ok, dubbing=%v", dubbing.String())
				return nil
			}

			f, err := os.Create(absDubbingAudioFile)
			if err != nil {
				return errors.Wrapf(err, "create %v", absDubbingAudioFile)
			}
			defer f.Close()

			// 100KHZ, each frame is 10ms.
			buf := &audio.IntBuffer{Data: make([]int, 100000*48), Format: &audio.Format{SampleRate: 100000, NumChannels: 1}}
			enc := wav.NewEncoder(f, buf.Format.SampleRate, 16, buf.Format.NumChannels, 1)
			defer enc.Close()

			insertSilent := func(duration float64) error {
				if duration >= 0.01 {
					logger.Tf(ctx, "Write wav ok, silent=%v", duration)
					return enc.Write(&audio.IntBuffer{
						Data:   make([]int, int(100000*duration)),
						Format: &audio.Format{SampleRate: 100000, NumChannels: 1},
					})
				}
				return nil
			}

			for _, g := range asrResponse.Groups {
				if g.FirstSegment() == nil {
					continue
				}

				var gap float64
				if previous := asrResponse.PreviousGroup(g); previous != nil && g.FirstSegment() != nil {
					gap = g.FirstSegment().Start - previous.LastSegment().End
				}
				logger.Tf(ctx, "Dubbing artifact generate segment %v, time %v~%v",
					g.UUID, g.FirstSegment().Start, g.LastSegment().End)

				if err := insertSilent(gap); err != nil {
					return errors.Wrapf(err, "insert silent %v", gap)
				}

				var wavDuration float64
				if err := func() error {
					ttsFile := path.Join(conf.Pwd, aiDubbingWorkDir, g.TTS)
					wavFile := path.Join(conf.Pwd, aiDubbingWorkDir, dubbing.UUID, fmt.Sprintf("tts-%v.wav", g.UUID))
					logger.Tf(ctx, "Dubbing artifact convert tts %v to wav", ttsFile)
					if true {
						if err := exec.CommandContext(ctx, "ffmpeg",
							"-i", ttsFile,
							"-vn", "-c:a", "pcm_s16le", "-ac", "1", "-ar", "100000", "-ab", "300k",
							"-y", wavFile,
						).Run(); err != nil {
							return errors.Errorf("Error converting the file")
						}
					}

					wf, err := os.Open(wavFile)
					if err != nil {
						return errors.Wrapf(err, "open %v", wavFile)
					}
					defer wf.Close()

					dec := wav.NewDecoder(wf)
					bufWav, err := dec.FullPCMBuffer()
					if err != nil {
						return errors.Wrapf(err, "decode %v", wavFile)
					}
					if err = enc.Write(bufWav); err != nil {
						return errors.Wrapf(err, "write %v", wavFile)
					}

					wavDuration = float64(len(bufWav.Data)) / 100000.
					logger.Tf(ctx, "Dubbing artifact write wav ok, duration=%v, data=%.3f", g.TTSDuration, wavDuration)
					return nil
				}(); err != nil {
					return errors.Wrapf(err, "merge")
				}

				if err := insertSilent(g.LastSegment().End - g.FirstSegment().Start - wavDuration); err != nil {
					return errors.Wrapf(err, "insert silent %v", g.LastSegment().End-g.FirstSegment().Start-wavDuration)
				}
			}

			enc.Close()
			logger.Tf(ctx, "Dubbing artifact all segments are converted")

			// Merge original source file and dubbing audio file to a new file.
			if dubbing.SourceFormat != nil && dubbing.SourceFormat.HasVideo {
				absSourcePath := path.Join(conf.Pwd, aiDubbingWorkDir, dubbing.SourcePath)
				if err := exec.CommandContext(ctx, "ffmpeg",
					"-i", absSourcePath,
					"-i", absDubbingAudioFile,
					"-map", "0:v", "-map", "1:a",
					"-c:v", "copy",
					"-c:a", "aac", "-ac", "2", "-ar", "44100", "-ab", "120k",
					"-y", absExportFile,
				).Run(); err != nil {
					return errors.Errorf("Error converting the file")
				}
				logger.Tf(ctx, "Dubbing artifact convert %v and %v to aac %v ok",
					absSourcePath, absDubbingAudioFile, absExportFile)
			} else {
				if err := exec.CommandContext(ctx, "ffmpeg",
					"-i", absDubbingAudioFile,
					"-vn", "-c:a", "aac", "-ac", "2", "-ar", "44100", "-ab", "120k",
					"-y", absExportFile,
				).Run(); err != nil {
					return errors.Errorf("Error converting the file")
				}
				logger.Tf(ctx, "Dubbing artifact convert %v to aac %v ok", absDubbingAudioFile, absExportFile)
			}
			logger.Tf(ctx, "Dubbing artifact download AAC ok")

			w.Header().Set("Content-Type", "video/mp4")
			http.ServeFile(w, r, absExportFile)

			logger.Tf(ctx, "srs dubbing artifact download ok, dubbing=%v, export=%v", dubbing.String(), absExportFile)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/dubbing/task-start"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var dubbingUUID string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				UUID  *string `json:"uuid"`
			}{
				Token: &token, UUID: &dubbingUUID,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			dubbing := &SrsDubbingProject{UUID: dubbingUUID}
			if err := dubbing.Load(ctx); err != nil {
				return errors.Wrapf(err, "load dubbing project %v", dubbingUUID)
			}

			if err := dubbing.CheckAIConfig(ctx); err != nil {
				return errors.Wrapf(err, "check ai config")
			}

			task := dubbingServer.QueryTask(dubbing.TaskUUID)
			if task == nil {
				// TODO: Use transaction to ensure the task and project are consistent.
				task = NewSrsDubbingTask(func(task *SrsDubbingTask) {
					task.UUID = dubbing.UUID
					task.project = dubbing
					task.status = SrsDubbingTaskStatusInit
				})

				// Add task to server, memory object in this server.
				dubbingServer.AddTask(task)

				// Load the task from redis if exists.
				if dubbing.TaskUUID != "" {
					if err := task.Load(ctx); err != nil {
						logger.Wf(ctx, "ignore load dubbing task %v, err %+v", task.String(), err)
					}
				}

				// Save task to redis.
				if err := task.Save(ctx); err != nil {
					return errors.Wrapf(err, "save dubbing task %v", task.String())
				}

				// Set task to project and save project.
				dubbing.TaskUUID = task.UUID
				if err := dubbing.Save(ctx); err != nil {
					return errors.Wrapf(err, "save dubbing project %v", dubbing.String())
				}

				// Start the task for dubbing project.
				if err := task.Start(ctx); err != nil {
					return errors.Wrapf(err, "start dubbing task %v", task.String())
				}
			}

			// Wait for a very short time for starting task.
			select {
			case <-ctx.Done():
			case <-time.After(300 * time.Millisecond):
			}

			ohttp.WriteData(ctx, w, r, &struct {
				// Dubbing task UUID, equals to the dubbing project UUID.
				UUID string `json:"uuid"`
				// The task session UUID.
				SessionUUID string `json:"session"`
				// The task status, indicates the task is running or not.
				Status SrsDubbingTaskStatus `json:"status"`
			}{
				UUID: task.UUID, SessionUUID: task.SessionUUID, Status: task.status,
			})
			logger.Tf(ctx, "srs dubbing start task ok, dubbing=%v", dubbing.String())
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/dubbing/task-rephrase"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var dubbingUUID, taskUUID, groupUUID string
			if err := ParseBody(ctx, r.Body, &struct {
				Token     *string `json:"token"`
				UUID      *string `json:"uuid"`
				TaskUUID  *string `json:"task"`
				GroupUUID *string `json:"group"`
			}{
				Token: &token, UUID: &dubbingUUID, TaskUUID: &taskUUID,
				GroupUUID: &groupUUID,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			dubbing := NewSrsDubbingProject(func(dubbing *SrsDubbingProject) {
				dubbing.UUID = dubbingUUID
			})
			if err := dubbing.Load(ctx); err != nil {
				return errors.Wrapf(err, "load dubbing project %v", dubbingUUID)
			}

			// TODO: Allow task use different UUID from project.
			if dubbing.TaskUUID != taskUUID {
				return errors.Errorf("invalid task %v, should be %v", taskUUID, dubbing.TaskUUID)
			}

			// User should started the task.
			task := dubbingServer.QueryTask(taskUUID)
			if task == nil {
				return errors.Errorf("task %v not exists", taskUUID)
			}
			if task.AsrResponse == nil {
				return errors.Errorf("task %v not asr response", task.UUID)
			}

			// Query audio group.
			group := task.AsrResponse.QueryGroup(groupUUID)
			if group == nil {
				return errors.Errorf("group %v not exists", groupUUID)
			}
			previous := task.AsrResponse.PreviousGroup(group)

			// Setup the length of output. Note that for the first time, we translate again so we do not
			// use shorter length.
			maxOutputLength := len(group.Translated) / (1 << group.manuallyRephrasedCount)
			group.manuallyRephrasedCount++

			// Translate the group for the first manually rephrasing.
			if group.manuallyRephrasedCount == 1 {
				// Reset the translated and rephrased text.
				group.Translated, group.Rephrased = "", ""

				if dubbing.Translation.AIChatEnabled && group.Text() != "" {
					if err := group.Translate(ctx, dubbing.Translation, previous); err != nil {
						return errors.Wrapf(err, "translate group %v", group)
					}
				} else {
					logger.Tf(ctx, "Dubbing: Ignore translate for group %v", group)
				}

				if dubbing.TTS.AITTSEnabled && group.SourceTextForTTS() != "" {
					if err := group.GenerateTTS(ctx, dubbing.TTS, dubbingUUID); err != nil {
						return errors.Wrapf(err, "generate tts group %v", group)
					}
				} else {
					logger.Tf(ctx, "Dubbing: Ignore tts for group %v", group)
				}
			}

			// Rephrase and regenerate TTS of group.
			if group.manuallyRephrasedCount > 1 {
				if dubbing.Rephrase.AIChatEnabled && group.Translated != "" {
					if err := group.RephraseGroup(ctx, dubbing.Rephrase, previous, maxOutputLength); err != nil {
						return errors.Wrapf(err, "rephrase group %v", group)
					}
				} else {
					logger.Tf(ctx, "Dubbing: Ignore rephrase for group %v", group)
				}

				if dubbing.TTS.AITTSEnabled && group.SourceTextForTTS() != "" {
					if err := group.GenerateTTS(ctx, dubbing.TTS, dubbingUUID); err != nil {
						return errors.Wrapf(err, "generate tts group %v", group)
					}
				} else {
					logger.Tf(ctx, "Dubbing: Ignore tts for group %v", group)
				}
			}

			// Save the task.
			if err := task.Save(ctx); err != nil {
				return errors.Wrapf(err, "save")
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Status SrsDubbingTaskStatus `json:"status"`
				*SrsDubbingTask
			}{
				SrsDubbingTask: task, Status: task.status,
			})
			logger.Tf(ctx, "srs dubbing merge group ok, dubbing=%v, task=%v, group=%v",
				dubbing.String(), task.String(), group)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/dubbing/task-merge"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var dubbingUUID, taskUUID, groupUUID, direction string
			if err := ParseBody(ctx, r.Body, &struct {
				Token     *string `json:"token"`
				UUID      *string `json:"uuid"`
				TaskUUID  *string `json:"task"`
				GroupUUID *string `json:"group"`
				Direction *string `json:"direction"`
			}{
				Token: &token, UUID: &dubbingUUID, TaskUUID: &taskUUID,
				GroupUUID: &groupUUID, Direction: &direction,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			dubbing := NewSrsDubbingProject(func(dubbing *SrsDubbingProject) {
				dubbing.UUID = dubbingUUID
			})
			if err := dubbing.Load(ctx); err != nil {
				return errors.Wrapf(err, "load dubbing project %v", dubbingUUID)
			}

			// TODO: Allow task use different UUID from project.
			if dubbing.TaskUUID != taskUUID {
				return errors.Errorf("invalid task %v, should be %v", taskUUID, dubbing.TaskUUID)
			}

			// User should started the task.
			task := dubbingServer.QueryTask(taskUUID)
			if task == nil {
				return errors.Errorf("task %v not exists", taskUUID)
			}
			if task.AsrResponse == nil {
				return errors.Errorf("task %v not asr response", task.UUID)
			}

			// Query audio group.
			group := task.AsrResponse.QueryGroup(groupUUID)
			if group == nil {
				return errors.Errorf("group %v not exists", groupUUID)
			}

			if direction != "next" {
				return errors.Errorf("invalid direction %v", direction)
			}

			next := task.AsrResponse.NextGroup(group)
			if err := task.AsrResponse.MergeGroup(ctx, dubbing, next, group); err != nil {
				return errors.Wrapf(err, "merge group %v to %v", next, group)
			}

			// Save the task.
			if err := task.Save(ctx); err != nil {
				return errors.Wrapf(err, "save")
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Status SrsDubbingTaskStatus `json:"status"`
				*SrsDubbingTask
			}{
				SrsDubbingTask: task, Status: task.status,
			})
			logger.Tf(ctx, "srs dubbing merge group ok, dubbing=%v, task=%v, group=%v",
				dubbing.String(), task.String(), group)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/dubbing/task-query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var dubbingUUID, taskUUID string
			if err := ParseBody(ctx, r.Body, &struct {
				Token    *string `json:"token"`
				UUID     *string `json:"uuid"`
				TaskUUID *string `json:"task"`
			}{
				Token: &token, UUID: &dubbingUUID, TaskUUID: &taskUUID,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			dubbing := NewSrsDubbingProject(func(dubbing *SrsDubbingProject) {
				dubbing.UUID = dubbingUUID
			})
			if err := dubbing.Load(ctx); err != nil {
				return errors.Wrapf(err, "load dubbing project %v", dubbingUUID)
			}

			// TODO: Allow task use different UUID from project.
			if dubbing.TaskUUID != taskUUID {
				return errors.Errorf("invalid task %v, should be %v", taskUUID, dubbing.TaskUUID)
			}

			// User should started the task.
			task := dubbingServer.QueryTask(taskUUID)
			if task == nil {
				return errors.Errorf("task %v not exists", taskUUID)
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Status SrsDubbingTaskStatus `json:"status"`
				*SrsDubbingTask
			}{
				SrsDubbingTask: task, Status: task.status,
			})
			logger.Tf(ctx, "srs dubbing query task ok, dubbing=%v, task=%v", dubbing.String(), task.String())
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/dubbing/task-tts"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			q := r.URL.Query()

			token := q.Get("token")
			if token == "" {
				return errors.Errorf("empty token")
			}
			dubbingUUID := q.Get("uuid")
			if dubbingUUID == "" {
				return errors.Errorf("empty uuid")
			}
			groupUUID := q.Get("group")
			if groupUUID == "" {
				return errors.Errorf("empty group")
			}

			// Convert the token in query to header Bearer token.
			r.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, "", r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			dubbing := &SrsDubbingProject{UUID: dubbingUUID}
			if err := dubbing.Load(ctx); err != nil {
				return errors.Wrapf(err, "load dubbing project %v", dubbingUUID)
			}

			// User should started the task.
			task := dubbingServer.QueryTask(dubbing.TaskUUID)
			if task == nil {
				return errors.Errorf("task %v not exists", dubbing.TaskUUID)
			}
			if task.AsrResponse == nil {
				return errors.Errorf("task %v not asr response", task.UUID)
			}

			// Query audio group.
			group := task.AsrResponse.QueryGroup(groupUUID)
			if group == nil {
				return errors.Errorf("group %v not exists", groupUUID)
			}

			// The source file to play.
			filename := path.Join(conf.Pwd, aiDubbingWorkDir, group.TTS)
			ext := strings.Trim(path.Ext(filename), ".")
			contentType := fmt.Sprintf("audio/%v", ext)
			logger.Tf(ctx, "Serve example file=%v, ext=%v, contentType=%v", filename, ext, contentType)

			w.Header().Set("Content-Type", contentType)
			http.ServeFile(w, r, filename)

			ohttp.WriteData(ctx, w, r, &dubbing)
			logger.Tf(ctx, "srs dubbing get tts src ok, uuid=%v, group=%v, dubbing=%v",
				dubbingUUID, groupUUID, dubbing.String())
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/dubbing/source"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			type DubbingTempFile struct {
				// The file name.
				Name string `json:"name"`
				// The file path.
				Path string `json:"path"`
				// The file size in bytes.
				Size int64 `json:"size"`
				// The UUID for file.
				UUID string `json:"uuid"`
				// The target file name.
				Target string `json:"target"`
				// The source type.
				Type FFprobeSourceType `json:"type"`
			}

			var token string
			var files []*DubbingTempFile
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string             `json:"token"`
				Files *[]*DubbingTempFile `json:"files"`
			}{
				Token: &token, Files: &files,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if len(files) == 0 {
				return errors.New("no files")
			}

			// Always cleanup the files in upload.
			var tempFiles []string
			for _, f := range files {
				if f.Type != FFprobeSourceTypeStream {
					tempFiles = append(tempFiles, f.Target)
				}
			}
			defer func() {
				for _, tempFile := range tempFiles {
					if _, err := os.Stat(tempFile); err == nil {
						os.Remove(tempFile)
						logger.Tf(ctx, "Dubbing: Cleanup %v", tempFile)
					}
				}
			}()

			// Check files.
			for _, f := range files {
				if f.Target == "" {
					return errors.New("no target")
				}
				if f.Type != FFprobeSourceTypeStream {
					if _, err := os.Stat(f.Target); err != nil {
						return errors.Wrapf(err, "no file %v", f.Target)
					}
					if !strings.HasPrefix(f.Target, dirUploadPath) {
						return errors.Errorf("invalid target %v", f.Target)
					}
				}
			}

			// Parsed source files.
			var parsedFiles []*FFprobeSource

			// Parse file information and move file from dirUploadPath to dirDubbingPath.
			for _, file := range files {
				// Probe file information.
				toCtx, toCancelFunc := context.WithTimeout(ctx, 15*time.Second)
				defer toCancelFunc()

				args := []string{
					"-show_error", "-show_private_data", "-v", "quiet", "-find_stream_info", "-print_format", "json",
					"-show_format", "-show_streams",
				}
				// For RTSP stream source, always use TCP transport.
				if strings.HasPrefix(file.Target, "rtsp://") {
					args = append(args, "-rtsp_transport", "tcp")
				}
				// Rebuild the stream url, because it may contain special characters.
				if strings.Contains(file.Target, "://") {
					if u, err := RebuildStreamURL(file.Target); err != nil {
						return errors.Wrapf(err, "rebuild %v", file.Target)
					} else {
						args = append(args, "-i", u.String())
					}
				} else {
					args = append(args, "-i", file.Target)
				}

				// TODO: FIXME: Use FFprobeFileFormat.
				stdout, err := exec.CommandContext(toCtx, "ffprobe", args...).Output()
				if err != nil {
					return errors.Wrapf(err, "probe %v with ffprobe %v", file.Target, args)
				}

				format := struct {
					Format FFprobeFormat `json:"format"`
				}{}
				if err = json.Unmarshal([]byte(stdout), &format); err != nil {
					return errors.Wrapf(err, "parse format %v", stdout)
				}

				videos := struct {
					Streams []FFprobeVideo `json:"streams"`
				}{}
				if err = json.Unmarshal([]byte(stdout), &videos); err != nil {
					return errors.Wrapf(err, "parse video streams %v", stdout)
				}
				var matchVideo *FFprobeVideo
				for _, video := range videos.Streams {
					if video.CodecType == "video" {
						matchVideo = &video
						format.Format.HasVideo = true
						break
					}
				}

				audios := struct {
					Streams []FFprobeAudio `json:"streams"`
				}{}
				if err = json.Unmarshal([]byte(stdout), &audios); err != nil {
					return errors.Wrapf(err, "parse audio streams %v", stdout)
				}
				var matchAudio *FFprobeAudio
				for _, audio := range audios.Streams {
					if audio.CodecType == "audio" {
						matchAudio = &audio
						format.Format.HasAudio = true
						break
					}
				}

				parsedFile := &FFprobeSource{
					Name: file.Name, Path: file.Path, Size: uint64(file.Size), UUID: file.UUID,
					Target: file.Target, Type: file.Type,
					Format: &format.Format, Video: matchVideo, Audio: matchAudio,
				}
				if file.Type != FFprobeSourceTypeStream {
					parsedFile.Target = path.Join(dirDubbingPath, fmt.Sprintf("%v%v", file.UUID, path.Ext(file.Target)))
					if err = os.Rename(file.Target, parsedFile.Target); err != nil {
						return errors.Wrapf(err, "rename %v to %v", file.Target, parsedFile.Target)
					}
				}

				parsedFiles = append(parsedFiles, parsedFile)
				logger.Tf(ctx, "Dubbing: Process file %v", parsedFile.String())
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Files []*FFprobeSource `json:"files"`
			}{
				Files: parsedFiles,
			})
			logger.Tf(ctx, "Dubbing: Update dubbing ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	return nil
}

type SrsDubbingServer struct {
	// All dubbing tasks.
	tasks []*SrsDubbingTask
	// The lock to protect the tasks.
	lock sync.Mutex
}

func NewDubbingServer() *SrsDubbingServer {
	return &SrsDubbingServer{}
}

func (v *SrsDubbingServer) Close() error {
	// Copy tasks.
	var tasks []*SrsDubbingTask
	func() {
		v.lock.Lock()
		defer v.lock.Unlock()

		tasks = v.tasks
		v.tasks = nil
	}()

	// Cleanup all tasks.
	for _, task := range tasks {
		_ = task.Close()
	}
	return nil
}

func (v *SrsDubbingServer) AddTask(task *SrsDubbingTask) {
	v.lock.Lock()
	defer v.lock.Unlock()

	v.tasks = append(v.tasks, task)
}

func (v *SrsDubbingServer) RemoveTask(uuid string) {
	v.lock.Lock()
	defer v.lock.Unlock()

	for i, task := range v.tasks {
		if task.UUID == uuid {
			v.tasks = append(v.tasks[:i], v.tasks[i+1:]...)
			break
		}
	}
}

func (v *SrsDubbingServer) QueryTask(uuid string) *SrsDubbingTask {
	v.lock.Lock()
	defer v.lock.Unlock()

	for _, task := range v.tasks {
		if task.UUID == uuid {
			return task
		}
	}
	return nil
}

type AITime time.Time

func (v *AITime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*v).Format(time.RFC3339))
}

func (v *AITime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return errors.Wrapf(err, "unmarshal")
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return errors.Wrapf(err, "parse")
	}

	*v = AITime(t)
	return nil
}

type AudioSegment struct {
	// The id of whisper response.
	ID int `json:"id"`
	// The seek of whisper response.
	Seek int `json:"seek"`
	// The start time of whisper response, in seconds.
	Start float64 `json:"start"`
	// The end time of whisper response, in seconds.
	End float64 `json:"end"`
	// The text of whisper response.
	Text string `json:"text"`
	// The tokens of whisper response.
	Tokens []int `json:"tokens"`
	// The temperature of whisper response.
	Temperature float64 `json:"temperature"`
	// The average logprob of whisper response.
	AvgLogprob float64 `json:"avg_logprob"`
	// The compression ratio of whisper response.
	CompressionRatio float64 `json:"compression_ratio"`
	// The no speech prob of whisper response.
	NoSpeechProb float64 `json:"no_speech_prob"`
	// The transient of whisper response.
	Transient bool `json:"transient"`

	// Original starttime, allow recover and skip exists segments.
	OriginalStart float64 `json:"original_start"`
	// The UUID generated by system.
	UUID string `json:"uuid"`
}

func (v *AudioSegment) Duration() time.Duration {
	return time.Duration((v.End - v.Start) * float64(time.Second))
}

func (v *AudioSegment) Diff(to *AudioSegment) time.Duration {
	if v.Start < to.Start {
		return time.Duration((to.Start - v.End) * float64(time.Second))
	} else {
		return time.Duration((v.Start - to.End) * float64(time.Second))
	}
}

type AudioGroup struct {
	// The group UUID.
	UUID string `json:"uuid"`
	// The id of group index.
	ID int `json:"id"`
	// The audio segments in the same group.
	Segments []*AudioSegment `json:"segments"`

	// Translated text.
	Translated string `json:"translated"`
	// Translate time.
	TranslatedAt AITime `json:"translated_at"`

	// Rephrased text.
	Rephrased string `json:"rephrased"`
	// Rephrased time.
	RephrasedAt AITime `json:"rephrased_at"`

	// TTS filename, without the main dir.
	TTS string `json:"tts"`
	// Convert TTS time.
	TTSAt AITime `json:"tts_at"`
	// The TTS audio duration, in seconds.
	TTSDuration float64 `json:"tts_duration"`

	// How many times the group rephrased manually.
	manuallyRephrasedCount uint32
}

func (v *AudioGroup) Translate(ctx context.Context, trans *SrsAssistant, previous *AudioGroup) error {
	systemPrompt := fmt.Sprintf("%v. Never answer questions but directly translate text.",
		trans.AIChatPrompt)
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
	}

	if previous != nil && previous.Translated != "" && previous.Text() != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleUser, Content: previous.Text(),
		})
		messages = append(messages, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleAssistant, Content: previous.Translated,
		})
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleUser, Content: v.Text(),
	})

	// Initialize the AI services.
	aiConfig := openai.DefaultConfig(trans.AISecretKey)
	aiConfig.OrgID = trans.AIOrganization
	aiConfig.BaseURL = trans.AIBaseURL

	client := openai.NewClientWithConfig(aiConfig)
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    trans.AIChatModel,
		Messages: messages,
	})
	if err != nil {
		return errors.Wrapf(err, "translate")
	}

	v.Translated = resp.Choices[0].Message.Content
	v.TranslatedAt = AITime(time.Now())
	logger.Tf(ctx, "dubbing translate group %v to %v, model=%v, prompt is %v",
		v.Text(), v.Translated, trans.AIChatModel, systemPrompt)
	return nil
}

func (v *AudioGroup) GenerateTTS(ctx context.Context, tts *SrsAssistant, projectUUID string) error {
	if err := func() error {
		// Initialize the AI services.
		aiConfig := openai.DefaultConfig(tts.AISecretKey)
		aiConfig.OrgID = tts.AIOrganization
		aiConfig.BaseURL = tts.AIBaseURL

		client := openai.NewClientWithConfig(aiConfig)
		resp, err := client.CreateSpeech(ctx, openai.CreateSpeechRequest{
			Model:          openai.TTSModel1,
			Input:          v.SourceTextForTTS(),
			Voice:          openai.VoiceNova,
			ResponseFormat: openai.SpeechResponseFormatAac,
		})
		if err != nil {
			return errors.Wrapf(err, "create speech")
		}
		defer resp.Close()

		ttsFilename := path.Join(projectUUID, fmt.Sprintf("tts-%v.aac", v.UUID))
		ttsFile := path.Join(conf.Pwd, aiDubbingWorkDir, ttsFilename)
		out, err := os.Create(ttsFile)
		if err != nil {
			return errors.Errorf("Unable to create the file %v for writing", ttsFile)
		}
		defer out.Close()

		if _, err = io.Copy(out, resp); err != nil {
			return errors.Errorf("Error writing the file")
		}

		v.TTS = ttsFilename
		v.TTSAt = AITime(time.Now())
		logger.Tf(ctx, "dubbing generate TTS for group %v to %v", v.SourceTextForTTS(), v.TTS)
		return nil
	}(); err != nil {
		return errors.Wrapf(err, "generate tts")
	}

	// Detect the TTS duration.
	if format, _, _, err := FFprobeFileFormat(ctx, path.Join(conf.Pwd, aiDubbingWorkDir, v.TTS)); err != nil {
		return errors.Wrapf(err, "ffprobe %v", v.TTS)
	} else {
		v.TTSDuration = format.Duration
	}

	return nil
}

func (v *AudioGroup) RephraseGroup(ctx context.Context, rephrase *SrsAssistant, previous *AudioGroup, maxOutputLength int) error {
	systemPrompt := fmt.Sprintf("%v. Limit output length to %v. Never answer questions but directly rephrase text.",
		rephrase.AIChatPrompt, maxOutputLength)
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
	}

	if previous != nil && previous.Translated != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleUser, Content: previous.Translated,
		})
		messages = append(messages, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleAssistant, Content: previous.SourceTextForTTS(),
		})
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleUser, Content: v.Translated,
	})

	// Initialize the AI services.
	aiConfig := openai.DefaultConfig(rephrase.AISecretKey)
	aiConfig.OrgID = rephrase.AIOrganization
	aiConfig.BaseURL = rephrase.AIBaseURL

	client := openai.NewClientWithConfig(aiConfig)
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    rephrase.AIChatModel,
		Messages: messages,
	})
	if err != nil {
		return errors.Wrapf(err, "translate")
	}

	v.Rephrased = resp.Choices[0].Message.Content
	v.RephrasedAt = AITime(time.Now())
	logger.Tf(ctx, "dubbing rephrase group %vB to %vB, %v to %v, rephrased=%v, prompt is %v",
		len(v.Translated), len(v.Rephrased), v.Translated, v.Rephrased, v.manuallyRephrasedCount, systemPrompt)

	return nil
}

func (v *AudioGroup) SourceTextForTTS() string {
	if v.Rephrased != "" {
		return v.Rephrased
	}
	return v.Translated
}

func (v *AudioGroup) ASRDuration() float64 {
	if len(v.Segments) == 0 {
		return 0
	}

	return v.LastSegment().End - v.FirstSegment().Start
}

func (v *AudioGroup) ASRDuration2() time.Duration {
	return time.Duration(v.ASRDuration() * float64(time.Second))
}

func (v *AudioGroup) FirstSegment() *AudioSegment {
	if len(v.Segments) > 0 {
		return v.Segments[0]
	}
	return nil
}

func (v *AudioGroup) LastSegment() *AudioSegment {
	if len(v.Segments) > 0 {
		return v.Segments[len(v.Segments)-1]
	}
	return nil
}

func (v *AudioGroup) RemoveSegment(segment *AudioSegment) *AudioSegment {
	for i, s := range v.Segments {
		if s.UUID == segment.UUID {
			v.Segments = append(v.Segments[:i], v.Segments[i+1:]...)
			return s
		}
	}
	return nil
}

func (v *AudioGroup) Text() string {
	var text string
	for _, s := range v.Segments {
		text += s.Text
	}
	return text
}

type AudioResponse struct {
	Task     string        `json:"task"`
	Language string        `json:"language"`
	Duration float64       `json:"duration"`
	Groups   []*AudioGroup `json:"groups"`
	Text     string        `json:"text"`
}

func NewAudioResponse() *AudioResponse {
	return &AudioResponse{}
}

func (v *AudioResponse) MergeGroup(ctx context.Context, dubbing *SrsDubbingProject, from, to *AudioGroup) error {
	if from == nil || to == nil {
		return errors.Errorf("invalid merge group from=%v, to=%v", from, to)
	}

	to.Translated += from.Translated
	to.Rephrased += from.Rephrased
	to.Segments = append(to.Segments, from.Segments...)

	// Remove the source from group.
	for i, g := range v.Groups {
		if g == from {
			v.Groups = append(v.Groups[:i], v.Groups[i+1:]...)
			break
		}
	}

	// Regenerate the group id.
	for i, g := range v.Groups {
		g.ID = i
	}

	// Rephrase the merged group.
	if dubbing.Rephrase.AIChatEnabled && to.Translated != "" {
		if err := to.RephraseGroup(ctx, dubbing.Rephrase, v.PreviousGroup(to), len(to.Translated)/3); err != nil {
			return errors.Wrapf(err, "rephrase group")
		}
	} else {
		logger.Tf(ctx, "Dubbing: Ignore rephrase for group %v", to)
	}

	// Regenerate TTS for the group.
	if dubbing.TTS.AITTSEnabled && to.SourceTextForTTS() != "" {
		if err := to.GenerateTTS(ctx, dubbing.TTS, dubbing.UUID); err != nil {
			return errors.Wrapf(err, "generate tts")
		}
	} else {
		logger.Tf(ctx, "Dubbing: Ignore tts for group %v", to)
	}

	return nil
}

func (v *AudioResponse) PreviousGroup(g *AudioGroup) *AudioGroup {
	for i, group := range v.Groups {
		if group == g && i > 0 {
			return v.Groups[i-1]
		}
	}
	return nil
}

func (v *AudioResponse) NextGroup(g *AudioGroup) *AudioGroup {
	for i, group := range v.Groups {
		if group == g && i < len(v.Groups)-1 {
			return v.Groups[i+1]
		}
	}
	return nil
}

func (v *AudioResponse) QueryGroup(uuid string) *AudioGroup {
	for _, group := range v.Groups {
		if group.UUID == uuid {
			return group
		}
	}
	return nil
}

func (v *AudioResponse) MatchGroups(starttime float64, max int) []*AudioGroup {
	var matched []*AudioGroup
	for _, g := range v.Groups {
		if first := g.FirstSegment(); first != nil && starttime <= first.Start {
			matched = append(matched, g)
		}

		if len(matched) >= max {
			break
		}
	}
	return matched
}

func (v *AudioResponse) AppendSegment(resp openai.AudioResponse, starttime float64) {
	v.Task = resp.Task
	v.Language = resp.Language
	v.Duration += resp.Duration
	if len(v.Text) < 1024 {
		v.Text += " " + resp.Text
	} else if !strings.HasSuffix(v.Text, "...") {
		v.Text += "..."
	}

	for _, s := range resp.Segments {
		v.Groups = append(v.Groups, &AudioGroup{
			// The group index.
			ID: len(v.Groups),
			// The uuid of group.
			UUID: uuid.NewString(),
			// The segments in the group.
			Segments: []*AudioSegment{
				&AudioSegment{
					// To identify the segments.
					OriginalStart: starttime,
					// ASR Segment.
					ID:               s.ID,
					Seek:             s.Seek,
					Start:            starttime + s.Start,
					End:              starttime + s.End,
					Text:             s.Text,
					Tokens:           s.Tokens,
					Temperature:      s.Temperature,
					AvgLogprob:       s.AvgLogprob,
					CompressionRatio: s.CompressionRatio,
					NoSpeechProb:     s.NoSpeechProb,
					Transient:        s.Transient,
					// UUID.
					UUID: uuid.NewString(),
				},
			},
		})
	}
}

type SrsDubbingTaskStatus string

const (
	// SrsDubbingTaskStatusInit is the init status.
	SrsDubbingTaskStatusInit SrsDubbingTaskStatus = "init"
	// SrsDubbingTaskStatusRunning means User start the task.
	SrsDubbingTaskStatusRunning SrsDubbingTaskStatus = "running"
	// SrsDubbingTaskStatusExtractAudio means the task is extract audio from source input.
	SrsDubbingTaskStatusExtractAudio SrsDubbingTaskStatus = "extract"
	// SrsDubbingTaskStatusAsrGenerating means the task is doing ASR, convert audio to text.
	SrsDubbingTaskStatusAsrGenerating SrsDubbingTaskStatus = "asr"
	// SrsDubbingTaskStatusTranslating means the task is translating the ASR text.
	SrsDubbingTaskStatusTranslating SrsDubbingTaskStatus = "translating"
	// SrsDubbingTaskStatusTTSGenerating means the task is generating TTS audio, the first pass.
	SrsDubbingTaskStatusTTSGenerating SrsDubbingTaskStatus = "tts"
	// SrsDubbingTaskStatusRephrasing means the task is rephrasing the translated text.
	SrsDubbingTaskStatusRephrasing SrsDubbingTaskStatus = "rephrasing"
	// SrsDubbingTaskStatusMerging means the task is merging small groups.
	SrsDubbingTaskStatusMerging SrsDubbingTaskStatus = "merging"
	// SrsDubbingTaskStatusDone means the task is done.
	SrsDubbingTaskStatusDone SrsDubbingTaskStatus = "done"
)

type SrsDubbingTask struct {
	// Dubbing task UUID, equals to the dubbing project UUID.
	UUID string `json:"uuid"`
	// The task session UUID.
	SessionUUID string `json:"session"`
	// The ASR input audio file.
	AsrInputAudio string `json:"asr_input_audio"`
	// The ASR input bitrate in bps.
	AsrInputBitrate int `json:"asr_input_bitrate"`
	// The ASR response object, text and segemnts.
	AsrResponse *AudioResponse `json:"asr_response"`

	// Whether running the task.
	running bool
	// The task status, indicates the task is running or not.
	status SrsDubbingTaskStatus
	// The owner dubbing project.
	project *SrsDubbingProject
	// The context to control the task.
	ctx context.Context
	// The cancel function to stop the task.
	cancel context.CancelFunc
	// The cleanup wait group.
	wg sync.WaitGroup
	// The lock to protect the task.
	lock sync.Mutex
}

func NewSrsDubbingTask(opts ...func(task *SrsDubbingTask)) *SrsDubbingTask {
	v := &SrsDubbingTask{
		// The uuid, normally equals to the project.
		UUID: uuid.NewString(),
		// The session uuid, created for eachtime.
		SessionUUID: uuid.NewString(),
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

func (v *SrsDubbingTask) Close() error {
	if v.cancel != nil {
		v.cancel()
	}

	v.wg.Wait()

	// TODO: FIXME: Cleanup project files.

	return nil
}

func (v *SrsDubbingTask) String() string {
	return fmt.Sprintf("uuid=%v, asr=%v", v.UUID, v.AsrInputAudio)
}

func (v *SrsDubbingTask) UpdateProject(dubbing *SrsDubbingProject) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.project = dubbing
}

func (v *SrsDubbingTask) Load(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	dubbingUUID := v.UUID

	if r0, err := rdb.HGet(ctx, SRS_DUBBING_TASKS, dubbingUUID).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_DUBBING_TASKS, dubbingUUID)
	} else if r0 == "" {
		return errors.Errorf("dubbing project %v not exists", dubbingUUID)
	} else if err = json.Unmarshal([]byte(r0), v); err != nil {
		return errors.Wrapf(err, "unmarshal %v %v", dubbingUUID, r0)
	}

	return nil
}

func (v *SrsDubbingTask) Save(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	dubbing := v

	if b, err := json.Marshal(dubbing); err != nil {
		return errors.Wrapf(err, "marshal dubbing task")
	} else if err := rdb.HSet(ctx, SRS_DUBBING_TASKS, dubbing.UUID, string(b)).Err(); err != nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_DUBBING_TASKS, dubbing.UUID, string(b))
	}

	return nil
}

func (v *SrsDubbingTask) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	v.ctx, v.cancel = ctx, cancel
	logger.Tf(ctx, "start dubbing worker for task %v", v.UUID)

	// Always set the constant fields.
	v.AsrInputBitrate = 30 * 1000 // in N * kbps.

	// Extract audio from source file.
	extractAudio := func() error {
		if v.AsrInputAudio != "" {
			if _, err := os.Stat(path.Join(conf.Pwd, aiDubbingWorkDir, v.AsrInputAudio)); err != nil {
				v.AsrResponse = nil
				logger.Tf(ctx, "Dubbing: reset the asr response for no input %v", v.AsrInputAudio)
			} else {
				logger.Tf(ctx, "Dubbing: ignore for asr ready, task=%v, asr=%v", v.UUID, v.AsrInputAudio)
				return nil
			}
		}

		// We are extracting the audio from source file.
		v.status = SrsDubbingTaskStatusExtractAudio

		// Regenerate the ASR response if input changed.
		asrInputAudio := path.Join(v.project.UUID, fmt.Sprintf("asr-%v.m4a", v.SessionUUID))
		absSourcePath := path.Join(conf.Pwd, aiDubbingWorkDir, v.project.SourcePath)
		if err := exec.CommandContext(ctx, "ffmpeg",
			"-i", absSourcePath,
			"-vn", "-c:a", "aac", "-ac", "1", "-ar", "16000", "-ab", fmt.Sprintf("%v", v.AsrInputBitrate),
			path.Join(conf.Pwd, aiDubbingWorkDir, asrInputAudio),
		).Run(); err != nil {
			return errors.Errorf("Error converting the file")
		}

		v.AsrInputAudio = asrInputAudio
		if err := v.Save(ctx); err != nil {
			return errors.Wrapf(err, "save task for asr audio, task is %v", v.String())
		}

		logger.Tf(ctx, "Dubbing: Convert task=%v, source=%v to m4a asr=%v ok", v.UUID, absSourcePath, asrInputAudio)
		return nil
	}

	// Generate ASR response by AI service.
	generateAsrResponse := func() error {
		if v.project == nil || v.project.SourceFormat == nil {
			return errors.Errorf("invalid project or source format")
		}
		if v.project.SourceFormat.Bitrate <= 0 {
			return errors.Errorf("invalid source bitrate %v", v.project.SourceFormat.Bitrate)
		}
		if v.project.SourceFormat.Duration <= 0 {
			return errors.Errorf("invalid source duration %v", v.project.SourceFormat.Duration)
		}

		// We are generating the ASR response.
		v.status = SrsDubbingTaskStatusAsrGenerating

		// Ignore if disabled.
		if !v.project.ASR.AIASREnabled {
			logger.Tf(ctx, "Dubbing: ignore for ASR disabled, task=%v, project=%v", v.UUID, v.project.UUID)
			return nil
		}

		// Whether force to generate ASR response.
		if alwaysForceRegenerateASRResponse {
			v.AsrResponse = nil
		}

		// Whether exists command file to regenerate ASR.
		regenerateASR := path.Join(conf.Pwd, aiDubbingWorkDir, v.project.UUID, "regenerate.txt")
		if _, err := os.Stat(regenerateASR); err == nil {
			os.Remove(regenerateASR)
			v.AsrResponse = nil
		}

		if v.AsrResponse == nil {
			v.AsrResponse = NewAudioResponse()
		}
		bitrate, duration := v.AsrInputBitrate, v.project.SourceFormat.Duration

		// Split the audio to segments, because each ASR is limited to 25MB by OpenAI,
		// see https://platform.openai.com/docs/guides/speech-to-text
		limitDuration := int(25*1024*1024*8/float64(bitrate)) / 10
		for starttime := float64(0); starttime < duration; starttime += float64(limitDuration) {
			// For debugging, only the first segment.
			if starttime > 0 && onlyRegenerateFirstSegment {
				break
			}

			// If there are some matched groups about start time, do not generate the ASR for it.
			ignoreIfMatched := 3
			if matches := v.AsrResponse.MatchGroups(starttime, ignoreIfMatched*2); len(matches) > ignoreIfMatched {
				continue
			}

			if err := func() error {
				tmpAsrInputAudio := path.Join(conf.Pwd, aiDubbingWorkDir, v.project.UUID, fmt.Sprintf("%v-%v.m4a", v.SessionUUID, starttime))
				defer os.Remove(tmpAsrInputAudio)

				absAsrInputAudio := path.Join(conf.Pwd, aiDubbingWorkDir, v.AsrInputAudio)
				if err := exec.CommandContext(ctx, "ffmpeg",
					"-i", absAsrInputAudio,
					"-ss", fmt.Sprintf("%v", starttime), "-t", fmt.Sprintf("%v", limitDuration),
					"-c", "copy", "-y", tmpAsrInputAudio,
				).Run(); err != nil {
					return errors.Errorf("Error converting the file %v to %v", absAsrInputAudio, tmpAsrInputAudio)
				}
				logger.Tf(ctx, "Convert %v to segment %v ok, starttime=%v", absAsrInputAudio, tmpAsrInputAudio, starttime)

				// Initialize the AI services.
				aiConfig := openai.DefaultConfig(v.project.ASR.AISecretKey)
				aiConfig.OrgID = v.project.ASR.AIOrganization
				aiConfig.BaseURL = v.project.ASR.AIBaseURL

				// Do ASR, convert to text.
				client := openai.NewClientWithConfig(aiConfig)
				resp, err := client.CreateTranscription(
					ctx,
					openai.AudioRequest{
						Model:    openai.Whisper1,
						FilePath: tmpAsrInputAudio,
						Format:   openai.AudioResponseFormatVerboseJSON,
						Language: v.project.ASR.AIASRLanguage,
					},
				)
				if err != nil {
					return errors.Wrapf(err, "transcription")
				}
				logger.Tf(ctx, "ASR ok, project=%v, resp is <%v>B, group=%v",
					v.project.UUID, len(resp.Text), len(v.AsrResponse.Groups))

				// Append the segment to ASR output object.
				v.AsrResponse.AppendSegment(resp, starttime)
				logger.Tf(ctx, "Save ASR output ok")

				return nil
			}(); err != nil {
				return errors.Wrapf(err, "split starttime=%v, duration=%v", starttime, limitDuration)
			}

			if err := v.Save(ctx); err != nil {
				return errors.Wrapf(err, "save task for asr response, task is %v", v.String())
			}
		}

		return nil
	}

	// Always try to merge very small segments to previous groups.
	mergeSegmentsBetweenGroups := func() error {
		if v.AsrResponse == nil {
			logger.Tf(ctx, "Dubbing: ignore merge for no ASR response, task=%v, project=%v", v.UUID, v.project.UUID)
			return nil
		}

		isEnglish := func(s string) bool {
			for _, r := range s {
				if r > unicode.MaxASCII {
					return false
				}
			}
			return true
		}

		// Merge the segments between groups.
		var anythingChanged bool
		for gIndex, g := range v.AsrResponse.Groups {
			// Find next group, ignore if no next.
			var nextGroup *AudioGroup
			if gIndex+1 < len(v.AsrResponse.Groups) {
				nextGroup = v.AsrResponse.Groups[gIndex+1]
			}
			if nextGroup == nil {
				continue
			}

			// Ignore if already translated.
			if g.Translated != "" || g.Rephrased != "" || g.TTS != "" {
				continue
			}

			// Find first segment of next group, ignore if no segment.
			lastSegment, nextFirstSegment := g.LastSegment(), nextGroup.FirstSegment()
			if lastSegment == nil || nextGroup == nil || nextFirstSegment == nil {
				continue
			}

			// If last segment is very close to the next group, and first segment of next group
			// has very few words, merge next segment to this group.
			var nextFewWords bool
			if nextFirstSegment.Diff(lastSegment) < MergeSegmentBetweenGroups {
				nextText := strings.TrimSpace(nextFirstSegment.Text)

				// Whether the next segment is very few words.
				if isEnglish(nextText) {
					nextFewWords = strings.Count(nextText, " ") < MergeSegmentForSmallWords
				} else {
					nextFewWords = utf8.RuneCount([]byte(nextText)) < MergeSegmentForSmallWords
				}
			}

			if !nextFewWords {
				continue
			}

			// Merge first segment of next group to this group.
			if s := nextGroup.RemoveSegment(nextFirstSegment); s != nil {
				anythingChanged = true
				g.Segments = append(g.Segments, s)
			}
		}

		if !anythingChanged {
			return nil
		}

		// Eliminate the empty group.
		var groups []*AudioGroup
		for _, g := range v.AsrResponse.Groups {
			if len(g.Segments) == 0 {
				continue
			}

			g.ID = len(groups)
			groups = append(groups, g)
		}

		logger.Tf(ctx, "Merge groups from %v to %v", len(v.AsrResponse.Groups), len(groups))
		v.AsrResponse.Groups = groups

		if err := v.Save(ctx); err != nil {
			return errors.Wrapf(err, "save task for asr response, task is %v", v.String())
		}

		return nil
	}

	// Translate all groups to target language.
	translateGroups := func() error {
		v.status = SrsDubbingTaskStatusTranslating

		if v.AsrResponse == nil {
			logger.Tf(ctx, "Dubbing: ignore translate for no ASR response, task=%v, project=%v", v.UUID, v.project.UUID)
			return nil
		}

		// Ignore if disabled.
		if !v.project.Translation.AIChatEnabled {
			logger.Tf(ctx, "Dubbing: ignore for translate disabled, task=%v, project=%v", v.UUID, v.project.UUID)
			return nil
		}

		for _, g := range v.AsrResponse.Groups {
			if g.Translated != "" {
				continue
			}

			previous := v.AsrResponse.PreviousGroup(g)
			if err := g.Translate(ctx, v.project.Translation, previous); err != nil {
				return errors.Wrapf(err, "translate group %v", g)
			}

			if err := v.Save(ctx); err != nil {
				return errors.Wrapf(err, "save task for translation, task is %v", v.String())
			}
		}

		logger.Tf(ctx, "Translate all %v groups to target language", len(v.AsrResponse.Groups))
		return nil
	}

	// Generate TTS audio for each group.
	generateTTS := func() error {
		v.status = SrsDubbingTaskStatusTTSGenerating

		if v.AsrResponse == nil {
			logger.Tf(ctx, "Dubbing: ignore TTS for no ASR response, task=%v, project=%v", v.UUID, v.project.UUID)
			return nil
		}

		// Ignore if disabled.
		if !v.project.TTS.AITTSEnabled {
			logger.Tf(ctx, "Dubbing: ignore for TTS disabled, task=%v, project=%v", v.UUID, v.project.UUID)
			return nil
		}

		for _, g := range v.AsrResponse.Groups {
			if g.TTS != "" {
				continue
			}
			if g.SourceTextForTTS() == "" {
				continue
			}

			if err := g.GenerateTTS(ctx, v.project.TTS, v.project.UUID); err != nil {
				return errors.Wrapf(err, "generate tts for group %v", g)
			}

			if err := v.Save(ctx); err != nil {
				return errors.Wrapf(err, "save task for tts, task is %v", v.String())
			}
		}

		logger.Tf(ctx, "Generate TTS for all %v groups", len(v.AsrResponse.Groups))
		return nil
	}

	rephraseAsShorterTranslations := func() error {
		v.status = SrsDubbingTaskStatusRephrasing

		if v.AsrResponse == nil {
			logger.Tf(ctx, "Dubbing: ignore rephrase for no ASR response, task=%v, project=%v", v.UUID, v.project.UUID)
			return nil
		}

		// Ignore if disabled.
		if !v.project.Rephrase.AIChatEnabled {
			logger.Tf(ctx, "Dubbing: ignore for rephrase disabled, task=%v, project=%v", v.UUID, v.project.UUID)
			return nil
		}

		for _, g := range v.AsrResponse.Groups {
			if g.Rephrased != "" && !alwaysRephraseTranslations {
				continue
			}

			// If TTS duration not exceed the group duration, or already rephrased, use translated text directly.
			if g.TTSDuration <= 0 || g.TTSDuration <= g.ASRDuration() {
				continue
			}

			previous := v.AsrResponse.PreviousGroup(g)
			if err := g.RephraseGroup(ctx, v.project.Rephrase, previous, len(g.Translated)/2); err != nil {
				return errors.Wrapf(err, "rephrase group %v", g)
			}

			if err := g.GenerateTTS(ctx, v.project.TTS, v.project.UUID); err != nil {
				return errors.Wrapf(err, "generate tts for group %v", g)
			}

			if err := v.Save(ctx); err != nil {
				return errors.Wrapf(err, "save task for rephrase, task is %v", v.String())
			}
		}

		logger.Tf(ctx, "Rephrase all %v groups to better text", len(v.AsrResponse.Groups))
		return nil
	}

	// The main cycle for worker.
	cycle := func() error {
		// Now, task is running.
		v.status = SrsDubbingTaskStatusRunning

		// Extract audio from source file.
		if err := extractAudio(); err != nil {
			return errors.Wrapf(err, "extract audio")
		}

		// Generate ASR response by AI service.
		if err := generateAsrResponse(); err != nil {
			return errors.Wrapf(err, "generate asr response")
		}

		// Always try to merge segments to groups.
		if err := mergeSegmentsBetweenGroups(); err != nil {
			return errors.Wrapf(err, "merge groups")
		}

		// Translate all groups to target language.
		if err := translateGroups(); err != nil {
			return errors.Wrapf(err, "translate groups")
		}

		// Generate TTS audio for each group, first pass.
		if err := generateTTS(); err != nil {
			return errors.Wrapf(err, "generate tts 1pass")
		}

		// Rephrase all groups to shorter text.
		if err := rephraseAsShorterTranslations(); err != nil {
			return errors.Wrapf(err, "rephrase groups")
		}

		v.status = SrsDubbingTaskStatusDone

		return nil
	}

	v.wg.Add(1)
	go func() {
		defer v.wg.Done()

		var running bool
		func() {
			v.lock.Lock()
			defer v.lock.Unlock()

			running = v.running
			if !v.running {
				v.running = true
			}
		}()

		// Ignore if already running.
		if running {
			logger.Tf(ctx, "ignore start dubbing worker for task %v", v.UUID)
			return
		}

		defer func() {
			v.running = false
		}()

		for i := 0; i < 30; i++ {
			if err := cycle(); err != nil {
				logger.Wf(ctx, "ignore dubbing task %v err %+v", v.UUID, err)
				select {
				case <-ctx.Done():
					return
				case <-time.After(10 * time.Second):
					logger.Tf(ctx, "retry dubbing task %v", v.UUID)
				}
			} else {
				logger.Tf(ctx, "dubbing task %v done", v.UUID)
				break
			}
		}
	}()

	return nil
}

type SrsDubbingProject struct {
	// Dubbing project UUID.
	UUID string `json:"uuid"`
	// Live room title.
	Title string `json:"title"`
	// File source type.
	FileType FFprobeSourceType `json:"filetype"`
	// File source path.
	FilePath string `json:"filepath"`
	// Create time.
	CreatedAt string `json:"created_at"`

	// The AI assistant settings for ASR.
	ASR *SrsAssistant `json:"asr"`
	// The AI assistant settings for Translation.
	Translation *SrsAssistant `json:"trans"`
	// The AI assistant settings for rephrase text to shorter.
	Rephrase *SrsAssistant `json:"rephrase"`
	// The AI assistant settings for TTS.
	TTS *SrsAssistant `json:"tts"`

	// Cached source file UUID.
	SourceUUID string `json:"source"`
	// Cached source filepath.
	SourcePath string `json:"source_path"`
	// Source file format.
	SourceFormat *MediaFormat `json:"format"`
	// Source file video format.
	SourceVideo *FFprobeVideo `json:"video"`
	// Source file audio format.
	SourceAudio *FFprobeAudio `json:"audio"`

	// The dubbing task uuid, should equals to the project uuid, if task exists.
	TaskUUID string `json:"task"`
}

func NewSrsDubbingProject(opts ...func(dubbing *SrsDubbingProject)) *SrsDubbingProject {
	v := &SrsDubbingProject{
		UUID: uuid.NewString(),
		// Create time.
		CreatedAt: time.Now().Format(time.RFC3339),
		// Create a default ASR assistant.
		ASR: NewAssistant(),
		// Create a default Translation assistant.
		Translation: NewAssistant(),
		// Create a default rephrase assistant.
		Rephrase: NewAssistant(),
		// Create a default TTS assistant.
		TTS: NewAssistant(),
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

func (v *SrsDubbingProject) String() string {
	return fmt.Sprintf("uuid=%v, title=%v, filetype=%v, filepath=%v, created=%v, source=%v",
		v.UUID, v.Title, v.FileType, v.FilePath, v.CreatedAt, v.SourceUUID)
}

func (v *SrsDubbingProject) Load(ctx context.Context) error {
	dubbingUUID := v.UUID

	if r0, err := rdb.HGet(ctx, SRS_DUBBING_PROJECTS, dubbingUUID).Result(); err != nil && err != redis.Nil {
		return errors.Wrapf(err, "hget %v %v", SRS_DUBBING_PROJECTS, dubbingUUID)
	} else if r0 == "" {
		return errors.Errorf("dubbing project %v not exists", dubbingUUID)
	} else if err = json.Unmarshal([]byte(r0), v); err != nil {
		return errors.Wrapf(err, "unmarshal %v %v", dubbingUUID, r0)
	}

	return nil
}

func (v *SrsDubbingProject) Save(ctx context.Context) error {
	dubbing := v

	if b, err := json.Marshal(dubbing); err != nil {
		return errors.Wrapf(err, "marshal dubbing project")
	} else if err := rdb.HSet(ctx, SRS_DUBBING_PROJECTS, dubbing.UUID, string(b)).Err(); err != nil {
		return errors.Wrapf(err, "hset %v %v %v", SRS_DUBBING_PROJECTS, dubbing.UUID, string(b))
	}

	return nil
}

func (v *SrsDubbingProject) CheckSource(ctx context.Context, target string) error {
	if v.FileType != FFprobeSourceTypeFile && v.FileType != FFprobeSourceTypeUpload && v.FileType != FFprobeSourceTypeYTDL {
		return errors.Errorf("unsupported file type %v", v.FileType)
	}

	fileAbsPath, err := filepath.Abs(target)
	if err != nil {
		return errors.Wrapf(err, "abs %v", target)
	}

	var validExtension bool
	for _, ext := range append(serverAllowVideoFiles, serverAllowAudioFiles...) {
		if strings.HasSuffix(fileAbsPath, ext) {
			validExtension = true
			break
		}
	}
	if !validExtension {
		return errors.Errorf("invalid file extension %v, should be %v",
			fileAbsPath, append(serverAllowVideoFiles, serverAllowAudioFiles...))
	}

	if _, err := os.Stat(fileAbsPath); err != nil {
		return errors.Wrapf(err, "stat %v", fileAbsPath)
	}

	logger.Tf(ctx, "check source ok, type=%v, file=%v", v.FileType, target)
	return nil
}

func (v *SrsDubbingProject) CheckAIConfig(ctx context.Context) error {
	if v.ASR == nil {
		return errors.Errorf("invalid asr")
	}
	if v.ASR.AISecretKey == "" {
		return errors.Errorf("invalid asr secret")
	}
	if v.ASR.AIBaseURL == "" {
		return errors.Errorf("invalid asr base url")
	}
	if v.ASR.AIASRLanguage == "" {
		return errors.Errorf("invalid asr language")
	}

	logger.Tf(ctx, "check AI config ok, asr=%v", v.ASR.String())
	return nil
}
