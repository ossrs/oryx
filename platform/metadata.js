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
      const cloud = await platform.cloud();

      let image = 'ossrs/srs';
      if (cloud === 'TENCENT') image = 'ossrs/lighthouse';
      if (cloud === 'DO') image = 'ossrs/droplet';
      if (cloud === 'BT') image = 'ossrs/srs';
      if (process.env.NODE_ENV === 'development') image = 'ossrs/srs';
      if (process.env.SRS_DOCKER === 'srs') image = 'ossrs/srs';

      const registry = await platform.registry();
      return `${registry}/${image}:4`;
    },
    tcpPorts: [1935, 1985, 8080],
    udpPorts: [8000, 10080],
    command: ['./objs/srs', '-c', 'conf/lighthouse.conf'],
    logConfig: [
      '--log-driver=json-file',
      '--log-opt', 'max-size=3g',
      '--log-opt', 'max-file=3',
    ],
    // Note that we should use platform.cwd() which is cwd of mgmt.
    volumes: () => [
      `${platform.cwd()}/containers/conf/srs.release.conf:/usr/local/srs/conf/lighthouse.conf`,
      // Note that we mount the whole www directory, so we must build the static files such as players.
      `${platform.cwd()}/containers/objs/nginx/html:/usr/local/srs/objs/nginx/html`,
      // We must mount the player and console because the HTTP home of SRS is overwrite by DVR.
      `${platform.cwd()}/containers/www/players:/usr/local/srs/www/players`,
      `${platform.cwd()}/containers/www/console:/usr/local/srs/www/console`,
      // For coredump, save to /cores.
      ...(platform.isDarwin ? [] : ['/cores:/cores']),
    ],
    extras: () => [
      ...(platform.isDarwin ? [] : ['--network=srs-cloud']),
      ...(platform.isDarwin ? [] : ['--ulimit', 'core=-1']),
    ],
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
    command: ['./objs/srs', '-c', 'conf/lighthouse.conf'],
    logConfig: [
      '--log-driver=json-file',
      '--log-opt', 'max-size=3g',
      '--log-opt', 'max-file=3',
    ],
    volumes: () => [
      `${platform.cwd()}/containers/conf/srs.dev.conf:/usr/local/srs/conf/lighthouse.conf`,
      // Note that we mount the whole www directory, so we must build the static files such as players.
      `${platform.cwd()}/containers/objs/nginx/html:/usr/local/srs/objs/nginx/html`,
      // We must mount the player and console because the HTTP home of SRS is overwrite by DVR.
      `${platform.cwd()}/containers/www/players:/usr/local/srs/www/players`,
      `${platform.cwd()}/containers/www/console:/usr/local/srs/www/console`,
      // For coredump, save to /cores.
      ...(platform.isDarwin ? [] : ['/cores:/cores']),
    ],
    extras: () => [
      ...(platform.isDarwin ? [] : ['--network=srs-cloud']),
      ...(platform.isDarwin ? [] : ['--ulimit', 'core=-1']),
    ],
  },
  hooks: {
    name: metadata.market.hooks.name,
    image: async () => {
      const registry = await platform.registry();
      return `${registry}/ossrs/srs-cloud:hooks-1`;
    },
    tcpPorts: [2021],
    udpPorts: [],
    command: ['node', '.'],
    logConfig: [
      '--log-driver=json-file',
      '--log-opt', 'max-size=1g',
      '--log-opt', 'max-file=3',
    ],
    volumes: () => [
      `${platform.cwd()}/.env:/usr/local/srs-cloud/hooks/.env`,
      // We mount the containers to mgmt in platform container, which links to platform.
      `${platform.cwd()}/containers:/usr/local/srs-cloud/mgmt/containers`,
      // We mount the containers to mgmt in hooks container, which links to hooks.
      `${platform.cwd()}/containers/objs/nginx/html:/usr/local/srs-cloud/mgmt/containers/objs/nginx/html`,
      `${platform.cwd()}/containers/data/dvr:/usr/local/srs-cloud/mgmt/containers/data/dvr`,
      `${platform.cwd()}/containers/data/vod:/usr/local/srs-cloud/mgmt/containers/data/vod`,
    ],
    extras: () => [
      ...(platform.isDarwin ? [] : ['--network=srs-cloud']),
      ...(platform.isDarwin ? ['--env=REDIS_HOST=host.docker.internal'] : ['--env=REDIS_HOST=redis']),
    ],
  },
  tencent: {
    name: metadata.market.tencent.name,
    image: async () => {
      const registry = await platform.registry();
      return `${registry}/ossrs/srs-cloud:tencent-1`;
    },
    tcpPorts: [2020],
    udpPorts: [],
    command: ['node', '.'],
    logConfig: [
      '--log-driver=json-file',
      '--log-opt', 'max-size=1g',
      '--log-opt', 'max-file=3',
    ],
    volumes: () => [
      `${platform.cwd()}/.env:/usr/local/srs-cloud/tencent/.env`,
      // We mount the containers to mgmt in platform container, which links to platform.
      `${platform.cwd()}/containers:/usr/local/srs-cloud/mgmt/containers`,
    ],
    extras: () => [
      ...(platform.isDarwin ? [] : ['--network=srs-cloud']),
      ...(platform.isDarwin ? ['--env=REDIS_HOST=host.docker.internal'] : ['--env=REDIS_HOST=redis']),
    ],
  },
  ffmpeg: {
    name: metadata.market.ffmpeg.name,
    image: async () => {
      const registry = await platform.registry();
      return `${registry}/ossrs/srs-cloud:ffmpeg-1`;
    },
    tcpPorts: [2019],
    udpPorts: [],
    command: ['node', '.'],
    logConfig: [
      '--log-driver=json-file',
      '--log-opt', 'max-size=1g',
      '--log-opt', 'max-file=3',
    ],
    volumes: () => [
      `${platform.cwd()}/.env:/usr/local/srs-cloud/ffmpeg/.env`,
      // We mount the containers to mgmt in platform container, which links to platform.
      `${platform.cwd()}/containers:/usr/local/srs-cloud/mgmt/containers`,
    ],
    extras: () => [
      ...(platform.isDarwin ? [] : ['--network=srs-cloud']),
      ...(platform.isDarwin ? ['--env=REDIS_HOST=host.docker.internal'] : ['--env=REDIS_HOST=redis']),
    ],
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
    logConfig: [
      '--log-driver=json-file',
      '--log-opt', 'max-size=1g',
      '--log-opt', 'max-file=3',
    ],
    volumes: () => {
      const config = platform.isDarwin ? 'prometheus.darwin.yml' : 'prometheus.yml';
      return [
        `${platform.cwd()}/containers/conf/${config}:/etc/prometheus/prometheus.yml`,
        `${platform.cwd()}/containers/data/prometheus:/prometheus`,
      ];
    },
    extras: () => [
      ...(platform.isDarwin ? [] : ['--network=srs-cloud']),
      ...(platform.isDarwin ? [] : ['--user=root']),
    ],
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
    logConfig: [
      '--log-driver=json-file',
      '--log-opt', 'max-size=1g',
      '--log-opt', 'max-file=3',
    ],
    volumes: () => {
      return platform.isDarwin ? [] : ['/:/host:ro,rslave'];
    },
    extras: () => [
      platform.isDarwin ? '--network=host' : '--network=srs-cloud',
      ...(platform.isDarwin ? [] : ['--pid=host']),
    ],
  },
  // The bellow configurations are only a hint, because they are managed by mgmt.
  platform: {name: 'platform'},
  redis: {name: 'redis'},
};

