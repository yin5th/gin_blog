package jwt

import (
	"gin_blog/pkg/e"
	"gin_blog/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.GetHeader("token")

		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			//解析token
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				//token过期
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": e.GetMsg(code),
				"data":    data,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
