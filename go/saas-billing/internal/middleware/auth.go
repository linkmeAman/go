package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linkmeAman/saas-billing/internal/auth"
)

// AuthRequired verifies JWT token and adds claims to context
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(bearerToken[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Add claims to context
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// RequireRole checks if the user has the required role in the organization
func RequireRole(orgService interface {
	CheckUserRole(userID, orgID string) (string, error)
}, requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		orgID := c.Param("orgID")

		if orgID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID required"})
			c.Abort()
			return
		}

		role, err := orgService.CheckUserRole(userID, orgID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is not a member of this organization"})
			c.Abort()
			return
		}

		// Check if user's role is in the required roles
		hasRole := false
		for _, r := range requiredRoles {
			if role == r {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Set("userRole", role)
		c.Next()
	}
}
