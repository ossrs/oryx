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

# We increase version from the platform-v* base.
RELEASE=$(git describe --tags --abbrev=0 --match platform-v*) &&
REVISION=$(echo $RELEASE|awk -F . '{print $3}') &&
let NEXT=$REVISION+1 &&
echo "Last release is $RELEASE, revision is $REVISION, next is $NEXT"
if [[ $? -ne 0 ]]; then echo "Release failed"; exit 1; fi

######################################################################
VERSION="1.0.$NEXT" &&
TAG="v$VERSION" &&
echo "publish version $VERSION as tag $TAG"
if [[ $? -ne 0 ]]; then echo "Release failed"; exit 1; fi

cat mgmt/version.go |sed "s|const\ version\ =.*|const version = \"v$VERSION\";|g" > tmp.go && mv tmp.go mgmt/version.go &&
cat platform/package.json |sed "s|\"version\":.*|\"version\":\"$VERSION\",|g" > tmp.json && mv tmp.json platform/package.json &&
cat releases/main.go |sed "s|const\ latest\ =.*|const latest = \"v$VERSION\";|g" > tmp.go && mv tmp.go releases/main.go &&
git ci -am "Update mgmt version to $TAG"
if [[ $? -ne 0 ]]; then echo "Release failed"; exit 1; fi

######################################################################
# Note that the mgmt should always use v1.2.3 without any prefix, to be compatible with previous upgrade script.
echo -e "\n\n"
git push
git tag -d $TAG 2>/dev/null && git push origin :$TAG
git tag $TAG; git push origin $TAG

git remote |grep -q gitee && git push gitee && git push gitee $TAG
git remote |grep -q cloud && git push cloud && git push cloud $TAG

echo "publish $TAG ok"

######################################################################
echo -e "\n\n"
TAG="platform-v$VERSION"
git tag -d $TAG 2>/dev/null && git push origin :$TAG
git tag $TAG; git push origin $TAG

git remote |grep -q gitee && git push gitee && git push gitee $TAG
git remote |grep -q cloud && git push cloud && git push cloud $TAG

######################################################################
echo -e "\n\n"
TAG="mgmt-v$VERSION"
git tag -d $TAG 2>/dev/null && git push origin :$TAG
git tag $TAG; git push origin $TAG

git remote |grep -q gitee && git push gitee && git push gitee $TAG
git remote |grep -q cloud && git push cloud && git push cloud $TAG

echo "publish $TAG ok"
echo "    Please test it after https://github.com/ossrs/srs-cloud/actions/workflows/platform.yml done"
echo "    Then release new latest version by: ./auto/releases_pub.sh"

