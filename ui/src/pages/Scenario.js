//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import {useSearchParams} from "react-router-dom";
import {Container, Tabs, Tab} from "react-bootstrap";
import React from "react";
import ScenarioLiveStreams from './ScenarioLive';
import useUrls from "../components/UrlGenerator";
import ScenarioForward from './ScenarioForward';
import {SrsErrorBoundary} from "../components/SrsErrorBoundary";
import ScenarioTutorials from './ScenarioTutorials';
import {useTranslation} from "react-i18next";
import {useSrsLanguage} from "../components/LanguageSwitch";
import ScenarioRecord from "./ScenarioRecord";
import ScenarioVLive from "./ScenarioVLive";
import {ScenarioVxOthers} from "./ScenarioOthers";
import ScenarioTranscode from "./ScenarioTranscode";
import ScenarioTranscript from "./ScenarioTranscript";
import ScenarioLiveRoom from "./ScenarioLiveRoom";

export default function Scenario() {
  const [searchParams] = useSearchParams();
  const [defaultActiveTab, setDefaultActiveTab] = React.useState();
  const language = useSrsLanguage();

  React.useEffect(() => {
    const tab = searchParams.get('tab') || 'tutorials';
    console.log(`?tab=tutorials|live|stream|record|vlive|transcode|transcript|others, current=${tab}, Select the tab to render`);
    setDefaultActiveTab(tab);
  }, [searchParams, language]);

  return (
    <SrsErrorBoundary>
      { defaultActiveTab && <ScenarioImpl {...{defaultActiveTab}} /> }
    </SrsErrorBoundary>
  );
}

function ScenarioImpl({defaultActiveTab}) {
  const [activeTab, setActiveTab] = React.useState(defaultActiveTab);
  const setSearchParams = useSearchParams()[1];
  const {t} = useTranslation();
  const urls = useUrls();

  const onSelectTab = React.useCallback((k) => {
    setSearchParams({'tab': k});
    setActiveTab(k);
  }, [setSearchParams, setActiveTab]);

  return (
    <>
      <p></p>
      <Container fluid>
        <Tabs defaultActiveKey={activeTab} id="tab0" className="mb-3" onSelect={(k) => onSelectTab(k)}>
          <Tab eventKey="tutorials" title={t('scenario.tutorials')}>
            {activeTab === 'tutorials' && <ScenarioTutorials/>}
          </Tab>
          <Tab eventKey="live" title={t('scenario.live')}>
            {activeTab === 'live' && <ScenarioLiveStreams {...{urls}} />}
          </Tab>
          <Tab eventKey="stream" title={t('scenario.stream')}>
            {activeTab === 'stream' && <ScenarioLiveRoom/>}
          </Tab>
          <Tab eventKey="forward" title={t('scenario.forward')}>
            {activeTab === 'forward' && <ScenarioForward/>}
          </Tab>
          <Tab eventKey="record" title={t('scenario.record')}>
            {activeTab === 'record' && <ScenarioRecord/>}
          </Tab>
          <Tab eventKey="vlive" title={t('scenario.vlive')}>
            {activeTab === 'vlive' && <ScenarioVLive/>}
          </Tab>
          <Tab eventKey="transcode" title={t('scenario.transcode')}>
            {activeTab === 'transcode' && <ScenarioTranscode {...{urls}} />}
          </Tab>
          <Tab eventKey="transcript" title={t('transcript.title')}>
            {activeTab === 'transcript' && <ScenarioTranscript/>}
          </Tab>
          <Tab eventKey="others" title={t('scenario.others')}>
            {activeTab === 'others' && <ScenarioVxOthers {...{urls}} />}
          </Tab>
        </Tabs>
      </Container>
    </>
  );
}

