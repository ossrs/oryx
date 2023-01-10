'use strict';

const utils = require("js-core/utils");
const pkg = require('../package.json');
const settings = require('./settings');

exports.handle = (router) => {
  // Mount hooks HTTP APIs.
  settings.handle(router);

  router.all('/terraform/v1/ffmpeg/versions', async (ctx) => {
    ctx.body = utils.asResponse(0, {version: pkg.version});
  });
}

exports.run = () => {
}

