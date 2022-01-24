'use strict';

const moment = require('moment');
const jwt = require('jsonwebtoken');
const utils = require('./utils');
const errs = require('./errs');

function createToken(config) {
  // Update the user info, @see https://www.npmjs.com/package/jsonwebtoken#usage
  const expire = moment.duration(10, 'seconds');
  const createAt = moment.utc().format(utils.MYSQL_DATETIME);
  const expireAt = moment.utc().add(expire).format(utils.MYSQL_DATETIME);
  const token = jwt.sign(
    {v: 1.0, t: createAt, d: expire},
    config.MGMT_PASSWORD, {expiresIn: expire.asSeconds()},
  );

  return {expire, expireAt, createAt, token};
}

exports.handle = (router) => {
  router.all('/terraform/v1/mgmt/init', async (ctx) => {
    const { MGMT_PASSWORD } = utils.loadConfig();
    ctx.body = utils.asResponse(0, {
      init: !!MGMT_PASSWORD,
    });
  });

  router.all('/terraform/v1/mgmt/token', async (ctx) => {
    const {token} = ctx.request.body;
    if (!token) throw utils.asError(errs.sys.empty, errs.status.auth, 'no token');

    const config = utils.loadConfig();
    if (!config.MGMT_PASSWORD) throw utils.asError(errs.auth.token, errs.status.auth, 'invalid token');

    // Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
    const decoded = await new Promise((resolve, reject) => {
      jwt.verify(token, config.MGMT_PASSWORD, function (err, decoded) {
        if (!err) return resolve(decoded);
        if (err.name === 'TokenExpiredError') throw utils.asError(errs.auth.token, errs.status.auth, `token expired, token=${token}, expiredAt=${err.expiredAt}, ${err.message}`);
        if (err.name === 'JsonWebTokenError') throw utils.asError(errs.auth.token, errs.status.auth, `token invalid, token=${token}, ${err.message}`);
        throw utils.asError(errs.auth.token, errs.status.auth, `token verify failed, ${err.message}`);
      });
    });

    const {expire, expireAt, createAt, token2} = createToken(config);
    console.log(`login by token ok, decoded=${JSON.stringify(decoded)}, duration=${expire}, create=${createAt}, expire=${expireAt}, token=${'*'.repeat(token.length)}`);
    ctx.body = utils.asResponse(0, {token:token2, createAt, expireAt});
  });

  router.all('/terraform/v1/mgmt/login', async (ctx) => {
    const {password} = ctx.request.body;
    if (!password) throw utils.asError(errs.sys.empty, errs.status.auth, 'no password');

    let config = utils.loadConfig();
    if (config.MGMT_PASSWORD && password !== config.MGMT_PASSWORD)
      throw utils.asError(errs.auth.password, errs.status.auth, 'invalid password');

    if (!config.MGMT_PASSWORD) {
      console.log(`init mgmt password ${'*'.repeat(password.length)} ok`);
      config = utils.saveConfig({...config, MGMT_PASSWORD: password});
    }

    const {expire, expireAt, createAt, token} = createToken(config);
    console.log(`login by password ok, duration=${expire}, create=${createAt}, expire=${expireAt}, password=${'*'.repeat(password.length)}`);
    ctx.body = utils.asResponse(0, {token, createAt, expireAt});
  });

  return router;
};

