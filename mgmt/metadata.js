'use strict';

const platform = require('./platform');

exports.upgrade = {
  name: 'upgrade',
  releases: {
    stable: null,
    latest: null,
  },
};

exports.market = {
  srs: {
    name: 'srs-server',
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
    volumes: [`${process.cwd()}/containers/conf/srs.conf:/usr/local/srs/conf/lighthouse.conf`],
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
    volumes: [`${process.cwd()}/.env:/srs-terraform/hooks/.env`],
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

