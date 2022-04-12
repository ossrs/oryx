'use strict';

// For mgmt, it's ok to connect to localhost.
const config = {
  redis:{
    host: 'localhost',
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const fs = require('fs');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const keys = require('js-core/keys');

async function handleFavicon(ctx, next) {
  ctx.type = 'image/x-icon';
  ctx.set('Cache-Control', 'public, max-age=31536000');
  ctx.body = fs.readFileSync('./ui/build/favicon.ico');
  return;
}

async function handleRoot(ctx, next) {
  const r0 = await redis.hget(keys.redis.SRS_HTTP_REWRITE, '/');
  return ctx.response.redirect(r0 || '/mgmt/');
}

exports.handle = (app) => {
  app.use(async (ctx, next) => {
    if (ctx.request.path === '/favicon.ico') {
      return await handleFavicon(ctx, next);
    } else if (ctx.request.path === '/' || ctx.request.path === '/index.html') {
      return await handleRoot(ctx, next);
    } else {
      await next();
    }
  });
};

