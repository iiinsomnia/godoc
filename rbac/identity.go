package rbac

import (
	"godoc/dao/mysql"
	"godoc/session"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

type Identity struct {
	ID            int       `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Email         string    `db:"email" json:"email"`
	Password      string    `db:"password" json:"password"`
	Salt          string    `db:"salt" json:"salt"`
	Role          int       `db:"role" json:"role"`
	LastLoginIP   string    `db:"last_login_ip" json:"last_login_ip"`
	LastLoginTime time.Time `db:"last_login_time" json:"last_login_time"`
}

// GetIdentity 获取用户登录信息
func GetIdentity(c *gin.Context) *Identity {
	identity := &Identity{}

	loginData, err := session.Get(c, "identity")

	if err == nil && loginData != nil {
		err = json.Unmarshal(loginData.([]byte), identity)

		if err != nil {
			yiigo.LogError(err.Error())
		}
	}

	return identity
}

// IsGuest 判断用户是否是Guest
func IsGuest(c *gin.Context) bool {
	loginData, err := session.Get(c, "identity")

	if err != nil || loginData == nil {
		return true
	}

	return false
}

// SignIn 用户登录
func SignIn(c *gin.Context, identity *Identity, duration ...int) bool {
	loginIP := c.ClientIP()
	loginTime := time.Now()

	userDao := mysql.NewUserDao()
	userDao.UpdateById(identity.ID, yiigo.X{
		"last_login_ip":   loginIP,
		"last_login_time": loginTime,
	})

	identity.LastLoginIP = loginIP
	identity.LastLoginTime = loginTime

	loginData, err := json.Marshal(identity)

	if err != nil {
		yiigo.LogError(err.Error())

		return false
	}

	err = session.Set(c, "identity", loginData, duration...)

	if err != nil {
		return false
	}

	return true
}
