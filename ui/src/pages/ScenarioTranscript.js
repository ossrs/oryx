import React from "react";
import {Accordion, Button, Form, Nav, Spinner, Table, Card} from "react-bootstrap";
import {useSrsLanguage} from "../components/LanguageSwitch";
import {useTranslation} from "react-i18next";
import {Token} from "../utils";
import axios from "axios";
import {useErrorHandler} from "react-error-boundary";
import PopoverConfirm from "../components/PopoverConfirm";
import {OpenAISecretSettings} from "../components/OpenAISettings";

export default function ScenarioTranscript(props) {
  const handleError = useErrorHandler();
  const [config, setConfig] = React.useState();
  const [uuid, setUuid] = React.useState();
  const [activeKey, setActiveKey] = React.useState();

  React.useEffect(() => {
    axios.post('/terraform/v1/ai/transcript/query', {
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

      console.log(`Transcript: Query ok, ${JSON.stringify(data)}`);
    }).catch(handleError);
  }, [handleError, setActiveKey, setConfig,  setUuid]);

  if (!activeKey) return <></>;
  return <ScenarioTranscriptImpl {...props} {...{
    activeKey, defaultEnabled: config?.all, defaultConf: config, defaultUuid: uuid,
  }}/>;
}

function ScenarioTranscriptImpl({activeKey, defaultEnabled, defaultConf, defaultUuid}) {
  const language = useSrsLanguage();
  const {t} = useTranslation();
  const handleError = useErrorHandler();

  const [operating, setOperating] = React.useState(false);
  const [refreshNow, setRefreshNow] = React.useState();
  const [transcriptEnabled, setTranscriptEnabled] = React.useState(defaultEnabled);
  const [secretKey, setSecretKey] = React.useState(defaultConf.secretKey);
  const [organization, setOrganization] = React.useState(defaultConf.organization);
  const [baseURL, setBaseURL] = React.useState(defaultConf.baseURL || (language === 'zh' ? '' : 'https://api.openai.com/v1'));
  const [targetLanguage, setTargetLanguage] = React.useState(defaultConf.lang || language);
  const [forceStyle, setForceStyle] = React.useState(defaultConf.forceStyle || 'Alignment=2,MarginV=20');
  const [videoCodecParams, setVideoCodecParams] = React.useState(defaultConf.videoCodecParams || '-c:v libx264 -profile:v main -preset:v medium -tune zerolatency -bf 0');
  const [overlayEnabled, setOverlayEnabled] = React.useState(defaultConf.overlayEnabled);
  const [webvttEnabled, setWebvttEnabled] = React.useState(defaultConf.webvttEnabled);

  const [liveQueue, setLiveQueue] = React.useState();
  const [asrQueue, setAsrQueue] = React.useState();
  const [fixQueue, setFixQueue] = React.useState();
  const [overlayQueue, setOverlayQueue] = React.useState();

  const [uuid, setUuid] = React.useState(defaultUuid);
  const [overlayHlsUrl, setOverlayHlsUrl] = React.useState();
  const [overlayHlsPreview, setOverlayHlsPreview] = React.useState();
  const [webvttHlsUrl, setWebvttHlsUrl] = React.useState();
  const [webvttHlsPreview, setWebvttHlsPreview] = React.useState();
  const [originalHlsUrl, setOriginalHlsUrl] = React.useState();
  const [originalHlsPreview, setOriginalHlsPreview] = React.useState();

  const [configItem, setConfigItem] = React.useState('provider');

  React.useEffect(() => {
    if (secretKey) return;

    axios.post('/terraform/v1/mgmt/openai/query', null, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      const data = res.data.data;
      setSecretKey(data.aiSecretKey);
      setBaseURL(data.aiBaseURL);
      setOrganization(data.aiOrganization);
      console.log(`Transcript: Query open ai ok, data=${JSON.stringify(data)}`);
    }).catch(handleError);
  }, [handleError, secretKey, setSecretKey, setBaseURL, setOrganization]);

  const changeConfigItem = React.useCallback((e, t) => {
    e.preventDefault();
    setConfigItem(t);
  }, [setConfigItem]);

  React.useEffect(() => {
    const l = window.location;
    const schema = l.protocol.replace(':', '');
    const httpPort = l.port || (l.protocol === 'http:' ? 80 : 443);

    setOverlayHlsUrl(`${l.protocol}//${l.host}/terraform/v1/ai/transcript/hls/overlay/${uuid}.m3u8`);
    setOverlayHlsPreview(`/players/srs_player.html?schema=${schema}&port=${httpPort}&autostart=true&app=terraform/v1/ai/transcript/hls/overlay&stream=${uuid}.m3u8`);

    setWebvttHlsUrl(`${l.protocol}//${l.host}/terraform/v1/ai/transcript/hls/webvtt/${uuid}/index.m3u8`);
    setWebvttHlsPreview(`/players/srs_player.html?schema=${schema}&port=${httpPort}&autostart=true&app=terraform/v1/ai/transcript/hls/webvtt/${uuid}&stream=index.m3u8`);

    setOriginalHlsUrl(`${l.protocol}//${l.host}/terraform/v1/ai/transcript/hls/original/${uuid}.m3u8`);
    setOriginalHlsPreview(`/players/srs_player.html?schema=${schema}&port=${httpPort}&autostart=true&app=terraform/v1/ai/transcript/hls/original&stream=${uuid}.m3u8`);
  }, [uuid, setOverlayHlsUrl, setOverlayHlsPreview, setWebvttHlsUrl, setWebvttHlsPreview, setOriginalHlsUrl, setOriginalHlsPreview]);

  const updateAiService = React.useCallback((enabled, success) => {
    if (!secretKey) return alert(`Invalid secret key ${secretKey}`);
    if (!baseURL) return alert(`Invalid base url ${baseURL}`);

    axios.post('/terraform/v1/ai/transcript/apply', {
      uuid, all: !!enabled, secretKey, organization, baseURL, lang: targetLanguage,
      overlayEnabled: !!overlayEnabled, forceStyle, videoCodecParams,
      webvttEnabled: !!webvttEnabled,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      alert(t('helper.setOk'));
      console.log(`Transcript: Apply config ok, uuid=${uuid}.`);
      success && success();
    }).catch(handleError);
  }, [t, handleError, secretKey, baseURL, targetLanguage, overlayEnabled, forceStyle, videoCodecParams, webvttEnabled, uuid, organization]);

  const resetTask = React.useCallback(() => {
    setOperating(true);

    axios.post('/terraform/v1/ai/transcript/reset', {
      uuid,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      alert(t('helper.setOk'));
      const data = res.data.data;
      setUuid(data.uuid);
      console.log(`Transcript: Reset task ${uuid} ok: ${JSON.stringify(data)}`);
    }).catch(handleError).finally(setOperating);
  }, [t, handleError, uuid, setUuid, setOperating]);

  const clearText = React.useCallback((segment) => {
    setOperating(true);

    axios.post('/terraform/v1/ai/transcript/clear-subtitle', {
      uuid, tsid: segment.tsid,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      alert(t('helper.setOk'));
      const data = res.data.data;
      setRefreshNow(!refreshNow);
      console.log(`Transcript: Clear subtitle of task ${uuid} segment ${segment.tsid} ok: ${JSON.stringify(data)}`);
    }).catch(handleError).finally(setOperating);
  }, [t, handleError, setOperating, uuid, refreshNow, setRefreshNow]);

  React.useEffect(() => {
    const refreshLiveQueueTask = () => {
      axios.post('/terraform/v1/ai/transcript/live-queue', {
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
        console.log(`Transcript: Query live queue ${JSON.stringify(queue)}`);
      }).catch(handleError);
    };

    refreshLiveQueueTask();
    const timer = setInterval(() => refreshLiveQueueTask(), 3 * 1000);
    return () => clearInterval(timer);
  }, [handleError, setLiveQueue]);

  React.useEffect(() => {
    const refreshAsrQueueTask = () => {
      axios.post('/terraform/v1/ai/transcript/asr-queue', {
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        const queue = res.data.data;
        queue.segments = queue?.segments?.map(segment => {
          return {
            ...segment,
            duration: Number(segment.duration),
            size: Number(segment.size / 1024.0),
            eac: Number(segment.eac),
          };
        });
        setAsrQueue(queue);
        console.log(`Transcript: Query asr queue ${JSON.stringify(queue)}`);
      }).catch(handleError);
    };

    refreshAsrQueueTask();
    const timer = setInterval(() => refreshAsrQueueTask(), 3 * 1000);
    return () => clearInterval(timer);
  }, [handleError, setAsrQueue]);

  React.useEffect(() => {
    const refreshFixQueueTask = () => {
      axios.post('/terraform/v1/ai/transcript/fix-queue', {
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        const queue = res.data.data;
        queue.segments = queue?.segments?.map((segment, index) => {
          return {
            ...segment,
            duration: Number(segment.duration),
            size: Number(segment.size / 1024.0),
            eac: Number(segment.eac),
            asrc: Number(segment.asrc),
            // The max length of the subtitle text in asrs array.
            asrsMaxLength: Math.max(...segment.asrs.map(asr => asr.text.length)),
            asrsMaxWords: Math.max(...segment.asrs.map(asr => asr.text.split(' ').length)),
            // Rules:
            // 1. Always allow to clear the first segment, that is only one segment in the queue.
            // 2. Prevent the first segment from clearing subtitles, as it may have already been added
            //    to the overlay queue and not be able to modify it.
            // 3. If already cleared, the uca(User Clear ASR) is set to true.
            allowClearSubtitle: (queue.segments.length <= 1 || index !== 0) && !segment.uca,
          };
        });
        setFixQueue(queue);
        console.log(`Transcript: Query fix queue ${JSON.stringify(queue)}`);
      }).catch(handleError);
    };

    refreshFixQueueTask();
    const timer = setInterval(() => refreshFixQueueTask(), 3 * 1000);
    return () => clearInterval(timer);
  }, [handleError, setFixQueue, refreshNow]);

  React.useEffect(() => {
    const refreshOverlayQueueTask = () => {
      axios.post('/terraform/v1/ai/transcript/overlay-queue', {
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        const queue = res.data.data;
        queue.segments = queue?.segments?.map(segment => {
          return {
            ...segment,
            duration: Number(segment.duration),
            size: Number(segment.size / 1024.0 / 1024.0),
            eac: Number(segment.eac),
            asrc: Number(segment.asrc),
            olc: Number(segment.olc),
            // The max length of the subtitle text in asrs array.
            asrsMaxLength: Math.max(...segment.asrs.map(asr => asr.text.length)),
            asrsMaxWords: Math.max(...segment.asrs.map(asr => asr.text.split(' ').length)),
          };
        });
        setOverlayQueue(queue);
        console.log(`Transcript: Query overlay queue ${JSON.stringify(queue)}`);
      }).catch(handleError);
    };

    refreshOverlayQueueTask();
    const timer = setInterval(() => refreshOverlayQueueTask(), 3 * 1000);
    return () => clearInterval(timer);
  }, [handleError, setOverlayQueue]);

  return (
    <Accordion defaultActiveKey={activeKey}>
      <React.Fragment>
        {language === 'zh' ?
          <Accordion.Item eventKey="0">
            <Accordion.Header>场景介绍</Accordion.Header>
            <Accordion.Body>
              <div>
                AI字幕使用人工智能将实时语音转换成文本，然后允许人工编辑和校正文本，并将其翻译成多种语言，
                并将修改后的多语言文本合并在视频流中，最终生成一个新的直播流。
                <p></p>
              </div>
              <p>可应用的具体场景包括：</p>
              <ul>
                <li>直播时，使用AI生成的自动字幕，让听力受限的观众，在听不到声音时，可以看视频的字幕。</li>
                <li>为不同语言的观众提供多语言字幕。直播时，由AI翻译成各种语言，从而生成多个流，每个流都有特定语言的字幕。
                  例如，如果直播源是英语的，那么会有带有英语、中文、法语等其他语言字幕的输出流。</li>
                <li>为多个直播平台提供一致的字幕体验。因为一些平台支持自动字幕，而其他平台则不支持。通过在源直播中加入自动字幕，
                  我们可以确保在各种直播平台上的一致性，确保所有平台都有一致的字幕。</li>
              </ul>
            </Accordion.Body>
          </Accordion.Item> :
          <Accordion.Item eventKey="0">
            <Accordion.Header>Introduction</Accordion.Header>
            <Accordion.Body>
              <div>
                Transcription uses AI to convert live speech into text, then allows you to edit and correct the
                text, translates it into multiple languages, and overlay the multilingual text onto the video,
                ultimately generating a new live stream.
                <p></p>
              </div>
              <p>Specific scenarios where this can be applied include:</p>
              <ul>
                <li>AI-generated automatic subtitles for live streams are provided for audiences with hearing
                  impairments, allowing them to read the subtitles even if they are unable to hear the speech.</li>
                <li>Multilingual subtitles are provided for audiences who speak different languages. The live
                  stream is translated by AI into various languages, resulting in multiple streams, each with
                  subtitles in a specific language. For example, if the source stream is in English, there will
                  be output streams with subtitles in English, Chinese, French, and other languages.</li>
                <li>Automatic subtitles are provided for multiple live stream platforms. This is because some
                  platforms offer automatic subtitles, while others do not. By incorporating automatic subtitles
                  into the source stream, we can ensure consistency across various live streaming platforms,
                  ensuring that all have subtitles.</li>
              </ul>
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
                    <Nav.Link href="#asr" onClick={(e) => changeConfigItem(e, 'asr')}>{t('lr.room.asr')}</Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link href="#overlay" onClick={(e) => changeConfigItem(e, 'overlay')}>{t('transcript.overlay2')}</Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link href="#webvtt" onClick={(e) => changeConfigItem(e, 'webvtt')}>{t('transcript.vtt')}</Nav.Link>
                  </Nav.Item>
                </Nav>
              </Card.Header>
              {configItem === 'provider' && <Card.Body>
                <OpenAISecretSettings {...{
                  baseURL, setBaseURL, secretKey, setSecretKey,
                  organization, setOrganization,
                }} />
              </Card.Body>}
              {configItem === 'asr' && <Card.Body>
                <Form.Group className="mb-3">
                  <Form.Label>{t('transcript.lang')}</Form.Label>
                  <Form.Text> * {t('transcript.lang2')}. &nbsp;
                    {t('helper.eg')} <code>en, zh, fr, de, ja, ru </code>, ... &nbsp;
                    {t('helper.see')} <a href='https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes' target='_blank' rel='noreferrer'>ISO-639-1</a>.
                  </Form.Text>
                  <Form.Control as="input" defaultValue={targetLanguage} onChange={(e) => setTargetLanguage(e.target.value)} />
                </Form.Group>
              </Card.Body>}
              {configItem === 'overlay' && <Card.Body>
                <Form.Group className="mb-3">
                  <Form.Group className="mb-3" controlId="formOverlayEnabledCheckbox">
                    <Form.Check type="checkbox" label={t('transcript.ole')} defaultChecked={overlayEnabled} onClick={() => setOverlayEnabled(!overlayEnabled)} />
                  </Form.Group>
                </Form.Group>
                <Form.Group className="mb-3">
                  <Form.Label>{t('transcript.fstyle')}</Form.Label>
                  <Form.Text> * {t('transcript.fstyle2')}. &nbsp;
                    {t('helper.see')} <a href={t('transcript.fstyle3')} target='_blank' rel='noreferrer'>FFmpeg: force_style</a>.
                  </Form.Text>
                  <Form.Control as="input" defaultValue={forceStyle} onChange={(e) => setForceStyle(e.target.value)} />
                </Form.Group>
                <Form.Group className="mb-3">
                  <Form.Label>{t('transcript.trans0')}</Form.Label>
                  <Form.Text> * {t('transcript.trans1')}. &nbsp;
                    {t('helper.see')} <a href={t('transcript.trans2')} target='_blank' rel='noreferrer'>FFmpeg: video codec</a>.
                  </Form.Text>
                  <Form.Control as="input" defaultValue={videoCodecParams} onChange={(e) => setVideoCodecParams(e.target.value)} />
                </Form.Group>
              </Card.Body>}
              {configItem === 'webvtt' && <Card.Body>
                <Form.Group className="mb-3">
                  <Form.Group className="mb-3" controlId="formWebvttEnabledCheckbox">
                    <Form.Check type="checkbox" label={t('transcript.vtt2')} defaultChecked={webvttEnabled} onClick={() => setWebvttEnabled(!webvttEnabled)} />
                  </Form.Group>
                </Form.Group>
              </Card.Body>}
            </Card>
            <p></p>
            <Button ariant="primary" type="submit" onClick={(e) => {
              e.preventDefault();
              updateAiService(!transcriptEnabled, () => {
                setTranscriptEnabled(!transcriptEnabled);
              });
            }}>
              {!transcriptEnabled ? t('transcript.start') : t('transcript.stop')}
            </Button> &nbsp;
            {!transcriptEnabled && <React.Fragment>
              <Button ariant="primary" type="submit" disabled={operating} onClick={(e) => {
                e.preventDefault();
                resetTask();
              }}>
                {t('transcript.reset')}
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
        <Accordion.Header>{t('transcript.asr')}</Accordion.Header>
        <Accordion.Body>
          {asrQueue?.segments?.length ? (
            <Table striped bordered hover>
              <thead>
              <tr>
                <th>#</th>
                <th>URL</th>
                <th>Duration</th>
                <th title={t('transcript.eac')}>EAC</th>
                <th>Size</th>
              </tr>
              </thead>
              <tbody>
              {asrQueue?.segments?.map((segment, index) => {
                return <tr key={segment.tsid}>
                  <td>{segment.seqno}</td>
                  <td>{segment.url}</td>
                  <td>{`${segment.duration.toFixed(1)}`}s</td>
                  <td>{`${segment.eac.toFixed(1)}`}ms</td>
                  <td>{`${segment.size.toFixed(1)}`}KB</td>
                </tr>;
              })}
              </tbody>
            </Table>
          ) : t('transcript.noasr')}
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="4">
        <Accordion.Header>{t('transcript.fix')}</Accordion.Header>
        <Accordion.Body>
          {fixQueue?.segments?.length ? (
            <Table striped bordered hover>
              <thead>
              <tr>
                <th>#</th>
                <th>URL</th>
                <th>Duration</th>
                <th title={t('transcript.eac')}>EAC</th>
                <th title={t('transcript.asrc')}>ASRC</th>
                <th>Segments</th>
                <th>Size</th>
                <th>Text</th>
                <th>{t('transcript.action')}</th>
              </tr>
              </thead>
              <tbody>
              {fixQueue?.segments?.map((segment, index) => {
                return <tr key={segment.tsid}>
                  <td>{segment.seqno}</td>
                  <td>{segment.url}</td>
                  <td>{`${segment.duration.toFixed(1)}`}s</td>
                  <td>{`${segment.eac.toFixed(1)}`}ms</td>
                  <td>{`${segment.asrc.toFixed(1)}`}ms</td>
                  <td title={`There are ${segment.asrs.length} segments, max text length is ${segment.asrsMaxLength} bytes, max words is ${segment.asrsMaxWords}`}>
                    {segment.asrs.length}/{segment.asrsMaxLength}/{segment.asrsMaxWords}
                  </td>
                  <td>{`${segment.size.toFixed(1)}`}KB</td>
                  <td style={{textDecoration: segment.uca ? "line-through" : ''}}>{segment.asr}</td>
                  <td>
                    <PopoverConfirm placement='top'
                                    trigger={ <a href={`#${segment.tsid}`} hidden={!segment.allowClearSubtitle}>{t('transcript.clear')}</a> }
                                    onClick={() => clearText(segment)}>
                      <p>
                        {t('transcript.clear2')}
                      </p>
                    </PopoverConfirm>
                  </td>
                </tr>;
              })}
              </tbody>
            </Table>
          ) : t('transcript.nofix')}
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="5">
        <Accordion.Header>{t('transcript.overlay')}</Accordion.Header>
        <Accordion.Body>
          {overlayQueue?.segments?.length ? (
            <Table striped bordered hover>
              <thead>
              <tr>
                <th>#</th>
                <th>URL</th>
                <th>Duration</th>
                <th title={t('transcript.eac')}>EAC</th>
                <th title={t('transcript.asrc')}>ASRC</th>
                <th title={t('transcript.olc')}>OLC</th>
                <th>Segments</th>
                <th>Size</th>
                <th>Text</th>
              </tr>
              </thead>
              <tbody>
              {overlayQueue?.segments?.map((segment, index) => {
                return <tr key={segment.tsid}>
                  <td>{segment.seqno}</td>
                  <td>{segment.url}</td>
                  <td>{`${segment.duration.toFixed(1)}`}s</td>
                  <td>{`${segment.eac.toFixed(1)}`}ms</td>
                  <td>{`${segment.asrc.toFixed(1)}`}ms</td>
                  <td>{`${segment.olc.toFixed(1)}`}ms</td>
                  <td>{`${segment.size.toFixed(1)}`}MB</td>
                  <td title={`There are ${segment.asrs.length} segments, max text length is ${segment.asrsMaxLength} bytes, max words is ${segment.asrsMaxWords}`}>
                    {segment.asrs.length}/{segment.asrsMaxLength}/{segment.asrsMaxWords}
                  </td>
                  <td style={{textDecoration: segment.uca ? "line-through" : ''}}>{segment.asr}</td>
                </tr>;
              })}
              </tbody>
            </Table>
          ) : t('transcript.nooverlay')}
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="6">
        <Accordion.Header>{t('transcript.ops')}</Accordion.Header>
        <Accordion.Body>
          {t('transcript.porg')}: <a href={originalHlsPreview} target='_blank' rel='noreferrer'>{originalHlsUrl}</a><br/>
          {t('transcript.pol')}: <a href={overlayHlsPreview} target='_blank' rel='noreferrer'>{overlayHlsUrl}</a><br/>
          {webvttEnabled && <>
            {t('transcript.pvtt')}: <a href={webvttHlsPreview} target='_blank' rel='noreferrer'>{webvttHlsUrl}</a><br/>
          </>}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}
