package service

import (
	"database/sql"
	"godoc/rbac"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

type HistoryService struct {
	Identity *rbac.Identity
}

type History struct {
	ID        int       `db:"id"`
	UserName  string    `db:"username"`
	Flag      int       `db:"flag"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewHistoryService(c *gin.Context) *HistoryService {
	return &HistoryService{
		Identity: rbac.GetIdentity(c),
	}
}

func (h *HistoryService) GetHistory(docID int) ([]History, error) {
	defer yiigo.Flush()

	data := []History{}

	query := "SELECT a.id, a.flag, a.created_at, a.updated_at, b.name AS username FROM go_history AS a LEFT JOIN go_user AS b ON a.user_id = b.id WHERE doc_id = ? ORDER BY a.id DESC"
	err := yiigo.DB.Select(&data, query, docID)

	if err != nil && err != sql.ErrNoRows {
		yiigo.Err(err.Error())
	}

	return data, err
}
