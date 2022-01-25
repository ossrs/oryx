'use strict';

const dotenv = require('dotenv');
const fs = require('fs');
const path = require('path');
['.', '..', '../..'].map(envDir => {
  if (fs.existsSync(path.join(envDir, '.env'))) {
    dotenv.config({path: path.join(envDir, '.env')});
  }
});

const Koa = require('koa');
const proxy = require('koa-proxies');
const Router = require('koa-router');
const Cors = require('koa2-cors');
const BodyParser = require('koa-bodyparser');
const serve = require('koa-static');
const mount  = require('koa-mount');
const releases = require('./releases');
const auth = require('./auth');
const utils = require('./utils');
const system = require('./system');
const threads = require('./threads');
const hooks = require('./hooks');

const app = new Koa();

// Always enable CORS and parse body.
app.use(Cors());
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

// For backend APIs, with specified path.
if (true) {
  const router = new Router();
  releases.handle(auth.handle(system.handle(router)));
  hooks.handle(router);
  app.use(router.routes());
}

// For react-router, can't match the files.
// See https://stackoverflow.com/a/52464577/17679565
app.use(async (ctx, next) => {
  // The react-router should never with extensions, and without multiple path.
  const isReactRouter = ctx.request.path.indexOf('/mgmt/') === 0
    && ctx.request.path.indexOf('.') === -1
    && ctx.request.path.match(/\//g).length === 2
  if (isReactRouter) {
    ctx.type = 'text/html';
    ctx.body = fs.readFileSync('./ui/build/index.html');
    return;
  }

  await next();
});

// For react, static files server.
app.use(mount('/mgmt', serve('./ui/build')));
// For homepage, use mgmt.
app.use(async (ctx, next) => {
  if (ctx.request.path === '/') return ctx.response.redirect(['/mgmt/', ctx.request.querystring].filter(e => e).join('?'));
  if (ctx.request.path === '/index.html') return ctx.response.redirect(['/mgmt/', ctx.request.querystring].filter(e => e).join('?'));
  await next();
});

// Proxy for special path of SRS>
app.use(async (ctx, next) => {
  if (ctx.request.path === '/console') return ctx.response.redirect(['/console/', ctx.request.querystring].filter(e => e).join('?'));
  if (ctx.request.path === '/players') return ctx.response.redirect(['/players/', ctx.request.querystring].filter(e => e).join('?'));
  await next();
});

// Proxy to SRS HTTP streaming, console and player.
// See https://github.com/vagusX/koa-proxies
app.use(proxy('/api/', {
  target: 'http://127.0.0.1:1985/',
}));
app.use(proxy('/rtc/', {
  target: 'http://127.0.0.1:1985/',
}));
app.use(proxy('/', {
  target: 'http://127.0.0.1:8080',
}));

// Start all workers threads.
threads.run();

const config = {
  port: process.env.PORT || 2022,
}
app.listen(config.port, () => {
  console.log(`Server start on http://localhost:${config.port}`);
});

