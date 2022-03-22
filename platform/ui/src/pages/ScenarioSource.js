import {Accordion} from "react-bootstrap";
import React from "react";

export default function ScenarioSource() {
  return (
    <Accordion defaultActiveKey="0">
      <Accordion.Item eventKey="0">
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
      <Accordion.Item eventKey="1">
        <Accordion.Header>SRS/GB28181</Accordion.Header>
        <Accordion.Body>
          <div>以<a href='https://github.com/ossrs/srs-gb28181' target='_blank' rel='noreferrer'>srs-gb28181</a>为例：</div>
          <div>1. 点击下载SRS源码：<a href='/terraform/v1/sources/srs.tar.gz'>下载</a></div>
          <div>2. 在本机设置为srs-gb28181的源：</div>
          <div><code>git remote set-url origin https://github.com/ossrs/srs-gb28181.git</code></div>
          <div><code>git fetch origin</code></div>
          <div><code>git checkout -b feature/gb28181 origin/feature/gb28181</code></div>
          <div>3. 后续就只需要从github增量更新代码就可以：</div>
          <div><code>git pull</code></div>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

