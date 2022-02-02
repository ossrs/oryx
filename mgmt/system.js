'use strict';

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

exports.handle = (router) => {
  router.all('/terraform/v1/mgmt/status', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    console.log(`status ok, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      version: `v${pkg.version}`,
      releases: {
        stable: metadata.releases.releases?.stable,
        latest: metadata.releases.releases?.latest,
      },
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
    metadata.releases.releases = releases;

    let target = releases?.latest || 'lighthouse';
    console.log(`Start upgrade to target=${target}, current=${pkg.version}, releases=${JSON.stringify(releases)}`);

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

  router.all('/terraform/v1/mgmt/srs', async (ctx) => {
    const {token, action} = ctx.request.body;
    const decoded = await utils.verifyToken(jwt, token);

    if (action === 'restart') {
      // We must rm the container to get a new ID.
      await exec(`docker rm -f ${metadata.market.srs.name}`);

      const previousContainerID = metadata.market.srs.container.ID;
      for (let i = 0; i < 20; i++) {
        // Wait util running and got another container ID.
        const [all, running] = await market.queryContainer(metadata.market.srs.name);
        // Please note that we don't update the metadata of SRS, client must request the updated status.
        if (all && all.ID && running && running.ID && running.ID !== previousContainerID) break;
        await new Promise(resolve => setTimeout(resolve, 500));
      }
    }

    console.log(`srs ok, action=${action} decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      name: metadata.market.srs.name,
      container: {
        ID: metadata.market.srs.container.ID,
        State: metadata.market.srs.container.State,
        Status: metadata.market.srs.container.Status,
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

