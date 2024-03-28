package main

import (
	"cloud-platform/internal/exec"
	"cloud-platform/pkg/router"
	"flag"
)

func main() {
	flag.Parse()
	// 初始化工程
	exec.Init()
	// 工程进入前夕，释放资源
	defer exec.Eve()
	/** gin 启动要放在最后*/
	router.Run()
}
