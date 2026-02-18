package handlers

import (
	"net/http"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/repositories"
	"github.com/gin-gonic/gin"
)

func GetAllChallengesHandler(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	challenges, currentPage, totalPages, err := repositories.GetAllChallenges(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"currentPage": currentPage,
		"totalPages":  totalPages,
		"challenges":  challenges,
	})
}

func GetChallengeByIdHandler(c *gin.Context) {
	id := c.Param("id")

	challengeType, err := repositories.GetChallengeType(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if challengeType == "1" {
		dsaChallenge, err := repositories.GetDSAChallenge(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"challenge": dsaChallenge,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Challenge type not supported",
	})

}
