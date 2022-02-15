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
* `/terraform/v1/mgmt/strategy` Toggle the upgrade strategy.
* `/terraform/v1/mgmt/token` System auth with token.
* `/terraform/v1/mgmt/login` System auth with password.
* `/terraform/v1/mgmt/ssl` Config the system SSL config.
* `/terraform/v1/mgmt/pubkey` Update the access for platform administrator pubkey.
* `/terraform/v1/mgmt/letsencrypt` Config the let's encrypt SSL.
* `/terraform/v1/mgmt/containers` Query and upgrade SRS container.
* `/terraform/v1/mgmt/bilibili` Query the video information.

Releases:

* `/terraform/v1/releases` Version management for all components.

Market:

* `/terraform/v1/hooks/srs/verify` Hooks: Verify the stream request URL of SRS.
* `/terraform/v1/hooks/srs/secret/query` Hooks: Query the secret to generate stream URL.
* `/terraform/v1/hooks/srs/secret/update` Hooks: Update the secret to generate stream URL.
* `/prometheus` Prometheus: Time-series database and monitor.

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
  * Config: `mgmt/containers/conf/srs.conf`
* [srs-hooks](https://github.com/ossrs/srs-terraform/tree/lighthouse/hooks), `docker --name srs-hooks`
* [Prometheus](https://github.com/prometheus/prometheus#install), `docker --name prometheus`
  * Config: `mgmt/containers/conf/prometheus.yml`
  * Data directory: `mgmt/containers/data/prometheus`
* [NodeExporter](https://github.com/prometheus/node_exporter), `docker --name node-exporter`

