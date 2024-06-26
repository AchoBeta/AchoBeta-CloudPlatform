package api

import (
	"cloud-platform/global"
	"cloud-platform/internal/base"
	"cloud-platform/internal/base/cloud"
	"cloud-platform/internal/handle"
	"cloud-platform/internal/router"
	"cloud-platform/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func init() {
	router.Register(func(router gin.IRoutes) {
		router.POST("/containers", createContainer)
		router.GET("/containers", getContainers)
	}, router.V1)

	router.Register(func(router gin.IRoutes) {
		router.GET("/containers/:id", getContainer)
		router.DELETE("/containers/:id", removeContainer)
		router.GET("/containers/:id/start", startContainer)
		router.GET("/containers/:id/stop", stopContainer)
		router.GET("/containers/:id/restart", restartContainer)
		router.POST("/containers/:id/make-image", makeImage)
		router.GET("/containers/:id/log", getContainerLog)
	}, router.V3)
}

// 创建容器
func createContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	container := &cloud.Container{}
	err := c.ShouldBind(&container)
	user, _ := c.Get("user")
	if err != nil || container.Image == "" || container.Name == "" {
		r.Error(handle.PARAM_NOT_VALID)
		return
	}
	code, err := service.CreateContainer(c.GetHeader("Authorization"), container, user.(*base.User))
	if code == 0 {
		r.Success(nil)
	} else if code == 1 {
		glog.Errorf("[cmd] run container error ! msg: %s\n", err.Error())
		r.Error(handle.CONTAINER_CREATE_FAIL)
		return
	} else if code == 2 {
		glog.Errorf("[cmd] set container ssh pwd error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	} else if code == 3 || code == 4 {
		glog.Errorf("[db] insert container error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}

// 根据 id 获取容器信息
func getContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	id := c.Param("id")
	container := &cloud.Container{}
	code, err := service.GetContainer(id, container)
	if code == 0 {
		r.Success(container)
	} else if code == 1 {
		r.Error(handle.CONTAINER_NOT_FOUND)
	} else if code == 2 {
		glog.Errorf("[db] find container error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	} else if code == 3 {
		glog.Errorf("decode container error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}

// 根据用户获取容器
func getContainers(c *gin.Context) {
	r := handle.NewResponse(c)
	user, _ := c.Get("user")
	code, containers, err := service.GetContainers(user.(*base.User))
	if code == 0 {
		r.Success(containers)
	} else if code == 1 {
		glog.Errorf("[db] find containers error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	} else if code == 2 {
		glog.Errorf("decode container error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}

// 删除容器
func removeContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	id := c.Param("id")
	user, _ := c.Get("user")
	code, err := service.RemoveContainer(c.GetHeader("Authorization"), id, user.(*base.User))
	if code == 0 {
		r.Success(nil)
	} else if code == 1 {
		glog.Errorf("[cmd] delete container error ! msg: %s\n", err.Error())
		r.Error(handle.CONTAINER_REMOVE_FAIL)
	} else if code == 2 {
		glog.Errorf("[db] delete container error ! msg: %s\n", err.Error())
		r.Error(handle.CONTAINER_REMOVE_FAIL)
	} else if code == 3 || code == 4 {
		glog.Errorf("[db] update user containers error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}

// 开启容器
func startContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	id := c.Param("id")
	if global.Config.App.Engine != "docker" {
		r.Error(handle.INTERNAL_ERROR)
		return
	}
	code, err := service.StartDockerContainer(id)
	if code == 0 {
		r.Success(nil)
	} else if code == 1 {
		glog.Errorf("[cmd] start container error ! msg: %s\n", err.Error())
		r.Error(handle.CONTAINER_START_FAIL)
	} else if code == 2 {
		glog.Errorf("[db] start container error ! msg: %s\n", err.Error())
		r.Error(handle.CONTAINER_START_FAIL)
	}
}

// 停止容器
func stopContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	id := c.Param("id")
	if global.Config.App.Engine != "docker" {
		r.Error(handle.INTERNAL_ERROR)
		return
	}
	code, err := service.StopDockerContainer(id)
	if code == 0 {
		r.Success(nil)
	} else if code == 1 {
		glog.Errorf("[cmd] stop container error ! msg: %s\n", err.Error())
		r.Error(handle.CONTAINER_STOP_FAIL)
	} else if code == 2 {
		glog.Errorf("[db] stop container error ! msg: %s\n", err.Error())
		r.Error(handle.CONTAINER_STOP_FAIL)
	}
}

// 重启容器
func restartContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	id := c.Param("id")
	if global.Config.App.Engine != "docker" {
		r.Error(handle.INTERNAL_ERROR)
		return
	}
	code, err := service.RestartDockerContainer(id)
	if code == 0 {
		r.Success(nil)
	} else if code == 1 {
		glog.Errorf("[cmd] restart container error ! msg: %s\n", err.Error())
		r.Error(handle.CONTAINER_RESTART_FAIL)
	} else if code == 2 {
		glog.Errorf("[db] restart container error ! msg: %s\n", err.Error())
		r.Error(handle.CONTAINER_RESTART_FAIL)
	}
}

// 根据容器制作镜像
func makeImage(c *gin.Context) {
	r := handle.NewResponse(c)
	image := &cloud.Image{}
	id := c.Param("id")
	c.ShouldBind(image)
	image.Name = fmt.Sprintf("%s/%s", global.Config.App.Name, image.Name)
	if id == "" || image.Name == "" || image.Desc == "" || image.Author == "" {
		r.Error(handle.PARAM_NOT_COMPLETE)
		return
	}
	var err error
	var code int8
	if global.Config.App.Engine == "docker" {
		code, err = service.MakeDockerImage(id, image)
	} else {
		code, err = service.MakeK8SImage(id, image)
	}
	fmt.Print(code)
	if code == 0 {
		r.Success(nil)
	} else if code == 1 {
		glog.Errorf("[cmd] make image from container error ! msg: %s\n", err.Error())
		r.Error(handle.IMAGE_CREATE_FAIL)
	} else if code == 2 {
		glog.Errorf("[db] make image from container error ! msg: %s\n", err.Error())
		r.Error(handle.IMAGE_CREATE_FAIL)
	} else if code == 3 {
		glog.Errorf("[cmd] make image fail, image has exist!")
		r.Error(handle.IMAGE_CREATE_FAIL)
	}
}

// 获取容器日志
func getContainerLog(c *gin.Context) {
	r := handle.NewResponse(c)
	id := c.Param("id")
	code, out, err := service.GetContainerLog(id)
	if code == 0 {
		r.Success(out)
	} else if code == 1 {
		glog.Errorf("[cmd] get container log error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}
