package service

import (
	"database/sql"
	"godoc/models"
	"godoc/rbac"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

type History struct {
	Identity *models.Identity
}

func NewHistory(c *gin.Context) *History {
	return &History{
		Identity: rbac.GetIdentity(c),
	}
}

func (h *History) GetHistory(docID int) ([]models.History, error) {
	defer yiigo.Flush()

	data := []models.History{}

	query := "SELECT a.id, a.flag, a.created_at, a.updated_at, b.name AS username FROM go_history AS a LEFT JOIN go_user AS b ON a.user_id = b.id WHERE doc_id = ? ORDER BY a.id DESC"
	err := yiigo.DB.Select(&data, query, docID)

	if err != nil && err != sql.ErrNoRows {
		yiigo.Errf("%s, SQL: %s, Args: [%d]", err.Error(), query, docID)
	}

	return data, err
}
