#!/bin/bash

REALPATH=$(realpath $0)
WORK_DIR=$(cd $(dirname $REALPATH)/../../.. && pwd)

TMP_DIR=/tmp/zip-for-bt-srs_cloud
PLUGIN=srs_cloud
ZIP_FILE=bt-$PLUGIN.zip
echo "Zip at TMP_DIR=$TMP_DIR, PLUGIN=$PLUGIN, ZIP_FILE=$ZIP_FILE, WORK_DIR=$WORK_DIR"

rm -rf $TMP_DIR/$PLUGIN && mkdir -p $TMP_DIR/$PLUGIN &&
echo "Zip TMP_DIR=$TMP_DIR"
if [[ $? -ne 0 ]]; then echo "Setup temporary directory failed"; exit 1; fi

mkdir -p $TMP_DIR/source && cd $TMP_DIR/source &&
if [[ ! -d srs-cloud ]]; then
  echo "Clone srs-cloud to $TMP_DIR/source/srs-cloud"
  if [[ -z $GITHUB_ACTIONS ]]; then
    git clone https://gitee.com/ossrs/srs-cloud.git;
  else
    git clone --depth=1 https://github.com/ossrs/srs-cloud.git
  fi
fi
if [[ $? -ne 0 ]]; then echo "Cache source failed"; exit 1; fi

# For development, not GitHub Actions, set to gitee and refresh the code.
if [[ -z $GITHUB_ACTIONS ]]; then
  cd $TMP_DIR/source/srs-cloud && git remote set-url origin https://gitee.com/ossrs/srs-cloud.git &&
  cd $TMP_DIR/source/srs-cloud && git reset --hard HEAD~10 >/dev/null && git pull | grep files &&
  git branch -vv |grep '*' &&
  echo "Cache at $TMP_DIR/source/srs-cloud"
  if [[ $? -ne 0 ]]; then echo "Cache source failed"; exit 1; fi
fi

# For aaPanel, reset the repository to GitHub.
cd $TMP_DIR/source/srs-cloud && git remote set-url origin https://gitee.com/ossrs/srs-cloud.git
if [[ $? -ne 0 ]]; then echo "Setup source remote failed"; exit 1; fi

mkdir -p $TMP_DIR/$PLUGIN/srs-cloud && cd $TMP_DIR/$PLUGIN/srs-cloud &&
ln -sf $TMP_DIR/source/srs-cloud source &&
echo "Source at $TMP_DIR/$PLUGIN/srs-cloud/source"
if [[ $? -ne 0 ]]; then echo "Setup source directory failed"; exit 1; fi

cp -r source/.git . &&
echo "Copy code to $TMP_DIR/$PLUGIN/srs-cloud"
if [[ $? -ne 0 ]]; then echo "Copy data failed"; exit 1; fi

cd $TMP_DIR/$PLUGIN && cp -r srs-cloud/source/scripts/setup-bt/* . &&
echo "Copy plugin to $TMP_DIR/$PLUGIN"
if [[ $? -ne 0 ]]; then echo "Copy plugin failed"; exit 1; fi

rm -f srs-cloud/source && echo "Remove source files"
if [[ $? -ne 0 ]]; then echo "Remove source failed"; exit 1; fi

cd $TMP_DIR/ &&
zip -q -r $PLUGIN.zip $PLUGIN &&
echo "Zip generated at $TMP_DIR/$PLUGIN.zip"
if [[ $? -ne 0 ]]; then echo "Zip plugin failed"; exit 1; fi

mv $TMP_DIR/$PLUGIN.zip $TMP_DIR/$ZIP_FILE &&
echo "Rename to $TMP_DIR/$ZIP_FILE"

if [[ -f ~/Downloads/$ZIP_FILE ]]; then
  mv $TMP_DIR/$ZIP_FILE ~/Downloads/$ZIP_FILE
  echo "Plugin OK:"
  echo "    Zip moved to ~/Downloads/$ZIP_FILE"
else
  echo "You could move it by:"
  echo "    mv $TMP_DIR/$ZIP_FILE ~/Downloads/$ZIP_FILE"
fi

