'use strict';

const dotenv = require('dotenv');
const fs = require('fs');
const path = require('path');
['.', '..', '../..', '../mgmt', '../../mgmt'].map(envDir => {
  if (fs.existsSync(path.join(envDir, '.env'))) {
    dotenv.config({path: path.join(envDir, '.env')});
  }
});
console.log(`load envs MGMT_PASSWORD=${'*'.repeat(process.env.MGMT_PASSWORD?.length)}`);

const Koa = require('koa');
const Router = require('koa-router');
const Cors = require('koa2-cors');
const BodyParser = require('koa-bodyparser');
const settings = require('./settings');
const pkg = require('./package.json');
const utils = require('js-core/utils');

const app = new Koa();

app.use(Cors());
app.use(BodyParser());

const router = new Router();

settings.handle(router);

router.all('/terraform/v1/tencent/versions', async (ctx) => {
  ctx.body = utils.asResponse(0, {version: pkg.version});
});

app.use(router.routes());

app.listen(2020, () => {
  console.log(`Server start on http://localhost:2020`);
});

