package web

import (
	//_ "CloudPlatform/conf/secret"
	"CloudPlatform/pkg/router"
	_ "CloudPlatform/pkg/router/logic"
	_ "CloudPlatform/util"
)

func Run() {
	router.Run()
}
