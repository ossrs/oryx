import {useSearchParams} from "react-router-dom";
import Container from "react-bootstrap/Container";
import React from "react";
import {Token} from "../utils";
import axios from "axios";
import {Row, Col, Card, Button, Form} from "react-bootstrap";
import UpgradeConfirmButton from '../components/UpgradeConfirmButton';
import SwitchConfirmButton from '../components/SwitchConfirmButton';
import * as semver from 'semver';
import * as moment from 'moment';
import {SrsErrorBoundary} from "../components/SrsErrorBoundary";
import {useErrorHandler} from "react-error-boundary";
import {useTranslation} from "react-i18next";

export default function Components() {
  return (
    <SrsErrorBoundary>
      <ComponentsImpl />
    </SrsErrorBoundary>
  );
}

function ComponentsImpl() {
  const [status, setStatus] = React.useState();
  const [srsRelease, setSrsRelease] = React.useState();
  const [srsDev, setSrsDev] = React.useState();
  const [platform, setPlatform] = React.useState();
  const [prometheus, setPrometheus] = React.useState();
  const [nodeExporter, setNodeExporter] = React.useState();
  const [redisServer, setRedisServer] = React.useState();
  const [searchParams] = useSearchParams();
  const [allowDisableContainer, setAllowDisableContainer] = React.useState();
  const [refreshContainers, setRefreshContainers] = React.useState();
  const [allowSwitchContainer, setAllowSwitchContainer] = React.useState(true);
  const handleError = useErrorHandler();
  const {t} = useTranslation();

  React.useEffect(() => {
    const allowDisableContainer = searchParams.get('allow-disable') === 'true';
    console.log(`?allow-disable=true|false, current=${allowDisableContainer}, Whether allow disable container`);
    if (!searchParams.get('allow-disable')) return; // Ignore if not specified.
    setAllowDisableContainer(allowDisableContainer);
  }, [searchParams]);

  React.useEffect(() => {
    const allowSwitchContainer = searchParams.get('allow-switch') === 'true';
    console.log(`?allow-switch=true|false, current=${allowSwitchContainer}, Whether allow switch srs server`);
    if (!searchParams.get('allow-switch')) return; // Ignore if not specified.
    setAllowSwitchContainer(allowSwitchContainer);
  }, [searchParams]);

  React.useEffect(() => {
    const refreshMgmtStatus = () => {
      const token = Token.load();
      axios.post('/terraform/v1/mgmt/status', {
        ...token,
      }).then(res => {
        const status = res.data.data;

        // Normally state.
        setStatus(status);

        console.log(`${moment().format()}: Status: Query ok, status=${JSON.stringify(status)}`);
      }).catch(e => {
        console.log('ignore any error during status', e);
      });
    };

    refreshMgmtStatus();
    const timer = setInterval(() => refreshMgmtStatus(), 10 * 1000);
    return () => clearInterval(timer);
  }, [setStatus]);

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/mgmt/containers', {
      ...token, action: 'query',
    }).then(res => {
      const containers = res.data.data;
      containers.filter(m => {
        if (m.container.State || m.container.Status) {
          m.StatusMessage = `${m.container.State || ''} ${m.container.Status || ''}`.trim();
        } else {
          m.StatusMessage = 'Stopped';
        }
        if (!m.enabled) {
          m.StatusMessage = 'Disabled';
        }

        if (m.name === 'srs-server') setSrsRelease(m);
        if (m.name === 'srs-dev') setSrsDev(m);
        if (m.name === 'platform') setPlatform(m);
        if (m.name === 'prometheus') setPrometheus(m);
        if (m.name === 'node-exporter') setNodeExporter(m);
        if (m.name === 'redis') setRedisServer(m);

        return null;
      });
      console.log(`SRS: Query ok, containers are ${JSON.stringify(containers)}`);
    }).catch(handleError);
  }, [refreshContainers, handleError]);

  const handleContainerChange = React.useCallback((container) => {
    const token = Token.load();
    axios.post('/terraform/v1/mgmt/containers', {
      ...token, action: 'enabled', name: container.name, enabled: !container.enabled,
    }).then(res => {
      console.log(`SRS: Update ok, enabled=${!container.enabled}`);
      setRefreshContainers(Math.random());
    }).catch(handleError);
  }, [handleError]);

  const handleSwitch = React.useCallback((container) => {
    const token = Token.load();
    axios.post('/terraform/v1/mgmt/containers', {
      ...token, action: 'switch', name: container.name,
    }).then(res => {
      console.log(`SRS: Switch ok, name=${container.name}`);
      setRefreshContainers(Math.random());
    }).catch(handleError);
  }, [handleError]);

  return (
    <>
      <Container>
        <Row>
          <Col xs lg={3}>
            <Card style={{ width: '18rem', marginTop: '16px' }}>
              <Card.Header>{t('coms.srs4')}</Card.Header>
              <Card.Body>
                <Card.Text as={Col}>
                  {t('coms.containerName')}：{srsRelease?.name} <br/>
                  {t('coms.containerId')}：{srsRelease?.container?.ID ? srsRelease.container.ID : 'No Container'} <br/>
                  {t('coms.containerState')}：{srsRelease?.StatusMessage}
                  <p></p>
                </Card.Text>
                <div style={{display: 'inline-block'}}>
                  {srsDev?.enabled || <>
                    <Button className='disabled'>{t('helper.restart')}</Button> &nbsp;
                  </>}
                  <MgmtUpdateContainer
                    allow={allowDisableContainer && srsRelease?.name}
                    enabled={srsRelease?.enabled}
                    onClick={() => handleContainerChange(srsRelease)}
                  /> &nbsp;
                  <SwitchConfirmButton
                    enabled={srsDev?.enabled}
                    onClick={() => handleSwitch(srsRelease)}
                    allowSwitchContainer={allowSwitchContainer}
                  >
                    <p>{t('coms.switchConfirm')}</p>
                  </SwitchConfirmButton>
                </div>
              </Card.Body>
            </Card>
          </Col>
          <Col xs lg={3}>
            <Card style={{ width: '18rem', marginTop: '16px' }}>
              <Card.Header>{t('coms.srs5')}</Card.Header>
              <Card.Body>
                <Card.Text as={Col}>
                  {t('coms.containerName')}：{srsDev?.name} <br/>
                  {t('coms.containerId')}：{srsDev?.container?.ID ? srsDev.container.ID : 'No Container'} <br/>
                  {t('coms.containerState')}：{srsDev?.StatusMessage}
                  <p></p>
                </Card.Text>
                <div style={{display: 'inline-block'}}>
                  {srsRelease?.enabled || <>
                    <Button className='disabled'>{t('helper.restart')}</Button> &nbsp;
                  </>}
                  <MgmtUpdateContainer
                    allow={allowDisableContainer && srsDev?.name}
                    enabled={srsDev?.enabled}
                    onClick={() => handleContainerChange(srsDev)}
                  /> &nbsp;
                  <SwitchConfirmButton
                    enabled={srsRelease?.enabled}
                    onClick={() => handleSwitch(srsDev)}
                    allowSwitchContainer={allowSwitchContainer}
                  >
                    <p>
                      {t('coms.switchConfirm1')}
                      <font color='red'>{t('coms.switchConfirm2')}</font>
                      {t('coms.switchConfirm3')}
                    </p>
                  </SwitchConfirmButton>
                </div>
              </Card.Body>
            </Card>
          </Col>
          <Col xs lg={3}>
            <Card style={{ width: '18rem', marginTop: '16px' }}>
              <Card.Header>{t('coms.prometheus')}</Card.Header>
              <Card.Body>
                <Card.Text as={Col}>
                  {t('coms.containerName')}：{prometheus?.name} <br/>
                  {t('coms.containerId')}：{prometheus?.container?.ID} <br/>
                  {t('coms.containerState')}：{prometheus?.StatusMessage}
                  <p></p>
                </Card.Text>
                <div style={{display: 'inline-block'}}>
                  <Button className='disabled'>
                    {t('helper.restart')}
                  </Button> &nbsp;
                  <MgmtUpdateContainer
                    allow={allowDisableContainer && prometheus?.name}
                    enabled={prometheus?.enabled}
                    onClick={() => handleContainerChange(prometheus)}
                  />
                </div>
              </Card.Body>
            </Card>
          </Col>
          <Col xs lg={3}>
            <Card style={{ width: '18rem', marginTop: '16px' }}>
              <Card.Header>{t('coms.node')}</Card.Header>
              <Card.Body>
                <Card.Text as={Col}>
                  {t('coms.containerName')}：{nodeExporter?.name} <br/>
                  {t('coms.containerId')}：{nodeExporter?.container?.ID} <br/>
                  {t('coms.containerState')}：{nodeExporter?.StatusMessage}
                  <p></p>
                </Card.Text>
                <div style={{display: 'inline-block'}}>
                  <Button className='disabled'>
                    {t('helper.restart')}
                  </Button> &nbsp;
                  <MgmtUpdateContainer
                    allow={allowDisableContainer && nodeExporter?.name}
                    enabled={nodeExporter?.enabled}
                    onClick={() => handleContainerChange(nodeExporter)}
                  />
                </div>
              </Card.Body>
            </Card>
          </Col>
          <Col xs lg={3}>
            <Card style={{ width: '18rem', marginTop: '16px' }}>
              <Card.Header>{t('coms.redis')}</Card.Header>
              <Card.Body>
                <Card.Text as={Col}>
                  {t('coms.containerName')}：{redisServer?.name} <br/>
                  {t('coms.containerId')}：{redisServer?.container?.ID} <br/>
                  {t('coms.containerState')}：{redisServer?.StatusMessage}
                  <p></p>
                </Card.Text>
                <div style={{display: 'inline-block'}}>
                  <Button className='disabled'>
                    {t('helper.restart')}
                  </Button> &nbsp;
                  <MgmtUpdateContainer
                    allow={allowDisableContainer && redisServer?.name}
                    enabled={redisServer?.enabled}
                    onClick={() => handleContainerChange(redisServer)}
                  />
                </div>
              </Card.Body>
            </Card>
          </Col>
          <Col xs lg={3}>
            <Card style={{ width: '18rem', marginTop: '16px' }}>
              <Card.Header>{t('coms.platform')}</Card.Header>
              <Card.Body>
                <Card.Text as={Col}>
                  {t('coms.containerName')}：{platform?.name} <br/>
                  {t('coms.containerId')}：{platform?.container?.ID} <br/>
                  {t('coms.containerState')}：{platform?.StatusMessage}
                  <p></p>
                </Card.Text>
                <div style={{display: 'inline-block'}}>
                  <Button className='disabled'>
                    {t('helper.restart')}
                  </Button> &nbsp;
                  <MgmtUpdateContainer
                    allow={allowDisableContainer && platform?.name}
                    enabled={platform?.enabled}
                    onClick={() => handleContainerChange(platform)}
                  />
                </div>
              </Card.Body>
            </Card>
          </Col>
          <Col xs lg={3}>
            <Card style={{ width: '18rem', marginTop: '16px' }}>
              <Card.Header>{t('coms.host')}</Card.Header>
              <Card.Body>
                <Card.Text as={Col}>
                  <p>
                  {t('coms.version')}: {status?.version} <br/>
                  {t('coms.stable')}: {status?.releases?.stable}<br/>
                  {t('coms.latest')}: <a href='https://github.com/ossrs/srs/issues/2856#changelog' target='_blank' rel='noreferrer'>{status?.releases?.latest}</a>
                  </p>
                  <footer className="blockquote-footer">
                    {t('coms.upgradeManually')}
                  </footer>
                </Card.Text>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </Container>
    </>
  );
}

function MgmtUpdateContainer({allow, enabled, onClick}) {
  const {t} = useTranslation();

  const handleClick = (e) => {
    if (enabled && !window.confirm(t('coms.disableContainer'))) {
      e.preventDefault();
      return;
    }
    onClick();
  };

  return (
    <Button
      className={allow ? '' : 'disabled'}
      variant={enabled ? 'danger' : 'success'}
      onClick={(e) => handleClick(e)}
    >
      {enabled ? t('helper.disable') : t('helper.enable')}
    </Button>
  );
}

