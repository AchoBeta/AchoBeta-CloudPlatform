package web

import (
<<<<<<< HEAD
	_ "CloudPlatform/conf/secret"
	"CloudPlatform/pkg/router"
	_ "CloudPlatform/pkg/router/logic"
	_ "CloudPlatform/util"
=======
	"CloudPlatform/pkg/router"
	_ "CloudPlatform/pkg/router/logic"
>>>>>>> master
)

func Run() {
	router.Run()
}
