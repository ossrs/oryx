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
const manager = require('./manager');
const fs = require('fs');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const keys = require('js-core/keys');
const moment = require('moment');

exports.handle = (router) => {
  // TODO: FIXME: Fixed token.
  // See https://github.com/ossrs/srs/wiki/v4_EN_HTTPCallback
  router.all('/terraform/v1/hooks/srs/hls', async (ctx) => {
    const {action, file, duration, m3u8_url, url, seq_no: seqno} = ctx.request.body;

    if (action !== 'on_hls') throw utils.asError(errs.sys.invalid, errs.status.args, `invalid action=${action}`);
    if (!fs.existsSync(file)) throw utils.asError(errs.sys.invalid, errs.status.args, `invalid ts file ${file}`);

    // Create a DVR task if enabled.
    const dvrAll = await redis.hget(keys.redis.SRS_DVR_PATTERNS, 'all');
    let dvr = 'ignore';
    if (dvrAll === 'true') {
      dvr = 'task_created';
      console.log(`create dvr task file=${file}, duration=${duration}, seqno=${seqno}, m3u8_url=${m3u8_url}, url=${url}`);
      manager.postMessage({
        action: 'on_dvr_file', file, duration, seqno, m3u8_url, url, params: ctx.request.body,
      });
    }

    // Create a VoD task if enabled.
    const vodAll = await redis.hget(keys.redis.SRS_VOD_PATTERNS, 'all');
    let vod = 'ignore';
    if (vodAll === 'true') {
      vod = 'task_created';
      console.log(`create vod task file=${file}, duration=${duration}, seqno=${seqno}, m3u8_url=${m3u8_url}, url=${url}`);
      manager.postMessage({
        action: 'on_vod_file', file, duration, seqno, m3u8_url, url, params: ctx.request.body,
      });
    }

    // Create a Record task if enabled.
    const recordAll = await redis.hget(keys.redis.SRS_RECORD_PATTERNS, 'all');
    let record = 'ignore';
    if (recordAll === 'true') {
      record = 'task_created';
      console.log(`create record task file=${file}, duration=${duration}, seqno=${seqno}, m3u8_url=${m3u8_url}, url=${url}`);
      manager.postMessage({
        action: 'on_record_file', file, duration, seqno, m3u8_url, url, params: ctx.request.body,
      });
    }

    const update = moment().format();
    const r0 = await redis.hset(keys.redis.SRS_DVR_PATTERNS, m3u8_url, JSON.stringify({dvr, update}));
    const r1 = await redis.hset(keys.redis.SRS_VOD_PATTERNS, m3u8_url, JSON.stringify({vod, update}));
    const r2 = await redis.hset(keys.redis.SRS_RECORD_PATTERNS, m3u8_url, JSON.stringify({record, update}));

    console.log(`srs hooks ok, dvr=${dvrAll}/${dvr}/${r0}, vod=${vodAll}/${vod}/${r1}, record=${recordAll}/${record}/${r2}, ${JSON.stringify(ctx.request.body)}`);
    ctx.body = utils.asResponse(0);
  });

  return router;
};

