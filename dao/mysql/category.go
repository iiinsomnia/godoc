package mysql

import "github.com/iiinsomnia/yiigo"

type CategoryDao struct {
	yiigo.MySQL
}

func NewCategoryDao() *CategoryDao {
	return &CategoryDao{
		yiigo.MySQL{
			Table: "category",
		},
	}
}

func (c *CategoryDao) GetByID(id int, data interface{}) error {
	query := yiigo.X{
		"where": "id = ?",
		"binds": []interface{}{id},
	}

	err := c.MySQL.FindOne(query, data)

	if err != nil {
		if err.Error() != "not found" {
			yiigo.LogError(err.Error())
		}

		return err
	}

	return nil
}

func (c *CategoryDao) GetAll(data interface{}) error {
	query := yiigo.X{
		"order": "updated_at DESC",
	}
	err := c.MySQL.Find(query, data)

	if err != nil {
		if err.Error() != "not found" {
			yiigo.LogError(err.Error())
		}

		return err
	}

	return nil
}

func (c *CategoryDao) AddNewRecord(data yiigo.X) (int64, error) {
	id, err := c.MySQL.Insert(data)

	if err != nil {
		yiigo.LogError(err.Error())
		return 0, err
	}

	return id, nil
}

func (c *CategoryDao) UpdateByID(id int, data yiigo.X) error {
	query := yiigo.X{
		"where": "id = ?",
		"binds": []interface{}{id},
	}

	_, err := c.MySQL.Update(query, data)

	if err != nil {
		yiigo.LogError(err.Error())
		return err
	}

	return nil
}

func (c *CategoryDao) DeleteByID(id int) error {
	operations := []yiigo.X{
		yiigo.X{
			"type": "delete",
			"query": yiigo.X{
				"table": "history",
				"where": "category_id = ?",
				"binds": []interface{}{id},
			},
		},
		yiigo.X{
			"type": "delete",
			"query": yiigo.X{
				"table": "doc",
				"where": "category_id = ?",
				"binds": []interface{}{id},
			},
		},
		yiigo.X{
			"type": "delete",
			"query": yiigo.X{
				"table": "project",
				"where": "category_id = ?",
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

	err := c.MySQL.DoTransactions(operations)

	if err != nil {
		yiigo.LogError(err.Error())
		return err
	}

	return nil
}
