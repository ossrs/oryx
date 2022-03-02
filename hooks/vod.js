'use strict';

// For components in docker, connect by host.
const config = {
  redis:{
    host: process.env.NODE_ENV === 'development' ? 'localhost' : 'mgmt.srs.local',
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const utils = require('js-core/utils');
const errs = require('js-core/errs');
const jwt = require('jsonwebtoken');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const keys = require('js-core/keys');
const VodClient = require("tencentcloud-sdk-nodejs").vod.v20180717.Client;
const m3u8Generator = require('./m3u8Generator');

exports.handle = (router) => {
  // Query the VoD patterns.
  router.all('/terraform/v1/hooks/vod/query', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    const all = await redis.hget(keys.redis.SRS_VOD_PATTERNS, 'all');

    const appId = await redis.hget(keys.redis.SRS_TENCENT_CAM, 'appId');
    const secretId = await redis.hget(keys.redis.SRS_TENCENT_CAM, 'secretId');
    const secretKey = await redis.hget(keys.redis.SRS_TENCENT_CAM, 'secretKey');

    console.log(`vod apply ok, all=${all}, appId=${appId}, secretId=${secretId}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      all: all === 'true',
      secret: !!(appId && secretId && secretKey),
    });
  });

  // Setup the VoD patterns.
  router.all('/terraform/v1/hooks/vod/apply', async (ctx) => {
    const {token, all} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    if (all !== true && all !== false) throw utils.asError(errs.sys.invalid, errs.status.args, `invalid all=${all}`);

    const r0 = await redis.hset(keys.redis.SRS_VOD_PATTERNS, 'all', all);

    console.log(`vod apply ok, all=${all}, r0=${r0}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  // List the VoD files.
  router.all('/terraform/v1/hooks/vod/files', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    const files = [];
    const fileObjs = {};
    const [cursor, fileKVs] = await redis.hscan(keys.redis.SRS_VOD_M3U8_METADATA, 0, '*', 100);
    for (let i = 0; i < fileKVs.length; i += 2) {
      const file = fileKVs[i + 1];
      if (file) {
        const fileObj = JSON.parse(file);
        files.push(fileObj);
        fileObjs[fileObj.fileId] = fileObj;
      }
    }

    // TODO: FIXME: Query in worker thread.
    await queryMediaInfo(files, fileObjs);

    const r0 = files.map(e => {
      return {
        uuid: e.uuid,
        vhost: e.vhost,
        app: e.app,
        stream: e.stream,
        file: e.fileId,
        media: e.mediaUrl,
        task: e.taskObj,
        progress: e.progress,
        update: e.update,
        nn: e.files.length,
        duration: e.files.reduce((p, c) => p + (c.duration || 0), 0),
        size: e.files.reduce((p, c) => p + (c.size || 0), 0),
      };
    });

    console.log(`vod files ok, cursor=${cursor}, files=${files.length}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, r0);
  });

  // Generate m3u8 to play.
  router.all('/terraform/v1/hooks/vod/hls/:uuid.m3u8', async (ctx) => {
    const {uuid} = ctx.params;
    if (!uuid) throw utils.asError(errs.sys.empty, errs.status.args, `no param uuid`);

    const domain = await redis.hget(keys.redis.SRS_TENCENT_VOD, 'domain');
    const metadata = await redis.hget(keys.redis.SRS_VOD_M3U8_METADATA, uuid);
    if (!metadata) throw utils.asError(errs.sys.invalid, errs.status.args, `no hls for uuid=${uuid}`);

    const metadataObj = JSON.parse(metadata);
    const [contentType, m3u8Body, duration] = m3u8Generator.buildVodM3u8(metadataObj, true, domain);
    console.log(`vod generate m3u8 ok, uuid=${uuid}, duration=${duration}`);

    ctx.type = contentType;
    ctx.body = m3u8Body;
  });

  return router;
}

async function queryMediaInfo(files, fileObjs) {
  const fileIds = files.map(e => {
    return (e.definition && e.taskId && !e.taskObj) ? e.fileId : null;
  }).filter(e => e);

  if (!fileIds?.length) return;

  const secretId = await redis.hget(keys.redis.SRS_TENCENT_CAM, 'secretId');
  const secretKey = await redis.hget(keys.redis.SRS_TENCENT_CAM, 'secretKey');
  const region = await redis.hget(keys.redis.SRS_TENCENT_LH, 'region');
  const vod = new VodClient({
    credential: {secretId, secretKey},
    region,
    profile: {
      httpProfile: {
        endpoint: "vod.tencentcloudapi.com",
      },
    },
  });

  // See https://cloud.tencent.com/document/product/266/33769
  const {MediaInfoSet} = await new Promise((resolve, reject) => {
    vod.DescribeMediaInfos({
      FileIds: fileIds,
      Filters: ['transcodeInfo'],
    }).then(
      (data) => {
        resolve(data);
      },
      (err) => {
        reject(err);
      },
    );
  });

  for (const i in MediaInfoSet) {
    const e = MediaInfoSet[i];
    if (!e.TranscodeInfo?.TranscodeSet?.length) continue;

    const fileObj = fileObjs[e.FileId];
    for (const j in e.TranscodeInfo.TranscodeSet) {
      const f = e.TranscodeInfo.TranscodeSet[j];
      if (f.Definition !== fileObj.definition) continue;

      fileObj.taskObj = {
        url: f.Url,
        bitrate: f.Bitrate,
        height: f.Height,
        width: f.Width,
        size: f.Size,
        duration: f.Duration,
        md5: f.Md5,
      };
      await redis.hset(keys.redis.SRS_VOD_M3U8_METADATA, fileObj.uuid, JSON.stringify(fileObj));
    }
  }
}

