package service

import (
	"database/sql"
	"godoc/models"
	"godoc/rbac"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

type Project struct {
	Identity *models.Identity
}

func NewProject(c *gin.Context) *Project {
	return &Project{
		Identity: rbac.GetIdentity(c),
	}
}

func (p *Project) GetProjects(categoryID int) ([]models.Project, error) {
	defer yiigo.Flush()

	data := []models.Project{}

	query := "SELECT * FROM go_project WHERE category_id = ? ORDER BY updated_at DESC"
	err := yiigo.DB.Select(&data, query, categoryID)

	if err != nil && err != sql.ErrNoRows {
		yiigo.Errf("%s, SQL: %s, Args: [%d]", err.Error(), query, categoryID)
	}

	return data, err
}

func (p *Project) GetDetail(id int) (*models.Project, error) {
	defer yiigo.Flush()

	data := &models.Project{}

	query := "SELECT a.*, b.name AS category_name FROM go_project AS a LEFT JOIN go_category AS b ON a.category_id = b.id WHERE a.id = ?"
	err := yiigo.DB.Get(data, query, id)

	if err != nil && err != sql.ErrNoRows {
		yiigo.Errf("%s, SQL: %s, Args: [%d]", err.Error(), query, id)
	}

	return data, err
}

func (p *Project) Add(data yiigo.X) (int64, error) {
	defer yiigo.Flush()

	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()

	sql, binds := yiigo.InsertSQL("go_project", data)
	r, err := yiigo.DB.Exec(sql, binds...)

	if err != nil {
		yiigo.Errf("%s, SQL: %s, Args: %v", err.Error(), sql, binds)

		return 0, err
	}

	id, _ := r.LastInsertId()

	return id, nil
}

func (p *Project) Edit(id int, data yiigo.X) error {
	defer yiigo.Flush()

	data["updated_at"] = time.Now()

	sql, binds := yiigo.UpdateSQL("UPDATE go_project SET ? WHERE id = ?", data, id)

	_, err := yiigo.DB.Exec(sql, binds...)

	if err != nil {
		yiigo.Errf("%s, SQL: %s, Args: %v", err.Error(), sql, binds)
	}

	return err
}

func (c *Project) Delete(id int) error {
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

	sql = "DELETE FROM go_project WHERE id = ?"
	_, err = tx.Exec(sql, id)

	if err != nil {
		tx.Rollback()
		yiigo.Errf("%s, SQL: %s, Args: [%d]", err.Error(), sql, id)

		return err
	}

	tx.Commit()

	return nil
}
