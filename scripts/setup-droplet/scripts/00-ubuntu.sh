#!/bin/bash

# Ignore darwin
if [[ $(uname -s) == 'Darwin' ]]; then
  echo "Mac is not supported"; exit 1;
fi

echo "Install depends"
apt-get install -y git gcc g++ gdb make tree dstat docker docker.io redis nginx curl net-tools &&
apt-get -qqy clean
if [[ $? -ne 0 ]]; then echo "Install depends failed"; exit 1; fi

echo "Create nginx directory"
mkdir -p /etc/nginx/default.d
if [[ $? -ne 0 ]]; then echo "Copy srs-cloud failed"; exit 1; fi

echo "Install srs-cloud"
mkdir -p /usr/local/lighthouse/softwares && cd /usr/local/lighthouse/softwares
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

cd /usr/local/lighthouse/softwares && rm -rf srs-terraform && ln -sf srs-cloud srs-terraform
if [[ $? -ne 0 ]]; then echo "Link srs-cloud failed"; exit 1; fi

echo "Install nodejs 16"
curl -sL https://deb.nodesource.com/setup_16.x -o /tmp/nodesource_setup.sh && bash /tmp/nodesource_setup.sh &&
apt-get install -y nodejs
if [[ $? -ne 0 ]]; then echo "Install nodejs failed"; exit 1; fi

# User should install nodejs, because we can't do it.
cd /usr/local/lighthouse/softwares/srs-cloud &&
(cd scripts/check-node-version && npm install && node .)
if [[ $? -ne 0 ]]; then echo "Please install node from https://nodejs.org"; exit 1; fi

