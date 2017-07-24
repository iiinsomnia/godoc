package service

import (
	"godoc/dao/mysql"
	"time"

	"github.com/gin-gonic/gin"
)

type HistoryService struct {
	*service
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
		construct(c),
	}
}

func (p *HistoryService) GetHistory(docID int) ([]History, error) {
	data := []History{}

	historyDao := mysql.NewHistoryDao()
	err := historyDao.GetByDocID(docID, &data)

	return data, err
}
