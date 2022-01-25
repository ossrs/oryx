'use strict';

const dotenv = require('dotenv');
const Koa = require('koa');
const Router = require('koa-router');
const Cors = require('koa2-cors');
const BodyParser = require('koa-bodyparser');
const releases = require('./releases');

// Default to local development env.
process.env.STAGE = process.env.STAGE || 'local';

// Try to read .env manually, for directly run node.
dotenv.config({path: `.env.${process.env.STAGE}`});

const app = new Koa();

app.use(Cors());
app.use(BodyParser());

const router = new Router();
router.all('/terraform/v1/releases', async (ctx) => {
  releases.handle(ctx);
});
app.use(router.routes());

// Redirect /${stage}/xxx to /xxx
const prefixedRouter = new Router({
  prefix: `/${process.env.STAGE}`
});
prefixedRouter.use(router.routes());
app.use(prefixedRouter.routes());

app.listen(9000, () => {
  console.log(`Server start on http://localhost:9000`);
});

