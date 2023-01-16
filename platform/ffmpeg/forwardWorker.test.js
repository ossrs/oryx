'use strict';

const forwardWorker = require('./forwardWorker');

test('build output', () => {
  expect(forwardWorker.generateOutput('rtmp://server/live', 'livestream'))
    .toStrictEqual('rtmp://server/live/livestream');
  expect(forwardWorker.generateOutput('rtmp://server/live', 'livestream?secret=abc'))
    .toStrictEqual('rtmp://server/live/livestream?secret=abc');
  expect(forwardWorker.generateOutput('rtmp://server/live', 'livestream?secret=abc&k=v'))
    .toStrictEqual('rtmp://server/live/livestream?secret=abc&k=v');
});

test('with slash', () => {
  expect(forwardWorker.generateOutput('rtmp://server/live/', 'livestream'))
    .toStrictEqual('rtmp://server/live/livestream');
  expect(forwardWorker.generateOutput('rtmp://server/live', '/livestream'))
    .toStrictEqual('rtmp://server/live/livestream');
});

test('empty secret', () => {
  expect(forwardWorker.generateOutput('rtmp://server/live', ''))
    .toStrictEqual('rtmp://server/live');
  expect(forwardWorker.generateOutput('rtmp://server/live', null))
    .toStrictEqual('rtmp://server/live');
  expect(forwardWorker.generateOutput('rtmp://server/live'))
    .toStrictEqual('rtmp://server/live');
});

// Should be one of localhost or mgmt.srs.local
test('localhost', () => {
  expect(forwardWorker.generateOutput('rtmp://localhost/live', 'livestream').replace(/mgmt.srs.local/g, 'localhost'))
    .toStrictEqual('rtmp://localhost/live/livestream');
  expect(forwardWorker.generateOutput('rtmp://localhost/live', 'livestream').replace(/localhost/g, 'mgmt.srs.local'))
    .toStrictEqual('rtmp://mgmt.srs.local/live/livestream');
});

