package router

import (
	"cloud-platform/global"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
)

type Router func(h *route.RouterGroup)

type RouteLevel struct {
	Level int16
	Url   string
}

var (
	V1Router = RouteLevel{1, "/api/v1"}
	V2Router = RouteLevel{2, "/api/v2"}
	V3Router = RouteLevel{3, "/api/v3"}
)

var (
	routersSize = 0
	mapRouters  = map[RouteLevel][]Router{
		V1Router: make([]Router, 0),
		V2Router: make([]Router, 0),
		V3Router: make([]Router, 0),
	}
)

func RegisterRouter(router Router, level RouteLevel) {
	routersSize++
	mapRouters[level] = append(mapRouters[level], router)
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
	registerRouter(h)
	return h, nil
}

func registerRouter(h *server.Hertz) {
	global.Logger.Infof("register router start, size: %d", routersSize)
	for level, routers := range mapRouters {
		v := h.Group(level.Url)
		for _, router := range routers {
			router(v)
		}
	}
}
