'use strict';

const pkg = require('./package.json');

const stable = 'v1.0.200';
const latest = 'v1.0.220';

// Build the version and docker image url.
function buildVersion(q, version) {
  return {
    stable,
    latest,
    api: `v${pkg.version}`,
  };
}

// Build features query.
function buildFeatures(q, version, res) {
}

// Filter the version from querystring.
function filterVersion(event) {
  let q = event?.queryString || {}

  let version = q.version? q.version :  "v0.0.0"
  if (version.indexOf('v') !== 0) {
    version = "v" + version
  }
  if (version.indexOf('.') === -1) {
    version += ".0.0"
  }

  return {q, version};
}
exports.filterVersion = filterVersion;

// See GetOriginalClientIP of https://github.com/winlinvip/http-gif-sls-writer/blob/master/main.go
function getOriginalClientIP(q, headers, sourceIp) {
  if (q && q.clientip) return q.clientip;

  const fwd = headers && headers['x-forwarded-for'];
  if (fwd) {
    const index = fwd.indexOf(',')
    if (index !== -1) return fwd.substr(0, index);
    return fwd;
  }

  const rip = headers && headers['x-real-ip'];
  if (rip) return rip;

  return sourceIp;
}

// Filter headers.
function filterHeaders(event) {
  event.headers = event.headers || {}
  Object.keys(event.headers).map(e => {
    event.headers[e.toLowerCase()] = event.headers[e];
  });
}

exports.handle = (router) => {
  router.all('/terraform/v1/releases', async (ctx) => {
    // Filter the querystring.
    let {q, version} = filterVersion({
      queryString: ctx.request.query
    });

    // Parse headers to lower case.
    filterHeaders({
      headers: ctx.headers,
    });
    console.log(`api q=${JSON.stringify(q)}, headers=${JSON.stringify(ctx.headers)}`);

    // Build response.
    let res = buildVersion(q, version);
    buildFeatures(q, version, res);

    // See GetOriginalClientIP of https://github.com/winlinvip/http-gif-sls-writer/blob/master/main.go
    q.rip = getOriginalClientIP(q, ctx.headers, ctx.request.ip);
    q.fwd = ctx.headers['x-forwarded-for'];

    // Add the feed back address.
    if (q.feedback) {
      res.addr = {rip: q.rip, fwd: q.fwd};
    }

    console.log(`srs-cloud id=${q.id}, version=${version}, eip=${q.eip}, rip=${q.rip}, fwd=${q.fwd}, res=${JSON.stringify(res)}`);
    ctx.body = res;
  });

  return router;
};

