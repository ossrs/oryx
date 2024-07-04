import React from "react";
import {Accordion, Button, Form, Nav, Spinner, Table, Card} from "react-bootstrap";
import {useSrsLanguage} from "../components/LanguageSwitch";
import {useTranslation} from "react-i18next";
import {Locale, Token} from "../utils";
import axios from "axios";
import {useErrorHandler} from "react-error-boundary";
import {OpenAISecretSettings} from "../components/OpenAISettings";
import {useLocation, useNavigate} from "react-router-dom";

export default function ScenarioOCR(props) {
  const handleError = useErrorHandler();
  const [config, setConfig] = React.useState();
  const [uuid, setUuid] = React.useState();
  const [activeKey, setActiveKey] = React.useState();

  React.useEffect(() => {
    axios.post('/terraform/v1/ai/ocr/query', {
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      const data = res.data.data;

      setConfig(data.config);
      setUuid(data.task.uuid);

      if (data.config.all) {
        setActiveKey(['2', '3', '4', '5', '6']);
      } else {
        setActiveKey(['1']);
      }

      console.log(`OCR: Query ok, ${JSON.stringify(data)}`);
    }).catch(handleError);
  }, [handleError, setActiveKey, setConfig,  setUuid]);

  if (!activeKey) return <></>;
  return <ScenarioOCRImpl {...props} {...{
    activeKey, defaultEnabled: config?.all, defaultConf: config, defaultUuid: uuid,
  }}/>;
}

function ScenarioOCRImpl({activeKey, defaultEnabled, defaultConf, defaultUuid}) {
  const language = useSrsLanguage();
  const {t} = useTranslation();
  const handleError = useErrorHandler();
  const navigate = useNavigate();
  const location = useLocation();

  const [operating, setOperating] = React.useState(false);
  const [ocrEnabled, setOcrEnabled] = React.useState(defaultEnabled);
  const [aiProvider, setAiProvider] = React.useState(defaultConf.aiProvider || 'openai');
  const [aiSecretKey, setAiSecretKey] = React.useState(defaultConf.aiSecretKey);
  const [aiOrganization, setAiOrganization] = React.useState(defaultConf.aiOrganization);
  const [aiBaseURL, setAiBaseURL] = React.useState(defaultConf.aiBaseURL || (language === 'zh' ? '' : 'https://api.openai.com/v1'));
  const [aiChatEnabled, setAiChatEnabled] = React.useState(defaultConf.aiChatEnabled);
  const [aiChatModel, setAiChatModel] = React.useState(defaultConf.aiChatModel || 'gpt-4o');
  const [aiChatPrompt, setAiChatPrompt] = React.useState(defaultConf.aiChatPrompt || 'Recognize the text in the image. Output the identified text directly.');
  const [aiChatMaxWindow, setAiChatMaxWindow] = React.useState(defaultConf.aiChatMaxWindow || 5);
  const [aiChatMaxWords, setAiChatMaxWords] = React.useState(defaultConf.aiChatMaxWords || 300);

  const [liveQueue, setLiveQueue] = React.useState();
  const [ocrQueue, setOcrQueue] = React.useState();
  const [callbackQueue, setCallbackQueue] = React.useState();
  const [cleanupQueue, setCleanupQueue] = React.useState();
  const [lastObject, setLastObject] = React.useState();
  const [uuid, setUuid] = React.useState(defaultUuid);

  const [configItem, setConfigItem] = React.useState('provider');
  const changeConfigItem = React.useCallback((e, t) => {
    e.preventDefault();
    setConfigItem(t);
  }, [setConfigItem]);

  React.useEffect(() => {
    if (aiSecretKey) return;

    axios.post('/terraform/v1/mgmt/openai/query', null, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      const data = res.data.data;
      setAiSecretKey(data.aiSecretKey);
      setAiBaseURL(data.aiBaseURL);
      setAiOrganization(data.aiOrganization);
      console.log(`OCR: Query open ai ok, data=${JSON.stringify(data)}`);
    }).catch(handleError);
  }, [handleError, aiSecretKey, setAiSecretKey, setAiBaseURL, setAiOrganization]);

  const updateOcrService = React.useCallback((enabled, success) => {
    if (!aiSecretKey) return alert(`Invalid secret key ${aiSecretKey}`);
    if (!aiBaseURL) return alert(`Invalid base url ${aiBaseURL}`);

    axios.post('/terraform/v1/ai/ocr/apply', {
      uuid, all: !!enabled, aiProvider, aiSecretKey, aiOrganization, aiBaseURL,
      aiChatEnabled: !!aiChatEnabled, aiChatModel, aiChatPrompt,
      aiChatMaxWindow: parseInt(aiChatMaxWindow), aiChatMaxWords: parseInt(aiChatMaxWords),
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      alert(t('helper.setOk'));
      console.log(`OCR: Apply config ok, uuid=${uuid}.`);
      success && success();
    }).catch(handleError);
   }, [t, handleError, uuid, aiProvider, aiSecretKey, aiBaseURL, aiOrganization, aiChatEnabled, aiChatModel, aiChatPrompt, aiChatMaxWindow, aiChatMaxWords]);

  const resetTask = React.useCallback(() => {
    setOperating(true);

    axios.post('/terraform/v1/ai/ocr/reset', {
      uuid,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      alert(t('helper.setOk'));
      const data = res.data.data;
      setUuid(data.uuid);
      console.log(`OCR: Reset task ${uuid} ok: ${JSON.stringify(data)}`);
    }).catch(handleError).finally(setOperating);
  }, [t, handleError, uuid, setUuid, setOperating]);

  React.useEffect(() => {
    const refreshLiveQueueTask = () => {
      axios.post('/terraform/v1/ai/ocr/live-queue', {
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        const queue = res.data.data;
        queue.segments = queue?.segments?.map(segment => {
          return {
            ...segment,
            duration: Number(segment.duration),
            size: Number(segment.size / 1024.0 / 1024),
          };
        });
        setLiveQueue(queue);
        console.log(`OCR: Query live queue ${JSON.stringify(queue)}`);
      }).catch(handleError);
    };

    refreshLiveQueueTask();
    const timer = setInterval(() => refreshLiveQueueTask(), 3 * 1000);
    return () => clearInterval(timer);
  }, [handleError, setLiveQueue]);

  React.useEffect(() => {
    const refreshOcrQueueTask = () => {
      axios.post('/terraform/v1/ai/ocr/ocr-queue', {
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        const queue = res.data.data;
        queue.segments = queue?.segments?.map(segment => {
          return {
            ...segment,
            duration: Number(segment.duration),
            size: Number(segment.size / 1024.0),
            eic: Number(segment.eic),
          };
        });
        setOcrQueue(queue);
        console.log(`OCR: Query ocr queue ${JSON.stringify(queue)}`);
      }).catch(handleError);
    };

    refreshOcrQueueTask();
    const timer = setInterval(() => refreshOcrQueueTask(), 3 * 1000);
    return () => clearInterval(timer);
  }, [handleError, setOcrQueue]);

  React.useEffect(() => {
    const refreshOcrQueueTask = () => {
      axios.post('/terraform/v1/ai/ocr/callback-queue', {
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        const queue = res.data.data;
        queue.segments = queue?.segments?.map(segment => {
          return {
            ...segment,
            duration: Number(segment.duration),
            size: Number(segment.size / 1024.0),
            eic: Number(segment.eic),
            ocrc: Number(segment.ocrc),
          };
        });
        setCallbackQueue(queue);
        console.log(`OCR: Query callback queue ${JSON.stringify(queue)}`);
      }).catch(handleError);
    };

    refreshOcrQueueTask();
    const timer = setInterval(() => refreshOcrQueueTask(), 3 * 1000);
    return () => clearInterval(timer);
  }, [handleError, setCallbackQueue]);

  React.useEffect(() => {
    const refreshOcrQueueTask = () => {
      axios.post('/terraform/v1/ai/ocr/cleanup-queue', {
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        const queue = res.data.data;
        queue.segments = queue?.segments?.map(segment => {
          return {
            ...segment,
            duration: Number(segment.duration),
            size: Number(segment.size / 1024.0),
            eic: Number(segment.eic),
            ocrc: Number(segment.ocrc),
            cbc: Number(segment.cbc),
          };
        });
        setCleanupQueue(queue);
        setLastObject(queue?.segments?.length > 0 ? queue.segments[queue.segments.length - 1] : null);
        console.log(`OCR: Query cleanup queue ${JSON.stringify(queue)}`);
      }).catch(handleError);
    };

    refreshOcrQueueTask();
    const timer = setInterval(() => refreshOcrQueueTask(), 3 * 1000);
    return () => clearInterval(timer);
  }, [handleError, setCleanupQueue, setLastObject]);

  return (
    <Accordion defaultActiveKey={activeKey}>
      <React.Fragment>
        {language === 'zh' ?
          <Accordion.Item eventKey="0">
            <Accordion.Header>场景介绍</Accordion.Header>
            <Accordion.Body>
              <div>
                OCR识别，使用AI识别视频中的对象。先将视频流转成图片，然后使用AI识别图片中的对象。
                比如视频中的文字识别，人物识别，场景描述等。
                <p></p>
              </div>
            </Accordion.Body>
          </Accordion.Item> :
          <Accordion.Item eventKey="0">
            <Accordion.Header>Introduction</Accordion.Header>
            <Accordion.Body>
              <div>
                OCR recognition involves using AI to identify objects in videos. First, convert the video
                stream into images, then use AI to recognize objects in those images. For example, recognizing
                text, identifying people, and describing scenes in the video.
                <p></p>
              </div>
            </Accordion.Body>
          </Accordion.Item>
        }
      </React.Fragment>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{t('transcript.service')}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Card>
              <Card.Header>
                <Nav variant="tabs" defaultActiveKey="#provider">
                  <Nav.Item>
                    <Nav.Link href="#provider" onClick={(e) => changeConfigItem(e, 'provider')}>{t('lr.room.provider')}</Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link href="#chat" onClick={(e) => changeConfigItem(e, 'chat')}>{t('lr.room.chat')}</Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link href="#callback" onClick={(e) => changeConfigItem(e, 'callback')}>{t('ocr.callback')}</Nav.Link>
                  </Nav.Item>
                </Nav>
              </Card.Header>
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
                  baseURL: aiBaseURL, setBaseURL: setAiBaseURL, secretKey: aiSecretKey, setSecretKey: setAiSecretKey,
                  organization: aiOrganization, setOrganization: setAiOrganization,
                }} />
              </Card.Body>}
              {configItem === 'chat' && <Card.Body>
                <Form.Group className="mb-3">
                  <Form.Group className="mb-3" controlId="formAiChatEnabledCheckbox">
                    <Form.Check type="checkbox" label={t('lr.room.chate')} defaultChecked={aiChatEnabled} onClick={() => setAiChatEnabled(!aiChatEnabled)} />
                  </Form.Group>
                </Form.Group>
                <Form.Group className="mb-3">
                  <Form.Label>{t('lr.room.model')}</Form.Label>
                  <Form.Text> * {t('ocr.model')}</Form.Text>
                  <Form.Control as="input" type='input' defaultValue={aiChatModel} onChange={(e) => setAiChatModel(e.target.value)} />
                </Form.Group>
                <Form.Group className="mb-3">
                  <Form.Label>{t('lr.room.prompt')}</Form.Label>
                  <Form.Text> * {t('lr.room.prompt2')}</Form.Text>
                  <Form.Control as="textarea" type='text' rows={7}  defaultValue={aiChatPrompt} onChange={(e) => setAiChatPrompt(e.target.value)} />
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
              {configItem === 'callback' && <Card.Body>
                <Button variant='link' onClick={()=>{
                  const search = new URLSearchParams(location.search);
                  search.set('tab', 'callback');
                  navigate({pathname: `/${Locale.current()}/routers-settings`, search: search.toString()});
                }}>
                  {t('ocr.link')}
                </Button> &nbsp;
                <Form.Text> * {t('ocr.link2')}</Form.Text>
              </Card.Body>}
            </Card>
            <p></p>
            <Button variant="primary" type="submit" onClick={(e) => {
              e.preventDefault();
              updateOcrService(!ocrEnabled, () => {
                setOcrEnabled(!ocrEnabled);
              });
            }}>
              {!ocrEnabled ? t('ocr.start') : t('ocr.stop')}
            </Button> &nbsp;
            {!ocrEnabled && <React.Fragment>
              <Button ariant="primary" type="submit" disabled={operating} onClick={(e) => {
                e.preventDefault();
                resetTask();
              }}>
                {t('ocr.reset')}
              </Button>  &nbsp;
              {operating && <Spinner animation="border" variant="success" style={{verticalAlign: 'middle'}} />}
            </React.Fragment>}
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>{t('transcript.live')}</Accordion.Header>
        <Accordion.Body>
          {liveQueue?.segments?.length ? (
            <Table striped bordered hover>
              <thead>
              <tr>
                <th>#</th>
                <th>Seq</th>
                <th>URL</th>
                <th>Duration</th>
                <th>Size</th>
              </tr>
              </thead>
              <tbody>
              {liveQueue?.segments?.map((segment, index) => {
                return <tr key={segment.tsid}>
                  <td>{segment.tsid}</td>
                  <td>{segment.seqno}</td>
                  <td>{segment.url}</td>
                  <td>{`${segment.duration.toFixed(1)}`}s</td>
                  <td>{`${segment.size.toFixed(1)}`}MB</td>
                </tr>;
              })}
              </tbody>
            </Table>
          ) : t('transcript.nolive')}
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>{t('ocr.ocr')}</Accordion.Header>
        <Accordion.Body>
          {ocrQueue?.segments?.length ? (
            <Table striped bordered hover>
              <thead>
              <tr>
                <th>#</th>
                <th>URL</th>
                <th>Duration</th>
                <th title={t('ocr.eic')}>EIC</th>
                <th>Size</th>
              </tr>
              </thead>
              <tbody>
              {ocrQueue?.segments?.map((segment, index) => {
                return <tr key={segment.tsid}>
                  <td>{segment.seqno}</td>
                  <td>{segment.url}</td>
                  <td>{`${segment.duration.toFixed(1)}`}s</td>
                  <td>{`${segment.eic.toFixed(1)}`}ms</td>
                  <td>{`${segment.size.toFixed(1)}`}KB</td>
                </tr>;
              })}
              </tbody>
            </Table>
          ) : t('ocr.noocr')}
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="4">
        <Accordion.Header>{t('ocr.cbq')}</Accordion.Header>
        <Accordion.Body>
          {callbackQueue?.segments?.length ? (
            <Table striped bordered hover>
              <thead>
              <tr>
                <th>#</th>
                <th>URL</th>
                <th>Duration</th>
                <th title={t('ocr.eic')}>EIC</th>
                <th title={t('ocr.ocrc')}>OCRC</th>
                <th>Size</th>
              </tr>
              </thead>
              <tbody>
              {callbackQueue?.segments?.map((segment, index) => {
                return <tr key={segment.tsid}>
                  <td>{segment.seqno}</td>
                  <td>{segment.url}</td>
                  <td>{`${segment.duration.toFixed(1)}`}s</td>
                  <td>{`${segment.eic.toFixed(1)}`}ms</td>
                  <td>{`${segment.ocrc.toFixed(1)}`}ms</td>
                  <td>{`${segment.size.toFixed(1)}`}KB</td>
                </tr>;
              })}
              </tbody>
            </Table>
          ) : t('ocr.nocb')}
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="5">
        <Accordion.Header>{t('ocr.clq')}</Accordion.Header>
        <Accordion.Body>
          {cleanupQueue?.segments?.length ? (
            <Table striped bordered hover>
              <thead>
              <tr>
                <th>#</th>
                <th>URL</th>
                <th>Duration</th>
                <th title={t('ocr.eic')}>EIC</th>
                <th title={t('ocr.ocrc')}>OCRC</th>
                <th title={t('ocr.cbc')}>CBC</th>
                <th>Size</th>
              </tr>
              </thead>
              <tbody>
              {cleanupQueue?.segments?.map((segment, index) => {
                return <tr key={segment.tsid}>
                  <td>{segment.seqno}</td>
                  <td>{segment.url}</td>
                  <td>{`${segment.duration.toFixed(1)}`}s</td>
                  <td>{`${segment.eic.toFixed(1)}`}ms</td>
                  <td>{`${segment.ocrc.toFixed(1)}`}ms</td>
                  <td>{`${segment.cbc.toFixed(1)}`}ms</td>
                  <td>{`${segment.size.toFixed(1)}`}KB</td>
                </tr>;
              })}
              </tbody>
            </Table>
          ) : t('ocr.nocl')}
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="6">
        <Accordion.Header>{t('ocr.last')}</Accordion.Header>
        <Accordion.Body>
          {lastObject ? <>
            <p>TS: {lastObject.tsid}</p>
            <p>
              Duration: {lastObject.duration.toFixed(1)}s,
              Cost: {Number(lastObject.eic + lastObject.ocrc + lastObject.cbc).toFixed(1)}ms,
              Size: {lastObject.size.toFixed(1)}KB
            </p>
            <p>
              <img src={`/terraform/v1/ai/ocr/image/${lastObject.tsid}.jpg`} alt='Preview the last object for OCR' />
            </p>
            <p>
              <strong>AI Model</strong>: {aiChatModel} <br/>
              <strong>AI Prompt</strong>: {aiChatPrompt} <br/>
              <strong>Result:</strong>
              <pre className='ai-ocr-result'>
                {lastObject.ocr}
              </pre>
            </p>
          </> : t('ocr.nolast')}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}
