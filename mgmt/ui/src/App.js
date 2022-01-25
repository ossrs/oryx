import React from 'react';
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
import Dashboard from "./Dashboard";

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
            <Route path="/dashboard" element={<System/>}/>
            <Route path="/system" element={<Dashboard/>}/>
            <Route path="/logout" element={<Logout onLogout={() => setTokenUpdated(!tokenUpdated)} />}/>
          </>}
        </Routes>
      </BrowserRouter>
      <Footer />
    </>
  );
}

export default App;
