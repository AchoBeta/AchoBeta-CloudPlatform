package middleware

import (
	"cloud-platform/global"
	"cloud-platform/pkg/base"
	"cloud-platform/pkg/base/config"
	"cloud-platform/pkg/handle"
	commonx "cloud-platform/pkg/handle/common"
	"cloud-platform/pkg/load/tlog"
	"cloud-platform/pkg/router/manager"
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

func init() {
	manager.RouteHandler.RegisterMiddleware(manager.LEVEL_GLOBAL, AddTraceId, false)
}

func TokenVer() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		r := handle.NewResponse(ctx)
		token := ctx.GetHeader("Authorization")
		if token == nil {
			// 无权限
			r.Ctx.Header("WWW-Authenticate", "Basic")
			r.Ctx.Status(401)
			ctx.Abort()
			return
		}
		/** 验证token是否合法与过期 */
		cmd := global.Rdb.Get(fmt.Sprintf(base.TOKEN, token))
		if cmd.Err() != nil {
			if cmd.Err() == redis.Nil {
				r.Error(handle.TOKEN_IS_EXPIRED)
				ctx.Abort()
				return
			}
			tlog.Errorf("redis get token error ! msg: %s", cmd.Err().Error())
			ctx.Abort()
			return
		}
		user := &base.User{}
		commonx.JsonToStruct(cmd.Val(), user)
		cmd1 := global.Rdb.Expire(fmt.Sprintf(base.TOKEN, token), 30*time.Minute)
		if cmd1.Err() != nil {
			tlog.Errorf("token extension of time error ! msg: %s\n", cmd1.Err().Error())
			ctx.Abort()
			return
		}
		// 将 user 放到 context
		ctx.Set("user", user)
	}
}

func AdminVer() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		user, _ := ctx.Get("user")
		if user.(*base.User).Pow != config.ADMIN_POW {
			r := handle.NewResponse(ctx)
			r.Error(handle.INSUFFICENT_PERMISSIONS)
			ctx.Abort()
		}
	}
}

func ContainerVer() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		r := handle.NewResponse(ctx)
		containerId := ctx.Param("id")
		user, _ := ctx.Get("user")
		if containerId == "" || user == nil {
			r.Error(handle.PARAM_NOT_COMPLETE)
			ctx.Abort()
			return
		}
		// 查看此用户是否拥有此容器
		for _, v := range user.(*base.User).Containers {
			if v == containerId {
				return
			}
		}
		r.Error(handle.INSUFFICENT_PERMISSIONS)
		ctx.Abort()
	}
}

func AddTraceId() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 假设 Trace ID 存在于 HTTP Header "X-Trace-ID" 中
		traceID := ctx.Request.Header.Get("X-Request-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c = tlog.NewContext(c, zap.String("traceId", traceID))
		ctx.Next(c)
	}
}
