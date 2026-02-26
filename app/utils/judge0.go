package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/config"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/logger"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/types"
)

var judge0HTTPClient = &http.Client{
	Timeout: 15 * time.Second,
	Transport: &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func TestDSAChallenge(ctx context.Context, payload any) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		logger.Log.Error("Failed to marshal payload for Judge0 submission", "error", err, "payload", payload)
		return nil, err
	}

	url := config.Get().Judge0API + "/submissions?base64_encoded=false&wait=true"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		logger.Log.Error("Failed to create Judge0 request", "error", err, "url", url)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := judge0HTTPClient.Do(req)
	if err != nil {
		logger.Log.Error("Judge0 API request failed", "error", err, "url", url)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error("Failed to read Judge0 response body", "error", err)
		return nil, err
	}
	if resp.StatusCode >= 400 {
		logger.Log.Error("Judge0 API returned error status", "status", resp.StatusCode, "body", string(respBody))
		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, string(respBody))
	}
	logger.Log.Info("Judge0 submission successful", "status", resp.StatusCode)
	return respBody, nil
}
func SubmitDSAChallenge(ctx context.Context, testCases []types.DSAChallengeTestCase, payload types.SubmitDSAChallengeRequestType) (bool, error) {
	submissions := make([]types.Judge0SubmissionRequest, len(testCases))
	for i, tc := range testCases {
		submissions[i] = types.Judge0SubmissionRequest{
			LanguageID:     payload.LanguageID,
			SourceCode:     payload.SourceCode,
			CallbackURL:    "https://webhook.site/74e67551-b71a-437e-97dd-696340c70efd",
			Stdin:          Base64Encode(tc.TestInput),
			ExpectedOutput: Base64Encode(tc.TestOutput),
		}
	}

	batchPayload := types.Judge0BatchSubmissionRequest{
		Submissions: submissions,
	}

	body, err := json.Marshal(batchPayload)
	if err != nil {
		logger.Log.Error("Failed to marshal Judge0 batch payload", "error", err, "payload", batchPayload)
		return false, err
	}

	// url := "https://webhook.site/74e67551-b71a-437e-97dd-696340c70efd"
	// To use Judge0 API, uncomment below:
	url := config.Get().Judge0API + "/submissions/batch?base64_encoded=true&wait=false"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		logger.Log.Error("Failed to create Judge0 batch request", "error", err, "url", url)
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := judge0HTTPClient.Do(req)
	if err != nil {
		logger.Log.Error("Judge0 batch API request failed", "error", err, "url", url)
		return false, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error("Failed to read Judge0 batch response body", "error", err)
		return false, err
	}
	if resp.StatusCode >= 400 {
		logger.Log.Error("Judge0 batch API returned error status", "status", resp.StatusCode, "body", string(respBody))
		return false, fmt.Errorf("http %d: %s", resp.StatusCode, string(respBody))
	}
	logger.Log.Info("Judge0 batch submission successful", "status", resp.StatusCode)
	return true, nil
}
