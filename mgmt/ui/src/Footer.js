import React from "react";
import Container from "react-bootstrap/Container";

export default function Footer() {
  const [versions, setVersions] = React.useState();
  React.useEffect(() => {
    fetch('/terraform/v1/mgmt/versions')
      .then((res) => res.json())
      .then(res => setVersions(res));
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
