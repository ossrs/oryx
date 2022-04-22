#!/bin/bash

echo "bash setup.sh $*"

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH) && pwd)
echo "Run setup at $WORK_DIR from $0"
cd $WORK_DIR

R0_FILE=$1
if [[ -z $R0_FILE ]]; then echo "No r0 file"; exit 1; fi

bash do_setup.sh $* && echo "Execute OK"
r0=$?; if [[ $r0 -ne 0 ]]; then
  # Rollback the install.
  systemctl disable srs-cloud >/dev/null 2>&1
  rm -f /usr/lib/systemd/system/srs-cloud.service

  echo "Execute failed"
  echo "R0=$r0" > $R0_FILE
  echo "TS=$(date +%s)" >> $R0_FILE
  echo "DATE=\"$(date)\"" >> $R0_FILE
  echo "ARGS=\"bash $0 $*\"" >> $R0_FILE
  echo "DESC=\"srs-cloud install script\"" >> $R0_FILE
  exit 1
fi

