#!/bin/bash

cd `dirname $0`

rm -rf timer
mkdir timer
protoc --go_out=plugins=grpc:./timer timer.proto
cd - > /dev/null 2>&1
