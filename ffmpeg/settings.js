'use strict';

// For components in docker, connect by host.
const config = {
  redis:{
    host: process.env.NODE_ENV === 'development' ? 'localhost' : 'mgmt.srs.local',
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const errs = require('js-core/errs');
const utils = require('js-core/utils');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const jwt = require('jsonwebtoken');
const keys = require('js-core/keys');

exports.handle = (router) => {
  router.all('/terraform/v1/ffmpeg/forward/secret', async (ctx) => {
    const {token, action, platform, server, secret, enabled} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    const allowedActions = ['update'];
    const allowedPlatforms = ['wx', 'bilibili', 'kuaishou'];
    if (action) {
      if (!allowedActions.includes(action)) {
        return utils.asError(errs.sys.invalid, errs.status.args, `invalid action ${action}`);
      }

      if (!platform) return utils.asError(errs.sys.empty, errs.status.args, 'no platform');
      if (!allowedPlatforms.includes(platform)) {
        return utils.asError(errs.sys.invalid, errs.status.args, `invalid platform ${platform}`);
      }

      if (!server) return utils.asError(errs.sys.empty, errs.status.args, 'no server');
      if (!secret) return utils.asError(errs.sys.empty, errs.status.args, 'no secret');
      if (enabled === undefined) return utils.asError(errs.sys.empty, errs.status.args, 'no enabled');
    }

    let res = null;
    if (action === 'update') {
      const r0 = await redis.hset(keys.redis.SRS_FORWARD_CONFIG, platform, JSON.stringify({
        platform, server, secret, enabled,
      }));
      console.log(`FFmpeg: Forward update secret ok, action=${action}, platform=${platform}, r0=${r0}`);
    } else {
      const configs = await redis.hgetall(keys.redis.SRS_FORWARD_CONFIG);
      for (const k in configs) {
        configs[k] = JSON.parse(configs[k]);
      }
      res = configs;
    }

    console.log(`FFmpeg: Forward secret ok, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, res);
  });

  router.all('/terraform/v1/ffmpeg/forward/streams', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    const configs = await redis.hgetall(keys.redis.SRS_FORWARD_CONFIG);
    const maps = await redis.hgetall(keys.redis.SRS_FORWARD_MAP);
    const frames = await redis.hgetall(keys.redis.SRS_FORWARD_FRAME);

    const res = Object.values(configs).map(e => {
      const conf = JSON.parse(e);

      const stream = maps[conf.platform];
      const frame = stream && frames[`${conf.platform}@${stream}`];
      return {
        platform: conf.platform,
        enabled: conf.enabled,
        stream,
        frame: frame ? JSON.parse(frame) : null,
      };
    });

    console.log(`FFmpeg: Query streams ok, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, res);
  });

  return router;
};

