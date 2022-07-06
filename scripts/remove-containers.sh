#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/.. && pwd)
cd $WORK_DIR

CONTAINERS=$(cd mgmt/containers/names && ls | grep -v redis)
for NAME in $CONTAINERS; do
  if [[ $(docker ps -f name=$NAME --format '{{json .}}' |wc -l) -eq 1 ]]; then
    if [[ $NAME == 'redis' ]]; then
      # For some container, we must wait for it to save data to disk.
      echo "Stop container $NAME" && docker stop -t 30 $NAME &&
      if [[ $(docker ps -f name=$NAME --format '{{json .}}' |wc -l) -eq 1 ]]; then
        echo "Remove container $NAME" && docker rm -f $NAME
      fi
    else
      echo "Remove container $NAME" && docker rm -f $NAME
    fi
  fi
done

