'use strict';

const dotenv = require('dotenv');
const fs = require('fs');
const path = require('path');
['.', '..', '../..'].map(envDir => {
  if (fs.existsSync(path.join(envDir, '.env'))) {
    dotenv.config({path: path.join(envDir, '.env')});
  }
});
console.log(`load envs MGMT_PASSWORD=${'*'.repeat(process.env.MGMT_PASSWORD?.length)}`);

const Koa = require('koa');
const proxy = require('koa-proxies');
const Router = require('koa-router');
const Cors = require('koa2-cors');
const BodyParser = require('koa-bodyparser');
const serve = require('koa-static');
const mount  = require('koa-mount');
const auth = require('./auth');
const utils = require('js-core/utils');
const system = require('./system');
const threads = require('./threads');
const consts = require('./consts');
const pkg = require('./package.json');
const staticCache = require('koa-static-cache');

// Start all workers threads first.
threads.run();

// Create koa webserver.
const app = new Koa();

// Always enable CORS for statics or apis.
app.use(Cors());

///////////////////////////////////////////////////////////////////////////////////////////
//   Proxy Server sections.
///////////////////////////////////////////////////////////////////////////////////////////
function withLogs(options) {
  return {
    ...options,
    logs: (ctx, target) => {
      const r = new URL(ctx.req.url, target);
      console.log('%s - %s %s proxy to -> %s', new Date().toISOString(), ctx.req.method, ctx.req.oldPath, r);
    },
  };
}

// For version management.
app.use(proxy('/terraform/v1/releases', withLogs({target: 'http://127.0.0.1:2023/'})));

// For prometheus.
// TODO: FIXME: Do authentication for api.
app.use(proxy('/prometheus', withLogs({target: 'http://127.0.0.1:9090/'})));

// For registered modules, by /terraform/v1/hooks/
app.use(proxy('/terraform/v1/hooks/', withLogs({target: 'http://127.0.0.1:2021/'})));
// Compatible with old APIs, to work with running SRS wihtout restart them.
// TODO: FIXME: Remove it when all SRS container restarted.
app.use(proxy('/terraform/v1/mgmt/srs/hooks', withLogs({
  target: 'http://127.0.0.1:2021/',
  rewrite: path => path.replace('/terraform/v1/mgmt/srs/hooks', '/terraform/v1/hooks/srs/verify'),
})));

// Proxy to SRS HTTP streaming, console and player, by /api/, /rtc/, /live/, /console/, /players/
// See https://github.com/vagusX/koa-proxies
// TODO: FIXME: Do authentication for api.
app.use(proxy('/api/', withLogs({target: 'http://127.0.0.1:1985/'})));
app.use(proxy('/rtc/', withLogs({target: 'http://127.0.0.1:1985/'})));
app.use(proxy('/*/*.(flv|m3u8|ts|aac|mp3)', withLogs({target: 'http://127.0.0.1:8080/'})));
app.use(proxy('/console/', withLogs({target: 'http://127.0.0.1:8080/'})));
app.use(proxy('/players/', withLogs({target: 'http://127.0.0.1:8080/'})));

///////////////////////////////////////////////////////////////////////////////////////////
//   Static File Server sections.
///////////////////////////////////////////////////////////////////////////////////////////
// For source files like srs.tar.gz, by /terraform/v1/sources/
app.use(mount('/terraform/v1/sources/', serve('./sources')));

// For automatic HTTPS by letsencrypt, for certbot to verify the domain.
// Note that should never create the directory .well-known/acme-challenge/ because it's auto created by certbot.
// See https://eff-certbot.readthedocs.io/en/stable/using.html#webroot
// See https://github.com/ossrs/srs/issues/2864#issuecomment-1027944527
app.use(mount('/.well-known/acme-challenge/', serve('./containers/www/.well-known/acme-challenge/')));

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

// For react, static files server, by /mgmt/
if (true) {
  const reactFiles = {};

  app.use(staticCache(path.join(__dirname, 'ui/build'), {
    // Cache for a year for it never changes.
    maxAge: 365 * 24 * 3600,
    // It's important to set to dynamic, because the js might changed.
    dynamic: true,
    // If not set, NOT FOUND.
    alias: {
      '/mgmt/': '/mgmt/index.html',
    },
    // The baseUrl for react.
    prefix: '/mgmt/',
  }, reactFiles));

  // Disable the index.html cache, because need to load the correct latest js files,
  // see https://github.com/koajs/static-cache#editing-the-files-object
  reactFiles['/mgmt/index.html'].maxAge = 0;
}

// For /favicon.ico
// For homepage from root, use mgmt.
app.use(async (ctx, next) => {
  if (ctx.request.path === '/favicon.ico') {
    ctx.type = 'image/x-icon';
    ctx.set('Cache-Control', 'public, max-age=31536000');
    ctx.body = fs.readFileSync('./ui/build/favicon.ico');
    return;
  }

  if (ctx.request.path === '/') return ctx.response.redirect('/mgmt/');
  if (ctx.request.path === '/index.html') return ctx.response.redirect('/mgmt/');
  await next();
});

///////////////////////////////////////////////////////////////////////////////////////////
//   API sections.
///////////////////////////////////////////////////////////////////////////////////////////
// !!! Start body-parser after proxies, see https://github.com/vagusX/koa-proxies/issues/55
// Start body-parser only for APIs, which requires the body.
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

auth.handle(system.handle(router));

router.all('/terraform/v1/mgmt/versions', async (ctx) => {
  ctx.body = utils.asResponse(0, {version: pkg.version});
});

app.use(router.routes());

///////////////////////////////////////////////////////////////////////////////////////////
console.log(`Run with cwd=${process.cwd()}, USE_DOCKER=${process.env.USE_DOCKER}`);
app.listen(consts.config.port, () => {
  console.log(`Server start on http://localhost:${consts.config.port}`);
});

