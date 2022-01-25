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

RELEASE=$(git describe --tags --abbrev=0 --exclude release-*)
REVISION=$(echo $RELEASE|awk -F . '{print $3}')
let NEXT=$REVISION+1
echo "Last release is $RELEASE, revision is $REVISION, next is $NEXT"

VERSION="1.0.$NEXT"
TAG="v$VERSION"
echo "publish version $VERSION as tag $TAG"

cat mgmt/package.json |sed "s|\"version\":.*|\"version\":\"$VERSION\",|g" > tmp.json && mv tmp.json mgmt/package.json &&
cat releases/package.json |sed "s|\"version\":.*|\"version\":\"$VERSION\",|g" > tmp.json && mv tmp.json releases/package.json &&
cat releases/releases.js |sed "s|const\ latest\ =.*|const latest = '$TAG';|g" > tmp.js && mv tmp.js releases/releases.js
git ci -am "Release $VERSION as $TAG"

git push
git tag -d $TAG 2>/dev/null && git push origin :$TAG
git tag $TAG
git push origin $TAG

git remote |grep -q gitee &&
git push gitee &&
git push gitee $TAG

echo "publish $TAG ok"
echo "    https://github.com/ossrs/srs-terraform/actions"
