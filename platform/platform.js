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
const utils = require('js-core/utils');

exports.init = async () => {
  // Request and cache the apiSecret.
  const apiSecret = await utils.apiSecret(redis);
  console.log(`Initialize apiSecret=${apiSecret.length}B`);
};

