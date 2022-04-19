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
  const [hooks, setHooks] = React.useState();
  const [tencent, setTencent] = React.useState();
  const [ffmpeg, setFFmpeg] = React.useState();
  const [platform, setPlatform] = React.useState();
  const [prometheus, setPrometheus] = React.useState();
  const [nodeExporter, setNodeExporter] = React.useState();
  const [redisServer, setRedisServer] = React.useState();
  const [strategyAutoUpgrade, setStrategyAutoUpgrade] = React.useState();
  const [userToggleStrategy, setUserToggleStrategy] = React.useState();
  const [searchParams] = useSearchParams();
  const [allowManuallyUpgrade, setAllowManuallyUpgrade] = React.useState();
  const [allowDisableContainer, setAllowDisableContainer] = React.useState();
  const [refreshContainers, setRefreshContainers] = React.useState();
  const [allowSwitchContainer, setAllowSwitchContainer] = React.useState();
  const handleError = useErrorHandler();
  const {t} = useTranslation();

  React.useEffect(() => {
    const allowManuallyUpgrade = searchParams.get('allow-manual') === 'true';
    console.log(`?allow-manual=true|false, current=${allowManuallyUpgrade}, Whether allow manually upgrade`);
    setAllowManuallyUpgrade(allowManuallyUpgrade);
  }, [searchParams]);

  React.useEffect(() => {
    const allowDisableContainer = searchParams.get('allow-disable') === 'true';
    console.log(`?allow-disable=true|false, current=${allowDisableContainer}, Whether allow disable container`);
    setAllowDisableContainer(allowDisableContainer);
  }, [searchParams]);

  React.useEffect(() => {
    const allowSwitchContainer = searchParams.get('allow-switch') === 'true';
    console.log(`?allow-switch=true|false, current=${allowSwitchContainer}, Whether allow switch srs server`);
    setAllowSwitchContainer(allowSwitchContainer);
  }, [searchParams]);

  // Because the onStatus always change during rendering, so we use a callback so that the useEffect() could depends on
  // it to avoid infinitely loops. That is callback is not changed, while onStatus changed(not null) mnay times during
  // each rendering of components.
  const onStatus = React.useCallback((status) => {
    setStrategyAutoUpgrade(status.strategy === 'auto');
    setStatus(status);
  }, []);

  const handleUpgradeStrategyChange = React.useCallback((e) => {
    if (strategyAutoUpgrade && !window.confirm(t('coms.disableUpgrade'))) {
      e.preventDefault();
      return;
    }

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/strategy', {
      ...token,
    }).then(res => {
      setUserToggleStrategy(!userToggleStrategy);
      console.log(`Strategy: Change ok`);
    }).catch(handleError);
  }, [handleError, strategyAutoUpgrade, userToggleStrategy, t]);

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

        if (m.name === 'srs-server') setSrsRelease(m);
        if (m.name === 'srs-dev') setSrsDev(m);
        if (m.name === 'srs-hooks') setHooks(m);
        if (m.name === 'tencent-cloud') setTencent(m);
        if (m.name === 'ffmpeg') setFFmpeg(m);
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
                  <Button className='disabled'>{t('helper.upgrade')}</Button> &nbsp;
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
                  <Button className='disabled'>
                    {t('helper.upgrade')}
                  </Button> &nbsp;
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
              <Card.Header>{t('coms.hooks')}</Card.Header>
              <Card.Body>
                <Card.Text as={Col}>
                  {t('coms.containerName')}：{hooks?.name} <br/>
                  {t('coms.containerId')}：{hooks?.container?.ID} <br/>
                  {t('coms.containerState')}：{hooks?.StatusMessage}
                  <p></p>
                </Card.Text>
                <div style={{display: 'inline-block'}}>
                  <Button className='disabled'>
                    {t('helper.restart')}
                  </Button> &nbsp;
                  <Button className='disabled'>
                    {t('helper.upgrade')}
                  </Button> &nbsp;
                  <MgmtUpdateContainer
                    allow={allowDisableContainer && hooks?.name}
                    enabled={hooks?.enabled}
                    onClick={() => handleContainerChange(hooks)}
                  />
                </div>
              </Card.Body>
            </Card>
          </Col>
          <Col xs lg={3}>
            <Card style={{ width: '18rem', marginTop: '16px' }}>
              <Card.Header>{t('coms.ffmpeg')}</Card.Header>
              <Card.Body>
                <Card.Text as={Col}>
                  {t('coms.containerName')}：{ffmpeg?.name} <br/>
                  {t('coms.containerId')}：{ffmpeg?.container?.ID} <br/>
                  {t('coms.containerState')}：{ffmpeg?.StatusMessage}
                  <p></p>
                </Card.Text>
                <div style={{display: 'inline-block'}}>
                  <Button className='disabled'>
                    {t('helper.restart')}
                  </Button> &nbsp;
                  <Button className='disabled'>
                    {t('helper.upgrade')}
                  </Button> &nbsp;
                  <MgmtUpdateContainer
                    allow={allowDisableContainer && ffmpeg?.name}
                    enabled={ffmpeg?.enabled}
                    onClick={() => handleContainerChange(ffmpeg)}
                  />
                </div>
              </Card.Body>
            </Card>
          </Col>
          <Col xs lg={3}>
            <Card style={{ width: '18rem', marginTop: '16px' }}>
              <Card.Header>{t('coms.tencent')}</Card.Header>
              <Card.Body>
                <Card.Text as={Col}>
                  {t('coms.containerName')}：{tencent?.name} <br/>
                  {t('coms.containerId')}：{tencent?.container?.ID} <br/>
                  {t('coms.containerState')}：{tencent?.StatusMessage}
                  <p></p>
                </Card.Text>
                <div style={{display: 'inline-block'}}>
                  <Button className='disabled'>
                    {t('helper.restart')}
                  </Button> &nbsp;
                  <Button className='disabled'>
                    {t('helper.upgrade')}
                  </Button> &nbsp;
                  <MgmtUpdateContainer
                    allow={allowDisableContainer && tencent?.name}
                    enabled={tencent?.enabled}
                    onClick={() => handleContainerChange(tencent)}
                  />
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
                  <Button className='disabled'>
                    {t('helper.upgrade')}
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
                  <Button className='disabled'>
                    {t('helper.upgrade')}
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
                  <Button className='disabled'>
                    {t('helper.upgrade')}
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
                  <Button className='disabled'>
                    {t('helper.upgrade')}
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
                  {t('coms.version')}: {status?.version} <br/>
                  {t('coms.stable')}: {status?.releases?.stable} &nbsp;
                  <Form.Check
                    type='switch'
                    label={t('coms.autoUpgrade')}
                    style={{display: 'inline-block'}}
                    title={t('coms.autoUpgradeTip')}
                    disabled={!allowManuallyUpgrade}
                    defaultChecked={strategyAutoUpgrade}
                    onClick={(e) => handleUpgradeStrategyChange(e)}
                  />
                  <br/>
                  {t('coms.latest')}: <a href='https://github.com/ossrs/srs/issues/2856#changelog' target='_blank' rel='noreferrer'>{status?.releases?.latest}</a>
                  <p></p>
                </Card.Text>
                <MgmtUpgradeButton onStatus={onStatus}/>
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

const upgradeProgress = 300;

function MgmtUpgradeButton({onStatus}) {
  const [startingUpgrade, setStartingUpgrade] = React.useState();
  const [isUpgrading, setIsUpgrading] = React.useState();
  const [releaseAvailable, setReleaseAvailable] = React.useState();
  const [upgradeDone, setUpgradeDone] = React.useState();
  const [progress, setProgress] = React.useState(upgradeProgress);
  const {t} = useTranslation();

  // For callback to use state.
  const ref = React.useRef({});
  React.useEffect(() => {
    ref.current.startingUpgrade = startingUpgrade;
    ref.current.progress = progress;
    ref.current.upgradeDone = upgradeDone;
  }, [startingUpgrade, progress, upgradeDone]);

  React.useEffect(() => {
    const refreshMgmtStatus = () => {
      const token = Token.load();
      axios.post('/terraform/v1/mgmt/status', {
        ...token,
      }).then(res => {
        const status = res.data.data;

        // Normally state.
        setIsUpgrading(status.upgrading);
        onStatus(status);

        // Whether upgrade is available.
        if (status && status.releases && status.releases.latest) {
          setReleaseAvailable(semver.lt(status.version, status.releases.latest));
        }

        // If upgradeDone is false, we're in the upgrading progress, so it's done when upgrading changed to false.
        if (ref.current.upgradeDone === false && !status.upgrading && ref.current.progress < upgradeProgress) {
          setStartingUpgrade(false);
          setUpgradeDone(true);
        }

        // If state not set, but already upgrading, it's restore from the previous state.
        if (status.upgrading && ref.current.upgradeDone === undefined && ref.current.startingUpgrade === undefined) {
          setStartingUpgrade(true);
          setUpgradeDone(false);
        }

        console.log(`${moment().format()}: Status: Query ok, startingUpgrade=${ref.current.startingUpgrade}, upgradeDone=${ref.current.upgradeDone}, status=${JSON.stringify(status)}`);
      }).catch(e => {
        console.log('ignore any error during status', e);
      });
    };

    refreshMgmtStatus();
    const timer = setInterval(() => refreshMgmtStatus(), 10 * 1000);
    return () => clearInterval(timer);
  }, [startingUpgrade, onStatus]);

  const handleStartUpgrade = () => {
    if (isUpgrading) return;

    setUpgradeDone(false);
    setStartingUpgrade(true);
    setProgress(upgradeProgress);

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/upgrade', {
      ...token,
    }).then(res => {
      console.log(`upgrade ok, ${JSON.stringify(res.data.data)}`);
    }).catch(e => {
      console.log('ignore any error during upgrade', e);
    });
  };

  React.useEffect(() => {
    if (!isUpgrading) return;
    const timer = setInterval(() => {
      if (ref.current.progress <= 0) return;
      if (ref.current.progress <= 10) setUpgradeDone(true);
      setProgress(ref.current.progress - 1);
    }, 1000);
    return () => clearInterval(timer);
  }, [isUpgrading]);

  React.useEffect(() => {
    if (!upgradeDone) return;
    alert(t('coms.upgradeOk'));
  }, [upgradeDone, t]);

  return (
    <UpgradeConfirmButton
      releaseAvailable={releaseAvailable}
      upgrading={startingUpgrade || isUpgrading}
      progress={`${progress}s`}
      onClick={handleStartUpgrade}
      text={t('helper.upgrade')}
    >
      <p>
        {t('coms.upgradeTip1')}
        <span className='text-danger'><strong>
          {t('coms.upgradeTip2')}
        </strong></span>
        {t('coms.upgradeTip3')}
      </p>
    </UpgradeConfirmButton>
  );
}

