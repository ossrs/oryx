'use strict';

// For components in docker, connect by host.
const config = {
  redis:{
    host: process.env.NODE_ENV === 'development' ? 'localhost' : 'mgmt.srs.local',
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const utils = require('js-core/utils');
const errs = require('js-core/errs');
const jwt = require('jsonwebtoken');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const moment = require('moment');
const keys = require('js-core/keys');
const helper = require('./helper');
const metadata = require('js-core/metadata');
const axios = require('axios');

exports.handle = (router) => {
  router.all('/terraform/v1/mgmt/upgrade', async (ctx) => {
    const {token} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    const uwStart = await redis.hget(keys.redis.SRS_UPGRADE_WINDOW, 'start');
    const uwDuration = await redis.hget(keys.redis.SRS_UPGRADE_WINDOW, 'duration');
    const inUpgradeWindow = helper.inUpgradeWindow(uwStart, uwDuration, moment());

    const {version, latest} = await helper.execApi('refreshVersion');
    const target = latest || 'lighthouse';
    const upgradingMessage = `upgrade to target=${target}, current=${version}, latest=${latest}, window=${inUpgradeWindow}`;
    console.log(`Start ${upgradingMessage}`);

    const r0 = await redis.hget(keys.redis.SRS_UPGRADING, 'upgrading');
    if (r0 === "1") {
      const r1 = await redis.hget(keys.redis.SRS_UPGRADING, 'desc');
      throw utils.asError(errs.sys.upgrading, errs.status.sys, `already upgrading ${r0}, ${r1}`);
    }

    // Set the upgrading to avoid others.
    await redis.hset(keys.redis.SRS_UPGRADING, 'upgrading', 1);
    await redis.hset(keys.redis.SRS_UPGRADING, 'desc', `${upgradingMessage}`);

    try {
      await helper.execApi('execUpgrade', [target]);
    } finally {
      await redis.hset(keys.redis.SRS_UPGRADING, 'upgrading', 0);
    }

    console.log(`upgrade ok, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      version,
    });
  });

  router.all('/terraform/v1/mgmt/strategy', async (ctx) => {
    const {token} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    const r3 = await helper.execApi('refreshVersion');

    const upgrading = await redis.hget(keys.redis.SRS_UPGRADING, 'upgrading');
    const r0 = await redis.hget(keys.redis.SRS_UPGRADE_STRATEGY, 'strategy');
    const strategy = r0 || 'auto';
    const newStrategy = strategy === 'auto' ? 'manual' : 'auto';
    const r1 = await redis.hset(keys.redis.SRS_UPGRADE_STRATEGY, 'strategy', newStrategy);
    const r2 = await redis.hset(keys.redis.SRS_UPGRADE_STRATEGY, 'desc', `${moment().format()} changed, upgrading=${upgrading}, r0=${r0}/${strategy}, r1=${r1}/${newStrategy}`);
    console.log(`status ok, upgrading=${upgrading}, r0=${r0}/${strategy}, r1=${r1}/${newStrategy}, r2=${r2}, releases=${JSON.stringify(r3)}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  router.all('/terraform/v1/mgmt/ssl', async (ctx) => {
    const {token, key, crt} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    if (!key) throw utils.asError(errs.sys.empty, errs.status.args, 'no key');
    if (!crt) throw utils.asError(errs.sys.empty, errs.status.args, 'no crt');
    await helper.execApi('updateSslFile', [key, crt]);

    // Setup the HTTPS information.
    await redis.set(keys.redis.SRS_HTTPS, 'ssl');

    console.log(`ssl ok, key=${key.length}B, crt=${crt.length}B, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  router.all('/terraform/v1/mgmt/letsencrypt', async (ctx) => {
    const {token, domain} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    if (!domain) throw utils.asError(errs.sys.empty, errs.status.args, 'no domain');
    await helper.execApi('updateLetsEncrypt', [domain]);

    // Setup the HTTPS information.
    await redis.set(keys.redis.SRS_HTTPS, 'lets');

    const keyFile = `${process.cwd()}/containers/etc/letsencrypt/live/${domain}/privkey.pem`;
    const crtFile = `${process.cwd()}/containers/etc/letsencrypt/live/${domain}/cert.pem`;

    console.log(`let's encrypt ok, domain=${domain}, key=${keyFile}, crt=${crtFile}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  router.all('/terraform/v1/mgmt/pubkey', async (ctx) => {
    const {token, enabled} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    const enabledValue = enabled ? 'enable' : 'disable';
    await helper.execApi('accessSsh', [enabledValue]);

    console.log(`pubkey ok, enable=${enabled}/${enabledValue}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  router.all('/terraform/v1/mgmt/containers', async (ctx) => {
    const {token, action, name, enabled} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    if (!action) throw utils.asError(errs.sys.empty, errs.status.args, `no param action`);

    const validActions = ['query', 'enabled', 'switch'];
    if (!validActions.includes(action)) throw utils.asError(errs.sys.invalid, errs.status.args, `invalid action ${action}, should be ${validActions}`);

    if (action === 'enabled') {
      if (!name) throw utils.asError(errs.sys.empty, errs.status.args, `no param name`);
      if (enabled !== true && enabled !== false) throw utils.asError(errs.sys.empty, errs.status.args, `no param enabled`);

      const r0 = await redis.hset(keys.redis.SRS_CONTAINER_DISABLED, name, !enabled);
      if (!enabled && name) await helper.execApi('rmContainer', [name]);

      console.log(`srs ok, action=${action}, name=${name}, enabled=${enabled}, r0=${r0}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
      return ctx.body = utils.asResponse(0);
    } else if (action === 'switch') {
      const market = metadata.market;
      if (!name) throw utils.asError(errs.sys.empty, errs.status.args, `no param name`);
      if (name !== market.srs.name && name !== market.srsDev.name) {
        throw utils.asError(errs.sys.invalid, errs.status.args, `invalid name ${name}`);
      }

      const disable = (name === market.srsDev.name) ? market.srs.name : market.srsDev.name;
      const r0 = await redis.hset(keys.redis.SRS_CONTAINER_DISABLED, name, false);
      const r1 = await redis.hset(keys.redis.SRS_CONTAINER_DISABLED, disable, true);
      await helper.execApi('rmContainer', [disable]);

      console.log(`switch ok, action=${action}, name=${name}/${r0}, disable=${disable}/${r1}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
      return ctx.body = utils.asResponse(0);
    }

    // Query containers
    const {containers} = await helper.execApi('fetchContainers', [name]);

    console.log(`srs ok, action=${action}, name=${name}, containers=${containers.length}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, containers);
  });

  router.all('/terraform/v1/mgmt/bilibili', async (ctx) => {
    const {token, bvid} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    if (!bvid) throw utils.asError(errs.sys.empty, errs.status.args, 'no bvid');

    const bilibili = await redis.hget(keys.redis.SRS_CACHE_BILIBILI, bvid);
    const bilibiliObj = bilibili ? JSON.parse(bilibili) : {};

    let cacheExpired = false;
    if (bilibiliObj.update) {
      const expired = moment(bilibiliObj.update).add(process.env.NODE_ENV === 'development' ? 300 : 3 * 3600, 's');
      if (expired.isBefore(moment())) cacheExpired = true;
    }

    if (!bilibiliObj.res || cacheExpired) {
      bilibiliObj.update = moment().format();
      bilibiliObj.res = await new Promise((resolve, reject) => {
        axios.get(`https://api.bilibili.com/x/web-interface/view?bvid=${bvid}`).then(res => {
          resolve(res.data.data);
        }).catch(err => reject);
      });

      await redis.hset(keys.redis.SRS_CACHE_BILIBILI, bvid, JSON.stringify(bilibiliObj));
      console.log(`bilibili cache bvid=${bvid}, update=${bilibiliObj.update}`);
    }

    console.log(`bilibili query ok, bvid=${bvid}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, bilibiliObj.res);
  });

  router.all('/terraform/v1/mgmt/beian/query', async (ctx) => {
    const icp = await redis.hget(keys.redis.SRS_BEIAN, 'icp');

    console.log(`beian: query ok, miit=${JSON.stringify(icp)}`);
    ctx.body = utils.asResponse(0, {icp});
  });

  router.all('/terraform/v1/mgmt/beian/update', async (ctx) => {
    const {token, beian, text} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    if (!beian) throw utils.asError(errs.sys.empty, errs.status.args, 'no beian');
    if (!text) throw utils.asError(errs.sys.empty, errs.status.args, 'no text');

    const r0 = await redis.hset(keys.redis.SRS_BEIAN, beian, text);
    console.log(`beian: update ok, beian=${beian}, text=${text}, r0=${JSON.stringify(r0)}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, r0);
  });
  router.all('/terraform/v1/mgmt/window/query', async (ctx) => {
    const {token} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    const start = await redis.hget(keys.redis.SRS_UPGRADE_WINDOW, 'start');
    const duration = await redis.hget(keys.redis.SRS_UPGRADE_WINDOW, 'duration');

    console.log(`query upgrade window ok, start=${start}, duration=${duration}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      start: parseInt(start),
      duration: parseInt(duration),
    });
  });

  router.all('/terraform/v1/mgmt/window/update', async (ctx) => {
    const {token, start, duration} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    if (isNaN(parseInt(start))) throw utils.asError(errs.sys.empty, errs.status.args, `no param start`);
    if (!duration) throw utils.asError(errs.sys.empty, errs.status.args, `no param end`);
    if (duration <= 3) return utils.asError(errs.sys.invalid, errs.status.args, `window should greater than 3 hours`);
    if (duration > 24) return utils.asError(errs.sys.invalid, errs.status.args, `window should smaller than 24 hours`);

    const r0 = await redis.hset(keys.redis.SRS_UPGRADE_WINDOW, 'start', start);
    const r1 = await redis.hset(keys.redis.SRS_UPGRADE_WINDOW, 'duration', duration);
    const r2 = await redis.hset(keys.redis.SRS_UPGRADE_WINDOW, 'update', moment().format());

    console.log(`update upgrade window ok, start=${start}, duration=${duration}, r0=${r0}, r1=${r1}, r2=${r2}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  return router;
};

