package repositories

import (
	"errors"

	"github.com/AndikaPrasetia/pos-cafee/internal/db"
	"github.com/AndikaPrasetia/pos-cafee/internal/models"
)

// stockTransactionRepo implements the StockTransactionRepo interface
type stockTransactionRepo struct {
	queries *db.Queries
}

// CreateStockTransaction creates a new stock transaction
func (r *stockTransactionRepo) CreateStockTransaction(transaction *models.StockTransaction) (*models.StockTransaction, error) {
	// TODO: Implement this method
	return nil, errors.New("method not implemented")
}

// ListStockTransactions retrieves a list of stock transactions based on filter
func (r *stockTransactionRepo) ListStockTransactions(filter models.StockTransactionFilter) ([]*models.StockTransaction, error) {
	// TODO: Implement this method
	return nil, errors.New("method not implemented")
}