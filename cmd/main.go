package main

import (
	"cloud-platform/pkg/load"
	"cloud-platform/pkg/router"
	"flag"
)

func main() {
	flag.Parse()
	// 初始化工程
	load.Init()
	// 工程进入前夕，释放资源
	defer load.Eve()
	/** gin 启动要放在最后*/
	router.RunS()
}
