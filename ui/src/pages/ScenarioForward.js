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
import {SrsEnvContext} from "../components/SrsEnvContext";

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

    setActiveKey('1');
    for (const key in secrets) {
      const e = secrets[key];
      if (e.enabled && e.server && e.secret) {
        setActiveKey('99');
      }
    }
  }, [init, secrets]);

  if (!activeKey) return <></>;
  return <ScenarioForwardImpl defaultActiveKey={activeKey} defaultSecrets={secrets}/>;
}

function ScenarioForwardImpl({defaultActiveKey, defaultSecrets}) {
  const language = useSrsLanguage();
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const env = React.useContext(SrsEnvContext)[0];

  const [configs, setConfigs] = React.useState([]);
  const [forwards, setForwards] = React.useState();
  const [submiting, setSubmiting] = React.useState();

  // Convert default config from kv to objects in array.
  React.useEffect(() => {
    if (!defaultSecrets) return;

    let index = 1;
    const confs = [{
      platform: 'wx', enabled: false, index: String(index++), allowCustom: true,
      ...defaultSecrets?.wx,
      locale: {
        label: null, link: t('plat.wx.link'), link2: t('plat.wx.link2'),
        generate: (e) => {
          e.locale.label = e.custom ? t('plat.com.custom') : t('plat.wx.title');
        },
      },
    }, {
      platform: 'bilibili', enabled: false, index: String(index++), allowCustom: true,
      ...defaultSecrets?.bilibili,
      locale: {
        label: null, link: t('plat.bl.link'), link2: t('plat.bl.link2'),
        generate: (e) => {
          e.locale.label = e.custom ? t('plat.com.custom') : t('plat.bl.title');
        },
      },
    }, {
      platform: 'kuaishou', enabled: false, index: String(index++), allowCustom: true, custom: language === 'zh',
      ...defaultSecrets?.kuaishou,
      locale: {
        label: null, link: t('plat.ks.link'), link2: t('plat.ks.link2'),
        generate: (e) => {
          e.locale.label = e.custom ? t('plat.com.custom') : t('plat.ks.title');
        },
      },
    }];

    // Regenerate the locale label, because it may change after created from defaults.
    confs.forEach((e) => {
      e?.locale?.generate && e.locale.generate(e);
    });

    // Generate more forwarding configures.
    while (confs.length < env.forwardLimit) {
      const rindex = index++;
      const rid = Math.random().toString(16).slice(-6);

      // Load the configured forwarding from defaults.
      const existsConf = Object.values(defaultSecrets).find(e => e.platform.indexOf(`forwarding-${rindex}-`) === 0);
      if (existsConf) {
        confs.push(existsConf);
      } else {
        confs.push({
          platform: `forwarding-${rindex}-${rid}`, enabled: false, index: String(rindex), allowCustom: false,
          server: null, secret: null, custom: true, label: `Forwarding #${rindex}`,
        });
      }
    }

    setConfigs(confs);
    console.log(`Forward: Init configs ${JSON.stringify(confs)}`);
  }, [defaultSecrets, setConfigs, env, language, t]);

  // Fetch the forwarding streams from server.
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

  // Update config object in array.
  const updateConfigObject = React.useCallback((conf) => {
    const confs = configs.map((e) => {
      if (e.platform === conf.platform) {
        return conf;
      }
      return e;
    })
    setConfigs(confs);
    console.log(`Forward: Update config ${JSON.stringify(conf)} to ${JSON.stringify(confs)}`);
  }, [configs, setConfigs]);

  // Update the forward config to server.
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
      {configs.map((conf) => {
        return (
          <Accordion.Item eventKey={conf.index} key={conf.platform}>
            <Accordion.Header>{conf?.locale?.label} {conf.label}</Accordion.Header>
            <Accordion.Body>
              <Form>
                <Form.Group className="mb-3">
                  <Form.Label>{t('plat.com.name')}</Form.Label>
                  <Form.Text> * {conf.custom ? `(${t('helper.required')})` : `(${t('helper.optional')})`} {t('plat.com.name2')}</Form.Text>
                  <Form.Control as="input" defaultValue={conf.label} onChange={(e) => updateConfigObject({...conf, label: e.target.value})}/>
                </Form.Group>
                <Form.Group className="mb-3">
                  <Form.Label>{conf.custom ? t('plat.com.server') : t('plat.com.server2')}</Form.Label>
                  {!conf.custom && <Form.Text> * {t('plat.com.server3')} <a href={conf?.locale?.link} target='_blank' rel='noreferrer'>{conf?.locale?.link2}</a>, {t('plat.com.server4')}</Form.Text>}
                  <Form.Control as="input" defaultValue={conf.server} onChange={(e) => updateConfigObject({...conf, server: e.target.value})}/>
                </Form.Group>
                <Form.Group className="mb-3">
                  <Form.Label>{t('plat.com.key')}</Form.Label>
                  {!conf.custom && <Form.Text> * {t('plat.com.server3')} <a href={conf?.locale?.link} target='_blank' rel='noreferrer'>{conf?.locale?.link2}</a>, {t('plat.com.key2')}</Form.Text>}
                  <Form.Control as="input" defaultValue={conf.secret} onChange={(e) => updateConfigObject({...conf, secret: e.target.value})}/>
                </Form.Group>
                {conf?.allowCustom && (
                  <Row>
                    <Col xs='auto'>
                      <Form.Group className="mb-3" controlId={`formCustomCheckbox-${conf.platform}`}>
                        <Form.Check type="checkbox" label={t('plat.com.custom')} defaultChecked={conf.custom} onClick={() => updateConfigObject({...conf, custom: !conf.custom})} />
                      </Form.Group>
                    </Col>
                  </Row>
                )}
                <Button
                  variant="primary"
                  type="submit"
                  disabled={submiting}
                  onClick={(e) => {
                    updateSecrets(e, 'update', conf.platform, conf.server, conf.secret, !conf.enabled, conf.custom, conf.label, () => {
                      updateConfigObject({...conf, enabled: !conf.enabled});
                    });
                  }}
                >
                  {conf.enabled ? t('plat.com.stop') : t('plat.com.start')}
                </Button> &nbsp;
                <Form.Text> * {t('forward.tip')}</Form.Text>
              </Form>
            </Accordion.Body>
          </Accordion.Item>
        );
      })}
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
