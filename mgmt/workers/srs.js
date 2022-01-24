'use strict';

const { isMainThread, parentPort } = require("worker_threads");
const util = require('util');
const exec = util.promisify(require('child_process').exec);

const metadata = {
  name: 'srs-server',
  major: '4',
  container: {
    id: null,
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
    await new Promise(resolve => setTimeout(resolve, 3000));
  }
}

async function thread_main() {
  // Query the container from docker.
  let container = await queryContainer();
  if (container.ID) {
    metadata.container = container;
    console.log(`Thread #${metadata.name}: query ID=${container.ID}, State=${container.State}`);
  }

  // Restart the SRS container.
  if (!container.ID || container.State !== 'running') {
    await exec(`docker rm -f ${metadata.name}`);
    await exec(`docker run -d -it --restart always --privileged --name ${metadata.name} -p 1935:1935 -p 1985:1985 -p 8080:8080 -p 8000:8000/udp registry.cn-hangzhou.aliyuncs.com/ossrs/lighthouse:${metadata.major} ./objs/srs -c conf/docker.conf`);

    let container = await queryContainer();
    metadata.container = container;
    console.log(`Thread #${metadata.name}: create ID=${container.ID}, State=${container.State}`);
  }

  // Update the metadata to main thread.
  parentPort.postMessage({metadata});
}

async function queryContainer() {
  // See https://docs.docker.com/config/formatting/
  const {stdout} = await exec(`docker ps -a -f name=${metadata.name} --format '{{json .}}'`);
  return stdout ? JSON.parse(stdout) : {};
}
