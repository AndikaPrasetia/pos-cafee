package middleware

import (
	"net/http"
	"strings"

	"github.com/AndikaPrasetia/pos-cafee/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a middleware for JWT authentication
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if the header follows the "Bearer {token}" format
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header must be in the format 'Bearer {token}'",
			})
			c.Abort()
			return
		}

		tokenString := authParts[1]

		// Parse and validate the token
		claims, err := utils.ParseJWT(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Check if token is expired
		isExpired, err := utils.IsTokenExpired(tokenString, jwtSecret)
		if err != nil || isExpired {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token has expired",
			})
			c.Abort()
			return
		}

		// Set user ID, username and role in the context for use in handlers
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		// Continue with the request
		c.Next()
	}
}

// RoleAuthMiddleware creates a middleware that checks user role
func RoleAuthMiddleware(jwtSecret string, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header must be in the format 'Bearer {token}'",
			})
			c.Abort()
			return
		}

		tokenString := authParts[1]
		claims, err := utils.ParseJWT(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Check role
		userRole := claims.Role
		switch requiredRole {
		case "admin":
			if userRole != "admin" {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"message": "Admin role required for this action",
				})
				c.Abort()
				return
			}
		case "manager":
			if userRole != "admin" && userRole != "manager" {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"message": "Manager or admin role required for this action",
				})
				c.Abort()
				return
			}
		case "cashier":
			// Cashier role or higher (manager, admin) allowed
			if userRole != "cashier" && userRole != "manager" && userRole != "admin" {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"message": "Cashier, manager, or admin role required for this action",
				})
				c.Abort()
				return
			}
		}

		// Set user ID, username and role in the context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		// Continue with the request
		c.Next()
	}
}