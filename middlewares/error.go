package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorMiddleware 用户登录验证
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Writer.Status() {
		case http.StatusNotFound:
			c.Redirect(http.StatusFound, "/404")
			c.Abort()

			return
		case http.StatusInternalServerError:
			c.Redirect(http.StatusFound, "/500")
			c.Abort()

			return
		}

		c.Next()
	}
}
