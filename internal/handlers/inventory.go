package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/internal/services"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// InventoryHandler handles inventory-related HTTP requests
type InventoryHandler struct {
	inventoryService *services.InventoryService
	validate         *validator.Validate
}

// NewInventoryHandler creates a new inventory handler
func NewInventoryHandler(inventoryService *services.InventoryService) *InventoryHandler {
	validate := validator.New()
	types.RegisterValidatorRegistrations(validate)

	return &InventoryHandler{
		inventoryService: inventoryService,
		validate:         validate,
	}
}

// GetInventory retrieves inventory information for a specific menu item
func (h *InventoryHandler) GetInventory(c *gin.Context) {
	menuItemID := c.Query("menu_item_id")
	if menuItemID == "" {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("menu_item_id parameter is required"))
		return
	}

	result, err := h.inventoryService.GetInventoryByMenuItem(menuItemID)
	if err != nil {
		c.JSON(http.StatusNotFound, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// ListInventory retrieves a list of inventory items with optional filtering
func (h *InventoryHandler) ListInventory(c *gin.Context) {
	var filter models.InventoryFilter

	// Get query parameters
	lowStockOnlyStr := c.DefaultQuery("low_stock_only", "false")
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	lowStockOnly, err := strconv.ParseBool(lowStockOnlyStr)
	if err != nil {
		lowStockOnly = false
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	filter.LowStockOnly = lowStockOnly
	filter.Limit = limit
	filter.Offset = offset

	result, err := h.inventoryService.ListInventory(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateInventory handles manual stock adjustment requests
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.APIResponseWithError("Unauthorized"))
		return
	}

	var updateData models.InventoryUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid request data: "+err.Error()))
		return
	}

	if err := h.validate.Struct(updateData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Validation error: "+err.Error()))
		return
	}

	// Check that quantity is not zero
	if updateData.Quantity == 0 {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Quantity must not be zero"))
		return
	}

	result, err := h.inventoryService.UpdateStock(userID.(string), &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// ListStockTransactions retrieves a list of stock transactions with optional filtering
func (h *InventoryHandler) ListStockTransactions(c *gin.Context) {
	var filter models.StockTransactionFilter

	// Get query parameters
	menuItemID := c.Query("menu_item_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	// Set default values for optional parameters
	if menuItemID != "" {
		filter.MenuItemID = &menuItemID
	} else {
		// Use default UUID when not provided
		defaultUUID := "00000000-0000-0000-0000-000000000000"
		filter.MenuItemID = &defaultUUID
	}

	// Parse date parameters if provided
	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid start date format, expected YYYY-MM-DD"))
			return
		}
		filter.StartDate = &startDate
	} else {
		// Use default date when not provided
		defaultDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
		filter.StartDate = &defaultDate
	}

	if endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid end date format, expected YYYY-MM-DD"))
			return
		}
		filter.EndDate = &endDate
	} else {
		// Use default date when not provided
		defaultDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
		filter.EndDate = &defaultDate
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	filter.Limit = limit
	filter.Offset = offset

	result, err := h.inventoryService.ListStockTransactions(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return

}

	c.JSON(http.StatusOK, result)
}

// GetLowStockItems retrieves items with stock below minimum threshold
func (h *InventoryHandler) GetLowStockItems(c *gin.Context) {
	var filter models.InventoryFilter
	filter.LowStockOnly = true

	// Get pagination parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	filter.Limit = limit
	filter.Offset = offset

	result, err := h.inventoryService.ListInventory(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}
