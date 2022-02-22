import React from "react";
import {Button, Form} from "react-bootstrap";
import {Errors, Token} from "../utils";
import axios from "axios";
import {useNavigate} from "react-router-dom";

export default function SetupCamSecret() {
  const navigate = useNavigate();
  const [secretId, setSecretId] = React.useState();
  const [secretKey, setSecretKey] = React.useState();

  const updateTencentSecret = (e) => {
    e.preventDefault();

    const token = Token.load();
    axios.post('/terraform/v1/tencent/cam/secret', {
      ...token, secretId, secretKey,
    }).then(res => {
      alert('腾讯云访问密钥设置成功');
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/routers-logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  };

  return (<>
    <Form>
      <Form.Group className="mb-3">
        <Form.Label>SecretId</Form.Label>
        <Form.Text> * 腾讯云的SecretId, <a href='https://console.cloud.tencent.com/cam' target='_blank' rel='noreferrer'>获取密钥</a></Form.Text>
        <Form.Control as="input" rows={2} defaultValue={secretId} onChange={(e) => setSecretId(e.target.value)} />
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>SecretKey</Form.Label>
        <Form.Text> * 腾讯云的SecretKey, <a href='https://console.cloud.tencent.com/cam' target='_blank' rel='noreferrer'>获取密钥</a></Form.Text>
        <Form.Control as="input" type='password' rows={2} defaultValue={secretKey} onChange={(e) => setSecretKey(e.target.value)} />
      </Form.Group>
      <Button variant="primary" type="submit" onClick={(e) => updateTencentSecret(e)}>
        设置账号
      </Button>
    </Form>
  </>);
}

