package routes

import (
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/handlers"
	"github.com/gin-gonic/gin"
)

func ChallengeRoutes(r *gin.RouterGroup) {
	challenge := r.Group("/challenge")
	{
		challenge.GET("/get", handlers.GetChallengesHandler)
	}
}
