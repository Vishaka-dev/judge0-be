package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetChallengesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"test": "test",
	})
}
