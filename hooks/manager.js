'use strict';

const { Worker } = require("worker_threads");

let worker = null;

exports.run = async () => {
  worker = new Worker("./worker.js");

  new Promise((resolve, reject) => {
    worker.on('message', (msg) => {
      console.log('Thread #manager:', msg);
    });
    worker.on('error', reject);
    worker.on('exit', (code) => {
      console.log(`Thread #manager: exit with ${code}`);
      if (code !== 0) {
        return reject(new Error(`Thread #manager: stopped with exit code ${code}`));
      }
      resolve();
    });
  });
};

exports.postMessage = (msg) => {
  worker?.postMessage(msg);
};

