import {useNavigate} from "react-router-dom";
import {Container, Tabs, Tab, Accordion, Form, Row, Col, Button} from "react-bootstrap";
import React from "react";
import {Token, Errors} from "./utils";
import axios from "axios";

export default function Dashboard() {
  const navigate = useNavigate();
  const [rtmpServer, setRtmpServer] = React.useState();
  const [rtmpStreamKey, setRtmpStreamKey] = React.useState();
  const [cnConsole, setCnConsole] = React.useState();
  const [cnPlayer, setCnPlayer] = React.useState();

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

  React.useEffect(() => {
    // Build RTMP url.
    if (true) {
      setRtmpServer(`rtmp://${window.location.hostname}/live/`);
      setRtmpStreamKey(`livestream`);
    }

    // Build console url.
    if (true) {
      // If not 8080, user should proxy to the default port.
      const query = parseInt(window.location.port) === 8080 ? `?port=1985` : '';
      const cnUrl = `console/ng_index.html#/summaries${query}`;
      setCnConsole(cnUrl);
    }

    // The player url.
    if (true) {
      const prefix = `players/?schema=${window.location.protocol.replace(':', '')}`;
      const httpPort = window.location.port || (window.location.protocol === 'http:' ? 80 : 443);
      // If not 8080, user should proxy both stream and API to the default port.
      const query = parseInt(window.location.port) === 8080 ? '' : `&port=${httpPort}&api=${httpPort}`;
      setCnPlayer(`${prefix}${query}`);
    }
  }, []);

  return (
    <>
      <p></p>
      <Container>
        <Tabs defaultActiveKey="live" id="uncontrolled-tab-example" className="mb-3">
          <Tab eventKey="live" title="SRS控制台">
            <Accordion defaultActiveKey="0">
              <Accordion.Item eventKey="0">
                <Accordion.Header>OBS推流</Accordion.Header>
                <Accordion.Body>
                  <div>
                    先<a href="https://obsproject.com/download">下载OBS</a>
                  </div>
                  <div>
                    推流地址 <code>{rtmpServer}</code>
                  </div>
                  <div>
                    推流密钥 <code>{rtmpStreamKey}</code>
                  </div>
                  <div>
                    点击进入<a id="cnPlayer" href={cnPlayer}>SRS播放器</a>
                  </div>
                  <div>
                    点击进入<a id="cnConsole" href={cnConsole}>SRS控制台</a>
                  </div>
                </Accordion.Body>
              </Accordion.Item>
              <Accordion.Item eventKey="1">
                <Accordion.Header>FFmpeg推流</Accordion.Header>
                <Accordion.Body>
                  <div>
                    先<a href="https://ffmpeg.org/download.html">下载FFmpeg</a>
                  </div>
                  <div>
                    推流地址 <code>{rtmpServer}{rtmpStreamKey}</code>
                  </div>
                  <div>
                    点击进入<a id="cnPlayer" href={cnPlayer}>SRS播放器</a>
                  </div>
                  <div>
                    点击进入<a id="cnConsole" href={cnConsole}>SRS控制台</a>
                  </div>
                </Accordion.Body>
              </Accordion.Item>
            </Accordion>
          </Tab>
          <Tab eventKey="source" title="音视频开发套件">
            <Accordion defaultActiveKey="0">
              <Accordion.Item eventKey="0">
                <Accordion.Header>SRS</Accordion.Header>
                <Accordion.Body>
                  <code>cd ~lighthouse/git/srs && ./configure && make</code>
                </Accordion.Body>
              </Accordion.Item>
            </Accordion>
          </Tab>
        </Tabs>
      </Container>
    </>
  );
}

