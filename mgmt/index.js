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

// Static file server for UI.
app.use(mount('/mgmt', serve('./ui/build')));
router.all('/', async (ctx) => {
  ctx.response.redirect('/mgmt/');
});

releases.handle(auth.handle(router));
app.use(router.routes());

const config = {
  port: process.env.PORT || 2022,
}
app.listen(config.port, () => {
  console.log(`Server start on http://localhost:${config.port}`);
});

