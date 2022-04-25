#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/../.. && pwd)
echo "Run setup at $WORK_DIR from $0"
cd $WORK_DIR

# The main directory.
DEPLOY_HOME=/usr/local/lighthouse/softwares
SRS_HOME=${DEPLOY_HOME}/srs-cloud
INSTALL_HOME=/usr/local/srs-cloud

########################################################################################################################
# Check OS first, only support CentOS 7.
yum --version >/dev/null 2>&1 && rpm --version >/dev/null 2>&1
if [[ $? -ne 0 ]]; then echo "Only support CentOS 7"; exit 1; fi

# Check CentOS version.
CentOS_VERSION=$(rpm --eval '%{centos_ver}')
if [[ $CentOS_VERSION -ne 7 ]]; then echo "Only support CentOS 7, yours is $CentOS_VERSION"; exit 1; fi

# User should install nodejs, because we can't do it.
(cd scripts/check-node-version && npm install && node .)
if [[ $? -ne 0 ]]; then echo "Please install node from https://nodejs.org"; exit 1; fi

# Check user lighthouse
if [[ $(id -un lighthouse 2>/dev/null) == '' ]]; then
  echo "No user lighthouse"; exit 1;
fi

# Check user lighthouse home directory.
if [[ ! -d ~lighthouse ]]; then
  echo "No home directory ~lighthouse"; exit 1;
fi

# Ignore darwin
if [[ $(uname -s) == 'Darwin' ]]; then
  echo "Mac is not supported"; exit 1;
fi

########################################################################################################################
# Install depends services, except nodejs.
yum install -y git gcc-c++ gdb make tree dstat docker nginx &&
systemctl enable docker nginx
if [[ $? -ne 0 ]]; then echo "Install dependencies failed"; exit 1; fi

# Install files to lighthouse directory.
mkdir -p $DEPLOY_HOME &&
rm -rf ${SRS_HOME} &&
(cd $(dirname $WORK_DIR) && cp -r $(basename $WORK_DIR) ${SRS_HOME}) &&
cd ${SRS_HOME} &&
make build && make install
if [[ $? -ne 0 ]]; then echo "Copy srs-cloud failed"; exit 1; fi

cd $DEPLOY_HOME && rm -rf srs-terraform && ln -sf srs-cloud srs-terraform
if [[ $? -ne 0 ]]; then echo "Link srs-cloud failed"; exit 1; fi

########################################################################################################################
# Cache the docker images for srs-cloud to startup faster.
systemctl start docker &&
echo "Cache docker images" &&
docker pull registry.cn-hangzhou.aliyuncs.com/ossrs/srs:4 &&
docker pull registry.cn-hangzhou.aliyuncs.com/ossrs/lighthouse:4 &&
docker pull registry.cn-hangzhou.aliyuncs.com/ossrs/node:slim &&
docker pull registry.cn-hangzhou.aliyuncs.com/ossrs/srs-cloud:hooks-1 &&
docker pull registry.cn-hangzhou.aliyuncs.com/ossrs/srs-cloud:tencent-1 &&
docker pull registry.cn-hangzhou.aliyuncs.com/ossrs/srs-cloud:ffmpeg-1 &&
docker pull registry.cn-hangzhou.aliyuncs.com/ossrs/srs-cloud:platform-1 &&
docker pull registry.cn-hangzhou.aliyuncs.com/ossrs/prometheus &&
docker pull registry.cn-hangzhou.aliyuncs.com/ossrs/redis_exporter &&
docker pull registry.cn-hangzhou.aliyuncs.com/ossrs/node-exporter &&
docker pull registry.cn-hangzhou.aliyuncs.com/ossrs/certbot &&
docker pull registry.cn-hangzhou.aliyuncs.com/ossrs/redis
if [[ $? -ne 0 ]]; then echo "Cache docker images failed"; exit 1; fi

# If install ok, the directory should exists.
if [[ ! -d ${INSTALL_HOME} || ! -d ${INSTALL_HOME}/mgmt ]]; then
  echo "Install srs-cloud failed"; exit 1;
fi

# Compatible with previous version.
cd $(dirname $INSTALL_HOME) && rm -rf srs-terraform && ln -sf srs-cloud srs-terraform
if [[ $? -ne 0 ]]; then echo "Link srs-cloud failed"; exit 1; fi

# Create srs-cloud service, and the credential file.
# Remark: Never start the service, because the IP will change for new machine created.
cd ${INSTALL_HOME} &&
cp -f usr/lib/systemd/system/srs-cloud.service /usr/lib/systemd/system/srs-cloud.service &&
touch ${INSTALL_HOME}/mgmt/.env &&
systemctl enable srs-cloud
if [[ $? -ne 0 ]]; then echo "Install srs-cloud failed"; exit 1; fi

# Choose default language.
cat << END > ${SRS_HOME}/mgmt/.env
REACT_APP_LOCALE=zh
END
if [[ $? -ne 0 ]]; then echo "Setup language failed"; exit 1; fi

# Setup the nginx configuration.
rm -f /etc/nginx/nginx.conf &&
ln -sf ${SRS_HOME}/mgmt/containers/conf/nginx.conf /etc/nginx/nginx.conf
if [[ $? -ne 0 ]]; then echo "Setup nginx config failed"; exit 1; fi

# Build the mgmt/containers/conf/conf.d/nginx.vhost.conf
cd ${SRS_HOME}/mgmt && bash auto/setup_vhost
if [[ $? -ne 0 ]]; then echo "Build nginx vhost failed"; exit 1; fi

cd ${SRS_HOME}/mgmt &&
rm -f /etc/nginx/conf.d/nginx.vhost.conf /etc/nginx/conf.d/server.conf &&
ln -sf ${SRS_HOME}/mgmt/containers/conf/conf.d/nginx.vhost.conf /etc/nginx/conf.d/vhost.conf
if [[ $? -ne 0 ]]; then echo "Reload nginx failed"; exit 1; fi

# Setup git alias to make it convenient.
cd ${SRS_HOME}/mgmt &&
echo "Setup git alias to make it more convenient" &&
git config --local alias.co checkout &&
git config --local alias.br branch &&
git config --local alias.ci commit &&
git config --local alias.st status
if [[ $? -ne 0 ]]; then echo "Setup git alias failed"; exit 1; fi

# Update sysctl.conf and add if not exists. For example:
#   update_sysctl net.ipv4.ip_forward 1 0 "# Controls IP packet forwarding"
function update_sysctl() {
    SYSCTL_KEY=$1 && SYSCTL_VALUE=$2 && SYSCTL_EMPTY_LINE=$3 && SYSCTL_COMMENTS=$4
    echo "Update with sysctl $SYSCTL_KEY=$SYSCTL_VALUE, empty-line=$SYSCTL_EMPTY_LINE, comment=$SYSCTL_COMMENTS"

    grep -q "^${SYSCTL_KEY}[ ]*=" /etc/sysctl.conf
    if [[ $? == 0 ]]; then
      sed -i "s/^${SYSCTL_KEY}[ ]*=.*$/${SYSCTL_KEY} = ${SYSCTL_VALUE}/g" /etc/sysctl.conf
    else
      if [[ $SYSCTL_EMPTY_LINE == 1 ]]; then echo '' >> /etc/sysctl.conf; fi &&
      if [[ $SYSCTL_COMMENTS != '' ]]; then echo "$SYSCTL_COMMENTS" >> /etc/sysctl.conf; fi &&
      echo "${SYSCTL_KEY} = ${SYSCTL_VALUE}" >> /etc/sysctl.conf
    fi
    if [[ $? -ne 0 ]]; then echo "Failed to sysctl $SYSCTL_KEY = $SYSCTL_VALUE $SYSCTL_COMMENTS"; exit 1; fi

    RESULT=$(grep "^${SYSCTL_KEY}[ ]*=" /etc/sysctl.conf)
    echo "Update done: ${RESULT}"
}

# Allow network forwarding, required by docker.
# See https://stackoverflow.com/a/41453306/17679565
update_sysctl net.ipv4.ip_forward 1 1 "# Controls IP packet forwarding"

# Setup the UDP buffer for WebRTC and SRT.
# See https://www.jianshu.com/p/6d4a89359352
update_sysctl net.core.rmem_max 16777216 1 "# For RTC/SRT over UDP"
update_sysctl net.core.rmem_default 16777216
update_sysctl net.core.wmem_max 16777216
update_sysctl net.core.wmem_default 16777216

########################################################################################################################
# Note that we keep files as root, because we run srs-cloud as root, see https://stackoverflow.com/a/70953525/17679565
chown lighthouse:lighthouse ${INSTALL_HOME}/mgmt/.env
if [[ $? -ne 0 ]]; then echo "Link files failed"; exit 1; fi

rm -rf ~lighthouse/credentials.txt && ln -sf ${INSTALL_HOME}/mgmt/.env ~lighthouse/credentials.txt &&
if [[ $? -ne 0 ]]; then echo "Link files failed"; exit 1; fi

