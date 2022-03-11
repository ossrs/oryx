'use strict';

const errs = require('./errs');
const utils = require('./utils');

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
const create = ({config, redis}) => {
  const connect = function () {
    if (!redis) throw utils.asError(errs.sys.empty, errs.status.sys, `redis required`);
    if (!config) throw utils.asError(errs.sys.empty, errs.status.sys, `config required`);
    if (!config.host) throw utils.asError(errs.sys.empty, errs.status.sys, `config.host required`);
    if (!config.port) throw utils.asError(errs.sys.empty, errs.status.sys, `config.port required`);

    const dbConfig = {
      port: parseInt(config.port),
      host: config.host,
      family: 4,
      db: 0,
      password: config.password
    };

    const Redis = redis;
    const client = new Redis(dbConfig);
    return client;
  };

  // Lazy init client util use it.
  let redisClient = null;
  const buildClient = function() {
    if (!redisClient) {
      redisClient = connect();
    }
    return redisClient;
  };

  return {
    del: async function (key) {
      return await buildClient().del(key);
    },
    // @see https://redis.io/commands/set
    set: async function (key, value) {
      return await buildClient().set(key, value);
    },
    get: async function (key) {
      return await buildClient().get(key);
    },
    // @see https://redis.io/commands/hset
    hset: async function (key, field, value) {
      return await buildClient().hset(key, field, value);
    },
    hget: async function (key, field) {
      return await buildClient().hget(key, field);
    },
    hdel: async function (key, field) {
      return await buildClient().hdel(key, field);
    },
    hscan: async function (key, cursor, match, count) {
      return await buildClient().hscan(key, cursor, 'MATCH', match, 'COUNT', count);
    },
    hkeys: async function (key) {
      return await buildClient().hkeys(key);
    },
    hgetall: async function (key) {
      return await buildClient().hgetall(key);
    },
    time: async function () {
      return await buildClient().time();
    },
    hlen: async function (key) {
      return await buildClient().hlen(key);
    },
  };
};
exports.create = create;

