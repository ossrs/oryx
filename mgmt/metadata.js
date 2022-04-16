'use strict';

const platform = require('./platform');
const metadata = require('js-core/metadata');
const pkg = require('./package.json');

exports.market = {
  platform: {
    name: 'platform',
    image: async () => {
      const registry = await platform.registry();
      // Note that we use the same version of mgmt for image, to support latest and stable version of container, and to
      // avoid upgrading to latest for stable mgmt. The mgmt is locked with the correct version of platform.
      return `${registry}/ossrs/srs-cloud:platform-v${pkg.version}`;
    },
    tcpPorts: [2024],
    udpPorts: [],
    command: ['node', '.'],
    logConfig: [
      '--log-driver=json-file',
      '--log-opt=max-size=1g',
      '--log-opt=max-file=3',
    ],
    volumes: [
      `${process.cwd()}/.env:/usr/local/srs-cloud/platform/.env`,
      // We mount the containers to mgmt in platform container, which links to platform.
      `${process.cwd()}/containers:/usr/local/srs-cloud/mgmt/containers`,
    ],
    extras: [
      `--network=srs-cloud`,
      `--env=SRS_DOCKER=${process.env.SRS_DOCKER || ''}`,
      `--env=USE_DOCKER=${process.env.USE_DOCKER || ''}`,
      // If use docker, should always use production to connect to redis.
      `--env=NODE_ENV=production`,
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

