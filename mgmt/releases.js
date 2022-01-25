'use strict';

const { isMainThread, parentPort } = require("worker_threads");
const { spawn } = require('child_process');
const pkg = require('./package.json');
const axios = require('axios');
const semver = require('semver');

const metadata = {
  name: 'mgmt-vers',
  releases: {
    stable: null,
    latest: null,
  },
};
exports.metadata = metadata;

if (!isMainThread) {
  _thread_main();
}

async function _thread_main() {
  while (true) {
    try {
      await thread_main();
    } catch (e) {
      console.error(`thread ${metadata.name} err`, e);
    }
    await new Promise(resolve => setTimeout(resolve, 3600 * 1000));
  }
}

async function thread_main() {
  console.log(`Thread #${metadata.name}: current version=v${pkg.version}`);

  // Wait for a while to request version.
  await new Promise(resolve => setTimeout(resolve, 5 * 1000));
  console.log(`Thread #${metadata.name}: request by version=v${pkg.version}`);

  const {data} = await axios.get('http://api.ossrs.net/terraform/v1/releases', {
    params: {
      version: `v${pkg.version}`,
      ts: new Date().getTime(),
    }
  });
  metadata.releases = data;
  console.log(`Thread #${metadata.name}: request, version=v${pkg.version}, response=${JSON.stringify(data)}`);

  if (metadata.releases && metadata.releases.stable && semver.lt(`v${pkg.version}`, metadata.releases.stable)) {
    console.log(`Thread #${metadata.name}: upgrade from v${pkg.version} to stable ${metadata.releases.stable}`);

    await new Promise((resolve, reject) => {
      const child = spawn('bash', ['upgrade', metadata.releases.stable]);
      child.stdout.on('data', (chunk) => {
        console.log(`Thread #${metadata.name}: ${chunk.toString()}`);
      });
      child.stderr.on('data', (chunk) => {
        console.log(`Thread #${metadata.name}: ${chunk.toString()}`);
      });
      child.on('close', (code) => {
        console.log(`Thread #${metadata.name}: upgrading exited with code ${code}`);
        if (code !== 0) return reject(code);
        resolve();
      });
    });
  }

  // Update the metadata to main thread.
  parentPort.postMessage({metadata});
}

