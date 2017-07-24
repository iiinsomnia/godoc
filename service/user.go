package service

import (
	"godoc/dao/mysql"
	"math"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	*service
}

type User struct {
	ID            int       `db:"id"`
	Name          string    `db:"name"`
	Email         string    `db:"email"`
	Role          int       `db:"role"`
	LastLoginIP   string    `db:"last_login_ip"`
	LastLoginTime string    `db:"last_login_time"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func NewUserService(c *gin.Context) *UserService {
	return &UserService{
		construct(c),
	}
}

func (u *UserService) GetUserList(query url.Values, size ...int) (int, int, int, []User, error) {
	curPage := 1
	limit := 10

	if len(size) > 0 {
		limit = size[0]
	}

	offset := 0

	if v, ok := query["page"]; ok {
		curPage, _ = strconv.Atoi(v[0])
		offset = (curPage - 1) * limit
	}

	data := []User{}

	userDao := mysql.NewUserDao()
	count, err := userDao.GetByPagination(query, limit, offset, &data)

	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	return count, curPage, totalPage, data, err
}
