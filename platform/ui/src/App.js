import React from 'react';
import axios from "axios";
import './App.css'
import {Container} from "react-bootstrap";
import {BrowserRouter, Routes, Route, useParams, Outlet, useNavigate, useLocation} from "react-router-dom";
import Footer from './pages/Footer';
import Login from './pages/Login';
import Logout from './pages/Logout';
import Navigator from './pages/Navigator';
import Setup from './pages/Setup';
import {Locale, Token} from "./utils";
import Components from "./pages/Components";
import Scenario from "./pages/Scenario";
import Settings from "./pages/Settings";
import Dashboard from './pages/Dashboard';
import Contact from "./pages/Contact";
import {ErrorBoundary, useErrorHandler} from 'react-error-boundary';
import {SrsErrorBoundary} from "./components/SrsErrorBoundary";
import resources from "./resources/locale.json";
import {SrsEnvContext} from "./components/SrsEnvContext";

function App() {
  const [env, setEnv] = React.useState();

  return (
    <SrsEnvContext.Provider value={[env, setEnv]}>
      <ErrorBoundary FallbackComponent={(RootError)}>
        <SrsErrorBoundary>
          <AppImpl/>
        </SrsErrorBoundary>
      </ErrorBoundary>
    </SrsEnvContext.Provider>
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
  const setEnv = React.useContext(SrsEnvContext)[1];
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
    if (!setEnv) return;

    axios.get('/terraform/v1/mgmt/envs').then(res => {
      setEnv(res.data.data);
      console.log(`Env ok, ${JSON.stringify(res.data)}`);
    }).catch(handleError);
  }, [handleError, setEnv]);

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
          <Navigator {...{initialized, token}} />
          <Routes>
            <Route path="/" element={<AppRoot />}/>
            <Route path=':locale' element={<AppLocale />}>
              {!initialized && <>
                <Route path="*" element={<Setup onInit={onInit} />}/>
                <Route path="routers-setup" element={<Setup onInit={onInit} />}/>
              </>}
              {initialized && !token && <>
                <Route path="*" element={<Login onLogin={() => setTokenUpdated(!tokenUpdated)}/>}/>
              </>}
              {initialized && token && <>
                <Route path="*" element={<Login onLogin={() => setTokenUpdated(!tokenUpdated)}/>}/>
                <Route path="routers-login" element={<Login onLogin={() => setTokenUpdated(!tokenUpdated)}/>}/>
                <Route path="routers-dashboard" element={<Dashboard/>}/>
                <Route path="routers-scenario" element={<Scenario/>}/>
                <Route path="routers-settings" element={<Settings/>}/>
                <Route path="routers-contact" element={<Contact/>}/>
                <Route path="routers-components" element={<Components/>}/>
                <Route path="routers-logout" element={<Logout onLogout={() => setTokenUpdated(!tokenUpdated)} />}/>
              </>}
            </Route>
          </Routes>
          <Footer />
        </BrowserRouter>
      </>}
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

  return <Outlet />;
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
