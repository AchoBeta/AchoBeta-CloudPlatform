#!/bin/bash
command=$1
cur=$(pwd)
updateAbcp() {
    if [ -e "$cur/go.mod" ]; then
        rm -r $cur/go.mod
    fi
    startAbcp
}

startAbcp() {
    if [ ! -e "$cur/go.mod" ]; then
        echo "=====生成 go.mod====="
        touch $cur/go.mod
        echo -e 'module CloudPlatform\n\ngo 1.18' > $cur/go.mod
        echo "=====更新依赖====="
        go get -u all
    fi
    if [ ! -d "$cur/log" ]; then
        mkdir $cur/log
    fi
    echo "=====运行 abcp====="
    go run $cur/cmd/main.go -log_dir=$cur/log > $cur/logs/log.out 2>&1 &
}
stopAbcp() {
    echo "=====停止 abcp====="
    lsof -i :1210 | awk 'NR>1 {print $2}' | xargs kill -9
}
restartAbcp() {
    stopAbcp
    startAbcp
}

if [ $command == "start" ]; then
    startAbcp
elif [ $command == "stop" ]; then
    stopAbcp
elif [ $command == "restart" ]; then
    restartAbcp
elif [ $command == "update" ]; then
    updateAbcp
else
    echo "=====命令错误，请重试====="
fi