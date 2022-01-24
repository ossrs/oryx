'use strict';

const dotenv = require('dotenv');
const fs = require('fs');
const os = require('os');

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

// MySQL日期字段格式化字符串 @see https://stackoverflow.com/a/27381633
exports.MYSQL_DATETIME = 'YYYY-MM-DD HH:mm:ss';

