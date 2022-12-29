#!/bin/bash

# Install srs-cloud, for example:
#   bash /www/server/panel/plugin/srs_cloud/setup.sh /tmp/srs_cloud_install.r0 /www/server/nodejs/v16.9.0/bin/node /www/server/nginx/logs/nginx.pid /www/wwwroot srs.cloud.local
# If ok, we will create systemctl service at:
#   /usr/lib/systemd/system/srs-cloud.service

R0_FILE=$1; NODEJS=$2; NGINX_PID=$3; WWW_HOME=$4; SITE_NAME=$5
if [[ -z $NODEJS ]]; then echo "The nodejs is empty"; exit 1; fi
if [[ ! -f $NODEJS ]]; then echo "No nodejs installed"; exit 1; fi
if [[ -z $NGINX_PID ]]; then echo "The nginx is empty"; exit 1; fi
if [[ ! -f $NGINX_PID ]]; then echo "No nginx pid at $NGINX_PID"; exit 1; fi
if [[ ! -f /etc/init.d/nginx ]]; then echo "No nginx file at /etc/init.d/nginx"; exit 1; fi
if [[ -z $WWW_HOME ]]; then echo "No www home"; exit 1; fi
if [[ -z $SITE_NAME ]]; then echo "No site name"; exit 1; fi

# Setup the path.
install_path=/www/server/panel/plugin/srs_cloud
DEPLOY_HOME=/usr/local/lighthouse/softwares
SRS_HOME=${DEPLOY_HOME}/srs-cloud
INSTALL_HOME=/usr/local/srs-cloud
echo "Setup SRS at install_path=$install_path, SRS_HOME=$SRS_HOME, INSTALL_HOME=$INSTALL_HOME, NODEJS=$NODEJS, NGINX_PID=$NGINX_PID, WWW_HOME=$WWW_HOME, SITE_NAME=$SITE_NAME"

source do_os.sh
if [[ $? -ne 0 ]]; then echo "Setup OS failed"; exit 1; fi

########################################################################################################################
# Setup the PATH for nodejs.
export PATH=$PATH:$(dirname $NODEJS)

# Generate PATH for node.
mkdir -p $SRS_HOME/mgmt/containers/bin &&
cat << END > $SRS_HOME/mgmt/containers/bin/bootstrap
#!/bin/bash
NODEJS=$(dirname $NODEJS)
export PATH=\$PATH:\$NODEJS

NGINX_PID=$NGINX_PID
export NGINX_PID=\$NGINX_PID
END
if [[ $? -ne 0 ]]; then echo "Setup bootstrap failed"; exit 1; fi

# For BT, we use special env, to disable discover of platform.
cat << END > ${SRS_HOME}/mgmt/.env
CLOUD=AAPANEL
REACT_APP_LOCALE=en
END
if [[ $? -ne 0 ]]; then echo "Setup .env failed"; exit 1; fi

# Setup extra env.
mkdir -p $SRS_HOME/mgmt/containers/bin &&
cat << END > $SRS_HOME/mgmt/containers/bin/.env
# Please use BT to configure the domain and HTTPS.
SRS_HTTPS=off
END
if [[ $? -ne 0 ]]; then echo "Setup extra env failed"; exit 1; fi

########################################################################################################################
# User should install nodejs, because we can't do it.
(cd $SRS_HOME/scripts/check-node-version && npm install && node .)
if [[ $? -ne 0 ]]; then echo "Please install node from https://nodejs.org"; exit 1; fi

cd ${SRS_HOME} && make install
if [[ $? -ne 0 ]]; then echo "Copy srs-cloud failed"; exit 1; fi

cd $DEPLOY_HOME && rm -rf srs-terraform && ln -sf srs-cloud srs-terraform
if [[ $? -ne 0 ]]; then echo "Link srs-cloud failed"; exit 1; fi

# Setup git alias to make it convenient.
cd ${SRS_HOME}/mgmt &&
echo "Setup git alias to make it more convenient" &&
git config --local alias.co checkout &&
git config --local alias.br branch &&
git config --local alias.ci commit &&
git config --local alias.st status
if [[ $? -ne 0 ]]; then echo "Setup git alias failed"; exit 1; fi

########################################################################################################################
# Update the docker images.
echo "Cache docker images" &&
docker pull docker.io/ossrs/srs:4 &&
docker pull docker.io/ossrs/node:slim &&
docker pull docker.io/ossrs/srs-cloud:hooks-1 &&
docker pull docker.io/ossrs/srs-cloud:tencent-1 &&
docker pull docker.io/ossrs/srs-cloud:ffmpeg-1 &&
docker pull docker.io/ossrs/srs-cloud:platform-1 &&
docker pull docker.io/ossrs/prometheus &&
docker pull docker.io/ossrs/redis_exporter &&
docker pull docker.io/ossrs/node-exporter &&
docker pull docker.io/ossrs/certbot &&
docker pull docker.io/ossrs/redis
if [[ $? -ne 0 ]]; then echo "Cache docker images failed"; exit 1; fi

# If install ok, the directory should exists.
if [[ ! -d ${INSTALL_HOME} || ! -d ${INSTALL_HOME}/mgmt ]]; then
  echo "Install srs-cloud failed"; exit 1;
fi

# Link the www root to container.
WWW_FILES=$(ls ${SRS_HOME}/mgmt/containers/www)
for file in $WWW_FILES; do
  rm -rf $WWW_HOME/$SITE_NAME/$file &&
  ln -sf ${SRS_HOME}/mgmt/containers/www/$file $WWW_HOME/$SITE_NAME/$file
done
if [[ $? -ne 0 ]]; then echo "Link www root failed"; exit 1; fi

# Execute script for each run.
cd ${SRS_HOME}/mgmt && bash auto/foreach_run
if [[ $? -ne 0 ]]; then echo "Execute for each run script failed"; exit 1; fi

# Create init.d script.
rm -f /etc/init.d/srs_cloud &&
cp $install_path/init.d.sh /etc/init.d/srs_cloud &&
chmod +x /etc/init.d/srs_cloud
if [[ $? -ne 0 ]]; then echo "Setup init.d script failed"; exit 1; fi

# Create srs-cloud service, and the credential file.
# Remark: Never start the service, because the IP will change for new machine created.
cd ${INSTALL_HOME} &&
cp -f usr/lib/systemd/system/srs-cloud.service /usr/lib/systemd/system/srs-cloud.service &&
systemctl daemon-reload && systemctl enable srs-cloud
if [[ $? -ne 0 ]]; then echo "Install srs-cloud failed"; exit 1; fi

/etc/init.d/srs_cloud restart srs-cloud
if [[ $? -ne 0 ]]; then echo "Start srs-cloud failed"; exit 1; fi

