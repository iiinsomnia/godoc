package service

import (
	"database/sql"
	"godoc/rbac"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

type CategoryService struct {
	Identity *rbac.Identity
}

type Category struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewCategoryService(c *gin.Context) *CategoryService {
	return &CategoryService{
		Identity: rbac.GetIdentity(c),
	}
}

func (c *CategoryService) GetAll() ([]Category, error) {
	defer yiigo.Flush()

	data := []Category{}

	err := yiigo.DB.Select(&data, "SELECT * FROM go_category ORDER BY updated_at DESC")

	if err != nil && err != sql.ErrNoRows {
		yiigo.Err(err.Error())
	}

	return data, err
}

func (c *CategoryService) GetDetail(id int) (*Category, error) {
	defer yiigo.Flush()

	data := &Category{}

	err := yiigo.DB.Get(data, "SELECT * FROM go_category WHERE id = ?", id)

	if err != nil && err != sql.ErrNoRows {
		yiigo.Err(err.Error())
	}

	return data, err
}

func (c *CategoryService) Add(data yiigo.X) (int64, error) {
	defer yiigo.Flush()

	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()

	sql, binds := yiigo.InsertSQL("go_category", data)
	r, err := yiigo.DB.Exec(sql, binds...)

	if err != nil {
		yiigo.Err(err.Error())

		return 0, err
	}

	id, _ := r.LastInsertId()

	return id, err
}

func (c *CategoryService) Edit(id int, data yiigo.X) error {
	defer yiigo.Flush()

	data["updated_at"] = time.Now()

	sql, binds := yiigo.UpdateSQL("UPDATE go_category SET ? WHERE id = ?", data, id)

	_, err := yiigo.DB.Exec(sql, binds...)

	if err != nil {
		yiigo.Err(err.Error())

		return err
	}

	return err
}

func (c *CategoryService) Delete(id int) error {
	defer yiigo.Flush()

	tx, err := yiigo.DB.Beginx()

	if err != nil {
		yiigo.Err(err.Error())

		return err
	}

	_, err = tx.Exec("DELETE FROM go_history WHERE category_id = ?", id)

	if err != nil {
		tx.Rollback()
		yiigo.Err(err.Error())

		return err
	}

	_, err = tx.Exec("DELETE FROM go_doc WHERE category_id = ?", id)

	if err != nil {
		tx.Rollback()
		yiigo.Err(err.Error())

		return err
	}

	_, err = tx.Exec("DELETE FROM go_project WHERE category_id = ?", id)

	if err != nil {
		tx.Rollback()
		yiigo.Err(err.Error())

		return err
	}

	_, err = tx.Exec("DELETE FROM go_category WHERE id = ?", id)

	if err != nil {
		tx.Rollback()
		yiigo.Err(err.Error())

		return err
	}

	tx.Commit()

	return nil
}
