package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AndikaPrasetia/pos-cafee/internal/db"
	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/google/uuid"
)

// userRepo implements the UserRepo interface
type userRepo struct {
	queries *db.Queries
}

// GetUser retrieves a user by ID
func (r *userRepo) GetUser(id string) (*models.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := r.queries.GetUser(context.Background(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &models.User{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.PasswordHash,
		Role:      types.UserRole(user.Role),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUserByUsername retrieves a user by username
func (r *userRepo) GetUserByUsername(username string) (*models.User, error) {
	user, err := r.queries.GetUserByUsername(context.Background(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &models.User{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.PasswordHash,
		Role:      types.UserRole(user.Role),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUserByEmail retrieves a user by email
func (r *userRepo) GetUserByEmail(email string) (*models.User, error) {
	user, err := r.queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &models.User{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.PasswordHash,
		Role:      types.UserRole(user.Role),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// CreateUser creates a new user
func (r *userRepo) CreateUser(user *models.User) (*models.User, error) {
	createdUser, err := r.queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:  user.Username,
		Email:     user.Email,
		PasswordHash: user.Password,
		Role:      string(user.Role),
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        createdUser.ID.String(),
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		Role:      types.UserRole(createdUser.Role),
		FirstName: createdUser.FirstName,
		LastName:  createdUser.LastName,
		IsActive:  createdUser.IsActive,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}, nil
}

// UpdateUser updates an existing user
func (r *userRepo) UpdateUser(user *models.User) (*models.User, error) {
	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	updatedUser, err := r.queries.UpdateUser(context.Background(), db.UpdateUserParams{
		ID:        userID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      string(user.Role),
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        updatedUser.ID.String(),
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		Role:      types.UserRole(updatedUser.Role),
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		IsActive:  updatedUser.IsActive,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}

// UpdateUserPassword updates a user's password
func (r *userRepo) UpdateUserPassword(userID, hashedPassword string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	return r.queries.UpdateUserPassword(context.Background(), db.UpdateUserPasswordParams{
		ID:           userUUID,
		PasswordHash: hashedPassword,
	})
}

// UpdateUserStatus updates a user's active status
func (r *userRepo) UpdateUserStatus(userID string, isActive bool) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	return r.queries.UpdateUserStatus(context.Background(), db.UpdateUserStatusParams{
		ID:       userUUID,
		IsActive: isActive,
	})
}