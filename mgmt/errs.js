'use strict';

exports.status = {
  auth: 401,
  sys: 500,
};

exports.sys = {
  empty: 1000,
  boot: 1001,
};

exports.auth = {
  password: 2000,
  token: 2001,
  init: 2002,
};

exports.srs = {
  verify: 3000,
};

