# SRS-Cloud

A lightweight open-source video cloud based on Nodejs, SRS, FFmpeg, WebRTC, etc.

## Usage: LightHouse

- [x] [Getting Started](https://mp.weixin.qq.com/s/fWmdkw-2AoFD_pEmE_EIkA).
- [x] [Live Streaming](https://mp.weixin.qq.com/s/AKqVWIdk3SBD-6uiTMliyA).
- [x] [Realtime SRT Streaming](https://mp.weixin.qq.com/s/HQb3gLRyJHHu56pnyHerxA).
- [x] [Automatical HTTPS](https://mp.weixin.qq.com/s/O70Fz-mxNedZpxgGXQ8DsA).
- [x] [Dashboard by Prometheus](https://mp.weixin.qq.com/s/ub9ZGmntOy_-S11oxFkxvg).
- [x] [DVR to Cloud Storage or VoD](https://mp.weixin.qq.com/s/UXR5EBKZ-LnthwKN_rlIjg).

Other more use scenarios is on the way, please read [this post](https://github.com/ossrs/srs/issues/2856#lighthouse).

## Architecture

The architecture of [srs-cloud](https://github.com/ossrs/srs-cloud#architecture) by 
[mermaid](https://mermaid.live/edit/#pako:eNpdkU1vwjAMhv9KlMMEEh87bJduQkIUhLQv1DIuLYe0MW1Hk1Spy0CI_74kZWxwqO3Yj_M67pGmigP16KZU32nONJLX4CmWhMiskHvS74-IyAR2rHlO9HCUHIg0LV9112F1k2SaVTkJg9AliI2CyBjyMLh3LSGypIT1X9l35cdz2YcdlKpq6yC5dVbOiVvQ-rlS27rNoAYmxg3m5I74q8DYlfJdsz1ZZPLRznK5ZTYTFWQu117Y_5-6YJ-VeQsHFwdQAqvhGpiAxkS18Xy5XNzILEGmIHFSqoa3_PgtMt_QDDQ0Q66v8YVWAjCHpnYTvZu1TveV0gj6GvRVur3NTeUu6gxA7twOUw3cKBesrAe4x-6NUgC8qKOOc-RlZcq0RwVowQpufv7RwjE1owiIqWdCDhvWlBjTWJ4M2lScIUx5gUpTb2NUoEdZgyo8yJR6qBv4hfyCmSWKM3X6AZm9vuQ)

```mermaid
flowchart LR;
  nginx --> mgmt(mgmt<br/>by nodejs);
  subgraph SRS;
    SRSR[SRS 4.0<br/>Stable];
    SRSD[SRS 5.0<br/>Develop];
  end
  mgmt --> SRS --> Hooks --> StreamAuth & DVR & VoD;
  DVR --> COS;
  mgmt --> FFmpeg;
  SRS --- FFmpeg;
  mgmt --> Upgrade --> Release;
  mgmt --> Certbot --> HTTPS;
  mgmt --> TencentCloud --> CAM[CAM/COS/VoD];
  mgmt --> Prometheus --- NodeExporter;
  mgmt --> Docker;
  mgmt --> Env[(.env<br/>credentials.txt)];
  mgmt --> Redis[(Redis KV)];
```

> Note: It's a single node, also light-weighted, video cloud for tiny company, personal user and starter.

## Ports

The ports allocated:

| Module | TCP Ports | UDP Ports | Notes |
| ------ | --------- | --------- | ----- |
| SRS | 1935, 1985, 8080,<br/> 8088, 1990, 554,<br/> 8936 | 8000, 8935, 10080,<br/> 1989 | See [SRS ports](https://github.com/ossrs/srs/blob/develop/trunk/doc/Resources.md#ports) |
| releases | 2023 |  - | Mount at `/terraform/v1/releases` |
| mgmt | 2022 |  - | Mount at `/mgmt/` and `/terraform/v1/mgmt/` |
| hooks | 2021 |  - | Mount at `/terraform/v1/hooks/` |
| tencent-cloud | 2020 |  - | Mount at `/terraform/v1/tencent/` |
| ffmpeg | 2019 |  - | Mount at `/terraform/v1/ffmpeg/` |
| prometheus | 9090 | - | Mount at `/prometheus` |
| node-exporter | 9100 | - | - |

## Features

The features that we're developing:

* [x] A mgmt support authentication and automatic updates.
* [x] Run SRS in docker, query status by docker and SRS API.
* [x] Support publish by RTMP/WebRTC, play by RTMP/HTTP-FLV/HLS/WebRTC.
* [x] SRS container use docker logs `json-file` and rotate for logging.
* [x] Support high-resolution and realtime(200~500ms) live streaming by SRT.
* [x] Run SRS hooks in docker, to callback by SRS server.
* [x] Support publish by SRT, play by RTMP/HTTP-FLV/HLS/WebRTC/SRT.
* [x] Integrate with prometheus and node-exporter.
* [x] Support DVR to tencent cloud storage, see [#1193](https://github.com/ossrs/srs/issues/1193).
* [x] Change redis port and use randomly password.
* [x] Support integrity with tencent cloud VoD.
* [x] Forward stream to multiple platforms, see [#2676](https://github.com/ossrs/srs/issues/2676).
* [ ] Support GB28181 by SRS 5.0 container.
* [ ] Support live streaming transcoding by FFmpeg, see [#2869](https://github.com/ossrs/srs/issues/2869).
* [ ] Support virtual live streaming, covert file or other resource to live.
* [ ] Support WebRTC face to face chat, see [#2857](https://github.com/ossrs/srs/issues/2857).
* [ ] Support WebRTC video chat room, see [#2924](https://github.com/ossrs/srs/issues/2924).
* [ ] Support a set of tools for developer, see [#2891](https://github.com/ossrs/srs/issues/2891).
* [ ] Collect logs of mgmt and containers together.
* [ ] Stop, restart and upgrade containers.
* [ ] Support logrotate to manage the logs.
* [ ] Enhance prometheus API with authentication.

## APIs

Platform:

* `/terraform/v1/mgmt/versions` Public version api.
* `/terraform/v1/mgmt/init` Whether mgmt initialized.
* `/terraform/v1/mgmt/check` Check whether system is ok.
* `/terraform/v1/mgmt/status` Query the version of mgmt.
* `/terraform/v1/mgmt/upgrade` Upgrade the mgmt to latest version.
* `/terraform/v1/mgmt/strategy` Toggle the upgrade strategy.
* `/terraform/v1/mgmt/token` System auth with token.
* `/terraform/v1/mgmt/login` System auth with password.
* `/terraform/v1/mgmt/ssl` Config the system SSL config.
* `/terraform/v1/mgmt/pubkey` Update the access for platform administrator pubkey.
* `/terraform/v1/mgmt/letsencrypt` Config the let's encrypt SSL.
* `/terraform/v1/mgmt/containers` Query and upgrade SRS container.
* `/terraform/v1/mgmt/bilibili` Query the video information.
* `/terraform/v1/mgmt/beian/query` Query the beian information.
* `/terraform/v1/mgmt/beian/update` Update the beian information.
* `/terraform/v1/mgmt/window/query` Query the upgrade time window.
* `/terraform/v1/mgmt/window/update` Update the upgrade time window.
* `/tools/` A set of H5 tools, like simple player.
* `/console/` The SRS console, serve by mgmt.
* `/players/` The SRS player, serve by mgmt.
* `/terraform/v1/releases` Version management for all components.
* `/.well-known/acme-challenge/` HTTPS verify mount for letsencrypt.

Market:

* `/terraform/v1/hooks/srs/verify` Hooks: Verify the stream request URL of SRS.
* `/terraform/v1/hooks/srs/secret/query` Hooks: Query the secret to generate stream URL.
* `/terraform/v1/hooks/srs/secret/update` Hooks: Update the secret to generate stream URL.
* `/terraform/v1/hooks/srs/hls` Hooks: Handle the `on_hls` event.
* `/terraform/v1/hooks/dvr/apply` Hooks: Apply the DVR pattern.
* `/terraform/v1/hooks/dvr/query` Hooks: Query the DVR pattern.
* `/terraform/v1/hooks/dvr/files` Hooks: List the DVR files.
* `/terraform/v1/hooks/dvr/hls` Hooks: Generate HLS/m3u8 url to preview or download.
* `/terraform/v1/hooks/vod/query` Hooks: Query the VoD pattern.
* `/terraform/v1/hooks/vod/apply` Hooks: Apply the VoD pattern.
* `/terraform/v1/hooks/vod/files` Hooks: List the VoD files.
* `/terraform/v1/hooks/vod/hls` Hooks: Generate HLS/m3u8 url to preview or download.
* `/terraform/v1/tencent/cam/secret` Tencent: Setup the CAM SecretId and SecretKey.
* `/terraform/v1/ffmpeg/forward/secret` FFmpeg: Setup the forward secret to live streaming platforms.
* `/terraform/v1/ffmpeg/forward/streams` FFmpeg: Query the forwarding streams.
* `/prometheus` Prometheus: Time-series database and monitor.
* `/api/` SRS: HTTP API of SRS media server.
* `/rtc/` SRS: HTTP API for WebERTC of SRS media server.
* `/*/*.(flv|m3u8|ts|aac|mp3)` SRS: Media stream for HTTP-FLV, HLS, HTTP-TS, HTTP-AAC, HTTP-MP3.

## Depends

The software we depend on:

* Docker, `yum install -y docker`
* Redis, `yum install -y redis`
* Nginx, `yum install -y nginx`
  * SSL: `/etc/nginx/ssl`
* [Certbot](https://github.com/ossrs/srs/issues/2864#lets-encrypt), `docker --name certbot`
  * Verify webroot: `mgmt/containers/www/.well-known/acme-challenge/`
  * Cert files: `mgmt/containers/etc/letsencrypt/live/`
* [SRS](https://github.com/ossrs/srs), `docker --name srs-server`
  * Config: `mgmt/containers/conf/srs.conf` mount as `/usr/local/srs/conf/lighthouse.conf`
  * Volume: `mgmt/containers/objs/nginx/html` mount as `/usr/local/srs/objs/nginx/html`
* [srs-hooks](https://github.com/ossrs/srs-cloud/tree/lighthouse/hooks), `docker --name srs-hooks`
  * Volume: `mgmt/containers/objs/nginx/html` mount as `/usr/local/mgmt/containers/objs/nginx/html`
* [tencent-cloud](https://github.com/ossrs/srs-cloud/tree/lighthouse/tencent), `docker --name tencent-cloud`
  * [CAM](https://console.cloud.tencent.com/cam/overview) Authentication by secretId and secretKey.
* [ffmpeg](https://github.com/ossrs/srs-cloud/tree/lighthouse/ffmpeg), `docker --name ffmpeg`
  * [FFmpeg and ffprobe](https://ffmpeg.org) tools in `ossrs/srs:node-av`
* [Prometheus](https://github.com/prometheus/prometheus#install), `docker --name prometheus`
  * Config: `mgmt/containers/conf/prometheus.yml`
  * Data directory: `mgmt/containers/data/prometheus`
* [NodeExporter](https://github.com/prometheus/node_exporter), `docker --name node-exporter`

## Upgrade Workflow

When upgrading automatically or manually by user:

* `bash upgrade` for each upgrade.

When system start, check the flag `SRS_FIRST_BOOT` in redis, if not set:

* `bash auto/upgrade_prepare` do upgrade for previous images.
* Restart container srs-server and srs-hooks, for config changed.

They are not mutually exclusive.

## System Boot

When system boot:

* Restart the mgmt service by `systemctl start srs-terraform`
* Execute script `bootstrap` at mgmt
* Run script `auto/foreach_run` at mgmt
* Start application by `node .` at mgmt

## System Setup

When user setup the system, the admin password for the first boot:

* Setup the `MGMT_PASSWORD` in `.env`
* Restat all containers that depends on `.env`

## Environments

The optional environments defined by `mgmt/.env`:

* `MGMT_PASSWORD`: The mgmt administrator password.
* `REGION`: `ap-guangzhou|ap-singapore`, the region for upgrade source.
* `PLATFORM`: The platform name.

For testing the specified service:

* `NODE_ENV`: `development|production`, if development, use local redis; otherwise, use `mgmt.srs.local` in docker.
* `LOCAL_RELEASE`: `true|false`, whether use local release service.

For github actions to control the containers:

* `SRS_DOCKER`: `srs` to enfore use `ossrs/srs` docker image.
* `USE_DOCKER`: `true|false`, if false, disable all docker containers.

For mgmt and containers to connect to redis:

* `REDIS_PASSWORD`: The redis password.
* `REDIS_PORT`: The redis port.

## Develop

Install dependencies:

```bash
cd mgmt && npm install
```

Run the mgmt backend:

```
cd mgmt
npm start
```

Run the mgmt react ui:

```
cd mgmt/ui
npm start
```

Access the browser: http://localhost:3000

