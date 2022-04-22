#!/bin/bash

TMP_DIR=/tmp/zip-for-bt-srs_cloud
PLUGIN=srs_cloud
echo "Zip at TMP_DIR=$TMP_DIR, PLUGIN=$PLUGIN"

rm -rf $TMP_DIR/$PLUGIN && mkdir -p $TMP_DIR/$PLUGIN &&
echo "Zip TMP_DIR=$TMP_DIR"
if [[ $? -ne 0 ]]; then echo "Setup temporary directory failed"; exit 1; fi

mkdir -p $TMP_DIR/source && cd $TMP_DIR/source &&
if [[ ! -d srs-cloud ]]; then git clone https://gitee.com/ossrs/srs-cloud.git; fi &&
cd $TMP_DIR/source/srs-cloud && git reset --hard HEAD~10 >/dev/null && git pull | grep files &&
git branch -vv |grep '*' &&
echo "Cache at $TMP_DIR/source/srs-cloud"
if [[ $? -ne 0 ]]; then echo "Cache source failed"; exit 1; fi

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

if [[ -f ~/Downloads/$PLUGIN.zip ]]; then
  mv $TMP_DIR/$PLUGIN.zip ~/Downloads/$PLUGIN.zip
  echo "Plugin OK:"
  echo "    Zip moved to ~/Downloads/$PLUGIN.zip"
else
  echo "You could move it by:"
  echo "    mv $TMP_DIR/$PLUGIN.zip ~/Downloads/$PLUGIN.zip"
fi

