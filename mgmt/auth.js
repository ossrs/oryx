'use strict';

const utils = require('./utils');

exports.handle = (router) => {
  router.all('/terraform/v1/mgmt/init', async (ctx) => {
    const { MGMT_PASSWORD } = utils.loadConfig();
    ctx.body = utils.asResponse(0, {
      init: !!MGMT_PASSWORD,
    });
  });

  router.all('/terraform/v1/mgmt/login', async (ctx) => {
    const {password} = ctx.request.body;
    if (!password) throw utils.asError(102, 401, 'invalid password');

    const config = utils.loadConfig();
    if (!config.MGMT_PASSWORD) {
      console.log(`init mgmt password ${'*'.repeat(password.length)} ok`);
      utils.saveConfig({...config, MGMT_PASSWORD: password});
      return ctx.body = utils.asResponse(0);
    }

    if (password !== config.MGMT_PASSWORD) throw utils.asError(100, 401, 'invalid password');

    console.log(`login by password ${'*'.repeat(password.length)} ok`);
    ctx.body = utils.asResponse(0);
  });

  return router;
};

