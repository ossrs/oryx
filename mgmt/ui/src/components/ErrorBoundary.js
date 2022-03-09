import React from "react";
import {ErrorBoundary} from 'react-error-boundary';
import {Alert, Container} from "react-bootstrap";
import {Errors} from "../utils";
import {NavLink} from "react-router-dom";

export function SrsErrorBoundary({children}) {
  return (
    <ErrorBoundary FallbackComponent={ErrorFallback}>
      <>{children}</>
    </ErrorBoundary>
  );
}

export function ErrorFallback({error}) {
  const [show, setShow] = React.useState(true);

  if (!show) return <></>;
  return (
    <Container>
      <Alert variant="danger" onClose={() => setShow(false)} dismissible>
        <Alert.Heading>You got an error!</Alert.Heading>
        <ErrorDetail error={error} />
      </Alert>
    </Container>
  );
}

function ErrorDetail({error}) {
  if (!error) return (
    <p>Empty unknown error</p>
  );

  const err = error?.response?.data;
  if (err?.code === Errors.auth) {
    return <p>Token过期，请<NavLink to='/routers-logout'>重新登录</NavLink>，{err?.code}: {err?.data?.message}</p>;
  }

  if (err?.code) return (
    <p>
      Request: {`${error.request?.responseURL}`} <br/>
      Status: {`${error.response?.status}`} {`${error.response?.statusText}`} <br/>
      Code: {`${err?.code}`} <br/>
      Message: {`${err?.data?.message}`} <br/> <br/>
      <pre>
        {JSON.stringify(error.response.data, null, 2)}
      </pre>
    </p>
  );

  return (
    <p>
      Request: {`${error.request?.responseURL}`} <br/>
      Status: {`${error.response?.status}`} {`${error.response?.statusText}`} <br/>
      Data: {`${err}`}
    </p>
  );
}

