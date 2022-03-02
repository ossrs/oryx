'use strict';

// For mgmt, it's ok to connect to localhost.
const config = {
  redis:{
    host: 'localhost',
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const os = require('os');
const fs = require('fs');
const dotenv = require('dotenv');
const utils = require('js-core/utils');
const errs = require('js-core/errs');
const moment = require('moment');
const jwt = require('jsonwebtoken');
const util = require('util');
const exec = util.promisify(require('child_process').exec);
const metadata = require('./metadata');
const market = require('./market');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const keys = require('js-core/keys');
const consts = require('./consts');

// MySQL日期字段格式化字符串 @see https://stackoverflow.com/a/27381633
const MYSQL_DATETIME = 'YYYY-MM-DD HH:mm:ss';

function loadConfig() {
  dotenv.config({path: '.env', override: true});
  return {
    MGMT_PASSWORD: process.env.MGMT_PASSWORD,
    REDIS_PORT: process.env.REDIS_PORT,
    REDIS_PASSWORD: process.env.REDIS_PASSWORD,
  };
}

function saveConfig(config) {
  const envVars = Object.keys(config).map(k => {
    const v = config[k];
    return v ? `${k}=${v}` : '';
  }).filter(e => e);

  // Append an empty line.
  envVars.push('');

  fs.writeFileSync('.env', envVars.join(os.EOL));
  return config;
}

function createToken(moment, jwt) {
  // Update the user info, @see https://www.npmjs.com/package/jsonwebtoken#usage
  const expire = moment.duration(1, 'years');
  const createAt = moment.utc().format(MYSQL_DATETIME);
  const expireAt = moment.utc().add(expire).format(MYSQL_DATETIME);
  const token = jwt.sign(
    {v: 1.0, t: createAt, d: expire},
    process.env.MGMT_PASSWORD, {expiresIn: expire.asSeconds()},
  );

  return {expire, expireAt, createAt, token};
}

exports.handle = (router) => {
  router.all('/terraform/v1/mgmt/init', async (ctx) => {
    const {password} = ctx.request.body;
    if (!process.env.MGMT_PASSWORD && password) {
      console.log(`init mgmt password ${'*'.repeat(password.length)} ok`);
      const config = loadConfig();
      saveConfig({...config, MGMT_PASSWORD: password});
      loadConfig();

      const [allHooks, runningHooks] = await market.queryContainer(metadata.market.hooks.name);
      const [allTencent, runningTencent] = await market.queryContainer(metadata.market.tencent.name);

      try {
        await exec(`docker rm -f ${metadata.market.hooks.name}`);
      } catch (e) {
      }

      try {
        await exec(`docker rm -f ${metadata.market.tencent.name}`);
      } catch (e) {
      }

      if (allHooks?.ID && runningHooks?.ID) {
        // We must restart the hooks, which depends on the .env
        for (let i = 0; i < 60; i++) {
          // Wait util running and got another container ID.
          const [all, running] = await market.queryContainer(metadata.market.hooks.name);
          // Please note that we don't update the metadata of SRS, client must request the updated status.
          if (all && all.ID && running && running.ID && running.ID !== runningHooks.ID) break;
          await new Promise(resolve => setTimeout(resolve, 500));
        }
        console.log(`restart ${metadata.market.hooks.name} ${metadata.market.hooks.container.ID} when .env updated`);
      }

      if (allTencent?.ID && runningTencent?.ID) {
        // We must restart the tencent, which depends on the .env
        for (let i = 0; i < 60; i++) {
          // Wait util running and got another container ID.
          const [all, running] = await market.queryContainer(metadata.market.tencent.name);
          // Please note that we don't update the metadata of SRS, client must request the updated status.
          if (all && all.ID && running && running.ID && running.ID !== runningTencent.ID) break;
          await new Promise(resolve => setTimeout(resolve, 500));
        }
        console.log(`restart ${metadata.market.tencent.name} ${metadata.market.tencent.container.ID} when .env updated`);
      }

      const {expire, expireAt, createAt, token} = createToken(moment, jwt);
      console.log(`init password ok, duration=${expire}, create=${createAt}, expire=${expireAt}, password=${'*'.repeat(password.length)}`);
      return ctx.body = utils.asResponse(0, {
        token,
        createAt,
        expireAt,
      });
    }

    ctx.body = utils.asResponse(0, {
      init: !!process.env.MGMT_PASSWORD,
    });
  });

  router.all('/terraform/v1/mgmt/check', async (ctx) => {
    // Check whether redis is ok.
    const r0 = await redis.get(keys.redis.SRS_SECRET_PUBLISH);
    const r1 = await redis.hlen(keys.redis.SRS_FIRST_BOOT);
    const r2 = await redis.hlen(keys.redis.SRS_TENCENT_LH);
    if (!r0 || !r1 || !r2) throw utils.asError(errs.sys.redis, errs.status.sys, `redis corrupt`);

    const upgrading = await redis.hget(consts.SRS_UPGRADING, 'upgrading');

    console.log(`system check ok, r0=${r0}, r1=${r1}, r2=${r2}`);
    ctx.body = utils.asResponse(0, {
      upgrading: upgrading === "1",
    });
  });

  router.all('/terraform/v1/mgmt/token', async (ctx) => {
    const {token} = ctx.request.body;
    if (!token) throw utils.asError(errs.sys.empty, errs.status.auth, 'no token');

    const decoded = await utils.verifyToken(jwt, token);
    const {expire, expireAt, createAt, token2} = createToken(moment, jwt);
    console.log(`login by token ok, decoded=${JSON.stringify(decoded)}, duration=${expire}, create=${createAt}, expire=${expireAt}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {token:token2, createAt, expireAt});
  });

  router.all('/terraform/v1/mgmt/login', async (ctx) => {
    if (!process.env.MGMT_PASSWORD) throw utils.asError(errs.auth.init, errs.status.auth, 'not init');

    const {password} = ctx.request.body;
    if (!password) throw utils.asError(errs.sys.empty, errs.status.auth, 'no password');

    if (password !== process.env.MGMT_PASSWORD)
      throw utils.asError(errs.auth.password, errs.status.auth, 'invalid password');

    const {expire, expireAt, createAt, token} = createToken(moment, jwt);
    console.log(`login by password ok, duration=${expire}, create=${createAt}, expire=${expireAt}, password=${'*'.repeat(password.length)}`);
    ctx.body = utils.asResponse(0, {token, createAt, expireAt});
  });

  return router;
};

