#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/.. && pwd)
cd $WORK_DIR

# Please update the release version for each major version.
TAG=publication-v4.6
echo "Publication TAG=$TAG, WORK_DIR=$WORK_DIR"

git tag -d $TAG 2>/dev/null
git push origin :$TAG 2>/dev/null
git tag $TAG
git push origin $TAG
echo "Publication ok, please visit"
echo "    https://github.com/ossrs/srs-cloud/releases"

