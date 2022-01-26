'use strict';

const dotenv = require('dotenv');
const fs = require('fs');
const os = require('os');
const moment = require('moment');
const jwt = require('jsonwebtoken');
const errs = require('./errs');
const consts = require('./consts');

const asResponse = (code, data) => {
  return {
    code,
    ... (data? {data} : {}),
  };
};
exports.asResponse = asResponse;

const asError = (code, status, message) => {
  return {
    code,
    status,
    err: new Error(message),
  };
};
exports.asError = asError;

const loadConfig = () => {
  dotenv.config({path: '.env', override: true});
  return {
    MGMT_PASSWORD: process.env.MGMT_PASSWORD,
  };
};
exports.loadConfig = loadConfig;

const saveConfig = (config) => {
  const envVars = Object.keys(config).map(k => {
    const v = config[k];
    return `${k}=${v}`;
  });

  // Append an empty line.
  envVars.push('');

  fs.writeFileSync('.env', envVars.join(os.EOL));
  return config;
};
exports.saveConfig = saveConfig;

const createToken = () => {
  // Update the user info, @see https://www.npmjs.com/package/jsonwebtoken#usage
  const expire = moment.duration(1, 'years');
  const createAt = moment.utc().format(consts.MYSQL_DATETIME);
  const expireAt = moment.utc().add(expire).format(consts.MYSQL_DATETIME);
  const token = jwt.sign(
    {v: 1.0, t: createAt, d: expire},
    process.env.MGMT_PASSWORD, {expiresIn: expire.asSeconds()},
  );

  return {expire, expireAt, createAt, token};
};
exports.createToken = createToken;

const verifyToken = async (token) => {
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
exports.verifyToken = verifyToken;

/*
The config SHOULD be config:Object for Redis db, with bellow fields:
    {host, port, password}
well, the password is optional, which might for redis without password.

For example:
    const config = {
        redis: {
            host: process.env.REDIS_HOST,
            port: process.env.REDIS_PORT,
            password: process.env.REDIS_PASSWORD,
        }
    };
    const ioredis = require('ioredis');
    const redis = require('utils').redis({config: config.redis, redis: ioredis});

    const r0 = await redis.set('KEY', 'VALUE');
 */
const redis = ({config, redis}) => {
  const connect = function () {
    if (!redis) throw asError(errs.sys.empty, errs.status.sys, `redis required`);
    if (!config) throw asError(errs.sys.empty, errs.status.sys, `config required`);
    if (!config.host) throw asError(errs.sys.empty, errs.status.sys, `config.host required`);
    if (!config.port) throw asError(errs.sys.empty, errs.status.sys, `config.port required`);

    const dbConfig = {
      port: config.port,
      host: config.host,
      family: 4,
      db: 0,
      password: config.password
    };

    const Redis = redis;
    const client = new Redis(dbConfig);
    return client;
  };

  const client = connect();
  return {
    del: async function (key) {
      return await client.del(key);
    },
    // @see https://redis.io/commands/set
    set: async function (key, value) {
      return await client.set(key, value);
    },
    get: async function (key) {
      return await client.get(key);
    },
    // @see https://redis.io/commands/hset
    hset: async function (key, field, value) {
      return await client.hset(key, field, value);
    },
    hget: async function (key, field) {
      return await client.hget(key, field);
    },
    hdel: async function (key, field) {
      return await client.hdel(key, field);
    },
    hscan: async function (key, cursor, match, count) {
      return await client.hscan(key, cursor, 'MATCH', match, 'COUNT', count);
    },
    hkeys: async function (key) {
      return await client.hkeys(key);
    },
    hgetall: async function (key) {
      return await client.hgetall(key);
    },
    time: async function () {
      return await client.time();
    },
    hlen: async function (key) {
      return await client.hlen(key);
    },
  };
};
exports.redis = redis;

