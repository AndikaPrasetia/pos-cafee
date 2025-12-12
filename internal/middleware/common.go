package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/AndikaPrasetia/pos-cafee/pkg/utils"
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
		// Start timer to measure request duration
		start := time.Now()

		// Read request body if needed (be careful with this in production for performance)
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			// Restore the io.ReadCloser to its original state
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Log the request
		logRequest(c, bodyBytes)

		// Process the request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log the response
		logResponse(c, duration)
	}
}

// logRequest logs the incoming request details
func logRequest(c *gin.Context, body []byte) {
	utils.LogInfo("Incoming request", map[string]any{
		"method":      c.Request.Method,
		"uri":         c.Request.RequestURI,
		"remote_addr": c.ClientIP(),
		"user_agent":  c.Request.UserAgent(),
		"body":        string(body),
		"timestamp":   time.Now().Format(time.RFC3339),
		"request_id":  generateRequestID(), // You can implement request ID generation if needed
	})
}

// logResponse logs the response details
func logResponse(c *gin.Context, duration time.Duration) {
	// Get response size from context or headers
	responseSize := c.Writer.Size()
	if responseSize < 0 {
		responseSize = 0
	}

	utils.LogInfo("Request completed", map[string]any{
		"method":         c.Request.Method,
		"uri":            c.Request.RequestURI,
		"status":         c.Writer.Status(),
		"response_size":  responseSize,
		"duration_ms":    duration.Milliseconds(),
		"duration":       duration.String(),
		"remote_addr":    c.ClientIP(),
		"timestamp":      time.Now().Format(time.RFC3339),
		"request_id":     generateRequestID(),
	})
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	return utils.GenerateUUID()
}

// RateLimitMiddleware provides basic rate limiting
func RateLimitMiddleware(maxRequests int, windowSizeInSeconds int) gin.HandlerFunc {
	// In a complete implementation, this would track requests per user/IP and enforce limits
	// For now, we'll just return a basic middleware that passes through
	return func(c *gin.Context) {
		c.Next()
	}
}