'use strict';

const config = {
  redis: {
    host: 'localhost',
    port: 6379,
    password: '',
  }
};

const errs = require('./errs');
const utils = require('./utils');
const ioredis = require('ioredis');
const redis = utils.redis({config: config.redis, redis: ioredis});

exports.handle = (router) => {
  // See https://github.com/ossrs/srs/wiki/v4_EN_HTTPCallback
  router.all('/terraform/v1/mgmt/srs/hooks', async (ctx) => {
    const {action, param} = ctx.request.body;
    if (action === 'on_publish') {
      const publish = await redis.get(utils.SRS_SECRET_PUBLISH);
      if (publish && param.indexOf(publish) === -1) {
        throw utils.asError(errs.srs.verify, errs.status.auth, `invalid params=${param} action=${action}`);
      }
    }

    console.log(`srs hooks ok, ${JSON.stringify(ctx.request.body)}`);
    ctx.body = utils.asResponse(0);
  });

  router.all('/terraform/v1/mgmt/srs/secret', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(token);

    let publish = await redis.get(utils.SRS_SECRET_PUBLISH);
    if (!publish) {
      publish = Math.random().toString(16).slice(-8);
      const r0 = await redis.set(utils.SRS_SECRET_PUBLISH, publish);
      console.log(`srs secret create ok, key=${utils.SRS_SECRET_PUBLISH}, r0=${r0}`);
    }

    console.log(`srs secret ok, key=${utils.SRS_SECRET_PUBLISH}, value=${'*'.repeat(publish.length)}, decoded=${JSON.stringify(decoded)}`);
    ctx.body = utils.asResponse(0, {publish});
  });

  return router;
};

