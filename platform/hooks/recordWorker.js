'use strict';

// For components in docker, connect by host.
const config = {
  redis:{
    host: process.env.NODE_ENV === 'development' ? 'localhost' : (process.env.REDIS_HOST || 'mgmt.srs.local'),
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const { isMainThread, parentPort } = require("worker_threads");
const util = require('util');
const execFile = util.promisify(require('child_process').execFile);
const { v4: uuidv4 } = require('uuid');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const keys = require('js-core/keys');
const moment = require('moment');
const fs = require('fs');
const m3u8Generator = require('./m3u8Generator');
const { spawn } = require('child_process');

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  // We must initialize the thread first.
  console.log(`Thread #recordWorker: initialize`);

  parentPort.on('message', (msg) => {
    handleMessage(msg);
  });

  while (true) {
    try {
      await doThreadMain();
    } catch (e) {
      console.error(`Thread #recordWorker: err`, e);
    } finally {
      await new Promise(resolve => setTimeout(resolve, 30 * 1000));
    }
  }
}

async function handleMessage(msg) {
  const {action, m3u8_url, file, duration, seqno} = msg;
  if (action !== 'on_record_file') return console.log(`Thread #recordWorker: ignore ${JSON.stringify(msg)}`);

  // Copy the ts file to temporary cache dir.
  const tsid = uuidv4();
  const tsfile = `record/${tsid}.ts`;
  // Always use execFile when params contains user inputs, see https://auth0.com/blog/preventing-command-injection-attacks-in-node-js-apps/
  // Note that should never use fs.copyFileSync(file, tsfile, fs.constants.COPYFILE_FICLONE_FORCE) which fails in macOS.
  await execFile('cp', ['-f', file, tsfile]);

  // Create or update active m3u8 object, for worker to scan.
  const m3u8 = await redis.hget(keys.redis.SRS_RECORD_M3U8_ACTIVE, m3u8_url);
  const m3u8Obj = m3u8 ? JSON.parse(m3u8) : {update: null, uuid: uuidv4()};
  m3u8Obj.update = moment().format();

  // Append ts files to local m3u8 object.
  const local = await redis.hget(keys.redis.SRS_RECORD_M3U8_LOCAL, m3u8_url);
  const localObj = local ? JSON.parse(local) : {
    nn: 0,
    update: null,
    done: null,
    uuid: m3u8Obj.uuid,
    uuids: [],
    files: [],
  };
  if (!local || localObj.uuid !== m3u8Obj.uuid) {
    localObj.done = null;
    localObj.uuid = m3u8Obj.uuid;
    localObj.m3u8_url = m3u8_url;
    localObj.uuids.push(m3u8Obj.uuid);
    console.log(`Thread #recordWorker: local start new m3u8=${m3u8_url}, uuid=${m3u8Obj.uuid}, uuids=${localObj.uuids.length}`);
  }
  localObj.files.push({m3u8_url, tsid, tsfile, ...msg});
  localObj.nn = localObj.files.length;
  localObj.update = moment().format();

  // We update the local m3u8 object, then update the active m3u8 stream list, to make sure the localObj is ready.
  const r1 = await redis.hset(keys.redis.SRS_RECORD_M3U8_LOCAL, m3u8_url, JSON.stringify(localObj));
  const r0 = await redis.hset(keys.redis.SRS_RECORD_M3U8_ACTIVE, m3u8_url, JSON.stringify(m3u8Obj));
  console.log(`Thread #recordWorker: local record task m3u8=${m3u8_url}, files=${localObj.files.length}, uuid=${m3u8Obj.uuid}, file=${file}, tsfile=${tsfile}, duration=${duration}, seqno=${seqno}, r0=${r0}, r1=${r1}`);
}

async function doThreadMain() {
  while (true) {
    // We query the local active keys from m3u8 status.
    const activeKeys = await redis.hkeys(keys.redis.SRS_RECORD_M3U8_ACTIVE);
    if (!activeKeys || !activeKeys.length) {
      await new Promise(resolve => setTimeout(resolve, 9 * 1000));
      continue;
    }
    console.log(`Thread #recordWorker: local active keys ${JSON.stringify(activeKeys)}`);

    for (const i in activeKeys) {
      // Get the local object by the active key.
      const localKey = activeKeys[i];
      const local = await redis.hget(keys.redis.SRS_RECORD_M3U8_LOCAL, localKey);
      await handleLocalObject(localKey, local ? JSON.parse(local) : null);
    }

    await new Promise(resolve => setTimeout(resolve, 3000));
  }
}

async function handleLocalObject(localKey, localObj) {
  if (!localObj || !localObj.files.length) {
    await finishLocalObject(localKey, localObj);
    return;
  }

  const localFiles = [...localObj.files];
  for (const i in localFiles) {
    const localFile = localFiles[i];
    await handleLocalFile(localKey, localObj, localFile);
  }

  const metadata = await redis.hget(keys.redis.SRS_RECORD_M3U8_METADATA, localObj.uuid);
  const metadataObj = metadata && JSON.parse(metadata);

  console.log(`Thread #recordWorker: Finished files=${localFiles.length}, left=${localObj.files.length}, metadata=${metadataObj?.files?.length}`);
}

async function handleLocalFile(localKey, localObj, localFile) {
  // Ignore file if not exists.
  if (!fs.existsSync(localFile.tsfile)) {
    await updateLocalObject(localKey, localObj, localFile, null);
    console.warn(`Thread #recordWorker: Ignore m3u8=${localKey}, ts=${localFile.url} for TsNotFount`);
    return;
  }

  const stats = fs.statSync(localFile.tsfile);
  const key = `record/${localObj.uuid}/${localFile.tsid}.ts`;
  fs.mkdirSync(`record/${localObj.uuid}`, {recursive: true});
  fs.renameSync(localFile.tsfile, key)

  // Update the metadata for m3u8.
  await updateMetadataObject(localKey, localObj, localFile, key, stats);

  // Update the left local local files.
  await updateLocalObject(localKey, localObj, localFile, key);

  // Local record, already move the tsfile to key.
  console.log(`Thread #recordWorker: Finish local for m3u8=${localKey}, ts=${localFile.url}, as=${key}`);
}

async function updateMetadataObject(localKey, localObj, localFile, key, stats) {
  const metadata = await redis.hget(keys.redis.SRS_RECORD_M3U8_METADATA, localObj.uuid);
  const metadataObj = metadata ? JSON.parse(metadata) : {
    nn: 0,
    update: moment().format(),
    uuid: localObj.uuid,
    m3u8_url: localObj.m3u8_url,
    vhost: localFile.params.vhost,
    app: localFile.params.app,
    stream: localFile.params.stream,
    // The Record is progressing, use local m3u8 address to preview or download.
    progress: true,
    done: null,
    m3u8: null,
    // The ts files.
    files: [],
  };

  // Reduce the uploaded files by uuid.
  metadataObj.files = metadataObj.files.filter(e => e.tsid !== localFile.tsid);
  metadataObj.files.push({
    key,
    tsid: localFile.tsid,
    url: localFile.url,
    seqno: localFile.seqno,
    duration: localFile.duration,
    size: stats.size,
  });
  metadataObj.nn = metadataObj.files.length;
  metadataObj.update = moment().format();

  await redis.hset(keys.redis.SRS_RECORD_M3U8_METADATA, localObj.uuid, JSON.stringify(metadataObj));
  console.log(`Thread #recordWorker: Update metadata for m3u8=${localKey}, uuid=${localObj.uuid}, files=${metadataObj.nn}`);
}

// Note that key is optional, logging only.
async function updateLocalObject(localKey, localObj, localFile, key) {
  // @remark Note that the local.files might changed by other asyncs, so we must reload it before save it.
  const localRef = await redis.hget(keys.redis.SRS_RECORD_M3U8_LOCAL, localKey);
  const localRefObj = JSON.parse(localRef);

  // Warning if files changed.
  if (localObj.files.length !== localRefObj.files.length) {
    console.warn(`Thread #recordWorker: LocalTsFiles changed, m3u8=${localKey}, before=${localObj.files.length}, ref=${localRefObj.files.length}`);
  }

  // Filter the left files.
  const leftFiles = localRefObj.files.filter(e => e.tsid !== localFile.tsid);
  console.log(`Thread #recordWorker: Update local for m3u8=${localKey}, ts=${localFile.url}, as=${key}, before=${localRefObj.files.length}, left=${leftFiles.length}`);

  // Update the local realtime reference object.
  localObj.files = localRefObj.files = leftFiles;
  localObj.nn = localRefObj.nn = leftFiles.length;
  localObj.update = localRefObj.update = moment().format();

  // Write to redis.
  await redis.hset(keys.redis.SRS_RECORD_M3U8_LOCAL, localKey, JSON.stringify(localRefObj));
}

async function finishLocalObject(localKey, localObj) {
  if (!localObj || !localObj.update) {
    await redis.hdel(keys.redis.SRS_RECORD_M3U8_LOCAL, localKey);
    return;
  }

  // If stream expired, finish the Record.
  const expired = moment(localObj.update).add(process.env.NODE_ENV === 'development' ? 30 : 300, 's');
  if (expired.isAfter(moment())) return;

  // Try to finish the m3u8 first.
  const duration = await finishM3u8(localKey, localObj);

  // Keep the local status, to allow query all uuids of uploaded.
  localObj.done = moment().format();
  await redis.hset(keys.redis.SRS_RECORD_M3U8_LOCAL, localKey, JSON.stringify(localObj));
  console.log(`Thread #recordWorker: Record expired, key=${localKey}, update=${moment(localObj.update).format()}, expired=${expired.format()}, now=${moment().format()}`);

  // Remove the status, to create new m3u8 next Record.
  await redis.hdel(keys.redis.SRS_RECORD_M3U8_ACTIVE, localKey);

  console.log(`Thread #recordWorker: Record done, key=${localKey}, duration=${duration}`);
}

async function finishM3u8(localKey, localObj) {
  // Update the metadata, and keep the uploaded.
  const metadata = await redis.hget(keys.redis.SRS_RECORD_M3U8_METADATA, localObj.uuid);
  if (!metadata) return;

  const metadataObj = metadata && JSON.parse(metadata);
  const [contentType, m3u8Body, duration] = m3u8Generator.buildVodM3u8(
    metadataObj, false, null, false,
  );

  const hls = `record/${localObj.uuid}/index.m3u8`;
  const mp4 = `record/${localObj.uuid}/index.mp4`;
  fs.writeFileSync(hls, m3u8Body);
  console.log(`Thread #recordWorker: Record m3u8=${hls}, mp4=${mp4}, contentType=${contentType}`);

  // Start a child process to transmux HLS to MP4.
  const child = spawn('ffmpeg', ['-i', hls, '-c', 'copy', mp4]);
  console.log(`Thread #forwardWorker: Start task=${child.pid}, input=${hls}, output=${mp4}`);

  await new Promise((resolve, reject) => {
    let nnLogs = 0;
    child.stdout.on('data', (chunk) => {
      console.log(chunk.toString());
    });
    child.stderr.on('data', (chunk) => {
      if ((nnLogs++ % 3) !== 0) return;
      const log = chunk.toString().trim().replace(/= +/g, '=');
      if (log.indexOf('frame=') < 0) return;
      console.log(`Thread #forwardWorker: Active task=${child.pid}, hls=${hls}, ${log}`);
    });
    child.on('close', (code) => {
      console.log(`Thread #forwardWorker: Close task=${child.pid}, hls=${hls}, code=${code}`);
      if (code !== 0) {
        reject(code);
      } else {
        resolve();
      }
    });
  });

  metadataObj.progress = false;
  metadataObj.done = moment().format();
  await redis.hset(keys.redis.SRS_RECORD_M3U8_METADATA, localObj.uuid, JSON.stringify(metadataObj));
  return duration;
}

