import React from 'react';
import Container from "react-bootstrap/Container";
import axios from "axios";
import './App.css'
import {BrowserRouter, Routes, Route} from "react-router-dom";
import Footer from './Footer';
import Login from './Login';
import Logout from './Logout';
import Navigator from './Navigator';

function Status() {
  return <Container>Status</Container>;
}

function System() {
  return <Container>System</Container>;
}

function Software() {
  return <Container>Software</Container>;
}

function App() {
  const [initialized, setInitialized] = React.useState();
  const [expire, setExpire] = React.useState();

  React.useEffect(() => {
    axios.get('/terraform/v1/mgmt/init').then(res => {
      setInitialized(res.data.data.init);
    }).catch(e => {
      alert(e.response.data);
      console.error(e);
    });
  }, []);

  const onLogin = () => {
    setInitialized(true);
    setExpire(!expire);
  };
  const onLogout = () => {
    setExpire(!expire);
  };

  return (
    <>
      <BrowserRouter basename={window.PUBLIC_URL}>
        <Navigator initialized={initialized} expire={expire} />
        <Routes>
          <Route path="/" element={<Login initialized={initialized} onLogin={onLogin}/>}/>
          <Route path="/status" element={<Status/>}/>
          <Route path="/system" element={<System/>}/>
          <Route path="/software" element={<Software/>}/>
          <Route path="/logout" element={<Logout onLogout={onLogout} />}/>
        </Routes>
      </BrowserRouter>
      <Footer />
    </>
  );
}

export default App;
