#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/.. && pwd)
echo "Run pub at $WORK_DIR from $0"
cd $WORK_DIR

help=false
refresh=false

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -h|--help) help=true; shift ;;
        -refresh|--refresh) refresh=true; shift ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
done

if [ "$help" = true ]; then
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -h, --help           Show this help message and exit"
    echo "  -refresh, --refresh  Refresh current tag. Default: false"
    exit 0
fi

# We increase version from the platform-v* base.
RELEASE=$(git describe --tags --abbrev=0 --match platform-v*) &&
REVISION=$(echo $RELEASE|awk -F . '{print $3}')
if [[ $? -ne 0 ]]; then echo "Release failed"; exit 1; fi

let NEXT=$REVISION+1
if [[ $refresh == true && $REVISION != "-1" ]]; then
  let NEXT=$REVISION
fi
echo "Last release is $RELEASE, revision is $REVISION, next is $NEXT"

VERSION="1.0.$NEXT" &&
TAG="v$VERSION" &&
echo "publish version $VERSION as tag $TAG"
if [[ $? -ne 0 ]]; then echo "Release failed"; exit 1; fi

######################################################################
if [[ $(grep -q "const version = \"$TAG\"" mgmt/version.go || echo no) == no ]]; then
    echo "Failed: Please update mgmt/version.go to $TAG"
    exit 1
fi
if [[ $(grep -q "const version = \"$TAG\"" platform/version.go || echo no) == no ]]; then
    echo "Failed: Please update platform/version.go to $TAG"
    exit 1
fi
if [[ $(grep -q "const latest = \"$TAG\"" releases/version.go || echo no) == no ]]; then
    echo "Failed: Please update releases/version.go to $TAG"
    exit 1
fi

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

git tag $TAG && git push origin $TAG && git push gitee
echo "Publish OK: $TAG"

######################################################################
PLATFORM_TAG="platform-v$VERSION"

git tag -d $PLATFORM_TAG 2>/dev/null; git push origin :$PLATFORM_TAG 2>/dev/null; git push gitee :$PLATFORM_TAG 2>/dev/null
echo "Delete tag OK: $PLATFORM_TAG"

git tag $PLATFORM_TAG && git push origin $PLATFORM_TAG && git push gitee
echo "Publish OK: $PLATFORM_TAG"

echo -e "\n\n"
echo "Publish OK: $TAG $PLATFORM_TAG"
echo "    https://github.com/ossrs/srs-cloud/actions/workflows/platform.yml"
