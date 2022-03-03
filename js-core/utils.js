'use strict';

const errs = require('./errs');

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

const verifyToken = async (jwt, token) => {
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

const streamURL = (vhost, app, stream) => {
  if (vhost === '__defaultVhost__') {
    return `${app}/${stream}`;
  } else {
    return `${vhost}/${app}/${stream}`;
  }
};
exports.streamURL = streamURL;

