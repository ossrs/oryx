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

const app = new Koa();

app.use(Cors());
app.use(BodyParser());

const router = new Router();
app.use(mount('/mgmt', serve('./ui/build')));
router.all('/terraform/v1/mgmt/versions', async (ctx) => {
  releases.handle(ctx);
});
app.use(router.routes());

const config = {
  port: process.env.PORT || 2022,
}
app.listen(config.port, () => {
  console.log(`Server start on http://localhost:${config.port}`);
});

