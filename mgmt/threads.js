'use strict';

const { Worker } = require("worker_threads");
const srs = require('./srs');
const releases = require('./releases');

exports.run = async () => {
  new Promise((resolve, reject) => {
    const worker = new Worker("./srs.js");
    worker.on('message', (msg) => {
      srs.metadata = msg.metadata;
    });
    worker.on('error', reject);
    worker.on('exit', (code) => {
      console.log(`thread #${srs.metadata.name}: exit with ${code}`)
      if (code !== 0) {
        return reject(new Error(`Worker #${srs.metadata.name}: stopped with exit code ${code}`));
      }
      resolve();
    });
  });

  new Promise((resolve, reject) => {
    const worker = new Worker("./releases.js");
    worker.on('message', (msg) => {
      releases.metadata = msg.metadata;
    });
    worker.on('error', reject);
    worker.on('exit', (code) => {
      console.log(`thread #${releases.metadata.name}: exit with ${code}`)
      if (code !== 0) {
        return reject(new Error(`Worker #${releases.metadata.name}: stopped with exit code ${code}`));
      }
      resolve();
    });
  });
};

