package api

import (
	"cloud-platform/pkg/load/tlog"
	router "cloud-platform/pkg/router/manager"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	lgr "github.com/sirupsen/logrus"
)

func init() {
	router.RouteHandler.RegisterRouter(router.LEVEL_GLOBAL, func(r *route.RouterGroup) {
		r.GET("/test2", Test2)
		r.GET("/test", Test)
	})
}

func Test(ctx context.Context, c *app.RequestContext) {
	tlog.Infof("load - test")
	tlog.CtxInfof(ctx, "load ctx info - test")
	tlog.CtxWarnf(ctx, "load ctx warn - test")
	tlog.CtxErrorf(ctx, "load ctx error - test")
	c.JSON(consts.StatusOK, utils.H{"ping": "pong"})
}

func Test2(ctx context.Context, c *app.RequestContext) {
	lgr.Info("logrus - test2222222222222xxx")
	lgr.Warn("logrus - warn22222222222")
	c.JSON(consts.StatusOK, utils.H{"ping": "pong"})
}
