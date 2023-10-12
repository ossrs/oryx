#!/bin/bash

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

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/../.. && pwd)
echo "Run setup at $WORK_DIR from $0"
cd $WORK_DIR

# The main directory.
SRS_HOME=/usr/local/srs-stack
DATA_HOME=/data
IMAGE_URL=registry.cn-hangzhou.aliyuncs.com/ossrs/srs-stack:5
SOURCE=$WORK_DIR

mkdir -p /usr/local/lighthouse/softwares/srs-stack &&
rm -rf $SRS_HOME && ln -sf /usr/local/lighthouse/softwares/srs-stack $SRS_HOME &&
(cd /usr/local/lighthouse/softwares && rm -rf srs-terraform && ln -sf srs-stack srs-terraform)
ret=$?; if [[ 0 -ne $ret ]]; then echo "Failed to create $SRS_HOME"; exit $ret; fi

if [[ $(id -un lighthouse 2>/dev/null) == '' ]]; then
    bash scripts/setup-lighthouse/create_lighthouse_user.sh
    if [[ $? -ne 0 ]]; then echo "Create user lighthouse failed"; exit 1; fi
fi

########################################################################################################################
# Check OS first, only support CentOS 7.
apt-get --version >/dev/null 2>&1 && OS_NAME='Ubuntu'
if [[ -z $OS_NAME ]]; then echo "Only support Ubuntu"; exit 1; fi

if [[ $OS_NAME == 'Ubuntu' ]]; then
  # Check Ubuntu version.
  Ubuntu_VERSION=$(cat /etc/os-release |grep VERSION_ID |awk -F '"' '{print $2}' |awk -F '.' '{print $1}')
  if [[ $Ubuntu_VERSION -lt 18 ]]; then echo "Only support Ubuntu 18+, yours is $Ubuntu_VERSION"; exit 1; fi
fi

# Check user lighthouse
if [[ $(id -un lighthouse 2>/dev/null) == '' ]]; then
  echo "No user lighthouse"; exit 1;
fi

# Check user lighthouse home directory.
if [[ ! -d ~lighthouse ]]; then
  echo "No home directory ~lighthouse"; exit 1;
fi

########################################################################################################################
# Install depends services. Retry because apt-get might be busy.
for ((i=0; i<3; i++)); do
    apt-get update -y &&
    apt-get install -y git gcc g++ gdb make tree dstat docker docker.io nginx curl net-tools &&
    apt-get -qqy clean
    ret=$?; if [[ $ret -eq 0 ]]; then break; fi
    sleep 5;
done
if [[ $ret -ne 0 ]]; then echo "Install dependencies failed"; exit 1; fi

echo "Enable service" &&
systemctl enable docker nginx &&
systemctl start docker
if [[ $? -ne 0 ]]; then echo "Enable service failed"; exit 1; fi

# Allow network forwarding, required by docker.
# See https://stackoverflow.com/a/41453306/17679565
update_sysctl net.ipv4.ip_forward 1 1 "# Controls IP packet forwarding"

# Setup the UDP buffer for WebRTC and SRT.
# See https://www.jianshu.com/p/6d4a89359352
update_sysctl net.core.rmem_max 16777216 1 "# For RTC/SRT over UDP"
update_sysctl net.core.rmem_default 16777216
update_sysctl net.core.wmem_max 16777216
update_sysctl net.core.wmem_default 16777216

# Install files to lighthouse directory.
cp -r ${SOURCE}/usr ${SRS_HOME}/usr &&
cp ${SOURCE}/LICENSE ${SRS_HOME}/LICENSE &&
cp ${SOURCE}/README.md ${SRS_HOME}/README.md &&
mkdir -p ${SRS_HOME}/mgmt && cp ${SOURCE}/mgmt/bootstrap ${SRS_HOME}/mgmt/bootstrap
if [[ $? -ne 0 ]]; then echo "Copy srs-stack failed"; exit 1; fi

########################################################################################################################
echo "Start to create data and config files"
mkdir -p ${DATA_HOME}/config && touch ${DATA_HOME}/config/.env
if [[ $? -ne 0 ]]; then echo "Create /data/config failed"; exit 1; fi
echo "Create data and config files ok"

# Setup the nginx configuration.
rm -f /etc/nginx/nginx.conf &&
cp ${SOURCE}/platform/containers/conf/nginx.conf /etc/nginx/nginx.conf &&
sed -i "s/user nginx;/user www-data;/g" /etc/nginx/nginx.conf &&
touch ${DATA_HOME}/config/nginx.http.conf ${DATA_HOME}/config/nginx.server.conf
if [[ $? -ne 0 ]]; then echo "Setup nginx config failed"; exit 1; fi

echo "Start to update bootstrap"
sed -i "s|^DATA_HOME=.*|DATA_HOME=${DATA_HOME}|g" ${SRS_HOME}/mgmt/bootstrap &&
sed -i "s|^IMAGE=.*|IMAGE=${IMAGE_URL}|g" ${SRS_HOME}/mgmt/bootstrap
if [[ $? -ne 0 ]]; then echo "Update bootstrap failed"; exit 1; fi
echo "Update bootstrap ok"

# Choose default language.
echo "Start to setup .env"
if [[ -f ${DATA_HOME}/config/.env && -s ${DATA_HOME}/config/.env ]]; then
    echo "The .env already exists, skip"
else
    cat << END > ${DATA_HOME}/config/.env
CLOUD=TENCENT
REACT_APP_LOCALE=zh
IMAGE=${IMAGE_URL}
END
    if [[ $? -ne 0 ]]; then echo "Setup .env failed"; exit 1; fi
fi

# Update the docker images.
echo "Cache docker images" &&
if [[ $(docker images --format "{{.Repository}}:{{.Tag}}" ${IMAGE_URL} |wc -l) -eq 1 ]]; then
    echo "Docker images ${IMAGE_URL} exists, skip pull"
else
    docker pull ${IMAGE_URL}
    if [[ $? -ne 0 ]]; then echo "Cache docker images failed"; exit 1; fi
fi

# Create srs-stack service, and the credential file.
# Remark: Never start the service, because the IP will change for new machine created.
cp -f ${SRS_HOME}/usr/lib/systemd/system/srs-stack.service /usr/lib/systemd/system/srs-stack.service &&
systemctl daemon-reload && systemctl enable srs-stack
if [[ $? -ne 0 ]]; then echo "Install srs-stack failed"; exit 1; fi

########################################################################################################################
# Create srs-cloud soft link, to keep compatible with lighthouse HTTPS management. If user user lighthouse to
# setup the HTTPS and domain, the lighthouse will create a file like /etc/nginx/init.d/yourdomain.conf, which
# includes the file /usr/local/lighthouse/softwares/srs-cloud/mgmt/containers/conf/default.d/proxy.conf, which
# finally proxy to srs-stack.
(cd /usr/local/lighthouse/softwares && rm -rf srs-cloud && ln -sf srs-stack srs-cloud) &&
mkdir -p /usr/local/lighthouse/softwares/srs-cloud/mgmt/containers/conf/default.d &&
cat << END > /usr/local/lighthouse/softwares/srs-cloud/mgmt/containers/conf/default.d/proxy.conf
  location / {
    proxy_pass http://127.0.0.1:2022;
    proxy_set_header Host \$host;
  }
END
if [[ $? -ne 0 ]]; then echo "Compatible lighthouse HTTPS failed"; exit 1; fi

########################################################################################################################
# Note that we keep files as root, because we run srs-stack as root, see https://stackoverflow.com/a/70953525/17679565
chown lighthouse:lighthouse ${DATA_HOME}/config/.env
if [[ $? -ne 0 ]]; then echo "Link files failed"; exit 1; fi

rm -rf ~lighthouse/credentials.txt && ln -sf ${DATA_HOME}/config/.env ~lighthouse/credentials.txt
if [[ $? -ne 0 ]]; then echo "Link files failed"; exit 1; fi

echo "Install srs-stack ok"
