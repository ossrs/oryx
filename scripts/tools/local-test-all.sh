#!/bin/bash

# Execute by: bash xxx.sh or bash zzz/yyy/xxx.sh or ./xxx.sh or ./zzz/yyy/xxx.sh source xxx.sh
REALPATH=$(realpath ${BASH_SOURCE[0]})
SCRIPT_DIR=$(cd $(dirname ${REALPATH}) && pwd)
WORK_DIR=$(cd $(dirname ${REALPATH})/../.. && pwd)
cd ${WORK_DIR}

# Check OS, must be darwin.
OS=$(uname -s)
if [[ ${OS} != Darwin ]]; then
    echo "Must run on macOS, current os is ${OS}"
    exit 1
fi

HELP=no
TARGET=all

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -h|--help) HELP=yes; shift ;;
        --target) TARGET=$2; shift 2;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
done

if [[ "$HELP" == yes ]]; then
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -h, --help    Show this help message and exit"
    echo "  --target      Test special target: all, script, aapanel, bt. default: $TARGET"
    exit 0
fi

if [[ $TARGET != all && $TARGET != script && $TARGET != aapanel && $TARGET != bt ]]; then
    echo "Unknown target $TARGET, should be script, aapanel, bt, or all"
    exit 1
fi

CONTAINERS="script bt aapanel"
echo "Test TARGET=$TARGET, CONTAINERS=$CONTAINERS"

echo "Remove all docker containers"
docker rm -f $CONTAINERS 2>/dev/null || echo 'OK'
echo "Remove all docker containers OK"

#####################################################################################
echo "Rebuild platform docker image" &&
docker rmi platform:latest 2>/dev/null || echo OK &&
docker build -t platform:latest -f Dockerfile . &&
docker save -o platform.tar platform:latest
ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Rebuild platform docker image failed, ret=$ret"; exit $ret; fi

#####################################################################################
if [[ $TARGET == all || $TARGET == script ]]; then
    echo "Test script installer"

    echo "Build script dev docker image"
    docker rmi srs-script-dev 2>/dev/null || echo 'OK' &&
    docker build -t srs-script-dev -f scripts/setup-ubuntu/Dockerfile.script .
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Build script dev docker image failed, ret=$ret"; exit $ret; fi

    echo "Run script dev docker image" &&
    docker rm -f $CONTAINERS 2>/dev/null || echo 'OK' &&
    docker run -p 2022:2022 -p 1935:1935/tcp -p 1985:1985/tcp \
        -p 8080:8080/tcp -p 8000:8000/udp -p 10080:10080/udp \
        --privileged -v /sys/fs/cgroup:/sys/fs/cgroup:rw --cgroupns=host \
        -d --rm -it -v $(pwd):/g -w /g --name=script srs-script-dev
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Run script dev docker image failed, ret=$ret"; exit $ret; fi

    echo "Load platform image to docker" &&
    version=$(bash scripts/version.sh) &&
    docker exec -it script docker load -i platform.tar &&
    docker exec -it script docker tag platform:latest ossrs/srs-stack:$version &&
    docker exec -it script docker tag platform:latest registry.cn-hangzhou.aliyuncs.com/ossrs/srs-stack:$version &&
    docker exec -it script docker images
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Load platform image to docker failed, ret=$ret"; exit $ret; fi

    echo "Setup script installer" &&
    docker exec -it bt rm -rf /data/* &&
    docker exec -it script bash build/srs_cloud/scripts/setup-ubuntu/uninstall.sh || echo OK &&
    bash scripts/setup-ubuntu/build.sh --output $(pwd)/build --extract &&
    docker exec -it script bash build/srs_cloud/scripts/setup-ubuntu/install.sh --verbose
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Setup script installer failed, ret=$ret"; exit $ret; fi

    echo "Test script installer" &&
    docker exec -it script make -j -C test &&
    docker exec -it script ./test/srs-stack.test -test.v -endpoint http://localhost:2022 \
        -srs-log=true -wait-ready=true -init-password=true \
        -check-api-secret=false -test.run TestApi_Empty &&
    bash scripts/tools/secret.sh --output test/.env &&
    docker exec -it script ./test/srs-stack.test -test.v -wait-ready -endpoint http://localhost:2022 \
        -srs-log=true -wait-ready=true -init-password=false \
        -check-api-secret=true \
        -test.parallel 1
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Test script installer failed, ret=$ret"; exit $ret; fi

    echo "Test script installer OK"
fi

#####################################################################################
if [[ $TARGET == all || $TARGET == aapanel ]]; then
    echo "Test aaPanel installer"

    echo "Build aaPanel dev docker image"
    docker rm -f $CONTAINERS 2>/dev/null || echo 'OK' &&
    AAPANEL_KEY=$(cat $HOME/.bt/api.json |awk -F token_crypt '{print $2}' |cut -d'"' -f3)
    docker run -p 80:80 -p 7800:7800 \
        -v $(pwd)/build/srs_cloud:/www/server/panel/plugin/srs_cloud \
        -v $HOME/.bt/api.json:/www/server/panel/config/api.json -e BT_KEY=$AAPANEL_KEY \
        --privileged -v /sys/fs/cgroup:/sys/fs/cgroup:rw --cgroupns=host \
        --add-host srs.cloud.local:127.0.0.1 \
        -d --rm -it -v $(pwd):/g -w /g --name=aapanel ossrs/aapanel-plugin-dev:1
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Build aaPanel dev docker image failed, ret=$ret"; exit $ret; fi

    echo "Load platform image to docker" &&
    version=$(bash scripts/version.sh) &&
    docker exec -it aapanel docker load -i platform.tar &&
    docker exec -it aapanel docker tag platform:latest ossrs/srs-stack:$version &&
    docker exec -it aapanel docker tag platform:latest registry.cn-hangzhou.aliyuncs.com/ossrs/srs-stack:$version &&
    docker exec -it aapanel docker images
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Load platform image to docker failed, ret=$ret"; exit $ret; fi

    echo "Setup aaPanel installer" &&
    docker exec -it aapanel rm -rf /data/* &&
    docker exec -it aapanel bash /www/server/panel/plugin/srs_cloud/install.sh uninstall || echo 'OK' &&
    bash scripts/setup-aapanel/auto/zip.sh --output $(pwd)/build --extract &&
    docker exec -it aapanel bash /www/server/panel/plugin/srs_cloud/install.sh install
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Setup aaPanel installer failed, ret=$ret"; exit $ret; fi

    echo "Test aaPanel installer" &&
    docker exec -it aapanel python3 /www/server/panel/plugin/srs_cloud/bt_api_remove_site.py &&
    docker exec -it aapanel python3 /www/server/panel/plugin/srs_cloud/bt_api_create_site.py &&
    docker exec -it aapanel python3 /www/server/panel/plugin/srs_cloud/bt_api_setup_site.py &&
    docker exec -it aapanel bash /www/server/panel/plugin/srs_cloud/setup.sh \
        --r0 /tmp/srs_cloud_install.r0 --nginx /www/server/nginx/logs/nginx.pid \
        --www /www/wwwroot --site srs.cloud.local
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Test aaPanel installer failed, ret=$ret"; exit $ret; fi

    echo "Test aaPanel installer" &&
    docker exec -it aapanel make -j -C test &&
    docker exec -it aapanel ./test/srs-stack.test -test.v -endpoint http://srs.cloud.local:80 \
        -srs-log=true -wait-ready=true -init-password=true \
        -check-api-secret=false -test.run TestApi_Empty &&
    bash scripts/tools/secret.sh --output test/.env &&
    docker exec -it aapanel ./test/srs-stack.test -test.v -wait-ready -endpoint http://srs.cloud.local:80 \
        -srs-log=true -wait-ready=true -init-password=false \
        -check-api-secret=true \
        -test.parallel 1
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Test aaPanel installer failed, ret=$ret"; exit $ret; fi

    echo "Test aaPanel installer OK"
fi

#####################################################################################
if [[ $TARGET == all || $TARGET == bt ]]; then
    echo "Test bt installer"

    echo "Build bt dev docker image"
    docker rm -f $CONTAINERS 2>/dev/null || echo 'OK' &&
    BT_KEY=$(cat $HOME/.bt/api.json |awk -F token_crypt '{print $2}' |cut -d'"' -f3)
    docker run -p 80:80 -p 7800:7800 \
        -v $(pwd)/build/srs_cloud:/www/server/panel/plugin/srs_cloud \
        -v $HOME/.bt/userInfo.json:/www/server/panel/data/userInfo.json \
        -v $HOME/.bt/api.json:/www/server/panel/config/api.json -e BT_KEY=$BT_KEY \
        --privileged -v /sys/fs/cgroup:/sys/fs/cgroup:rw --cgroupns=host \
        --add-host srs.cloud.local:127.0.0.1 \
        -d --rm -it -v $(pwd):/g -w /g --name=bt ossrs/bt-plugin-dev:1
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Build bt dev docker image failed, ret=$ret"; exit $ret; fi

    echo "Load platform image to docker" &&
    version=$(bash scripts/version.sh) &&
    docker exec -it bt docker load -i platform.tar &&
    docker exec -it bt docker tag platform:latest ossrs/srs-stack:$version &&
    docker exec -it bt docker tag platform:latest registry.cn-hangzhou.aliyuncs.com/ossrs/srs-stack:$version &&
    docker exec -it bt docker images
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Load platform image to docker failed, ret=$ret"; exit $ret; fi

    echo "Install bt installer" &&
    docker exec -it bt rm -rf /data/* &&
    docker exec -it bt bash /www/server/panel/plugin/srs_cloud/install.sh uninstall || echo 'OK' &&
    bash scripts/setup-bt/auto/zip.sh --output $(pwd)/build --extract &&
    docker exec -it bt bash /www/server/panel/plugin/srs_cloud/install.sh install
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Setup bt installer failed, ret=$ret"; exit $ret; fi

    echo "Setup bt installer" &&
    docker exec -it bt python3 /www/server/panel/plugin/srs_cloud/bt_api_remove_site.py &&
    docker exec -it bt python3 /www/server/panel/plugin/srs_cloud/bt_api_create_site.py &&
    docker exec -it bt python3 /www/server/panel/plugin/srs_cloud/bt_api_setup_site.py &&
    docker exec -it bt bash /www/server/panel/plugin/srs_cloud/setup.sh \
        --r0 /tmp/srs_cloud_install.r0 --nginx /www/server/nginx/logs/nginx.pid \
        --www /www/wwwroot --site srs.cloud.local
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Test bt installer failed, ret=$ret"; exit $ret; fi

    echo "Test bt installer" &&
    docker exec -it bt make -j -C test &&
    docker exec -it bt ./test/srs-stack.test -test.v -endpoint http://srs.cloud.local:80 \
        -srs-log=true -wait-ready=true -init-password=true \
        -check-api-secret=false -test.run TestApi_Empty &&
    bash scripts/tools/secret.sh --output test/.env &&
    docker exec -it bt ./test/srs-stack.test -test.v -wait-ready -endpoint http://srs.cloud.local:80 \
          -srs-log=true -wait-ready=true -init-password=false \
          -check-api-secret=true \
          -test.parallel 1
    ret=$?; if [[ 0 -ne ${ret} ]]; then echo "Test bt installer failed, ret=$ret"; exit $ret; fi

    echo "Test bt installer OK"
fi

#####################################################################################
docker rm -f $CONTAINERS 2>/dev/null

echo ""
echo "All tests OK"
