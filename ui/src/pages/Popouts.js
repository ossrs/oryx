//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from 'react';
import {useSearchParams} from "react-router-dom";

export default function Popouts() {
  const [searchParams] = useSearchParams();

  React.useEffect(() => {
    const app = searchParams.get('app');
    const popout = searchParams.get('popout');
    const room = searchParams.get('room');
    const stage = searchParams.get('stage');
    console.log(`?app=ai-talk, current=${app}, The popout application`);
    console.log(`?popout=1, current=${popout}, Whether enable popout mode.`);
    if (app === 'ai-talk') {
      console.log(`?room=room-uuid, current=${room}, The room uuid for ai-talk.`);
      console.log(`?stage=stage-uuid, current=${stage}, The stage uuid for ai-talk.`);
    }
  }, [searchParams]);

  const app = searchParams.get('app');
  if (app === 'ai-talk') {
    return <PopoutAITalk {...{roomUuid: searchParams.get('room'), stageUuid: searchParams.get('stage')}}/>;
  } else {
    return <>Invalid app {app}</>;
  }
}

function PopoutAITalk({roomUuid, stageUuid}) {
  return <>On the way...</>;
}
