'use strict';

const utils = require('./utils');
const pkg = require('./package.json');
const { spawn } = require('child_process');
const srs = require('./srs');
const releases = require('./releases');

exports.handle = (router) => {
  router.all('/terraform/v1/mgmt/status', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(token);

    console.log(`status ok, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      version: pkg.version,
      releases: releases.metadata,
    });
  });

  router.all('/terraform/v1/mgmt/upgrade', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(token);

    let target = releases.metadata.versions?.stable || 'lighthouse';
    console.log(`Start upgrade to target=${target}, current=${pkg.version}, releases=${JSON.stringify(releases.metadata.versions)}`);

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
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(token);

    console.log(`srs ok, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      ...srs.metadata,
      container: {
        ID: srs.metadata.container.ID,
        State: srs.metadata.container.State,
        Status: srs.metadata.container.Status,
      },
    });
  });

  return router;
};

