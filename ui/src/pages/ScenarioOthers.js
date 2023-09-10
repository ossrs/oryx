//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {useSrsLanguage} from "../components/LanguageSwitch";
import {Accordion, Tab, Tabs} from "react-bootstrap";
import {useSearchParams} from "react-router-dom";
import {useTranslation} from "react-i18next";
import {useErrorHandler} from "react-error-boundary";
import useUrls from "../components/UrlGenerator";
import {Clipboard, Token} from "../utils";
import axios from "axios";
import ScenarioSrt from "./ScenarioSrt";
import ScenarioRecordCos from "./ScenarioRecordCos";
import ScenarioRecordVod from "./ScenarioRecordVod";

export function ScenarioVxOthers() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [activeChildTab, setActiveChildTab] = React.useState();
  const language = useSrsLanguage();
  const {t} = useTranslation();

  const [secret, setSecret] = React.useState();
  const [streamName, setStreamName] = React.useState('livestream');
  const handleError = useErrorHandler();
  const urls = useUrls({secret, streamName});

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/hooks/srs/secret/query', {
      ...token,
    }).then(res => {
      setSecret(res.data.data);
      console.log(`Status: Query ok, secret=${JSON.stringify(res.data.data)}`);
    }).catch(handleError);
  }, [handleError]);

  const updateStreamName = React.useCallback(() => {
    setStreamName(Math.random().toString(16).slice(-6).split('').map(e => {
      return (e >= '0' && e <= '9') ? String.fromCharCode('a'.charCodeAt(0) + (parseInt(Math.random() * 16 + e) % 25)) : e;
    }).join(''));
    alert(t('helper.changeStream'));
  }, [t]);

  const copyToClipboard = React.useCallback((e, text) => {
    e.preventDefault();

    Clipboard.copy(text).then(() => {
      alert(t('helper.copyOk'));
    }).catch((err) => {
      alert(`${t('helper.copyFail')} ${err}`);
    });
  }, [t]);

  React.useEffect(() => {
    const ctab = searchParams.get('ctab') || 'other';
    console.log(`?ctab=other|srt|dvr|vod, current=${ctab}, Select the child tab to render`);
    setActiveChildTab(ctab);
  }, [searchParams, language, setActiveChildTab]);

  const onSelectChildTab = React.useCallback((k) => {
    setSearchParams({...searchParams, 'tab': 'ogroup', 'ctab': k});
    setActiveChildTab(k);
  }, [searchParams, setSearchParams, setActiveChildTab]);

  return <>
    {activeChildTab &&
      <Tabs defaultActiveKey={activeChildTab} id="ctab0" className="mb-3"
            onSelect={(k) => onSelectChildTab(k)}>
        <Tab eventKey="other" title={t('scenario.other')}>
          {activeChildTab === 'other' && <ScenarioOther/>}
        </Tab>
        <Tab eventKey="srt" title={t('scenario.srt')}>
          {activeChildTab === 'srt' && <ScenarioSrt {...{updateStreamName, copyToClipboard, urls}} />}
        </Tab>
        <Tab eventKey="dvr" title={t('scenario.dvr')}>
          {activeChildTab === 'dvr' && <ScenarioRecordCos/>}
        </Tab>
        <Tab eventKey="vod" title={t('scenario.vod')}>
          {activeChildTab === 'vod' && <ScenarioRecordVod/>}
        </Tab>
      </Tabs>
    }
  </>;
}

function ScenarioOther() {
  const language = useSrsLanguage();
  if (language === 'zh') {
    return <>
      <Accordion defaultActiveKey='0'>
        <Accordion.Item eventKey="0">
          <Accordion.Header>场景介绍</Accordion.Header>
          <Accordion.Body>
            <div>
              其他非常用场景。
              <p></p>
            </div>
          </Accordion.Body>
        </Accordion.Item>
      </Accordion>
    </>;
  }
  return <>
    <Accordion defaultActiveKey='0'>
      <Accordion.Item eventKey="0">
        <Accordion.Header>Introduction</Accordion.Header>
        <Accordion.Body>
          <div>
            Other less common scenarios.
            <p></p>
          </div>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  </>;
}

