package service

import (
	"godoc/rbac"

	"github.com/gin-gonic/gin"
)

type service struct {
	Identity *rbac.Identity
}

// 构造函数
func construct(c *gin.Context) *service {
	return &service{
		Identity: rbac.GetIdentity(c),
	}
}
