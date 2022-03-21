'use strict';

// For components in docker, connect by host.
const config = {
  redis:{
    host: process.env.NODE_ENV === 'development' ? 'localhost' : 'mgmt.srs.local',
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const { isMainThread } = require("worker_threads");
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const platform = require('./platform');
const keys = require('js-core/keys');
const COS = require('cos-nodejs-sdk-v5');
const cos = require('js-core/cos');
const vod = require('js-core/vod');
const {AbstractClient} = require('./sdk-internal/common/abstract_client');
const VodClient = require("tencentcloud-sdk-nodejs").vod.v20180717.Client;
const utils = require('js-core/utils');
const { v4: uuidv4 } = require('uuid');
const helper = require('./helper');
const metadata = require('js-core/metadata');
const moment = require('moment');

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  // We must initialize the thread first.
  await platform.init();
  console.log(`Thread #upgrade: initialize`);

  while (true) {
    try {
      await doThreadMain();
    } catch (e) {
      console.error(`Thread #upgrade: err`, e);
    } finally {
      await new Promise(resolve => setTimeout(resolve, 3600 * 1000));
    }
  }
}

async function doThreadMain() {
  // Run only once for each process.
  await resetUpgrading();

  // Setup the upgrade window if not set.
  await setupUpgradeWindow();

  // Try to create cloud service, we ignore error, because user might delete the secret directly on console of cloud
  // platform, so it might throw exception when create cloud resource.
  try {
    const region = await redis.hget(keys.redis.SRS_TENCENT_LH, 'region');
    await cos.createCosBucket(redis, COS, region);
    await vod.createVodService(redis, VodClient, AbstractClient, region);
  } catch (e) {
    console.warn(`Thread #upgrade: Ignore cloud service err`, e);
  }

  // Run only once for a special version.
  await firstRun();
}

async function resetUpgrading() {
  // When restart, reset the upgrading.
  const r1 = await redis.hget(keys.redis.SRS_UPGRADING, 'upgrading');
  if (r1) {
    const r2 = await redis.hget(keys.redis.SRS_UPGRADING, 'desc');
    const r3 = await redis.del(keys.redis.SRS_UPGRADING);
    console.log(`Thread #upgrade: reset upgrading for r1=${r1}, r2=${r2}, r3=${r3}`);
  }
}

async function setupUpgradeWindow() {
  // If user not setup it, we will set the default value for each time, because we could change it apparently.
  const update = await redis.hget(keys.redis.SRS_UPGRADE_WINDOW, 'update');
  if (!update) {
    await redis.hset(keys.redis.SRS_UPGRADE_WINDOW, 'start', 23);
    await redis.hset(keys.redis.SRS_UPGRADE_WINDOW, 'duration', 6);
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
  //    SRS_FIRST_BOOT.v5, For current release, restart containers for .env changed.
  const SRS_FIRST_BOOT = keys.redis.SRS_FIRST_BOOT;
  const bootRelease = 'v11';

  // Run once, record in redis.
  const r0 = await redis.hget(SRS_FIRST_BOOT, bootRelease);
  await redis.hset(SRS_FIRST_BOOT, bootRelease, 1);
  if (r0) {
    console.log(`Thread #upgrade: boot already done, r0=${r0}`);
    return false;
  }

  // To prevent boot again and again.
  console.log(`Thread #upgrade: boot start to setup, v=${SRS_FIRST_BOOT}.${bootRelease}, r0=${r0}`);

  // For the second time, we prepare the os, for previous version which does not run the upgrade living.
  console.log(`Thread #upgrade: Prepare OS for first run, r0=${r0}`);
  await helper.execApi('executeUpgradePrepare');

  // Setup the api secret.
  const [token, created] = await utils.setupApiSecret(redis, uuidv4, moment);
  console.log(`Thread #upgrade: Platform api secret, token=${token.length}B, created=${created}`);

  // Remove containers.
  await helper.execApi('rmContainer', [metadata.market.srs.name]);
  await helper.execApi('rmContainer', [metadata.market.hooks.name]);
  await helper.execApi('rmContainer', [metadata.market.tencent.name]);
  await helper.execApi('rmContainer', [metadata.market.ffmpeg.name]);

  console.log(`Thread #upgrade: boot done`);
  return true;
}

