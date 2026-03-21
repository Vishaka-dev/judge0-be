package utils_test

import (
	"context"
	"os"
	"testing"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/database"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/logger"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/repositories"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/utils"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/types"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load("e:/Coding/clubs/mozilla/official/judge0-be/.env")
	database.Init()
	code := m.Run()
	database.Close()
	os.Exit(code)
}

func TestSubmitDSAChallenge(t *testing.T) {
	testCases, err := repositories.GetDSAChallengeTestCases(context.Background(), 6)
	if err != nil {
		t.Fatalf("Failed to get DSA challenge test cases: %v", err)
	}
	payload := types.SubmitDSAChallengeRequestType{
		ChallengeID: 6,
		SourceCode:  "print('Hello, World!')",
		LanguageID:  71,
	}
	utils.SubmitDSAChallenge(context.Background(), testCases, payload, utils.GenerateSubmissionID())
}

func TestBase64Encode(t *testing.T) {
	testCases, _ := repositories.GetDSAChallengeTestCases(context.Background(), 6)

	input := testCases[0].TestInput
	expected := "NQoxIDIgMyA0IDU="
	output := utils.Base64Encode(input)
	logger.Log.Info("input", input)
	logger.Log.Info("output", output)
	logger.Log.Info("expected", expected)
}
