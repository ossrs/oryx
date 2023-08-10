import Container from "react-bootstrap/Container";
import React from "react";
import {Token} from "../utils";
import axios from "axios";
import {Row, Col, Card} from "react-bootstrap";
import * as moment from 'moment';
import {SrsErrorBoundary} from "../components/SrsErrorBoundary";
import {useTranslation} from "react-i18next";
import {SrsEnvContext} from "../components/SrsEnvContext";

export default function Components() {
  return (
    <SrsErrorBoundary>
      <ComponentsImpl />
    </SrsErrorBoundary>
  );
}

function ComponentsImpl() {
  const [status, setStatus] = React.useState();
  const {t} = useTranslation();
  const env = React.useContext(SrsEnvContext)[0];

  React.useEffect(() => {
    const refreshMgmtStatus = () => {
      const token = Token.load();
      axios.post('/terraform/v1/mgmt/status', {
        ...token,
      }).then(res => {
        const status = res.data.data;

        // Normally state.
        setStatus(status);

        console.log(`${moment().format()}: Status: Query ok, status=${JSON.stringify(status)}`);
      }).catch(e => {
        console.log('ignore any error during status', e);
      });
    };

    refreshMgmtStatus();
    const timer = setInterval(() => refreshMgmtStatus(), 10 * 1000);
    return () => clearInterval(timer);
  }, [setStatus, env]);

  return (
    <>
      <Container>
        <Row>
          <Col xs lg={3}>
            <Card style={{ width: '18rem', marginTop: '16px' }}>
              <Card.Header>{t('coms.host')}</Card.Header>
              <Card.Body>
                <Card.Text as={Col}>
                  {t('coms.version')}: {status?.version} <br/>
                  {t('coms.stable')}: {status?.releases?.stable}<br/>
                  {t('coms.latest')}: <a href='https://github.com/ossrs/srs-stack/issues/4#changelog' target='_blank' rel='noreferrer'>{status?.releases?.latest}</a>
                </Card.Text>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </Container>
    </>
  );
}

