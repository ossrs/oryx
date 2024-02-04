import React from "react";
import {useSrsLanguage} from "../components/LanguageSwitch";
import {Accordion, Badge, Button, Form, Stack} from "react-bootstrap";
import {useTranslation} from "react-i18next";
import {useErrorHandler} from "react-error-boundary";
import {Token} from "../utils";
import axios from "axios";
import {buildUrls} from "../components/UrlGenerator";
import moment from "moment/moment";
import {SrsEnvContext} from "../components/SrsEnvContext";

export default function ScenarioTranscode(props) {
  const handleError = useErrorHandler();
  const [config, setConfig] = React.useState();
  const [activeKey, setActiveKey] = React.useState();

  React.useEffect(() => {
    axios.post('/terraform/v1/ffmpeg/transcode/query', {
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      const data = res.data.data;

      setConfig(data);
      if (data.all) {
        setActiveKey('2');
      } else {
        setActiveKey('1');
      }
      console.log(`Transcode: Query ok, ${JSON.stringify(data)}`);
    }).catch(handleError);
  }, [handleError, setActiveKey]);

  if (!props?.urls?.rtmpServer || !activeKey) return <></>;
  return <ScenarioTranscodeImpl {...props} {...{
    activeKey, defaultEnabled: config?.all, defaultConf: config
  }}/>;
}

function ScenarioTranscodeImpl({activeKey, urls, defaultEnabled, defaultConf}) {
  const language = useSrsLanguage();
  const {t} = useTranslation();
  const handleError = useErrorHandler();

  const [transcodeEnabled, setTranscodeEnabled] = React.useState(defaultEnabled);
  const [vbitrate, setVbitrate] = React.useState(defaultConf.vbitrate || 1200);
  const [abitrate, setAbitrate] = React.useState(defaultConf.abitrate || 64);
  const [vcodec, setVcodec] = React.useState(defaultConf.vcodec || 'libx264');
  const [vprofile, setVprofile] = React.useState(defaultConf.vprofile || 'baseline');
  const [vpreset, setVpreset] = React.useState(defaultConf.vpreset || 'faster');
  const [acodec, setAcodec] = React.useState(defaultConf.acodec || 'aac');
  const [achannels, setAchannels] = React.useState(defaultConf.achannels || 0);
  const [server, setServer] = React.useState(defaultConf.server || urls.rtmpServer);
  const [secret, setSecret] = React.useState((defaultConf.server || defaultConf.secret) ? defaultConf.secret : urls.transcodeStreamKey);

  const [task, setTask] = React.useState();
  const [taskInputUrls, setTaskInputUrls] = React.useState();
  const [taskOutputUrls, setTaskOutputUrls] = React.useState();
  const env = React.useContext(SrsEnvContext)[0];

  React.useEffect(() => {
    const refreshTask = () => {
      axios.post('/terraform/v1/ffmpeg/transcode/task', {
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        const task = res.data.data;

        // Mask the secret.
        if (task?.input?.indexOf(urls?.secret?.publish) > 0) {
          task.input_ = task.input;
          task.input = task.input.replaceAll(urls.secret.publish, '******');
        }
        if (task?.output?.indexOf(urls?.secret?.publish) > 0) {
          task.output_ = task.output;
          task.output = task.output.replaceAll(urls.secret.publish, '******');
        }

        // Convert to moment.
        if (task?.frame) {
          task.frame.update = task?.frame?.update ? moment(task.frame.update) : null;
        }

        setTask(task);
        if (task?.input) setTaskInputUrls(buildUrls(task.input, urls.secret, env));
        if (task?.output) setTaskOutputUrls(buildUrls(task.output, urls.secret, env));
        console.log(`Transcode: Query task ${JSON.stringify(task)}`);
      }).catch(handleError);
    };

    refreshTask();
    const timer = setInterval(() => refreshTask(), 10 * 1000);
    return () => clearInterval(timer);
  }, [handleError, urls, setTask, setTaskInputUrls, setTaskOutputUrls, env]);

  const updateTranscodeStatus = React.useCallback((enabled, success) => {
    if (!vbitrate || vbitrate < 100 || vbitrate > 100*1000) return alert(`Invalid vbitrate ${vbitrate}, should be in [100, 100000] Kbps`);
    if (!abitrate || abitrate < 10 || abitrate > 1000) return alert(`Invalid abitrate ${abitrate}, should be in [10, 1000] Kbps`);
    if (!vcodec || vcodec !== 'libx264') return alert(`Invalid vcodec ${vcodec}, should be libx264`);
    if (!vprofile || !['baseline', 'main', 'high'].includes(vprofile)) return alert(`Invalid vprofile ${vprofile}, should be in [baseline, main, high]`);
    if (!vpreset || !['ultrafast', 'superfast', 'veryfast', 'faster', 'fast', 'medium', 'slow'].includes(vpreset)) return alert(`Invalid vpreset ${vpreset}, should be in [ultrafast, superfast, veryfast, faster, fast, medium, slow]`);
    if (!acodec || acodec !== 'aac') return alert(`Invalid acodec ${acodec}, should be aac`);
    if (achannels === undefined || achannels === null || achannels === '') return alert(`Invalid achannels ${achannels}, should not empty`);
    if (![0, 1, 2].includes(achannels)) return alert(`Invalid achannels ${achannels}, should be in [0, 1, 2]`);
    if (!server && !secret) return alert(`Invalid server ${server} and key ${secret}`);

    axios.post('/terraform/v1/ffmpeg/transcode/apply', {
      all: enabled, vcodec, acodec, vbitrate, abitrate, achannels: achannels, vprofile, vpreset,
      server, secret,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      alert(t('helper.setOk'));
      console.log(`Transcode: Apply patterns ok, all=${enabled}, vbitrate=${vbitrate}, abitrate=${abitrate}, vcodec=${vcodec}, vprofile=${vprofile}, vpreset=${vpreset}, acodec=${acodec}, server=${server}, secret=${secret}`);
      success && success();
    }).catch(handleError);
  }, [handleError, t, vbitrate, abitrate, achannels, vcodec, vprofile, vpreset, acodec, server, secret]);

  return (
    <Accordion defaultActiveKey={[activeKey]} alwaysOpen>
      <React.Fragment>
        {language === 'zh' ?
          <Accordion.Item eventKey="0">
            <Accordion.Header>场景介绍</Accordion.Header>
            <Accordion.Body>
              <div>
                直播转码，是将SRS Stack的直播流，用FFmpeg转成不同码率和清晰度的直播流，再推送到SRS Stack。
                <p></p>
              </div>
              <p>可应用的具体场景包括：</p>
              <ul>
                <li>降低带宽，保持相同清晰度，降低流的码率，从而降低整体的观看带宽</li>
                <li>提高清晰度，保持同样输出码率和带宽，提高原始码流的码率，提升输出流的清晰度和画质</li>
                <li>多清晰度多码率，从高分辨率的原始码流，转成多个不同清晰度和码率的流，给不同设备观看</li>
              </ul>
              <p>使用说明：</p>
              <ul>
                <li>首先先将需要转码的原始流，推送到SRS Stack</li>
                <li>然后配置转码的信息，选择不同的转码模板，调整转码参数</li>
                <li>开始转码后，将生成新的不同的转码的流，推送到SRS Stack上</li>
              </ul>
            </Accordion.Body>
          </Accordion.Item> :
          <Accordion.Item eventKey="0">
            <Accordion.Header>Introduction</Accordion.Header>
            <Accordion.Body>
              <div>
                Live streaming transcoding is the process of converting the live stream from SRS Stack using FFmpeg into different bitrates and resolutions, and then pushing it back to SRS Stack.
                <p></p>
              </div>
              <p>Specific scenarios where this can be applied include:</p>
              <ul>
                <li>Reducing bandwidth while maintaining the same resolution by lowering the bitrate of the stream, thus reducing the overall viewing bandwidth.</li>
                <li>Improving resolution while keeping the same output bitrate and bandwidth by increasing the bitrate of the original stream, enhancing the output stream's clarity and quality.</li>
                <li>Creating multiple resolutions and bitrates from a high-resolution original stream for viewing on different devices.</li>
              </ul>
              <p>Instructions for use:</p>
              <ul>
                <li>First, push the original stream that needs to be transcoded to SRS Stack.</li>
                <li>Then, configure the transcoding information, choose different transcoding templates, and adjust the transcoding parameters.</li>
                <li>After starting the transcoding, new streams with different transcoded resolutions and bitrates will be generated and pushed to SRS Stack.</li>
              </ul>
            </Accordion.Body>
          </Accordion.Item>
        }
      </React.Fragment>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{t('transcode.config.header')}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>{t('transcode.config.vbitrate')}</Form.Label>
              <Form.Text> * {t('transcode.config.vbitrate2')}</Form.Text>
              <Form.Control as="input" defaultValue={vbitrate} onChange={(e) => setVbitrate(parseInt(e.target.value))} />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{t('transcode.config.abitrate')}</Form.Label>
              <Form.Text> * {t('transcode.config.abitrate2')}</Form.Text>
              <Form.Control as="input" defaultValue={abitrate} onChange={(e) => setAbitrate(parseInt(e.target.value))} />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{t('transcode.config.vcodec')}</Form.Label>
              <Form.Text> * {t('transcode.config.vcodec2')}</Form.Text>
              <Form.Select defaultValue={vcodec} onChange={(e) => setVcodec(e.target.value)}>
                <option value="">--{t('helper.noSelect')}--</option>
                <option value="libx264">{t('transcode.config.vcodec3')}</option>
              </Form.Select>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{t('transcode.config.vprofile')}</Form.Label>
              <Form.Text> * {t('transcode.config.vprofile2')}</Form.Text>
              <Form.Select defaultValue={vprofile} onChange={(e) => setVprofile(e.target.value)}>
                <option value="">--{t('helper.noSelect')}--</option>
                <option value="baseline">{t('transcode.config.vprofile3')}</option>
                <option value="main">{t('transcode.config.vprofile4')}</option>
                <option value="high">{t('transcode.config.vprofile5')}</option>
              </Form.Select>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{t('transcode.config.vpreset')}</Form.Label>
              <Form.Text> * {t('transcode.config.vpreset2')}</Form.Text>
              <Form.Select defaultValue={vpreset} onChange={(e) => setVpreset(e.target.value)}>
                <option value="">--{t('helper.noSelect')}--</option>
                <option value="ultrafast">Ultrafast({t('transcode.config.vpreset3')})</option>
                <option value="superfast">Superfast</option>
                <option value="veryfast">Veryfast</option>
                <option value="faster">Faster({t('transcode.config.vpreset4')})</option>
                <option value="fast">Fast</option>
                <option value="medium">Medium</option>
                <option value="slow">Slow({t('transcode.config.vpreset5')})</option>
              </Form.Select>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{t('transcode.config.acodec')}</Form.Label>
              <Form.Text> * {t('transcode.config.acodec2')}</Form.Text>
              <Form.Select defaultValue={acodec} onChange={(e) => setAcodec(e.target.value)}>
                <option value="">--{t('helper.noSelect')}--</option>
                <option value="aac">{t('transcode.config.acodec3')}</option>
              </Form.Select>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{t('transcode.config.achannel')}</Form.Label>
              <Form.Text> * {t('transcode.config.achannelN')}</Form.Text>
              <Form.Select defaultValue={achannels} onChange={(e) => setAchannels(parseInt(e.target.value))}>
                <option value="">--{t('helper.noSelect')}--</option>
                <option value="0">{t('transcode.config.achannel0')}</option>
                <option value="1">{t('transcode.config.achannel1')}</option>
                <option value="2">{t('transcode.config.achannel2')}</option>
              </Form.Select>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{t('transcode.config.server')}</Form.Label>
              <Form.Text> * {t('transcode.config.server2')}</Form.Text>
              <Form.Control as="input" defaultValue={server} onChange={(e) => setServer(e.target.value)} />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{t('transcode.config.key')}</Form.Label>
              <Form.Text> * {t('transcode.config.key2')}</Form.Text>
              <Form.Control as="input" defaultValue={secret} onChange={(e) => setSecret(e.target.value)} />
            </Form.Group>
            <Button ariant="primary" type="submit" onClick={(e) => {
              e.preventDefault();
              updateTranscodeStatus(!transcodeEnabled, () => {
                setTranscodeEnabled(!transcodeEnabled);
              });
            }}>
              {!transcodeEnabled ? t('transcode.status.start') : t('transcode.status.stop')}
            </Button>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>{t('transcode.status.header')}</Accordion.Header>
        <Accordion.Body>
          <Stack gap={1}>
            <div>
              <Badge bg={task?.enabled ? (task?.frame ? 'success' : 'primary') : 'secondary'}>
                {task?.enabled ? (task?.frame ? t('transcode.status.transcoding') : t('transcode.status.waiting')) : t('transcode.status.inactive')}
              </Badge>
            </div>
            <div>
              Update: {task?.frame?.update?.format('YYYY-MM-DD HH:mm:ss')}
            </div>
            <div>
              Log: {task?.frame?.log}
            </div>
            <div>
              Input: {task?.input} &nbsp;
              <a href={taskInputUrls?.flvPlayer} target='_blank' rel='noreferrer'>{t('transcode.status.preview')}</a>
            </div>
            <div>
              Output: {task?.output} &nbsp;
              <a href={taskOutputUrls?.flvPlayer} target='_blank' rel='noreferrer'>{t('transcode.status.preview')}</a>
            </div>
          </Stack>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

