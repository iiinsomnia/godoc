package service

import (
	"godoc/dao/mysql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

type ProjectService struct {
	*service
}

type Project struct {
	ID           int       `db:"id"`
	Name         string    `db:"name"`
	CategoryID   int       `db:"category_id"`
	CategoryName string    `db:"category_name"`
	Description  string    `db:"description"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func NewProjectService(c *gin.Context) *ProjectService {
	return &ProjectService{
		construct(c),
	}
}

func (p *ProjectService) GetProjects(categoryID int) ([]Project, error) {
	data := []Project{}

	projectDao := mysql.NewProjectDao()
	err := projectDao.GetByCategoryID(categoryID, &data)

	return data, err
}

func (p *ProjectService) GetDetail(id int) (*Project, error) {
	data := &Project{}

	projectDao := mysql.NewProjectDao()
	err := projectDao.GetByID(id, data)

	return data, err
}

func (p *ProjectService) Add(data yiigo.X) (int64, error) {
	projectDao := mysql.NewProjectDao()
	id, err := projectDao.AddNewRecord(data)

	return id, err
}

func (p *ProjectService) Edit(id int, data yiigo.X) error {
	projectDao := mysql.NewProjectDao()
	err := projectDao.UpdateByID(id, data)

	return err
}

func (p *ProjectService) Delete(id int) error {
	projectDao := mysql.NewProjectDao()
	err := projectDao.DeleteByID(id)

	return err
}
