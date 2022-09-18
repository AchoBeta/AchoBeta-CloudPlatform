#!/bin/bash

command=$1
orginaze=achobeta
name=abcp_webssh
tag=0.1
port=8888

# 启动 webssh
startWebssh() {
    docker ps -a | grep $name &> /dev/null
    if [ $? -eq 0 ]; then
        echo "=====运行 database====="
        docker start $name
    else
        echo "=====创建并运行 database====="
        docker run -d --name $name --restart unless-stopped \
            --net bridge -p $port:$port $orginaze/$name:$tag
    fi
    echo "=====运行成功，请访问：xxxx:8888====="
}

# 停止 webssh
stopWebssh() {
    echo "=====停止 webssh 容器====="
    docker stop $name
}

# 重启 webssh
restartWebssh() {
    stopWebssh
    startWebssh
}

# 更新 webssh
updateWebssh() {
    echo "=====删除 webssh 容器====="
    docker rm -f $(docker ps -a | grep $name | awk '{print $1}')
    echo "=====删除 webssh 镜像====="
    docker rmi -f  $(docker images | grep $orginaze/$name | awk '{print $3}')
    echo "=====拉取最新 webssh 镜像====="
    docker pull $orginaze/$name::$tag
    startWebssh
}

if [ $command == "start" ]; then
    startWebssh
elif [ $command == "stop" ]; then
    stopWebssh
elif [ $command == "restart" ]; then
    restartWebssh
elif [ $command == "update" ]; then
    updateWebssh
else
    echo "======命令错误，请重试====="
fi
