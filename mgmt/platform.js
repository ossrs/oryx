'use strict';

const axios = require('axios');

exports.isDarwin = process.platform === 'darwin';

// We must mark these fields as async, to notice user not to use it before it initialized.
const conf = {region: null, registry: null};
exports.region = async () => {
  return conf.region;
};
exports.registry = async () => {
  return conf.registry;
};

// Initialize the platform before thread run.
exports.init = async () => {
  const isDarwin = exports.isDarwin;

  const region = await discoverRegion();
  conf.region = region;

  const registry = await discoverRegistry(region);
  conf.registry = registry;

  console.log(`Initialize region=${region}, registry=${registry}, isDarwin=${isDarwin}`);
  return {region, registry, isDarwin};
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

async function discoverRegistry(region) {
  if (exports.isDarwin) {
    return 'ccr.ccs.tencentyun.com';
  }

  let registry = 'sgccr.ccs.tencentyun.com';
  ['ap-guangzhou', 'ap-shanghai', 'ap-nanjing', 'ap-beijing', 'ap-chengdu', 'ap-chongqing'].filter(v => {
    if (region.startsWith(v)) registry = 'ccr.ccs.tencentyun.com';
    return null;
  });

  return registry;
}

