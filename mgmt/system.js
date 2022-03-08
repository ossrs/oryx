'use strict';

// For mgmt, it's ok to connect to localhost.
const config = {
  redis:{
    host: 'localhost',
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const utils = require('js-core/utils');
const pkg = require('./package.json');
const { spawn } = require('child_process');
const metadata = require('./metadata');
const util = require('util');
const exec = util.promisify(require('child_process').exec);
const execFile = util.promisify(require('child_process').execFile);
const axios = require('axios');
const consts = require('./consts');
const fs = require('fs');
const errs = require('js-core/errs');
const jwt = require('jsonwebtoken');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const moment = require('moment');
const platform = require('./platform');
const {queryLatestVersion} = require('./releases');
const keys = require('js-core/keys');
const helper = require('./helper');

exports.handle = (router) => {
  router.all('/terraform/v1/mgmt/status', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    const upgrading = await redis.hget(consts.SRS_UPGRADING, 'upgrading');
    const r0 = await redis.hget(consts.SRS_UPGRADE_STRATEGY, 'strategy');
    const strategy = r0 || 'auto';
    console.log(`status ok, upgrading=${upgrading}, strategy=${strategy}, releases=${JSON.stringify(metadata.upgrade.releases)}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      version: `v${pkg.version}`,
      releases: {
        stable: metadata.upgrade.releases?.stable,
        latest: metadata.upgrade.releases?.latest,
      },
      upgrading: upgrading === "1",
      strategy,
    });
  });

  router.all('/terraform/v1/mgmt/strategy', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    const releases = await queryLatestVersion(redis, axios);
    metadata.upgrade.releases = releases;

    const upgrading = await redis.hget(consts.SRS_UPGRADING, 'upgrading');
    const r0 = await redis.hget(consts.SRS_UPGRADE_STRATEGY, 'strategy');
    const strategy = r0 || 'auto';
    const newStrategy = strategy === 'auto' ? 'manual' : 'auto';
    const r1 = await redis.hset(consts.SRS_UPGRADE_STRATEGY, 'strategy', newStrategy);
    const r2 = await redis.hset(consts.SRS_UPGRADE_STRATEGY, 'desc', `${moment().format()} changed, upgrading=${upgrading}, r0=${r0}/${strategy}, r1=${r1}/${newStrategy}`);
    console.log(`status ok, upgrading=${upgrading}, r0=${r0}/${strategy}, r1=${r1}/${newStrategy}, r2=${r2}, releases=${JSON.stringify(metadata.upgrade.releases)}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  router.all('/terraform/v1/mgmt/window/query', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    const start = await redis.hget(keys.redis.SRS_UPGRADE_WINDOW, 'start');
    const duration = await redis.hget(keys.redis.SRS_UPGRADE_WINDOW, 'duration');

    console.log(`query upgrade window ok, start=${start}, duration=${duration}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      start: parseInt(start),
      duration: parseInt(duration),
    });
  });

  router.all('/terraform/v1/mgmt/window/update', async (ctx) => {
    const {token, start, duration} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    if (isNaN(parseInt(start))) throw utils.asError(errs.sys.empty, errs.status.args, `no param start`);
    if (!duration) throw utils.asError(errs.sys.empty, errs.status.args, `no param end`);
    if (duration <= 3) return utils.asError(errs.sys.invalid, errs.status.args, `window should greater than 3 hours`);
    if (duration > 24) return utils.asError(errs.sys.invalid, errs.status.args, `window should smaller than 24 hours`);

    const r0 = await redis.hset(keys.redis.SRS_UPGRADE_WINDOW, 'start', start);
    const r1 = await redis.hset(keys.redis.SRS_UPGRADE_WINDOW, 'duration', duration);
    const r2 = await redis.hset(keys.redis.SRS_UPGRADE_WINDOW, 'update', moment().format());

    console.log(`update upgrade window ok, start=${start}, duration=${duration}, r0=${r0}, r1=${r1}, r2=${r2}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  router.all('/terraform/v1/mgmt/upgrade', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    const releases = await queryLatestVersion(redis, axios);
    metadata.upgrade.releases = releases;

    const uwStart = await redis.hget(keys.redis.SRS_UPGRADE_WINDOW, 'start');
    const uwDuration = await redis.hget(keys.redis.SRS_UPGRADE_WINDOW, 'duration');
    const inUpgradeWindow = helper.inUpgradeWindow(uwStart, uwDuration, moment());

    const target = releases?.latest || 'lighthouse';
    const upgradingMessage = `upgrade to target=${target}, current=${pkg.version}, releases=${JSON.stringify(releases)}, window=${inUpgradeWindow}`;
    console.log(`Start ${upgradingMessage}`);

    const r0 = await redis.hget(consts.SRS_UPGRADING, 'upgrading');
    if (r0 === "1") {
      const r1 = await redis.hget(consts.SRS_UPGRADING, 'desc');
      throw utils.asError(errs.sys.upgrading, errs.status.sys, `already upgrading ${r0}, ${r1}`);
    }

    // Set the upgrading to avoid others.
    await redis.hset(consts.SRS_UPGRADING, 'upgrading', 1);
    await redis.hset(consts.SRS_UPGRADING, 'desc', `${upgradingMessage}`);

    try {
      await new Promise((resolve, reject) => {
        const child = spawn('bash', ['upgrade', target]);
        child.stdout.on('data', (chunk) => {
          console.log(chunk.toString());
        });
        child.stderr.on('data', (chunk) => {
          console.log(chunk.toString());
        });
        child.on('close', (code) => {
          console.log(`upgrading exited with code ${code}`);
          if (code !== 0) return reject(code);
          resolve();
        });
      });
    } finally {
      // Upgrade done.
      await redis.hset(consts.SRS_UPGRADING, 'upgrading', 0);
    }

    console.log(`upgrade ok, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      version: pkg.version,
    });
  });

  router.all('/terraform/v1/mgmt/containers', async (ctx) => {
    const {token, action, name, enabled} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    if (!action) throw utils.asError(errs.sys.empty, errs.status.args, `no param action`);

    const validActions = ['query', 'enabled'];
    if (!validActions.includes(action)) throw utils.asError(errs.sys.invalid, errs.status.args, `invalid action ${action}, should be ${validActions}`);

    if (action === 'enabled') {
      if (!name) throw utils.asError(errs.sys.empty, errs.status.args, `no param name`);
      if (enabled !== true && enabled !== false) throw utils.asError(errs.sys.empty, errs.status.args, `no param enabled`);
      const r0 = await redis.hset(consts.SRS_CONTAINER_DISABLED, name, !enabled);
      if (!enabled && name) await exec(`docker rm -f ${name}`);
      console.log(`srs ok, action=${action}, name=${name}, enabled=${enabled}, r0=${r0}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
      return ctx.body = utils.asResponse(0);
    }

    // Query containers
    const containers = [];
    for (const k in metadata.market) {
      if (name && k !== name) return;

      const container = metadata.market[k];
      if (!container) throw utils.asError(errs.sys.resource, errs.status.not, `no container --name=${k}`);

      // Query container enabled status from redis.
      const disabled = await redis.hget(consts.SRS_CONTAINER_DISABLED, container.name);

      containers.push({
        name: container.name,
        enabled: disabled !== 'true',
        container: {
          ID: container.container.ID,
          State: container.container.State,
          Status: container.container.Status,
        },
      });
    }

    console.log(`srs ok, action=${action}, name=${name}, containers=${containers.length}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, containers);
  });

  router.all('/terraform/v1/mgmt/pubkey', async (ctx) => {
    const {token, enabled} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    await exec(`bash auto/update_access ${enabled ? 'enable' : 'disable'}`);

    console.log(`pubkey ok, enable=${enabled}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  router.all('/terraform/v1/mgmt/ssl', async (ctx) => {
    const {token, key, crt} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    if (!key) throw utils.asError(errs.sys.empty, errs.status.args, 'no key');
    if (!crt) throw utils.asError(errs.sys.empty, errs.status.args, 'no crt');
    if (!fs.existsSync('/etc/nginx/ssl/nginx.key')) throw utils.asError(errs.sys.ssl, errs.status.sys, 'no key file');
    if (!fs.existsSync('/etc/nginx/ssl/nginx.crt')) throw utils.asError(errs.sys.ssl, errs.status.sys, 'no crt file');

    // Remove the ssl file, because it might link to other file.
    await exec(`rm -f /etc/nginx/ssl/nginx.key /etc/nginx/ssl/nginx.crt`);

    // Write the ssl key and cert, and reload nginx when ready.
    fs.writeFileSync('/etc/nginx/ssl/nginx.key', key);
    fs.writeFileSync('/etc/nginx/ssl/nginx.crt', crt);
    await exec(`systemctl reload nginx.service`);

    // Setup the HTTPS information.
    await redis.set(keys.redis.SRS_HTTPS, 'ssl');

    console.log(`ssl ok, key=${key.length}B, crt=${crt.length}B, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  router.all('/terraform/v1/mgmt/letsencrypt', async (ctx) => {
    const {token, domain} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    if (!domain) throw utils.asError(errs.sys.empty, errs.status.args, 'no domain');
    if (!fs.existsSync('/etc/nginx/ssl/nginx.key')) throw utils.asError(errs.sys.ssl, errs.status.sys, 'no key file');
    if (!fs.existsSync('/etc/nginx/ssl/nginx.crt')) throw utils.asError(errs.sys.ssl, errs.status.sys, 'no crt file');

    // We run always with "-n Run non-interactively"
    // Note that it's started by nodejs, so never use '-it' or failed for 'the input device is not a TTY'.
    const registry = await platform.registry();
    const dockerArgs = `docker run --rm --name certbot-certonly \\
      -v "${process.cwd()}/containers/etc/letsencrypt:/etc/letsencrypt" \\
      -v "${process.cwd()}/containers/var/lib/letsencrypt:/var/lib/letsencrypt" \\
      -v "${process.cwd()}/containers/var/log/letsencrypt:/var/log/letsencrypt" \\
      -v "${process.cwd()}/containers/www:/www" \\
      ${registry}/ossrs/certbot \\
      certonly --webroot -w /www \\
      -d ${domain} --register-unsafely-without-email --agree-tos --preferred-challenges http \\
      -n`;
    await exec(dockerArgs);
    console.log(`certbot request ssl ok ${dockerArgs}`);

    const keyFile = `${process.cwd()}/containers/etc/letsencrypt/live/${domain}/privkey.pem`;
    if (!fs.existsSync(keyFile)) throw utils.asError(errs.sys.ssl, errs.status.sys, `issue key file ${keyFile}`);

    const crtFile = `${process.cwd()}/containers/etc/letsencrypt/live/${domain}/cert.pem`;
    if (!fs.existsSync(crtFile)) throw utils.asError(errs.sys.ssl, errs.status.sys, `issue crt file ${crtFile}`);

    // Remove the ssl file, because it might link to other file.
    await exec(`rm -f /etc/nginx/ssl/nginx.key /etc/nginx/ssl/nginx.crt`);

    // Always use execFile when params contains user inputs, see https://auth0.com/blog/preventing-command-injection-attacks-in-node-js-apps/
    await execFile('ln', ['-sf', keyFile, '/etc/nginx/ssl/nginx.key']);
    await execFile('ln', ['-sf', crtFile, '/etc/nginx/ssl/nginx.crt']);

    // Restart the nginx service to reload the SSL files.
    await exec(`systemctl reload nginx.service`);

    // Setup the HTTPS information.
    await redis.set(keys.redis.SRS_HTTPS, 'lets');

    console.log(`let's encrypt ok, domain=${domain}, key=${keyFile}, crt=${crtFile}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  router.all('/terraform/v1/mgmt/bilibili', async (ctx) => {
    const {token, bvid} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    if (!bvid) throw utils.asError(errs.sys.empty, errs.status.args, 'no bvid');

    const res = await new Promise((resolve, reject) => {
      axios.get(`https://api.bilibili.com/x/web-interface/view?bvid=${bvid}`).then(res => {
        resolve(res.data.data);
      }).catch(err => reject);
    });

    console.log(`bilibili query ok, bvid=${bvid}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, res);
  });

  router.all('/terraform/v1/mgmt/beian/query', async (ctx) => {
    const icp = await redis.hget(keys.redis.SRS_BEIAN, 'icp');

    console.log(`beian: query ok, miit=${JSON.stringify(icp)}`);
    ctx.body = utils.asResponse(0, {icp});
  });

  router.all('/terraform/v1/mgmt/beian/update', async (ctx) => {
    const {token, beian, text} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    if (!beian) throw utils.asError(errs.sys.empty, errs.status.args, 'no beian');
    if (!text) throw utils.asError(errs.sys.empty, errs.status.args, 'no text');

    const r0 = await redis.hset(keys.redis.SRS_BEIAN, beian, text);
    console.log(`beian: update ok, beian=${beian}, text=${text}, r0=${JSON.stringify(r0)}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, r0);
  });

  return router;
};

