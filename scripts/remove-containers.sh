#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/.. && pwd)
cd $WORK_DIR

CONTAINERS=$(cd mgmt/containers/names && ls)
for NAME in $CONTAINERS; do
  if [[ $(docker ps -f name=$NAME --format '{{json .}}' |wc -l) -eq 1 ]]; then
    echo "Remove container $NAME" && docker rm -f $NAME
  fi
done

