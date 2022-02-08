import React from "react";
import Container from "react-bootstrap/Container";
import axios from "axios";

export default function Footer() {
  const [versions, setVersions] = React.useState();
  React.useEffect(() => {
    axios.get('/terraform/v1/mgmt/versions')
      .then(res => setVersions(res.data));
  }, []);

  return (
    <Container>
      <p></p>
      <p className="text-center">
        <a href='https://github.com/ossrs/srs-terraform' target='_blank' rel='noreferrer'>
          &copy;ossrs mgmt/{versions?.data?.version}
        </a>
      </p>
    </Container>
  );
}
