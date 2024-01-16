//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {Accordion, Badge, Button, Col, Form, ListGroup, Row, Table} from "react-bootstrap";
import {Token} from "../utils";
import axios from "axios";
import moment from "moment";
import {useErrorHandler} from "react-error-boundary";
import {useSrsLanguage} from "../components/LanguageSwitch";
import {useTranslation} from "react-i18next";
import {SrsErrorBoundary} from "../components/SrsErrorBoundary";

export default function ScenarioCamera() {
  const [init, setInit] = React.useState();
  const [activeKey, setActiveKey] = React.useState();
  const [secrets, setSecrets] = React.useState();
  const handleError = useErrorHandler();
  const language = useSrsLanguage();

  React.useEffect(() => {
    axios.post('/terraform/v1/ffmpeg/camera/secret', {
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      const secrets = res.data.data;
      setInit(true);
      setSecrets(secrets || {});
      console.log(`Camera: Secret query ok ${JSON.stringify(secrets)}`);
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
    return <ScenarioCameraImplCn defaultActiveKey={activeKey} defaultSecrets={secrets}/>;
  }
  return <ScenarioCameraImplEn defaultActiveKey={activeKey} defaultSecrets={secrets}/>;
}

function ScenarioCameraImplCn({defaultActiveKey, defaultSecrets}) {
  const [wxEnabled, setWxEnabled] = React.useState(defaultSecrets?.wx?.enabled);
  const [wxServer, setWxServer] = React.useState(defaultSecrets?.wx?.server);
  const [wxSecret, setWxSecret] = React.useState(defaultSecrets?.wx?.secret);
  const [wxCustom, setWxCustom] = React.useState(defaultSecrets?.wx?.custom);
  const [wxLabel, setWxLabel] = React.useState(defaultSecrets?.wx?.label);
  const [wxFiles, setWxFiles] = React.useState(defaultSecrets?.wx?.files);
  const [wxExtraAudio, setWxExtraAudio] = React.useState(defaultSecrets?.wx?.extraAudio);
  const [bilibiliEnabled, setBilibiliEnabled] = React.useState(defaultSecrets?.bilibili?.enabled);
  const [bilibiliServer, setBilibiliServer] = React.useState(defaultSecrets?.bilibili?.server);
  const [bilibiliSecret, setBilibiliSecret] = React.useState(defaultSecrets?.bilibili?.secret);
  const [bilibiliCustom, setBilibiliCustom] = React.useState(defaultSecrets?.bilibili?.custom);
  const [bilibiliLabel, setBilibiliLabel] = React.useState(defaultSecrets?.bilibili?.label);
  const [bilibiliFiles, setBilibiliFiles] = React.useState(defaultSecrets?.bilibili?.files);
  const [bilibiliExtraAudio, setBilibiliExtraAudio] = React.useState(defaultSecrets?.bilibili?.extraAudio);
  const [kuaishouEnabled, setKuaishouEnabled] = React.useState(defaultSecrets?.kuaishou?.enabled);
  const [kuaishouServer, setKuaishouServer] = React.useState(defaultSecrets?.kuaishou?.server);
  const [kuaishouSecret, setKuaishouSecret] = React.useState(defaultSecrets?.kuaishou?.secret);
  const [kuaishouCustom, setKuaishouCustom] = React.useState(defaultSecrets?.kuaishou?.custom);
  const [kuaishouLabel, setKuaishouLabel] = React.useState(defaultSecrets?.kuaishou?.label);
  const [kuaishouFiles, setKuaishouFiles] = React.useState(defaultSecrets?.kuaishou?.files);
  const [kuaishouExtraAudio, setKuaishouExtraAudio] = React.useState(defaultSecrets?.kuaishou?.extraAudio);
  const [cameras, setCameras] = React.useState();
  const [submiting, setSubmiting] = React.useState();
  const handleError = useErrorHandler();

  React.useEffect(() => {
    const refreshStreams = () => {
      axios.post('/terraform/v1/ffmpeg/camera/streams', {
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        setCameras(res.data.data.map((e, i) => {
          const item = {
            ...e,
            name: {wx: '视频号直播间', bilibili: 'Bilibili直播间', kuaishou: '快手直播间'}[e.platform],
            update: e.frame?.update ? moment(e.frame.update) : null,
            i,
          };

          // Find file source object by uuid(item.source).
          const sources = item.files?.filter(e => e?.uuid === item?.source);
          item.sourceObj = sources?.length ? sources[0] : null;
          return item;
        }));
        console.log(`Camera: Query streams ${JSON.stringify(res.data.data)}`);
      }).catch(handleError);
    };

    refreshStreams();
    const timer = setInterval(() => refreshStreams(), 10 * 1000);
    return () => clearInterval(timer);
  }, [handleError]);

  const updateSecrets = React.useCallback((e, action, platform, server, secret, enabled, custom, label, files, extraAudio, onSuccess) => {
    e.preventDefault();
    if (!files?.length) return alert('请上传视频源');
    if (!server) return alert('请输入推流地址');
    if (custom && !label) return alert('自定义平台请输入名称，否则不好区分直播状态');

    try {
      setSubmiting(true);

      axios.post('/terraform/v1/ffmpeg/camera/secret', {
        action, platform, server, secret, enabled: !!enabled, custom: !!custom, label, files,
        extraAudio,
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        alert('直播设置成功');
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
            摄像头直播，是将一个摄像头的流，用FFmpeg转成直播流，推送到SRS Stack或其他平台。
            <p></p>
          </div>
          <p>可应用的具体场景包括：</p>
          <ul>
            <li>无人直播间，7x24小时获得直播收益</li>
          </ul>
          <p>使用说明：</p>
          <ul>
            <li>首先设置摄像头信息</li>
            <li>然后设置直播流信息</li>
          </ul>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{wxCustom ? '自定义平台' : '视频号直播间'} {wxLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>名称</Form.Label>
              <Form.Text> * {wxCustom ? '(必选)' : '(可选)'} 起一个好记的名字</Form.Text>
              <Form.Control as="input" defaultValue={wxLabel} onChange={(e) => setWxLabel(e.target.value)}/>
            </Form.Group>
            <SrsErrorBoundary>
              <ChooseVideoSourceCn platform='wx' cameraFiles={wxFiles} setCameraFiles={setWxFiles} />
              <CameraExtraAudioTrack extraAudio={wxExtraAudio} setExtraAudio={setWxExtraAudio} />
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
                updateSecrets(e, 'update', 'wx', wxServer, wxSecret, !wxEnabled, wxCustom, wxLabel, wxFiles, wxExtraAudio, () => {
                  setWxEnabled(!wxEnabled);
                });
              }}
            >
              {wxEnabled ? '停止直播' : '开始直播'}
            </Button> &nbsp;
            <Form.Text> * 将摄像头转直播流</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>{bilibiliCustom ? '自定义平台' : 'Bilibili直播间'} {bilibiliLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>名称</Form.Label>
              <Form.Text> * {bilibiliCustom ? '(必选)' : '(可选)'} 起一个好记的名字</Form.Text>
              <Form.Control as="input" defaultValue={bilibiliLabel} onChange={(e) => setBilibiliLabel(e.target.value)}/>
            </Form.Group>
            <SrsErrorBoundary>
              <ChooseVideoSourceCn platform='bilibili' cameraFiles={bilibiliFiles} setCameraFiles={setBilibiliFiles} />
              <CameraExtraAudioTrack extraAudio={bilibiliExtraAudio} setExtraAudio={setBilibiliExtraAudio} />
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
                updateSecrets(e, 'update', 'bilibili', bilibiliServer, bilibiliSecret, !bilibiliEnabled, bilibiliCustom, bilibiliLabel, bilibiliFiles, bilibiliExtraAudio, () => {
                  setBilibiliEnabled(!bilibiliEnabled);
                });
              }}
            >
              {bilibiliEnabled ? '停止直播' : '开始直播'}
            </Button> &nbsp;
            <Form.Text> * 将摄像头转直播流</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>{kuaishouCustom ? '自定义平台' : '快手直播间'} {kuaishouLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>名称</Form.Label>
              <Form.Text> * {kuaishouCustom ? '(必选)' : '(可选)'} 起一个好记的名字</Form.Text>
              <Form.Control as="input" defaultValue={kuaishouLabel} onChange={(e) => setKuaishouLabel(e.target.value)}/>
            </Form.Group>
            <SrsErrorBoundary>
              <ChooseVideoSourceCn platform='kuaishou' cameraFiles={kuaishouFiles} setCameraFiles={setKuaishouFiles} />
              <CameraExtraAudioTrack extraAudio={kuaishouExtraAudio} setExtraAudio={setKuaishouExtraAudio} />
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
                updateSecrets(e, 'update', 'kuaishou', kuaishouServer, kuaishouSecret, !kuaishouEnabled, kuaishouCustom, kuaishouLabel, kuaishouFiles, kuaishouExtraAudio, () => {
                  setKuaishouEnabled(!kuaishouEnabled);
                });
              }}
            >
              {kuaishouEnabled ? '停止直播' : '开始直播'}
            </Button> &nbsp;
            <Form.Text> * 将摄像头转直播流</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="99">
        <Accordion.Header>摄像头直播状态</Accordion.Header>
        <Accordion.Body>
          {
            cameras?.length ? (
              <Table striped bordered hover>
                <thead>
                <tr>
                  <th>#</th>
                  <th>平台</th>
                  <th>状态</th>
                  <th>更新时间</th>
                  <th>视频源</th>
                  <th>额外音频</th>
                  <th>日志</th>
                </tr>
                </thead>
                <tbody>
                {
                  cameras?.map(file => {
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
                        <CameraFileFormatInfo file={file.sourceObj}/>
                      </td>
                      <td>{file?.extraAudio}</td>
                      <td>{file.frame?.log}</td>
                    </tr>;
                  })
                }
                </tbody>
              </Table>
            ) : ''
          }
          {!cameras?.length ? '没有流。请开启直播间后，等待大约30秒左右，列表会自动更新' : ''}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

function ScenarioCameraImplEn({defaultActiveKey, defaultSecrets}) {
  const [wxEnabled, setWxEnabled] = React.useState(defaultSecrets?.wx?.enabled);
  const [wxServer, setWxServer] = React.useState(defaultSecrets?.wx?.server);
  const [wxSecret, setWxSecret] = React.useState(defaultSecrets?.wx?.secret);
  const [wxCustom, setWxCustom] = React.useState(defaultSecrets?.wx?.custom);
  const [wxLabel, setWxLabel] = React.useState(defaultSecrets?.wx?.label);
  const [wxFiles, setWxFiles] = React.useState(defaultSecrets?.wx?.files);
  const [wxExtraAudio, setWxExtraAudio] = React.useState(defaultSecrets?.wx?.extraAudio);
  const [bilibiliEnabled, setBilibiliEnabled] = React.useState(defaultSecrets?.bilibili?.enabled);
  const [bilibiliServer, setBilibiliServer] = React.useState(defaultSecrets?.bilibili?.server || 'rtmp://live.twitch.tv/app');
  const [bilibiliSecret, setBilibiliSecret] = React.useState(defaultSecrets?.bilibili?.secret);
  const [bilibiliCustom, setBilibiliCustom] = React.useState(defaultSecrets?.bilibili?.custom);
  const [bilibiliLabel, setBilibiliLabel] = React.useState(defaultSecrets?.bilibili?.label);
  const [bilibiliFiles, setBilibiliFiles] = React.useState(defaultSecrets?.bilibili?.files);
  const [bilibiliExtraAudio, setBilibiliExtraAudio] = React.useState(defaultSecrets?.bilibili?.extraAudio);
  const [kuaishouEnabled, setKuaishouEnabled] = React.useState(defaultSecrets?.kuaishou?.enabled);
  const [kuaishouServer, setKuaishouServer] = React.useState(defaultSecrets?.kuaishou?.server || 'rtmps://live-api-s.facebook.com:443/rtmp');
  const [kuaishouSecret, setKuaishouSecret] = React.useState(defaultSecrets?.kuaishou?.secret);
  const [kuaishouCustom, setKuaishouCustom] = React.useState(defaultSecrets?.kuaishou?.custom);
  const [kuaishouLabel, setKuaishouLabel] = React.useState(defaultSecrets?.kuaishou?.label);
  const [kuaishouFiles, setKuaishouFiles] = React.useState(defaultSecrets?.kuaishou?.files);
  const [kuaishouExtraAudio, setKuaishouExtraAudio] = React.useState(defaultSecrets?.kuaishou?.extraAudio);
  const [vLives, setVLives] = React.useState();
  const [submiting, setSubmiting] = React.useState();
  const handleError = useErrorHandler();

  React.useEffect(() => {
    const refreshStreams = () => {
      axios.post('/terraform/v1/ffmpeg/camera/streams', {
      }, {
        headers: Token.loadBearerHeader(),
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
        console.log(`Camera: Query streams ${JSON.stringify(res.data.data)}`);
      }).catch(handleError);
    };

    refreshStreams();
    const timer = setInterval(() => refreshStreams(), 10 * 1000);
    return () => clearInterval(timer);
  }, [handleError]);

  const updateSecrets = React.useCallback((e, action, platform, server, secret, enabled, custom, label, files, extraAudio, onSuccess) => {
    e.preventDefault();
    if (!files?.length) return alert('Please upload video source');
    if (!server) return alert('Please input stream URL');
    if (custom && !label) return alert('For custom platform, please input label to identify the stream');

    try {
      setSubmiting(true);

      axios.post('/terraform/v1/ffmpeg/camera/secret', {
        action, platform, server, secret, enabled: !!enabled, custom: !!custom, label, files,
        extraAudio,
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        alert('IP camera stream setup ok');
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
            IP camera streaming is the process of converting the stream from IP Camera into a live stream using FFmpeg and pushing it to the SRS Stack or other platforms.
            <p></p>
          </div>
          <p>Specific application scenarios include:</p>
          <ul>
            <li>Unmanned live streaming rooms, 7x24 hours of live streaming revenue</li>
          </ul>
          <p>Instructions for use:</p>
          <ul>
            <li>First, setup the IP Camera</li>
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
              <Form.Text> * {wxCustom ? '(Required)' : '(Optional)'} IP camera stream label</Form.Text>
              <Form.Control as="input" defaultValue={wxLabel} onChange={(e) => setWxLabel(e.target.value)}/>
            </Form.Group>
            <SrsErrorBoundary>
              <ChooseVideoSourceEn platform='wx' cameraFiles={wxFiles} setCameraFiles={setWxFiles} />
              <CameraExtraAudioTrack extraAudio={wxExtraAudio} setExtraAudio={setWxExtraAudio} />
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
                updateSecrets(e, 'update', 'wx', wxServer, wxSecret, !wxEnabled, wxCustom, wxLabel, wxFiles, wxExtraAudio, () => {
                  setWxEnabled(!wxEnabled);
                });
              }}
            >
              {wxEnabled ? 'Stop Camera Live' : 'Start Camera Live'}
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
              <Form.Text> * {bilibiliCustom ? '(Required)' : '(Optional)'} IP camera stream label</Form.Text>
              <Form.Control as="input" defaultValue={bilibiliLabel} onChange={(e) => setBilibiliLabel(e.target.value)}/>
            </Form.Group>
            <SrsErrorBoundary>
              <ChooseVideoSourceEn platform='bilibili' cameraFiles={bilibiliFiles} setCameraFiles={setBilibiliFiles} />
              <CameraExtraAudioTrack extraAudio={bilibiliExtraAudio} setExtraAudio={setBilibiliExtraAudio} />
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
                updateSecrets(e, 'update', 'bilibili', bilibiliServer, bilibiliSecret, !bilibiliEnabled, bilibiliCustom, bilibiliLabel, bilibiliFiles, bilibiliExtraAudio, () => {
                  setBilibiliEnabled(!bilibiliEnabled);
                });
              }}
            >
              {bilibiliEnabled ? 'Stop Camera Live' : 'Start Camera Live'}
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
              <Form.Text> * {kuaishouCustom ? '(Required)' : '(Optional)'} IP camera stream label</Form.Text>
              <Form.Control as="input" defaultValue={kuaishouLabel} onChange={(e) => setKuaishouLabel(e.target.value)}/>
            </Form.Group>
            <SrsErrorBoundary>
              <ChooseVideoSourceEn platform='kuaishou' cameraFiles={kuaishouFiles} setCameraFiles={setKuaishouFiles} />
              <CameraExtraAudioTrack extraAudio={kuaishouExtraAudio} setExtraAudio={setKuaishouExtraAudio} />
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
                updateSecrets(e, 'update', 'kuaishou', kuaishouServer, kuaishouSecret, !kuaishouEnabled, kuaishouCustom, kuaishouLabel, kuaishouFiles, kuaishouExtraAudio, () => {
                  setKuaishouEnabled(!kuaishouEnabled);
                });
              }}
            >
              {kuaishouEnabled ? 'Stop Camera Live' : 'Start Camera Live'}
            </Button> &nbsp;
            <Form.Text> * Convert file to live stream</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="99">
        <Accordion.Header>Camera Live Status</Accordion.Header>
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
                  <th>Extra Audio</th>
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
                        <CameraFileFormatInfo file={file.sourceObj}/>
                      </td>
                      <td>{file?.extraAudio}</td>
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

function CameraFileList({files, onChangeFiles}) {
  const {t} = useTranslation();
  return (
    <Row>
      <Col xs='auto'>
        <ListGroup>
          {files.map((f, index) => {
            return <ListGroup.Item key={index}>
              {f.name} &nbsp;
              <CameraFileFormatInfo file={f}/> &nbsp;
              <CameraFileVideoInfo file={f}/> &nbsp;
              <CameraFileAudioInfo file={f}/>
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

function CameraExtraAudioTrack({extraAudio, setExtraAudio}) {
  const {t} = useTranslation();

  return <>
    <Form.Group className="mb-3">
      <Form.Label>{t('camera.silent2')}</Form.Label>
      <Form.Text> * {t('camera.silent3')}</Form.Text>
      <Form.Select defaultValue={extraAudio} onChange={(e) => setExtraAudio(e.target.value)}>
        <option value="">--{t('helper.noSelect')}--</option>
        <option value="silent">{t('camera.silent')}</option>
      </Form.Select>
    </Form.Group>
  </>;
}

function ChooseVideoSourceCn({platform, cameraFiles, setCameraFiles}) {
  return (<>
    <Form.Group className="mb-3">
      <Form.Label>视频源</Form.Label>
      <Form.Text> * 流地址支持 rtmp, http, https, 或 rtsp 等格式</Form.Text>
      <SrsErrorBoundary>
        <CameraStreamSelectorCn platform={platform} cameraFiles={cameraFiles} setCameraFiles={setCameraFiles}/>
      </SrsErrorBoundary>
    </Form.Group>
  </>);
}

function ChooseVideoSourceEn({platform, cameraFiles, setCameraFiles}) {
  return (<>
    <Form.Group className="mb-3">
      <Form.Label>Live Stream Source</Form.Label>
      <Form.Text> * The stream URL should start with rtmp, http, https, or rtsp.</Form.Text>
      <SrsErrorBoundary>
        <CameraStreamSelectorEn platform={platform} cameraFiles={cameraFiles} setCameraFiles={setCameraFiles}/>
      </SrsErrorBoundary>
    </Form.Group>
  </>);
}

function CameraStreamSelectorCn({platform, cameraFiles, setCameraFiles}) {
  const handleError = useErrorHandler();
  const [inputStream, setInputStream] = React.useState(cameraFiles?.length ? cameraFiles[0].target : '');
  const [submiting, setSubmiting] = React.useState();

  const checkStreamUrl = React.useCallback(async () => {
    if (!inputStream) return alert('请输入流地址');
    const isHTTP = inputStream.startsWith('http://') || inputStream.startsWith('https://');
    if (!inputStream.startsWith('rtmp://') && !inputStream.startsWith('rtsp://') && !isHTTP) return alert('流地址必须是 rtmp/http/https/rtsp 格式');
    if (isHTTP && inputStream.indexOf('.flv') < 0 && inputStream.indexOf('.m3u8') < 0) return alert('HTTP流必须是 http-flv或hls 格式');

    setSubmiting(true);
    try {
      const res = await new Promise((resolve, reject) => {
        axios.post(`/terraform/v1/ffmpeg/camera/stream-url`, {
          url: inputStream,
        }, {
          headers: Token.loadBearerHeader(),
        }).then(res => {
          resolve(res);
        }).catch(reject);
      });

      await new Promise((resolve, reject) => {
        console.log(`检查流地址成功，${JSON.stringify(res.data.data)}`);
        const streamObj = res.data.data;
        const files = [{name: streamObj.name, size: 0, uuid: streamObj.uuid, target: streamObj.target, type: "stream"}];
        axios.post('/terraform/v1/ffmpeg/camera/source', {
          platform, files,
        }, {
          headers: Token.loadBearerHeader(),
        }).then(res => {
          console.log(`更新虚拟直播源为流地址成功，${JSON.stringify(res.data.data)}`);
          setCameraFiles(res.data.data.files);
          resolve();
        }).catch(reject);
      });
    } catch (e) {
      handleError(e);
    } finally {
      setSubmiting(false);
    }
  }, [inputStream, handleError, platform, setCameraFiles, setSubmiting]);

  return (<>
    <Form.Control as="div">
      {!cameraFiles?.length && <>
        <Row>
          <Col>
            <Form.Control type="text" defaultValue={inputStream} placeholder="请输入流地址" onChange={e => setInputStream(e.target.value)} />
          </Col>
          <Col xs="auto">
            <Button variant="primary" disabled={submiting} onClick={checkStreamUrl}>确认</Button>
          </Col>
        </Row></>
      }
      {cameraFiles?.length && <CameraFileList files={cameraFiles} onChangeFiles={(e) => setCameraFiles(null)}/>}
    </Form.Control>
  </>);
}

function CameraStreamSelectorEn({platform, cameraFiles, setCameraFiles}) {
  const handleError = useErrorHandler();
  const [inputStream, setInputStream] = React.useState(cameraFiles?.length ? cameraFiles[0].target : '');
  const [submiting, setSubmiting] = React.useState();

  const checkStreamUrl = React.useCallback(async () => {
    if (!inputStream) return alert('Please input stream URL');
    const isHTTP = inputStream.startsWith('http://') || inputStream.startsWith('https://');
    if (!inputStream.startsWith('rtmp://') && !inputStream.startsWith('rtsp://') && !isHTTP) return alert('The stream must be rtmp/http/https/rtsp');
    if (isHTTP && inputStream.indexOf('.flv') < 0 && inputStream.indexOf('.m3u8') < 0) return alert('The HTTP stream must be http-flv/hls');

    setSubmiting(true);
    try {
      const res = await new Promise((resolve, reject) => {
        axios.post(`/terraform/v1/ffmpeg/camera/stream-url`, {
          url: inputStream,
        }, {
          headers: Token.loadBearerHeader(),
        }).then(res => {
          resolve(res);
        }).catch(reject);
      });

      await new Promise((resolve, reject) => {
        console.log(`Check stream url ok，${JSON.stringify(res.data.data)}`);
        const streamObj = res.data.data;
        const files = [{name: streamObj.name, size: 0, uuid: streamObj.uuid, target: streamObj.target, type: "stream"}];
        axios.post('/terraform/v1/ffmpeg/camera/source', {
          platform, files,
        }, {
          headers: Token.loadBearerHeader(),
        }).then(res => {
          console.log(`Setup the virtual live stream ok，${JSON.stringify(res.data.data)}`);
          setCameraFiles(res.data.data.files);
          resolve();
        }).catch(reject);
      });
    } catch (e) {
      handleError(e);
    } finally {
      setSubmiting(false);
    }
  }, [inputStream, handleError, platform, setCameraFiles, setSubmiting]);

  return (<>
    <Form.Control as="div">
      {!cameraFiles?.length && <>
        <Row>
          <Col>
            <Form.Control type="text" defaultValue={inputStream} placeholder="please input stream URL" onChange={e => setInputStream(e.target.value)} />
          </Col>
          <Col xs="auto">
            <Button variant="primary" disabled={submiting} onClick={checkStreamUrl}>Submit</Button>
          </Col>
        </Row></>
      }
      {cameraFiles?.length && <CameraFileList files={cameraFiles} onChangeFiles={(e) => setCameraFiles(null)}/>}
    </Form.Control>
  </>)
}

function CameraFileFormatInfo({file}) {
  const f = file;
  if (!f?.format) return <></>;
  return <>
    {f?.type !== 'stream' &&
      <>
        File &nbsp;
        {Number(f?.size/1024/1024).toFixed(1)}MB &nbsp;
        {Number(f?.format?.duration).toFixed(0)}s &nbsp;
      </>
    }
    {f?.type === 'stream' &&
      <>
        Stream &nbsp;
      </>
    }
    {Number(f?.format?.bit_rate/1000).toFixed(1)}Kbps
  </>;
}

function CameraFileVideoInfo({file}) {
  const f = file;
  if (!f?.video) return <>NoVideo</>;
  return <>Video({f?.video?.codec_name} {f?.video?.profile} {f?.video?.width}x{f?.video?.height})</>;
}

function CameraFileAudioInfo({file}) {
  const f = file;
  if (!f?.audio) return <>NoAudio</>;
  return <>Audio({f?.audio?.codec_name} {f?.audio?.sample_rate}HZ {f?.audio?.channels}CH)</>;
}

