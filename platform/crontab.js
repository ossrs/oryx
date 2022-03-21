'use strict';

const { isMainThread } = require("worker_threads");
const util = require('util');
const exec = util.promisify(require('child_process').exec);
const helper = require('./helper');
const platform = require('./platform');

if (!isMainThread) {
  threadMain();
}

async function threadMain() {
  // We must initialize the thread first.
  await platform.init();
  console.log(`Thread #crontab: initialize`);

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
  const {stdout: liveDomains} = await exec('ls -d containers/etc/letsencrypt/live/*/ |wc -l');
  const nnLiveDomains = liveDomains && parseInt(liveDomains.trim());
  if (!nnLiveDomains) return console.log(`Thread #crontab: No domains in containers/etc/letsencrypt/live/`);

  const {stdout, renewOk} = await helper.execApi('renewLetsEncrypt');
  console.log(`Thread #crontab: renew ssl updated=${renewOk}, message is ${stdout}`);
}

