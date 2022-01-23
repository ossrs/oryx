'use strict';

const pkg = require('./package.json');

exports.all = (router) => {
  router.all('/terraform/v1/mgmt/versions', async (ctx) => {
    ctx.body = {
      code: 0,
      data: {
        version: pkg.version
      }
    };
  });
};

