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

  // Report about FFmpeg forwarding.
  const forward = await redis.hlen(keys.redis.SRS_FORWARD_STREAM);
  if (forward) params.forward = forward;

  // Report about active streams.
  const streams = await redis.hlen(keys.redis.SRS_STREAM_ACTIVE);
  if (streams) params.streams = streams;

  // Report about SRT stream.
  const srt = await redis.hlen(keys.redis.SRS_STREAM_SRT_ACTIVE);
  if (srt) params.srt = srt;

  // Report about WebRTC stream.
  const rtc = await redis.hlen(keys.redis.SRS_STREAM_RTC_ACTIVE);
  if (rtc) params.rtc = rtc;

  // Report about beian feature.
  const beian = await redis.hlen(keys.redis.SRS_BEIAN);
  if (beian) params.beian = beian;

  // Report about HTTPS feature.
  const ssl = await redis.get(keys.redis.SRS_HTTPS);
  if (ssl) params.https = ssl;

  // Request the release service API.
  const releaseServer = process.env.LOCAL_RELEASE === 'true' ? `http://localhost:${consts.config.port}` : 'https://api.ossrs.net';
  console.log(`Query ${releaseServer} with ${JSON.stringify(params)}`);

  const {data: releases} = await axios.get(`${releaseServer}/terraform/v1/releases`, {params});
  return releases;
}

exports.queryLatestVersion = queryLatestVersion;

