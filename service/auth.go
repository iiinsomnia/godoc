package service

import (
	"database/sql"
	"errors"

	"godoc/rbac"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

type AuthService struct {
	Identity *rbac.Identity
}

func NewAuthService(c *gin.Context) *AuthService {
	return &AuthService{
		Identity: rbac.GetIdentity(c),
	}
}

// Login 用户登录
func (a *AuthService) Login(c *gin.Context, account string, password string) error {
	identity := &rbac.Identity{}

	query := "SELECT id, name, email, password, salt, role, last_login_ip, last_login_time FROM go_user WHERE name = ? OR email = ?"
	err := yiigo.DB.Get(identity, query, account, account)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("帐号不存在")
		}

		yiigo.Err(err.Error())

		return errors.New("登录失败")
	}

	if yiigo.MD5(password+identity.Salt) != identity.Password {
		return errors.New("帐号或密码错误")
	}

	success := rbac.SignIn(c, identity, 12*3600)

	if !success {
		return errors.New("登录失败")
	}

	return nil
}
