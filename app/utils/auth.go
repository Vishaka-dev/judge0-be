package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/config"
)

type CurrentUserConnection struct {
	LinkedAt             string `json:"linkedAt"`
	Provider             string `json:"provider"`
	ProviderAccountEmail string `json:"providerAccountEmail"`
	ProviderUserID       string `json:"providerUserId"`
	ProviderUserName     string `json:"providerUserName"`
}

type CurrentUser struct {
	ID          string                  `json:"id"`
	Email       string                  `json:"email"`
	Name        string                  `json:"name"`
	Roles       []string                `json:"roles"`
	Private     bool                    `json:"private"`
	CreatedAt   string                  `json:"createdAt"`
	UpdatedAt   string                  `json:"updatedAt"`
	Connections []CurrentUserConnection `json:"connections"`
}

type GetCurrentUserResponse struct {
	Data CurrentUser `json:"data"`
}

type HTTPStatusError struct {
	StatusCode int
	Message    string
}

func (e *HTTPStatusError) Error() string {
	return e.Message
}

var authHTTPClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func GetCurrentUser(ctx context.Context, authorizationHeader string) (*CurrentUser, error) {
	headerParts := strings.Fields(authorizationHeader)
	if len(headerParts) != 2 || !strings.EqualFold(headerParts[0], "Bearer") {
		return nil, &HTTPStatusError{StatusCode: http.StatusUnauthorized, Message: "invalid authorization header format"}
	}

	endpoint := strings.TrimRight(config.Get().AUTH_API, "/") + "/users/me"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, &HTTPStatusError{StatusCode: http.StatusBadGateway, Message: "authentication service unavailable"}
	}

	req.Header.Set("Authorization", authorizationHeader)
	req.Header.Set("Accept", "application/json")

	resp, err := authHTTPClient.Do(req)
	if err != nil {
		return nil, &HTTPStatusError{StatusCode: http.StatusBadGateway, Message: "authentication service unavailable"}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, &HTTPStatusError{StatusCode: http.StatusUnauthorized, Message: "not logged in or invalid token"}
	}

	if resp.StatusCode == http.StatusBadRequest {
		return nil, &HTTPStatusError{StatusCode: http.StatusBadRequest, Message: "user not found"}
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, &HTTPStatusError{StatusCode: http.StatusBadGateway, Message: fmt.Sprintf("authentication service error: %d", resp.StatusCode)}
	}

	var payload GetCurrentUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, &HTTPStatusError{StatusCode: http.StatusBadGateway, Message: "invalid response from authentication service"}
	}

	if payload.Data.ID == "" {
		return nil, &HTTPStatusError{StatusCode: http.StatusBadGateway, Message: "invalid user payload from auth service"}
	}

	if payload.Data.Roles == nil {
		payload.Data.Roles = []string{}
	}

	return &payload.Data, nil
}
