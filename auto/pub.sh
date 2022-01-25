#!/bin/bash

git st |grep -q 'nothing to commit'
if [[ $? -ne 0 ]]; then
  echo "Failed: Please commit before release";
  exit 1
fi

RELEASE=$(git describe --tags --abbrev=0 --exclude release-*)
REVISION=$(echo $RELEASE|awk -F . '{print $3}')
let NEXT=$REVISION+1
echo "Last release is $RELEASE, revision is $REVISION, next is $NEXT"

TAG="v1.0.$NEXT"
echo "publish $TAG"

cat mgmt/package.json |sed "s|\"version\":.*|\"version\":\"$TAG\",|g" > t.json
mv t.json mgmt/package.json
git ci -am "Release $TAG"

git push
git tag -d $TAG 2>/dev/null && git push origin :$TAG
git tag $TAG
git push origin $TAG
echo "publish $TAG ok"
echo "    https://github.com/ossrs/srs-terraform/actions"
