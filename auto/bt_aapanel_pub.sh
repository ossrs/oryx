#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/.. && pwd)
cd $WORK_DIR

# Please update the release version for each major version.
TAG=publication-v4.6.7
echo "Publication TAG=$TAG, WORK_DIR=$WORK_DIR"

git tag -d $TAG 2>/dev/null
git push origin :$TAG 2>/dev/null
git tag $TAG
git push origin $TAG
echo "Publication ok, please visit"
echo "    Please test it after https://github.com/ossrs/srs-cloud/actions/workflows/publication.yml done"
echo "    Download bt-srs_cloud.zip from https://github.com/ossrs/srs-cloud/releases"
echo "    Then submit it to https://www.bt.cn/developer/details.html?id=600801805"

