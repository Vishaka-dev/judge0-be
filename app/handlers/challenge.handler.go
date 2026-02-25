package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/repositories"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/types"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/utils"
	"github.com/gin-gonic/gin"
)

func GetAllChallengesHandler(c *gin.Context) {
	ctx := c.Request.Context()
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	challenges, currentPage, totalPages, err := repositories.GetAllChallenges(ctx, page, pageSize)
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
	ctx := c.Request.Context()
	id := c.Param("id")

	challengeType, err := repositories.GetChallengeType(ctx, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	switch challengeType {
	case 1:
		dsaChallenge, err := repositories.GetDSAChallenge(ctx, id)
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

	c.JSON(http.StatusNotFound, gin.H{
		"error": "Challenge not found",
	})

}

func AddChallengeHandler(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read body"})
		return
	}

	var challenge types.AddChallengeRequestType
	if err := json.Unmarshal(body, &challenge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if err := utils.ValidateAddChallengeRequest(challenge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if challenge.TypeID == 1 {
		ctx := c.Request.Context()
		var dsaChallenge types.AddDSAChallengeRequestType
		if err := json.Unmarshal(body, &dsaChallenge); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		if err := utils.ValidateAddDSAChallengeRequest(dsaChallenge); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		challengeID, err := repositories.AddDSAChallenge(ctx, dsaChallenge)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":      challengeID,
			"message": "DSA challenge created successfully",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Unsupported challenge type",
	})
}

func TestDSAChallengeHandler(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read body"})
		return
	}

	var challenge types.TestDSAChallengeRequestType
	if err := json.Unmarshal(body, &challenge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if err := utils.ValidateTestDSAChallengeRequest(challenge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	byteResp, err := utils.TestDSAChallengeHandler(ctx, challenge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result types.TestDSAChallengeResponse
	if err := json.Unmarshal(byteResp, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid response from judge0"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}
