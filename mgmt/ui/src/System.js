import {useNavigate} from "react-router-dom";
import Container from "react-bootstrap/Container";
import React from "react";
import {Token, Errors} from "./utils";
import axios from "axios";
import {Button, Spinner, Card, OverlayTrigger, Popover} from "react-bootstrap";

function PopoverConfirmButton({upgrading, handleClick, text, children}) {
  const [showUpgrading, setShowUpgrading] = React.useState();

  const onHandleClick = () => {
    setShowUpgrading(false);
    handleClick();
  };

  const popover = (
    <Popover id="popover-basic">
      <Popover.Header as="h3">Confirm</Popover.Header>
      <Popover.Body>
        <p>
          {children}
        </p>
        <div className='row row-cols-lg-auto g-3 align-items-center'>
          <div className="col-12">
            <Button
              variant="danger"
              disabled={upgrading}
              onClick={!upgrading ? onHandleClick : null}
            >
              确认升级
            </Button>
          </div>
          <div className="col-12">
            <Button
              variant="primary"
              onClick={() => setShowUpgrading(false)}
            >
              取消
            </Button>
          </div>
        </div>
      </Popover.Body>
    </Popover>
  );

  return (
    <div className='row row-cols-lg-auto g-3 align-items-center'>
      <div className="col-12">
        <OverlayTrigger trigger="click" placement="right" overlay={popover} show={showUpgrading}>
          <Button variant="primary" onClick={() => setShowUpgrading(!showUpgrading)}>
            {upgrading ? '正在升级中...' : text}
          </Button>
        </OverlayTrigger>
      </div>
      <div className="col-12">
        {upgrading && <Spinner animation="border" />}
      </div>
    </div>
  );
}

export default function System() {
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
      if (e.response.status === 502) {
        alert(`升级完成，请刷新页面`);
      } else {
        alert(`未知错误, ${e.message}`);
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
          <Card.Header>管理后台</Card.Header>
          <Card.Body>
            <Card.Title>当前版本</Card.Title>
            <Card.Text>
              {status?.version}
            </Card.Text>
            <PopoverConfirmButton upgrading={upgrading} handleClick={handleClick} text='升级管理后台'>
                升级管理后台，需要较长时间（1分钟左右），并且可能造成<span className='text-danger'><strong>系统不可用</strong></span>，
                确认继续升级么？
            </PopoverConfirmButton>
          </Card.Body>
        </Card>
      </Container>
    </>
  );
}

