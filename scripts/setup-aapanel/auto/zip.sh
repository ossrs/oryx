#!/bin/bash

# Execute by: bash xxx.sh or bash zzz/yyy/xxx.sh or ./xxx.sh or ./zzz/yyy/xxx.sh source xxx.sh
REALPATH=$(realpath ${BASH_SOURCE[0]})
SCRIPT_DIR=$(cd $(dirname ${REALPATH}) && pwd)
WORK_DIR=$(cd $(dirname ${REALPATH})/../../.. && pwd)
echo "BASH_SOURCE=${BASH_SOURCE}, REALPATH=${REALPATH}, SCRIPT_DIR=${SCRIPT_DIR}, WORK_DIR=${WORK_DIR}"
cd ${WORK_DIR}

HELP=no
EXTRACT=no
OUTPUT=${WORK_DIR}/build
VERSION=$(bash ${WORK_DIR}/scripts/version.sh)

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -h|--help) HELP=yes; shift ;;
        --version) VERSION=$2; shift 2 ;;
        --output) OUTPUT=$2; shift 2 ;;
        --extract) EXTRACT=yes; shift 1 ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
done

if [[ "$HELP" == yes ]]; then
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -h, --help    Show this help message and exit"
    echo "  --version     The image version to use. Default: ${VERSION}"
    echo "  --output      The output directory to save the final install file. Default: ${OUTPUT}"
    echo "  --extract     Whether to extract the final install file. Default: ${EXTRACT}"
    exit 0
fi

OUTPUT=$(cd ${WORK_DIR} && mkdir -p ${OUTPUT} && cd ${OUTPUT} && pwd)
echo "Install with options: VERSION=${VERSION}, OUTPUT=${OUTPUT}, EXTRACT=${EXTRACT}"

TMP_DIR="/tmp/srs-stack-$(date +%s)" && TARGET_DIR="${TMP_DIR}/srs_cloud" && mkdir -p ${TARGET_DIR}
ret=$?; if [[ 0 -ne ${ret} ]]; then echo "mkdir ${TARGET_DIR} failed, ret=$ret"; exit $ret; fi
echo "Create tmp dir ${TARGET_DIR}"

mkdir -p $TARGET_DIR &&
cp -r scripts/setup-aapanel/* $TARGET_DIR && rm -rf $TARGET_DIR/auto &&
rm -f $TARGET_DIR/bt_tools.py && cp -f scripts/tools/bt_*.py $TARGET_DIR &&
cp -r ${WORK_DIR}/usr ${TARGET_DIR}/usr &&
cp ${WORK_DIR}/LICENSE ${TARGET_DIR}/LICENSE &&
mkdir -p ${TARGET_DIR}/mgmt && cp ${WORK_DIR}/mgmt/bootstrap ${TARGET_DIR}/mgmt/bootstrap
echo "Copy files to $TARGET_DIR"
if [[ $? -ne 0 ]]; then echo "Copy files failed"; exit 1; fi

# For aaPanel, should never use .env, because it will be removed when install.
cat << END > $TARGET_DIR/config
LANGUAGE=en
IMAGE=ossrs/srs-stack:${VERSION}
END
if [[ $? -ne 0 ]]; then echo "Generate config failed"; exit 1; fi
echo "Generate config to $TARGET_DIR/config"

INSTALL_FILE=aapanel-srs_cloud.zip
(cd $TMP_DIR/ && zip -q -r $INSTALL_FILE srs_cloud) &&
echo "Zip generated at $TMP_DIR/$INSTALL_FILE"
if [[ $? -ne 0 ]]; then echo "Zip plugin failed"; exit 1; fi

mkdir -p ${OUTPUT} &&
mv ${TMP_DIR}/$INSTALL_FILE ${OUTPUT}
ret=$?; if [[ 0 -ne ${ret} ]]; then echo "mv failed, ret=$ret"; exit $ret; fi
echo "Move $INSTALL_FILE to ${OUTPUT}"

if [[ $EXTRACT == yes ]]; then
    (cd ${OUTPUT} && rm -rf srs_cloud && unzip -q $INSTALL_FILE)
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "tar failed, ret=$ret"; exit $ret; fi
    echo "Extract $INSTALL_FILE to ${OUTPUT}"
fi

rm -rf $TMP_DIR
echo "Remove old ${TMP_DIR}"

