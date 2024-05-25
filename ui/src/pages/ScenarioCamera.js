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
import {SrsEnvContext} from "../components/SrsEnvContext";
import ChooseVideoSource from "../components/VideoSourceSelector";

export default function ScenarioCamera() {
  const [init, setInit] = React.useState();
  const [activeKey, setActiveKey] = React.useState();
  const [secrets, setSecrets] = React.useState();
  const handleError = useErrorHandler();

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

    setActiveKey('1');
    for (const key in secrets) {
      const e = secrets[key];
      if (e.enabled && e.server && e.secret) {
        setActiveKey('99');
      }
    }
  }, [init, secrets]);

  if (!activeKey) return <></>;
  return <ScenarioCameraImpl defaultActiveKey={activeKey} defaultSecrets={secrets}/>;
}

function ScenarioCameraImpl({defaultActiveKey, defaultSecrets}) {
  const language = useSrsLanguage();
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const env = React.useContext(SrsEnvContext)[0];

  const [cameras, setCameras] = React.useState();
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

    // Generate more camera streaming configures.
    while (confs.length < env.cameraLimit) {
      const rindex = index++;
      const rid = Math.random().toString(16).slice(-6);

      // Load the configured virtual live from defaults.
      const existsConf = Object.values(defaultSecrets).find(e => e.platform.indexOf(`camera-${rindex}-`) === 0);
      if (existsConf) {
        confs.push(existsConf);
      } else {
        confs.push({
          platform: `camera-${rindex}-${rid}`, enabled: false, index: String(rindex), allowCustom: false,
          server: null, secret: null, custom: true, label: `Camera #${rindex}`, files: [], extraAudio: '',
        });
      }
    }

    setConfigs(confs);
    console.log(`Camera: Init configs ${JSON.stringify(confs)}`);
  }, [defaultSecrets, setConfigs, language, t]);

  React.useEffect(() => {
    const refreshStreams = () => {
      axios.post('/terraform/v1/ffmpeg/camera/streams', {
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        setCameras(res.data.data.map((e, i) => {
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
        console.log(`Camera: Query streams ${JSON.stringify(res.data.data)}`);
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
    console.log(`Camera: Update config ${JSON.stringify(conf)} to ${JSON.stringify(confs)}`);
  }, [configs, setConfigs]);

  const updateSecrets = React.useCallback((e, action, platform, server, secret, enabled, custom, label, files, extraAudio, onSuccess) => {
    e.preventDefault();
    if (!files?.length) return alert(t('camera.source'));
    if (!server) return alert(t('plat.com.addr'));
    if (custom && !label) return alert(t('plat.com.label'));

    try {
      setSubmiting(true);

      axios.post('/terraform/v1/ffmpeg/camera/secret', {
        action, platform, server, secret, enabled: !!enabled, custom: !!custom, label, files,
        extraAudio,
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
                摄像头直播，是将一个摄像头的流，用FFmpeg转成直播流，推送到Oryx或其他平台。
                <p></p>
              </div>
              <p>可应用的具体场景包括：</p>
              <ul>
                <li>无人直播间，7x24小时获得直播收益</li>
              </ul>
              <p>使用说明：</p>
              <ul>
                <li>首先设置摄像头拉流信息</li>
                <li>然后设置直播流信息</li>
              </ul>
            </Accordion.Body>
          </Accordion.Item> :
          <Accordion.Item eventKey="0">
            <Accordion.Header>Introduction</Accordion.Header>
            <Accordion.Body>
              <div>
                Camera streaming is the process of converting the stream from RTSP/IP Camera into a live stream using FFmpeg and pushing it to the Oryx or other platforms.
                <p></p>
              </div>
              <p>Specific application scenarios include:</p>
              <ul>
                <li>Unmanned live streaming rooms, 7x24 hours of live streaming revenue</li>
              </ul>
              <p>Instructions for use:</p>
              <ul>
                <li>First, setup the Camera URL to pull stream from</li>
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
                  <ChooseCameraSource platform='kuaishou' cameraFiles={conf.files} setCameraFiles={(files) => updateConfigObject({...conf, files: files})} />
                  <CameraExtraAudioTrack extraAudio={conf.extraAudio} setExtraAudio={(extraAudio) => updateConfigObject({...conf, extraAudio: extraAudio})} />
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
                    updateSecrets(e, 'update', conf.platform, conf.server, conf.secret, !conf.enabled, conf.custom, conf.label, conf.files, conf.extraAudio, () => {
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
            cameras?.length ? (
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
                  <th>{t('camera.extra')}</th>
                  <th>{t('plat.com.log')}</th>
                </tr>
                </thead>
                <tbody>
                {
                  cameras?.map(file => {
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
                        <CameraFileFormatInfo file={file.sourceObj}/>
                      </td>
                      <td>{file?.extraAudio || 'OFF'}</td>
                      <td>{file.frame?.log}</td>
                    </tr>;
                  })
                }
                </tbody>
              </Table>
            ) : ''
          }
          {!cameras?.length ? t('camera.s3') : ''}
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

function ChooseCameraSource({platform, cameraFiles, setCameraFiles}) {
  const {t} = useTranslation();

  return (<>
    <Form.Group className="mb-3">
      <Form.Label>{t('plat.tool.source')}</Form.Label>
      <Form.Text> * {t('plat.tool.stream2')}</Form.Text>
      <SrsErrorBoundary>
        <CameraStreamSelector platform={platform} cameraFiles={cameraFiles} setCameraFiles={setCameraFiles}/>
      </SrsErrorBoundary>
    </Form.Group>
  </>);
}

function CameraStreamSelector({platform, cameraFiles, setCameraFiles}) {
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const [inputStream, setInputStream] = React.useState(cameraFiles?.length ? cameraFiles[0].target : '');
  const [submiting, setSubmiting] = React.useState();

  const checkStreamUrl = React.useCallback(async () => {
    if (!inputStream) return alert(t('plat.tool.stream3'));
    const isHTTP = inputStream.startsWith('http://') || inputStream.startsWith('https://');
    if (!inputStream.startsWith('rtmp://') && !inputStream.startsWith('srt://') && !inputStream.startsWith('rtsp://') && !isHTTP) return alert(t('plat.tool.stream2'));
    if (isHTTP && inputStream.indexOf('.flv') < 0 && inputStream.indexOf('.m3u8') < 0) return alert(t('plat.tool.stream4'));

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
        console.log(`${t('plat.tool.stream5')}，${JSON.stringify(res.data.data)}`);
        const streamObj = res.data.data;
        const files = [{name: streamObj.name, size: 0, uuid: streamObj.uuid, target: streamObj.target, type: "stream"}];
        axios.post('/terraform/v1/ffmpeg/camera/source', {
          platform, files,
        }, {
          headers: Token.loadBearerHeader(),
        }).then(res => {
          console.log(`${t('plat.tool.stream6')}，${JSON.stringify(res.data.data)}`);
          setCameraFiles(res.data.data.files);
          resolve();
        }).catch(reject);
      });
    } catch (e) {
      handleError(e);
    } finally {
      setSubmiting(false);
    }
  }, [t, inputStream, handleError, platform, setCameraFiles, setSubmiting]);

  return (<>
    <Form.Control as="div">
      {!cameraFiles?.length ? <>
        <Row>
          <Col>
            <Form.Control type="text" defaultValue={inputStream} placeholder={t('plat.tool.stream3')} onChange={e => setInputStream(e.target.value)} />
          </Col>
          <Col xs="auto">
            <Button variant="primary" disabled={submiting} onClick={checkStreamUrl}>{t('helper.submit')}</Button>
          </Col>
        </Row></> : <></>
      }
      {cameraFiles?.length ? <CameraFileList files={cameraFiles} onChangeFiles={(e) => setCameraFiles(null)}/> : <></>}
    </Form.Control>
  </>);
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

