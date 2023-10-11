#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/../.. && pwd)
echo "Run setup at $WORK_DIR from $0"
cd $WORK_DIR

help=no
ip=
os=
user=
password=
cleanup=yes
SOURCE=$WORK_DIR
VERSION=$(bash ${WORK_DIR}/scripts/version.sh)

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -h|--help) help=yes; shift ;;
        --ip) ip=$2; shift 2;;
        --os) os=$2; shift 2;;
        --user) user=$2; shift 2;;
        --password) password=$2; shift 2;;
        --cleanup) cleanup=$2; shift 2;;
        --version) VERSION=$2; shift 2 ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
done

if [[ "$help" == yes ]]; then
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -h, --help      Show this help message and exit"
    echo "  --ip            The ip of server to build. For example: xx.xxx.xxx.xxx"
    echo "  --os            The os of server to build. For example: ubuntu"
    echo "  --user          The user of server to build. For example: ubuntu"
    echo "  --password      The password of server to build."
    echo "  --cleanup       Whether do cleanup. yes or no. Default: $cleanup"
    echo "  --version     The image version to use. Default: ${VERSION}"
    exit 0
fi

if [[ -z $ip ]]; then echo "No ip"; exit 1; fi
if [[ -z $os ]]; then echo "No os"; exit 1; fi
if [[ -z $user ]]; then echo "No user"; exit 1; fi
if [[ -z $password ]]; then echo "No password"; exit 1; fi
if [[ -z $VERSION ]]; then echo "No VERSION"; exit 1; fi

IMAGE_URL="registry.cn-hangzhou.aliyuncs.com/ossrs/srs-stack:$VERSION"
echo "SOURCE=$SOURCE, ip=$ip, os=$os, user=$user, password=${#password}B, cleanup=$cleanup, VERSION=$VERSION, IMAGE_URL=$IMAGE_URL"

sshCmd="sshpass -p $password ssh -o StrictHostKeyChecking=no"
scpCmd="sshpass -p $password scp -o StrictHostKeyChecking=no"

$sshCmd -t $user@$ip "hostname" && echo "Check sshpass ok"
if [[ $ret -ne 0 ]]; then echo "Check sshpass failed"; echo "See https://stackoverflow.com/a/32258393/17679565"; exit 1; fi

SRS_HOME=/tmp/lighthouse/srs-stack &&
rm -rf $(dirname $SRS_HOME) && mkdir -p $SRS_HOME &&
echo "mkdir $SRS_HOME ok"
ret=$?; if [[ 0 -ne $ret ]]; then echo "mkdir $SRS_HOME failed, ret=$ret"; exit $ret; fi

rm -rf ${SRS_HOME} &&
mkdir -p $SRS_HOME/mgmt ${SRS_HOME}/scripts ${SRS_HOME}/platform/containers/conf &&
cp -r ${SOURCE}/usr ${SRS_HOME}/usr &&
cp -r ${SOURCE}/scripts/setup-lighthouse ${SRS_HOME}/scripts/setup-lighthouse &&
cp ${SOURCE}/LICENSE ${SRS_HOME}/LICENSE &&
cp ${SOURCE}/README.md ${SRS_HOME}/README.md &&
cp ${SOURCE}/mgmt/bootstrap ${SRS_HOME}/mgmt/bootstrap &&
cp ${SOURCE}/platform/containers/conf/nginx.conf ${SRS_HOME}/platform/containers/conf/nginx.conf
if [[ $? -ne 0 ]]; then echo "Copy srs-stack failed"; exit 1; fi

echo "Start to update bootstrap"
sed -i '' "s|^IMAGE=.*|IMAGE=${IMAGE_URL}|g" ${SRS_HOME}/mgmt/bootstrap
if [[ $? -ne 0 ]]; then echo "Update bootstrap failed"; exit 1; fi
echo "Update bootstrap ok"

tgzName=/tmp/lighthouse/srs-stack.zip &&
(cd $(dirname $tgzName) && rm -f $tgzName && zip -q -r $tgzName $(basename $SRS_HOME)) &&
echo "Package $tgzName ok" && ls -lh $tgzName
if [[ $? -ne 0 ]]; then echo "Package $tgzName failed"; exit 1; fi

$scpCmd $tgzName $user@$ip:~
if [[ $? -ne 0 ]]; then echo "Copy failed"; exit 1; fi
echo "Copy $tgzName to $ip ok"

tgzFile=$(basename $tgzName) &&
SRS_NAME=$(basename $SRS_HOME) &&
echo "Run command on server: $ip" &&
echo "  unzip -q $tgzFile"
echo "  bash ~/$SRS_NAME/scripts/setup-lighthouse/setup_lighthouse.sh"

$sshCmd -t $user@$ip "
    rm -rf $SRS_NAME && unzip -q $tgzFile && \
    sudo bash $SRS_NAME/scripts/setup-lighthouse/setup_lighthouse.sh &&
    if [[ $cleanup == yes ]]; then
        sudo bash $SRS_NAME/scripts/setup-lighthouse/post_build.sh &&
        sudo rm -rf $tgzFile $SRS_NAME; \
    fi && \
    echo OK; \
"
if [[ $? -ne 0 ]]; then echo "Run setup failed"; exit 1; fi

echo ""
echo "OK"
