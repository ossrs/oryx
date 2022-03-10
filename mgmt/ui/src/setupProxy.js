'use strict';

const { createProxyMiddleware } = require('http-proxy-middleware');

console.log('setupProxy for development reactjs');

// See https://create-react-app.dev/docs/proxying-api-requests-in-development/
// See https://create-react-app.dev/docs/proxying-api-requests-in-development/#configuring-the-proxy-manually
module.exports = function(app) {
  // Proxy all default mounts to mgmt.
  // Note that should never mount the root, because the mgmt will be overwrite.
  app.use('/console/', createProxyMiddleware({target: 'http://127.0.0.1:2022/'}));
  app.use('/players/', createProxyMiddleware({target: 'http://127.0.0.1:2022/'}));
  app.use('/prometheus/', createProxyMiddleware({target: 'http://127.0.0.1:2022/'}));
  app.use('/terraform/', createProxyMiddleware({target: 'http://127.0.0.1:2022/'}));
  app.use('/tools/', createProxyMiddleware({target: 'http://127.0.0.1:2022/'}));

  // Proxy to SRS API and streaming.
  const withLogs = (options) => {
    return createProxyMiddleware(options);
  };
  app.use('/api/', withLogs({target: 'http://127.0.0.1:1985/'}));
  app.use('/rtc/', withLogs({target: 'http://127.0.0.1:1985/'}));
  app.use('/*/*.(flv|m3u8|ts|aac|mp3)', withLogs({target: 'http://127.0.0.1:8080/'}));
};

