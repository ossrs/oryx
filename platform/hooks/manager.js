'use strict';

const { Worker } = require("worker_threads");

let dvrWorker, vodWorker, recordWorker;

exports.run = async () => {
  // The DVR worker, for cloud storage.
  new Promise((resolve, reject) => {
    // Note that current work directory is platform, so we use ./hooks/xxx.js
    dvrWorker = new Worker("./hooks/dvrWorker.js");

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
    // Note that current work directory is platform, so we use ./hooks/xxx.js
    vodWorker = new Worker("./hooks/vodWorker.js");

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

  // The Record worker, for cloud storage.
  new Promise((resolve, reject) => {
    // Note that current work directory is platform, so we use ./hooks/xxx.js
    recordWorker = new Worker("./hooks/recordWorker.js");

    recordWorker.on('message', (msg) => {
      console.log('Thread #manager:', msg);
    });

    recordWorker.on('error', reject);

    recordWorker.on('exit', (code) => {
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
  } else if (msg.action === 'on_record_file') {
    recordWorker?.postMessage(msg);
  } else {
    console.error(`Thread #manager: Ignore message ${JSON.stringify(msg)}`);
  }
};

