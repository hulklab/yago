#!/bin/bash

cd `dirname $0`
cd ..
pwd
mkdir -p bin
cp conf/app.toml bin/
APP_NAME=`cat bin/app.toml|grep -E "^app_name=|^app_name "|awk -F "=" '{print $2}'|xargs echo`
go build -o bin/$APP_NAME main.go
go build -o bin/$APP_NAME"-cmd" cmd/main.go
