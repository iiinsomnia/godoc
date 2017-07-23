package mysql

import "github.com/iiinsomnia/yiigo"

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

func (u *UserDao) GetAll(data interface{}) error {
	err := u.MySQL.FindAll(data)

	if err != nil {
		yiigo.LogError(err.Error())
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
