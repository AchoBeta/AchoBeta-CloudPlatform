package manager

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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
	// 路由级别, 规定数字越大级别越高
	// 在某种情况下, 高级别的路由会适用于低级别的路由
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
	routeCount, middlewareCount := 0, 0
	for _, route := range rm.Routes {
		v := h.Group(route.Url)
		// 中间件注册
		for _, middleware := range route.Middlewares {
			middlewareCount++
			v.Use(middleware())
		}
		// 路由注册
		for _, router := range route.Path {
			routeCount++
			router(v)
		}
	}
	hlog.Infof("Registering routes, total routes: %d, total middlewares: %d", routeCount, middlewareCount)
}

func (rm *RouteManager) RegisterRouter(level RouteLevel, router PathHandler) {
	rm.checkRoute(level)
	rm.Routes[level].Path = append(rm.Routes[level].Path, router)
}

// @title RegisterMiddleware
// @description 注册中间件
// @param level RouteLevel 路由级别
// @param middleware Middleware 中间件
// @param iteration bool 是否迭代, 即 v3 等级的中间件会同时注册到 v2, v1, v0 等级 (向下兼容)
func (rm *RouteManager) RegisterMiddleware(level RouteLevel, middleware Middleware, iteration bool) {
	rm.checkRoute(level)
	if iteration {
		for e := level; e >= 0; e-- {
			rm.Routes[e].Middlewares = append(rm.Routes[e].Middlewares, middleware)
		}
		return
	}
	rm.Routes[level].Middlewares = append(rm.Routes[level].Middlewares, middleware)
}

// @title requestGlobalMiddleware
// @description 注册全局中间件
// func requestGlobalMiddleware(v *route.RouterGroup) {
// 	v.Use(requestid.New())
// }
