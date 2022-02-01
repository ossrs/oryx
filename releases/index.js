'use strict';

const dotenv = require('dotenv');
const Koa = require('koa');
const Router = require('koa-router');
const Cors = require('koa2-cors');
const BodyParser = require('koa-bodyparser');
const releases = require('./releases');
const pkg = require('./package.json');
const utils = require('js-core/utils');

// Default to local development env.
process.env.STAGE = process.env.STAGE || 'local';

// Try to read .env manually, for directly run node.
dotenv.config({path: `.env.${process.env.STAGE}`});

const app = new Koa();

app.use(Cors());
app.use(BodyParser());

const router = new Router();

releases.handle(router);

router.all('/terraform/v1/releases/versions', async (ctx) => {
  ctx.body = utils.asResponse(0, {version: pkg.version});
});

app.use(router.routes());

// Redirect /${stage}/xxx to /xxx
const prefixedRouter = new Router({
  prefix: `/${process.env.STAGE}`
});
prefixedRouter.use(router.routes());
app.use(prefixedRouter.routes());

const listenPort = process.env.PORT ? parseInt(process.env.PORT) : 2023;
app.listen(listenPort, () => {
  console.log(`Server start on http://localhost:${listenPort}`);
});

