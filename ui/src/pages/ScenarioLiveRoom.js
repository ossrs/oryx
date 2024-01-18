//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {useSrsLanguage} from "../components/LanguageSwitch";
import {Accordion, Alert, Button, Card, Col, Form, Nav, Row, Spinner, Table} from "react-bootstrap";
import {useTranslation} from "react-i18next";
import axios from "axios";
import {Clipboard, Token} from "../utils";
import {useErrorHandler} from "react-error-boundary";
import {useSearchParams} from "react-router-dom";
import {buildUrls} from "../components/UrlGenerator";
import {SrsEnvContext} from "../components/SrsEnvContext";
import * as Icon from "react-bootstrap-icons";
import PopoverConfirm from "../components/PopoverConfirm";
import {OpenAIWhisperSettings} from "../components/OpenAISettings";

export default function ScenarioLiveRoom() {
  const [searchParams] = useSearchParams();
  // The room id, to maintain a specified room.
  const [roomId, setRoomId] = React.useState();

  React.useEffect(() => {
    const id = searchParams.get('roomid') || null;
    console.log(`?roomid=xxx, current=${id}, Set the roomid to manage.`);
    setRoomId(id);
  }, [searchParams, setRoomId]);

  if (roomId) return <ScenarioLiveRoomManager {...{setRoomId, roomId}} />;
  return <ScenarioLiveRoomList {...{setRoomId}} />;
}

function ScenarioLiveRoomList({setRoomId}) {
  const language = useSrsLanguage();
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const [searchParams, setSearchParams] = useSearchParams();
  const [name, setName] = React.useState('My Live Room');
  const [rooms, setRooms] = React.useState([]);
  const [refreshNow, setRefreshNow] = React.useState();

  const createLiveRoom = React.useCallback((e) => {
    e.preventDefault();

    axios.post('/terraform/v1/live/room/create', {
      title: name,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      const {uuid} = res.data.data;
      searchParams.set('roomid', uuid); setSearchParams(searchParams);
      setRoomId(uuid);
      console.log(`Status: Create ok, name=${name}, data=${JSON.stringify(res.data.data)}`);
    }).catch(handleError);
  }, [handleError, name, setRoomId, searchParams, setSearchParams]);

  const removeRoom = React.useCallback((uuid) => {
    axios.post('/terraform/v1/live/room/remove', {
      uuid: uuid,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      setRefreshNow(!refreshNow);
      console.log(`Status: Remove ok, uuid=${uuid}, data=${JSON.stringify(res.data.data)}`);
    }).catch(handleError);
  }, [handleError, refreshNow, setRefreshNow]);

  const manageRoom = React.useCallback((room) => {
    const uuid = room.uuid;
    searchParams.set('roomid', uuid); setSearchParams(searchParams);
    setRoomId(room.uuid);
  }, [searchParams, setSearchParams, setRoomId]);

  React.useEffect(() => {
    const refreshLiveRoomsTask = () => {
      axios.post('/terraform/v1/live/room/list', {
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        const {rooms} = res.data.data;
        setRooms(rooms || []);
        console.log(`Status: List ok, data=${JSON.stringify(res.data.data)}`);
      }).catch(handleError);
    };

    refreshLiveRoomsTask();
    const timer = setInterval(() => refreshLiveRoomsTask(), 3 * 1000);
    return () => {
      clearInterval(timer);
      setRooms([]);
    }
  }, [handleError, setRooms, refreshNow]);

  return (
    <Accordion defaultActiveKey={['1', '2']}>
      <React.Fragment>
        {language === 'zh' ?
          <Accordion.Item eventKey="0">
            <Accordion.Header>场景介绍</Accordion.Header>
            <Accordion.Body>
              <div>直播间，提供了按每个流鉴权的能力，并支持直播间的业务功能。</div>
              <p></p>
              <p>可应用的具体场景包括：</p>
              <ul>
                <li>自建直播间，私域直播，仅限私域会员能观看的直播。</li>
                <li>企业直播，企业内部的直播间，仅限企业内部人员观看。</li>
                <li>电商直播，仅限电商特定买家可观看的直播。</li>
              </ul>
            </Accordion.Body>
          </Accordion.Item> :
          <Accordion.Item eventKey="0">
            <Accordion.Header>Scenario Introduction</Accordion.Header>
            <Accordion.Body>
              <div>Live room, which provides the ability to authenticate each stream and supports business functions of live room.</div>
              <p></p>
              <p>The specific scenarios that can be applied include:</p>
              <ul>
                <li>Self-built live room, private domain live broadcast, live broadcast that can only be watched by private domain members.</li>
                <li>Enterprise live broadcast, live room within the enterprise, only for internal personnel of the enterprise.</li>
                <li>E-commerce live broadcast, live broadcast that can only be watched by specific buyers of e-commerce.</li>
              </ul>
            </Accordion.Body>
          </Accordion.Item>}
      </React.Fragment>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{t('lr.create.title')}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>{t('lr.create.name')}</Form.Label>
              <Form.Text> * {t('lr.create.name2')}</Form.Text>
              <Form.Control as="input" defaultValue={name} onChange={(e) => setName(e.target.value)} />
            </Form.Group>
            <Button ariant="primary" type="submit" onClick={(e) => createLiveRoom(e)}>
              {t('helper.create')}
            </Button>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>{t('lr.list.title')}</Accordion.Header>
        <Accordion.Body>
          {rooms?.length ? <Table striped bordered hover>
            <thead>
            <tr>
              <th>#</th>
              <th>UUID</th>
              <th>Title</th>
              <th>Created At</th>
              <th>Actions</th>
            </tr>
            </thead>
            <tbody>
            {rooms?.map((room, index) => {
              return <tr key={room.uuid}>
                <td>{index}</td>
                <td>
                  <a href="#!" onClick={(e) => {
                    e.preventDefault();
                    manageRoom(room);
                  }}>{room.uuid}</a>
                </td>
                <td>{room.title}</td>
                <td>{room.created_at}</td>
                <td>
                  <a href="#!" onClick={(e) => {
                    e.preventDefault();
                    manageRoom(room);
                  }}>{t('helper.manage')}</a> &nbsp;
                  <PopoverConfirm placement='top' trigger={ <a href='#!'>{t('helper.delete')}</a> } onClick={() => removeRoom(room.uuid)}>
                    <p>
                      {t('lr.list.delete')}
                    </p>
                  </PopoverConfirm>
                </td>
              </tr>;
            })}
            </tbody>
          </Table> : t('lr.list.empty')}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

function ScenarioLiveRoomManager({roomId, setRoomId}) {
  const {t} = useTranslation();
  const handleError = useErrorHandler();

  const [requesting, setRequesting] = React.useState(false);
  const [room, setRoom] = React.useState();

  React.useEffect(() => {
    axios.post('/terraform/v1/live/room/query', {
      uuid: roomId,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      setRoom(res.data.data);
      console.log(`Room: Query ok, uuid=${roomId}, data=${JSON.stringify(res.data.data)}`);
    }).catch(handleError);
  }, [handleError, roomId, setRoom]);

  const updateRoom = React.useCallback((room) => {
    setRequesting(true);
    try {
      axios.post('/terraform/v1/live/room/update', {
        uuid: room.uuid, ...room,
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        alert(t('helper.setOk'));
        setRoom(res.data.data);
        console.log(`Room: Update ok, uuid=${room.uuid}, data=${JSON.stringify(res.data.data)}`);
      }).catch(handleError);
    } finally {
      setRequesting(false);
    }
  }, [t, handleError, setRequesting, setRoom]);

  // TODO: FIXME: Change to ['0', '1', '2', '3']
  const defaultActiveKey = ['3'];
  return <>
    <Accordion defaultActiveKey={defaultActiveKey} alwaysOpen>
      <Accordion.Item eventKey="0">
        <Accordion.Header>{t('lr.room.nav')}</Accordion.Header>
        <Accordion.Body>
          <Button variant="link" onClick={() => setRoomId(null)}>Back to Rooms</Button>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{t('lr.room.stream')}</Accordion.Header>
        <Accordion.Body>
          {room ? <LiveRoomStreamer {...{room}}/> : ''}
        </Accordion.Body>
      </Accordion.Item>
      {room ? <Accordion.Item eventKey="2">
        <Accordion.Header>{t('lr.room.aic')}</Accordion.Header>
        <Accordion.Body>
          <LiveRoomAssistantConfiguration {...{room, requesting, updateRoom}}/>
        </Accordion.Body>
      </Accordion.Item> : ''}
      {room && room.assistant ? <Accordion.Item eventKey="3">
        <Accordion.Header>{t('lr.room.aiw')}</Accordion.Header>
        <Accordion.Body>
          <LiveRoomAssistantWorker {...{room}}/>
        </Accordion.Body>
      </Accordion.Item> : ''}
    </Accordion>
  </>;
}

function LiveRoomStreamer({room}) {
  const {t} = useTranslation();
  const env = React.useContext(SrsEnvContext)[0];

  const [urls, setUrls] = React.useState({});
  const [streamType, setStreamType] = React.useState('rtmp');

  const copyToClipboard = React.useCallback((e, text) => {
    e.preventDefault();

    Clipboard.copy(text).then(() => {
      alert(t('helper.copyOk'));
    }).catch((err) => {
      alert(`${t('helper.copyFail')} ${err}`);
    });
  }, [t]);

  const changeStreamType = React.useCallback((e, t) => {
    e.preventDefault();
    setStreamType(t);
  }, [setStreamType]);

  React.useEffect(() => {
    if (!room?.secret) return;
    const urls = buildUrls(`live/${room.uuid}`, {publish: room.secret}, env);
    setUrls(urls);
  }, [room, env, setUrls]);

  const {
    rtmpServer, rtmpStreamKey, hlsPlayer, m3u8Url, srtPublishUrl,
  } = urls;

  return (
    <Card>
      <Card.Header>
        <Nav variant="tabs" defaultActiveKey="#rtmp">
          <Nav.Item>
            <Nav.Link href="#rtmp" onClick={(e) => changeStreamType(e, 'rtmp')}>{t('live.obs.title')}</Nav.Link>
          </Nav.Item>
          <Nav.Item>
            <Nav.Link href="#srt" onClick={(e) => changeStreamType(e, 'srt')}>{t('live.srt.title')}</Nav.Link>
          </Nav.Item>
        </Nav>
      </Card.Header>
      {streamType === 'rtmp' ? <Card.Body>
          <div>
            {t('live.obs.server')} <code>{rtmpServer}</code> &nbsp;
            <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
              <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtmpServer)} />
            </div>
          </div>
          <div>
            {t('live.obs.key')} <code>{rtmpStreamKey}</code> &nbsp;
            <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
              <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, rtmpStreamKey)} />
            </div>
          </div>
          <div>
            {t('live.share.hls')}&nbsp;
            <a href={hlsPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>,&nbsp;
            <code>{m3u8Url}</code> &nbsp;
            <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
              <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8Url)} />
            </div>
          </div>
        </Card.Body> :
        <Card.Body>
          <div>
            {t('live.obs.server')} <code>{srtPublishUrl}</code> &nbsp;
            <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
              <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, srtPublishUrl)} />
            </div>
          </div>
          <div>
            {t('live.obs.key')} <code>{t('live.obs.nokey')}</code>
          </div>
          <div>
            {t('live.share.hls')}&nbsp;
            <a href={hlsPlayer} target='_blank' rel='noreferrer'>{t('live.share.simple')}</a>,&nbsp;
            <code>{m3u8Url}</code> &nbsp;
            <div role='button' style={{display: 'inline-block'}} title={t('helper.copy')}>
              <Icon.Clipboard size={20} onClick={(e) => copyToClipboard(e, m3u8Url)} />
            </div>
          </div>
        </Card.Body>}
    </Card>
  );
}

function LiveRoomAssistantConfiguration({room, requesting, updateRoom}) {
  const language = useSrsLanguage();
  const {t} = useTranslation();

  const [aiName, setAiName] = React.useState(room.aiName);
  const [aiProvider, setAiProvider] = React.useState(room.aiProvider || 'openai');
  const [aiSecretKey, setAiSecretKey] = React.useState(room.aiSecretKey);
  const [aiBaseURL, setAiBaseURL] = React.useState(room.aiBaseURL || (language === 'zh' ? '' : 'https://api.openai.com/v1'));
  const [aiAsrLanguage, setAiAsrLanguage] = React.useState(room.aiAsrLanguage || language);
  const [aiChatModel, setAiChatModel] = React.useState(room.aiChatModel || 'gpt-3.5-turbo-1106');
  const [aiChatPrompt, setAiChatPrompt] = React.useState(room.aiChatPrompt || 'You are a helpful assistant.');
  const [aiChatMaxWindow, setAiChatMaxWindow] = React.useState(room.aiChatMaxWindow || 5);
  const [aiChatMaxWords, setAiChatMaxWords] = React.useState(room.aiChatMaxWords || 30);

  if (!room.assistant) {
    return (
      <Button variant="primary" type="button" disabled={requesting}
              onClick={(e) => updateRoom({...room, assistant: true})}>
        {t('lr.room.enable')}
      </Button>
    );
  }
  return (
    <Form>
      <Form.Group className="mb-3">
        <Form.Label>{t('lr.room.name')}</Form.Label>
        <Form.Text> * {t('lr.room.name2')}</Form.Text>
        <Form.Control as="input" type='input' defaultValue={aiName} onChange={(e) => setAiName(e.target.value)} />
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>{t('lr.room.provider')}</Form.Label>
        <Form.Text> * {t('lr.room.provider2')}</Form.Text>
        <Form.Select defaultValue={aiProvider} onChange={(e) => setAiProvider(e.target.value)}>
          <option value="">--{t('helper.noSelect')}--</option>
          <option value="openai">OpenAI</option>
        </Form.Select>
      </Form.Group>
      <OpenAIWhisperSettings {...{
        baseURL: aiBaseURL, setBaseURL: setAiBaseURL,
        secretKey: aiSecretKey, setSecretKey: setAiSecretKey,
        targetLanguage: aiAsrLanguage, setTargetLanguage: setAiAsrLanguage
      }} />
      <Form.Group className="mb-3">
        <Form.Label>{t('lr.room.model')}</Form.Label>
        <Form.Text> * {t('lr.room.model2')}</Form.Text>
        <Form.Control as="input" type='input' defaultValue={aiChatModel} onChange={(e) => setAiChatModel(e.target.value)} />
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>{t('lr.room.prmpt')}</Form.Label>
        <Form.Text> * {t('lr.room.prmpt2')}</Form.Text>
        <Form.Control as="textarea" type='text' rows={3}  defaultValue={aiChatPrompt} onChange={(e) => setAiChatPrompt(e.target.value)} />
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>{t('lr.room.window')}</Form.Label>
        <Form.Text> * {t('lr.room.window2')}</Form.Text>
        <Form.Control as="input" type='input' defaultValue={aiChatMaxWindow} onChange={(e) => setAiChatMaxWindow(e.target.value)} />
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>{t('lr.room.words')}</Form.Label>
        <Form.Text> * {t('lr.room.words2')}</Form.Text>
        <Form.Control as="input" type='input' defaultValue={aiChatMaxWords} onChange={(e) => setAiChatMaxWords(e.target.value)} />
      </Form.Group>
      <Button variant="primary" type="button" disabled={requesting}
              onClick={(e) => updateRoom({
                ...room, assistant: true,
                aiName, aiProvider, aiSecretKey, aiBaseURL, aiAsrLanguage, aiChatModel,
                aiChatPrompt, aiChatMaxWindow, aiChatMaxWords,
              })}>
        {t('lr.room.update')}
      </Button> &nbsp;
      <Button variant="primary" type="button" disabled={requesting}
              onClick={(e) => updateRoom({...room, assistant: false})}>
        {t('lr.room.disable')}
      </Button>
    </Form>
  );
}

function useIsMobile() {
  const [isMobile, setIsMobile] = React.useState(false);

  function handleWindowSizeChange() {
    setIsMobile(window.innerWidth <= 768);
  }
  React.useEffect(() => {
    handleWindowSizeChange();
    window.addEventListener('resize', handleWindowSizeChange);
    return () => {
      window.removeEventListener('resize', handleWindowSizeChange);
    }
  }, [setIsMobile]);

  return isMobile;
}

function LiveRoomAssistantWorker({room}) {
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const isMobile = useIsMobile();

  // The timeout in milliseconds.
  const timeoutForMicrophoneTestToRun = 50;
  const timeoutWaitForMicrophoneToClose = 600;
  const timeoutWaitForLastVoice = 700;
  const durationRequiredUserInput = 600;

  // The player ref, to access the audio player.
  const playerRef = React.useRef(null);
  const [robotReady, setRobotReady] = React.useState(false);
  const [processing, setProcessing] = React.useState(false);
  const [micWorking, setMicWorking] = React.useState(false);

  // The uuid and robot in stage, which is unchanged after stage started.
  const [stageRobot, setStageRobot] = React.useState(null);
  const [stageUUID, setStageUUID] = React.useState(null);

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
        last.msg = `${last.msg} ${msg}`;
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
  const logPanelRef = React.useRef(null);
  const endPanelRef = React.useRef(null);
  React.useEffect(() => {
    if (!logPanelRef?.current) return;
    console.log(`Logs scroll to end, height=${logPanelRef.current.scrollHeight}, logs=${traceLogs.length}, count=${traceCount}`);
    logPanelRef.current.scrollTo(0, logPanelRef.current.scrollHeight);
  }, [traceLogs, logPanelRef, traceCount]);
  React.useEffect(() => {
    if (!robotReady || !endPanelRef?.current) return;
    console.log(`Logs setup to end, height=${endPanelRef.current.scrollHeight}, tips=${tipLogs.length}`);
    endPanelRef.current.scrollIntoView({ behavior: "smooth", block: "end" });
  }, [robotReady, endPanelRef, tipLogs]);

  // The application is started now.
  React.useEffect(() => {
    // Only allow localhost or https to access microphone.
    const isLo = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
    const isHttps = window.location.protocol === 'https:';
    const securityAllowed = isLo || isHttps;
    securityAllowed || errorLog(t('lr.room.https'));
    console.log(`App started, allowed=${securityAllowed}, lo=${isLo}, https=${isHttps}`);
    if (!securityAllowed) return;

    // Try to open the microphone to request permission.
    new Promise(resolve => {
      console.log(`Start: Open microphone`);

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
  }, [t, errorLog, setBooting]);

  // Request server to create a new stage.
  React.useEffect(() => {
    if (booting) return;

    console.log(`Start: Create a new stage`);
    axios.post('/terraform/v1/ai-talk/stage/start', {
      room: room.uuid,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      console.log(`Start: Create stage success: ${JSON.stringify(res.data.data)}`);
      setStageUUID(res.data.data.sid);
      setStageRobot(res.data.data.robot);
    }).catch(handleError);
  }, [handleError, booting, room, setStageUUID, setStageRobot]);

  // Start to chat, set the robot to ready.
  const startChatting = React.useCallback(() => {
    const listener = () => {
      playerRef.current.removeEventListener('ended', listener);

      setRobotReady(true);
      console.log(`Stage started, AI is ready, sid=${stageUUID}`);
    };
    playerRef.current.addEventListener('ended', listener);

    playerRef.current.src = `/terraform/v1/ai-talk/stage/examples/${stageRobot.voice}?sid=${stageUUID}`;
    playerRef.current.play().catch(error => errorLog(`${t('lr.room.speaker')}: ${error}`));
  }, [t, errorLog, stageUUID, setRobotReady, stageRobot]);

  // For test only, append some logs.
  const appendTestLogs = React.useCallback((logs) => {
    traceLog('You', 'Hello', 'primary');
    traceLog('Bot', `World ${new Date()}`, 'success');
  }, [traceLog]);

  // When robot is ready, open the microphone ASAP to accept user input.
  React.useEffect(() => {
    if (!robotReady) return;
    if (ref.current.mediaStream) return;

    console.log(`Robot is ready, open microphone.`)
    navigator.mediaDevices.getUserMedia(
      { audio: true }
    ).then((stream) => {
      ref.current.mediaStream = stream;
      console.log(`Robot is ready, microphone opened.`);
    }).catch(error => errorLog(`${t('lr.room.mic')}: ${error}`));
  }, [errorLog, t, robotReady, ref]);

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
  }, [robotReady, ref, setMicWorking]);

  // User click stop button, we delay some time to allow cancel the stopping event.
  const stopRecording = React.useCallback(async () => {
    if (!robotReady) return;

    const processUserInput = async(userMayInput) => {
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
      const requestUUID = await new Promise((resolve, reject) => {
        console.log(`ASR: Uploading ${ref.current.audioChunks.length} chunks, robot=${stageRobot.uuid}`);
        ref.current.audioChunks = [];

        axios.post('/terraform/v1/ai-talk/stage/upload', {
          sid: stageUUID, robot: stageRobot.uuid, umi: userMayInput, audio: audioBase64Data,
        }, {
          headers: Token.loadBearerHeader(),
        }).then(res => {
          console.log(`ASR: Upload success: ${res.data.data.rid} ${res.data.data.asr}`);
          traceLog('You', res.data.data.asr, 'primary');
          resolve(res.data.data.rid);
        }).catch((error) => reject(error));
      });

      // Get the AI generated audio from the server.
      while (true) {
        console.log(`TTS: Requesting ${requestUUID} response audios, rid=${requestUUID}`);
        let audioSegmentUUID = null;
        while (!audioSegmentUUID) {
          const resp = await new Promise((resolve, reject) => {
            axios.post('/terraform/v1/ai-talk/stage/query', {
              sid: stageUUID, rid: requestUUID,
            }, {
              headers: Token.loadBearerHeader(),
            }).then(res => {
              if (res.data?.data?.asid) {
                console.log(`TTS: Audio ready: ${res.data.data.asid} ${res.data.data.tts}`);
                traceLog('Bot', res.data.data.tts, 'success');
              }
              resolve(res.data.data);
            }).catch(error => reject(error));
          });

          if (!resp.asid) {
            break;
          }

          if (resp.processing) {
            await new Promise((resolve) => setTimeout(resolve, 300));
            continue;
          }

          audioSegmentUUID = resp.asid;
        }

        // All audios are played.
        if (!audioSegmentUUID) {
          console.log(`TTS: All audios are played, rid=${requestUUID}`);
          console.log("=============");
          break;
        }

        // Play the AI generated audio.
        await new Promise(resolve => {
          const url = `/terraform/v1/ai-talk/stage/tts?sid=${stageUUID}&rid=${requestUUID}&asid=${audioSegmentUUID}`;
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
          axios.post('/terraform/v1/ai-talk/stage/remove', {
            sid: stageUUID, rid: requestUUID, asid: audioSegmentUUID,
          }, {
            headers: Token.loadBearerHeader(),
          }).then(res => {
            console.log(`TTS: Audio removed: ${audioSegmentUUID}`);
            resolve();
          }).catch(error => reject(error));
        });
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

        // Reopen the microphone.
        console.log(`Robot is ready, open microphone.`)
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
  }, [playerRef, stageUUID, stageRobot, robotReady, ref, setProcessing, setMicWorking, traceLog]);

  // Setup the keyboard event, for PC browser.
  React.useEffect(() => {
    if (!robotReady) return;

    const handleKeyDown = (e) => {
      if (processing) return;
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
  }, [robotReady, startRecording, stopRecording, processing]);

  if (!room.assistant || !room.aiProvider || !room.aiSecretKey || !room.aiBaseURL || !room.aiAsrLanguage
    || !room.aiChatModel || !room.aiChatPrompt) {
    return <>{t('lr.room.aiwe')}</>;
  }
  return (
    <div>
      {errorLogs.map((log) => {
        return (
          <Alert key={log.id} onClose={() => removeErrorLog(log)} variant='danger' dismissible>
            <Alert.Heading>Error!</Alert.Heading>
            <p>{log.msg}</p>
          </Alert>
        );
      })}
      {tipLogs.map((log) => {
        return (
          <Alert key={log.id} onClose={() => removeTipLog(log)} variant='success' dismissible>
            <Alert.Heading>{log.title}</Alert.Heading>
            <p>{log.msg}</p>
          </Alert>
        );
      })}
      <div><audio ref={playerRef} controls={true} hidden='hidden' /></div>
      {booting ? <>Booting ...</> : ''}
      {stageUUID && !robotReady ? <Button variant="primary" type="submit" onClick={startChatting}>{t('lr.room.talk')}</Button> : ''}
      {stageUUID && robotReady ? <Button variant="primary" type="submit" onClick={appendTestLogs} hidden={true}>Test Logs</Button> : ''}
      {robotReady && !isMobile ?
        <Row>
          <Col>
            <div className='ai-talk-container-pc' onTouchStart={startRecording} onTouchEnd={stopRecording} disabled={processing}>
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
          </Col>
          <Col>
            <div className='ai-talk-trace-logs-pc' ref={logPanelRef}>
              <p></p>
              <div>
                {traceLogs.map((log) => {
                  return (
                    <Alert key={log.id} variant={log.variant}>
                      {log.role}: {log.msg}
                    </Alert>
                  );
                })}
              </div>
            </div>
            <div ref={endPanelRef}></div>
          </Col>
        </Row> : ''}
      {robotReady && isMobile ?
        <Row>
          <Col>
            <div className='ai-talk-trace-logs-mobile' ref={logPanelRef}>
              <p></p>
              <div>
                {traceLogs.map((log) => {
                  return (
                    <Alert key={log.id} variant={log.variant}>
                      {log.role}: {log.msg}
                    </Alert>
                  );
                })}
              </div>
            </div>
            <div className="ai-talk-container-mobile" onTouchStart={startRecording} onTouchEnd={stopRecording} disabled={processing}>
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
            <div ref={endPanelRef}></div>
          </Col>
        </Row> : ''}
    </div>
  );
}
