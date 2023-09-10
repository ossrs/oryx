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

export default function ScenarioRecord() {
  const language = useSrsLanguage();
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
  if (language === 'zh') {
    return <ScenarioRecordImplCn activeKey={activeKey} defaultApplyAll={recordStatus.all} recordHome={recordStatus.home} />;
  }
  return <ScenarioRecordImplEn activeKey={activeKey} defaultApplyAll={recordStatus.all} recordHome={recordStatus.home} />;
}

function ScenarioRecordImplCn({activeKey, defaultApplyAll, recordHome}) {
  const [recordAll, setRecordAll] = React.useState(defaultApplyAll);
  const [recordFiles, setRecordFiles] = React.useState();
  const [refreshNow, setRefreshNow] = React.useState();
  const handleError = useErrorHandler();
  const {t} = useTranslation();

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
      alert('设置录制规则成功');
      console.log(`Record: Apply patterns ok, all=${recordAll}`);
    }).catch(handleError);
  }, [handleError]);

  const removeRecord = React.useCallback((file) => {
    const token = Token.load();
    axios.post('/terraform/v1/hooks/record/remove', {
      ...token, uuid: file.uuid,
    }).then(res => {
      setRefreshNow(!refreshNow);
      console.log(`Record: Remove file ok, file=${JSON.stringify(file)}`);
    }).catch(handleError);
  }, [refreshNow, handleError]);

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
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>录制文件夹</Accordion.Header>
        <Accordion.Body>
          保存路径： <code>{recordHome}</code> &nbsp;
          <div role='button' style={{display: 'inline-block'}} title='拷贝'>
            <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, recordHome)} />
          </div> &nbsp;
          <TutorialsText prefixLine={true} title='如何修改目录?'>
            你可以使用软链接，将本目录软链到其他磁盘的目录，然后<font color='red'>重启服务</font>(必须)
          </TutorialsText>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>设置录制规则</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Button variant="primary" type="submit" onClick={(e) => {
              setRecordAll(!recordAll);
              setupRecordPattern(e, !recordAll);
            }}>
              {recordAll ? '停止录制' : '开始录制'}
            </Button> &nbsp;
            <Form.Text> * 录制所有流</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>录制任务列表</Accordion.Header>
        <Accordion.Body>
          {
            recordFiles?.length ? (
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
                  <th>地址</th>
                  <th>操作</th>
                </tr>
                </thead>
                <tbody>
                {
                  recordFiles?.map(file => {
                    return <tr key={file.uuid} name={file.uuid}>
                      <td>{file.i}</td>
                      <td title='若300秒没有流，会自动完成录制'>{file.progress ? '录制中' : '已完成'}</td>
                      <td>{`${file.update.format('YYYY-MM-DD HH:mm:ss')}`}</td>
                      <td>{file.url}</td>
                      <td>{`${file.duration.toFixed(1)}`}秒</td>
                      <td>{`${file.size.toFixed(1)}`}MB</td>
                      <td>{file.nn}</td>
                      <td><a href={file.location} onClick={(e) => copyToClipboard(e, file.location)} target='_blank' rel='noreferrer'>复制</a></td>
                      <td>
                        <a href={file.preview} target='_blank' rel='noreferrer'>预览</a> &nbsp;
                        <PopoverConfirm placement='top' trigger={ <a href={`#${file.uuid}`} hidden={file.progress}>删除</a> } onClick={() => removeRecord(file)}>
                          <p>
                            {t('scenario.rmFileTip1')}
                            <span className='text-danger'><strong>
                              {t('scenario.rmFileTip2')}
                            </strong></span>
                            {t('scenario.rmFileTip3')}
                          </p>
                        </PopoverConfirm>
                      </td>
                    </tr>;
                  })
                }
                </tbody>
              </Table>
            ) : ''
          }
          {!recordFiles?.length ? '没有流。请开启录制并推流后，等待大约60秒左右，录制列表会自动更新' : ''}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

function ScenarioRecordImplEn({activeKey, defaultApplyAll, recordHome}) {
  const [recordAll, setRecordAll] = React.useState(defaultApplyAll);
  const [recordFiles, setRecordFiles] = React.useState();
  const [refreshNow, setRefreshNow] = React.useState();
  const handleError = useErrorHandler();
  const {t} = useTranslation();

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
      alert('Setup OK');
      console.log(`Record: Apply patterns ok, all=${recordAll}`);
    }).catch(handleError);
  }, [handleError]);

  const removeRecord = React.useCallback((file) => {
    const token = Token.load();
    axios.post('/terraform/v1/hooks/record/remove', {
      ...token, uuid: file.uuid,
    }).then(res => {
      setRefreshNow(!refreshNow);
      console.log(`Record: Remove file ok, file=${JSON.stringify(file)}`);
    }).catch(handleError);
  }, [refreshNow, handleError]);

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
      <Accordion.Item eventKey="1">
        <Accordion.Header>Record Directory</Accordion.Header>
        <Accordion.Body>
          Work Directory： <code>{recordHome}</code> &nbsp;
          <div role='button' style={{display: 'inline-block'}} title='Copy'>
            <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, recordHome)} />
          </div> &nbsp;
          <TutorialsText prefixLine={true} title='How to modify the directory?'>
            You can use a symbolic link to link this directory to a directory on another disk, and then <font color='red'>restart the service</font> (required).
          </TutorialsText>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>Setup Record Rules</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Button variant="primary" type="submit" onClick={(e) => {
              setRecordAll(!recordAll);
              setupRecordPattern(e, !recordAll);
            }}>
              {recordAll ? 'Stop Record' : 'Start Record'}
            </Button> &nbsp;
            <Form.Text> * Record all streams</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>Record Tasks</Accordion.Header>
        <Accordion.Body>
          {
            recordFiles?.length ? (
              <Table striped bordered hover>
                <thead>
                <tr>
                  <th>#</th>
                  <th>Status</th>
                  <th>Update</th>
                  <th>Source Stream</th>
                  <th>Duration</th>
                  <th>Size</th>
                  <th>Slices</th>
                  <th>URL</th>
                  <th>Operation</th>
                </tr>
                </thead>
                <tbody>
                {
                  recordFiles?.map(file => {
                    return <tr key={file.uuid} name={file.uuid}>
                      <td>{file.i}</td>
                      <td title='If there is no stream for 300 seconds, the recording will be automatically completed.'>{file.progress ? 'Recording' : 'Done'}</td>
                      <td>{`${file.update.format('YYYY-MM-DD HH:mm:ss')}`}</td>
                      <td>{file.url}</td>
                      <td>{`${file.duration.toFixed(1)}`}s</td>
                      <td>{`${file.size.toFixed(1)}`}MB</td>
                      <td>{file.nn}</td>
                      <td><a href={file.location} onClick={(e) => copyToClipboard(e, file.location)} target='_blank' rel='noreferrer'>Copy</a></td>
                      <td>
                        <a href={file.preview} target='_blank' rel='noreferrer'>Preview</a> &nbsp;
                        <PopoverConfirm placement='top' trigger={ <a href={`#${file.uuid}`} hidden={file.progress}>Delete</a> } onClick={() => removeRecord(file)}>
                          <p>
                            {t('scenario.rmFileTip1')}
                            <span className='text-danger'><strong>
                              {t('scenario.rmFileTip2')}
                            </strong></span>
                            {t('scenario.rmFileTip3')}
                          </p>
                        </PopoverConfirm>
                      </td>
                    </tr>;
                  })
                }
                </tbody>
              </Table>
            ) : ''
          }
          {!recordFiles?.length ? 'There is no stream. Please start the recording and push the stream, then wait for about 60 seconds, and the recording list will be automatically updated.' : ''}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

