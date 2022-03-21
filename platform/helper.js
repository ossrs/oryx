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
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const jwt = require('jsonwebtoken');
const axios = require('axios');
const moment = require('moment');

exports.inUpgradeWindow = (uwStart, uwDuration, now) => {
  if (uwStart === undefined || uwStart === null || !uwDuration || !now) {
    return true;
  }

  const [start, duration] = [parseInt(uwStart), parseInt(uwDuration)];
  const end = duration >= 24 ? start + 24 : (start + duration) % 24;

  let inWindow;
  if (start < end) {
    inWindow = (start <= now.hours() && now.hours() <= end);
  } else {
    inWindow = !(end < now.hours() && now.hours() < start);
  }
  console.log(`Upgrade window=${inWindow}, start=${start}, duration=${duration}, end=${end}, now=${now.format()}, hours=${now.hours()}`);

  return inWindow;
}

async function execApi(action, args) {
  if (args !== null && args !== undefined && !Array.isArray(args)) {
    throw new Error(`args is not array, ${args}`);
  }

  const apiSecret = await utils.apiSecret(redis);
  const token = utils.createToken(moment, jwt, apiSecret);
  const server = process.env.NODE_ENV === 'development' ? 'localhost' : 'mgmt.srs.local';

  try {
    const data = await axios.post(
      `http://${server}:2022/terraform/v1/host/exec`,
      {
        ...token, action, args: args || [],
      },
    );
    return data?.data?.data;
  } catch (e) {
    console.error(`exec server=${server}, action=${action}, args=${JSON.stringify(args)}, err`, e);
  }
}
exports.execApi = execApi;

