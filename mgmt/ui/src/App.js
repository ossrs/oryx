import React from 'react';
import Container from "react-bootstrap/Container";
import axios from "axios";
import './App.css'
import {BrowserRouter, Routes, Route} from "react-router-dom";
import Footer from './Footer';
import Login from './Login';
import Logout from './Logout';
import Navigator from './Navigator';
import Init from './Init';

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
  const [tokenUpdated, setTokenUpdated] = React.useState();

  React.useEffect(() => {
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
        <Navigator initialized={initialized} tokenUpdated={tokenUpdated} />
        <Routes>
          {!initialized && <>
            <Route path="*" element={<Init onInit={()=>setInitialized(true)} />}/>
          </>}
          {initialized && <>
            <Route path="/" element={<Login />} />
            <Route path="/login" element={<Login onLogin={() => setTokenUpdated(!tokenUpdated)}/>}/>
            <Route path="/status" element={<Status/>}/>
            <Route path="/system" element={<System/>}/>
            <Route path="/software" element={<Software/>}/>
            <Route path="/logout" element={<Logout onLogout={() => setTokenUpdated(!tokenUpdated)} />}/>
          </>}
        </Routes>
      </BrowserRouter>
      <Footer />
    </>
  );
}

export default App;
