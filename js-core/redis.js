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
exports.create = create;

