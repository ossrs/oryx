'use strict';

// For mgmt, it's ok to connect to localhost.
const config = {
  redis:{
    host: 'localhost',
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const keys = require('js-core/keys');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const moment = require('moment');

async function inUpgradeWindow() {
  const uwStart = await redis.hget(keys.redis.SRS_UPGRADE_WINDOW, 'start');
  const uwDuration = await redis.hget(keys.redis.SRS_UPGRADE_WINDOW, 'duration');
  if (uwStart === undefined || uwStart === null || !uwDuration) return true;

  const [start, duration] = [parseInt(uwStart), parseInt(uwDuration)];
  const end = (start + duration) % 24;

  let inWindow;
  if (start < end) {
    inWindow = (start <= moment().hours() && moment().hours() <= end);
  } else {
    inWindow = !(end <= moment().hours() && moment().hours() <= start);
  }
  console.log(`Upgrade window=${inWindow}, start=${start}, duration=${duration}, end=${end}, now=${moment().format()}, hours=${moment().hours()}`);

  return inWindow;
}
exports.inUpgradeWindow = inUpgradeWindow;

