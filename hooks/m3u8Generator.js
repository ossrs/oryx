'use strict';

// See https://github.com/ossrs/srs/wiki/v4_EN_DeliveryHLS#vodm3u8
// See https://developer.apple.com/documentation/http_live_streaming/example_playlists_for_http_live_streaming/video_on_demand_playlist_construction
exports.buildVodM3u8 = (metadataObj, absUrl) => {
  if (!metadataObj) throw new Error('no object');
  if (!metadataObj.files) throw new Error('no files');
  if (!metadataObj.bucket) throw new Error('no bucket');
  if (!metadataObj.region) throw new Error('no region');

  const duration = metadataObj.files.reduce((p, c) => Math.max(p, c.duration || 0), 0);
  const m3u8 = [
    '#EXTM3U',
    '#EXT-X-VERSION:3',
    '#EXT-X-ALLOW-CACHE:YES',
    '#EXT-X-PLAYLIST-TYPE:VOD',
    `#EXT-X-TARGETDURATION:${Math.ceil(duration)}`,
    '#EXT-X-MEDIA-SEQUENCE:0',
    ...metadataObj.files.map((e, index) => {
      const desc = [];

      // TODO: FIXME: Identify discontinuity by callback.
      if (index < metadataObj.files.length - 2) {
        const next = metadataObj.files[index + 1];
        if (e.seqno + 1 !== next.seqno) desc.push('#EXT-X-DISCONTINUITY');
      }

      desc.push(`#EXTINF:${Number(e.duration).toFixed(2)}, no desc`);
      desc.push(absUrl ? `https://${metadataObj.bucket}.cos.${metadataObj.region}.myqcloud.com/${e.key}` : `${e.tsid}.ts`);
      return desc;
    }).flat(),
    '#EXT-X-ENDLIST',
  ];

  return ['application/vnd.apple.mpegurl', m3u8.join('\n'), duration];
};

// See https://github.com/ossrs/srs/wiki/v4_EN_DeliveryHLS#eventm3u8
// See https://developer.apple.com/documentation/http_live_streaming/example_playlists_for_http_live_streaming/event_playlist_construction
exports.buildEventM3u8 = (metadataObj, absUrl) => {
  if (!metadataObj) throw new Error('no object');
  if (!metadataObj.files) throw new Error('no files');
  if (!metadataObj.bucket) throw new Error('no bucket');
  if (!metadataObj.region) throw new Error('no region');

  const duration = metadataObj.files.reduce((p, c) => Math.max(p, c.duration || 0), 0);
  const m3u8 = [
    '#EXTM3U',
    '#EXT-X-VERSION:3',
    '#EXT-X-ALLOW-CACHE:YES',
    '#EXT-X-PLAYLIST-TYPE:EVENT',
    `#EXT-X-TARGETDURATION:${Math.ceil(duration)}`,
    '#EXT-X-MEDIA-SEQUENCE:0',
    ...metadataObj.files.map((e, index) => {
      const desc = [];

      // TODO: FIXME: Identify discontinuity by callback.
      if (index < metadataObj.files.length - 2) {
        const next = metadataObj.files[index + 1];
        if (e.seqno + 1 !== next.seqno) desc.push('#EXT-X-DISCONTINUITY');
      }

      desc.push(`#EXTINF:${Number(e.duration).toFixed(2)}, no desc`);
      desc.push(absUrl ? `https://${metadataObj.bucket}.cos.${metadataObj.region}.myqcloud.com/${e.key}` : `${e.tsid}.ts`);
      return desc;
    }).flat(),
  ];

  return ['application/vnd.apple.mpegurl', m3u8.join('\n'), duration];
};

