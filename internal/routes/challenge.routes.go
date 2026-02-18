package routes

import "github.com/gin-gonic/gin"

func ChallengeRoutes(r *gin.RouterGroup) {
	challenge := r.Group("/challenge")
}
