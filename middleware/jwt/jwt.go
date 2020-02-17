package jwt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wudidayan/go-blog/pkg/app"
	"github.com/wudidayan/go-blog/pkg/e"
	"github.com/wudidayan/go-blog/pkg/util"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		appResp := app.Gin{c}
		code := e.SUCCESS
		data := make(map[string]string)

		token := c.Query("token")
		if token == "" {
			code = e.INVALID_PARAMS
			appResp.Response(http.StatusOK, code, data)
			c.Abort()
			return
		}

		claims, err := util.ParseToken(token)
		if err != nil {
			code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
		}

		if code != e.SUCCESS {
			appResp.Response(http.StatusOK, code, data)
			c.Abort()
			return
		}

		c.Next()
	}
}
