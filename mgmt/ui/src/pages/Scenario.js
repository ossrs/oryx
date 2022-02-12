import {useNavigate} from "react-router-dom";
import {Container, Tabs, Tab, Accordion, Form} from "react-bootstrap";
import React from "react";
import {Token, Errors} from "../utils";
import axios from "axios";
import {TutorialsButton, useTutorials} from '../components/TutorialsButton';

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

  const srtTutorials = useTutorials([
    {author: '崔国栋', id: 'BV1aS4y1G7iG'},
  ]);
  const movieTutorials = useTutorials([
    {author: '徐光磊', id: 'BV1RS4y1G7tb'},
  ]);

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/hooks/srs/secret', {
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
            <Accordion defaultActiveKey="1">
              <Accordion.Item eventKey="0">
                <Accordion.Header>场景介绍</Accordion.Header>
                <Accordion.Body>
                  <div>
                    超清实时直播<TutorialsButton prefixLine={true} tutorials={srtTutorials} />，指码率很高（测试过2~8Mbps），延迟很低（200~500ms）且无累计延迟的直播。
                    <p></p>
                  </div>
                  <p>可应用的具体场景包括：</p>
                  <ul>
                    <li>超高清视频会议，使用专业导播台，把直播流投屏到大屏，注意需要专门的硬件做降噪和回声消除</li>
                    <li>远距离和弱网推流直播，比如跨国推流，注意推RTMP后并不使用SRT播放而是普通直播播放（HTTP-FLV/HLS/WebRTC）</li>
                  </ul>
                  <p>使用说明：</p>
                  <ul>
                    <li>延迟和每个环节都相关，我们在这个后台简化了配置，具体可以参考<a href='https://github.com/ossrs/srs/issues/1147#lagging' target='_blank' rel='noreferrer'>这里</a>。 </li>
                    <li>推荐使用<a href='http://www.sinsam.com/' target='_blank' rel='noreferrer'>芯象直播(Windows)</a>推流，其次是<a href='https://obsproject.com/download' target='_blank' rel='noreferrer'>OBS</a>和vmix</li>
                    <li>推荐使用<a href='https://ffmpeg.org/download.html' target='_blank' rel='noreferrer'>ffplay</a>播放，其次是vmix，<a href='http://www.sinsam.com/' target='_blank' rel='noreferrer'>芯象直播(Windows)</a></li>
                  </ul>
                </Accordion.Body>
              </Accordion.Item>
              <Accordion.Item eventKey="1">
                <Accordion.Header>芯象推流+ffplay播放 ~= 230ms延迟</Accordion.Header>
                <Accordion.Body>
                  <div>
                    <p style={{display: 'inline-block'}}><strong>前提：</strong></p>
                    <TutorialsButton prefixLine={false} tutorials={srtTutorials} />
                  </div>
                  <ol>
                    <li>
                      请检查网络延迟，RTT必须在<code>60ms</code>之内，请执行命令：<br/>
                      <code>ping {window.location.hostname}</code>
                    </li>
                    <li>
                      请关注网络质量，推荐使用网线和专线接入服务器，网络丢包不能超过<code>10%</code>
                    </li>
                    <li>
                      不建议使用WiFi，信号强度很好噪声很低时，偶然RTT和丢包也会比较大，若只能WiFi：<br/>
                      <ul>
                        <li>请先用<a href='https://www.intuitibits.com/products/wifiexplorer/' target='_blank' rel='noreferrer'>WiFi Explorer</a>确认信号良好</li>
                        <li>不要有任何屏障，比如和WiFi不要有水泥墙，不要太远距离</li>
                        <li>不要使用共享的WiFi，避免干扰</li>
                      </ul>
                    </li>
                    <li>
                      请关注电脑的CPU使用率，不能超过<code>80%</code>
                    </li>
                  </ol>
                  <p><strong>推流操作步骤：</strong></p>
                  <ol>
                    <li>先在服务器防火墙开启<code>UDP/10080</code>端口</li>
                    <li>下载<a href='http://www.sinsam.com/' target='_blank' rel='noreferrer'>芯象直播Windows版</a>，注意一定要<code>Windows版</code>，若你是Mac请用其他方案</li>
                    <li>
                      配置芯象推流，可以参考<a href='https://github.com/ossrs/srs/issues/1147#lagging-encoder'>链接</a>：
                      <ol>
                        <li>类型：<code>自定义推流</code></li>
                        <li>推流地址：<br/><code>{srtPublishUrl}</code></li>
                        <li>传输模式：<code>单一网络</code></li>
                        <li>编码方式：<code>软件编码</code></li>
                        <li>配置文件：<code>基线配置</code></li>
                        <li>速率控制：<code>CBR</code></li>
                      </ol>
                    </li>
                    <li>点击推流按钮</li>
                  </ol>
                  <p><strong>播放操作步骤：</strong></p>
                  <ol>
                    <li>SRT流播放地址：<br/><code>{srtPlayUrl}</code></li>
                    <li>下载<a href='https://ffmpeg.org/download.html' target='_blank' rel='noreferrer'>ffplay</a>，FFmpeg自带的低延迟播放器</li>
                    <li>
                      Windows，执行命令：<br/>
                      <code>
                        ffplay -fflags nobuffer -flags low_delay -i "{srtPlayUrl}"
                      </code>
                    </li>
                    <li>
                      Mac或Linux，执行命令：<br/>
                      <code>
                        ffplay -fflags nobuffer -flags low_delay -i '{srtPlayUrl}'
                      </code>
                    </li>
                    <li>SRT流画面出来较慢，请稍安勿躁</li>
                    <li>
                      也可以快速预览其他格式的流，注意延迟比直接播放SRT流会高很多：<br/>
                      <ul>
                        <li>播放<a href={flvPlayer} target='_blank' rel='noreferrer'>HTTP-FLV流</a> <code>{flvUrl}</code></li>
                        <li>播放<a href={hlsPlayer} target='_blank' rel='noreferrer'>HLS流</a> <code>{m3u8Url}</code></li>
                        <li>播放<a href={rtcPlayer} target='_blank' rel='noreferrer'>WebRTC流</a></li>
                      </ul>
                    </li>
                  </ol>
                  <p>若需要测量延迟请参考<a href='https://github.com/ossrs/srs/issues/1147#lagging-benchmark' target='_blank' rel='noreferrer'>这里</a></p>
                </Accordion.Body>
              </Accordion.Item>
              <Accordion.Item eventKey="2">
                <Accordion.Header>OBS推流+ffplay播放 ~= 300ms延迟</Accordion.Header>
                <Accordion.Body>
                  <div>
                    <p style={{display: 'inline-block'}}><strong>前提：</strong></p>
                    <TutorialsButton prefixLine={false} tutorials={srtTutorials} />
                  </div>
                  <ol>
                    <li>
                      请检查网络延迟，RTT必须在<code>60ms</code>之内，请执行命令：<br/>
                      <code>ping {window.location.hostname}</code>
                    </li>
                    <li>
                      请关注网络质量，推荐使用网线和专线接入服务器，网络丢包不能超过<code>10%</code>
                    </li>
                    <li>
                      不建议使用WiFi，信号强度很好噪声很低时，偶然RTT和丢包也会比较大，若只能WiFi：<br/>
                      <ul>
                        <li>请先用<a href='https://www.intuitibits.com/products/wifiexplorer/' target='_blank' rel='noreferrer'>WiFi Explorer</a>确认信号良好</li>
                        <li>不要有任何屏障，比如和WiFi不要有水泥墙，不要太远距离</li>
                        <li>不要使用共享的WiFi，避免干扰</li>
                      </ul>
                    </li>
                    <li>
                      请关注电脑的CPU使用率，不能超过<code>80%</code>
                    </li>
                  </ol>
                  <p><strong>推流操作步骤：</strong></p>
                  <ol>
                    <li>先在服务器防火墙开启<code>UDP/10080</code>端口</li>
                    <li>请从<a href='https://obsproject.com/download' target='_blank' rel='noreferrer'>下载OBS</a>并安装</li>
                    <li>
                      配置OBS推流，可以参考<a href='https://github.com/ossrs/srs/issues/1147#lagging-encoder'>链接</a>：
                      <ol>
                        <li>服务： <code>自定义</code></li>
                        <li>推流地址（服务器）： <br/><code>{srtPublishUrl}</code></li>
                        <li>推流密钥（串流密钥）：<code>无，注意请不要填任何字符串</code></li>
                      </ol>
                    </li>
                    <li>
                      配置OBS的输出，可以参考<a href='https://github.com/ossrs/srs/issues/1147#lagging-encoder'>链接</a>：
                      <ol>
                        <li>输出模式：<code>高级</code></li>
                        <li>编码器：<code>x264</code></li>
                        <li>码率控制：<code>CBR</code></li>
                        <li>关键帧间隔： <code>3</code></li>
                        <li>CPU使用预设：<code>veryfast</code></li>
                        <li>配置（Profile）：<code>baseline</code></li>
                        <li>微调（Tune）： <code>zerolatency</code></li>
                      </ol>
                    </li>
                    <li>点击开始推流</li>
                  </ol>
                  <p><strong>播放操作步骤：</strong></p>
                  <ol>
                    <li>SRT流播放地址：<br/><code>{srtPlayUrl}</code></li>
                    <li>下载<a href='https://ffmpeg.org/download.html' target='_blank' rel='noreferrer'>ffplay</a>，FFmpeg自带的低延迟播放器</li>
                    <li>
                      Windows，执行命令：<br/>
                      <code>
                        ffplay -fflags nobuffer -flags low_delay -i "{srtPlayUrl}"
                      </code>
                    </li>
                    <li>
                      Mac或Linux，执行命令：<br/>
                      <code>
                        ffplay -fflags nobuffer -flags low_delay -i '{srtPlayUrl}'
                      </code>
                    </li>
                    <li>SRT流画面出来较慢，请稍安勿躁</li>
                    <li>
                      也可以快速预览其他格式的流，注意延迟比直接播放SRT流会高很多：<br/>
                      <ul>
                        <li>播放<a href={flvPlayer} target='_blank' rel='noreferrer'>HTTP-FLV流</a> <code>{flvUrl}</code></li>
                        <li>播放<a href={hlsPlayer} target='_blank' rel='noreferrer'>HLS流</a> <code>{m3u8Url}</code></li>
                        <li>播放<a href={rtcPlayer} target='_blank' rel='noreferrer'>WebRTC流</a></li>
                      </ul>
                    </li>
                  </ol>
                  <p>若需要测量延迟请参考<a href='https://github.com/ossrs/srs/issues/1147#lagging-benchmark' target='_blank' rel='noreferrer'>这里</a></p>
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
                  <div><code>cd ~lighthouse/git/srs/trunk</code></div>
                  <div><code>git pull</code></div>
                  <div><code>./configure</code></div>
                  <div><code>make</code></div>
                </Accordion.Body>
              </Accordion.Item>
              <Accordion.Item eventKey="1">
                <Accordion.Header>SRS/GB28181</Accordion.Header>
                <Accordion.Body>
                  <div>以<a href='https://github.com/ossrs/srs-gb28181' target='_blank' rel='noreferrer'>srs-gb28181</a>为例：</div>
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

