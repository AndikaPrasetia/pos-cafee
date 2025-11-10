package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/internal/repositories"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/AndikaPrasetia/pos-cafee/pkg/utils"
)

// AuthService handles user authentication and authorization
type AuthService struct {
	userRepo repositories.UserRepo
	jwtSecret string
	jwtExpiry time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repositories.UserRepo, jwtSecret string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(loginData *models.UserLogin) (*types.APIResponse, error) {
	user, err := s.userRepo.GetUserByUsername(loginData.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if !utils.CheckPasswordHash(loginData.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, string(user.Role), s.jwtSecret, s.jwtExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data: map[string]any{
			"user": map[string]any{
				"id":         user.ID,
				"username":   user.Username,
				"email":      user.Email,
				"role":       user.Role,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
			},
			"token":      token,
			"expires_in": int(s.jwtExpiry.Seconds()),
		},
	}, nil
}

// Register creates a new user account
func (s *AuthService) Register(registerData *models.UserRegister) (*types.APIResponse, error) {
	// Validate inputs
	if err := utils.ValidateUsername(registerData.Username); err != nil {
		return nil, err
	}

	if err := utils.ValidateEmail(registerData.Email); err != nil {
		return nil, err
	}

	// Check if username already exists
	_, err := s.userRepo.GetUserByUsername(registerData.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	_, err = s.userRepo.GetUserByEmail(registerData.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(registerData.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	user := &models.User{
		Username:  registerData.Username,
		Email:     registerData.Email,
		Password:  hashedPassword,
		Role:      registerData.Role,
		FirstName: registerData.FirstName,
		LastName:  registerData.LastName,
		IsActive:  true,
	}

	createdUser, err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	token, err := utils.GenerateJWT(createdUser.ID, string(createdUser.Role), s.jwtSecret, s.jwtExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data: map[string]any{
			"user": map[string]any{
				"id":         createdUser.ID,
				"username":   createdUser.Username,
				"email":      createdUser.Email,
				"role":       createdUser.Role,
				"first_name": createdUser.FirstName,
				"last_name":  createdUser.LastName,
			},
			"token":      token,
			"expires_in": int(s.jwtExpiry.Seconds()),
		},
	}, nil
}

// GetUserProfile returns user profile information
func (s *AuthService) GetUserProfile(userID string) (*types.APIResponse, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return &types.APIResponse{
		Success: true,
		Data: map[string]any{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"role":       user.Role,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"is_active":  user.IsActive,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	}, nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(userID string, changePasswordData *models.UserChangePassword) error {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(changePasswordData.CurrentPassword, user.Password) {
		return errors.New("current password is incorrect")
	}

	newHashedPassword, err := utils.HashPassword(changePasswordData.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %v", err)
	}

	return s.userRepo.UpdateUserPassword(userID, newHashedPassword)
}

// ValidateUserAccess validates if a user has access to a specific resource based on role
func (s *AuthService) ValidateUserAccess(userID, requiredRole string) (bool, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return false, err
	}

	if !user.IsActive {
		return false, errors.New("user account is deactivated")
	}

	// Check if the user has the required role or higher
	userRole := string(user.Role)

	// Admin can access everything
	if userRole == "admin" {
		return true, nil
	}

	// Manager can access most things except admin-specific features
	if userRole == "manager" && requiredRole != "admin" {
		return true, nil
	}

	// Cashier has limited access, mostly to POS functions
	if userRole == "cashier" && (requiredRole == "cashier") {
		return true, nil
	}

	return false, errors.New("insufficient permissions")
}
