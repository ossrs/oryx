#!/bin/bash

PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

install_path=/www/server/panel/plugin/srs_cloud
DEPLOY_HOME=/usr/local/lighthouse/softwares
SRS_HOME=${DEPLOY_HOME}/srs-cloud

# Update sysctl.conf and add if not exists. For example:
#   update_sysctl net.ipv4.ip_forward 1 0 "# Controls IP packet forwarding"
function update_sysctl() {
    SYSCTL_KEY=$1 && SYSCTL_VALUE=$2 && SYSCTL_EMPTY_LINE=$3 && SYSCTL_COMMENTS=$4
    echo "Update with sysctl $SYSCTL_KEY=$SYSCTL_VALUE, empty-line=$SYSCTL_EMPTY_LINE, comment=$SYSCTL_COMMENTS"

    grep -q "^${SYSCTL_KEY}[ ]*=" /etc/sysctl.conf
    if [[ $? == 0 ]]; then
      sed -i "s/^${SYSCTL_KEY}[ ]*=.*$/${SYSCTL_KEY} = ${SYSCTL_VALUE}/g" /etc/sysctl.conf
    else
      if [[ $SYSCTL_EMPTY_LINE == 1 ]]; then echo '' >> /etc/sysctl.conf; fi &&
      if [[ $SYSCTL_COMMENTS != '' ]]; then echo "$SYSCTL_COMMENTS" >> /etc/sysctl.conf; fi &&
      echo "${SYSCTL_KEY} = ${SYSCTL_VALUE}" >> /etc/sysctl.conf
    fi
    if [[ $? -ne 0 ]]; then echo "Failed to sysctl $SYSCTL_KEY = $SYSCTL_VALUE $SYSCTL_COMMENTS"; exit 1; fi

    RESULT=$(grep "^${SYSCTL_KEY}[ ]*=" /etc/sysctl.conf)
    echo "Update done: ${RESULT}"
}

Install() {
  echo "Installing to $install_path, pwd: $(pwd)"

  source do_os.sh
  if [[ $? -ne 0 ]]; then echo "Setup OS failed"; exit 1; fi

  # Restore files from git.
  cd $install_path/srs-cloud && git reset --hard HEAD
  if [[ $? -ne 0 ]]; then echo "Reset files failed"; exit 1; fi

  # Change file permissions.
  find $install_path -type d -exec chmod 0755 {} \; &&
  find $install_path -type f -exec chmod 0644 {} \; &&
  cd $install_path/srs-cloud && chmod 755 mgmt/bootstrap mgmt/upgrade scripts/remove-containers.sh
  if [[ $? -ne 0 ]]; then echo "Change file permissions failed"; exit 1; fi

  # Restore files from git again, after changing file permisisons.
  cd $install_path/srs-cloud && git reset --hard HEAD
  if [[ $? -ne 0 ]]; then echo "Reset files failed"; exit 1; fi

  # We also process for git clone --depth=1, see https://stackoverflow.com/a/23987039/17679565
  GIT_DEPTH=$(git rev-list --all --count)
  if [[ $GIT_DEPTH -eq 1 ]]; then
    git pull --unshallow
    if [[ $? -ne 0 ]]; then echo "Git unshallow failed"; exit 1; fi
  fi

  # Move srs-cloud to its home.
  mkdir -p $DEPLOY_HOME
  if [[ -d $install_path/srs-cloud && ! -d $SRS_HOME/.git ]]; then
    rm -rf $SRS_HOME && mv $install_path/srs-cloud $SRS_HOME &&
    ln -sf $SRS_HOME $install_path/srs-cloud
    if [[ $? -ne 0 ]]; then echo "Create srs-cloud failed"; exit 1; fi
  fi

  # We must create the .env to avoid docker mountint as a dir.
  touch ${SRS_HOME}/mgmt/.env &&
  ln -sf ${SRS_HOME}/mgmt/.env ~/credentials.txt
  if [[ $? -ne 0 ]]; then echo "Create ${SRS_HOME}/mgmt/.env failed"; exit 1; fi

  # Allow network forwarding, required by docker.
  # See https://stackoverflow.com/a/41453306/17679565
  update_sysctl net.ipv4.ip_forward 1 1 "# Controls IP packet forwarding"

  # Setup the UDP buffer for WebRTC and SRT.
  # See https://www.jianshu.com/p/6d4a89359352
  update_sysctl net.core.rmem_max 16777216 1 "# For RTC/SRT over UDP"
  update_sysctl net.core.rmem_default 16777216
  update_sysctl net.core.wmem_max 16777216
  update_sysctl net.core.wmem_default 16777216

  # Now, we're ready to install by aaPanel.
  echo 'Wait for srs-cloud plugin ready...'; sleep 10;
  touch ${install_path}/.bt_ready

  echo 'Install OK'
}

Uninstall() {
  /etc/init.d/srs_cloud stop
  echo "Stop srs-cloud service ok"

  INIT_D=/etc/init.d/srs_cloud &&
  rm -f $INIT_D
  echo "Remove init.d script $INIT_D ok"

  systemctl disable srs-cloud
  rm -f /usr/lib/systemd/system/srs-cloud.service
  systemctl daemon-reload
  systemctl reset-failed
  echo "Remove srs-cloud.service ok"

  INSTALL_HOME=/usr/local/srs-cloud
  rm -rf $INSTALL_HOME
  echo "Remove install $INSTALL_HOME ok"

  SRS_ALIAS=/usr/local/lighthouse/softwares/srs-terraform
	rm -rf $SRS_HOME $SRS_ALIAS
	echo "Remove srs home $SRS_HOME ok"

  rm -rf $install_path
  echo "Remove plugin path $install_path ok"

  LOGS=$(ls /tmp/srs_cloud_install.*)
  rm -f $LOGS
  echo "Remove install flag files $LOGS ok"
}

if [ "${1}" == 'install' ];then
	Install
elif [ "${1}" == 'uninstall' ];then
	Uninstall
else
	echo 'Error!'; exit 1;
fi

