#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/.. && pwd)
cd $WORK_DIR

REGISTRY=ccr.ccs.tencentyun.com
echo "Update containers, REGISTRY: $REGISTRY"

bash mgmt/auto/upgrade_containers $REGISTRY
if [[ $? -ne 0 ]]; then echo "Update docker containers failed"; exit 1; fi

