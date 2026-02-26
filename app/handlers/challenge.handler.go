package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/logger"
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
		logger.Log.Error("Failed to get all challenges", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	logger.Log.Info("Fetched all challenges", "currentPage", currentPage, "totalPages", totalPages, "challenge_count", len(challenges))
	c.JSON(http.StatusOK, gin.H{
		"currentPage": currentPage,
		"totalPages":  totalPages,
		"challenges":  challenges,
	})
}

func GetChallengeByIdHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	logger.Log.Info("Fetching challenge by ID", "id", id)

	challengeType, err := repositories.GetChallengeType(ctx, id)

	if err != nil {
		logger.Log.Error("Failed to get challenge type", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	switch challengeType {
	case 1:
		dsaChallenge, err := repositories.GetDSAChallenge(ctx, id)
		if err != nil {
			logger.Log.Error("Failed to get DSA challenge", "id", id, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		logger.Log.Info("Fetched DSA challenge", "id", id)

		c.JSON(http.StatusOK, gin.H{
			"challenge": dsaChallenge,
		})
		return
	}

	logger.Log.Warn("Challenge not found", "id", id)
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Challenge not found",
	})

}

func AddChallengeHandler(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		logger.Log.Warn("Failed to read request body in AddChallengeHandler", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read body"})
		return
	}

	var challenge types.AddChallengeRequestType
	if err := json.Unmarshal(body, &challenge); err != nil {
		logger.Log.Warn("Invalid JSON in AddChallengeHandler", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if err := utils.ValidateAddChallengeRequest(challenge); err != nil {
		logger.Log.Warn("Validation failed in AddChallengeHandler", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if challenge.TypeID == 1 {
		ctx := c.Request.Context()
		var dsaChallenge types.AddDSAChallengeRequestType
		if err := json.Unmarshal(body, &dsaChallenge); err != nil {
			logger.Log.Warn("Invalid JSON in AddChallengeHandler (DSA)", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		if err := utils.ValidateAddDSAChallengeRequest(dsaChallenge); err != nil {
			logger.Log.Warn("Validation failed in AddChallengeHandler (DSA)", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		challengeID, err := repositories.AddDSAChallenge(ctx, dsaChallenge)
		if err != nil {
			logger.Log.Error("Failed to add DSA challenge", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logger.Log.Info("DSA challenge created successfully", "challenge_id", challengeID)

		c.JSON(http.StatusCreated, gin.H{
			"id":      challengeID,
			"message": "DSA challenge created successfully",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Unsupported challenge type",
	})
	logger.Log.Warn("Unsupported challenge type in AddChallengeHandler", "type_id", challenge.TypeID)
}

func TestDSAChallengeHandler(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		logger.Log.Warn("Failed to read request body in TestDSAChallengeHandler", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read body"})
		return
	}

	var challenge types.TestDSAChallengeRequestType
	if err := json.Unmarshal(body, &challenge); err != nil {
		logger.Log.Warn("Invalid JSON in TestDSAChallengeHandler", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := utils.ValidateTestDSAChallengeRequest(challenge); err != nil {
		logger.Log.Warn("Validation failed in TestDSAChallengeHandler", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	byteResp, err := utils.TestDSAChallenge(ctx, challenge)
	if err != nil {
		logger.Log.Error("Judge0 test submission failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result types.TestDSAChallengeResponse
	if err := json.Unmarshal(byteResp, &result); err != nil {
		logger.Log.Error("Failed to unmarshal Judge0 response", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid response from judge0"})
		return
	}
	logger.Log.Info("Judge0 test submission succeeded", "challenge_id", challenge.ChallengeID)

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func SubmitDSAChallengeHandler(c *gin.Context) {

	body, err := c.GetRawData()
	if err != nil {
		logger.Log.Warn("Failed to read request body in SubmitDSAChallengeHandler", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read body"})
		return
	}

	var challenge types.SubmitDSAChallengeRequestType
	if err := json.Unmarshal(body, &challenge); err != nil {
		logger.Log.Warn("Invalid JSON in SubmitDSAChallengeHandler", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := utils.ValidateSubmitDSAChallengeRequest(challenge); err != nil {
		logger.Log.Warn("Validation failed in SubmitDSAChallengeHandler", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	testCases, err := repositories.GetDSAChallengeTestCases(c.Request.Context(), challenge.ChallengeID)
	if err != nil {
		logger.Log.Error("Failed to get DSA challenge test cases", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = utils.SubmitDSAChallenge(c, testCases, challenge)
	if err != nil {
		logger.Log.Error("Failed to submit DSA challenge", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userEmail, _ := c.Get("user_email")
	logger.Log.Info("Received request to submit DSA challenge solution", "user_email", userEmail)
	submissionId := utils.GenerateSubmissionID()

	c.JSON(http.StatusOK, gin.H{
		"submission_id": submissionId,
	})
}
