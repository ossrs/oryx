//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import {Accordion, Badge, Button, Col, Form, Row, Table} from "react-bootstrap";
import React from "react";
import {Token} from "../utils";
import axios from "axios";
import moment from "moment";
import {TutorialsButton, useTutorials} from "../components/TutorialsButton";
import {useErrorHandler} from "react-error-boundary";
import {useSrsLanguage} from "../components/LanguageSwitch";

export default function ScenarioForward() {
  const [init, setInit] = React.useState();
  const [activeKey, setActiveKey] = React.useState();
  const [secrets, setSecrets] = React.useState();
  const handleError = useErrorHandler();
  const language = useSrsLanguage();

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/ffmpeg/forward/secret', {
      ...token,
    }).then(res => {
      const secrets = res.data.data;
      setInit(true);
      setSecrets(secrets || {});
      console.log(`Forward: Secret query ok ${JSON.stringify(secrets)}`);
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
    return <ScenarioForwardImplCn defaultActiveKey={activeKey} defaultSecrets={secrets}/>;
  }
  return <ScenarioForwardImplEn defaultActiveKey={activeKey} defaultSecrets={secrets}/>;
}

function ScenarioForwardImplCn({defaultActiveKey, defaultSecrets}) {
  const [wxEnabled, setWxEnabled] = React.useState(defaultSecrets?.wx?.enabled);
  const [wxServer, setWxServer] = React.useState(defaultSecrets?.wx?.server);
  const [wxSecret, setWxSecret] = React.useState(defaultSecrets?.wx?.secret);
  const [wxCustom, setWxCustom] = React.useState(defaultSecrets?.wx?.custom);
  const [wxLabel, setWxLabel] = React.useState(defaultSecrets?.wx?.label);
  const [bilibiliEnabled, setBilibiliEnabled] = React.useState(defaultSecrets?.bilibili?.enabled);
  const [bilibiliServer, setBilibiliServer] = React.useState(defaultSecrets?.bilibili?.server);
  const [bilibiliSecret, setBilibiliSecret] = React.useState(defaultSecrets?.bilibili?.secret);
  const [bilibiliCustom, setBilibiliCustom] = React.useState(defaultSecrets?.bilibili?.custom);
  const [bilibiliLabel, setBilibiliLabel] = React.useState(defaultSecrets?.bilibili?.label);
  const [kuaishouEnabled, setKuaishouEnabled] = React.useState(defaultSecrets?.kuaishou?.enabled);
  const [kuaishouServer, setKuaishouServer] = React.useState(defaultSecrets?.kuaishou?.server);
  const [kuaishouSecret, setKuaishouSecret] = React.useState(defaultSecrets?.kuaishou?.secret);
  const [kuaishouCustom, setKuaishouCustom] = React.useState(defaultSecrets?.kuaishou?.custom || true);
  const [kuaishouLabel, setKuaishouLabel] = React.useState(defaultSecrets?.kuaishou?.label);
  const [forwards, setForwards] = React.useState();
  const [submiting, setSubmiting] = React.useState();
  const handleError = useErrorHandler();

  const forwardTutorials = useTutorials({
    bilibili: React.useRef([
      {author: 'SRS', id: 'BV1KY411V7uc'},
    ])
  });

  React.useEffect(() => {
    const refreshStreams = () => {
      const token = Token.load();
      axios.post('/terraform/v1/ffmpeg/forward/streams', {
        ...token,
      }).then(res => {
        setForwards(res.data.data.map((e, i) => ({
          ...e,
          name: {wx: '视频号直播', bilibili: 'Bilibili直播间', kuaishou: '快手云直播'}[e.platform],
          update: e.frame?.update ? moment(e.frame.update) : null,
          i,
        })));
        console.log(`Forward: Query streams ${JSON.stringify(res.data.data)}`);
      }).catch(handleError);
    };

    refreshStreams();
    const timer = setInterval(() => refreshStreams(), 10 * 1000);
    return () => clearInterval(timer);
  }, [handleError]);

  const updateSecrets = React.useCallback((e, action, platform, server, secret, enabled, custom, label, onSuccess) => {
    e.preventDefault();
    if (!server) return alert('请输入推流地址');
    if (custom && !label) return alert('自定义平台请输入名称，否则不好区分转播状态');

    try {
      setSubmiting(true);

      const token = Token.load();
      axios.post('/terraform/v1/ffmpeg/forward/secret', {
        ...token, action, platform, server, secret, enabled: !!enabled, custom: !!custom, label,
      }).then(res => {
        alert('转播设置成功');
        onSuccess && onSuccess();
      }).catch(handleError);
    } finally {
      new Promise(resolve => setTimeout(resolve, 3000)).then(() => setSubmiting(false));
    }
  }, [handleError, setSubmiting]);

  return (
    <Accordion defaultActiveKey={[defaultActiveKey]} alwaysOpen>
      <Accordion.Item eventKey="0">
        <Accordion.Header>场景介绍</Accordion.Header>
        <Accordion.Body>
          <div>
            多平台转播<TutorialsButton prefixLine={true} tutorials={forwardTutorials} />，将流转播给其他平台，比如视频号直播、快手、B站等。
            <p></p>
          </div>
          <p>可应用的具体场景包括：</p>
          <ul>
            <li>节约上行带宽，避免客户端推多路流，服务器转播更有保障</li>
          </ul>
          <p>使用说明：</p>
          <ul>
            <li>首先使用适合你的场景推流</li>
            <li>然后设置转播的平台</li>
          </ul>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{wxCustom ? '自定义平台' : '视频号直播'} {wxLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>名称</Form.Label>
              <Form.Text> * {wxCustom ? '(必选)' : '(可选)'} 起一个好记的名字</Form.Text>
              <Form.Control as="input" defaultValue={wxLabel} onChange={(e) => setWxLabel(e.target.value)}/>
            </Form.Group>
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
                updateSecrets(e, 'update', 'wx', wxServer, wxSecret, !wxEnabled, wxCustom, wxLabel, () => {
                  setWxEnabled(!wxEnabled);
                });
              }}
            >
              {wxEnabled ? '停止转播' : '开始转播'}
            </Button> &nbsp;
            <TutorialsButton prefixLine={true} tutorials={forwardTutorials} /> &nbsp;
            <Form.Text> * 若有多个流，随机选择一个</Form.Text>
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
                updateSecrets(e, 'update', 'bilibili', bilibiliServer, bilibiliSecret, !bilibiliEnabled, bilibiliCustom, bilibiliLabel, () => {
                  setBilibiliEnabled(!bilibiliEnabled);
                });
              }}
            >
              {bilibiliEnabled ? '停止转播' : '开始转播'}
            </Button> &nbsp;
            <TutorialsButton prefixLine={true} tutorials={forwardTutorials} /> &nbsp;
            <Form.Text> * 若有多个流，随机选择一个</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>{kuaishouCustom ? '自定义平台' : '快手云直播'} {kuaishouLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>名称</Form.Label>
              <Form.Text> * {kuaishouCustom ? '(必选)' : '(可选)'} 起一个好记的名字</Form.Text>
              <Form.Control as="input" defaultValue={kuaishouLabel} onChange={(e) => setKuaishouLabel(e.target.value)}/>
            </Form.Group>
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
                updateSecrets(e, 'update', 'kuaishou', kuaishouServer, kuaishouSecret, !kuaishouEnabled, kuaishouCustom, kuaishouLabel, () => {
                  setKuaishouEnabled(!kuaishouEnabled);
                });
              }}
            >
              {kuaishouEnabled ? '停止转播' : '开始转播'}
            </Button> &nbsp;
            <TutorialsButton prefixLine={true} tutorials={forwardTutorials} /> &nbsp;
            <Form.Text> * 若有多个流，随机选择一个</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="99">
        <Accordion.Header>转播状态</Accordion.Header>
        <Accordion.Body>
          {
            forwards?.length ? (
              <Table striped bordered hover>
                <thead>
                <tr>
                  <th>#</th>
                  <th>平台</th>
                  <th>状态</th>
                  <th>更新时间</th>
                  <th>转播流</th>
                  <th>日志</th>
                </tr>
                </thead>
                <tbody>
                {
                  forwards?.map(file => {
                    return <tr key={file.platform} style={{verticalAlign: 'middle'}}>
                      <td>{file.i}</td>
                      <td>{file.custom ? (file.label ? '' : '自定义平台') : file.name} {file.label}</td>
                      <td>
                        <Badge bg={file.enabled ? (file.frame ? 'success' : 'primary') : 'secondary'}>
                          {file.enabled ? (file.frame ? '转播中' : '等待中') : '未开启'}
                        </Badge>
                      </td>
                      <td>{file.update && `${file.update?.format('YYYY-MM-DD HH:mm:ss')}`}</td>
                      <td>{file.stream}</td>
                      <td>{file.frame?.log}</td>
                    </tr>;
                  })
                }
                </tbody>
              </Table>
            ) : ''
          }
          {!forwards?.length ? '没有流。请开启转播并推流后，等待大约30秒左右，转播列表会自动更新' : ''}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

function ScenarioForwardImplEn({defaultActiveKey, defaultSecrets}) {
  const [wxEnabled, setWxEnabled] = React.useState(defaultSecrets?.wx?.enabled);
  const [wxServer, setWxServer] = React.useState(defaultSecrets?.wx?.server);
  const [wxSecret, setWxSecret] = React.useState(defaultSecrets?.wx?.secret);
  const [wxCustom, setWxCustom] = React.useState(defaultSecrets?.wx?.custom);
  const [wxLabel, setWxLabel] = React.useState(defaultSecrets?.wx?.label);
  const [bilibiliEnabled, setBilibiliEnabled] = React.useState(defaultSecrets?.bilibili?.enabled);
  const [bilibiliServer, setBilibiliServer] = React.useState(defaultSecrets?.bilibili?.server || 'rtmp://live.twitch.tv/app');
  const [bilibiliSecret, setBilibiliSecret] = React.useState(defaultSecrets?.bilibili?.secret);
  const [bilibiliCustom, setBilibiliCustom] = React.useState(defaultSecrets?.bilibili?.custom);
  const [bilibiliLabel, setBilibiliLabel] = React.useState(defaultSecrets?.bilibili?.label);
  const [kuaishouEnabled, setKuaishouEnabled] = React.useState(defaultSecrets?.kuaishou?.enabled);
  const [kuaishouServer, setKuaishouServer] = React.useState(defaultSecrets?.kuaishou?.server || 'rtmps://live-api-s.facebook.com:443/rtmp');
  const [kuaishouSecret, setKuaishouSecret] = React.useState(defaultSecrets?.kuaishou?.secret);
  const [kuaishouCustom, setKuaishouCustom] = React.useState(defaultSecrets?.kuaishou?.custom);
  const [kuaishouLabel, setKuaishouLabel] = React.useState(defaultSecrets?.kuaishou?.label);
  const [forwards, setForwards] = React.useState();
  const [submiting, setSubmiting] = React.useState();
  const handleError = useErrorHandler();

  React.useEffect(() => {
    const refreshStreams = () => {
      const token = Token.load();
      axios.post('/terraform/v1/ffmpeg/forward/streams', {
        ...token,
      }).then(res => {
        setForwards(res.data.data.map((e, i) => ({
          ...e,
          name: {wx: 'YouTube', bilibili: 'Twitch', kuaishou: 'Facebook'}[e.platform],
          update: e.frame?.update ? moment(e.frame.update) : null,
          i,
        })));
        console.log(`Forward: Query streams ${JSON.stringify(res.data.data)}`);
      }).catch(handleError);
    };

    refreshStreams();
    const timer = setInterval(() => refreshStreams(), 10 * 1000);
    return () => clearInterval(timer);
  }, [handleError]);

  const updateSecrets = React.useCallback((e, action, platform, server, secret, enabled, custom, label, onSuccess) => {
    e.preventDefault();
    if (!server) return alert('Please input streaming server URL');
    if (custom && !label) return alert('Please enter a name for the custom platform, to distinguish the streaming status.');

    try {
      setSubmiting(true);

      const token = Token.load();
      axios.post('/terraform/v1/ffmpeg/forward/secret', {
        ...token, action, platform, server, secret, enabled: !!enabled, custom: !!custom, label,
      }).then(res => {
        alert('Setup OK');
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
            Multi-platform streaming, forward to other platforms, such as YouTube, Twitch, TikTok, etc.
            <p></p>
          </div>
          <p>Specific application scenarios include:</p>
          <ul>
            <li>Save on upstream bandwidth, avoid pushing multiple streams from the client, and server streaming is more secure</li>
          </ul>
          <p>Usage instructions:</p>
          <ul>
            <li>First, use the appropriate streaming method for your scenario</li>
            <li>Then set the platform for streaming</li>
          </ul>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{wxCustom ? 'Custom' : 'YouTube'} {wxLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>Label</Form.Label>
              <Form.Text> * {wxCustom ? '(Required)' : '(Optional)'} Stream event label</Form.Text>
              <Form.Control as="input" defaultValue={wxLabel} onChange={(e) => setWxLabel(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{wxCustom ? 'Server' : 'Stream URL'}</Form.Label>
              {!wxCustom && <Form.Text> * Please click <a href='https://studio.youtube.com/channel/UC/livestreaming' target='_blank' rel='noreferrer'>Go live</a>, then copy the Stream URL</Form.Text>}
              <Form.Control as="input" defaultValue={wxServer} onChange={(e) => setWxServer(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Stream Key</Form.Label>
              {!wxCustom && <Form.Text> * Please click <a href='https://channels.weixin.qq.com/platform/live/liveBuild' target='_blank' rel='noreferrer'>Go live</a>, then copy the Stream Key</Form.Text>}
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
                updateSecrets(e, 'update', 'wx', wxServer, wxSecret, !wxEnabled, wxCustom, wxLabel, () => {
                  setWxEnabled(!wxEnabled);
                });
              }}
            >
              {wxEnabled ? 'Stop Forward' : 'Start Forward'}
            </Button> &nbsp;
            <Form.Text> * If there are multiple streams, randomly select one.</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>{bilibiliCustom ? 'Custom' : 'Twitch'} {bilibiliLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>Label</Form.Label>
              <Form.Text> * {bilibiliCustom ? '(Required)' : '(Optional)'} Stream event label</Form.Text>
              <Form.Control as="input" defaultValue={bilibiliLabel} onChange={(e) => setBilibiliLabel(e.target.value)}/>
            </Form.Group>
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
                updateSecrets(e, 'update', 'bilibili', bilibiliServer, bilibiliSecret, !bilibiliEnabled, bilibiliCustom, bilibiliLabel, () => {
                  setBilibiliEnabled(!bilibiliEnabled);
                });
              }}
            >
              {bilibiliEnabled ? 'Stop Forward' : 'Start Forward'}
            </Button> &nbsp;
            <Form.Text> * If there are multiple streams, randomly select one.</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>{kuaishouCustom ? 'Custom' : 'Facebook'} {kuaishouLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>Label</Form.Label>
              <Form.Text> * {kuaishouCustom ? '(Required)' : '(Optional)'} Stream event label</Form.Text>
              <Form.Control as="input" defaultValue={kuaishouLabel} onChange={(e) => setKuaishouLabel(e.target.value)}/>
            </Form.Group>
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
                updateSecrets(e, 'update', 'kuaishou', kuaishouServer, kuaishouSecret, !kuaishouEnabled, kuaishouCustom, kuaishouLabel, () => {
                  setKuaishouEnabled(!kuaishouEnabled);
                });
              }}
            >
              {kuaishouEnabled ? 'Stop Forward' : 'Start Forward'}
            </Button> &nbsp;
            <Form.Text> * If there are multiple streams, randomly select one.</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="99">
        <Accordion.Header>Forward Status</Accordion.Header>
        <Accordion.Body>
          {
            forwards?.length ? (
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
                  forwards?.map(file => {
                    return <tr key={file.platform} style={{verticalAlign: 'middle'}}>
                      <td>{file.i}</td>
                      <td>{file.custom ? (file.label ? '' : 'Custom') : file.name} {file.label}</td>
                      <td>
                        <Badge bg={file.enabled ? (file.frame ? 'success' : 'primary') : 'secondary'}>
                          {file.enabled ? (file.frame ? 'Forwarding' : 'Waiting') : 'Inactive'}
                        </Badge>
                      </td>
                      <td>{file.update && `${file.update?.format('YYYY-MM-DD HH:mm:ss')}`}</td>
                      <td>{file.stream}</td>
                      <td>{file.frame?.log}</td>
                    </tr>;
                  })
                }
                </tbody>
              </Table>
            ) : ''
          }
          {!forwards?.length ? 'No stream. Please turn on the broadcast and push the stream, then wait for about 30 seconds, and the broadcast list will automatically update.' : ''}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

