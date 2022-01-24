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
import {Token} from "./utils";
import System from "./System";

function App() {
  const [initialized, setInitialized] = React.useState();
  const [tokenUpdated, setTokenUpdated] = React.useState();
  const [token, setToken] = React.useState();

  React.useEffect(() => {
    axios.get('/terraform/v1/mgmt/init').then(res => {
      setInitialized(res.data.data.init);
    }).catch(e => {
      alert(e.response.data);
      console.error(e);
    });
  }, []);

  React.useEffect(() => {
    setToken(Token.load());
  }, [tokenUpdated]);

  return (
    <>
      <BrowserRouter basename={window.PUBLIC_URL}>
        <Navigator initialized={initialized} token={token} />
        <Routes>
          {!initialized && <>
            <Route path="*" element={<Init onInit={()=>setInitialized(true)} />}/>
          </>}
          {initialized && !token && <>
            <Route path="*" element={<Login onLogin={() => setTokenUpdated(!tokenUpdated)}/>}/>
          </>}
          {initialized && token && <>
            <Route path="*" element={<Login onLogin={() => setTokenUpdated(!tokenUpdated)}/>}/>
            <Route path="/login" element={<Login onLogin={() => setTokenUpdated(!tokenUpdated)}/>}/>
            <Route path="/system" element={<System/>}/>
            <Route path="/logout" element={<Logout onLogout={() => setTokenUpdated(!tokenUpdated)} />}/>
          </>}
        </Routes>
      </BrowserRouter>
      <Footer />
    </>
  );
}

export default App;
