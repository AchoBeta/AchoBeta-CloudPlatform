#!/bin/bash

command=$1

# 安装 docker
installDocker() {
    docker -v
    if [ $? -eq 0 ]; then
        echo "=====已有 docker===="
        exit
    fi
	echo "=====安装 docker====="
	# 设置 docker 的 yum 源地址
	yum-config-manager \
		--add-repo \
		https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
	sed -i 's/download.docker.com/mirrors.aliyun.com\/docker-ce/g' /etc/yum.repos.d/docker-ce.repo
	yum makecache fast
	yum install -y docker-ce
    # 更换 docker 源
    mkdir -p /etc/docker
    tee /etc/docker/daemon.json <<-'EOF'
    {
        "registry-mirrors": ["https://e46bdzxc.mirror.aliyuncs.com"]
    }
EOF
    service docker start
    systemctl daemon-reload
    systemctl restart docker
}

# 更新 docker
updateDocker() {
	echo "=====删除 docker====="
	yum remove -y docker docker-common docker-selinux docker-engine
	installDocker
}

if [ $command == "install" ]; then
	installDocker
elif [ $command == "update" ]; then
	updateDocker
else
	echo "======命令错误，请重试====="
fi