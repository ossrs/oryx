'use strict';

// For mgmt, it's ok to connect to localhost.
const config = {
  redis:{
    host: 'localhost',
    port: 6379,
    password: '',
  },
};

const utils = require('js-core/utils');
const pkg = require('./package.json');
const { spawn } = require('child_process');
const metadata = require('./metadata');
const util = require('util');
const exec = util.promisify(require('child_process').exec);
const axios = require('axios');
const market = require('./market');
const consts = require('./consts');
const fs = require('fs');
const errs = require('js-core/errs');
const jwt = require('jsonwebtoken');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});

exports.handle = (router) => {
  router.all('/terraform/v1/mgmt/status', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    const upgrading = await redis.hget(consts.SRS_UPGRADING, 'upgrading');
    console.log(`status ok, upgrading=${upgrading}, releases=${JSON.stringify(metadata.upgrade.releases)}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      version: `v${pkg.version}`,
      releases: {
        stable: metadata.upgrade.releases?.stable,
        latest: metadata.upgrade.releases?.latest,
      },
      upgrading,
    });
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
    const dockerArgs = `docker run --rm --name certbot-certonly \\
      -v "${process.cwd()}/containers/etc/letsencrypt:/etc/letsencrypt" \\
      -v "${process.cwd()}/containers/var/lib/letsencrypt:/var/lib/letsencrypt" \\
      -v "${process.cwd()}/containers/var/log/letsencrypt:/var/log/letsencrypt" \\
      -v "${process.cwd()}/containers/www:/www" \\
      ccr.ccs.tencentyun.com/ossrs/certbot \\
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
    await exec(`ln -sf ${keyFile} /etc/nginx/ssl/nginx.key`);
    await exec(`ln -sf ${crtFile} /etc/nginx/ssl/nginx.crt`);
    await exec(`systemctl reload nginx.service`);

    console.log(`let's encrypt ok, domain=${domain}, key=${keyFile}, crt=${crtFile}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0);
  });

  router.all('/terraform/v1/mgmt/upgrade', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    const releaseServer = process.env.NODE_ENV === 'development' ? `http://localhost:${consts.config.port}` : 'http://api.ossrs.net';
    const {data: releases} = await axios.get(`${releaseServer}/terraform/v1/releases`, {
      params: {
        version: `v${pkg.version}`,
        ts: new Date().getTime(),
      }
    });
    metadata.upgrade.releases = releases;

    const target = releases?.latest || 'lighthouse';
    const upgradingMessage = `upgrade to target=${target}, current=${pkg.version}, releases=${JSON.stringify(releases)}`;
    console.log(`Start ${upgradingMessage}`);

    const r0 = await redis.hget(consts.SRS_UPGRADING, 'upgrading');
    if (r0) {
      const r1 = await redis.hget(consts.SRS_UPGRADING, 'desc');
      throw utils.asError(errs.sys.upgrading, errs.status.sys, `already upgrading ${r0}, ${r1}`);
    }

    // Set the upgrading to avoid others.
    await redis.hset(consts.SRS_UPGRADING, 'upgrading', 1);
    await redis.hset(consts.SRS_UPGRADING, 'desc', `${upgradingMessage}`);

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

    console.log(`upgrade ok, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      version: pkg.version,
    });
  });

  router.all('/terraform/v1/mgmt/container', async (ctx) => {
    const {token, action, name} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    if (!action) throw utils.asError(errs.sys.empty, errs.status.args, `no param action`);
    if (!name) throw utils.asError(errs.sys.empty, errs.status.args, `no param name`);

    const container = metadata.market[name];
    if (!container) throw utils.asError(errs.sys.resource, errs.status.not, `no container --name=${name}`);

    const validActions = ['query'];
    if (!validActions.includes(action)) throw utils.asError(errs.sys.invalid, errs.status.args, `invalid action ${action}, should be ${validActions}`);

    console.log(`srs ok, action=${action}, name=${name}, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      name: container.name,
      container: {
        ID: container.container.ID,
        State: container.container.State,
        Status: container.container.Status,
      },
    });
  });

  router.all('/terraform/v1/mgmt/hooks', async (ctx) => {
    const {token, action} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    console.log(`srs ok, action=${action} decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      name: metadata.market.hooks.name,
      container: {
        ID: metadata.market.hooks.container.ID,
        State: metadata.market.hooks.container.State,
        Status: metadata.market.hooks.container.Status,
      },
    });
  });

  return router;
};

