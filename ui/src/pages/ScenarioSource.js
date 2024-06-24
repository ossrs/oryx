//
// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
//
import {Accordion} from "react-bootstrap";
import React from "react";
import {useSrsLanguage} from "../components/LanguageSwitch";

export default function ScenarioSource() {
  const language = useSrsLanguage();
  return language === 'zh' ? <ScenarioSourceCn /> : <ScenarioSourceEn />;
}

function ScenarioSourceCn() {
  return (
    <Accordion defaultActiveKey="0">
      <Accordion.Item eventKey="0">
        <Accordion.Header>Oryx</Accordion.Header>
        <Accordion.Body>
          <div>Oryx是开源项目：</div>
          <div>1. Github：<a href='https://github.com/ossrs/oryx'>https://github.com/ossrs/oryx</a></div>
          <div>2. Gitee：<a href='https://gitee.com/ossrs/oryx'>https://gitee.com/ossrs/oryx</a></div>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>SRS 4.0</Accordion.Header>
        <Accordion.Body>
          <div>LightHouse云服务器自带了SRS 4.0的源码，你可以选择：</div>
          <div>1. 点击下载SRS源码：<a href='/terraform/v1/sources/srs.tar.gz'>下载</a></div>
          <div>2. 直接在云服务器编译SRS：</div>
          <div><code>cd ~lighthouse/git/srs/trunk</code></div>
          <div><code>git pull</code></div>
          <div><code>./configure</code></div>
          <div><code>make</code></div>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

function ScenarioSourceEn() {
  return (
    <span>On the way...</span>
  );
}

