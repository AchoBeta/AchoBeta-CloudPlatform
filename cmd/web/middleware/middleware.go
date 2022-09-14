package middleware

import (
	"CloudPlatform/pkg/handle"

	"github.com/gin-gonic/gin"
)

func TokenVer() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			// 无权限
			r := handle.NewResponse(c)
			r.Error(handle.USER_NOT_LOGIN)
			c.Abort()
			return
		} else {
			/** 权限检验逻辑 todo */
		}
	}
}
