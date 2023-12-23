#!/bin/bash

sudo apt-get update -y &&
sudo apt-get install -y curl docker.io
if [[ $? -ne 0 ]]; then echo "Install curl failed."; exit 1; fi

rm -f linux-srs_stack-en.tar.gz &&
curl -O -fsSL https://github.com/ossrs/srs-stack/releases/latest/download/linux-srs_stack-en.tar.gz
if [[ $? -ne 0 ]]; then echo "Download srs-stack failed."; exit 1; fi

tar xf linux-srs_stack-en.tar.gz
if [[ $? -ne 0 ]]; then echo "Unpack srs-stack failed."; exit 1; fi

sed -i 's|MGMT_PORT=2022|MGMT_PORT=80|g' srs_stack/mgmt/bootstrap &&
sed -i 's|HTTPS_PORT=2443|HTTPS_PORT=443|g' srs_stack/mgmt/bootstrap
if [[ $? -ne 0 ]]; then echo "Update bootstrap failed."; exit 1; fi

sudo bash srs_stack/scripts/setup-ubuntu/install.sh
if [[ $? -ne 0 ]]; then echo "Install srs-stack failed."; exit 1; fi
