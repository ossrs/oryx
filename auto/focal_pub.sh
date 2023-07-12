#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/.. && pwd)
echo "Run pub at $WORK_DIR from $0"
cd $WORK_DIR

help=false

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -h|--help) help=true; shift ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
done

if [ "$help" = true ]; then
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -h, --help           Show this help message and exit"
    exit 0
fi

RELEASE=$(git describe --tags --abbrev=0 --match focal-*)
REVISION=$(echo $RELEASE|awk -F . '{print $3}')
let NEXT=$REVISION+1
echo "Last release is $RELEASE, revision is $REVISION, next is $NEXT"

VERSION="1.0.$NEXT"
TAG="focal-v$VERSION"
echo "publish version $VERSION as tag $TAG"

######################################################################
git st |grep -q 'nothing to commit'
if [[ $? -ne 0 ]]; then
  echo "Failed: Please commit before release";
  exit 1
fi

git fetch origin
if [[ $(git status |grep -q 'Your branch is up to date' || echo 'no') == no ]]; then
  git status
  echo "Failed: Please sync before release";
  exit 1
fi
echo "Sync OK"

git fetch gitee
if [[ $(git diff origin/main gitee/main |grep -q diff && echo no) == no ]]; then
  git diff origin/main gitee/main |grep diff
  echo "Failed: Please sync gitee before release";
  exit 1
fi
echo "Sync gitee OK"

######################################################################
git tag -d $TAG 2>/dev/null; git push origin :$TAG 2>/dev/null; git push gitee :$TAG 2>/dev/null
echo "Delete tag OK: $TAG"

git tag $TAG && git push origin $TAG && git push gitee $TAG
echo "Publish OK: $TAG"

echo "publish $TAG ok"
echo "    Please test it after https://github.com/ossrs/srs-cloud/actions/workflows/focal.yml done"

