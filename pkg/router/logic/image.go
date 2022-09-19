package router

import (
	"CloudPlatform/base"
	"CloudPlatform/pkg/handle"
	"CloudPlatform/pkg/router"
	"fmt"
	"os/exec"
	"strings"

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
	// TODO: 查数据库
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 根据 id 获取镜像信息（暂不需要）
func getImage(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
}

// 通过上传 Dockerfile 制作镜像
func buildImage(c *gin.Context) {
	r := handle.NewResponse(c)
	imageName := c.BindQuery("image-name")
	// TODO: 上传文件并保存到 `/docker/` 目录
	fileName := "TestDockerFile" // 生成的文件名
	cmd := fmt.Sprintf(base.IMAGE_BUILD, fileName, imageName)
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		r.Error(handle.IMAGE_CREATE_FAIL)
		return
	}
	// TODO: 保存到数据库
	r.Success(string(out))
}

// 删除镜像（逻辑删除）
func deleteImage(c *gin.Context) {
	r := handle.NewResponse(c)
	// TODO: 删除数据库
	r.Success(nil)
}

// 在 DockerHub 搜索镜像
func searchImages(c *gin.Context) {
	r := handle.NewResponse(c)
	image := c.Query("image")
	tag := c.Query("tag")
	cmd := fmt.Sprintf(base.IMAGE_SEARCH, image+":"+tag)
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		r.Error(handle.IMAGE_PUSH_FAIL)
		return
	}
	ss := strings.Split(out, "\n")
	var data []string
	for i := 1; i < len(ss)-1; i++ {
		data = append(data, ss[i][0:strings.Index(ss[i], " ")])
	}
	r.Success(data)
}

// 上传镜像到 DockerHub
func pushImage(c *gin.Context) {
	r := handle.NewResponse(c)
	image := c.Query("image")
	tag := c.Query("tag")
	cmd := fmt.Sprintf(base.IMAGE_PUSH, image+":"+tag)
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		r.Error(handle.IMAGE_PUSH_FAIL)
	} else {
		r.Success(string(out))
	}
}

func executor(name, arg string) (string, error) {
	out, err := exec.Command(name, strings.Split(arg, " ")...).Output()
	return string(out), err
}
