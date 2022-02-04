'use strict';

const { isMainThread } = require("worker_threads");
const util = require('util');
const exec = util.promisify(require('child_process').exec);

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  while (true) {
    try {
      await doThreadMain();
    } catch (e) {
      console.error(`Thread #crontab: err`, e);
    } finally {
      await new Promise(resolve => setTimeout(resolve, 3600 * 1000));
    }
  }
}

async function doThreadMain() {
  if (process.platform === 'darwin') {
    console.log('Thread #crontab: ignore for Darwin');
    return;
  }

  console.log(`Thread #crontab: auto renew the Let's Encrypt ssl`);
  const {stdout} = await exec(`bash letsencrypt/crontab_renew`);
  console.log(`Thread #crontab: renew ssl ${stdout}`);
}

