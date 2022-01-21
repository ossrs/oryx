import React from 'react';
import logo from './logo.svg';
import './App.css';

function App() {
  const [versions, setVersions] = React.useState();
  React.useEffect(() => {
    fetch('/terraform/v1/mgmt/versions')
      .then((res) => res.json())
      .then(res => setVersions(res));
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        Backend api {versions?.data?.version}
      </header>
    </div>
  );
}

export default App;
