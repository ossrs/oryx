import React from "react";
import {Alert, Button, Card, Col, Dropdown, Form, InputGroup, Row, Spinner} from "react-bootstrap";
import {useTranslation} from "react-i18next";
import {useErrorHandler} from "react-error-boundary";
import useIsMobile from "./IsMobile";
import axios from "axios";
import {Locale, Token} from "../utils";
import * as Icon from "react-bootstrap-icons";
import Container from "react-bootstrap/Container";

export function AITalkAssistantPanel({roomUUID, roomToken, fullscreen}) {
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const isMobile = useIsMobile();

  // The timeout in milliseconds.
  const timeoutForMicrophoneTestToRun = 50;
  const timeoutWaitForMicrophoneToClose = 900;
  const timeoutWaitForLastVoice = 700;
  const durationRequiredUserInput = 600;

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

  // When robot is ready, open the microphone ASAP to accept user input.
  React.useEffect(() => {
    const fnImpl = async () => {
      if (!robotReady) return;
      if (ref.current.mediaStream) return;

      // Wait util no audio is playing.
      do {
        await new Promise((resolve) => setTimeout(resolve, 300));
      } while(refRequest.current.playingAudio);

      console.log(`Robot is ready, open microphone.`)
      navigator.mediaDevices.getUserMedia(
        { audio: true }
      ).then((stream) => {
        ref.current.mediaStream = stream;
        console.log(`Robot is ready, microphone opened.`);
      }).catch(error => errorLog(`${t('lr.room.mic')}: ${error}`));
    };
    fnImpl();
  }, [errorLog, t, robotReady, ref, refRequest]);

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

  // User start a conversation, by recording input.
  const startRecording = React.useCallback(async () => {
    if (!robotReady) return;
    if (!ref.current.mediaStream) return;
    if (ref.current.stopHandler) clearTimeout(ref.current.stopHandler);
    if (ref.current.mediaRecorder) return;
    if (ref.current.isRecording) return;
    if (refRequest.current.playingAudio) return;
    ref.current.recordStarttime = new Date();
    ref.current.isRecording = true;
    ref.current.count += 1;

    console.log("=============");

    // The stream is already opened when robot ready, or all answers are played.
    // See https://www.sitelint.com/lab/media-recorder-supported-mime-type/
    ref.current.mediaRecorder = new MediaRecorder(ref.current.mediaStream);
    ref.current.mediaStream = null;

    // See https://developer.mozilla.org/en-US/docs/Web/API/MediaRecorder#events
    ref.current.mediaRecorder.addEventListener("start", () => {
      console.log(`Event: Recording start to record`);
      setMicWorking(true);
    });

    ref.current.mediaRecorder.addEventListener("dataavailable", ({ data }) => {
      ref.current.audioChunks.push(data);
      console.log(`Event: Device dataavailable event ${data.size} bytes`);
    });

    ref.current.mediaRecorder.start();
    console.log(`Event: Recording started`);
  }, [robotReady, ref, setMicWorking, refRequest]);

  // User click stop button, we delay some time to allow cancel the stopping event.
  const stopRecording = React.useCallback(async () => {
    if (!robotReady) return;

    const processUserInput = async(userMayInput) => {
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

        const audioBlob = new Blob(ref.current.audioChunks);
        reader.readAsDataURL(audioBlob);
      });

      // Upload the user input audio to the server.
      await new Promise((resolve, reject) => {
        console.log(`ASR: Uploading ${ref.current.audioChunks.length} chunks, stage=${stageUUID}, user=${userID}`);
        ref.current.audioChunks = [];

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
    };

    const stopRecordingImpl = async () => {
      if (!ref.current.mediaRecorder) return;

      try {
        const userMayInput = new Date() - ref.current.recordStarttime - timeoutWaitForLastVoice;
        console.log(`Event: User stop record, duration=${userMayInput}ms, state=${ref.current.mediaRecorder.state}`);

        await new Promise(resolve => {
          ref.current.mediaRecorder.addEventListener("stop", () => {
            const stream = ref.current.mediaRecorder.stream;
            stream.getTracks().forEach(track => track.stop());
            setTimeout(resolve, 30);
          });

          console.log(`Event: Recorder stop, chunks=${ref.current.audioChunks.length}, state=${ref.current.mediaRecorder.state}`);
          ref.current.mediaRecorder.stop();
        });

        setMicWorking(false);
        setProcessing(true);
        console.log(`Event: Recoder stopped, chunks=${ref.current.audioChunks.length}`);

        if (userMayInput < durationRequiredUserInput) {
          console.warn(`System: You didn't say anything!`);
          alert(`Warning: You didn't say anything!`);
        } else {
          try {
            await processUserInput(userMayInput);
          } catch (e) {
            console.warn(`System: Server error ${e}`);
            console.warn(`System: Please try again.`);
            alert(`System: Server error ${e}`);
          }
        }

        // Wait util no audio is playing.
        do {
          await new Promise((resolve) => setTimeout(resolve, 600));
        } while(refRequest.current.playingAudio);

        // Reopen the microphone if no audio is playing.
        console.log(`Conversation is ended, open microphone.`)
        new Promise((resolve, reject) => {
          navigator.mediaDevices.getUserMedia(
            { audio: true }
          ).then((stream) => {
            ref.current.mediaStream = stream;
            console.log(`All audios is played, microphone opened.`);
            resolve();
          }).catch(error => reject(error));
        });
      } catch (e) {
        alert(e);
      } finally {
        setProcessing(false);
        ref.current.mediaRecorder = null;
        ref.current.isRecording = false;
      }
    };

    if (ref.current.stopHandler) clearTimeout(ref.current.stopHandler);
    ref.current.stopHandler = setTimeout(() => {
      stopRecordingImpl();
    }, timeoutWaitForLastVoice);
  }, [roomUUID, roomToken, stageUUID, userID, robotReady, ref, setProcessing, setMicWorking, refRequest, setUserAsrText]);

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
    } finally {
      onFinished && onFinished();
    }
  }, [robotReady, handleError, roomUUID, roomToken, stageUUID, userID]);

  // Setup the keyboard event, for PC browser.
  React.useEffect(() => {
    if (!robotReady) return;
    if (!aiAsrEnabled) return;

    const handleKeyDown = (e) => {
      if (processing) return;
      // Ignore the input event.
      const tagName = e.target.tagName.toLowerCase();
      if (tagName === 'input' || tagName === 'textarea') return;
      if (e.key !== 'r' && e.key !== '\\' && e.key !== ' ') return;
      startRecording();
    };
    const handleKeyUp = (e) => {
      if (processing) return;
      if (e.key !== 'r' && e.key !== '\\' && e.key !== ' ') return;
      stopRecording();
    };

    window.addEventListener('keydown', handleKeyDown);
    window.addEventListener('keyup', handleKeyUp);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
      window.removeEventListener('keyup', handleKeyUp);
    };
  }, [robotReady, startRecording, stopRecording, processing, aiAsrEnabled]);

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
            traceLog(msg.username || 'You', msg.msg, 'primary');
            continue;
          }

          const audioSegmentUUID = msg.asid;
          traceLog(msg.username, msg.msg, 'success');

          // No audio file, skip it.
          if (!msg.hasAudio) {
            console.log(`TTS: Consume text message done, ${JSON.stringify(msg)}`);
          } else {
            // Play the AI generated audio.
            await new Promise(resolve => {
              const url = `/terraform/v1/ai-talk/subscribe/tts?sid=${stageUUID}&spid=${stagePopoutUUID}&asid=${audioSegmentUUID}&room=${roomUUID}&roomToken=${roomToken}`;
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
          }

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

  // When we got any messages from server, set to playing mode and last for a while.
  React.useEffect(() => {
    const timer = setInterval(async () => {
      if (refRequest.current.gotMessages || refRequest.current.hasPendingMessages) {
        if (!refRequest.current.playingAudio) {
          console.log(`Subscribe: Got messages, start to play, mic=${micWorking}`);
          if (!micWorking) {
            setProcessing(true);
          }
        }
        refRequest.current.playingAudio = 5;
        return;
      }

      if (refRequest.current.playingAudio) {
        refRequest.current.playingAudio -= 1;
        if (!refRequest.current.playingAudio) {
          console.log(`Subscribe: No messages, stop playing, mic=${micWorking}`);
          if (!micWorking) {
            setProcessing(false);
          }
        }
      }
    }, 300);
    return () => clearInterval(timer);
  }, [refRequest, micWorking, setProcessing]);

  // A long time interval logger.
  React.useEffect(() => {
    const timer = setInterval(async () => {
      const ts = new Date().toISOString().split('T')[1].split('Z')[0];
      const r0 = refRequest.current;
      console.log(`Timer: ${ts}, messages got=${r0.gotMessages}, pending=${r0.hasPendingMessages}, playing=${r0.playingAudio}, requesting=${r0.requesting}`);
    }, 10000);
    return () => clearInterval(timer);
  }, [refRequest]);

  if (booting) {
    return <><Spinner animation="border" variant="primary" size='sm'></Spinner> Booting...</>;
  }
  return (
    <div>
      <div><audio ref={playerRef} controls={true} hidden='hidden' /></div>
      {stageUUID && !robotReady ? <>
        <AITalkUserConfig {...{roomUUID, roomToken, stageUUID, userID, disabled: requesting, label: t('lr.room.talk'), onSubmit: startChatting}} />
      </> : ''}
      {robotReady && !isMobile ?
        <Row>
          <Col>
            <AITalkAssistantImpl {...{
              processing, micWorking, startRecording, stopRecording, sendText, roomUUID, roomToken,
              stageUUID, stageUser, aiAsrEnabled, userAsrText
            }} />
          </Col>
          <Col>
            <AITalkTraceLogPC {...{traceLogs, traceCount, roomUUID, roomToken, fullscreen}}>
              <AITalkErrorLog {...{errorLogs, removeErrorLog}} />
              <AITalkTipLog {...{tipLogs, removeTipLog}} />
            </AITalkTraceLogPC>
          </Col>
        </Row> : ''}
      {robotReady && isMobile ?
        <div>
          <AITalkTraceLogMobile {...{traceLogs, traceCount, fullscreen}} />
          <AITalkErrorLog {...{errorLogs, removeErrorLog}} />
          <AITalkTipLog {...{tipLogs, removeTipLog}} />
          <AITalkAssistantImpl {...{
            processing, micWorking, startRecording, stopRecording, sendText, roomUUID, roomToken,
            stageUUID, stageUser, aiAsrEnabled, userAsrText
          }} />
        </div> : ''}
      <div ref={endPanelRef}></div>
    </div>
  );
}

export function AITalkChatPanel({roomUUID, roomToken}) {
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const isMobile = false; // For popout, always PC, not mobile.

  // The player ref, to access the audio player.
  const playerRef = React.useRef(null);
  const [requesting, setRequesting] = React.useState(false);
  const [robotReady, setRobotReady] = React.useState(false);

  // The uuid and robot in stage, which is unchanged after stage started.
  const [stageUUID, setStageUUID] = React.useState(null);
  const [stageVoice, setStageVoice] = React.useState(null);
  const [stagePopoutUUID, setStagePopoutUUID] = React.useState(null);

  // Possible value is 1: yes, -1: no, 0: undefined.
  const [obsAutostart, setObsAutostart] = React.useState(0);
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
    axios.post('/terraform/v1/ai-talk/subscribe/start', {
      room: roomUUID, roomToken,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      console.log(`Start: Create popout success: ${JSON.stringify(res.data.data)}`);
      setStageUUID(res.data.data.sid);
      setStagePopoutUUID(res.data.data.spid);
      setStageVoice(res.data.data.voice);
    }).catch(handleError);
  }, [handleError, roomUUID, roomToken, setStagePopoutUUID, setStageVoice, setStageUUID]);

  // Try to start the robot automatically, for OBS.
  React.useEffect(() => {
    if (!stageUUID || !stageVoice) return;

    const listener = () => {
      playerRef.current.removeEventListener('ended', listener);

      setObsAutostart(1);
      setRobotReady(true);
      console.log(`Stage started, AI is ready, sid=${stageUUID}`);
    };
    playerRef.current.addEventListener('ended', listener);

    playerRef.current.src = `/terraform/v1/ai-talk/stage/hello-voices/${stageVoice}?sid=${stageUUID}`;
    playerRef.current.play().catch((error) => {
      setObsAutostart(-1);
    });
  }, [t, errorLog, stageUUID, stageVoice, setObsAutostart]);

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

    playerRef.current.src = `/terraform/v1/ai-talk/stage/hello-voices/${stageVoice}?sid=${stageUUID}`;
    playerRef.current.play().catch(error => {
      errorLog(`${t('lr.room.speaker')}: ${error}`);
      setRequesting(false);
    });
  }, [t, errorLog, stageUUID, stageVoice, setRobotReady, setRequesting]);

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
  // Try to request messages of stage util end.
  React.useEffect(() => {
    if (!robotReady) return;

    const requestMessages = async () => {
      if (!robotReady || !stageUUID || !stagePopoutUUID) return;
      if (refRequest.current.requesting) return;
      refRequest.current.requesting = true;

      try {
        const msgs = await new Promise((resolve, reject) => {
          axios.post('/terraform/v1/ai-talk/subscribe/query', {
            room: roomUUID, roomToken, sid: stageUUID, spid: stagePopoutUUID,
          }, {
            headers: Token.loadBearerHeader(),
          }).then(res => {
            // Don't show detail logs for pulling.
            //const ts = new Date().toISOString().split('T')[1].split('Z')[0];
            //console.log(`Start: Query popout success at ${ts}: ${JSON.stringify(res.data.data)}`);
            resolve(res.data.data.msgs);
          }).catch(handleError);
        });

        if (!msgs?.length) return;
        for (let i = 0; i < msgs.length; i++) {
          const msg = msgs[i];
          if (msg.role === 'user') {
            traceLog(msg.username, msg.msg, 'primary');
            return;
          }

          const audioSegmentUUID = msg.asid;
          traceLog(msg.username, msg.msg, 'success');

          // Play the AI generated audio.
          await new Promise(resolve => {
            const url = `/terraform/v1/ai-talk/subscribe/tts?sid=${stageUUID}&spid=${stagePopoutUUID}&asid=${audioSegmentUUID}&room=${roomUUID}&roomToken=${roomToken}`;
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
  }, [robotReady, handleError, stageUUID, stagePopoutUUID, traceLog, refRequest, roomUUID, roomToken]);

  return (
    <Container fluid>
      <p></p>
      <div>
        {obsAutostart === -1 ?
          <Button disabled={requesting} variant="primary" type="submit" onClick={startChatting}>
            {t('lr.room.talk')}
          </Button> : ''}
        <div><audio ref={playerRef} controls={true} hidden='hidden' /></div>
        <AITalkErrorLog {...{errorLogs, removeErrorLog}} />
        <AITalkTipLog {...{tipLogs, removeTipLog}} />
        <AITalkTraceLogPopout {...{traceLogs, traceCount}} />
        <div ref={endPanelRef}></div>
      </div>
    </Container>
  );
}

function AITalkUserConfig({roomUUID, roomToken, stageUUID, userID, disabled, label, onSubmit, onCancel}) {
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
      setUser({...res.data.data, userId: userID});
      console.log(`Start: Query stage user success, ${JSON.stringify(res.data.data)}`);
    }).catch(handleError);
  }, [handleError, setUser, setLoading, roomUUID, roomToken, stageUUID, userID]);

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

function AITalkAssistantImpl({processing, micWorking, startRecording, stopRecording, userAsrText, sendText, roomUUID, roomToken, stageUUID, stageUser, aiAsrEnabled}) {
  const {t} = useTranslation();
  const isMobile = useIsMobile();
  const [showSettings, setShowSettings] = React.useState(false);
  const [popoutUrl, setPopoutUrl] = React.useState(null);
  const [showUserConfig, setShowUserConfig] = React.useState(false);
  const [user, setUser] = React.useState(stageUser);
  const [description, setDescription] = React.useState();
  const [userText, setUserText] = React.useState('');

  React.useEffect(() => {
    if (!roomUUID) return;

    const r0 = Math.random().toString(16).slice(-8);
    const created = new Date().toISOString();
    const url = `${window.PUBLIC_URL}/${Locale.current()}/routers-popout?app=ai-talk&popout=1&assistant=1&room=${roomUUID}&created=${created}&random=${r0}&roomToken=${roomToken}`;
    setPopoutUrl(url);
    console.log(`Generated popout URL: ${url}`);
  }, [roomUUID, setPopoutUrl, roomToken]);

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
              <Dropdown.Item href="#!" onClick={openPopout}>{t('lr.room.popchat2')}</Dropdown.Item>
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
            <div className={isMobile ? 'ai-talk-container-mobile' : 'ai-talk-container-pc'}
                 onTouchStart={startRecording} onTouchEnd={stopRecording} disabled={processing || showUserConfig}>
              {!processing ?
                <div>
                  <div className={micWorking ? 'ai-talk-gn-active' : 'ai-talk-gn-normal'}>
                    <div className='ai-talk-mc'></div>
                  </div>
                </div> :
                <div>
                  <Spinner animation="border" variant="light" className='ai-talk-spinner'></Spinner>
                </div>}
            </div>
          </>}
        </Card.Body>
      </Card>
    </div>
  );
}

function AITalkTraceLogPC({traceLogs, traceCount, children, roomUUID, roomToken, fullscreen}) {
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

    const r0 = Math.random().toString(16).slice(-8);
    const created = new Date().toISOString();
    const url = `${window.PUBLIC_URL}/${Locale.current()}/routers-popout?app=ai-talk&popout=1&assistant=0&room=${roomUUID}&created=${created}&random=${r0}&roomToken=${roomToken}`;
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
          <div className={fullscreen ? 'ai-talk-trace-logs-pcfs' : 'ai-talk-trace-logs-pc'} ref={logPanelRef}>
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

function AITalkTraceLogMobile({traceLogs, traceCount, fullscreen}) {
  // Scroll the log panel.
  const logPanelRef = React.useRef(null);
  React.useEffect(() => {
    if (!logPanelRef?.current) return;
    console.log(`Logs scroll to end, height=${logPanelRef.current.scrollHeight}, logs=${traceLogs.length}, count=${traceCount}`);
    logPanelRef.current.scrollTo(0, logPanelRef.current.scrollHeight);
  }, [traceLogs, logPanelRef, traceCount]);

  return (
    <div className={fullscreen ? 'ai-talk-trace-logs-mobilefs' : 'ai-talk-trace-logs-mobile'} ref={logPanelRef}>
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

function AITalkTraceLogPopout({traceLogs, traceCount}) {
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
