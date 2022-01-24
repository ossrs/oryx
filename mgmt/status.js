'use strict';

const utils = require('./utils');
const pkg = require('./package.json');

exports.handle = (router) => {
  router.all('/terraform/v1/mgmt/status', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(token);

    console.log(`status ok, decoded=${JSON.stringify(decoded)}, token=${'*'.repeat(token.length)}`);
    ctx.body = utils.asResponse(0, {
      version: pkg.version,
    });
  });

  router.all('/terraform/v1/mgmt/software', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(token);

    console.log(`software ok, decoded=${JSON.stringify(decoded)}, token=${'*'.repeat(token.length)}`);
    ctx.body = utils.asResponse(0, {
      version: pkg.version,
    });
  });

  return router;
};

