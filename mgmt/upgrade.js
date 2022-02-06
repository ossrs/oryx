'use strict';

// For mgmt, it's ok to connect to localhost.
const config = {
  redis:{
    host: 'localhost',
    port: 6379,
    password: '',
  },
};

const { isMainThread, parentPort } = require("worker_threads");
const { spawn } = require('child_process');
const pkg = require('./package.json');
const axios = require('axios');
const semver = require('semver');
const consts = require('./consts');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const util = require('util');
const exec = util.promisify(require('child_process').exec);
const metadata = require('./metadata');

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  while (true) {
    try {
      await doThreadMain();
    } catch (e) {
      console.error(`Thread #${metadata.upgrade.name}: err`, e);
    } finally {
      await new Promise(resolve => setTimeout(resolve, 3600 * 1000));
    }
  }
}

async function doThreadMain() {
  await firstRun();

  // Wait for a while to request version.
  console.log(`Thread #${metadata.upgrade.name}: query current version=v${pkg.version}`);
  await new Promise(resolve => setTimeout(resolve, 1 * 1000));
  console.log(`Thread #${metadata.upgrade.name}: query request by version=v${pkg.version}`);

  // For development, request the releases from itself which proxy to the releases service.
  const releaseServer = process.env.NODE_ENV === 'development' ? `http://localhost:${consts.config.port}` : 'http://api.ossrs.net';
  const {data} = await axios.get(`${releaseServer}/terraform/v1/releases`, {
    params: {
      version: `v${pkg.version}`,
      ts: new Date().getTime(),
    }
  });
  metadata.upgrade.releases = data;
  console.log(`Thread #${metadata.upgrade.name}: query done, version=v${pkg.version}, response=${JSON.stringify(data)}`);

  // Update the metadata to main thread.
  parentPort.postMessage({
    metadata: {
      upgrade: metadata.upgrade,
    },
  });

  // Try to upgrade terraform itself.
  const higherStable = semver.lt(`v${pkg.version}`, metadata.upgrade.releases.stable);
  if (metadata.upgrade.releases && metadata.upgrade.releases.stable && higherStable) {
    const upgradingMessage = `upgrade from v${pkg.version} to stable ${metadata.upgrade.releases.stable}`;
    console.log(`Thread #${metadata.upgrade.name}: ${upgradingMessage}`);

    const r0 = await redis.hget(consts.SRS_UPGRADING, 'upgrading');
    if (r0) {
      const r1 = await redis.hget(consts.SRS_UPGRADING, 'desc');
      console.log(`Thread #${metadata.upgrade.name}: already upgrading r0=${r0} ${r1}`);
      return;
    }

    // Set the upgrading to avoid others.
    await redis.hset(consts.SRS_UPGRADING, 'upgrading', 1);
    await redis.hset(consts.SRS_UPGRADING, 'desc', `${upgradingMessage}`);

    await new Promise((resolve, reject) => {
      const child = spawn('bash', ['upgrade', metadata.upgrade.releases.stable]);
      child.stdout.on('data', (chunk) => {
        console.log(`Thread #${metadata.upgrade.name}: upgrade ${chunk.toString()}`);
      });
      child.stderr.on('data', (chunk) => {
        console.log(`Thread #${metadata.upgrade.name}: upgrade ${chunk.toString()}`);
      });
      child.on('close', (code) => {
        console.log(`Thread #${metadata.upgrade.name}: upgrade exited with code ${code}`);
        if (code !== 0) return reject(code);
        resolve();
      });
    });
  }
}

async function firstRun() {
  // When restart, reset the upgrading.
  const r1 = await redis.hget(consts.SRS_UPGRADING, 'upgrading');
  if (r1) {
    const r2 = await redis.hget(consts.SRS_UPGRADING, 'desc');
    const r3 = await redis.del(consts.SRS_UPGRADING);
    console.log(`Thread #${metadata.upgrade.name}: reset upgrading for r1=${r1}, r2=${r2}, r3=${r3}`);
  }

  // For each init stage changed, we could use a different redis key, to identify this special init workflow.
  // However, keep in mind that previous defined workflow always be executed, so these operations should be idempotent.
  // History:
  //    SRS_FIRST_BOOT_DONE, For release 4.1, to restart srs.
  //    SRS_FIRST_BOOT_DONE_v1, For current release, to restart srs, exec upgrade_prepare.
  const SRS_FIRST_BOOT_DONE = 'SRS_FIRST_BOOT_DONE_v1';

  // Run once, record in redis.
  const r0 = await redis.get(SRS_FIRST_BOOT_DONE);
  await redis.set(SRS_FIRST_BOOT_DONE, 1);
  if (r0) {
    console.log(`Thread #${metadata.upgrade.name}: boot already done, r0=${r0}`);
    return false;
  }

  // To prevent boot again and again.
  console.log(`Thread #${metadata.upgrade.name}: boot start to setup, v=${SRS_FIRST_BOOT_DONE}, r0=${r0}`);

  // For the second time, we prepare the os, for previous version which does not run the upgrade living.
  console.log(`Thread #${metadata.upgrade.name}: Prepare OS for first run, r0=${r0}`);
  await exec(`bash upgrade_prepare`);

  try {
    // Because we already create the container, and cached the last SRS 4.0 image, also set the hosts for hooks by
    // --add-host which is incorrect for new machine, so we must delete the container and restart it when first run.
    await exec(`docker rm -f ${metadata.market.srs.name}`);
    console.log(`Thread #${metadata.upgrade.name}: boot remove docker ${metadata.market.srs.name}`);
  } catch (e) {
    console.log(`Thread #${metadata.upgrade.name}: boot ignore rm docker ${metadata.market.srs.name} error ${e.message}`);
  }

  console.log(`Thread #${metadata.upgrade.name}: boot done`);
  return true;
}

