//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {Accordion, Container, Form, Button, Tabs, Tab, Spinner, Stack, Badge} from "react-bootstrap";
import {Clipboard, Token} from "../utils";
import axios from "axios";
import {useSearchParams} from "react-router-dom";
import {SrsErrorBoundary} from "../components/SrsErrorBoundary";
import {useErrorHandler} from "react-error-boundary";
import {useTranslation} from "react-i18next";
import {TutorialsButton, useTutorials} from "../components/TutorialsButton";

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
    console.log(`?tab=https|hls|auth|beian|callback|platform, current=${tab}, Select the tab to render`);
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
      alert(t('helper.copyOk'));
    }).catch((err) => {
      alert(`${t('helper.copyFail')} ${err}`);
    });
  }, [t]);

  return (
    <>
      <p></p>
      <Container>
        <Tabs defaultActiveKey={activeTab} id="uncontrolled-tab-example" className="mb-3" onSelect={(k) => onSelectTab(k)}>
          <Tab eventKey="auth" title={t('settings.tabAuth')}>
            <SettingAuth />
          </Tab>
          <Tab eventKey="https" title="HTTPS">
            <SettingHttps />
          </Tab>
          <Tab eventKey="hls" title="HLS">
            <SettingHighPerformanceHLS />
          </Tab>
          <Tab eventKey="beian" title={t('settings.tabFooter')}>
            <SettingBeian />
          </Tab>
          <Tab eventKey="callback" title={t('settings.tabCallback')}>
            <SettingCallback />
          </Tab>
          <Tab eventKey="api" title="OpenAPI">
            <SettingOpenApi {...{copyToClipboard}}/>
          </Tab>
        </Tabs>
      </Container>
    </>
  );
}

function SettingHighPerformanceHLS() {
  const [noHlsCtx, setNoHlsCtx] = React.useState();
  const [hlsLL, setHlsLL] = React.useState();
  const handleError = useErrorHandler();
  const {t} = useTranslation();

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/mgmt/hphls/query', {
      ...token,
    }).then(res => {
      setNoHlsCtx(res.data.data.noHlsCtx === true);
      console.log(`Status: Query ok, hlsDelivery=${JSON.stringify(res.data.data)}`);
    }).catch(handleError);
  }, [handleError, setNoHlsCtx]);

  const updateHlsDelivery = React.useCallback((e) => {
    e.preventDefault();

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/hphls/update', {
      ...token, noHlsCtx: noHlsCtx,
    }).then(res => {
      alert(t('helper.setOk'));
    }).catch(handleError);
  }, [handleError, noHlsCtx, t]);

  React.useEffect(() => {
    axios.post('/terraform/v1/mgmt/hlsll/query', {
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      setHlsLL(res.data.data.hlsLowLatency === true);
      console.log(`Status: Query ok, hlsLowLatency=${JSON.stringify(res.data.data)}`);
    }).catch(handleError);
  }, [handleError, setHlsLL]);

  const updateHlsLL = React.useCallback((e) => {
    e.preventDefault();

    axios.post('/terraform/v1/mgmt/hlsll/update', {
      hlsLowLatency: hlsLL,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      alert(t('helper.setOk'));
    }).catch(handleError);
  }, [handleError, hlsLL, t]);

  return (
    <Accordion defaultActiveKey={['0','1']} alwaysOpen>
      <Accordion.Item eventKey="0">
        <Accordion.Header>{t('settings.nginxHlsTitle')}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3" controlId="formNginxHlsCheckbox">
              <Form.Check type="checkbox" label={t('settings.nginxHlsTip')} defaultChecked={noHlsCtx} onClick={() => setNoHlsCtx(!noHlsCtx)} />
            </Form.Group>
            <Button variant="primary" type="submit" onClick={(e) => updateHlsDelivery(e)}>
              {t('helper.submit')}
            </Button>
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{t('settings.hlsLL')}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3" controlId="formHlsLLCheckbox">
              <Form.Check type="checkbox" label={t('settings.hlsLLTip')} defaultChecked={hlsLL} onClick={() => setHlsLL(!hlsLL)} />
            </Form.Group>
            <Button variant="primary" type="submit" onClick={(e) => updateHlsLL(e)}>
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

  const curlCode = `
curl ${window.location.protocol}//${window.location.host}/terraform/v1/hooks/srs/secret/query \\
  -X POST -H 'Authorization: Bearer ${apiSecret}' \\
  -H 'Content-Type: application/json' --data '{}'
  `;

  const jQueryCode = `
$.ajax({
  url: '${window.location.protocol}//${window.location.host}/terraform/v1/hooks/srs/secret/query',
  type: 'POST',
  headers: {
    "Authorization": "Bearer ${apiSecret}",
  },
  dataType: 'json',
  contentType: "application/json",
  data: {}, 
  success: function () {},
  error: function () {},
});
  `;

  return (
    <Accordion defaultActiveKey={["0", "1", "2", "3", "4"]} alwaysOpen>
      <Accordion.Item eventKey="0">
        <Accordion.Header>{t('openapi.title')}</Accordion.Header>
        <Accordion.Body>
          <div>
            {t('openapi.summary')}
          </div>
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
      <Accordion.Item eventKey="2">
        <Accordion.Header>{t('openapi.apiPublishSecret')}</Accordion.Header>
        <Accordion.Body>
          <RunOpenAPI apiSecret={apiSecret} api='/terraform/v1/hooks/srs/secret/query' />
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>{t('openapi.apiPublishSecret2')}</Accordion.Header>
        <Accordion.Body>
          <pre>{curlCode}</pre>
          <Button variant="primary" type="submit" onClick={(e) => copyApiSecret(e, curlCode)}>
            {t('openapi.curlCopy')}
          </Button>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="4">
        <Accordion.Header>{t('openapi.apiPublishSecret3')}</Accordion.Header>
        <Accordion.Body>
          <pre>{jQueryCode}</pre>
          <Button variant="primary" type="submit" onClick={(e) => copyApiSecret(e, jQueryCode)}>
            {t('openapi.codeCopy')}
          </Button>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

function SettingCallback() {
  const handleError = useErrorHandler();
  const [config, setConfig] = React.useState();
  const [activeKey, setActiveKey] = React.useState();

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/mgmt/hooks/query', {
      ...token
    }).then(res => {
      const data = res.data.data;

      setConfig(data);
      if (data.all) {
        setActiveKey('2');
      } else {
        setActiveKey('1');
      }
      console.log(`Hooks: Query ok, ${JSON.stringify(data)}`);
    }).catch(handleError);
  }, [handleError]);

  if (!activeKey) return <></>;
  return <SettingCallbackImpl {...{
    activeKey, defaultEnabled: config?.all, defaultConf: config
  }}/>;
}

function SettingCallbackImpl({activeKey, defaultEnabled, defaultConf}) {
  const defaultUrl = `${window.location.protocol}//${window.location.host}/terraform/v1/mgmt/hooks/example?fail=false`;
  const handleError = useErrorHandler();
  const {t} = useTranslation();
  const [target, setTarget] = React.useState(defaultConf.target || defaultUrl);
  const [opaque, setOpaque] = React.useState(defaultConf.opaque);
  const [allEvents, setAllEvents] = React.useState(defaultEnabled);
  const [task, setTask] = React.useState();

  React.useEffect(() => {
    const refreshTask = () => {
      const token = Token.load();
      axios.post('/terraform/v1/mgmt/hooks/query', {
        ...token,
      }).then(res => {
        const task = res.data.data;
        if (task?.req) task.req = JSON.parse(task.req);
        if (task?.res) task.res = JSON.parse(task.res);
        setTask(task);
        console.log(`Hooks: Query task ${JSON.stringify(task)}`);
      }).catch(handleError);
    };

    refreshTask();
    const timer = setInterval(() => refreshTask(), 10 * 1000);
    return () => clearInterval(timer);
  }, [handleError, setTask]);

  const onUpdateCallback = React.useCallback((e) => {
    e.preventDefault();

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/hooks/apply', {
      ...token, all: !!allEvents, target, opaque,
    }).then(res => {
      alert(t('helper.setOk'));
      console.log(`Hooks apply ok, all=${allEvents}, target=${target}, opaque=${opaque}, response=${JSON.stringify(res.data.data)}`);
    }).catch(handleError);
  }, [handleError, t, allEvents, target, opaque]);

  return <Accordion defaultActiveKey={[activeKey]} alwaysOpen>
    <Accordion.Item eventKey="0">
      <Accordion.Header>{t('cb.title')}</Accordion.Header>
      <Accordion.Body>
        <div>
          {t('cb.summary')}
          <p></p>
        </div>
        <p>Usage:</p>
        <ul>
          <li> {t('cb.usage1')} </li>
          <li> {t('cb.usage2')} </li>
        </ul>
      </Accordion.Body>
    </Accordion.Item>
    <Accordion.Item eventKey="1">
      <Accordion.Header>{t('cb.setting')}</Accordion.Header>
      <Accordion.Body>
        <Form>
          <Form.Group className="mb-3">
            <Form.Label>{t('cb.target')}</Form.Label>
            <Form.Text>* {t('cb.target2')}</Form.Text>
            <Form.Control as="input" placeholder='For example: http://your-server/callback' defaultValue={target}
                          onChange={(e) => setTarget(e.target.value)} />
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>{t('cb.opaque')}</Form.Label>
            <Form.Text>* {t('cb.opaque2')}</Form.Text>
            <Form.Control as="input" placeholder='For example: authentication secret for hooks' defaultValue={opaque}
                          onChange={(e) => setOpaque(e.target.value)} />
          </Form.Group>
          <Form.Group className="mb-3" controlId="formCbAllCheckbox">
            <Form.Label>{t('cb.event')}</Form.Label>
            <Form.Text>
              * {t('cb.event2')}: &nbsp;
              <a href={t('helper.doc')+'#http-callback-on_publish'} target='_blank' rel='noreferrer'>publish,</a> &nbsp;
              <a href={t('helper.doc')+'#http-callback-on_unpublish'} target='_blank' rel='noreferrer'>unpublish</a> &nbsp;
            </Form.Text>
            <Form.Check type="checkbox" defaultChecked={allEvents} label={t('cb.event3')}
                        onChange={(e) => setAllEvents(!allEvents)} />
          </Form.Group>
          <Button variant="primary" type="submit" onClick={(e) => onUpdateCallback(e)}>
            {t('settings.footerSubmit')}
          </Button>
        </Form>
      </Accordion.Body>
    </Accordion.Item>
    <Accordion.Item eventKey="2">
      <Accordion.Header>{t('cb.show')}</Accordion.Header>
      <Accordion.Body>
        <Stack gap={1}>
          <div>
            <Badge bg={allEvents ? 'success' : 'secondary'}>
              {allEvents ? t('cb.active') : t('cb.inactive')}
            </Badge>
          </div>
          <div>{t('cb.show2')}</div>
          <div>{t('cb.target')}: {target}</div>
          <div>
            {t('cb.req')}:
            <pre>{task?.req && JSON.stringify(task.req, null, 2)}</pre>
          </div>
          <div>
            {t('cb.res')}:
            <pre>{task?.res && JSON.stringify(task.res, null, 2)}</pre>
          </div>
        </Stack>
      </Accordion.Body>
    </Accordion.Item>
  </Accordion>;
}

function SettingBeian() {
  const [beian, setBeian] = React.useState();
  const [siteTitle, setSiteTitle] = React.useState();
  const handleError = useErrorHandler();
  const {t} = useTranslation();

  React.useEffect(() => {
    axios.get('/terraform/v1/mgmt/beian/query')
      .then(res => {
        setSiteTitle(res.data.data.title);
        setBeian(res.data.data.icp);
        console.log(`Beian: query ${JSON.stringify(res.data.data)}`);
      }).catch(handleError);
  }, [handleError]);

  // Update the footer for beian.
  const updateBeian = React.useCallback((e) => {
    e.preventDefault();

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/beian/update', {
      ...token, beian: 'icp', text: beian,
    }).then(res => {
      alert(t('settings.footer'));
    }).catch(handleError);
  }, [handleError, beian, t]);

  // Update the title for site.
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
    <Accordion defaultActiveKey={["1", "2"]} alwaysOpen>
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

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/hooks/srs/secret/query', {
      ...token,
    }).then(res => {
      setSecret(res.data.data.publish);
      console.log(`Status: Query ok, secret=${JSON.stringify(res.data.data)}`);
    }).catch(handleError);
  }, [handleError]);

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
    <Accordion defaultActiveKey={["0"]} alwaysOpen>
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

function SettingHttps() {
  const [config, setConfig] = React.useState({});
  const [loading, setLoading] = React.useState(true);
  const handleError = useErrorHandler();

  React.useEffect(() => {
    setLoading(true);

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/cert/query', {
      ...token,
    }).then(res => {
      setConfig(res?.data?.data || {});
      console.log(`SSL: Query ok, provider=${res?.data?.data?.provider}`);
    }).catch(handleError).finally(setLoading);
  }, [handleError, setLoading, setConfig]);

  return !loading ? <SettingHttpsImpl config={config} /> : <></>;
}

function SettingHttpsImpl({config}) {
  const [key, setKey] = React.useState(config.key);
  const [crt, setCrt] = React.useState(config.crt);
  const [domain, setDomain] = React.useState(config.domain);
  const [operating, setOperating] = React.useState(false);
  const handleError = useErrorHandler();
  const {t} = useTranslation();

  const domainRegex = React.useMemo(() => {
    return /^(?=.{1,253})(?:(?!-)[A-Za-z0-9-]{1,63}(?<!-)\.?)+[A-Za-z]{2,6}$/;
  }, []);

  React.useEffect(() => {
    if (!domainRegex.test(window.location.hostname)) return;
    if (domain) return;
    setDomain(window.location.hostname);
  }, [domain, domainRegex, setDomain]);

  const sslTutorials = useTutorials({
    bilibili: React.useRef([
      {author: '程晓龙', id: 'BV1tZ4y1R7qp'},
    ]),
    medium: React.useRef([
      {id: 'cb618777639f'},
    ])
  });

  const updateSSL = React.useCallback(async (e) => {
    e.preventDefault();

    if (!key || !crt) {
      alert(t('settings.sslNoFile'));
      return;
    }

    setOperating(true);

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/ssl', {
      ...token, key, crt,
    }).then(res => {
      alert(t('settings.sslOk'));
      console.log(`SSL: Update ok`);
    }).catch(handleError).finally(setOperating);
  }, [handleError, key, crt, t, setOperating]);

  const requestLetsEncrypt = React.useCallback((e) => {
    e.preventDefault();

    const port = window.location.port;
    if (port && port !== '80' && port !== '3000') {
      return alert(`${t('settings.sslInvalidPort')} ${port} (${window.location.host})`);
    }

    if (!domainRegex.test(window.location.hostname)) {
      return alert(`${t('settings.sslInvalidHost')} "${window.location.hostname}", ${t('settings.sslInvalidHost3')}`);
    }

    if (!domain) {
      return alert(t('settings.sslNoDomain'));
    }

    if (!domainRegex.test(domain)) {
      return alert(t('settings.sslInvalidDomain'));
    }

    if (window.location.hostname !== domain) {
      return alert(`${t('settings.sslInvalidHost')} "${window.location.hostname}", ${t('settings.sslInvalidHost2')} "${domain}"`);
    }

    setOperating(true);

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/letsencrypt', {
      ...token, domain,
    }).then(res => {
      alert(t('settings.sslLetsOk'));
      console.log(`SSL: Let's Encrypt SSL ok`);
    }).catch(handleError).finally(setOperating);
  }, [handleError, domain, t, setOperating, domainRegex]);

  const defaultKey = config?.provider === 'ssl' ? '1' : '0';
  return (
    <Accordion defaultActiveKey={[defaultKey]} alwaysOpen>
      <Accordion.Item eventKey="0">
        <Accordion.Header>{t('settings.letsTitle')}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>{t('settings.letsDomain')}</Form.Label>
              <Form.Text> * {t('settings.letsDomainTip')}</Form.Text>
              <Form.Control as="input" defaultValue={domain} onChange={(e) => setDomain(e.target.value)} />
            </Form.Group>
            <Button variant="primary" type="submit" disabled={operating} onClick={(e) => requestLetsEncrypt(e)}>
              {t('settings.letsDomainSubmit')}
            </Button> &nbsp;
            <TutorialsButton prefixLine={true} tutorials={sslTutorials} /> &nbsp;
            {operating && <Spinner animation="border" variant="success" style={{verticalAlign: 'middle'}} />}
          </Form>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>{t('settings.sslFileTitle')}</Accordion.Header>
        <Accordion.Body>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>{t('settings.sslFileKey')}</Form.Label>
              <Form.Text> * {t('settings.sslFileKeyTip')}</Form.Text>
              <Form.Control as="textarea" rows={5} defaultValue={key} onChange={(e) => setKey(e.target.value)} />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>{t('settings.sslFileCert')}</Form.Label>
              <Form.Text> * {t('settings.sslFileCertTip')}</Form.Text>
              <Form.Control as="textarea" rows={5} defaultValue={crt} onChange={(e) => setCrt(e.target.value)} />
            </Form.Group>
            <Button variant="primary" type="submit" disabled={operating} onClick={(e) => updateSSL(e)}>
              {t('settings.sslFileSubmit')}
            </Button> &nbsp;
            <TutorialsButton prefixLine={true} tutorials={sslTutorials} /> &nbsp;
            {operating && <Spinner animation="border" variant="success" style={{verticalAlign: 'middle'}} />}
          </Form>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
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
        {showResult ? 'Reset' : 'Try it out'}
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

