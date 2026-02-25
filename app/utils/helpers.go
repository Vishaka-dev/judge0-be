package utils

import (
	"errors"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/types"
)

func ValidateAddChallengeRequest(request types.AddChallengeRequestType) error {
	if request.Title == "" {
		return errors.New("title is required")
	}
	if request.Description == "" {
		return errors.New("description is required")
	}
	if request.TypeID <= 0 {
		return errors.New("type_id is required")
	}
	if request.StatusID <= 0 {
		return errors.New("status_id is required")
	}
	return nil
}

func ValidateAddDSAChallengeRequest(request types.AddDSAChallengeRequestType) error {
	if request.Title == "" {
		return errors.New("title is required")
	}
	if request.Description == "" {
		return errors.New("description is required")
	}
	if request.TypeID <= 0 {
		return errors.New("type_id is required")
	}
	if request.StatusID <= 0 {
		return errors.New("status_id is required")
	}
	if request.SampleInput == "" {
		return errors.New("sample_input is required")
	}
	if request.SampleOutput == "" {
		return errors.New("sample_output is required")
	}
	return nil
}

func ValidateTestDSAChallengeRequest(request types.TestDSAChallengeRequestType) error {
	if request.ChallengeID <= 0 {
		return errors.New("challenge_id is required")
	}
	if request.SourceCode == "" {
		return errors.New("source_code is required")
	}
	if request.Stdin == "" {
		return errors.New("stdin is required")
	}
	if request.ExpectedOutput == "" {
		return errors.New("expected_output is required")
	}
	return nil
}
