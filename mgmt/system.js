'use strict';

// For mgmt, it's ok to connect to localhost.
const config = {
  redis:{
    host: 'localhost',
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD || '',
  },
};

const fs = require('fs');
const path = require('path');
const dotenv = require('dotenv');
const utils = require('js-core/utils');
const errs = require('js-core/errs');
const jwt = require('jsonwebtoken');
const util = require('util');
const execFile = util.promisify(require('child_process').execFile);
const market = require('./market');
const ioredis = require('ioredis');
const redis = require('js-core/redis').create({config: config.redis, redis: ioredis});
const pkg = require('./package.json');
const metadata = require('./metadata');
const { spawn } = require('child_process');
const axios = require('axios');
const {queryLatestVersion} = require('./releases');
const platform = require('./platform');
const keys = require('js-core/keys');

const handlers = {
  // Query the containers information.
  queryContainer: async ({ctx, action, args}) => {
    const name = args[0];
    if (!name) throw utils.asError(errs.sys.empty, errs.status.args, `no name`);

    const [all, running] = await market.queryContainer(name);
    ctx.body = utils.asResponse(0, {all, running});
  },

  // Fecth the containers information.
  fetchContainers: async ({ctx, action, args}) => {
    const name = args[0];

    const containers = [];
    for (const k in metadata.market) {
      if (name && k !== name) return;

      const container = metadata.market[k];
      if (!container) throw utils.asError(errs.sys.resource, errs.status.not, `no container --name=${k}`);

      // Query container enabled status from redis.
      const disabled = await redis.hget(keys.redis.SRS_CONTAINER_DISABLED, container.name);

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

    ctx.body = utils.asResponse(0, {containers});
  },

  // Remove the specified container.
  rmContainer: async ({ctx, action, args}) => {
    const name = args[0];
    if (!name) throw utils.asError(errs.sys.empty, errs.status.args, `no name`);

    await utils.removeContainerQuiet(execFile, name);
    ctx.body = utils.asResponse(0);
  },

  // Reload the env from .env
  reloadEnv: async ({ctx, action, args}) => {
    utils.reloadEnv(dotenv, fs, path);
    ctx.body = utils.asResponse(0);
  },

  // Query the cached version.
  queryVersion: async ({ctx, action, args}) => {
    ctx.body = utils.asResponse(0, {
      version: `v${pkg.version}`,
      stable: metadata.upgrade.releases?.stable,
      latest: metadata.upgrade.releases?.latest,
    });
  },

  // Refresh the version from api.
  refreshVersion: async ({ctx, action, args}) => {
    const releases = await queryLatestVersion(redis, axios);
    metadata.upgrade.releases = releases;

    ctx.body = utils.asResponse(0, {
      version: `v${pkg.version}`,
      stable: metadata.upgrade.releases?.stable,
      latest: metadata.upgrade.releases?.latest,
    });
  },

  // Start upgrade.
  execUpgrade: async ({ctx, action, args}) => {
    const target = args[0];
    if (!target) throw utils.asError(errs.sys.empty, errs.status.args, `no target`);

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

    ctx.body = utils.asResponse(0);
  },

  // Write SSL cert and key files.
  updateSslFile: async ({ctx, action, args}) => {
    const [key, crt] = args;

    if (!key) throw utils.asError(errs.sys.empty, errs.status.args, 'no key');
    if (!crt) throw utils.asError(errs.sys.empty, errs.status.args, 'no crt');

    if (!fs.existsSync('/etc/nginx/ssl/nginx.key')) throw utils.asError(errs.sys.ssl, errs.status.sys, 'no key file');
    if (!fs.existsSync('/etc/nginx/ssl/nginx.crt')) throw utils.asError(errs.sys.ssl, errs.status.sys, 'no crt file');

    // Remove the ssl file, because it might link to other file.
    await execFile('rm', ['-f', '/etc/nginx/ssl/nginx.key', '/etc/nginx/ssl/nginx.crt']);

    // Write the ssl key and cert, and reload nginx when ready.
    fs.writeFileSync('/etc/nginx/ssl/nginx.key', key);
    fs.writeFileSync('/etc/nginx/ssl/nginx.crt', crt);
    await execFile('systemctl', ['reload', 'nginx.service']);

    ctx.body = utils.asResponse(0);
  },

  // Update SSL by let's encrypt.
  updateLetsEncrypt: async ({ctx, action, args}) => {
    const [domain] = args;

    if (!domain) throw utils.asError(errs.sys.empty, errs.status.args, 'no domain');

    if (!fs.existsSync('/etc/nginx/ssl/nginx.key')) throw utils.asError(errs.sys.ssl, errs.status.sys, 'no key file');
    if (!fs.existsSync('/etc/nginx/ssl/nginx.crt')) throw utils.asError(errs.sys.ssl, errs.status.sys, 'no crt file');

    // We run always with "-n Run non-interactively"
    // Note that it's started by nodejs, so never use '-it' or failed for 'the input device is not a TTY'.
    //
    // The www root to verify while requesting the SSL file:
    //    /.well-known/acme-challenge/
    // which mount as:
    //    ./containers/www/.well-known/acme-challenge/
    const registry = await platform.registry();
    const dockerArgs = ['run', '--rm', '--name', 'certbot-certonly',
      '-v', `${process.cwd()}/containers/etc/letsencrypt:/etc/letsencrypt`,
      '-v', `${process.cwd()}/containers/var/lib/letsencrypt:/var/lib/letsencrypt`,
      '-v', `${process.cwd()}/containers/var/log/letsencrypt:/var/log/letsencrypt`,
      '-v', `${process.cwd()}/containers/www:/www`,
      `${registry}/ossrs/certbot`,
      'certonly', '--webroot', '-w', '/www',
      '-d', `${domain}`, '--register-unsafely-without-email', '--agree-tos', '--preferred-challenges', 'http',
      '-n',
    ];
    await execFile('docker', dockerArgs);
    console.log(`certbot request ssl ok docker ${dockerArgs.join(' ')}`);

    const keyFile = `${process.cwd()}/containers/etc/letsencrypt/live/${domain}/privkey.pem`;
    if (!fs.existsSync(keyFile)) throw utils.asError(errs.sys.ssl, errs.status.sys, `issue key file ${keyFile}`);

    const crtFile = `${process.cwd()}/containers/etc/letsencrypt/live/${domain}/cert.pem`;
    if (!fs.existsSync(crtFile)) throw utils.asError(errs.sys.ssl, errs.status.sys, `issue crt file ${crtFile}`);

    // Remove the ssl file, because it might link to other file.
    await execFile('rm', ['-f', '/etc/nginx/ssl/nginx.key', '/etc/nginx/ssl/nginx.crt']);

    // Always use execFile when params contains user inputs, see https://auth0.com/blog/preventing-command-injection-attacks-in-node-js-apps/
    await execFile('ln', ['-sf', keyFile, '/etc/nginx/ssl/nginx.key']);
    await execFile('ln', ['-sf', crtFile, '/etc/nginx/ssl/nginx.crt']);

    // Restart the nginx service to reload the SSL files.
    await execFile('systemctl', ['reload', 'nginx.service']);

    ctx.body = utils.asResponse(0);
  },

  // Renew the lets encrypt SSL files.
  renewLetsEncrypt: async ({ctx, action, args}) => {
    // Whether force to renew.
    const [force] = args;

    // Remove the ssl file, because it might link to other file.
    const signalFile = `${process.cwd()}/containers/var/log/letsencrypt/CERTBOT_HOOK_RELOAD_NGINX`;
    await execFile('rm', ['-f', signalFile]);

    // We run always with "-n Run non-interactively"
    // Note that it's started by nodejs, so never use '-it' or failed for 'the input device is not a TTY'.
    const registry = await platform.registry();
    const dockerArgs = ['run', '--rm', '--name', 'certbot-renew',
      '-v', `${process.cwd()}/containers/etc/letsencrypt:/etc/letsencrypt`,
      '-v', `${process.cwd()}/containers/var/lib/letsencrypt:/var/lib/letsencrypt`,
      '-v', `${process.cwd()}/containers/var/log/letsencrypt:/var/log/letsencrypt`,
      `${registry}/ossrs/certbot`,
      'renew', '--post-hook', 'touch /var/log/letsencrypt/CERTBOT_HOOK_RELOAD_NGINX',
      // See https://github.com/ossrs/srs/issues/2864#issuecomment-1027944527
      // Use --force-renewal and --no-random-sleep-on-renew to always renew a cert,
      // see https://community.letsencrypt.org/t/disabling-random-sleep-of-certbot/83201
      ...(force ? ['--force-renewal', '--no-random-sleep-on-renew'] : []),
      '-n',
    ];
    const {stdout} = await execFile('docker', dockerArgs);
    console.log(`certbot renew ssl ok, docker ${dockerArgs.join(' ')}`);

    // Restart the nginx service to reload the SSL files.
    const renewOk = fs.existsSync(signalFile);
    if (renewOk) await execFile('systemctl', ['reload', 'nginx.service']);
    console.log(`certbot renew updated=${renewOk}, args=${args}, message is ${stdout}`);

    ctx.body = utils.asResponse(0, {stdout, renewOk});
  },

  // Update access for ssh keys.
  accessSsh: async ({ctx, action, args}) => {
    const [enabled] = args;
    const {stdout} = await execFile('bash', ['auto/update_access', enabled]);
    console.log(`Access: SSH enabled=${enabled}, message is ${stdout}`);

    ctx.body = utils.asResponse(0);
  },

  // Execute the upgrade prepare script when platform startup.
  executeUpgradePrepare: async ({ctx, action, args}) => {
    const {stdout} = await execFile('bash', ['auto/upgrade_prepare']);
    console.log(`Upgrade: Prepare message is ${stdout}`);

    ctx.body = utils.asResponse(0);
  },

  // RPC template.
  xxx: async ({ctx, action, args}) => {
  },
};

exports.handle = (router) => {
  // In the container, we can't neither manage other containers, nor execute command, so we must request this api to
  // execute the command on host machine.
  router.all('/terraform/v1/host/exec', async (ctx) => {
    const {action, token, args} = ctx.request.body;

    if (!token) throw utils.asError(errs.sys.empty, errs.status.auth, `no token`);
    if (!action) throw utils.asError(errs.sys.empty, errs.status.args, `no action`);

    if (!Object.keys(handlers).includes(action)) {
      throw utils.asError(errs.sys.invalid, errs.status.args, `invalid action=${action}`);
    }

    const apiSecret = await utils.apiSecret(redis);
    const decoded = await utils.verifyToken(jwt, token, apiSecret);
    console.log(`exec verify action=${action}, args=${JSON.stringify(args)}, token=${token.length}B, decoded=${JSON.stringify(decoded)}`);

    if (handlers[action]) {
      await handlers[action]({ctx, action, args});
    } else {
      throw utils.asError(errs.sys.invalid, errs.status.args, `invalid action ${action}`);
    }
  });

  return router;
};

