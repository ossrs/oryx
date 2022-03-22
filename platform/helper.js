'use strict';

// For components in docker, connect by host.
const config = {
  redis:{
    host: process.env.NODE_ENV === 'development' ? 'localhost' : 'mgmt.srs.local',
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const utils = require('js-core/utils');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const jwt = require('jsonwebtoken');
const axios = require('axios');
const moment = require('moment');
const keys = require('js-core/keys');
const metadata = require('js-core/metadata');

exports.inUpgradeWindow = (uwStart, uwDuration, now) => {
  if (uwStart === undefined || uwStart === null || !uwDuration || !now) {
    return true;
  }

  const [start, duration] = [parseInt(uwStart), parseInt(uwDuration)];
  const end = duration >= 24 ? start + 24 : (start + duration) % 24;

  let inWindow;
  if (start < end) {
    inWindow = (start <= now.hours() && now.hours() <= end);
  } else {
    inWindow = !(end < now.hours() && now.hours() < start);
  }
  console.log(`Upgrade window=${inWindow}, start=${start}, duration=${duration}, end=${end}, now=${now.format()}, hours=${now.hours()}`);

  return inWindow;
}

async function execApi(action, args) {
  if (args !== null && args !== undefined && !Array.isArray(args)) {
    throw new Error(`args is not array, ${args}`);
  }

  const apiSecret = await utils.apiSecret(redis);
  const token = utils.createToken(moment, jwt, apiSecret);
  const server = process.env.NODE_ENV === 'development' ? 'localhost' : 'mgmt.srs.local';

  try {
    const data = await axios.post(
      `http://${server}:2022/terraform/v1/host/exec`,
      {
        ...token, action, args: args || [],
      },
    );
    return data?.data?.data;
  } catch (e) {
    console.error(`exec server=${server}, action=${action}, args=${JSON.stringify(args)}, err`, e);
  }
}
exports.execApi = execApi;

async function queryLatestVersion() {
  // Request release api with params.
  const params = {};

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
  const platform = await redis.hget(keys.redis.SRS_TENCENT_LH, 'platform');
  if (platform) params.plat = platform;

  return await execApi('refreshVersion', [params]);
}
exports.queryLatestVersion = queryLatestVersion;

