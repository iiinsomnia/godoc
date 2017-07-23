package mysql

import "github.com/iiinsomnia/yiigo"

type DocDao struct {
	yiigo.MySQL
}

func NewDocDao() *DocDao {
	return &DocDao{
		yiigo.MySQL{
			Table: "doc",
		},
	}
}

func (d *DocDao) GetByID(id int, data interface{}) error {
	query := yiigo.X{
		"select": "a.id, a.title, a.category_id, a.project_id, a.label, a.markdown, a.created_at, a.updated_at, b.name AS category_name, c.name AS project_name",
		"join": []string{
			"LEFT JOIN go_category AS b ON a.category_id = b.id",
			"LEFT JOIN go_project AS c ON a.project_id = c.id",
		},
		"where": "a.id = ?",
		"binds": []interface{}{id},
	}

	err := d.MySQL.FindOne(query, data)

	if err != nil {
		if err.Error() != "not found" {
			yiigo.LogError(err.Error())
		}

		return err
	}

	return nil
}

func (d *DocDao) GetByProjectID(projectID int, data interface{}) error {
	query := yiigo.X{
		"where": "project_id = ?",
		"binds": []interface{}{projectID},
		"order": "updated_at DESC",
	}

	err := d.MySQL.Find(query, data)

	if err != nil {
		if err.Error() != "not found" {
			yiigo.LogError(err.Error())
		}

		return err
	}

	return nil
}

func (d *DocDao) AddNewRecord(data yiigo.X) (int64, error) {
	id, err := d.MySQL.Insert(data)

	if err != nil {
		yiigo.LogError(err.Error())
		return 0, err
	}

	return id, nil
}

func (d *DocDao) UpdateByID(id int, data yiigo.X) error {
	query := yiigo.X{
		"where": "id = ?",
		"binds": []interface{}{id},
	}

	_, err := d.MySQL.Update(query, data)

	if err != nil {
		yiigo.LogError(err.Error())
		return err
	}

	return nil
}

func (d *DocDao) DeleteByID(id int) error {
	operations := []yiigo.X{
		yiigo.X{
			"type": "delete",
			"query": yiigo.X{
				"table": "log",
				"where": "doc_id = ?",
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

	err := d.MySQL.DoTransactions(operations)

	if err != nil {
		yiigo.LogError(err.Error())
		return err
	}

	return nil
}
