package router

import (
	"cloud-platform/global"
	"fmt"

	_ "cloud-platform/pkg/router/api"
	router "cloud-platform/pkg/router/manager"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func RunServer() {
	h, err := listen()
	if err != nil {
		global.Logger.Errorf("Listen error: %v", err)
		panic(err.Error())
	}
	h.Spin()
}

func listen() (*server.Hertz, error) {
	h := server.Default(server.WithHostPorts(fmt.Sprintf("%s:%d", global.Config.App.Host, global.Config.App.Port)))
	router.RouteHandler.Register(h)
	return h, nil
}
