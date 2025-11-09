package models

import (
	"time"

	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
)

// Category represents a menu category
type Category struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required,min=1,max=100"`
	Description *string   `json:"description,omitempty" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CategoryCreate represents data to create a category
type CategoryCreate struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
}

// CategoryUpdate represents data to update a category
type CategoryUpdate struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// MenuItem represents a menu item
type MenuItem struct {
	ID          string            `json:"id" db:"id"`
	Name        string            `json:"name" db:"name" validate:"required,min=1,max=255"`
	CategoryID  string            `json:"category_id" db:"category_id" validate:"required,uuid"`
	Description *string           `json:"description,omitempty" db:"description"`
	Price       types.DecimalText `json:"price" db:"price" validate:"required,gt=0"`
	Cost        types.DecimalText `json:"cost" db:"cost" validate:"required,gt=0,ltefield=Price"`
	IsAvailable bool              `json:"is_available" db:"is_available"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
}

// MenuItemCreate represents data to create a menu item
type MenuItemCreate struct {
	Name        string            `json:"name" validate:"required,min=1,max=255"`
	CategoryID  string            `json:"category_id" validate:"required,uuid"`
	Description *string           `json:"description,omitempty" validate:"omitempty,max=500"`
	Price       types.DecimalText `json:"price" validate:"required,gt=0"`
	Cost        types.DecimalText `json:"cost" validate:"required,gt=0,ltefield=Price"`
	IsAvailable bool              `json:"is_available,omitempty"`
}

// MenuItemUpdate represents data to update a menu item
type MenuItemUpdate struct {
	Name        *string            `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	CategoryID  *string            `json:"category_id,omitempty" validate:"omitempty,uuid"`
	Description *string            `json:"description,omitempty" validate:"omitempty,max=500"`
	Price       *types.DecimalText `json:"price,omitempty" validate:"omitempty,gt=0"`
	Cost        *types.DecimalText `json:"cost,omitempty" validate:"omitempty,gt=0,ltefield=Price"`
	IsAvailable *bool              `json:"is_available,omitempty"`
}

// MenuItemWithCategory represents a menu item with its category name
type MenuItemWithCategory struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	CategoryID    string            `json:"category_id"`
	CategoryName  string            `json:"category_name"`
	Description   *string           `json:"description,omitempty"`
	Price         types.DecimalText `json:"price"`
	Cost          types.DecimalText `json:"cost"`
	IsAvailable   bool              `json:"is_available"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}