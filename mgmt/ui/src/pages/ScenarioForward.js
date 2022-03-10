import {Accordion, Badge, Button, Form, Table} from "react-bootstrap";
import React from "react";
import {Token} from "../utils";
import axios from "axios";
import moment from "moment";
import {TutorialsButton, useTutorials} from "../components/TutorialsButton";
import {useErrorHandler} from "react-error-boundary";

export default function ScenarioForward() {
  const [init, setInit] = React.useState();
  const [activeKey, setActiveKey] = React.useState();
  const [secrets, setSecrets] = React.useState();
  const handleError = useErrorHandler();

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/ffmpeg/forward/secret', {
      ...token,
    }).then(res => {
      const secrets = res.data.data;
      setInit(true);
      setSecrets(secrets || {});
      console.log(`Forward: Secret query ok ${JSON.stringify(secrets)}`);
    }).catch(handleError);
  }, [handleError]);

  React.useEffect(() => {
    if (!init || !secrets) return;

    if (secrets.wx?.enabled || secrets.bilibili?.enabled || secrets.kuaishou?.enabled) {
      return setActiveKey('99');
    }

    if (!secrets.wx?.server || !secrets.wx?.secret || !secrets.wx?.enabled) {
      setActiveKey('1');
    } else if (!secrets.bilibili?.server || !secrets.bilibili?.secret || !secrets.bilibili?.enabled) {
      setActiveKey('2');
    } else if (!secrets.kuaishou?.server || !secrets.kuaishou?.secret || !secrets.kuaishou?.enabled) {
      setActiveKey('3');
    } else {
      setActiveKey('99');
    }
  }, [init, secrets]);

  return <>
    {activeKey && <ScenarioForwardImpl defaultActiveKey={activeKey} defaultSecrets={secrets}/>}
  </>;
}

function ScenarioForwardImpl({defaultActiveKey, defaultSecrets}) {
  const [wxEnabled, setWxEnabled] = React.useState(defaultSecrets?.wx?.enabled);
  const [wxServer, setWxServer] = React.useState(defaultSecrets?.wx?.server);
  const [wxSecret, setWxSecret] = React.useState(defaultSecrets?.wx?.secret);
  const [bilibiliEnabled, setBilibiliEnabled] = React.useState(defaultSecrets?.bilibili?.enabled);
  const [bilibiliServer, setBilibiliServer] = React.useState(defaultSecrets?.bilibili?.server);
  const [bilibiliSecret, setBilibiliSecret] = React.useState(defaultSecrets?.bilibili?.secret);
  const [kuaishouEnabled, setKuaishouEnabled] = React.useState(defaultSecrets?.kuaishou?.enabled);
  const [kuaishouServer, setKuaishouServer] = React.useState(defaultSecrets?.kuaishou?.server);
  const [kuaishouSecret, setKuaishouSecret] = React.useState(defaultSecrets?.kuaishou?.secret);
  const [forwards, setForwards] = React.useState();
  const handleError = useErrorHandler();

  const forwardTutorials = useTutorials(React.useRef([
    {author: 'SRS', id: 'BV1KY411V7uc'},
  ]));

  React.useEffect(() => {
    const refreshStreams = () => {
      const token = Token.load();
      axios.post('/terraform/v1/ffmpeg/forward/streams', {
        ...token,
      }).then(res => {
        setForwards(res.data.data.map((e, i) => ({
          ...e,
          name: {wx: '视频号直播', bilibili: 'Bilibili直播间', kuaishou: '快手云直播'}[e.platform],
          update: e.frame?.update ? moment(e.frame.update) : null,
          i,
        })));
        console.log(`Forward: Query streams ${JSON.stringify(res.data.data)}`);
      }).catch(handleError);
    };

    refreshStreams();
    const timer = setInterval(() => refreshStreams(), 10 * 1000);
    return () => clearInterval(timer);
  }, [handleError]);

  const updateSecrets = (e, action, platform, server, secret, enabled) => {
    e.preventDefault();
    if (!server) return alert('请输入推流地址');

    const token = Token.load();
    axios.post('/terraform/v1/ffmpeg/forward/secret', {
      ...token, action, platform, server, secret, enabled: !!enabled,
    }).then(res => {
      alert('转推设置成功');
    }).catch(handleError);
  };

  return (
    <Accordion defaultActiveKey={defaultActiveKey}>
      <Accordion.Item eventKey="0">
        <Accordion.Header>场景介绍</Accordion.Header>
        <Accordion.Body>
          <div>
            多平台转播<TutorialsButton prefixLine={true} tutorials={forwardTutorials} />，将流转发给其他平台，比如视频号直播、快手、B站等。
            <p></p>
          </div>
          <p>可应用的具体场景包括：</p>
          <ul>
            <li>节约上行带宽，避免客户端推多路流，服务器转发更有保障</li>
          </ul>
          <p>使用说明：</p>
          <ul>
            <li>首先使用适合你的场景推流</li>
            <li>然后设置转发的平台</li>
          </ul>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>视频号直播</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>推流地址</Form.Label>
              <Form.Text> * 请先<a href='https://channels.weixin.qq.com/platform/live/liveBuild' target='_blank' rel='noreferrer'>创建直播</a>，然后获取推流地址</Form.Text>
              <Form.Control as="input" defaultValue={wxServer} onChange={(e) => setWxServer(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>推流密钥</Form.Label>
              <Form.Text> * 请先<a href='https://channels.weixin.qq.com/platform/live/liveBuild' target='_blank' rel='noreferrer'>创建直播</a>，然后获取推流密钥</Form.Text>
              <Form.Control as="input" defaultValue={wxSecret} onChange={(e) => setWxSecret(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3" controlId="formWxEnabledCheckbox">
              <Form.Check type="checkbox" label="开启转推" defaultChecked={wxEnabled} onClick={() => setWxEnabled(!wxEnabled)} />
            </Form.Group>
            <Button
              variant="primary"
              type="submit"
              onClick={(e) => updateSecrets(e, 'update', 'wx', wxServer, wxSecret, wxEnabled)}
            >
              更新配置
            </Button> &nbsp;
            <TutorialsButton prefixLine={true} tutorials={forwardTutorials} /> &nbsp;
            <Form.Text> * 若有多个流，随机选择一个</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>Bilibili直播间</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>推流地址</Form.Label>
              <Form.Text> * 请先<a href='https://link.bilibili.com/p/center/index#/my-room/start-live' target='_blank' rel='noreferrer'>开始直播</a>，然后获取推流地址</Form.Text>
              <Form.Control as="input" defaultValue={bilibiliServer} onChange={(e) => setBilibiliServer(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>推流密钥</Form.Label>
              <Form.Text> * 请先<a href='https://link.bilibili.com/p/center/index#/my-room/start-live' target='_blank' rel='noreferrer'>开始直播</a>，然后获取推流密钥</Form.Text>
              <Form.Control as="input" defaultValue={bilibiliSecret} onChange={(e) => setBilibiliSecret(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3" controlId="formBilibiliEnabledCheckbox">
              <Form.Check type="checkbox" label="开启转推" defaultChecked={bilibiliEnabled} onClick={() => setBilibiliEnabled(!bilibiliEnabled)} />
            </Form.Group>
            <Button
              variant="primary"
              type="submit"
              onClick={(e) => updateSecrets(e, 'update', 'bilibili', bilibiliServer, bilibiliSecret, bilibiliEnabled)}
            >
              更新配置
            </Button> &nbsp;
            <TutorialsButton prefixLine={true} tutorials={forwardTutorials} /> &nbsp;
            <Form.Text> * 若有多个流，随机选择一个</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>快手云直播</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>推流地址</Form.Label>
              <Form.Text> * 请先<a href='https://studio.kuaishou.com/live/list' target='_blank' rel='noreferrer'>创建直播</a>，进入直播详情页，然后获取推流地址</Form.Text>
              <Form.Control as="input" defaultValue={kuaishouServer} onChange={(e) => setKuaishouServer(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>推流密钥</Form.Label>
              <Form.Text> * 请先<a href='https://studio.kuaishou.com/live/list' target='_blank' rel='noreferrer'>创建直播</a>，进入直播详情页，然后获取推流密钥</Form.Text>
              <Form.Control as="input" defaultValue={kuaishouSecret} onChange={(e) => setKuaishouSecret(e.target.value)}/>
            </Form.Group>
            <Form.Group className="mb-3" controlId="formKuaishouEnabledCheckbox">
              <Form.Check type="checkbox" label="开启转推" defaultChecked={kuaishouEnabled} onClick={() => setKuaishouEnabled(!kuaishouEnabled)} />
            </Form.Group>
            <Button
              variant="primary"
              type="submit"
              onClick={(e) => updateSecrets(e, 'update', 'kuaishou', kuaishouServer, kuaishouSecret, kuaishouEnabled)}
            >
              更新配置
            </Button> &nbsp;
            <TutorialsButton prefixLine={true} tutorials={forwardTutorials} /> &nbsp;
            <Form.Text> * 若有多个流，随机选择一个</Form.Text>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="99">
        <Accordion.Header>转推状态</Accordion.Header>
        <Accordion.Body>
          {
            forwards?.length ? (
              <Table striped bordered hover>
                <thead>
                <tr>
                  <th>#</th>
                  <th>平台</th>
                  <th>状态</th>
                  <th>更新时间</th>
                  <th>转发流</th>
                  <th>日志</th>
                </tr>
                </thead>
                <tbody>
                {
                  forwards?.map(file => {
                    return <tr key={file.platform} style={{verticalAlign: 'middle'}}>
                      <td>{file.i}</td>
                      <td>{file.name}</td>
                      <td>
                        <Badge bg={file.enabled ? 'success' : 'secondary'}>
                          {file.enabled ? '转发中' : '未开启'}
                        </Badge>
                      </td>
                      <td>{file.update && `${file.update?.format('YYYY-MM-DD HH:mm:ss')}`}</td>
                      <td>{file.stream}</td>
                      <td>{file.frame?.log}</td>
                    </tr>;
                  })
                }
                </tbody>
              </Table>
            ) : ''
          }
          {!forwards?.length ? '没有流。请开启转发并推流后，等待大约30秒左右，转发列表会自动更新' : ''}
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

