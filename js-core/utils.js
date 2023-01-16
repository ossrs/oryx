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
    if (passwd) {
      return await doVerifyToken(jwt, token, passwd);
    }
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
    const token = await redis.hget(keys.redis.SRS_PLATFORM_SECRET, 'token');
    if (token) {
      process.env.SRS_PLATFORM_SECRET = token;
      console.log(`Update api secret to ${process.env.SRS_PLATFORM_SECRET?.length}B`);
    }
  }

  return process.env.SRS_PLATFORM_SECRET;
};
exports.apiSecret = apiSecret;

// Create api secret if not exists.
const setupApiSecret = async (redis, uuidv4, moment) => {
  const r0 = await redis.hget(keys.redis.SRS_PLATFORM_SECRET, 'token');
  if (r0) return [r0, false];

  const token = `srs-v1-${uuidv4().replace(/-/g, '')}`;
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
const removeContainerQuiet = async (execFile, name, disableErrorLog) => {
  try {
    await execFile('docker', ['rm', '-f', name]);
  } catch (e) {
    if (disableErrorLog) return;
    console.log('utils ignore remove container err', e);
  }
};
exports.removeContainerQuiet = removeContainerQuiet;

// Stop container, ignore any error.
const stopContainerQuiet = async (execFile, name, disableErrorLog, time) => {
  try {
    await execFile('docker', [
      'stop',
      ...(time ? ['-t', time] : []),
      name,
    ]);
  } catch (e) {
    if (disableErrorLog) return;
    console.log('utils ignore stop container err', e);
  }
};
exports.stopContainerQuiet = stopContainerQuiet;

function reloadEnv(dotenv, fs, path) {
  // The path ../mgmt is for development, others is for release.
  ['../mgmt', '.', './containers/bin'].map(envDir => {
    if (fs.existsSync(path.join(envDir, '.env'))) {
      dotenv.config({path: path.join(envDir, '.env'), override: true});
    }
  });
}
exports.reloadEnv = reloadEnv;

function saveEnvs(fs, os, dotenv, filename, config) {
  const {parsed: envs} = dotenv.config({path: filename, override: true});

  const merged = {...envs, ...config};

  const envVars = [];
  if (process.env.MGMT_PASSWORD || merged.MGMT_PASSWORD) {
    envVars.push(`MGMT_PASSWORD=${process.env.MGMT_PASSWORD || merged.MGMT_PASSWORD}`);
    envVars.push('');
  }

  envVars.push(
    ...Object.keys(merged).map(k => {
      const v = merged[k];
      if (k === 'MGMT_PASSWORD') return '';
      return v ? `${k}=${v}` : '';
    }).filter(e => e)
  );

  // Append an empty line.
  envVars.push('');

  fs.writeFileSync(filename, envVars.join(os.EOL));
  return merged;
}
exports.saveEnvs = saveEnvs;

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

async function generateDockerArgs(platform, ipv4, conf) {
  const evalValue = (e, defaults) => {
    if (!e) return defaults || '';
    if (typeof(e) === 'function') return e();
    return e;
  };

  const tcpPorts = evalValue(conf.tcpPorts, []).map(e => {
    return e.toString().indexOf(':') > 0 ? ['-p', `${e}/tcp`] : ['-p', `${e}:${e}/tcp`];
  }).flat();
  const udpPorts = evalValue(conf.udpPorts, []).map(e => {
    return e.toString().indexOf(':') > 0 ? ['-p', `${e}/udp`] : ['-p', `${e}:${e}/udp`];
  }).flat();
  const volumes = evalValue(conf.volumes, []).map(e => ['-v', e]).flat();
  const command = evalValue(conf.command, []);
  const extras = evalValue(conf.extras, []);
  const logConfig = evalValue(conf.logConfig, []);

  // The image depends on the registry, which is discovered by platform.
  const image = await conf.image();
  const region = await platform.region();
  const source = await platform.source();

  // Note that it's started by nodejs, so never use '-it'.
  const dockerArgs = [
    'run', '-d', '--restart=always', '--privileged', `--name=${evalValue(conf.name)}`,
    ...(ipv4 ? [`--add-host=mgmt.srs.local:${ipv4.address}`] : []),
    ...tcpPorts,
    ...udpPorts,
    ...logConfig,
    ...volumes,
    ...extras,
    '--env', `SRS_REGION=${region}`,
    '--env', `SRS_SOURCE=${source}`,
    image,
    ...command,
  ];

  return dockerArgs;
}
exports.generateDockerArgs = generateDockerArgs;

// Reload nginx, try to use systemctl, or kill -1 {pid}, or killall -1 nginx.
async function reloadNginx(fs, execFile) {
  if (process.platform === 'darwin') return;

  const nginxServiceExists = fs.existsSync('/usr/lib/systemd/system/nginx.service');
  const nginxPidExists = process.env.NGINX_PID && fs.existsSync(process.env.NGINX_PID);
  if (!nginxServiceExists && !nginxPidExists) {
    throw new Error(`Can't reload NGINX, no service or pid=${process.env.NGINX_PID}`);
  }

  const reloadByService = async () => {
    return await execFile('systemctl', ['reload', 'nginx.service']);
  };
  const reloadByPid = async () => {
    const pid = fs.readFileSync(process.env.NGINX_PID).toString().trim();
    if (pid) return await execFile('kill', ['-s', 'SIGHUP', pid]);
  };

  // Try to reload by service if exists, try pid if failed.
  try {
    if (nginxServiceExists) return await reloadByService();
  } catch (e) {
    if (nginxPidExists) return await reloadByPid();
    throw e;
  }

  // Try to reload by pid if no service.
  if (nginxPidExists) return await reloadByPid();
}
exports.reloadNginx = reloadNginx;

// Copy object, with optional extras fields, for example:
//    copy({id: 0}, ['msg': 'hi'])
// Return an object:
//    {id: 0, msg: 'hi'}
function copy(from, extras) {
  let cp = merge({}, from);

  for (let i = 0; i < extras?.length; i += 2) {
    const k = extras[i];
    const v = extras[i + 1];
    const ov = cp[k];

    const obj = {};
    obj[k] = merge(ov, v);
    cp = merge(cp, obj);
  }
  return cp;
}
exports.copy = copy;

// Merge two object, rewrite dst by src fields.
function merge(dst, src) {
  if (typeof dst !== 'object') return src;
  if (typeof src !== 'object') return src;

  const cp = {};
  for (const k in dst) {
    cp[k] = dst[k];
  }
  for (const k in src) {
    cp[k] = src[k];
  }
  return cp;
}
exports.merge = merge;

