package router

import (
	"CloudPlatform/pkg/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

const (
	V0   uint8  = 0
	V1   uint8  = 1
	V2   uint8  = 2
	HOST string = ":1210"
)

var (
	_hooks_V0, _hooks_V1, _hooks_V2 []Hook
)

type Hook func(router gin.IRoutes)

func Register(hook Hook, hookType uint8) {
	switch hookType {
	case V0:
		_hooks_V0 = append(_hooks_V0, hook)
		break
	case V1:
		_hooks_V1 = append(_hooks_V1, hook)
		break
	case V2:
		_hooks_V2 = append(_hooks_V2, hook)
		break
	default:
		glog.Error("Register Error")
	}
}

func Run() {
	r := Listen()
	fmt.Println("run!")
	// 监听端口
	r.Run(HOST)
}

// gin 配置
func Listen() *gin.Engine {
	r := gin.New()
	/** 中间件部分 */
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	/** 路由登记部分 */
	// v0 模块, 无需权限校验
	v0 := r.Group("/api")
	{
		RegisterRouter(_hooks_V0, v0)
	}

	// v1 模块, 1级权限(需要登陆), 使用Token鉴权中间件
	v1 := r.Group("/api", middleware.TokenVer())
	{
		RegisterRouter(_hooks_V1, v1)
	}

	// v2 模块, 2级权限(校验token和是否有权操作容器)
	v2 := r.Group("/api")
	{
		r.Use(middleware.TokenVer())
		r.Use(middleware.ContainerVer())
		RegisterRouter(_hooks_V2, v2)
	}
	return r
}

func RegisterRouter(hooks []Hook, v *gin.RouterGroup) {
	for _, hook := range hooks {
		hook(v)
	}
}
