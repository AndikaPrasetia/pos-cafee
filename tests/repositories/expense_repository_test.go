package repositories

import (
	"testing"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/db"
	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestExpenseRepository(t *testing.T) {
	// Setup test database connection
	dbConn := setupTestDB(t)
	defer dbConn.Close()

	queries := db.New(dbConn)
	expenseRepo := &expenseRepo{queries: queries}

	// Create test user
	userID := createTestUser(t, dbConn)

	// Create test expense
	testExpense := &models.Expense{
		ID:          "test-expense-id", // This would normally be generated
		Category:    "Utilities",
		Description: "Electricity bill",
		Amount:      types.DecimalText(decimal.NewFromFloat(150.75)),
		Date:        time.Now(),
		UserID:      &userID,
		CreatedAt:   time.Now(),
	}

	t.Run("CreateExpense", func(t *testing.T) {
		createdExpense, err := expenseRepo.CreateExpense(testExpense)
		assert.NoError(t, err)
		assert.Equal(t, testExpense.Category, createdExpense.Category)
		assert.Equal(t, testExpense.Description, createdExpense.Description)
		assert.Equal(t, testExpense.Amount, createdExpense.Amount)
	})

	t.Run("GetExpense", func(t *testing.T) {
		createdExpense, err := expenseRepo.CreateExpense(testExpense)
		assert.NoError(t, err)

		retrievedExpense, err := expenseRepo.GetExpense(createdExpense.ID)
		assert.NoError(t, err)
		assert.Equal(t, createdExpense.ID, retrievedExpense.ID)
		assert.Equal(t, createdExpense.Category, retrievedExpense.Category)
	})

	t.Run("ListExpenses", func(t *testing.T) {
		// Create multiple expenses
		expense1 := &models.Expense{
			ID:          "expense-1",
			Category:    "Food",
			Description: "Groceries",
			Amount:      types.DecimalText(decimal.NewFromFloat(50.00)),
			Date:        time.Now(),
			UserID:      &userID,
			CreatedAt:   time.Now(),
		}

		expense2 := &models.Expense{
			ID:          "expense-2",
			Category:    "Transport",
			Description: "Gas",
			Amount:      types.DecimalText(decimal.NewFromFloat(40.00)),
			Date:        time.Now().AddDate(0, 0, -1), // Yesterday
			UserID:      &userID,
			CreatedAt:   time.Now(),
		}

		_, err := expenseRepo.CreateExpense(expense1)
		assert.NoError(t, err)

		_, err = expenseRepo.CreateExpense(expense2)
		assert.NoError(t, err)

		filter := models.ExpenseFilter{
			Limit:  10,
			Offset: 0,
		}

		expenses, err := expenseRepo.ListExpenses(filter)
		assert.NoError(t, err)
		// At least the 2 expenses we just created should be returned
		assert.GreaterOrEqual(t, len(expenses), 2)
	})

	t.Run("UpdateExpense", func(t *testing.T) {
		createdExpense, err := expenseRepo.CreateExpense(testExpense)
		assert.NoError(t, err)

		// Update the expense
		updatedCategory := "Office Supplies"
		updatedDescription := "New description"
		updatedAmount := types.DecimalText(decimal.NewFromFloat(75.25))

		createdExpense.Category = updatedCategory
		createdExpense.Description = updatedDescription
		createdExpense.Amount = updatedAmount

		updatedExpense, err := expenseRepo.UpdateExpense(createdExpense)
		assert.NoError(t, err)
		assert.Equal(t, updatedCategory, updatedExpense.Category)
		assert.Equal(t, updatedDescription, updatedExpense.Description)
		assert.Equal(t, updatedAmount, updatedExpense.Amount)
	})

	t.Run("DeleteExpense", func(t *testing.T) {
		createdExpense, err := expenseRepo.CreateExpense(testExpense)
		assert.NoError(t, err)

		err = expenseRepo.DeleteExpense(createdExpense.ID)
		assert.NoError(t, err)

		// Try to get the deleted expense - should return error
		_, err = expenseRepo.GetExpense(createdExpense.ID)
		assert.Error(t, err)
	})
}
