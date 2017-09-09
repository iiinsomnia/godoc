package main

import (
	"fmt"
	"godoc/middlewares"
	"godoc/routes"
	"godoc/views"
	"runtime"

	"godoc/public/assets"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	err := yiigo.Bootstrap(true, false, false)

	if err != nil {
		yiigo.Error(err.Error())
		yiigo.Flush()
	}

	loadStaticResource()

	fmt.Println("app start, version", yiigo.EnvString("app", "version", "1.0.0"))

	run()
}

// load static resource
func loadStaticResource() {
	assets.LoadAssets()
	views.LoadViews()
}

// load routes
func loadRoutes(r *gin.Engine) {
	routes.LoadHomeRoutes(r)
	routes.LoadUserRoutes(r)
	routes.LoadCategoryRoutes(r)
	routes.LoadProjectRoutes(r)
	routes.LoadDocRoutes(r)
}

func run() {
	debug := yiigo.EnvBool("app", "debug", true)

	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	if debug {
		r.Use(gin.Logger(), gin.Recovery())
	}

	r.Use(middlewares.ErrorMiddleware())

	r.StaticFS("/assets", assets.AssetBox.HTTPBox())
	r.StaticFile("/favicon.ico", "./favicon.ico")

	loadRoutes(r)
	r.Run(fmt.Sprintf(":%d", yiigo.EnvInt("app", "port", 8000)))
}
