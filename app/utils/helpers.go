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
