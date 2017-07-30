package mysql

import "github.com/iiinsomnia/yiigo"

type ProjectDao struct {
	yiigo.MySQL
}

func NewProjectDao() *ProjectDao {
	return &ProjectDao{
		yiigo.MySQL{
			Table: "project",
		},
	}
}

func (p *ProjectDao) GetByID(id int, data interface{}) error {
	query := yiigo.X{
		"select": "a.id, a.name, a.category_id, a.description, a.created_at, a.updated_at, b.name AS category_name",
		"join":   []string{"LEFT JOIN go_category AS b ON a.category_id = b.id"},
		"where":  "a.id = ?",
		"binds":  []interface{}{id},
	}

	err := p.MySQL.FindOne(query, data)

	if err != nil {
		if err.Error() != "not found" {
			yiigo.LogError(err.Error())
		}

		return err
	}

	return nil
}

func (p *ProjectDao) GetByCategoryID(categoryID int, data interface{}) error {
	query := yiigo.X{
		"where": "category_id = ?",
		"binds": []interface{}{categoryID},
		"order": "updated_at DESC",
	}
	err := p.MySQL.Find(query, data)

	if err != nil {
		if err.Error() != "not found" {
			yiigo.LogError(err.Error())
		}

		return err
	}

	return nil
}

func (p *ProjectDao) AddNewRecord(data yiigo.X) (int64, error) {
	id, err := p.MySQL.Insert(data)

	if err != nil {
		yiigo.LogError(err.Error())
		return 0, err
	}

	return id, nil
}

func (p *ProjectDao) UpdateByID(id int, data yiigo.X) error {
	query := yiigo.X{
		"where": "id = ?",
		"binds": []interface{}{id},
	}

	_, err := p.MySQL.Update(query, data)

	if err != nil {
		yiigo.LogError(err.Error())
		return err
	}

	return nil
}

func (p *ProjectDao) DeleteByID(id int) error {
	operations := []yiigo.X{
		yiigo.X{
			"type": "delete",
			"query": yiigo.X{
				"table": "history",
				"where": "project_id = ?",
				"binds": []interface{}{id},
			},
		},
		yiigo.X{
			"type": "delete",
			"query": yiigo.X{
				"table": "doc",
				"where": "project_id = ?",
				"binds": []interface{}{id},
			},
		},
		yiigo.X{
			"type": "delete",
			"query": yiigo.X{
				"where": "id = ?",
				"binds": []interface{}{id},
			},
		},
	}

	err := p.MySQL.DoTransactions(operations)

	if err != nil {
		yiigo.LogError(err.Error())
		return err
	}

	return nil
}
