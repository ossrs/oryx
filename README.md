# srs-terraform

A control panel for SRS, to terraform the open-source video platform

## Usage

Install dependencies:

```bash
cd releases && npm install
cd hooks && npm install
cd mgmt && npm install
```

Run the releases:

```
cd releases
npm start
```

Run the hooks:

```
cd hooks
npm start
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
| releases | 2023 |  - | Mount at `/terraform/v1/releases` |
| mgmt | 2022 |  - | Mount at `/mgmt/` and `/terraform/v1/mgmt/` |
| hooks | 2021 |  - | Mount at `/terraform/v1/hooks/` |

## Features

The features that we're developing:

* [x] A mgmt support authentication and automatic updates.
* [x] Run SRS in docker, query status by docker and SRS API.
* [x] Support publish by RTMP/WebRTC, play by RTMP/HTTP-FLV/HLS/WebRTC.
* [x] SRS container write containers/objs/srs.log for logging.
* [x] Support high-resolution and realtime(200~500ms) live streaming by SRT.
* [x] Run SRS hooks in docker, to callback by SRS server.
* [ ] Support publish by SRT, play by RTMP/HTTP-FLV/HLS/WebRTC/SRT.
* [ ] Collect logs of mgmt and containers together.
* [ ] Stop, restart and upgrade containers.
* [ ] Support logrotate to manage the logs.

## APIs

Platform:

* `/terraform/v1/mgmt/versions` Public version api.
* `/terraform/v1/mgmt/init` Whether mgmt initialized.
* `/terraform/v1/mgmt/status` Query the version of mgmt.
* `/terraform/v1/mgmt/upgrade` Upgrade the mgmt to latest version.
* `/terraform/v1/mgmt/token` System auth with token.
* `/terraform/v1/mgmt/login` System auth with password.

Releases:

* `/terraform/v1/releases` Version management for all components.
* `/terraform/v1/mgmt/srs` SRS: Query and upgrade SRS container.
* `/terraform/v1/mgmt/hooks` Hooks: Query the hooks container.

Market:

* `/terraform/v1/hooks/srs/verify` Hooks: Verify the stream request URL of SRS.
* `/terraform/v1/hooks/srs/secret` Hooks: Query the secret to generate stream URL.

