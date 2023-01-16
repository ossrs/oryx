'use strict';

// For components in docker, connect by host.
const config = {
  redis: {
    host: process.env.NODE_ENV === 'development' ? 'localhost' : (process.env.REDIS_HOST || 'mgmt.srs.local'),
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const errs = require('js-core/errs');
const utils = require('js-core/utils');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const jwt = require('jsonwebtoken');
const keys = require('js-core/keys');
const {koaBody} = require('koa-body');
const fs = require("fs");
const {v4: uuidv4} = require('uuid');
const path = require('path');
const util = require("util");
const execFile = util.promisify(require('child_process').execFile);

const dirUploadPath = path.join('.', 'upload');
const dirVLivePath = path.join('.', 'vlive');

exports.handle = (router) => {
  router.all('/terraform/v1/ffmpeg/forward/secret', async (ctx) => {
    const {token, action, platform, server, secret, enabled, custom, label} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    const allowedActions = ['update'];
    const allowedPlatforms = ['wx', 'bilibili', 'kuaishou'];
    if (action) {
      if (!allowedActions.includes(action)) {
        throw utils.asError(errs.sys.invalid, errs.status.args, `invalid action ${action}`);
      }

      if (!platform) throw utils.asError(errs.sys.empty, errs.status.args, 'no platform');
      if (!allowedPlatforms.includes(platform)) {
        throw utils.asError(errs.sys.invalid, errs.status.args, `invalid platform ${platform}`);
      }

      if (!server) throw utils.asError(errs.sys.empty, errs.status.args, 'no server');
      if (!server && !secret) throw utils.asError(errs.sys.empty, errs.status.args, 'no secret');
      if (enabled === undefined) throw utils.asError(errs.sys.empty, errs.status.args, 'no enabled');
      if (custom === undefined) throw utils.asError(errs.sys.empty, errs.status.args, 'no custom');
    }

    let res = null;
    if (action === 'update') {
      const conf = await redis.hget(keys.redis.SRS_FORWARD_CONFIG, platform);
      const confObj = conf ? JSON.parse(conf) : {};
      const r0 = await redis.hset(keys.redis.SRS_FORWARD_CONFIG, platform, JSON.stringify({
        ...confObj, platform, server, secret, enabled, custom, label,
      }));

      // Restart the forwarding if exists.
      const stream = await redis.hget(keys.redis.SRS_FORWARD_MAP, platform);
      if (stream) {
        const activeKey = `${platform}@${stream}`;
        const forward = await redis.hget(keys.redis.SRS_FORWARD_STREAM, activeKey);
        const forwardObj = forward && JSON.parse(forward);
        if (forwardObj?.task) {
          try {
            process.kill(forwardObj.task, 'SIGKILL');
          } catch (e) {
          }
          console.log(`FFmpeg: Forward kill pid=${forwardObj.task}, stream=${activeKey}`);
        }
      }
      console.log(`FFmpeg: Forward update secret ok, action=${action}, platform=${platform}, r0=${r0}`);
    } else {
      const configs = await redis.hgetall(keys.redis.SRS_FORWARD_CONFIG);
      for (const k in configs) {
        configs[k] = JSON.parse(configs[k]);
      }
      res = configs;
    }

    console.log(`FFmpeg: Forward secret ok, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, res);
  });

  router.all('/terraform/v1/ffmpeg/forward/streams', async (ctx) => {
    const {token} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    const configs = await redis.hgetall(keys.redis.SRS_FORWARD_CONFIG);
    const maps = await redis.hgetall(keys.redis.SRS_FORWARD_MAP);
    const frames = await redis.hgetall(keys.redis.SRS_FORWARD_FRAME);

    const res = Object.values(configs).map(e => {
      const conf = JSON.parse(e);

      const stream = maps[conf.platform];
      const frame = stream && frames[`${conf.platform}@${stream}`];
      return {
        platform: conf.platform,
        enabled: conf.enabled,
        custom: !!conf.custom,
        label: conf.label,
        stream,
        frame: frame ? JSON.parse(frame) : null,
      };
    });

    console.log(`FFmpeg: Query forward streams ok, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, res);
  });

  router.all('/terraform/v1/ffmpeg/vlive/secret', async (ctx) => {
    const {token, action, platform, server, secret, enabled, custom, label, files} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    const allowedActions = ['update'];
    const allowedPlatforms = ['wx', 'bilibili', 'kuaishou'];
    if (action) {
      if (!allowedActions.includes(action)) {
        throw utils.asError(errs.sys.invalid, errs.status.args, `invalid action ${action}`);
      }

      if (!platform) throw utils.asError(errs.sys.empty, errs.status.args, 'no platform');
      if (!allowedPlatforms.includes(platform)) {
        throw utils.asError(errs.sys.invalid, errs.status.args, `invalid platform ${platform}`);
      }

      if (!server) throw utils.asError(errs.sys.empty, errs.status.args, 'no server');
      if (!server && !secret) throw utils.asError(errs.sys.empty, errs.status.args, 'no secret');
      if (enabled === undefined) throw utils.asError(errs.sys.empty, errs.status.args, 'no enabled');
      if (custom === undefined) throw utils.asError(errs.sys.empty, errs.status.args, 'no custom');
      if (!files?.length) throw utils.asError(errs.sys.empty, errs.status.args, 'no files');
    }

    let res = null;
    if (action === 'update') {
      const conf = await redis.hget(keys.redis.SRS_VLIVE_CONFIG, platform);
      const confObj = conf ? JSON.parse(conf) : {files:[]};
      const r0 = await redis.hset(keys.redis.SRS_VLIVE_CONFIG, platform, JSON.stringify({
        ...confObj, platform, server, secret, enabled, custom, label, files,
      }));

      // Restart the virtual live stream if exists.
      const source = await redis.hget(keys.redis.SRS_VLIVE_MAP, platform);
      if (source) {
        const activeKey = `${platform}@${source}`;
        const vLive = await redis.hget(keys.redis.SRS_VLIVE_STREAM, activeKey);
        const vLiveObj = vLive && JSON.parse(vLive);
        if (vLiveObj?.task) {
          try {
            process.kill(vLiveObj.task, 'SIGKILL');
          } catch (e) {
          }
          console.log(`FFmpeg: vLive kill pid=${vLiveObj.task}, stream=${activeKey}`);
        }
      }
      console.log(`FFmpeg: vLive update secret ok, action=${action}, platform=${platform}, r0=${r0}`);
    } else {
      const configs = await redis.hgetall(keys.redis.SRS_VLIVE_CONFIG);
      for (const k in configs) {
        configs[k] = JSON.parse(configs[k]);
      }
      res = configs;
    }

    console.log(`FFmpeg: vLive secret ok, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, res);
  });

  router.all('/terraform/v1/ffmpeg/vlive/streams', async (ctx) => {
    const {token} = ctx.request.body;

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    const configs = await redis.hgetall(keys.redis.SRS_VLIVE_CONFIG);
    const maps = await redis.hgetall(keys.redis.SRS_VLIVE_MAP);
    const frames = await redis.hgetall(keys.redis.SRS_VLIVE_FRAME);

    const res = Object.values(configs).map(e => {
      const conf = JSON.parse(e);

      const source = maps[conf.platform];
      const frame = source && frames[`${conf.platform}@${source}`];
      return {
        platform: conf.platform,
        enabled: conf.enabled,
        custom: !!conf.custom,
        label: conf.label,
        files: conf.files,
        source,
        frame: frame ? JSON.parse(frame) : null,
      };
    });

    console.log(`FFmpeg: Query vlive streams ok, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, res);
  });

  // Serve ts to play.
  router.all('/terraform/v1/ffmpeg/vlive/upload/:name', koaBody(
    {
      multipart: true,
      formidable: {
        maxFileSize: 99 * 1024 * 1024 * 1024, // 99GB
        uploadDir: dirUploadPath,
      },
    }
  ), async (ctx) => {
    const {name} = ctx.params;
    if (!name) throw utils.asError(errs.sys.empty, errs.status.args, `no param name`);

    const {filepath, newFilename, originalFilename, mimetype, size} = ctx.request.files[name];
    const uuid = uuidv4();
    const target = path.join(dirUploadPath, `${uuid}${path.extname(originalFilename)}`);
    console.log(`FFmpeg: Got vlive filepath=${filepath}, newFilename=${newFilename}, originalFilename=${originalFilename}, mimetype=${mimetype}, size=${size}, target=${target}`);

    fs.renameSync(filepath, target);
    console.log(`FFmpeg: Rename ${filepath} to ${target}`);

    ctx.body = utils.asResponse(0, {uuid, target});
  });

  const refreshSourceOfVLive = async (ctx, token, platform, sourceFiles) => {
    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);

    // Copy the source files, because we need to modify the target.
    const files = sourceFiles.map(f => utils.copy(f));

    const allowedPlatforms = ['wx', 'bilibili', 'kuaishou'];
    if (true) {
      if (!platform) throw utils.asError(errs.sys.empty, errs.status.args, 'no platform');
      if (!allowedPlatforms.includes(platform)) {
        throw utils.asError(errs.sys.invalid, errs.status.args, `invalid platform ${platform}`);
      }
    }

    // Move file from dirUploadPath to dirVLivePath.
    for (const file of files) {
      // Probe file information.
      const {stdout} = await execFile('ffprobe', [
        '-show_error', '-show_private_data', '-v', 'quiet', '-find_stream_info', '-print_format', 'json',
        '-show_format', '-show_streams', file.target,
      ]);
      const fi = JSON.parse(stdout);
      file.format = {
        duration: fi?.format?.duration, bit_rate: fi?.format?.bit_rate,
        nb_streams: fi?.format?.nb_streams, probe_score: fi?.format?.probe_score,
        has_video: !!fi?.streams?.filter(e => e?.codec_type === 'video')?.length,
        has_audio: !!fi?.streams?.filter(e => e?.codec_type === 'audio')?.length,
      };
      if (file.format.has_video) {
        const v = fi?.streams?.filter(e => e?.codec_type === 'video')[0];
        file.video = {
          codec_name: v?.codec_name, profile: v?.profile, width: v?.width, height: v?.height,
          pix_fmt: v?.pix_fmt, level: v?.level, duration: v?.duration, bit_rate: v?.bit_rate,
          nb_frames: v?.nb_frames,
        };
      }
      if (file.format.has_audio) {
        const a = fi?.streams?.filter(e => e?.codec_type === 'audio')[0];
        file.audio = {
          codec_name: a?.codec_name, profile: a?.profile, sample_fmt: a?.sample_fmt,
          sample_rate: a?.sample_rate, channels: a?.channels, channel_layout: a?.channel_layout,
          duration: a?.duration, bit_rate: a?.bit_rate, nb_frames: a?.nb_frames,
        };
      }

      // Move file from upload to vlive.
      const newTarget = path.join(dirVLivePath, `${file.uuid}${path.extname(file.target)}`);
      fs.renameSync(file.target, newTarget);
      console.log(`FFmpeg: vLive move file from ${JSON.stringify(file)} to ${newTarget}, platform=${platform}`);
      file.target = newTarget;
    }

    // Update redis object.
    const conf = await redis.hget(keys.redis.SRS_VLIVE_CONFIG, platform);
    const confObj = conf ? JSON.parse(conf) : {files:[]};
    const r0 = await redis.hset(keys.redis.SRS_VLIVE_CONFIG, platform, JSON.stringify({
      ...confObj, platform, files,
    }));

    // Remove old files.
    for (const file of (confObj.files || [])) {
      if (!file.target) continue;
      if (fs.existsSync(file.target)) fs.rmSync(file.target);
      console.log(`FFmpeg: vLive remove old file ${JSON.stringify(file)} for platform=${platform}`);
    }

    console.log(`FFmpeg: vLive update secret ok, decoded=${JSON.stringify(decoded)}, platform=${platform}, r0=${r0}`);
    ctx.body = utils.asResponse(0, {platform, files});
  };

  // Refresh the source of vlive.
  router.all('/terraform/v1/ffmpeg/vlive/source', async (ctx) => {
    const {token, platform, files} = ctx.request.body;

    if (!files?.length) throw utils.asError(errs.sys.empty, errs.status.args, 'no files');
    for (const file of files) {
      if (!file?.target) throw utils.asError(errs.sys.empty, errs.status.args, `empty file ${JSON.stringify(file)}`);
      if (!fs.existsSync(file?.target)) throw utils.asError(errs.sys.empty, errs.status.args, `no file ${JSON.stringify(file)}`);
      const requiredPrefix = path.join(path.normalize(dirUploadPath), path.sep);
      if (path.normalize(file?.target).indexOf(requiredPrefix) !== 0) throw utils.asError(errs.sys.empty, errs.status.args, `invalid file ${JSON.stringify(file)}`);
    }

    try {
      await refreshSourceOfVLive(ctx, token, platform, files);
    } finally {
      // Always cleanup the files in upload.
      if (!files?.length) return;
      for (const file of files) {
        if (!fs.existsSync(file?.target)) continue;
        fs.rmSync(file?.target);
        console.warn(`FFmpeg: vLive cleanup ${JSON.stringify(file)}`);
      }
    }
  });

  return router;
};

