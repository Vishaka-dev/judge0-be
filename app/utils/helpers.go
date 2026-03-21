package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/logger"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/types"
)

func ValidateAddChallengeRequest(request types.AddChallengeRequestType) error {
	if request.Title == "" {
		logger.Log.Warn("Validation failed: title is required", "request", request)
		return errors.New("title is required")
	}
	if request.Description == "" {
		logger.Log.Warn("Validation failed: description is required", "request", request)
		return errors.New("description is required")
	}
	if request.TypeID <= 0 {
		logger.Log.Warn("Validation failed: type_id is required", "request", request)
		return errors.New("type_id is required")
	}
	if request.StatusID <= 0 {
		logger.Log.Warn("Validation failed: status_id is required", "request", request)
		return errors.New("status_id is required")
	}
	if request.Marks <= 0 {
		logger.Log.Warn("Validation failed: marks is required", "request", request)
		return errors.New("marks is required")
	}
	return nil
}

func ValidateAddDSAChallengeRequest(request types.AddDSAChallengeRequestType) error {
	if request.Title == "" {
		logger.Log.Warn("Validation failed: title is required", "request", request)
		return errors.New("title is required")
	}
	if request.Description == "" {
		logger.Log.Warn("Validation failed: description is required", "request", request)
		return errors.New("description is required")
	}
	if request.TypeID <= 0 {
		logger.Log.Warn("Validation failed: type_id is required", "request", request)
		return errors.New("type_id is required")
	}
	if request.StatusID <= 0 {
		logger.Log.Warn("Validation failed: status_id is required", "request", request)
		return errors.New("status_id is required")
	}
	if request.Marks <= 0 {
		logger.Log.Warn("Validation failed: marks is required", "request", request)
		return errors.New("marks is required")
	}
	if request.SampleInput == "" {
		logger.Log.Warn("Validation failed: sample_input is required", "request", request)
		return errors.New("sample_input is required")
	}
	if request.SampleOutput == "" {
		logger.Log.Warn("Validation failed: sample_output is required", "request", request)
		return errors.New("sample_output is required")
	}
	if len(request.TestCases) == 0 {
		logger.Log.Warn("Validation failed: test_cases is required", "request", request)
		return errors.New("test_cases is required")
	}

	for i, testCase := range request.TestCases {
		if strings.TrimSpace(testCase.TestInput) == "" {
			logger.Log.Warn("Validation failed: test case input is required", "index", i)
			return fmt.Errorf("test_cases[%d].test_input is required", i)
		}
		if strings.TrimSpace(testCase.TestOutput) == "" {
			logger.Log.Warn("Validation failed: test case output is required", "index", i)
			return fmt.Errorf("test_cases[%d].test_output is required", i)
		}
	}
	return nil
}

func ValidateTestDSAChallengeRequest(request types.TestDSAChallengeRequestType) error {
	if request.ChallengeID <= 0 {
		logger.Log.Warn("Validation failed: challenge_id is required", "request", request)
		return errors.New("challenge_id is required")
	}
	if request.SourceCode == "" {
		logger.Log.Warn("Validation failed: source_code is required", "request", request)
		return errors.New("source_code is required")
	}
	if request.Stdin == "" {
		logger.Log.Warn("Validation failed: stdin is required", "request", request)
		return errors.New("stdin is required")
	}
	if request.ExpectedOutput == "" {
		logger.Log.Warn("Validation failed: expected_output is required", "request", request)
		return errors.New("expected_output is required")
	}
	return nil
}

func ValidateSubmitDSAChallengeRequest(request types.SubmitDSAChallengeRequestType) error {
	if request.ChallengeID <= 0 {
		logger.Log.Warn("Validation failed: challenge_id is required", "request", request)
		return errors.New("challenge_id is required")
	}
	if request.SourceCode == "" {
		logger.Log.Warn("Validation failed: source_code is required", "request", request)
		return errors.New("source_code is required")
	}
	if request.LanguageID <= 0 {
		logger.Log.Warn("Validation failed: language_id is required", "request", request)
		return errors.New("language_id is required")
	}
	return nil
}

func GenerateSubmissionID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}
