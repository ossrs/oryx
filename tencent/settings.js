'use strict';

// For components in docker, connect by host.
const config = {
  redis:{
    host: process.env.NODE_ENV === 'development' ? 'localhost' : 'mgmt.srs.local',
    port: 6379,
    password: '',
  },
};

const errs = require('js-core/errs');
const utils = require('js-core/utils');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const jwt = require('jsonwebtoken');

// The redis key.
const SRS_TENCENT_CAM = 'SRS_TENCENT_CAM';

exports.handle = (router) => {
  // See https://console.cloud.tencent.com/cam
  router.all('/terraform/v1/tencent/cam/secret', async (ctx) => {
    const {token, secretId, secretKey} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    if (!secretId) throw utils.asError(errs.sys.empty, errs.status.args, `no param secretId`);
    if (!secretKey) throw utils.asError(errs.sys.empty, errs.status.args, `no param secretKey`);

    const r0 = await redis.hset(SRS_TENCENT_CAM, 'secretId', secretId);
    const r1 = await redis.hset(SRS_TENCENT_CAM, 'secretKey', secretKey);

    console.log(`tencent cam ok, r0=${r0}, r1=${r1}, secretId=${secretId}, secretKey=${secretKey.length}B, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  return router;
};

