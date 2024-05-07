#!/bin/bash

SUDO=sudo
if [[ $EUID -eq 0 ]]; then
    SUDO=
fi
echo "Use SUDO=$SUDO"

if [[ $(curl --version 1>/dev/null 2>/dev/null && docker --version 1>/dev/null 2>/dev/null && echo yes) != yes ]]; then
    $SUDO apt-get update -y &&
    $SUDO apt-get install -y curl docker.io
    if [[ $? -ne 0 ]]; then echo "Error: Install curl failed."; exit 1; fi
fi
echo "Check curl and docker ok."

rm -f linux-oryx-en.tar.gz &&
curl -O -fsSL https://github.com/ossrs/oryx/releases/latest/download/linux-oryx-en.tar.gz
if [[ $? -ne 0 ]]; then echo "Error: Download oryx failed."; exit 1; fi
echo "Download Oryx ok."

tar xf linux-oryx-en.tar.gz
if [[ $? -ne 0 ]]; then echo "Error: Unpack oryx failed."; exit 1; fi
if [[ ! -f oryx/mgmt/bootstrap ]]; then
    echo "Error: No bootstrap found."; exit 1;
fi
echo "Unpack Oryx ok."

sed -i 's|MGMT_PORT=2022|MGMT_PORT=80|g' oryx/mgmt/bootstrap &&
sed -i 's|HTTPS_PORT=2443|HTTPS_PORT=443|g' oryx/mgmt/bootstrap
if [[ $? -ne 0 ]]; then echo "Error: Update bootstrap failed."; exit 1; fi
echo "Update bootstrap ok."

$SUDO bash oryx/scripts/setup-ubuntu/install.sh
if [[ $? -ne 0 ]]; then echo "Error: Install oryx failed."; exit 1; fi
echo "Install Oryx ok."
