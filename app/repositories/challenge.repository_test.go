package repositories

import (
	"context"
	"os"
	"testing"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/database"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/logger"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load("e:/Coding/clubs/mozilla/official/judge0-be/.env")
	database.Init()
	code := m.Run()
	database.Close()
	os.Exit(code)
}

func TestGetDSAChallengeTestCases(t *testing.T) {
	challenges, err := GetDSAChallengeTestCases(context.Background(), 6)
	logger.Log.Info("TestGetDSAChallengeTestCases", "challenges", challenges, "error", err)
}
