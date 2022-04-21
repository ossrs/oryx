import {Accordion} from "react-bootstrap";
import React from "react";
import {TutorialsButton, useTutorials} from "../components/TutorialsButton";
import SrsQRCode from "../components/SrsQRCode";
import * as Icon from "react-bootstrap-icons";
import {useSrsLanguage} from "../components/LanguageSwitch";

export default function ScenarioSrt(props) {
  const language = useSrsLanguage();
  return language === 'zh' ? <ScenarioSrtCn {...props} /> : <ScenarioSrtEn {...props} />;
}

function ScenarioSrtCn({updateStreamName, copyToClipboard, urls}) {
  const {srtPublishUrl, srtPlayUrl, flvPlayer, hlsPlayer, flvUrl, m3u8Url, rtcPlayer} = urls;
  const [hostname, setHostname] = React.useState();
  const [srtPort, setSrtPort] = React.useState();
  const [srtPublishStreamId, setPublishStreamId] = React.useState();
  const [srtPlayStreamId, setPlayStreamId] = React.useState();
  const ffplayWindows = `ffplay -fflags nobuffer -flags low_delay -i "${srtPlayUrl}"`;
  const ffplayMac = `ffplay -fflags nobuffer -flags low_delay -i '${srtPlayUrl}'`;

  const srtTutorials = useTutorials({
    bilibili: React.useRef([
      {author: '崔国栋', id: 'BV1aS4y1G7iG'},
      {author: '马景瑞', id: 'BV1c341177e7'},
      {author: 'SRS', id: 'BV1Nb4y1t7ij'},
      {author: '瓦全', id: 'BV1SF411t7Li'},
      {author: '王大江', id: 'BV16r4y1q7ZT'},
    ])
  });

  React.useEffect(() => {
    if (!srtPublishUrl) return;
    const u = new URL(srtPublishUrl.replace('srt://', 'http://'));
    setHostname(u.hostname);
    setSrtPort(u.port);
    setPublishStreamId(`${u.hash}`);
  }, [srtPublishUrl]);

  React.useEffect(() => {
    if (!srtPublishUrl || !srtPlayUrl) return;
    const u = new URL(srtPlayUrl.replace('srt://', 'http://'));
    setPlayStreamId(`${u.hash.split('&')[0]}`);
  }, [srtPublishUrl, srtPlayUrl]);

  return (
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
            <li>远程制作和导播，户外直播用手机或摄像头推流到SRS云服务器，用OBS/vMix/芯象制作后再播出，编辑不用在直播现场</li>
          </ul>
          <p>使用说明：</p>
          <ul>
            <li>延迟和每个环节都相关，我们在这个后台简化了配置，具体可以参考<a href='https://github.com/ossrs/srs/issues/1147#lagging' target='_blank' rel='noreferrer'>这里</a>。 </li>
            <li>推荐使用<a href='http://www.sinsam.com/' target='_blank' rel='noreferrer'>芯象直播(Windows)</a>推流，其次是<a href='https://obsproject.com/download' target='_blank' rel='noreferrer'>OBS</a>和vMix</li>
            <li>推荐使用<a href='https://ffmpeg.org/download.html' target='_blank' rel='noreferrer'>ffplay</a>播放，其次是vMix，<a href='http://www.sinsam.com/' target='_blank' rel='noreferrer'>芯象直播(Windows)</a></li>
          </ul>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>vMix推拉流 ~= 300ms延迟</Accordion.Header>
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
            <li>下载<a href='https://www.vmix.com/software/download.aspx' target='_blank' rel='noreferrer'>vMix</a>，只有Windows版本</li>
            <li>
              配置vMix推流，可以参考<a href='https://github.com/ossrs/srs/issues/1147#lagging-encoder'>链接</a>：
              <ol>
                <li>
                  点右上角 <code>Settings(设置)</code> => <code>Outputs/NDI/SRT(输出/SRT)</code> / 点<code>Output(1)</code>的设置 <br/>
                  或者点下方的 <code>External(外部)</code> => <code>Outputs/NDI/SRT Settings</code> / 点<code>Output(1)</code>的设置
                </li>
                <li>勾选上：<code>Enable SRT</code></li>
                <li>Type：<code>Caller</code></li>
                <li>
                  Hostname：<code>{hostname}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, hostname)} />
                  </div>
                </li>
                <li>
                  Port：<code>{srtPort}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPort)} />
                  </div>
                </li>
                <li>
                  Latency：<code>20</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, '20')} />
                  </div>
                </li>
                <li>
                  Stream ID：<code>{srtPublishStreamId}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='更换流名称'>
                    <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
                  </div> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPublishStreamId)} />
                  </div>
                </li>
              </ol>
            </li>
            <li>点OK就开始推流</li>
          </ol>
          <p><strong>拉流操作步骤：</strong></p>
          <ol>
            <li>下载<a href='https://www.vmix.com/software/download.aspx' target='_blank' rel='noreferrer'>vMix</a>，只有Windows版本</li>
            <li>
              配置vMix拉流：
              <ol>
                <li>点下方的 <code>Add Input(添加输入)</code> => <code>Stream/SRT(流/SRT)</code></li>
                <li>Stream Type(码流类型)：<code>SRT Caller</code></li>
                <li>
                  Hostname：<code>{hostname}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, hostname)} />
                  </div>
                </li>
                <li>
                  Port：<code>{srtPort}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPort)} />
                  </div>
                </li>
                <li>
                  Latency(延迟)：<code>20</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, '20')} />
                  </div>
                </li>
                <li>
                  Stream ID(流ID)：<code>{srtPlayStreamId}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPlayStreamId)} />
                  </div>
                </li>
                <li>注意：若无法播放，请取消勾选<code>Use Hardware Decoder(使用硬件解码器)</code></li>
              </ol>
            </li>
            <li>点OK就开始拉流</li>
            <li>
              也可以快速预览其他格式的流，注意延迟比直接播放SRT流会高很多：<br/>
              <ul>
                <li>可以用FFplay播放，参考下面的<code>ffplay播放</code>部分</li>
                <li>
                  播放<a href={flvPlayer} target='_blank' rel='noreferrer'>HTTP-FLV流</a> <code>{flvUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrl)} />
                  </div>
                </li>
                <li>
                  播放<a href={hlsPlayer} target='_blank' rel='noreferrer'>HLS流</a> <code>{m3u8Url}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8Url)} />
                  </div>
                </li>
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
                <li>服务：<code>自定义</code></li>
                <li>
                  推流地址（服务器）：<br/><code>{srtPublishUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='更换流名称'>
                    <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
                  </div> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPublishUrl)} />
                  </div>
                </li>
                <li>推流密钥（串流密钥）：<code>无，注意请不要填任何字符串</code></li>
              </ol>
            </li>
            <li>
              配置OBS的输出，可以参考<a href='https://github.com/ossrs/srs/issues/1147#lagging-encoder'>链接</a>：
              <ol>
                <li>输出模式：<code>高级</code></li>
                <li>编码器：<code>x264</code></li>
                <li>码率控制：<code>CBR</code></li>
                <li>
                  关键帧间隔： <code>3</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, '3')} />
                  </div>
                </li>
                <li>CPU使用预设：<code>veryfast</code></li>
                <li>配置（Profile）：<code>baseline</code></li>
                <li>微调（Tune）： <code>zerolatency</code></li>
              </ol>
            </li>
            <li>点击开始推流</li>
          </ol>
          <p><strong>播放操作步骤：</strong></p>
          <ol>
            <li>
              SRT流播放地址：<br/><code>{srtPlayUrl}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPlayUrl)} />
              </div>
            </li>
            <li>下载<a href='https://ffmpeg.org/download.html' target='_blank' rel='noreferrer'>ffplay</a>，FFmpeg自带的低延迟播放器</li>
            <li>
              Windows，执行命令：<br/>
              <code>{ffplayWindows}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, ffplayWindows)} />
              </div>
            </li>
            <li>
              Mac或Linux，执行命令：<br/>
              <code>{ffplayMac}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, ffplayMac)} />
              </div>
            </li>
            <li>SRT流画面出来较慢，请稍安勿躁</li>
            <li>
              也可以快速预览其他格式的流，注意延迟比直接播放SRT流会高很多：<br/>
              <ul>
                <li>
                  播放<a href={flvPlayer} target='_blank' rel='noreferrer'>HTTP-FLV流</a> <code>{flvUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrl)} />
                  </div>
                </li>
                <li>
                  播放<a href={hlsPlayer} target='_blank' rel='noreferrer'>HLS流</a> <code>{m3u8Url}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8Url)} />
                  </div>
                </li>
                <li>播放<a href={rtcPlayer} target='_blank' rel='noreferrer'>WebRTC流</a></li>
              </ul>
            </li>
          </ol>
          <p>若需要测量延迟请参考<a href='https://github.com/ossrs/srs/issues/1147#lagging-benchmark' target='_blank' rel='noreferrer'>这里</a></p>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
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
            <li>下载<a href='http://www.sinsam.com/' target='_blank' rel='noreferrer'>芯象直播Windows版</a>，注意一定要<code>Windows版</code>，或者下载<a href='http://www.sinsam.com/sm/download/?t=2' target='_blank' rel='noreferrer'>芯象 APP</a>，若你是Mac请用其他方案</li>
            <li>
              配置芯象推流，可以参考<a href='https://github.com/ossrs/srs/issues/1147#lagging-encoder'>链接</a>：
              <ol>
                <li>类型：<code>自定义推流</code></li>
                <li>
                  推流地址：<br/><code>{srtPublishUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='更换流名称'>
                    <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
                  </div> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPublishUrl)} />
                  </div>
                  <br/>
                  <SrsQRCode url={srtPublishUrl} />
                </li>
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
            <li>
              SRT流播放地址：<br/><code>{srtPlayUrl}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPlayUrl)} />
              </div>
            </li>
            <SrsQRCode url={srtPlayUrl} />
            <li>下载<a href='https://ffmpeg.org/download.html' target='_blank' rel='noreferrer'>ffplay</a>，FFmpeg自带的低延迟播放器</li>
            <li>
              Windows，执行命令：<br/>
              <code>{ffplayWindows}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, ffplayWindows)} />
              </div>
            </li>
            <li>
              Mac或Linux，执行命令：<br/>
              <code>{ffplayMac}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, ffplayMac)} />
              </div>
            </li>
            <li>SRT流画面出来较慢，请稍安勿躁</li>
            <li>
              也可以快速预览其他格式的流，注意延迟比直接播放SRT流会高很多：<br/>
              <ul>
                <li>
                  播放<a href={flvPlayer} target='_blank' rel='noreferrer'>HTTP-FLV流</a> <code>{flvUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, flvUrl)} />
                  </div>
                </li>
                <li>
                  播放<a href={hlsPlayer} target='_blank' rel='noreferrer'>HLS流</a> <code>{m3u8Url}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='拷贝'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8Url)} />
                  </div>
                </li>
                <li>播放<a href={rtcPlayer} target='_blank' rel='noreferrer'>WebRTC流</a></li>
              </ul>
            </li>
          </ol>
          <p>若需要测量延迟请参考<a href='https://github.com/ossrs/srs/issues/1147#lagging-benchmark' target='_blank' rel='noreferrer'>这里</a></p>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

function ScenarioSrtEn({updateStreamName, copyToClipboard, urls}) {
  const {srtPublishUrl, srtPlayUrl, flvPlayer, hlsPlayer, flvUrl, m3u8Url, rtcPlayer} = urls;
  const [hostname, setHostname] = React.useState();
  const [srtPort, setSrtPort] = React.useState();
  const [srtPublishStreamId, setPublishStreamId] = React.useState();
  const [srtPlayStreamId, setPlayStreamId] = React.useState();
  const ffplayWindows = `ffplay -fflags nobuffer -flags low_delay -i "${srtPlayUrl}"`;
  const ffplayMac = `ffplay -fflags nobuffer -flags low_delay -i '${srtPlayUrl}'`;

  React.useEffect(() => {
    if (!srtPublishUrl) return;
    const u = new URL(srtPublishUrl.replace('srt://', 'http://'));
    setHostname(u.hostname);
    setSrtPort(u.port);
    setPublishStreamId(`${u.hash}`);
  }, [srtPublishUrl]);

  React.useEffect(() => {
    if (!srtPublishUrl || !srtPlayUrl) return;
    const u = new URL(srtPlayUrl.replace('srt://', 'http://'));
    setPlayStreamId(`${u.hash.split('&')[0]}`);
  }, [srtPublishUrl, srtPlayUrl]);

  return (
    <Accordion defaultActiveKey="2">
      <Accordion.Item eventKey="0">
        <Accordion.Header>Introduction</Accordion.Header>
        <Accordion.Body>
          <p>
            Low latency live streaming service.
          </p>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>vMix ~= 300ms</Accordion.Header>
        <Accordion.Body>
          <div>
            <p style={{display: 'inline-block'}}><strong>Requires:</strong></p>
          </div>
          <ol>
            <li>
              The RTT must <code>&lt;=60ms</code>, please check by<br/>
              <code>ping {window.location.hostname}</code>
            </li>
            <li>
              The network packet loss must <code>&lt;=10%</code>
            </li>
            <li>
              Please use wire network, WiFi is not recommend.
            </li>
            <li>
              The CPU of client and server must <code>&lt;=80%</code>
            </li>
          </ol>
          <p><strong>Usage for Publisher:</strong></p>
          <ol>
            <li>Allow <code>UDP/10080</code> by firewall</li>
            <li>Down vMix from <a href='https://www.vmix.com/software/download.aspx' target='_blank' rel='noreferrer'>here</a></li>
            <li>
              Config vMix as bellow
              <ol>
                <li>
                  Click <code>Settings</code> => <code>Outputs/NDI/SRT</code> / Click <code>Output(1)</code> <br/>
                  Or click <code>External</code> => <code>Outputs/NDI/SRT Settings</code> / Click <code>Output(1)</code>
                </li>
                <li>Check <code>Enable SRT</code></li>
                <li>Type: <code>Caller</code></li>
                <li>
                  Hostname：<code>{hostname}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, hostname)} />
                  </div>
                </li>
                <li>
                  Port：<code>{srtPort}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPort)} />
                  </div>
                </li>
                <li>
                  Latency：<code>20</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, '20')} />
                  </div>
                </li>
                <li>
                  Stream ID：<code>{srtPublishStreamId}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Change'>
                    <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
                  </div> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPublishStreamId)} />
                  </div>
                </li>
              </ol>
            </li>
            <li>Click OK</li>
          </ol>
          <p><strong>Usage for Player:</strong></p>
          <ol>
            <li>Down vMix from <a href='https://www.vmix.com/software/download.aspx' target='_blank' rel='noreferrer'>here</a></li>
            <li>
              Config vMix as bellow
              <ol>
                <li>Click <code>Add Input</code> => <code>Stream/SRT</code></li>
                <li>Stream Type: <code>SRT Caller</code></li>
                <li>
                  Hostname：<code>{hostname}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, hostname)} />
                  </div>
                </li>
                <li>
                  Port：<code>{srtPort}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPort)} />
                  </div>
                </li>
                <li>
                  Latency(延迟)：<code>20</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, '20')} />
                  </div>
                </li>
                <li>
                  Stream ID(流ID)：<code>{srtPlayStreamId}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPlayStreamId)} />
                  </div>
                </li>
                <li>Note: If play failed, uncheck the <code>Use Hardware Decoder</code> and try</li>
              </ol>
            </li>
            <li>Click OK</li>
            <li>
              You're able to play by H5<br/>
              <ul>
                <li>For ffplay, please read other section</li>
                <li>
                  <a href={flvPlayer} target='_blank' rel='noreferrer'>HTTP-FLV</a> <code>{flvUrl}</code> &nbsp;
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
          </ol>
          <p>For latency benchmark, please read <a href='https://github.com/ossrs/srs/issues/1147#lagging-benchmark' target='_blank' rel='noreferrer'>here</a></p>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>OBS+ffplay ~= 300ms</Accordion.Header>
        <Accordion.Body>
          <div>
            <p style={{display: 'inline-block'}}><strong>Requires:</strong></p>
          </div>
          <ol>
            <li>
              The RTT must <code>&lt;=60ms</code>, please check by<br/>
              <code>ping {window.location.hostname}</code>
            </li>
            <li>
              The network packet loss must <code>&lt;=10%</code>
            </li>
            <li>
              Please use wire network, WiFi is not recommend.
            </li>
            <li>
              The CPU of client and server must <code>&lt;=80%</code>
            </li>
          </ol>
          <p><strong>Usage for Publisher:</strong></p>
          <ol>
            <li>Allow <code>UDP/10080</code> by firewall</li>
            <li>Download OBS from <a href='https://obsproject.com/download' target='_blank' rel='noreferrer'>here</a> and install</li>
            <li>
              Config OBS Stream:
              <ol>
                <li>Service: <code>Custom</code></li>
                <li>
                  Server: <br/><code>{srtPublishUrl}</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Change'>
                    <Icon.ArrowRepeat size={20} onClick={updateStreamName}/>
                  </div> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPublishUrl)} />
                  </div>
                </li>
                <li>Stream Key: <code>Empty. Please keep it empty.</code></li>
              </ol>
            </li>
            <li>
              Config OBS Output:
              <ol>
                <li>Output Mode: <code>Advanced</code></li>
                <li>Encoder: <code>x264</code></li>
                <li>Rate Control: <code>CBR</code></li>
                <li>
                  Keyframe Interval: <code>3</code> &nbsp;
                  <div role='button' style={{display: 'inline-block'}} title='Copy'>
                    <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, '3')} />
                  </div>
                </li>
                <li>CPU Usage Preset: <code>veryfast</code></li>
                <li>Profile: <code>baseline</code></li>
                <li>Tune: <code>zerolatency</code></li>
              </ol>
            </li>
            <li>Click Start Streaming</li>
          </ol>
          <p><strong>Usage for Player:</strong></p>
          <ol>
            <li>
              SRT URL: <br/><code>{srtPlayUrl}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='Copy'>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPlayUrl)} />
              </div>
            </li>
            <li>Download ffplay from <a href='https://ffmpeg.org/download.html' target='_blank' rel='noreferrer'>here</a>, a low latency player by FFmpeg</li>
            <li>
              For Windows:<br/>
              <code>{ffplayWindows}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='Copy'>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, ffplayWindows)} />
              </div>
            </li>
            <li>
              For Mac or Linux:<br/>
              <code>{ffplayMac}</code> &nbsp;
              <div role='button' style={{display: 'inline-block'}} title='Copy'>
                <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, ffplayMac)} />
              </div>
            </li>
            <li>It takes a while to render the SRT stream, please wait.</li>
            <li>
              You're able to play by H5<br/>
              <ul>
                <li>For ffplay, please read other section</li>
                <li>
                  <a href={flvPlayer} target='_blank' rel='noreferrer'>HTTP-FLV</a> <code>{flvUrl}</code> &nbsp;
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
          </ol>
          <p>For latency benchmark, please read <a href='https://github.com/ossrs/srs/issues/1147#lagging-benchmark' target='_blank' rel='noreferrer'>here</a></p>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

