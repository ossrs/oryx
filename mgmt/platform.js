'use strict';

// For mgmt, it's ok to connect to localhost.
const config = {
  redis:{
    host: 'localhost', // For mgmt, we always use localhost, rather than docker container name.
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const os = require('os');
const fs = require('fs');
const dotenv = require('dotenv');
const axios = require('axios');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const keys = require('js-core/keys');
const utils = require('js-core/utils');
const { v4: uuidv4 } = require('uuid');
const moment = require('moment');

exports.isDarwin = process.platform === 'darwin';

// We must mark these fields as async, to notice user not to use it before it initialized.
const conf = {cloud: null, region: null, source: null, registry: null};
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

// Initialize the source for redis, note that we don't change the env.
exports.initOs = async() => {
  // The redis is not available when os startup, so we must directly discover from env or network.
  const {cloud, region} = await discoverRegion();
  conf.cloud = cloud; conf.region = region;

  // Always update the source, because it might change.
  const source = await discoverSource(conf.cloud, conf.region);
  conf.source = source;

  // Always update the registry, because it might change.
  discoverRegistry(source);

  // Create directories for data, allow user to link it.
  [
    "containers/data/dvr", "containers/data/prometheus", "containers/data/record", "containers/data/vod",
    "containers/data/upload", "containers/data/vlive",
  ].forEach((f) => {
    if (!fs.existsSync(f)) {
      fs.mkdirSync(f, {recursive: true})
    }
  })
};

function discoverRegistry(source) {
  const registry = (source === 'github') ? 'docker.io' : 'registry.cn-hangzhou.aliyuncs.com';
  conf.registry = registry;
  return registry;
}

// Initialize the platform before thread run.
exports.init = async () => {
  const isDarwin = exports.isDarwin;

  // Initialize the node id.
  let nid = await redis.hget(keys.redis.SRS_TENCENT_LH, 'node');
  if (!nid) {
    nid = uuidv4();
    await redis.hset(keys.redis.SRS_TENCENT_LH, 'node', nid);
  }

  // Setup the api secret.
  // Remark: Should do it before any helper.execApi, which depends on it.
  await utils.setupApiSecret(redis, uuidv4, moment);

  // Load the region first, because it never changed.
  conf.cloud = await redis.hget(keys.redis.SRS_TENCENT_LH, 'cloud');
  conf.region = await redis.hget(keys.redis.SRS_TENCENT_LH, 'region');
  if (!conf.cloud || !conf.region) {
    const {cloud, region} = await discoverRegion();
    conf.cloud = cloud; conf.region = region;
    await redis.hset(keys.redis.SRS_TENCENT_LH, 'cloud', cloud);
    await redis.hset(keys.redis.SRS_TENCENT_LH, 'region', region);
  }

  // Always update the source, because it might change.
  const source = await discoverSource(conf.cloud, conf.region);
  conf.source = source;
  await redis.hset(keys.redis.SRS_TENCENT_LH, 'source', source);

  // Refresh the env file.
  utils.saveEnvs(fs, os, dotenv, '.env', {
    CLOUD: conf.cloud,
    REGION: conf.region,
    SOURCE: conf.source,
  });

  // Always update the registry, because it might change.
  const registry = discoverRegistry(source);
  conf.registry = registry;
  await redis.hset(keys.redis.SRS_TENCENT_LH, 'registry', registry);

  // Load the cvm first, because it never changed.
  let platform = await redis.hget(keys.redis.SRS_TENCENT_LH, 'platform');
  if (!platform) {
    platform = await discoverPlatform(conf.cloud);
    await redis.hset(keys.redis.SRS_TENCENT_LH, 'platform', platform);
  }

  // Request and cache the apiSecret.
  const apiSecret = await utils.apiSecret(redis);

  console.log(`Initialize node=${nid}, cloud=${conf.cloud}, region=${conf.region}, source=${source}, registry=${registry}, platform=${platform}, isDarwin=${isDarwin}, apiSecret=${apiSecret.length}B`);
  return {region: conf.region, registry, isDarwin};
};

async function discoverRegion() {
  if (exports.isDarwin) {
    return {cloud: 'DEV', region: 'ap-beijing'};
  }

  if (process.env.CLOUD === 'BT') {
    return {cloud: 'BT', region: 'ap-beijing'};
  }

  if (process.env.CLOUD === 'AAPANEL') {
    return {cloud: 'AAPANEL', region: 'ap-singapore'};
  }

  if (process.env.CLOUD && process.env.REGION) {
    return {cloud: process.env.CLOUD, region: process.env.REGION};
  }

  console.log(`Initialize start to discover region`);
  return await new Promise((resolve, reject) => {
    axios.get(`http://metadata.tencentyun.com/latest/meta-data/placement/region`).then((data) => {
      if (data?.data) resolve({cloud: 'TENCENT', region: data.data});
    }).catch(e => {
      console.warn(`Ignore tencent region error, ${e.message}`);
    });

    // See https://docs.digitalocean.com/reference/api/metadata-api/#operation/getRegion
    axios.get(`http://169.254.169.254/metadata/v1/region`).then((data) => {
      if (data?.data) resolve({cloud: 'DO', region: data.data});
    }).catch(e => {
      console.warn(`Ignore do region error, ${e.message}`);
    });
  });
}

async function discoverSource(cloud, region) {
  if (cloud === 'DEV') return 'gitee';
  if (cloud === 'DO') return 'github';
  if (cloud === 'BT') return 'gitee';
  if (cloud === 'AAPANEL') return 'github';

  let source = 'github';
  if (cloud === 'TENCENT') {
    ['ap-guangzhou', 'ap-shanghai', 'ap-nanjing', 'ap-beijing', 'ap-chengdu', 'ap-chongqing'].filter(v => {
      if (region.startsWith(v)) source = 'gitee';
      return null;
    });
  }

  return source;
}

async function discoverPlatform(cloud) {
  if (cloud === 'DEV') return 'dev';
  if (cloud === 'DO') return 'droplet';
  if (cloud === 'BT') return 'bt';
  if (cloud === 'AAPANEL') return 'aapanel';

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
  if (privateIPv4?.name?.indexOf('en') !== 0 && privateIPv4?.name?.indexOf('eth') !== 0) {
    Object.keys(networks).map(e => {
      if (e.indexOf('en') === 0 || e.indexOf('eth') === 0) {
        privateIPv4 = networks[e];
      }
    });
  }
  console.log(`discover ip privateIPv4=${JSON.stringify(privateIPv4)}, update=${privateIPv4Update.format()}`);

  return privateIPv4;
}

