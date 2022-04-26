#!/bin/bash

### BEGIN INIT INFO
# Provides:          ossrs(srs-cloud)
# Required-Start:    $all
# Required-Stop:     $all
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: ossrs(srs-cloud)
# Description:       https://github.com/ossrs/srs-cloud
### END INIT INFO

SRS_HOME=/usr/local/lighthouse/softwares/srs-cloud

start() {
  systemctl start srs-cloud.service
}

stop() {
  systemctl stop srs-cloud.service
  bash $SRS_HOME/scripts/remove-containers.sh
}

status() {
  systemctl status srs-cloud.service
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

