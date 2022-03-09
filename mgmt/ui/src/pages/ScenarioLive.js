import {Accordion} from "react-bootstrap";
import React from "react";
import {TutorialsButton, useTutorials} from "../components/TutorialsButton";
import SrsQRCode from "../components/SrsQRCode";

export default function ScenarioLive({urls}) {
  const {flvPlayer, rtmpServer, flvUrl, rtmpStreamKey, hlsPlayer, m3u8Url, rtcPlayer, cnConsole, rtcPublisher, flvPlayer2, flvUrl2, hlsPlayer2, m3u8Url2, rtcPlayer2} = urls;
  const rtmpPublishUrl = `${rtmpServer}${rtmpStreamKey}`;

  const movieTutorials = useTutorials(React.useRef([
    {author: '徐光磊', id: 'BV1RS4y1G7tb'},
    {author: 'SRS', id: 'BV1Nb4y1t7ij'},
    {author: '瓦全', id: 'BV1SF411t7Li'},
    {author: '王大江', id: 'BV16r4y1q7ZT'},
  ]));

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
        <Accordion.Header>OBS/vMix推流</Accordion.Header>
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
                <li>推流地址（服务器） <code>{rtmpServer}</code></li>
                <li>推流密钥（串流密钥） <code>{rtmpStreamKey}</code></li>
              </ul>
            </li>
            <li>
              请选择播放的流：
              <ul>
                <li>播放<a href={flvPlayer} target='_blank' rel='noreferrer'>HTTP-FLV流</a> <code>{flvUrl}</code></li>
                <li>播放<a href={hlsPlayer} target='_blank' rel='noreferrer'>HLS流</a> <code>{m3u8Url}</code></li>
                <li>播放<a href={rtcPlayer} target='_blank' rel='noreferrer'>WebRTC流</a></li>
              </ul>
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
              <code>
                ffmpeg -re -i ~/git/srs/trunk/doc/source.flv -c copy -f flv {rtmpServer}{rtmpStreamKey}
              </code>
            </li>
            <li>推流地址：<br/><code>{rtmpPublishUrl}</code></li>
            <SrsQRCode url={rtmpPublishUrl} />
            <li>
              请选择播放的流：
              <ul>
                <li>播放<a href={flvPlayer} target='_blank' rel='noreferrer'>HTTP-FLV流</a> <code>{flvUrl}</code></li>
                <li>播放<a href={hlsPlayer} target='_blank' rel='noreferrer'>HLS流</a> <code>{m3u8Url}</code></li>
                <li>播放<a href={rtcPlayer} target='_blank' rel='noreferrer'>WebRTC流</a></li>
              </ul>
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
            <li>打开页面推<a href={rtcPublisher} target='_blank' rel='noreferrer'>WebRTC流</a>。注意先停止掉FFmpeg/OBS推流。</li>
            <li>
              请选择播放的流：
              <ul>
                <li>播放<a href={flvPlayer2} target='_blank' rel='noreferrer'>HTTP-FLV流</a> <code>{flvUrl2}</code></li>
                <li>播放<a href={hlsPlayer2} target='_blank' rel='noreferrer'>HLS流</a> <code>{m3u8Url2}</code></li>
                <li>播放<a href={rtcPlayer2} target='_blank' rel='noreferrer'>WebRTC流</a></li>
              </ul>
            </li>
            <li>可选，点击进入<a id="cnConsole" href={cnConsole}>SRS控制台</a>查看流信息</li>
          </ol>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

