import {useNavigate} from "react-router-dom";
import Container from "react-bootstrap/Container";
import React from "react";
import {Token, Errors} from "./utils";
import axios from "axios";

export default function Dashboard() {
  const navigate = useNavigate();

  React.useEffect(() => {
    const token = Token.load();
    axios.post('/terraform/v1/mgmt/status', {
      ...token,
    }).then(res => {
      console.log(`Status: Query ok, status=${JSON.stringify(res.data.data)}`);
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
      <p></p>
      <Container>
      </Container>
    </>
  );
}

