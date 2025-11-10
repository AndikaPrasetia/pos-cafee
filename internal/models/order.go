package models

import (
	"time"

	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
)

// Order represents a customer order
type Order struct {
	ID             string               `json:"id" db:"id"`
	OrderNumber    string               `json:"order_number" db:"order_number"`
	UserID         string               `json:"user_id" db:"user_id"`
	Status         types.OrderStatus    `json:"status" db:"status"`
	TotalAmount    types.DecimalText    `json:"total_amount" db:"total_amount"`
	DiscountAmount types.DecimalText    `json:"discount_amount" db:"discount_amount"`
	TaxAmount      types.DecimalText    `json:"tax_amount" db:"tax_amount"`
	PaymentMethod  *types.PaymentMethod `json:"payment_method,omitempty" db:"payment_method"`
	PaymentStatus  types.PaymentStatus  `json:"payment_status" db:"payment_status"`
	CompletedAt    *time.Time           `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt      time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at" db:"updated_at"`
}

// OrderCreate represents data to create a draft order
type OrderCreate struct {
	Items []OrderItemCreate `json:"items" validate:"required,min=1,dive"`
}

// OrderUpdate represents data to update an order
type OrderUpdate struct {
	PaymentMethod  *types.PaymentMethod `json:"payment_method,omitempty" validate:"omitempty,oneof=cash card qris transfer"`
	DiscountAmount *types.DecimalText   `json:"discount_amount,omitempty" validate:"omitempty,gt=0"`
	TaxAmount      *types.DecimalText   `json:"tax_amount,omitempty" validate:"omitempty,gt=0"`
	Reason         *string              `json:"reason,omitempty"` // For cancellation
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID         string            `json:"id" db:"id"`
	OrderID    string            `json:"order_id" db:"order_id"`
	MenuItemID string            `json:"menu_item_id" db:"menu_item_id"`
	Quantity   int               `json:"quantity" db:"quantity" validate:"required,gt=0"`
	UnitPrice  types.DecimalText `json:"unit_price" db:"unit_price"`
	TotalPrice types.DecimalText `json:"total_price" db:"total_price"`
}

// OrderItemCreate represents data to create an order item
type OrderItemCreate struct {
	MenuItemID string `json:"menu_item_id" validate:"required,uuid"`
	Quantity   int    `json:"quantity" validate:"required,gt=0"`
}

// OrderItemWithDetails represents an order item with menu item details
type OrderItemWithDetails struct {
	ID           string            `json:"id"`
	OrderID      string            `json:"order_id"`
	MenuItemID   string            `json:"menu_item_id"`
	MenuItemName string            `json:"menu_item_name"`
	Quantity     int               `json:"quantity"`
	UnitPrice    types.DecimalText `json:"unit_price"`
	TotalPrice   types.DecimalText `json:"total_price"`
}

// OrderWithDetails represents an order with user and item details
type OrderWithDetails struct {
	ID             string                 `json:"id"`
	OrderNumber    string                 `json:"order_number"`
	UserID         string                 `json:"user_id"`
	UserName       string                 `json:"user_name"`
	Status         types.OrderStatus      `json:"status"`
	TotalAmount    types.DecimalText      `json:"total_amount"`
	DiscountAmount types.DecimalText      `json:"discount_amount"`
	TaxAmount      types.DecimalText      `json:"tax_amount"`
	PaymentMethod  *types.PaymentMethod   `json:"payment_method,omitempty"`
	PaymentStatus  types.PaymentStatus    `json:"payment_status"`
	CompletedAt    *time.Time             `json:"completed_at,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	Items          []OrderItemWithDetails `json:"items"`
}

