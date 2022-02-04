'use strict';

// HTTP Status code.
exports.status = {
  auth: 401,
  args: 422, // See https://stackoverflow.com/a/42171674/17679565
  not: 404,
  sys: 500,
};

// Error codes for system.
exports.sys = {
  // Param or query is empty but it's required.
  empty: 1000,
  // System is booting, not available.
  boot: 1001,
  // System ssl error.
  ssl: 1002,
  // Resource not found.
  resource: 1003,
  // Arguments is invalid.
  invalid: 1004,
};

// Error codes for auth.
exports.auth = {
  password: 2000,
  token: 2001,
  init: 2002,
};

// Error codes for srs.
exports.srs = {
  // Verify the secret in hooks.
  verify: 3000,
};

