import React from "react";
import {Alert, Button, Card, Col, Dropdown, Form, InputGroup, Row, Spinner} from "react-bootstrap";
import {useTranslation} from "react-i18next";
import {useErrorHandler} from "react-error-boundary";
import {useIsMobile} from "./IsMobile";
import axios from "axios";
import {Locale, Token} from "../utils";
import * as Icon from "react-bootstrap-icons";

export function AITalkDictationPanel({roomUUID, roomToken, username, userLanguage}) {
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const isMobile = useIsMobile();

  // The timeout in milliseconds.
  const timeoutForMicrophoneTestToRun = 50;
  const timeoutWaitForMicrophoneToClose = 900;
  const timeoutWaitForLastVoice = 700;
  const maxSegmentTime = 10 * 1000; // in ms.

  // The player ref, to access the audio player.
  const playerRef = React.useRef(null);
  const [requesting, setRequesting] = React.useState(false);
  const [robotReady, setRobotReady] = React.useState(false);
  const [processing, setProcessing] = React.useState(false);
  const [micWorking, setMicWorking] = React.useState(false);

  // The uuid and robot in stage, which is unchanged after stage started.
  const [stageUUID, setStageUUID] = React.useState(null);
  const [userID, setUserID] = React.useState(null);
  const [aiAsrEnabled, setAiAsrEnabled] = React.useState();
  const [stageUser, setStageUser] = React.useState(null);
  const [stagePopoutUUID, setStagePopoutUUID] = React.useState(null);
  // Last user input text, from ASR, set to input for user to update it.
  const [userAsrText, setUserAsrText] = React.useState(null);

  const [booting, setBooting] = React.useState(true);
  const [errorLogs, setErrorLogs] = React.useState([]);
  const [traceCount, setTraceCount] = React.useState(0);
  const [traceLogs, setTraceLogs] = React.useState([]);
  const [tipLogs, setTipLogs] = React.useState([]);

  // The refs, about the logs and audio chunks model.
  const ref = React.useRef({
    shouldStop: false,
    mediaStream: null,
    errorLogs: [],
    traceLogs: [],
    tipsLogs: [],
    traceCount: 0
  });
  // When robot is ready, start query and play all text and voices.
  const refRequest = React.useRef({
    requesting: false,
    gotMessages: 0,
    hasPendingMessages: false,
    playingAudio: 0,
  });

  const errorLog = React.useCallback((msg) => {
    const rid = `id-${Math.random().toString(16).slice(-4)}${new Date().getTime().toString(16).slice(-4)}`;
    ref.current.errorLogs = [...ref.current.errorLogs, {id: rid, msg}];
    setErrorLogs(ref.current.errorLogs);
  }, [setErrorLogs, ref]);

  const traceLog = React.useCallback((role, msg, variant, ignoreMerge) => {
    setTraceCount(++ref.current.traceCount);

    // Merge to last log with the same role.
    if (ref.current.traceLogs.length > 0 && !ignoreMerge) {
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

  // The application is started now.
  React.useEffect(() => {
    const fnImpl = async () => {
      // Only allow localhost or https to access microphone.
      const isLo = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
      const isHttps = window.location.protocol === 'https:';
      const securityAllowed = isLo || isHttps;
      securityAllowed || errorLog(t('lr.room.https'));
      console.log(`App started, allowed=${securityAllowed}, lo=${isLo}, https=${isHttps}`);
      if (!securityAllowed) return;

      // Wait util no audio is playing.
      do {
        await new Promise((resolve) => setTimeout(resolve, 300));
      } while (refRequest.current.playingAudio);

      // Try to open the microphone to request permission.
      new Promise(resolve => {
        console.log(`Start: Startup open microphone`);

        navigator.mediaDevices.getUserMedia(
          {audio: true}
        ).then((stream) => {
          console.log(`Start: Microphone opened, try to record`);
          const recorder = new MediaRecorder(stream);

          const audioChunks = [];
          recorder.addEventListener("dataavailable", ({data}) => {
            audioChunks.push(data);
          });
          recorder.addEventListener("stop", async () => {
            // Stop the microphone.
            console.log(`Start: Microphone ok, chunks=${audioChunks.length}, state=${recorder.state}`);
            stream.getTracks().forEach(track => track.stop());
            setTimeout(() => {
              console.log(`Start: Microphone test ok.`);
              resolve();
            }, timeoutWaitForMicrophoneToClose);
          });

          recorder.start();
          setTimeout(() => {
            recorder.stop();
            console.log(`Start: Microphone stopping, state is ${recorder.state}`);
          }, timeoutForMicrophoneTestToRun);
        }).catch(error => errorLog(`${t('lr.room.mic')}: ${error}`));
      }).then(() => {
        setBooting(false);
      });
    };
    fnImpl();
  }, [t, errorLog, setBooting, refRequest]);

  // Request server to create a new stage.
  React.useEffect(() => {
    if (booting) return;

    console.log(`Start: Create a new stage`);
    axios.post('/terraform/v1/ai-talk/stage/start', {
      room: roomUUID, roomToken,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      console.log(`Start: Create stage success: ${JSON.stringify(res.data.data)}`);
      setStageUUID(res.data.data.sid);
      setUserID(res.data.data.userId);
      setAiAsrEnabled(res.data.data.aiAsrEnabled);
    }).catch(handleError);
  }, [handleError, booting, roomUUID, roomToken, setStageUUID, setUserID, setAiAsrEnabled]);

  // Start to chat, set the robot to ready.
  const startChatting = React.useCallback(async (user) => {
    setStageUser(user);

    setRequesting(true);
    const listener = () => {
      playerRef.current.removeEventListener('ended', listener);

      setRobotReady(true);
      setRequesting(false);
      refRequest.current.playingAudio = false;
      console.log(`Stage started, AI is ready, sid=${stageUUID}`);
    };
    playerRef.current.addEventListener('ended', listener);

    console.log('Stage start to play demo audio');
    refRequest.current.playingAudio = true;
    playerRef.current.src = `/terraform/v1/ai-talk/stage/hello-voices/${user.voice}?sid=${stageUUID}`;
    playerRef.current.play().catch(error => {
      errorLog(`${t('lr.room.speaker')}: ${error}`);
      setRequesting(false);
    });
  }, [t, errorLog, stageUUID, setRobotReady, setRequesting, refRequest, setStageUser]);

  // When robot is ready, show tip logs, and cleanup timeout tips.
  React.useEffect(() => {
    if (!robotReady) return;
    tipLog('Usage', isMobile ? t('lr.room.usage') : t('lr.room.usage2'));

    const timer = setInterval(() => {
      const tipsLogs = [...ref.current.tipsLogs];
      tipsLogs.forEach((log) => {
        if (new Date() - log.created > 10 * 1000) {
          removeTipLog(log);
        }
      });
    }, 1000);
    return () => clearInterval(timer);
  }, [t, robotReady, tipLog, isMobile, ref, removeTipLog]);

  // Handle the user input audio data, upload to server.
  const processUserInput = React.useCallback(async(artifact) => {
    const userMayInput = artifact.duration();

    // End conversation, for stat the elapsed time cost accurately.
    const requestUUID = await new Promise((resolve, reject) => {
      axios.post('/terraform/v1/ai-talk/stage/conversation', {
        room: roomUUID, roomToken, sid: stageUUID,
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        console.log(`ASR: Start conversation success, rid=${res.data.data.rid}`);
        resolve(res.data.data.rid);
      }).catch((error) => reject(error));
    });

    // Convert audio from binary to base64 in text.
    const audioBase64Data = await new Promise((resolve) => {
      const reader = new FileReader();
      reader.onload = function() {
        // Remove the data URL prefix, for example, result is:
        //    application/octet-stream;base64,GkXfo59ChoEB.............
        const base64Audio = reader.result.split(',')[1];
        resolve(base64Audio);
      };

      const audioBlob = new Blob(artifact.audioChunks);
      reader.readAsDataURL(audioBlob);
    });

    // Upload the user input audio to the server.
    await new Promise((resolve, reject) => {
      console.log(`ASR: Uploading ${artifact.audioChunks.length} chunks, stage=${stageUUID}, user=${userID}`);

      axios.post('/terraform/v1/ai-talk/stage/upload', {
        room: roomUUID, roomToken, sid: stageUUID, rid: requestUUID, userId: userID,
        umi: userMayInput, audio: audioBase64Data,
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        console.log(`ASR: Upload success: ${res.data.data.rid} ${res.data.data.asr}`);
        resolve(res.data.data.rid);
        setUserAsrText(res.data.data.asr);
      }).catch((error) => reject(error));
    });

    // Get the AI generated audio from the server.
    console.log(`TTS: Requesting ${requestUUID} response audios, rid=${requestUUID}`);
    while (true) {
      const resp = await new Promise((resolve, reject) => {
        axios.post('/terraform/v1/ai-talk/stage/query', {
          room: roomUUID, roomToken, sid: stageUUID, rid: requestUUID,
        }, {
          headers: Token.loadBearerHeader(),
        }).then(res => {
          if (res.data?.data?.finished) {
          }
          resolve(res.data.data);
        }).catch(error => reject(error));
      });

      if (resp.finished) {
        console.log(`TTS: Conversation finished.`);
        break;
      }

      await new Promise((resolve) => setTimeout(resolve, 1000));
    }
  }, [roomUUID, roomToken, stageUUID, userID, setUserAsrText]);

  // When user start dictation, open the microphone, start the media recorder, and send the audio
  // every some seconds.
  const startDictation = React.useCallback(async (onWorking) => {
    if (!robotReady) return;
    if (ref.current.mediaStream) return;

    console.log(`Start dictation, open microphone.`)
    await new Promise(resolve => {
      navigator.mediaDevices.getUserMedia(
        {audio: true}
      ).then((stream) => {
        ref.current.mediaStream = stream;
        console.log(`Start dictation, microphone opened.`);
        resolve();
      }).catch(error => errorLog(`${t('lr.room.mic')}: ${error}`));
    });

    // Now dictation is working.
    ref.current.shouldStop = false;
    setProcessing(true);
    onWorking && onWorking();

    const startRecorder = async () => {
      console.log("=============");

      // Media artifact, a piece of audio segment generated by media recorder.
      const artifact = {
        // The time elappsed.
        recordStarttime: new Date(),
        recordStoptime: null,
        duration: () => {
          if (!artifact.recordStoptime) return null;
          return artifact.recordStoptime - artifact.recordStarttime;
        },
        // For dictation, we never stop the media stream util user press the stop button.
        // TODO: FIXME: Set the codec, see https://developer.mozilla.org/en-US/docs/Web/API/MediaRecorder/mimeType
        // The stream is already opened when robot ready, or all answers are played.
        // See https://www.sitelint.com/lab/media-recorder-supported-mime-type/
        mediaRecorder: new MediaRecorder(ref.current.mediaStream),
        audioChunks: [],
      };

      // See https://developer.mozilla.org/en-US/docs/Web/API/MediaRecorder#events
      artifact.mediaRecorder.addEventListener("start", () => {
        console.log(`Event: Recording start to record`);
        setMicWorking(true);
      });

      artifact.mediaRecorder.addEventListener("dataavailable", ({ data }) => {
        artifact.audioChunks.push(data);
        const userMayInput = new Date() - artifact.recordStarttime - timeoutWaitForLastVoice;
        const ts = new Date().toISOString().split('T')[1].split('Z')[0];
        console.log(`${ts} Event: Device dataavailable event ${data.size} bytes, duration=${userMayInput}ms`);
      });

      // For dictation, we use shorter timeslice, to trigger dataavailable event ASAP.
      artifact.mediaRecorder.start(1500);
      console.log(`Event: Recording started`);
      return artifact;
    };

    const stopRecorder = async (artifact, dispose) => {
      // Now, stop the dictation.
      const closeMicrophone = () => {
        if (!dispose) return;
        console.log(`Stop dictation, close microphone.`);
        const stream = ref.current.mediaStream;
        ref.current.mediaStream = null;
        stream.getTracks().forEach(track => track.stop());
        setMicWorking(false);
      };

      if (!artifact) {
        closeMicrophone();
        return;
      }

      console.log(`Stop dictation, stop recorder.`)
      await new Promise(resolve => {
        artifact.mediaRecorder.addEventListener("stop", () => {
          if (dispose) closeMicrophone();

          // To help estimate the duration of artifact.
          artifact.recordStoptime = new Date();

          const ts = new Date().toISOString().split('T')[1].split('Z')[0];
          console.log(`${ts} Event: Recorder stopped, chunks=${artifact.audioChunks.length}, duration=${artifact.duration()}ms`);
          
          resolve();
        });

        console.log(`Event: Recorder stop, chunks=${artifact.audioChunks.length}, state=${artifact.mediaRecorder.state}`);
        artifact.mediaRecorder.stop();
      });

      return artifact;
    };

    // Handle all events during dictation.
    try {
      let artifact = await startRecorder();
      while (true) {
        if (ref.current.shouldStop) {
          console.log(`Stop dictation, user reset it.`);
          await stopRecorder(artifact, true);
          return;
        }

        if (new Date() - artifact.recordStarttime >= maxSegmentTime) {
          const readyArtifact = await stopRecorder(artifact, false);
          artifact = await startRecorder();
          const ts = new Date().toISOString().split('T')[1].split('Z')[0];
          console.log(`${ts} Restart dictation ok, start to send artifact.`);

          // Start a async task to process the audio segment.
          processUserInput(readyArtifact).catch((e) => {
            console.warn(`Dictation ignore any segment ${JSON.stringify(readyArtifact)} error ${e}`);
          });
        }

        await new Promise(resolve => setTimeout(resolve, 1000));
      }
    } finally {
      setProcessing(false);
      setMicWorking(false);
    }
  }, [ref, robotReady, t, errorLog, maxSegmentTime, processUserInput]);

  // When user stop dictation, close the microphone, stop the media recorder, and send the last
  // audio data.
  const stopDictation = React.useCallback(async () => {
    if (!robotReady) return;
    if (!ref.current.mediaStream) return;

    // Notify to stop.
    ref.current.shouldStop = true;

    // Wait util stopped.
    while (ref.current.mediaStream) {
      await new Promise(resolve => setTimeout(resolve, 300));
    }
  }, [ref, robotReady]);

  // User directly send text message to assistant.
  const sendText = React.useCallback(async (text, onFinished) => {
    if (!robotReady) return;
    if (!text) return;

    try {
      const requestUUID = await new Promise((resolve, reject) => {
        axios.post('/terraform/v1/ai-talk/stage/conversation', {
          room: roomUUID, roomToken, sid: stageUUID,
        }, {
          headers: Token.loadBearerHeader(),
        }).then(res => {
          console.log(`ASR: Start conversation success, rid=${res.data.data.rid}`);
          resolve(res.data.data.rid);
        }).catch(handleError);
      });

      await new Promise((resolve, reject) => {
        axios.post('/terraform/v1/ai-talk/stage/upload', {
          room: roomUUID, roomToken, sid: stageUUID, rid: requestUUID, userId: userID,
          text: text,
        }, {
          headers: Token.loadBearerHeader(),
        }).then(res => {
          console.log(`ASR: Send text success: ${res.data.data.rid} ${res.data.data.asr}`);
          resolve(res.data.data.rid);
        }).catch(handleError);
      });

      // For text message, don't wait for conversation to finish, because user may send text
      // message continuously.
    } finally {
      onFinished && onFinished();
    }
  }, [robotReady, handleError, roomUUID, roomToken, stageUUID, userID]);

  // Request server to create a new popout for subscribing all events.
  React.useEffect(() => {
    if (!stageUUID) return; // For assistant, we must want stage to be ready.
    if (!roomUUID) return;

    console.log(`Start: Create a new stage`);
    axios.post('/terraform/v1/ai-talk/subscribe/start', {
      room: roomUUID, roomToken,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      console.log(`Start: Create popout success: ${JSON.stringify(res.data.data)}`);
      setStagePopoutUUID(res.data.data.spid);
    }).catch(handleError);
  }, [handleError, roomUUID, roomToken, setStagePopoutUUID, stageUUID]);

  // Try to request messages of stage util end.
  React.useEffect(() => {
    if (!robotReady) return;

    const requestMessages = async () => {
      if (!robotReady || !stageUUID || !stagePopoutUUID) return;
      if (refRequest.current.requesting) return;
      refRequest.current.requesting = true;

      try {
        const {msgs, pending} = await new Promise((resolve, reject) => {
          axios.post('/terraform/v1/ai-talk/subscribe/query', {
            room: roomUUID, roomToken, sid: stageUUID, spid: stagePopoutUUID, userId: userID,
          }, {
            headers: Token.loadBearerHeader(),
          }).then(res => {
            // Don't show detail logs for pulling.
            //const ts = new Date().toISOString().split('T')[1].split('Z')[0];
            //console.log(`Start: Query popout success at ${ts}: ${JSON.stringify(res.data.data)}`);
            resolve(res.data.data);
          }).catch(handleError);
        });

        refRequest.current.hasPendingMessages = pending;
        refRequest.current.gotMessages = msgs?.length || 0;
        if (!msgs?.length) return;

        for (let i = 0; i < msgs.length; i++) {
          const msg = msgs[i];
          if (msg.role === 'user') {
            traceLog(msg.username || 'You', msg.msg, 'primary', true);
            continue;
          }

          const audioSegmentUUID = msg.asid;
          traceLog(msg.username, msg.msg, 'success');

          // For dictation pattern, we always ignore TTS audio files.
          // No audio file, skip it.
          console.log(`TTS: Dictation always consume as text, ${JSON.stringify(msg)}`);

          // Remove the AI generated audio.
          await new Promise((resolve, reject) => {
            axios.post('/terraform/v1/ai-talk/subscribe/remove', {
              room: roomUUID, roomToken, sid: stageUUID, spid: stagePopoutUUID, asid: audioSegmentUUID,
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
    };

    const timer = setInterval(async () => {
      requestMessages().catch(handleError);
    }, 1000);
    return () => clearInterval(timer);
  }, [robotReady, handleError, stageUUID, stagePopoutUUID, traceLog, refRequest, roomUUID, roomToken, userID]);

  if (booting) {
    return <><Spinner animation="border" variant="primary" size='sm'></Spinner> Booting...</>;
  }
  return (
    <div>
      <div><audio ref={playerRef} controls={true} hidden='hidden' /></div>
      {stageUUID && !robotReady ? <>
        <AITalkUserConfig {...{roomUUID, roomToken,
          username, userLanguage, stageUUID, userID, disabled: requesting, label: t('lr.room.startd'),
          onSubmit: startChatting}} />
      </> : ''}
      {robotReady && !isMobile ?
        <Row>
          <Col>
            <AITalkDictationImpl {...{
              processing, micWorking, startDictation, stopDictation, sendText, roomUUID, roomToken,
              stageUUID, stageUser, aiAsrEnabled, userAsrText
            }} />
          </Col>
          <Col>
            <AITalkTraceLogPC {...{traceLogs, traceCount, roomUUID, roomToken}}>
              <AITalkErrorLog {...{errorLogs, removeErrorLog}} />
              <AITalkTipLog {...{tipLogs, removeTipLog}} />
            </AITalkTraceLogPC>
          </Col>
        </Row> : ''}
      {robotReady && isMobile ?
        <div>
          <AITalkTraceLogMobile {...{traceLogs, traceCount}} />
          <AITalkErrorLog {...{errorLogs, removeErrorLog}} />
          <AITalkTipLog {...{tipLogs, removeTipLog}} />
          <AITalkDictationImpl {...{
            processing, micWorking, startDictation, stopDictation, sendText, roomUUID, roomToken,
            stageUUID, stageUser, aiAsrEnabled, userAsrText
          }} />
        </div> : ''}
      <div ref={endPanelRef}></div>
    </div>
  );
}

function AITalkUserConfig({roomUUID, roomToken, username, userLanguage, stageUUID, userID, disabled, label, onSubmit, onCancel}) {
  const handleError = useErrorHandler();
  const [loading, setLoading] = React.useState(true);
  const [user, setUser] = React.useState(null);

  React.useEffect(() => {
    if (!userID) return;

    axios.post('/terraform/v1/ai-talk/user/query', {
      room: roomUUID, roomToken, sid: stageUUID, userId: userID,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      setLoading(false);
      const u = res.data.data;
      setUser({
        ...u, userId: userID,
        username: username || u.username,
        language: userLanguage || u.language,
      });
      console.log(`Start: Query stage user success, ${JSON.stringify(res.data.data)}`);
    }).catch(handleError);
  }, [handleError, setUser, setLoading, roomUUID, roomToken, username, userLanguage, stageUUID, userID]);

  if (loading || !user) {
    return <><Spinner animation="border" variant="primary" size='sm'></Spinner> Loading...</>;
  }
  return <AITalkUserConfigImpl {...{roomUUID, roomToken, stageUUID, user, disabled, label, onSubmit, onCancel}} />;
}

function AITalkUserConfigImpl({roomUUID, roomToken, stageUUID, user, disabled, label, onSubmit, onCancel}) {
  const {t} = useTranslation();
  const handleError = useErrorHandler();

  const [requesting, setRequesting] = React.useState(false);
  const [userName, setUserName] = React.useState(user.username);
  const [userLanguage, setUserLanguage] = React.useState(user.language);

  const updateConfig = React.useCallback((e) => {
    setRequesting(true);

    axios.post('/terraform/v1/ai-talk/user/update', {
      room: roomUUID, roomToken, sid: stageUUID, userId: user.userId,
      name: userName, lang: userLanguage
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      setRequesting(false);
      console.log(`Start: Update stage user success`);
      onSubmit && onSubmit(res.data.data);
    }).catch(handleError);
  }, [handleError, roomUUID, roomToken, stageUUID, user, userName, userLanguage, onSubmit, setRequesting]);

  return <>
    <Form>
      <Form.Group className="mb-3">
        <Form.Label>{t('lr.room.uname')}</Form.Label>
        <Form.Text> * {t('lr.room.uname2')}</Form.Text>
        <Form.Control as="input" type='input' defaultValue={userName} onChange={(e) => {
          e.preventDefault();
          setUserName(e.target.value);
        }} />
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>{t('transcript.lang')}</Form.Label>
        <Form.Text> * {t('transcript.lang3')}. &nbsp;
          {t('helper.eg')} <code>en, zh, fr, de, ja, ru </code>, ... &nbsp;
          {t('helper.see')} <a href='https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes' target='_blank' rel='noreferrer'>ISO-639-1</a>.
        </Form.Text>
        <Form.Control as="input" defaultValue={userLanguage} onChange={(e) => {
          e.preventDefault();
          setUserLanguage(e.target.value);
        }} />
      </Form.Group>
      <Button variant="primary" type="button" disabled={requesting || disabled} onClick={updateConfig}>
        {label || t('helper.submit')}
      </Button>
      {onCancel ? <>
        &nbsp;
        <Button variant="primary" type="button" disabled={requesting || disabled} onClick={onCancel}>
          {t('helper.cancel')}
        </Button>
      </> : ''}
      <p></p>
    </Form>
  </>;
}

function AITalkDictationImpl({processing, micWorking, userAsrText, sendText, startDictation, stopDictation, roomUUID, roomToken, stageUUID, stageUser, aiAsrEnabled}) {
  const {t} = useTranslation();
  const [showSettings, setShowSettings] = React.useState(false);
  const [popoutUrl, setPopoutUrl] = React.useState(null);
  const [showUserConfig, setShowUserConfig] = React.useState(false);
  const [user, setUser] = React.useState(stageUser);
  const [description, setDescription] = React.useState();
  const [userText, setUserText] = React.useState('');
  const [takingDictation, setTakingDictation] = React.useState(processing);
  const [stillWorking, setStillWorking] = React.useState(false);

  React.useEffect(() => {
    if (!roomUUID) return;

    const params = [
      'app=ai-talk',
      'popout=1',
      'assistant=1',
      `created=${new Date().toISOString()}`,
      `random=${Math.random().toString(16).slice(-8)}`,
      `room=${roomUUID}`,
      `roomToken=${roomToken}`,
    ];
    const url = `${window.PUBLIC_URL}/${Locale.current()}/routers-popout?${params.join('&')}`;
    setPopoutUrl(url);
    console.log(`Generated popout URL: ${url}`);
  }, [roomUUID, setPopoutUrl, roomToken, user]);

  const openPopout = React.useCallback((e) => {
    e.preventDefault();
    if (!popoutUrl) return;
    setShowSettings(false);
    window.open(popoutUrl, '_blank', 'noopener,noreferrer,width=1024,height=768');
  }, [popoutUrl, setShowSettings]);

  const openSettings = React.useCallback((e) => {
    e.preventDefault();
    setShowUserConfig(true);
    setShowSettings(false);
  }, [setShowUserConfig, setShowSettings]);

  const onFinishUserConfig = React.useCallback((user) => {
    setShowUserConfig(false);
    setUser({...user});
  }, [setShowUserConfig, setUser]);

  React.useEffect(() => {
    if (!user) return;
    setDescription(` for ${user.username || 'You'} (${user.language})`);
  }, [user, setDescription]);

  const onSendText = React.useCallback((e) => {
    e.preventDefault();
    if (!userText) return;

    sendText && sendText(userText, () => {
      setUserText('');
    });
  }, [userText, sendText, setUserText]);

  const onUserPressKey = React.useCallback((e) => {
    if (e.key === 'Enter') {
      onSendText(e);
    }
  }, [onSendText]);

  React.useEffect(() => {
    if (!userAsrText) return;
    setUserText(userAsrText);
  }, [userAsrText]);

  const onUserStartDictation = React.useCallback(async (e) => {
    e.preventDefault();
    try {
      setStillWorking(true);
      // Note that we should never wait for all dictation finished, but only wait for
      // callback onWorking.
      await new Promise((resolve, reject) => {
        startDictation(resolve).catch(reject);
      });
      await new Promise((resolve) => setTimeout(resolve, 1200));
      setTakingDictation(true);
    } finally {
      setStillWorking(false);
    }
  }, [setStillWorking, startDictation, setTakingDictation]);

  const onUserStopDictation = React.useCallback(async (e) => {
    e.preventDefault();
    try {
      setStillWorking(true);
      await stopDictation();
      await new Promise((resolve) => setTimeout(resolve, 800));
      setTakingDictation(false);
    } finally {
      setStillWorking(false);
    }
  }, [setStillWorking, stopDictation, setTakingDictation]);

  return (
    <div>
      <Card>
        <Card.Header>
          {t('lr.room.ait')}{showSettings ? '' : description}
          <div role='button' className='ai-talk-settings-btn'>
            <Icon.Gear size={20} onClick={(e) => setShowSettings(!showSettings)} />
          </div>
          {showSettings && <div className='ai-talk-settings-menu2'>
            <Dropdown.Menu show={true}>
              {false && <Dropdown.Item href="#!" onClick={openPopout}>{t('lr.room.popchat2')}</Dropdown.Item>}
              <Dropdown.Item href="#!" onClick={openSettings}>{t('lr.room.settings')}</Dropdown.Item>
            </Dropdown.Menu>
          </div>}
        </Card.Header>
        <Card.Body>
          {showUserConfig ? <AITalkUserConfig {...{
            roomUUID, roomToken, stageUUID, userID: stageUser.userId, onSubmit: onFinishUserConfig,
            onCancel: () => setShowUserConfig(false),
          }} /> : ''}
          <InputGroup className="mb-3">
            <Form.Control
              as="input" placeholder={t('lr.room.text')} aria-describedby="basic-addon2" value={userText}
              onChange={(e) => setUserText(e.target.value)}
              onKeyPress={onUserPressKey}/>
            <Button variant="primary" id="button-addon2" onClick={onSendText}>{t('helper.send')}</Button>
          </InputGroup>
          {aiAsrEnabled && <>
            {!takingDictation ? <>
                <Button variant="primary" type="button" disabled={stillWorking} onClick={onUserStartDictation}>
                  {t('helper.start')}
                </Button> &nbsp;
                {(processing || micWorking) && <Spinner animation="border" variant="primary" size='sm'></Spinner>}
              </> :
              <>
                <Button variant="primary" type="button" disabled={stillWorking} onClick={onUserStopDictation}>
                  {t('helper.stop')}
                </Button> &nbsp;
                <Spinner animation="border" variant="primary" size='sm'></Spinner>
              </>}
          </>}
        </Card.Body>
      </Card>
    </div>
  );
}

function AITalkTraceLogPC({traceLogs, traceCount, children, roomUUID, roomToken}) {
  const {t} = useTranslation();
  const [showSettings, setShowSettings] = React.useState(false);
  const [popoutUrl, setPopoutUrl] = React.useState(null);

  // Scroll the log panel.
  const logPanelRef = React.useRef(null);
  React.useEffect(() => {
    if (!logPanelRef?.current) return;
    console.log(`Logs scroll to end, height=${logPanelRef.current.scrollHeight}, logs=${traceLogs.length}, count=${traceCount}`);
    logPanelRef.current.scrollTo(0, logPanelRef.current.scrollHeight);
  }, [traceLogs, logPanelRef, traceCount]);

  React.useEffect(() => {
    if (!roomUUID) return;

    const params = [
      'app=ai-talk',
      'popout=1',
      'assistant=0', // Without assistant.
      `created=${new Date().toISOString()}`,
      `random=${Math.random().toString(16).slice(-8)}`,
      `room=${roomUUID}`,
      `roomToken=${roomToken}`,
    ];
    const url = `${window.PUBLIC_URL}/${Locale.current()}/routers-popout?${params.join('&')}`;
    setPopoutUrl(url);
    console.log(`Generated popout URL: ${url}`);
  }, [roomUUID, setPopoutUrl, roomToken]);

  const openPopout = React.useCallback((e) => {
    e.preventDefault();
    if (!popoutUrl) return;
    setShowSettings(false);
    window.open(popoutUrl, '_blank', 'noopener,noreferrer,width=1024,height=768');
  }, [popoutUrl, setShowSettings]);

  return (
    <div>
      <Card>
        <Card.Header>
          {t('lr.room.ait2')}
          <div role='button' className='ai-talk-settings-btn'>
            <Icon.Gear size={20} onClick={(e) => setShowSettings(!showSettings)} />
          </div>
          {showSettings && <div className='ai-talk-settings-menu'>
            <Dropdown.Menu show={true}>
              <Dropdown.Item href="#!" onClick={openPopout}>{t('lr.room.popchat')}</Dropdown.Item>
            </Dropdown.Menu>
          </div>}
        </Card.Header>
        <Card.Body>
          <div className='ai-talk-trace-logs-pcfs' ref={logPanelRef}>
            {children}
            {traceLogs.map((log) => {
              return (
                <Alert key={log.id} variant={log.variant}>
                  {log.role}: {log.msg}
                </Alert>
              );
            })}
          </div>
        </Card.Body>
      </Card>
    </div>
  );
}

function AITalkTraceLogMobile({traceLogs, traceCount}) {
  // Scroll the log panel.
  const logPanelRef = React.useRef(null);
  React.useEffect(() => {
    if (!logPanelRef?.current) return;
    console.log(`Logs scroll to end, height=${logPanelRef.current.scrollHeight}, logs=${traceLogs.length}, count=${traceCount}`);
    logPanelRef.current.scrollTo(0, logPanelRef.current.scrollHeight);
  }, [traceLogs, logPanelRef, traceCount]);

  return (
    <div className='ai-talk-trace-logs-mobilefs-dictation' ref={logPanelRef}>
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

function AITalkErrorLog({errorLogs, removeErrorLog}) {
  return (
    <React.Fragment>
      {errorLogs.map((log) => {
        return (
          <Alert key={log.id} onClose={() => removeErrorLog(log)} variant='danger' dismissible>
            <Alert.Heading>Error!</Alert.Heading>
            <p>{log.msg}</p>
          </Alert>
        );
      })}
    </React.Fragment>
  );
}

function AITalkTipLog({tipLogs, removeTipLog}) {
  return (
    <React.Fragment>
      {tipLogs.map((log) => {
        return (
          <Alert key={log.id} onClose={() => removeTipLog(log)} variant='success' dismissible>
            <Alert.Heading>{log.title}</Alert.Heading>
            <p>{log.msg}</p>
          </Alert>
        );
      })}
    </React.Fragment>
  );
}
