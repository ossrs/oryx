#!/bin/bash

### BEGIN INIT INFO
# Provides:          ossrs(oryx)
# Required-Start:    $all
# Required-Stop:     $all
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: ossrs(oryx)
# Description:       https://github.com/ossrs/oryx
### END INIT INFO

SRS_HOME=/usr/local/oryx

start() {
  systemctl start oryx.service
}

stop() {
  systemctl stop oryx.service
}

status() {
  systemctl status oryx.service
}

menu() {
    case "$1" in
        start)
            start
            ;;
        stop)
            stop
            ;;
        restart)
            stop
            start
            ;;
        status)
            status
            ;;
        *)
            echo "Usage: $0 {start|stop|status|restart}"
            return 1
            ;;
    esac
}

menu $1

code=$?
exit ${code}

