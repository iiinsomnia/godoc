package routes

import (
	"godoc/controllers"
	"godoc/middlewares"

	"github.com/gin-gonic/gin"
)

func LoadCategoryRoutes(r *gin.Engine) {
	categoryController := controllers.NewCategoryController(r)

	g := r.Group("")
	g.Use(middlewares.AuthMiddleware())
	{
		g.GET("/categories/view/:id", categoryController.View)
		g.POST("/categories/add", categoryController.Add)
		g.POST("/categories/edit/:id", categoryController.Edit)
		g.GET("/categories/delete/:id", categoryController.Delete)
	}
}
