package routes

import (
	"godoc/controllers"
	"godoc/middlewares"

	"github.com/gin-gonic/gin"
)

func LoadDocRoutes(r *gin.Engine) {
	docController := controllers.NewDocController(r)

	g := r.Group("")
	g.Use(middlewares.AuthMiddleware())
	{
		g.GET("/docs/view/:id", docController.View)
		g.GET("/docs/add/:project", docController.Add)
		g.POST("/docs/add/:project", docController.Add)
		g.GET("/docs/edit/:id", docController.Edit)
		g.POST("/docs/edit/:id", docController.Edit)
		g.GET("/docs/delete/:id", docController.Delete)
	}
}
