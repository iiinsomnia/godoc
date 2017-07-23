package routes

import (
	"godoc/controllers"
	"godoc/middlewares"

	"github.com/gin-gonic/gin"
)

func LoadProjectRoutes(r *gin.Engine) {
	projectController := controllers.NewProjectController(r)

	g := r.Group("")
	g.Use(middlewares.AuthMiddleware())
	{
		g.GET("/projects/view/:id", projectController.View)
		g.GET("/projects/add/:category", projectController.Add)
		g.POST("/projects/add/:category", projectController.Add)
		g.GET("/projects/edit/:id", projectController.Edit)
		g.POST("/projects/edit/:id", projectController.Edit)
		g.GET("/projects/delete/:id", projectController.Delete)
	}
}
