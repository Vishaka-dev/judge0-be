package middlewares

import (
	"errors"
	"net/http"
	"strings"

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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		user, err := utils.GetCurrentUser(c.Request.Context(), auth)
		if err != nil {
			var statusErr *utils.HTTPStatusError
			if errors.As(err, &statusErr) {
				c.AbortWithStatusJSON(statusErr.StatusCode, gin.H{"error": statusErr.Message})
				return
			}

			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "authentication service unavailable"})
			return
		}

		if len(allowedRoles) > 0 && !hasAllowedRole(user.Roles, allowedRoles) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Set("user_name", user.Name)
		c.Set("user_roles", user.Roles)
		c.Next()
	}
}
