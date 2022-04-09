import {useSearchParams} from "react-router-dom";
import {Container, Tabs, Tab} from "react-bootstrap";
import React from "react";
import {Clipboard, Token} from "../utils";
import axios from "axios";
import ScenarioDvr from './ScenarioDvr';
import ScenarioSource from './ScenarioSource';
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

export default function Scenario() {
  const [searchParams] = useSearchParams();
  const [defaultActiveTab, setDefaultActiveTab] = React.useState();
  const language = useSrsLanguage();

  React.useEffect(() => {
    const tab = searchParams.get('tab') || 'tutorials';
    console.log(`?tab=tutorials|live|srt|dvr|source, current=${tab}, Select the tab to render`);
    setDefaultActiveTab(tab);
  }, [searchParams, language]);

  return (
    <SrsErrorBoundary>
      { defaultActiveTab && <ScenarioImpl defaultActiveTab={defaultActiveTab} /> }
    </SrsErrorBoundary>
  );
}

function ScenarioImpl({defaultActiveTab}) {
  const [secret, setSecret] = React.useState();
  const [streamName, setStreamName] = React.useState('livestream');
  const [activeTab, setActiveTab] = React.useState(defaultActiveTab);
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
        <Tabs defaultActiveKey={activeTab} id="uncontrolled-tab-example" className="mb-3" onSelect={(k) => onSelectTab(k)}>
          <Tab eventKey="tutorials" title={t('scenario.tutorials')}>
            { activeTab === 'tutorials' && <ScenarioTutorials /> }
          </Tab>
          <Tab eventKey="live" title={t('scenario.live')}>
            { activeTab === 'live' && <ScenarioLive {...{updateStreamName, copyToClipboard, urls}} /> }
          </Tab>
          <Tab eventKey="srt" title={t('scenario.srt')}>
            { activeTab === 'srt' && <ScenarioSrt {...{updateStreamName, copyToClipboard, urls}} /> }
          </Tab>
          <Tab eventKey="forward" title={t('scenario.restream')}>
            { activeTab === 'forward' && <ScenarioForward /> }
          </Tab>
          <Tab eventKey="dvr" title={t('scenario.dvr')}>
            { activeTab === 'dvr' && <ScenarioDvr /> }
          </Tab>
          <Tab eventKey="vod" title={t('scenario.vod')}>
            { activeTab === 'vod' && <ScenarioVod /> }
          </Tab>
          <Tab eventKey="source" title={t('scenario.code')}>
            { activeTab === 'source' && <ScenarioSource /> }
          </Tab>
        </Tabs>
      </Container>
    </>
  );
}

