package logic

import (
	"CloudPlatform/cmd/web"
	"flag"

	"github.com/golang/glog"
)

func Run() {
	// 日志启动要放在最开始
	flag.Parse()
	defer glog.Flush()
	// 读配置相关的可以放在这

	/** gin 启动要放在最后*/
	web.Start()
}
