#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/.. && pwd)
cd $WORK_DIR

# Please update the release version for each major version.
TAG=publication-v4.6.16
echo "Publication TAG=$TAG, WORK_DIR=$WORK_DIR"

RELEASE=$(git describe --tags --abbrev=0 --match publication-*)
if [[ $TAG == $RELEASE ]]; then
  echo "Failed: Release $TAG already published."
  echo "Please update the TAG in $0 then run again.";
  exit 1
fi

VERSION=$(echo $TAG| sed 's/publication-v//g')
cat scripts/setup-aapanel/info.json |sed "s|\"versions\": .*|\"versions\": \"$VERSION\",|g" > tmp.json && mv tmp.json scripts/setup-aapanel/info.json &&
cat scripts/setup-bt/info.json |sed "s|\"versions\": .*|\"versions\": \"$VERSION\",|g" > tmp.json && mv tmp.json scripts/setup-bt/info.json

git st |grep -q 'nothing to commit'
if [[ $? -ne 0 ]]; then
  echo "Failed: Please commit before release";
  exit 1
fi

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
