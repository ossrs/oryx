import React from "react";

export default function useUrls({secret}) {
  const [rtmpServer, setRtmpServer] = React.useState();
  const [rtmpStreamKey, setRtmpStreamKey] = React.useState();
  const [srtPublishUrl, setSrtPublishUrl] = React.useState();
  const [srtPlayUrl, setSrtPlayUrl] = React.useState();
  const [flvUrl, setFlvUrl] = React.useState();
  const [m3u8Url, setM3u8Url] = React.useState();
  const [cnConsole, setCnConsole] = React.useState();
  const [flvPlayer, setFlvPlayer] = React.useState();
  const [hlsPlayer, setHlsPlayer] = React.useState();
  const [rtcPlayer, setRtcPlayer] = React.useState();
  const [rtcPublisher, setRtcPublisher] = React.useState();
  const [flvPlayer2, setFlvPlayer2] = React.useState();
  const [hlsPlayer2, setHlsPlayer2] = React.useState();
  const [rtcPlayer2, setRtcPlayer2] = React.useState();
  const [flvUrl2, setFlvUrl2] = React.useState();
  const [m3u8Url2, setM3u8Url2] = React.useState();

  React.useEffect(() => {
    // Build RTMP url.
    if (true) {
      setRtmpServer(`rtmp://${window.location.hostname}/live/`);
      setRtmpStreamKey(secret ? `livestream?secret=${secret.publish}` : 'livestream');
    }

    // Build SRT url.
    if (true) {
      const secretQuery = secret ? `,secret=${secret.publish}` : '';
      setSrtPublishUrl(`srt://${window.location.hostname}:10080?streamid=#!::r=live/livestream${secretQuery},m=publish`);
      setSrtPlayUrl(`srt://${window.location.hostname}:10080?streamid=#!::r=live/livestream,latency=20,m=request`);
    }

    // Build console url.
    if (true) {
      const httpPort = window.location.port || (window.location.protocol === 'http:' ? 80 : 443);
      setCnConsole(`/console/ng_index.html#/summaries?port=${httpPort}&http=${httpPort}`);
    }

    // The player url.
    if (true) {
      const schema = window.location.protocol.replace(':', '');
      const httpPort = window.location.port || (window.location.protocol === 'http:' ? 80 : 443);
      setFlvUrl(`${schema}://${window.location.hostname}/live/livestream.flv`);
      setM3u8Url(`${schema}://${window.location.hostname}/live/livestream.m3u8`);
      // /tools/player.html?url=http://localhost:3000/live/livestream.m3u8
      setFlvPlayer(`/tools/player.html?url=${schema}://${window.location.host}/live/livestream.flv`);
      setHlsPlayer(`/tools/player.html?url=${schema}://${window.location.host}/live/livestream.m3u8`);
      setRtcPlayer(`/players/rtc_player.html?schema=${schema}&port=${httpPort}&api=${httpPort}&autostart=true&stream=livestream`);
    }

    // For WebRTC url.
    if (true) {
      const secretQuery = secret ? `&&secret=${secret.publish}` : '';
      setFlvUrl2(`https://${window.location.hostname}/live/livestream.flv`);
      setM3u8Url2(`https://${window.location.hostname}/live/livestream.m3u8`);
      setRtcPublisher(`/players/rtc_publisher.html?schema=https&port=443&api=443&autostart=true&stream=livestream${secretQuery}`);
      setFlvPlayer2(`/players/srs_player.html?schema=https&port=443&api=443&autostart=true&stream=livestream.flv`);
      setHlsPlayer2(`/players/srs_player.html?schema=https&port=443&api=443&autostart=true&stream=livestream.m3u8`);
      setRtcPlayer2(`/players/rtc_player.html?schema=https&port=443&api=443&autostart=true&stream=livestream`);
    }
  }, [secret]);

  return {
    rtmpServer,
    rtmpStreamKey,
    srtPublishUrl,
    srtPlayUrl,
    flvUrl,
    m3u8Url,
    cnConsole,
    flvPlayer,
    hlsPlayer,
    rtcPlayer,
    rtcPublisher,
    flvPlayer2,
    hlsPlayer2,
    rtcPlayer2,
    flvUrl2,
    m3u8Url2,
  };
}

