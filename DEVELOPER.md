# Developer

This guide is for developers and covers topics such as OpenAPI, environment variables, 
resources, and ports, as well as development on Mac or using Docker.

## Develop All in macOS

Start redis and SRS by docker, set the candidate explicitly:

```bash
docker rm -f redis srs 2>/dev/null &&
docker run --name redis --rm -it -v $HOME/data/redis:/data -p 6379:6379 -d redis &&
touch platform/containers/data/config/srs.server.conf platform/containers/data/config/srs.vhost.conf &&
docker run --name srs --rm -it \
    -v $(pwd)/platform/containers/data/config:/usr/local/srs/containers/data/config \
    -v $(pwd)/platform/containers/conf/srs.release-mac.conf:/usr/local/srs/conf/docker.conf \
    -v $(pwd)/platform/containers/objs/nginx:/usr/local/srs/objs/nginx \
    -p 1935:1935 -p 1985:1985 -p 8080:8080 -p 8000:8000/udp -p 10080:10080/udp \
    -d ossrs/srs:5
```

> Note: Use the intranet IP for WebRTC to set the candidate.

> Note: Stop service by `docker rm -f redis srs`

> Note: Also, you can run SRS by `(cd platform && ~/git/srs/trunk/objs/srs -c containers/conf/srs.release-local.conf)`

> Note: You can set the candidate for WebRTC by `--env CANDIDATE=$(ifconfig en0 |grep 'inet ' |awk '{print $2}')`

Run the platform backend, or run in GoLand:

```bash
(cd platform && go run .)
```

> Note: Set `AUTO_SELF_SIGNED_CERTIFICATE=on` if need to generate self-signed certificate.

Run all tests:

```bash
bash scripts/tools/secret.sh --output test/.env &&
(cd test && go test -v --endpoint=http://localhost:2022 -init-self-signed-cert=true) &&
(cd test && go test -v --endpoint=https://localhost:2443 -init-self-signed-cert=false)
```

Run the platform react ui, or run in WebStorm:

```bash
(cd ui && npm install && npm start)
```

Access the browser: http://localhost:3000

## Develop the Script Installer

> Note: Please note that BT plugin will use the current branch version, including develop version.

Build a docker image:

```bash
docker rm -f script 2>/dev/null &&
docker rmi srs-script-dev 2>/dev/null || echo OK &&
docker build -t srs-script-dev -f scripts/setup-ubuntu/Dockerfile.script .
```

Create a docker container in daemon:

```bash
docker rm -f script 2>/dev/null &&
docker run -p 2022:2022 -p 2443:2443 -p 1935:1935 \
    -p 8080:8080 -p 8000:8000/udp -p 10080:10080/udp \
    --env CANDIDATE=$(ifconfig en0 |grep 'inet ' |awk '{print $2}') \
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
docker exec -it script rm -f /data/config/.env &&
docker exec -it script bash build/srs_stack/scripts/setup-ubuntu/uninstall.sh 2>/dev/null || echo OK &&
bash scripts/setup-ubuntu/build.sh --output $(pwd)/build --extract &&
docker exec -it script bash build/srs_stack/scripts/setup-ubuntu/install.sh --verbose
```

Run test for script:

```bash
rm -f test/srs-stack.test &&
docker exec -it script make -j -C test &&
bash scripts/tools/secret.sh --output test/.env &&
docker exec -it script ./test/srs-stack.test -test.v -endpoint http://localhost:2022 \
    -srs-log=true -wait-ready=true -init-password=true -check-api-secret=true -init-self-signed-cert=true \
    -test.run TestSystem_Empty &&
bash scripts/tools/secret.sh --output test/.env &&
docker exec -it script ./test/srs-stack.test -test.v -wait-ready -endpoint http://localhost:2022 \
    -srs-log=true -wait-ready=true -init-password=false -check-api-secret=true \
    -test.parallel 3 &&
docker exec -it script ./test/srs-stack.test -test.v -wait-ready -endpoint https://localhost:2443 \
    -srs-log=true -wait-ready=true -init-password=false -check-api-secret=true \
    -test.parallel 3
```

Access the browser: [http://localhost:2022](http://localhost:2022)

## Develop the aaPanel Plugin

> Note: Please note that BT plugin will use the current branch version, including develop version.

Start a container and mount as plugin:

```bash
docker rm -f bt aapanel 2>/dev/null &&
AAPANEL_KEY=$(cat $HOME/.bt/api.json |awk -F token_crypt '{print $2}' |cut -d'"' -f3) &&
docker run -p 80:80 -p 443:443 -p 7800:7800 \
    -p 1935:1935 -p 8080:8080 -p 8000:8000/udp -p 10080:10080/udp \
    --env CANDIDATE=$(ifconfig en0 |grep 'inet ' |awk '{print $2}') \
    -v $(pwd)/build/srs_stack:/www/server/panel/plugin/srs_stack \
    -v $HOME/.bt/api.json:/www/server/panel/config/api.json -e BT_KEY=$AAPANEL_KEY \
    --privileged -v /sys/fs/cgroup:/sys/fs/cgroup:rw --cgroupns=host \
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
major=$(echo $version |awk -F '.' '{print $1}' |sed 's/v//g') &&
docker exec -it aapanel docker load -i platform.tar && 
docker exec -it aapanel docker tag platform:latest ossrs/srs-stack:$version &&
docker exec -it aapanel docker tag platform:latest ossrs/srs-stack:$major &&
docker exec -it aapanel docker tag platform:latest registry.cn-hangzhou.aliyuncs.com/ossrs/srs-stack:$version &&
docker exec -it aapanel docker tag platform:latest registry.cn-hangzhou.aliyuncs.com/ossrs/srs-stack:$major &&
docker exec -it aapanel docker images
```

Next, build the aaPanel plugin and install it:

```bash
docker exec -it aapanel rm -f /data/config/.env &&
docker exec -it aapanel bash /www/server/panel/plugin/srs_stack/install.sh uninstall 2>/dev/null || echo OK &&
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

Setup the dns lookup for domain `srs.stack.local`:

```bash
PIP=$(docker exec -it aapanel ifconfig eth0 |grep 'inet ' |awk '{print $2}') &&
docker exec -it aapanel bash -c "echo '$PIP srs.stack.local' >> /etc/hosts" &&
docker exec -it aapanel cat /etc/hosts && echo OK &&
docker exec -it aapanel docker exec -it srs-stack bash -c "echo '$PIP srs.stack.local' >> /etc/hosts" &&
docker exec -it aapanel docker exec -it srs-stack cat /etc/hosts
```
> Note: We add host `srs.stack.local` to the ip of eth0, because we need to access it in the
> srs-stack docker in docker.

Run test for aaPanel:

```bash
rm -f test/srs-stack.test &&
docker exec -it aapanel make -j -C test &&
bash scripts/tools/secret.sh --output test/.env &&
docker exec -it aapanel ./test/srs-stack.test -test.v -endpoint http://srs.stack.local:80 \
    -srs-log=true -wait-ready=true -init-password=true -check-api-secret=true -init-self-signed-cert=true \
    -test.run TestSystem_Empty &&
bash scripts/tools/secret.sh --output test/.env &&
docker exec -it aapanel ./test/srs-stack.test -test.v -wait-ready -endpoint http://srs.stack.local:80 \
    -srs-log=true -wait-ready=true -init-password=false -check-api-secret=true \
    -test.parallel 3 &&
docker exec -it aapanel ./test/srs-stack.test -test.v -wait-ready -endpoint https://srs.stack.local:443 \
    -srs-log=true -wait-ready=true -init-password=false -check-api-secret=true \
    -test.parallel 3
```

Open [http://localhost:7800/srsstack](http://localhost:7800/srsstack) to install plugin.

> Note: Or you can use `docker exec -it aapanel bt default` to show the login info.

In the [application store](http://localhost:7800/soft), there is a `srs_stack` plugin. After test, you can install the plugin
`build/aapanel-srs_stack.zip` to production aaPanel panel.

## Develop the BT Plugin

> Note: Please note that BT plugin will use the current branch version, including develop version.

Start a container and mount as plugin:

```bash
docker rm -f bt aapanel 2>/dev/null &&
BT_KEY=$(cat $HOME/.bt/api.json |awk -F token_crypt '{print $2}' |cut -d'"' -f3) &&
docker run -p 80:80 -p 443:443 -p 7800:7800 \
    -p 1935:1935 -p 8080:8080 -p 8000:8000/udp -p 10080:10080/udp \
    --env CANDIDATE=$(ifconfig en0 |grep 'inet ' |awk '{print $2}') \
    -v $(pwd)/build/srs_stack:/www/server/panel/plugin/srs_stack \
    -v $HOME/.bt/userInfo.json:/www/server/panel/data/userInfo.json \
    -v $HOME/.bt/api.json:/www/server/panel/config/api.json -e BT_KEY=$BT_KEY \
    --privileged -v /sys/fs/cgroup:/sys/fs/cgroup:rw --cgroupns=host \
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
major=$(echo $version |awk -F '.' '{print $1}' |sed 's/v//g') &&
docker exec -it bt docker load -i platform.tar && 
docker exec -it bt docker tag platform:latest ossrs/srs-stack:$version &&
docker exec -it bt docker tag platform:latest ossrs/srs-stack:$major &&
docker exec -it bt docker tag platform:latest registry.cn-hangzhou.aliyuncs.com/ossrs/srs-stack:$version &&
docker exec -it bt docker tag platform:latest registry.cn-hangzhou.aliyuncs.com/ossrs/srs-stack:$major &&
docker exec -it bt docker images
```

Next, build the BT plugin and install it:

```bash
docker exec -it bt bash /www/server/panel/plugin/srs_stack/install.sh uninstall 2>/dev/null || echo OK &&
docker exec -it bt rm -f /data/config/.env &&
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

Setup the dns lookup for domain `srs.stack.local`:

```bash
PIP=$(docker exec -it bt ifconfig eth0 |grep 'inet ' |awk '{print $2}') &&
docker exec -it bt bash -c "echo '$PIP srs.stack.local' >> /etc/hosts" &&
docker exec -it bt cat /etc/hosts && echo OK &&
docker exec -it bt docker exec -it srs-stack bash -c "echo '$PIP srs.stack.local' >> /etc/hosts" &&
docker exec -it bt docker exec -it srs-stack cat /etc/hosts
```
> Note: We add host `srs.stack.local` to the ip of eth0, because we need to access it in the
> srs-stack docker in docker.

Run test for BT:

```bash
rm -f test/srs-stack.test &&
docker exec -it bt make -j -C test &&
bash scripts/tools/secret.sh --output test/.env &&
docker exec -it bt ./test/srs-stack.test -test.v -endpoint http://srs.stack.local:80 \
    -srs-log=true -wait-ready=true -init-password=true -check-api-secret=true -init-self-signed-cert=true \
    -test.run TestSystem_Empty &&
bash scripts/tools/secret.sh --output test/.env &&
docker exec -it bt ./test/srs-stack.test -test.v -wait-ready -endpoint http://srs.stack.local:80 \
    -srs-log=true -wait-ready=true -init-password=false -check-api-secret=true \
    -test.parallel 3 &&
docker exec -it bt ./test/srs-stack.test -test.v -wait-ready -endpoint https://srs.stack.local:443 \
    -srs-log=true -wait-ready=true -init-password=false -check-api-secret=true \
    -test.parallel 3
```

Open [http://localhost:7800/srsstack](http://localhost:7800/srsstack) to install plugin.

> Note: Or you can use `docker exec -it bt bt default` to show the login info.

In the [application store](http://localhost:7800/soft), there is a `srs_stack` plugin. After test, you can install the plugin 
`build/bt-srs_stack.zip` to production BT panel.

## Develop the Droplet Image

> Note: Please note that BT plugin will use the current branch version, including develop version.

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

Please check the [snapshot](https://cloud.digitalocean.com/images/snapshots/droplets), and create a test droplet.

```bash
IMAGE=$(doctl compute snapshot list --context market --format ID --no-header) &&
sshkey=$(doctl compute ssh-key list --context market --no-header |grep srs |awk '{print $1}') &&
doctl compute droplet create srs-stack-test --context market --image $IMAGE \
    --region sgp1 --size s-2vcpu-2gb --ssh-keys $sshkey --wait &&
SRS_DROPLET_EIP=$(doctl compute droplet get srs-stack-test --context market --format PublicIPv4 --no-header)
```

Prepare test environment:

```bash
ssh root@$SRS_DROPLET_EIP sudo mkdir -p /data/upload test scripts/tools &&
ssh root@$SRS_DROPLET_EIP sudo chmod 777 /data/upload &&
cp ~/git/srs/trunk/doc/source.200kbps.768x320.flv test/ &&
scp ./test/source.200kbps.768x320.flv root@$SRS_DROPLET_EIP:/data/upload/ &&
docker run --rm -it -v $(pwd):/g -w /g ossrs/srs:ubuntu20 make -C test clean default &&
scp ./test/srs-stack.test ./test/source.200kbps.768x320.flv root@$SRS_DROPLET_EIP:~/test/ &&
scp ./scripts/tools/secret.sh root@$SRS_DROPLET_EIP:~/scripts/tools &&
ssh root@$SRS_DROPLET_EIP docker run --rm -v /usr/bin:/g ossrs/srs:tools \
    cp /usr/local/bin/ffmpeg /usr/local/bin/ffprobe /g/
```

Test the droplet instance:

```bash
ssh root@$SRS_DROPLET_EIP bash scripts/tools/secret.sh --output test/.env &&
ssh root@$SRS_DROPLET_EIP ./test/srs-stack.test -test.v -endpoint http://$SRS_DROPLET_EIP:2022 \
    -srs-log=true -wait-ready=true -init-password=true -check-api-secret=true -init-self-signed-cert=true \
    -test.run TestSystem_Empty &&
ssh root@$SRS_DROPLET_EIP bash scripts/tools/secret.sh >test/.env &&
ssh root@$SRS_DROPLET_EIP ./test/srs-stack.test -test.v -wait-ready -endpoint http://$SRS_DROPLET_EIP:2022 \
    -endpoint-rtmp rtmp://$SRS_DROPLET_EIP -endpoint-http http://$SRS_DROPLET_EIP -endpoint-srt srt://$SRS_DROPLET_EIP:10080 \
    -srs-log=true -wait-ready=true -init-password=false -check-api-secret=true \
    -test.parallel 1 &&
ssh root@$SRS_DROPLET_EIP ./test/srs-stack.test -test.v -wait-ready -endpoint https://$SRS_DROPLET_EIP:2443 \
    -endpoint-rtmp rtmp://$SRS_DROPLET_EIP -endpoint-http https://$SRS_DROPLET_EIP -endpoint-srt srt://$SRS_DROPLET_EIP:10080 \
    -srs-log=true -wait-ready=true -init-password=false -check-api-secret=true \
    -test.parallel 1
```

Remove the droplet instance:

```bash
doctl compute droplet delete srs-stack-test --context market --force
```

After submit to [marketplace](https://cloud.digitalocean.com/vendorportal/624145d53da4ad68de259945/10/edit), cleanup the snapshot:

```bash
doctl compute snapshot delete $IMAGE --context market --force
```

> Note: The snapshot should be removed if submit to marketplace, so you don't need to delete it.

## Develop the Lighthouse Image

> Note: Please note that BT plugin will use the current branch version, including develop version.

To build SRS image for [TencentCloud Lighthouse](https://cloud.tencent.com/product/lighthouse).

For the first run, please create a [TencentCloud Secret](https://console.cloud.tencent.com/cam/capi) and save
to `~/.lighthouse/.env` file:

```bash
LH_ACCOUNT=xxxxxx
LH_PROD=xxxxxx
SECRET_ID=xxxxxx
SECRET_KEY=xxxxxx
```

> Note: Share the image to `LH_ACCOUNT` to publish it.

Create a CVM instance:

```bash
rm -f .tmp/lh-*.txt &&
echo "$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 16)A0" >.tmp/lh-token.txt &&
VM_TOKEN=$(cat .tmp/lh-token.txt) bash scripts/tools/tencent-cloud/helper.sh create-cvm.py --id $(pwd)/.tmp/lh-instance.txt &&
bash scripts/tools/tencent-cloud/helper.sh query-cvm-ip.py --instance $(cat .tmp/lh-instance.txt) --id $(pwd)/.tmp/lh-ip.txt &&
echo "Instance: $(cat .tmp/lh-instance.txt), IP: ubuntu@$(cat .tmp/lh-ip.txt), Password: $(cat .tmp/lh-token.txt)" && sleep 5 &&
bash scripts/setup-lighthouse/build.sh --ip $(cat .tmp/lh-ip.txt) --os ubuntu --user ubuntu --password $(cat .tmp/lh-token.txt) &&
bash scripts/tools/tencent-cloud/helper.sh create-image.py --instance $(cat .tmp/lh-instance.txt) --id $(pwd)/.tmp/lh-image.txt &&
bash scripts/tools/tencent-cloud/helper.sh share-image.py --image $(cat .tmp/lh-image.txt) &&
echo "Image: $(cat .tmp/lh-image.txt) created and shared." &&
bash scripts/tools/tencent-cloud/helper.sh remove-cvm.py --instance $(cat .tmp/lh-instance.txt)
```

Next, create a test CVM instance with the image:

```bash
echo "$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 16)A0" >.tmp/lh-token2.txt &&
VM_TOKEN=$(cat .tmp/lh-token2.txt) bash scripts/tools/tencent-cloud/helper.sh create-verify.py --image $(cat .tmp/lh-image.txt) --id $(pwd)/.tmp/lh-test.txt &&
bash scripts/tools/tencent-cloud/helper.sh query-cvm-ip.py --instance $(cat .tmp/lh-test.txt) --id $(pwd)/.tmp/lh-ip2.txt && 
echo "IP: ubuntu@$(cat .tmp/lh-ip2.txt), Password: $(cat .tmp/lh-token2.txt)" &&
echo "http://$(cat .tmp/lh-ip2.txt)"
```

Prepare test environment:

```bash
sshCmd="sshpass -p $(cat .tmp/lh-token2.txt) ssh -o StrictHostKeyChecking=no -t" &&
scpCmd="sshpass -p $(cat .tmp/lh-token2.txt) scp -o StrictHostKeyChecking=no" &&
$sshCmd ubuntu@$(cat .tmp/lh-ip2.txt) sudo mkdir -p /data/upload &&
$sshCmd ubuntu@$(cat .tmp/lh-ip2.txt) mkdir -p test scripts/tools &&
$sshCmd ubuntu@$(cat .tmp/lh-ip2.txt) sudo chmod 777 /data/upload &&
cp ~/git/srs/trunk/doc/source.200kbps.768x320.flv test/ &&
$scpCmd test/source.200kbps.768x320.flv ubuntu@$(cat .tmp/lh-ip2.txt):/data/upload/ &&
docker run --rm -it -v $(pwd):/g -w /g ossrs/srs:ubuntu20 make -C test clean default &&
$scpCmd ./test/srs-stack.test ./test/source.200kbps.768x320.flv ubuntu@$(cat .tmp/lh-ip2.txt):~/test/ &&
$scpCmd ./scripts/tools/secret.sh ubuntu@$(cat .tmp/lh-ip2.txt):~/scripts/tools &&
$sshCmd ubuntu@$(cat .tmp/lh-ip2.txt) sudo docker run --rm -v /usr/bin:/g \
    registry.cn-hangzhou.aliyuncs.com/ossrs/srs:tools \
    cp /usr/local/bin/ffmpeg /usr/local/bin/ffprobe /g/
```

Test the CVM instance:

```bash
$sshCmd ubuntu@$(cat .tmp/lh-ip2.txt) sudo bash scripts/tools/secret.sh --output test/.env &&
$sshCmd ubuntu@$(cat .tmp/lh-ip2.txt) ./test/srs-stack.test -test.v -endpoint http://$(cat .tmp/lh-ip2.txt):2022 \
    -srs-log=true -wait-ready=true -init-password=true -check-api-secret=true -init-self-signed-cert=true \
    -test.run TestSystem_Empty &&
$sshCmd ubuntu@$(cat .tmp/lh-ip2.txt) sudo bash scripts/tools/secret.sh --output test/.env &&
$sshCmd ubuntu@$(cat .tmp/lh-ip2.txt) ./test/srs-stack.test -test.v -wait-ready -endpoint http://$(cat .tmp/lh-ip2.txt):2022 \
    -endpoint-rtmp rtmp://$(cat .tmp/lh-ip2.txt) -endpoint-http http://$(cat .tmp/lh-ip2.txt) -endpoint-srt srt://$(cat .tmp/lh-ip2.txt):10080 \
    -srs-log=true -wait-ready=true -init-password=false -check-api-secret=true \
    -test.parallel 3 &&
ssh ubuntu@$(cat .tmp/lh-ip2.txt) ./test/srs-stack.test -test.v -wait-ready -endpoint https://$(cat .tmp/lh-ip2.txt):2443 \
    -endpoint-rtmp rtmp://$(cat .tmp/lh-ip2.txt) -endpoint-http https://$(cat .tmp/lh-ip2.txt) -endpoint-srt srt://$(cat .tmp/lh-ip2.txt):10080 \
    -srs-log=true -wait-ready=true -init-password=false -check-api-secret=true \
    -test.parallel 3
```

Verify then cleanup the test CVM instance:

```bash
bash scripts/tools/tencent-cloud/helper.sh remove-cvm.py --instance $(cat .tmp/lh-test.txt)
```

After publish to lighthouse, cleanup the CVM, disk images, and snapshots:

```bash
bash scripts/tools/tencent-cloud/helper.sh remove-image.py --image $(cat .tmp/lh-image.txt)
```

If need to test the domain of lighthouse, create a domain `lighthouse.ossrs.net`:

```bash
# Create the test domain for lighthouse
doctl compute domain records create ossrs.net \
    --record-type A --record-name lighthouse --record-data $(cat .tmp/lh-ip2.txt) \
    --record-ttl 300 &&
echo "https://lighthouse.ossrs.net"

# Remove the test domain for lighthouse
doctl compute domain records delete ossrs.net -f \
    $(doctl compute domain records list ossrs.net --no-header |grep lighthouse |awk '{print $1}') &&
echo "Record lighthouse.ossrs.net removed"
```

## Develop the SSL Cert for HTTPS 

Create a domain for HTTPS:

```bash
LNAME=lego && LDOMAIN=ossrs.net &&
doctl compute domain records create $LDOMAIN \
    --record-type A --record-name $LNAME --record-data $(dig +short $LDOMAIN) \
    --record-ttl 3600
```

Build and save the script image to file:

```bash
docker rmi platform:latest 2>/dev/null || echo OK &&
docker build -t platform:latest -f Dockerfile . &&
docker save -o platform.tar platform:latest
```

Copy and load the image to server:

```bash
ssh root@$LNAME.$LDOMAIN rm -f platform.tar* 2>/dev/null &&
rm -f platform.tar.gz && tar zcf platform.tar.gz platform.tar &&
scp platform.tar.gz root@$LNAME.$LDOMAIN:~ &&
ssh root@$LNAME.$LDOMAIN tar xf platform.tar.gz &&
version=$(bash scripts/version.sh) &&
ssh root@$LNAME.$LDOMAIN docker load -i platform.tar &&
ssh root@$LNAME.$LDOMAIN docker tag platform:latest ossrs/srs-stack:$version &&
ssh root@$LNAME.$LDOMAIN docker tag platform:latest registry.cn-hangzhou.aliyuncs.com/ossrs/srs-stack:$version &&
ssh root@$LNAME.$LDOMAIN docker image prune -f &&
ssh root@$LNAME.$LDOMAIN docker images
```

Next, build the BT plugin and install it:

```bash
ssh root@$LNAME.$LDOMAIN bash /www/server/panel/plugin/srs_stack/install.sh uninstall 2>/dev/null || echo OK &&
bash scripts/setup-bt/auto/zip.sh --output $(pwd)/build --extract &&
scp build/bt-srs_stack.zip root@$LNAME.$LDOMAIN:~ &&
ssh root@$LNAME.$LDOMAIN unzip -q bt-srs_stack.zip -d /www/server/panel/plugin &&
ssh root@$LNAME.$LDOMAIN bash /www/server/panel/plugin/srs_stack/install.sh install
```

On the server, setup the `.bashrc`:

```bash
export BT_KEY=xxxxxx
export PYTHONIOENCODING=UTF-8
```

You can use BT panel to install the plugin, or by command:

```bash
ssh root@$LNAME.$LDOMAIN python3 /www/server/panel/plugin/srs_stack/bt_api_remove_site.py &&
ssh root@$LNAME.$LDOMAIN DOMAIN=$LNAME.$LDOMAIN python3 /www/server/panel/plugin/srs_stack/bt_api_create_site.py &&
ssh root@$LNAME.$LDOMAIN python3 /www/server/panel/plugin/srs_stack/bt_api_setup_site.py &&
ssh root@$LNAME.$LDOMAIN bash /www/server/panel/plugin/srs_stack/setup.sh \
    --r0 /tmp/srs_stack_install.r0 --nginx /www/server/nginx/logs/nginx.pid \
    --www /www/wwwroot --site srs.stack.local
```

Cleanup, remove the files and domain:

```bash
ssh root@$LNAME.$LDOMAIN rm -f platform.tar* bt-srs_stack.zip 2>/dev/null &&
ssh root@$LNAME.$LDOMAIN python3 /www/server/panel/plugin/srs_stack/bt_api_remove_site.py &&
ssh root@$LNAME.$LDOMAIN bash /www/server/panel/plugin/srs_stack/install.sh uninstall 2>/dev/null || echo OK &&
domains=$(doctl compute domain records ls $LDOMAIN --no-header |grep $LNAME) && echo "Cleanup domains: $domains" &&
doctl compute domain records delete $LDOMAIN $(echo $domains |awk '{print $1}') -f
```

Query domain and droplet:

```bash
doctl compute domain records ls ossrs.io |grep lego &&
doctl compute droplet ls |grep lego
```

## Develop the NGINX HLS CDN

Run SRS Stack by previous steps, such as [Develop All in macOS](#develop-all-in-macos), publish stream 
and there should be a HLS stream:

* [http://localhost:2022/live/livestream.m3u8](http://localhost:2022/tools/player.html?url=http://localhost:2022/live/livestream.m3u8)

Build the image of nginx:

```bash
docker rm -f nginx 2>/dev/null &&
docker rmi scripts/nginx-hls-cdn 2>/dev/null || echo OK &&
docker build -t ossrs/srs-stack:nginx-hls-cdn scripts/nginx-hls-cdn
```

> Note: The official image is build by [workflow](https://github.com/ossrs/srs-stack/actions/runs/5970907929) 
> which is triggered manually.

If you want to use NGINX as proxy, run by docker:

```bash
SRS_STACK_SERVER=$(ifconfig en0 |grep 'inet ' |awk '{print $2}') &&
docker run --rm -it -p 80:80 --name nginx -e SRS_STACK_SERVER=${SRS_STACK_SERVER}:2022 \
    ossrs/srs-stack:nginx-hls-cdn
```

There should be a new HLS stream, cached by NGINX:

* [http://localhost/live/livestream.m3u8](http://localhost:2022/tools/player.html?url=http://localhost/live/livestream.m3u8)

To test the CROS with `OPTIONS`, use [HTTP-REST](http://ossrs.net/http-rest/) tool, or by curl:

```bash
curl 'http://localhost/live/livestream.m3u8' -X 'OPTIONS' -H 'Origin: http://ossrs.net' -v
curl 'http://localhost/live/livestream.m3u8' -X 'GET' -H 'Origin: http://ossrs.net' -v
```

To start a [srs-bench](https://github.com/ossrs/srs-bench) to test the performance:

```bash
docker run --rm -d ossrs/srs:sb ./objs/sb_hls_load \
    -c 100 -r http://host.docker.internal/live/livestream.m3u8
```

The load should be taken by NGINX, not the SRS Stack.

## Product the NGINX HLS CDN

Install SRS Stack by BT or aaPanel or docker, assume the domain is `bt.ossrs.net`, publish
a RTMP stream to SRS Stack:

```bash
ffmpeg -re -i ~/git/srs/trunk/doc/source.flv -c copy \
    -f flv rtmp://bt.ossrs.net/live/livestream?secret=xxx
```

Open the [http://bt.ossrs.net/live/livestream.m3u8](http://bt.ossrs.net/tools/player.html?url=http://bt.ossrs.net/live/livestream.m3u8) 
to Check it.

Create a new domain `bt2.ossrs.net` for the same server:

```bash
doctl compute domain records create ossrs.net \
    --record-type A --record-name bt2 --record-data 39.100.79.15 \
    --record-ttl 3600
```

Create a new WebSite `bt2.ossrs.net` by BT or aaPanel, proxy to NGINX HLS Edge server:

```nginx
location /tools/ {
  proxy_pass http://localhost:2022;
}
location / {
  proxy_no_cache 1;
  proxy_cache_bypass 1;
  add_header X-Cache-Status-Proxy $upstream_cache_status;
  proxy_pass http://localhost:23080;
}
```

Start a NGINX HLS Edge server:

```bash
docker rm -f srs-stack-nginx01 || echo OK &&
PIP=$(ifconfig eth0 |grep 'inet ' |awk '{print $2}') &&
docker run --rm -it -e SRS_STACK_SERVER=$PIP:2022 \
    -p 23080:80 --name srs-stack-nginx01 -d \
    ossrs/srs-stack:nginx-hls-cdn
```

Open the [http://bt2.ossrs.net/live/livestream.m3u8](http://bt.ossrs.net/tools/player.html?url=http://bt2.ossrs.net/live/livestream.m3u8)
to Check it.

Use curl to test the HLS cache:

```bash
curl -v http://bt.ossrs.net/live/livestream.m3u8
curl -v http://bt.ossrs.net/live/livestream.m3u8 -H 'Origin: http://test.com'
curl -v http://bt2.ossrs.net/live/livestream.m3u8
curl -v http://bt2.ossrs.net/live/livestream.m3u8 -H 'Origin: http://test.com'
```

Be aware that the cache will store the CORS headers as well. This means that if you query 
and obtain HLS without CORS, it will remain without CORS even when a request includes an 
Origin header that necessitates CORS.

## Use HELM to Install SRS Stack

Install [HELM](https://helm.sh/docs/intro/install/) and [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/),
then add repo of SRS Stack:

```bash
helm repo add srs http://helm.ossrs.io/stable
```

Install the latest SRS Stack:

```bash
helm install srs srs/srs-stack
```

Or, install from file:

```bash
helm install srs ~/git/srs-helm/stable/srs-stack-1.0.5.tgz
```

Or, setup the persistence directory:

```bash
helm install srs ~/git/srs-helm/stable/srs-stack-1.0.5.tgz \
  --set persistence.path=$HOME/data
```

Finally, open [http://localhost](http://localhost) to check it.

## Run test in Goland

Prepare the .env:

```bash
bash scripts/tools/secret.sh --output test/.env &&
cp ~/git/srs/trunk/doc/source.200kbps.768x320.flv test/
```

Run testcase in Goland.

## Config SRS Container

The SRS container is configured by environment variables, which loads the `/data/config/.srs.env` 
file. To build a test image:

```bash
docker rmi srs-stack-env 2>/dev/null || echo OK &&
docker build -t srs-stack-env -f Dockerfile .
```

Setup the logging to file:

```bash
echo 'SRS_LOG_TANK=file' > $HOME/data/config/.srs.env
```

Run SRS Stack by docker:

```bash
docker run --rm -it -p 2022:2022 -p 2443:2443 -p 1935:1935 \
  -p 8080:8080 -p 8000:8000/udp -p 10080:10080/udp --name srs-stack \
    --env CANDIDATE=$(ifconfig en0 |grep 'inet ' |awk '{print $2}') \
  -v $HOME/data:/data srs-stack-env
```

Note that the logs should be written to file, there is no log `write log to console`, instead there
should be a log like `you can check log by`.

## WebRTC Candidate

SRS Stack follows the rules for WebRTC candidate, see [CANDIDATE](https://ossrs.io/lts/en-us/docs/v5/doc/webrtc#config-candidate),
but also has extra improvements for we can do more after proxy the API.

1. Disable `use_auto_detect_network_ip` and `api_as_candidates` in SRS config.
1. Always use `?eip=xxx` and ignore any other config, if user force to use the specified IP.
1. If `NAME_LOOKUP` (default is `on`) isn't `off`, try to resolve the candidate from `Host` of HTTP API by SRS Stack.
  1. If access SRS Stack by `localhost` for debugging or run in localhost.
    1. If `PLATFORM_DOCKER` is `off`, such as directly run in host, not in docker, use the private ip of SRS Stack.
    1. If not set `CANDIDATE`, use `127.0.0.1` for OBS WHIP or native client to access SRS Stack by localhost.
  1. Use `Host` if it's a valid IP address, for example, to access SRS Stack by public ip address.
  1. Use DNS lookup if `Host` is a domain, for example, to access SRS Stack by domain name.
1. If no candidate, use docker IP address discovered by SRS. 

> Note: Client can also set the header `X-Real-Host` to set the candidate.

> Note: Never use `host.docker.internal` because it's only available in docker, not in host server.

## Docker Allocated Ports

The ports allocated:

| Module | TCP Ports                                         | UDP Ports | Notes                                                                                                                                            |
| ------ |---------------------------------------------------| --------- |--------------------------------------------------------------------------------------------------------------------------------------------------|
| SRS | 1935, 1985, 8080,<br/> 8088, 1990, 554,<br/> 8936 | 8000, 8935, 10080,<br/> 1989 | See [SRS ports](https://github.com/ossrs/srs/blob/develop/trunk/doc/Resources.md#ports)                                                          |
| platform | 2022                                              |  - | Mount at `/mgmt/`, `/terraform/v1/mgmt/`, `/terraform/v1/hooks/`, `/terraform/v1/ffmpeg/` and `/terraform/v1/tencent/` |

> Note: FFmpeg(2019), TencentCloud(2020), Hooks(2021), Mgmt(2022), Platform(2024) has been migrated to platform(2024).

## HTTP OpenAPI

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
* `/terraform/v1/mgmt/hphls/update` HLS delivery in high performance mode.
* `/terraform/v1/mgmt/hphls/query` Query HLS delivery in high performance mode.
* `/terraform/v1/mgmt/ssl` Config the system SSL config.
* `/terraform/v1/mgmt/auto-self-signed-certificate` Create the self-signed certificate if no cert.
* `/terraform/v1/mgmt/letsencrypt` Config the let's encrypt SSL.
* `/terraform/v1/mgmt/cert/query` Query the key and cert for HTTPS.
* `/terraform/v1/mgmt/hooks/apply` Update the HTTP callback.
* `/terraform/v1/mgmt/hooks/query` Query the HTTP callback.
* `/terraform/v1/mgmt/hooks/example` Example target for HTTP callback.
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
* `/terraform/v1/ffmpeg/forward/secret` FFmpeg: Setup the forward secret to live streaming platforms.
* `/terraform/v1/ffmpeg/forward/streams` FFmpeg: Query the forwarding streams.
* `/terraform/v1/ffmpeg/vlive/secret` Setup the Virtual Live streaming secret.
* `/terraform/v1/ffmpeg/vlive/streams` Query the Virtual Live streaming streams.
* `/terraform/v1/ffmpeg/vlive/upload/` Upload Virtual Live source file.
* `/terraform/v1/ffmpeg/vlive/server/` Use server file as Virtual Live source.
* `/terraform/v1/ffmpeg/vlive/source` Setup Virtual Live source file.
* `/terraform/v1/ffmpeg/transcode/query` Query transcode config.
* `/terraform/v1/ffmpeg/transcode/apply` Apply transcode config.
* `/terraform/v1/ffmpeg/transcode/task` Query transcode task.

Also provided by platform for market:

* `/api/` SRS: HTTP API of SRS media server.
* `/rtc/` SRS: HTTP API for WebERTC of SRS media server.
* `/*/*.(flv|m3u8|ts|aac|mp3)` SRS: Media stream for HTTP-FLV, HLS, HTTP-TS, HTTP-AAC, HTTP-MP3.
* `/.well-known/acme-challenge/` HTTPS verify mount for letsencrypt.

Also provided by platform for static Files:

* `/tools/` A set of H5 tools, like simple player, xgplayer, etc, serve by mgmt.
* `/console/` The SRS console, serve by mgmt.
* `/players/` The SRS player, serve by mgmt.
* `/mgmt/` The ui for mgmt, serve by mgmt.

**Deprecated** API:

* `/terraform/v1/tencent/cam/secret` Tencent: Setup the CAM SecretId and SecretKey.
* `/terraform/v1/hooks/dvr/apply` Hooks: Apply the DVR pattern.
* `/terraform/v1/hooks/dvr/query` Hooks: Query the DVR pattern.
* `/terraform/v1/hooks/dvr/files` Hooks: List the DVR files.
* `/terraform/v1/hooks/dvr/hls/:uuid.m3u8` Hooks: Generate HLS/m3u8 url to preview or download.
* `/terraform/v1/hooks/vod/query` Hooks: Query the VoD pattern.
* `/terraform/v1/hooks/vod/apply` Hooks: Apply the VoD pattern.
* `/terraform/v1/hooks/vod/files` Hooks: List the VoD files.
* `/terraform/v1/hooks/vod/hls/:uuid.m3u8` Hooks: Generate HLS/m3u8 url to preview or download.

**Removed** API:

* `/terraform/v1/mgmt/strategy` Toggle the upgrade strategy.
* `/prometheus` Prometheus: Time-series database and monitor.
* `/terraform/v1/mgmt/nginx/proxy` Setup a reverse proxy location.
* `/terraform/v1/mgmt/dns/lb` HTTP-DNS for hls load balance.
* `/terraform/v1/mgmt/dns/backend/update` HTTP-DNS: Update the backend servers for hls load balance.
* `/terraform/v1/mgmt/nginx/homepage` Setup the homepage redirection.
* `/terraform/v1/mgmt/window/query` Query the upgrade time window.
* `/terraform/v1/mgmt/window/update` Update the upgrade time window.
* `/terraform/v1/mgmt/pubkey` Update the access for platform administrator pubkey.
* `/terraform/v1/mgmt/upgrade` Upgrade the mgmt to latest version.
* `/terraform/v1/mgmt/containers` Query SRS container.
* `/terraform/v1/host/exec` Exec command sync, response the stdout and stderr.
* `/terraform/v1/mgmt/secret/token` Create token for OpenAPI.

## Depends Softwares

The software we depend on:

* Docker, `apt-get install -y docker.io`
* Nginx, `apt-get install -y nginx`
    * Conf: `platform/containers/conf/nginx.conf`
    * Include: `platform/containers/data/config/nginx.http.conf`
    * Include: `platform/containers/data/config/nginx.server.conf`
    * SSL Key: `platform/containers/data/config/nginx.key`
    * Certificate: `platform/containers/data/config/nginx.crt`
* [LEGO](https://github.com/go-acme/lego)
    * Verify webroot: `platform/containers/data/.well-known/acme-challenge/`
    * Cert files: `platform/containers/data/lego/.lego/certificates/`
* [SRS](https://github.com/ossrs/srs)
    * Config: `platform/containers/conf/srs.release.conf` mount as `/usr/local/srs/conf/srs.conf`
    * Include: `platform/containers/data/config/srs.server.conf`
    * Include: `platform/containers/data/config/srs.vhost.conf`
    * Volume: `platform/containers/objs/nginx/` mount as `/usr/local/srs/objs/nginx/`
* FFmpeg:
    * [FFmpeg and ffprobe](https://ffmpeg.org) tools in `ossrs/srs:ubuntu20`

## Environment Variables

The optional environments defined by `platform/containers/data/config/.env`:

* `CLOUD`: `dev|bt|aapanel|droplet|docker`, The cloud platform name, DEV for development.
* `REGION`: `ap-guangzhou|ap-singapore|sgp1`, The region for upgrade source.
* `REGISTRY`: `docker.io|registry.cn-hangzhou.aliyuncs.com`, The docker registry.
* `MGMT_LISTEN`: The listen port for mgmt HTTP server. Default: 2022
* `PLATFORM_LISTEN`: The listen port for platform HTTP server. Default: 2024
* `HTTPS_LISTEN`: The listen port for HTTPS server. Default: 2443

For feature control:

* `NAME_LOOKUP`: `on|off`, whether enable the host name lookup, on or off. Default: off

For testing the specified service:

* `NODE_ENV`: `development|production`, if development, use local redis; otherwise, use `mgmt.srs.local` in docker.
* `LOCAL_RELEASE`: `on|off`, whether use local release service. Default: off
* `PLATFORM_DOCKER`: `on|off`, whether run platform in docker. Default: off

For mgmt and containers to connect to redis:

* `REDIS_PASSWORD`: The redis password. Default: empty.
* `REDIS_PORT`: The redis port. Default: 6379.

Environments for react ui:

* `PUBLIC_URL`: The mount prefix.
* `BUILD_PATH`: The output build path, default to `build`.

> Note: The env for react must start with `REACT_APP_`, please read [this post](https://create-react-app.dev/docs/adding-custom-environment-variables/#referencing-environment-variables-in-the-html).

Removed variables in .env:

* `SRS_PLATFORM_SECRET`: The mgmt api secret for token generating and verifying.

For HTTPS, automatically generate a self-signed certificate:

* `AUTO_SELF_SIGNED_CERTIFICATE`: `on|off`, whether generate self-signed certificate.

Deprecated and unused variables:

* `SRS_DOCKERIZED`: `on|off`, indicates the OS is in docker.
* `SRS_DOCKER`: `srs` to enfore use `ossrs/srs` docker image.
* `MGMT_DOCKER`: `on|off`, whether run mgmt in docker. Default: false
* `USE_DOCKER`: `on|off`, if false, disable all docker containers.
* `SRS_UTEST`: `on|off`, if on, running in utest mode.
* `SOURCE`: `github|gitee`, The source code for upgrading.

Please restart service when `.env` changed.
