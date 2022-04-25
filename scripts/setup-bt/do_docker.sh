#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH) && pwd)
echo "Run install at $WORK_DIR from $0"
cd $WORK_DIR

source $WORK_DIR/do_os.sh
if [[ $? -ne 0 ]]; then echo "Setup OS failed"; exit 1; fi

if [[ $OS_NAME == 'Ubuntu' ]]; then
  # Install docker if not installed. Note that we must install docker here because the bt.soft.install can't install it,
  # and we could get around of it after installed manually.
  if [[ ! -f /usr/lib/systemd/system/docker.service ]];then
    apt-get install -y docker docker.io
    if [[ $? -ne 0 ]]; then echo "Install docker failed"; exit 1; fi
  fi
fi

