package api

import (
	"net/http"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/config"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/database"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/routes"
	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func init() {
	config.Get()
	database.Init()

	app = gin.New()
	app.Use(gin.Logger(), gin.Recovery())

	app.Static("/public", "./public")
	app.StaticFile("/", "./public/index.html")
	app.StaticFile("/openapi.yaml", "./public/openapi.yaml")

	api := app.Group("/api")
	routes.RegisterAllRoutes(api)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
