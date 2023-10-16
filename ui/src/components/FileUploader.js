//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {useErrorHandler} from "react-error-boundary";
import {Tools} from "../utils";
import {useTranslation} from "react-i18next";
import moment from "moment";
import {Button, Col, Form, Row, ListGroup} from "react-bootstrap";
import axios from "axios";

export default function FileUploader({onFilesUploaded}) {
  const accept = React.useMemo(() => ['.mp4', '.flv', '.ts'], []);
  const multiple = React.useMemo(() => false, []);

  const handleError = useErrorHandler();
  const {t} = useTranslation();
  const [uploading, setUploading] = React.useState(false);
  const [filesToUpload, setFilesToUpload] = React.useState([]);
  const [filesUploading, setFilesUploading] = React.useState({});

  // For callback to update state, because in callback we can only get the copy, so we need a ref to point to the latest
  // copy of state of variant objects.
  const ref = React.useRef({});
  React.useEffect(() => {
    ref.current.filesProgressing = filesUploading;
  }, [filesUploading]);

  const onUploadFile = React.useCallback((e) => {
    const files = [];
    for (const file of e.target.files) files.push(file);
    setFilesToUpload(files);

    const size = files.map(f => f.size).reduce((a, b) => a + b, 0) / 1024 / 1024;
    console.log(`Uploader: Got ${files.length} files ${size.toFixed(1)}MB to upload, ${JSON.stringify(files.map(f => f.name))}`);
  }, [setFilesToUpload]);

  const uploaVideoFiles = React.useCallback(async () => {
    setUploading(true);

    const uploadFile = async (file, fileInfo) => {
      const starttime = moment();
      return await new Promise((resolve, reject) => {
        console.log(`Uploader: start to upload ${file.name} ${JSON.stringify(fileInfo)}`);

        const formData = new FormData();
        const sessionName = file.name;
        formData.append(sessionName, file);

        axios.post(`/terraform/v1/ffmpeg/vlive/upload/${sessionName}`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          },
          onUploadProgress: (progressData) => {
            //{loaded: 112967680, total: 169659461, timeStamp:257985, type:"progress"};
            progressData.cost = moment() - starttime;
            if (progressData.cost) progressData.speed = progressData.loaded * 1000.0 / progressData.cost; // In Kbps
            if (progressData.total) progressData.percent = progressData.loaded * 1.0 / progressData.total; // [0, 1];
            fileInfo = Tools.merge(fileInfo, {...progressData, cost: (moment() - starttime) || 1});
            setFilesUploading(Tools.copy(ref.current.filesProgressing, [file.name, fileInfo]));
          },
        }).then(res => {
          fileInfo.percent = 1;
          fileInfo.cost = moment() - starttime;
          if (fileInfo.cost) fileInfo.speed = file.size * 1000.0 / fileInfo.cost;

          // Merge the server information.
          fileInfo.uuid = res.data.data.uuid;
          fileInfo.target = res.data.data.target;

          console.log(`Uploader: Upload ${file.name} ok, info=${JSON.stringify(fileInfo)}, cost=${moment() - starttime}ms`);
          resolve(fileInfo);
        }).catch(handleError);
      });
    };

    try {
      const uploadedFiles = [];
      for (const file of filesToUpload) {
        let fileInfo = Tools.merge(ref.current.filesProgressing[file.name], {
          name: file.name, size: file.size,
        });
        try {
          fileInfo = await uploadFile(file, {...fileInfo});
          uploadedFiles.push(fileInfo);
        } finally {
          setFilesUploading(Tools.copy(ref.current.filesProgressing, [file.name, fileInfo]));
        }
      }

      console.log(`Uploader: All files done, callback with ${JSON.stringify(uploadedFiles)}`);
      onFilesUploaded && onFilesUploaded(uploadedFiles);
    } catch (e) {
      handleError(e);
    } finally {
      setUploading(false);
    }
  }, [setUploading, filesToUpload, setFilesUploading, onFilesUploaded, ref, handleError]);

  return (<>
    <Row>
      <Col xs='auto'>
        <FilePicker {...{accept, multiple, disabled: uploading, onChange: onUploadFile}} />
      </Col>
      <Col xs='auto'>
        <Button variant="primary" type="button" disabled={uploading || !filesToUpload?.length}
                onClick={uploaVideoFiles}>{t('helper.upload')}</Button>
      </Col>
    </Row>
    <UploadingFiles filesUploading={filesUploading} filesToUpload={filesToUpload} />
  </>);
}

function FilePicker({accept, multiple, disabled, onChange}) {
  const {t} = useTranslation();
  const [filesToUpload, setFilesToUpload] = React.useState([]);
  const [hover, setHover] = React.useState(false);

  const onUploadFile = React.useCallback((e) => {
    const files = [];
    for (const file of e.target.files) files.push(file);
    setFilesToUpload(files);
    onChange && onChange(e);
  }, [setFilesToUpload, onChange]);

  const filePickerId = `file-picker-${Math.random().toString(16).slice(-6)}`;
  return <>
    <label htmlFor={filePickerId} style={{
      cursor: 'pointer',
      padding: '7px 12px 7px 12px',
      backgroundColor: hover ? '#d9d9d9' : '#e9e9e9',
      borderBottomLeftRadius: '5px',
      borderTopLeftRadius: '5px'
    }} onMouseEnter={() => setHover(true)} onMouseLeave={() => setHover(false)}>
      {t('fp.label')}
    </label>
    <label htmlFor={filePickerId} style={{
      cursor: 'pointer',
      padding: '6px 130px 6px 10px',
      backgroundColor: '#fefefe',
      border: '1px solid #e0e0e0',
      borderBottomRightRadius: '5px',
      borderTopRightRadius: '5px'
    }} onMouseEnter={() => setHover(true)} onMouseLeave={() => setHover(false)}>
      {filesToUpload?.length ? filesToUpload.map(f => f.name).join(', ') : t('fp.nofile')}
    </label>
    <Form.Control style={{display: 'none'}} id={filePickerId} type="file" accept={accept}
                  multiple={multiple} disabled={disabled} onChange={onUploadFile}/>
  </>;
}

function UploadingFiles({filesUploading, filesToUpload}) {
  if (!filesToUpload?.length) return <></>;

  return <Row>
    <Col xs='auto'>
      <ListGroup>
        {filesToUpload.map((f, index) => {
          return <ListGroup.Item key={f.name}>
            {index + 1}: {f.name} {Number(f.size / 1024 / 1024).toFixed(1)}MB &nbsp;
            {filesUploading[f.name]?.speed && `${Number(filesUploading[f.name].speed * 8.0 / 1000 / 1000).toFixed(2)}Mbps`} &nbsp;
            {filesUploading[f.name]?.percent && `${Number(filesUploading[f.name].percent * 100).toFixed(0)}%`} &nbsp;
            {filesUploading[f.name]?.cost && `${Number(filesUploading[f.name].cost / 1000).toFixed(0)}s`} &nbsp;
          </ListGroup.Item>;
        })}
        <ListGroup.Item>
          <b>Total: {Number(filesToUpload.map(f => f.size).reduce((a, b) => a + b, 0) / 1024 / 1024).toFixed(1)}MB</b>
        </ListGroup.Item>
      </ListGroup>
    </Col>
  </Row>;
}
