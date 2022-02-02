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

// The redis key we used.
const SRS_FIRST_BOOT_DONE = 'SRS_FIRST_BOOT_DONE';

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  while (true) {
    try {
      await doThreadMain();
    } catch (e) {
      console.error(`Thread #${metadata.releases.name}: err`, e);
    } finally {
      await new Promise(resolve => setTimeout(resolve, 3600 * 1000));
    }
  }
}

async function doThreadMain() {
  await firstRun();

  console.log(`Thread #${metadata.releases.name}: query current version=v${pkg.version}`);

  // Wait for a while to request version.
  await new Promise(resolve => setTimeout(resolve, 1 * 1000));
  console.log(`Thread #${metadata.releases.name}: query request by version=v${pkg.version}`);

  // For development, request the releases from itself which proxy to the releases service.
  const releaseServer = process.env.NODE_ENV === 'development' ? `http://localhost:${consts.config.port}` : 'http://api.ossrs.net';
  const {data} = await axios.get(`${releaseServer}/terraform/v1/releases`, {
    params: {
      version: `v${pkg.version}`,
      ts: new Date().getTime(),
    }
  });
  metadata.releases.releases = data;
  console.log(`Thread #${metadata.releases.name}: query done, version=v${pkg.version}, response=${JSON.stringify(data)}`);

  // Update the metadata to main thread.
  parentPort.postMessage({
    metadata: {
      releases: metadata.releases,
    },
  });

  // Try to upgrade terraform itself.
  const higherStable = semver.lt(`v${pkg.version}`, metadata.releases.releases.stable);
  if (metadata.releases.releases && metadata.releases.releases.stable && higherStable) {
    console.log(`Thread #${metadata.releases.name}: upgrade from v${pkg.version} to stable ${metadata.releases.releases.stable}`);

    await new Promise((resolve, reject) => {
      const child = spawn('bash', ['upgrade', metadata.releases.releases.stable]);
      child.stdout.on('data', (chunk) => {
        console.log(`Thread #${metadata.releases.name}: upgrade ${chunk.toString()}`);
      });
      child.stderr.on('data', (chunk) => {
        console.log(`Thread #${metadata.releases.name}: upgrade ${chunk.toString()}`);
      });
      child.on('close', (code) => {
        console.log(`Thread #${metadata.releases.name}: upgrade exited with code ${code}`);
        if (code !== 0) return reject(code);
        resolve();
      });
    });
  }
}

async function firstRun() {
  const r0 = await redis.get(SRS_FIRST_BOOT_DONE);
  await redis.set(SRS_FIRST_BOOT_DONE, r0 ? parseInt(r0) + 1 : 1);

  // We do the first run for the first N times.
  // 1. The first time, for the startup version in image, to init the whole system.
  // 2. The second time, upgraded to stable version, to do additional init, such as change container args.
  if (parseInt(r0) >= 2) {
    console.log(`Thread #${metadata.releases.name}: boot already done, r0=${r0}`);
    return;
  }

  // To prevent boot again and again.
  console.log(`Thread #${metadata.releases.name}: boot start to setup`);

  try {
    // Because we already create the container, and cached the last SRS 4.0 image, also set the hosts for hooks by
    // --add-host which is incorrect for new machine, so we must delete the container and restart it when first run.
    await exec(`docker rm -f ${metadata.market.srs.name}`);
    console.log(`Thread #${metadata.releases.name}: boot remove docker ${metadata.market.srs.name}`);
  } catch (e) {
    console.log(`Thread #${metadata.releases.name}: boot ignore rm docker ${metadata.market.srs.name} error ${e.message}`);
  }

  console.log(`Thread #${metadata.releases.name}: boot done`);
}

