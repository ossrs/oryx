import {useNavigate} from "react-router-dom";
import {Container, Tabs, Tab, Accordion} from "react-bootstrap";
import React from "react";
import {Token, Errors} from "./utils";
import axios from "axios";

export default function Dashboard() {
  const navigate = useNavigate();
  const [rtmpServer, setRtmpServer] = React.useState();
  const [rtmpStreamKey, setRtmpStreamKey] = React.useState();
  const [cnConsole, setCnConsole] = React.useState();
  const [flvPlayer, setFlvPlayer] = React.useState();
  const [hlsPlayer, setHlsPlayer] = React.useState();
  const [rtcPlayer, setRtcPlayer] = React.useState();
  const [rtcPublisher, setRtcPublisher] = React.useState();
  const [flvPlayer2, setFlvPlayer2] = React.useState();
  const [hlsPlayer2, setHlsPlayer2] = React.useState();
  const [rtcPlayer2, setRtcPlayer2] = React.useState();
  const [secret, setSecret] = React.useState();

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/mgmt/srs/secret', {
      ...token,
    }).then(res => {
      setSecret(res.data.data);
      console.log(`Status: Query ok, secret=${JSON.stringify(res.data.data)}`);
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/logout');
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

    // Build console url.
    setCnConsole('/console/ng_index.html#/summaries');

    // The player url.
    if (true) {
      const schema = window.location.protocol.replace(':', '');
      const httpPort = window.location.port || (window.location.protocol === 'http:' ? 80 : 443);
      setFlvPlayer(`/players/srs_player.html?schema=${schema}&port=${httpPort}&autostart=true&stream=livestream.flv`);
      setHlsPlayer(`/players/srs_player.html?schema=${schema}&port=${httpPort}&autostart=true&stream=livestream.m3u8`);
      setRtcPlayer(`/players/rtc_player.html?schema=${schema}&port=${httpPort}&api=${httpPort}&autostart=true&stream=livestream`);
    }

    // For WebRTC url.
    if (true) {
      setRtcPublisher(`/players/rtc_publisher.html?schema=https&port=443&api=443&autostart=true&stream=livestream`);
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
          <Tab eventKey="live" title="直播间">
            <Accordion defaultActiveKey="0">
              <Accordion.Item eventKey="0">
                <Accordion.Header>OBS推流</Accordion.Header>
                <Accordion.Body>
                  <div>
                    1. 先在防火墙开启<code>TCP/1935</code>端口
                  </div>
                  <div>
                    2. 然后<a href="https://obsproject.com/download">下载OBS</a>并安装
                  </div>
                  <div>
                    3. 在OBS输入：
                    <ul>
                      <li>推流地址 <code>{rtmpServer}</code></li>
                      <li>推流密钥 <code>{rtmpStreamKey}</code></li>
                    </ul>
                  </div>
                  <div>
                    4. 请选择播放的流：
                    <ul>
                      <li>播放<a href={flvPlayer}>HTTP-FLV流</a></li>
                      <li>播放<a href={hlsPlayer}>HLS流</a></li>
                      <li>播放<a href={rtcPlayer}>WebRTC流</a></li>
                    </ul>
                  </div>
                  <div>
                    5. 点击进入<a id="cnConsole" href={cnConsole}>SRS控制台</a>
                  </div>
                </Accordion.Body>
              </Accordion.Item>
              <Accordion.Item eventKey="1">
                <Accordion.Header>FFmpeg推流</Accordion.Header>
                <Accordion.Body>
                  <div>
                    1. 先在防火墙开启<code>TCP/1935</code>端口
                  </div>
                  <div>
                    2. 然后<a href="https://ffmpeg.org/download.html">下载FFmpeg</a>
                  </div>
                  <div>
                    3. FFmpeg推流命令：<br/>
                    <code>
                      ffmpeg -re -i doc/source.flv -c copy -f flv {rtmpServer}{rtmpStreamKey}
                    </code>
                  </div>
                  <div>
                    4. 请选择播放的流：
                    <ul>
                      <li>播放<a href={flvPlayer}>HTTP-FLV流</a></li>
                      <li>播放<a href={hlsPlayer}>HLS流</a></li>
                      <li>播放<a href={rtcPlayer}>WebRTC流</a></li>
                    </ul>
                  </div>
                  <div>
                    5. 点击进入<a id="cnConsole" href={cnConsole}>SRS控制台</a>
                  </div>
                </Accordion.Body>
              </Accordion.Item>
              <Accordion.Item eventKey="2">
                <Accordion.Header>WebRTC推流</Accordion.Header>
                <Accordion.Body>
                  <div>
                    1. 先在防火墙开启<code>UDP/8000</code>端口
                  </div>
                  <div>
                    2. 请使用<code>https</code>访问管理后台。若使用自签名证书，请点页面空白处然后敲<code>thisisunsafe</code>
                  </div>
                  <div>
                    3. 打开页面推<a href={rtcPublisher}>WebRTC流</a>。注意先停止掉FFmpeg/OBS推流。
                  </div>
                  <div>
                    4. 请选择播放的流：
                    <ul>
                      <li>播放<a href={flvPlayer2}>HTTP-FLV流</a></li>
                      <li>播放<a href={hlsPlayer2}>HLS流</a></li>
                      <li>播放<a href={rtcPlayer2}>WebRTC流</a></li>
                    </ul>
                  </div>
                  <div>
                    5. 点击进入<a id="cnConsole" href={cnConsole}>SRS控制台</a>
                  </div>
                </Accordion.Body>
              </Accordion.Item>
            </Accordion>
          </Tab>
          <Tab eventKey="source" title="源代码">
            <Accordion defaultActiveKey="0">
              <Accordion.Item eventKey="0">
                <Accordion.Header>SRS 4.0</Accordion.Header>
                <Accordion.Body>
                  <div>LightHouse云服务器自带了SRS 4.0的源码，你可以选择：</div>
                  <div>1. 点击下载SRS源码：<a href='/terraform/v1/sources/srs.tar.gz'>下载</a></div>
                  <div>2. 直接在云服务器编译SRS：</div>
                  <div><code>cd ~lighthouse/git/srs</code></div>
                  <div><code>git pull</code></div>
                  <div><code>./configure</code></div>
                  <div><code>make</code></div>
                </Accordion.Body>
              </Accordion.Item>
              <Accordion.Item eventKey="1">
                <Accordion.Header>SRS/GB28181</Accordion.Header>
                <Accordion.Body>
                  <div>以[srs-gb28181](https://github.com/ossrs/srs-gb28181)为例：</div>
                  <div>1. 点击下载SRS源码：<a href='/terraform/v1/sources/srs.tar.gz'>下载</a></div>
                  <div>2. 在本机设置为srs-gb28181的源：</div>
                  <div><code>git remote set-url origin https://github.com/ossrs/srs-gb28181.git</code></div>
                  <div><code>git fetch origin</code></div>
                  <div><code>git checkout -b feature/gb28181 origin/feature/gb28181</code></div>
                  <div>3. 后续就只需要从github增量更新代码就可以：</div>
                  <div><code>git pull</code></div>
                </Accordion.Body>
              </Accordion.Item>
            </Accordion>
          </Tab>
        </Tabs>
      </Container>
    </>
  );
}

