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

RELEASE=$(git describe --tags --abbrev=0 --match releases-*)
REVISION=$(echo $RELEASE|awk -F . '{print $3}')
let NEXT=$REVISION+1
echo "Last release is $RELEASE, revision is $REVISION, next is $NEXT"

VERSION="1.0.$NEXT"
TAG="releases-v$VERSION"
echo "publish version $VERSION as tag $TAG"

cat package.json |sed "s|\"version\":.*|\"version\":\"$VERSION\",|g" > tmp.json && mv tmp.json package.json &&
echo "Update version ok"
if [[ $? -ne 0 ]]; then echo "Release failed"; exit 1; fi

