# srs-terraform

A control panel for SRS, to terraform the open-source video platform

## Usage

Install dependencies:

```bash
cd mgmt
npm install
```

Build UI:

```bash
cd mgmt
npm run build
```

Run server with UI:

```bash
cd mgmt
npm start
```

Access the browser: http://localhost:2022/mgmt

## Development

Install dependencies:

```bash
cd mgmt
npm install
```

Run the backend:

```
cd mgmt
npm start
```

Run the ui:

```
cd mgmt/ui
npm start
```

Access the browser: http://localhost:3000

## Ports

The ports allocated:

| Module | TCP Ports | UDP Ports | Notes |
| ------ | --------- | --------- | ----- |
| SRS | 1935, 1985, 8080,<br/> 8088, 1990, 554,<br/> 8936 | 8000, 8935, 10080,<br/> 1989 | See [SRS ports](https://github.com/ossrs/srs/blob/develop/trunk/doc/Resources.md#ports) |
| mgmt | 2022 |  - | Mount at `/mgmt/` |
| hooks | 2021 |  - | Mount at `/terraform/v1/hooks/` |

## Features

The features that we're developing:

* [x] A mgmt support authentication and automatic updates.
* [x] Run SRS in docker, query status by docker and SRS API.
* [x] Support publish by RTMP/WebRTC, play by RTMP/HTTP-FLV/HLS/WebRTC.
* [x] SRS container write containers/objs/srs.log for logging.
* [ ] Run SRS hooks in docker, to callback by SRS server.
* [ ] Support publish by SRT, play by RTMP/HTTP-FLV/HLS/WebRTC/SRT.
* [ ] Collect logs of mgmt and containers together.
* [ ] Stop, restart and upgrade containers.

