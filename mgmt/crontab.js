'use strict';

const { isMainThread } = require("worker_threads");
const util = require('util');
const exec = util.promisify(require('child_process').exec);
const platform = require('./platform');

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  // We must initialize the thread first.
  const {region, registry} = await platform.init();
  console.log(`Thread #crontab: initialize region=${region}, registry=${registry}`);

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
  const {stdout} = await exec(`bash containers/bin/letsencrypt_renew`);
  console.log(`Thread #crontab: renew ssl ${stdout}`);
}

