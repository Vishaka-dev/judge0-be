package handlers

import (
	"net/http"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/repositories"
	"github.com/gin-gonic/gin"
)

func GetAllChallengesHandler(c *gin.Context) {
	challenges, err := repositories.GetAllChallenges()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"challenges": challenges,
	})
}
