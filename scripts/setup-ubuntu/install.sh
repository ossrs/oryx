#!/bin/bash

# Execute by: bash xxx.sh or bash zzz/yyy/xxx.sh or ./xxx.sh or ./zzz/yyy/xxx.sh source xxx.sh
REALPATH=$(realpath ${BASH_SOURCE[0]})
SCRIPT_DIR=$(cd $(dirname ${REALPATH}) && pwd)
WORK_DIR=$(cd $(dirname ${REALPATH})/../.. && pwd)
echo "BASH_SOURCE=${BASH_SOURCE}, REALPATH=${REALPATH}, SCRIPT_DIR=${SCRIPT_DIR}, WORK_DIR=${WORK_DIR}"
cd ${WORK_DIR}

DATA_HOME=/data
SRS_HOME=/usr/local/oryx

HELP=no
VERBOSE=no
LANGUAGE=zh
REGISTRY=auto
REGION=auto
IMAGE=ossrs/oryx:5

# Allow use .env to override the default values.
if [[ -f ${SCRIPT_DIR}/.env ]]; then source ${SCRIPT_DIR}/.env; fi

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -h|--help) HELP=yes; shift ;;
        --verbose) VERBOSE=yes; shift ;;
        --language) LANGUAGE=$2; shift 2 ;;
        --registry) REGISTRY=$2; shift 2 ;;
        --image) IMAGE=$2; shift 2 ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
done

function help() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -h, --help    Show this help message and exit"
    echo "  --verbose     Whether to show verbose log. Default: ${VERBOSE}"
    echo "  --language    The language to use. zh or en. Default: ${LANGUAGE}"
    echo "  --registry    The registry for docker images. Default: ${REGISTRY}"
    echo "  --image       The image to run. Default: ${IMAGE}"
}

# Guess the registry automatically by language.
if [[ $REGISTRY == auto ]]; then
    REGISTRY=$([[ $LANGUAGE == zh ]] && echo registry.cn-hangzhou.aliyuncs.com || echo docker.io)
    REGION=$([[ $LANGUAGE == zh ]] && echo ap-beijing || echo ap-singapore)
    IMAGE_URL=$([[ $REGISTRY == docker.io ]] && echo ${IMAGE} || echo ${REGISTRY}/${IMAGE})
fi

if [[ "$HELP" == yes ]]; then
    help
    exit 0
fi
echo -n "Install with options: VERBOSE=${VERBOSE}, LANGUAGE=${LANGUAGE}, SRS_HOME=${SRS_HOME}, DATA_HOME=${DATA_HOME}"
echo ", REGISTRY=${REGISTRY}, IMAGE=${IMAGE}, REGION=${REGION}, IMAGE_URL=${IMAGE_URL}"

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

# Check OS first, only support Ubuntu
apt-get --version >/dev/null 2>&1 && OS_NAME='Ubuntu'
if [[ -z $OS_NAME ]]; then echo "Only support Ubuntu"; exit 1; fi

# ubuntu:xenial is 16
# ubuntu:bionic is 18
# ubuntu:focal is 20
# ubuntu:jammy is 22
OS_VERSION=$(source /etc/os-release && echo $VERSION_ID |cut -d . -f 1)
if [[ $OS_VERSION -lt 18 ]]; then echo "Only support Ubuntu 18 and higher, yours is $OS_VERSION"; exit 1; fi
echo "Check OS_VERSION=$OS_VERSION ok"

# Requires systemd and docker.io
if [[ $(systemctl --version >/dev/null 2>&1 || echo no) == no ]]; then
    echo "Requires systemd, please install by:"
    echo "    sudo apt-get install systemd"
    exit 1
fi

echo "Start to install files"
mkdir -p ${SRS_HOME} ${DATA_HOME} && rm -rf ${SRS_HOME}/* &&
cp -r ${WORK_DIR}/usr ${SRS_HOME}/usr &&
cp -r ${WORK_DIR}/mgmt ${SRS_HOME}/mgmt &&
cp -r ${WORK_DIR}/LICENSE ${SRS_HOME}/LICENSE
ret=$?; if [[ $ret -ne 0 ]]; then echo "Copy files failed, ret=$ret"; exit $ret; fi
echo "Install files at ${SRS_HOME} ok"

if [[ ${VERBOSE} == yes ]]; then
    echo "Source files:"
    echo $files
    echo "Total size:"
    du -sh ${SRS_HOME}
    echo "Detail files:"
    du -sh ${SRS_HOME}/*
fi

echo "Start to create data and config files"
mkdir -p ${DATA_HOME}/config && touch ${DATA_HOME}/config/.env
if [[ $? -ne 0 ]]; then echo "Create /data/config/.env failed"; exit 1; fi
echo "Create data and config files ok"

echo "Start to update bootstrap"
sed -i "s|^DATA_HOME=.*|DATA_HOME=${DATA_HOME}|g" ${SRS_HOME}/mgmt/bootstrap &&
sed -i "s|^IMAGE=.*|IMAGE=${IMAGE_URL}|g" ${SRS_HOME}/mgmt/bootstrap
if [[ $? -ne 0 ]]; then echo "Update bootstrap failed"; exit 1; fi
echo "Update bootstrap ok"

# Allow network forwarding, required by docker.
# See https://stackoverflow.com/a/41453306/17679565
echo "Start controls IP packet forwarding"
update_sysctl net.ipv4.ip_forward 1 1 "# Controls IP packet forwarding"
echo "Controls IP packet forwarding ok"

# Setup the UDP buffer for WebRTC and SRT.
# See https://www.jianshu.com/p/6d4a89359352
echo "Start setup kernel UDP buffer"
update_sysctl net.core.rmem_max 16777216 1 "# For RTC/SRT over UDP"
update_sysctl net.core.rmem_default 16777216
update_sysctl net.core.wmem_max 16777216
update_sysctl net.core.wmem_default 16777216
echo "Setup kernel UDP buffer ok"

# For BT, we use special env, to disable discover of platform.
echo "Start to setup .env"
if [[ -f ${DATA_HOME}/config/.env && -s ${DATA_HOME}/config/.env ]]; then
    echo "The .env already exists, skip"
else
    cat << END > ${DATA_HOME}/config/.env
CLOUD=BIN
REGION=${REGION}
REACT_APP_LOCALE=${LANGUAGE}
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

# If install ok, the directory should exists.
if [[ ! -d ${SRS_HOME} ]]; then
  echo "Install oryx failed"; exit 1;
fi

# Create init.d script.
rm -f /etc/init.d/oryx &&
cp ${SCRIPT_DIR}/init.d.sh /etc/init.d/oryx &&
chmod +x /etc/init.d/oryx
if [[ $? -ne 0 ]]; then echo "Setup init.d script failed"; exit 1; fi

# Create oryx service.
# Remark: Never start the service, because the IP will change for new machine created.
cd ${SRS_HOME} &&
cp -f usr/lib/systemd/system/oryx.service /usr/lib/systemd/system/oryx.service &&
systemctl daemon-reload && systemctl enable oryx
if [[ $? -ne 0 ]]; then echo "Install oryx failed"; exit 1; fi

/etc/init.d/oryx restart oryx
if [[ $? -ne 0 ]]; then echo "Start oryx failed"; exit 1; fi

echo 'Install OK'
