'use strict';

const { createProxyMiddleware } = require('http-proxy-middleware');

console.log('setupProxy for development reactjs');

// See https://create-react-app.dev/docs/proxying-api-requests-in-development/
// See https://create-react-app.dev/docs/proxying-api-requests-in-development/#configuring-the-proxy-manually
module.exports = function(app) {
  const withLogs = (options) => {
    return createProxyMiddleware(options);
  };
  app.use('/api/', withLogs({target: 'http://127.0.0.1:1985/'}));
  app.use('/rtc/', withLogs({target: 'http://127.0.0.1:1985/'}));
  app.use('/*/*.(flv|m3u8|ts|aac|mp3)', withLogs({target: 'http://127.0.0.1:8080/'}));
  app.use('/console/', withLogs({target: 'http://127.0.0.1:8080/'}));
  app.use('/players/', withLogs({target: 'http://127.0.0.1:8080/'}));
};

