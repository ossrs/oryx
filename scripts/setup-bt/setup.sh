#!/bin/bash

echo "bash setup.sh $*"

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH) && pwd)
echo "Run setup at $WORK_DIR from $0"
cd $WORK_DIR

ARGS=$*

HELP=no
R0_FILE=

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -h|--help) HELP=yes; shift ;;
        --r0) R0_FILE=$2; shift 2 ;;
        *) shift 2 ;;
    esac
done

if [[ "$HELP" == yes ]]; then
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -h, --help    Show this help message and exit"
    echo "  --r0          The install error file. Default: ${R0_FILE}"
    exit 0
fi

if [[ -z $R0_FILE ]]; then echo "No r0 file"; exit 1; fi

bash do_setup.sh $ARGS && echo "Execute OK"
r0=$?; if [[ $r0 -ne 0 ]]; then
  # Rollback the install.
  systemctl disable oryx >/dev/null 2>&1
  rm -f /usr/lib/systemd/system/oryx.service

  echo "Execute failed"
  echo "R0=$r0" > $R0_FILE
  echo "TS=$(date +%s)" >> $R0_FILE
  echo "DATE=\"$(date)\"" >> $R0_FILE
  echo "ARGS=\"bash $0 $ARGS\"" >> $R0_FILE
  echo "DESC=\"oryx install script\"" >> $R0_FILE
  exit 1
fi

