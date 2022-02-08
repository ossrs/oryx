import Container from "react-bootstrap/Container";
import React from "react";
import {Form, Button} from 'react-bootstrap';
import axios from "axios";
import {useNavigate} from "react-router-dom";
import {Token, Tools} from '../utils';

export default function Login({onLogin}) {
  const [plaintext, setPlaintext] = React.useState(false);
  const [password, setPassword] = React.useState();
  const navigate = useNavigate();
  const passwordRef = React.useRef();
  const plaintextRef = React.useRef();

  // Verify the token if exists.
  React.useEffect(() => {
    const token = Token.load();
    if (!token || !token.token) return;

    console.log(`Login: Verify, token is ${Tools.mask(token)}`);
    axios.post('/terraform/v1/mgmt/token', {
      ...token,
    }).then(res => {
      console.log(`Login: Done, token is ${Tools.mask(token)}`);
      navigate('/routers-scenario');
    }).catch(e => {
      const err = e.response.data;
      alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
    });
  }, [navigate]);

  // Focus to password input.
  React.useEffect(() => {
    plaintext ? plaintextRef.current?.focus() : passwordRef.current?.focus();
  }, [plaintext]);

  // User click login button.
  const handleLogin = (e) => {
    e.preventDefault();

    axios.post('/terraform/v1/mgmt/login', {
      password,
    }).then(res => {
      const data = res.data.data;
      console.log(`Login: OK, token is ${Tools.mask(data)}`);
      Token.save(data);

      onLogin && onLogin();
      navigate('/routers-scenario');
    }).catch(e => {
      const err = e.response.data;
      alert(`${err.code}: ${err.data.message}`);
      console.error(e);
    });
  };

  return (
    <>
      <Container>
        <Form>
          <Form.Group className="mb-3" controlId="formBasicPassword">
            <Form.Label>请输入密码</Form.Label>
            {
              !plaintext &&
              <Form.Control type="password" placeholder="Password" ref={passwordRef} defaultValue={password}
                onChange={(e) => setPassword(e.target.value)}/>
            }
            {
              plaintext &&
              <Form.Control type="text" placeholder="Password" ref={plaintextRef} defaultValue={password}
                onChange={(e) => setPassword(e.target.value)}/>
            }
            <Form.Text className="text-muted">
              * 忘记密码？可登录机器执行 <code>cat ~lighthouse/credentials.txt</code>
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formBasicCheckbox">
            <Form.Check type="checkbox" label="显示密码" defaultChecked={plaintext}
              onClick={() => setPlaintext(!plaintext)}/>
          </Form.Group>
          <Button variant="primary" type="submit" onClick={(e) => handleLogin(e)}>
            登录
          </Button>
        </Form>
      </Container>
    </>
  );
}

