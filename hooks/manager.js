'use strict';

const { Worker } = require("worker_threads");

let dvrWorker, vodWorker;

exports.run = async () => {
  // The DVR worker, for cloud storage.
  new Promise((resolve, reject) => {
    dvrWorker = new Worker("./dvrWorker.js");

    dvrWorker.on('message', (msg) => {
      console.log('Thread #manager:', msg);
    });

    dvrWorker.on('error', reject);

    dvrWorker.on('exit', (code) => {
      console.log(`Thread #manager: exit with ${code}`);
      if (code !== 0) {
        return reject(new Error(`Thread #manager: stopped with exit code ${code}`));
      }
      resolve();
    });
  });

  // The VoD worker, for cloud VoD.
  new Promise((resolve, reject) => {
    vodWorker = new Worker("./vodWorker.js");

    vodWorker.on('message', (msg) => {
      console.log('Thread #manager:', msg);
    });

    vodWorker.on('error', reject);

    vodWorker.on('exit', (code) => {
      console.log(`Thread #manager: exit with ${code}`);
      if (code !== 0) {
        return reject(new Error(`Thread #manager: stopped with exit code ${code}`));
      }
      resolve();
    });
  });
};

exports.postMessage = (msg) => {
  if (msg.action === 'on_dvr_file') {
    dvrWorker?.postMessage(msg);
  } else if (msg.action === 'on_vod_file') {
    vodWorker?.postMessage(msg);
  } else {
    console.error(`Thread #manager: Ignore message ${JSON.stringify(msg)}`);
  }
};

