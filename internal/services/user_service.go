package services

import (
	"errors"
	"fmt"

	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/internal/repositories"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/AndikaPrasetia/pos-cafee/pkg/utils"
	"github.com/google/uuid"
)

// UserService handles user management business logic
type UserService struct {
	userRepo repositories.UserRepo
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// ListUsers retrieves a paginated list of users
func (s *UserService) ListUsers(page, limit int) (*types.APIResponse, error) {
	// Validate inputs
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get users with pagination
	users, err := s.userRepo.ListUsers(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %v", err)
	}

	// Count total users for pagination metadata
	total, err := s.userRepo.CountUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"users": users,
			"pagination": map[string]interface{}{
				"page":  page,
				"limit": limit,
				"total": total,
				"pages": (total + limit - 1) / limit,
			},
		},
	}, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id string) (*types.APIResponse, error) {
	// Validate UUID
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := s.userRepo.GetUser(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data:    user,
	}, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(id string, updateData *models.UserUpdate) (*types.APIResponse, error) {
	// Validate UUID
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := s.userRepo.GetUser(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Update fields if provided in updateData
	if updateData.FirstName != nil {
		user.FirstName = *updateData.FirstName
	}
	if updateData.LastName != nil {
		user.LastName = *updateData.LastName
	}
	if updateData.Email != nil {
		// Check if email already exists for another user
		_, err := s.userRepo.GetUserByEmail(*updateData.Email)
		if err == nil {
			// User with this email already exists
			return nil, errors.New("email already exists")
		}
		user.Email = *updateData.Email
	}
	if updateData.Role != nil {
		user.Role = *updateData.Role
	}

	updatedUser, err := s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data:    updatedUser,
	}, nil
}

// DeactivateUser deactivates a user account
func (s *UserService) DeactivateUser(id string) error {
	// Validate UUID
	_, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	user, err := s.userRepo.GetUser(id)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	if !user.IsActive {
		return errors.New("user is already deactivated")
	}

	user.IsActive = false
	_, err = s.userRepo.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %v", err)
	}

	return nil
}

// ActivateUser activates a user account
func (s *UserService) ActivateUser(id string) error {
	// Validate UUID
	_, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	user, err := s.userRepo.GetUser(id)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	if user.IsActive {
		return errors.New("user is already active")
	}

	user.IsActive = true
	_, err = s.userRepo.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("failed to activate user: %v", err)
	}

	return nil
}

// ChangeUserPassword changes a user's password (admin function)
func (s *UserService) ChangeUserPassword(userID string, newPassword string) error {
	// Validate UUID
	_, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	err = s.userRepo.UpdateUserPassword(userID, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update user password: %v", err)
	}

	return nil
}