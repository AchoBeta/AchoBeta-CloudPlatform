package api

import (
	"CloudPlatform/global"
	"CloudPlatform/internal/base"
	"CloudPlatform/internal/handle"
	"CloudPlatform/internal/router"
	"CloudPlatform/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func init() {
	router.Register(func(router gin.IRoutes) {
		router.GET("/images", getImages)
		router.GET("/images/:id", getImageInfo)
		//router.POST("/images/:id/build", buildImage)
		//router.GET("/images/search", searchImages)
		router.GET("/images/:id/push", pushImage)
	}, router.V1)
	router.Register(func(router gin.IRoutes) {
		router.DELETE("/images/:id", deleteImage)
	}, router.V2)
}

// 获取本地所有镜像
func getImages(c *gin.Context) {
	// TODO: 查数据库
	r := handle.NewResponse(c)
	err, code, images := service.GetImages()
	if code == 0 {
		r.Success(images)
	} else if code == 1 {
		glog.Errorf("[db] find images error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	} else if code == 2 {
		glog.Errorf("decode image error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}

// 根据 id 获取镜像信息
func getImageInfo(c *gin.Context) {
	r := handle.NewResponse(c)
	id := c.Param("id")
	image := &base.Image{}
	err, code := service.GetImageInfo(id, image)
	if code == 0 {
		r.Success(image)
	} else if code == 1 {
		r.Error(handle.IMAGE_NOT_FIND)
	} else if code == 2 {
		glog.Errorf("[db] find image by id error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}

// 通过上传 Dockerfile 制作镜像
// func buildImage(c *gin.Context) {
// 	r := handle.NewResponse(c)
// 	imageName := c.BindQuery("image-name")
// 	// TODO: 上传文件并保存到 `/docker/` 目录
// 	fileName := "TestDockerFile" // 生成的文件名
// 	cmd := fmt.Sprintf(base.IMAGE_BUILD, fileName, imageName)
// 	out, err := exec.Command(base.DOCKER, cmd)
// 	if err != nil {
// 		r.Error(handle.IMAGE_CREATE_FAIL)
// 		return
// 	}
// 	// TODO: 保存到数据库
// 	r.Success(string(out))
// }

// 删除镜像（逻辑删除）
func deleteImage(c *gin.Context) {
	r := handle.NewResponse(c)
	id := c.Param("id")
	err, code := service.DeleteImage(id)
	if code == 0 {
		r.Success(nil)
	} else if code == 1 {
		r.Error(handle.IMAGE_NOT_FIND)
	} else if code == 2 {
		glog.Errorf("[db] delete image error ! msg: %v\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}

// 在 DockerHub 搜索镜像
// func searchImages(c *gin.Context) {
// 	r := handle.NewResponse(c)
// 	image := c.Param("image")
// 	tag := c.Param("tag")
// 	cmd := fmt.Sprintf(base.IMAGE_SEARCH, image+":"+tag)
// 	out, err := executor(base.DOCKER, cmd)
// 	if err != nil {
// 		r.Error(handle.IMAGE_PUSH_FAIL)
// 		return
// 	}
// 	ss := strings.Split(out, "\n")
// 	var data []string
// 	for i := 1; i < len(ss)-1; i++ {
// 		data = append(data, ss[i][0:strings.Index(ss[i], " ")])
// 	}
// 	r.Success(data)
// }

// TODO: 上传镜像到 DockerHub
func pushImage(c *gin.Context) {
	r := handle.NewResponse(c)
	image := c.Param("image")
	tag := c.Param("tag")
	var err error
	var code int8
	if global.Config.App.Type == "docker" {
		err, code = service.PushDockerImage(image + ":" + tag)
	} else {
		err, code = service.PushK8SImage(image + ":" + tag)
	}
	if code == 0 {
		r.Success(nil)
	} else if code == 1 {
		glog.Errorf("[cmd] push image error ! msg: %s/n", err.Error())
		r.Error(handle.IMAGE_PUSH_FAIL)
	}
}
