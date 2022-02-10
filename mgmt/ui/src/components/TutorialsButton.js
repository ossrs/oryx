import React from "react";
import {Button, Toast} from "react-bootstrap";
import logo from '../resources/logo.svg';
import * as Icon from 'react-bootstrap-icons';

export default function TutorialsButton({tutorials, prefixLine}) {
  const [show, setShow] = React.useState(false);

  return (
    <>
      <div role='button' style={{display: 'inline-block'}}>
        <Icon.PatchQuestion onClick={() => setShow(!show)} />
      </div>
      {show && tutorials.map((tutorial, index) => (
        <div key={index}>
          {prefixLine && <p></p>}
          <Toast show={show} onClose={() => setShow(false)}>
            <Toast.Header>
              <img src={logo} className="rounded me-2" width={32}/>
              <strong className="me-auto">SRS云服务器</strong>
              <small>by {tutorial.author}</small>
            </Toast.Header>
            <Toast.Body>
              <a href={tutorial.link} target='_blank' rel='noreferrer'>{tutorial.title}</a>
            </Toast.Body>
          </Toast>
        </div>
      ))}
      {show && <p></p>}
    </>
  );
}
