import Container from "react-bootstrap/Container";
import React from "react";
import {Form, Button} from 'react-bootstrap';
import axios from "axios";
import {useNavigate} from "react-router-dom";
import {Token, Tools} from './utils';

export default function Login({initialized, onLogin}) {
  const [token, setToken] = React.useState();
  const [verify, setVerify] = React.useState(false);
  const [plaintext, setPlaintext] = React.useState(false);
  const [password, setPassword] = React.useState();
  const navigate = useNavigate();
  const passwordRef = React.useRef();
  const plaintextRef = React.useRef();

  // Load the token if page initialized or verify changed.
  React.useEffect(() => {
    if (initialized === undefined) return ;

    console.log(`Login: Load token, initialized=${initialized}, verify=${verify}`);
    Token.load((data) => {
      if (!data || !data.token) return;
      setToken(data);
    });
  }, [initialized, verify]);

  // Verify the token if token changed.
  React.useEffect(() => {
    if (!token || !token.token) return;

    console.log(`Login: Verify, token is ${Tools.mask(token)}`);
    axios.post('/terraform/v1/mgmt/token', {
      ...token,
    }).then(res => {
      setVerify(true);
    }).catch(e => {
      const err = e.response.data;
      alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
    });
  }, [token]);

  // Redirect if token verify done.
  React.useEffect(() => {
    if (!token || !token.token) return;
    if (verify !== true) return;

    console.log(`Login: Done, verify=${verify}, token is ${Tools.mask(token)}`);
    navigate('/status');
  }, [token, verify]);

  // Generate password if not initialized.
  React.useEffect(() => {
    if (initialized !== false) return;

    setPassword(Math.random().toString(16).slice(-6));
    setPlaintext(true);
  }, [initialized]);

  // Focus to password input.
  React.useEffect(() => {
    plaintext ? plaintextRef.current?.focus() : passwordRef.current?.focus();
  }, [initialized, plaintext]);

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
      navigate('/status');
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
            <Form.Label>
              {initialized ? '请输入密码' : '请设置初始密码'}
            </Form.Label>
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
              {!initialized ? '* 自动生成的初始管理员密码，你可以修改' : '* 若忘记密码，可登录机器查看文件 ~lighthouse/credentials.txt'}
            </Form.Text>
          </Form.Group>
          {
            initialized &&
            <Form.Group className="mb-3" controlId="formBasicCheckbox">
              <Form.Check type="checkbox" label="显示密码" defaultChecked={plaintext}
                onClick={() => setPlaintext(!plaintext)}/>
            </Form.Group>
          }
          <Button variant="primary" type="submit" onClick={(e) => handleLogin(e)}>
            {initialized ? '登录' : '设置管理员密码'}
          </Button>
        </Form>
      </Container>
    </>
  );
}

