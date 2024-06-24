//
// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
//
import React from "react";
import {Toast} from "react-bootstrap";
import logo from '../resources/logo.svg';
import * as Icon from 'react-bootstrap-icons';

function TutorialsText({title, children, prefixLine}) {
  const [show, setShow] = React.useState(false);

  return (
    <>
      <div role='button' style={{display: 'inline-block'}}>
        <Icon.PatchQuestion onClick={() => setShow(!show)}/>
      </div>
      {show && prefixLine && <p></p>}
      {show &&
        <>
          <Toast onClose={() => setShow(false)}>
            <Toast.Header>
              <img src={logo} className="rounded me-2" width={56} alt=''/>
              <strong className="me-auto">{title}</strong>
            </Toast.Header>
            <Toast.Body>
              {children}
            </Toast.Body>
          </Toast>
          <p></p>
        </>
      }
    </>
  );
}

export default TutorialsText;

