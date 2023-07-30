# SRS-Cloud

[![](https://img.shields.io/twitter/follow/srs_server?style=social)](https://twitter.com/srs_server)
[![](https://badgen.net/discord/members/bQUPDRqy79)](https://discord.gg/bQUPDRqy79)
[![](https://ossrs.net/wiki/images/wechat-badge4.svg)](https://ossrs.net/lts/zh-cn/contact#discussion)
[![](https://ossrs.net/wiki/images/do-btn-srs-125x20.svg)](https://cloud.digitalocean.com/droplets/new?appId=133468816&size=s-1vcpu-512mb-10gb&region=sgp1&image=ossrs-srs&type=applications)

A lightweight open-source video cloud based on Nodejs, SRS, FFmpeg, WebRTC, etc.

## Usage

Run srs-cloud in one docker:

```bash
docker run --rm -it -p 2022:2022 -p 1935:1935/tcp -p 1985:1985/tcp \
  -p 8080:8080/tcp -p 8000:8000/udp -p 10080:10080/udp --name srs-cloud \
  -v $HOME/db:/data ossrs/srs-cloud:1
```

> Note: Please use `registry.cn-hangzhou.aliyuncs.com/ossrs/srs-cloud:1` in China.

> Note: All data will be saved in `$HOME/db` directory, please change it to your directory.

Then open http://localhost:2022 in browser.

You have the option to modify the volumes for srs-cloud and direct them to different directories.

* `/data` The global data directory.
    * `redis` The redis data directory, the publish secret and record configuration.
    * `config` The mgmt password and cloud configuration.
    * `record` The record storage directory, save record files.
    * `vlive` The storage directory for virtual live, save video files.

You can use environment variables to modify the settings.

* `MGMT_PASSWORD`: The mgmt administrator password.
* `REACT_APP_LOCALE`: The i18n config for ui, `en` or `zh`, default to `zh`.

> Note: Kindly utilize the command such as `docker run -e REACT_APP_LOCALE=en` to modify the environment.

To access additional environment variables, please refer to the [Environments](DEVELOPER.md#environments) section.

## FAQ

1. [English FAQ](https://ossrs.io/lts/en-us/faq-srs-cloud)
1. [中文 FAQ](https://ossrs.net/lts/zh-cn/faq-srs-cloud)

## Tutorials

- [x] [Getting Started](https://mp.weixin.qq.com/s/fWmdkw-2AoFD_pEmE_EIkA).
- [x] [Live Streaming](https://mp.weixin.qq.com/s/AKqVWIdk3SBD-6uiTMliyA).
- [x] [Realtime SRT Streaming](https://mp.weixin.qq.com/s/HQb3gLRyJHHu56pnyHerxA).
- [x] [DVR to Cloud Storage or VoD](https://mp.weixin.qq.com/s/UXR5EBKZ-LnthwKN_rlIjg).
- [x] [Support WordPress Plugin](https://mp.weixin.qq.com/s/YjTkcJLkErMcZYHIjzsW_w) or [here](https://wordpress.org/plugins/srs-player).
- [x] [Support Typecho Plugin](https://github.com/ossrs/Typecho-Plugin-SrsPlayer).
- [x] [Support aaPanel to install on any linux](https://github.com/ossrs/srs-cloud/issues/29).
- [x] [Support DVR to local disk](https://github.com/ossrs/srs-cloud/issues/42).
- [x] [Support Virtual Live Streaming](https://mp.weixin.qq.com/s/I0Kmxtc24txpngO-PiR_tQ).
- [ ] [Automatical HTTPS](https://mp.weixin.qq.com/s/O70Fz-mxNedZpxgGXQ8DsA).
- [ ] [Dashboard by Prometheus](https://mp.weixin.qq.com/s/ub9ZGmntOy_-S11oxFkxvg).

Other more use scenarios is on the way, please read [this post](https://github.com/ossrs/srs/issues/2856#lighthouse).

## Features

The features that we're developing:

- [x] A mgmt support authentication and automatic updates.
- [x] Run SRS in docker, query status by docker and SRS API.
- [x] Support publish by RTMP/WebRTC, play by RTMP/HTTP-FLV/HLS/WebRTC.
- [x] SRS container use docker logs `json-file` and rotate for logging.
- [x] Support high-resolution and realtime(200~500ms) live streaming by SRT.
- [x] Run SRS hooks in docker, to callback by SRS server.
- [x] Support publish by SRT, play by RTMP/HTTP-FLV/HLS/WebRTC/SRT.
- [x] Support DVR to tencent cloud storage, see [#1193](https://github.com/ossrs/srs/issues/1193).
- [x] Change redis port and use randomly password.
- [x] Support integrity with tencent cloud VoD.
- [x] Forward stream to multiple platforms, see [#2676](https://github.com/ossrs/srs/issues/2676).
- [x] [Support WordPress Plugin](https://mp.weixin.qq.com/s/YjTkcJLkErMcZYHIjzsW_w) or [here](https://wordpress.org/plugins/srs-player).
- [x] [Support Typecho Plugin](https://github.com/ossrs/Typecho-Plugin-SrsPlayer).
- [x] [Support aaPanel to install on any linux](https://github.com/ossrs/srs-cloud/issues/29).
- [x] [Support DVR to local disk](https://github.com/ossrs/srs-cloud/issues/42).
- [ ] Support upgrade to latest version manually.
- [ ] Support HTTPS by let's encrypt with certbot.
- [ ] Support GB28181 by SRS 5.0 container.
- [ ] Support live streaming transcoding by FFmpeg, see [#2869](https://github.com/ossrs/srs/issues/2869).
- [ ] Support virtual live streaming, covert file or other resource to live.
- [ ] Support WebRTC face to face chat, see [#2857](https://github.com/ossrs/srs/issues/2857).
- [ ] Support WebRTC video chat room, see [#2924](https://github.com/ossrs/srs/issues/2924).
- [ ] Support a set of tools for developer, see [#2891](https://github.com/ossrs/srs/issues/2891).
- [ ] Collect logs of mgmt and containers together.
- [ ] Stop, restart and upgrade containers.
- [ ] Support logrotate to manage the logs.
- [ ] Enhance prometheus API with authentication.
- [ ] Integrate with prometheus and node-exporter.

## Developer

For development, please refer to the [Environments](DEVELOPER.md) about the API and architecture.

2022.11

