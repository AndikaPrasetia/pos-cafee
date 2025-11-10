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

// OrderHandler handles order-related HTTP requests
type OrderHandler struct {
	orderService *services.OrderService
	validate    *validator.Validate
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	validate := validator.New()
	types.RegisterValidatorRegistrations(validate)
	
	return &OrderHandler{
		orderService: orderService,
		validate:    validate,
	}
}

// CreateOrder handles draft order creation requests
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.APIResponseWithError("Unauthorized"))
		return
	}

	var orderData models.OrderCreate
	if err := c.ShouldBindJSON(&orderData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid request data: " + err.Error()))
		return
	}

	if err := h.validate.Struct(orderData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Validation error: " + err.Error()))
		return
	}

	result, err := h.orderService.CreateOrder(userID.(string), &orderData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetOrder handles order retrieval requests
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	result, err := h.orderService.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// ListOrders handles order listing requests with optional filtering
func (h *OrderHandler) ListOrders(c *gin.Context) {
	var filter types.OrderFilter

	// Get filter parameters from query
	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}
	
	if userID := c.Query("user_id"); userID != "" {
		filter.UserID = &userID
	}
	
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}
	
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = &endDate
		}
	}

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

	result, err := h.orderService.ListOrders(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// AddItemToOrder handles adding an item to an existing order
func (h *OrderHandler) AddItemToOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.APIResponseWithError("Unauthorized"))
		return
	}

	var itemData models.OrderItemCreate
	if err := c.ShouldBindJSON(&itemData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid request data: " + err.Error()))
		return
	}

	if err := h.validate.Struct(itemData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Validation error: " + err.Error()))
		return
	}

	result, err := h.orderService.AddItemToOrder(orderID, userID.(string), &itemData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// CompleteOrder handles order completion and payment processing
func (h *OrderHandler) CompleteOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.APIResponseWithError("Unauthorized"))
		return
	}

	var updateData models.OrderUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid request data: " + err.Error()))
		return
	}

	// Validation for completion is different, so we'll do minimal validation here
	if updateData.PaymentMethod != nil {
		pm := *updateData.PaymentMethod
		if pm != types.PaymentMethodCash && pm != types.PaymentMethodCard && 
		   pm != types.PaymentMethodQris && pm != types.PaymentMethodTransfer {
			c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid payment method"))
			return
		}
	}

	result, err := h.orderService.CompleteOrder(orderID, userID.(string), &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// CancelOrder handles order cancellation requests
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.APIResponseWithError("Unauthorized"))
		return
	}

	var updateData models.OrderUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid request data: " + err.Error()))
		return
	}

	result, err := h.orderService.CancelOrder(orderID, userID.(string), &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}