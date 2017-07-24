package mysql

import (
	"net/url"
	"strings"

	"github.com/iiinsomnia/yiigo"
)

type UserDao struct {
	yiigo.MySQL
}

func NewUserDao() *UserDao {
	return &UserDao{
		yiigo.MySQL{
			Table: "user",
		},
	}
}

func (u *UserDao) GetByPagination(query url.Values, limit int, offset int, data interface{}) (int, error) {
	where := []string{}
	binds := []interface{}{}

	for k, v := range query {
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

	err := u.MySQL.Find(yiigo.X{
		"select": "id, name, email, role, last_login_ip, last_login_time, created_at, updated_at",
		"where":  strings.Join(where, " AND "),
		"binds":  binds,
		"limit":  limit,
		"offset": offset,
	}, data)

	if err != nil {
		yiigo.LogError(err.Error())
		return 0, err
	}

	count, err := u.MySQL.Count(yiigo.X{
		"where": strings.Join(where, " AND "),
		"binds": binds,
	})

	if err != nil {
		yiigo.LogError(err.Error())
	}

	return count, nil
}

func (u *UserDao) GetById(id int, data interface{}) error {
	query := yiigo.X{
		"where": "id = ?",
		"binds": []interface{}{id},
	}

	err := u.MySQL.FindOne(query, data)

	if err != nil {
		if err.Error() != "not found" {
			yiigo.LogError(err.Error())
		}

		return err
	}

	return nil
}

func (u *UserDao) GetByAccount(account string, data interface{}) error {
	query := yiigo.X{
		"select": "id, name, email, password, salt, role, last_login_ip, last_login_time",
		"where":  "name = ? OR email = ?",
		"binds":  []interface{}{account, account},
	}

	err := u.MySQL.FindOne(query, data)

	if err != nil {
		if err.Error() != "not found" {
			yiigo.LogError(err.Error())
		}

		return err
	}

	return nil
}

func (u *UserDao) AddNewRecord(data yiigo.X) (int64, error) {
	id, err := u.MySQL.Insert(data)

	if err != nil {
		yiigo.LogError(err.Error())
		return 0, err
	}

	return id, nil
}

func (u *UserDao) UpdateById(id int, data yiigo.X) error {
	query := yiigo.X{
		"where": "id = ?",
		"binds": []interface{}{id},
	}

	_, err := u.MySQL.Update(query, data)

	if err != nil {
		yiigo.LogError(err.Error())
		return err
	}

	return nil
}

func (u *UserDao) DeleteById(id int) error {
	query := yiigo.X{
		"where": "id = ?",
		"binds": []interface{}{id},
	}

	_, err := u.MySQL.Delete(query)

	if err != nil {
		yiigo.LogError(err.Error())
		return err
	}

	return nil
}
