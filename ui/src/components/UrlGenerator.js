//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {Token, StreamName} from "../utils";
import axios from "axios";
import {useTranslation} from "react-i18next";
import {useErrorHandler} from "react-error-boundary";

export function buildUrls(defaultUrl, secret) {
  if (!defaultUrl) defaultUrl = `live/livestream`;
  if (defaultUrl.indexOf('://') > 0) {
    const a0 = document.createElement("a");
    a0.href = defaultUrl.replace('rtmp:', 'http:');
    defaultUrl = a0.pathname.substring(1);
  }

  // Parse the default url to vhost, app, and stream.
  // @see: http://stackoverflow.com/questions/10469575/how-to-use-location-object-to-parse-url-without-redirecting-the-page-in-javascri
  const a = document.createElement("a");
  a.href = `${window.location.protocol}//${window.location.host}/${defaultUrl}`;

  const defaultSchema = a.protocol.replace(':', '');
  const defaultHost = a.host;
  const defaultHostname = a.hostname;
  const defaultPort = a.port || (a.protocol === 'http:' ? 80 : 443);
  const defaultApp = a.pathname.substring(1, a.pathname.lastIndexOf("/"));
  const defaultStream = a.pathname.slice(a.pathname.lastIndexOf("/") + 1);

  const urls = {};

  // Build RTMP url.
  if (true) {
    urls.rtmpServer = `rtmp://${defaultHostname}/${defaultApp}/`;
    urls.rtmpStreamKey = secret ? `${defaultStream}?secret=${secret.publish}` : defaultStream;
  }

  // Build SRT url.
  if (true) {
    const secretQuery = secret ? `?secret=${secret.publish}` : '';
    urls.srtPublishUrl = `srt://${defaultHostname}:10080?streamid=#!::r=${defaultApp}/${defaultStream}${secretQuery},m=publish`;
    urls.srtPlayUrl = `srt://${defaultHostname}:10080?streamid=#!::r=${defaultApp}/${defaultStream},latency=20,m=request`;
  }

  // Build console url.
  if (true) {
    const httpPort = defaultPort;
    urls.cnConsole = `/console/ng_index.html#/summaries?port=${httpPort}&http=${httpPort}`;
    urls.enConsole = `/console/en_index.html#/summaries?port=${httpPort}&http=${httpPort}`;
  }

  // The player url.
  if (true) {
    const secretQuery = secret ? `?secret=${secret.publish}` : '';
    const schema = defaultSchema;
    const httpPort = defaultPort;
    urls.flvUrl = `${schema}://${defaultHostname}/${defaultApp}/${defaultStream}.flv`;
    urls.m3u8Url = `${schema}://${defaultHostname}/${defaultApp}/${defaultStream}.m3u8`;
    urls.rtcUrl = `webrtc://${defaultHostname}/${defaultApp}/${defaultStream}`;
    urls.rtcPublishUrl = `webrtc://${defaultHostname}/${defaultApp}/${defaultStream}${secretQuery}`;
    // /tools/player.html?url=http://localhost:3000/live/livestream.m3u8
    urls.flvPlayer = `/tools/player.html?url=${schema}://${defaultHost}/${defaultApp}/${defaultStream}.flv`;
    urls.hlsPlayer = `/tools/player.html?url=${schema}://${defaultHost}/${defaultApp}/${defaultStream}.m3u8`;
    urls.rtcPlayer = `/players/whep.html?schema=${schema}&port=${httpPort}&api=${httpPort}&autostart=true&stream=${defaultStream}`;
  }

  // For WebRTC url.
  if (true) {
    const secretQuery = secret ? `&&secret=${secret.publish}` : '';
    urls.flvUrl2 = `https://${defaultHostname}/${defaultApp}/${defaultStream}.flv`;
    urls.m3u8Url2 = `https://${defaultHostname}/${defaultApp}/${defaultStream}.m3u8`;
    urls.rtcPublisher = `/players/rtc_publisher.html?schema=https&port=443&api=443&autostart=true&stream=${defaultStream}${secretQuery}`;
    urls.flvPlayer2 = `/players/srs_player.html?schema=https&port=443&api=443&autostart=true&stream=${defaultStream}.flv`;
    urls.hlsPlayer2 = `/players/srs_player.html?schema=https&port=443&api=443&autostart=true&stream=${defaultStream}.m3u8`;
    urls.rtcPlayer2 = `/players/whep.html?schema=https&port=443&api=443&autostart=true&stream=${defaultStream}`;
  }

  // For transcode stream and urls.
  if (true) {
    const transcodeStreamName = `${defaultStream}-trans`;
    urls.transcodeStreamName = transcodeStreamName;
    urls.transcodeStreamKey = secret ? `${transcodeStreamName}?secret=${secret.publish}` : transcodeStreamName;
    const schema = defaultSchema;
    urls.transcodeFlvPlayer = `/tools/player.html?url=${schema}://${defaultHost}/${defaultApp}/${transcodeStreamName}.flv`;
  }

  return urls;
}

export default function useUrls() {
  const [rtmpServer, setRtmpServer] = React.useState();
  const [rtmpStreamName, setRtmpStreamName] = React.useState(StreamName.load());
  const [rtmpStreamKey, setRtmpStreamKey] = React.useState();
  const [srtPublishUrl, setSrtPublishUrl] = React.useState();
  const [srtPlayUrl, setSrtPlayUrl] = React.useState();
  const [flvUrl, setFlvUrl] = React.useState();
  const [m3u8Url, setM3u8Url] = React.useState();
  const [rtcUrl, setRtcUrl] = React.useState();
  const [rtcPublishUrl, setRtcPublishUrl] = React.useState();
  const [cnConsole, setCnConsole] = React.useState();
  const [enConsole, setEnConsole] = React.useState();
  const [flvPlayer, setFlvPlayer] = React.useState();
  const [hlsPlayer, setHlsPlayer] = React.useState();
  const [rtcPlayer, setRtcPlayer] = React.useState();
  const [rtcPublisher, setRtcPublisher] = React.useState();
  const [flvPlayer2, setFlvPlayer2] = React.useState();
  const [hlsPlayer2, setHlsPlayer2] = React.useState();
  const [rtcPlayer2, setRtcPlayer2] = React.useState();
  const [flvUrl2, setFlvUrl2] = React.useState();
  const [m3u8Url2, setM3u8Url2] = React.useState();

  const [transcodeStreamName, setTranscodeStreamName] = React.useState();
  const [transcodeStreamKey, setTranscodeStreamKey] = React.useState();
  const [transcodeFlvPlayer, setTranscodeFlvPlayer] = React.useState();

  const [loading, setLoading] = React.useState(true);
  const [secret, setSecret] = React.useState();
  const {t} = useTranslation();
  const handleError = useErrorHandler();

  const updateStreamName = React.useCallback(() => {
    const name = Math.random().toString(16).slice(-6).split('').map(e => {
      return (e >= '0' && e <= '9') ? String.fromCharCode('a'.charCodeAt(0) + (parseInt(Math.random() * 16 + e) % 25)) : e;
    }).join('');
    StreamName.save(name);
    setRtmpStreamName(name);
    alert(t('helper.changeStream'));
  }, [t, setRtmpStreamName]);

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/hooks/srs/secret/query', {
      ...token,
    }).then(res => {
      const secret = res.data.data;
      setSecret(secret);
      console.log(`Status: Query ok, secret=${JSON.stringify(res.data.data)}`);
    }).catch(handleError).finally(() => {
      setLoading(false);
    });
  }, [handleError, setSecret, setLoading]);

  React.useEffect(() => {
    // Ignore if not loaded the secret.
    if (loading) return;

    const urls = buildUrls(`live/${rtmpStreamName}`, secret);

    // Build RTMP url.
    if (true) {
      setRtmpServer(urls.rtmpServer);
      setRtmpStreamKey(urls.rtmpStreamKey);
    }

    // Build SRT url.
    if (true) {
      setSrtPublishUrl(urls.srtPublishUrl);
      setSrtPlayUrl(urls.srtPlayUrl);
    }

    // Build console url.
    if (true) {
      setCnConsole(urls.cnConsole);
      setEnConsole(urls.enConsole);
    }

    // The player url.
    if (true) {
      setFlvUrl(urls.flvUrl);
      setM3u8Url(urls.m3u8Url);
      setRtcUrl(urls.rtcUrl);
      setRtcPublishUrl(urls.rtcPublishUrl);
      // /tools/player.html?url=http://localhost:3000/live/livestream.m3u8
      setFlvPlayer(urls.flvPlayer);
      setHlsPlayer(urls.hlsPlayer);
      setRtcPlayer(urls.rtcPlayer);
    }

    // For WebRTC url.
    if (true) {
      setFlvUrl2(urls.flvUrl2);
      setM3u8Url2(urls.m3u8Url2);
      setRtcPublisher(urls.rtcPublisher);
      setFlvPlayer2(urls.flvPlayer2);
      setHlsPlayer2(urls.hlsPlayer2);
      setRtcPlayer2(urls.rtcPlayer2);
    }

    // For transcode stream and urls.
    if (true) {
      setTranscodeStreamName(urls.transcodeStreamName);
      setTranscodeStreamKey(urls.transcodeStreamKey);
      setTranscodeFlvPlayer(urls.transcodeFlvPlayer);
    }
  }, [loading, secret, rtmpStreamName])

  return {
    // For basic stream.
    rtmpServer,
    rtmpStreamName,
    rtmpStreamKey,
    secret,
    // Basic URLs.
    srtPublishUrl,
    srtPlayUrl,
    flvUrl,
    m3u8Url,
    rtcUrl,
    rtcPublishUrl,
    // Web URLs.
    cnConsole,
    enConsole,
    flvPlayer,
    hlsPlayer,
    rtcPlayer,
    rtcPublisher,
    // HTTPS URLS.
    flvPlayer2,
    hlsPlayer2,
    rtcPlayer2,
    flvUrl2,
    m3u8Url2,
    // For transcode.
    transcodeStreamName,
    transcodeStreamKey,
    transcodeFlvPlayer,
    // For update stream name.
    updateStreamName,
  };
}

