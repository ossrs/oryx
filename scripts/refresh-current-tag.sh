#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/.. && pwd)
echo "Run pub at $WORK_DIR from $0"
cd $WORK_DIR

git st |grep -q 'nothing to commit'
if [[ $? -ne 0 ]]; then
  echo "Failed: Please commit before release";
  exit 1
fi

######################################################################
VERSION=$(cat mgmt/version.go|grep 'const version'|awk '{print $4}'| sed 's/"//g') &&
TAG="mgmt-$VERSION" && echo "Refresh $TAG"
if [[ $? -ne 0 ]]; then echo "Release failed"; exit 1; fi

git tag -d $TAG 2>/dev/null && git push origin :$TAG && git push gitee :$TAG
git tag $TAG; git push origin $TAG; git push gitee $TAG

echo "publish $TAG ok"

# Release v1.0.xxx for mgmt.
TAG=$VERSION && echo "Refresh $TAG"
git tag -d $TAG 2>/dev/null && git push origin :$TAG && git push gitee :$TAG
git tag $TAG; git push origin $TAG; git push gitee $TAG

echo "publish $TAG ok"

######################################################################
VERSION=$(cat platform/version.go|grep 'const version'|awk '{print $4}'| sed 's/"//g') &&
TAG="platform-$VERSION" && echo "Refresh $TAG"
if [[ $? -ne 0 ]]; then echo "Release failed"; exit 1; fi

git tag -d $TAG 2>/dev/null && git push origin :$TAG && git push gitee :$TAG
git tag $TAG; git push origin $TAG; git push gitee $TAG

echo "publish $TAG ok"

echo "    Please test it after https://github.com/ossrs/srs-cloud/actions/workflows/mgmt.yml done"
echo "    Please test it after https://github.com/ossrs/srs-cloud/actions/workflows/platform.yml done"

