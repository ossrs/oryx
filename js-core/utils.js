'use strict';

const errs = require('./errs');
const keys = require('./keys');

// The string format for MySQL, see https://stackoverflow.com/a/27381633
const MYSQL_DATETIME = 'YYYY-MM-DD HH:mm:ss';
exports.MYSQL_DATETIME = MYSQL_DATETIME;

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

async function doVerifyToken(jwt, token, passwd) {
  const utils = exports;

  if (!token) throw utils.asError(errs.sys.empty, errs.status.auth, 'no token');
  if (!passwd) throw utils.asError(errs.auth.init, errs.status.auth, 'no mgmt password');

  // Verify token first, @see https://www.npmjs.com/package/jsonwebtoken#errors--codes
  return await new Promise((resolve, reject) => {
    jwt.verify(token, passwd, function (err, decoded) {
      if (!err) return resolve(decoded);
      if (err?.name === 'TokenExpiredError') throw utils.asError(errs.auth.token, errs.status.auth, `token expired, token=${token}, expiredAt=${err?.expiredAt}, ${err?.message}`);
      if (err?.name === 'JsonWebTokenError') throw utils.asError(errs.auth.token, errs.status.auth, `token invalid, token=${token}, ${err?.message}`);
      throw utils.asError(errs.auth.token, errs.status.auth, `token verify failed, ${err?.message}`);
    });
  });
};

// Verify the token.
const verifyToken = async (jwt, token, passwd) => {
  try {
    return await doVerifyToken(jwt, token, passwd);
  } catch (e) {
  }

  // Fallback, use MGMT_PASSWORD as passwd to verify it.
  // TODO: FIXME: We could remove it when all clients refreshed.
  return await doVerifyToken(jwt, token, process.env.MGMT_PASSWORD);
};
exports.verifyToken = verifyToken;

// Create token.
function createToken(moment, jwt, passwd) {
  if (!passwd) throw exports.asError(errs.auth.init, errs.status.auth, 'no mgmt password');

  // Update the user info, @see https://www.npmjs.com/package/jsonwebtoken#usage
  const expire = moment.duration(1, 'years');
  const createAt = moment.utc().format(MYSQL_DATETIME);
  const expireAt = moment.utc().add(expire).format(MYSQL_DATETIME);
  const nonce = Math.random().toString(16).slice(-6);
  const token = jwt.sign(
    {v: 1.0, t: createAt, d: expire, n: nonce},
    passwd, {expiresIn: expire.asSeconds()},
  );

  return {expire, expireAt, createAt, nonce, token};
}
exports.createToken = createToken;

// Query the api secret from redis, cache it to env.
const apiSecret = async (redis) => {
  if (!process.env.SRS_PLATFORM_SECRET) {
    process.env.SRS_PLATFORM_SECRET = await redis.hget(keys.redis.SRS_PLATFORM_SECRET, 'token');
    console.log(`Update api secret to ${process.env.SRS_PLATFORM_SECRET?.length}B`);
  }

  return process.env.SRS_PLATFORM_SECRET;
};
exports.apiSecret = apiSecret;

// Create api secret if not exists.
const setupApiSecret = async (redis, uuidv4, moment) => {
  const r0 = await redis.hget(keys.redis.SRS_PLATFORM_SECRET, 'token');
  if (r0) return [r0, false];

  const token = uuidv4();
  const r1 = await redis.hset(keys.redis.SRS_PLATFORM_SECRET, 'token', token);
  const r2 = await redis.hset(keys.redis.SRS_PLATFORM_SECRET, 'update', moment().format());

  console.log(`Platform api secret update, token=${token.length}B, r1=${r1}, r2=${r2}`);
  return [token, true];
};
exports.setupApiSecret = setupApiSecret;

const streamURL = (vhost, app, stream) => {
  if (vhost === '__defaultVhost__') {
    return `${app}/${stream}`;
  } else {
    return `${vhost}/${app}/${stream}`;
  }
};
exports.streamURL = streamURL;

// Remove container, ignore any error.
const removeContainerQuiet = async (execFile, name) => {
  try {
    await execFile('docker', ['rm', '-f', name]);
  } catch (e) {
    console.log('utils ignore remove container err', e);
  }
};
exports.removeContainerQuiet = removeContainerQuiet;

function reloadEnv(dotenv, fs, path) {
  ['.', '..', '../..', '../mgmt', '../../mgmt'].map(envDir => {
    if (fs.existsSync(path.join(envDir, '.env'))) {
      dotenv.config({path: path.join(envDir, '.env'), override: true});
    }
  });
}
exports.reloadEnv = reloadEnv;

function srsProxy(staticCache, app, home, prefix, noCaches, alias) {
  const reactFiles = {};

  app.use(staticCache(home, {
    // Cache for a year for it never changes.
    maxAge: 365 * 24 * 3600,
    // It's important to set to dynamic, because the js might changed.
    dynamic: true,
    // If not set, NOT FOUND.
    alias,
    // The baseUrl to mount.
    prefix,
  }, reactFiles));

  noCaches && noCaches.map(f => {
    if (reactFiles[f]) reactFiles[f].maxAge = 0;
  });
}
exports.srsProxy = srsProxy;

