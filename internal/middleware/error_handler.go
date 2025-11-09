package middleware

import (
	"net/http"

	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware handles errors and formats them consistently
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are errors in the gin context
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last()
			
			// Log the error (in a real implementation, you'd log to a proper logger)
			// log.Printf("Error occurred: %v", err)

			// Return a consistent error response
			c.JSON(http.StatusInternalServerError, types.APIResponse{
				Success: false,
				Error:   err.Error(),
			})
		}
	}
}