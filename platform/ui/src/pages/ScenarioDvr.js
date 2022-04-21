import React from "react";
import {Accordion, Form, Button, Table} from "react-bootstrap";
import {Token, StreamURL, Clipboard} from "../utils";
import axios from "axios";
import SetupCamSecret from '../components/SetupCamSecret';
import moment from "moment";
import {TutorialsButton, useTutorials} from "../components/TutorialsButton";
import useDvrVodStatus from "../components/DvrVodStatus";
import {useErrorHandler} from "react-error-boundary";
import {useTranslation} from "react-i18next";
import {useSrsLanguage} from "../components/LanguageSwitch";

export default function ScenarioDvr() {
  const language = useSrsLanguage();
  return language === 'zh' ? <ScenarioDvrCn /> : <ScenarioDvrEn />;
}

function ScenarioDvrCn() {
  const [dvrStatus, vodStatus] = useDvrVodStatus();
  const [activeKey, setActiveKey] = React.useState();

  // We must init the activeKey, because the defaultActiveKey only apply when init for Accordion.
  // See https://stackoverflow.com/q/61324259/17679565
  React.useEffect(() => {
    if (!vodStatus || !dvrStatus) return;

    if (vodStatus.all) {
      setActiveKey('2');
      return;
    }

    if (dvrStatus.secret) {
      if (dvrStatus.all) {
        setActiveKey('3');
      } else {
        setActiveKey('2');
      }
    } else {
      setActiveKey('1');
    }
  }, [dvrStatus, vodStatus]);

  return (
    <>
      { activeKey && <ScenarioDvrImpl activeKey={activeKey} defaultApplyAll={dvrStatus.all} enabled={!vodStatus?.all || dvrStatus.all} /> }
    </>
  );
}

function ScenarioDvrImpl({activeKey, defaultApplyAll, enabled}) {
  const [dvrAll, setDvrAll] = React.useState(defaultApplyAll);
  const [dvrFiles, setDvrFiles] = React.useState();
  const handleError = useErrorHandler();
  const {t} = useTranslation();

  const dvrTutorials = useTutorials({
    bilibili: React.useRef([
      {author: '唐为', id: 'BV14S4y1k7gr'},
    ])
  });

  React.useEffect(() => {
    const refreshDvrFiles = () => {
      const token = Token.load();
      axios.post('/terraform/v1/hooks/dvr/files', {
        ...token,
      }).then(res => {
        console.log(`DVR: Files ok, ${JSON.stringify(res.data.data)}`);
        setDvrFiles(res.data.data.map(file => {
          if (file.progress) {
            const l = window.location;
            const schema = l.protocol.replace(':', '');
            const httpPort = l.port || (l.protocol === 'http:' ? 80 : 443);
            file.location = `${l.protocol}//${l.host}/terraform/v1/hooks/dvr/hls/${file.uuid}.m3u8`;
            file.preview = `/players/srs_player.html?schema=${schema}&port=${httpPort}&autostart=true&app=terraform/v1/hooks/dvr/hls&stream=${file.uuid}.m3u8`;
          } else {
            const host = `${file.bucket}.cos.${file.region}.myqcloud.com`;
            file.location = `https://${host}/${file.uuid}/index.m3u8`;
            file.preview = `/players/srs_player.html?schema=https&port=443&autostart=true&vhost=${host}&server=${host}&app=${file.uuid}&stream=index.m3u8`;
          }

          return {
            ...file,
            url: StreamURL.build(file.vhost, file.app, file.stream),
            update: moment(file.update),
            duration: Number(file.duration),
            size: Number(file.size / 1024.0 / 1024),
          };
        }).sort((a, b) => {
          return b.update - a.update;
        }).map((file, i) => {
          return {...file, i: i + 1};
        }));
      }).catch(handleError);
    };

    refreshDvrFiles();
    const timer = setInterval(() => refreshDvrFiles(), 10 * 1000);
    return () => clearInterval(timer);
  }, [handleError]);

  const setupDvrPattern = (e) => {
    e.preventDefault();

    if (!enabled) return;

    const token = Token.load();
    axios.post('/terraform/v1/hooks/dvr/apply', {
      ...token, all: !!dvrAll,
    }).then(res => {
      alert('设置录制规则成功');
      console.log(`DVR: Apply patterns ok, all=${dvrAll}`);
    }).catch(handleError);
  };

  const copyToClipboard = React.useCallback((e, text) => {
    e.preventDefault();

    Clipboard.copy(text).then(() => {
      alert(t('helper.copyOk'));
    }).catch((err) => {
      alert(`${t('helper.copyFail')} ${err}`);
    });
  }, [t]);

  return (
    <Accordion defaultActiveKey={activeKey}>
      <Accordion.Item eventKey="0">
        <Accordion.Header>场景介绍</Accordion.Header>
        <Accordion.Body>
          <div>
            云录制<TutorialsButton prefixLine={true} tutorials={dvrTutorials} />，指录制视频流到云存储，只要推送到服务器的流都可以录制。
            <p></p>
          </div>
          <p>可应用的具体场景包括：</p>
          <ul>
            <li>直播转点播，录制直播流成为一个HLS文件，存储在云存储上，可以下载</li>
          </ul>
          <p>使用说明：</p>
          <ul>
            <li>云存储无法使用开源方案搭建，依赖公有云的云存储（<a href='https://buy.cloud.tencent.com/price/cos/calculator' target='_blank' rel='noreferrer'>计费</a>），设置密钥后将自动开通<a href='https://console.cloud.tencent.com/cos' target='_blank' rel='noreferrer'>腾讯云COS</a>云存储服务</li>
            <li>第一次使用，需要先设置云存储的访问密钥，我们会自动创建<code>srs-lighthouse</code>开头的存储桶</li>
            <li>具体使用步骤，请根据下面引导操作</li>
          </ul>
          <p>和云录制差别：</p>
          <ul>
            <li>云录制是无限磁盘，生成的是一个HLS文件，适用于直接下载和预览HLS的用户</li>
            <li>云点播是点播媒资系统，会产生一个多格式的点播文件，适用于短视频的用户，比如下载HLS和MP4、多码率、极致高清、内容处理、对视频号、分类和搜索等能力</li>
          </ul>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>设置云密钥</Accordion.Header>
        <Accordion.Body>
          <SetupCamSecret>
            <TutorialsButton prefixLine={true} tutorials={dvrTutorials} />
          </SetupCamSecret>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>设置录制规则</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3" controlId="formDvrAllCheckbox">
              <Form.Check type="checkbox" label="录制所有流" disabled={!enabled} defaultChecked={dvrAll} onClick={() => setDvrAll(!dvrAll)} />
            </Form.Group>
            <Button variant="primary" type="submit" disabled={!enabled} onClick={(e) => setupDvrPattern(e)}>
              提交
            </Button> &nbsp;
            <TutorialsButton prefixLine={true} tutorials={dvrTutorials} /> &nbsp;
            {!enabled && <Form.Text> * 若需要开启云录制，请关闭云点播(<font color='red'>云点播 / 设置点播规则 / 取消录制流</font>)</Form.Text>}
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>录制任务列表</Accordion.Header>
        <Accordion.Body>
          {
            dvrFiles?.length ? (
              <Table striped bordered hover>
                <thead>
                <tr>
                  <th>#</th>
                  <th>状态</th>
                  <th>更新时间</th>
                  <th>媒体流</th>
                  <th>时长</th>
                  <th>大小</th>
                  <th>切片</th>
                  <th>COS Bucket</th>
                  <th>地址</th>
                  <th>预览</th>
                </tr>
                </thead>
                <tbody>
                {
                  dvrFiles?.map(file => {
                    return <tr key={file.uuid}>
                      <td>{file.i}</td>
                      <td title='若300秒没有流，会自动完成录制'>{file.progress ? '录制中' : '已完成'}</td>
                      <td>{`${file.update.format('YYYY-MM-DD HH:mm:ss')}`}</td>
                      <td>{file.url}</td>
                      <td>{`${file.duration.toFixed(1)}`}秒</td>
                      <td>{`${file.size.toFixed(1)}`}MB</td>
                      <td>{file.nn}</td>
                      <td><a href={`https://console.cloud.tencent.com/cos/bucket?bucket=${file.bucket}&region=${file.region}&path=%252F${file.uuid}%252F`} target='_blank' rel='noreferrer'>{file.uuid.slice(0, 13)}</a></td>
                      <td><a href={file.location} onClick={(e) => copyToClipboard(e, file.location)} target='_blank' rel='noreferrer'>复制</a></td>
                      <td><a href={file.preview} target='_blank' rel='noreferrer'>预览</a></td>
                    </tr>;
                  })
                }
                </tbody>
              </Table>
            ) : ''
          }
          {!dvrFiles?.length ? '没有流。请开启录制并推流后，等待大约60秒左右，录制列表会自动更新' : ''}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

function ScenarioDvrEn() {
  return (
    <span>On the way...</span>
  );
}

