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
const m3u8Generator = require('./m3u8Generator');

exports.handle = (router) => {
  // Setup the DVR patterns.
  router.all('/terraform/v1/hooks/dvr/apply', async (ctx) => {
    const {token, all} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    if (all !== true && all !== false) throw utils.asError(errs.sys.invalid, errs.status.args, `invalid all=${all}`);

    const r0 = await redis.hset(keys.redis.SRS_DVR_PATTERNS, 'all', all);

    console.log(`dvr apply ok, all=${all}, r0=${r0}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  // Query the DVR patterns.
  router.all('/terraform/v1/hooks/dvr/query', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    const all = await redis.hget(keys.redis.SRS_DVR_PATTERNS, 'all');

    const appId = await redis.hget(keys.redis.SRS_TENCENT_CAM, 'appId');
    const secretId = await redis.hget(keys.redis.SRS_TENCENT_CAM, 'secretId');
    const secretKey = await redis.hget(keys.redis.SRS_TENCENT_CAM, 'secretKey');

    console.log(`dvr apply ok, all=${all}, appId=${appId}, secretId=${secretId}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      all: all === 'true',
      secret: !!(appId && secretId && secretKey),
    });
  });

  // List the DVR files.
  router.all('/terraform/v1/hooks/dvr/files', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    const files = [];
    const [cursor, fileKVs] = await redis.hscan(keys.redis.SRS_DVR_M3U8_METADATA, 0, '*', 100);
    for (let i = 0; i < fileKVs.length; i += 2) {
      const file = fileKVs[i + 1];
      file && files.push(JSON.parse(file));
    }

    const r0 = files.map(e => {
      return {
        bucket: e.bucket,
        region: e.region,
        uuid: e.uuid,
        vhost: e.vhost,
        app: e.app,
        stream: e.stream,
        progress: e.progress,
        update: e.update,
        nn: e.files.length,
        duration: e.files.reduce((p, c) => p + (c.duration || 0), 0),
        size: e.files.reduce((p, c) => p + (c.size || 0), 0),
      };
    });

    console.log(`dvr files ok, cursor=${cursor}, files=${files.length}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, r0);
  });

  // Generate m3u8 to play.
  router.all('/terraform/v1/hooks/dvr/hls/:uuid.m3u8', async (ctx) => {
    const {uuid} = ctx.params;
    if (!uuid) throw utils.asError(errs.sys.empty, errs.status.args, `no param uuid`);

    const metadata = await redis.hget(keys.redis.SRS_DVR_M3U8_METADATA, uuid);
    if (!metadata) throw utils.asError(errs.sys.invalid, errs.status.args, `no hls for uuid=${uuid}`);

    const metadataObj = JSON.parse(metadata);
    const [contentType, m3u8Body, duration] = m3u8Generator.buildVodM3u8(metadataObj, true);
    console.log(`dvr generate m3u8 ok, uuid=${uuid}, duration=${duration}`);

    ctx.type = contentType;
    ctx.body = m3u8Body;
  });

  return router;
};

