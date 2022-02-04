'use strict';

const { isMainThread, parentPort } = require("worker_threads");
const util = require('util');
const exec = util.promisify(require('child_process').exec);
const metadata = require('./metadata');
const os = require('os');

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  while (true) {
    try {
      await doThreadMain();
    } catch (e) {
      console.error(`Thread #market: err`, e);
      await new Promise(resolve => setTimeout(resolve, 300 * 1000));
    } finally {
      await new Promise(resolve => setTimeout(resolve, 3 * 1000));
    }
  }
}

async function doThreadMain() {
  for (const e in metadata.market) {
    const conf = metadata.market[e];
    const container = await doContainerMain(conf);

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

  if (true) {
    const {stdout} = await exec(`docker ps -a -f name=${name} --format '{{json .}}'`);
    all = stdout ? JSON.parse(stdout) : {};
  }

  if (true) {
    const {stdout} = await exec(`docker ps -f name=${name} --format '{{json .}}'`);
    running = stdout ? JSON.parse(stdout) : {};
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
  const dockerArgs = `-d --restart always --privileged -it --name ${evalValue(conf.name)} \\
    --add-host=mgmt.srs.local:${privateIPv4.address} \\
    ${tcpPorts} ${udpPorts} \\
    ${evalValue(conf.logConfig)} \\
    ${volumes} ${evalValue(conf.extras)} \\
    ${evalValue(conf.image)} \\
    ${evalValue(conf.command)}`;
  console.log(`Thread #market: docker run args ip=${privateIPv4.name}/${privateIPv4.address}, docker run ${dockerArgs}`);

  // Only remove the container when got ID, to avoid fail for CentOS.
  const all = (await queryContainer(conf.name))[0];
  if (all && all.ID) {
    await exec(`docker rm -f ${conf.name}`);
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

