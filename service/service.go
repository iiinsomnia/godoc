package service

import (
	"crypto/md5"
	"fmt"
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

func (s *service) md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))

	return fmt.Sprintf("%x", h.Sum(nil))
}
