'use strict';

// For components in docker, connect by host.
const config = {
  redis:{
    host: process.env.NODE_ENV === 'development' ? 'localhost' : (process.env.REDIS_HOST || 'mgmt.srs.local'),
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const errs = require('js-core/errs');
const utils = require('js-core/utils');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const jwt = require('jsonwebtoken');
const keys = require('js-core/keys');
const { v4: uuidv4 } = require('uuid');

exports.handle = (router) => {
  // Init the secrets.
  init();

  // See https://github.com/ossrs/srs/wiki/v4_EN_HTTPCallback
  router.all('/terraform/v1/hooks/srs/verify', async (ctx) => {
    const noAuth = await redis.hget(keys.redis.SRS_AUTH_SECRET, 'pubNoAuth');
    if (noAuth === 'true') {
      ctx.body = utils.asResponse(0);
      return console.log(`srs hooks disabled, ${JSON.stringify(ctx.request.body)}`);
    }

    const {action, param, vhost, app, stream, server_id, client_id} = ctx.request.body;
    if (action === 'on_publish') {
      const publish = await redis.hget(keys.redis.SRS_AUTH_SECRET, 'pubSecret');

      // Note that we allow pass secret by params or in stream name, for example, some encoder does not support params
      // with ?secret=xxx, so it will fail when url is:
      //      rtmp://ip/live/livestream?secret=xxx
      // so user could change the url to bellow to get around of it:
      //      rtmp://ip/live/livestreamxxx
      // or simply use secret as stream:
      //      rtmp://ip/live/xxx
      // in this situation, the secret is part of stream name.
      if (publish && param.indexOf(publish) === -1 && stream.indexOf(publish) === -1) {
        throw utils.asError(errs.srs.verify, errs.status.auth, `invalid params=${param} action=${action}`);
      }
    }

    let active = null;
    const srt = param && param.indexOf('upstream=srt') >= 0;
    const rtc = param && param.indexOf('upstream=rtc') >= 0;
    const url = vhost && app && stream && utils.streamURL(vhost, app, stream);
    if (action === 'on_publish') {
      const streamObj = {vhost, app, stream, server: server_id, client: client_id,};
      active = await redis.hset(keys.redis.SRS_STREAM_ACTIVE, url, JSON.stringify(streamObj));
      await redis.hincrby(keys.redis.SRS_STAT_COUNTER, 'publish', 1);
      if (srt) await redis.hset(keys.redis.SRS_STREAM_SRT_ACTIVE, url, JSON.stringify(streamObj));
      if (rtc) await redis.hset(keys.redis.SRS_STREAM_RTC_ACTIVE, url, JSON.stringify(streamObj));
    } else if (action === 'on_unpublish') {
      active = await redis.hdel(keys.redis.SRS_STREAM_ACTIVE, url);
      await redis.hdel(keys.redis.SRS_STREAM_SRT_ACTIVE, url);
      await redis.hdel(keys.redis.SRS_STREAM_RTC_ACTIVE, url);
    } else if (action === 'on_play') {
      await redis.hincrby(keys.redis.SRS_STAT_COUNTER, 'play', 1);
    }

    console.log(`srs hooks ok, action=${action}, active=${active}, srt=${srt}, rtc=${rtc}, url=${url}, ${JSON.stringify(ctx.request.body)}`);
    ctx.body = utils.asResponse(0);
  });

  const handleSecretQuery = async (ctx) => {
    const {token} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    const publish = await redis.hget(keys.redis.SRS_AUTH_SECRET, 'pubSecret');
    if (!publish) throw utils.asError(errs.sys.boot, errs.status.sys, `system not boot yet`);

    console.log(`srs secret ok, key=${keys.redis.SRS_AUTH_SECRET}, field=pubSecret, value=${'*'.repeat(publish.length)}, decoded=${JSON.stringify(decoded)}`);
    ctx.body = utils.asResponse(0, {publish});
  };
  // Compatible with previous mgmt.
  router.all('/terraform/v1/hooks/srs/secret', handleSecretQuery);
  router.all('/terraform/v1/hooks/srs/secret/query', handleSecretQuery);

  router.all('/terraform/v1/hooks/srs/secret/update', async (ctx) => {
    const { token, secret} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    if (!secret) throw utils.asError(errs.sys.empty, errs.status.args, 'no secret');

    const r0 = await redis.hset(keys.redis.SRS_AUTH_SECRET, 'pubSecret', secret);
    await redis.set(keys.redis.SRS_SECRET_PUBLISH, secret);

    console.log(`hooks update secret, key=${keys.redis.SRS_AUTH_SECRET}, field=pubSecret, value=${'*'.repeat(secret.length)}, r0=${r0}, decoded=${JSON.stringify(decoded)}`);
    ctx.body = utils.asResponse(0, {});
  });

  router.all('/terraform/v1/hooks/srs/secret/disable', async (ctx) => {
    const { token, pubNoAuth} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    const r0 = await redis.hset(keys.redis.SRS_AUTH_SECRET, 'pubNoAuth', pubNoAuth);

    console.log(`hooks disable secret, pubNoAuth=${pubNoAuth}, r0=${r0}, decoded=${JSON.stringify(decoded)}`);
    ctx.body = utils.asResponse(0);
  });

  return router;
};

async function init() {
  let publish = await redis.hget(keys.redis.SRS_AUTH_SECRET, 'pubSecret');

  // Migrate from previous key.
  // TODO: FIXME: Remove SRS_SECRET_PUBLISH in the next major release.
  if (!publish) {
    let previous = await redis.get(keys.redis.SRS_SECRET_PUBLISH);
    if (previous) {
      publish = previous;
      await redis.hset(keys.redis.SRS_AUTH_SECRET, 'pubSecret', publish);
    }
  }

  // Setup the publish secret for first run.
  if (!publish) {
    publish = uuidv4().replace(/-/g, '');
    const r0 = await redis.hset(keys.redis.SRS_AUTH_SECRET, 'pubSecret', publish);
    await redis.set(keys.redis.SRS_SECRET_PUBLISH, publish);
    console.log(`hooks create secret, key=${keys.redis.SRS_AUTH_SECRET}, field=pubSecret, value=${'*'.repeat(publish.length)}, r0=${r0}`);
  }
}

