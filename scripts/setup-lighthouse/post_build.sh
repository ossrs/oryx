#!/bin/bash

if [ -f /etc/centos-release ]
then
    passwd -d root
elif [ -f /etc/lsb-release ]
then
    passwd -d ubuntu
else
    echo "failed to delete password"
    exit 1
fi

