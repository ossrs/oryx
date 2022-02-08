import React from "react";
import {Accordion, Container, Form, Button} from "react-bootstrap";
import {Errors, Token} from "../utils";
import axios from "axios";
import {useNavigate} from "react-router-dom";

export default function Config() {
  const navigate = useNavigate();
  const [key, setKey] = React.useState();
  const [crt, setCrt] = React.useState();
  const [domain, setDomain] = React.useState();

  const updateSSL = (e) => {
    e.preventDefault();

    if (!key || !crt) {
      alert('请输入SSL密钥和证书');
      return;
    }

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/ssl', {
      ...token, key, crt,
    }).then(res => {
      alert(`SSL证书更新成功`);
      console.log(`SSL: Update ok`);
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

  const requestLetsEncrypt = (e) => {
    e.preventDefault();

    if (!domain) {
      alert('请输入你域名');
      return;
    }

    const token = Token.load();
    axios.post('/terraform/v1/mgmt/letsencrypt', {
      ...token, domain,
    }).then(res => {
      alert(`Let's Encrypt SSL证书更新成功`);
      console.log(`SSL: Let's Encrypt SSL ok`);
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

  return (
    <>
      <p></p>
      <Container>
        <Accordion defaultActiveKey="0">
          <Accordion.Item eventKey="0">
            <Accordion.Header>HTTPS: Let's Encrypt</Accordion.Header>
            <Accordion.Body>
              <Form>
                <Form.Group className="mb-3">
                  <Form.Label>域名</Form.Label>
                  <Form.Text> * 你的域名，请先解析到本服务器的公网IP，例如 your-domain.com</Form.Text>
                  <Form.Control as="input" defaultValue={domain} onChange={(e) => setDomain(e.target.value)} />
                </Form.Group>
                <Button variant="primary" type="submit" onClick={(e) => requestLetsEncrypt(e)}>
                  申请证书
                </Button>
              </Form>
            </Accordion.Body>
          </Accordion.Item>
          <Accordion.Item eventKey="1">
            <Accordion.Header>HTTPS: 上传证书</Accordion.Header>
            <Accordion.Body>
              <Form>
                <Form.Group className="mb-3">
                  <Form.Label>密钥(KEY)</Form.Label>
                  <Form.Text> * Nginx格式的SSL密钥内容，例如 your-domain.com.key</Form.Text>
                  <Form.Control as="textarea" rows={5} defaultValue={key} onChange={(e) => setKey(e.target.value)} />
                </Form.Group>
                <Form.Group className="mb-3">
                  <Form.Label>证书(PEM格式)</Form.Label>
                  <Form.Text> * Nginx格式的SSL证书内容，例如 your-domain.com.pem</Form.Text>
                  <Form.Control as="textarea" rows={5} defaultValue={crt} onChange={(e) => setCrt(e.target.value)} />
                </Form.Group>
                <Button variant="primary" type="submit" onClick={(e) => updateSSL(e)}>
                  更新证书
                </Button>
              </Form>
            </Accordion.Body>
          </Accordion.Item>
        </Accordion>
      </Container>
    </>
  );
}

