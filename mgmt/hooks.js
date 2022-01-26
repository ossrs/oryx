'use strict';

const errs = require('./errs');
const utils = require('./utils');
const consts = require('./consts');
const ioredis = require('ioredis');
const redis = utils.redis({config: consts.redis, redis: ioredis});

exports.handle = (router) => {
  // See https://github.com/ossrs/srs/wiki/v4_EN_HTTPCallback
  router.all('/terraform/v1/mgmt/srs/hooks', async (ctx) => {
    const {action, param} = ctx.request.body;
    if (action === 'on_publish') {
      const publish = await redis.get(consts.SRS_SECRET_PUBLISH);
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

    const publish = await redis.get(consts.SRS_SECRET_PUBLISH);
    if (!publish) throw utils.asError(errs.sys.boot, errs.status.sys, `system not boot yet`);

    console.log(`srs secret ok, key=${consts.SRS_SECRET_PUBLISH}, value=${'*'.repeat(publish.length)}, decoded=${JSON.stringify(decoded)}`);
    ctx.body = utils.asResponse(0, {publish});
  });

  return router;
};

