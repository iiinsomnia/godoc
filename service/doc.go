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

	apiDao := mysql.NewDocDao()
	err := apiDao.GetByProjectID(projectID, &data)

	return data, err
}

func (d *DocService) GetDetail(id int) (*Doc, error) {
	data := &Doc{}

	apiDao := mysql.NewDocDao()
	err := apiDao.GetByID(id, data)

	return data, err
}

func (d *DocService) Add(data yiigo.X) (int64, error) {
	apiDao := mysql.NewDocDao()
	id, err := apiDao.AddNewRecord(data)

	return id, err
}

func (d *DocService) Edit(id int, data yiigo.X) error {
	apiDao := mysql.NewDocDao()
	err := apiDao.UpdateByID(id, data)

	return err
}

func (d *DocService) Delete(id int) error {
	apiDao := mysql.NewDocDao()
	err := apiDao.DeleteByID(id)

	return err
}
