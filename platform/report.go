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
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"

	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
)

// queryLatestVersion is to query the latest and stable version from SRS Stack API.
func queryLatestVersion(ctx context.Context) (*Versions, error) {
	// Request release api with params.
	params := make(map[string]string)

	// Generate and setup the node id.
	if r0, err := rdb.HGet(ctx, SRS_TENCENT_LH, "node").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v node", SRS_TENCENT_LH)
	} else if r0 != "" {
		params["nid"] = r0
	}

	// Report about local Reocrd.
	if r0, err := rdb.HGet(ctx, SRS_RECORD_PATTERNS, "all").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v all", SRS_RECORD_PATTERNS)
	} else if r0 == "true" {
		params["rkd"] = "1"
	}
	if r0, err := rdb.HLen(ctx, SRS_RECORD_M3U8_ARTIFACT).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_RECORD_M3U8_ARTIFACT)
	} else if r0 > 0 {
		params["rkdn"] = fmt.Sprintf("%v", r0)
	}

	// Report about COS and resource usage.
	if r0, err := rdb.HGet(ctx, SRS_TENCENT_COS, "bucket").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v bucket", SRS_TENCENT_COS)
	} else if r0 == "true" {
		params["cos"] = "1"
	}
	if r0, err := rdb.HLen(ctx, SRS_DVR_M3U8_ARTIFACT).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_DVR_M3U8_ARTIFACT)
	} else if r0 > 0 {
		params["cosn"] = fmt.Sprintf("%v", r0)
	}

	// Report about VoD and resource usage.
	if r0, err := rdb.HGet(ctx, SRS_TENCENT_VOD, "storage").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v storage", SRS_TENCENT_VOD)
	} else if r0 == "true" {
		params["vod"] = "1"
	}
	if r0, err := rdb.HLen(ctx, SRS_DVR_M3U8_ARTIFACT).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_DVR_M3U8_ARTIFACT)
	} else if r0 > 0 {
		params["vodn"] = fmt.Sprintf("%v", r0)
	}

	// Report about FFmpeg forwarding.
	if r0, err := rdb.HLen(ctx, SRS_FORWARD_TASK).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_FORWARD_TASK)
	} else if r0 > 0 {
		params["forward"] = fmt.Sprintf("%v", r0)
	}

	// Report about FFmpeg virtual live from file or other source.
	if r0, err := rdb.HLen(ctx, SRS_VLIVE_TASK).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_VLIVE_TASK)
	} else if r0 > 0 {
		params["vfile"] = fmt.Sprintf("%v", r0)
	}
	if configs, err := rdb.HGetAll(ctx, SRS_VLIVE_CONFIG).Result(); err == nil {
		for _, v := range configs {
			var obj VLiveConfigure
			if err = json.Unmarshal([]byte(v), &obj); err == nil {
				for _, vFile := range obj.Files {
					if vFile.Type == SrsVLiveSourceTypeFile {
						params["vft0"] = "1"
					} else if vFile.Type == SrsVLiveSourceTypeUpload {
						params["vft1"] = "1"
					} else if vFile.Type == SrsVLiveSourceTypeStream {
						params["vft2"] = "1"
					}
				}
			}
		}
	}

	// Report about active streams.
	if r0, err := rdb.HGet(ctx, SRS_STAT_COUNTER, "publish").Int64(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v publish", SRS_STAT_COUNTER)
	} else if r0 > 0 {
		if err = rdb.HSet(ctx, SRS_STAT_COUNTER, "publish", 0).Err(); err != nil && err != redis.Nil {
			return nil, errors.Wrapf(err, "hset %v publish", SRS_STAT_COUNTER)
		}
		params["streams"] = fmt.Sprintf("%v", r0)
	}

	// Report about active players.
	if r0, err := rdb.HGet(ctx, SRS_STAT_COUNTER, "play").Int64(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v play", SRS_STAT_COUNTER)
	} else if r0 > 0 {
		if err = rdb.HSet(ctx, SRS_STAT_COUNTER, "play", 0).Err(); err != nil && err != redis.Nil {
			return nil, errors.Wrapf(err, "hset %v play", SRS_STAT_COUNTER)
		}
		params["players"] = fmt.Sprintf("%v", r0)
	}

	// Report about SRT stream.
	if r0, err := rdb.HLen(ctx, SRS_STREAM_SRT_ACTIVE).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_STREAM_SRT_ACTIVE)
	} else if r0 > 0 {
		params["srt"] = fmt.Sprintf("%v", r0)
	}

	// Report about WebRTC stream.
	if r0, err := rdb.HLen(ctx, SRS_STREAM_RTC_ACTIVE).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_STREAM_RTC_ACTIVE)
	} else if r0 > 0 {
		params["rtc"] = fmt.Sprintf("%v", r0)
	}

	// Report about beian feature.
	if r0, err := rdb.HLen(ctx, SRS_BEIAN).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hlen %v", SRS_BEIAN)
	} else if r0 > 0 {
		params["beian"] = fmt.Sprintf("%v", r0)
	}

	// Report about HTTPS feature.
	if r0, err := rdb.Get(ctx, SRS_HTTPS).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "get %v", SRS_HTTPS)
	} else if r0 != "" {
		params["https"] = r0
	}

	// Report about locale feature.
	if r0, err := rdb.Get(ctx, SRS_LOCALE).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "get %v", SRS_LOCALE)
	} else if r0 != "" {
		params["lan"] = r0
	}

	// Report about upgrade window feature.
	if r0, err := rdb.HGet(ctx, SRS_UPGRADE_WINDOW, "update").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v update", SRS_UPGRADE_WINDOW)
	} else if r0 == "true" {
		params["uwin"] = "1"
	}

	// Report whether start as develop environment.
	if os.Getenv("NODE_ENV") == "development" {
		params["dev"] = "1"
	}

	// Report whether enable SRS development version.
	if r0, err := rdb.HGet(ctx, SRS_CONTAINER_DISABLED, srsDevDockerName).Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v %v", SRS_CONTAINER_DISABLED, srsDevDockerName)
	} else if r0 == "false" {
		params["srsd"] = "1"
	}

	// Report about the platform.
	if r0, err := rdb.HGet(ctx, SRS_TENCENT_LH, "platform").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v platform", SRS_TENCENT_LH)
	} else if r0 != "" {
		params["plat"] = r0
	}

	if r0, err := rdb.HGet(ctx, SRS_TENCENT_LH, "cloud").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v cloud", SRS_TENCENT_LH)
	} else if r0 != "" {
		params["cloud"] = r0
	}

	if r0, err := rdb.HGet(ctx, SRS_TENCENT_LH, "region").Result(); err != nil && err != redis.Nil {
		return nil, errors.Wrapf(err, "hget %v region", SRS_TENCENT_LH)
	} else if r0 != "" {
		params["region"] = r0
	}

	// Refresh the version from api.
	refreshVersion := func(ctx context.Context, params map[string]string, resObj interface{}) error {
		if params == nil {
			return errors.Errorf("no params")
		}

		params["version"] = version
		params["ts"] = fmt.Sprintf("%v", time.Now().UnixNano()/int64(time.Millisecond))
		releaseServer := "https://api.ossrs.net"
		if os.Getenv("LOCAL_RELEASE") != "" {
			releaseServer = "http://localhost:2023"
		}
		logger.Tf(ctx, "Query %v with %v", releaseServer, params)

		queries := []string{}
		for k, v := range params {
			queries = append(queries, fmt.Sprintf("%v=%v", k, v))
		}
		requestURL := fmt.Sprintf("%v/terraform/v1/releases?%v", releaseServer, strings.Join(queries, "&"))
		res, err := http.Get(requestURL)
		if err != nil {
			return errors.Wrapf(err, "request %v", requestURL)
		}
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrapf(err, "read %v", requestURL)
		}

		if err := json.Unmarshal(b, resObj); err != nil {
			return errors.Wrapf(err, "parse %v of %v", string(b), requestURL)
		}

		logger.Tf(ctx, "execApi req=%v, res is %v", params, resObj)
		return nil
	}

	versions := &Versions{
		Version: version,
	}
	if err := refreshVersion(ctx, params, versions); err != nil {
		return nil, errors.Wrapf(err, "refresh version with %v", params)
	}
	return versions, nil
}
