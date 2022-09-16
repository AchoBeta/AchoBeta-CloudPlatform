package middleware

import (
	"CloudPlatform/pkg/handle"
<<<<<<< HEAD
	"CloudPlatform/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/golang/glog"
=======

	"github.com/gin-gonic/gin"
>>>>>>> master
)

func TokenVer() gin.HandlerFunc {
	return func(c *gin.Context) {
<<<<<<< HEAD
		r := handle.NewResponse(c)
		rdb := util.GetRDBClient()
		token := c.GetHeader("Authorization")
		if token == "" {
			// 无权限
=======
		token := c.GetHeader("Authorization")
		if token == "" {
			// 无权限
			r := handle.NewResponse(c)
>>>>>>> master
			r.Error(handle.USER_NOT_LOGIN)
			c.Abort()
			return
		} else {
<<<<<<< HEAD
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
			/** 续期 5分钟 */
			result, err := rdb.Set(token, uid, 5*60*time.Second).Result()
			if err != nil {
				glog.Errorf("token extension of time error ! msg: %s", err.Error())
				c.Abort()
				return
			}
			r.Success(result)
=======
			/** 权限检验逻辑 todo */
>>>>>>> master
		}
	}
}
