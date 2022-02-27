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
const util = require('util');
const execFile = util.promisify(require('child_process').execFile);
const { v4: uuidv4 } = require('uuid');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const keys = require('js-core/keys');
const moment = require('moment');
const COS = require('cos-nodejs-sdk-v5');
const fs = require('fs');
const m3u8Generator = require('./m3u8Generator');

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  // We must initialize the thread first.
  console.log(`Thread #dvrWorker: initialize`);

  parentPort.on('message', (msg) => {
    handleMessage(msg);
  });

  while (true) {
    try {
      await doThreadMain();
    } catch (e) {
      console.error(`Thread #dvrWorker: err`, e);
    } finally {
      await new Promise(resolve => setTimeout(resolve, 30 * 1000));
    }
  }
}

async function handleMessage(msg) {
  const {action, m3u8_url, file, duration, seqno} = msg;
  if (action !== 'on_dvr_file') return console.log(`Thread #dvrWorker: ignore ${JSON.stringify(msg)}`);

  // Copy the ts file to temporary cache dir.
  const tsid = uuidv4();
  const tsfile = `dvr/${tsid}.ts`;
  // Always use execFile when params contains user inputs, see https://auth0.com/blog/preventing-command-injection-attacks-in-node-js-apps/
  await execFile('cp', ['-f', file, tsfile]);

  // Create or update active m3u8 object, for worker to scan.
  const m3u8 = await redis.hget(keys.redis.SRS_DVR_M3U8_ACTIVE, m3u8_url);
  const m3u8Obj = m3u8 ? JSON.parse(m3u8) : {update: null, uuid: uuidv4()};
  m3u8Obj.update = moment().format();
  const r0 = await redis.hset(keys.redis.SRS_DVR_M3U8_ACTIVE, m3u8_url, JSON.stringify(m3u8Obj));

  // Append ts files to local m3u8 object.
  const local = await redis.hget(keys.redis.SRS_DVR_M3U8_LOCAL, m3u8_url);
  const localObj = local ? JSON.parse(local) : {nn: 0, update: null, done: null, uuid: m3u8Obj.uuid, uuids: [], files: []};
  if (!local || localObj.uuid !== m3u8Obj.uuid) {
    localObj.done = null;
    localObj.uuid = m3u8Obj.uuid;
    localObj.uuids.push(m3u8Obj.uuid);
    console.log(`Thread #dvrWorker: local start new m3u8=${m3u8_url}, uuid=${m3u8Obj.uuid}, uuids=${localObj.uuids.length}`);
  }
  localObj.files.push({m3u8_url, tsid, tsfile, ...msg});
  localObj.nn = localObj.files.length;
  localObj.update = moment().format();
  const r1 = await redis.hset(keys.redis.SRS_DVR_M3U8_LOCAL, m3u8_url, JSON.stringify(localObj));
  console.log(`Thread #dvrWorker: local dvr task m3u8=${m3u8_url}, files=${localObj.files.length}, uuid=${m3u8Obj.uuid}, file=${file}, tsfile=${tsfile}, duration=${duration}, seqno=${seqno}, r0=${r0}, r1=${r1}`);
}

async function doThreadMain() {
  const secretId = await redis.hget(keys.redis.SRS_TENCENT_CAM, 'secretId');
  const secretKey = await redis.hget(keys.redis.SRS_TENCENT_CAM, 'secretKey');
  const cos = new COS({SecretId: secretId, SecretKey: secretKey});

  while (true) {
    // We query the local active keys from m3u8 status.
    const activeKeys = await redis.hkeys(keys.redis.SRS_DVR_M3U8_ACTIVE);
    if (!activeKeys || !activeKeys.length) {
      await new Promise(resolve => setTimeout(resolve, 9 * 1000));
      continue;
    }
    console.log(`Thread #dvrWorker: local active keys ${JSON.stringify(activeKeys)}`);

    const bucket = await redis.hget(keys.redis.SRS_TENCENT_COS, 'bucket');
    const region = await redis.hget(keys.redis.SRS_TENCENT_LH, 'region');
    if (!bucket || !region) {
      console.log(`Thread #dvrWorker: ignore for bucket=${bucket}, region=${region}`);
      await new Promise(resolve => setTimeout(resolve, 3 * 1000));
      continue;
    }

    for (const i in activeKeys) {
      // Get the local object by the active key.
      const localKey = activeKeys[i];
      const local = await redis.hget(keys.redis.SRS_DVR_M3U8_LOCAL, localKey);
      await handleLocalObject(cos, bucket, region, localKey, local ? JSON.parse(local) : null);
    }

    await new Promise(resolve => setTimeout(resolve, 3000));
  }
}

async function handleLocalObject(cos, bucket, region, localKey, localObj) {
  if (!localObj || !localObj.files.length) {
    await finishLocalObject(cos, bucket, region, localKey, localObj);
    return;
  }

  const localFiles = [...localObj.files];
  for (const i in localFiles) {
    const localFile = localFiles[i];
    await handleLocalFile(cos, bucket, region, localKey, localObj, localFile);
  }

  const uploaded = await redis.hget(keys.redis.SRS_DVR_M3U8_UPLOADED, localObj.uuid);
  const uploadedObj = uploaded && JSON.parse(uploaded);

  const metadata = await redis.hget(keys.redis.SRS_DVR_M3U8_METADATA, localObj.uuid);
  const metadataObj = metadata && JSON.parse(metadata);

  console.log(`Thread #dvrWorker: Finished files=${localFiles.length}, left=${localObj.files.length}, uploaded=${uploadedObj?.files?.length}, metadata=${metadataObj?.files?.length}`);
}

async function handleLocalFile(cos, bucket, region, localKey, localObj, localFile) {
  // Ignore file if not exists.
  if (!fs.existsSync(localFile.tsfile)) {
    await updateLocalObject(localKey, localObj, localFile, null);
    console.warn(`Thread #dvrWorker: Ignore m3u8=${localKey}, ts=${localFile.url} for TsNotFount`);
    return;
  }

  // Upload the ts file to COS.
  const stats = fs.statSync(localFile.tsfile);
  const key = `${localObj.uuid}/${localFile.tsid}.ts`;
  await uploadToCos(cos, bucket, region, localObj, localFile, key, stats);

  // Update the uploaded ts files.
  await updateUploadedObject(localObj, localFile, key, stats);

  // Update the metadata for m3u8.
  await updateMetadataObject(bucket, region, localKey, localObj, localFile, key, stats);

  // Update the left local local files.
  await updateLocalObject(localKey, localObj, localFile, key);

  // Remove the tsfile.
  fs.unlinkSync(localFile.tsfile);
  console.log(`Thread #dvrWorker: Remove local for m3u8=${localKey}, ts=${localFile.url}, as=${key}`);
}

async function uploadToCos(cos, bucket, region, localObj, localFile, key, stats) {
  await new Promise((resolve, reject) => {
    // See https://cloud.tencent.com/document/product/436/64980
    cos.putObject({
      Bucket: bucket,
      Region: region,
      Key: key,
      StorageClass: 'STANDARD',
      Body: fs.createReadStream(localFile.tsfile),
      ContentLength: stats.size,
      ContentType: 'video/MP2T',
      onProgress: function(progressData) {
        //console.log(`Thread #dvrWorker: progress ${JSON.stringify(progressData)}`);
      },
    }, function(err, data) {
      if (err) return reject(err);
      resolve(data);
    });
  });
}

async function updateUploadedObject(localObj, localFile, key, stats) {
  const uploaded = await redis.hget(keys.redis.SRS_DVR_M3U8_UPLOADED, localObj.uuid);
  const uploadedObj = uploaded ? JSON.parse(uploaded) : {
    nn: 0,
    update: moment().format(),
    uuid: localObj.uuid,
    files: [],
  };

  // Reduce the uploaded files by uuid.
  uploadedObj.files = uploadedObj.files.filter(e => e.tsid !== localFile.tsid);
  // Append the uploaded ts file.
  uploadedObj.files.push({
    key,
    tsid: localFile.tsid,
    size: stats.size,
    ...localFile,
  });
  uploadedObj.nn = uploadedObj.files.length;

  await redis.hset(keys.redis.SRS_DVR_M3U8_UPLOADED, localObj.uuid, JSON.stringify(uploadedObj));
}

async function updateMetadataObject(bucket, region, localKey, localObj, localFile, key, stats) {
  const metadata = await redis.hget(keys.redis.SRS_DVR_M3U8_METADATA, localObj.uuid);
  const metadataObj = metadata ? JSON.parse(metadata) : {
    nn: 0,
    update: moment().format(),
    bucket,
    region,
    uuid: localObj.uuid,
    vhost: localFile.params.vhost,
    app: localFile.params.app,
    stream: localFile.params.stream,
    // The DVR is progressing, use local m3u8 address to preview or download.
    progress: true,
    done: null,
    m3u8: null,
    // The ts files in COS.
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

  await redis.hset(keys.redis.SRS_DVR_M3U8_METADATA, localObj.uuid, JSON.stringify(metadataObj));
  console.log(`Thread #dvrWorker: Update metadata for m3u8=${localKey}, uuid=${localObj.uuid}, files=${metadataObj.nn}`);
}

// Note that key is optional, logging only.
async function updateLocalObject(localKey, localObj, localFile, key) {
  // @remark Note that the local.files might changed by other asyncs, so we must reload it before save it.
  const localRef = await redis.hget(keys.redis.SRS_DVR_M3U8_LOCAL, localKey);
  const localRefObj = JSON.parse(localRef);

  // Warning if files changed.
  if (localObj.files.length !== localRefObj.files.length) {
    console.warn(`Thread #dvrWorker: LocalTsFiles changed, m3u8=${localKey}, before=${localObj.files.length}, ref=${localRefObj.files.length}`);
  }

  // Filter the left files.
  const leftFiles = localRefObj.files.filter(e => e.tsid !== localFile.tsid);
  console.log(`Thread #dvrWorker: Update local for m3u8=${localKey}, ts=${localFile.url}, as=${key}, before=${localRefObj.files.length}, left=${leftFiles.length}`);

  // Update the local realtime reference object.
  localObj.files = localRefObj.files = leftFiles;
  localObj.nn = localRefObj.nn = leftFiles.length;
  localObj.update = localRefObj.update = moment().format();

  // Write to redis.
  await redis.hset(keys.redis.SRS_DVR_M3U8_LOCAL, localKey, JSON.stringify(localRefObj));
}

async function finishLocalObject(cos, bucket, region, localKey, localObj) {
  if (!localObj || !localObj.update) {
    await redis.hdel(keys.redis.SRS_DVR_M3U8_LOCAL, localKey);
    return;
  }

  // If stream expired, finish the DVR.
  const expired = moment(localObj.update).add(process.env.NODE_ENV === 'development' ? 30 : 300, 's');
  if (expired.isAfter(moment())) return;

  // Try to finish the m3u8 first.
  const duration = await finishM3u8(cos, bucket, region, localKey, localObj);

  // Keep the local status, to allow query all uuids of uploaded.
  localObj.done = moment().format();
  await redis.hset(keys.redis.SRS_DVR_M3U8_LOCAL, localKey, JSON.stringify(localObj));
  console.log(`Thread #dvrWorker: DVR expired, key=${localKey}, update=${moment(localObj.update).format()}, expired=${expired.format()}, now=${moment().format()}`);

  // Remove the status, to create new m3u8 next DVR.
  await redis.hdel(keys.redis.SRS_DVR_M3U8_ACTIVE, localKey);

  console.log(`Thread #dvrWorker: DVR done, key=${localKey}, duration=${duration}`);
}

async function finishM3u8(cos, bucket, region, localKey, localObj) {
  // Update the metadata, and keep the uploaded.
  const metadata = await redis.hget(keys.redis.SRS_DVR_M3U8_METADATA, localObj.uuid);
  if (!metadata) return;

  const metadataObj = metadata && JSON.parse(metadata);
  const [contentType, m3u8Body, duration] = m3u8Generator.buildVodM3u8(metadataObj, false);

  // Upload the ts file to COS.
  const key = `${localObj.uuid}/index.m3u8`;
  await new Promise((resolve, reject) => {
    // See https://cloud.tencent.com/document/product/436/64980
    cos.putObject({
      Bucket: bucket,
      Region: region,
      Key: key,
      StorageClass: 'STANDARD',
      Body: m3u8Body,
      ContentLength: m3u8Body.length,
      ContentType: contentType,
      onProgress: function(progressData) {
        //console.log(`Thread #dvrWorker: progress ${JSON.stringify(progressData)}`);
      },
    }, function(err, data) {
      if (err) return reject(err);
      resolve(data);
    });
  });

  metadataObj.progress = false;
  metadataObj.done = moment().format();
  await redis.hset(keys.redis.SRS_DVR_M3U8_METADATA, localObj.uuid, JSON.stringify(metadataObj));
  return duration;
}

