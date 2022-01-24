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
const Router = require('koa-router');
const Cors = require('koa2-cors');
const BodyParser = require('koa-bodyparser');
const serve = require('koa-static');
const mount  = require('koa-mount');
const releases = require('./releases');
const auth = require('./auth');
const utils = require('./utils');
const status = require('./status');

const app = new Koa();

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

app.use(Cors());
app.use(BodyParser());

const router = new Router();

// For react, static files server.
app.use(mount('/mgmt', serve('./ui/build')));
router.all('/', async (ctx) => {
  ctx.response.redirect('/mgmt/');
});

// For backend APIs.
releases.handle(auth.handle(status.handle(router)));
app.use(router.routes());

// For react, static files server.
// See https://stackoverflow.com/a/52464577/17679565
app.use(async (ctx, next) => {
  if (ctx.request.path.indexOf('/mgmt/') === 0) {
    ctx.type = 'text/html';
    ctx.body = fs.readFileSync('./ui/build/index.html');
    return;
  }
  await next();
});

const config = {
  port: process.env.PORT || 2022,
}
app.listen(config.port, () => {
  console.log(`Server start on http://localhost:${config.port}`);
});

