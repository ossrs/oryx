'use strict';

const platform = require('./platform');
const metadata = require('js-core/metadata');

exports.market = {
  platform: {
    name: 'platform',
    image: async () => {
      const registry = await platform.registry();
      return `${registry}/ossrs/srs-terraform:platform-1`;
    },
    tcpPorts: [2024],
    udpPorts: [],
    command: ['node', '.'],
    logConfig: [
      '--log-driver=json-file',
      '--log-opt', 'max-size=1g',
      '--log-opt', 'max-file=3',
    ],
    volumes: [
      `${process.cwd()}/.env:/usr/local/srs-terraform/platform/.env`,
      // We mount the containers to mgmt in platform container, which links to platform.
      `${process.cwd()}/containers:/usr/local/srs-terraform/mgmt/containers`,
    ],
    extras: [
      '--env', `SRS_DOCKER=${process.env.SRS_DOCKER || ''}`,
      '--env', `USE_DOCKER=${process.env.USE_DOCKER || ''}`,
      '--env', `NODE_ENV=${process.env.NODE_ENV || ''}`,
    ],
    // The running status, updated by fetchContainer.
    container: {
      ID: null,
      State: null,
      Status: null,
    },
  },
  // It's only a hint. It will be built by fetchContainer.
  srs: {
    name: metadata.market.srs.name,
    // The running status, updated by fetchContainer.
    container: {
      ID: null,
      State: null,
      Status: null,
    },
  },
};

