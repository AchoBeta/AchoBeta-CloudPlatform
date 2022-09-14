package web

import (
	"CloudPlatform/cmd/web/middleware"
	"CloudPlatform/pkg/handle"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Start() {
	c := NewRouter()
	fmt.Println("run!")
	c.Run(":1210")
}

// gin 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()
	/** 中间件都在这个部分*/
	r.GET("/test", func(ctx *gin.Context) {
		r := handle.NewResponse(ctx)
		r.Success(nil)
		return
	})

	r.POST("/login", func(c *gin.Context) {
		r := handle.NewResponse(c)
		r.Success(nil)
		return
	})
	// v1 模块
	v1 := r.Group("/api/v1", middleware.TokenVer())
	{
		v1.GET("/test1", func(ctx *gin.Context) {
			r := handle.NewResponse(ctx)
			r.Success(nil)
			return
		})

		v1.POST("/test")
	}
	return r
}
