package middleware

import (
	"cloud-platform/global"
	"cloud-platform/internal/base"
	"cloud-platform/internal/base/config"
	"cloud-platform/internal/handle"
	commonx "cloud-platform/internal/pkg/common"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"github.com/golang/glog"
)

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,x-token")
		context.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE,PATCH,PUT")
		context.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Allow-Headers,Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		context.Header("Access-Control-Expose-Headers", "Authorization")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
	}
}

func TokenVer() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := handle.NewResponse(c)
		token := c.GetHeader("Authorization")
		if token == "" {
			// 无权限
			r.Ctx.Header("WWW-Authenticate", "Bearer")
			r.Ctx.Status(401)
			c.Abort()
			return
		}
		/** 验证token是否合法与过期 */
		cmd := global.Rdb.Get(fmt.Sprintf(base.TOKEN, token))
		if cmd.Err() != nil {
			if cmd.Err() == redis.Nil {
				r.Error(handle.TOKEN_IS_EXPIRED)
				c.Abort()
				return
			}
			glog.Errorf("redis get token error ! msg: %s", cmd.Err().Error())
			c.Abort()
			return
		}
		user := &base.User{}
		commonx.JsonToStruct(cmd.Val(), user)
		cmd1 := global.Rdb.Expire(context.TODO(), fmt.Sprintf(base.TOKEN, token), 30*time.Minute)
		if cmd1.Err() != nil {
			glog.Errorf("token extension of time error ! msg: %s\n", cmd1.Err().Error())
			c.Abort()
			return
		}
		// 将 user 放到 context
		c.Set("user", user)
	}
}

func AdminVer() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user")
		fmt.Print(user.(*base.User).Pow)
		if user.(*base.User).Pow != config.ADMIN_POW {
			r := handle.NewResponse(c)
			r.Error(handle.INSUFFICENT_PERMISSIONS)
			c.Abort()
		}
	}
}

func ContainerVer() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Print("abcd")
		r := handle.NewResponse(c)
		containerId := c.Param("id")
		user, _ := c.Get("user")
		if containerId == "" || user == nil {
			r.Error(handle.PARAM_NOT_COMPLETE)
			c.Abort()
			return
		}
		// 查看此用户是否拥有此容器
		for _, v := range user.(*base.User).Containers {
			if v == containerId {
				return
			}
		}
		r.Error(handle.INSUFFICENT_PERMISSIONS)
		c.Abort()
	}
}
