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
import Settings from "./pages/Settings";
import Dashboard from './pages/Dashboard';
import Contact from "./pages/Contact";
import {ErrorBoundary, useErrorHandler} from 'react-error-boundary';
import {SrsErrorBoundary} from "./components/ErrorBoundary";

function App() {
  return (
    <ErrorBoundary FallbackComponent={(RootError)}>
      <SrsErrorBoundary>
        <AppImpl />
      </SrsErrorBoundary>
    </ErrorBoundary>
  );
}

function RootError({error}) {
  return <Container>{error?.message}</Container>;
}

function AppImpl() {
  const [loading, setLoading] = React.useState(true);
  const [initialized, setInitialized] = React.useState();
  const [tokenUpdated, setTokenUpdated] = React.useState();
  const [token, setToken] = React.useState();
  const handleError = useErrorHandler();

  React.useEffect(() => {
    axios.get('/terraform/v1/mgmt/init').then(res => {
      setInitialized(res.data.data.init);
    }).catch(handleError).finally(() => {
      setLoading(false);
    });
  }, [handleError]);

  React.useEffect(() => {
    axios.get('/terraform/v1/mgmt/check').then(res => {
      console.log(`Check ok, ${JSON.stringify(res.data)}`);
    }).catch(handleError);
  }, [handleError]);

  React.useEffect(() => {
    setToken(Token.load());
  }, [tokenUpdated]);

  const onInit = React.useCallback((token) => {
    setInitialized(true);
    setTokenUpdated(true);
  }, []);

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
              <Route path="/routers-settings" element={<Settings/>}/>
              <Route path="/routers-contact" element={<Contact/>}/>
              <Route path="/routers-system" element={<System/>}/>
              <Route path="/routers-logout" element={<Logout onLogout={() => setTokenUpdated(!tokenUpdated)} />}/>
            </>}
          </Routes>
          <Footer />
        </BrowserRouter>
      </>}
    </>
  );
}

export default App;
