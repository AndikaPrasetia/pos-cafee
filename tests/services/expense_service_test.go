package services

import (
	"testing"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExpenseRepo is a mock implementation of the ExpenseRepo interface
type MockExpenseRepo struct {
	mock.Mock
}

func (m *MockExpenseRepo) CreateExpense(expense *models.Expense) (*models.Expense, error) {
	args := m.Called(expense)
	return args.Get(0).(*models.Expense), args.Error(1)
}

func (m *MockExpenseRepo) GetExpense(id string) (*models.Expense, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Expense), args.Error(1)
}

func (m *MockExpenseRepo) ListExpenses(filter models.ExpenseFilter) ([]*models.Expense, error) {
	args := m.Called(filter)
	return args.Get(0).([]*models.Expense), args.Error(1)
}

func (m *MockExpenseRepo) UpdateExpense(expense *models.Expense) (*models.Expense, error) {
	args := m.Called(expense)
	return args.Get(0).(*models.Expense), args.Error(1)
}

func (m *MockExpenseRepo) DeleteExpense(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestExpenseService(t *testing.T) {
	mockRepo := new(MockExpenseRepo)
	expenseService := NewExpenseService(mockRepo)

	userID := "test-user-id"

	t.Run("CreateExpense", func(t *testing.T) {
		expenseData := &models.ExpenseCreate{
			Category:    "Utilities",
			Description: "Electricity bill",
			Amount:      types.DecimalText(decimal.NewFromFloat(150.75)),
			Date:        time.Now(),
		}

		expectedExpense := &models.Expense{
			ID:          "test-expense-id",
			Category:    expenseData.Category,
			Description: expenseData.Description,
			Amount:      expenseData.Amount,
			Date:        expenseData.Date,
			UserID:      &userID,
		}

		mockRepo.On("CreateExpense", mock.AnythingOfType("*models.Expense")).Return(expectedExpense, nil)

		result, err := expenseService.CreateExpense(userID, expenseData)
		assert.NoError(t, err)
		assert.True(t, result.Success)
		assert.NotNil(t, result.Data)

		mockRepo.AssertExpectations(t)
	})

	t.Run("GetExpense", func(t *testing.T) {
		expenseID := "test-expense-id"
		expectedExpense := &models.Expense{
			ID:       expenseID,
			Category: "Utilities",
			Amount:   types.DecimalText(decimal.NewFromFloat(150.75)),
			Date:     time.Now(),
		}

		mockRepo.On("GetExpense", expenseID).Return(expectedExpense, nil)

		result, err := expenseService.GetExpense(expenseID)
		assert.NoError(t, err)
		assert.True(t, result.Success)
		assert.NotNil(t, result.Data)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ListExpenses", func(t *testing.T) {
		expectedExpenses := []*models.Expense{
			{
				ID:       "expense-1",
				Category: "Utilities",
				Amount:   types.DecimalText(decimal.NewFromFloat(150.75)),
				Date:     time.Now(),
			},
			{
				ID:       "expense-2",
				Category: "Food",
				Amount:   types.DecimalText(decimal.NewFromFloat(50.25)),
				Date:     time.Now().AddDate(0, 0, -1),
			},
		}

		filter := models.ExpenseFilter{
			Limit:  10,
			Offset: 0,
		}

		mockRepo.On("ListExpenses", filter).Return(expectedExpenses, nil)

		result, err := expenseService.ListExpenses(filter)
		assert.NoError(t, err)
		assert.True(t, result.Success)
		assert.Len(t, result.Data, 2)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateExpense", func(t *testing.T) {
		expenseID := "test-expense-id"
		userID := "test-user-id"
		updateData := &models.ExpenseUpdate{
			Category:    stringPtr("Office Supplies"),
			Description: stringPtr("New office supplies"),
			Amount:      &types.DecimalText(decimal.NewFromFloat(75.50)),
			Date:        &time.Now(),
		}

		existingExpense := &models.Expense{
			ID:          expenseID,
			Category:    "Utilities",
			Description: "Old description",
			Amount:      types.DecimalText(decimal.NewFromFloat(150.75)),
			Date:        time.Now(),
		}

		updatedExpense := &models.Expense{
			ID:          expenseID,
			Category:    "Office Supplies",
			Description: "New office supplies",
			Amount:      types.DecimalText(decimal.NewFromFloat(75.50)),
			Date:        time.Now(),
		}

		mockRepo.On("GetExpense", expenseID).Return(existingExpense, nil)
		mockRepo.On("UpdateExpense", mock.AnythingOfType("*models.Expense")).Return(updatedExpense, nil)

		result, err := expenseService.UpdateExpense(expenseID, userID, updateData)
		assert.NoError(t, err)
		assert.True(t, result.Success)
		assert.NotNil(t, result.Data)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteExpense", func(t *testing.T) {
		expenseID := "test-expense-id"
		userID := "test-user-id"

		expectedExpense := &models.Expense{
			ID:       expenseID,
			Category: "Utilities",
			Amount:   types.DecimalText(decimal.NewFromFloat(150.75)),
			Date:     time.Now(),
		}

		mockRepo.On("GetExpense", expenseID).Return(expectedExpense, nil)
		mockRepo.On("DeleteExpense", expenseID).Return(nil)

		result, err := expenseService.DeleteExpense(expenseID, userID)
		assert.NoError(t, err)
		assert.True(t, result.Success)
		assert.Equal(t, "Expense deleted successfully", result.Message)

		mockRepo.AssertExpectations(t)
	})

	t.Run("GetExpenseSummary", func(t *testing.T) {
		startDate := time.Now().AddDate(0, -1, 0) // One month ago
		endDate := time.Now()

		expectedExpenses := []*models.Expense{
			{
				ID:       "expense-1",
				Category: "Utilities",
				Amount:   types.DecimalText(decimal.NewFromFloat(150.75)),
				Date:     time.Now().AddDate(0, 0, -15),
			},
			{
				ID:       "expense-2",
				Category: "Utilities",
				Amount:   types.DecimalText(decimal.NewFromFloat(49.25)),
				Date:     time.Now().AddDate(0, 0, -10),
			},
			{
				ID:       "expense-3",
				Category: "Food",
				Amount:   types.DecimalText(decimal.NewFromFloat(75.00)),
				Date:     time.Now().AddDate(0, 0, -5),
			},
		}

		mockRepo.On("ListExpenses", mock.MatchedBy(func(f models.ExpenseFilter) bool {
			return f.StartDate != nil && f.EndDate != nil
		})).Return(expectedExpenses, nil)

		result, err := expenseService.GetExpenseSummary(startDate, endDate)
		assert.NoError(t, err)
		assert.True(t, result.Success)
		assert.NotNil(t, result.Data)

		// Verify the summary calculation
		summary := result.Data.(*models.ExpenseSummary)
		expectedTotal := types.DecimalText(decimal.NewFromFloat(275.00)) // 150.75 + 49.25 + 75.00
		assert.Equal(t, expectedTotal, summary.TotalAmount)
		assert.Len(t, summary.CategoryCount, 2) // Utilities and Food categories

		mockRepo.AssertExpectations(t)
	})
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}