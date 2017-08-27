package rbac

import (
	"encoding/json"
	"godoc/session"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

var Roles = map[int]string{
	1: "普通用户",
	2: "高级用户",
	3: "超级用户",
}

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

// GetRoleName 获取角色名称
func GetRoleName(role int) string {
	if v, ok := Roles[role]; ok {
		return v
	}

	return ""
}

// GetIdentity 获取用户登录信息
func GetIdentity(c *gin.Context) *Identity {
	identity := &Identity{}

	loginData, err := session.Get(c, "identity")

	if err == nil && loginData != nil {
		err = json.Unmarshal(loginData.([]byte), identity)

		if err != nil {
			yiigo.Err(err.Error())
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

	sql, binds := yiigo.UpdateSQL("UPDATE go_user SET ? WHERE id = ?", yiigo.X{
		"last_login_ip":   loginIP,
		"last_login_time": loginTime,
	}, identity.ID)

	yiigo.DB.Exec(sql, binds...)

	identity.LastLoginIP = loginIP
	identity.LastLoginTime = loginTime

	loginData, err := json.Marshal(identity)

	if err != nil {
		yiigo.Err(err.Error())

		return false
	}

	err = session.Set(c, "identity", loginData, duration...)

	if err != nil {
		return false
	}

	return true
}

// GenerateSalt 生成随机加密盐
func GenerateSalt() string {
	salt := []string{}
	pattern := "abcdef!ghijklm@nopqrst#uvwxyz$12345%67890^ABCDEFGH&IJKLMNOP*QRSTUVWXYZ"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := len(pattern)

	for i := 0; i < 16; i++ {
		n := r.Intn(length)
		salt = append(salt, pattern[n:n+1])
	}

	return strings.Join(salt, "")
}
