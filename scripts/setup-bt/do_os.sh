#!/bin/bash

# Check OS first, only support CentOS or Ubuntu
# Modification from mxdyeah.
yum --version >/dev/null 2>&1 && rpm --version >/dev/null 2>&1 && OS_NAME='CentOS'
apt-get --version >/dev/null 2>&1 && OS_NAME='Ubuntu' && OS_NAME='Debian'
if [[ -z $OS_NAME ]]; then echo "Only support CentOS/Ubuntu/Debian"; exit 1; fi

if [[ $OS_NAME == 'CentOS' ]]; then
  # Check CentOS version.
  CentOS_VERSION=$(rpm --eval '%{centos_ver}')
  if [[ $CentOS_VERSION -lt 7 ]]; then echo "Only support CentOS 7+, yours is $CentOS_VERSION"; exit 1; fi
fi

if [[ $OS_NAME == 'Ubuntu' ]]; then
  # Check Ubuntu version.
  Ubuntu_VERSION=$(cat /etc/os-release |grep VERSION_ID |awk -F '"' '{print $2}' |awk -F '.' '{print $1}')
  if [[ $Ubuntu_VERSION -lt 18 ]]; then echo "Only support Ubuntu 18+, yours is $Ubuntu_VERSION"; exit 1; fi
  echo "OS is $OS_NAME $Ubuntu_VERSION"
fi

if [[ $OS_NAME == 'Debian' ]]; then
  # Check Debian version.
  Debian_VERSION=$(cat /etc/os-release |grep VERSION_ID |awk -F '"' '{print $2}' |awk -F '.' '{print $1}')
  if [[ $Debian_VERSION -lt 11 ]]; then echo "Only support Debian 11+, yours is $Debian_VERSION"; exit 1; fi
  echo "OS is $OS_NAME $Debian_VERSION"
fi

