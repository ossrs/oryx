//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {Accordion} from "react-bootstrap";
import {TutorialsButton, useTutorials} from "../components/TutorialsButton";
import SrsQRCode from "../components/SrsQRCode";
import * as Icon from 'react-bootstrap-icons';
import {useSrsLanguage} from "../components/LanguageSwitch";
import {Clipboard} from "../utils";
import {useTranslation} from "react-i18next";

export default function ScenarioLive({urls}) {
  const {t} = useTranslation();
  const copyToClipboard = React.useCallback((e, text) => {
    e.preventDefault();

    Clipboard.copy(text).then(() => {
      alert(t('helper.copyOk'));
    }).catch((err) => {
      alert(`${t('helper.copyFail')} ${err}`);
    });
  }, [t]);

  return <ScenarioLiveImpl {...{urls, copyToClipboard}} />;
}

function ScenarioLiveImpl({copyToClipboard, urls}) {
  const language = useSrsLanguage();
  const {t} = useTranslation();
  const {
    flvPlayer, rtmpServer, flvUrl, rtmpStreamKey, hlsPlayer, m3u8Url, rtcUrl, rtcPlayer, cnConsole, enConsole, rtcPublisher,
    srtPublishUrl, srtPlayUrl, rtcPublishUrl, updateStreamName,
  } = urls;

  const rtmpPublishUrl = `${rtmpServer}${rtmpStreamKey}`;
  const xgFlvPlayerUrl = flvPlayer?.replace('player.html', 'xgplayer.html');
  const xgHlsPlayerUrl = hlsPlayer?.replace('player.html', 'xgplayer.html');
  const ffmpegPublishCli = `ffmpeg -re -i ~/git/srs/trunk/doc/source.flv -c copy -f flv ${rtmpPublishUrl}`;
  const ffmpegSrtCli = `ffmpeg -re -i ~/git/srs/trunk/doc/source.flv -c copy -pes_payload_size 0 -f mpegts '${srtPublishUrl}'`;
  const ffplayWindows = `ffplay -fflags nobuffer -flags low_delay -i "${srtPlayUrl}"`;
  const ffplayMac = `ffplay -fflags nobuffer -flags low_delay -i '${srtPlayUrl}'`;

  // Shortcodes of WordPress.
  const flvUrlShortCode = `[srs_player url="${flvUrl}"]`;
  const m3u8UrlShortCode = `[srs_player url="${m3u8Url}"]`;
  const rtcUrlShortCode = `[srs_player url="${rtcUrl}"]`;
  const rtc2UrlShortCode = `[srs_publisher url="${rtcPublishUrl}"]`;

  const movieTutorials = useTutorials({
    bilibili: React.useRef([
      {author: '徐光磊', id: 'BV1RS4y1G7tb'},
      {author: 'SRS', id: 'BV1Nb4y1t7ij'},
      {author: '王大江', id: 'BV16r4y1q7ZT'},
      {author: '周亮', id: 'BV1gT4y1U76d'},
    ])
  });

  return (
    <Accordion defaultActiveKey="1">
      <React.Fragment>
        {
          language === 'zh' ?
            <Accordion.Item eventKey="0">
              <Accordion.Header>场景介绍</Accordion.Header>
              <Accordion.Body>
                <div>
                  推拉直播流<TutorialsButton prefixLine={true} tutorials={movieTutorials} />，公网可以直接使用的直播间，带鉴权只有自己能推流。
                  <p></p>
                </div>
                <p>可应用的具体场景包括：</p>
                <ul>
                  <li>一起看电影，异地恋的情侣，或者三五个好朋友，一起看看自己喜欢的电影</li>
                  <li>远程制作和导播，户外直播用手机或摄像头推流到SRS Stack，用OBS/vMix/芯象制作后再播出，编辑不用在直播现场</li>
                </ul>
                <p>使用说明：</p>
                <ul>
                  <li>推流一般OBS比较好操作，也可以选择FFmpeg或WebRTC</li>
                  <li>播放可以直接复制播放链接，使用Chrome浏览器观看，也可以选择VLC播放流地址</li>
                </ul>
              </Accordion.Body>
            </Accordion.Item> :
            <Accordion.Item eventKey="0">
              <Accordion.Header>Introduction</Accordion.Header>
              <Accordion.Body>
                <div>
                  Private live streaming room, a live streaming room that can be used directly on the public network, with authentication so that only you can push the stream.
                  <p></p>
                </div>
                <p>Specific application scenarios include:</p>
                <ul>
                  <li>Watching movies together, long-distance couples, or a few good friends, watching your favorite movies together</li>
                  <li>Remote production and directing, outdoor live streaming using mobile phones or cameras to push the stream to SRS Stack, then broadcast after production with OBS/vMix, editors don't need to be at the live scene</li>
                </ul>
                <p>Usage instructions:</p>
                <ul>
                  <li>For pushing the stream, OBS is generally easier to operate, but you can also choose FFmpeg or WebRTC</li>
                  <li>For playback, you can directly copy the playback link and watch it using the Chrome browser, or you can choose VLC to play the stream address</li>
                </ul>
              </Accordion.Body>
            </Accordion.Item>
        }
      </React.Fragment>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{t('live.obs.title')}</Accordion.Header>
        <Accordion.Body>
          <div>
            <p style={{display: 'inline-block'}}><strong>{t('live.share.step')}</strong></p>
            {language === 'zh' && <TutorialsButton prefixLine={false} tutorials={movieTutorials}/>}
          </div>
          <ol>
            <li>{t('live.share.fw')} <code>TCP/1935</code></li>
            <li>{t('live.obs.download')} <a href='https://obsproject.com/download' target='_blank' rel='noreferrer'>{t('helper.link')}</a></li>
            <li>
              {t('live.obs.config')}:
              <ul>
                <li>
                  {t('live.obs.server')} <code>{rtmpServer}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtmpServer)} />
                  </div>
                </li>
                <li>
                  {t('live.obs.key')} <code>{rtmpStreamKey}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.switchStream')}>
                    <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
                  </div> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtmpStreamKey)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>
              {t('live.share.title')}
              <ul>
                <li>
                  {t('live.share.flv')}&nbsp;
                  <a href={flvPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>,&nbsp;
                  <a href={xgFlvPlayerUrl} target='_blank' rel='noreferrer'>{t('live.share.xg')}</a>&nbsp;
                  <code>{flvUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrl)} />
                  </div>
                </li>
                <li>
                  {t('live.share.hls')}&nbsp;
                  <a href={hlsPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>,&nbsp;
                  <a href={xgHlsPlayerUrl} target='_blank' rel='noreferrer'>{t('live.share.xg')}</a>&nbsp;
                  <code>{m3u8Url}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8Url)} />
                  </div>
                </li>
                <li>
                  {t('live.share.rtc')}&nbsp;
                  <a href={rtcPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>
                </li>
              </ul>
            </li>
            <li>
              {t('live.share.wp')}&nbsp;
              <ul>
                <li>
                  {t('live.share.wpflv')} &nbsp;
                  <code>{flvUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrlShortCode)} />
                  </div>
                </li>
                <li>
                  {t('live.share.wphls')} &nbsp;
                  <code>{m3u8UrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8UrlShortCode)} />
                  </div>
                </li>
                <li>
                  {t('live.share.wprtc')} &nbsp;
                  <code>{rtcUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtcUrlShortCode)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>
              {t('live.share.console')} &nbsp;
              <a id="console" href={language === 'zh' ? cnConsole : enConsole}>{t('helper.link')}</a>
            </li>
          </ol>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>{t('live.ffmpeg.title')}</Accordion.Header>
        <Accordion.Body>
          <div>
            <p style={{display: 'inline-block'}}><strong>{t('live.share.step')}</strong></p>
            {language === 'zh' && <TutorialsButton prefixLine={false} tutorials={movieTutorials} />}
          </div>
          <ol>
            <li>{t('live.share.fw')} <code>TCP/1935</code></li>
            <li>{t('live.ffmpeg.download')} <a href='https://ffmpeg.org/download.html' target='_blank' rel='noreferrer'>{t('helper.link')}</a></li>
            <li>
              {t('live.ffmpeg.cli')} <br/>
              <code>{ffmpegPublishCli}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
              </div> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, ffmpegPublishCli)} />
              </div>
            </li>
            <li>
              {t('live.ffmpeg.url')} <br/>
              <code>{rtmpPublishUrl}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
              </div> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtmpPublishUrl)} />
              </div>
              <br/>
              {language === 'zh' && <SrsQRCode url={rtmpPublishUrl} />}
            </li>
            <li>
              {t('live.share.title')}
              <ul>
                <li>
                  {t('live.share.flv')}&nbsp;
                  <a href={flvPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>,&nbsp;
                  <a href={xgFlvPlayerUrl} target='_blank' rel='noreferrer'>{t('live.share.xg')}</a>&nbsp;
                  <code>{flvUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrl)} />
                  </div>
                </li>
                <li>
                  {t('live.share.hls')}&nbsp;
                  <a href={hlsPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>,&nbsp;
                  <a href={xgHlsPlayerUrl} target='_blank' rel='noreferrer'>{t('live.share.xg')}</a>&nbsp;
                  <code>{m3u8Url}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8Url)} />
                  </div>
                </li>
                <li>
                  {t('live.share.rtc')}&nbsp;
                  <a href={rtcPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>
                </li>
              </ul>
            </li>
            <li>
              {t('live.share.wp')}&nbsp;
              <ul>
                <li>
                  {t('live.share.wpflv')} &nbsp;
                  <code>{flvUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrlShortCode)} />
                  </div>
                </li>
                <li>
                  {t('live.share.wphls')} &nbsp;
                  <code>{m3u8UrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8UrlShortCode)} />
                  </div>
                </li>
                <li>
                  {t('live.share.wprtc')} &nbsp;
                  <code>{rtcUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtcUrlShortCode)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>
              {t('live.share.console')} &nbsp;
              <a id="console" href={language === 'zh' ? cnConsole : enConsole}>{t('helper.link')}</a>
            </li>
          </ol>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>{t('live.srt.title')}</Accordion.Header>
        <Accordion.Body>
          <div>
            <p style={{display: 'inline-block'}}><strong>{t('live.share.step')}</strong></p>
            {language === 'zh' && <TutorialsButton prefixLine={false} tutorials={movieTutorials} />}
          </div>
          <ol>
            <li>{t('live.share.fw')} <code>UDP/10080</code></li>
            <li>{t('live.obs.download')} <a href='https://obsproject.com/download' target='_blank' rel='noreferrer'>{t('helper.link')}</a></li>
            <li>
              {t('live.obs.config')}. &nbsp;
              {t('helper.see')} <a href='https://github.com/ossrs/srs/issues/1147#lagging-encoder'>{t('helper.link')}</a>:
              <ul>
                <li>{t('live.obs.service')}: <code>{t('live.obs.custom')}</code></li>
                <li>
                  {t('live.obs.server')} <code>{srtPublishUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
                  </div> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPublishUrl)} />
                  </div>
                </li>
                <li>{t('live.obs.key')} <code>{t('live.obs.nokey')}</code></li>
              </ul>
            </li>
            <li>
              {t('live.share.title')}
              <ul>
                <li>
                  {t('live.srt.url')} <br/><code>{srtPlayUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPlayUrl)} />
                  </div>
                </li>
                <li>
                  {t('live.srt.ffplay')} &nbsp;
                  <a href='https://ffmpeg.org/download.html' target='_blank' rel='noreferrer'>{t('helper.link')}</a>
                </li>
                <li>
                  {t('live.srt.win')}<br/>
                  <code>{ffplayWindows}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, ffplayWindows)} />
                  </div>
                </li>
                <li>
                  {t('live.srt.mac')}<br/>
                  <code>{ffplayMac}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, ffplayMac)} />
                  </div>
                </li>
                <li>{t('live.srt.wait')}</li>
                <li>
                  {t('live.srt.wait2')}<br/>
                  <ul>
                    <li>
                      {t('live.share.flv')}&nbsp;
                      <a href={flvPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>,&nbsp;
                      <a href={xgFlvPlayerUrl} target='_blank' rel='noreferrer'>{t('live.share.xg')}</a>&nbsp;
                      <code>{flvUrl}</code> &nbsp;
                      <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                        <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrl)} />
                      </div>
                    </li>
                    <li>
                      {t('live.share.hls')}&nbsp;
                      <a href={hlsPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>,&nbsp;
                      <a href={xgHlsPlayerUrl} target='_blank' rel='noreferrer'>{t('live.share.xg')}</a>&nbsp;
                      <code>{m3u8Url}</code> &nbsp;
                      <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                        <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8Url)} />
                      </div>
                    </li>
                    <li>
                      {t('live.share.rtc')}&nbsp;
                      <a href={rtcPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>
                    </li>
                  </ul>
                </li>
              </ul>
            </li>
            <li>
              {t('live.share.wp')}&nbsp;
              <ul>
                <li>
                  {t('live.share.wpflv')} &nbsp;
                  <code>{flvUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrlShortCode)} />
                  </div>
                </li>
                <li>
                  {t('live.share.wphls')} &nbsp;
                  <code>{m3u8UrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8UrlShortCode)} />
                  </div>
                </li>
                <li>
                  {t('live.share.wprtc')} &nbsp;
                  <code>{rtcUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtcUrlShortCode)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>
              {t('live.share.console')} &nbsp;
              <a id="console" href={language === 'zh' ? cnConsole : enConsole}>{t('helper.link')}</a>
            </li>
          </ol>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="4">
        <Accordion.Header>{t('live.srt.ff')}</Accordion.Header>
        <Accordion.Body>
          <div>
            <p style={{display: 'inline-block'}}><strong>{t('live.share.step')}</strong></p>
            {language === 'zh' && <TutorialsButton prefixLine={false} tutorials={movieTutorials} />}
          </div>
          <ol>
            <li>{t('live.share.fw')} <code>UDP/10080</code></li>
            <li>{t('live.ffmpeg.download')} <a href='https://ffmpeg.org/download.html' target='_blank' rel='noreferrer'>{t('helper.link')}</a></li>
            <li>
              {t('live.ffmpeg.cli')} <br/>
              <code>{ffmpegSrtCli}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
              </div> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, ffmpegSrtCli)} />
              </div>
            </li>
            <li>
              {t('live.ffmpeg.url')} <br/>
              <code>{srtPublishUrl}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
              </div> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPublishUrl)} />
              </div>
              <br/>
              {language === 'zh' && <SrsQRCode url={srtPublishUrl} />}
            </li>
            <li>
              {t('live.share.title')}
              <ul>
                <li>
                  {t('live.srt.url')} <br/><code>{srtPlayUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPlayUrl)} />
                  </div>
                </li>
                <li>
                  {t('live.srt.ffplay')} &nbsp;
                  <a href='https://ffmpeg.org/download.html' target='_blank' rel='noreferrer'>{t('helper.link')}</a>
                </li>
                <li>
                  {t('live.srt.win')}<br/>
                  <code>{ffplayWindows}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, ffplayWindows)} />
                  </div>
                </li>
                <li>
                  {t('live.srt.mac')}<br/>
                  <code>{ffplayMac}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, ffplayMac)} />
                  </div>
                </li>
                <li>{t('live.srt.wait')}</li>
                <li>
                  {t('live.srt.wait2')}<br/>
                  <ul>
                    <li>
                      {t('live.share.flv')}&nbsp;
                      <a href={flvPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>,&nbsp;
                      <a href={xgFlvPlayerUrl} target='_blank' rel='noreferrer'>{t('live.share.xg')}</a>&nbsp;
                      <code>{flvUrl}</code> &nbsp;
                      <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                        <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrl)} />
                      </div>
                    </li>
                    <li>
                      {t('live.share.hls')}&nbsp;
                      <a href={hlsPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>,&nbsp;
                      <a href={xgHlsPlayerUrl} target='_blank' rel='noreferrer'>{t('live.share.xg')}</a>&nbsp;
                      <code>{m3u8Url}</code> &nbsp;
                      <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                        <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8Url)} />
                      </div>
                    </li>
                    <li>
                      {t('live.share.rtc')}&nbsp;
                      <a href={rtcPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>
                    </li>
                  </ul>
                </li>
              </ul>
            </li>
            <li>
              {t('live.share.wp')}&nbsp;
              <ul>
                <li>
                  {t('live.share.wpflv')} &nbsp;
                  <code>{flvUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrlShortCode)} />
                  </div>
                </li>
                <li>
                  {t('live.share.wphls')} &nbsp;
                  <code>{m3u8UrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8UrlShortCode)} />
                  </div>
                </li>
                <li>
                  {t('live.share.wprtc')} &nbsp;
                  <code>{rtcUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtcUrlShortCode)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>
              {t('live.share.console')} &nbsp;
              <a id="console" href={language === 'zh' ? cnConsole : enConsole}>{t('helper.link')}</a>
            </li>
          </ol>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="5">
        <Accordion.Header>{t('live.rtc.title')}</Accordion.Header>
        <Accordion.Body>
          <div>
            <p style={{display: 'inline-block'}}><strong>{t('live.share.step')}</strong></p>
            {language === 'zh' && <TutorialsButton prefixLine={false} tutorials={movieTutorials} />}
          </div>
          <ol>
            <li>{t('live.share.fw')} <code>UDP/8000</code></li>
            {window.location.protocol === 'http:' && <li>{t('live.rtc.https')} <code>thisisunsafe</code></li>}
            <li>
              {t('live.rtc.tip')} <a href={rtcPublisher} target='_blank' rel='noreferrer'>{t('helper.link')}</a> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
              </div>
              <br/>
              <code>{t('live.rtc.tip2')}</code>
            </li>
            <li>
              {t('live.rtc.wp')}
              <ul>
                <li>
                  {t('live.rtc.sc')} &nbsp;
                  <code>{rtc2UrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
                  </div> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtc2UrlShortCode)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>
              {t('live.share.title')}
              <ul>
                <li>
                  {t('live.share.flv')}&nbsp;
                  <a href={flvPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>,&nbsp;
                  <a href={xgFlvPlayerUrl} target='_blank' rel='noreferrer'>{t('live.share.xg')}</a>&nbsp;
                  <code>{flvUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrl)} />
                  </div>
                </li>
                <li>
                  {t('live.share.hls')}&nbsp;
                  <a href={hlsPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>,&nbsp;
                  <a href={xgHlsPlayerUrl} target='_blank' rel='noreferrer'>{t('live.share.xg')}</a>&nbsp;
                  <code>{m3u8Url}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8Url)} />
                  </div>
                </li>
                <li>
                  {t('live.share.rtc')}&nbsp;
                  <a href={rtcPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>
                </li>
              </ul>
            </li>
            <li>
              {t('live.share.wp')}&nbsp;
              <ul>
                <li>
                  {t('live.share.wpflv')} &nbsp;
                  <code>{flvUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrlShortCode)} />
                  </div>
                </li>
                <li>
                  {t('live.share.wphls')} &nbsp;
                  <code>{m3u8UrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8UrlShortCode)} />
                  </div>
                </li>
                <li>
                  {t('live.share.wprtc')} &nbsp;
                  <code>{rtcUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtcUrlShortCode)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>
              {t('live.share.console')} &nbsp;
              <a id="console" href={language === 'zh' ? cnConsole : enConsole}>{t('helper.link')}</a>
            </li>
          </ol>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

