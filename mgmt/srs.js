'use strict';

const os = require('os');
const { isMainThread, parentPort } = require("worker_threads");
const util = require('util');
const exec = util.promisify(require('child_process').exec);
const fs = require('fs');
const metadata = require('./metadata');

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  while (true) {
    try {
      await doThreadMain();
    } catch (e) {
      console.error(`Thread #${metadata.srs.name}: err`, e);
      await new Promise(resolve => setTimeout(resolve, 300 * 1000));
    } finally {
      await new Promise(resolve => setTimeout(resolve, 3 * 1000));
    }
  }
}

async function doThreadMain() {
  // Query the container from docker.
  let [all, running] = await queryContainer();
  if (all && all.ID) {
    metadata.srs.container = all;
    console.log(`Thread #${metadata.srs.name}: query ID=${all.ID}, State=${all.State}, Status=${all.Status}, running=${running?.ID}`);
  }

  // Restart the SRS container.
  if (!all || !all.ID || !running || !running.ID) {
    await startContainer();

    all = (await queryContainer())[0];
    if (all && all.ID) metadata.srs.container = all;
    console.log(`Thread #${metadata.srs.name}: create ID=${all.ID}, State=${all.State}, Status=${all.Status}`);
  }

  // Update the metadata to main thread.
  parentPort.postMessage({
    metadata:{
      srs: metadata.srs,
    },
  });
}

// See https://docs.docker.com/config/formatting/
async function queryContainer() {
  let all, running;

  if (true) {
    const {stdout} = await exec(`docker ps -a -f name=${metadata.srs.name} --format '{{json .}}'`);
    all = stdout ? JSON.parse(stdout) : {};
  }

  if (true) {
    const {stdout} = await exec(`docker ps -f name=${metadata.srs.name} --format '{{json .}}'`);
    running = stdout ? JSON.parse(stdout) : {};
  }

  return [all, running];
}
exports.queryContainer = queryContainer;

async function startContainer() {
  console.log(`Thread #${metadata.srs.name}: start container`);

  // We must create the file, or docker will mount it as directory, see https://stackoverflow.com/a/44950494/17679565
  const logFile = `${process.cwd()}/containers/objs/srs.log`;
  if (!fs.existsSync(logFile)) fs.createWriteStream(logFile, {overwrite: false});

  const privateIPv4 = await discoverPrivateIPv4();
  const confFile = `${process.cwd()}/containers/conf/srs.conf`;
  const dockerArgs = `-d -it --restart always --privileged --name ${metadata.srs.name} \\
    --add-host=mgmt.srs.local:${privateIPv4.address} \\
    -v ${confFile}:/usr/local/srs/conf/lighthouse.conf \\
    -v ${logFile}:/usr/local/srs/objs/srs.log \\
    -p 1935:1935 -p 1985:1985 -p 8080:8080 -p 8000:8000/udp -p 10080:10080/udp \\
    registry.cn-hangzhou.aliyuncs.com/ossrs/lighthouse:${metadata.srs.major} \\
    ./objs/srs -c conf/lighthouse.conf`;
  console.log(`Thread #${metadata.srs.name}: docker run args ip=${privateIPv4.name}/${privateIPv4.address}, conf=${confFile}, docker run ${dockerArgs}`);

  // Only remove the container when got ID, to avoid fail for CentOS.
  const all = (await queryContainer())[0];
  if (all && all.ID) {
    await exec(`docker rm -f ${metadata.srs.name}`);
    console.log(`Thread #${metadata.srs.name}: docker run remove ID=${all.ID}`);
  }

  await exec(`docker run ${dockerArgs}`);
  console.log(`Thread #${metadata.srs.name}: docker run ok`);
}

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
  console.log(`Thread #${metadata.srs.name}: discover ip networks=${JSON.stringify(networks)}`);

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
  console.log(`Thread #${metadata.srs.name}: discover ip privateIPv4=${JSON.stringify(privateIPv4)}`);

  return privateIPv4;
}
