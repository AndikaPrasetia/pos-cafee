package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/internal/repositories"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ExpenseService handles expense-related business logic
type ExpenseService struct {
	expenseRepo repositories.ExpenseRepo
}

// NewExpenseService creates a new expense service
func NewExpenseService(
	expenseRepo repositories.ExpenseRepo,
) *ExpenseService {
	return &ExpenseService{
		expenseRepo: expenseRepo,
	}
}

// CreateExpense creates a new expense record
func (s *ExpenseService) CreateExpense(userID string, expenseData *models.ExpenseCreate) (*types.APIResponse, error) {
	// Validate user ID
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Validate expense data
	if expenseData.Category == "" {
		return nil, errors.New("category is required")
	}

	if expenseData.Amount.Cmp(types.DecimalText(decimal.Zero)) <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	if expenseData.Date.IsZero() {
		return nil, errors.New("date is required")
	}

	// Create expense model
	expense := &models.Expense{
		ID:          uuid.New().String(),
		Category:    expenseData.Category,
		Description: expenseData.Description,
		Amount:      expenseData.Amount,
		Date:        expenseData.Date,
		UserID:      &userID,
		CreatedAt:   time.Now(),
	}

	// Save to database
	createdExpense, err := s.expenseRepo.CreateExpense(expense)
	if err != nil {
		return nil, fmt.Errorf("failed to create expense: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data:    createdExpense,
	}, nil
}

// GetExpense retrieves an expense by ID
func (s *ExpenseService) GetExpense(id string) (*types.APIResponse, error) {
	// Validate expense ID
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid expense ID")
	}

	expense, err := s.expenseRepo.GetExpense(id)
	if err != nil {
		return nil, fmt.Errorf("expense not found: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data:    expense,
	}, nil
}

// ListExpenses retrieves a list of expenses based on filter criteria
func (s *ExpenseService) ListExpenses(filter models.ExpenseFilter) (*types.APIResponse, error) {
	expenses, err := s.expenseRepo.ListExpenses(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list expenses: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data:    expenses,
	}, nil
}

// UpdateExpense updates an existing expense record
func (s *ExpenseService) UpdateExpense(id string, userID string, expenseData *models.ExpenseUpdate) (*types.APIResponse, error) {
	// Validate expense ID
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid expense ID")
	}

	// Validate user ID
	_, err = uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get existing expense
	existingExpense, err := s.expenseRepo.GetExpense(id)
	if err != nil {
		return nil, fmt.Errorf("expense not found: %v", err)
	}

	// Update fields if provided
	if expenseData.Category != nil {
		existingExpense.Category = *expenseData.Category
	}
	if expenseData.Description != nil {
		existingExpense.Description = *expenseData.Description
	}
	if expenseData.Amount != nil && expenseData.Amount.Cmp(types.DecimalText(decimal.Zero)) > 0 {
		existingExpense.Amount = *expenseData.Amount
	}
	if expenseData.Date != nil {
		existingExpense.Date = *expenseData.Date
	}

	// Save updated expense
	updatedExpense, err := s.expenseRepo.UpdateExpense(existingExpense)
	if err != nil {
		return nil, fmt.Errorf("failed to update expense: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data:    updatedExpense,
	}, nil
}

// DeleteExpense deletes an expense by ID
func (s *ExpenseService) DeleteExpense(id string, userID string) (*types.APIResponse, error) {
	// Validate expense ID
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid expense ID")
	}

	// Validate user ID
	_, err = uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Check if expense exists
	_, err = s.expenseRepo.GetExpense(id)
	if err != nil {
		return nil, fmt.Errorf("expense not found: %v", err)
	}

	// Delete the expense
	err = s.expenseRepo.DeleteExpense(id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete expense: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Message: "Expense deleted successfully",
	}, nil
}

// GetExpenseSummary generates a summary of expenses for reporting
func (s *ExpenseService) GetExpenseSummary(startDate, endDate time.Time) (*types.APIResponse, error) {
	filter := models.ExpenseFilter{
		StartDate: &startDate,
		EndDate:   &endDate,
		Limit:     10000, // Reasonable limit for expenses
		Offset:    0,
	}

	expenses, err := s.expenseRepo.ListExpenses(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expenses: %v", err)
	}

	// Calculate summary data
	var totalAmount types.DecimalText = types.DecimalText(decimal.Zero)
	categoryCount := make(map[string]types.DecimalText)

	for _, expense := range expenses {
		totalAmount = totalAmount.Add(expense.Amount)
		
		// Add to category total
		if current, exists := categoryCount[expense.Category]; exists {
			categoryCount[expense.Category] = current.Add(expense.Amount)
		} else {
			categoryCount[expense.Category] = expense.Amount
		}
	}

	summary := &models.ExpenseSummary{
		TotalAmount:   totalAmount,
		CategoryCount: categoryCount,
	}

	return &types.APIResponse{
		Success: true,
		Data:    summary,
	}, nil
}