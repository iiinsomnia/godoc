package middlewares

import "github.com/gin-gonic/gin"

// RBACMiddleware 用户rbac权限验证
func RBACMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
