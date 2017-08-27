package service

import (
	"database/sql"
	"fmt"
	"godoc/models"
	"godoc/rbac"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

type User struct {
	Identity *models.Identity
}

func NewUser(c *gin.Context) *User {
	return &User{
		Identity: rbac.GetIdentity(c),
	}
}

func (u *User) GetUserList(query url.Values, size ...int) (int, int, int, []models.User, error) {
	defer yiigo.Flush()

	data := []models.User{}

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

	where := []string{}
	binds := []interface{}{}

	for k, v := range query {
		if strings.TrimSpace(v[0]) != "" {
			switch k {
			case "name":
				where = append(where, "name = ?")
				binds = append(binds, v[0])
			case "email":
				where = append(where, "email = ?")
				binds = append(binds, v[0])
			case "role":
				where = append(where, "role = ?")
				binds = append(binds, v[0])
			}
		}
	}

	countSQL := "SELECT COUNT(*) FROM go_user"
	querySQL := "SELECT id, name, email, role, last_login_ip, last_login_time, created_at, updated_at FROM go_user"

	if len(where) > 0 {
		countSQL = fmt.Sprintf("%s WHERE %s", countSQL, strings.Join(where, " AND "))
		querySQL = fmt.Sprintf("%s WHERE %s", querySQL, strings.Join(where, " AND "))
	}

	count := 0

	err := yiigo.DB.Get(&count, countSQL, binds...)

	if err != nil {
		yiigo.Errf("%s, SQL: %s, Args: %v", err.Error(), countSQL, binds)

		return 0, curPage, 0, data, err
	}

	binds = append(binds, offset, limit)
	querySQL += " LIMIT ?, ?"

	err = yiigo.DB.Select(&data, querySQL, binds...)

	if err != nil && err != sql.ErrNoRows {
		yiigo.Errf("%s, SQL: %s, Args: %v", err.Error(), querySQL, binds)
	}

	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	return count, curPage, totalPage, data, err
}

func (u *User) GetDetail(id int) (*models.User, error) {
	defer yiigo.Flush()

	data := &models.User{}

	query := "SELECT id, name, email, role, last_login_ip, last_login_time, created_at, updated_at FROM go_user WHERE id = ?"

	err := yiigo.DB.Get(data, query, id)

	if err != nil && err != sql.ErrNoRows {
		yiigo.Errf("%s, SQL: %s, Args: %d", err.Error(), query, id)
	}

	return data, err
}

func (u *User) Add(data yiigo.X) (int64, error) {
	defer yiigo.Flush()

	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()

	sql, binds := yiigo.InsertSQL("go_user", data)
	r, err := yiigo.DB.Exec(sql, binds...)

	if err != nil {
		yiigo.Errf("%s, SQL: %s, Args: %v", err.Error(), sql, binds)

		return 0, err
	}

	id, _ := r.LastInsertId()

	return id, nil
}

func (u *User) Edit(id int, data yiigo.X) error {
	defer yiigo.Flush()

	data["updated_at"] = time.Now()

	sql, binds := yiigo.UpdateSQL("UPDATE go_user SET ? WHERE id = ?", data, id)

	_, err := yiigo.DB.Exec(sql, binds...)

	if err != nil {
		yiigo.Errf("%s, SQL: %s, Args: %v", err.Error(), sql, binds)
	}

	return err
}

func (u *User) Delete(id int) error {
	defer yiigo.Flush()

	sql := "DELETE FROM go_user WHERE id = ?"
	_, err := yiigo.DB.Exec(sql, id)

	if err != nil {
		yiigo.Errf("%s, SQL: %s, Args: [%d]", err.Error(), sql, id)
	}

	return err
}

func (u *User) CheckUnique(name string, email string, id ...int) (bool, error) {
	defer yiigo.Flush()

	data := &models.User{}
	binds := []interface{}{name, email}

	query := "SELECT id FROM go_user WHERE name = ? OR email = ?"

	if len(id) > 0 {
		query = "SELECT id FROM go_user WHERE id <> ? AND (name = ? OR email = ?)"
		binds = append(binds, id[0])
	}

	err := yiigo.DB.Get(data, query, binds...)

	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		yiigo.Errf("%s, SQL: %s, Args: %v", err.Error(), query, binds)

		return false, err
	}

	return false, nil
}
