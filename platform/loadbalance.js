'use strict';

// For components in docker, connect by host.
const config = {
  redis:{
    host: process.env.NODE_ENV === 'development' ? 'localhost' : 'mgmt.srs.local',
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const keys = require('js-core/keys');
const errs = require('js-core/errs');
const utils = require('js-core/utils');

exports.handle = (router) => {
  router.all('/terraform/v1/mgmt/dns/lb/:app/:stream', async (ctx) => {
    const {app, stream} = ctx.params;

    if (!app) throw utils.asError(errs.sys.empty, errs.status.args, 'no app');
    if (!stream) throw utils.asError(errs.sys.empty, errs.status.args, 'no stream');

    const url = `/${app}/${stream}`;
    const r0 = await redis.hget(keys.redis.SRS_DNS_LB_BACKENDS, url);
    if (!r0) throw utils.asError(errs.sys.invalid, errs.status.args, `no lb for ${url}`);

    const backends = JSON.parse(r0);
    const seed = parseInt(Math.random() * 1000);
    const selected = backends[seed % backends.length];
    const protocol = selected.indexOf('://') === -1 ? `${ctx.request.protocol}://` : '';

    const redirectPath = `${protocol}${selected}${url}`;
    const redirectUrl = ctx.request.querystring ? `${redirectPath}?${ctx.request.querystring}` : redirectPath;
    console.log(`DNS LB url=${url}, seed=${seed}, redirectUrl=${redirectUrl}`);

    ctx.response.redirect(redirectUrl);
  });
};

