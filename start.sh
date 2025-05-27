#!/bin/bash

config_filename="config.yaml"

if [ $1 ]
then
    config_filename=($1)
fi

# 启动 broker grpc 服务
bin/server -conf "configs/${config_filename}"

# 等待所有后台进程
wait
