'use strict';

exports.srs = {
  name: 'srs-server',
  major: '4',
  container: {
    ID: null,
    State: null,
    Status: null,
  },
};

exports.releases = {
  name: 'mgmt-vers',
  releases: {
    stable: null,
    latest: null,
  },
};

exports.market = {
  hooks: {
    name: 'srs-hooks',
    image: 'registry.cn-hangzhou.aliyuncs.com/ossrs/srs-terraform:hooks-1',
    port: 2021,
    releases: {
      stable: null,
      latest: null,
    },
    container: {
      ID: null,
      State: null,
      Status: null,
    },
  },
};

