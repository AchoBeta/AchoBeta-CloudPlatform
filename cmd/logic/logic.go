package logic

import (
	"CloudPlatform/cmd/web"
	"CloudPlatform/util/redis"
	"context"
	"flag"
	"log"

	"github.com/golang/glog"
)

func Run() {
	// 日志启动要放在最开始
	flag.Parse()
	defer glog.Flush()
	err := redis.Connect(context.Background())
	if err != nil {
		log.Fatal("redis连接失败")
	}
	defer redis.Rdb.Close()
	// 读配置相关的可以放在这

	/** gin 启动要放在最后*/
	web.Run()
}
