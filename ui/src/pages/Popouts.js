//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from 'react';
import {useSearchParams} from "react-router-dom";
import {useErrorHandler} from "react-error-boundary";
import {Spinner} from "react-bootstrap";
import Container from "react-bootstrap/Container";
import axios from "axios";
import {Token} from "../utils";
import {AITalkAssistantPanel, AITalkChatPanel} from "../components/AITalk";

export default function Popouts() {
  const handleError = useErrorHandler();
  const [searchParams] = useSearchParams();
  const [loading, setLoading] = React.useState(true);

  React.useEffect(() => {
    const app = searchParams.get('app');
    const roomToken = searchParams.get('roomToken');
    const popout = searchParams.get('popout');
    const room = searchParams.get('room');
    const assistant = searchParams.get('assistant');
    console.log(`?app=ai-talk, current=${app}, The popout application`);
    console.log(`?roomToken=xxx, current=${roomToken?.length}B, The popout token for each room`);
    console.log(`?popout=1, current=${popout}, Whether enable popout mode.`);
    if (app === 'ai-talk') {
      console.log(`?room=room-uuid, current=${room}, The room uuid for ai-talk.`);
      console.log(`?assistant=0, current=${assistant}, Whether popout the assistant, allow user to talk.`);
    }

    if (!app) throw new Error(`no app`);
    if (!roomToken) throw new Error(`no room token`);
    if (app === 'ai-talk' && !room) throw new Error(`no room id`);

    axios.post('/terraform/v1/ai-talk/stage/verify', {
      room: searchParams.get('room'), roomToken: searchParams.get('roomToken'),
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      const {token} = res.data.data;
      Token.updateBearer(token);
      setLoading(false);
      console.log(`Verify room token ok, data=${JSON.stringify(res.data.data)}`);
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
    return (
      <Container fluid>
        <p></p>
        {assistant ?
          <AITalkAssistantPanel {...{roomUUID, fullscreen: true}}/> :
          <AITalkChatPanel {...{roomUUID, roomToken}}/>}
      </Container>
    );
  } else {
    return <>Invalid app {app}</>;
  }
}
