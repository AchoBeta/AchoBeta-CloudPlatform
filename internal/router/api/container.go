package api

import (
	"CloudPlatform/global"
	"CloudPlatform/internal/base"
	"CloudPlatform/internal/handle"
	"CloudPlatform/internal/router"
	"CloudPlatform/internal/service"
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
		router.POST("/containers/:id/makeImage", makeImage)
		router.GET("/containers/:id/log", getContainerLog)
	}, router.V3)
}

// 创建容器
func createContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	container := &base.Container{}
	err := c.BindJSON(&container)
	user, _ := c.Get("user")
	if err != nil || container.Image == "" || container.Name == "" {
		r.Error(handle.PARAM_NOT_VALID)
		return
	}
	var code int8
	if global.Config.App.Type == "docker" {
		err, code = service.CreateDockerContainer(container, user.(*base.User))
	} else {
		err, code = service.CreateK8SContainer(container)
	}
	if code == 0 {
		r.Success(nil)
	} else if code == 1 {
		glog.Errorf("[cmd] run container error ! msg: %s\n", err.Error())
		r.Error(handle.CONTAINER_CREATE_FAIL)
		return
	} else if code == 2 {
		glog.Errorf("[cmd] set container ssh pwd error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	} else if code == 3 {
		glog.Errorf("[db] insert container error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}

// 根据 id 获取容器信息
func getContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	id := c.Param("id")
	container := &base.Container{}
	err, code := service.GetContainer(id, container)
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
	err, code, containers := service.GetContainers(user.(*base.User))
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
	var err error
	var code int8
	if global.Config.App.Type == "docker" {
		err, code = service.RemoveDockerContainer(id, user.(*base.User).Id)
	} else {
		err, code = service.RemoveK8SContainer(id)
	}
	if code == 0 {
		r.Success(nil)
	} else if code == 1 {
		glog.Errorf("[cmd] delete container error ! msg: %s\n", err.Error())
		r.Error(handle.CONTAINER_REMOVE_FAIL)
	} else if code == 2 {
		glog.Errorf("[db] delete container error ! msg: %s\n", err.Error())
		r.Error(handle.CONTAINER_REMOVE_FAIL)
	} else if code == 3 {
		glog.Errorf("[db] update user containers error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}

// 开启容器
func startContainer(c *gin.Context) {
	r := handle.NewResponse(c)
	id := c.Param("id")
	var err error
	var code int8
	if global.Config.App.Type == "docker" {
		err, code = service.StartDockerContainer(id)
	} else {
		err, code = service.StartK8SContainer(id)
	}
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
	var err error
	var code int8
	if global.Config.App.Type == "docker" {
		err, code = service.StopDockerContainer(id)
	} else {
		err, code = service.StopK8SContainer(id)
	}
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
	var err error
	var code int8
	if global.Config.App.Type == "docker" {
		err, code = service.RestartDockerContainer(id)
	} else {
		err, code = service.RestartK8SContainer(id)
	}
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
	image := &base.Image{}
	id := c.Param("id")
	c.BindJSON(image)
	image.Name = fmt.Sprintf("%s/%s", global.Config.App.Name, image.Name)
	if id == "" || image.Name == "" || image.Desc == "" || image.Author == "" {
		r.Error(handle.PARAM_NOT_COMPLETE)
		return
	}
	var err error
	var code int8
	if global.Config.App.Type == "docker" {
		err, code = service.MakeDockerImage(id, image)
	} else {
		err, code = service.MakeK8SImage(id, image)
	}
	if code == 0 {
		r.Success(nil)
	} else if code == 1 {
		glog.Errorf("[cmd] make image from container error ! msg: %s\n", err.Error())
		r.Error(handle.IMAGE_CREATE_FAIL)
	} else if code == 2 {
		glog.Errorf("[db] make image from container error ! msg: %s\n", err.Error())
		r.Error(handle.IMAGE_CREATE_FAIL)
	}
}

// 获取容器日志
func getContainerLog(c *gin.Context) {
	r := handle.NewResponse(c)
	id := c.Param("id")
	var err error
	var code int8
	var out string
	if global.Config.App.Type == "docker" {
		err, code, out = service.GetDockerContainerLog(id)
	} else {
		err, code, out = service.GetK8SContainerLog(id)
	}
	if code == 0 {
		r.Success(out)
	} else if code == 1 {
		glog.Errorf("[cmd] get container log error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}
