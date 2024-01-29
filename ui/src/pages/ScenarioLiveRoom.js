//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {useSrsLanguage} from "../components/LanguageSwitch";
import {Accordion, Button, Card, Form, Nav, Spinner, Table} from "react-bootstrap";
import {useTranslation} from "react-i18next";
import axios from "axios";
import {Clipboard, Token} from "../utils";
import {useErrorHandler} from "react-error-boundary";
import {useSearchParams} from "react-router-dom";
import {buildUrls} from "../components/UrlGenerator";
import {SrsEnvContext} from "../components/SrsEnvContext";
import * as Icon from "react-bootstrap-icons";
import PopoverConfirm from "../components/PopoverConfirm";
import {OpenAISecretSettings} from "../components/OpenAISettings";
import {AITalkAssistantPanel} from "../components/AITalk";

export default function ScenarioLiveRoom() {
  const [searchParams] = useSearchParams();
  // The room id, to maintain a specified room.
  const [roomId, setRoomId] = React.useState();

  React.useEffect(() => {
    const id = searchParams.get('roomid') || null;
    console.log(`?roomid=xxx, current=${id}, Set the roomid to manage.`);
    setRoomId(id);
  }, [searchParams, setRoomId]);

  if (roomId) return <ScenarioLiveRoomImpl {...{setRoomId, roomId}} />;
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
              <th>Stream</th>
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
                <td>{room.stream}</td>
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

function ScenarioLiveRoomImpl({roomId, setRoomId}) {
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

  if (!room) return <Spinner animation="border" variant="primary" />;
  return <>
    <Accordion defaultActiveKey={room.assistant ? ['3'] : ['0', '1', '2', '3']} alwaysOpen>
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
          <LiveRoomAssistantStage {...{room}}/>
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
    const urls = buildUrls(`live/${room.stream}`, {publish: room.secret}, env);
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
  const [aiAsrEnabled, setAiAsrEnabled] = React.useState(room.aiAsrEnabled);
  const [aiChatEnabled, setAiChatEnabled] = React.useState(room.aiChatEnabled);
  const [aiTtsEnabled, setAiTtsEnabled] = React.useState(room.aiTtsEnabled);
  const [aiAsrLanguage, setAiAsrLanguage] = React.useState(room.aiAsrLanguage || language);
  const [aiChatModel, setAiChatModel] = React.useState(room.aiChatModel || 'gpt-4-turbo-preview');
  const [aiChatPrompt, setAiChatPrompt] = React.useState(room.aiChatPrompt || 'You are a helpful assistant.');
  const [aiChatMaxWindow, setAiChatMaxWindow] = React.useState(room.aiChatMaxWindow || 5);
  const [aiChatMaxWords, setAiChatMaxWords] = React.useState(room.aiChatMaxWords || 300);

  const [configItem, setConfigItem] = React.useState('basic');

  const changeConfigItem = React.useCallback((e, t) => {
    e.preventDefault();
    setConfigItem(t);
  }, [setConfigItem]);

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
      <Card>
        <Card.Header>
          <Nav variant="tabs" defaultActiveKey="#basic">
            <Nav.Item>
              <Nav.Link href="#basic" onClick={(e) => changeConfigItem(e, 'basic')}>{t('lr.room.basic')}</Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link href="#provider" onClick={(e) => changeConfigItem(e, 'provider')}>{t('lr.room.provider')}</Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link href="#asr" onClick={(e) => changeConfigItem(e, 'asr')}>{t('lr.room.asr')}</Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link href="#chat" onClick={(e) => changeConfigItem(e, 'chat')}>{t('lr.room.chat')}</Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link href="#tts" onClick={(e) => changeConfigItem(e, 'tts')}>{t('lr.room.tts')}</Nav.Link>
            </Nav.Item>
          </Nav>
        </Card.Header>
        {configItem === 'basic' && <Card.Body>
          <Form.Group className="mb-3">
            <Form.Label>{t('lr.room.name')}</Form.Label>
            <Form.Text> * {t('lr.room.name2')}</Form.Text>
            <Form.Control as="input" type='input' defaultValue={aiName} onChange={(e) => setAiName(e.target.value)} />
          </Form.Group>
        </Card.Body>}
        {configItem === 'provider' && <Card.Body>
          <Form.Group className="mb-3">
            <Form.Label>{t('lr.room.provider')}</Form.Label>
            <Form.Text> * {t('lr.room.provider2')}</Form.Text>
            <Form.Select defaultValue={aiProvider} onChange={(e) => setAiProvider(e.target.value)}>
              <option value="">--{t('helper.noSelect')}--</option>
              <option value="openai">OpenAI</option>
            </Form.Select>
          </Form.Group>
          <OpenAISecretSettings {...{
            baseURL: aiBaseURL, setBaseURL: setAiBaseURL,
            secretKey: aiSecretKey, setSecretKey: setAiSecretKey,
            targetLanguage: aiAsrLanguage, setTargetLanguage: setAiAsrLanguage
          }} />
        </Card.Body>}
        {configItem === 'asr' && <Card.Body>
          <Form.Group className="mb-3">
            <Form.Group className="mb-3" controlId="formAiAsrEnabledCheckbox">
              <Form.Check type="checkbox" label={t('lr.room.asre')} defaultChecked={aiAsrEnabled} onClick={() => setAiAsrEnabled(!aiAsrEnabled)} />
            </Form.Group>
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>{t('transcript.lang')}</Form.Label>
            <Form.Text> * {t('transcript.lang2')}. &nbsp;
              {t('helper.eg')} <code>en, zh, fr, de, ja, ru </code>, ... &nbsp;
              {t('helper.see')} <a href='https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes' target='_blank' rel='noreferrer'>ISO-639-1</a>.
            </Form.Text>
            <Form.Control as="input" defaultValue={aiAsrLanguage} onChange={(e) => setAiAsrLanguage(e.target.value)} />
          </Form.Group>
        </Card.Body>}
        {configItem === 'chat' && <Card.Body>
          <Form.Group className="mb-3">
            <Form.Group className="mb-3" controlId="formAiChatEnabledCheckbox">
              <Form.Check disabled={true} type="checkbox" label={t('lr.room.chate')} defaultChecked={aiChatEnabled} onClick={() => setAiChatEnabled(!aiChatEnabled)} />
            </Form.Group>
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>{t('lr.room.model')}</Form.Label>
            <Form.Text> * {t('lr.room.model2')}</Form.Text>
            <Form.Control as="input" type='input' defaultValue={aiChatModel} onChange={(e) => setAiChatModel(e.target.value)} />
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>{t('lr.room.prompt')}</Form.Label>
            <Form.Text> * {t('lr.room.prompt2')}</Form.Text>
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
        </Card.Body>}
        {configItem === 'tts' && <Card.Body>
          <Form.Group className="mb-3">
            <Form.Group className="mb-3" controlId="formAiTtsEnabledCheckbox">
              <Form.Check type="checkbox" label={t('lr.room.ttse')} defaultChecked={aiTtsEnabled} onClick={() => setAiTtsEnabled(!aiTtsEnabled)} />
            </Form.Group>
          </Form.Group>
        </Card.Body>}
      </Card>
      <p></p>
      <Button variant="primary" type="button" disabled={requesting}
              onClick={(e) => updateRoom({
                ...room, assistant: true,
                aiName, aiProvider, aiSecretKey, aiBaseURL, aiAsrLanguage, aiChatModel,
                aiChatPrompt, aiChatMaxWindow: parseInt(aiChatMaxWindow),
                aiChatMaxWords: parseInt(aiChatMaxWords), aiAsrEnabled: !!aiAsrEnabled,
                aiChatEnabled: !!aiChatEnabled, aiTtsEnabled: !!aiTtsEnabled,
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

function LiveRoomAssistantStage({room}) {
  const {t} = useTranslation();

  if (!room.assistant || !room.aiProvider || !room.aiSecretKey || !room.aiBaseURL || !room.aiAsrLanguage
    || !room.aiChatModel || !room.aiChatPrompt) {
    return <>{t('lr.room.aiwe')}</>;
  }
  return <AITalkAssistantPanel {...{roomUUID: room.uuid, roomToken: room.roomToken, fullscreen: false}} />;
}
