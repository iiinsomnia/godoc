package service

import (
	"godoc/rbac"
	"crypto/md5"
	"fmt"

	"github.com/gin-gonic/gin"
)

type service struct {
	User *rbac.Identity
}

// 构造函数
func construct(c *gin.Context) *service {
	return &service{
		User: rbac.GetIdentity(c),
	}
}

func (s *service) md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))

	return fmt.Sprintf("%x", h.Sum(nil))
}
