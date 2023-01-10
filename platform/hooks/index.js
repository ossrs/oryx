'use strict';

const utils = require("js-core/utils");
const pkg = require('../package.json');
const hooks = require('./hooks');
const hls = require('./hls');
const record = require('./record');
const dvr = require('./dvr');
const vod = require('./vod');
const manager = require('./manager');

exports.handle = (router) => {
  // Mount hooks HTTP APIs.
  record.handle(dvr.handle(vod.handle(hls.handle(hooks.handle(router)))));

  router.all('/terraform/v1/hooks/versions', async (ctx) => {
    ctx.body = utils.asResponse(0, {version: pkg.version});
  });
}

exports.run = () => {
  manager.run();
}
