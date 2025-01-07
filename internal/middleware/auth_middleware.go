package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"point-system-api/pkg/utils"
)

// AuthMiddleware is a middleware that checks for a valid JWT token in the request header.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Check if the header is in the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		// Extract the token
		token := parts[1]

		// Validate the token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Set the user ID and role in the context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		// Continue to the next handler
		c.Next()
	}
}

// RoleMiddleware is a middleware that checks if the user has the required role to access a route.
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user role from the context
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			return
		}

		// Check if the user has the required role
		if userRole != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		// Continue to the next handler
		c.Next()
	}
}
