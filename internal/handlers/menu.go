package handlers

import (
	"net/http"
	"strconv"

	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/internal/services"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// MenuHandler handles menu-related HTTP requests
type MenuHandler struct {
	menuService *services.MenuService
	validate    *validator.Validate
}

// NewMenuHandler creates a new menu handler
func NewMenuHandler(menuService *services.MenuService) *MenuHandler {
	validate := validator.New()
	types.RegisterValidatorRegistrations(validate)
	
	return &MenuHandler{
		menuService: menuService,
		validate:    validate,
	}
}

// CreateCategory handles category creation requests
func (h *MenuHandler) CreateCategory(c *gin.Context) {
	var categoryData models.CategoryCreate
	if err := c.ShouldBindJSON(&categoryData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid request data: " + err.Error()))
		return
	}

	if err := h.validate.Struct(categoryData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Validation error: " + err.Error()))
		return
	}

	result, err := h.menuService.CreateCategory(&categoryData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetCategory handles category retrieval requests
func (h *MenuHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	result, err := h.menuService.GetCategory(id)
	if err != nil {
		c.JSON(http.StatusNotFound, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// ListCategories handles category listing requests
func (h *MenuHandler) ListCategories(c *gin.Context) {
	// Get query parameters
	isActiveStr := c.DefaultQuery("active", "true")
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	isActive, err := strconv.ParseBool(isActiveStr)
	if err != nil {
		isActive = true // Default to true if not provided or invalid
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50 // Default to 50 if not provided or invalid
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0 // Default to 0 if not provided or invalid
	}

	result, err := h.menuService.ListCategories(isActive, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateCategory handles category update requests
func (h *MenuHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	
	var updateData models.CategoryUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid request data: " + err.Error()))
		return
	}

	result, err := h.menuService.UpdateCategory(id, &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteCategory handles category deletion requests
func (h *MenuHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	
	result, err := h.menuService.DeleteCategory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// CreateMenuItem handles menu item creation requests
func (h *MenuHandler) CreateMenuItem(c *gin.Context) {
	var itemData models.MenuItemCreate
	if err := c.ShouldBindJSON(&itemData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid request data: " + err.Error()))
		return
	}

	if err := h.validate.Struct(itemData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Validation error: " + err.Error()))
		return
	}

	result, err := h.menuService.CreateMenuItem(&itemData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetMenuItem handles menu item retrieval requests
func (h *MenuHandler) GetMenuItem(c *gin.Context) {
	id := c.Param("id")
	result, err := h.menuService.GetMenuItem(id)
	if err != nil {
		c.JSON(http.StatusNotFound, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// ListMenuItems handles menu item listing requests
func (h *MenuHandler) ListMenuItems(c *gin.Context) {
	// Get query parameters
	categoryID := c.Query("category_id")
	isAvailableStr := c.DefaultQuery("is_available", "true")
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	isAvailable, err := strconv.ParseBool(isAvailableStr)
	if err != nil {
		isAvailable = true // Default to true if not provided or invalid
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50 // Default to 50 if not provided or invalid
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0 // Default to 0 if not provided or invalid
	}

	var result *types.APIResponse
	if categoryID != "" {
		// List items by category
		result, err = h.menuService.ListMenuItemsByCategory(categoryID, limit, offset)
	} else {
		// List all items
		result, err = h.menuService.ListMenuItems(isAvailable, limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateMenuItem handles menu item update requests
func (h *MenuHandler) UpdateMenuItem(c *gin.Context) {
	id := c.Param("id")
	
	var updateData models.MenuItemUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid request data: " + err.Error()))
		return
	}

	result, err := h.menuService.UpdateMenuItem(id, &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteMenuItem handles menu item deletion requests
func (h *MenuHandler) DeleteMenuItem(c *gin.Context) {
	id := c.Param("id")
	
	result, err := h.menuService.DeleteMenuItem(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}