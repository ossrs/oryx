#!/bin/bash

# Ignore darwin
if [[ $(uname -s) == 'Darwin' ]]; then
  echo "Mac is not supported"; exit 1;
fi

# The main directory.
SRS_HOME=/usr/local/srs-cloud
DEPLOY_HOME=$(dirname ${SRS_HOME})

echo "Install depends"
apt-get install -y git gcc g++ gdb make tree dstat docker docker.io nginx curl net-tools &&
apt-get -qqy clean
if [[ $? -ne 0 ]]; then echo "Install depends failed"; exit 1; fi

echo "Create nginx directory"
mkdir -p /etc/nginx/default.d /etc/nginx/conf.d
if [[ $? -ne 0 ]]; then echo "Copy srs-cloud failed"; exit 1; fi

echo "Install srs-cloud"
mkdir -p ${DEPLOY_HOME} && cd ${DEPLOY_HOME}
if [[ $? -ne 0 ]]; then echo "Copy srs-cloud failed"; exit 1; fi

# When droplet created, it might fail as:
#   gnutls_handshake() failed: The TLS connection was non-properly terminated.
# so we try to wait for a while and try later.
for ((i=0; i<30; i++)); do
  rm -rf srs-cloud && git clone -b main https://github.com/ossrs/srs-cloud.git && GIT_DONE=YES
  if [[ $? -eq 0 ]]; then break; fi
  echo "Ignore error and try later..."; sleep 3;
done
if [[ $GIT_DONE != YES ]]; then
  echo "Clone srs-cloud failed"; exit 1;
fi

cd ${DEPLOY_HOME} && rm -rf srs-terraform && ln -sf srs-cloud srs-terraform
if [[ $? -ne 0 ]]; then echo "Link srs-cloud failed"; exit 1; fi

