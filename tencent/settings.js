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
const COS = require('cos-nodejs-sdk-v5');
const cos = require('js-core/cos');
const sdkInternal = require('./sdk-internal/common/abstract_client');

exports.handle = (router) => {
  // See https://console.cloud.tencent.com/cam
  router.all('/terraform/v1/tencent/cam/secret', async (ctx) => {
    const {token, secretId, secretKey} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    if (!secretId) throw utils.asError(errs.sys.empty, errs.status.args, `no param secretId`);
    if (!secretKey) throw utils.asError(errs.sys.empty, errs.status.args, `no param secretKey`);

    const res = await apiCloud({
      secretId,
      secretKey,
      endpoint: 'cam.tencentcloudapi.com',
      version: '2019-01-16',
      action: 'GetUserAppId',
      req: {},
    });
    const {AppId: appId, OwnerUin: uin} = res;
    if (!appId) throw utils.asError(errs.sys.auth, errs.status.args, `query appId failed`);

    const r0 = await redis.hset(keys.redis.SRS_TENCENT_CAM, 'appId', appId);
    const r1 = await redis.hset(keys.redis.SRS_TENCENT_CAM, 'secretId', secretId);
    const r2 = await redis.hset(keys.redis.SRS_TENCENT_CAM, 'secretKey', secretKey);
    const r3 = await redis.hset(keys.redis.SRS_TENCENT_CAM, 'uin', uin);

    const region = await redis.hget(keys.redis.SRS_TENCENT_LH, 'region');
    await cos.createCosBucket(redis, COS, region);

    console.log(`CAM: Update ok, appId=${appId}, uin=${uin}, secretId=${secretId}, secretKey=${secretKey.length}B, r0=${r0}, r1=${r1}, r2=${r2}, r3=${r3}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  return router;
};

async function apiCloud({secretId, secretKey, endpoint, region, version, action, req}) {
  const profile = {
    credential: {
      secretId: secretId,
      secretKey: secretKey,
    },
    profile: {
      httpProfile: {
        endpoint: endpoint,
      },
    },
  };
  if (region) profile.region = region;

  const sdk = new sdkInternal.AbstractClient(endpoint, version, profile);
  return await sdk.request(action, req);
}

