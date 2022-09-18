package router

import (
	"CloudPlatform/pkg/handle"
	"CloudPlatform/pkg/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.Register(func(router gin.IRoutes) {
		router.GET("/images", getImages)
		router.GET("/images/:id", getImage)
		router.POST("/images/:id/build", buildImage)
		router.DELETE("/images/:id", deleteImage)
		router.GET("/images/search", searchImages)
		router.GET("/images/:id/push", pushImage)
	}, router.V0)
}

// 获取本地所有镜像
func getImages(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 根据 id 获取镜像信息
func getImage(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 通过上传 Dockerfile 制作镜像
func buildImage(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 删除镜像（逻辑删除）
func deleteImage(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 在 DockerHub 搜索镜像
func searchImages(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 上传镜像到 DockerHub
func pushImage(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}
