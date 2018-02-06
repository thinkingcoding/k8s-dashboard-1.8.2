#!/bin/bash

export CUR_DIR=$(cd "$(dirname "$0")";pwd)

docker run -v $CUR_DIR:/dashboard -w /dashboard lxc968/ubuntu-dev:2.1 npm install --unsafe-perm
if [ $? -ne 0 ]; then
    exit 1
fi
