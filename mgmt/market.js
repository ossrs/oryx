'use strict';

// For mgmt, it's ok to connect to localhost.
const config = {
  redis:{
    host: 'localhost',
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const { isMainThread, parentPort } = require("worker_threads");
const util = require('util');
const exec = util.promisify(require('child_process').exec);
const execFile = util.promisify(require('child_process').execFile);
const metadata = require('./metadata');
const os = require('os');
const platform = require('./platform');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const utils = require('js-core/utils');
const keys = require('js-core/keys');

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  // We must initialize the thread first.
  const {region, registry} = await platform.init();
  console.log(`Thread #market: initialize region=${region}, registry=${registry}`);

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
    const container = await doContainerMain(conf);

    // Ignore if not running.
    if (!container) continue;

    // Update the metadata to main thread.
    const msg = {metadata: {}};
    msg.metadata[e] = container;
    parentPort.postMessage(msg);
  }
}

async function doContainerMain(conf) {
  let container = null;

  // Query the container from docker.
  let [all, running] = await queryContainer(conf.name);
  if (all && all.ID) {
    container = all;
    console.log(`Thread #market: query ID=${all.ID}, State=${all.State}, Status=${all.Status}, running=${running?.ID}`);
  }

  // Restart the SRS container.
  if (!all || !all.ID || !running || !running.ID) {
    // Query container enabled status from redis.
    const disabled = await redis.hget(keys.redis.SRS_CONTAINER_DISABLED, conf.name);
    if (disabled === 'true') {
      console.log(`Thread #market: container ${conf.name} disable`);
      return container;
    }

    await startContainer(conf);

    all = (await queryContainer(conf.name))[0];
    if (all && all.ID) container = all;
    console.log(`Thread #market: create ID=${all.ID}, State=${all.State}, Status=${all.Status}`);
  }

  return container;
}

// See https://docs.docker.com/config/formatting/
async function queryContainer(name) {
  let all, running;

  try {
    const {stdout} = await exec(`docker ps -a -f name=${name} --format '{{json .}}'`);
    all = stdout ? JSON.parse(stdout) : {};
  } catch (e) {
    console.log(`Thread #market: Ignore query container ${name} err`, e);
  }

  try {
    const {stdout} = await exec(`docker ps -f name=${name} --format '{{json .}}'`);
    running = stdout ? JSON.parse(stdout) : {};
  } catch (e) {
    console.log(`Thread #market: Ignore query container ${name} err`, e);
  }

  return [all, running];
}
exports.queryContainer = queryContainer;

async function startContainer(conf) {
  console.log(`Thread #market: start container`);

  const evalValue = (e, defaults) => {
    if (!e) return defaults || '';
    if (typeof(e) === 'function') return e();
    return e;
  };

  const privateIPv4 = await discoverPrivateIPv4();
  const tcpPorts = evalValue(conf.tcpPorts, []).map(e => `-p ${e}:${e}/tcp`).join(' ');
  const udpPorts = evalValue(conf.udpPorts, []).map(e => `-p ${e}:${e}/udp`).join(' ');
  const volumes = evalValue(conf.volumes, []).map(e => `-v "${e}"`).join(' ');
  const command = evalValue(conf.command, []).join(' ');
  const extras = evalValue(conf.extras, []).join(' ');
  // The image depends on the registry, which is discovered by platform.
  const image = await conf.image();
  const region = await platform.region();
  const source = await platform.source();
  // Note that it's started by nodejs, so never use '-it'.
  const dockerArgs = `-d --restart always --privileged --name ${evalValue(conf.name)} \\
    --add-host=mgmt.srs.local:${privateIPv4.address} \\
    ${tcpPorts} ${udpPorts} \\
    ${evalValue(conf.logConfig)} \\
    ${volumes} ${extras} \\
    --env SRS_REGION=${region} --env SRS_SOURCE=${source} \\
    ${image} \\
    ${command}`;
  console.log(`Thread #market: docker run args ip=${privateIPv4.name}/${privateIPv4.address}, docker run ${dockerArgs}`);

  // Only remove the container when got ID, to avoid fail for CentOS.
  const all = (await queryContainer(conf.name))[0];
  if (all && all.ID) {
    await utils.removeContainerQuiet(execFile, conf.name);
    console.log(`Thread #market: docker run remove ID=${all.ID}`);
  }

  await exec(`docker run ${dockerArgs}`);
  console.log(`Thread #market: docker run ok`);
}

// Discover the private ip of machine.
let privateIPv4 = null;
async function discoverPrivateIPv4() {
  if (privateIPv4) return privateIPv4;

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

  // Best match the en or eth network, for example, eth0 or en0.
  Object.keys(networks).map(e => {
    if (e.indexOf('en') === 0 || e.indexOf('eth') === 0) {
      privateIPv4 = networks[e];
    }
  });
  console.log(`discover ip privateIPv4=${JSON.stringify(privateIPv4)}`);

  return privateIPv4;
}

