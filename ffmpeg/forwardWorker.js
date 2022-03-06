'use strict';

// For components in docker, connect by host.
const config = {
  redis:{
    host: process.env.NODE_ENV === 'development' ? 'localhost' : 'mgmt.srs.local',
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
  console.log(`Thread #forwardWorker: initialize`);

  parentPort.on('message', (msg) => {
    handleMessage(msg);
  });

  while (true) {
    try {
      await doThreadMain();
    } catch (e) {
      console.error(`Thread #forwardWorker: err`, e);
    } finally {
      await new Promise(resolve => setTimeout(resolve, 30 * 1000));
    }
  }
}

async function handleMessage(msg) {
  console.error(`Thread #forwardWorker: Ignore msg ${JSON.stringify(msg)}`);
}

async function doThreadMain() {
  while (true) {
    await generateForwardRules();
    await handleForwardTasks();
    await new Promise(resolve => setTimeout(resolve, 3000));
  }
}

async function generateForwardRules() {
  const configs = await redis.hgetall(keys.redis.SRS_FORWARD_CONFIG);
  let nnPlatformEnabled = 0;
  for (const k in configs) {
    const config = JSON.parse(configs[k]);
    if (config.enabled) nnPlatformEnabled++;
    configs[k] = config;
  }

  const activeKeys = await redis.hkeys(keys.redis.SRS_STREAM_ACTIVE);
  const platforms = nnPlatformEnabled ? Object.values(configs).map(e => `${e.platform}:${e.enabled}`) : Object.keys(configs);
  console.log(`Thread #forwardWorker: Active streams ${JSON.stringify(activeKeys)}, enabled=${nnPlatformEnabled}, platforms=${JSON.stringify(platforms)}`);

  if (!nnPlatformEnabled) return;
  if (!activeKeys || !activeKeys.length) return;

  for (const i in activeKeys) {
    // Get the local object by the active key.
    const activeKey = activeKeys[i];
    const stream = await redis.hget(keys.redis.SRS_STREAM_ACTIVE, activeKey);
    const streamObj = stream ? JSON.parse(stream) : null;
    if (!streamObj) continue;

    const logs = [];
    for (const k in configs) {
      const configObj = configs[k];
      if (!configObj?.enabled) continue;

      const map = await redis.hget(keys.redis.SRS_FORWARD_MAP, k);
      if (map) continue;
      const r1 = await redis.hset(keys.redis.SRS_FORWARD_MAP, k, activeKey);

      const forward = await redis.hget(keys.redis.SRS_FORWARD_STREAM, `${k}@${activeKey}`);
      const forwardObj = forward ? JSON.parse(forward) : {
        uuid: uuidv4(),
        platform: k,
        input: `rtmp://${config.redis.host}/${streamObj.app}/${streamObj.stream}`,
      };

      // Update the forward object.
      forwardObj.stream = streamObj;
      forwardObj.update = moment().format();

      const r0 = await redis.hset(keys.redis.SRS_FORWARD_STREAM, `${k}@${activeKey}`, JSON.stringify(forwardObj));
      logs.push(`${k}@${activeKey}, platform=${k}, uuid=${forwardObj.uuid} enabled=${configObj?.enabled}, r0=${r0}, r1=${r1}`);
    }

    if (logs?.length) console.log(`Thread #forwardWorker: Update the forwards, activeKey=${activeKey}, map is ${JSON.stringify(logs)}`);
  }
}

async function handleForwardTasks() {
  const activeKeys = await redis.hkeys(keys.redis.SRS_FORWARD_STREAM);
  if (!activeKeys || !activeKeys.length) return;

  let skipNext;
  for (const i in activeKeys) {
    // Get the local object by the active key.
    const activeKey = activeKeys[i];
    const forward = await redis.hget(keys.redis.SRS_FORWARD_STREAM, activeKey);
    const forwardObj = forward && JSON.parse(forward);
    if (!forwardObj) continue;

    // Detect the enabled again.
    const platform = await redis.hget(keys.redis.SRS_FORWARD_CONFIG, forwardObj.platform);
    const configObj = platform && JSON.parse(platform);
    if (!configObj) continue;

    // Remove task if disabled.
    skipNext = await removeDisabledTask(activeKey, forwardObj, configObj);
    if (skipNext) continue;

    // Start new task if not exists.
    skipNext = await startNewTask(activeKey, forwardObj, configObj);
    if (skipNext) continue;

    // Restart the dead task.
    skipNext = await restartDeadTasks(activeKey, forwardObj, configObj);
    if (skipNext) continue;

    // Terminate task if no stream.
    skipNext = await terminateTaskForNoStream(activeKey, forwardObj, configObj);
    if (skipNext) continue;
  }
}

async function startNewTask(activeKey, forwardObj, configObj) {
  // Ignore if not enabled.
  if (!configObj.enabled) return;

  // Ignore if already generated.
  // Note that if worker pid changed, we should regenerate it.
  if (forwardObj.task && forwardObj.pid === process.pid) return;

  // Build the output stream url.
  const server = configObj.server.trim();
  const seperator = (server.endsWith('/') || configObj.secret.startsWith('/') || !configObj.secret) ? '' : '/';
  forwardObj.output = `${server}${seperator}${configObj.secret}`;

  // Start a child process to forward stream.
  const child = spawn('ffmpeg', ['-f', 'flv', '-i', forwardObj.input, '-c', 'copy', '-f', 'flv', forwardObj.output]);
  forwardObj.task = child.pid;
  const previousPid = forwardObj.pid;
  forwardObj.pid = process.pid;
  const r0 = await redis.hset(keys.redis.SRS_FORWARD_STREAM, activeKey, JSON.stringify(forwardObj));
  const r1 = await redis.hdel(keys.redis.SRS_FORWARD_CODE, activeKey);
  const r2 = await redis.hdel(keys.redis.SRS_FORWARD_FRAME, activeKey);
  console.log(`Thread #forwardWorker: Start task=${child.pid}, pid=${previousPid}/${forwardObj.pid}, stream=${activeKey}, input=${forwardObj.input}, output=${forwardObj.output}, r0=${r0}, r1=${r1}, r2=${r2}`);

  let nnLogs = 0;
  child.stdout.on('data', (chunk) => {
    console.log(chunk.toString());
  });
  child.stderr.on('data', (chunk) => {
    if ((nnLogs++ % 3) !== 0) return;
    const log = chunk.toString().trim().replace(/= +/g, '=');
    if (log.indexOf('frame=') < 0) return;
    redis.hset(keys.redis.SRS_FORWARD_FRAME, activeKey, JSON.stringify({log, update: moment().format()}));
    console.log(`Thread #forwardWorker: Active task=${child.pid}, stream=${activeKey}, ${log}`);
  });
  child.on('close', (code) => {
    redis.hset(keys.redis.SRS_FORWARD_CODE, activeKey, JSON.stringify({close: true, code, update: moment().format()}));
    console.log(`Thread #forwardWorker: Close task=${child.pid}, stream=${activeKey}, code=${code}`);
  });
}

async function restartDeadTasks(activeKey, forwardObj, configObj) {
  // Ignore if not enabled.
  if (!configObj.enabled) return;

  // Ignore if not started.
  if (!forwardObj.task) return;

  // Check whether exited.
  const code = await redis.hget(keys.redis.SRS_FORWARD_CODE, activeKey);
  if (!code) return;

  const codeObj = JSON.parse(code);
  if (!codeObj?.close) return;

  // Reset the task.
  const previousPid = forwardObj.task;
  forwardObj.task = null;
  const r0 = await redis.hset(keys.redis.SRS_FORWARD_STREAM, activeKey, JSON.stringify(forwardObj));
  console.log(`Thread #forwardWorker: Reset task=${previousPid}, stream=${activeKey}, code=${codeObj.code}, at=${codeObj?.update}, r0=${r0}`);
}

async function removeDisabledTask(activeKey, forwardObj, configObj) {
  // Ignore if not disabled.
  if (configObj.enabled) return;

  // Check whether exited.
  // If process changed, ignore.
  while (forwardObj.pid === process.pid) {
    const code = await redis.hget(keys.redis.SRS_FORWARD_CODE, activeKey);
    const codeObj = code && JSON.parse(code);
    if (codeObj?.close) break;

    // Kill the process if not exited.
    try {
      process.kill(forwardObj.task, 'SIGKILL');
    } catch (e) {
    }
    console.log(`Thread #forwardWorker: Kill task=${forwardObj.task}, stream=${activeKey}, code=${codeObj?.code}, at=${codeObj?.update}`);

    await new Promise(resolve => setTimeout(resolve, 800));
  }

  // Cleanup redis.
  const r0 = await redis.hdel(keys.redis.SRS_FORWARD_STREAM, activeKey);
  const r1 = await redis.hdel(keys.redis.SRS_FORWARD_FRAME, activeKey);
  const r2 = await redis.hdel(keys.redis.SRS_FORWARD_CODE, activeKey);
  const r3 = await redis.hdel(keys.redis.SRS_FORWARD_MAP, forwardObj.platform);
  console.log(`Thread #forwardWorker: Cleanup task=${forwardObj.task}, platform=${forwardObj.platform}, stream=${activeKey}, r0=${r0}, r1=${r1}, r2=${r2}, r3=${r3}`);

  // Skip any options.
  return true;
}

async function terminateTaskForNoStream(activeKey, forwardObj, configObj) {
  // Ignore if not enabled.
  if (!configObj.enabled) return;

  // Ignore if not started.
  if (!forwardObj.task) return;

  // Check whether exited.
  const frame = await redis.hget(keys.redis.SRS_FORWARD_FRAME, activeKey);
  if (!frame) return;

  const frameObj = JSON.parse(frame);
  if (!frameObj?.update) return;

  // If not expired, ignore.
  const expired = moment(frameObj.update).add(process.env.NODE_ENV === 'development' ? 9 : 30, 's');
  if (expired.isAfter(moment())) return;

  // Expired, terminate the task.
  configObj.enabled = false;
  console.log(`Thread #forwardWorker: Expire task=${forwardObj.task}, platform=${forwardObj.platform}, stream=${activeKey}, update=${forwardObj.update}, expired=${expired.format()}, now=${moment().format()}`);

  return await removeDisabledTask(activeKey, forwardObj, configObj);
}

