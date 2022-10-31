#!/bin/bash
command=$1
cur=$(pwd)

startAbcp() {
	if [ ! -d "$cur/log" ]; then
		mkdir $cur/log
	fi
	echo "=====运行 abcp====="
	go run $cur/cmd/main.go -log_dir=$cur/log >$cur/log/log.out 2>&1 &
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
else
	echo "=====命令错误，请重试====="
fi
