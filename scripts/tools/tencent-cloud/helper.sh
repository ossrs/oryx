#!/bin/bash

# Execute by: bash xxx.sh or bash zzz/yyy/xxx.sh or ./xxx.sh or ./zzz/yyy/xxx.sh source xxx.sh
REALPATH=$(realpath ${BASH_SOURCE[0]})
SCRIPT_DIR=$(cd $(dirname ${REALPATH}) && pwd)
WORK_DIR=$(cd $(dirname ${REALPATH}) && pwd)
cd ${WORK_DIR}

if [[ ! -d venv ]]; then
    python3 -m venv venv
    source venv/bin/activate
    pip install -r requirements.txt
else
    source venv/bin/activate
fi
ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Failed to activate venv."; exit ${ret}; fi

echo "Execute: python $@"
python $@
