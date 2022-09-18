package router

import (
	"CloudPlatform/pkg/handle"
	"CloudPlatform/pkg/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.Register(func(router gin.IRoutes) {
		router.POST("/containers", createContainer)
		router.GET("/containers/:id", getContainer)
		router.GET("/containers", getContainers)
		router.DELETE("/containers/:id", deleteContainer)
		router.GET("/containsers/:id/start", startContainer)
		router.GET("/containers/:id/stop", stopContainer)
		router.GET("/containers/:id/restart", restartContainer)
		router.GET("/containers/:id/connect", connectContainer)
		router.GET("/containers/:id/makeImage", makeImage)
		router.POST("/containers/:id/upload", uploadToContainer)
	}, router.V0)
}

// 创建容器
func createContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 根据 id 获取容器信息
func getContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 按条件获取容器
func getContainers(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 删除容器
func deleteContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 开启容器
func startContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 停止容器
func stopContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 重启容器
func restartContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 连接容器
func connectContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 根据容器制作镜像
func makeImage(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 将文件上传到容器里
func uploadToContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}
