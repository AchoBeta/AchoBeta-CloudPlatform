package api

import (
	"cloud-platform/global"
	"cloud-platform/pkg/base"
	"cloud-platform/pkg/base/cloud"
	"cloud-platform/pkg/handle"
	"cloud-platform/pkg/router/manager"
	"cloud-platform/pkg/service"
	"context"

	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/golang/glog"
)

func init() {
	manager.RouteHandler.RegisterRouter(manager.LEVEL_V1, func(router *route.RouterGroup) {
		router.POST("/containers", createContainer)
		router.GET("/containers", getContainers)
	})

	manager.RouteHandler.RegisterRouter(manager.LEVEL_V3, func(router *route.RouterGroup) {
		router.GET("/containers/:id", getContainer)
		router.DELETE("/containers/:id", removeContainer)
		router.GET("/containers/:id/start", startContainer)
		router.GET("/containers/:id/stop", stopContainer)
		router.GET("/containers/:id/restart", restartContainer)
		router.POST("/containers/:id/makeImage", makeImage)
		router.GET("/containers/:id/log", getContainerLog)
	})
}

// 创建容器
func createContainer(ctx context.Context, c *app.RequestContext) {
	r := handle.NewResponse(c)
	container := &cloud.Container{}
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
func getContainer(ctx context.Context, c *app.RequestContext) {
	r := handle.NewResponse(c)
	id := c.Param("id")
	container := &cloud.Container{}
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
func getContainers(ctx context.Context, c *app.RequestContext) {
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
func removeContainer(ctx context.Context, c *app.RequestContext) {
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
func startContainer(ctx context.Context, c *app.RequestContext) {
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
func stopContainer(ctx context.Context, c *app.RequestContext) {
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
func restartContainer(ctx context.Context, c *app.RequestContext) {
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
func makeImage(ctx context.Context, c *app.RequestContext) {
	r := handle.NewResponse(c)
	image := &cloud.Image{}
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
func getContainerLog(ctx context.Context, c *app.RequestContext) {
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
