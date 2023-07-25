import React from "react";
import {Accordion} from "react-bootstrap";
import {TutorialsButton, useTutorials} from "../components/TutorialsButton";
import SrsQRCode from "../components/SrsQRCode";
import * as Icon from 'react-bootstrap-icons';
import {useSrsLanguage} from "../components/LanguageSwitch";

export default function ScenarioLive(props) {
  const language = useSrsLanguage();
  return language === 'zh' ? <ScenarioLiveCn {...props} /> : <ScenarioLiveEn {...props} />;
}

function ScenarioLiveCn({updateStreamName, copyToClipboard, urls}) {
  const {
    flvPlayer, rtmpServer, flvUrl, rtmpStreamKey, hlsPlayer, m3u8Url, rtcUrl, rtcPlayer, cnConsole, rtcPublisher,
    srtPublishUrl, flvPlayer2, flvUrl2, hlsPlayer2, m3u8Url2, rtcPlayer2, rtcPublishUrl,
  } = urls;
  const rtmpPublishUrl = `${rtmpServer}${rtmpStreamKey}`;
  const xgFlvPlayerUrl = flvPlayer?.replace('player.html', 'xgplayer.html');
  const xgHlsPlayerUrl = hlsPlayer?.replace('player.html', 'xgplayer.html');
  const ffmpegPublishCli = `ffmpeg -re -i ~/git/srs/trunk/doc/source.flv -c copy -f flv ${rtmpPublishUrl}`;
  const ffmpegSrtCli = `ffmpeg -re -i ~/git/srs/trunk/doc/source.flv -c copy -pes_payload_size 0 -f mpegts '${srtPublishUrl}'`;

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
      <Accordion.Item eventKey="0">
        <Accordion.Header>场景介绍</Accordion.Header>
        <Accordion.Body>
          <div>
            私人直播间<TutorialsButton prefixLine={true} tutorials={movieTutorials} />，公网可以直接使用的直播间，带鉴权只有自己能推流。
            <p></p>
          </div>
          <p>可应用的具体场景包括：</p>
          <ul>
            <li>一起看电影，异地恋的情侣，或者三五个好朋友，一起看看自己喜欢的电影</li>
            <li>远程制作和导播，户外直播用手机或摄像头推流到SRS云服务器，用OBS/vMix/芯象制作后再播出，编辑不用在直播现场</li>
          </ul>
          <p>使用说明：</p>
          <ul>
            <li>推流一般OBS比较好操作，也可以选择FFmpeg或WebRTC</li>
            <li>播放可以直接复制播放链接，使用Chrome浏览器观看，也可以选择VLC播放流地址</li>
          </ul>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>OBS或vMix推流</Accordion.Header>
        <Accordion.Body>
          <div>
            <p style={{display: 'inline-block'}}><strong>操作步骤：</strong></p>
            <TutorialsButton prefixLine={false} tutorials={movieTutorials} />
          </div>
          <ol>
            <li>在服务器防火墙开启<code>TCP/1935</code>端口</li>
            <li>请从<a href='https://obsproject.com/download' target='_blank' rel='noreferrer'>下载OBS</a>并安装</li>
            <li>
              在OBS输入：
              <ul>
                <li>
                  推流地址（服务器） <code>{rtmpServer}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtmpServer)} />
                  </div>
                </li>
                <li>
                  推流密钥（串流密钥）<code>{rtmpStreamKey}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='更换流名称'>
                    <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
                  </div> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtmpStreamKey)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>
              请选择播放的流：
              <ul>
                <li>
                  播放HTTP-FLV流, 请选择
                  <a href={flvPlayer} target='_blank' rel='noreferrer'>简易</a>或
                  <a href={xgFlvPlayerUrl} target='_blank' rel='noreferrer'>西瓜</a>播放器&nbsp;
                  <code>{flvUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrl)} />
                  </div>
                </li>
                <li>
                  播放HLS流, 请选择
                  <a href={hlsPlayer} target='_blank' rel='noreferrer'>简易</a>或
                  <a href={xgHlsPlayerUrl} target='_blank' rel='noreferrer'>西瓜</a>播放器&nbsp;
                  <code>{m3u8Url}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8Url)} />
                  </div>
                </li>
                <li>播放<a href={rtcPlayer} target='_blank' rel='noreferrer'>WebRTC流</a></li>
              </ul>
            </li>
            <li>
              你也可以嵌入到WordPress：
              <ul>
                <li>
                  嵌入HTTP-FLV流 &nbsp;
                  <code>{flvUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrlShortCode)} />
                  </div>
                </li>
                <li>
                  嵌入HLS流 &nbsp;
                  <code>{m3u8UrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8UrlShortCode)} />
                  </div>
                </li>
                <li>
                  嵌入WebRTC流 &nbsp;
                  <code>{rtcUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtcUrlShortCode)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>
              可选，SRT推流地址（服务器）：<br/>
              <code>{srtPublishUrl}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='更换流名称'>
                <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
              </div> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPublishUrl)} />
              </div>
            </li>
            <li>可选，点击进入<a id="cnConsole" href={cnConsole}>SRS控制台</a>查看流信息</li>
          </ol>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>FFmpeg/芯象推流</Accordion.Header>
        <Accordion.Body>
          <div>
            <p style={{display: 'inline-block'}}><strong>操作步骤：</strong></p>
            <TutorialsButton prefixLine={false} tutorials={movieTutorials} />
          </div>
          <ol>
            <li>先在防火墙开启<code>TCP/1935</code>端口</li>
            <li>请<a href='https://ffmpeg.org/download.html' target='_blank' rel='noreferrer'>下载FFmpeg</a>工具</li>
            <li>
              FFmpeg推流命令：<br/>
              <code>{ffmpegPublishCli}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='更换流名称'>
                <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
              </div> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, ffmpegPublishCli)} />
              </div>
            </li>
            <li>
              推流地址：<br/>
              <code>{rtmpPublishUrl}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='更换流名称'>
                <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
              </div> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtmpPublishUrl)} />
              </div>
              <br/>
              <SrsQRCode url={rtmpPublishUrl} />
            </li>
            <li>
              请选择播放的流：
              <ul>
                <li>
                  播放HTTP-FLV流, 请选择
                  <a href={flvPlayer} target='_blank' rel='noreferrer'>简易</a>或
                  <a href={xgFlvPlayerUrl} target='_blank' rel='noreferrer'>西瓜</a>播放器&nbsp;
                  <code>{flvUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrl)} />
                  </div>
                </li>
                <li>
                  播放HLS流, 请选择
                  <a href={hlsPlayer} target='_blank' rel='noreferrer'>简易</a>或
                  <a href={xgHlsPlayerUrl} target='_blank' rel='noreferrer'>西瓜</a>播放器&nbsp;
                  <code>{m3u8Url}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8Url)} />
                  </div>
                </li>
                <li>播放<a href={rtcPlayer} target='_blank' rel='noreferrer'>WebRTC流</a></li>
              </ul>
            </li>
            <li>
              你也可以嵌入到WordPress：
              <ul>
                <li>
                  嵌入HTTP-FLV流 &nbsp;
                  <code>{flvUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrlShortCode)} />
                  </div>
                </li>
                <li>
                  嵌入HLS流 &nbsp;
                  <code>{m3u8UrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8UrlShortCode)} />
                  </div>
                </li>
                <li>
                  嵌入WebRTC流 &nbsp;
                  <code>{rtcUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtcUrlShortCode)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>
              可选，SRT推流：<br/>
              <code>{ffmpegSrtCli}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='更换流名称'>
                <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
              </div> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, ffmpegSrtCli)} />
              </div>
            </li>
            <li>可选，点击进入<a id="cnConsole" href={cnConsole}>SRS控制台</a>查看流信息</li>
          </ol>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>WebRTC推流</Accordion.Header>
        <Accordion.Body>
          <div>
            <p style={{display: 'inline-block'}}><strong>操作步骤：</strong></p>
            <TutorialsButton prefixLine={false} tutorials={movieTutorials} />
          </div>
          <ol>
            <li>先在服务器防火墙开启<code>UDP/8000</code>端口</li>
            <li>请使用<code>https</code>访问管理后台。若使用自签名证书，请点页面空白处然后敲<code>thisisunsafe</code></li>
            <li>
              打开页面推<a href={rtcPublisher} target='_blank' rel='noreferrer'>WebRTC流</a>。 &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='更换流名称'>
                <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
              </div>
              <br/>
              <code>注意先停止掉FFmpeg/OBS推流。</code>
            </li>
            <li>
              使用WordPress页面推流：
              <ul>
                <li>
                  简码 &nbsp;
                  <code>{rtc2UrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='更换流名称'>
                    <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
                  </div> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtc2UrlShortCode)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>
              请选择播放的流：
              <ul>
                <li>播放<a href={flvPlayer2} target='_blank' rel='noreferrer'>HTTP-FLV流</a> <code>{flvUrl2}</code></li>
                <li>播放<a href={hlsPlayer2} target='_blank' rel='noreferrer'>HLS流</a> <code>{m3u8Url2}</code></li>
                <li>播放<a href={rtcPlayer2} target='_blank' rel='noreferrer'>WebRTC流</a></li>
              </ul>
            </li>
            <li>
              你也可以嵌入到WordPress：
              <ul>
                <li>
                  嵌入HTTP-FLV流 &nbsp;
                  <code>{flvUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrlShortCode)} />
                  </div>
                </li>
                <li>
                  嵌入HLS流 &nbsp;
                  <code>{m3u8UrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8UrlShortCode)} />
                  </div>
                </li>
                <li>
                  嵌入WebRTC流 &nbsp;
                  <code>{rtcUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtcUrlShortCode)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>可选，点击进入<a id="cnConsole" href={cnConsole}>SRS控制台</a>查看流信息</li>
          </ol>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

function ScenarioLiveEn({updateStreamName, copyToClipboard, urls}) {
  const {
    flvPlayer, rtmpServer, flvUrl, rtmpStreamKey, hlsPlayer, m3u8Url, rtcUrl, rtcPlayer, enConsole, rtcPublisher,
    flvPlayer2, flvUrl2, hlsPlayer2, m3u8Url2, rtcPlayer2, rtcPublishUrl,
  } = urls;
  const rtmpPublishUrl = `${rtmpServer}${rtmpStreamKey}`;
  const ffmpegPublishCli = `ffmpeg -re -i ~/git/srs/trunk/doc/source.flv -c copy -f flv ${rtmpPublishUrl}`;

  // Shortcodes of WordPress.
  const flvUrlShortCode = `[srs_player url="${flvUrl}"]`;
  const m3u8UrlShortCode = `[srs_player url="${m3u8Url}"]`;
  const rtcUrlShortCode = `[srs_player url="${rtcUrl}"]`;
  const rtc2UrlShortCode = `[srs_publisher url="${rtcPublishUrl}"]`;

  return (
    <Accordion defaultActiveKey="1">
      <Accordion.Item eventKey="0">
        <Accordion.Header>Introduction</Accordion.Header>
        <Accordion.Body>
          <p>
            Build a live streaming service.
          </p>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>OBS or vMix</Accordion.Header>
        <Accordion.Body>
          <div>
            <p style={{display: 'inline-block'}}><strong>Usage:</strong></p>
          </div>
          <ol>
            <li>Download OBS from <a href='https://obsproject.com/download' target='_blank' rel='noreferrer'>here</a> and install.</li>
            <li>
              OBS configuration:
              <ul>
                <li>
                  Server: <code>{rtmpServer}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtmpServer)} />
                  </div>
                </li>
                <li>
                  Stream Key: <code>{rtmpStreamKey}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Change Stream'>
                    <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
                  </div> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtmpStreamKey)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>
              Play the stream by:
              <ul>
                <li>
                  HTTP-FLV by <a href={flvPlayer} target='_blank' rel='noreferrer'>H5</a> &nbsp;
                  <code>{flvUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrl)} />
                  </div>
                </li>
                <li>
                  HLS by <a href={hlsPlayer} target='_blank' rel='noreferrer'>H5</a> &nbsp;
                  <code>{m3u8Url}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8Url)} />
                  </div>
                </li>
                <li>WebRTC by <a href={rtcPlayer} target='_blank' rel='noreferrer'>H5</a></li>
              </ul>
            </li>
            <li>
              Embed in WordPress post/page：
              <ul>
                <li>
                  For HTTP-FLV &nbsp;
                  <code>{flvUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrlShortCode)} />
                  </div>
                </li>
                <li>
                  For HLS &nbsp;
                  <code>{m3u8UrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8UrlShortCode)} />
                  </div>
                </li>
                <li>
                  For WebRTC &nbsp;
                  <code>{rtcUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtcUrlShortCode)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>Optional, check by <a href={enConsole} target='_blank' rel='noreferrer'>console</a></li>
          </ol>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>FFmpeg</Accordion.Header>
        <Accordion.Body>
          <div>
            <p style={{display: 'inline-block'}}><strong>Usage:</strong></p>
          </div>
          <ol>
            <li>Download FFmpeg from <a href='https://ffmpeg.org/download.html' target='_blank' rel='noreferrer'>here</a>.</li>
            <li>
              FFmpeg cli: <br/>
              <code>
                ffmpeg -re -i ~/git/srs/trunk/doc/source.flv -c copy -f flv {rtmpPublishUrl}
              </code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='Change Stream'>
                <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
              </div> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='Copy'>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, ffmpegPublishCli)} />
              </div>
            </li>
            <li>
              Stream URL:<br/>
              <code>{rtmpPublishUrl}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='Change Stream'>
                <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
              </div> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtmpPublishUrl)} />
              </div>
            </li>
            <li>
              Play the stream by:
              <ul>
                <li>
                  HTTP-FLV by <a href={flvPlayer} target='_blank' rel='noreferrer'>H5</a> &nbsp;
                  <code>{flvUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrl)} />
                  </div>
                </li>
                <li>
                  HLS by <a href={hlsPlayer} target='_blank' rel='noreferrer'>H5</a> &nbsp;
                  <code>{m3u8Url}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8Url)} />
                  </div>
                </li>
                <li>WebRTC by <a href={rtcPlayer} target='_blank' rel='noreferrer'>H5</a></li>
              </ul>
            </li>
            <li>
              Embed in WordPress post/page：
              <ul>
                <li>
                  For HTTP-FLV &nbsp;
                  <code>{flvUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrlShortCode)} />
                  </div>
                </li>
                <li>
                  For HLS &nbsp;
                  <code>{m3u8UrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8UrlShortCode)} />
                  </div>
                </li>
                <li>
                  For WebRTC &nbsp;
                  <code>{rtcUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtcUrlShortCode)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>Optional, check by <a href={enConsole} target='_blank' rel='noreferrer'>console</a></li>
          </ol>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>WebRTC</Accordion.Header>
        <Accordion.Body>
          <div>
            <p style={{display: 'inline-block'}}><strong>Usage:</strong></p>
          </div>
          <ol>
            <li>Allow <code>UDP/8000</code> by firewall</li>
            <li>Please use <code>https</code>.</li>
            <li>
              Publish by <a href={rtcPublisher} target='_blank' rel='noreferrer'>H5</a>. &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='Change'>
                <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
              </div>
              <br/>
              <code>Please stop FFmpeg/OBS before publishing.</code>
            </li>
            <li>
              Publish by WordPress:
              <ul>
                <li>
                  Shortcode &nbsp;
                  <code>{rtc2UrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Change'>
                    <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
                  </div> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtc2UrlShortCode)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>
              Play stream by:
              <ul>
                <li><a href={flvPlayer2} target='_blank' rel='noreferrer'>HTTP-FLV</a> <code>{flvUrl2}</code></li>
                <li><a href={hlsPlayer2} target='_blank' rel='noreferrer'>HLS</a> <code>{m3u8Url2}</code></li>
                <li><a href={rtcPlayer2} target='_blank' rel='noreferrer'>WebRTC</a></li>
              </ul>
            </li>
            <li>
              Embed in WordPress post/page：
              <ul>
                <li>
                  For HTTP-FLV &nbsp;
                  <code>{flvUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrlShortCode)} />
                  </div>
                </li>
                <li>
                  For HLS &nbsp;
                  <code>{m3u8UrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8UrlShortCode)} />
                  </div>
                </li>
                <li>
                  For WebRTC &nbsp;
                  <code>{rtcUrlShortCode}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtcUrlShortCode)} />
                  </div>
                </li>
              </ul>
            </li>
            <li>Optional, check by <a href={enConsole} target='_blank' rel='noreferrer'>console</a></li>
          </ol>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

