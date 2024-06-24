//
// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
//
import React from "react";
import {useTranslation} from "react-i18next";
import {Button, Col, Form, InputGroup, ListGroup, Row} from "react-bootstrap";
import {SrsErrorBoundary} from "./SrsErrorBoundary";
import {useErrorHandler} from "react-error-boundary";
import axios from "axios";
import {MediaSource, Token} from "../utils";
import FileUploader from "./FileUploader";

export default function ChooseVideoSource({platform, endpoint, vLiveFiles, setVLiveFiles, hideStreamSource}) {
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
      <Form.Label>{endpoint === 'dubbing' ? t('plat.tool.source2') : t('plat.tool.source')}</Form.Label>
      <Form.Text> * {endpoint === 'dubbing' ? t('dubb.create.source') : t('vle.source2')}</Form.Text>
      <Form.Check type="radio" label={t('plat.tool.upload')} id={'upload-' + platform} checked={checkType === 'upload'}
                  name={'chooseSource-' + platform} onChange={e => setCheckType('upload')}
      />
      {checkType === 'upload' &&
        <SrsErrorBoundary>
          <VLiveFileUploader {...{platform, endpoint, vLiveFiles, setVLiveFiles}} />
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
          <VLiveFileSelector {...{platform, endpoint, vLiveFiles, setVLiveFiles}} />
        </SrsErrorBoundary>
      }
    </Form.Group>
    {!hideStreamSource && <Form.Group className="mb-3">
      <InputGroup>
        <Form.Check type="radio" label={t('plat.tool.stream')} id={'stream-' + platform}
                    checked={checkType === 'stream'}
                    name={'chooseSource' + platform} onChange={e => setCheckType('stream')}
        /> &nbsp;
        <Form.Text> * {t('plat.tool.stream2')}</Form.Text>
      </InputGroup>
      {checkType === 'stream' &&
        <SrsErrorBoundary>
          <VLiveStreamSelector {...{platform, endpoint, vLiveFiles, setVLiveFiles}} />
        </SrsErrorBoundary>
      }
    </Form.Group>}
  </>);
}

function VLiveStreamSelector({platform, endpoint, vLiveFiles, setVLiveFiles}) {
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const [inputStream, setInputStream] = React.useState(vLiveFiles?.length ? vLiveFiles[0].target :'');
  const [submiting, setSubmiting] = React.useState();

  const checkStreamUrl = React.useCallback(async () => {
    if (!inputStream) return alert(t('plat.tool.stream3'));
    const isHTTP = inputStream.startsWith('http://') || inputStream.startsWith('https://');
    if (!inputStream.startsWith('rtmp://') && !inputStream.startsWith('srt://') && !inputStream.startsWith('rtsp://') && !isHTTP) return alert(t('plat.tool.stream2'));
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
  }, [t, inputStream, handleError, platform, setVLiveFiles, setSubmiting]);

  return (<>
    <Form.Control as="div">
      {!vLiveFiles?.length ? <>
        <Row>
          <Col>
            <Form.Control type="text" defaultValue={inputStream} placeholder={t('plat.tool.stream3')} onChange={e => setInputStream(e.target.value)} />
          </Col>
          <Col xs="auto">
            <Button variant="primary" disabled={submiting} onClick={checkStreamUrl}>{t('helper.submit')}</Button>
          </Col>
        </Row></> : <></>
      }
      {vLiveFiles?.length ? <VLiveFileList files={vLiveFiles} onChangeFiles={(e) => setVLiveFiles(null)}/> : <></>}
    </Form.Control>
  </>);
}

function VLiveFileSelector({platform, endpoint, vLiveFiles, setVLiveFiles}) {
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  // TODO: FIXME: As the file path is changed after used, so we can not use te target.
  const [inputFile, setInputFile] = React.useState('');

  const CheckLocalFile = React.useCallback(() => {
    if (!inputFile) return alert(t('plat.tool.file3'));
    if (!inputFile.startsWith('/data') && !inputFile.startsWith('upload/') && !inputFile.startsWith('./upload/')) {
      return alert(t('plat.tool.file2'));
    }

    const fileExtension = inputFile.slice(inputFile.lastIndexOf('.'));
    if (!MediaSource.exts.includes(fileExtension)) return alert(`${t('plat.tool.file4')}: ${MediaSource.exts.join(', ')}`);

    axios.post(`/terraform/v1/ffmpeg/vlive/server`, {
      file: inputFile,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      let apiUrl = '/terraform/v1/ffmpeg/vlive/source';
      if (endpoint === 'dubbing') apiUrl = '/terraform/v1/dubbing/source';

      console.log(`${t('plat.tool.file5')}，${JSON.stringify(res.data.data)}`);
      const localFileObj = res.data.data;
      const files = [{name: localFileObj.name, path: inputFile, size: localFileObj.size, uuid: localFileObj.uuid, target: localFileObj.target, type: "file"}];
      axios.post(apiUrl, {
        platform, files,
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        console.log(`${t('plat.tool.file6')}，${JSON.stringify(res.data.data)}`);
        setVLiveFiles(res.data.data.files);
      }).catch(handleError);
    }).catch(handleError);
  }, [t, inputFile, handleError, platform, setVLiveFiles, endpoint]);

  return (<>
    <Form.Control as="div">
      {!vLiveFiles?.length ? <>
        <Row>
          <Col>
            <Form.Control type="text" defaultValue={inputFile} placeholder={t('plat.tool.file7')} onChange={e => setInputFile(e.target.value)} />
          </Col>
          <Col xs="auto">
            <Button variant="primary" onClick={CheckLocalFile}>{t('helper.submit')}</Button>
          </Col>
        </Row></> : <></>
      }
      {vLiveFiles?.length ? <VLiveFileList files={vLiveFiles} onChangeFiles={(e) => setVLiveFiles(null)}/> : <></>}
    </Form.Control>
  </>);
}

function VLiveFileUploader({platform, endpoint, vLiveFiles, setVLiveFiles}) {
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const updateSources = React.useCallback((platform, files, setFiles) => {
    if (!files?.length) return alert(t('plat.tool.upload2'));

    let apiUrl = '/terraform/v1/ffmpeg/vlive/source';
    if (endpoint === 'dubbing') apiUrl = '/terraform/v1/dubbing/source';

    axios.post(apiUrl, {
      platform, files: files.map(f => {
        return {name: f.name, path: f.name, size: f.size, uuid: f.uuid, target: f.target, type: "upload"};
      }),
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      console.log(`${t('plat.tool.upload3')}, ${JSON.stringify(res.data.data)}`);
      setFiles(res.data.data.files);
    }).catch(handleError);
  }, [t, handleError, endpoint]);

  return (<>
    <Form.Control as='div'>
      {!vLiveFiles?.length ? <FileUploader onFilesUploaded={(files) => updateSources(platform, files, setVLiveFiles)}/> : <></>}
      {vLiveFiles?.length ? <VLiveFileList files={vLiveFiles} onChangeFiles={(e) => setVLiveFiles(null)}/> : <></>}
    </Form.Control>
  </>);
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

export function VLiveFileFormatInfo({file}) {
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
