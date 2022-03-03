import React from "react";
import Container from "react-bootstrap/Container";
import axios from "axios";

export default function Footer() {
  const [versions, setVersions] = React.useState();
  const [beian, setBeian] = React.useState();

  React.useEffect(() => {
    axios.get('/terraform/v1/mgmt/versions')
      .then(res => setVersions(res.data));
  }, []);

  React.useEffect(() => {
    axios.get('/terraform/v1/mgmt/beian/query')
      .then(res => {
        setBeian(res.data.data);
        console.log(`Beian: query ${JSON.stringify(res.data.data)}`);
      });
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
