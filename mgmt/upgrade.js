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
const { spawn } = require('child_process');
const pkg = require('./package.json');
const axios = require('axios');
const semver = require('semver');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const util = require('util');
const metadata = require('./metadata');
const platform = require('./platform');
const keys = require('js-core/keys');
const {queryLatestVersion} = require('./releases');
const moment = require('moment');

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  // We must initialize the thread first.
  const {region, registry} = await platform.init();

  const srs = await metadata.market.srs.image();
  await redis.hset(keys.redis.SRS_TENCENT_LH, 'srs', srs);

  console.log(`Thread #${metadata.upgrade.name}: initialize region=${region}, registry=${registry}, srs=${srs}`);

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
  // Wait for a while to request version.
  console.log(`Thread #${metadata.upgrade.name}: query current version=v${pkg.version}`);
  await new Promise(resolve => setTimeout(resolve, 1 * 1000));
  console.log(`Thread #${metadata.upgrade.name}: query request by version=v${pkg.version}`);

  // For development, request the releases from itself which proxy to the releases service.
  const releases = await queryLatestVersion(redis, axios);
  metadata.upgrade.releases = releases;
  console.log(`Thread #${metadata.upgrade.name}: query done, version=v${pkg.version}, response=${JSON.stringify(releases)}`);

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

    const uwStart = await redis.hget(keys.redis.SRS_UPGRADE_WINDOW, 'start');
    const uwDuration = await redis.hget(keys.redis.SRS_UPGRADE_WINDOW, 'duration');
    if (!helper.inUpgradeWindow(uwStart, uwDuration, moment())) {
      console.log(`Thread #${metadata.upgrade.name}: Ignore for not in window`);
      return;
    }

    const r0 = await redis.hget(keys.redis.SRS_UPGRADING, 'upgrading');
    if (r0 === "1") {
      const r1 = await redis.hget(keys.redis.SRS_UPGRADING, 'desc');
      console.log(`Thread #${metadata.upgrade.name}: already upgrading r0=${r0} ${r1}`);
      return;
    }

    const r2 = await redis.hget(keys.redis.SRS_UPGRADE_STRATEGY, 'strategy');
    const strategy = r2 || 'auto';
    if (strategy !== 'auto') {
      const r3 = await redis.hget(keys.redis.SRS_UPGRADE_STRATEGY, 'desc');
      console.log(`Thread #${metadata.upgrade.name}: ignore for strategy=${r2}/${strategy} ${r3}`);
      return;
    }

    // Set the upgrading to avoid others.
    await redis.hset(keys.redis.SRS_UPGRADING, 'upgrading', 1);
    await redis.hset(keys.redis.SRS_UPGRADING, 'desc', `${upgradingMessage}`);

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

