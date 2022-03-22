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
const { isMainThread } = require("worker_threads");
const platform = require('./platform');
const metadata = require('./metadata');
const helper = require('./helper');
const keys = require('js-core/keys');

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  // We must initialize the thread first.
  await platform.init();
  console.log(`Thread #market: initialize`);

  while (true) {
    try {
      await doThreadMain();
    } catch (e) {
      console.error(`Thread #market: err`, e);
      await new Promise(resolve => setTimeout(resolve, 30 * 1000));
    } finally {
      await new Promise(resolve => setTimeout(resolve, 10 * 1000));
    }
  }
}

async function doThreadMain() {
  // For SRS, if release enabled, disable dev automatically.
  const srsReleaseDisabled = await redis.hget(keys.redis.SRS_CONTAINER_DISABLED, metadata.market.srs.name);
  const srsDevDisabled = await redis.hget(keys.redis.SRS_CONTAINER_DISABLED, metadata.market.srsDev.name);
  if (srsReleaseDisabled !== 'true' && srsDevDisabled !== 'true') {
    const r0 = await redis.hset(keys.redis.SRS_CONTAINER_DISABLED, metadata.market.srsDev.name, true);
    await utils.removeContainerQuiet(execFile, metadata.market.srsDev.name);
    console.log(`Thread #market: Disable srs dev for release enabled, r0=${r0}`);
  }

  for (const e in metadata.market) {
    const conf = metadata.market[e];

    // If no image, it's only a hint, for example, the container is managed by others.
    if (!conf.image) continue;

    // Try to restart the container.
    await doContainerMain(e, conf);
  }
}

async function doContainerMain(marketName, conf) {
  let container = null;

  // Query the container from docker.
  let {all, running} = await helper.execApi('fetchContainer', [conf.name, marketName]);
  if (all && all.ID) {
    container = all;
    console.log(`Thread #market: query ID=${all.ID}, State=${all.State}, Status=${all.Status}, running=${running?.ID}`);
  }

  // Restart the SRS container.
  if (!all?.ID || !running?.ID) {
    // Query container enabled status from redis.
    const disabled = await redis.hget(keys.redis.SRS_CONTAINER_DISABLED, conf.name);
    if (disabled === 'true') {
      console.log(`Thread #market: container ${conf.name} disable`);
      return container;
    }

    console.log(`Thread #market: start container`);
    const privateIPv4 = await platform.ipv4();
    const dockerArgs = await utils.generateDockerArgs(platform, privateIPv4, conf);
    await helper.execApi('startContainer', [conf.name, dockerArgs]);

    let {all} = await helper.execApi('fetchContainer', [conf.name, marketName]);
    if (all && all.ID) container = all;
    console.log(`Thread #market: create ID=${all.ID}, State=${all.State}, Status=${all.Status}`);
  }

  return container;
}

async function startContainer(conf) {
  console.log(`Thread #market: start container`);

  const evalValue = (e, defaults) => {
    if (!e) return defaults || '';
    if (typeof(e) === 'function') return e();
    return e;
  };

  const tcpPorts = evalValue(conf.tcpPorts, []).map(e => `-p ${e}:${e}/tcp`);
  const udpPorts = evalValue(conf.udpPorts, []).map(e => `-p ${e}:${e}/udp`);
  const volumes = evalValue(conf.volumes, []).map(e => `-v "${e}"`);
  const command = evalValue(conf.command, []);
  const extras = evalValue(conf.extras, []);
  const logConfig = evalValue(conf.logConfig, []);
  // The image depends on the registry, which is discovered by platform.
  const image = await conf.image();
  const region = await platform.region();
  const source = await platform.source();
  // Note that it's started by nodejs, so never use '-it'.
  const dockerArgs = `-d --restart=always --privileged --name ${evalValue(conf.name)} \\
    ${tcpPorts} ${udpPorts} \\
    ${logConfig} \\
    ${volumes} ${extras} \\
    --env SRS_REGION=${region} --env SRS_SOURCE=${source} \\
    ${image} \\
    ${command}`;
  console.log(`Thread #market: docker run args ${dockerArgs}`);

  // Only remove the container when got ID, to avoid fail for CentOS.
  const all = (await queryContainer(conf.name))[0];
  if (all && all.ID) {
    await utils.removeContainerQuiet(execFile, conf.name);
    console.log(`Thread #market: docker run remove ID=${all.ID}`);
  }

  await exec(`docker run ${dockerArgs}`);
  console.log(`Thread #market: docker run ok`);
}

