package middleware

import (
	"net/http"
	"strconv"

	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware adds CORS headers to responses
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// PaginationMiddleware extracts pagination parameters from query strings
func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get limit and offset from query parameters
		limitStr := c.DefaultQuery("limit", "50")
		offsetStr := c.DefaultQuery("offset", "0")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 50
		}
		// Set a maximum limit to prevent abuse
		if limit > 1000 {
			limit = 1000
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			offset = 0
		}

		// Store in context for handlers to use
		c.Set("limit", limit)
		c.Set("offset", offset)

		c.Next()
	}
}

// ErrorMiddleware handles errors and formats them consistently
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are errors in the gin context
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		}
	}
}

// RequestLoggingMiddleware logs incoming requests
func RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log the request
		logRequest(c)

		// Process the request
		c.Next()

		// Log the response
		logResponse(c)
	}
}

// logRequest logs the incoming request details
func logRequest(c *gin.Context) {
	// In a real implementation, you would log to a file or structured logger
	// This is a placeholder implementation
}

// logResponse logs the response details
func logResponse(c *gin.Context) {
	// In a real implementation, you would log to a file or structured logger
	// This is a placeholder implementation
}

// RateLimitMiddleware provides basic rate limiting
func RateLimitMiddleware(maxRequests int, windowSizeInSeconds int) gin.HandlerFunc {
	// In a complete implementation, this would track requests per user/IP and enforce limits
	// For now, we'll just return a basic middleware that passes through
	return func(c *gin.Context) {
		c.Next()
	}
}