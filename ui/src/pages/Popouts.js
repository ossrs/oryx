//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from 'react';
import {useSearchParams} from "react-router-dom";
import {useTranslation} from "react-i18next";
import {useErrorHandler} from "react-error-boundary";
import {Alert, Button} from "react-bootstrap";
import Container from "react-bootstrap/Container";
import axios from "axios";
import {Token} from "../utils";
import {AITalkErrorLogPanel, AITalkTipLogPanel} from "../components/AITalk";

export default function Popouts() {
  const [searchParams] = useSearchParams();

  React.useEffect(() => {
    const app = searchParams.get('app');
    const popout = searchParams.get('popout');
    const room = searchParams.get('room');
    console.log(`?app=ai-talk, current=${app}, The popout application`);
    console.log(`?popout=1, current=${popout}, Whether enable popout mode.`);
    if (app === 'ai-talk') {
      console.log(`?room=room-uuid, current=${room}, The room uuid for ai-talk.`);
    }
  }, [searchParams]);

  const app = searchParams.get('app');
  if (app === 'ai-talk') {
    return <PopoutAITalk {...{roomUUID: searchParams.get('room')}}/>;
  } else {
    return <>Invalid app {app}</>;
  }
}

function PopoutAITalk({roomUUID}) {
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const isMobile = false; // For popout, always PC, not mobile.

  // The player ref, to access the audio player.
  const playerRef = React.useRef(null);
  const [requesting, setRequesting] = React.useState(false);
  const [robotReady, setRobotReady] = React.useState(false);

  // The uuid and robot in stage, which is unchanged after stage started.
  const [stageUUID, setStageUUID] = React.useState(null);
  const [stageRobot, setStageRobot] = React.useState(null);
  const [stagePopoutUUID, setStagePopoutUUID] = React.useState(null);

  // Possible value is 1: yes, -1: no, 0: undefined.
  const [needUserStart, setNeedUserStart] = React.useState(0);
  const [errorLogs, setErrorLogs] = React.useState([]);
  const [traceCount, setTraceCount] = React.useState(0);
  const [traceLogs, setTraceLogs] = React.useState([]);
  const [tipLogs, setTipLogs] = React.useState([]);

  // The refs, about the logs and audio chunks model.
  const ref = React.useRef({
    count: 0,
    isRecording: false,
    recordStarttime: null,
    stopHandler: null,
    mediaStream: null,
    mediaRecorder: null,
    audioChunks: [],
    errorLogs: [],
    traceLogs: [],
    tipsLogs: [],
    traceCount: 0
  });

  const errorLog = React.useCallback((msg) => {
    const rid = `id-${Math.random().toString(16).slice(-4)}${new Date().getTime().toString(16).slice(-4)}`;
    ref.current.errorLogs = [...ref.current.errorLogs, {id: rid, msg}];
    setErrorLogs(ref.current.errorLogs);
  }, [setErrorLogs, ref]);

  const traceLog = React.useCallback((role, msg, variant) => {
    setTraceCount(++ref.current.traceCount);

    // Merge to last log with the same role.
    if (ref.current.traceLogs.length > 0) {
      const last = ref.current.traceLogs[ref.current.traceLogs.length - 1];
      if (last.role === role) {
        last.msg = `${last.msg}${msg}`;
        setTraceLogs([...ref.current.traceLogs]);
        return;
      }
    }

    const rid = `id-${Math.random().toString(16).slice(-4)}${new Date().getTime().toString(16).slice(-4)}`;
    ref.current.traceLogs = [...ref.current.traceLogs, {id: rid, role, msg, variant}];
    setTraceLogs(ref.current.traceLogs);
  }, [setTraceLogs, ref, setTraceCount]);

  const tipLog = React.useCallback((title, msg) => {
    const rid = `id-${Math.random().toString(16).slice(-4)}${new Date().getTime().toString(16).slice(-4)}`;
    ref.current.tipsLogs = [...ref.current.tipsLogs, {id: rid, title, msg, created: new Date()}];
    setTipLogs(ref.current.tipsLogs);
  }, [setTipLogs, ref]);

  const removeTipLog = React.useCallback((log) => {
    const index = ref.current.tipsLogs.findIndex((l) => l.id === log.id);
    ref.current.tipsLogs.splice(index, 1);
    setTipLogs([...ref.current.tipsLogs]);
  }, [setTipLogs, ref]);

  const removeErrorLog = React.useCallback((log) => {
    const index = ref.current.errorLogs.findIndex((l) => l.id === log.id);
    ref.current.errorLogs.splice(index, 1);
    setErrorLogs([...ref.current.errorLogs]);
  }, [setErrorLogs, ref]);

  // Scroll the log panel.
  const endPanelRef = React.useRef(null);
  React.useEffect(() => {
    if (!robotReady || !endPanelRef?.current) return;
    console.log(`Logs setup to end, height=${endPanelRef.current.scrollHeight}, tips=${tipLogs.length}`);
    endPanelRef.current.scrollIntoView({ behavior: "smooth", block: "end" });
  }, [robotReady, endPanelRef, tipLogs]);

  // Request server to create a new popout from source stage.
  React.useEffect(() => {
    if (!roomUUID) return;

    console.log(`Start: Create a new stage`);
    axios.post('/terraform/v1/ai-talk/popout/start', {
      room: roomUUID,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      console.log(`Start: Create popout success: ${JSON.stringify(res.data.data)}`);
      setStageUUID(res.data.data.sid);
      setStagePopoutUUID(res.data.data.spid);
      setStageRobot(res.data.data.robot);
    }).catch(handleError);
  }, [handleError, roomUUID, setStagePopoutUUID, setStageRobot]);

  // Try to start the robot automatically, for OBS.
  React.useEffect(() => {
    if (!stageUUID || !stageRobot) return;

    const listener = () => {
      playerRef.current.removeEventListener('ended', listener);

      setNeedUserStart(1);
      setRobotReady(true);
      console.log(`Stage started, AI is ready, sid=${stageUUID}`);
    };
    playerRef.current.addEventListener('ended', listener);

    playerRef.current.src = `/terraform/v1/ai-talk/stage/examples/${stageRobot.voice}?sid=${stageUUID}`;
    playerRef.current.play().catch((error) => {
      setNeedUserStart(-1);
    });
  }, [t, errorLog, stageUUID, stageRobot, setNeedUserStart]);

  // Requires user to start the robot manually, for Chrome.
  const startChatting = React.useCallback(() => {
    setRequesting(true);

    const listener = () => {
      playerRef.current.removeEventListener('ended', listener);

      setRobotReady(true);
      setRequesting(false);
      console.log(`Stage started, AI is ready, sid=${stageUUID}`);
    };
    playerRef.current.addEventListener('ended', listener);

    playerRef.current.src = `/terraform/v1/ai-talk/stage/examples/${stageRobot.voice}?sid=${stageUUID}`;
    playerRef.current.play().catch(error => {
      errorLog(`${t('lr.room.speaker')}: ${error}`);
      setRequesting(false);
    });
  }, [t, errorLog, stageUUID, stageRobot, setRobotReady, setRequesting]);

  // When robot is ready, show tip logs, and cleanup timeout tips.
  React.useEffect(() => {
    if (!robotReady) return;
    tipLog('Usage', t('lr.room.popout'));

    const timer = setInterval(() => {
      const tipsLogs = [...ref.current.tipsLogs];
      tipsLogs.forEach((log) => {
        if (new Date() - log.created > 10 * 1000) {
          removeTipLog(log);
        }
      });
    }, 500);
    return () => clearInterval(timer);
  }, [t, robotReady, tipLog, isMobile, ref, removeTipLog]);

  // When robot is ready, start query and play all text and voices.
  const refRequest = React.useRef({
    requesting: false,
  });
  // The popouts request worker.
  const requestPopouts = React.useCallback(async () => {
    if (!robotReady || !stageUUID || !stagePopoutUUID) return;
    if (refRequest.current.requesting) return;
    refRequest.current.requesting = true;

    try {
      const msgs = await new Promise((resolve, reject) => {
        axios.post('/terraform/v1/ai-talk/popout/query', {
          sid: stageUUID, spid: stagePopoutUUID,
        }, {
          headers: Token.loadBearerHeader(),
        }).then(res => {
          const ts = new Date().toISOString().split('T')[1].split('Z')[0];
          console.log(`Start: Query popout success at ${ts}: ${JSON.stringify(res.data.data)}`);
          resolve(res.data.data.msgs);
        }).catch(handleError);
      });

      if (!msgs?.length) return;
      for (let i = 0; i < msgs.length; i++) {
        const msg = msgs[i];
        if (msg.role === 'user') {
          traceLog('You', msg.msg, 'primary');
          return;
        }

        const audioSegmentUUID = msg.asid;
        traceLog('Bot', msg.msg, 'success');

        // Play the AI generated audio.
        await new Promise(resolve => {
          const url = `/terraform/v1/ai-talk/popout/tts?sid=${stageUUID}&spid=${stagePopoutUUID}&asid=${audioSegmentUUID}`;
          console.log(`TTS: Playing ${url}`);

          const listener = () => {
            playerRef.current.removeEventListener('ended', listener);
            console.log(`TTS: Played ${url} done.`);
            resolve();
          };
          playerRef.current.addEventListener('ended', listener);

          playerRef.current.src = url;
          playerRef.current.play().catch(error => {
            console.log(`TTS: Play ${url} failed: ${error}`);
            resolve();
          });
        });

        // Remove the AI generated audio.
        await new Promise((resolve, reject) => {
          axios.post('/terraform/v1/ai-talk/popout/remove', {
            sid: stageUUID, spid: stagePopoutUUID, asid: audioSegmentUUID,
          }, {
            headers: Token.loadBearerHeader(),
          }).then(res => {
            console.log(`TTS: Audio removed: ${audioSegmentUUID}`);
            resolve();
          }).catch(error => reject(error));
        });
      }
    } finally {
      refRequest.current.requesting = false;
    }
  }, [handleError, stageUUID, stagePopoutUUID, traceLog, refRequest, robotReady]);
  // Try to request popouts util end.
  React.useEffect(() => {
    if (!robotReady) return;
    const timer = setInterval(async () => {
      requestPopouts().catch(handleError);
    }, 1000);
    return () => clearInterval(timer);
  }, [robotReady, requestPopouts, handleError]);

  return (
    <Container fluid>
      <p></p>
      <div>
        {needUserStart === -1 ?
          <Button disabled={requesting} variant="primary" type="submit" onClick={startChatting}>
            {t('lr.room.talk')}
          </Button> : ''}
        <div><audio ref={playerRef} controls={true} hidden='hidden' /></div>
        <AITalkErrorLogPanel {...{errorLogs, removeErrorLog}} />
        <AITalkTipLogPanel {...{tipLogs, removeTipLog}} />
        <AITalkTraceLogPanelPopout {...{traceLogs, traceCount}} />
        <div ref={endPanelRef}></div>
      </div>
    </Container>
  );
}

export function AITalkTraceLogPanelPopout({traceLogs, traceCount}) {
  // Scroll the log panel.
  const logPanelRef = React.useRef(null);
  React.useEffect(() => {
    if (!logPanelRef?.current) return;
    console.log(`Logs scroll to end, height=${logPanelRef.current.scrollHeight}, logs=${traceLogs.length}, count=${traceCount}`);
    logPanelRef.current.scrollTo(0, logPanelRef.current.scrollHeight);
  }, [traceLogs, logPanelRef, traceCount]);

  return (
    <div className='ai-talk-trace-logs-popout' ref={logPanelRef}>
      {traceLogs.map((log) => {
        return (
          <Alert key={log.id} variant={log.variant}>
            {log.role}: {log.msg}
          </Alert>
        );
      })}
    </div>
  );
}
