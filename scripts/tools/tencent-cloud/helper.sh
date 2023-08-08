#!/bin/bash

# Execute by: bash xxx.sh or bash zzz/yyy/xxx.sh or ./xxx.sh or ./zzz/yyy/xxx.sh source xxx.sh
REALPATH=$(realpath ${BASH_SOURCE[0]})
SCRIPT_DIR=$(cd $(dirname ${REALPATH}) && pwd)
WORK_DIR=$(cd $(dirname ${REALPATH}) && pwd)
cd ${WORK_DIR}

source venv/bin/activate
ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Failed to activate venv."; exit ${ret}; fi

echo "Execute: python $@"
python $@
