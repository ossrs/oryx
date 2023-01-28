'use strict';

// For components in docker, connect by host.
const config = {
  node: {
    host: process.env.NODE_ENV === 'development' ? 'localhost' : 'mgmt.srs.local',
  },
  redis:{
    host: process.env.NODE_ENV === 'development' ? 'localhost' : (process.env.REDIS_HOST || 'mgmt.srs.local'),
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const { isMainThread, parentPort } = require("worker_threads");
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const keys = require('js-core/keys');
const moment = require('moment');
const { spawn } = require('child_process');
const { v4: uuidv4 } = require('uuid');

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  // We must initialize the thread first.
  console.log(`Thread #vFileWorker: initialize`);

  parentPort.on('message', (msg) => {
    handleMessage(msg);
  });

  while (true) {
    try {
      await doThreadMain();
    } catch (e) {
      console.error(`Thread #vFileWorker: err`, e);
    } finally {
      await new Promise(resolve => setTimeout(resolve, 30 * 1000));
    }
  }
}

async function handleMessage(msg) {
  console.error(`Thread #vFileWorker: Ignore msg ${JSON.stringify(msg)}`);
}

async function doThreadMain() {
  while (true) {
    await generateVLiveRules();
    await handleVLiveTasks();
    await new Promise(resolve => setTimeout(resolve, 3000));
  }
}

async function generateVLiveRules() {
  const configs = await redis.hgetall(keys.redis.SRS_VLIVE_CONFIG);
  let nnPlatformEnabled = 0;
  for (const k in configs) {
    const config = JSON.parse(configs[k]);
    if (config.enabled) nnPlatformEnabled++;
    configs[k] = config;
  }

  const platforms = nnPlatformEnabled ? Object.values(configs).map(e => `${e.platform}:${e.enabled}`) : Object.keys(configs);
  console.log(`Thread #vFileWorker: Active channels enabled=${nnPlatformEnabled}, platforms=${JSON.stringify(platforms)}`);

  if (!nnPlatformEnabled) return;

  const logs = [];
  for (const k in configs) {
    const configObj = configs[k];
    if (!configObj?.enabled) continue;
    if (!configObj?.files?.length) continue;

    // The active or selected source file.
    const sourceObj = configObj.files[0];

    // Update the map, to find by source uuid. Refresh if uuid changed.
    const map = await redis.hget(keys.redis.SRS_VLIVE_MAP, k);
    if (map && map === sourceObj.uuid) continue;
    const r1 = await redis.hset(keys.redis.SRS_VLIVE_MAP, k, sourceObj.uuid);

    // Update the stream, to iterate all vLive tasks.
    const vLive = await redis.hget(keys.redis.SRS_VLIVE_STREAM, `${k}@${sourceObj.uuid}`);
    const vLiveObj = vLive ? JSON.parse(vLive) : {
      uuid: uuidv4(),
      platform: k,
      input: sourceObj.target,
    };

    // Update the vLive object of stream.
    vLiveObj.stream = sourceObj;
    vLiveObj.update = moment().format();

    const r0 = await redis.hset(keys.redis.SRS_VLIVE_STREAM, `${k}@${sourceObj.uuid}`, JSON.stringify(vLiveObj));
    logs.push(`${k}@${sourceObj.uuid}, platform=${k}, uuid=${vLiveObj.uuid} enabled=${configObj?.enabled}, r0=${r0}, r1=${r1}`);
  }

  if (logs?.length) console.log(`Thread #vFileWorker: Update the vLives, map is ${JSON.stringify(logs)}`);
}

async function handleVLiveTasks() {
  const activeKeys = await redis.hkeys(keys.redis.SRS_VLIVE_STREAM);
  if (!activeKeys || !activeKeys.length) return;

  let skipNext;
  for (const i in activeKeys) {
    // Get the local object by the active key.
    const activeKey = activeKeys[i];
    const vLive = await redis.hget(keys.redis.SRS_VLIVE_STREAM, activeKey);
    const vLiveObj = vLive && JSON.parse(vLive);
    if (!vLiveObj) continue;

    // Detect the enabled again.
    const platform = await redis.hget(keys.redis.SRS_VLIVE_CONFIG, vLiveObj.platform);
    const configObj = platform && JSON.parse(platform);
    if (!configObj) continue;

    // Remove task if disabled.
    skipNext = await removeDisabledTask(activeKey, vLiveObj, configObj);
    if (skipNext) continue;

    // Start new task if not exists.
    skipNext = await startNewTask(activeKey, vLiveObj, configObj);
    if (skipNext) continue;

    // Restart the dead task.
    skipNext = await restartDeadTasks(activeKey, vLiveObj, configObj);
    if (skipNext) continue;

    // Terminate task if no stream.
    skipNext = await terminateTaskForNoStream(activeKey, vLiveObj, configObj);
    if (skipNext) continue;
  }
}

function generateOutput(svr, secret) {
  const server = svr.trim().replace(/localhost/g, config.node.host);
  const seperator = (server.endsWith('/') || !secret || secret.startsWith('/')) ? '' : '/';
  return `${server}${seperator}${secret || ''}`;
}
exports.generateOutput = generateOutput;

async function startNewTask(activeKey, vLiveObj, configObj) {
  // Ignore if not enabled.
  if (!configObj.enabled) return;

  // Ignore if already generated.
  // Note that if worker pid changed, we should regenerate it.
  if (vLiveObj.task && vLiveObj.pid === process.pid) return;

  // Build the output stream url.
  vLiveObj.output = generateOutput(configObj.server, configObj.secret);

  // Start a child process to vLive stream.
  const child = spawn('ffmpeg', [
    '-stream_loop', '-1', '-re', '-i', vLiveObj.input, '-c', 'copy', '-f', 'flv', vLiveObj.output,
  ]);
  vLiveObj.task = child.pid;
  const previousPid = vLiveObj.pid;
  vLiveObj.pid = process.pid;
  const r0 = await redis.hset(keys.redis.SRS_VLIVE_STREAM, activeKey, JSON.stringify(vLiveObj));
  const r1 = await redis.hdel(keys.redis.SRS_VLIVE_CODE, activeKey);
  const r2 = await redis.hdel(keys.redis.SRS_VLIVE_FRAME, activeKey);
  console.log(`Thread #vFileWorker: Start task=${child.pid}, pid=${previousPid}/${vLiveObj.pid}, stream=${activeKey}, input=${vLiveObj.input}, output=${vLiveObj.output}, r0=${r0}, r1=${r1}, r2=${r2}`);

  let nnLogs = 0;
  child.stdout.on('data', (chunk) => {
    console.log(chunk.toString());
  });
  child.stderr.on('data', (chunk) => {
    if ((nnLogs++ % 3) !== 0) return;
    const log = chunk.toString().trim().replace(/= +/g, '=');
    if (log.indexOf('frame=') < 0) return;
    redis.hset(keys.redis.SRS_VLIVE_FRAME, activeKey, JSON.stringify({log, update: moment().format()}));
    console.log(`Thread #vFileWorker: Active task=${child.pid}, stream=${activeKey}, ${log}`);
  });
  child.on('close', (code) => {
    redis.hset(keys.redis.SRS_VLIVE_CODE, activeKey, JSON.stringify({close: true, code, update: moment().format()}));
    console.log(`Thread #vFileWorker: Close task=${child.pid}, stream=${activeKey}, code=${code}`);
  });
}

async function restartDeadTasks(activeKey, vLiveObj, configObj) {
  // Ignore if not enabled.
  if (!configObj.enabled) return;

  // Ignore if not started.
  if (!vLiveObj.task) return;

  // Check whether exited.
  const code = await redis.hget(keys.redis.SRS_VLIVE_CODE, activeKey);
  if (!code) return;

  const codeObj = JSON.parse(code);
  if (!codeObj?.close) return;

  // Reset the task.
  const previousPid = vLiveObj.task;
  vLiveObj.task = null;
  const r0 = await redis.hset(keys.redis.SRS_VLIVE_STREAM, activeKey, JSON.stringify(vLiveObj));
  console.log(`Thread #vFileWorker: Reset task=${previousPid}, stream=${activeKey}, code=${codeObj.code}, at=${codeObj?.update}, r0=${r0}`);
}

async function removeDisabledTask(activeKey, vLiveObj, configObj) {
  // Ignore if not disabled.
  if (configObj.enabled) return;

  // Check whether exited.
  // If process changed, ignore.
  while (vLiveObj.pid === process.pid) {
    const code = await redis.hget(keys.redis.SRS_VLIVE_CODE, activeKey);
    const codeObj = code && JSON.parse(code);
    if (codeObj?.close) break;

    // Ignore if process not exists. Because if not exists, the code never update.
    // See https://stackoverflow.com/a/21296291/17679565
    try {
      process.kill(vLiveObj.task, 0);
    } catch (e) {
      console.warn(`Thread #vFileWorker: Ignore task=${vLiveObj.task} for process not exists`);
      break;
    }

    // Kill the process if not exited.
    try {
      process.kill(vLiveObj.task, 'SIGKILL');
    } catch (e) {
    }
    console.log(`Thread #vFileWorker: Kill task=${vLiveObj.task}, stream=${activeKey}, code=${codeObj?.code}, at=${codeObj?.update}`);

    await new Promise(resolve => setTimeout(resolve, 800));
  }

  // Cleanup redis.
  const r0 = await redis.hdel(keys.redis.SRS_VLIVE_STREAM, activeKey);
  const r1 = await redis.hdel(keys.redis.SRS_VLIVE_FRAME, activeKey);
  const r2 = await redis.hdel(keys.redis.SRS_VLIVE_CODE, activeKey);
  const r3 = await redis.hdel(keys.redis.SRS_VLIVE_MAP, vLiveObj.platform);
  console.log(`Thread #vFileWorker: Cleanup task=${vLiveObj.task}, platform=${vLiveObj.platform}, stream=${activeKey}, r0=${r0}, r1=${r1}, r2=${r2}, r3=${r3}`);

  // Skip any options.
  return true;
}

async function terminateTaskForNoStream(activeKey, vLiveObj, configObj) {
  // Ignore if not enabled.
  if (!configObj.enabled) return;

  // Ignore if not started.
  if (!vLiveObj.task) return;

  // Check whether exited.
  const frame = await redis.hget(keys.redis.SRS_VLIVE_FRAME, activeKey);
  if (!frame) return;

  const frameObj = JSON.parse(frame);
  if (!frameObj?.update) return;

  // If not expired, ignore.
  const expired = moment(frameObj.update).add(process.env.NODE_ENV === 'development' ? 9 : 30, 's');
  if (expired.isAfter(moment())) return;

  // Expired, terminate the task.
  configObj.enabled = false;
  console.log(`Thread #vFileWorker: Expire task=${vLiveObj.task}, platform=${vLiveObj.platform}, stream=${activeKey}, update=${vLiveObj.update}, expired=${expired.format()}, now=${moment().format()}`);

  return await removeDisabledTask(activeKey, vLiveObj, configObj);
}
