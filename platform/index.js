'use strict';

const dotenv = require('dotenv');
const fs = require('fs');
const path = require('path');
const utils = require('js-core/utils');
utils.reloadEnv(dotenv, fs, path);
console.log(`load envs MGMT_PASSWORD=${'*'.repeat(process.env.MGMT_PASSWORD?.length)}`);

const Koa = require('koa');
const Router = require('koa-router');
const Cors = require('koa2-cors');
const BodyParser = require('koa-bodyparser');
const token = require('./token');
const pkg = require('./package.json');
const thread = require('./thread');
const system = require('./system');
const platform = require('./platform');
const staticCache = require('koa-static-cache');
const serve = require('koa-static');
const mount  = require('koa-mount');

const app = new Koa();

// Always enable CORS for statics or apis.
app.use(Cors());

///////////////////////////////////////////////////////////////////////////////////////////
//   Static File Server sections.
///////////////////////////////////////////////////////////////////////////////////////////
// For react, static files server, by /mgmt/
utils.srsProxy(
  staticCache,
  app,
  path.join(__dirname, 'ui/build'),
  '/mgmt/',
  ['/mgmt/index.html'],
  {'/mgmt/': '/mgmt/index.html'},
);

// For react-router pages, by /mgmt/routers-*
app.use(async (ctx, next) => {
  // Compatible with old react routes.
  // TODO: FIXME: Remove it in next large release.
  const isPreviousReactRoutes = [
    '/mgmt/login',
    '/mgmt/dashboard',
    '/mgmt/scenario',
    '/mgmt/config',
    '/mgmt/system',
    '/mgmt/logout',
  ].includes(ctx.request.path);

  // Directly serve the react routes by index.html
  // See https://stackoverflow.com/a/52464577/17679565
  if (isPreviousReactRoutes || ctx.request.path.indexOf('/mgmt/routers-') === 0) {
    ctx.type = 'text/html';
    ctx.set('Cache-Control', 'public, max-age=0');
    ctx.body = fs.readFileSync('./ui/build/index.html');
    return;
  }

  await next();
});

// For automatic HTTPS by letsencrypt, for certbot to verify the domain.
// Note that should never create the directory .well-known/acme-challenge/ because it's auto created by certbot.
// See https://eff-certbot.readthedocs.io/en/stable/using.html#webroot
// See https://github.com/ossrs/srs/issues/2864#issuecomment-1027944527
app.use(
  mount(
    '/.well-known/acme-challenge/',
    serve('./containers/www/.well-known/acme-challenge/'),
  ),
);

///////////////////////////////////////////////////////////////////////////////////////////
//   API sections.
///////////////////////////////////////////////////////////////////////////////////////////
app.use(BodyParser());

// For Error handler.
app.use(async (ctx, next) => {
  try {
    await next();
  } catch (e) {
    ctx.status = e.status || 500;
    ctx.body = utils.asResponse(e.code || 1, {
      message: e.message || e.err?.message || 'unknown error',
    });
    console.error(e);
  }
});

// For backend APIs, with specified path, by /terraform/v1/mgmt/
const router = new Router();

token.handle(system.handle(router));

router.all('/terraform/v1/mgmt/versions', async (ctx) => {
  ctx.body = utils.asResponse(0, {version: pkg.version});
});

app.use(router.routes());

///////////////////////////////////////////////////////////////////////////////////////////
const run = async () => {
  await platform.init();
  console.log(`Run with cwd=${process.cwd()}`);

  thread.run();

  app.listen(2024, () => {
    console.log(`Server start on http://localhost:2024`);
  });
};
run();

