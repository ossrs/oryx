import {useNavigate} from "react-router-dom";
import {Container, Tabs, Tab, Accordion} from "react-bootstrap";
import React from "react";
import {Token, Errors} from "../utils";
import axios from "axios";
import {TutorialsButton, useTutorials} from '../components/TutorialsButton';
import ScenarioDvr from './ScenarioDvr';
import ScenarioSource from './ScenarioSource';
import ScenarioSrt from './ScenarioSrt';

export default function Dashboard() {
  const navigate = useNavigate();
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
  const [secret, setSecret] = React.useState();

  const movieTutorials = useTutorials(React.useRef([
    {author: '徐光磊', id: 'BV1RS4y1G7tb'},
    {author: 'SRS', id: 'BV1Nb4y1t7ij'},
    {author: '瓦全', id: 'BV1SF411t7Li'},
  ]));

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/hooks/srs/secret/query', {
      ...token,
    }).then(res => {
      setSecret(res.data.data);
      console.log(`Status: Query ok, secret=${JSON.stringify(res.data.data)}`);
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/routers-logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  }, [navigate]);

  React.useEffect(() => {
    // Build RTMP url.
    if (true) {
      setRtmpServer(`rtmp://${window.location.hostname}/live/`);
      setRtmpStreamKey(secret ? `livestream?secret=${secret.publish}` : 'livestream');
    }

    // Build SRT url.
    if (true) {
      const secretQuery = secret ? `?secret=${secret.publish}` : '';
      setSrtPublishUrl(`srt://${window.location.hostname}:10080?streamid=#!::h=live/livestream${secretQuery},m=publish`);
      setSrtPlayUrl(`srt://${window.location.hostname}:10080?streamid=#!::h=live/livestream${secretQuery},m=request&latency=20`);
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
      setFlvPlayer(`/players/srs_player.html?schema=${schema}&port=${httpPort}&autostart=true&stream=livestream.flv`);
      setHlsPlayer(`/players/srs_player.html?schema=${schema}&port=${httpPort}&autostart=true&stream=livestream.m3u8`);
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

  return (
    <>
      <p></p>
      <Container>
        <Tabs defaultActiveKey="live" id="uncontrolled-tab-example" className="mb-3">
          <Tab eventKey="live" title="私人直播间">
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
                    <li>远程制作和导播，户外直播用手机或摄像头推流到SRS云服务器，用OBS/Vmix/芯象制作后再播出，编辑不用在直播现场</li>
                  </ul>
                  <p>使用说明：</p>
                  <ul>
                    <li>推流一般OBS比较好操作，也可以选择FFmpeg或WebRTC</li>
                    <li>播放可以直接复制播放链接，使用Chrome浏览器观看，也可以选择VLC播放流地址</li>
                  </ul>
                </Accordion.Body>
              </Accordion.Item>
              <Accordion.Item eventKey="1">
                <Accordion.Header>OBS推流</Accordion.Header>
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
                <Accordion.Header>FFmpeg推流</Accordion.Header>
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
          </Tab>
          <Tab eventKey="srt" title="超清实时直播">
            <ScenarioSrt urls={{srtPublishUrl, srtPlayUrl, flvPlayer, hlsPlayer, flvUrl, m3u8Url, rtcPlayer}}/>
          </Tab>
          <Tab eventKey="dvr" title="云录制">
            <ScenarioDvr />
          </Tab>
          <Tab eventKey="source" title="源代码">
            <ScenarioSource />
          </Tab>
        </Tabs>
      </Container>
    </>
  );
}

