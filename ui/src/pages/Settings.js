import React from "react";
import {Accordion, Container, Form, Button, Tabs, Tab} from "react-bootstrap";
import {Clipboard, Token} from "../utils";
import axios from "axios";
import {useSearchParams} from "react-router-dom";
import SetupCamSecret from '../components/SetupCamSecret';
import {SrsErrorBoundary} from "../components/SrsErrorBoundary";
import {useErrorHandler} from "react-error-boundary";
import {useTranslation} from "react-i18next";
import {SrsEnvContext} from "../components/SrsEnvContext";

export default function Systems() {
  return (
    <SrsErrorBoundary>
      <SystemsImpl />
    </SrsErrorBoundary>
  );
}

function SystemsImpl() {
  const [searchParams] = useSearchParams();
  const [defaultActiveTab, setDefaultActiveTab] = React.useState();

  React.useEffect(() => {
    const tab = searchParams.get('tab') || 'auth';
    console.log(`?tab=https|auth|tencent|beian|platform, current=${tab}, Select the tab to render`);
    setDefaultActiveTab(tab);
  }, [searchParams]);

  return (<>
    {
      defaultActiveTab &&
      <SettingsImpl2 defaultActiveTab={defaultActiveTab} />
    }
  </>);
}

function SettingsImpl2({defaultActiveTab}) {
  const [activeTab, setActiveTab] = React.useState(defaultActiveTab);
  const setSearchParams = useSearchParams()[1];
  const {t} = useTranslation();

  const onSelectTab = React.useCallback((k) => {
    setSearchParams({'tab': k});
    setActiveTab(k);
  }, [setSearchParams]);

  const copyToClipboard = React.useCallback((e, text) => {
    e.preventDefault();

    Clipboard.copy(text).then(() => {
      alert(`已经复制到剪切板`);
    }).catch((err) => {
      alert(`复制失败，请右键复制链接 ${err}`);
    });
  }, []);

  return (
    <>
      <p></p>
      <Container>
        <Tabs defaultActiveKey={activeTab} id="uncontrolled-tab-example" className="mb-3" onSelect={(k) => onSelectTab(k)}>
          <Tab eventKey="auth" title={t('settings.tabAuth')}>
            <SettingAuth />
          </Tab>
          <Tab eventKey="https" title="HTTPS">
            <SettingHttpsDisabled />
          </Tab>
          <Tab eventKey="nginx" title="NGINX">
            <SettingNginx />
          </Tab>
          <Tab eventKey="beian" title={t('settings.tabFooter')}>
            <SettingBeian />
          </Tab>
          <Tab eventKey="tencent" title={t('settings.tabTencent')}>
            <SettingTencent />
          </Tab>
          <Tab eventKey="api" title="OpenAPI">
            <SettingOpenApi {...{copyToClipboard}}/>
          </Tab>
        </Tabs>
      </Container>
    </>
  );
}

function SettingNginx() {
  const [hlsDelivery, setHlsDelivery] = React.useState();
  const handleError = useErrorHandler();
  const {t} = useTranslation();
  const env = React.useContext(SrsEnvContext)[0];

  const updateHlsDelivery = React.useCallback((e) => {
    e.preventDefault();

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/nginx/hls', {
      ...token, enabled: hlsDelivery,
    }).then(res => {
      alert(t('helper.setOk'));
    }).catch(handleError);
  }, [handleError, hlsDelivery, t]);

  if (env?.mgmtDocker) {
    return <span style={{color: 'red'}}>{t('errs.noNginx')}</span>;
  }

  return (
    <Accordion defaultActiveKey="0">
      <Accordion.Item eventKey="0">
        <Accordion.Header>{t('settings.nginxHlsTitle')}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3" controlId="formDvrAllCheckbox">
              <Form.Check type="checkbox" label={t('settings.nginxHlsTip')} defaultChecked={hlsDelivery} onClick={() => setHlsDelivery(!hlsDelivery)} />
            </Form.Group>
            <Button variant="primary" type="submit" onClick={(e) => updateHlsDelivery(e)}>
              {t('helper.submit')}
            </Button>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

function SettingOpenApi({copyToClipboard}) {
  const [apiSecret, setAPISecret] = React.useState();
  const handleError = useErrorHandler();
  const {t} = useTranslation();

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/mgmt/secret/query', {
      ...token,
    }).then(res => {
      setAPISecret(res.data.data);
      console.log(`Status: Query ok, apiSecret=${JSON.stringify(res.data.data)}`);
    }).catch(handleError);
  }, [handleError]);

  const copyApiSecret = React.useCallback((e, apiSecret) => {
    copyToClipboard(e, apiSecret);
  }, [copyToClipboard]);

  return (
    <Accordion defaultActiveKey="0">
      <Accordion.Item eventKey="0">
        <Accordion.Header>{t('openapi.title')}</Accordion.Header>
        <Accordion.Body>
          <div>
            {t('openapi.summary')}
            <p></p>
          </div>
          <p>Usage:</p>
          <ul>
            <li> {t('openapi.usage1')} </li>
            <li> {t('openapi.usage2')} </li>
          </ul>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{t('openapi.secret')}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>ApiSecret</Form.Label>
              <Form.Text> * {t('openapi.secretTip')}</Form.Text>
              <Form.Control as="input" type='input' rows={1} defaultValue={apiSecret} readOnly={true}/>
            </Form.Group>
            <Button variant="primary" type="submit" onClick={(e) => copyApiSecret(e, apiSecret)}>
              {t('openapi.secretCopy')}
            </Button>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>{t('openapi.apiPublishSecret')}</Accordion.Header>
        <Accordion.Body>
          <RunOpenAPI apiSecret={apiSecret} api='/terraform/v1/hooks/srs/secret/query' />
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

function SettingBeian() {
  const [beian, setBeian] = React.useState();
  const [siteTitle, setSiteTitle] = React.useState();
  const handleError = useErrorHandler();
  const {t} = useTranslation();

  const updateBeian = React.useCallback((e) => {
    e.preventDefault();

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/beian/update', {
      ...token, beian: 'icp', text: beian,
    }).then(res => {
      alert(t('settings.footer'));
    }).catch(handleError);
  }, [handleError, beian, t]);

  const updateSiteTitle = React.useCallback((e) => {
    e.preventDefault();

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/beian/update', {
      ...token, beian: 'title', text: siteTitle,
    }).then(res => {
      alert(t('settings.header'));
    }).catch(handleError);
  }, [handleError, siteTitle, t]);

  return (
    <Accordion defaultActiveKey="1">
      <Accordion.Item eventKey="1">
        <Accordion.Header>{t('settings.footerTitle')}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>{t('settings.footerIcp')}</Form.Label>
              <Form.Control as="input" defaultValue={beian} placeholder={t('settings.footerHolder')} onChange={(e) => setBeian(e.target.value)}/>
            </Form.Group>
            <Button variant="primary" type="submit" onClick={(e) => updateBeian(e)}>
              {t('settings.footerSubmit')}
            </Button>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>{t('settings.headerTitle')}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>{t('settings.headerIcp')}</Form.Label>
              <Form.Control as="input" defaultValue={siteTitle} placeholder={t('settings.headerHolder')} onChange={(e) => setSiteTitle(e.target.value)}/>
            </Form.Group>
            <Button variant="primary" type="submit" onClick={(e) => updateSiteTitle(e)}>
              {t('settings.headerSubmit')}
            </Button>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

function SettingTencent() {
  const {t} = useTranslation();

  return (
    <Accordion defaultActiveKey="0">
      <Accordion.Item eventKey="0">
        <Accordion.Header>{t('settings.tecentTitle')}</Accordion.Header>
        <Accordion.Body>
          <SetupCamSecret />
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

function SettingAuth() {
  const [secret, setSecret] = React.useState();
  const [allowNoAuth, setAllowNoAuth] = React.useState();
  const [noAuth, setNoAuth] = React.useState();
  const [searchParams] = useSearchParams();
  const handleError = useErrorHandler();
  const {t} = useTranslation();

  React.useEffect(() => {
    const allowNoAuth = searchParams.get('allow-noauth') === 'true';
    console.log(`?allow-noauth=true|false, current=${allowNoAuth}, Whether allow disable auth`);
    setAllowNoAuth(allowNoAuth);
  }, [searchParams]);

  const updateSecret = React.useCallback((e) => {
    e.preventDefault();

    if (!secret) {
      alert(t('settings.secretNoValue'));
      return;
    }

    const token = Token.load();
    axios.post('/terraform/v1/hooks/srs/secret/update', {
      ...token, secret,
    }).then(res => {
      alert(t('settings.secretOk'));
      console.log(`Secret: Update ok`);
    }).catch(handleError);
  }, [handleError, secret, t]);

  const updateNoAuth = React.useCallback((e) => {
    const token = Token.load();
    axios.post('/terraform/v1/hooks/srs/secret/disable', {
      ...token, pubNoAuth: !!noAuth,
    }).then(res => {
      alert(t('helper.setOk'));
      console.log(`Disable: Update ok, noAuth=${noAuth}`);
    }).catch(handleError);
  }, [handleError, t, noAuth])

  return (
    <Accordion defaultActiveKey="0">
      <Accordion.Item eventKey="0">
        <Accordion.Header>{t('settings.authTitle')}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>{t('settings.authSecret')}</Form.Label>
              <Form.Text> * {t('settings.authSecretTip')}</Form.Text>
              <Form.Control as="input" defaultValue={secret} onChange={(e) => setSecret(e.target.value)}/>
            </Form.Group>
            <Button variant="primary" type="submit" onClick={(e) => updateSecret(e, true)}>
              {t('settings.authSubmit')}
            </Button> &nbsp;
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      {allowNoAuth && <>
        <Accordion.Item eventKey="1">
          <Accordion.Header>{t('settings.noAuthTitle')}</Accordion.Header>
          <Accordion.Body>
            <Form.Group className="mb-3" controlId="formDisableAuthCheckbox">
              <Form.Check type="checkbox" label={t('settings.noAuthTips')} defaultChecked={noAuth} onClick={() => setNoAuth(!noAuth)} />
            </Form.Group>
            <Button variant="primary" type="submit" onClick={(e) => updateNoAuth(e)}>
              {t('helper.submit')}
            </Button>
          </Accordion.Body>
        </Accordion.Item>
      </>}
    </Accordion>
  );
}

function SettingHttpsDisabled() {
  const {t} = useTranslation();
  return <span style={{color: 'red'}}>{t('errs.btHttps1')}</span>;
}

function RunOpenAPI(props) {
  const [showResult, setShowResult] = React.useState();
  const {t} = useTranslation();
  const {apiSecret, api, data} = props;

  const onClick = React.useCallback((e) => {
    e.preventDefault();
    setShowResult(!showResult);
  }, [showResult]);

  if (!apiSecret) {
    return (
      <div>
        {t('openapi.secretEmpty')}<code>{t('openapi.secret')}</code>
      </div>
    );
  }

  return (
    <Form>
      <Form.Group className="mb-3">
        <Form.Label>URL</Form.Label>
        <Form.Control as="textarea" rows={1} defaultValue={`POST ${api}`} readOnly={true} />
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Headers</Form.Label>
        <Form.Control as="textarea" rows={1} defaultValue={`Authorization: Bearer ${apiSecret}`} readOnly={true} />
      </Form.Group>
      {data &&
        <Form.Group className="mb-3">
          <Form.Label>Body</Form.Label>
          <Form.Control as="textarea" rows={5} defaultValue={JSON.stringify(data, null, 2)} readOnly={true} />
          <pre>

          </pre>
        </Form.Group>
      }
      <Form.Group className="mb-3">
        { showResult && <SrsErrorBoundary><OpenAPIResult {...props} /></SrsErrorBoundary> }
      </Form.Group>
      <Button variant="primary" type="submit" onClick={(e) => onClick(e)}>
        {showResult ? 'Reset' : 'Run'}
      </Button> &nbsp;
    </Form>
  );
}

function OpenAPIResult({apiSecret, api, data}) {
  const [openAPIRes, setOpenAPIRes] = React.useState();
  const handleError = useErrorHandler();

  React.useEffect(() => {
    axios.post(api, data, {
      headers: {
        'Authorization': `Bearer ${apiSecret}`,
      }
    }).then(res => {
      setOpenAPIRes(res.data);
      console.log(`OpenAPI: Run api=${api} ok, data=${JSON.stringify(res.data.data)}`);
    }).catch(handleError);
  }, [handleError, apiSecret, api, data]);

  return (
    <>
      <Form.Label>Response</Form.Label>
      <pre>
      {JSON.stringify(openAPIRes, null, 2)}
      </pre>
    </>
  );
}

