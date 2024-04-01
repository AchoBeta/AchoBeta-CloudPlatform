package router

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
)

func init() {
	RegisterRouter(func(r *route.RouterGroup) {
		r.GET("/test", Test)
	}, LEVEL_GLOBAL)
}

func Test(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{"ping": "pong"})
}
