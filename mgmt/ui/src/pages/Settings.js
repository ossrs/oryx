import React from "react";
import {Accordion, Container, Form, Button, Tabs, Tab} from "react-bootstrap";
import {Errors, Token, PlatformPublicKey} from "../utils";
import axios from "axios";
import {useNavigate, useSearchParams} from "react-router-dom";
import {TutorialsButton, useTutorials} from '../components/TutorialsButton';
import SetupCamSecret from '../components/SetupCamSecret';

function SettingsImpl({defaultActiveTab}) {
  const navigate = useNavigate();
  const [key, setKey] = React.useState();
  const [crt, setCrt] = React.useState();
  const [domain, setDomain] = React.useState();
  const [secret, setSecret] = React.useState();
  const [beian, setBeian] = React.useState();
  const [activeTab, setActiveTab] = React.useState(defaultActiveTab);
  const setSearchParams = useSearchParams()[1];

  const sslTutorials = useTutorials(React.useRef([
    {author: '程晓龙', id: 'BV1tZ4y1R7qp'},
  ]));

  const updateBeian = (e) => {
    e.preventDefault();

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/beian/update', {
      ...token, beian: 'icp', text: beian,
    }).then(res => {
      alert('设置备案信息成功，请刷新页面');
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/routers-logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  };

  const enablePlatformAccess = (e, enabled) => {
    e.preventDefault();

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/pubkey', {
      ...token, enabled,
    }).then(res => {
      alert(enabled ? '授权平台管理员访问成功' : '取消授权成功');
      console.log(`PublicKey: Update ok, enabled=${enabled}`);
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/routers-logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  };

  const updateSSL = (e) => {
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
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/routers-logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  };

  const requestLetsEncrypt = (e) => {
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
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/routers-logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  };

  const updateSecret = (e) => {
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
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/routers-logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  };

  const onSelectTab = (k) => {
    setSearchParams({'tab': k});
    setActiveTab(k);
  };

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
                  <SetupCamSecret submitTips=' * 会自动创建依赖的云资源' />
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
            </Accordion>
          </Tab>
        </Tabs>
      </Container>
    </>
  );
}

export default function Settings() {
  const [searchParams] = useSearchParams();
  const [defaultActiveTab, setDefaultActiveTab] = React.useState();

  React.useEffect(() => {
    const tab = searchParams.get('tab') || 'auth';
    console.log(`?tab=https|auth|tencent|beian|platform, current=${tab}, Select the tab to render`);
    setDefaultActiveTab(tab);
  }, [searchParams]);

  return (<>
    { defaultActiveTab && <SettingsImpl defaultActiveTab={defaultActiveTab} /> }
  </>);
}

