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
import {Clipboard} from "../utils";
import ScenarioSrt from "./ScenarioSrt";
import ScenarioRecordCos from "./ScenarioRecordCos";
import ScenarioRecordVod from "./ScenarioRecordVod";

export function ScenarioVxOthers({urls}) {
  const [searchParams, setSearchParams] = useSearchParams();
  const [activeChildTab, setActiveChildTab] = React.useState();
  const language = useSrsLanguage();
  const {t} = useTranslation();

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
    setSearchParams({...searchParams, 'tab': 'others', 'ctab': k});
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
          {activeChildTab === 'srt' && <ScenarioSrt {...{copyToClipboard, urls}} />}
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

