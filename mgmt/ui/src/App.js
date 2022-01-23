import React from 'react';
import './App.css'
import Container from "react-bootstrap/Container";
// See https://react-bootstrap.github.io/components/navbar/
import {Navbar, Nav} from 'react-bootstrap';
// See https://reactrouter.com/docs/en/v6/getting-started/tutorial
import { BrowserRouter, Routes, Link, Route } from "react-router-dom";
import Footer from './Footer';
import logo from './logo.svg';

function MyNavbar() {
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
        <Nav className="me-auto">
          <Nav.Link as={Link} to='/'>Login</Nav.Link>
          <Nav.Link as={Link} to='/system'>System</Nav.Link>
          <Nav.Link as={Link} to='/components'>Components</Nav.Link>
        </Nav>
      </Container>
    </Navbar>
  );
}

function Navigator() {
  return (
    <BrowserRouter basename={window.PUBLIC_URL}>
      <MyNavbar />
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="*" element={<Login />} />
        <Route path="/login" element={<Login />} />
        <Route path="/system" element={<System />} />
        <Route path="/components" element={<Components />} />
      </Routes>
    </BrowserRouter>
  );
}

function Login() {
  return <Container>Login</Container>;
}

function System() {
  return <Container>System</Container>;
}

function Components() {
  return <Container>Components</Container>;
}

function App() {
  return (
    <>
      <Navigator />
      <Footer />
    </>
  );
}

export default App;
