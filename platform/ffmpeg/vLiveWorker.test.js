'use strict';

const vLiveWorker = require('./vLiveWorker');

test('build output', () => {
  expect(vLiveWorker.generateOutput('rtmp://server/live', 'livestream'))
    .toStrictEqual('rtmp://server/live/livestream');
  expect(vLiveWorker.generateOutput('rtmp://server/live', 'livestream?secret=abc'))
    .toStrictEqual('rtmp://server/live/livestream?secret=abc');
  expect(vLiveWorker.generateOutput('rtmp://server/live', 'livestream?secret=abc&k=v'))
    .toStrictEqual('rtmp://server/live/livestream?secret=abc&k=v');
});

test('with slash', () => {
  expect(vLiveWorker.generateOutput('rtmp://server/live/', 'livestream'))
    .toStrictEqual('rtmp://server/live/livestream');
  expect(vLiveWorker.generateOutput('rtmp://server/live', '/livestream'))
    .toStrictEqual('rtmp://server/live/livestream');
});

test('empty secret', () => {
  expect(vLiveWorker.generateOutput('rtmp://server/live', ''))
    .toStrictEqual('rtmp://server/live');
  expect(vLiveWorker.generateOutput('rtmp://server/live', null))
    .toStrictEqual('rtmp://server/live');
  expect(vLiveWorker.generateOutput('rtmp://server/live'))
    .toStrictEqual('rtmp://server/live');
});

// Should be one of localhost or mgmt.srs.local
test('localhost', () => {
  expect(vLiveWorker.generateOutput('rtmp://localhost/live', 'livestream').replace(/mgmt.srs.local/g, 'localhost'))
    .toStrictEqual('rtmp://localhost/live/livestream');
  expect(vLiveWorker.generateOutput('rtmp://localhost/live', 'livestream').replace(/localhost/g, 'mgmt.srs.local'))
    .toStrictEqual('rtmp://mgmt.srs.local/live/livestream');
});

