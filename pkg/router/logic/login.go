package logic

import (
	"CloudPlatform/pkg/handle"
	"CloudPlatform/pkg/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.Register(func(router gin.IRoutes) {
		router.GET("/test", test)
	}, router.V0)

	router.Register(func(router gin.IRoutes) {
		router.GET("/test")
	}, router.V1)
}

func test(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}
