#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH) && pwd)
echo "Run install at $WORK_DIR from $0"
cd $WORK_DIR

bash do_install.sh $* &&
echo "Execute OK" && sleep 3
if [[ $? -ne 0 ]]; then echo "Execute failed"; sleep 30; exit 1; fi

