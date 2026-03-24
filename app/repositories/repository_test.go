package repositories

import (
	"context"
	"os"
	"testing"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/database"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/logger"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/types"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/utils"
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

func TestGetDSATestCaseCount(t *testing.T) {
	count, err := GetDSATestCaseCount(context.Background(), 6)
	logger.Log.Info("TestGetDSATestCaseCount", "count", count, "error", err)
	if count != 4 {
		t.Fatal("Expected 4 test cases, got", count)
	}
}

func TestAddDSASubmission(t *testing.T) {
	testCount, err := GetDSATestCaseCount(context.Background(), 6)
	status, err := AddDSASubmission(context.Background(), utils.GenerateSubmissionID(), 6, "7ef624e2-790c-468c-91c6-bf97d7232620", int(testCount))

	if status {
		logger.Log.Info("Submission Added Subccessfully")
	} else {
		logger.Log.Info("error ", err)
	}
}

func TestRegisterUser(t *testing.T) {
	err := RegisterUser(context.Background(), types.RegisterUserRequestType{
		UserID: "7ef624e2-790c-468c-91c6-bf97d7232620",
		Email:  "test@example.com",
		Name:   "Test User",
	})
	if err != nil {
		t.Fatal("Failed to register user", err)
	}
}

func TestCheckUserExists(t *testing.T) {
	exists, err := CheckUserExists(context.Background(), "7ef624e2-790c-468c-91c6-bf97d7232620")
	logger.Log.Info("TestCheckUserExists", "exists", exists, "error", err)
	if !exists {
		t.Fatal("Expected user to exist, but it does not")
	}
}

func TestAddUserToLeaderboard(t *testing.T) {
	err := AddUserToLeaderboard(context.Background(), "7ef624e2-790c-468c-91c6-bf97d7232620")
	if err != nil {
		t.Fatal("Failed to add user to leaderboard", err)
	}
}
