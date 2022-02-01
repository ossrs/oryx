'use strict';

const consts = require('./consts');
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

const createToken = (moment, jwt) => {
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
