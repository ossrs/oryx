#!/bin/bash

# Execute by: bash xxx.sh or bash zzz/yyy/xxx.sh or ./xxx.sh or ./zzz/yyy/xxx.sh source xxx.sh
REALPATH=$(realpath ${BASH_SOURCE[0]})
SCRIPT_DIR=$(cd $(dirname ${REALPATH}) && pwd)
WORK_DIR=$(cd $(dirname ${REALPATH})/../.. && pwd)
cd ${WORK_DIR}

# Start by script.
SRS_PLATFORM_SECRET=$(docker exec srs-cloud redis-cli hget SRS_PLATFORM_SECRET token 2>/dev/null)
MGMT_PASSWORD=$(docker exec srs-cloud bash -c '. /data/config/.env && echo $MGMT_PASSWORD' 2>/dev/null)

# Start by BT or aaPanel.
if [[ -z $SRS_PLATFORM_SECRET ]]; then
    SRS_PLATFORM_SECRET=$(docker exec bt docker exec srs-cloud redis-cli hget SRS_PLATFORM_SECRET token 2>/dev/null)
    MGMT_PASSWORD=$(docker exec bt docker exec srs-cloud bash -c '. /data/config/.env && echo $MGMT_PASSWORD')
fi

# Start by develop.
if [[ -z $SRS_PLATFORM_SECRET ]]; then
    SRS_PLATFORM_SECRET=$(redis-cli hget SRS_PLATFORM_SECRET token 2>/dev/null)
    MGMT_PASSWORD=$(. ${WORK_DIR}/platform/containers/data/config/.env && echo $MGMT_PASSWORD)
fi

if [[ ! -z $SRS_PLATFORM_SECRET ]]; then
    echo "SRS_PLATFORM_SECRET=$SRS_PLATFORM_SECRET"
    echo "MGMT_PASSWORD=$MGMT_PASSWORD"
    exit 0
fi

echo "Get SRS_PLATFORM_SECRET failed"
exit 1
