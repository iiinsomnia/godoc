package service

import (
	"database/sql"
	"godoc/models"
	"godoc/rbac"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

type Category struct {
	Identity *models.Identity
}

func NewCategory(c *gin.Context) *Category {
	return &Category{
		Identity: rbac.GetIdentity(c),
	}
}

func (c *Category) GetAll() ([]models.Category, error) {
	defer yiigo.Flush()

	data := []models.Category{}

	query := "SELECT * FROM go_category ORDER BY updated_at DESC"
	err := yiigo.DB.Select(&data, query)

	if err != nil && err != sql.ErrNoRows {
		yiigo.Errf("%s, SQL: %s", err.Error(), query)
	}

	return data, err
}

func (c *Category) GetDetail(id int) (*models.Category, error) {
	defer yiigo.Flush()

	data := &models.Category{}

	query := "SELECT * FROM go_category WHERE id = ?"
	err := yiigo.DB.Get(data, query, id)

	if err != nil && err != sql.ErrNoRows {
		yiigo.Errf("%s, SQL: %s, Args: [%d]", err.Error(), query, id)
	}

	return data, err
}

func (c *Category) Add(data yiigo.X) (int64, error) {
	defer yiigo.Flush()

	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()

	sql, binds := yiigo.InsertSQL("go_category", data)
	r, err := yiigo.DB.Exec(sql, binds...)

	if err != nil {
		yiigo.Errf("%s, SQL: %s, Args: %v", err.Error(), sql, binds)

		return 0, err
	}

	id, _ := r.LastInsertId()

	return id, err
}

func (c *Category) Edit(id int, data yiigo.X) error {
	defer yiigo.Flush()

	data["updated_at"] = time.Now()

	sql, binds := yiigo.UpdateSQL("UPDATE go_category SET ? WHERE id = ?", data, id)

	_, err := yiigo.DB.Exec(sql, binds...)

	if err != nil {
		yiigo.Errf("%s, SQL: %s, Args: %v", err.Error(), sql, binds)

		return err
	}

	return err
}

func (c *Category) Delete(id int) error {
	defer yiigo.Flush()

	tx, err := yiigo.DB.Beginx()

	if err != nil {
		yiigo.Err(err.Error())

		return err
	}

	sql := "DELETE FROM go_history WHERE category_id = ?"
	_, err = tx.Exec(sql, id)

	if err != nil {
		tx.Rollback()
		yiigo.Errf("%s, SQL: %s, Args: [%d]", err.Error(), sql, id)

		return err
	}

	sql = "DELETE FROM go_doc WHERE category_id = ?"
	_, err = tx.Exec(sql, id)

	if err != nil {
		tx.Rollback()
		yiigo.Errf("%s, SQL: %s, Args: [%d]", err.Error(), sql, id)

		return err
	}

	sql = "DELETE FROM go_project WHERE category_id = ?"
	_, err = tx.Exec(sql, id)

	if err != nil {
		tx.Rollback()
		yiigo.Errf("%s, SQL: %s, Args: [%d]", err.Error(), sql, id)

		return err
	}

	sql = "DELETE FROM go_category WHERE id = ?"
	_, err = tx.Exec(sql, id)

	if err != nil {
		tx.Rollback()
		yiigo.Errf("%s, SQL: %s, Args: [%d]", err.Error(), sql, id)

		return err
	}

	tx.Commit()

	return nil
}
