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
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	ohttp "github.com/ossrs/go-oryx-lib/http"
	"github.com/ossrs/go-oryx-lib/logger"

	// Use v8 because we use Go 1.16+, while v9 requires Go 1.18+
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func handleDockerHooksService(ctx context.Context, handler *http.ServeMux) error {
	versionHandler := func(w http.ResponseWriter, r *http.Request) {
		ohttp.WriteData(ctx, w, r, &struct {
			Version string `json:"version"`
		}{
			Version: strings.TrimPrefix(version, "v"),
		})
	}

	ep := "/terraform/v1/tencent/versions"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, versionHandler)

	ep = "/terraform/v1/hooks/versions"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, versionHandler)

	// See https://github.com/ossrs/srs/wiki/v4_EN_HTTPCallback
	ep = "/terraform/v1/hooks/srs/verify"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			if noAuth, err := rdb.HGet(ctx, SRS_AUTH_SECRET, "pubNoAuth").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v pubNoAuth", SRS_AUTH_SECRET)
			} else if noAuth == "true" {
				ohttp.WriteData(ctx, w, r, nil)
				logger.Tf(ctx, "srs hooks disabled")
				return nil
			}

			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}
			requestBody := string(b)

			var action string
			var streamObj SrsStream
			if err := json.Unmarshal(b, &struct {
				Action   *string `json:"action"`
				*SrsStream
			}{
				Action: &action, SrsStream: &streamObj,
			}); err != nil {
				return errors.Wrapf(err, "json unmarshal %v", string(b))
			}

			if action == "on_publish" {
				publish, err := rdb.HGet(ctx, SRS_AUTH_SECRET, "pubSecret").Result()
				if err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hget %v pubSecret", SRS_AUTH_SECRET)
				}

				// Note that we allow pass secret by params or in stream name, for example, some encoder does not support params
				// with ?secret=xxx, so it will fail when url is:
				//      rtmp://ip/live/livestream?secret=xxx
				// so user could change the url to bellow to get around of it:
				//      rtmp://ip/live/livestreamxxx
				// or simply use secret as stream:
				//      rtmp://ip/live/xxx
				// in this situation, the secret is part of stream name.
				if publish != "" && !strings.Contains(streamObj.Param, publish) && !strings.Contains(streamObj.Stream, publish) {
					return errors.Errorf("invalid stream=%v, param=%v, action=%v", streamObj.Stream, streamObj.Param, action)
				}
			}

			// Automatically add by SRS.
			streamURL := streamObj.StreamURL()
			if action == "on_publish" {
				streamObj.Update = time.Now().Format(time.RFC3339)

				b, err := json.Marshal(&streamObj)
				if err != nil {
					return errors.Wrapf(err, "marshal json")
				} else if err = rdb.HSet(ctx, SRS_STREAM_ACTIVE, streamURL, string(b)).Err(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hset %v %v %v", SRS_STREAM_ACTIVE, streamURL, string(b))
				}

				if err := rdb.HIncrBy(ctx, SRS_STAT_COUNTER, "publish", 1).Err(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hincrby %v publish 1", SRS_STAT_COUNTER)
				}
				if streamObj.IsSRT() {
					if err := rdb.HSet(ctx, SRS_STREAM_SRT_ACTIVE, streamURL, string(b)).Err(); err != nil && err != redis.Nil {
						return errors.Wrapf(err, "hset %v %v %v", SRS_STREAM_SRT_ACTIVE, streamURL, string(b))
					}
				}
				if streamObj.IsRTC() {
					if err := rdb.HSet(ctx, SRS_STREAM_RTC_ACTIVE, streamURL, string(b)).Err(); err != nil && err != redis.Nil {
						return errors.Wrapf(err, "hset %v %v %v", SRS_STREAM_RTC_ACTIVE, streamURL, string(b))
					}
				}
			} else if action == "on_unpublish" {
				if err := rdb.HDel(ctx, SRS_STREAM_ACTIVE, streamURL).Err(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hset %v %v", SRS_STREAM_ACTIVE, streamURL)
				}
				if streamObj.IsSRT() {
					if err := rdb.HDel(ctx, SRS_STREAM_SRT_ACTIVE, streamURL).Err(); err != nil && err != redis.Nil {
						return errors.Wrapf(err, "hset %v %v", SRS_STREAM_SRT_ACTIVE, streamURL)
					}
				}
				if streamObj.IsRTC() {
					if err := rdb.HDel(ctx, SRS_STREAM_RTC_ACTIVE, streamURL).Err(); err != nil && err != redis.Nil {
						return errors.Wrapf(err, "hset %v %v", SRS_STREAM_RTC_ACTIVE, streamURL)
					}
				}
			} else if action == "on_play" {
				if err := rdb.HIncrBy(ctx, SRS_STAT_COUNTER, "play", 1).Err(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hincrby %v play 1", SRS_STAT_COUNTER)
				}
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "srs hooks ok, action=%v, %v, %v", action, streamObj.String(), requestBody)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	secretQueryHandler := func(w http.ResponseWriter, r *http.Request) {
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

			publish, err := rdb.HGet(ctx, SRS_AUTH_SECRET, "pubSecret").Result()
			if err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v pubSecret", SRS_AUTH_SECRET)
			}
			if publish == "" {
				return errors.New("system not boot yet")
			}

			ohttp.WriteData(ctx, w, r, &struct {
				Publish string `json:"publish"`
			}{
				Publish: publish,
			})
			logger.Tf(ctx, "srs secret ok ok, token=%vB", len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	}

	ep = "/terraform/v1/hooks/srs/secret"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, secretQueryHandler)

	ep = "/terraform/v1/hooks/srs/secret/query"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, secretQueryHandler)

	ep = "/terraform/v1/hooks/srs/secret/update"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var token, secret string
			if err := json.Unmarshal(b, &struct {
				Token  *string `json:"token"`
				Secret *string `json:"secret"`
			}{
				Token: &token, Secret: &secret,
			}); err != nil {
				return errors.Wrapf(err, "json unmarshal %v", string(b))
			}
			if secret == "" {
				return errors.New("no secret")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			// Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
			// See https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
			if _, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(apiSecret), nil
			}); err != nil {
				return errors.Wrapf(err, "verify token %v", token)
			}

			if err := rdb.HSet(ctx, SRS_AUTH_SECRET, "pubSecret", secret).Err(); err != nil {
				return errors.Wrapf(err, "hset %v pubSecret %v", SRS_AUTH_SECRET, secret)
			}
			if err := rdb.Set(ctx, SRS_SECRET_PUBLISH, secret, 0).Err(); err != nil {
				return errors.Wrapf(err, "set %v %v", SRS_SECRET_PUBLISH, secret)
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "hooks update secret, secret=%vB, token=%vB", len(secret), len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	ep = "/terraform/v1/hooks/srs/secret/disable"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var token string
			var pubNoAuth bool
			if err := json.Unmarshal(b, &struct {
				Token     *string `json:"token"`
				PubNoAuth *bool   `json:"pubNoAuth"`
			}{
				Token: &token, PubNoAuth: &pubNoAuth,
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

			if err := rdb.HSet(ctx, SRS_AUTH_SECRET, "pubNoAuth", fmt.Sprintf("%v", pubNoAuth)).Err(); err != nil {
				return errors.Wrapf(err, "hset %v pubSecret %v", SRS_AUTH_SECRET, pubNoAuth)
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "hooks disable secret, pubNoAuth=%v, token=%vB", pubNoAuth, len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	// See https://console.cloud.tencent.com/cam
	ep = "/terraform/v1/tencent/cam/secret"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var token, secretId, secretKey string
			if err := json.Unmarshal(b, &struct {
				Token     *string `json:"token"`
				SecretID  *string `json:"secretId"`
				SecretKey *string `json:"secretKey"`
			}{
				Token: &token, SecretID: &secretId, SecretKey: &secretKey,
			}); err != nil {
				return errors.Wrapf(err, "json unmarshal %v", string(b))
			}
			if secretId == "" {
				return errors.New("no secretId")
			}
			if secretKey == "" {
				return errors.New("no secretKey")
			}

			apiSecret := os.Getenv("SRS_PLATFORM_SECRET")
			// Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
			// See https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
			if _, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(apiSecret), nil
			}); err != nil {
				return errors.Wrapf(err, "verify token %v", token)
			}

			// Query, verify and setup the secret ID and key.
			var appID, ownerUIN string
			if true {
				cpf := profile.NewClientProfile()
				cpf.HttpProfile.Endpoint = TENCENT_CLOUD_CAM_ENDPOINT

				if client, err := cam.NewClient(common.NewCredential(secretId, secretKey), "", cpf); err != nil {
					return errors.Wrapf(err, "create tencent cloud sdk")
				} else if response, err := client.GetUserAppIdWithContext(ctx, cam.NewGetUserAppIdRequest()); err != nil {
					return errors.Wrapf(err, "tencent cloud api")
				} else if response.Response.AppId == nil || *response.Response.AppId == 0 {
					return errors.Wrapf(err, "invalid response %v", response.ToJsonString())
				} else {
					appID, ownerUIN = fmt.Sprintf("%v", *response.Response.AppId), *response.Response.OwnerUin
					logger.Tf(ctx, "CAM query appID=%v, ownerUIN=%v", appID, ownerUIN)
				}
			}

			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("appid=%v, ownerUIN=%v", appID, ownerUIN))

			if err := rdb.HSet(ctx, SRS_TENCENT_CAM, "appId", appID).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v appId %v", SRS_TENCENT_CAM, appID)
			}
			if err := rdb.HSet(ctx, SRS_TENCENT_CAM, "secretId", secretId).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v secretId %v", SRS_TENCENT_CAM, secretId)
			}
			if err := rdb.HSet(ctx, SRS_TENCENT_CAM, "secretKey", secretKey).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v secretKey %v", SRS_TENCENT_CAM, secretKey)
			}
			if err := rdb.HSet(ctx, SRS_TENCENT_CAM, "uin", ownerUIN).Err(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hset %v uin %v", SRS_TENCENT_CAM, ownerUIN)
			}

			// Create COS bucket if not exists.
			// Note that should never use headBucket to query it, because it cause always fail if bucket not found, even though
			// the bucket is created.
			var bucketName string
			var createBucket bool
			if bucket, err := rdb.HGet(ctx, SRS_TENCENT_COS, "bucket").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v bucket", SRS_TENCENT_COS)
			} else if bucketName = bucket; bucketName == "" {
				// Add nonce to bucket name, to avoid conflict on different region as bellow:
				//    The requested bucket has already existed in other region.
				nonce := fmt.Sprintf("%x", rand.Int63())[:4]
				bucketName = fmt.Sprintf("srs-lighthouse-%v-%v", nonce, appID)
				createBucket = true
			}

			var cosClient *cos.Client
			location := fmt.Sprintf("%v.cos.%v.myqcloud.com", bucketName, conf.Region)
			if true {
				u, err := url.Parse(fmt.Sprintf("https://%v", location))
				if err != nil {
					return errors.Wrapf(err, "parse %v", fmt.Sprintf("https://%v", location))
				}

				// Create the bucket with private ACL.
				// See https://cloud.tencent.com/document/product/436/36118
				cosClient = cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
					Transport: &cos.AuthorizationTransport{SecretID: secretId, SecretKey: secretKey},
				})
			}
			sb.WriteString(fmt.Sprintf(", bucket=%v, location=%v", bucketName, location))

			if createBucket {
				if _, err := cosClient.Bucket.Put(ctx, &cos.BucketPutOptions{XCosACL: "private"}); err != nil {
					return errors.Wrapf(err, "create bucket %v", bucketName)
				}

				// Save information to redis.
				if err := rdb.HSet(ctx, SRS_TENCENT_COS, "bucket", bucketName).Err(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hset %v bucket %v", SRS_TENCENT_COS, bucketName)
				}
				if err := rdb.HSet(ctx, SRS_TENCENT_COS, "location", location).Err(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hset %v location %v", SRS_TENCENT_COS, location)
				}
				logger.Tf(ctx, "COS create bucket=%v, location=%v", bucketName, location)
				sb.WriteString(", createBucket=true")
			}

			// Setup COS bucket if no policy.
			if policy, err := rdb.HGet(ctx, SRS_TENCENT_COS, "policy").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v policy", SRS_TENCENT_COS)
			} else if policy == "" {
				// Allow read without list files actions:
				//    cos:GetBucket
				//    cos:GetBucketObjectVersions
				// See https://cloud.tencent.com/document/product/436/43812
				opt := cos.BucketPutPolicyOptions{
					Version: "2.0",
					Statement: []cos.BucketStatement{
						{
							Effect: "Allow",
							Principal: map[string][]string{
								"qcs": []string{"qcs::cam::anyone:anyone"},
							},
							Action: []string{
								"name/cos:HeadBucket",
								"name/cos:ListMultipartUploads",
								"name/cos:ListParts",
								"name/cos:GetObject",
								"name/cos:HeadObject",
								"name/cos:OptionsObject",
							},
							Resource: []string{
								fmt.Sprintf("qcs::cos:%v:uid/%v:%v/*", conf.Region, appID, bucketName),
							},
						},
					},
				}
				if _, err = cosClient.Bucket.PutPolicy(ctx, &opt); err != nil {
					return errors.Wrapf(err, "put cos policy")
				}

				// Save information to redis.
				if err := rdb.HSet(ctx, SRS_TENCENT_COS, "policy", "read-without-list-files").Err(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hset %v policy %v", SRS_TENCENT_COS, "read-without-list-files")
				}
				logger.Tf(ctx, "COS create policy ok, bucket=%v", bucketName)
				sb.WriteString(", createPolicy=true")
			}

			// Setup the CORS of bucket.
			if cors, err := rdb.HGet(ctx, SRS_TENCENT_COS, "cors").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v cors", SRS_TENCENT_COS)
			} else if cors == "" {
				// See https://cloud.tencent.com/document/product/436/43811
				opt := cos.BucketPutCORSOptions{
					Rules: []cos.BucketCORSRule{
						{
							AllowedOrigins: []string{"*"},
							AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
							AllowedHeaders: []string{"*"},
							ExposeHeaders:  []string{"ETag", "Content-Length", "X-Cos-Request-Id"},
							MaxAgeSeconds:  0,
						},
					},
				}
				if _, err = cosClient.Bucket.PutCORS(ctx, &opt); err != nil {
					return errors.Wrapf(err, "put cos policy")
				}

				// Save information to redis.
				if err := rdb.HSet(ctx, SRS_TENCENT_COS, "cors", "true").Err(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hset %v cors %v", SRS_TENCENT_COS, "true")
				}
				logger.Tf(ctx, "COS create CORS ok, bucket=%v", bucketName)
				sb.WriteString(", createCORS=true")
			}

			// Create cloud VoD service and application if not exists.
			var vodAppID, vodAppName string
			var createVodApp bool
			if service, err := rdb.HGet(ctx, SRS_TENCENT_VOD, "service").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v service", SRS_TENCENT_VOD)
			} else if vodAppID = service; vodAppID == "" || vodAppID == "ok" {
				// Add nonce to bucket name, to avoid conflict on different region as bellow:
				//    The requested bucket has already existed in other region.
				nonce := fmt.Sprintf("%x", rand.Int63())[:4]
				vodAppName = fmt.Sprintf("srs-lighthouse-%v", nonce)
				// For previous platform, the service is set to string ok, so we also create an application.
				vodAppID = ""
				createVodApp = true
			}

			var vodClient *vod.Client
			if true {
				cpf := profile.NewClientProfile()
				cpf.HttpProfile.Endpoint = TENCENT_CLOUD_VOD_ENDPOINT

				client, err := vod.NewClient(common.NewCredential(secretId, secretKey), "", cpf)
				if err != nil {
					return errors.Wrapf(err, "create vod sdk client")
				}
				vodClient = client
			}

			if createVodApp {
				request := vod.NewCreateSubAppIdRequest()

				request.Name = common.StringPtr(vodAppName)
				request.Description = common.StringPtr("Application VoD for srs-cloud")

				if response, err := vodClient.CreateSubAppIdWithContext(ctx, request); err != nil {
					return errors.Wrapf(err, "create vod appid")
				} else {
					vodAppID = fmt.Sprintf("%v", *response.Response.SubAppId)
				}

				// Save information to redis.
				if err := rdb.HSet(ctx, SRS_TENCENT_VOD, "service", vodAppID).Err(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hset %v service %v", SRS_TENCENT_VOD, vodAppID)
				}
				logger.Tf(ctx, "VOD create appID=%v, name=%v", vodAppID, vodAppName)
				sb.WriteString(", createVodApp=true")
			}
			sb.WriteString(fmt.Sprintf(", vodAppID=%v", vodAppID))

			vodAppIDParsed, err := strconv.ParseInt(vodAppID, 10, 64)
			if err != nil {
				return errors.Wrapf(err, "parse %v", vodAppID)
			}
			sb.WriteString(fmt.Sprintf(", vodAppIDParsed=%v", vodAppIDParsed))

			// Create tencent vod storage region.
			if storage, err := rdb.HGet(ctx, SRS_TENCENT_VOD, "storage").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v storage", SRS_TENCENT_VOD)
			} else if storage == "" {
				// See https://cloud.tencent.com/document/product/266/72481
				if true {
					request := vod.NewCreateStorageRegionRequest()

					request.StorageRegion = common.StringPtr(conf.Region)
					request.SubAppId = common.Uint64Ptr(uint64(vodAppIDParsed))

					if _, err := vodClient.CreateStorageRegionWithContext(ctx, request); err != nil {
						return errors.Wrapf(err, "create storage region=%v, app=%v", conf.Region, vodAppIDParsed)
					}
				}

				// See https://cloud.tencent.com/document/product/266/72479
				if true {
					request := vod.NewModifyDefaultStorageRegionRequest()

					request.StorageRegion = common.StringPtr(conf.Region)
					request.SubAppId = common.Uint64Ptr(uint64(vodAppIDParsed))

					if _, err := vodClient.ModifyDefaultStorageRegionWithContext(ctx, request); err != nil {
						return errors.Wrapf(err, "set default storage region=%v, app=%v", conf.Region, vodAppIDParsed)
					}
				}

				// Save information to redis.
				if err := rdb.HSet(ctx, SRS_TENCENT_VOD, "storage", conf.Region).Err(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hset %v storage %v", SRS_TENCENT_VOD, conf.Region)
				}
				logger.Tf(ctx, "VOD create storage ok, appID=%v, region=%v", vodAppID, conf.Region)
				sb.WriteString(fmt.Sprintf(", createVodStorage=true, vodRegion=%v", conf.Region))
			}

			// Query templates.
			if transcode, err := rdb.HGet(ctx, SRS_TENCENT_VOD, "transcode").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v transcode", SRS_TENCENT_VOD)
			} else if transcode == "" {
				// See https://cloud.tencent.com/document/product/266/33769
				request := vod.NewDescribeTranscodeTemplatesRequest()

				request.SubAppId = common.Uint64Ptr(uint64(vodAppIDParsed))
				request.Type = common.StringPtr("Preset")
				request.ContainerType = common.StringPtr("Video")
				request.TEHDType = common.StringPtr("Common")
				request.Offset = common.Uint64Ptr(0)
				request.Limit = common.Uint64Ptr(100)

				var templates []*vod.TranscodeTemplate
				if response, err := vodClient.DescribeTranscodeTemplatesWithContext(ctx, request); err != nil {
					return errors.Wrapf(err, "set query transcode templates region=%v, app=%v", conf.Region, vodAppIDParsed)
				} else if *response.Response.TotalCount == 0 {
					return errors.Errorf("invalid transcode template %v", response.ToJsonString())
				} else {
					for _, template := range response.Response.TranscodeTemplateSet {
						if !strings.Contains(*template.Name, "Deprecated") {
							templates = append(templates, template)
						}
					}
				}
				if len(templates) == 0 {
					return errors.New("no vod transcode templates")
				}

				if b, err := json.Marshal(&struct {
					NN        int                      `json:"nn"`
					Templates []*vod.TranscodeTemplate `json:"templates"`
				}{
					NN: len(templates), Templates: templates,
				}); err != nil {
					return errors.Wrapf(err, "json marshal %v", templates)
				} else if err = rdb.HSet(ctx, SRS_TENCENT_VOD, "transcode", string(b)).Err(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hset %v transcode %v", SRS_TENCENT_VOD, string(b))
				}
				logger.Tf(ctx, "VOD query templates ok, nn=%v", len(templates))
				sb.WriteString(fmt.Sprintf(", createVodTranscode=true, templates=%v", len(templates)))
			}

			// Filter the remux template, covert to MP4.
			if remux, err := rdb.HGet(ctx, SRS_TENCENT_VOD, "remux").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v remux", SRS_TENCENT_VOD)
			} else if remux == "" {
				var templates []*vod.TranscodeTemplate
				if transcode, err := rdb.HGet(ctx, SRS_TENCENT_VOD, "transcode").Result(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hget %v transcode", SRS_TENCENT_VOD)
				} else if err = json.Unmarshal([]byte(transcode), &struct {
					Templates *[]*vod.TranscodeTemplate `json:"templates"`
				}{
					Templates: &templates,
				}); err != nil {
					return errors.Wrapf(err, "parse %v", transcode)
				}

				var remuxMp4 *vod.TranscodeTemplate
				for _, template := range templates {
					if *template.Container == "mp4" && *template.VideoTemplate.Codec == "copy" && *template.AudioTemplate.Codec == "copy" {
						remuxMp4 = template
						break
					}
				}
				if remuxMp4 == nil {
					return errors.Errorf("no remux template %v", templates)
				}

				target := VodTranscodeTemplate{
					Definition: *remuxMp4.Definition, Name: *remuxMp4.Name, Comment: *remuxMp4.Comment,
					Container: *remuxMp4.Container, Update: *remuxMp4.UpdateTime,
				}
				if b, err := json.Marshal(&target); err != nil {
					return errors.Wrapf(err, "json marshal %v", templates)
				} else if err = rdb.HSet(ctx, SRS_TENCENT_VOD, "remux", string(b)).Err(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hset %v remux %v", SRS_TENCENT_VOD, string(b))
				}
				logger.Tf(ctx, "VOD set remux templates ok, %v", target.String())
				sb.WriteString(fmt.Sprintf(", createVodRemux=true, remux=%v", target.Name))
			}

			// Query the global vod domain.
			if domain, err := rdb.HGet(ctx, SRS_TENCENT_VOD, "domain").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v domain", SRS_TENCENT_VOD)
			} else if domain == "" {
				// See https://cloud.tencent.com/document/product/266/54176
				request := vod.NewDescribeVodDomainsRequest()

				request.SubAppId = common.Uint64Ptr(uint64(vodAppIDParsed))

				var target *vod.DomainDetailInfo
				if response, err := vodClient.DescribeVodDomainsWithContext(ctx, request); err != nil {
					return errors.Wrapf(err, "set query vod domains region=%v, app=%v", conf.Region, vodAppIDParsed)
				} else if *response.Response.TotalCount == 0 {
					return errors.Errorf("invalid vod domains %v", response.ToJsonString())
				} else {
					for _, ds := range response.Response.DomainSet {
						if strings.Contains(*ds.DeployStatus, "Online") {
							target = ds
							break
						}
					}
				}

				if err = rdb.HSet(ctx, SRS_TENCENT_VOD, "domain", *target.Domain).Err(); err != nil && err != redis.Nil {
					return errors.Wrapf(err, "hset %v domain %v", SRS_TENCENT_VOD, *target.Domain)
				}
				logger.Tf(ctx, "VOD set domain ok, %v", *target.Domain)
				sb.WriteString(fmt.Sprintf(", createVodDomain=true, vodDomain=%v", *target.Domain))
			}

			ohttp.WriteData(ctx, w, r, nil)
			logger.Tf(ctx, "CAM: Update ok, %v, token=%vB", sb.String(), len(token))
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	if err := handleOnHls(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle hooks")
	}

	return nil
}

func handleOnHls(ctx context.Context, handler *http.ServeMux) error {
	// TODO: FIXME: Fixed token.
	// See https://github.com/ossrs/srs/wiki/v4_EN_HTTPCallback
	ep := "/terraform/v1/hooks/srs/hls"
	logger.Tf(ctx, "Handle %v", ep)
	handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrapf(err, "read body")
			}

			var msg SrsOnHlsMessage
			if err := json.Unmarshal(b, &msg); err != nil {
				return errors.Wrapf(err, "json unmarshal %v", string(b))
			}
			if msg.Action != "on_hls" {
				return errors.Errorf("invalid action=%v", msg.Action)
			}
			if _, err := os.Stat(msg.File); err != nil {
				return errors.Wrapf(err, "invalid ts file %v", msg.File)
			}
			logger.Tf(ctx, "on_hls ok, %v", string(b))

			// Create a Record task if enabled.
			if recordAll, err := rdb.HGet(ctx, SRS_RECORD_PATTERNS, "all").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v all", SRS_RECORD_PATTERNS)
			} else if recordAll == "true" {
				if err = recordWorker.OnHlsTsMessage(ctx, &msg); err != nil {
					return errors.Wrapf(err, "feed %v", msg.String())
				}
				logger.Tf(ctx, "record %v", msg.String())
			}

			// Create a DVR task if enabled.
			if dvrAll, err := rdb.HGet(ctx, SRS_DVR_PATTERNS, "all").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v all", SRS_DVR_PATTERNS)
			} else if dvrAll == "true" {
				if err = dvrWorker.OnHlsTsMessage(ctx, &msg); err != nil {
					return errors.Wrapf(err, "feed %v", msg.String())
				}
				logger.Tf(ctx, "dvr %v", msg.String())
			}

			// Create a VOD task if enabled.
			if vodAll, err := rdb.HGet(ctx, SRS_VOD_PATTERNS, "all").Result(); err != nil && err != redis.Nil {
				return errors.Wrapf(err, "hget %v all", SRS_VOD_PATTERNS)
			} else if vodAll == "true" {
				if err = vodWorker.OnHlsTsMessage(ctx, &msg); err != nil {
					return errors.Wrapf(err, "feed %v", msg.String())
				}
				logger.Tf(ctx, "vod %v", msg.String())
			}

			ohttp.WriteData(ctx, w, r, nil)
			return nil
		}(); err != nil {
			ohttp.WriteError(ctx, w, r, err)
		}
	})

	if err := recordWorker.Handle(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle record")
	}

	if err := dvrWorker.Handle(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle dvr")
	}

	if err := vodWorker.Handle(ctx, handler); err != nil {
		return errors.Wrapf(err, "handle vod")
	}

	return nil
}
