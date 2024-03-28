package router

import (
	"cloud-platform/global"
	"cloud-platform/internal/middleware"
	"fmt"

	"github.com/golang/glog"
)

const (
	V0 uint8 = 0
	V1 uint8 = 1
	V2 uint8 = 2
	V3 uint8 = 3
)

var (
	_hooks_V0, _hooks_V1, _hooks_V2, _hooks_V3 []Hook
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
	case V3:
		_hooks_V3 = append(_hooks_V3, hook)
		break
	default:
		glog.Error("Register Error")
	}
}

func Run() {
	r := Listen()
	fmt.Println("run!")
	// 监听端口
	r.Run(fmt.Sprintf("%s:%d", global.Config.App.Host, global.Config.App.Port))
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
	RegisterRouter(_hooks_V0, v0)

	// v1 模块, 使用Token鉴权中间件
	v1 := r.Group("/api", middleware.TokenVer())
	{
		RegisterRouter(_hooks_V1, v1)
	}

	// v2 模块，管理员模块
	v2 := r.Group("/api")
	{
		v2.Use(middleware.TokenVer(), middleware.AdminVer())
		RegisterRouter(_hooks_V2, v2)
	}

	// v3 模块, 是否有权操作容器
	v3 := r.Group("/api")
	{
		v3.Use(middleware.TokenVer(), middleware.ContainerVer())
		RegisterRouter(_hooks_V3, v3)
	}

	return r
}

func RegisterRouter(hooks []Hook, v *gin.RouterGroup) {
	for _, hook := range hooks {
		hook(v)
	}
}
