package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/ossrs/go-oryx-lib/errors"
	ohttp "github.com/ossrs/go-oryx-lib/http"
	"github.com/ossrs/go-oryx-lib/logger"
	"net/http"
	"os"
	"strings"
	"time"
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

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			room := &SrsLiveRoom{
				UUID: uuid.NewString(),
				// The title of live room.
				Title: title,
				// The secret of live room.
				Secret: strings.ReplaceAll(uuid.NewString(), "-", ""),
				// Create time.
				CreatedAt: time.Now().Format(time.RFC3339),
			}
			if b, err := json.Marshal(room); err != nil {
				return errors.Wrapf(err, "marshal room")
			} else if err := rdb.HSet(ctx, SRS_LIVE_ROOM, room.UUID, string(b)).Err(); err != nil {
				return errors.Wrapf(err, "hset %v %v %v", SRS_LIVE_ROOM, room.UUID, string(b))
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
			var token, rid string
			var assistant bool
			if err := ParseBody(ctx, r.Body, &struct {
				Token *string `json:"token"`
				// The UUID of room.
				RoomUUID *string `json:"uuid"`
				// Whether enable the AI assistant.
				Assistant *bool `json:"assistant"`
			}{
				Token: &token, RoomUUID: &rid, Assistant: &assistant,
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

			if room.Assistant != assistant {
				room.Assistant = assistant

				if b, err := json.Marshal(room); err != nil {
					return errors.Wrapf(err, "marshal room")
				} else if err := rdb.HSet(ctx, SRS_LIVE_ROOM, room.UUID, string(b)).Err(); err != nil {
					return errors.Wrapf(err, "hset %v %v %v", SRS_LIVE_ROOM, room.UUID, string(b))
				}
			}

			// Limit the changing rate for AI Assistant.
			select {
			case <-ctx.Done():
			case <-time.After(300 * time.Millisecond):
			}

			ohttp.WriteData(ctx, w, r, &room)
			logger.Tf(ctx, "srs live room update ok, uuid=%v, assistant=%v, room=%v", rid, assistant, room.String())
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

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
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
			logger.Tf(ctx, "srs live room create ok, rooms=%v", len(rooms))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/live/room/remove"
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

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			if err := Authenticate(ctx, apiSecret, token, r.Header); err != nil {
				return errors.Wrapf(err, "authenticate")
			}

			if r0, err := rdb.HGet(ctx, SRS_LIVE_ROOM, rid).Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v %v", SRS_LIVE_ROOM, rid)
			} else if r0 == "" {
				return errors.Errorf("live room %v not exists", rid)
			}

			if err := rdb.HDel(ctx, SRS_LIVE_ROOM, rid).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hdel %v %v", SRS_LIVE_ROOM, rid)
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "srs remove room ok, uuid=%v", rid)
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
	// Live room secret.
	Secret string `json:"secret"`

	// Whether enable the AI assistant.
	Assistant bool `json:"assistant"`

	// Create time.
	CreatedAt string `json:"created_at"`
}

func (v *SrsLiveRoom) String() string {
	return fmt.Sprintf("uuid=%v, title=%v, secret=%v, assistant=%v", v.UUID, v.Title, v.Secret, v.Assistant)
}
