'use strict';

const { Worker } = require("worker_threads");

let forwardWorker, vLiveWorker;

exports.run = async () => {
  // The DVR worker, for cloud storage.
  new Promise((resolve, reject) => {
    // Note that current work directory is platform, so we use ./ffmpeg/xxx.js
    forwardWorker = new Worker("./ffmpeg/forwardWorker.js");

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

  // The DVR worker, for cloud storage.
  new Promise((resolve, reject) => {
    // Note that current work directory is platform, so we use ./ffmpeg/xxx.js
    vLiveWorker = new Worker("./ffmpeg/vLiveWorker.js");

    vLiveWorker.on('message', (msg) => {
      console.log('Thread #manager:', msg);
    });

    vLiveWorker.on('error', reject);

    vLiveWorker.on('exit', (code) => {
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

