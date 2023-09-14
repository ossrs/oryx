# SRS-Stack

[![](https://img.shields.io/twitter/follow/srs_server?style=social)](https://twitter.com/srs_server)
[![](https://badgen.net/discord/members/bQUPDRqy79)](https://discord.gg/bQUPDRqy79)
[![](https://ossrs.net/wiki/images/wechat-badge4.svg)](https://ossrs.net/lts/zh-cn/contact#discussion)
[![](https://ossrs.net/wiki/images/do-btn-srs-125x20.svg)](https://cloud.digitalocean.com/droplets/new?appId=133468816&size=s-2vcpu-2gb&region=sgp1&image=ossrs-srs&type=applications)
[![](https://opencollective.com/srs-server/tiers/badge.svg)](https://opencollective.com/srs-server)

SRS Stack makes it easy for you to create an online video service with just one click, either on cloud 
platforms like DigitalOcean or AWS or by self-hosting. This open-source, simple video solution is made 
using Go, Reactjs, SRS, FFmpeg, and WebRTC. It supports protocols like RTMP, WebRTC, HLS, HTTP-FLV, and 
SRT. It offers features like authentication, streaming on multiple platforms, recording, virtual live 
events, automatic HTTPS, and an easy-to-use HTTP Open API.

[![](https://ossrs.io/lts/en-us/img/SRS-Stack-5-sd.png?v=1)](https://ossrs.io/lts/en-us/img/SRS-Stack-5-hd.png)

> Note: For more details on the SRS Stack, please visit the following [link](https://www.figma.com/file/Ju5h2DZeJMzUtx5k7D0Oak/SRS-Stack).

## Usage

Run srs-stack in one docker:

```bash
docker run --rm -it -p 2022:2022 -p 2443:2443 -p 1935:1935 \
  -p 8080:8080 -p 8000:8000/udp -p 10080:10080/udp --name srs-stack \
  -v $HOME/data:/data ossrs/srs-stack:5
```

> Note: Please use `-e REACT_APP_LOCALE=zh` and `registry.cn-hangzhou.aliyuncs.com/ossrs/srs-stack:5` in China.

> Note: All data will be saved in `$HOME/data` directory, please change it to your directory.

Then open http://localhost:2022 in browser.

You have the option to modify the volumes for srs-stack and direct them to different directories.

* `/data` The global data directory.
    * `.well-known` The directory for Let's Encrypt ACME challenge.
    * `config` The .env for password, srs/redis/nginx/prometheus config, and SSL files.
    * `dvr` The dvr storage directory, save dvr files.
    * `lego` The LEGO Let's Encrypt ACME challenge directory.
    * `record` The record storage directory, save record files.
    * `redis` The redis data directory, the publish secret and record configuration.
    * `signals` The signals storage directory, save signal files.
    * `upload` The upload storage directory, save upload files.
    * `vlive` The storage directory for virtual live, save video files.
    * `vod` The storage directory for VoD, save video files.
    * `nginx-cache` The storage directory for nginx cache, save cache files.

You can use environment variables to modify the settings.

* `MGMT_PASSWORD`: The mgmt administrator password.
* `REACT_APP_LOCALE`: The i18n config for ui, `en` or `zh`, default to `en`.

> Note: The `MGMT_PASSWORD` is also saved in `/data/config/.env`, you can modify it by yourself.

To access additional environment variables, please refer to the [Environments](DEVELOPER.md#environments) section.

## Sponsor

Would you like additional assistance from us? By becoming a sponsor or backer of SRS, we can provide you
with the support you need:

* Backer: $5 per month, online text chat support through Discord.
* Sponsor: $100 per month, online meeting support, 1 meeting per month in 1 hour.

Please visit [OpenCollective](https://opencollective.com/srs-server) to become a backer or sponsor, and send
us a direct message on [Discord](https://discord.gg/bQUPDRqy79). We are currently providing support to the 
developers listed below:

[![](https://opencollective.com/srs-server/backers.svg?width=800&button=false)](https://opencollective.com/srs-server)

We at SRS aim to establish a non-profit, open-source community that assists developers worldwide in creating
your own high-quality streaming and RTC platforms to support your businesses.

## FAQ

1. [English FAQ](https://ossrs.io/lts/en-us/faq-srs-stack)
1. [中文 FAQ](https://ossrs.net/lts/zh-cn/faq-srs-stack)

## Tutorials

- [x] Getting Started: [Blog](https://blog.ossrs.io/how-to-setup-a-video-streaming-service-by-1-click-e9fe6f314ac6), [EN](https://ossrs.io/lts/en-us/docs/v6/doc/getting-started-stack), [CN](https://ossrs.net/lts/zh-cn/docs/v5/doc/getting-started-stack).
- [x] Support WordPress Plugin: [Blog](https://blog.ossrs.io/publish-your-srs-livestream-through-wordpress-ec18dfae7d6f), [EN](https://ossrs.io/lts/en-us/blog/WordPress-Plugin), [CN](https://ossrs.net/lts/zh-cn/blog/WordPress-Plugin) or [WordPress Plugin](https://wordpress.org/plugins/srs-player).
- [x] Support Automatic HTTPS: [Blog](https://blog.ossrs.io/how-to-secure-srs-with-lets-encrypt-by-1-click-cb618777639f), [EN](https://ossrs.io/lts/en-us/blog/SRS-Stack-Tutorial), [CN](https://ossrs.net/lts/zh-cn/blog/SRS-Stack-HTTPS).
- [x] Support aaPanel to install on any linux: [Blog](https://blog.ossrs.io/how-to-setup-a-video-streaming-service-by-aapanel-9748ae754c8c), [EN](https://ossrs.io/lts/en-us/blog/BT-aaPanel), [CN](https://ossrs.net/lts/zh-cn/blog/BT-aaPanel).
- [x] Support DVR to local disk: [Blog](https://blog.ossrs.io/how-to-record-live-streaming-to-mp4-file-2aa792c35b25), [EN](https://ossrs.io/lts/en-us/blog/Record-Live-Streaming), [CN](https://mp.weixin.qq.com/s/axN_TPo-Gk_H7CbdqUud6g).
- [x] Support Virtual Live Streaming: [CN](https://mp.weixin.qq.com/s/I0Kmxtc24txpngO-PiR_tQ).
- [x] Support build small [HLS deliver CDN](https://github.com/ossrs/srs-stack/tree/main/scripts/nginx-hls-cdn) by Nginx.
- [x] Support Live Streaming: [CN](https://mp.weixin.qq.com/s/AKqVWIdk3SBD-6uiTMliyA).
- [x] Support Realtime SRT Streaming: [CN](https://mp.weixin.qq.com/s/HQb3gLRyJHHu56pnyHerxA).
- [x] Support DVR to Tencent Cloud Storage or VoD: [CN](https://mp.weixin.qq.com/s/UXR5EBKZ-LnthwKN_rlIjg).
- [x] Support Typecho Plugin: [CN](https://github.com/ossrs/Typecho-Plugin-SrsPlayer).

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
- [x] Change redis port and use randomly password.
- [x] Support integrity with tencent cloud VoD.
- [x] Forward stream to multiple platforms.
- [x] Support WordPress Plugin: SrsPlayer.
- [x] Support aaPanel to install on any linux.
- [x] Support DVR to local disk.
- [x] Support upgrade to latest version manually.
- [x] Support HTTPS by let's encrypt with LEGO.
- [x] Support virtual live streaming, covert file or other resource to live.
- [x] Support self-host HLS CDN, to serve 10k+ viewers.
- [x] Support Typecho Plugin: Typecho-Plugin-SrsPlayer.
- [x] Support DVR to TencentCloud storage.
- [ ] Support limit the streaming duration to limit the fee.
- [ ] Support GB28181 by SRS 5.0 container.
- [ ] Support live streaming transcoding by FFmpeg, see [#2869](https://github.com/ossrs/srs/issues/2869).
- [ ] Support WebRTC face to face chat, see [#2857](https://github.com/ossrs/srs/issues/2857).
- [ ] Support WebRTC video chat room, see [#2924](https://github.com/ossrs/srs/issues/2924).
- [ ] Support a set of tools for developer, see [#2891](https://github.com/ossrs/srs/issues/2891).
- [ ] Collect logs of mgmt and containers together.
- [ ] Stop, restart and upgrade containers.
- [ ] Support logrotate to manage the logs.
- [ ] Enhance prometheus API with authentication.
- [ ] Integrate with prometheus and node-exporter.

## License

SRS Stack is an open-source project, licensed under the [AGPL-3.0-or-later](https://spdx.org/licenses/AGPL-3.0-or-later.html) license.

## Developer

For development, please refer to the [Environments](DEVELOPER.md) about the API and architecture.

2022.11

