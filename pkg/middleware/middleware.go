package middleware

import (
	"CloudPlatform/pkg/handle"
	"CloudPlatform/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/golang/glog"
)

func TokenVer() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := handle.NewResponse(c)
		rdb := util.GetRDBClient()
		token := c.GetHeader("Authorization")
		if token == "" {
			// 无权限
			r.Error(handle.USER_NOT_LOGIN)
			c.Abort()
			return
		} else {
			/** 验证token是否合法与过期 */
			uid, err := rdb.Get(token).Result()
			if err != nil {
				if err == redis.Nil {
					r.Error(handle.TOKEN_IS_EXPIRED)
					c.Abort()
					return
				}
				glog.Errorf("redis get token error ! msg: %s", err.Error())
				c.Abort()
				return
			}
			/** 重置过期时间 30 分钟 */
			result, err := rdb.Set(token, uid, 30*60*time.Second).Result()
			if err != nil {
				glog.Errorf("token extension of time error ! msg: %s", err.Error())
				c.Abort()
				return
			}
			r.Success(result)
		}
	}
}

func ContainerVer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 校验 containerId —— userId
	}
}
