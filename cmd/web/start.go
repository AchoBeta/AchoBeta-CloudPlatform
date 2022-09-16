package web

import (
	_ "CloudPlatform/conf/secret"
	"CloudPlatform/pkg/router"
	_ "CloudPlatform/util"
)

func Run() {
	router.Run()
}
