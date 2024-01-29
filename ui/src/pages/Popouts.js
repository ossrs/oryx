//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from 'react';
import {useLocation, useSearchParams} from "react-router-dom";
import {useErrorHandler} from "react-error-boundary";
import {Spinner} from "react-bootstrap";
import Container from "react-bootstrap/Container";
import axios from "axios";
import {Locale, Token} from "../utils";
import {AITalkAssistantPanel, AITalkChatOnlyPanel} from "../components/AITalk";
import {useTranslation} from "react-i18next";
import resources from "../resources/locale";

export default function Popouts() {
  const location = useLocation();
  const {i18n} = useTranslation();
  const [initialized, setInitialized] = React.useState(false);

  // Switch language for popout, because it does not use navigator, so there is no
  // LanguageSwitch to do this.
  React.useEffect(() => {
    if (!i18n || !location) return;

    const lang = location.pathname.split('/')[1];

    // Ignore if invalid language.
    if (!lang || !Object.keys(resources).includes(lang)) {
      return;
    }

    // Change to language in url.
    if (Locale.current() !== lang) {
      i18n.changeLanguage(lang);
    }

    setInitialized(true);
  }, [setInitialized, i18n, location]);

  if (!initialized) {
    return <>
      <Spinner animation="border" variant="primary" size='sm'></Spinner>&nbsp;
      Initializing...
    </>;
  }
  return <PopoutsImpl/>;
}

function PopoutsImpl() {
  const handleError = useErrorHandler();
  const [searchParams] = useSearchParams();
  const [loading, setLoading] = React.useState(true);

  React.useEffect(() => {
    const app = searchParams.get('app');
    const roomToken = searchParams.get('roomToken');
    const popout = searchParams.get('popout');
    const room = searchParams.get('room');
    const assistant = searchParams.get('assistant');
    const username = searchParams.get('username');
    const userLanguage = searchParams.get('language');
    console.log(`?app=ai-talk, current=${app}, The popout application`);
    console.log(`?roomToken=xxx, current=${roomToken?.length}B, The popout token for each room`);
    console.log(`?popout=1, current=${popout}, Whether enable popout mode.`);
    if (app === 'ai-talk') {
      console.log(`?room=room-uuid, current=${room}, The room uuid for ai-talk.`);
      console.log(`?assistant=0, current=${assistant}, Whether popout the assistant, allow user to talk.`);
      if (assistant === '1') {
        console.log(`?username=xxx, current=${username}, The username of stream host.`);
        console.log(`?language=xxx, current=${userLanguage}, The language of user.`);
      }
    }

    if (!app) throw new Error(`no app`);
    if (!roomToken) throw new Error(`no room token`);
    if (app === 'ai-talk' && !room) throw new Error(`no room id`);

    axios.post('/terraform/v1/ai-talk/stage/verify', {
      room: searchParams.get('room'), roomToken: searchParams.get('roomToken'),
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      setLoading(false);
      console.log(`Verify room token ok`);
    }).catch(handleError);
  }, [handleError, searchParams, setLoading]);

  if (loading) {
    return <>
      <Spinner animation="border" variant="primary" size='sm'></Spinner>&nbsp;
      Loading...
    </>;
  }
  const app = searchParams.get('app');
  if (app === 'ai-talk') {
    const assistant = searchParams.get('assistant') === '1';
    const roomUUID = searchParams.get('room');
    const roomToken = searchParams.get('roomToken');
    const username = searchParams.get('username');
    const userLanguage = searchParams.get('language');
    return (
      <Container fluid>
        <p></p>
        {assistant ?
          <AITalkAssistantPanel {...{roomUUID, roomToken, username, userLanguage, fullscreen: true}}/> :
          <AITalkChatOnlyPanel {...{roomUUID, roomToken}}/>}
      </Container>
    );
  } else {
    return <>Invalid app {app}</>;
  }
}
