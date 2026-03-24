package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/logger"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/repositories"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/types"
	"github.com/gin-gonic/gin"
)

func RegisterUserHandler(c *gin.Context) {
	body, err := c.GetRawData()
	logger.Log.Info("Received request to register user", "path", c.Request.URL.Path, "method", c.Request.Method)

	if err != nil {
		logger.Log.Warn("Failed to read request body in RegisterUserHandler", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read body"})
		return
	}

	var user types.RegisterUserRequestType
	if err := json.Unmarshal(body, &user); err != nil {
		logger.Log.Warn("Invalid JSON in RegisterUserHandler", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	logger.Log.Info("Parsed register user request", "user_id", user.UserID, "email", user.Email, "name", user.Name)

	userExists, err := repositories.CheckUserExists(c.Request.Context(), user.UserID)
	if err != nil {
		logger.Log.Error("Error checking if user exists in RegisterUserHandler", "user_id", user.UserID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
		return
	}

	if userExists {
		logger.Log.Warn("User already exists in RegisterUserHandler", "user_id", user.UserID)
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	if err := repositories.RegisterUser(c.Request.Context(), user); err != nil {
		logger.Log.Error("Failed to register user in RegisterUserHandler", "user_id", user.UserID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	if err := repositories.AddUserToLeaderboard(c.Request.Context(), user.UserID); err != nil {
		logger.Log.Error("Failed to add user to leaderboard in RegisterUserHandler", "user_id", user.UserID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to leaderboard"})
		return
	}

	logger.Log.Info("User registered successfully", "user_id", user.UserID)
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
