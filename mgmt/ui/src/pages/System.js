import {useNavigate} from "react-router-dom";
import Container from "react-bootstrap/Container";
import React from "react";
import {Token, Errors} from "../utils";
import axios from "axios";
import {Row, Col, Card, Button, Form} from "react-bootstrap";
import UpgradeConfirmButton from '../components/UpgradeConfirmButton';
const semver = require('semver');

export default function System() {
  const navigate = useNavigate();
  const [status, setStatus] = React.useState();
  const [srs, setSRS] = React.useState();
  const [hooks, setHooks] = React.useState();
  const [prometheus, setPrometheus] = React.useState();
  const [nodeExporter, setNodeExporter] = React.useState();
  const [upgrading, setUpgrading] = React.useState();
  const [alreadyUpgrading, setAlreadyUpgrading] = React.useState();
  const [enableUpgrading, setEnableUpgrading] = React.useState();
  const [strategyAuto, setStrategyAuto] = React.useState();
  const [strategyChanged, setStrategyChanged] = React.useState();

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/mgmt/status', {
      ...token,
    }).then(res => {
      const status = res.data.data;
      if (status && status.releases && status.releases.latest) {
        if (semver.lt(status.version, status.releases.latest)) setEnableUpgrading(true);
      }
      if (status.upgrading) setAlreadyUpgrading(true);
      setStrategyAuto(status.strategy === 'auto');
      setStatus(status);
      console.log(`Status: Query ok, status=${JSON.stringify(status)}`);
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/routers-logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  }, [navigate, strategyChanged]);

  React.useEffect(() => {
    if (!upgrading || alreadyUpgrading) return;

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/upgrade', {
      ...token,
    }).then(res => {
      setUpgrading(false);
      console.log(`Status: Upgrade ok, status=${JSON.stringify(res.data.data)}`);
    }).catch(e => {
      if (e.response.status === 502) {
        alert(`升级完成，请刷新页面`);
      } else {
        alert(`未知错误, ${e.message}`);
      }
    });
  }, [upgrading, alreadyUpgrading]);

  const handleStrategyChange = () => {
    const token = Token.load();
    axios.post('/terraform/v1/mgmt/strategy', {
      ...token,
    }).then(res => {
      const status = res.data.data;
      setStrategyChanged(!strategyChanged);
      console.log(`Strategy: Change ok, status=${JSON.stringify(status)}`);
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

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/mgmt/containers', {
      ...token, action: 'query',
    }).then(res => {
      const containers = res.data.data;
      containers.map(container => {
        if (container.name === 'srs-server') setSRS(container);
        if (container.name === 'srs-hooks') setHooks(container);
        if (container.name === 'prometheus') setPrometheus(container);
        if (container.name === 'node-exporter') setNodeExporter(container);
      });
      console.log(`SRS: Query ok, status=${JSON.stringify(containers)}`);
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/routers-logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  }, [navigate]);

  return (
    <>
      <Container>
        <Row>
          <Col xs lg={3}>
            <Card style={{ width: '18rem', marginTop: '16px' }}>
              <Card.Header>SRS服务器</Card.Header>
              <Card.Body>
                <Card.Text>
                  容器名：{srs?.name} <br/>
                  容器ID：{srs?.container?.ID} <br/>
                  状态：{srs?.container.State} {srs?.container.Status}
                </Card.Text>
                <div style={{display: 'inline-block'}}>
                  <Button className='disabled'>
                    重启
                  </Button> &nbsp;
                  <Button className='disabled'>
                    升级
                  </Button>
                </div>
              </Card.Body>
            </Card>
          </Col>
          <Col xs lg={3}>
            <Card style={{ width: '18rem', marginTop: '16px' }}>
              <Card.Header>SRS回调</Card.Header>
              <Card.Body>
                <Card.Text>
                  容器名：{hooks?.name} <br/>
                  容器ID：{hooks?.container?.ID} <br/>
                  状态：{hooks?.container.State} {hooks?.container.Status}
                </Card.Text>
                <div style={{display: 'inline-block'}}>
                  <Button className='disabled'>
                    重启
                  </Button> &nbsp;
                  <Button className='disabled'>
                    升级
                  </Button>
                </div>
              </Card.Body>
            </Card>
          </Col>
          <Col xs lg={3}>
            <Card style={{ width: '18rem', marginTop: '16px' }}>
              <Card.Header>Prometheus监控</Card.Header>
              <Card.Body>
                <Card.Text>
                  容器名：{prometheus?.name} <br/>
                  容器ID：{prometheus?.container?.ID} <br/>
                  状态：{prometheus?.container.State} {prometheus?.container.Status}
                </Card.Text>
                <div style={{display: 'inline-block'}}>
                  <Button className='disabled'>
                    重启
                  </Button> &nbsp;
                  <Button className='disabled'>
                    升级
                  </Button>
                </div>
              </Card.Body>
            </Card>
          </Col>
          <Col xs lg={3}>
            <Card style={{ width: '18rem', marginTop: '16px' }}>
              <Card.Header>NodeExporter(节点监控)</Card.Header>
              <Card.Body>
                <Card.Text>
                  容器名：{nodeExporter?.name} <br/>
                  容器ID：{nodeExporter?.container?.ID} <br/>
                  状态：{nodeExporter?.container.State} {nodeExporter?.container.Status}
                </Card.Text>
                <div style={{display: 'inline-block'}}>
                  <Button className='disabled'>
                    重启
                  </Button> &nbsp;
                  <Button className='disabled'>
                    升级
                  </Button>
                </div>
              </Card.Body>
            </Card>
          </Col>
          <Col xs lg={3}>
            <Card style={{ width: '18rem', marginTop: '16px' }}>
              <Card.Header>管理后台</Card.Header>
              <Card.Body>
                <Card.Text as={Col}>
                  你的版本: {status?.version} {alreadyUpgrading && '升级中...'} <br/>
                  稳定版本: {status?.releases?.stable} &nbsp;
                  <Form.Check
                    type='switch'
                    label='自动升级'
                    style={{display: 'inline-block'}}
                    title='是否自动升级到稳定版本'
                    defaultChecked={strategyAuto}
                    onClick={handleStrategyChange}
                  />
                  <br/>
                  最新版本: <a href='https://github.com/ossrs/srs/issues/2856#changelog' target='_blank' rel='noreferrer'>{status?.releases?.latest}</a>
                  <p></p>
                </Card.Text>
                {
                  !enableUpgrading
                  ? <Button className='disabled'>升级</Button>
                  : <UpgradeConfirmButton upgrading={upgrading} handleClick={() => setUpgrading(true)} text='升级'>
                    <p>
                      升级管理后台，并且可能造成
                      <span className='text-danger'><strong>系统不可用</strong></span>，
                      确认继续升级么？
                    </p>
                  </UpgradeConfirmButton>
                } &nbsp;
                {!enableUpgrading &&
                  <UpgradeConfirmButton handleClick={() => setEnableUpgrading(true)} text='强制升级' operator='开启强制升级'>
                    <p>
                      你目前已经是最新版本，
                      <span className='text-warning'>没有必要强制升级</span>，
                      确认继续强制升级么？
                    </p>
                  </UpgradeConfirmButton>
                }
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </Container>
    </>
  );
}

