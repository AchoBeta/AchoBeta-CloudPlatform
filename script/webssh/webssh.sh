#!/bin/bash

echo "更新 yum 源"
cd /etc/yum.repos.d/
sed -i 's/mirrorlist/#mirrorlist/g' /etc/yum.repos.d/CentOS-*
sed -i 's|#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' /etc/yum.repos.d/CentOS-*
yum makecache
yum update -y

echo "安装 pip"
yum -y install --assumeyes python3-pip
# 清华源加速临时安装
pip3 install --upgrade pip -i https://pypi.tuna.tsinghua.edu.cn/simple/
# 清华源加速临时安装
pip3 install webssh -i https://pypi.tuna.tsinghua.edu.cn/simple/
