'use strict';

const axios = require('axios');

exports.upgrade = {
  name: 'upgrade',
  releases: {
    stable: null,
    latest: null,
  },
};

const isDarwin = process.platform === 'darwin';

let region = null;
exports.region = async () => {
  if (isDarwin) return null;
  if (region) return region;

  if (process.env.REGION) {
    region = process.env.REGION;
    return region;
  }

  const {data} = await axios.get(`http://metadata.tencentyun.com/latest/meta-data/placement/region`);
  region = data;

  console.log(`Request region, data=${JSON.stringify(data)}, region=${region}`);
  return region;
};

let registry = null;
exports.registry = async () => {
  await exports.region();

  if (!region) return null;
  if (registry) return registry;

  registry = 'sgccr.ccs.tencentyun.com';
  ['ap-guangzhou', 'ap-shanghai', 'ap-nanjing', 'ap-beijing', 'ap-chengdu', 'ap-chongqing'].filter(v => {
    if (region.startsWith(v)) registry = 'ccr.ccs.tencentyun.com';
    return null;
  });

  console.log(`Setup registry to ${registry}, region is ${region}`);
  return registry;
};

exports.market = {
  srs: {
    name: 'srs-server',
    image: () => {
      let image = 'ossrs/lighthouse';
      if (process.env.NODE_ENV === 'development') image = 'ossrs/srs';
      if (process.env.SRS_DOCKER === 'srs') image = 'ossrs/srs';
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
    image: '${registry}/ossrs/srs-terraform:hooks-1',
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
    image: '${registry}/ossrs/prometheus',
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
    extras: isDarwin ? [] : ['--user=root'],
    container: {
      ID: null,
      State: null,
      Status: null,
    },
  },
  node_exporter: {
    name: 'node-exporter',
    image: '${registry}/ossrs/node-exporter',
    tcpPorts: () => isDarwin ? [9100] : [],
    udpPorts: [],
    command: () => isDarwin ? [] : ['--path.rootfs=/host'],
    logConfig: '--log-driver json-file --log-opt max-size=1g --log-opt max-file=3',
    volumes: isDarwin ? [] : ['/:/host:ro,rslave'],
    extras: () => isDarwin ? [] : ['--net=host', '--pid=host'],
    container: {
      ID: null,
      State: null,
      Status: null,
    },
  },
};

