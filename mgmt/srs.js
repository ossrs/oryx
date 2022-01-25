'use strict';

const util = require('util');
const os = require('os');
const { isMainThread, parentPort } = require("worker_threads");
const exec = util.promisify(require('child_process').exec);

const metadata = {
  name: 'srs-server',
  major: '4',
  container: {
    id: null,
  },
};
exports.metadata = metadata;

if (!isMainThread) {
  _thread_main();
}

async function _thread_main() {
  while (true) {
    try {
      await thread_main();
    } catch (e) {
      console.error(`thread ${metadata.name} err`, e);
    }
    await new Promise(resolve => setTimeout(resolve, 3000));
  }
}

async function thread_main() {
  // Query the container from docker.
  let [all, running] = await queryContainer();
  if (all && all.ID) {
    metadata.container = all;
    console.log(`Thread #${metadata.name}: query ID=${all.ID}, State=${all.State}, Status=${all.Status}, running=${running?.ID}`);
  }

  // Restart the SRS container.
  if (!all || !all.ID || !running || !running.ID) {
    await startContainer();

    all = (await queryContainer())[0];
    if (all && all.ID) metadata.container = all;
    console.log(`Thread #${metadata.name}: create ID=${all.ID}, State=${all.State}, Status=${all.Status}`);
  }

  // Update the metadata to main thread.
  parentPort.postMessage({metadata});
}

// See https://docs.docker.com/config/formatting/
async function queryContainer() {
  let all, running;

  if (true) {
    const {stdout} = await exec(`docker ps -a -f name=${metadata.name} --format '{{json .}}'`);
    all = stdout ? JSON.parse(stdout) : {};
  }

  if (true) {
    const {stdout} = await exec(`docker ps -f name=${metadata.name} --format '{{json .}}'`);
    running = stdout ? JSON.parse(stdout) : {};
  }

  return [all, running];
}

async function startContainer() {
  console.log(`Thread #${metadata.name}: start container`);

  const privateIPv4 = await discoverPrivateIPv4();
  const confFile = `${process.cwd()}/containers/conf/srs.conf`;
  const dockerArgs = `-d -it --restart always --privileged --name ${metadata.name} \\
    --add-host=mgmt.srs.local:${privateIPv4.address} \\
    -v ${confFile}:/usr/local/srs/conf/lighthouse.conf \\
    -p 1935:1935 -p 1985:1985 -p 8080:8080 -p 8000:8000/udp \\
    registry.cn-hangzhou.aliyuncs.com/ossrs/lighthouse:${metadata.major} \\
    ./objs/srs -c conf/lighthouse.conf`;
  console.log(`Thread #${metadata.name}: docker run args ip=${privateIPv4.name}/${privateIPv4.address}, conf=${confFile}, docker run ${dockerArgs}`);

  // Only remove the container when got ID, to avoid fail for CentOS.
  const all = (await queryContainer())[0];
  if (all && all.ID) {
    await exec(`docker rm -f ${metadata.name}`);
    console.log(`Thread #${metadata.name}: docker run remove ID=${all.ID}`);
  }

  await exec(`docker run ${dockerArgs}`);
  console.log(`Thread #${metadata.name}: docker run ok`);
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
  console.log(`Thread #${metadata.name}: discover ip networks=${JSON.stringify(networks)}`);

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
  console.log(`Thread #${metadata.name}: discover ip privateIPv4=${JSON.stringify(privateIPv4)}`);

  return privateIPv4;
}
