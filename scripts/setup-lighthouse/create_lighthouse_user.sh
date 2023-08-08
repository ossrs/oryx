#!/bin/bash

username=lighthouse

useradd -m -s /bin/bash $username
echo "$username ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

