package server

import (
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/routes"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	app := gin.New()
	app.Use(gin.Logger(), gin.Recovery())

	app.Static("/public", "./public")
	app.StaticFile("/", "./public/index.html")
	app.StaticFile("/openapi.yaml", "./public/openapi.yaml")

	api := app.Group("/api")
	routes.RegisterAllRoutes(api)

	return app
}
