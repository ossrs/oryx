#!/bin/bash

# Execute by: bash xxx.sh or bash zzz/yyy/xxx.sh or ./xxx.sh or ./zzz/yyy/xxx.sh source xxx.sh
REALPATH=$(realpath ${BASH_SOURCE[0]})
SCRIPT_DIR=$(cd $(dirname ${REALPATH}) && pwd)
WORK_DIR=$(cd $(dirname ${REALPATH})/../.. && pwd)
echo "BASH_SOURCE=${BASH_SOURCE}, REALPATH=${REALPATH}, SCRIPT_DIR=${SCRIPT_DIR}, WORK_DIR=${WORK_DIR}"
cd ${WORK_DIR}

HELP=no
LANGUAGE=zh
EXTRACT=no
VERSION=$(bash ${SCRIPT_DIR}/version.sh)

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -h|--help) HELP=yes; shift ;;
        --language) LANGUAGE=$2; shift 2 ;;
        --version) VERSION=$2; shift 2 ;;
        --extract) EXTRACT=yes; shift 1 ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
done

if [[ "$HELP" == yes ]]; then
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -h, --help    Show this help message and exit"
    echo "  --language    The language to use. zh or en. Default: ${LANGUAGE}"
    echo "  --version     The image version to use. Default: ${VERSION}"
    echo "  --extract     Whether to extract the final install file. Default: ${EXTRACT}"
    exit 0
fi
echo "Install with options: LANGUAGE=${LANGUAGE}, VERSION=${VERSION}, EXTRACT=${EXTRACT}"

TMP_DIR="/tmp/srs-cloud-$(date +%s)" && TARGET_DIR="${TMP_DIR}/srs-cloud" &&
mkdir -p ${TARGET_DIR} && cd ${TARGET_DIR}
ret=$?; if [[ 0 -ne ${ret} ]]; then echo "mkdir ${TARGET_DIR} failed, ret=$ret"; exit $ret; fi
echo "Enter work directory ${TARGET_DIR}"

mkdir -p ${TARGET_DIR}/scripts/setup-ubuntu &&
cp -rf ${WORK_DIR}/scripts/setup-ubuntu/*.sh ${TARGET_DIR}/scripts/setup-ubuntu &&
ret=$?; if [[ 0 -ne ${ret} ]]; then echo "cp scripts failed, ret=$ret"; exit $ret; fi
echo "Copy scripts to ${TARGET_DIR}/scripts/setup-ubuntu"

cp -r ${WORK_DIR}/usr ${TARGET_DIR}/usr
ret=$?; if [[ 0 -ne ${ret} ]]; then echo "cp usr failed, ret=$ret"; exit $ret; fi
echo "Copy usr to ${TARGET_DIR}/usr"

cp ${WORK_DIR}/LICENSE ${TARGET_DIR}/LICENSE
ret=$?; if [[ 0 -ne ${ret} ]]; then echo "cp LICENSE failed, ret=$ret"; exit $ret; fi
echo "Copy LICENSE to ${TARGET_DIR}/LICENSE"

cp ${WORK_DIR}/README.md ${TARGET_DIR}/README.md
ret=$?; if [[ 0 -ne ${ret} ]]; then echo "cp README.md failed, ret=$ret"; exit $ret; fi
echo "Copy README.md to ${TARGET_DIR}/README.md"

mkdir -p ${TARGET_DIR}/mgmt &&
cp ${WORK_DIR}/mgmt/bootstrap ${TARGET_DIR}/mgmt/bootstrap
ret=$?; if [[ 0 -ne ${ret} ]]; then echo "cp bootstrap failed, ret=$ret"; exit $ret; fi
echo "Copy bootstrap to ${TARGET_DIR}/mgmt"

mkdir -p ${TARGET_DIR}/platform/containers &&
cp ${WORK_DIR}/platform/containers/.gitkeep ${TARGET_DIR}/platform/containers/.gitkeep
ret=$?; if [[ 0 -ne ${ret} ]]; then echo "cp containers failed, ret=$ret"; exit $ret; fi
echo "Copy containers to ${TARGET_DIR}/platform"

cat << END > ${TARGET_DIR}/scripts/setup-ubuntu/.env
LANGUAGE=${LANGUAGE}
IMAGE=ossrs/srs-cloud:${VERSION}
END

INSTALL_FILE=srs-cloud-${LANGUAGE}-${VERSION}.tar.gz
(cd ${TMP_DIR} && tar zcf $INSTALL_FILE srs-cloud) &&
ret=$?; if [[ 0 -ne ${ret} ]]; then echo "tar failed, ret=$ret"; exit $ret; fi
echo "Tar $INSTALL_FILE"

mkdir -p ${WORK_DIR}/build &&
mv ${TMP_DIR}/$INSTALL_FILE ${WORK_DIR}/build
ret=$?; if [[ 0 -ne ${ret} ]]; then echo "mv failed, ret=$ret"; exit $ret; fi
echo "Move $INSTALL_FILE to ${WORK_DIR}/build"

if [[ $EXTRACT == yes ]]; then
    (cd ${WORK_DIR}/build && tar zxf $INSTALL_FILE)
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "tar failed, ret=$ret"; exit $ret; fi
    echo "Extract $INSTALL_FILE to ${WORK_DIR}/build"
fi

rm -rf $TMP_DIR
echo "Remove old ${TMP_DIR}"

