package routes

import (
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/handlers"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/middlewares"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	challenge := r.Group("/admin")
	{
		challenge.GET("/submissions/dsa", middlewares.AuthMiddleware("Subcommittee"), handlers.GetDSASubmissionResultsHandler)
	}
}
