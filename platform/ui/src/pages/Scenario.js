import {useSearchParams} from "react-router-dom";
import {Container, Tabs, Tab} from "react-bootstrap";
import React from "react";
import {Clipboard, Token} from "../utils";
import axios from "axios";
import ScenarioDvr from './ScenarioDvr';
import ScenarioSrt from './ScenarioSrt';
import ScenarioLive from './ScenarioLive';
import useUrls from "../components/UrlGenerator";
import ScenarioVod from './ScenarioVod';
import ScenarioForward from './ScenarioForward';
import {useErrorHandler} from 'react-error-boundary';
import {SrsErrorBoundary} from "../components/SrsErrorBoundary";
import ScenarioTutorials from './ScenarioTutorials';
import {useTranslation} from "react-i18next";
import {useSrsLanguage} from "../components/LanguageSwitch";
import ScenarioRecord from "./ScenarioRecord";

export default function Scenario() {
  const [searchParams] = useSearchParams();
  const [defaultActiveTab, setDefaultActiveTab] = React.useState();
  const [defaultActiveChildTab, setDefaultActiveChildTab] = React.useState();
  const language = useSrsLanguage();

  React.useEffect(() => {
    const tab = searchParams.get('tab') || 'tutorials';
    const ctab = searchParams.get('ctab') || 'record';
    console.log(`?tab=tutorials|live|srt|rgroup|source, current=${tab}, Select the tab to render`);
    console.log(`?ctab=record|dvr|vod, current=${tab}, Select the child tab to render`);
    setDefaultActiveTab(tab);
    setDefaultActiveChildTab(ctab);
  }, [searchParams, language]);

  return (
    <SrsErrorBoundary>
      { defaultActiveTab && <ScenarioImpl {...{defaultActiveTab, defaultActiveChildTab}} /> }
    </SrsErrorBoundary>
  );
}

function ScenarioImpl({defaultActiveTab, defaultActiveChildTab}) {
  const [secret, setSecret] = React.useState();
  const [streamName, setStreamName] = React.useState('livestream');
  const [activeTab, setActiveTab] = React.useState(defaultActiveTab);
  const [activeChildTab, setActiveChildTab] = React.useState(defaultActiveChildTab);
  const setSearchParams = useSearchParams()[1];
  const handleError = useErrorHandler();
  const urls = useUrls({secret, streamName});
  const {t} = useTranslation();

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/hooks/srs/secret/query', {
      ...token,
    }).then(res => {
      setSecret(res.data.data);
      console.log(`Status: Query ok, secret=${JSON.stringify(res.data.data)}`);
    }).catch(handleError);
  }, [handleError]);

  const onSelectTab = React.useCallback((k) => {
    setSearchParams({'tab': k});
    setActiveTab(k);
  }, [setSearchParams]);

  const onSelectChildTab = React.useCallback((pk, k) => {
    setSearchParams({'tab': pk, 'ctab': k});
    setActiveTab(pk);
    setActiveChildTab(k);
  }, [setSearchParams]);

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

  return (
    <>
      <p></p>
      <Container>
        <Tabs defaultActiveKey={activeTab} id="tab0" className="mb-3" onSelect={(k) => onSelectTab(k)}>
          <Tab eventKey="tutorials" title={t('scenario.tutorials')}>
            {activeTab === 'tutorials' && <ScenarioTutorials/>}
          </Tab>
          <Tab eventKey="live" title={t('scenario.live')}>
            {activeTab === 'live' && <ScenarioLive {...{updateStreamName, copyToClipboard, urls}} />}
          </Tab>
          <Tab eventKey="srt" title={t('scenario.srt')}>
            {activeTab === 'srt' && <ScenarioSrt {...{updateStreamName, copyToClipboard, urls}} />}
          </Tab>
          <Tab eventKey="forward" title={t('scenario.restream')}>
            {activeTab === 'forward' && <ScenarioForward/>}
          </Tab>
          <Tab eventKey="rgroup" title={t('scenario.rgroup')}>
            <Tabs defaultActiveKey={activeChildTab} id="ctab0" className="mb-3"
                  onSelect={(k) => onSelectChildTab('rgroup', k)}>
              <Tab eventKey="record" title={t('scenario.record')}>
                {activeChildTab === 'record' && <ScenarioRecord/>}
              </Tab>
              <Tab eventKey="dvr" title={t('scenario.dvr')}>
                {activeChildTab === 'dvr' && <ScenarioDvr/>}
              </Tab>
              <Tab eventKey="vod" title={t('scenario.vod')}>
                {activeChildTab === 'vod' && <ScenarioVod/>}
              </Tab>
            </Tabs>
          </Tab>
        </Tabs>
      </Container>
    </>
  );
}

