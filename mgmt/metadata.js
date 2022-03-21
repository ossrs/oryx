'use strict';

const platform = require('./platform');
const metadata = require('js-core/metadata');

exports.upgrade = {
  name: 'upgrade',
  releases: {
    stable: null,
    latest: null,
  },
};

exports.market = {
  srs: {
    name: metadata.market.srs.name,
    // For China, see https://console.cloud.tencent.com/tcr/repository/details/ccr/ossrs/lighthouse/1
    // For Global, see https://console.cloud.tencent.com/tcr/repository/details/ccr/ossrs/lighthouse/9
    image: async () => {
      let image = 'ossrs/lighthouse';
      if (process.env.NODE_ENV === 'development') image = 'ossrs/srs';
      if (process.env.SRS_DOCKER === 'srs') image = 'ossrs/srs';
      const registry = await platform.registry();
      return `${registry}/${image}:4`;
    },
    tcpPorts: [1935, 1985, 8080],
    udpPorts: [8000, 10080],
    command: ['./objs/srs -c conf/lighthouse.conf'],
    logConfig: '--log-driver json-file --log-opt max-size=3g --log-opt max-file=3',
    volumes: [
      `${process.cwd()}/containers/conf/srs.release.conf:/usr/local/srs/conf/lighthouse.conf`,
      // Note that we mount the whole www directory, so we must build the static files such as players.
      `${process.cwd()}/containers/objs/nginx/html:/usr/local/srs/objs/nginx/html`,
      // We must mount the player and console because the HTTP home of SRS is overwrite by DVR.
      `${process.cwd()}/containers/www/players:/usr/local/srs/www/players`,
      `${process.cwd()}/containers/www/console:/usr/local/srs/www/console`,
    ],
    extras: [],
    container: {
      ID: null,
      State: null,
      Status: null,
    },
  },
  srsDev: {
    name: metadata.market.srsDev.name,
    // For China, see https://console.cloud.tencent.com/tcr/repository/details/ccr/ossrs/lighthouse/1
    // For Global, see https://console.cloud.tencent.com/tcr/repository/details/ccr/ossrs/lighthouse/9
    image: async () => {
      let image = 'ossrs/lighthouse';
      if (process.env.NODE_ENV === 'development') image = 'ossrs/srs';
      if (process.env.SRS_DOCKER === 'srs') image = 'ossrs/srs';
      const registry = await platform.registry();
      return `${registry}/${image}:5`;
    },
    tcpPorts: [1935, 1985, 8080],
    udpPorts: [8000, 10080],
    command: ['./objs/srs -c conf/lighthouse.conf'],
    logConfig: '--log-driver json-file --log-opt max-size=3g --log-opt max-file=3',
    volumes: [
      `${process.cwd()}/containers/conf/srs.dev.conf:/usr/local/srs/conf/lighthouse.conf`,
      // Note that we mount the whole www directory, so we must build the static files such as players.
      `${process.cwd()}/containers/objs/nginx/html:/usr/local/srs/objs/nginx/html`,
      // We must mount the player and console because the HTTP home of SRS is overwrite by DVR.
      `${process.cwd()}/containers/www/players:/usr/local/srs/www/players`,
      `${process.cwd()}/containers/www/console:/usr/local/srs/www/console`,
    ],
    extras: [],
    container: {
      ID: null,
      State: null,
      Status: null,
    },
  },
  hooks: {
    name: 'srs-hooks',
    image: async () => {
      const registry = await platform.registry();
      return `${registry}/ossrs/srs-terraform:hooks-1`;
    },
    tcpPorts: [2021],
    udpPorts: [],
    command: ['node .'],
    logConfig: '--log-driver json-file --log-opt max-size=1g --log-opt max-file=3',
    volumes: [
      `${process.cwd()}/.env:/usr/local/srs-terraform/hooks/.env`,
      // We mount the containers to mgmt in hooks container, which links to hooks.
      `${process.cwd()}/containers/objs/nginx/html:/usr/local/srs-terraform/mgmt/containers/objs/nginx/html`,
      `${process.cwd()}/containers/data/dvr:/usr/local/srs-terraform/mgmt/containers/data/dvr`,
      `${process.cwd()}/containers/data/vod:/usr/local/srs-terraform/mgmt/containers/data/vod`,
    ],
    extras: [],
    container: {
      ID: null,
      State: null,
      Status: null,
    },
  },
  tencent: {
    name: 'tencent-cloud',
    image: async () => {
      const registry = await platform.registry();
      return `${registry}/ossrs/srs-terraform:tencent-1`;
    },
    tcpPorts: [2020],
    udpPorts: [],
    command: ['node .'],
    logConfig: '--log-driver json-file --log-opt max-size=1g --log-opt max-file=3',
    volumes: [
      `${process.cwd()}/.env:/usr/local/srs-terraform/tencent/.env`,
    ],
    extras: [],
    container: {
      ID: null,
      State: null,
      Status: null,
    },
  },
  ffmpeg: {
    name: 'ffmpeg',
    image: async () => {
      const registry = await platform.registry();
      return `${registry}/ossrs/srs-terraform:ffmpeg-1`;
    },
    tcpPorts: [2019],
    udpPorts: [],
    command: ['node .'],
    logConfig: '--log-driver json-file --log-opt max-size=1g --log-opt max-file=3',
    volumes: [
      `${process.cwd()}/.env:/usr/local/srs-terraform/ffmpeg/.env`,
    ],
    extras: [],
    container: {
      ID: null,
      State: null,
      Status: null,
    },
  },
  platform: {
    name: 'platform',
    image: async () => {
      const registry = await platform.registry();
      return `${registry}/ossrs/srs-terraform:platform-1`;
    },
    tcpPorts: [2024],
    udpPorts: [],
    command: ['node .'],
    logConfig: '--log-driver json-file --log-opt max-size=1g --log-opt max-file=3',
    volumes: [
      `${process.cwd()}/.env:/usr/local/srs-terraform/platform/.env`,
      // We mount the containers to mgmt in platform container, which links to platform.
      `${process.cwd()}/containers:/usr/local/srs-terraform/mgmt/containers`,
    ],
    extras: [],
    container: {
      ID: null,
      State: null,
      Status: null,
    },
  },
  prometheus: {
    name: 'prometheus',
    image: async () => {
      const registry = await platform.registry();
      return `${registry}/ossrs/prometheus`;
    },
    tcpPorts: [9090],
    udpPorts: [],
    command: [
      '--storage.tsdb.path=/prometheus',
      '--config.file=/etc/prometheus/prometheus.yml',
      '--web.external-url=http://localhost:9090/prometheus/',
    ],
    logConfig: '--log-driver json-file --log-opt max-size=1g --log-opt max-file=3',
    volumes: [
      `${process.cwd()}/containers/conf/prometheus.yml:/etc/prometheus/prometheus.yml`,
      `${process.cwd()}/containers/data/prometheus:/prometheus`,
    ],
    extras: platform.isDarwin ? [] : ['--user=root'],
    container: {
      ID: null,
      State: null,
      Status: null,
    },
  },
  node_exporter: {
    name: 'node-exporter',
    image: async () => {
      const registry = await platform.registry();
      return `${registry}/ossrs/node-exporter`;
    },
    tcpPorts: () => platform.isDarwin ? [9100] : [],
    udpPorts: [],
    command: () => platform.isDarwin ? [] : ['--path.rootfs=/host'],
    logConfig: '--log-driver json-file --log-opt max-size=1g --log-opt max-file=3',
    volumes: platform.isDarwin ? [] : ['/:/host:ro,rslave'],
    extras: () => platform.isDarwin ? [] : ['--net=host', '--pid=host'],
    container: {
      ID: null,
      State: null,
      Status: null,
    },
  },
};

