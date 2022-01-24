'use strict';

const { Worker } = require("worker_threads");
const srs = require('./workers/srs');

exports.run = async () => {
  return new Promise((resolve, reject) => {
    const worker = new Worker("./workers/srs.js");
    worker.on('message', (msg) => {
      srs.metadata = msg.metadata;
    });
    worker.on('error', reject);
    worker.on('exit', (code) => {
      console.log(`thread: ${srs.metadata.name} exit with ${code}`)
      if (code !== 0) {
        return reject(new Error(`Worker stopped with exit code ${code}`));
      }
      resolve();
    });
  });
};

