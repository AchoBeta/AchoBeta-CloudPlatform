package router

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func init() {
	RegisterRouter(func(h *server.Hertz) {
		h.GET("/unverity_users", Test)
	})
}

func Test(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{"ping": "pong"})
}
