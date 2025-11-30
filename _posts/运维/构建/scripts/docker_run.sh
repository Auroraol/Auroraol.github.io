#!/bin/bash

BIN_NAME="goods_center"
IMAGE_NAME="registry.cn-hangzhou.aliyuncs.com/xiaoduoai/devops-golang-template:latest"
LOCAL_PATH="`pwd`"

docker rm -f ${BIN_NAME}
docker run -d --restart=always --net=host --name ${BIN_NAME} \
       -w /app \
       -v ${LOCAL_PATH}:/app/config \
       -v /etc/localtime:/etc/localtime \
       -v /var/log/xiaoduo:/var/log/xiaoduo \
       -e LANG="en_US.UTF-8" ${IMAGE_NAME} \
       ./${BIN_NAME} -c /app/config/${BIN_NAME}.conf
