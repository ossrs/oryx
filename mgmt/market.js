'use strict';

const os = require('os');
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
      console.error(`Thread #market: err`, e);
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
    metadata.market.hooks.container = all;
    console.log(`Thread #market: query ID=${all.ID}, State=${all.State}, Status=${all.Status}, running=${running?.ID}`);
  }

  // Restart the SRS container.
  if (!all || !all.ID || !running || !running.ID) {
    await startContainer();

    all = (await queryContainer())[0];
    if (all && all.ID) metadata.market.hooks.container = all;
    console.log(`Thread #market: create ID=${all.ID}, State=${all.State}, Status=${all.Status}`);
  }

  // Update the metadata to main thread.
  parentPort.postMessage({
    metadata:{
      hooks: metadata.market.hooks,
    },
  });
}

// See https://docs.docker.com/config/formatting/
async function queryContainer() {
  let all, running;

  if (true) {
    const {stdout} = await exec(`docker ps -a -f name=${metadata.market.hooks.name} --format '{{json .}}'`);
    all = stdout ? JSON.parse(stdout) : {};
  }

  if (true) {
    const {stdout} = await exec(`docker ps -f name=${metadata.market.hooks.name} --format '{{json .}}'`);
    running = stdout ? JSON.parse(stdout) : {};
  }

  return [all, running];
}
exports.queryContainer = queryContainer;

async function startContainer() {
  console.log(`Thread #market: start container`);

  const privateIPv4 = await utils.discoverPrivateIPv4();
  const envFile = `${process.cwd()}/.env`;
  const dockerArgs = `-d --restart always --privileged -it --name ${metadata.market.hooks.name} \\
    --add-host=mgmt.srs.local:${privateIPv4.address} \\
    -v ${envFile}:/srs-terraform/hooks/.env \\
    -p ${metadata.market.hooks.port}:${metadata.market.hooks.port} \\
    --log-driver json-file --log-opt max-size=1g --log-opt max-file=3 \\
    ${metadata.market.hooks.image} \\
    node .`;
  console.log(`Thread #market: docker run args ip=${privateIPv4.name}/${privateIPv4.address}, docker run ${dockerArgs}`);

  // Only remove the container when got ID, to avoid fail for CentOS.
  const all = (await queryContainer())[0];
  if (all && all.ID) {
    await exec(`docker rm -f ${metadata.market.hooks.name}`);
    console.log(`Thread #market: docker run remove ID=${all.ID}`);
  }

  await exec(`docker run ${dockerArgs}`);
  console.log(`Thread #market: docker run ok`);
}
