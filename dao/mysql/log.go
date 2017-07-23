package mysql

import "github.com/iiinsomnia/yiigo"

type LogDao struct {
	yiigo.MySQL
}

func NewLogDao() *LogDao {
	return &LogDao{
		yiigo.MySQL{
			Table: "log",
		},
	}
}

func (l *LogDao) GetByDocID(docID int, data interface{}) error {
	query := yiigo.X{
		"select": "a.id, a.user_id, a.category_id, a.project_id, a.doc_id, a.flag, a.created_at, a.updated_at, b.name AS username",
		"join":   []string{"LEFT JOIN go_user AS b ON a.user_id = b.id"},
		"where":  "a.doc_id = ?",
		"binds":  []interface{}{docID},
		"order":  "id DESC",
	}

	err := l.MySQL.Find(query, data)

	if err != nil {
		if err.Error() != "not found" {
			yiigo.LogError(err.Error())
		}

		return err
	}

	return nil
}

func (l *LogDao) AddNewRecord(data yiigo.X) (int64, error) {
	id, err := l.MySQL.Insert(data)

	if err != nil {
		yiigo.LogError(err.Error())
		return 0, err
	}

	return id, nil
}
