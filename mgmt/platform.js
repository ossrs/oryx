'use strict';

const axios = require('axios');

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

// Initialize the platform before thread run.
exports.init = async () => {
  const isDarwin = exports.isDarwin;

  const region = await discoverRegion();
  conf.region = region;

  const source = await discoverSource(region);
  conf.source = source;

  const registry = (source === 'github') ? 'sgccr.ccs.tencentyun.com' : 'ccr.ccs.tencentyun.com';
  conf.registry = registry;

  console.log(`Initialize region=${region}, source=${source}, registry=${registry}, isDarwin=${isDarwin}`);
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

