'use strict';

const moment = require('moment');
const jwt = require('jsonwebtoken');
const utils = require('./utils');
const errs = require('./errs');

function createToken() {
  // Update the user info, @see https://www.npmjs.com/package/jsonwebtoken#usage
  const expire = moment.duration(10, 'seconds');
  const createAt = moment.utc().format(utils.MYSQL_DATETIME);
  const expireAt = moment.utc().add(expire).format(utils.MYSQL_DATETIME);
  const token = jwt.sign(
    {v: 1.0, t: createAt, d: expire},
    process.env.MGMT_PASSWORD, {expiresIn: expire.asSeconds()},
  );

  return {expire, expireAt, createAt, token};
}

async function verifyToken(token) {
  // Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
  return await new Promise((resolve, reject) => {
    jwt.verify(token, process.env.MGMT_PASSWORD, function (err, decoded) {
      if (!err) return resolve(decoded);
      if (err.name === 'TokenExpiredError') throw utils.asError(errs.auth.token, errs.status.auth, `token expired, token=${token}, expiredAt=${err.expiredAt}, ${err.message}`);
      if (err.name === 'JsonWebTokenError') throw utils.asError(errs.auth.token, errs.status.auth, `token invalid, token=${token}, ${err.message}`);
      throw utils.asError(errs.auth.token, errs.status.auth, `token verify failed, ${err.message}`);
    });
  });
}

exports.handle = (router) => {
  router.all('/terraform/v1/mgmt/init', async (ctx) => {
    const {password} = ctx.request.body;
    if (!process.env.MGMT_PASSWORD && password) {
      console.log(`init mgmt password ${'*'.repeat(password.length)} ok`);
      const config = utils.loadConfig();
      utils.saveConfig({...config, MGMT_PASSWORD: password});
      utils.loadConfig();

      const {expire, expireAt, createAt, token} = createToken();
      console.log(`init password ok, duration=${expire}, create=${createAt}, expire=${expireAt}, password=${'*'.repeat(password.length)}`);
      return ctx.body = utils.asResponse(0, {token, createAt, expireAt});
    }

    ctx.body = utils.asResponse(0, {
      init: !!process.env.MGMT_PASSWORD,
    });
  });

  router.all('/terraform/v1/mgmt/token', async (ctx) => {
    if (!process.env.MGMT_PASSWORD) throw utils.asError(errs.auth.init, errs.status.auth, 'not init');

    const {token} = ctx.request.body;
    if (!token) throw utils.asError(errs.sys.empty, errs.status.auth, 'no token');

    if (!process.env.MGMT_PASSWORD) throw utils.asError(errs.auth.token, errs.status.auth, 'invalid token');

    const decoded = await verifyToken(token);
    const {expire, expireAt, createAt, token2} = createToken();
    console.log(`login by token ok, decoded=${JSON.stringify(decoded)}, duration=${expire}, create=${createAt}, expire=${expireAt}, token=${'*'.repeat(token.length)}`);
    ctx.body = utils.asResponse(0, {token:token2, createAt, expireAt});
  });

  router.all('/terraform/v1/mgmt/login', async (ctx) => {
    if (!process.env.MGMT_PASSWORD) throw utils.asError(errs.auth.init, errs.status.auth, 'not init');

    const {password} = ctx.request.body;
    if (!password) throw utils.asError(errs.sys.empty, errs.status.auth, 'no password');

    if (password !== process.env.MGMT_PASSWORD)
      throw utils.asError(errs.auth.password, errs.status.auth, 'invalid password');

    const {expire, expireAt, createAt, token} = createToken();
    console.log(`login by password ok, duration=${expire}, create=${createAt}, expire=${expireAt}, password=${'*'.repeat(password.length)}`);
    ctx.body = utils.asResponse(0, {token, createAt, expireAt});
  });

  return router;
};

