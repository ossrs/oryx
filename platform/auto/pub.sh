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

RELEASE=$(git describe --tags --abbrev=0 --match platform-*)
REVISION=$(echo $RELEASE|awk -F . '{print $3}')
let NEXT=$REVISION+1
echo "Last release is $RELEASE, revision is $REVISION, next is $NEXT"

VERSION="1.0.$NEXT"
TAG="platform-v$VERSION"
echo "publish version $VERSION as tag $TAG"

cat package.json |sed "s|\"version\":.*|\"version\":\"$VERSION\",|g" > tmp.json && mv tmp.json package.json &&
git ci -am "Update version to $TAG"

git push
git tag -d $TAG 2>/dev/null && git push origin :$TAG
git tag $TAG
git push origin $TAG

git remote |grep -q gitee && git push gitee && git push gitee $TAG
git remote |grep -q cloud && git push cloud && git push cloud $TAG

echo "publish $TAG ok"
echo "    https://github.com/ossrs/srs-cloud/actions?query=branch%3A$TAG"

echo "now, update the releases"
bash ../releases/auto/pub.sh

