package service

import (
	"godoc/dao/mysql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

type CategoryService struct {
	*service
}

type Category struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewCategoryService(c *gin.Context) *CategoryService {
	return &CategoryService{
		construct(c),
	}
}

func (c *CategoryService) GetAll() ([]Category, error) {
	data := []Category{}

	categoryDao := mysql.NewCategoryDao()
	err := categoryDao.GetAll(&data)

	return data, err
}

func (c *CategoryService) GetDetail(id int) (*Category, error) {
	data := &Category{}

	categoryDao := mysql.NewCategoryDao()
	err := categoryDao.GetByID(id, data)

	return data, err
}

func (c *CategoryService) Add(data yiigo.X) (int64, error) {
	categoryDao := mysql.NewCategoryDao()
	id, err := categoryDao.AddNewRecord(data)

	return id, err
}

func (c *CategoryService) Edit(id int, data yiigo.X) error {
	categoryDao := mysql.NewCategoryDao()
	err := categoryDao.UpdateByID(id, data)

	return err
}

func (c *CategoryService) Delete(id int) error {
	categoryDao := mysql.NewCategoryDao()
	err := categoryDao.DeleteByID(id)

	return err
}
