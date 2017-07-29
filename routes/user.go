package routes

import (
	"godoc/controllers"
	"godoc/middlewares"

	"github.com/gin-gonic/gin"
)

func LoadUserRoutes(r *gin.Engine) {
	userController := controllers.NewUserController(r)

	g := r.Group("")
	g.Use(middlewares.AuthMiddleware())
	{
		g.GET("/users", userController.Index)
		g.GET("/users/add", userController.Add)
		g.POST("/users/add", userController.Add)
		g.GET("/users/edit/:id", userController.Edit)
		g.POST("/users/edit/:id", userController.Edit)
		g.POST("/users/password", userController.Password)
		g.GET("/users/reset/:id", userController.Reset)
		g.GET("/users/delete/:id", userController.Delete)
	}
}
