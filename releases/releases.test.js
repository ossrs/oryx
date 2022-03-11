'use strict';

const releases = require('./releases');

test('normal version', () => {
  expect(releases.filterVersion({queryString: {version: 'v4.0.0'}}).version)
    .toStrictEqual('v4.0.0');
});

test('no version', () => {
  expect(releases.filterVersion().version)
    .toStrictEqual('v0.0.0');
});

