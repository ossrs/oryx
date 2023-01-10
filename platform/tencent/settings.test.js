'use strict';

const settings = require('./settings');

test('test consts', () => {
  expect(settings.GetUserAppId).toStrictEqual('GetUserAppId');
});

