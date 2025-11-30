#!/bin/bash

SERVER_NAME=goods_center
SERVER_GROUP=ecrobot
SRC_IMG=registry.cn-zhangjiakou.aliyuncs.com/xiaoduoai/golang:1.13-pulsar2.41

VERSION=$(git describe --tags 2>/dev/null)
COMMIT=$(git rev-parse --short HEAD)
TIME=$(date +%FT%T)
if [[ -z ${VERSION} ]]; then
    VERSION=${COMMIT}
fi

BUILD_IMG_NAME=go_build

docker stop $BUILD_IMG_NAME || echo "skip stop $BUILD_IMG_NAME"
docker rm -f $BUILD_IMG_NAME || echo "skip stop $BUILD_IMG_NAME"

echo "pack ..."

echo "build ..."

docker run -v $(pwd):/go \
  -w /go \
  --name $BUILD_IMG_NAME \
  ${SRC_IMG} \
  /bin/bash -c "export GO111MODULE=on && pwd && ls && env GOOS=linux GOARCH=amd64 go build && echo \"---build success---\" && exit 0 || echo \"---build failed---\" && exit 1"

docker_build_state=$?
if [[ ${docker_build_state} -ne 0 ]]; then
  echo "build failed"
  exit 1
fi

TAG=${SERVER_GROUP}-${SERVER_NAME}:${VERSION}
docker rmi ${TAG} --force || echo "skip docker rmi ${TAG}"

docker build -f ./Dockerfile -t ${TAG} ./

docker_build_state=$?
if [[ ${docker_build_state} -ne 0 ]]; then
  echo "pack failed"
  exit 1
fi

echo "build done"

echo "push to image hub"

docker tag ${TAG} registry.cn-zhangjiakou.aliyuncs.com/xiaoduoai/${TAG}

docker push registry.cn-zhangjiakou.aliyuncs.com/xiaoduoai/${TAG}

docker_build_state=$?
if [[ ${docker_build_state} -ne 0 ]]; then
  echo "push to registry failed"
  exit 1
fi

echo "push to image done"

echo "pack done"

