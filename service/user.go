package service

import (
	"godoc/dao/mysql"
	"math"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
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
	LastLoginTime time.Time `db:"last_login_time"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func NewUserService(c *gin.Context) *UserService {
	return &UserService{
		construct(c),
	}
}

func (u *UserService) GetUserList(query url.Values, size ...int) (int, int, int, []User, error) {
	data := []User{}

	curPage := 1
	limit := 10

	if len(size) > 0 {
		limit = size[0]
	}

	offset := 0

	if v, ok := query["page"]; ok {
		curPage, _ = strconv.Atoi(v[0])

		if curPage < 1 {
			return 0, curPage, 0, data, nil
		}

		offset = (curPage - 1) * limit
	}

	userDao := mysql.NewUserDao()
	count, err := userDao.GetByPagination(query, limit, offset, &data)

	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	return count, curPage, totalPage, data, err
}

func (u *UserService) GetDetail(id int) (*User, error) {
	data := &User{}

	userDao := mysql.NewUserDao()
	err := userDao.GetByID(id, data)

	return data, err
}

func (u *UserService) Add(data yiigo.X) (int64, error) {
	userDao := mysql.NewUserDao()
	id, err := userDao.AddNewRecord(data)

	return id, err
}

func (u *UserService) Edit(id int, data yiigo.X) error {
	userDao := mysql.NewUserDao()
	err := userDao.UpdateByID(id, data)

	return err
}

func (u *UserService) Delete(id int) error {
	userDao := mysql.NewUserDao()
	err := userDao.DeleteByID(id)

	return err
}
