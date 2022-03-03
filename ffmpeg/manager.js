'use strict';

const { Worker } = require("worker_threads");

let forwardWorker;

exports.run = async () => {
  // The DVR worker, for cloud storage.
  new Promise((resolve, reject) => {
    forwardWorker = new Worker("./forwardWorker.js");

    forwardWorker.on('message', (msg) => {
      console.log('Thread #manager:', msg);
    });

    forwardWorker.on('error', reject);

    forwardWorker.on('exit', (code) => {
      console.log(`Thread #manager: exit with ${code}`);
      if (code !== 0) {
        return reject(new Error(`Thread #manager: stopped with exit code ${code}`));
      }
      resolve();
    });
  });
};

exports.postMessage = (msg) => {
  console.error(`Thread #manager: Ignore message ${JSON.stringify(msg)}`);
};

