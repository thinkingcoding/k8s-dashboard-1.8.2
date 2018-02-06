#!/bin/bash

export CUR_DIR=$(cd "$(dirname "$0")";pwd)

docker run -v $CUR_DIR:/dashboard -w /dashboard lxc968/ubuntu-dev:2.1 ./node_modules/.bin/gulp build
if [ $? -ne 0 ]; then
    exit 1
fi

cd $CUR_DIR/dist/amd64
cp $CUR_DIR/src/deploy/cdf/Dockerfile $CUR_DIR/dist/amd64/
docker build -t localhost:5000/itom-k8s-dashboard:1.8.2 .
if [ $? -ne 0 ]; then
    exit 2
fi