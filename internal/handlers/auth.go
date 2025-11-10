package handlers

import (
	"net/http"

	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/internal/services"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService *services.AuthService
	validate    *validator.Validate
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	validate := validator.New()
	types.RegisterValidatorRegistrations(validate)
	
	return &AuthHandler{
		authService: authService,
		validate:    validate,
	}
}

// Login handles user login requests
func (h *AuthHandler) Login(c *gin.Context) {
	var loginData models.UserLogin
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid request data: " + err.Error()))
		return
	}

	if err := h.validate.Struct(loginData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Validation error: " + err.Error()))
		return
	}

	result, err := h.authService.Login(&loginData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// Register handles user registration requests
func (h *AuthHandler) Register(c *gin.Context) {
	var registerData models.UserRegister
	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid request data: " + err.Error()))
		return
	}

	if err := h.validate.Struct(registerData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Validation error: " + err.Error()))
		return
	}

	result, err := h.authService.Register(&registerData)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, result)
}

// Profile handles user profile requests
func (h *AuthHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.APIResponseWithError("Unauthorized"))
		return
	}

	result, err := h.authService.GetUserProfile(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError("Failed to retrieve profile: " + err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// ChangePassword handles password change requests
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.APIResponseWithError("Unauthorized"))
		return
	}

	var changePasswordData models.UserChangePassword
	if err := c.ShouldBindJSON(&changePasswordData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Invalid request data: " + err.Error()))
		return
	}

	if err := h.validate.Struct(changePasswordData); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("Validation error: " + err.Error()))
		return
	}

	err := h.authService.ChangePassword(userID.(string), &changePasswordData)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, types.APIResponseWithMessage("Password updated successfully"))
}

// Logout handles user logout requests
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a stateless JWT system, the server doesn't store session state
	// The client is responsible for discarding the token
	
	c.JSON(http.StatusOK, types.APIResponseWithMessage("Successfully logged out"))
}