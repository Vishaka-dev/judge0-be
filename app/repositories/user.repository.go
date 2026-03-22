package repositories

import (
	"context"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/database"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/logger"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/types"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/utils"
)

func RegisterUser(ctx context.Context, user types.RegisterUserRequestType) error {
	logger.Log.Info("Registering user", "user_id", user.UserID, "email", user.Email, "name", user.Name)

	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()

	_, err := pool.Exec(ctx, "INSERT INTO users (user_id, email, name) VALUES ($1, $2, $3)", user.UserID, user.Email, user.Name)
	if err != nil {
		logger.Log.Error("Failed to register user", "user_id", user.UserID, "error", err)
		return err
	}

	logger.Log.Info("User registered successfully", "user_id", user.UserID)
	return nil
}

func CheckUserExists(ctx context.Context, userID string) (bool, error) {
	logger.Log.Info("Checking if user exists", "user_id", userID)

	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()

	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE user_id = $1", userID).Scan(&count)
	if err != nil {
		logger.Log.Error("Failed to check if user exists", "user_id", userID, "error", err)
		return false, err
	}

	logger.Log.Info("User existence check complete", "user_id", userID, "exists", count > 0)
	return count > 0, nil
}

func AddUserToLeaderboard(ctx context.Context, userID string) error {
	logger.Log.Info("Adding user to leaderboard", "user_id", userID)
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()
	_, err := pool.Exec(ctx, "INSERT INTO leaderboard (user_id) VALUES ($1)", userID)
	if err != nil {
		logger.Log.Error("Failed to add user to leaderboard", "user_id", userID, "error", err)
		return err
	}
	logger.Log.Info("User added to leaderboard successfully", "user_id", userID)
	return nil
}
