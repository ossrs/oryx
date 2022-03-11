'use strict';

const m3u8Generator = require('./m3u8Generator');

test('build output', () => {
  expect(m3u8Generator.buildVodM3u8({
    files:[{seqno: 0, duration: 1, tsid: 'file0'}],
    bucket:'bucket',
    region:'region',
  }, false)[1]).toStrictEqual([
    '#EXTM3U',
    '#EXT-X-VERSION:3',
    '#EXT-X-ALLOW-CACHE:YES',
    '#EXT-X-PLAYLIST-TYPE:VOD',
    `#EXT-X-TARGETDURATION:1`,
    '#EXT-X-MEDIA-SEQUENCE:0',
    '#EXTINF:1.00, no desc',
    'file0.ts',
    '#EXT-X-ENDLIST',
  ].join('\n'));
});

test('build output2', () => {
  expect(m3u8Generator.buildVodM3u8({
    files:[
      {seqno: 0, duration: 1, tsid: 'file0'},
      {seqno: 1, duration: 2, tsid: 'file1'},
      {seqno: 2, duration: 3, tsid: 'file2'},
    ],
    bucket:'bucket',
    region:'region',
  }, false)[1]).toStrictEqual([
    '#EXTM3U',
    '#EXT-X-VERSION:3',
    '#EXT-X-ALLOW-CACHE:YES',
    '#EXT-X-PLAYLIST-TYPE:VOD',
    `#EXT-X-TARGETDURATION:3`,
    '#EXT-X-MEDIA-SEQUENCE:0',
    '#EXTINF:1.00, no desc',
    'file0.ts',
    '#EXTINF:2.00, no desc',
    'file1.ts',
    '#EXTINF:3.00, no desc',
    'file2.ts',
    '#EXT-X-ENDLIST',
  ].join('\n'));
});

test('build with abs', () => {
  expect(m3u8Generator.buildVodM3u8({
    files:[{seqno: 0, duration: 1, tsid: 'file0', key: 'file0.ts'}],
    bucket:'bucket',
    region:'region',
  }, true)[1]).toStrictEqual([
    '#EXTM3U',
    '#EXT-X-VERSION:3',
    '#EXT-X-ALLOW-CACHE:YES',
    '#EXT-X-PLAYLIST-TYPE:VOD',
    `#EXT-X-TARGETDURATION:1`,
    '#EXT-X-MEDIA-SEQUENCE:0',
    '#EXTINF:1.00, no desc',
    'https://bucket.cos.region.myqcloud.com/file0.ts',
    '#EXT-X-ENDLIST',
  ].join('\n'));
});

test('build with domain', () => {
  expect(m3u8Generator.buildVodM3u8({
    files:[{seqno: 0, duration: 1, tsid: 'file0', key: 'file0.ts'}],
    bucket:'bucket',
    region:'region',
  }, true, 'domain')[1]).toStrictEqual([
    '#EXTM3U',
    '#EXT-X-VERSION:3',
    '#EXT-X-ALLOW-CACHE:YES',
    '#EXT-X-PLAYLIST-TYPE:VOD',
    `#EXT-X-TARGETDURATION:1`,
    '#EXT-X-MEDIA-SEQUENCE:0',
    '#EXTINF:1.00, no desc',
    'https://domain/file0.ts',
    '#EXT-X-ENDLIST',
  ].join('\n'));
});

