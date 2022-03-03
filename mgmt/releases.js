'use strict';

const consts = require('./consts');
const pkg = require('./package.json');
const keys = require('js-core/keys');
const { v4: uuidv4 } = require('uuid');

async function queryLatestVersion(redis, axios) {
  // Request release api with params.
  const params = {
    version: `v${pkg.version}`,
    ts: new Date().getTime(),
  };

  // Generate and setup the node id.
  let nid = await redis.hget(keys.redis.SRS_TENCENT_LH, 'node');
  if (!nid) {
    nid = uuidv4();
    await redis.hset(keys.redis.SRS_TENCENT_LH, 'node', nid);
    console.log(`Node: Generate node nid=${nid}`);
  }
  params.nid = nid;

  // Report about COS and resource usage.
  const cos = await redis.hget(keys.redis.SRS_TENCENT_COS, 'bucket');
  if (cos) {
    params.cos = 1;

    const cosn = await redis.hlen(keys.redis.SRS_DVR_M3U8_METADATA);
    if (cosn) params.cosn = cosn;
  }

  // Report about VoD and resource usage.
  const vod = await redis.hget(keys.redis.SRS_TENCENT_VOD, 'storage');
  if (vod) {
    params.vod = 1;

    const vodn = await redis.hlen(keys.redis.SRS_VOD_M3U8_METADATA);
    if (vodn) params.vodn = vodn;
  }

  // Request the release service API.
  const releaseServer = process.env.LOCAL_RELEASE === 'true' ? `http://localhost:${consts.config.port}` : 'https://api.ossrs.net';
  const {data: releases} = await axios.get(`${releaseServer}/terraform/v1/releases`, {params});
  return releases;
}

exports.queryLatestVersion = queryLatestVersion;

