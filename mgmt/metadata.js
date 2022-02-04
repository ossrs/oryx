'use strict';

exports.releases = {
  name: 'releases',
  releases: {
    stable: null,
    latest: null,
  },
};

exports.market = {
  srs: {
    name: 'srs-server',
    image: () => {
      let image = 'ossrs/lighthouse';
      if (process.env.NODE_ENV === 'development') image = 'ossrs/srs';
      if (process.env.SRS_DOCKER === 'srs') image = 'ossrs/srs';
      return `ccr.ccs.tencentyun.com/${image}:4`;
    },
    tcpPorts: [1935, 1985, 8080],
    udpPorts: [8000, 10080],
    command: './objs/srs -c conf/lighthouse.conf',
    logConfig: '--log-driver json-file --log-opt max-size=3g --log-opt max-file=3',
    extras: `-v ${process.cwd()}/containers/conf/srs.conf:/usr/local/srs/conf/lighthouse.conf`,
    container: {
      ID: null,
      State: null,
      Status: null,
    },
  },
  hooks: {
    name: 'srs-hooks',
    image: 'ccr.ccs.tencentyun.com/ossrs/srs-terraform:hooks-1',
    tcpPorts: [2021],
    udpPorts: [],
    command: 'node .',
    logConfig: '--log-driver json-file --log-opt max-size=1g --log-opt max-file=3',
    extras: `-v ${process.cwd()}/.env:/srs-terraform/hooks/.env`,
    container: {
      ID: null,
      State: null,
      Status: null,
    },
  },
};

