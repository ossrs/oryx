#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/.. && pwd)
cd $WORK_DIR/mgmt

cd containers/names &&
docker rm -f *

