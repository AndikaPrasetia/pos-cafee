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

// UserHandler handles user management related HTTP requests
type UserHandler struct {
	authService *services.AuthService
	userService *services.UserService
	validate    *validator.Validate
}

// NewUserHandler creates a new user management handler
func NewUserHandler(authService *services.AuthService, userService *services.UserService) *UserHandler {
	validate := validator.New()
	types.RegisterValidatorRegistrations(validate)

	return &UserHandler{
		authService: authService,
		userService: userService,
		validate:    validate,
	}
}

// ListUsers handles requests to list all users with pagination
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Extract query parameters for pagination
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	result, err := h.userService.ListUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError("Failed to list users: " + err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUser handles requests to get a specific user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	result, err := h.userService.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, types.APIResponseWithError("User not found: " + err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateUser handles requests to update a user's information
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var updateData models.UserUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid request data: " + err.Error()))
		return
	}

	if err := h.validate.Struct(updateData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Validation error: " + err.Error()))
		return
	}

	result, err := h.userService.UpdateUser(userID, &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError("Failed to update user: " + err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeactivateUser handles requests to deactivate a user account
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	userID := c.Param("id")

	err := h.userService.DeactivateUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError("Failed to deactivate user: " + err.Error()))
		return
	}

	c.JSON(http.StatusOK, types.APIResponseWithMessage("User deactivated successfully"))
}

// ActivateUser handles requests to activate a user account
func (h *UserHandler) ActivateUser(c *gin.Context) {
	userID := c.Param("id")

	err := h.userService.ActivateUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError("Failed to activate user: " + err.Error()))
		return
	}

	c.JSON(http.StatusOK, types.APIResponseWithMessage("User activated successfully"))
}