package models

import (
	"time"

	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
)

// Expense represents a business expense
type Expense struct {
	ID          string            `json:"id" db:"id"`
	Category    string            `json:"category" db:"category" validate:"required,min=1,max=100"`
	Description string            `json:"description,omitempty" db:"description"`
	Amount      types.DecimalText `json:"amount" db:"amount" validate:"required,gt=0"`
	Date        time.Time         `json:"date" db:"date" validate:"required"`
	UserID      *string           `json:"user_id,omitempty" db:"user_id"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
}

// ExpenseCreate represents data to create a new expense
type ExpenseCreate struct {
	Category    string            `json:"category" validate:"required,min=1,max=100"`
	Description string            `json:"description,omitempty" validate:"omitempty,max=500"`
	Amount      types.DecimalText `json:"amount" validate:"required,gt=0"`
	Date        time.Time         `json:"date" validate:"required"`
}

// ExpenseUpdate represents data to update an existing expense
type ExpenseUpdate struct {
	Category    *string           `json:"category,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string           `json:"description,omitempty" validate:"omitempty,max=500"`
	Amount      *types.DecimalText `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Date        *time.Time        `json:"date,omitempty" validate:"omitempty"`
}

// ExpenseFilter represents filter options for listing expenses
type ExpenseFilter struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Category  *string    `json:"category,omitempty"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

// ExpenseSummary represents a summary of expenses for reporting
type ExpenseSummary struct {
	TotalAmount   types.DecimalText `json:"total_amount"`
	CategoryCount map[string]types.DecimalText `json:"category_count"`
}