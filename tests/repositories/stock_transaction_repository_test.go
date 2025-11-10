package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/db"
	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestStockTransactionRepository(t *testing.T) {
	// Setup test database connection
	dbConn := setupTestDB(t)
	defer dbConn.Close()
	
	queries := db.New(dbConn)
	stockTransactionRepo := &stockTransactionRepo{queries: queries}

	// Create test data
	menuItemID := createTestMenuItem(t, dbConn)
	userID := createTestUser(t, dbConn)

	// Create test stock transaction
	testTransaction := &models.StockTransaction{
		ID:              "test-transaction-id",
		MenuItemID:      menuItemID,
		TransactionType: types.TransactionTypeIn,
		Quantity:        10,
		PreviousStock:   50,
		CurrentStock:    60,
		Reason:          "Initial stock",
		UserID:          &userID,
		CreatedAt:       time.Now(),
	}

	t.Run("CreateStockTransaction", func(t *testing.T) {
		createdTransaction, err := stockTransactionRepo.CreateStockTransaction(testTransaction)
		assert.NoError(t, err)
		assert.Equal(t, testTransaction.MenuItemID, createdTransaction.MenuItemID)
		assert.Equal(t, testTransaction.TransactionType, createdTransaction.TransactionType)
		assert.Equal(t, testTransaction.Quantity, createdTransaction.Quantity)
	})

	t.Run("ListStockTransactions", func(t *testing.T) {
		menuItemID := createTestMenuItem(t, dbConn)
		userID := createTestUser(t, dbConn)
		
		// Create multiple stock transactions
		transaction1 := &models.StockTransaction{
			ID:              "transaction-1",
			MenuItemID:      menuItemID,
			TransactionType: types.TransactionTypeIn,
			Quantity:        5,
			PreviousStock:   0,
			CurrentStock:    5,
			Reason:          "Initial stock",
			UserID:          &userID,
			CreatedAt:       time.Now(),
		}
		
		transaction2 := &models.StockTransaction{
			ID:              "transaction-2",
			MenuItemID:      menuItemID,
			TransactionType: types.TransactionTypeOut,
			Quantity:        -2,
			PreviousStock:   5,
			CurrentStock:    3,
			Reason:          "Sale",
			UserID:          &userID,
			CreatedAt:       time.Now().Add(time.Hour), // Later time
		}
		
		_, err := stockTransactionRepo.CreateStockTransaction(transaction1)
		assert.NoError(t, err)
		
		_, err = stockTransactionRepo.CreateStockTransaction(transaction2)
		assert.NoError(t, err)

		filter := models.StockTransactionFilter{
			Limit:  10,
			Offset: 0,
		}

		transactions, err := stockTransactionRepo.ListStockTransactions(filter)
		assert.NoError(t, err)
		// At least the 2 transactions we just created should be returned
		assert.GreaterOrEqual(t, len(transactions), 2)
	})
}