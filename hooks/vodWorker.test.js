'use strict';

const vodWorker = require('./vodWorker');

test('build vod key', () => {
  expect(vodWorker.buildVodTsKey('a/b/f.m3u8', 'file')).toStrictEqual('a/b/file.ts');
  expect(vodWorker.buildVodTsKey('/a/b/f.m3u8', 'file')).toStrictEqual('/a/b/file.ts');
});

