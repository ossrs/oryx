import React from 'react';
import axios from "axios";
import './App.css'
import {BrowserRouter, Routes, Route} from "react-router-dom";
import Footer from './pages/Footer';
import Login from './pages/Login';
import Logout from './pages/Logout';
import Navigator from './pages/Navigator';
import Setup from './pages/Setup';
import {Token} from "./utils";
import System from "./pages/System";
import Scenario from "./pages/Scenario";
import {Container} from "react-bootstrap";
import Config from "./pages/Config";
import Dashboard from './pages/Dashboard';

function App() {
  const [loading, setLoading] = React.useState(true);
  const [initialized, setInitialized] = React.useState();
  const [tokenUpdated, setTokenUpdated] = React.useState();
  const [token, setToken] = React.useState();

  React.useEffect(() => {
    axios.get('/terraform/v1/mgmt/init').then(res => {
      setInitialized(res.data.data.init);
    }).catch(e => {
      alert(e.response.data);
      console.error(e);
    }).finally(() => {
      setLoading(false);
    });
  }, []);

  React.useEffect(() => {
    setToken(Token.load());
  }, [tokenUpdated]);

  const onInit = (token) => {
    setInitialized(true);
    setTokenUpdated(true);
  };

  return (
    <>
      {loading && <>
        <Container>Loading...</Container>
      </>}
      {!loading && <>
        <BrowserRouter basename={window.PUBLIC_URL}>
          <Navigator initialized={initialized} token={token} />
          <Routes>
            {!initialized && <>
              <Route path="*" element={<Setup onInit={onInit} />}/>
            </>}
            {initialized && !token && <>
              <Route path="*" element={<Login onLogin={() => setTokenUpdated(!tokenUpdated)}/>}/>
            </>}
            {initialized && token && <>
              <Route path="*" element={<Login onLogin={() => setTokenUpdated(!tokenUpdated)}/>}/>
              <Route path="/routers-login" element={<Login onLogin={() => setTokenUpdated(!tokenUpdated)}/>}/>
              <Route path="/routers-dashboard" element={<Dashboard/>}/>
              <Route path="/routers-scenario" element={<Scenario/>}/>
              <Route path="/routers-config" element={<Config/>}/>
              <Route path="/routers-system" element={<System/>}/>
              <Route path="/routers-logout" element={<Logout onLogout={() => setTokenUpdated(!tokenUpdated)} />}/>
            </>}
          </Routes>
        </BrowserRouter>
      </>}
      <Footer />
    </>
  );
}

export default App;
