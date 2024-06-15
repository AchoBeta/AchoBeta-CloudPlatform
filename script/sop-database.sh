#!/bin/bash

command=$1
current_dir=$(cd $(dirname $0); pwd)
organize=159.75.216.233:5000
name=abcp_database
tag=0.2
redis_port=6379
mongo_port=27017
mongo_username=root
mongo_password=123456

# 启动数据库
startDatabase() {
    docker ps -a | grep $name &> /dev/null
    if [ $? -eq 0 ]; then
        echo "=====运行 database====="
        docker start redis
	docker start mongo
    else
	echo "=====创建并运行 reids ====="
        docker run -d --restart unless-stopped --name redis \
            -p $redis_port:$redis_port --net bridge \
            -v $current_dir/database/data/redis:/data/redis \
            -v $current_dir/database/:/etc/redis/ \
            -v $current_dir/database/data/redis.log:/var/log/redis/redis-server.log \
            -d redis:latest redis-server /etc/redis/redis.conf

        echo "=====创建并运行 reids ====="
        docker run -d --name mongo --restart unless-stopped \
            --net bridge -p $mongo_port:$mongo_port \
            -e MONGO_INITDB_ROOT_USERNAME=$mongo_username -e MONGO_INITDB_ROOT_PASSWORD=$mongo_password \
            -v $current_dir/database/data/mongo:/data/mongo \
            --privileged=true mongo:4.4
    fi
}

# 停止数据库
stopDatabase() {
    echo "=====停止 database 容器====="
    docker stop redis
    docker stop mongo
}

# 重启数据库
restartDatabase() {
    stopDatabase
    startDatabase
}

# 更新数据库
updateDatabase() {
    # echo "=====删除 database 容器====="
    # docker rm -f $(docker ps -a | grep $name | awk '{print $1}')
    # echo "=====删除 database 镜像====="
    # docker rmi  $(docker images | grep $organize/$name | awk '{print $3}')
    # echo "=====拉取最新 database 镜像====="
    # docker pull $organize/$name:$tag
    # startDatabase
}


if [ $command == "start" ]; then
    startDatabase
elif [ $command == "stop" ]; then
    stopDatabase
elif [ $command == "restart" ]; then
    restartDatabase
elif [ $command == "update" ]; then
    updateDatabase
else
	echo "=====命令错误，请重试====="
fi
