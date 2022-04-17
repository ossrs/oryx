'use strict';

// For components in docker, connect by host.
const config = {
  redis:{
    host: process.env.NODE_ENV === 'development' ? 'localhost' : (process.env.REDIS_HOST || 'mgmt.srs.local'),
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const moment = require('moment');
const { v4: uuidv4 } = require('uuid');
const utils = require('js-core/utils');
const platform = require('./platform');

const { isMainThread } = require("worker_threads");

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  // We must initialize the thread first.
  await platform.init();
  console.log(`Thread #apiSecret: initialize`);

  while (true) {
    try {
      await doThreadMain();
    } catch (e) {
      console.error(`Thread #apiSecret: err`, e);
    } finally {
      await new Promise(resolve => setTimeout(resolve, 3600 * 1000));
    }
  }
}

async function doThreadMain() {
  const [token, created] = await utils.setupApiSecret(redis, uuidv4, moment);
  console.log(`Thread #apiSecret: Platform api secret, token=${token.length}B, created=${created}`);
}

