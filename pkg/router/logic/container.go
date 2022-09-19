package router

import (
	"CloudPlatform/base"
	"CloudPlatform/pkg/handle"
	"CloudPlatform/pkg/router"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func init() {
	router.Register(func(router gin.IRoutes) {
		router.POST("/containers", createContainer)
		router.GET("/containers/:id", getContainer)
		router.GET("/containers", getContainers)
		router.DELETE("/containers/:id", removeContainer)
		router.GET("/containsers/:id/start", startContainer)
		router.GET("/containers/:id/stop", stopContainer)
		router.GET("/containers/:id/restart", restartContainer)
		router.GET("/containers/:id/connect", connectContainer)
		router.POST("/containers/:id/makeImage", makeImage)
		router.POST("/containers/:id/upload", uploadToContainer)
	}, router.V0)
}

// 创建容器
func createContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	json := make(map[string]interface{})
	// 拼接指令参数
	c.BindJSON(&json)
	var build strings.Builder
	for k, v := range json {
		build.WriteString(k)
		build.WriteString(fmt.Sprintf(" %v ", v))
	}
	cmd := fmt.Sprintf(base.CONTAINER_RUN, build.String())
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		r.Error(handle.CONTAINER_CREATE_FAIL)
		return
	}
	// TODO: 添加到数据库
	r.Success(string(out))
}

// 根据 id 获取容器信息
func getContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 按条件(用户)获取容器
func getContainers(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 删除容器
func removeContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	container := c.Query("id")
	cmd := fmt.Sprintf(base.CONTAINER_RM, container)
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		r.Error(handle.CONTAINER_REMOVE_FAIL)
		return
	}
	// TODO: 删除数据库
	r.Success(string(out))
}

// 开启容器
func startContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	container := c.Query("id")
	// TODO: 判断容器是否开启
	cmd := fmt.Sprintf(base.CONTAINER_START, container)
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		r.Error(handle.CONTAINER_START_FAIL)
		return
	}
	// TODO: 更新数据库
	r.Success(string(out))
}

// 停止容器
func stopContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	container := c.Query("id")
	// TODO: 判断容器是否关闭
	cmd := fmt.Sprintf(base.CONTAINER_STOP, container)
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		r.Error(handle.CONTAINER_STOP_FAIL)
		return
	}
	// TODO: 更新数据库
	r.Success(string(out))
}

// 重启容器
func restartContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	container := c.Query("id")
	cmd := fmt.Sprintf(base.CONTAINER_RESTART, container)
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		r.Error(handle.CONTAINER_RESTART_FAIL)
		return
	}
	// TODO: 更新数据库
	r.Success(string(out))
}

// 连接容器
func connectContainer(c *gin.Context) {
	// TODO: 开启容器
	// TODO: 安装 SSH
	// TODO: 返回 SSH IP PORT USER PWD
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 根据容器制作镜像
func makeImage(c *gin.Context) {
	r := handle.NewResponse(c)
	container := c.Query("id")
	m := make(map[string]interface{})
	c.BindJSON(&m)
	cmd := fmt.Sprintf(base.CONTAINER_COMMIT, m["author"], m["desc"], container, m["name"])
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		r.Error(handle.IMAGE_CREATE_FAIL)
		return
	}
	// TODO: 添加数据库
	r.Success(string(out))
}

// 将文件上传到容器里
func uploadToContainer(c *gin.Context) {
	// TODO: 开启容器
	// TODO: 安装 FTP
	// TODO: 文件上传
	r := handle.NewResponse(c)
	r.Success(nil)
}
