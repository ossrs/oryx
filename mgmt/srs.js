'use strict';

const { isMainThread, parentPort } = require("worker_threads");
const util = require('util');
const exec = util.promisify(require('child_process').exec);
const fs = require('fs');
const metadata = require('./metadata');
const utils = require('./utils');

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

  const privateIPv4 = await utils.discoverPrivateIPv4();
  const confFile = `${process.cwd()}/containers/conf/srs.conf`;
  const image = (process.env.NODE_ENV === 'development' || process.env.SRS_DOCKER === 'srs') ? 'ossrs/srs' : 'ossrs/lighthouse';
  const dockerArgs = `-d -it --restart always --privileged --name ${metadata.srs.name} \\
    --add-host=mgmt.srs.local:${privateIPv4.address} \\
    -v ${confFile}:/usr/local/srs/conf/lighthouse.conf \\
    -p 1935:1935 -p 1985:1985 -p 8080:8080 -p 8000:8000/udp -p 10080:10080/udp \\
    --log-driver json-file --log-opt max-size=3g --log-opt max-file=3 \\
    registry.cn-hangzhou.aliyuncs.com/${image}:${metadata.srs.major} \\
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
