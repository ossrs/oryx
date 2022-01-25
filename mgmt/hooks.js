'use strict';

const utils = require('./utils');
const errs = require('./errs');

exports.handle = (router) => {
  // See https://github.com/ossrs/srs/wiki/v4_EN_HTTPCallback
  router.all('/terraform/v1/mgmt/srs/hooks', async (ctx) => {
    const {action} = ctx.request.body;
    console.log(`srs hooks ok, ${JSON.stringify(ctx.request.body)}`);
    ctx.body = utils.asResponse(0);
  });

  router.all('/terraform/v1/mgmt/srs/secret', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(token);

    console.log(`srs secret ok`);
    ctx.body = utils.asResponse(0);
  });

  return router;
};

