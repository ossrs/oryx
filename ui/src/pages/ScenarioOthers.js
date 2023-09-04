//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {useSrsLanguage} from "../components/LanguageSwitch";
import {Accordion} from "react-bootstrap";

export function ScenarioOther() {
  const language = useSrsLanguage();
  if (language === 'zh') {
    return <>
      <Accordion defaultActiveKey='0'>
        <Accordion.Item eventKey="0">
          <Accordion.Header>场景介绍</Accordion.Header>
          <Accordion.Body>
            <div>
              其他非常用场景。
              <p></p>
            </div>
          </Accordion.Body>
        </Accordion.Item>
      </Accordion>
    </>;
  }
  return <>
    <Accordion defaultActiveKey='0'>
      <Accordion.Item eventKey="0">
        <Accordion.Header>Introduction</Accordion.Header>
        <Accordion.Body>
          <div>
            Other less common scenarios.
            <p></p>
          </div>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  </>;
}

