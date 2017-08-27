package service

import (
	"database/sql"
	"godoc/models"
	"godoc/rbac"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

type Doc struct {
	Identity *models.Identity
}

func NewDoc(c *gin.Context) *Doc {
	return &Doc{
		Identity: rbac.GetIdentity(c),
	}
}

func (d *Doc) GetDocs(projectID int) ([]models.Doc, error) {
	defer yiigo.Flush()

	data := []models.Doc{}

	query := "SELECT * FROM go_doc WHERE project_id = ? ORDER BY updated_at DESC"
	err := yiigo.DB.Select(&data, query, projectID)

	if err != nil && err != sql.ErrNoRows {
		yiigo.Errf("%s, SQL: %s, Args: [%d]", err.Error(), query, projectID)
	}

	return data, err
}

func (d *Doc) GetDetail(id int) (*models.Doc, error) {
	defer yiigo.Flush()

	data := &models.Doc{}

	query := "SELECT a.*, b.name AS category_name, c.name AS project_name FROM go_doc AS a LEFT JOIN go_category AS b ON a.category_id = b.id LEFT JOIN go_project AS c ON a.project_id = c.id WHERE a.id = ?"
	err := yiigo.DB.Get(data, query, id)

	if err != nil && err != sql.ErrNoRows {
		yiigo.Errf("%s, SQL: %s, Args: [%d]", err.Error(), query, id)
	}

	return data, err
}

func (d *Doc) Add(data yiigo.X, history yiigo.X) (int64, error) {
	defer yiigo.Flush()

	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()

	tx, err := yiigo.DB.Beginx()

	if err != nil {
		yiigo.Err(err.Error())

		return 0, err
	}

	sql, binds := yiigo.InsertSQL("go_doc", data)
	r, err := yiigo.DB.Exec(sql, binds...)

	if err != nil {
		tx.Rollback()
		yiigo.Errf("%s, SQL: %s, Args: %v", err.Error(), sql, binds)

		return 0, err
	}

	id, _ := r.LastInsertId()

	// 记录操作历史
	history["user_id"] = d.Identity.ID
	history["doc_id"] = id
	history["flag"] = 1
	history["updated_at"] = time.Now()

	sql, binds = yiigo.InsertSQL("go_history", history)
	_, err = yiigo.DB.Exec(sql, binds...)

	if err != nil {
		tx.Rollback()
		yiigo.Errf("%s, SQL: %s, Args: %v", err.Error(), sql, binds)

		return 0, err
	}

	tx.Commit()

	return id, nil
}

func (d *Doc) Edit(id int, data yiigo.X) error {
	defer yiigo.Flush()

	doc := &models.Doc{}
	err := yiigo.DB.Get(doc, "SELECT id FROM go_doc WHERE id = ?", id)

	if err != nil {
		return err
	}

	tx, err := yiigo.DB.Beginx()

	if err != nil {
		yiigo.Err(err.Error())

		return err
	}

	sql, binds := yiigo.UpdateSQL("UPDATE go_doc SET ? WHERE id = ?", data, id)
	_, err = tx.Exec(sql, binds...)

	if err != nil {
		tx.Rollback()
		yiigo.Errf("%s, SQL: %s, Args: %v", err.Error(), sql, binds)

		return err
	}

	// 记录操作历史
	history := yiigo.X{
		"user_id":     d.Identity.ID,
		"category_id": doc.CategoryID,
		"project_id":  doc.ProjectID,
		"doc_id":      id,
		"flag":        2,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}

	sql, binds = yiigo.InsertSQL("go_history", history)
	_, err = tx.Exec(sql, binds...)

	if err != nil {
		tx.Rollback()
		yiigo.Errf("%s, SQL: %s, Args: %v", err.Error(), sql, binds)

		return err
	}

	tx.Commit()

	return nil
}

func (c *Doc) Delete(id int) error {
	defer yiigo.Flush()

	tx, err := yiigo.DB.Beginx()

	if err != nil {
		yiigo.Err(err.Error())

		return err
	}

	sql := "DELETE FROM go_history WHERE doc_id = ?"
	_, err = tx.Exec(sql, id)

	if err != nil {
		tx.Rollback()
		yiigo.Errf("%s, SQL: %s, Args: [%d]", err.Error(), sql, id)

		return err
	}

	sql = "DELETE FROM go_doc WHERE id = ?"
	_, err = tx.Exec(sql, id)

	if err != nil {
		tx.Rollback()
		yiigo.Errf("%s, SQL: %s, Args: [%d]", err.Error(), sql, id)

		return err
	}

	tx.Commit()

	return nil
}
