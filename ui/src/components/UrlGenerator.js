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
import {SrsEnvContext} from "./SrsEnvContext";

export function buildUrls(defaultUrl, secret, env) {
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
  const defaultHostname = a.hostname;
  const defaultPort = a.port || (a.protocol === 'http:' ? 80 : 443);
  const defaultApp = a.pathname.substring(1, a.pathname.lastIndexOf("/"));
  const defaultStream = a.pathname.slice(a.pathname.lastIndexOf("/") + 1);

  const urls = {};

  // Build RTMP url.
  if (true) {
    const rtmpPort = env.rtmpPort ? `:${env.rtmpPort}` : '';
    urls.rtmpServer = `rtmp://${defaultHostname}${rtmpPort}/${defaultApp}/`;
    urls.rtmpStreamKey = secret ? `${defaultStream}?secret=${secret.publish}` : defaultStream;
  }

  // Build SRT url.
  if (true) {
    const secretQuery = secret ? `?secret=${secret.publish}` : '';
    const srtPort = env.srtPort ? `:${env.srtPort}` : '';
    urls.srtPublishUrl = `srt://${defaultHostname}${srtPort}?streamid=#!::r=${defaultApp}/${defaultStream}${secretQuery},m=publish`;
    urls.srtPlayUrl = `srt://${defaultHostname}${srtPort}?streamid=#!::r=${defaultApp}/${defaultStream},latency=20,m=request`;
  }

  // Build console url.
  if (true) {
    const httpPort = env.httpPort ? env.httpPort : defaultPort;
    urls.cnConsole = `/console/ng_index.html#/summaries?port=${httpPort}&http=${httpPort}`;
    urls.enConsole = `/console/en_index.html#/summaries?port=${httpPort}&http=${httpPort}`;
  }

  // The player url.
  if (true) {
    const secretQuery = secret ? `?secret=${secret.publish}` : '';
    const schema = defaultSchema;
    const httpPort = env.httpPort ? env.httpPort : defaultPort;
    const httpUrlPort = `:${httpPort}`;
    urls.flvUrl = `${schema}://${defaultHostname}${httpUrlPort}/${defaultApp}/${defaultStream}.flv`;
    urls.m3u8Url = `${schema}://${defaultHostname}${httpUrlPort}/${defaultApp}/${defaultStream}.m3u8`;
    urls.rtcUrl = `webrtc://${defaultHostname}${httpUrlPort}/${defaultApp}/${defaultStream}`;
    urls.rtcPublishUrl = `webrtc://${defaultHostname}${httpUrlPort}/${defaultApp}/${defaultStream}${secretQuery}`;
    // /tools/player.html?url=http://localhost:3000/live/livestream.m3u8
    urls.flvPlayer = `/tools/player.html?url=${schema}://${defaultHostname}${httpUrlPort}/${defaultApp}/${defaultStream}.flv`;
    urls.hlsPlayer = `/tools/player.html?url=${schema}://${defaultHostname}${httpUrlPort}/${defaultApp}/${defaultStream}.m3u8`;
    urls.rtcPlayer = `/players/whep.html?schema=${schema}&port=${httpPort}&api=${httpPort}&autostart=true&stream=${defaultStream}`;
  }

  // For WebRTC url.
  if (true) {
    const secretQuery = secret ? `&secret=${secret.publish}` : '';
    const httpPort = env.httpPort ? env.httpPort : defaultPort;
    urls.rtcPublisher = `/players/whip.html?schema=https&port=${httpPort}&api=${httpPort}&autostart=true&stream=${defaultStream}${secretQuery}`;
    urls.whipUrl = `${defaultSchema}://${defaultHostname}:${httpPort}/rtc/v1/whip/?app=${defaultApp}&stream=${defaultStream}${secretQuery}`;
    urls.whepUrl = `${defaultSchema}://${defaultHostname}:${httpPort}/rtc/v1/whep/?app=${defaultApp}&stream=${defaultStream}`;
  }

  // For transcode stream and urls.
  if (true) {
    const transcodeStreamName = `${defaultStream}-trans`;
    urls.transcodeStreamName = transcodeStreamName;
    urls.transcodeStreamKey = secret ? `${transcodeStreamName}?secret=${secret.publish}` : transcodeStreamName;
    const schema = defaultSchema;
    const httpPort = env.httpPort ? env.httpPort : defaultPort;
    const httpUrlPort = `:${httpPort}`;
    urls.transcodeFlvPlayer = `/tools/player.html?url=${schema}://${defaultHostname}${httpUrlPort}/${defaultApp}/${transcodeStreamName}.flv`;
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
  const [whipUrl, setWhipUrl] = React.useState();
  const [whepUrl, setWhepUrl] = React.useState();

  const [transcodeStreamName, setTranscodeStreamName] = React.useState();
  const [transcodeStreamKey, setTranscodeStreamKey] = React.useState();
  const [transcodeFlvPlayer, setTranscodeFlvPlayer] = React.useState();

  const [loading, setLoading] = React.useState(true);
  const [secret, setSecret] = React.useState();
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const env = React.useContext(SrsEnvContext)[0];

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

    const urls = buildUrls(`live/${rtmpStreamName}`, secret, env);

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
      setRtcPublisher(urls.rtcPublisher);
      setWhipUrl(urls.whipUrl);
      setWhepUrl(urls.whepUrl);
    }

    // For transcode stream and urls.
    if (true) {
      setTranscodeStreamName(urls.transcodeStreamName);
      setTranscodeStreamKey(urls.transcodeStreamKey);
      setTranscodeFlvPlayer(urls.transcodeFlvPlayer);
    }
  }, [loading, secret, rtmpStreamName, env])

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
    whipUrl,
    whepUrl,
    // For transcode.
    transcodeStreamName,
    transcodeStreamKey,
    transcodeFlvPlayer,
    // For update stream name.
    updateStreamName,
  };
}

