#!/bin/bash

### BEGIN INIT INFO
# Provides:          ossrs(srs-stack)
# Required-Start:    $all
# Required-Stop:     $all
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: ossrs(srs-stack)
# Description:       https://github.com/ossrs/srs-stack
### END INIT INFO

SRS_HOME=/usr/local/srs-stack

start() {
  systemctl start srs-stack.service
}

stop() {
  systemctl stop srs-stack.service
}

status() {
  systemctl status srs-stack.service
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

