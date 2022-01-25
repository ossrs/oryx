import {useNavigate} from "react-router-dom";
import Container from "react-bootstrap/Container";
import React from "react";
import {Token, Errors} from "./utils";
import axios from "axios";
import {Row, Col, Card, Button} from "react-bootstrap";
import PopoverConfirmButton from './PopoverConfirmButton';

export default function System() {
  const navigate = useNavigate();
  const [status, setStatus] = React.useState();
  const [srs, setSRS] = React.useState();
  const [upgrading, setUpgrading] = React.useState();

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/mgmt/status', {
      ...token,
    }).then(res => {
      setStatus(res.data.data);
      console.log(`Status: Query ok, status=${JSON.stringify(res.data.data)}`);
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  }, [navigate]);

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/mgmt/srs', {
      ...token,
    }).then(res => {
      setSRS(res.data.data);
      console.log(`SRS: Query ok, status=${JSON.stringify(res.data.data)}`);
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  }, [navigate]);

  React.useEffect(() => {
    if (!upgrading) return;

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
  }, [upgrading]);

  return (
    <>
      <p></p>
      <Container>
        <Row>
          <Col xs lg={3}>
            <Card style={{ width: '18rem' }}>
              <Card.Header>SRS Server</Card.Header>
              <Card.Body>
                <Card.Text>
                  {srs?.major} {srs?.container.State} {srs?.container.Status}
                </Card.Text>
                <Button className='disabled'>
                  升级SRS服务器
                </Button>
              </Card.Body>
            </Card>
          </Col>
          <Col xs lg={3}>
            <Card style={{ width: '18rem' }}>
              <Card.Header>管理后台</Card.Header>
              <Card.Body>
                <Card.Text>
                  Current: {status?.version} <br/>
                  Stable: {status?.releases?.versions?.stable} <br/>
                  Latest: {status?.releases?.versions?.latest}
                </Card.Text>
                <PopoverConfirmButton upgrading={upgrading} handleClick={() => setUpgrading(true)} text='升级管理后台'>
                  <p>
                    升级管理后台，并且可能造成
                    <span className='text-danger'><strong>系统不可用</strong></span>，
                    确认继续升级么？
                  </p>
                </PopoverConfirmButton>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </Container>
    </>
  );
}

