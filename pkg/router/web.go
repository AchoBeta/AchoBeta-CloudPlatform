package router

import (
	"cloud-platform/global"
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

type Hook func(ctx context.Context, c *app.RequestContext)
type Router func(h *server.Hertz)

var (
	routers []Router
)

func RegisterRouter(router Router) {
	routers = append(routers, router)
}

func RunS() {
	h, err := Listen()
	if err != nil {
		global.Logger.Errorf("Listen error: %v", err)
		panic(err.Error())
	}
	h.Spin()
}

func Listen() (*server.Hertz, error) {
	h := server.Default(server.WithHostPorts(fmt.Sprintf("%s:%d", global.Config.App.Host, global.Config.App.Port)))
	for _, router := range routers {
		router(h)
	}
	return h, nil
}
