#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/.. && pwd)
cd $WORK_DIR

# Please update the release version for each major version.
TAG=publication-v4.6.12
echo "Publication TAG=$TAG, WORK_DIR=$WORK_DIR"

git tag -d $TAG 2>/dev/null
git push origin :$TAG 2>/dev/null
git push gitee :$TAG 2>/dev/null
git tag $TAG
git push origin $TAG
git push gitee $TAG
echo "Publication ok, please visit"
echo "    Please test it after https://github.com/ossrs/srs-cloud/actions/workflows/publication.yml done"
echo "    Download bt-srs_cloud.zip from https://github.com/ossrs/srs-cloud/releases"
echo "    Then submit it to https://www.bt.cn/developer/details.html?id=600801805"
echo "    Finally, update release at https://gitee.com/ossrs/srs-cloud/releases/new"
