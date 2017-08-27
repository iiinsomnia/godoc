package service

import (
	"database/sql"
	"godoc/rbac"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

type ProjectService struct {
	Identity *rbac.Identity
}

type Project struct {
	ID           int       `db:"id"`
	Name         string    `db:"name"`
	CategoryID   int       `db:"category_id"`
	CategoryName string    `db:"category_name"`
	Description  string    `db:"description"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func NewProjectService(c *gin.Context) *ProjectService {
	return &ProjectService{
		Identity: rbac.GetIdentity(c),
	}
}

func (p *ProjectService) GetProjects(categoryID int) ([]Project, error) {
	defer yiigo.Flush()

	data := []Project{}

	err := yiigo.DB.Select(&data, "SELECT * FROM go_project WHERE category_id = ? ORDER BY updated_at DESC", categoryID)

	if err != nil && err != sql.ErrNoRows {
		yiigo.Err(err.Error())
	}

	return data, err
}

func (p *ProjectService) GetDetail(id int) (*Project, error) {
	defer yiigo.Flush()

	data := &Project{}

	query := "SELECT a.*, b.name AS category_name FROM go_project AS a LEFT JOIN go_category AS b ON a.category_id = b.id WHERE a.id = ?"
	err := yiigo.DB.Get(data, query, id)

	if err != nil && err != sql.ErrNoRows {
		yiigo.Err(err.Error())
	}

	return data, err
}

func (p *ProjectService) Add(data yiigo.X) (int64, error) {
	defer yiigo.Flush()

	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()

	sql, binds := yiigo.InsertSQL("go_project", data)
	r, err := yiigo.DB.Exec(sql, binds...)

	if err != nil {
		yiigo.Err(err.Error())

		return 0, err
	}

	id, _ := r.LastInsertId()

	return id, nil
}

func (p *ProjectService) Edit(id int, data yiigo.X) error {
	defer yiigo.Flush()

	data["updated_at"] = time.Now()

	sql, binds := yiigo.UpdateSQL("UPDATE go_project SET ? WHERE id = ?", data, id)

	_, err := yiigo.DB.Exec(sql, binds...)

	if err != nil {
		yiigo.Err(err.Error())
	}

	return err
}

func (c *ProjectService) Delete(id int) error {
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

	_, err = tx.Exec("DELETE FROM go_project WHERE id = ?", id)

	if err != nil {
		tx.Rollback()
		yiigo.Err(err.Error())

		return err
	}

	tx.Commit()

	return nil
}
