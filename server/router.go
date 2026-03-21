package server

import (
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/routes"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	app := gin.New()
	app.Use(gin.Logger(), gin.Recovery())

	api := app.Group("/api")
	routes.RegisterAllRoutes(api)

	return app
}
