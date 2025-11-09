package models

import (
	"time"

	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
)

// User represents a user in the POS system
type User struct {
	ID        string           `json:"id" db:"id"`
	Username  string           `json:"username" db:"username" validate:"required,min=3,max=50,alphanum"`
	Email     string           `json:"email" db:"email" validate:"required,email"`
	Password  string           `json:"-" db:"password_hash"`
	Role      types.UserRole   `json:"role" db:"role" validate:"required,oneof=cashier manager admin"`
	FirstName string           `json:"first_name" db:"first_name" validate:"required,min=1,max=100"`
	LastName  string           `json:"last_name" db:"last_name" validate:"required,min=1,max=100"`
	IsActive  bool             `json:"is_active" db:"is_active"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`
}

// UserLogin represents login credentials
type UserLogin struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

// UserRegister represents registration data
type UserRegister struct {
	Username  string         `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email     string         `json:"email" validate:"required,email"`
	Password  string         `json:"password" validate:"required,min=8"`
	Role      types.UserRole `json:"role" validate:"required,oneof=cashier manager admin"`
	FirstName string         `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string         `json:"last_name" validate:"required,min=1,max=100"`
}

// UserProfile represents public user information
type UserProfile struct {
	ID        string         `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Role      types.UserRole `json:"role"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// UserUpdate represents data to update a user
type UserUpdate struct {
	Username  *string        `json:"username,omitempty" validate:"omitempty,min=3,max=50,alphanum"`
	Email     *string        `json:"email,omitempty" validate:"omitempty,email"`
	Role      *types.UserRole `json:"role,omitempty" validate:"omitempty,oneof=cashier manager admin"`
	FirstName *string        `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName  *string        `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
}

// UserChangePassword represents password change data
type UserChangePassword struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}