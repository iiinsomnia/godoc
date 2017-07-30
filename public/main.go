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

	bootstrap()
	loadStaticResource()

	version := yiigo.GetEnvString("app", "version", "1.0.0")
	fmt.Println("app start, version", version)

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

func bootstrap() {
	b := yiigo.New()

	// b.EnableMongo()
	// b.EnableRedis()

	err := b.Bootstrap()

	if err != nil {
		yiigo.LogError(err.Error())
	}
}

func run() {
	debug := yiigo.GetEnvBool("app", "debug", false)
	mode := gin.ReleaseMode

	if debug {
		mode = gin.DebugMode
	}

	gin.SetMode(mode)

	r := gin.New()
	r.Use(middlewares.ErrorMiddleware())

	r.StaticFS("/assets", assets.Asset.HTTPBox())
	r.StaticFile("/favicon.ico", "./favicon.ico")
	// r.LoadHTMLGlob("../views/**/**/*")
	loadRoutes(r)
	r.Run(fmt.Sprintf(":%d", yiigo.GetEnvInt("app", "port", 8000)))
}
