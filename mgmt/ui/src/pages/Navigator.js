import React from 'react';
import Container from "react-bootstrap/Container";
import {Navbar, Nav} from 'react-bootstrap';
import {Link} from "react-router-dom";
import logo from '../resources/logo.svg';

export default function Navigator({initialized, token}) {
  return (
    <Navbar bg='light' variant='light'>
      <Container>
        <Navbar.Brand>
          <img
            src={logo}
            width="64"
            height="30"
            className="d-inline-block align-top"
            alt="SRS Terraform"
          />
        </Navbar.Brand>
        {
          initialized &&
          <Nav className="me-auto">
            {token && <>
              <Nav.Link as={Link} to='/routers-dashboard'>仪表盘</Nav.Link>
              <Nav.Link as={Link} to='/routers-scenario'>应用场景</Nav.Link>
              <Nav.Link as={Link} to='/routers-config'>系统配置</Nav.Link>
              <Nav.Link as={Link} to='/routers-system'>组件管理</Nav.Link>
            </>}
            {!token && <Nav.Link as={Link} to='/routers-login'>登录</Nav.Link>}
            {token && <Nav.Link as={Link} to='/routers-logout'>退出</Nav.Link>}
          </Nav>
        }
      </Container>
    </Navbar>
  );
}

