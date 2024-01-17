//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {Accordion, Badge, Button, Col, Form, InputGroup, ListGroup, Row, Table} from "react-bootstrap";
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
  return <ScenarioVLiveImpl defaultActiveKey={activeKey} defaultSecrets={secrets}/>;
}

function ScenarioVLiveImpl({defaultActiveKey, defaultSecrets}) {
  const language = useSrsLanguage();
  const {t} = useTranslation();
  const handleError = useErrorHandler();

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
  }, [handleError, setSubmiting]);

  return (
    <Accordion defaultActiveKey={defaultActiveKey}>
      <React.Fragment>
        {language === 'zh' ?
        <Accordion.Item eventKey="0">
          <Accordion.Header>场景介绍</Accordion.Header>
          <Accordion.Body>
            <div>
              虚拟直播，是将一个视频文件，用FFmpeg转成直播流，推送到SRS Stack或其他平台。
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
            <SrsErrorBoundary>
              <ChooseVideoSource platform='wx' vLiveFiles={wxFiles} setVLiveFiles={setWxFiles} />
            </SrsErrorBoundary>
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
                updateSecrets(e, 'update', 'wx', wxServer, wxSecret, !wxEnabled, wxCustom, wxLabel, wxFiles, () => {
                  setWxEnabled(!wxEnabled);
                });
              }}
            >
              {wxEnabled ? t('plat.com.stop') : t('plat.com.start')}
            </Button> &nbsp;
            <Form.Text> * {t('vle.tip')}</Form.Text>
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
            <SrsErrorBoundary>
              <ChooseVideoSource platform='bilibili' vLiveFiles={bilibiliFiles} setVLiveFiles={setBilibiliFiles} />
            </SrsErrorBoundary>
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
                updateSecrets(e, 'update', 'bilibili', bilibiliServer, bilibiliSecret, !bilibiliEnabled, bilibiliCustom, bilibiliLabel, bilibiliFiles, () => {
                  setBilibiliEnabled(!bilibiliEnabled);
                });
              }}
            >
              {bilibiliEnabled ? t('plat.com.stop') : t('plat.com.start')}
            </Button> &nbsp;
            <Form.Text> * {t('vle.tip')}</Form.Text>
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
            <SrsErrorBoundary>
              <ChooseVideoSource platform='kuaishou' vLiveFiles={kuaishouFiles} setVLiveFiles={setKuaishouFiles} />
            </SrsErrorBoundary>
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
                updateSecrets(e, 'update', 'kuaishou', kuaishouServer, kuaishouSecret, !kuaishouEnabled, kuaishouCustom, kuaishouLabel, kuaishouFiles, () => {
                  setKuaishouEnabled(!kuaishouEnabled);
                });
              }}
            >
              {kuaishouEnabled ? t('plat.com.stop') : t('plat.com.start')}
            </Button> &nbsp;
            <Form.Text> * {t('vle.tip')}</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
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

function ChooseVideoSource({platform, vLiveFiles, setVLiveFiles}) {
  const {t} = useTranslation();

  const [checkType, setCheckType] = React.useState('upload');
  React.useEffect(() => {
    if (vLiveFiles?.length) {
      const type = vLiveFiles[0].type;
      if (type === 'upload' || type === 'file' || type === 'stream') {
        setCheckType(type);
      }
    }
  }, [vLiveFiles]);
  return (<>
    <Form.Group className="mb-3">
      <Form.Label>{t('plat.tool.source')}</Form.Label>
      <Form.Text> * {t('vle.source2')}</Form.Text>
      <Form.Check type="radio" label={t('plat.tool.upload')} id={'upload-' + platform} checked={checkType === 'upload'}
        name={'chooseSource-' + platform} onChange={e => setCheckType('upload')}
      />
      {checkType === 'upload' && 
      <SrsErrorBoundary>
        <VLiveFileUploader platform={platform} vLiveFiles={vLiveFiles} setVLiveFiles={setVLiveFiles} />
      </SrsErrorBoundary>
      }
    </Form.Group>
    <Form.Group className="mb-3">
      <InputGroup>
        <Form.Check type="radio" label={t('plat.tool.file')} id={'server-' + platform} checked={checkType === 'file'}
                    name={'chooseSource' + platform} onChange={e => setCheckType('file')}
        /> &nbsp;
        <Form.Text> * {t('plat.tool.file2')}</Form.Text>
      </InputGroup>
      {checkType === 'file' &&
      <SrsErrorBoundary>
        <VLiveFileSelector platform={platform} vLiveFiles={vLiveFiles} setVLiveFiles={setVLiveFiles}/>
      </SrsErrorBoundary>
      }
    </Form.Group>
    <Form.Group className="mb-3">
      <InputGroup>
        <Form.Check type="radio" label={t('plat.tool.stream')} id={'stream-' + platform} checked={checkType === 'stream'}
                    name={'chooseSource' + platform} onChange={e => setCheckType('stream')}
        /> &nbsp;
        <Form.Text> * {t('plat.tool.stream2')}</Form.Text>
      </InputGroup>
      {checkType === 'stream' &&
      <SrsErrorBoundary>
        <VLiveStreamSelector platform={platform} vLiveFiles={vLiveFiles} setVLiveFiles={setVLiveFiles}/>
      </SrsErrorBoundary>
      }
    </Form.Group>
  </>);
}

function VLiveStreamSelector({platform, vLiveFiles, setVLiveFiles}) {
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const [inputStream, setInputStream] = React.useState(vLiveFiles?.length ? vLiveFiles[0].target :'');
  const [submiting, setSubmiting] = React.useState();

  const checkStreamUrl = React.useCallback(async () => {
    if (!inputStream) return alert(t('plat.tool.stream3'));
    const isHTTP = inputStream.startsWith('http://') || inputStream.startsWith('https://');
    if (!inputStream.startsWith('rtmp://') && !inputStream.startsWith('rtsp://') && !isHTTP) return alert(t('plat.tool.stream2'));
    if (isHTTP && inputStream.indexOf('.flv') < 0 && inputStream.indexOf('.m3u8') < 0) return alert(t('plat.tool.stream4'));

    setSubmiting(true);
    try {
      const res = await new Promise((resolve, reject) => {
        axios.post(`/terraform/v1/ffmpeg/vlive/stream-url`, {
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
        axios.post('/terraform/v1/ffmpeg/vlive/source', {
          platform, files,
        }, {
          headers: Token.loadBearerHeader(),
        }).then(res => {
          console.log(`${t('plat.tool.stream6')}，${JSON.stringify(res.data.data)}`);
          setVLiveFiles(res.data.data.files);
          resolve();
        }).catch(reject);
      });
    } catch (e) {
      handleError(e);
    } finally {
      setSubmiting(false);
    }
  }, [inputStream, handleError, platform, setVLiveFiles, setSubmiting]);

  return (<>
    <Form.Control as="div">
      {!vLiveFiles?.length && <>
        <Row>
          <Col>
            <Form.Control type="text" defaultValue={inputStream} placeholder={t('plat.tool.stream3')} onChange={e => setInputStream(e.target.value)} />
          </Col>
          <Col xs="auto">
            <Button variant="primary" disabled={submiting} onClick={checkStreamUrl}>{t('helper.submit')}</Button>
          </Col>
        </Row></>
      }
      {vLiveFiles?.length && <VLiveFileList files={vLiveFiles} onChangeFiles={(e) => setVLiveFiles(null)}/>}
    </Form.Control>
  </>);
}

function VLiveFileSelector({platform, vLiveFiles, setVLiveFiles}) {
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  // TODO: FIXME: As the file path is changed after used, so we can not use te target.
  const [inputFile, setInputFile] = React.useState('');

  const CheckLocalFile = React.useCallback(() => {
    if (!inputFile) return alert(t('plat.tool.file3'));
    if (!inputFile.startsWith('/data') && !inputFile.startsWith('upload/') && !inputFile.startsWith('./upload/')) return alert(t('plat.tool.file2'));

    const fileExtension = inputFile.slice(inputFile.lastIndexOf('.'));
    if (!['.mp4', '.flv', '.ts'].includes(fileExtension)) return alert(t('plat.tool.file4'));

    axios.post(`/terraform/v1/ffmpeg/vlive/server`, {
      file: inputFile,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      console.log(`${t('plat.tool.file5')}，${JSON.stringify(res.data.data)}`);
      const localFileObj = res.data.data;
      const files = [{name: localFileObj.name, size: localFileObj.size, uuid: localFileObj.uuid, target: localFileObj.target, type: "file"}];
      axios.post('/terraform/v1/ffmpeg/vlive/source', {
        platform, files,
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        console.log(`${t('plat.tool.file6')}，${JSON.stringify(res.data.data)}`);
        setVLiveFiles(res.data.data.files);
      }).catch(handleError);
    }).catch(handleError);
  }, [inputFile, handleError, platform, setVLiveFiles]);

  return (<>
    <Form.Control as="div">
      {!vLiveFiles?.length && <>
        <Row>
          <Col>
            <Form.Control type="text" defaultValue={inputFile} placeholder={t('plat.tool.file7')} onChange={e => setInputFile(e.target.value)} />
          </Col>
          <Col xs="auto">
            <Button variant="primary" onClick={CheckLocalFile}>{t('helper.submit')}</Button>
          </Col>
        </Row></>
      }
      {vLiveFiles?.length && <VLiveFileList files={vLiveFiles} onChangeFiles={(e) => setVLiveFiles(null)}/>}
    </Form.Control>
  </>);
}

function VLiveFileUploader({platform, vLiveFiles, setVLiveFiles}) {
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const updateSources = React.useCallback((platform, files, setFiles) => {
    if (!files?.length) return alert(t('plat.tool.upload2'));

    axios.post('/terraform/v1/ffmpeg/vlive/source', {
      platform, files: files.map(f => {
        return {name: f.name, size: f.size, uuid: f.uuid, target: f.target, type: "upload"};
      }),
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      console.log(`${t('plat.tool.upload3')}, ${JSON.stringify(res.data.data)}`);
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

function VLiveFileVideoInfo({file}) {
  const f = file;
  if (!f?.video) return <>NoVideo</>;
  return <>Video({f?.video?.codec_name} {f?.video?.profile} {f?.video?.width}x{f?.video?.height})</>;
}

function VLiveFileAudioInfo({file}) {
  const f = file;
  if (!f?.audio) return <>NoAudio</>;
  return <>Audio({f?.audio?.codec_name} {f?.audio?.sample_rate}HZ {f?.audio?.channels}CH)</>;
}

