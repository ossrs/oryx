'use strict';

// For mgmt, it's ok to connect to localhost.
const config = {
  redis:{
    host: 'localhost',
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const axios = require('axios');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const keys = require('js-core/keys');
const utils = require('js-core/utils');
const os = require('os');
const { v4: uuidv4 } = require('uuid');
const moment = require('moment');

exports.isDarwin = process.platform === 'darwin';

// We must mark these fields as async, to notice user not to use it before it initialized.
const conf = {region: null, source: null, registry: null};
exports.region = async () => {
  return conf.region;
};
exports.source = async () => {
  return conf.source;
};
exports.registry = async () => {
  return conf.registry;
};
exports.ipv4 = async () => {
  return await discoverPrivateIPv4();
};

// Initialize the platform before thread run.
exports.init = async () => {
  const isDarwin = exports.isDarwin;

  // Setup the api secret.
  // Remark: Should do it before any helper.execApi, which depends on it.
  await utils.setupApiSecret(redis, uuidv4, moment);

  // Load the region first, because it never changed.
  conf.region = await redis.hget(keys.redis.SRS_TENCENT_LH, 'region');
  if (!conf.region) {
    const region = await discoverRegion();
    conf.region = region;
    await redis.hset(keys.redis.SRS_TENCENT_LH, 'region', region);
  }

  // Always update the source, because it might change.
  const source = await discoverSource(conf.region);
  conf.source = source;
  await redis.hset(keys.redis.SRS_TENCENT_LH, 'source', source);

  // Always update the registry, because it might change.
  const registry = (source === 'github') ? 'sgccr.ccs.tencentyun.com' : 'registry.cn-hangzhou.aliyuncs.com';
  conf.registry = registry;
  await redis.hset(keys.redis.SRS_TENCENT_LH, 'registry', registry);

  // Load the cvm first, because it never changed.
  let platform = await redis.hget(keys.redis.SRS_TENCENT_LH, 'platform');
  if (!platform) {
    platform = await discoverPlatform();
    await redis.hset(keys.redis.SRS_TENCENT_LH, 'platform', platform);
  }

  // Request and cache the apiSecret.
  const apiSecret = await utils.apiSecret(redis);

  console.log(`Initialize region=${conf.region}, source=${source}, registry=${registry}, platform=${platform}, isDarwin=${isDarwin}, apiSecret=${apiSecret.length}B`);
  return {region: conf.region, registry, isDarwin};
};

async function discoverRegion() {
  if (exports.isDarwin) {
    return 'ap-beijing';
  }

  if (process.env.REGION) {
    return process.env.REGION;
  }

  const {data} = await axios.get(`http://metadata.tencentyun.com/latest/meta-data/placement/region`);
  return data;
}

async function discoverSource(region) {
  if (exports.isDarwin) {
    return 'gitee';
  }

  let source = 'github';
  ['ap-guangzhou', 'ap-shanghai', 'ap-nanjing', 'ap-beijing', 'ap-chengdu', 'ap-chongqing'].filter(v => {
    if (region.startsWith(v)) source = 'gitee';
    return null;
  });

  return source;
}

async function discoverPlatform() {
  if (exports.isDarwin) {
    return 'dev';
  }

  if (process.env.PLATFORM) {
    return process.env.PLATFORM;
  }

  const {data} = await axios.get(`http://metadata.tencentyun.com/latest/meta-data/instance-name`);
  return data.indexOf('-lhins-') > 0 ? 'lighthouse' : 'cvm';
}

// Discover the private ip of machine.
let privateIPv4 = null;
let privateIPv4Update = null;
async function discoverPrivateIPv4() {
  if (privateIPv4 && privateIPv4Update) {
    // If not expired, return the cache.
    const expired = moment(privateIPv4Update).add(process.env.NODE_ENV === 'development' ? 10 : 24 * 3600, 's');
    if (expired.isAfter(moment())) return privateIPv4;
  }

  const networks = {};

  const networkInterfaces = os.networkInterfaces();
  Object.keys(networkInterfaces).map(name => {
    for (const network of networkInterfaces[name]) {
      if (network.family === 'IPv4' && !network.internal) {
        networks[name] = {...network, name};
      }
    }
  });
  console.log(`discover ip networks=${JSON.stringify(networks)}`);

  if (!Object.keys(networks).length) {
    throw new Error(`no private address from ${JSON.stringify(networkInterfaces)}`);
  }

  // Default to the first one.
  privateIPv4 = networks[Object.keys(networks)[0]];
  privateIPv4Update = moment();

  // Best match the en or eth network, for example, eth0 or en0.
  Object.keys(networks).map(e => {
    if (e.indexOf('en') === 0 || e.indexOf('eth') === 0) {
      privateIPv4 = networks[e];
    }
  });
  console.log(`discover ip privateIPv4=${JSON.stringify(privateIPv4)}, update=${privateIPv4Update.format()}`);

  return privateIPv4;
}

