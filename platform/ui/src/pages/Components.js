import {useSearchParams} from "react-router-dom";
import Container from "react-bootstrap/Container";
import React from "react";
import {Token} from "../utils";
import axios from "axios";
import {Row, Col, Card, Button} from "react-bootstrap";
import UpgradeConfirmButton from '../components/UpgradeConfirmButton';
import * as moment from 'moment';
import {SrsErrorBoundary} from "../components/SrsErrorBoundary";
import {useErrorHandler} from "react-error-boundary";
import {useTranslation} from "react-i18next";
import * as semver from "semver";

export default function Components() {
  return (
    <SrsErrorBoundary>
      <ComponentsImpl />
    </SrsErrorBoundary>
  );
}

function ComponentsImpl() {
  const [status, setStatus] = React.useState();
  const [platform, setPlatform] = React.useState();
  const [searchParams] = useSearchParams();
  const [allowDisableContainer, setAllowDisableContainer] = React.useState();
  const [refreshContainers, setRefreshContainers] = React.useState();
  const handleError = useErrorHandler();
  const {t} = useTranslation();

  React.useEffect(() => {
    const allowDisableContainer = searchParams.get('allow-disable') === 'true';
    console.log(`?allow-disable=true|false, current=${allowDisableContainer}, Whether allow disable container`);
    if (!searchParams.get('allow-disable')) return; // Ignore if not specified.
    setAllowDisableContainer(allowDisableContainer);
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

        if (m.name === 'platform') setPlatform(m);

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

  // Because the onStatus always change during rendering, so we use a callback so that the useEffect() could depends on
  // it to avoid infinitely loops. That is callback is not changed, while onStatus changed(not null) mnay times during
  // each rendering of components.
  const onStatus = React.useCallback((status) => {
    setStatus(status);
  }, []);

  return (
    <>
      <Container>
        <Row>
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
                  {t('coms.latest')}: <a href='https://github.com/ossrs/srs-cloud/issues/4#changelog' target='_blank' rel='noreferrer'>{status?.releases?.latest}</a>
                  </p>
                  {status?.upgrading === undefined &&
                    <footer className="blockquote-footer">
                      {t('coms.upgradeManually')}
                    </footer>
                  }
                </Card.Text> &nbsp;
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
  const [requestStatus, setRequestStatus] = React.useState(1);
  const [isUpgrading, setIsUpgrading] = React.useState();
  const [releaseAvailable, setReleaseAvailable] = React.useState();
  const [upgradeDone, setUpgradeDone] = React.useState();
  const [progress, setProgress] = React.useState(upgradeProgress);
  const {t} = useTranslation();

  // For callback to update state, because in callback we can only get the copy, so we need a ref to point to the latest
  // copy of state of variant objects.
  const ref = React.useRef({});
  React.useEffect(() => {
    ref.current.startingUpgrade = startingUpgrade;
    ref.current.progress = progress;
    ref.current.upgradeDone = upgradeDone;
    ref.current.requestStatus = requestStatus;
  }, [startingUpgrade, progress, upgradeDone, requestStatus]);

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
    const timeout = startingUpgrade ? 1.3 * 1000 : 8.5 * 1000;
    const timer = setInterval(() => refreshMgmtStatus(), timeout);
    return () => clearInterval(timer);
  }, [startingUpgrade, requestStatus, onStatus]);

  const handleStartUpgrade =() => {
    if (isUpgrading) return;

    setUpgradeDone(false);
    setStartingUpgrade(true);
    setProgress(upgradeProgress);
    setTimeout(() => {
      setRequestStatus(ref.current.requestStatus + 1);
    }, 300);

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

