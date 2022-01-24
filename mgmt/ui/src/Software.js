import {useNavigate} from "react-router-dom";
import Container from "react-bootstrap/Container";
import React from "react";
import {Token, Errors} from "./utils";
import axios from "axios";

export default function Software() {
  const navigate = useNavigate();
  const [status, setStatus] = React.useState();

  // Verify the token if token changed.
  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/mgmt/software', {
      ...token,
    }).then(res => {
      setStatus(res.data.data);
      console.log(`Status: Query ok, software=${JSON.stringify(res.data.data)}`);
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  }, []);

  return (
    <>
      {status && <Container>System: {status.version}</Container>}
    </>
  );
}

