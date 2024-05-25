//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {Accordion, Badge, Button, Col, Form, Row, Table} from "react-bootstrap";
import {Token} from "../utils";
import axios from "axios";
import moment from "moment";
import {useErrorHandler} from "react-error-boundary";
import {useSrsLanguage} from "../components/LanguageSwitch";
import {useTranslation} from "react-i18next";
import {SrsErrorBoundary} from "../components/SrsErrorBoundary";
import ChooseVideoSource, {VLiveFileFormatInfo} from "../components/VideoSourceSelector";
import {SrsEnvContext} from "../components/SrsEnvContext";

export default function ScenarioVLive() {
  const [init, setInit] = React.useState();
  const [activeKey, setActiveKey] = React.useState();
  const [secrets, setSecrets] = React.useState();
  const handleError = useErrorHandler();

  React.useEffect(() => {
    axios.post('/terraform/v1/ffmpeg/vlive/secret', {
    }, {
      headers: Token.loadBearerHeader(),
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

    setActiveKey('1');
    for (const key in secrets) {
      const e = secrets[key];
      if (e.enabled && e.server && e.secret) {
        setActiveKey('99');
      }
    }
  }, [init, secrets]);

  if (!activeKey) return <></>;
  return <ScenarioVLiveImpl defaultActiveKey={activeKey} defaultSecrets={secrets}/>;
}

function ScenarioVLiveImpl({defaultActiveKey, defaultSecrets}) {
  const language = useSrsLanguage();
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const env = React.useContext(SrsEnvContext)[0];

  const [vLives, setVLives] = React.useState();
  const [submiting, setSubmiting] = React.useState();
  const [configs, setConfigs] = React.useState([]);

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

    // Generate more virtual live configures.
    while (confs.length < env.vLiveLimit) {
      const rindex = index++;
      const rid = Math.random().toString(16).slice(-6);

      // Load the configured virtual live from defaults.
      const existsConf = Object.values(defaultSecrets).find(e => e.platform.indexOf(`vlive-${rindex}-`) === 0);
      if (existsConf) {
        confs.push(existsConf);
      } else {
        confs.push({
          platform: `vlive-${rindex}-${rid}`, enabled: false, index: String(rindex), allowCustom: false,
          server: null, secret: null, custom: true, label: `VLive #${rindex}`, files: [],
        });
      }
    }

    setConfigs(confs);
    console.log(`VLive: Init configs ${JSON.stringify(confs)}`);
  }, [env, defaultSecrets, setConfigs, language, t]);

  React.useEffect(() => {
    const refreshStreams = () => {
      axios.post('/terraform/v1/ffmpeg/vlive/streams', {
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        setVLives(res.data.data.map((e, i) => {
          const item = {
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
  }, [t, handleError]);

  // Update config object in array.
  const updateConfigObject = React.useCallback((conf) => {
    const confs = configs.map((e) => {
      if (e.platform === conf.platform) {
        return conf;
      }
      return e;
    })
    setConfigs(confs);
    console.log(`VLive: Update config ${JSON.stringify(conf)} to ${JSON.stringify(confs)}`);
  }, [configs, setConfigs]);

  const updateSecrets = React.useCallback((e, action, platform, server, secret, enabled, custom, label, files, onSuccess) => {
    e.preventDefault();
    if (!files?.length) return alert(t('plat.com.video'));
    if (!server) return alert(t('plat.com.addr'));
    if (custom && !label) return alert(t('plat.com.label'));

    try {
      setSubmiting(true);

      axios.post('/terraform/v1/ffmpeg/vlive/secret', {
        action, platform, server, secret, enabled: !!enabled, custom: !!custom, label, files,
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
    <Accordion defaultActiveKey={defaultActiveKey}>
      <React.Fragment>
        {language === 'zh' ?
        <Accordion.Item eventKey="0">
          <Accordion.Header>场景介绍</Accordion.Header>
          <Accordion.Body>
            <div>
              虚拟直播，是将一个视频文件，用FFmpeg转成直播流，推送到Oryx或其他平台。
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
        </Accordion.Item> :
        <Accordion.Item eventKey="0">
          <Accordion.Header>Introduction</Accordion.Header>
          <Accordion.Body>
            <div>
              Virtual live streaming is the process of converting a video file into a live stream using FFmpeg and pushing it to the Oryx or other platforms.
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
                <SrsErrorBoundary>
                  <ChooseVideoSource platform={conf.platform} vLiveFiles={conf.files} setVLiveFiles={(files) => updateConfigObject({...conf, files: files})} endpoint='vlive' />
                </SrsErrorBoundary>
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
                    updateSecrets(e, 'update', conf.platform, conf.server, conf.secret, !conf.enabled, conf.custom, conf.label, conf.files, () => {
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
            vLives?.length ? (
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
                  vLives?.map(file => {
                    return <tr key={file.platform} style={{verticalAlign: 'middle'}}>
                      <td>{file.i}</td>
                      <td>{file.custom ? (file.label ? '' : t('plat.com.custom')) : file.name} {file.label}</td>
                      <td>
                        <Badge bg={file.enabled ? (file.frame ? 'success' : 'primary') : 'secondary'}>
                          {file.enabled ? (file.frame ? t('plat.com.s0') : t('plat.com.s1')) : t('plat.com.s2')}
                        </Badge>
                      </td>
                      <td>
                        {file.start && file.start?.format('YYYY-MM-DD')}<br/>
                        {file.start && file.start?.format('HH:mm:ss')}
                      </td>
                      <td>
                        {file.ready && file.ready?.format('YYYY-MM-DD')}<br/>
                        {file.ready && file.ready?.format('HH:mm:ss')}
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
          {!vLives?.length ? t('vle.s3') : ''}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

