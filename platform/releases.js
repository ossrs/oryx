'use strict';

// For components in docker, connect by host.
const config = {
  redis:{
    host: process.env.NODE_ENV === 'development' ? 'localhost' : (process.env.REDIS_HOST || 'mgmt.srs.local'),
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const keys = require('js-core/keys');
const metadata = require('js-core/metadata');
const platform = require('./platform');

async function queryLatestVersion() {
  // Request release api with params.
  const params = {};

  // Generate and setup the node id.
  const nid = await redis.hget(keys.redis.SRS_TENCENT_LH, 'node');
  if (nid) params.nid = nid;

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
  const streams = await redis.hget(keys.redis.SRS_STAT_COUNTER, 'publish');
  if (streams) {
    await redis.hset(keys.redis.SRS_STAT_COUNTER, 'publish', 0);
    params.streams = parseInt(streams);
  }

  // Report about active players.
  const players = await redis.hget(keys.redis.SRS_STAT_COUNTER, 'play');
  if (streams) {
    await redis.hset(keys.redis.SRS_STAT_COUNTER, 'play', 0);
    params.players = parseInt(players);
  }

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

  // Report about upgrade window feature.
  const uwin = await redis.hget(keys.redis.SRS_UPGRADE_WINDOW, 'update');
  if (uwin) params.uwin = 1;

  // Report whether start as develop environment.
  const dev = (process.env.NODE_ENV === 'development');
  if (dev) params.dev = 1;

  // Report whether enable SRS development version.
  const srsDev = await redis.hget(keys.redis.SRS_CONTAINER_DISABLED, metadata.market.srsDev.name)
  if (srsDev === 'false') params.srsd = 1;

  // Report about the platform.
  const plat = await redis.hget(keys.redis.SRS_TENCENT_LH, 'platform');
  if (plat) params.plat = plat;

  const cloud = await redis.hget(keys.redis.SRS_TENCENT_LH, 'cloud');
  if (cloud) params.cloud = cloud;

  const region = await redis.hget(keys.redis.SRS_TENCENT_LH, 'region');
  if (region) params.region = region;

  return await execApi('refreshVersion', [params]);
}
exports.queryLatestVersion = queryLatestVersion;

