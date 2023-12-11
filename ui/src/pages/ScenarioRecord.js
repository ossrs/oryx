//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {Accordion, Form, Button, Table} from "react-bootstrap";
import {Token, StreamURL, Clipboard} from "../utils";
import axios from "axios";
import moment from "moment";
import {useRecordStatus} from "../components/DvrStatus";
import {useErrorHandler} from "react-error-boundary";
import {useTranslation} from "react-i18next";
import {useSrsLanguage} from "../components/LanguageSwitch";
import * as Icon from "react-bootstrap-icons";
import PopoverConfirm from "../components/PopoverConfirm";
import TutorialsText from "../components/TutorialsText";
import { minimatch } from "minimatch";

export default function ScenarioRecord() {
  const recordStatus = useRecordStatus();
  const [activeKey, setActiveKey] = React.useState();

  // We must init the activeKey, because the defaultActiveKey only apply when init for Accordion.
  // See https://stackoverflow.com/q/61324259/17679565
  React.useEffect(() => {
    if (!recordStatus) return;

    if (recordStatus.all) {
      setActiveKey('3');
    } else {
      setActiveKey('2');
    }
  }, [recordStatus]);

  if (!activeKey) return <></>;
  return <ScenarioRecordImpl activeKey={activeKey} defaultApplyAll={recordStatus.all} defaultGlobs={recordStatus.globs} recordHome={recordStatus.home} />;
}

function ScenarioRecordImpl({activeKey, defaultApplyAll, defaultGlobs, recordHome}) {
  const language = useSrsLanguage();
  const {t} = useTranslation();
  const handleError = useErrorHandler();

  const [recordAll, setRecordAll] = React.useState(defaultApplyAll);
  const [recordFiles, setRecordFiles] = React.useState();
  const [refreshNow, setRefreshNow] = React.useState();

  const [showGlobTest, setShowGlobTest] = React.useState(false);
  const [globFilters, setGlobFilters] = React.useState(defaultGlobs ? defaultGlobs.join('\n') : '');
  const [targetUrl, setTargetUrl] = React.useState();

  const testGlobFilters = React.useCallback(() => {
    if (!targetUrl) return alert(t('record.urlEmpty'));
    if (!globFilters) return alert(t('record.globEmpty'));

    const a0 = document.createElement("a");
    a0.href = targetUrl.replace('rtmp:', 'http:');

    let matched, matchGlob;
    globFilters.split('\n').forEach(glob => {
      if (minimatch(a0.pathname, glob)) {
        matched = true;
        matchGlob = glob;
      }
    });

    if (matched) {
      alert(`OK! URL ${targetUrl} is matched by glob filter: ${matchGlob}`);
    } else {
      alert(`Failed! URL ${targetUrl} is not matched by any glob filters!`);
    }
  }, [t, targetUrl, globFilters]);

  const updateGlobFilters = React.useCallback(() => {
    const token = Token.load();
    axios.post('/terraform/v1/hooks/record/globs', {
      ...token, globs: globFilters.split('\n'),
    }).then(res => {
      alert(t('record.setupOk'));
      console.log(`Record: Update glob filters ok`);
    }).catch(handleError);
  }, [t, handleError, globFilters]);

  React.useEffect(() => {
    const refreshRecordFiles = () => {
      const token = Token.load();
      axios.post('/terraform/v1/hooks/record/files', {
        ...token,
      }).then(res => {
        console.log(`Record: Files ok, ${JSON.stringify(res.data.data)}`);
        setRecordFiles(res.data.data.map(file => {
          const l = window.location;
          const schema = l.protocol.replace(':', '');
          const httpPort = l.port || (l.protocol === 'http:' ? 80 : 443);
          if (file.progress) {
            file.location = `${l.protocol}//${l.host}/terraform/v1/hooks/record/hls/${file.uuid}.m3u8`;
            file.preview = `/players/srs_player.html?schema=${schema}&port=${httpPort}&autostart=true&app=terraform/v1/hooks/record/hls&stream=${file.uuid}.m3u8`;
          } else {
            file.location = `${l.protocol}//${l.host}/terraform/v1/hooks/record/hls/${file.uuid}/index.mp4`;
            file.preview = `/terraform/v1/hooks/record/hls/${file.uuid}/index.mp4`;
          }

          return {
            ...file,
            url: StreamURL.build(file.vhost, file.app, file.stream),
            update: moment(file.update),
            duration: Number(file.duration),
            size: Number(file.size / 1024.0 / 1024),
          };
        }).sort((a, b) => {
          return b.update - a.update;
        }).map((file, i) => {
          return {...file, i: i + 1};
        }));
      }).catch(handleError);
    };

    refreshRecordFiles();
    const timer = setInterval(() => refreshRecordFiles(), 10 * 1000);
    return () => clearInterval(timer);
  }, [handleError, refreshNow]);

  const setupRecordPattern = React.useCallback((e, recordAll) => {
    e.preventDefault();

    const token = Token.load();
    axios.post('/terraform/v1/hooks/record/apply', {
      ...token, all: !!recordAll,
    }).then(res => {
      alert(t('record.setupOk'));
      console.log(`Record: Apply patterns ok, all=${recordAll}`);
    }).catch(handleError);
  }, [handleError, t]);

  const removeRecord = React.useCallback((file) => {
    const token = Token.load();
    axios.post('/terraform/v1/hooks/record/remove', {
      ...token, uuid: file.uuid,
    }).then(res => {
      setRefreshNow(!refreshNow);
      console.log(`Record: Remove file ok, file=${JSON.stringify(file)}`);
    }).catch(handleError);
  }, [refreshNow, handleError, setRefreshNow]);

  const endRecord = React.useCallback((file) => {
    const token = Token.load();
    axios.post('/terraform/v1/hooks/record/end', {
      ...token, uuid: file.uuid,
    }).then(res => {
      setTimeout(() => {
        setRefreshNow(!refreshNow);
      }, 1000);
      console.log(`Record: End file ok, file=${JSON.stringify(file)}`);
    }).catch(handleError);
  }, [refreshNow, handleError, setRefreshNow]);

  const copyToClipboard = React.useCallback((e, text) => {
    e.preventDefault();

    Clipboard.copy(text).then(() => {
      alert(t('helper.copyOk'));
    }).catch((err) => {
      alert(`${t('helper.copyFail')} ${err}`);
    });
  }, [t]);

  return (
    <Accordion defaultActiveKey={[activeKey]} alwaysOpen>
      <React.Fragment>
        {language === 'zh' ?
          <Accordion.Item eventKey="0">
            <Accordion.Header>场景介绍</Accordion.Header>
            <Accordion.Body>
              <div>
                本地录制，指录制视频流到SRS Stack的本地磁盘，只要推送到服务器的流都可以录制。
                <p></p>
              </div>
              <p>可应用的具体场景包括：</p>
              <ul>
                <li>直播转点播，录制直播流成为一个HLS文件，存储在SRS Stack本地磁盘，可以下载</li>
              </ul>
              <p>特别注意：</p>
              <ul>
                <li>如果流的路数特别多，磁盘会很忙，特别是挂共享存储，需要监控磁盘IO和负载。</li>
                <li>虽然本地磁盘足够大，但云存储是真的无限大，而本地磁盘其实还是有限的，需要监控磁盘空间。</li>
                <li>暂时不支持本地文件的管理，比如删除和清理等。</li>
              </ul>
              <p>使用说明：</p>
              <ul>
                <li>具体使用步骤，请根据下面引导操作</li>
              </ul>
            </Accordion.Body>
          </Accordion.Item> :
          <Accordion.Item eventKey="0">
            <Accordion.Header>Introduction</Accordion.Header>
            <Accordion.Body>
              <div>
                Local recording refers to recording video streams to the local disk of the SRS Stack, and any stream pushed to the server can be recorded.
                <p></p>
              </div>
              <p>Specific application scenarios include:</p>
              <ul>
                <li>Live to VOD, recording live streams into an HLS file, stored on the SRS Stack local disk, and can be downloaded</li>
              </ul>
              <p>Special attention:</p>
              <ul>
                <li>If there are many streams, the disk will be very busy, especially when using shared storage, so it is necessary to monitor disk IO and load.</li>
                <li>Although the local disk is large enough, cloud storage is truly unlimited, while local disk space is actually limited, so it is necessary to monitor disk space.</li>
                <li>Temporary management of local files, such as deletion and cleanup, is not supported.</li>
              </ul>
              <p>Instructions for use:</p>
              <ul>
                <li>For specific usage steps, please follow the guide below</li>
              </ul>
            </Accordion.Body>
          </Accordion.Item>
        }
      </React.Fragment>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{t('record.dir')}</Accordion.Header>
        <Accordion.Body>
          {t('record.dir2')} <code>{recordHome}</code> &nbsp;
          <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
            <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, recordHome)} />
          </div> &nbsp;
          <TutorialsText prefixLine={true} title={t('record.dir3')}>
            {t('record.dir4')} <font color='red'>/data</font> {t('record.dir5')} &nbsp;
            <font color='red'>{t('record.dir6')}</font> {t('record.dir7')}
          </TutorialsText>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>{t('record.rule')}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>{t('record.glob')}</Form.Label>
              <Form.Text> * {t('record.glob2')} &nbsp;
                <a href='#noid' onClick={(e)=>{e.preventDefault(); setShowGlobTest(!showGlobTest);}}>{t('record.glob3')}</a> &nbsp;
                {t('record.glob4')} </Form.Text>
              <Form.Control as="textarea" type='text' rows={5} defaultValue={globFilters}
                            placeholder="For example: /live/livestream or /live/* or /*/*"
                            onChange={(e) => setGlobFilters(e.target.value)} />
            </Form.Group>
            <React.Fragment>
              <Button variant="primary" type="submit" onClick={(e) => {
                e.preventDefault();
                updateGlobFilters();
              }}>{t('record.update2')}</Button> &nbsp;
              <Button variant="primary" type="submit" onClick={(e) => {
                e.preventDefault();
                setRecordAll(!recordAll);
                setupRecordPattern(e, !recordAll);
              }}>
                {recordAll ? t('record.stop') : t('record.start')}
              </Button> &nbsp;
              <Form.Text> * {t('record.rule2')}</Form.Text>
            </React.Fragment>
            {showGlobTest && <React.Fragment>
              <p/>
              <Form.Group className="mb-3">
                <Form.Label>{t('record.test')}</Form.Label>
                <Form.Text> * {t('record.test2')} </Form.Text>
                <Form.Control as="input" type='text' placeholder="For example: rtmp://localhost/live/livestream"
                              onChange={(e) => setTargetUrl(e.target.value)} />
              </Form.Group>
              <Button variant="primary" type="submit" onClick={(e) => {
                e.preventDefault();
                testGlobFilters();
              }}>{t('record.test3')}</Button>
            </React.Fragment>}
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>{t('record.task')}</Accordion.Header>
        <Accordion.Body>
          {
            recordFiles?.length ? (
              <Table striped bordered hover>
                <thead>
                <tr>
                  <th>#</th>
                  <th>{t('record.status')}</th>
                  <th>{t('record.update')}</th>
                  <th>{t('record.source')}</th>
                  <th>{t('record.duration')}</th>
                  <th>{t('record.size')}</th>
                  <th>{t('record.slice')}</th>
                  <th>{t('record.url')}</th>
                  <th>{t('record.action')}</th>
                </tr>
                </thead>
                <tbody>
                {
                  recordFiles?.map(file => {
                    return <tr key={file.uuid} name={file.uuid}>
                      <td>{file.i}</td>
                      <td title={t('record.tip')}>{file.progress ? t('record.ing') : t('record.done')}</td>
                      <td>{`${file.update.format('YYYY-MM-DD HH:mm:ss')}`}</td>
                      <td>{file.url}</td>
                      <td>{`${file.duration.toFixed(1)}`} {t('helper.seconds')}</td>
                      <td>{`${file.size.toFixed(1)}`}MB</td>
                      <td>{file.nn}</td>
                      <td><a href={file.location} onClick={(e) => copyToClipboard(e, file.location)} target='_blank' rel='noreferrer'>{t('helper.copy2')}</a></td>
                      <td>
                        <a href={file.preview} target='_blank' rel='noreferrer'>
                          {t('helper.preview')}
                        </a> &nbsp;
                        <PopoverConfirm placement='top' trigger={ <a href={`#${file.uuid}`} hidden={!file.progress}>{t('helper.end')}</a> } onClick={() => endRecord(file)}>
                          <p>
                            {t('scenario.endTip')}
                          </p>
                        </PopoverConfirm> <PopoverConfirm placement='top' trigger={ <a href={`#${file.uuid}`} hidden={file.progress}>{t('helper.delete')}</a> } onClick={() => removeRecord(file)}>
                          <p>
                            {t('scenario.rmFileTip1')} &nbsp;
                            <span className='text-danger'><strong>
                              {t('scenario.rmFileTip2')}
                            </strong></span>
                            {t('scenario.rmFileTip3')}
                          </p>
                        </PopoverConfirm> &nbsp;
                      </td>
                    </tr>;
                  })
                }
                </tbody>
              </Table>
            ) : ''
          }
          {!recordFiles?.length ? t('record.none') : ''}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}
