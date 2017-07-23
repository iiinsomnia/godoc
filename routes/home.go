package routes

import (
	"godoc/controllers"
	"godoc/middlewares"

	"github.com/gin-gonic/gin"
)

func LoadHomeRoutes(r *gin.Engine) {
	homeController := controllers.NewHomeController(r)

	r.GET("/login", homeController.Login)
	r.POST("/login/:captchaID", homeController.Login)
	r.GET("/captcha/:id", homeController.Captcha)

	g := r.Group("")
	g.Use(middlewares.AuthMiddleware())
	{
		g.GET("/", homeController.Index)
		g.GET("/logout", homeController.Logout)
	}
}
