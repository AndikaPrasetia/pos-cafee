package handlers

import (
	"net/http"

	"github.com/AndikaPrasetia/pos-cafee/internal/cache"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/gin-gonic/gin"
)

// SystemHandler handles system-level operations and monitoring
type SystemHandler struct {
	redisManager *cache.RedisManager
}

// NewSystemHandler creates a new SystemHandler
func NewSystemHandler(redisManager *cache.RedisManager) *SystemHandler {
	return &SystemHandler{
		redisManager: redisManager,
	}
}

// GetRedisHealth returns the health status of Redis
func (h *SystemHandler) GetRedisHealth(c *gin.Context) {
	isHealthy := h.redisManager.IsHealthy()

	status := http.StatusOK
	if !isHealthy {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"success": true,
		"data": gin.H{
			"redis_healthy": isHealthy,
		},
	})
}

// GetRedisMetrics returns Redis metrics for monitoring
func (h *SystemHandler) GetRedisMetrics(c *gin.Context) {
	metrics, err := h.redisManager.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponse{
			Success: false,
			Message: "Failed to get Redis metrics",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    metrics,
	})
}