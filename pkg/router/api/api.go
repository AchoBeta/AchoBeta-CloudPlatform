package api

import (
	router "cloud-platform/pkg/router/manager"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
)

func init() {
	router.RouteHandler.RegisterRouter(router.LEVEL_GLOBAL, func(r *route.RouterGroup) {
		r.GET("/test", Test)
	})
}

func Test(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{"ping": "pong"})
}
