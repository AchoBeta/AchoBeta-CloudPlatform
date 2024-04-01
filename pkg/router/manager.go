package router

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
)

type PathHandler func(h *route.RouterGroup)
type Middleware func() app.HandlerFunc
type RouteLevel int32
type RouteManager struct {
	Routes map[RouteLevel]*Route
}
type Route struct {
	Url         string
	Path        []PathHandler // 注册的path
	Middlewares []Middleware  // 注册的中间件
}

const (
	LEVEL_GLOBAL RouteLevel = 0 // 匿名级别路由
	LEVEL_V1     RouteLevel = 1
	LEVEL_V2     RouteLevel = 2
	LEVEL_V3     RouteLevel = 3
)

var (
	RouteHandler = &RouteManager{
		Routes: make(map[RouteLevel]*Route, 0),
	}
)

func buildUrl(level RouteLevel) string {
	return fmt.Sprintf("/api/v%d", level)
}

func NewRoute(level RouteLevel) *Route {
	return &Route{
		Url:         buildUrl(level),
		Path:        make([]PathHandler, 0),
		Middlewares: make([]Middleware, 0),
	}
}

func (rm *RouteManager) checkRoute(level RouteLevel) {
	if _, ok := rm.Routes[level]; !ok {
		rm.Routes[level] = NewRoute(level)
	}
}

func (rm *RouteManager) Register(h *server.Hertz) {
	for _, route := range rm.Routes {
		v := h.Group(route.Url)
		// 中间件注册
		for _, middleware := range route.Middlewares {
			v.Use(middleware())
		}
		// 路由注册
		for _, router := range route.Path {
			router(v)
		}
	}
}

func (rm *RouteManager) RegisterRouter(level RouteLevel, router PathHandler) {
	rm.checkRoute(level)
	rm.Routes[level].Path = append(rm.Routes[level].Path, router)
}

func (rm *RouteManager) RegisterMiddleware(level RouteLevel, middleware Middleware) {
	rm.checkRoute(level)
	rm.Routes[level].Middlewares = append(rm.Routes[level].Middlewares, middleware)
}
