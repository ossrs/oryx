'use strict';

const dotenv = require('dotenv');
const fs = require('fs');
const os = require('os');
const moment = require('moment');
const jwt = require('jsonwebtoken');
const errs = require('./errs');

exports.asResponse = (code, data) => {
  return {
    code,
    ... (data? {data} : {}),
  };
};

exports.asError = (code, status, message) => {
  return {
    code,
    status,
    err: new Error(message),
  };
};

exports.loadConfig = () => {
  dotenv.config({path: '.env', override: true});
  return {
    MGMT_PASSWORD: process.env.MGMT_PASSWORD,
  };
};

exports.saveConfig = (config) => {
  const envVars = Object.keys(config).map(k => {
    const v = config[k];
    return `${k}=${v}`;
  });
  fs.writeFileSync('.env', envVars.join(os.EOL));
  return config;
};

exports.createToken = () => {
  const utils = exports;

  // Update the user info, @see https://www.npmjs.com/package/jsonwebtoken#usage
  const expire = moment.duration(1, 'years');
  const createAt = moment.utc().format(utils.MYSQL_DATETIME);
  const expireAt = moment.utc().add(expire).format(utils.MYSQL_DATETIME);
  const token = jwt.sign(
    {v: 1.0, t: createAt, d: expire},
    process.env.MGMT_PASSWORD, {expiresIn: expire.asSeconds()},
  );

  return {expire, expireAt, createAt, token};
};

exports.verifyToken = async (token) => {
  const utils = exports;

  if (!token) throw utils.asError(errs.sys.empty, errs.status.auth, 'no token');
  if (!process.env.MGMT_PASSWORD) throw utils.asError(errs.auth.init, errs.status.auth, 'not init');

  // Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
  return await new Promise((resolve, reject) => {
    jwt.verify(token, process.env.MGMT_PASSWORD, function (err, decoded) {
      if (!err) return resolve(decoded);
      if (err.name === 'TokenExpiredError') throw utils.asError(errs.auth.token, errs.status.auth, `token expired, token=${token}, expiredAt=${err.expiredAt}, ${err.message}`);
      if (err.name === 'JsonWebTokenError') throw utils.asError(errs.auth.token, errs.status.auth, `token invalid, token=${token}, ${err.message}`);
      throw utils.asError(errs.auth.token, errs.status.auth, `token verify failed, ${err.message}`);
    });
  });
};

// MySQL日期字段格式化字符串 @see https://stackoverflow.com/a/27381633
exports.MYSQL_DATETIME = 'YYYY-MM-DD HH:mm:ss';

