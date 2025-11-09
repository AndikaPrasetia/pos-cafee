package repositories

import (
	"github.com/AndikaPrasetia/pos-cafee/internal/db"
)

// expenseRepo implements the ExpenseRepo interface
type expenseRepo struct {
	queries *db.Queries
}

// Additional expense-related methods will be implemented here