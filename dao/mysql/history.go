package mysql

import (
	"time"

	"github.com/iiinsomnia/yiigo"
)

type HistoryDao struct {
	yiigo.MySQL
}

func NewHistoryDao() *HistoryDao {
	return &HistoryDao{
		yiigo.MySQL{
			Table: "history",
		},
	}
}

func (h *HistoryDao) GetByDocID(docID int, data interface{}) error {
	query := yiigo.X{
		"select": "a.id, a.flag, a.created_at, a.updated_at, b.name AS username",
		"join":   []string{"LEFT JOIN go_user AS b ON a.user_id = b.id"},
		"where":  "a.doc_id = ?",
		"binds":  []interface{}{docID},
		"order":  "a.id DESC",
	}

	err := h.MySQL.Find(query, data)

	if err != nil {
		if err.Error() != "not found" {
			yiigo.LogError(err.Error())
		}

		return err
	}

	return nil
}

func (h *HistoryDao) AddNewRecord(data yiigo.X) (int64, error) {
	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()

	id, err := h.MySQL.Insert(data)

	if err != nil {
		yiigo.LogError(err.Error())
		return 0, err
	}

	return id, nil
}
