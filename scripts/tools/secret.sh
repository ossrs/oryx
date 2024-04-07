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

if [[ ! -z  $OUTPUT ]]; then
    echo "Truncate output file: $OUTPUT"
    echo "" > $OUTPUT
fi

for ((i=0; i<3; i++)); do
    # Start by docker.
    SRS_PLATFORM_SECRET=$(docker exec oryx redis-cli hget SRS_PLATFORM_SECRET token 2>/dev/null)
    MGMT_PASSWORD=$(docker exec oryx bash -c '. /data/config/.env && echo $MGMT_PASSWORD' 2>/dev/null)

    # Start by script.
    if [[ -z $SRS_PLATFORM_SECRET ]]; then
        SRS_PLATFORM_SECRET=$(docker exec script docker exec oryx redis-cli hget SRS_PLATFORM_SECRET token 2>/dev/null)
        MGMT_PASSWORD=$(docker exec script docker exec oryx bash -c '. /data/config/.env && echo $MGMT_PASSWORD' 2>/dev/null)
    fi

    # Start by BT.
    if [[ -z $SRS_PLATFORM_SECRET ]]; then
        SRS_PLATFORM_SECRET=$(docker exec bt docker exec oryx redis-cli hget SRS_PLATFORM_SECRET token 2>/dev/null)
        MGMT_PASSWORD=$(docker exec bt docker exec oryx bash -c '. /data/config/.env && echo $MGMT_PASSWORD' 2>/dev/null)
    fi

    # Start by aaPanel.
    if [[ -z $SRS_PLATFORM_SECRET ]]; then
        SRS_PLATFORM_SECRET=$(docker exec aapanel docker exec oryx redis-cli hget SRS_PLATFORM_SECRET token 2>/dev/null)
        MGMT_PASSWORD=$(docker exec aapanel docker exec oryx bash -c '. /data/config/.env && echo $MGMT_PASSWORD' 2>/dev/null)
    fi

    # Start by develop.
    if [[ -z $SRS_PLATFORM_SECRET ]]; then
        SRS_PLATFORM_SECRET=$(redis-cli hget SRS_PLATFORM_SECRET token 2>/dev/null)
        if [[ -f ${WORK_DIR}/platform/containers/data/config/.env ]]; then
            MGMT_PASSWORD=$(. ${WORK_DIR}/platform/containers/data/config/.env && echo $MGMT_PASSWORD)
        fi
    fi

    if [[ ! -z $SRS_PLATFORM_SECRET ]]; then break; fi

    echo "Warning: Retry to get SRS_PLATFORM_SECRET."
    sleep 3
done

if [[ ! -z $SRS_PLATFORM_SECRET ]]; then
    echo "SRS_PLATFORM_SECRET=$SRS_PLATFORM_SECRET"
    echo "MGMT_PASSWORD=$MGMT_PASSWORD"
    if [[ ! -z  $OUTPUT ]]; then
        echo "SRS_PLATFORM_SECRET=$SRS_PLATFORM_SECRET" > $OUTPUT
        echo "MGMT_PASSWORD=$MGMT_PASSWORD" >> $OUTPUT
        # For local development OpenAI service test. Ignored by production, because there is no source .env file.
        if [[ -f ~/git/issues-translation/.env ]]; then cat ~/git/issues-translation/.env |grep OPENAI >> $OUTPUT; fi
    fi
    exit 0
fi

echo "Get SRS_PLATFORM_SECRET failed"
exit 1
