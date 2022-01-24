import {useNavigate} from "react-router-dom";
import Container from "react-bootstrap/Container";
import React from "react";
import {Token, Errors} from "./utils";
import axios from "axios";
import {Button, Spinner, Form, Row, Col, Card} from "react-bootstrap";

export default function Status() {
  const navigate = useNavigate();
  const [status, setStatus] = React.useState();
  const [upgrading, setUpgrading] = React.useState();

  // Verify the token if token changed.
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
  }, [upgrading]);

  React.useEffect(() => {
    if (!upgrading) return;

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/upgrade', {
      ...token,
    }).then(res => {
      setUpgrading(false);
      console.log(`Status: Upgrade ok, status=${JSON.stringify(res.data.data)}`);
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  }, [upgrading]);

  const handleClick = () => {
    setUpgrading(true);
  };

  return (
    <>
      <p></p>
      <Container>
        <Card style={{ width: '18rem' }}>
          <Card.Header>System</Card.Header>
          <Card.Body>
            <Card.Title>Version</Card.Title>
            <Card.Text>
              {status?.version}
            </Card.Text>
            <div className='row row-cols-lg-auto g-3 align-items-center'>
              <div className="col-12">
                <Button
                  variant="primary"
                  disabled={upgrading}
                  onClick={!upgrading ? handleClick : null}
                >
                  {upgrading ? '正在升级中...' : '升级管理后台'}
                </Button>
              </div>
              <div className="col-12">
                {upgrading && <Spinner animation="border" />}
              </div>
            </div>
          </Card.Body>
        </Card>
      </Container>
    </>
  );
}

