import React from 'react';
import Container from "react-bootstrap/Container";
import './App.css'
import {Navbar, Nav} from 'react-bootstrap';
import {Link} from "react-router-dom";
import logo from './logo.svg';
import {Token} from "./utils";

export default function Navigator({initialized, expire}) {
  const [token, setToken] = React.useState();
  React.useEffect(() => {
    Token.load((data) => {
      setToken(data);
    });
  }, [initialized, expire]);

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
            <Nav.Link as={Link} to='/status'>Status</Nav.Link>
            <Nav.Link as={Link} to='/system'>System</Nav.Link>
            <Nav.Link as={Link} to='/software'>Software</Nav.Link>
            {token && <Nav.Link as={Link} to='/logout'>Logout</Nav.Link>}
          </Nav>
        }
      </Container>
    </Navbar>
  );
}

