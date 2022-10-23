package logic

import (
	"CloudPlatform/cmd/web"
	_ "CloudPlatform/internal/router/api"
	"flag"
)

func Run() {
	flag.Parse()
	// 初始化工程
	Init("./config.yaml")
	// 工程进入前夕，释放资源
	defer Eve()
	/** gin 启动要放在最后*/
	web.Run()
}
