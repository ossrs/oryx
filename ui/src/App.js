//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from 'react';
import axios from "axios";
import './App.css';
import './ai-talk.css';
import {Container} from "react-bootstrap";
import {
  BrowserRouter,
  Routes,
  Route,
  useParams,
  Outlet,
  useNavigate,
  useLocation,
  useSearchParams
} from "react-router-dom";
import Footer from './pages/Footer';
import Login from './pages/Login';
import Logout from './pages/Logout';
import Navigator from './pages/Navigator';
import Setup from './pages/Setup';
import {Locale, Token} from "./utils";
import Components from "./pages/Components";
import Scenario from "./pages/Scenario";
import Settings from "./pages/Settings";
import Contact from "./pages/Contact";
import {ErrorBoundary, useErrorHandler} from 'react-error-boundary';
import {SrsErrorBoundary} from "./components/SrsErrorBoundary";
import resources from "./resources/locale.json";
import {SrsEnvContext} from "./components/SrsEnvContext";
import Popouts from "./pages/Popouts";

function App() {
  const [env, setEnv] = React.useState(null);

  return (
    <SrsEnvContext.Provider value={[env, setEnv]}>
      <ErrorBoundary FallbackComponent={(RootError)}>
        <SrsErrorBoundary>
          <AppPreImpl/>
        </SrsErrorBoundary>
      </ErrorBoundary>
    </SrsEnvContext.Provider>
  );
}

function RootError({error}) {
  return <Container fluid>{error?.message}</Container>;
}

function AppPreImpl() {
  const [env, setEnv] = React.useContext(SrsEnvContext);
  const handleError = useErrorHandler();

  React.useEffect(() => {
    if (!setEnv) return;

    axios.post('/terraform/v1/mgmt/envs', {
      locale: Locale.current()
    }).then(res => {
      setEnv(res.data.data);
      console.log(`Env ok, ${JSON.stringify(res.data)}`);
    }).catch(handleError);
  }, [handleError, setEnv]);

  return <>{env && <AppImpl/>}</>;
}

function AppImpl() {
  const [loading, setLoading] = React.useState(true);
  // Possible value is 1: yes, -1: no, 0: undefined.
  const [initialized, setInitialized] = React.useState(0);
  const handleError = useErrorHandler();

  React.useEffect(() => {
    axios.get('/terraform/v1/mgmt/check').then(res => {
      console.log(`Check ok, ${JSON.stringify(res.data)}`);
      axios.get('/terraform/v1/mgmt/init').then(res => {
        setInitialized(res.data.data.init ? 1 : -1);
      }).catch(handleError);
    }).catch(handleError).finally(setLoading);
  }, [handleError]);

  return (
    <>
      {loading && <>
        <Container fluid>Loading...</Container>
      </>}
      {!loading && <>
        <BrowserRouter basename={window.PUBLIC_URL}>
          <AppRoute {...{initialized, setInitialized}} />
        </BrowserRouter>
      </>}
    </>
  );
}

function AppRoute({initialized, setInitialized}) {
  const [tokenUpdated, setTokenUpdated] = React.useState();
  const [token, setToken] = React.useState();
  const [searchParams] = useSearchParams();
  // Possible value is 1: yes, -1: no, 0: undefined.
  const [isPopout, setIsPopout] = React.useState(0);

  React.useEffect(() => {
    if (!searchParams) return;
    setIsPopout(searchParams.get('popout') === '1' ? 1 : -1);
  }, [searchParams]);

  React.useEffect(() => {
    setToken(Token.load());
  }, [tokenUpdated]);

  const onInit = React.useCallback((token) => {
    setInitialized(true);
    setTokenUpdated(true);
  }, [setInitialized, setTokenUpdated]);

  return (
    <>
      {isPopout === -1 && <Navigator {...{initialized, token}} />}
      <Routes>
        {initialized === 0 ?
          <React.Fragment>
            <Route path="*" element={<React.Fragment/>}/>
          </React.Fragment> :
          <React.Fragment>
            <Route path="/" element={<AppRoot/>}/>
            <Route path=':locale' element={<AppLocale/>}>
              {initialized === -1 && <>
                <Route path="*" element={<Setup onInit={onInit}/>}/>
                <Route path="routers-setup" element={<Setup onInit={onInit}/>}/>
              </>}
              {initialized === 1 && !token && <>
                <Route path="*" element={<Login onLogin={() => setTokenUpdated(!tokenUpdated)}/>}/>
              </>}
              {initialized === 1 && token && <>
                <Route path="*" element={<Login onLogin={() => setTokenUpdated(!tokenUpdated)}/>}/>
                <Route path="routers-login" element={<Login onLogin={() => setTokenUpdated(!tokenUpdated)}/>}/>
                <Route path="routers-scenario" element={<Scenario/>}/>
                <Route path="routers-settings" element={<Settings/>}/>
                <Route path="routers-contact" element={<Contact/>}/>
                <Route path="routers-components" element={<Components/>}/>
                <Route path="routers-logout" element={<Logout onLogout={() => setTokenUpdated(!tokenUpdated)}/>}/>
              </>}
              {initialized === 1 && <Route path="routers-popout" element={<Popouts/>}/>}
            </Route>
          </React.Fragment>}
      </Routes>
      {isPopout === -1 && <Footer/> }
    </>
  );
}

function AppLocale() {
  const params = useParams();
  const navigate = useNavigate();
  const location = useLocation();

  // Prefix with default language.
  React.useEffect(() => {
    if (Object.keys(resources).includes(params.locale)) {
      return;
    }
    if (!Object.keys(resources).includes(Locale.current())) {
      return;
    }

    const to = {pathname: `/${Locale.current()}/${params.locale}`, search: location.search};
    console.log(`Jump to ${JSON.stringify(to)} by locale`);
    return navigate(to);
  }, [navigate, params, location]);

  return <Outlet/>;
}

function AppRoot() {
  const navigate = useNavigate();
  const location = useLocation();

  React.useEffect(() => {
    const to = {pathname: `/${Locale.current()}/routers-login`, search: location.search};
    console.log(`Jump to ${JSON.stringify(to)} by root`);
    navigate(to);
  }, [navigate, location]);

  return <></>;
}

export default App;
