package router

import (
	"cloud-platform/global"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
)

type Router func(h *route.RouterGroup)
type Middleware func() app.HandlerFunc
type RouteLevel int32
type Route struct {
	Url         string
	Routers     []Router
	Middlewares []Middleware
}

const (
	LEVEL_GLOBAL RouteLevel = 0
	LEVEL_V1     RouteLevel = 1
	LEVEL_V2     RouteLevel = 2
	LEVEL_V3     RouteLevel = 3
)

var (
	routersSize = 0
	mapUrl      = map[RouteLevel]string{
		LEVEL_GLOBAL: "",
		LEVEL_V1:     "/api/v1",
		LEVEL_V2:     "/api/v2",
		LEVEL_V3:     "/api/v3",
	}

	mapRouters = map[RouteLevel][]Router{
		LEVEL_GLOBAL: make([]Router, 0),
		LEVEL_V1:     make([]Router, 0),
		LEVEL_V2:     make([]Router, 0),
		LEVEL_V3:     make([]Router, 0),
	}

	mapMiddlewares = map[RouteLevel][]Middleware{
		LEVEL_GLOBAL: make([]Middleware, 0),
		LEVEL_V1:     make([]Middleware, 0),
		LEVEL_V2:     make([]Middleware, 0),
		LEVEL_V3:     make([]Middleware, 0),
	}
)

func RegisterMiddleware(middleware Middleware, level RouteLevel) {
	mapMiddlewares[level] = append(mapMiddlewares[level], middleware)
	switch level {
	case LEVEL_V2:
		RegisterMiddleware(middleware, LEVEL_V3)
	case LEVEL_V1:
		RegisterMiddleware(middleware, LEVEL_V2)
	case LEVEL_GLOBAL:
		RegisterMiddleware(middleware, LEVEL_V1)
	}
}

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
	registerMiddleware(h)
	return h, nil
}

func registerRouter(h *server.Hertz) {
	global.Logger.Infof("register router start, size: %d", routersSize)
	for level, routers := range mapRouters {
		v := h.Group(mapUrl[level])
		for _, router := range routers {
			router(v)
		}
	}
}

func registerMiddleware(h *server.Hertz) {
	for level, middlewares := range mapMiddlewares {
		v := h.Group(mapUrl[level])
		for _, middleware := range middlewares {
			v.Use(middleware())
		}
	}
}
