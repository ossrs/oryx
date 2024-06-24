//
// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
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

  React.useEffect(() => {
    const ctab = searchParams.get('ctab') || 'other';
    console.log(`?ctab=other|deprecated|srt|dvr|vod, current=${ctab}, Select the child tab to render`);
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
          {activeChildTab === 'other' && <ScenarioOther {...{urls}} />}
        </Tab>
        <Tab eventKey="deprecated" title={t('scenario.deprecated')}>
          {activeChildTab === 'deprecated' && <ScenarioDeprecated/>}
        </Tab>
      </Tabs>
    }
  </>;
}

function ScenarioOther({urls}) {
  const {t} = useTranslation();
  const language = useSrsLanguage();
  const isZh = language === 'zh';

  const copyToClipboard = React.useCallback((e, text) => {
    e.preventDefault();

    Clipboard.copy(text).then(() => {
      alert(t('helper.copyOk'));
    }).catch((err) => {
      alert(`${t('helper.copyFail')} ${err}`);
    });
  }, [t]);

  return <>
    <Accordion defaultActiveKey='0'>
      <Accordion.Item eventKey="0">
        <Accordion.Header>{isZh ? '场景介绍' : 'Introduction'}</Accordion.Header>
        <Accordion.Body>
          <div>
            {isZh ? '其他较少使用的低频场景，请用其他常用场景替代。' : 'Other less common used scenarios, please use other common scenarios instead.'}
            <p></p>
          </div>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{t('scenario.srt')}</Accordion.Header>
        <Accordion.Body>
          <p>{isZh ? '可使用基本的推拉流场景实现，已经支持了SRT协议的教程。' : 'You can implement basic streaming scenarios with a tutorial that already supports the SRT protocol.'}</p>
          <ScenarioSrt {...{copyToClipboard, urls}} />
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  </>;
}

function ScenarioDeprecated() {
  const {t} = useTranslation();
  const language = useSrsLanguage();
  const isZh = language === 'zh';

  return <>
    <Accordion defaultActiveKey='0'>
      <Accordion.Item eventKey="0">
        <Accordion.Header>{isZh ? '场景介绍' : 'Introduction'}</Accordion.Header>
        <Accordion.Body>
          <div>
            {isZh ? '其他废弃场景，未来会移除，请不要使用。' : 'Other deprecated scenarios, will be removed in the future, please do not use.'}
            <p></p>
          </div>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{t('scenario.dvr')}</Accordion.Header>
        <Accordion.Body>
          <ScenarioRecordCos/>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>{t('scenario.vod')}</Accordion.Header>
        <Accordion.Body>
          <ScenarioRecordVod/>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  </>;
}

