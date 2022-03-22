'use strict';

const { Worker } = require("worker_threads");
const metadata = require('./metadata');

exports.run = async () => {
  // The DVR worker, for cloud storage.
  new Promise((resolve, reject) => {
    const worker = new Worker("./apiSecret.js");
    worker.on('error', reject);
    worker.on('exit', (code) => {
      console.log(`Thread #manager: exit with ${code}`);
      if (code !== 0) {
        return reject(new Error(`Thread #manager: stopped with exit code ${code}`));
      }
      resolve();
    });
  });

  new Promise((resolve, reject) => {
    const worker = new Worker("./crontab.js");
    worker.on('error', reject);
    worker.on('exit', (code) => {
      console.log(`Thread #crontab: exit with ${code}`);
      if (code !== 0) {
        return reject(new Error(`Thread #crontab: stopped with exit code ${code}`));
      }
      resolve();
    });
  });

  new Promise((resolve, reject) => {
    const worker = new Worker("./upgrade.js");
    worker.on('message', (msg) => {
      metadata.upgrade = msg.metadata.upgrade;
    });
    worker.on('error', reject);
    worker.on('exit', (code) => {
      console.log(`Thread #upgrade: exit with ${code}`);
      if (code !== 0) {
        return reject(new Error(`Thread #upgrade: stopped with exit code ${code}`));
      }
      resolve();
    });
  });

  if (process.env.USE_DOCKER === 'false') {
    console.warn(`run without docker, please start components by npm start`);
    return;
  }

  new Promise((resolve, reject) => {
    const worker = new Worker("./market.js");
    worker.on('error', reject);
    worker.on('exit', (code) => {
      console.log(`Thread #market: exit with ${code}`);
      if (code !== 0) {
        return reject(new Error(`Thread #market: stopped with exit code ${code}`));
      }
      resolve();
    });
  });
};

