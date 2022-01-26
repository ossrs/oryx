'use strict';

const utils = require('./utils');
const pkg = require('./package.json');
const { spawn } = require('child_process');
const metadata = require('./metadata');

exports.handle = (router) => {
  router.all('/terraform/v1/mgmt/status', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(token);

    console.log(`status ok, decoded=${JSON.stringify(decoded)}, token=${token.length}B`);
    ctx.body = utils.asResponse(0, {
      version: `v${pkg.version}`,
      releases: {
        stable: metadata.releases.releases.stable,
        latest: metadata.releases.releases.latest,
      },
    });
  });

  router.all('/terraform/v1/mgmt/upgrade', async (ctx) => {
    const {token} = ctx.request.body;
    const decoded = await utils.verifyToken(token);

    let target = metadata.releases?.stable || 'lighthouse';
    console.log(`Start upgrade to target=${target}, current=${pkg.version}, releases=${JSON.stringify(metadata.releases)}`);

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
      name: metadata.srs.name,
      major: metadata.srs.major,
      container: {
        ID: metadata.srs.container.ID,
        State: metadata.srs.container.State,
        Status: metadata.srs.container.Status,
      },
    });
  });

  return router;
};

