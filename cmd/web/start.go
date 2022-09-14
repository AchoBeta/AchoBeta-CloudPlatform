package web

import (
	"CloudPlatform/pkg/router"
	_ "CloudPlatform/pkg/router/logic"
)

func Run() {
	router.Run()
}
