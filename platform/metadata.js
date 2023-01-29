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
    // For China, see https://cr.console.aliyun.com/repository/cn-hangzhou/ossrs/srs/images
    // For Global, see https://hub.docker.com/r/ossrs/srs/tags
    image: async () => {
      const cloud = await platform.cloud();

      const image = 'ossrs/srs';
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
    // For China, see https://cr.console.aliyun.com/repository/cn-hangzhou/ossrs/srs/images
    // For Global, see https://hub.docker.com/r/ossrs/srs/tags
    image: async () => {
      const image = 'ossrs/srs';
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
  // The bellow configurations are only a hint, because they are managed by mgmt.
  platform: {name: 'platform'},
  redis: {name: 'redis'},
};

