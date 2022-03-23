#!/usr/bin/env bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/.. && pwd)
cd $WORK_DIR

find . -name *.js |grep -v node_modules |grep -v sdk |grep -v containers |grep -v build |xargs wc -l

