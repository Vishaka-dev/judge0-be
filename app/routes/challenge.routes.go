package routes

import (
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/handlers"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/middlewares"
	"github.com/gin-gonic/gin"
)

func ChallengeRoutes(r *gin.RouterGroup) {
	challenge := r.Group("/challenge")
	{
		challenge.GET("/get", handlers.GetAllChallengesHandler)
		challenge.GET("/get/:id", handlers.GetChallengeByIdHandler)
		challenge.POST("/add", middlewares.AuthMiddleware("Subcommittee"), handlers.AddChallengeHandler)
		challenge.POST("/test", middlewares.AuthMiddleware("Subcommittee"), handlers.TestDSAChallengeHandler)
		challenge.POST("/submit/dsa", middlewares.AuthMiddleware("Subcommittee"), handlers.SubmitDSAChallengeHandler)
	}
}
