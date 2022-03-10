import React from "react";
import {Accordion, Container, Form, Button, Tabs, Tab, InputGroup} from "react-bootstrap";
import {Token, PlatformPublicKey} from "../utils";
import axios from "axios";
import {useSearchParams} from "react-router-dom";
import {TutorialsButton, useTutorials} from '../components/TutorialsButton';
import SetupCamSecret from '../components/SetupCamSecret';
import {SrsErrorBoundary} from "../components/ErrorBoundary";
import {useErrorHandler} from "react-error-boundary";

export default function Settings() {
  return (
    <SrsErrorBoundary>
      <SettingsImpl />
    </SrsErrorBoundary>
  );
}

function SettingsImpl() {
  const [searchParams] = useSearchParams();
  const [defaultActiveTab, setDefaultActiveTab] = React.useState();
  const [upgradeWindow, setUpgradeWindow] = React.useState();
  const handleError = useErrorHandler();

  React.useEffect(() => {
    const tab = searchParams.get('tab') || 'auth';
    console.log(`?tab=https|auth|tencent|beian|platform, current=${tab}, Select the tab to render`);
    setDefaultActiveTab(tab);
  }, [searchParams]);

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/mgmt/window/query', {
      ...token,
    }).then(res => {
      const data = res.data.data;
      const win = {
        ...data,
        end: (data.start + data.duration)%24,
      };

      setUpgradeWindow(win);
      console.log(`Query upgrade window ${JSON.stringify(win)}`);
    }).catch(handleError);
  }, [handleError]);

  return (<>
    {
      defaultActiveTab && upgradeWindow &&
      <SettingsImpl2 defaultActiveTab={defaultActiveTab} defaultWindow={upgradeWindow} />
    }
  </>);
}

function SettingsImpl2({defaultActiveTab, defaultWindow}) {
  const [key, setKey] = React.useState();
  const [crt, setCrt] = React.useState();
  const [domain, setDomain] = React.useState();
  const [secret, setSecret] = React.useState();
  const [beian, setBeian] = React.useState();
  const [activeTab, setActiveTab] = React.useState(defaultActiveTab);
  const setSearchParams = useSearchParams()[1];
  const handleError = useErrorHandler();

  const sslTutorials = useTutorials(React.useRef([
    {author: '程晓龙', id: 'BV1tZ4y1R7qp'},
  ]));

  const timeSeries = React.useRef([...Array(25).keys()]).current;
  const [upgradeWindowStart, setUpgradeWindowStart] = React.useState(defaultWindow.start);
  const [upgradeWindowEnd, setUpgradeWindowEnd] = React.useState(defaultWindow.end);

  const updateBeian = React.useCallback((e) => {
    e.preventDefault();

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/beian/update', {
      ...token, beian: 'icp', text: beian,
    }).then(res => {
      alert('设置备案信息成功，请刷新页面');
    }).catch(handleError);
  }, [handleError, beian]);

  const enablePlatformAccess = React.useCallback((e, enabled) => {
    e.preventDefault();

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/pubkey', {
      ...token, enabled,
    }).then(res => {
      alert(enabled ? '授权平台管理员访问成功' : '取消授权成功');
      console.log(`PublicKey: Update ok, enabled=${enabled}`);
    }).catch(handleError);
  }, [handleError]);

  const updateSSL = React.useCallback((e) => {
    e.preventDefault();

    if (!key || !crt) {
      alert('请输入SSL密钥和证书');
      return;
    }

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/ssl', {
      ...token, key, crt,
    }).then(res => {
      alert(`SSL证书更新成功`);
      console.log(`SSL: Update ok`);
    }).catch(handleError);
  }, [handleError, key, crt]);

  const requestLetsEncrypt = React.useCallback((e) => {
    e.preventDefault();

    if (!domain) {
      alert('请输入你域名');
      return;
    }

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/letsencrypt', {
      ...token, domain,
    }).then(res => {
      alert(`Let's Encrypt SSL证书更新成功`);
      console.log(`SSL: Let's Encrypt SSL ok`);
    }).catch(handleError);
  }, [handleError, domain]);

  const updateSecret = React.useCallback((e) => {
    e.preventDefault();

    if (!secret) {
      alert('请输入密钥');
      return;
    }

    const token = Token.load();
    axios.post('/terraform/v1/hooks/srs/secret/update', {
      ...token, secret,
    }).then(res => {
      alert(`推流密钥更新成功`);
      console.log(`Secret: Update ok`);
    }).catch(handleError);
  }, [handleError, secret]);

  const onSelectTab = React.useCallback((k) => {
    setSearchParams({'tab': k});
    setActiveTab(k);
  }, []);

  const setupUpgradeWindow = React.useCallback((e) => {
    e.preventDefault();

    const [start, end] = [parseInt(upgradeWindowStart || 0), parseInt(upgradeWindowEnd || 0)];

    const duration = start < end ? end - start : end + 24 - start;
    if (duration <= 3) return alert(`升级窗口不能小于3小时`);

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/window/update', {
      ...token, start, duration,
    }).then(res => {
      alert(`升级窗口[${start}点至${end}点]更新成功，窗口长度${duration}小时`);
      console.log(`Setup upgrade window start=${start}, end=${end}`);
    }).catch(handleError);
  }, [handleError, upgradeWindowStart, upgradeWindowEnd]);

  return (
    <>
      <p></p>
      <Container>
        <Tabs defaultActiveKey={activeTab} id="uncontrolled-tab-example" className="mb-3" onSelect={(k) => onSelectTab(k)}>
          <Tab eventKey="https" title="HTTPS">
            <Accordion defaultActiveKey="0">
              <Accordion.Item eventKey="0">
                <Accordion.Header>HTTPS: Let's Encrypt</Accordion.Header>
                <Accordion.Body>
                  <Form>
                    <Form.Group className="mb-3">
                      <Form.Label>域名</Form.Label>
                      <Form.Text> * 你的域名，请先解析到本服务器的公网IP，例如 your-domain.com</Form.Text>
                      <Form.Control as="input" defaultValue={domain} onChange={(e) => setDomain(e.target.value)} />
                    </Form.Group>
                    <Button variant="primary" type="submit" onClick={(e) => requestLetsEncrypt(e)}>
                      申请证书
                    </Button> &nbsp;
                    <TutorialsButton prefixLine={true} tutorials={sslTutorials} />
                  </Form>
                </Accordion.Body>
              </Accordion.Item>
              <Accordion.Item eventKey="1">
                <Accordion.Header>HTTPS: 上传证书</Accordion.Header>
                <Accordion.Body>
                  <Form>
                    <Form.Group className="mb-3">
                      <Form.Label>密钥(KEY)</Form.Label>
                      <Form.Text> * Nginx格式的SSL密钥内容，例如 your-domain.com.key</Form.Text>
                      <Form.Control as="textarea" rows={5} defaultValue={key} onChange={(e) => setKey(e.target.value)} />
                    </Form.Group>
                    <Form.Group className="mb-3">
                      <Form.Label>证书(PEM格式)</Form.Label>
                      <Form.Text> * Nginx格式的SSL证书内容，例如 your-domain.com.pem</Form.Text>
                      <Form.Control as="textarea" rows={5} defaultValue={crt} onChange={(e) => setCrt(e.target.value)} />
                    </Form.Group>
                    <Button variant="primary" type="submit" onClick={(e) => updateSSL(e)}>
                      更新证书
                    </Button> &nbsp;
                    <TutorialsButton prefixLine={true} tutorials={sslTutorials} />
                  </Form>
                </Accordion.Body>
              </Accordion.Item>
            </Accordion>
          </Tab>
          <Tab eventKey="auth" title="流鉴权">
            <Accordion defaultActiveKey="0">
              <Accordion.Item eventKey="0">
                <Accordion.Header>更新流密钥</Accordion.Header>
                <Accordion.Body>
                  <Form>
                    <Form.Group className="mb-3">
                      <Form.Label>密钥</Form.Label>
                      <Form.Text> * 推流鉴权的密钥</Form.Text>
                      <Form.Control as="input" defaultValue={secret} onChange={(e) => setSecret(e.target.value)}/>
                    </Form.Group>
                    <Button variant="primary" type="submit" onClick={(e) => updateSecret(e, true)}>
                      更新
                    </Button> &nbsp;
                  </Form>
                </Accordion.Body>
              </Accordion.Item>
            </Accordion>
          </Tab>
          <Tab eventKey="tencent" title="腾讯云">
            <Accordion defaultActiveKey="0">
              <Accordion.Item eventKey="0">
                <Accordion.Header>腾讯云密钥(Secret)</Accordion.Header>
                <Accordion.Body>
                  <SetupCamSecret />
                </Accordion.Body>
              </Accordion.Item>
            </Accordion>
          </Tab>
          <Tab eventKey="beian" title="备案">
            <Accordion defaultActiveKey="0">
              <Accordion.Item eventKey="0">
                <Accordion.Header>设置备案号</Accordion.Header>
                <Accordion.Body>
                  <Form>
                    <Form.Group className="mb-3">
                      <Form.Label>ICP备案号</Form.Label>
                      <Form.Text> * 请参考<a href='https://beian.miit.gov.cn' target='_blank' rel='noreferrer'>https://beian.miit.gov.cn</a></Form.Text>
                      <Form.Control as="input" defaultValue={beian} placeholder='例如：京ICP备XXXXXXXX号-X' onChange={(e) => setBeian(e.target.value)}/>
                    </Form.Group>
                    <Button variant="primary" type="submit" onClick={(e) => updateBeian(e)}>
                      设置
                    </Button>
                  </Form>
                </Accordion.Body>
              </Accordion.Item>
            </Accordion>
          </Tab>
          <Tab eventKey="platform" title="平台">
            <Accordion defaultActiveKey="0">
              <Accordion.Item eventKey="0">
                <Accordion.Header>授权平台管理员</Accordion.Header>
                <Accordion.Body>
                  <Form>
                    <Form.Group className="mb-3">
                      <Form.Label>公钥</Form.Label>
                      <Form.Text> * 平台管理员的公钥</Form.Text>
                      <Form.Control as="textarea" rows={2} defaultValue={PlatformPublicKey} readOnly={true} />
                    </Form.Group>
                    <Button variant="primary" type="submit" onClick={(e) => enablePlatformAccess(e, true)}>
                      授权访问
                    </Button> &nbsp;
                    <Button variant="primary" type="submit" onClick={(e) => enablePlatformAccess(e, false)}>
                      取消授权
                    </Button>
                  </Form>
                </Accordion.Body>
              </Accordion.Item>
              <Accordion.Item eventKey="1">
                <Accordion.Header>设置升级窗口</Accordion.Header>
                <Accordion.Body>
                  <Form>
                    <Form.Label htmlFor="basic-url">升级窗口</Form.Label>
                    <Form.Text> * 系统会在这个时间段，自动升级到最新的稳定版本</Form.Text>
                    <InputGroup className="mb-3">
                      <InputGroup.Text>开始时间</InputGroup.Text>
                      <Form.Select
                        aria-label="Start time"
                        defaultValue={upgradeWindowStart}
                        onChange={(e) => setUpgradeWindowStart(e.target.value)}
                      >
                        {timeSeries.map((e) => {
                          return <option key={e} value={e}>{`${String(e).padStart(2, '0')}:00`}</option>;
                        })}
                      </Form.Select>
                    </InputGroup>
                    <InputGroup className="mb-3">
                      <InputGroup.Text>结束时间</InputGroup.Text>
                      <Form.Select
                        aria-label="End time"
                        defaultValue={upgradeWindowEnd}
                        onChange={(e) => setUpgradeWindowEnd(e.target.value)}
                      >
                        {timeSeries.map((e) => {
                          return <option key={e} value={e}>{`${String(e).padStart(2, '0')}:00`}</option>;
                        })}
                      </Form.Select>
                    </InputGroup>
                    <Button variant="primary" type="submit" onClick={(e) => setupUpgradeWindow(e)}>
                      设置窗口
                    </Button>
                  </Form>
                </Accordion.Body>
              </Accordion.Item>
            </Accordion>
          </Tab>
        </Tabs>
      </Container>
    </>
  );
}

