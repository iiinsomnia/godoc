package service

import (
	"godoc/dao/mysql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

type DocService struct {
	*service
}

type Doc struct {
	ID           int       `db:"id"`
	Title        string    `db:"title"`
	CategoryID   int       `db:"category_id"`
	CategoryName string    `db:"category_name"`
	ProjectID    int       `db:"project_id"`
	ProjectName  string    `db:"project_name"`
	Label        string    `db:"label"`
	Markdown     string    `db:"markdown"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func NewDocService(c *gin.Context) *DocService {
	return &DocService{
		construct(c),
	}
}

func (d *DocService) GetDocs(projectID int) ([]Doc, error) {
	data := []Doc{}

	docDao := mysql.NewDocDao()
	err := docDao.GetByProjectID(projectID, &data)

	return data, err
}

func (d *DocService) GetDetail(id int) (*Doc, error) {
	data := &Doc{}

	docDao := mysql.NewDocDao()
	err := docDao.GetByID(id, data)

	return data, err
}

func (d *DocService) Add(data yiigo.X, history yiigo.X) (int64, error) {
	docDao := mysql.NewDocDao()
	id, err := docDao.AddNewRecord(data)

	if err == nil {
		historyDao := mysql.NewHistoryDao()

		history["user_id"] = d.Identity.ID
		history["doc_id"] = id
		history["flag"] = 1

		historyDao.AddNewRecord(history)
	}

	return id, err
}

func (d *DocService) Edit(id int, data yiigo.X) error {
	docDao := mysql.NewDocDao()

	doc := &Doc{}
	err := docDao.GetByID(id, doc)

	if err != nil {
		return err
	}

	err = docDao.UpdateByID(id, data)

	if err == nil {
		historyDao := mysql.NewHistoryDao()

		history := yiigo.X{
			"user_id":     d.Identity.ID,
			"category_id": doc.CategoryID,
			"project_id":  doc.ProjectID,
			"doc_id":      id,
			"flag":        2,
		}

		historyDao.AddNewRecord(history)
	}

	return err
}

func (d *DocService) Delete(id int) error {
	docDao := mysql.NewDocDao()
	err := docDao.DeleteByID(id)

	return err
}
