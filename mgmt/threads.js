'use strict';

const { Worker } = require("worker_threads");
const metadata = require('./metadata');

exports.run = async () => {
  if (process.env.USE_DOCKER === 'false') {
    console.warn(`run without docker, please start components by npm start`);
    return;
  }

  new Promise((resolve, reject) => {
    const worker = new Worker("./market.js");
    worker.on('message', (msg) => {
      Object.keys(msg.metadata).map(e => {
        metadata.market[e].container = msg.metadata[e];
      });
      //console.log(`update metadata by ${JSON.stringify(msg)} to ${JSON.stringify(metadata)}`);
    });
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

