'use strict';

const forwardWorker = require('./forwardWorker');

test('build output', () => {
  expect(forwardWorker.generateOutput('rtmp://localhost/live', 'livestream'))
    .toStrictEqual('rtmp://localhost/live/livestream');
  expect(forwardWorker.generateOutput('rtmp://localhost/live', 'livestream?secret=abc'))
    .toStrictEqual('rtmp://localhost/live/livestream?secret=abc');
  expect(forwardWorker.generateOutput('rtmp://localhost/live', 'livestream?secret=abc&k=v'))
    .toStrictEqual('rtmp://localhost/live/livestream?secret=abc&k=v');
});

test('with slash', () => {
  expect(forwardWorker.generateOutput('rtmp://localhost/live/', 'livestream'))
    .toStrictEqual('rtmp://localhost/live/livestream');
  expect(forwardWorker.generateOutput('rtmp://localhost/live', '/livestream'))
    .toStrictEqual('rtmp://localhost/live/livestream');
});

test('empty secret', () => {
  expect(forwardWorker.generateOutput('rtmp://localhost/live', ''))
    .toStrictEqual('rtmp://localhost/live');
  expect(forwardWorker.generateOutput('rtmp://localhost/live', null))
    .toStrictEqual('rtmp://localhost/live');
  expect(forwardWorker.generateOutput('rtmp://localhost/live'))
    .toStrictEqual('rtmp://localhost/live');
});

