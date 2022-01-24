import Container from "react-bootstrap/Container";
import React from "react";
import {Form, Button} from 'react-bootstrap';
import axios from "axios";
import {useNavigate} from "react-router-dom";
import {Token} from './utils';

export default function Login({initialized, onLogin}) {
  const [plaintext, setPlaintext] = React.useState(false);
  const [password, setPassword] = React.useState();
  const navigate = useNavigate();
  const passwordRef = React.useRef();
  const plaintextRef = React.useRef();

  React.useEffect(() => {
    Token.load((data) => {
      if (data) navigate('/status');
    });
  }, [initialized]);

  React.useEffect(() => {
    if (initialized !== false) return;
    setPassword(Math.random().toString(16).slice(-6));
    setPlaintext(true);
  }, [initialized]);

  React.useEffect(() => {
    plaintext ? plaintextRef.current?.focus() : passwordRef.current?.focus();
  }, [plaintext]);

  const handleLogin = (e) => {
    e.preventDefault();
    axios.post('/terraform/v1/mgmt/login', {
      password,
    }).then(res => {
      const token = res.data.data;
      const mask = `***${token.token.length}B***`;
      console.log(`Login: OK, token is ${JSON.stringify({...token, token: mask})}`);
      Token.save(token);

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

