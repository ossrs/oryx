import {useNavigate} from "react-router-dom";
import {Accordion, Form, Button, Table} from "react-bootstrap";
import React from "react";
import {Token, Errors, URL, Clipboard} from "../utils";
import axios from "axios";
import SetupCamSecret from '../components/SetupCamSecret';
import moment from "moment";

export default function ScenarioDvr() {
  const navigate = useNavigate();
  const [patternStatus, setPatternStatus] = React.useState();
  const [activeKey, setActiveKey] = React.useState();

  // We must init the activeKey, because the defaultActiveKey only apply when init for Accordion.
  // See https://stackoverflow.com/q/61324259/17679565
  React.useEffect(() => {
    if (!patternStatus) return;

    if (patternStatus.secret) {
      if (patternStatus.all) {
        setActiveKey('3');
      } else {
        setActiveKey('2');
      }
    } else {
      setActiveKey('1');
    }
  }, [patternStatus]);

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/hooks/dvr/query', {
      ...token,
    }).then(res => {
      console.log(`DvrPattern: Query ok, ${JSON.stringify(res.data.data)}`);
      setPatternStatus(res.data.data);
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

  return (
    <>
      { activeKey && <ScenarioDvrImpl activeKey={activeKey} defaultApplyAll={patternStatus.all} /> }
    </>
  );
}

function ScenarioDvrImpl({activeKey, defaultApplyAll}) {
  const navigate = useNavigate();
  const [dvrAll, setDvrAll] = React.useState(defaultApplyAll);
  const [dvrFiles, setDvrFiles] = React.useState();

  React.useEffect(() => {
    const refreshDvrFiles = () => {
      const token = Token.load();
      axios.post('/terraform/v1/hooks/dvr/files', {
        ...token,
      }).then(res => {
        console.log(`DVR: Files ok, ${JSON.stringify(res.data.data)}`);
        setDvrFiles(res.data.data.map(file => {
          if (file.progress) {
            const l = window.location;
            const schema = l.protocol.replace(':', '');
            const httpPort = l.port || (l.protocol === 'http:' ? 80 : 443);
            file.location = `${l.protocol}//${l.host}/terraform/v1/hooks/dvr/hls/${file.uuid}.m3u8`;
            file.preview = `/players/srs_player.html?schema=${schema}&port=${httpPort}&autostart=true&app=terraform/v1/hooks/dvr/hls&stream=${file.uuid}.m3u8`;
          } else {
            const host = `${file.bucket}.cos.${file.region}.myqcloud.com`;
            file.location = `https://${host}/${file.uuid}/index.m3u8`;
            file.preview = `/players/srs_player.html?schema=https&port=443&autostart=true&vhost=${host}&server=${host}&app=${file.uuid}&stream=index.m3u8`;
          }

          return {
            ...file,
            url: URL.build(file.vhost, file.app, file.stream),
            update: moment(file.update),
            duration: Number(file.duration),
            size: Number(file.size / 1024.0 / 1024),
          };
        }).sort((a, b) => {
          return b.update - a.update;
        }).map((file, i) => {
          return {...file, i: i + 1};
        }));
      }).catch(e => {
        const err = e.response.data;
        if (err.code === Errors.auth) {
          alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
          navigate('/routers-logout');
        } else {
          alert(`服务器错误，${err.code}: ${err.data.message}`);
        }
      });
    };

    refreshDvrFiles();
    const timer = setInterval(() => refreshDvrFiles(), 5000);
    return () => clearInterval(timer);
  }, [navigate]);

  const setupDvrPattern = (e) => {
    e.preventDefault();

    const token = Token.load();
    axios.post('/terraform/v1/hooks/dvr/apply', {
      ...token, all: !!dvrAll,
    }).then(res => {
      alert('设置录制规则成功');
      console.log(`DVR: Apply patterns ok, all=${dvrAll}`);
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/routers-logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  };

  const copyToClipboard = async (e, text) => {
    e.preventDefault();

    try {
      await Clipboard.copy(text);
      alert(`已经复制到剪切板`);
    } catch (e) {
      alert(`复制失败，请右键复制链接 ${e}`);
    }
  };

  return (
    <Accordion defaultActiveKey={activeKey}>
      <Accordion.Item eventKey="0">
        <Accordion.Header>场景介绍</Accordion.Header>
        <Accordion.Body>
          <div>
            云录制，指录制视频流到云存储，只要推送到服务器的流都可以录制。
            <p></p>
          </div>
          <p>可应用的具体场景包括：</p>
          <ul>
            <li>直播转点播，录制直播间内容，剪辑精彩片段后做成点播短视频</li>
          </ul>
          <p>使用说明：</p>
          <ul>
            <li>云存储无法使用开源方案搭建，依赖公有云的云存储（<a href='https://buy.cloud.tencent.com/price/cos/calculator' target='_blank' rel='noreferrer'>计费</a>），请先开通<a href='https://console.cloud.tencent.com/cos' target='_blank' rel='noreferrer'>腾讯云COS</a>云存储服务</li>
            <li>第一次使用，需要先设置云存储的访问密钥，我们会自动创建<code>srs-lighthouse</code>开头的存储桶</li>
            <li>具体使用步骤，请根据下面引导操作</li>
          </ul>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>设置云密钥</Accordion.Header>
        <Accordion.Body>
          <SetupCamSecret submitTips=' * 会自动创建云存储的存储桶(Bucket)' />
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>指定录制的流</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3" controlId="formDvrAllCheckbox">
              <Form.Check type="checkbox" label="录制所有流" defaultChecked={dvrAll} onClick={() => setDvrAll(!dvrAll)} />
            </Form.Group>
            <Button variant="primary" type="submit" onClick={(e) => setupDvrPattern(e)}>
              提交
            </Button>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>录制任务列表</Accordion.Header>
        <Accordion.Body>
          {
            dvrFiles?.length ? (
              <Table striped bordered hover>
                <thead>
                <tr>
                  <th>#</th>
                  <th>状态</th>
                  <th>更新时间</th>
                  <th>媒体流</th>
                  <th>时长</th>
                  <th>大小</th>
                  <th>切片</th>
                  <th>COS Bucket</th>
                  <th>地址</th>
                  <th>预览</th>
                </tr>
                </thead>
                <tbody>
                {
                  dvrFiles?.map(file => {
                    return <tr key={file.uuid}>
                      <td>{file.i}</td>
                      <td title='一定时间没有流，会自动完成录制'>{file.progress ? '录制中' : '已完成'}</td>
                      <td>{`${file.update.format('YYYY-MM-DD HH:mm:ss')}`}</td>
                      <td>{file.url}</td>
                      <td>{`${file.duration.toFixed(1)}`}秒</td>
                      <td>{`${file.size.toFixed(1)}`}MB</td>
                      <td>{file.nn}</td>
                      <td><a href={`https://console.cloud.tencent.com/cos/bucket?bucket=${file.bucket}&region=${file.region}&path=%252F${file.uuid}%252F`} target='_blank' rel='noreferrer'>{file.uuid.slice(0, 13)}</a></td>
                      <td><a href={file.location} onClick={(e) => copyToClipboard(e, file.location)} target='_blank' rel='noreferrer'>复制</a></td>
                      <td><a href={file.preview} target='_blank' rel='noreferrer'>预览</a></td>
                    </tr>;
                  })
                }
                </tbody>
              </Table>
            ) : ''
          }
          {!dvrFiles?.length ? '没有流。请开启录制并推流后，等待大约60秒左右，录制列表会自动更新' : ''}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

