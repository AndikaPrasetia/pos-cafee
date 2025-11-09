package models

import (
	"time"

	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
)

// Inventory represents inventory tracking for a menu item
type Inventory struct {
	ID             string    `json:"id" db:"id"`
	MenuItemID     string    `json:"menu_item_id" db:"menu_item_id"`
	CurrentStock   int       `json:"current_stock" db:"current_stock"`
	MinimumStock   int       `json:"minimum_stock" db:"minimum_stock"`
	Unit           string    `json:"unit" db:"unit"`
	LastUpdatedAt  time.Time `json:"last_updated_at" db:"last_updated_at"`
	LastUpdatedBy  *string   `json:"last_updated_by,omitempty" db:"last_updated_by"`
	MenuItemName   string    `json:"menu_item_name,omitempty"`
	IsLowStock     bool      `json:"is_low_stock,omitempty"`
	LastUpdatedByName *string `json:"last_updated_by_name,omitempty"`
}

// InventoryUpdate represents data to update inventory stock
type InventoryUpdate struct {
	MenuItemID    string `json:"menu_item_id" validate:"required,uuid"`
	Quantity      int    `json:"quantity" validate:"required,ne=0"` // Can be positive or negative
	Reason        string `json:"reason" validate:"required,min=1,max=255"`
}

// StockTransaction represents a stock transaction record
type StockTransaction struct {
	ID              string                    `json:"id" db:"id"`
	MenuItemID      string                    `json:"menu_item_id" db:"menu_item_id"`
	MenuItemName    string                    `json:"menu_item_name,omitempty" db:"menu_item_name"`
	TransactionType types.TransactionType     `json:"transaction_type" db:"transaction_type"`
	Quantity        int                       `json:"quantity" db:"quantity"`
	PreviousStock   int                       `json:"previous_stock" db:"previous_stock"`
	CurrentStock    int                       `json:"current_stock" db:"current_stock"`
	Reason          string                    `json:"reason" db:"reason"`
	ReferenceType   *string                   `json:"reference_type,omitempty" db:"reference_type"`
	ReferenceID     *string                   `json:"reference_id,omitempty" db:"reference_id"`
	UserID          *string                   `json:"user_id,omitempty" db:"user_id"`
	UserName        *string                   `json:"user_name,omitempty"`
	CreatedAt       time.Time                 `json:"created_at" db:"created_at"`
}

// InventoryFilter represents filter options for listing inventory
type InventoryFilter struct {
	LowStockOnly bool `json:"low_stock_only"`
	Limit        int  `json:"limit"`
	Offset       int  `json:"offset"`
}

// StockTransactionFilter represents filter options for listing stock transactions
type StockTransactionFilter struct {
	MenuItemID *string    `json:"menu_item_id,omitempty"`
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset"`
}