'use strict';

const pkg = require('./package.json');

exports.handle = (ctx) => {
  ctx.body = {
    code: 0,
    data: {
      version: pkg.version
    }
  };
};