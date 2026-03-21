package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/logger"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/utils"

	"github.com/gin-gonic/gin"
)

func hasAllowedRole(userRoles []string, allowedRoles []string) bool {
	for _, requiredRole := range allowedRoles {
		for _, userRole := range userRoles {
			if strings.EqualFold(requiredRole, userRole) {
				return true
			}
		}
	}
	return false
}

func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			logger.Log.Warn("Missing Authorization header", "path", c.Request.URL.Path, "method", c.Request.Method)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		user, err := utils.GetCurrentUser(c.Request.Context(), auth)
		if err != nil {
			var statusErr *utils.HTTPStatusError
			if errors.As(err, &statusErr) {
				logger.Log.Warn("Auth error", "message", statusErr.Message, "status", statusErr.StatusCode, "path", c.Request.URL.Path, "method", c.Request.Method)
				c.AbortWithStatusJSON(statusErr.StatusCode, gin.H{"error": statusErr.Message})
				return
			}
			logger.Log.Error("Authentication service unavailable", "error", err, "path", c.Request.URL.Path, "method", c.Request.Method)
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "authentication service unavailable"})
			return
		}

		if len(allowedRoles) > 0 && !hasAllowedRole(user.Roles, allowedRoles) {
			logger.Log.Warn("Insufficient permissions", "user_id", user.ID, "user_email", user.Email, "required_roles", allowedRoles, "user_roles", user.Roles, "path", c.Request.URL.Path, "method", c.Request.Method)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		logger.Log.Info("Authenticated user", "user_id", user.ID, "user_email", user.Email, "user_roles", user.Roles, "path", c.Request.URL.Path, "method", c.Request.Method)
		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Set("user_name", user.Name)
		c.Set("user_roles", user.Roles)
		c.Next()
	}
}
