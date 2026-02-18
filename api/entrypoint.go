package api

import (
	"net/http"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/config"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/routes"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/supabase"
	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func init() {
	config.Load()
	supabase.Init()

	app = gin.New()
	app.Use(gin.Logger(), gin.Recovery())
	api := app.Group("/api")
	routes.RegisterAllRoutes(api)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
