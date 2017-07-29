package service

import (
	"errors"
	"godoc/dao/mysql"
	"godoc/helpers"

	"godoc/rbac"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	*service
}

func NewAuthService(c *gin.Context) *AuthService {
	return &AuthService{
		construct(c),
	}
}

// Login 用户登录
func (a *AuthService) Login(c *gin.Context, account string, password string) error {
	userDao := mysql.NewUserDao()
	identity := &rbac.Identity{}

	err := userDao.GetByAccount(account, identity)

	if err != nil {
		if err.Error() == "not found" {
			return errors.New("帐号不存在")
		}

		return errors.New("登录失败")
	}

	if helpers.MD5(password+identity.Salt) != identity.Password {
		return errors.New("帐号或密码错误")
	}

	success := rbac.SignIn(c, identity, 12*3600)

	if !success {
		return errors.New("登录失败")
	}

	return nil
}
