package middlewares

import (
	"godoc/rbac"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

// AuthMiddleware 用户登录验证
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if rbac.IsGuest(c) {
			if yiigo.IsXhr(c) {
				c.JSON(http.StatusOK, gin.H{
					"success":  false,
					"msg":      "登录已过期",
					"data":     nil,
					"redirect": "/login",
				})
			} else {
				c.Redirect(http.StatusFound, "/login")
			}

			c.Abort()

			return
		}

		c.Next()
	}
}
