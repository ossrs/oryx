package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ossrs/go-oryx-lib/errors"
	ohttp "github.com/ossrs/go-oryx-lib/http"
	"github.com/ossrs/go-oryx-lib/logger"
	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
)

func handleLiveRoomService(ctx context.Context, handler *http.ServeMux) error {
	ep := "/terraform/v1/live/room/create"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token, title string
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				Title *string `json:"title"`
			}{
				Token: &token, Title: &title,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			room := NewLiveRoom(func(room *SrsLiveRoom) {
				room.Title = title
				// By default, we always enable the AI assistant for user.
				room.Assistant = true
			})
			if b, err := json.Marshal(room); err != nil {
				return errors.Wrapf(err, "marshal room")
			} else if err := rdb.HSet(ctx, SRS_LIVE_ROOM, room.UUID, string(b)).Err(); err != nil {
				return errors.Wrapf(err, "hset %v %v %v", SRS_LIVE_ROOM, room.UUID, string(b))
			}

			// Note that we need to update the auth secret, because we do not use room uuid as stream name.
			roomPublishAuthKey := GenerateRoomPublishKey(room.StreamName)
			if err := rdb.HSet(ctx, SRS_AUTH_SECRET, roomPublishAuthKey, room.Secret).Err(); err != nil {
				return errors.Wrapf(err, "hset %v %v %v", SRS_AUTH_SECRET, roomPublishAuthKey, room.Secret)
			}

			ohttp.WriteData(ctx, w, r, &room)
			logger.Tf(ctx, "srs live room create ok, title=%v, room=%v", title, room.String())
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/live/room/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token, rid string
			if err := ParseBody(ctx, r.Body, &struct {
				Token    *string `json:"token"`
				RoomUUID *string `json:"uuid"`
			}{
				Token: &token, RoomUUID: &rid,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
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

			ohttp.WriteData(ctx, w, r, &room)
			logger.Tf(ctx, "srs live room query ok, uuid=%v, room=%v", rid, room.String())
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/live/room/update"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token string
			var room SrsLiveRoom
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				*SrsLiveRoom
			}{
				Token: &token, SrsLiveRoom: &room,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			// As room is a template config, to create active stage. So if we update the template, we
			// need to update the active stage object.
			if err := room.UpdateStage(ctx); err != nil {
				return errors.Wrapf(err, "update stage")
			}

			// TODO: FIXME: Should load room from redis and merge the fields.
			if b, err := json.Marshal(room); err != nil {
				return errors.Wrapf(err, "marshal room")
			} else if err := rdb.HSet(ctx, SRS_LIVE_ROOM, room.UUID, string(b)).Err(); err != nil {
				return errors.Wrapf(err, "hset %v %v %v", SRS_LIVE_ROOM, room.UUID, string(b))
			}

			// Note that we need to update the auth secret, because we do not use room uuid as stream name.
			roomPublishAuthKey := GenerateRoomPublishKey(room.StreamName)
			if err := rdb.HSet(ctx, SRS_AUTH_SECRET, roomPublishAuthKey, room.Secret).Err(); err != nil {
				return errors.Wrapf(err, "hset %v %v %v", SRS_AUTH_SECRET, roomPublishAuthKey, room.Secret)
			}

			// Limit the changing rate for AI Assistant.
			select {
			case <-ctx.Done():
			case <-time.After(300 * time.Millisecond):
			}

			ohttp.WriteData(ctx, w, r, &room)
			logger.Tf(ctx, "srs live room update ok, room=%v", room.String())
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/live/room/list"
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

			var rooms []*SrsLiveRoom
			if configs, err := rdb.HGetAll(ctx, SRS_LIVE_ROOM).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hgetall %v", SRS_LIVE_ROOM)
			} else {
				for k, v := range configs {
					var obj SrsLiveRoom
					if err = json.Unmarshal([]byte(v), &obj); err != nil {
						return errors.Wrapf(err, "unmarshal %v %v", k, v)
					}
					rooms = append(rooms, &obj)
				}
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Rooms []*SrsLiveRoom `json:"rooms"`
			}{
				Rooms: rooms,
			})
			logger.Tf(ctx, "srs live room list ok, rooms=%v", len(rooms))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/live/room/remove"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			var token, roomUUID string
			if err := ParseBody(ctx, r.Body, &struct {
				Token    *string `json:"token"`
				RoomUUID *string `json:"uuid"`
			}{
				Token: &token, RoomUUID: &roomUUID,
			}); err != nil {
				return errors.Wrapf(err, "parse body")
			}

			apiSecret := envApiSecret()
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

			if err := rdb.HDel(ctx, SRS_LIVE_ROOM, roomUUID).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hdel %v %v", SRS_LIVE_ROOM, roomUUID)
			}

			// Note that we need to update the auth secret, because we do not use room uuid as stream name.
			roomPublishAuthKey := GenerateRoomPublishKey(room.StreamName)
			if err := rdb.HDel(ctx, SRS_AUTH_SECRET, roomPublishAuthKey).Err(); err != nil {
				return errors.Wrapf(err, "hdel %v %v", SRS_AUTH_SECRET, roomPublishAuthKey)
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "srs remove room ok, uuid=%v", roomUUID)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	return nil
}

type SrsLiveRoom struct {
	// Live room UUID.
	UUID string `json:"uuid"`
	// Live room title.
	Title string `json:"title"`
	// The stream name, should never use roomUUID because it's secret.
	StreamName string `json:"stream"`
	// Live room secret.
	Secret string `json:"secret"`
	// The AI assistant settings.
	SrsAssistant
	// The current AI assistant stage, might change to others.
	// TODO: FIXME: Should not return to the client.
	StageUUID string `json:"stage_uuid"`
	// The room level authentication token, for example, popout application with this token to verify
	// the room, to prevent leaking of the bearer token.
	RoomToken string `json:"roomToken"`
	// Create time.
	CreatedAt string `json:"created_at"`
}

func NewLiveRoom(opts ...func(room *SrsLiveRoom)) *SrsLiveRoom {
	v := &SrsLiveRoom{
		UUID: uuid.NewString(),
		// The stream name of room.
		StreamName: strings.ToLower(strings.ReplaceAll(uuid.NewString(), "-", ""))[:12],
		// The secret of live room.
		Secret: strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))[:16],
		// Create time.
		CreatedAt: time.Now().Format(time.RFC3339),
		// The stage level token for popout.
		RoomToken: uuid.NewString(),
		// Create a default assistant.
		SrsAssistant: *NewAssistant(),
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

func (v *SrsLiveRoom) String() string {
	return fmt.Sprintf("uuid=%v, title=%v, stream=%v, secret=%vB, roomToken=%vB, stage=%v, assistant=<%v>",
		v.UUID, v.Title, v.StreamName, len(v.Secret), len(v.RoomToken), v.StageUUID, v.SrsAssistant.String())
}

func (v *SrsLiveRoom) UpdateStage(ctx context.Context) error {
	if stage := talkServer.QueryStageOfRoom(v.UUID); stage != nil {
		stage.UpdateFromRoom(v)
	}

	return nil
}

type SrsAssistant struct {
	// Whether enable the AI assistant.
	Assistant bool `json:"assistant"`
	// The AI name.
	AIName string `json:"aiName"`
	// The AI provider.
	AIProvider string `json:"aiProvider"`
	// The AI secret key.
	AISecretKey string `json:"aiSecretKey"`
	// The AI organization.
	AIOrganization string `json:"aiOrganization"`
	// The AI base URL.
	AIBaseURL string `json:"aiBaseURL"`

	// Whether enable the AI ASR.
	AIASREnabled bool `json:"aiAsrEnabled"`
	// The AI asr language.
	AIASRLanguage string `json:"aiAsrLanguage"`
	// The AI asr prompt type. user or user-ai.
	AIASRPrompt string `json:"aiAsrPrompt"`

	// Whether enable the AI processing.
	AIChatEnabled bool `json:"aiChatEnabled"`
	// The AI model name.
	AIChatModel string `json:"aiChatModel"`
	// The AI chat system prompt.
	AIChatPrompt string `json:"aiChatPrompt"`
	// The AI chat max window.
	AIChatMaxWindow int `json:"aiChatMaxWindow"`
	// The AI chat max words.
	AIChatMaxWords int `json:"aiChatMaxWords"`

	// Whether enable the AI post processing.
	AIPostEnabled bool `json:"aiPostEnabled"`
	// The AI model name.
	AIPostModel string `json:"aiPostModel"`
	// The AI chat system prompt.
	AIPostPrompt string `json:"aiPostPrompt"`
	// The AI chat max window.
	AIPostMaxWindow int `json:"aiPostMaxWindow"`
	// The AI chat max words.
	AIPostMaxWords int `json:"aiPostMaxWords"`

	// Whether enable the AI TTS.
	AITTSEnabled bool `json:"aiTtsEnabled"`
}

func NewAssistant(opts ...func(*SrsAssistant)) *SrsAssistant {
	v := &SrsAssistant{
		AIASREnabled: true, AIChatEnabled: true, AIPostEnabled: true, AITTSEnabled: true,
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

func (v *SrsAssistant) String() string {
	return fmt.Sprintf("assistant=%v, name=%v, provider=%v, secretKey=%vB, baseURL=%v, asr=<enabled=%v,language=%v,prompt=%v>, chat=<enabled=%v,model=%v,prompt=%v,window=%v,words=%v>, post=<enabled=%v,model=%v,prompt=%v,window=%v,words=%v>, tts=<%v>",
		v.Assistant, v.AIName, v.AIProvider, len(v.AISecretKey), v.AIBaseURL, v.AIASREnabled,
		v.AIASRLanguage, v.AIASRPrompt, v.AIChatEnabled, v.AIChatModel, v.AIChatPrompt, v.AIChatMaxWindow,
		v.AIChatMaxWords, v.AIPostEnabled, v.AIPostModel, v.AIPostPrompt, v.AIPostMaxWindow,
		v.AIPostMaxWords, v.AITTSEnabled)
}
