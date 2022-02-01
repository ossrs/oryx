'use strict';

const { Worker } = require("worker_threads");
const metadata = require('./metadata');

exports.run = async () => {
  new Promise((resolve, reject) => {
    const worker = new Worker("./releases.js");
    worker.on('message', (msg) => {
      metadata.releases = msg.metadata.releases;
    });
    worker.on('error', reject);
    worker.on('exit', (code) => {
      console.log(`thread #${metadata.releases.name}: exit with ${code}`);
      if (code !== 0) {
        return reject(new Error(`Worker #${metadata.releases.name}: stopped with exit code ${code}`));
      }
      resolve();
    });
  });

  if (process.env.USE_DOCKER === 'false') {
    console.warn(`run without docker, please start components by npm start`);
    return;
  }

  new Promise((resolve, reject) => {
    const worker = new Worker("./srs.js");
    worker.on('message', (msg) => {
      metadata.srs = msg.metadata.srs;
    });
    worker.on('error', reject);
    worker.on('exit', (code) => {
      console.log(`thread #${metadata.srs.name}: exit with ${code}`);
      if (code !== 0) {
        return reject(new Error(`Worker #${metadata.srs.name}: stopped with exit code ${code}`));
      }
      resolve();
    });
  });

  new Promise((resolve, reject) => {
    const worker = new Worker("./market.js");
    worker.on('message', (msg) => {
      metadata.market.hooks = msg.metadata.hooks;
    });
    worker.on('error', reject);
    worker.on('exit', (code) => {
      console.log(`thread #market: exit with ${code}`);
      if (code !== 0) {
        return reject(new Error(`Worker #market: stopped with exit code ${code}`));
      }
      resolve();
    });
  });
};

