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
const SRS_SECRET_PUBLISH = 'SRS_SECRET_PUBLISH';

exports.handle = (router) => {
  // Init the secrets.
  init();

  // See https://github.com/ossrs/srs/wiki/v4_EN_HTTPCallback
  router.all('/terraform/v1/hooks/srs/verify', async (ctx) => {
    const {action, param} = ctx.request.body;
    if (action === 'on_publish') {
      const publish = await redis.get(SRS_SECRET_PUBLISH);
      if (publish && param.indexOf(publish) === -1) {
        throw utils.asError(errs.srs.verify, errs.status.auth, `invalid params=${param} action=${action}`);
      }
    }

    console.log(`srs hooks ok, ${JSON.stringify(ctx.request.body)}`);
    ctx.body = utils.asResponse(0);
  });

  router.all('/terraform/v1/hooks/srs/secret', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    const publish = await redis.get(SRS_SECRET_PUBLISH);
    if (!publish) throw utils.asError(errs.sys.boot, errs.status.sys, `system not boot yet`);

    console.log(`srs secret ok, key=${SRS_SECRET_PUBLISH}, value=${'*'.repeat(publish.length)}, decoded=${JSON.stringify(decoded)}`);
    ctx.body = utils.asResponse(0, {publish});
  });

  router.all('/terraform/v1/hooks/srs/secret/update', async (ctx) => {
    const { token, secret} = ctx.request.body;
    await utils.verifyToken(jwt, token);

    if (!secret) throw utils.asError(errs.sys.empty, errs.status.args, 'no secret');

    const r0 = await redis.set(SRS_SECRET_PUBLISH, secret);

    console.log(`hooks update secret, key=${SRS_SECRET_PUBLISH}, value=${'*'.repeat(secret.length)}, r0=${r0}`);
    ctx.body = utils.asResponse(0, {});
  });

  return router;
};

async function init() {
  // Setup the publish secret for first run.
  let publish = await redis.get(SRS_SECRET_PUBLISH);
  if (!publish) {
    publish = Math.random().toString(16).slice(-8);
    const r0 = await redis.set(SRS_SECRET_PUBLISH, publish);
    console.log(`hooks create secret, key=${SRS_SECRET_PUBLISH}, value=${'*'.repeat(publish.length)}, r0=${r0}`);
  }
}

