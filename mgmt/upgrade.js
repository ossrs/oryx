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
const platform = require('./platform');
const COS = require('cos-nodejs-sdk-v5');
const cos = require('js-core/cos');
const keys = require('js-core/keys');

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  // We must initialize the thread first.
  const {region, registry} = await platform.init();
  console.log(`Thread #${metadata.upgrade.name}: initialize region=${region}, registry=${registry}`);

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
  // Run only once for each process.
  await resetUpgrading();

  const region = await platform.region();
  await cos.createCosBucket(redis, COS, region);

  // Run only once for a special version.
  await firstRun();

  // Wait for a while to request version.
  console.log(`Thread #${metadata.upgrade.name}: query current version=v${pkg.version}`);
  await new Promise(resolve => setTimeout(resolve, 1 * 1000));
  console.log(`Thread #${metadata.upgrade.name}: query request by version=v${pkg.version}`);

  // For development, request the releases from itself which proxy to the releases service.
  const releaseServer = process.env.LOCAL_RELEASE === 'true' ? `http://localhost:${consts.config.port}` : 'http://api.ossrs.net';
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
    if (r0 === "1") {
      const r1 = await redis.hget(consts.SRS_UPGRADING, 'desc');
      console.log(`Thread #${metadata.upgrade.name}: already upgrading r0=${r0} ${r1}`);
      return;
    }

    const r2 = await redis.hget(consts.SRS_UPGRADE_STRATEGY, 'strategy');
    const strategy = r2 || 'auto';
    if (strategy !== 'auto') {
      const r3 = await redis.hget(consts.SRS_UPGRADE_STRATEGY, 'desc');
      console.log(`Thread #${metadata.upgrade.name}: ignore for strategy=${r2}/${strategy} ${r3}`);
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
  // For each init stage changed, we could use a different redis key, to identify this special init workflow.
  // However, keep in mind that previous defined workflow always be executed, so these operations should be idempotent.
  // History:
  //    SRS_FIRST_BOOT_DONE, For release 4.1, to restart srs.
  //    SRS_FIRST_BOOT_DONE_v1, For release 4.2, to restart srs, exec upgrade_prepare.
  //    SRS_FIRST_BOOT_DONE_v2, For release 4.2, to restart srs-server, update the hls hooks.
  //    SRS_FIRST_BOOT.v3, For current release, to restart srs-hooks, update the volumes.
  //    SRS_FIRST_BOOT.v4, For current release, to restart tencent-cloud, update the volumes.
  const SRS_FIRST_BOOT = keys.redis.SRS_FIRST_BOOT;
  const bootRelease = 'v4';

  // Run once, record in redis.
  const r0 = await redis.hget(SRS_FIRST_BOOT, bootRelease);
  await redis.hset(SRS_FIRST_BOOT, bootRelease, 1);
  if (r0) {
    console.log(`Thread #${metadata.upgrade.name}: boot already done, r0=${r0}`);
    return false;
  }

  // To prevent boot again and again.
  console.log(`Thread #${metadata.upgrade.name}: boot start to setup, v=${SRS_FIRST_BOOT}.${bootRelease}, r0=${r0}`);

  // For the second time, we prepare the os, for previous version which does not run the upgrade living.
  console.log(`Thread #${metadata.upgrade.name}: Prepare OS for first run, r0=${r0}`);
  await exec(`bash auto/upgrade_prepare`);

  // Remove containers.
  const removeContainer = async (name) => {
    try {
      await exec(`docker rm -f ${name}`);
      console.log(`Thread #${metadata.upgrade.name}: boot remove docker ${name} ok`);
    } catch (e) {
      console.log(`Thread #${metadata.upgrade.name}: boot remove docker ${name}, ignore err ${e}`);
    }
  };
  await removeContainer(metadata.market.srs.name);
  await removeContainer(metadata.market.hooks.name);
  await removeContainer(metadata.market.tencent.name);

  console.log(`Thread #${metadata.upgrade.name}: boot done`);
  return true;
}

async function resetUpgrading() {
  // When restart, reset the upgrading.
  const r1 = await redis.hget(consts.SRS_UPGRADING, 'upgrading');
  if (r1) {
    const r2 = await redis.hget(consts.SRS_UPGRADING, 'desc');
    const r3 = await redis.del(consts.SRS_UPGRADING);
    console.log(`Thread #${metadata.upgrade.name}: reset upgrading for r1=${r1}, r2=${r2}, r3=${r3}`);
  }
}

