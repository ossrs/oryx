#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/.. && pwd)
cd $WORK_DIR

TAG=publication-v4.4
echo "Publication TAG=$TAG, WORK_DIR=$WORK_DIR"

git tag -d $TAG
git push origin :$TAG
git tag $TAG
git push origin $TAG
echo "Publication ok, please visit"
echo "    https://github.com/ossrs/srs-cloud/releases"

