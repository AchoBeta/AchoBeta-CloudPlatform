#!/bin/bash

# 定义镜像名称和容器名称
name=webssh
tag=dev
port=8888


# 判断容器是否存在
docker ps -a | grep $name &> /dev/null
if [ $? -eq 0 ]; then
	echo "删除容器"
	docker rm -f $(docker ps -a | grep $name | awk "{print $1}")
fi

# 判断镜像是否存在
docker images | grep $name &> /dev/null
if [ $? -eq 0 ]; then
	echo "删除镜像"
	docker rmi $(docker images | grep $name | awk "{print $3}")
fi

docker build -f ./webssh/Dockerfile -t $name:$tag .

docker run -tid --net host --name $name -p $port:8888 $name:$tag /bin/bash
