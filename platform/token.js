'use strict';

// For components in docker, connect by host.
const config = {
  redis:{
    host: process.env.NODE_ENV === 'development' ? 'localhost' : 'mgmt.srs.local',
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const os = require('os');
const fs = require('fs');
const errs = require('js-core/errs');
const utils = require('js-core/utils');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const jwt = require('jsonwebtoken');
const keys = require('js-core/keys');
const moment = require('moment');
const dotenv = require('dotenv');
const path = require('path');
const helper = require('./helper');
const metadata = require('./metadata');

exports.handle = (router) => {
  router.all('/terraform/v1/mgmt/init', async (ctx) => {
    const {password} = ctx.request.body;

    // If no password, query the system init status.
    if (!password) {
      return ctx.body = utils.asResponse(0, {
        init: !!process.env.MGMT_PASSWORD,
      });
    }

    // If already initialized, never set it again.
    if (process.env.MGMT_PASSWORD) {
      throw utils.asError(errs.sys.invalid, errs.status.args, 'already initialized');
    }

    // Initialize the system password, save to env.
    console.log(`init mgmt password ${'*'.repeat(password.length)} ok`);
    const filename = process.env.NODE_ENV === 'development' ? path.join('..', 'mgmt', '.env') : '.env';
    utils.saveEnvs(fs, os, dotenv, filename, {MGMT_PASSWORD: password});

    // Refresh the local token.
    await helper.execApi('reloadEnv');
    utils.reloadEnv(dotenv, fs, path);

    // Note that we use api secret to generate token, not the user password.
    const apiSecret = await utils.apiSecret(redis);
    const {expire, expireAt, createAt, nonce, token} = utils.createToken(moment, jwt, apiSecret);

    console.log(`init password ok, duration=${expire}, create=${createAt}, expire=${expireAt}, nonce=${nonce}, password=${'*'.repeat(password.length)}`);
    return ctx.body = utils.asResponse(0, {
      token,
      createAt,
      expireAt,
    });
  });

  router.all('/terraform/v1/mgmt/check', async (ctx) => {
    // Check whether redis is ok.
    const r0 = await redis.hget(keys.redis.SRS_AUTH_SECRET, 'pubSecret');
    const r1 = await redis.hlen(keys.redis.SRS_FIRST_BOOT);
    const r2 = await redis.hlen(keys.redis.SRS_TENCENT_LH);
    if (!r0 || !r1 || !r2) throw utils.asError(errs.sys.redis, errs.status.sys, `redis corrupt`);

    const upgrading = await redis.hget(keys.redis.SRS_UPGRADING, 'upgrading');

    console.log(`system check ok, r0=${r0}, r1=${r1}, r2=${r2}`);
    ctx.body = utils.asResponse(0, {
      upgrading: upgrading === "1",
    });
  });

  router.all('/terraform/v1/mgmt/token', async (ctx) => {
    const {token} = ctx.request.body;
    if (!token) throw utils.asError(errs.sys.empty, errs.status.auth, 'no token');

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    await new Promise(resolve => setTimeout(resolve, 1000));

    const {expire, expireAt, createAt, nonce, token2} = utils.createToken(moment, jwt, apiSecret);
    console.log(`login by token ok, decoded=${JSON.stringify(decoded)}, duration=${expire}, create=${createAt}, expire=${expireAt}, nonce=${nonce}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {token:token2, createAt, expireAt});
  });

  router.all('/terraform/v1/mgmt/secret/token', async (ctx) => {
    const {apiSecret} = ctx.request.body;

    const apiSecret2 = await utils.apiSecret(redis);

    if (apiSecret !== apiSecret2) {
      throw utils.asError(errs.auth.password, errs.status.auth, 'apiSecret verify failed');
    }

    const { expire, token } = utils.createToken(moment, jwt, apiSecret);
    const expiresIn = expire.asSeconds()
    console.log(`create token by apiSecret ok, token=${token.length}B, expiresIn=${expiresIn}`);
    ctx.body = utils.asResponse(0, {token, expiresIn});
  });

  router.all('/terraform/v1/mgmt/login', async (ctx) => {
    if (!process.env.MGMT_PASSWORD) throw utils.asError(errs.auth.init, errs.status.auth, 'not init');

    const {password} = ctx.request.body;
    if (!password) throw utils.asError(errs.sys.empty, errs.status.auth, 'no password');

    if (password !== process.env.MGMT_PASSWORD) {
      throw utils.asError(errs.auth.password, errs.status.auth, 'invalid password');
    }

    await new Promise(resolve => setTimeout(resolve, 1000));

    const apiSecret = await utils.apiSecret(redis);
    const {expire, expireAt, createAt, nonce, token} = utils.createToken(moment, jwt, apiSecret);
    console.log(`login by password ok, duration=${expire}, create=${createAt}, expire=${expireAt}, nonce=${nonce}, password=${'*'.repeat(password.length)}`);
    ctx.body = utils.asResponse(0, {token, createAt, expireAt});
  });

  router.all('/terraform/v1/mgmt/status', async (ctx) => {
    const {token} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    const {version} = await helper.execApi('queryVersion');
    const {stable, latest} = metadata.upgrade.releases;

    const upgrading = await redis.hget(keys.redis.SRS_UPGRADING, 'upgrading');
    const r0 = await redis.hget(keys.redis.SRS_UPGRADE_STRATEGY, 'strategy');
    const strategy = r0 || 'auto';

    console.log(`status ok, upgrading=${upgrading}, strategy=${strategy}, version=${version}, stable=${stable}, latest=${latest}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      version,
      releases: {
        stable,
        latest,
      },
      upgrading: upgrading === "1",
      strategy,
    });
  });

  return router;
};

