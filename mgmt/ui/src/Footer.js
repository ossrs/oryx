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
      <p className="text-center">
        <a href='https://github.com/ossrs/srs-terraform'>
          srs-terraform/{versions?.data?.version}
        </a>
      </p>
    </Container>
  );
}
