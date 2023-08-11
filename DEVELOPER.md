# Developer

This guide is for developers and covers topics such as OpenAPI, environment variables, 
resources, and ports, as well as development on Mac or using Docker.

## Develop All in macOS

Start redis and SRS by docker:

```bash
docker rm -f redis srs 2>/dev/null &&
docker run --name redis --rm -it -v $HOME/db/redis:/data -p 6379:6379 -d redis &&
docker run --name srs --rm -it \
    -v $(pwd)/platform/containers/conf/srs.release-mac.conf:/usr/local/srs/conf/srs.conf \
    -v $(pwd)/platform/containers/objs/nginx:/usr/local/srs/objs/nginx \
    -p 1935:1935/tcp -p 1985:1985/tcp -p 8080:8080/tcp -p 8000:8000/udp -p 10080:10080/udp \
    -d ossrs/srs:5
```

> Note: Stop service by `docker rm -f redis srs`

> Note: Also, you can run SRS by `(cd platform && ~/git/srs/trunk/objs/srs -c containers/conf/srs.release-local.conf)`

Run the platform backend, or run in GoLand:

```bash
(cd platform && go run .)
```

Run the platform react ui, or run in WebStorm:

```bash
(cd ui && npm install && npm start)
```

Access the browser: http://localhost:3000

## Develop the Script Installer

Build a docker image:

```bash
docker rm -f script 2>/dev/null &&
docker rmi srs-script-dev 2>/dev/null &&
docker build -t srs-script-dev -f scripts/setup-ubuntu/Dockerfile.script .
```

Create a docker container in daemon:

```bash
docker rm -f script 2>/dev/null &&
docker run -p 2022:2022 -p 1935:1935/tcp -p 1985:1985/tcp \
    -p 8080:8080/tcp -p 8000:8000/udp -p 10080:10080/udp \
    --privileged -v /sys/fs/cgroup:/sys/fs/cgroup:rw --cgroupns=host \
    -d --rm -it -v $(pwd):/g -w /g --name=script srs-script-dev
```

> Note: For Linux server, please use `--privileged -v /sys/fs/cgroup:/sys/fs/cgroup:ro` to start docker.

Build and save the script image to file:

```bash
docker rmi platform:latest 2>/dev/null || echo OK &&
docker build -t platform:latest -f Dockerfile . &&
docker save -o platform.tar platform:latest
```

Enter the docker container:

```bash
version=$(bash scripts/version.sh) &&
docker exec -it script docker load -i platform.tar && 
docker exec -it script docker tag platform:latest ossrs/srs-stack:$version &&
docker exec -it script docker tag platform:latest registry.cn-hangzhou.aliyuncs.com/ossrs/srs-stack:$version &&
docker exec -it script docker images
```

Test the build script, in the docker container:

```bash
docker exec -it bt rm -rf /data/* &&
docker exec -it script bash build/srs_stack/scripts/setup-ubuntu/uninstall.sh || echo OK &&
bash scripts/setup-ubuntu/build.sh --output $(pwd)/build --extract &&
docker exec -it script bash build/srs_stack/scripts/setup-ubuntu/install.sh --verbose
```

Run test for script:

```bash
docker exec -it script make -j -C test &&
docker exec -it script ./test/srs-stack.test -test.v -endpoint http://localhost:2022 \
    -srs-log=true -wait-ready=true -init-password=true \
    -check-api-secret=false -test.run TestApi_Empty &&
bash scripts/tools/secret.sh --output test/.env &&
docker exec -it script ./test/srs-stack.test -test.v -wait-ready -endpoint http://localhost:2022 \
    -srs-log=true -wait-ready=true -init-password=false \
    -check-api-secret=true \
    -test.parallel 8
```

Access the browser: [http://localhost:2022](http://localhost:2022)

## Develop the aaPanel Plugin

Start a container and mount as plugin:

```bash
docker rm -f bt aapanel 2>/dev/null &&
AAPANEL_KEY=$(cat $HOME/.bt/api.json |awk -F token_crypt '{print $2}' |cut -d'"' -f3)
docker run -p 80:80 -p 7800:7800 \
    -p 1935:1935/tcp -p 1985:1985/tcp -p 8080:8080/tcp -p 8000:8000/udp -p 10080:10080/udp \
    -v $(pwd)/build/srs_stack:/www/server/panel/plugin/srs_stack \
    -v $HOME/.bt/api.json:/www/server/panel/config/api.json -e BT_KEY=$AAPANEL_KEY \
    --privileged -v /sys/fs/cgroup:/sys/fs/cgroup:rw --cgroupns=host \
    --add-host srs.stack.local:127.0.0.1 \
    -d --rm -it -v $(pwd):/g -w /g --name=aapanel ossrs/aapanel-plugin-dev:1
```

> Note: For Linux server, please use `--privileged -v /sys/fs/cgroup:/sys/fs/cgroup:ro` to start docker.

> Note: Enable the [HTTP API](https://www.bt.cn/bbs/thread-20376-1-1.html) and get the `api.json`,
> and save it to `$HOME/.bt/api.json`.

Build and save the platform image to file:

```bash
docker rmi platform:latest 2>/dev/null || echo OK &&
docker build -t platform:latest -f Dockerfile . &&
docker save -o platform.tar platform:latest
```

Enter the docker container:

```bash
version=$(bash scripts/version.sh) &&
docker exec -it aapanel docker load -i platform.tar && 
docker exec -it aapanel docker tag platform:latest ossrs/srs-stack:$version &&
docker exec -it aapanel docker tag platform:latest registry.cn-hangzhou.aliyuncs.com/ossrs/srs-stack:$version &&
docker exec -it aapanel docker images
```

Next, build the aaPanel plugin and install it:

```bash
docker exec -it aapanel rm -rf /data/* &&
docker exec -it aapanel bash /www/server/panel/plugin/srs_stack/install.sh uninstall || echo OK &&
bash scripts/setup-aapanel/auto/zip.sh --output $(pwd)/build --extract &&
docker exec -it aapanel bash /www/server/panel/plugin/srs_stack/install.sh install
```

You can use aaPanel panel to install the plugin, or by command:

```bash
docker exec -it aapanel python3 /www/server/panel/plugin/srs_stack/bt_api_remove_site.py &&
docker exec -it aapanel python3 /www/server/panel/plugin/srs_stack/bt_api_create_site.py &&
docker exec -it aapanel python3 /www/server/panel/plugin/srs_stack/bt_api_setup_site.py &&
docker exec -it aapanel bash /www/server/panel/plugin/srs_stack/setup.sh \
    --r0 /tmp/srs_stack_install.r0 --nginx /www/server/nginx/logs/nginx.pid \
    --www /www/wwwroot --site srs.stack.local
```

Run test for aaPanel:

```bash
docker exec -it aapanel make -j -C test &&
docker exec -it aapanel ./test/srs-stack.test -test.v -endpoint http://srs.stack.local:80 \
    -srs-log=true -wait-ready=true -init-password=true \
    -check-api-secret=false -test.run TestApi_Empty &&
bash scripts/tools/secret.sh --output test/.env &&
docker exec -it aapanel ./test/srs-stack.test -test.v -wait-ready -endpoint http://srs.stack.local:80 \
    -srs-log=true -wait-ready=true -init-password=false \
    -check-api-secret=true \
    -test.parallel 8
```

Open [http://localhost:7800/srscloud](http://localhost:7800/srscloud) to install plugin.

> Note: Or you can use `docker exec -it aapanel bt default` to show the login info.

In the [application store](http://localhost:7800/soft), there is a `srs_stack` plugin. After test, you can install the plugin
`build/aapanel-srs_stack.zip` to production aaPanel panel.

## Develop the BT Plugin

Start a container and mount as plugin:

```bash
docker rm -f bt aapanel 2>/dev/null &&
BT_KEY=$(cat $HOME/.bt/api.json |awk -F token_crypt '{print $2}' |cut -d'"' -f3)
docker run -p 80:80 -p 7800:7800 \
    -p 1935:1935/tcp -p 1985:1985/tcp -p 8080:8080/tcp -p 8000:8000/udp -p 10080:10080/udp \
    -v $(pwd)/build/srs_stack:/www/server/panel/plugin/srs_stack \
    -v $HOME/.bt/userInfo.json:/www/server/panel/data/userInfo.json \
    -v $HOME/.bt/api.json:/www/server/panel/config/api.json -e BT_KEY=$BT_KEY \
    --privileged -v /sys/fs/cgroup:/sys/fs/cgroup:rw --cgroupns=host \
    --add-host srs.stack.local:127.0.0.1 \
    -d --rm -it -v $(pwd):/g -w /g --name=bt ossrs/bt-plugin-dev:1
```

> Note: For Linux server, please use `--privileged -v /sys/fs/cgroup:/sys/fs/cgroup:ro` to start docker.

> Note: Should bind the docker to your BT account, then you will get the `userInfo.json`, 
> and save it to `$HOME/.bt/userInfo.json`.

> Note: Enable the [HTTP API](https://www.bt.cn/bbs/thread-20376-1-1.html) and get the `api.json`, 
> and save it to `$HOME/.bt/api.json`.

Build and save the platform image to file:

```bash
docker rmi platform:latest 2>/dev/null || echo OK &&
docker build -t platform:latest -f Dockerfile . &&
docker save -o platform.tar platform:latest
```

Enter the docker container:

```bash
version=$(bash scripts/version.sh) &&
docker exec -it bt docker load -i platform.tar && 
docker exec -it bt docker tag platform:latest ossrs/srs-stack:$version &&
docker exec -it bt docker tag platform:latest registry.cn-hangzhou.aliyuncs.com/ossrs/srs-stack:$version &&
docker exec -it bt docker images
```

Next, build the BT plugin and install it:

```bash
docker exec -it bt rm -rf /data/* &&
docker exec -it bt bash /www/server/panel/plugin/srs_stack/install.sh || echo OK &&
bash scripts/setup-bt/auto/zip.sh --output $(pwd)/build --extract &&
docker exec -it bt bash /www/server/panel/plugin/srs_stack/install.sh install
```

You can use BT panel to install the plugin, or by command:

```bash
docker exec -it bt python3 /www/server/panel/plugin/srs_stack/bt_api_remove_site.py &&
docker exec -it bt python3 /www/server/panel/plugin/srs_stack/bt_api_create_site.py &&
docker exec -it bt python3 /www/server/panel/plugin/srs_stack/bt_api_setup_site.py &&
docker exec -it bt bash /www/server/panel/plugin/srs_stack/setup.sh \
    --r0 /tmp/srs_stack_install.r0 --nginx /www/server/nginx/logs/nginx.pid \
    --www /www/wwwroot --site srs.stack.local
```

Run test for BT:

```bash
docker exec -it bt make -j -C test &&
docker exec -it bt ./test/srs-stack.test -test.v -endpoint http://srs.stack.local:80 \
    -srs-log=true -wait-ready=true -init-password=true \
    -check-api-secret=false -test.run TestApi_Empty &&
bash scripts/tools/secret.sh --output test/.env &&
docker exec -it bt ./test/srs-stack.test -test.v -wait-ready -endpoint http://srs.stack.local:80 \
      -srs-log=true -wait-ready=true -init-password=false \
      -check-api-secret=true \
      -test.parallel 8
```

Open [http://localhost:7800/srscloud](http://localhost:7800/srscloud) to install plugin.

> Note: Or you can use `docker exec -it bt bt default` to show the login info.

In the [application store](http://localhost:7800/soft), there is a `srs_stack` plugin. After test, you can install the plugin 
`build/bt-srs_stack.zip` to production BT panel.

## Develop the Droplet Image

To build SRS droplet image for [DigitalOcean Marketplace](https://marketplace.digitalocean.com/).

For the first run, please [install Packer](https://www.packer.io/intro/getting-started/install.html) and plugin:

```bash
brew tap hashicorp/tap &&
brew install hashicorp/tap/packer &&
PACKER_LOG=1 packer plugins install github.com/digitalocean/digitalocean v1.1.1
```

Start to build SRS image by:

```bash
(export DIGITALOCEAN_TOKEN=$(grep market "${HOME}/Library/Application Support/doctl/config.yaml" |grep -v context |awk '{print $2}') &&
cd scripts/setup-droplet && packer build srs.json)
```

> Note: You can also create a [token](https://cloud.digitalocean.com/account/api/tokens) and setup the env `DIGITALOCEAN_TOKEN`.

Please check the [snapshot](https://cloud.digitalocean.com/images/snapshots/droplets).

## Develop the TencentCloud Lighthouse Image

To build SRS image for [TencentCloud Lighthouse](https://cloud.tencent.com/product/lighthouse).

For the first run, please create a [TencentCloud Secret](https://console.cloud.tencent.com/cam/capi) and save 
to `~/.lighthouse/.env` file:

```bash
LH_ACCOUNT=xxxxxx
SECRET_ID=xxxxxx
SECRET_KEY=xxxxxx
```

> Note: Share the image to `LH_ACCOUNT` to publish it.

Create a CVM instance:

```bash
rm -f /tmp/lh-*.txt &&
echo $(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 32) >/tmp/lh-token.txt &&
VM_TOKEN=$(cat /tmp/lh-token.txt) bash scripts/tools/tencent-cloud/helper.sh create-cvm.py --id /tmp/lh-instance.txt
bash scripts/tools/tencent-cloud/helper.sh query-cvm-ip.py --instance $(cat /tmp/lh-instance.txt) --id /tmp/lh-ip.txt &&
echo "Instance: $(cat /tmp/lh-instance.txt), IP: ubuntu@$(cat /tmp/lh-ip.txt), Password: $(cat /tmp/lh-token.txt)" && sleep 5 &&
bash scripts/setup-lighthouse/build.sh --ip $(cat /tmp/lh-ip.txt) --os ubuntu --user ubuntu --password $(cat /tmp/lh-token.txt) &&
bash scripts/tools/tencent-cloud/helper.sh create-image.py --instance $(cat /tmp/lh-instance.txt) --id /tmp/lh-image.txt &&
bash scripts/tools/tencent-cloud/helper.sh share-image.py --image $(cat /tmp/lh-image.txt) &&
echo "Image: $(cat /tmp/lh-image.txt) created and shared." &&
bash scripts/tools/tencent-cloud/helper.sh remove-cvm.py --instance $(cat /tmp/lh-instance.txt)
```

Next, create a test CVM instance with the image:

```bash
echo $(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 32) >/tmp/lh-token2.txt &&
VM_TOKEN=$(cat /tmp/lh-token2.txt) bash scripts/tools/tencent-cloud/helper.sh create-verify.py --image $(cat /tmp/lh-image.txt) --id /tmp/lh-test.txt &&
bash scripts/tools/tencent-cloud/helper.sh query-cvm-ip.py --instance $(cat /tmp/lh-test.txt) --id /tmp/lh-ip2.txt && 
echo "IP: ubuntu@$(cat /tmp/lh-ip2.txt), Password: $(cat /tmp/lh-token2.txt)" &&
echo "http://$(cat /tmp/lh-ip2.txt)"
```

Verify then cleanup the test CVM instance:

```bash
bash scripts/tools/tencent-cloud/helper.sh remove-cvm.py --instance $(cat /tmp/lh-test.txt)
```

After publish to lighthouse, cleanup the CVM, disk images, and snapshots:

```bash
bash scripts/tools/tencent-cloud/helper.sh remove-image.py --image $(cat /tmp/lh-image.txt)
```

# Tips

## Release

Release bugfix:

* For mgmt: `./auto/platform.sh`
* Then test the specified version of mgmt.

> Note: The [features](https://github.com/ossrs/srs-stack/issues/4) might need to be updated.

Release version for BT and aaPanel:

* Then run `./auto/release.sh`
* Finally, download [bt-srs_stack.zip](https://github.com/ossrs/srs-stack/releases) then submit to [bt.cn](https://www.bt.cn/developer/details.html?id=600801805)

> Note: The [BT forum](https://www.bt.cn/bbs/thread-90890-1-1.html) and [FAQ](https://github.com/ossrs/srs-stack/issues/4) might need to be updated.

To refresh current tag for mgmt and platform:

* Run `./auto/platform.sh -refresh`

The upgrade feature has been disabled, which means we no longer update the
[version API](https://api.ossrs.net/terraform/v1/releases), nor do we update
the `releases/version.go` file, and we don't use `./auto/releases_pub.sh`.

## Ports

The ports allocated:

| Module | TCP Ports | UDP Ports | Notes                                                                                   |
| ------ | --------- | --------- |-----------------------------------------------------------------------------------------|
| SRS | 1935, 1985, 8080,<br/> 8088, 1990, 554,<br/> 8936 | 8000, 8935, 10080,<br/> 1989 | See [SRS ports](https://github.com/ossrs/srs/blob/develop/trunk/doc/Resources.md#ports) |
| platform | 2024 |  - | Mount at `/terraform/v1/mgmt/`, `/terraform/v1/hooks/`, `/terraform/v1/ffmpeg/` and `/terraform/v1/tencent/`    |
| releases | 2023 |  - | Mount at `/terraform/v1/releases`                                                       |
| mgmt | 2022 |  - | Mount at `/mgmt/` and `/terraform/v1/mgmt/`                                             |
| node-exporter | 9100 | - | -                                                                                       |
| redis | 56379 | - | -                                                                                       |

> Note: Hooks(2021) has been migrated to platform(2024).

> Note: FFmpeg(2019) has been migrated to platform(2024).

> Note: TencentCloud(2020) has been migrated to platform(2024).

> Note: Mgmt(2022) has been migrated to platform(2024).

## APIs

Platform:

* `/terraform/v1/mgmt/versions` Public version api.
* `/terraform/v1/mgmt/init` Whether mgmt initialized.
* `/terraform/v1/mgmt/check` Check whether system is ok.
* `/terraform/v1/mgmt/token` System auth with token.
* `/terraform/v1/mgmt/login` System auth with password.
* `/terraform/v1/mgmt/status` Query the version of mgmt.
* `/terraform/v1/mgmt/envs` Query the envs of mgmt.
* `/terraform/v1/mgmt/bilibili` Query the video information.
* `/terraform/v1/mgmt/beian/query` Query the beian information.
* `/terraform/v1/mgmt/beian/update` Update the beian information.
* `/terraform/v1/mgmt/secret/query` Query the api secret for OpenAPI.
* `/terraform/v1/mgmt/nginx/hls` Update NGINX config, to enable HLS delivery.
* `/terraform/v1/mgmt/ssl` Config the system SSL config.
* `/terraform/v1/mgmt/letsencrypt` Config the let's encrypt SSL.
* `/terraform/v1/host/versions` Public version api.
* `/terraform/v1/releases` Version management for all components.

Also by platform module:

* `/terraform/v1/hooks/srs/verify` Hooks: Verify the stream request URL of SRS.
* `/terraform/v1/hooks/srs/secret/query` Hooks: Query the secret to generate stream URL.
* `/terraform/v1/hooks/srs/secret/update` Hooks: Update the secret to generate stream URL.
* `/terraform/v1/hooks/srs/secret/disable` Hooks: Disable the secret for authentication.
* `/terraform/v1/hooks/srs/hls` Hooks: Handle the `on_hls` event.
* `/terraform/v1/hooks/record/query` Hooks: Query the Record pattern.
* `/terraform/v1/hooks/record/apply` Hooks: Apply the Record pattern.
* `/terraform/v1/hooks/record/remove` Hooks: Remove the Record files.
* `/terraform/v1/hooks/record/files` Hooks: List the Record files.
* `/terraform/v1/hooks/record/hls/:uuid.m3u8` Hooks: Generate HLS/m3u8 url to preview or download.
* `/terraform/v1/hooks/record/hls/:uuid/index.m3u8` Hooks: Serve HLS m3u8 files.
* `/terraform/v1/hooks/record/hls/:dir/:m3u8/:uuid.ts` Hooks: Serve HLS ts files.
* `/terraform/v1/hooks/dvr/apply` Hooks: Apply the DVR pattern.
* `/terraform/v1/hooks/dvr/query` Hooks: Query the DVR pattern.
* `/terraform/v1/hooks/dvr/files` Hooks: List the DVR files.
* `/terraform/v1/hooks/dvr/hls/:uuid.m3u8` Hooks: Generate HLS/m3u8 url to preview or download.
* `/terraform/v1/hooks/vod/query` Hooks: Query the VoD pattern.
* `/terraform/v1/hooks/vod/apply` Hooks: Apply the VoD pattern.
* `/terraform/v1/hooks/vod/files` Hooks: List the VoD files.
* `/terraform/v1/hooks/vod/hls/:uuid.m3u8` Hooks: Generate HLS/m3u8 url to preview or download.
* `/terraform/v1/tencent/cam/secret` Tencent: Setup the CAM SecretId and SecretKey.
* `/terraform/v1/ffmpeg/forward/secret` FFmpeg: Setup the forward secret to live streaming platforms.
* `/terraform/v1/ffmpeg/forward/streams` FFmpeg: Query the forwarding streams.

Also provided by platform for market:

* `/api/` SRS: HTTP API of SRS media server.
* `/rtc/` SRS: HTTP API for WebERTC of SRS media server.
* `/*/*.(flv|m3u8|ts|aac|mp3)` SRS: Media stream for HTTP-FLV, HLS, HTTP-TS, HTTP-AAC, HTTP-MP3.

Also provided by platform for static Files:

* `/tools/` A set of H5 tools, like simple player, xgplayer, etc, serve by mgmt.
* `/console/` The SRS console, serve by mgmt.
* `/players/` The SRS player, serve by mgmt.
* `/mgmt/` The ui for mgmt, serve by mgmt.

**Removed** API:

* `/terraform/v1/mgmt/strategy` Toggle the upgrade strategy.
* `/prometheus` Prometheus: Time-series database and monitor.
* `/terraform/v1/mgmt/nginx/proxy` Setup a reverse proxy location.
* `/terraform/v1/mgmt/dns/lb` HTTP-DNS for hls load balance.
* `/terraform/v1/mgmt/dns/backend/update` HTTP-DNS: Update the backend servers for hls load balance.
* `/terraform/v1/mgmt/nginx/homepage` Setup the homepage redirection.
* `/terraform/v1/mgmt/window/query` Query the upgrade time window.
* `/terraform/v1/mgmt/window/update` Update the upgrade time window.
* `/.well-known/acme-challenge/` HTTPS verify mount for letsencrypt.
* `/terraform/v1/mgmt/pubkey` Update the access for platform administrator pubkey.
* `/terraform/v1/mgmt/upgrade` Upgrade the mgmt to latest version.
* `/terraform/v1/mgmt/containers` Query SRS container.
* `/terraform/v1/host/exec` Exec command sync, response the stdout and stderr.
* `/terraform/v1/mgmt/secret/token` Create token for OpenAPI.

## Depends

The software we depend on:

* Docker, `apt-get install -y docker.io`
    * Redis, `apt-get install -y redis`
    * Nginx, `apt-get install -y nginx`
        * SSL: `mgmt/containers/ssl`
* [LEGO](https://github.com/go-acme/lego)
    * Verify webroot: `mgmt/containers/www/.well-known/acme-challenge/`
    * Cert files: `mgmt/containers/etc/letsencrypt/live/`
* [SRS](https://github.com/ossrs/srs)
    * Config: `mgmt/containers/conf/srs.conf` mount as `/usr/local/srs/conf/lighthouse.conf`
    * Volume: `mgmt/containers/objs/nginx/html` mount as `/usr/local/srs/objs/nginx/html`
* [srs-hooks](https://github.com/ossrs/srs-stack/tree/lighthouse/hooks)
    * Volume: `mgmt/containers/objs/nginx/html` mount as `/usr/local/mgmt/containers/objs/nginx/html`
* [tencent-cloud](https://github.com/ossrs/srs-stack/tree/lighthouse/tencent)
    * [CAM](https://console.cloud.tencent.com/cam/overview) Authentication by secretId and secretKey.
* [ffmpeg](https://github.com/ossrs/srs-stack/tree/lighthouse/ffmpeg)
    * [FFmpeg and ffprobe](https://ffmpeg.org) tools in `ossrs/srs:node-av`

## Environments

The optional environments defined by `mgmt/.env`:

* `CLOUD`: `dev|bt|aapanel|droplet|docker`, The cloud platform name, DEV for development.
* `REGION`: `ap-guangzhou|ap-singapore|sgp1`, The region for upgrade source.
* `SOURCE`: `github|gitee`, The source code for upgrading.
* `REGISTRY`: `docker.io|registry.cn-hangzhou.aliyuncs.com`, The docker registry.
* `MGMT_LISTEN`: The listen port for mgmt HTTP server. Default: 2022
* `PLATFORM_LISTEN`: The listen port for platform HTTP server. Default: 2024
* `SRS_DOCKERIZED`: `true|false` Indicates the OS is in docker.

For mgmt to start platform in docker, because it can't access redis which is started by platform:

* `PLATFORM_DOCKER`: Whether run platform in docker. Default: true
* `MGMT_DOCKER`: Whether run mgmt in docker. Default: false

For testing the specified service:

* `NODE_ENV`: `development|production`, if development, use local redis; otherwise, use `mgmt.srs.local` in docker.
* `LOCAL_RELEASE`: `true|false`, whether use local release service.

For github actions to control the containers:

* `SRS_DOCKER`: `srs` to enfore use `ossrs/srs` docker image.
* `USE_DOCKER`: `true|false`, if false, disable all docker containers.
* `SRS_UTEST`: `true|false`, if true, running in utest mode.

For mgmt and containers to connect to redis:

* `REDIS_PASSWORD`: The redis password.
* `REDIS_PORT`: The redis port.

Environments for react ui:

* `PUBLIC_URL`: The mount prefix.
* `BUILD_PATH`: The output build path, default to `build`.

> Note: The env for react must start with `REACT_APP_`, please read [this post](https://create-react-app.dev/docs/adding-custom-environment-variables/#referencing-environment-variables-in-the-html).

Removed variables in .env:

* `SRS_PLATFORM_SECRET`: The mgmt api secret for token generating and verifying.

Please restart service when `.env` changed.

## Architecture

The architecture of [srs-stack](https://github.com/ossrs/srs-stack#architecture) by
[mermaid](https://mermaid.live/edit#pako:eNqNkctuwjAQRX_F8qIKEiH7tEKqCLRSXyhp2ZAuTDx5iNiOnDEFIf69jtMWwqoL2zPjozt37CPNFAca0rxWX1nJNJLn-DaVhMiiknvi-1MiCoFet91tdDDdHMiDGjmkqzkiiRN3Piq1bb2mZpgrLYKyS0c9gRqYuDdYkhsSrWK7r1TkVLqsQ2ZvyVB1sRANFGe5PO_yQWufLH9u_zByYbD35f9HaUo-mkIzDi6OoQbWwhB4B5mBxFmtDD9rXVb7WWf3L2u7AjtQYIf8HKostRKAJZjWWXu1zz_fN0oj6CEYqWx7XZvL3dqbgNy5r8g0cNu6YnU7wT2OrjrFwKt27bmDPK3sNR1TAVqwitsfP3ZwSq0VASkNbcghZ6bGlKbyZFHTcIYw5xUqTUPUBsaUGVTJQWa_ec9EFbNvJ2iYWytw-gba_8FA)

```mermaid
flowchart LR;
  nginx --> mgmt(mgmt<br/>by Go);
  mgmt --> SRS --> Hooks(platform/hooks) --> StreamAuth & DVR & VoD;
  DVR --> COS;
  mgmt --> FFmpeg(platform/ffmpeg);
  mgmt --- Platform(platform by Go);
  SRS --- FFmpeg(platform/ffmpeg);
  mgmt --> Upgrade --> Release;
  mgmt --> TencentCloud(platform/TencentCloud) --> CAM[CAM/COS/VoD];
  mgmt --> Prometheus --- NodeExporter;
  mgmt --> Docker;
  mgmt --> Env[(.env<br/>credentials.txt)];
  mgmt --> Redis[(Redis KV)];
```

> Note: It's a single node, also light-weighted, video cloud for tiny company, personal user and starter.

```mermaid
flowchart LR;
  nginx --> mgmt(mgmt<br/>by Go);
  aaPanel --> mgmt;
  aaPanel --> nginx;
```

> Note: This is an optional workflow for user to use aaPanel to deploy srs-stack.
