'use strict';

// For components in docker, connect by host.
const config = {
  redis:{
    host: process.env.NODE_ENV === 'development' ? 'localhost' : (process.env.REDIS_HOST || 'mgmt.srs.local'),
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
const fs = require("fs");

exports.handle = (router) => {
  // Query the record patterns.
  router.all('/terraform/v1/hooks/record/query', async (ctx) => {
    const {token} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    const all = await redis.hget(keys.redis.SRS_RECORD_PATTERNS, 'all');
    const home = '/usr/local/srs-cloud/mgmt/containers/data/record';

    console.log(`record query ok, home=${home}, all=${all}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      all: all === 'true',
      home,
    });
  });

  // Setup the record patterns.
  router.all('/terraform/v1/hooks/record/apply', async (ctx) => {
    const {token, all} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    if (all !== true && all !== false) throw utils.asError(errs.sys.invalid, errs.status.args, `invalid all=${all}`);

    const r0 = await redis.hset(keys.redis.SRS_RECORD_PATTERNS, 'all', all);

    console.log(`record apply ok, all=${all}, r0=${r0}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  // List the record files.
  router.all('/terraform/v1/hooks/record/files', async (ctx) => {
    const {token} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    const files = [];
    const [cursor, fileKVs] = await redis.hscan(keys.redis.SRS_RECORD_M3U8_METADATA, 0, '*', 100);
    for (let i = 0; i < fileKVs.length; i += 2) {
      const file = fileKVs[i + 1];
      file && files.push(JSON.parse(file));
    }

    const r0 = files.map(e => {
      return {
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

    console.log(`record files ok, cursor=${cursor}, files=${files.length}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, r0);
  });

  // Generate m3u8 to play.
  const handleM3u8 = async (ctx) => {
    const {uuid} = ctx.params;
    if (!uuid) throw utils.asError(errs.sys.empty, errs.status.args, `no param uuid`);

    const metadata = await redis.hget(keys.redis.SRS_RECORD_M3U8_METADATA, uuid);
    if (!metadata) throw utils.asError(errs.sys.invalid, errs.status.args, `no hls for uuid=${uuid}`);

    const metadataObj = JSON.parse(metadata);
    const [contentType, m3u8Body, duration] = m3u8Generator.buildVodM3u8(
      metadataObj, false, null, true, '/terraform/v1/hooks/record/hls',
    );
    console.log(`record generate m3u8 ok, uuid=${uuid}, duration=${duration}`);

    ctx.type = contentType;
    ctx.body = m3u8Body;
  };
  router.all('/terraform/v1/hooks/record/hls/:uuid.m3u8', handleM3u8);
  router.all('/terraform/v1/hooks/record/hls/:uuid/index.m3u8', handleM3u8);

  // Serve ts to play.
  router.all('/terraform/v1/hooks/record/hls/:dir/:m3u8/:uuid.ts', async (ctx) => {
    const {dir, m3u8, uuid} = ctx.params;
    if (!dir) throw utils.asError(errs.sys.empty, errs.status.args, `no param dir`);
    if (!m3u8) throw utils.asError(errs.sys.empty, errs.status.args, `no param m3u8`);
    if (!uuid) throw utils.asError(errs.sys.empty, errs.status.args, `no param uuid`);

    const tsfile = `${dir}/${m3u8}/${uuid}.ts`;
    if ((!fs.existsSync(tsfile))) {
      throw utils.asError(errs.sys.invalid, errs.status.not, `no ts file ${tsfile}`)
    }

    ctx.type = 'application/vnd.apple.mpegurl';
    ctx.body = fs.readFileSync(tsfile);
  });

  return router;
};

