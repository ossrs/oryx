//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import {Accordion, Badge, Button, Col, Form, InputGroup, ListGroup, Row, Table} from "react-bootstrap";
import React from "react";
import {Token} from "../utils";
import axios from "axios";
import moment from "moment";
import {useErrorHandler} from "react-error-boundary";
import {useSrsLanguage} from "../components/LanguageSwitch";
import FileUploader from "../components/FileUploader";
import {useTranslation} from "react-i18next";
import {SrsErrorBoundary} from "../components/SrsErrorBoundary";
import {TutorialsButton, useTutorials} from "../components/TutorialsButton";

export default function ScenarioVLive() {
  const [init, setInit] = React.useState();
  const [activeKey, setActiveKey] = React.useState();
  const [secrets, setSecrets] = React.useState();
  const handleError = useErrorHandler();
  const language = useSrsLanguage();

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/ffmpeg/vlive/secret', {
      ...token,
    }).then(res => {
      const secrets = res.data.data;
      setInit(true);
      setSecrets(secrets || {});
      console.log(`VLive: Secret query ok ${JSON.stringify(secrets)}`);
    }).catch(handleError);
  }, [handleError]);

  React.useEffect(() => {
    if (!init || !secrets) return;

    if (secrets.wx?.enabled || secrets.bilibili?.enabled || secrets.kuaishou?.enabled) {
      return setActiveKey('99');
    }

    if (!secrets.wx?.server || !secrets.wx?.secret || !secrets.wx?.enabled) {
      setActiveKey('1');
    } else if (!secrets.bilibili?.server || !secrets.bilibili?.secret || !secrets.bilibili?.enabled) {
      setActiveKey('2');
    } else if (!secrets.kuaishou?.server || !secrets.kuaishou?.secret || !secrets.kuaishou?.enabled) {
      setActiveKey('3');
    } else {
      setActiveKey('99');
    }
  }, [init, secrets]);

  if (!activeKey) return <></>;
  if (language === 'zh') {
    return <ScenarioVLiveImplCn defaultActiveKey={activeKey} defaultSecrets={secrets}/>;
  }
  return <ScenarioVLiveImplEn defaultActiveKey={activeKey} defaultSecrets={secrets}/>;
}

function ScenarioVLiveImplCn({defaultActiveKey, defaultSecrets}) {
  const [wxEnabled, setWxEnabled] = React.useState(defaultSecrets?.wx?.enabled);
  const [wxServer, setWxServer] = React.useState(defaultSecrets?.wx?.server);
  const [wxSecret, setWxSecret] = React.useState(defaultSecrets?.wx?.secret);
  const [wxCustom, setWxCustom] = React.useState(defaultSecrets?.wx?.custom);
  const [wxLabel, setWxLabel] = React.useState(defaultSecrets?.wx?.label);
  const [wxFiles, setWxFiles] = React.useState(defaultSecrets?.wx?.files);
  const [bilibiliEnabled, setBilibiliEnabled] = React.useState(defaultSecrets?.bilibili?.enabled);
  const [bilibiliServer, setBilibiliServer] = React.useState(defaultSecrets?.bilibili?.server);
  const [bilibiliSecret, setBilibiliSecret] = React.useState(defaultSecrets?.bilibili?.secret);
  const [bilibiliCustom, setBilibiliCustom] = React.useState(defaultSecrets?.bilibili?.custom);
  const [bilibiliLabel, setBilibiliLabel] = React.useState(defaultSecrets?.bilibili?.label);
  const [bilibiliFiles, setBilibiliFiles] = React.useState(defaultSecrets?.bilibili?.files);
  const [kuaishouEnabled, setKuaishouEnabled] = React.useState(defaultSecrets?.kuaishou?.enabled);
  const [kuaishouServer, setKuaishouServer] = React.useState(defaultSecrets?.kuaishou?.server);
  const [kuaishouSecret, setKuaishouSecret] = React.useState(defaultSecrets?.kuaishou?.secret);
  const [kuaishouCustom, setKuaishouCustom] = React.useState(defaultSecrets?.kuaishou?.custom);
  const [kuaishouLabel, setKuaishouLabel] = React.useState(defaultSecrets?.kuaishou?.label);
  const [kuaishouFiles, setKuaishouFiles] = React.useState(defaultSecrets?.kuaishou?.files);
  const [vLives, setVLives] = React.useState();
  const [submiting, setSubmiting] = React.useState();
  const handleError = useErrorHandler();

  const vLiveTutorials = useTutorials({
    bilibili: React.useRef([
      {author: '宝哥', id: 'BV1G3411d7Gb'},
    ])
  });

  React.useEffect(() => {
    const refreshStreams = () => {
      const token = Token.load();
      axios.post('/terraform/v1/ffmpeg/vlive/streams', {
        ...token,
      }).then(res => {
        setVLives(res.data.data.map((e, i) => {
          const item = {
            ...e,
            name: {wx: '视频号虚拟直播间', bilibili: 'Bilibili虚拟直播间', kuaishou: '快手虚拟直播间'}[e.platform],
            update: e.frame?.update ? moment(e.frame.update) : null,
            i,
          };

          // Find file source object by uuid(item.source).
          const sources = item.files?.filter(e => e?.uuid === item?.source);
          item.sourceObj = sources?.length ? sources[0] : null;
          return item;
        }));
        console.log(`VLive: Query streams ${JSON.stringify(res.data.data)}`);
      }).catch(handleError);
    };

    refreshStreams();
    const timer = setInterval(() => refreshStreams(), 10 * 1000);
    return () => clearInterval(timer);
  }, [handleError]);

  const updateSecrets = React.useCallback((e, action, platform, server, secret, enabled, custom, label, files, onSuccess) => {
    e.preventDefault();
    if (!files?.length) return alert('请上传视频源');
    if (!server) return alert('请输入推流地址');
    if (custom && !label) return alert('自定义平台请输入名称，否则不好区分虚拟直播状态');

    try {
      setSubmiting(true);

      const token = Token.load();
      axios.post('/terraform/v1/ffmpeg/vlive/secret', {
        ...token, action, platform, server, secret, enabled: !!enabled, custom: !!custom, label, files,
      }).then(res => {
        alert('虚拟直播设置成功');
        onSuccess && onSuccess();
      }).catch(handleError);
    } finally {
      new Promise(resolve => setTimeout(resolve, 3000)).then(() => setSubmiting(false));
    }
  }, [handleError, setSubmiting]);

  return (
    <Accordion defaultActiveKey={defaultActiveKey}>
      <Accordion.Item eventKey="0">
        <Accordion.Header>场景介绍</Accordion.Header>
        <Accordion.Body>
          <div>
            虚拟直播<TutorialsButton prefixLine={true} tutorials={vLiveTutorials} />，是将一个视频文件，用FFmpeg转成直播流，推送到SRS Stack或其他平台。
            <p></p>
          </div>
          <p>可应用的具体场景包括：</p>
          <ul>
            <li>无人直播间，7x24小时获得直播收益</li>
          </ul>
          <p>使用说明：</p>
          <ul>
            <li>首先上传视频文件</li>
            <li>然后设置直播流信息</li>
          </ul>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{wxCustom ? '自定义平台' : '视频号虚拟直播间'} {wxLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>名称</Form.Label>
              <Form.Text> * {wxCustom ? '(必选)' : '(可选)'} 起一个好记的名字</Form.Text>
              <Form.Control as="input" defaultValue={wxLabel} onChange={(e) => setWxLabel(e.target.value)}/>
            </Form.Group>
            <SrsErrorBoundary>
              <ChooseVideoSourceCn platform='wx' vLiveFiles={wxFiles} setVLiveFiles={setWxFiles} />
            </SrsErrorBoundary>
            <Form.Group className="mb-3">
              <Form.Label>推流地址</Form.Label>
              {!wxCustom && <Form.Text> * 请先<a href='https://channels.weixin.qq.com/platform/live/liveBuild' target='_blank' rel='noreferrer'>创建直播</a>，然后获取推流地址</Form.Text>}
              <Form.Control as="input" defaultValue={wxServer} onChange={(e) => setWxServer(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>推流密钥</Form.Label>
              {!wxCustom && <Form.Text> * 请先<a href='https://channels.weixin.qq.com/platform/live/liveBuild' target='_blank' rel='noreferrer'>创建直播</a>，然后获取推流密钥</Form.Text>}
              <Form.Control as="input" defaultValue={wxSecret} onChange={(e) => setWxSecret(e.target.value)}/>
            </Form.Group>
            <Row>
              <Col xs='auto'>
                <Form.Group className="mb-3" controlId="formWxCustomCheckbox">
                  <Form.Check type="checkbox" label="自定义平台" defaultChecked={wxCustom} onClick={() => setWxCustom(!wxCustom)} />
                </Form.Group>
              </Col>
            </Row>
            <Button
              variant="primary"
              type="submit"
              disabled={submiting}
              onClick={(e) => {
                updateSecrets(e, 'update', 'wx', wxServer, wxSecret, !wxEnabled, wxCustom, wxLabel, wxFiles, () => {
                  setWxEnabled(!wxEnabled);
                });
              }}
            >
              {wxEnabled ? '停止直播' : '开始直播'}
            </Button> &nbsp;
            <Form.Text> * 将文件转直播流</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>{bilibiliCustom ? '自定义平台' : 'Bilibili虚拟直播间'} {bilibiliLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>名称</Form.Label>
              <Form.Text> * {bilibiliCustom ? '(必选)' : '(可选)'} 起一个好记的名字</Form.Text>
              <Form.Control as="input" defaultValue={bilibiliLabel} onChange={(e) => setBilibiliLabel(e.target.value)}/>
            </Form.Group>
            <SrsErrorBoundary>
              <ChooseVideoSourceCn platform='bilibili' vLiveFiles={bilibiliFiles} setVLiveFiles={setBilibiliFiles} />
            </SrsErrorBoundary>
            <Form.Group className="mb-3">
              <Form.Label>推流地址</Form.Label>
              {!bilibiliCustom && <Form.Text> * 请先<a href='https://link.bilibili.com/p/center/index#/my-room/start-live' target='_blank' rel='noreferrer'>开始直播</a>，然后获取推流地址</Form.Text>}
              <Form.Control as="input" defaultValue={bilibiliServer} onChange={(e) => setBilibiliServer(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>推流密钥</Form.Label>
              {!bilibiliCustom && <Form.Text> * 请先<a href='https://link.bilibili.com/p/center/index#/my-room/start-live' target='_blank' rel='noreferrer'>开始直播</a>，然后获取推流密钥</Form.Text>}
              <Form.Control as="input" defaultValue={bilibiliSecret} onChange={(e) => setBilibiliSecret(e.target.value)}/>
            </Form.Group>
            <Row>
              <Col xs='auto'>
                <Form.Group className="mb-3" controlId="formBilibiliCustomCheckbox">
                  <Form.Check type="checkbox" label="自定义平台" defaultChecked={bilibiliCustom} onClick={() => setBilibiliCustom(!bilibiliCustom)} />
                </Form.Group>
              </Col>
            </Row>
            <Button
              variant="primary"
              type="submit"
              disabled={submiting}
              onClick={(e) => {
                updateSecrets(e, 'update', 'bilibili', bilibiliServer, bilibiliSecret, !bilibiliEnabled, bilibiliCustom, bilibiliLabel, bilibiliFiles, () => {
                  setBilibiliEnabled(!bilibiliEnabled);
                });
              }}
            >
              {bilibiliEnabled ? '停止直播' : '开始直播'}
            </Button> &nbsp;
            <Form.Text> * 将文件转直播流</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>{kuaishouCustom ? '自定义平台' : '快手虚拟直播间'} {kuaishouLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>名称</Form.Label>
              <Form.Text> * {kuaishouCustom ? '(必选)' : '(可选)'} 起一个好记的名字</Form.Text>
              <Form.Control as="input" defaultValue={kuaishouLabel} onChange={(e) => setKuaishouLabel(e.target.value)}/>
            </Form.Group>
            <SrsErrorBoundary>
              <ChooseVideoSourceCn platform='kuaishou' vLiveFiles={kuaishouFiles} setVLiveFiles={setKuaishouFiles} />
            </SrsErrorBoundary>
            <Form.Group className="mb-3">
              <Form.Label>推流地址</Form.Label>
              {!kuaishouCustom && <Form.Text> * 请先<a href='https://studio.kuaishou.com/live/list' target='_blank' rel='noreferrer'>创建直播</a>，进入直播详情页，然后获取推流地址</Form.Text>}
              <Form.Control as="input" defaultValue={kuaishouServer} onChange={(e) => setKuaishouServer(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>推流密钥</Form.Label>
              {!kuaishouCustom && <Form.Text> * 请先<a href='https://studio.kuaishou.com/live/list' target='_blank' rel='noreferrer'>创建直播</a>，进入直播详情页，然后获取推流密钥</Form.Text>}
              <Form.Control as="input" defaultValue={kuaishouSecret} onChange={(e) => setKuaishouSecret(e.target.value)}/>
            </Form.Group>
            <Row>
              <Col xs='auto'>
                <Form.Group className="mb-3" controlId="formKuaishouCustomCheckbox">
                  <Form.Check type="checkbox" label="自定义平台" defaultChecked={kuaishouCustom} onClick={() => setKuaishouCustom(!kuaishouCustom)} />
                </Form.Group>
              </Col>
            </Row>
            <Button
              variant="primary"
              type="submit"
              disabled={submiting}
              onClick={(e) => {
                updateSecrets(e, 'update', 'kuaishou', kuaishouServer, kuaishouSecret, !kuaishouEnabled, kuaishouCustom, kuaishouLabel, kuaishouFiles, () => {
                  setKuaishouEnabled(!kuaishouEnabled);
                });
              }}
            >
              {kuaishouEnabled ? '停止直播' : '开始直播'}
            </Button> &nbsp;
            <Form.Text> * 将文件转直播流</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="99">
        <Accordion.Header>虚拟直播状态</Accordion.Header>
        <Accordion.Body>
          {
            vLives?.length ? (
              <Table striped bordered hover>
                <thead>
                <tr>
                  <th>#</th>
                  <th>平台</th>
                  <th>状态</th>
                  <th>更新时间</th>
                  <th>视频源</th>
                  <th>日志</th>
                </tr>
                </thead>
                <tbody>
                {
                  vLives?.map(file => {
                    return <tr key={file.platform} style={{verticalAlign: 'middle'}}>
                      <td>{file.i}</td>
                      <td>{file.custom ? (file.label ? '' : '自定义平台') : file.name} {file.label}</td>
                      <td>
                        <Badge bg={file.enabled ? (file.frame ? 'success' : 'primary') : 'secondary'}>
                          {file.enabled ? (file.frame ? '直播中' : '等待中') : '未开启'}
                        </Badge>
                      </td>
                      <td>
                        {file.update && file.update?.format('YYYY-MM-DD')}<br/>
                        {file.update && file.update?.format('HH:mm:ss')}
                      </td>
                      <td>
                        {file.sourceObj?.name}<br/>
                        <VLiveFileFormatInfo file={file.sourceObj}/>
                      </td>
                      <td>{file.frame?.log}</td>
                    </tr>;
                  })
                }
                </tbody>
              </Table>
            ) : ''
          }
          {!vLives?.length ? '没有流。请开启虚拟直播间后，等待大约30秒左右，列表会自动更新' : ''}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

function ScenarioVLiveImplEn({defaultActiveKey, defaultSecrets}) {
  const [wxEnabled, setWxEnabled] = React.useState(defaultSecrets?.wx?.enabled);
  const [wxServer, setWxServer] = React.useState(defaultSecrets?.wx?.server);
  const [wxSecret, setWxSecret] = React.useState(defaultSecrets?.wx?.secret);
  const [wxCustom, setWxCustom] = React.useState(defaultSecrets?.wx?.custom);
  const [wxLabel, setWxLabel] = React.useState(defaultSecrets?.wx?.label);
  const [wxFiles, setWxFiles] = React.useState(defaultSecrets?.wx?.files);
  const [bilibiliEnabled, setBilibiliEnabled] = React.useState(defaultSecrets?.bilibili?.enabled);
  const [bilibiliServer, setBilibiliServer] = React.useState(defaultSecrets?.bilibili?.server || 'rtmp://live.twitch.tv/app');
  const [bilibiliSecret, setBilibiliSecret] = React.useState(defaultSecrets?.bilibili?.secret);
  const [bilibiliCustom, setBilibiliCustom] = React.useState(defaultSecrets?.bilibili?.custom);
  const [bilibiliLabel, setBilibiliLabel] = React.useState(defaultSecrets?.bilibili?.label);
  const [bilibiliFiles, setBilibiliFiles] = React.useState(defaultSecrets?.bilibili?.files);
  const [kuaishouEnabled, setKuaishouEnabled] = React.useState(defaultSecrets?.kuaishou?.enabled);
  const [kuaishouServer, setKuaishouServer] = React.useState(defaultSecrets?.kuaishou?.server || 'rtmps://live-api-s.facebook.com:443/rtmp');
  const [kuaishouSecret, setKuaishouSecret] = React.useState(defaultSecrets?.kuaishou?.secret);
  const [kuaishouCustom, setKuaishouCustom] = React.useState(defaultSecrets?.kuaishou?.custom);
  const [kuaishouLabel, setKuaishouLabel] = React.useState(defaultSecrets?.kuaishou?.label);
  const [kuaishouFiles, setKuaishouFiles] = React.useState(defaultSecrets?.kuaishou?.files);
  const [vLives, setVLives] = React.useState();
  const [submiting, setSubmiting] = React.useState();
  const handleError = useErrorHandler();

  React.useEffect(() => {
    const refreshStreams = () => {
      const token = Token.load();
      axios.post('/terraform/v1/ffmpeg/vlive/streams', {
        ...token,
      }).then(res => {
        setVLives(res.data.data.map((e, i) => {
          const item = {
            ...e,
            name: {wx: 'YouTube', bilibili: 'Twitch', kuaishou: 'Facebook'}[e.platform],
            update: e.frame?.update ? moment(e.frame.update) : null,
            i,
          };

          // Find file source object by uuid(item.source).
          const sources = item.files?.filter(e => e?.uuid === item?.source);
          item.sourceObj = sources?.length ? sources[0] : null;
          return item;
        }));
        console.log(`VLive: Query streams ${JSON.stringify(res.data.data)}`);
      }).catch(handleError);
    };

    refreshStreams();
    const timer = setInterval(() => refreshStreams(), 10 * 1000);
    return () => clearInterval(timer);
  }, [handleError]);

  const updateSecrets = React.useCallback((e, action, platform, server, secret, enabled, custom, label, files, onSuccess) => {
    e.preventDefault();
    if (!files?.length) return alert('Please upload video source');
    if (!server) return alert('Please input stream URL');
    if (custom && !label) return alert('For custom platform, please input label to identify the stream');

    try {
      setSubmiting(true);

      const token = Token.load();
      axios.post('/terraform/v1/ffmpeg/vlive/secret', {
        ...token, action, platform, server, secret, enabled: !!enabled, custom: !!custom, label, files,
      }).then(res => {
        alert('Virtual live stream setup ok');
        onSuccess && onSuccess();
      }).catch(handleError);
    } finally {
      new Promise(resolve => setTimeout(resolve, 3000)).then(() => setSubmiting(false));
    }
  }, [handleError, setSubmiting]);

  return (
    <Accordion defaultActiveKey={[defaultActiveKey]} alwaysOpen>
      <Accordion.Item eventKey="0">
        <Accordion.Header>Introduction</Accordion.Header>
        <Accordion.Body>
          <div>
            Virtual live streaming is the process of converting a video file into a live stream using FFmpeg and pushing it to the SRS Stack or other platforms.
            <p></p>
          </div>
          <p>Specific application scenarios include:</p>
          <ul>
            <li>Unmanned live streaming rooms, 7x24 hours of live streaming revenue</li>
          </ul>
          <p>Instructions for use:</p>
          <ul>
            <li>First, upload the video file</li>
            <li>Then, set the live stream information</li>
          </ul>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{wxCustom ? 'Custom' : 'YouTube'} {wxLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>Label</Form.Label>
              <Form.Text> * {wxCustom ? '(Required)' : '(Optional)'} Virtual live stream label</Form.Text>
              <Form.Control as="input" defaultValue={wxLabel} onChange={(e) => setWxLabel(e.target.value)}/>
            </Form.Group>
            <SrsErrorBoundary>
              <ChooseVideoSourceEn platform='wx' vLiveFiles={wxFiles} setVLiveFiles={setWxFiles} />
            </SrsErrorBoundary>
            <Form.Group className="mb-3">
              <Form.Label>{wxCustom ? 'Server' : 'Stream URL'}</Form.Label>
              {!wxCustom && <Form.Text> * Please click <a href='https://studio.youtube.com/channel/UC/livestreaming' target='_blank' rel='noreferrer'>Go live</a>, then copy the Stream URL</Form.Text>}
              <Form.Control as="input" defaultValue={wxServer} onChange={(e) => setWxServer(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Stream Key</Form.Label>
              {!wxCustom && <Form.Text> * Please click <a href='https://studio.youtube.com/channel/UC/livestreaming' target='_blank' rel='noreferrer'>Go live</a>, then copy the Stream Key</Form.Text>}
              <Form.Control as="input" defaultValue={wxSecret} onChange={(e) => setWxSecret(e.target.value)}/>
            </Form.Group>
            <Row>
              <Col xs='auto'>
                <Form.Group className="mb-3" controlId="formWxCustomCheckbox">
                  <Form.Check type="checkbox" label="Custom" defaultChecked={wxCustom} onClick={() => setWxCustom(!wxCustom)} />
                </Form.Group>
              </Col>
            </Row>
            <Button
              variant="primary"
              type="submit"
              disabled={submiting}
              onClick={(e) => {
                updateSecrets(e, 'update', 'wx', wxServer, wxSecret, !wxEnabled, wxCustom, wxLabel, wxFiles, () => {
                  setWxEnabled(!wxEnabled);
                });
              }}
            >
              {wxEnabled ? 'Stop Virtual Live' : 'Start Virtual Live'}
            </Button> &nbsp;
            <Form.Text> * Convert file to live stream</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>{bilibiliCustom ? 'Custom' : 'Twitch'} {bilibiliLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>Label</Form.Label>
              <Form.Text> * {bilibiliCustom ? '(Required)' : '(Optional)'} Virtual live stream label</Form.Text>
              <Form.Control as="input" defaultValue={bilibiliLabel} onChange={(e) => setBilibiliLabel(e.target.value)}/>
            </Form.Group>
            <SrsErrorBoundary>
              <ChooseVideoSourceEn platform='bilibili' vLiveFiles={bilibiliFiles} setVLiveFiles={setBilibiliFiles} />
            </SrsErrorBoundary>
            <Form.Group className="mb-3">
              <Form.Label>Server</Form.Label>
              <Form.Control as="input" defaultValue={bilibiliServer} onChange={(e) => setBilibiliServer(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Stream Key</Form.Label>
              {!bilibiliCustom && <Form.Text> * Please click <a href='https://www.twitch.tv/dashboard/settings' target='_blank' rel='noreferrer'>Dashboard</a>, then click Settings Stream and Copy the stream key</Form.Text>}
              <Form.Control as="input" defaultValue={bilibiliSecret} onChange={(e) => setBilibiliSecret(e.target.value)}/>
            </Form.Group>
            <Row>
              <Col xs='auto'>
                <Form.Group className="mb-3" controlId="formBilibiliCustomCheckbox">
                  <Form.Check type="checkbox" label="Custom" defaultChecked={bilibiliCustom} onClick={() => setBilibiliCustom(!bilibiliCustom)} />
                </Form.Group>
              </Col>
            </Row>
            <Button
              variant="primary"
              type="submit"
              disabled={submiting}
              onClick={(e) => {
                updateSecrets(e, 'update', 'bilibili', bilibiliServer, bilibiliSecret, !bilibiliEnabled, bilibiliCustom, bilibiliLabel, bilibiliFiles, () => {
                  setBilibiliEnabled(!bilibiliEnabled);
                });
              }}
            >
              {bilibiliEnabled ? 'Stop Virtual Live' : 'Start Virtual Live'}
            </Button> &nbsp;
            <Form.Text> * Convert file to live stream</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>{kuaishouCustom ? 'Custom' : 'Facebook'} {kuaishouLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>Label</Form.Label>
              <Form.Text> * {kuaishouCustom ? '(Required)' : '(Optional)'} Virtual live stream label</Form.Text>
              <Form.Control as="input" defaultValue={kuaishouLabel} onChange={(e) => setKuaishouLabel(e.target.value)}/>
            </Form.Group>
            <SrsErrorBoundary>
              <ChooseVideoSourceEn platform='kuaishou' vLiveFiles={kuaishouFiles} setVLiveFiles={setKuaishouFiles} />
            </SrsErrorBoundary>
            <Form.Group className="mb-3">
              <Form.Label>Server</Form.Label>
              <Form.Control as="input" defaultValue={kuaishouServer} onChange={(e) => setKuaishouServer(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Stream Key</Form.Label>
              {!kuaishouCustom && <Form.Text> * Please click <a href='https://www.facebook.com/live/producer?ref=OBS' target='_blank' rel='noreferrer'>Live Producer</a>, then click Go live then select Streaming software, copy the Stream key</Form.Text>}
              <Form.Control as="input" defaultValue={kuaishouSecret} onChange={(e) => setKuaishouSecret(e.target.value)}/>
            </Form.Group>
            <Row>
              <Col xs='auto'>
                <Form.Group className="mb-3" controlId="formKuaishouCustomCheckbox">
                  <Form.Check type="checkbox" label="Custom" defaultChecked={kuaishouCustom} onClick={() => setKuaishouCustom(!kuaishouCustom)} />
                </Form.Group>
              </Col>
            </Row>
            <Button
              variant="primary"
              type="submit"
              disabled={submiting}
              onClick={(e) => {
                updateSecrets(e, 'update', 'kuaishou', kuaishouServer, kuaishouSecret, !kuaishouEnabled, kuaishouCustom, kuaishouLabel, kuaishouFiles, () => {
                  setKuaishouEnabled(!kuaishouEnabled);
                });
              }}
            >
              {kuaishouEnabled ? 'Stop Virtual Live' : 'Start Virtual Live'}
            </Button> &nbsp;
            <Form.Text> * Convert file to live stream</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="99">
        <Accordion.Header>Virtual Live Status</Accordion.Header>
        <Accordion.Body>
          {
            vLives?.length ? (
              <Table striped bordered hover>
                <thead>
                <tr>
                  <th>#</th>
                  <th>Platform</th>
                  <th>Status</th>
                  <th>Update</th>
                  <th>Source Stream</th>
                  <th>Logging</th>
                </tr>
                </thead>
                <tbody>
                {
                  vLives?.map(file => {
                    return <tr key={file.platform} style={{verticalAlign: 'middle'}}>
                      <td>{file.i}</td>
                      <td>{file.custom ? (file.label ? '' : 'Custom') : file.name} {file.label}</td>
                      <td>
                        <Badge bg={file.enabled ? (file.frame ? 'success' : 'primary') : 'secondary'}>
                          {file.enabled ? (file.frame ? 'Streaming' : 'Waiting') : 'Inactive'}
                        </Badge>
                      </td>
                      <td>
                        {file.update && file.update?.format('YYYY-MM-DD')}<br/>
                        {file.update && file.update?.format('HH:mm:ss')}
                      </td>
                      <td>
                        {file.sourceObj?.name}<br/>
                        <VLiveFileFormatInfo file={file.sourceObj}/>
                      </td>
                      <td>{file.frame?.log}</td>
                    </tr>;
                  })
                }
                </tbody>
              </Table>
            ) : ''
          }
          {!vLives?.length ? 'There is no stream. Please start the virtual live room and wait for about 30 seconds, the list will update automatically.' : ''}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

function VLiveFileList({files, onChangeFiles}) {
  const {t} = useTranslation();
  return (
    <Row>
      <Col xs='auto'>
        <ListGroup>
          {files.map((f, index) => {
            return <ListGroup.Item key={index}>
              {f.name} &nbsp;
              <VLiveFileFormatInfo file={f}/> &nbsp;
              <VLiveFileVideoInfo file={f}/> &nbsp;
              <VLiveFileAudioInfo file={f}/>
            </ListGroup.Item>;
          })}
        </ListGroup>
      </Col>
      <Col>
        <Button variant="primary" type="button" onClick={onChangeFiles}>{t('helper.changeFiles')}</Button>
      </Col>
    </Row>
  );
}

function ChooseVideoSourceCn({platform, vLiveFiles, setVLiveFiles}) {
  const [checkType, setCheckType] = React.useState('upload');
  return (<>
    <Form.Group className="mb-2">
      <Form.Label>视频源</Form.Label>
      <Form.Text> * 虚拟直播就是将视频源(文件)转换成直播流</Form.Text>
      <Form.Check type="radio" label="上传本地文件" id={'upload-' + platform} checked={checkType === 'upload'}
        name={'chooseSource-' + platform} onChange={e => setCheckType('upload')}
      />
      {checkType === 'upload' && 
      <SrsErrorBoundary>
        <VLiveFileUploaderCn platform={platform} vLiveFiles={vLiveFiles} setVLiveFiles={setVLiveFiles} />
      </SrsErrorBoundary>
      }
    </Form.Group>
    <Form.Group className="mb-3">
      <InputGroup>
        <Form.Check type="radio" label="指定服务器文件" id={'server-' + platform} checked={checkType === 'server'}
                    name={'chooseSource' + platform} onChange={e => setCheckType('server')}
        /> &nbsp;
        <Form.Text> * 文件必须在 /data 目录下面</Form.Text>
      </InputGroup>
      {checkType === 'server' &&
      <SrsErrorBoundary>
        <VLiveFileSelectorCn platform={platform} vLiveFiles={vLiveFiles} setVLiveFiles={setVLiveFiles}/>
      </SrsErrorBoundary>
      }
    </Form.Group>
  </>);
}

function ChooseVideoSourceEn({platform, vLiveFiles, setVLiveFiles}) {
  const [checkType, setCheckType] = React.useState('upload');
  return (<>
    <Form.Group className="mb-2">
      <Form.Label>Live Stream Source</Form.Label>
      <Form.Text> * Virtual live streaming is the process of converting a video source (file) into a live stream.</Form.Text>
      <Form.Check type="radio" label="Upload local file" id={'upload-' + platform} checked={checkType === 'upload'}
                  name={'chooseSource-' + platform} onChange={e => setCheckType('upload')}
      />
      {checkType === 'upload' &&
        <SrsErrorBoundary>
          <VLiveFileUploaderEn platform={platform} vLiveFiles={vLiveFiles} setVLiveFiles={setVLiveFiles} />
        </SrsErrorBoundary>
      }
    </Form.Group>
    <Form.Group className="mb-3">
      <InputGroup>
        <Form.Check type="radio" label="Use server file" id={'server-' + platform} checked={checkType === 'server'}
                    name={'chooseSource' + platform} onChange={e => setCheckType('server')}
        /> &nbsp;
        <Form.Text> * The file must be in the /data directory.</Form.Text>
      </InputGroup>
      {checkType === 'server' &&
        <SrsErrorBoundary>
          <VLiveFileSelectorEn platform={platform} vLiveFiles={vLiveFiles} setVLiveFiles={setVLiveFiles}/>
        </SrsErrorBoundary>
      }
    </Form.Group>
  </>);
}

function VLiveFileSelectorCn({platform, vLiveFiles, setVLiveFiles}) {
  const handleError = useErrorHandler();
  const [inputFile, setInputFile] = React.useState('');
  
  const CheckLocalFile = function() {
    if (!inputFile) return alert('请输入文件路径');
    if (!inputFile.startsWith('/data') && !inputFile.startsWith('upload/')) return alert('文件必须在 /data 目录下');

    const fileExtension = inputFile.slice(inputFile.lastIndexOf('.'));
    if (!['.mp4', '.flv', '.ts'].includes(fileExtension)) return alert('文件必须是 mp4/flv/ts 格式');

    const token = Token.load();
    axios.post(`/terraform/v1/ffmpeg/vlive/server?file=${inputFile}`).then(res => {
      console.log(`检查服务器文件成功，${JSON.stringify(res.data.data)}`);
      const localFileObj = res.data.data;
      const files = [{name: localFileObj.name, size: localFileObj.size, uuid: localFileObj.uuid, target: localFileObj.target}];
      axios.post('/terraform/v1/ffmpeg/vlive/source', {
        ...token, platform, files,
      }).then(res => {
        console.log(`更新虚拟直播源为服务器文件成功，${JSON.stringify(res.data.data)}`);
        setVLiveFiles(res.data.data.files);
      }).catch(handleError);
    }).catch(handleError);
  };

  return (<>
    <Form.Control as="div">
      {!vLiveFiles?.length && <>
        <Row>
          <Col>
            <Form.Control type="text" value={inputFile} placeholder="请输入文件路径" onChange={e => setInputFile(e.target.value)} />
          </Col>
          <Col xs="auto">
            <Button variant="primary" onClick={CheckLocalFile}>确认</Button>
          </Col>
        </Row></>
      }
      {vLiveFiles?.length && <VLiveFileList files={vLiveFiles} onChangeFiles={(e) => setVLiveFiles(null)}/>}
    </Form.Control>
  </>);
}

function VLiveFileSelectorEn({platform, vLiveFiles, setVLiveFiles}) {
  const handleError = useErrorHandler();
  const [inputFile, setInputFile] = React.useState('');

  const CheckLocalFile = function() {
    if (!inputFile) return alert('Please input file path');
    if (!inputFile.startsWith('/data') && !inputFile.startsWith('upload/')) return alert('The file must be in the /data directory.');

    const fileExtension = inputFile.slice(inputFile.lastIndexOf('.'));
    if (!['.mp4', '.flv', '.ts'].includes(fileExtension)) return alert('The file must be in mp4/flv/ts format.');

    const token = Token.load();
    axios.post(`/terraform/v1/ffmpeg/vlive/server?file=${inputFile}`).then(res => {
      console.log(`Check server file ok，${JSON.stringify(res.data.data)}`);
      const localFileObj = res.data.data;
      const files = [{name: localFileObj.name, size: localFileObj.size, uuid: localFileObj.uuid, target: localFileObj.target}];
      axios.post('/terraform/v1/ffmpeg/vlive/source', {
        ...token, platform, files,
      }).then(res => {
        console.log(`Setup the virtual live file ok，${JSON.stringify(res.data.data)}`);
        setVLiveFiles(res.data.data.files);
      }).catch(handleError);
    }).catch(handleError);
  };

  return (<>
    <Form.Control as="div">
      {!vLiveFiles?.length && <>
        <Row>
          <Col>
            <Form.Control type="text" value={inputFile} placeholder="Please input file path" onChange={e => setInputFile(e.target.value)} />
          </Col>
          <Col xs="auto">
            <Button variant="primary" onClick={CheckLocalFile}>Submit</Button>
          </Col>
        </Row></>
      }
      {vLiveFiles?.length && <VLiveFileList files={vLiveFiles} onChangeFiles={(e) => setVLiveFiles(null)}/>}
    </Form.Control>
  </>);
}

function VLiveFileUploaderCn({platform, vLiveFiles, setVLiveFiles}) {
  const handleError = useErrorHandler();
  const updateSources = React.useCallback((platform, files, setFiles) => {
    if (!files?.length) return alert('无上传文件');

    const token = Token.load();
    axios.post('/terraform/v1/ffmpeg/vlive/source', {
      ...token, platform, files: files.map(f => {
        return {name: f.name, size: f.size, uuid: f.uuid, target: f.target};
      }),
    }).then(res => {
      console.log(`虚拟直播文件源设置成功, ${JSON.stringify(res.data.data)}`);
      setFiles(res.data.data.files);
    }).catch(handleError);
  }, [handleError]);

  return (<>
    <Form.Control as='div'>
      {!vLiveFiles?.length && <FileUploader onFilesUploaded={(files) => updateSources(platform, files, setVLiveFiles)}/>}
      {vLiveFiles?.length && <VLiveFileList files={vLiveFiles} onChangeFiles={(e) => setVLiveFiles(null)}/>}
    </Form.Control>
  </>);
}

function VLiveFileUploaderEn({platform, vLiveFiles, setVLiveFiles}) {
  const handleError = useErrorHandler();
  const updateSources = React.useCallback((platform, files, setFiles) => {
    if (!files?.length) return alert('No file selected');

    const token = Token.load();
    axios.post('/terraform/v1/ffmpeg/vlive/source', {
      ...token, platform, files: files.map(f => {
        return {name: f.name, size: f.size, uuid: f.uuid, target: f.target};
      }),
    }).then(res => {
      console.log(`Set file source ok, ${JSON.stringify(res.data.data)}`);
      setFiles(res.data.data.files);
    }).catch(handleError);
  }, [handleError]);

  return (<>
    <Form.Control as='div'>
      {!vLiveFiles?.length && <FileUploader onFilesUploaded={(files) => updateSources(platform, files, setVLiveFiles)}/>}
      {vLiveFiles?.length && <VLiveFileList files={vLiveFiles} onChangeFiles={(e) => setVLiveFiles(null)}/>}
    </Form.Control>
  </>);
}

function VLiveFileFormatInfo({file}) {
  const f = file;
  if (!f?.format) return <></>;
  return <>
    {Number(f?.size/1024/1024).toFixed(1)}MB &nbsp;
    {Number(f?.format?.duration).toFixed(0)}s &nbsp;
    {Number(f?.format?.bit_rate/1000).toFixed(1)}Kbps
  </>;
}

function VLiveFileVideoInfo({file}) {
  const f = file;
  if (!f?.video) return <></>;
  return <>Video({f?.video?.codec_name} {f?.video?.profile} {f?.video?.width}x{f?.video?.height})</>;
}

function VLiveFileAudioInfo({file}) {
  const f = file;
  if (!f?.audio) return <></>;
  return <>Audio({f?.audio?.codec_name} {f?.audio?.sample_rate}HZ {f?.audio?.channels}CH)</>;
}

