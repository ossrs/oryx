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
import {useErrorHandler} from "react-error-boundary";
import {useSrsLanguage} from "../components/LanguageSwitch";
import {useTranslation} from "react-i18next";

export default function ScenarioForward() {
  const [init, setInit] = React.useState();
  const [activeKey, setActiveKey] = React.useState();
  const [secrets, setSecrets] = React.useState();
  const handleError = useErrorHandler();

  React.useEffect(() => {
    axios.post('/terraform/v1/ffmpeg/forward/secret', {
    }, {
      headers: Token.loadBearerHeader(),
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
  return <ScenarioForwardImpl defaultActiveKey={activeKey} defaultSecrets={secrets}/>;
}

function ScenarioForwardImpl({defaultActiveKey, defaultSecrets}) {
  const language = useSrsLanguage();
  const {t} = useTranslation();
  const handleError = useErrorHandler();

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

  React.useEffect(() => {
    const refreshStreams = () => {
      axios.post('/terraform/v1/ffmpeg/forward/streams', {
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        setForwards(res.data.data.map((e, i) => ({
          ...e,
          name: {
            wx: t('plat.wx.title'),
            bilibili: t('plat.bl.title'),
            kuaishou: t('plat.ks.title')
          }[e.platform],
          start: e.start ? moment(e.start) : null,
          ready: e.ready ? moment(e.ready) : null,
          update: e.frame?.update ? moment(e.frame.update) : null,
          i,
        })));
        console.log(`Forward: Query streams ${JSON.stringify(res.data.data)}`);
      }).catch(handleError);
    };

    refreshStreams();
    const timer = setInterval(() => refreshStreams(), 10 * 1000);
    return () => clearInterval(timer);
  }, [t, handleError, setForwards]);

  const updateSecrets = React.useCallback((e, action, platform, server, secret, enabled, custom, label, onSuccess) => {
    e.preventDefault();
    if (!server) return alert(t('plat.com.addr'));
    if (custom && !label) return alert(t('plat.com.label'));

    try {
      setSubmiting(true);

      axios.post('/terraform/v1/ffmpeg/forward/secret', {
        action, platform, server, secret, enabled: !!enabled, custom: !!custom, label,
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        alert(t('plat.com.ok'));
        onSuccess && onSuccess();
      }).catch(handleError);
    } finally {
      new Promise(resolve => setTimeout(resolve, 3000)).then(() => setSubmiting(false));
    }
  }, [t, handleError, setSubmiting]);

  return (
    <Accordion defaultActiveKey={[defaultActiveKey]}>
      <React.Fragment>
        {language === 'zh' ?
          <Accordion.Item eventKey="0">
            <Accordion.Header>场景介绍</Accordion.Header>
            <Accordion.Body>
              <div>
                多平台转播，将流转播给其他平台，比如视频号直播、快手、B站等。
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
          </Accordion.Item> :
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
          </Accordion.Item>}
      </React.Fragment>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{wxCustom ? t('plat.com.custom') : t('plat.wx.title')} {wxLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>{t('plat.com.name')}</Form.Label>
              <Form.Text> * {wxCustom ? `(${t('helper.required')})` : `(${t('helper.optional')})`} {t('plat.com.name2')}</Form.Text>
              <Form.Control as="input" defaultValue={wxLabel} onChange={(e) => setWxLabel(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{wxCustom ? t('plat.com.server') : t('plat.com.server2')}</Form.Label>
              {!wxCustom && <Form.Text> * {t('plat.com.server3')} <a href={t('plat.wx.link')} target='_blank' rel='noreferrer'>{t('plat.wx.link2')}</a>, {t('plat.com.server4')}</Form.Text>}
              <Form.Control as="input" defaultValue={wxServer} onChange={(e) => setWxServer(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{t('plat.com.key')}</Form.Label>
              {!wxCustom && <Form.Text> * {t('plat.com.server3')} <a href={t('plat.wx.link')} target='_blank' rel='noreferrer'>{t('plat.wx.link2')}</a>, {t('plat.com.key2')}</Form.Text>}
              <Form.Control as="input" defaultValue={wxSecret} onChange={(e) => setWxSecret(e.target.value)}/>
            </Form.Group>
            <Row>
              <Col xs='auto'>
                <Form.Group className="mb-3" controlId="formWxCustomCheckbox">
                  <Form.Check type="checkbox" label={t('plat.com.custom')} defaultChecked={wxCustom} onClick={() => setWxCustom(!wxCustom)} />
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
              {wxEnabled ? t('plat.com.stop') : t('plat.com.start')}
            </Button> &nbsp;
            <Form.Text> * {t('forward.tip')}</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>{bilibiliCustom ? t('plat.com.custom') : t('plat.bl.title')} {bilibiliLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>{t('plat.com.name')}</Form.Label>
              <Form.Text> * {bilibiliCustom ? `(${t('helper.required')})` : `(${t('helper.optional')})`} {t('plat.com.name2')}</Form.Text>
              <Form.Control as="input" defaultValue={bilibiliLabel} onChange={(e) => setBilibiliLabel(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{bilibiliCustom ? t('plat.com.server') : t('plat.com.server2')}</Form.Label>
              {!bilibiliCustom && <Form.Text> * {t('plat.com.server3')} <a href={t('plat.bl.link')} target='_blank' rel='noreferrer'>{t('plat.bl.link2')}</a>, {t('plat.com.server4')}</Form.Text>}
              <Form.Control as="input" defaultValue={bilibiliServer} onChange={(e) => setBilibiliServer(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{t('plat.com.key')}</Form.Label>
              {!bilibiliCustom && <Form.Text> * {t('plat.com.server3')} <a href={t('plat.bl.link')} target='_blank' rel='noreferrer'>{t('plat.bl.link2')}</a>, {t('plat.com.key2')}</Form.Text>}
              <Form.Control as="input" defaultValue={bilibiliSecret} onChange={(e) => setBilibiliSecret(e.target.value)}/>
            </Form.Group>
            <Row>
              <Col xs='auto'>
                <Form.Group className="mb-3" controlId="formBilibiliCustomCheckbox">
                  <Form.Check type="checkbox" label={t('plat.com.custom')} defaultChecked={bilibiliCustom} onClick={() => setBilibiliCustom(!bilibiliCustom)} />
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
              {bilibiliEnabled ? t('plat.com.stop') : t('plat.com.start')}
            </Button> &nbsp;
            <Form.Text> * {t('forward.tip')}</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>{kuaishouCustom ? t('plat.com.custom') : t('plat.ks.title')} {kuaishouLabel}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>{t('plat.com.name')}</Form.Label>
              <Form.Text> * {kuaishouCustom ? `(${t('helper.required')})` : `(${t('helper.optional')})`} {t('plat.com.name2')}</Form.Text>
              <Form.Control as="input" defaultValue={kuaishouLabel} onChange={(e) => setKuaishouLabel(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{kuaishouCustom ? t('plat.com.server') : t('plat.com.server2')}</Form.Label>
              {!kuaishouCustom && <Form.Text> * {t('plat.com.server3')} <a href={t('plat.ks.link')} target='_blank' rel='noreferrer'>{t('plat.ks.link2')}</a>, {t('plat.com.server4')}</Form.Text>}
              <Form.Control as="input" defaultValue={kuaishouServer} onChange={(e) => setKuaishouServer(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{t('plat.com.key')}</Form.Label>
              {!kuaishouCustom && <Form.Text> * {t('plat.com.server3')} <a href={t('plat.ks.link')} target='_blank' rel='noreferrer'>{t('plat.ks.link2')}</a>, {t('plat.com.key2')}</Form.Text>}
              <Form.Control as="input" defaultValue={kuaishouSecret} onChange={(e) => setKuaishouSecret(e.target.value)}/>
            </Form.Group>
            <Row>
              <Col xs='auto'>
                <Form.Group className="mb-3" controlId="formKuaishouCustomCheckbox">
                  <Form.Check type="checkbox" label={t('plat.com.custom')} defaultChecked={kuaishouCustom} onClick={() => setKuaishouCustom(!kuaishouCustom)} />
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
              {kuaishouEnabled ? t('plat.com.stop') : t('plat.com.start')}
            </Button> &nbsp;
            <Form.Text> * {t('forward.tip')}</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="99">
        <Accordion.Header>{t('plat.com.status')}</Accordion.Header>
        <Accordion.Body>
          {
            forwards?.length ? (
              <Table striped bordered hover>
                <thead>
                <tr>
                  <th>#</th>
                  <th>{t('plat.com.platform')}</th>
                  <th>{t('plat.com.status2')}</th>
                  <th>Start</th>
                  <th>Ready</th>
                  <th>{t('plat.com.update')}</th>
                  <th>{t('plat.com.source')}</th>
                  <th>{t('plat.com.log')}</th>
                </tr>
                </thead>
                <tbody>
                {
                  forwards?.map(file => {
                    return <tr key={file.platform} style={{verticalAlign: 'middle'}}>
                      <td>{file.i}</td>
                      <td>{file.custom ? (file.label ? '' : t('plat.com.custom')) : file.name} {file.label}</td>
                      <td>
                        <Badge bg={file.enabled ? (file.frame ? 'success' : 'primary') : 'secondary'}>
                          {file.enabled ? (file.frame ? t('plat.com.s0') : t('plat.com.s1')) : t('plat.com.s2')}
                        </Badge>
                      </td>
                      <td>{file.start && `${file.start?.format('YYYY-MM-DD HH:mm:ss')}`}</td>
                      <td>{file.ready && `${file.ready?.format('YYYY-MM-DD HH:mm:ss')}`}</td>
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
          {!forwards?.length ? t('forward.s3') : ''}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}
