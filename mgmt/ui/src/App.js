import React from 'react';
import './App.css'
import Container from "react-bootstrap/Container";
import {Navbar, Nav} from 'react-bootstrap';
import { BrowserRouter, Routes, Link, Route } from "react-router-dom";
import Footer from './Footer';
import Login from './Login';
import logo from './logo.svg';
import axios from "axios";

function Navigator({initialized}) {
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
            <Nav.Link as={Link} to='/logout'>Logout</Nav.Link>
          </Nav>
        }
      </Container>
    </Navbar>
  );
}

function Status() {
  return <Container>Status</Container>;
}

function System() {
  return <Container>System</Container>;
}

function Software() {
  return <Container>Software</Container>;
}

function Logout() {
  return <Container>Logout</Container>;
}

function App() {
  const [initialized, setInitialized] = React.useState();

  React.useEffect(() => {
    console.log('xxx', initialized);
    axios.get('/terraform/v1/mgmt/init').then(res => {
      setInitialized(res.data.data.init);
    }).catch(e => {
      alert(e.response.data);
      console.error(e);
    });
  }, []);

  return (
    <>
      <BrowserRouter basename={window.PUBLIC_URL}>
        <Navigator initialized={initialized} />
        <Routes>
          <Route path="/" element={<Login initialized={initialized} onLoginSuccess={() => setInitialized(true)}/>}/>
          <Route path="/status" element={<Status/>}/>
          <Route path="/system" element={<System/>}/>
          <Route path="/software" element={<Software/>}/>
          <Route path="/logout" element={<Logout/>}/>
        </Routes>
      </BrowserRouter>
      <Footer />
    </>
  );
}

export default App;
