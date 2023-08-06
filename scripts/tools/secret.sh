#!/bin/bash

# Execute by: bash xxx.sh or bash zzz/yyy/xxx.sh or ./xxx.sh or ./zzz/yyy/xxx.sh source xxx.sh
REALPATH=$(realpath ${BASH_SOURCE[0]})
SCRIPT_DIR=$(cd $(dirname ${REALPATH}) && pwd)
WORK_DIR=$(cd $(dirname ${REALPATH})/../.. && pwd)
cd ${WORK_DIR}

HELP=no
OUTPUT=

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -h|--help) HELP=yes; shift ;;
        --output) OUTPUT=$2; shift 2 ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
done

if [[ "$HELP" == yes ]]; then
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -h, --help    Show this help message and exit"
    echo "  --output      The output directory to save the final install file. Default: ${OUTPUT}"
    exit 0
fi

# Start by script.
SRS_PLATFORM_SECRET=$(docker exec srs-cloud redis-cli hget SRS_PLATFORM_SECRET token 2>/dev/null)
MGMT_PASSWORD=$(docker exec srs-cloud bash -c '. /data/config/.env && echo $MGMT_PASSWORD' 2>/dev/null)

# Start by BT.
if [[ -z $SRS_PLATFORM_SECRET ]]; then
    SRS_PLATFORM_SECRET=$(docker exec bt docker exec srs-cloud redis-cli hget SRS_PLATFORM_SECRET token 2>/dev/null)
    MGMT_PASSWORD=$(docker exec bt docker exec srs-cloud bash -c '. /data/config/.env && echo $MGMT_PASSWORD')
fi

# Start by aaPanel.
if [[ -z $SRS_PLATFORM_SECRET ]]; then
    SRS_PLATFORM_SECRET=$(docker exec aapanel docker exec srs-cloud redis-cli hget SRS_PLATFORM_SECRET token 2>/dev/null)
    MGMT_PASSWORD=$(docker exec aapanel docker exec srs-cloud bash -c '. /data/config/.env && echo $MGMT_PASSWORD')
fi

# Start by develop.
if [[ -z $SRS_PLATFORM_SECRET ]]; then
    SRS_PLATFORM_SECRET=$(redis-cli hget SRS_PLATFORM_SECRET token 2>/dev/null)
    MGMT_PASSWORD=$(. ${WORK_DIR}/platform/containers/data/config/.env && echo $MGMT_PASSWORD)
fi

if [[ ! -z $SRS_PLATFORM_SECRET ]]; then
    echo "SRS_PLATFORM_SECRET=$SRS_PLATFORM_SECRET"
    echo "MGMT_PASSWORD=$MGMT_PASSWORD"
    if [[ ! -z  $OUTPUT ]]; then
        echo "SRS_PLATFORM_SECRET=$SRS_PLATFORM_SECRET" > $OUTPUT
        echo "MGMT_PASSWORD=$MGMT_PASSWORD" >> $OUTPUT
    fi
    exit 0
fi

echo "Get SRS_PLATFORM_SECRET failed"
exit 1
