import React from "react";
import Container from "react-bootstrap/Container";
import axios from "axios";
import {SrsErrorBoundary} from "../components/ErrorBoundary";
import {useErrorHandler} from "react-error-boundary";

export default function Footer() {
  return (
    <SrsErrorBoundary>
      <FooterImpl />
    </SrsErrorBoundary>
  );
}

function FooterImpl() {
  const [versions, setVersions] = React.useState();
  const [beian, setBeian] = React.useState();
  const handleError = useErrorHandler();

  React.useEffect(() => {
    axios.get('/terraform/v1/mgmt/versions')
      .then(res => setVersions(res.data)).catch(handleError);
  }, []);

  React.useEffect(() => {
    axios.get('/terraform/v1/mgmt/beian/query')
      .then(res => {
        setBeian(res.data.data);
        console.log(`Beian: query ${JSON.stringify(res.data.data)}`);
      }).catch(handleError);
  }, []);

  return (
    <Container>
      <p></p>
      <p className="text-center">
        <a href='https://github.com/ossrs/srs-cloud' target='_blank' rel='noreferrer'>
          &copy;srs-cloud/v{versions?.data?.version}
        </a>
        &nbsp; <a href='https://beian.miit.gov.cn' target='_blank' rel='noreferrer'>{beian?.icp}</a>
      </p>
    </Container>
  );
}
