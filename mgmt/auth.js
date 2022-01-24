'use strict';

const moment = require('moment');
const jwt = require('jsonwebtoken');
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

    let config = utils.loadConfig();
    if (config.MGMT_PASSWORD && password !== config.MGMT_PASSWORD) throw utils.asError(100, 401, 'invalid password');

    if (!config.MGMT_PASSWORD) {
      console.log(`init mgmt password ${'*'.repeat(password.length)} ok`);
      config = utils.saveConfig({...config, MGMT_PASSWORD: password});
    }

    // Update the user info, @see https://www.npmjs.com/package/jsonwebtoken#usage
    const expire = moment.duration(10, 'years');
    const createAt = moment.utc().format(utils.MYSQL_DATETIME);
    const expireAt = moment.utc().add(expire).format(utils.MYSQL_DATETIME);
    const token = jwt.sign(
      {v: 1.0, t: createAt, d: expire},
      config.MGMT_PASSWORD, {expiresIn: expire.asSeconds()},
    );

    console.log(`login by password ${'*'.repeat(password.length)} ok, duration=${expire}, create=${createAt}, expire=${expireAt}`);
    ctx.body = utils.asResponse(0, {token, createAt, expireAt});
  });

  return router;
};

