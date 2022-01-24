'use strict';

const utils = require('./utils');
const errs = require('./errs');

exports.handle = (router) => {
  router.all('/terraform/v1/mgmt/init', async (ctx) => {
    const {password} = ctx.request.body;
    if (!process.env.MGMT_PASSWORD && password) {
      console.log(`init mgmt password ${'*'.repeat(password.length)} ok`);
      const config = utils.loadConfig();
      utils.saveConfig({...config, MGMT_PASSWORD: password});
      utils.loadConfig();

      const {expire, expireAt, createAt, token} = utils.createToken();
      console.log(`init password ok, duration=${expire}, create=${createAt}, expire=${expireAt}, password=${'*'.repeat(password.length)}`);
      return ctx.body = utils.asResponse(0, {token, createAt, expireAt});
    }

    ctx.body = utils.asResponse(0, {
      init: !!process.env.MGMT_PASSWORD,
    });
  });

  router.all('/terraform/v1/mgmt/token', async (ctx) => {
    const {token} = ctx.request.body;
    if (!token) throw utils.asError(errs.sys.empty, errs.status.auth, 'no token');

    const decoded = await utils.verifyToken(token);
    const {expire, expireAt, createAt, token2} = utils.createToken();
    console.log(`login by token ok, decoded=${JSON.stringify(decoded)}, duration=${expire}, create=${createAt}, expire=${expireAt}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {token:token2, createAt, expireAt});
  });

  router.all('/terraform/v1/mgmt/login', async (ctx) => {
    if (!process.env.MGMT_PASSWORD) throw utils.asError(errs.auth.init, errs.status.auth, 'not init');

    const {password} = ctx.request.body;
    if (!password) throw utils.asError(errs.sys.empty, errs.status.auth, 'no password');

    if (password !== process.env.MGMT_PASSWORD)
      throw utils.asError(errs.auth.password, errs.status.auth, 'invalid password');

    const {expire, expireAt, createAt, token} = utils.createToken();
    console.log(`login by password ok, duration=${expire}, create=${createAt}, expire=${expireAt}, password=${'*'.repeat(password.length)}`);
    ctx.body = utils.asResponse(0, {token, createAt, expireAt});
  });

  return router;
};

