'use strict';

const pkg = require('./package.json');
const utils = require('./utils');

exports.handle = (router) => {
  router.all('/terraform/v1/mgmt/versions', async (ctx) => {
    ctx.body = utils.asResponse(0, {version: pkg.version});
  });

  return router;
};

