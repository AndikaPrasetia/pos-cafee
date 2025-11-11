package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/db"
	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// expenseRepo implements the ExpenseRepo interface
type expenseRepo struct {
	queries *db.Queries
}

// CreateExpense creates a new expense record
func (r *expenseRepo) CreateExpense(expense *models.Expense) (*models.Expense, error) {
	var userID uuid.NullUUID
	if expense.UserID != nil {
		parsedUUID, err := uuid.Parse(*expense.UserID)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID: %w", err)
		}
		userID = uuid.NullUUID{
			UUID:  parsedUUID,
			Valid: true,
		}
	} else {
		userID = uuid.NullUUID{
			UUID:  uuid.UUID{},
			Valid: false,
		}
	}

	dbExpense, err := r.queries.CreateExpense(context.Background(), db.CreateExpenseParams{
		Category:    expense.Category,
		Description: sql.NullString{String: expense.Description, Valid: expense.Description != ""},
		Amount:      expense.Amount.String(),
		Date:        expense.Date,
		UserID:      userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create expense in database: %w", err)
	}

	amount, err := decimal.NewFromString(dbExpense.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expense amount %s: %w", dbExpense.Amount, err)
	}

	createdExpense := &models.Expense{
		ID:          dbExpense.ID.String(),
		Category:    dbExpense.Category,
		Description: dbExpense.Description.String,
		Amount:      types.DecimalText(amount),
		Date:        dbExpense.Date,
		CreatedAt:   dbExpense.CreatedAt,
	}

	if dbExpense.UserID.Valid {
		userIDStr := dbExpense.UserID.UUID.String()
		createdExpense.UserID = &userIDStr
	}

	return createdExpense, nil
}

// GetExpense retrieves an expense by ID
func (r *expenseRepo) GetExpense(id string) (*models.Expense, error) {
	expenseID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid expense ID: %w", err)
	}

	dbExpense, err := r.queries.GetExpense(context.Background(), expenseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("expense not found")
		}
		return nil, fmt.Errorf("failed to fetch expense from database: %w", err)
	}

	amount, err := decimal.NewFromString(dbExpense.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expense amount %s: %w", dbExpense.Amount, err)
	}

	expense := &models.Expense{
		ID:          dbExpense.ID.String(),
		Category:    dbExpense.Category,
		Description: dbExpense.Description.String,
		Amount:      types.DecimalText(amount),
		Date:        dbExpense.Date,
		CreatedAt:   dbExpense.CreatedAt,
	}

	if dbExpense.UserID.Valid {
		userIDStr := dbExpense.UserID.UUID.String()
		expense.UserID = &userIDStr
	}

	return expense, nil
}

// ListExpenses retrieves a list of expenses based on filter
func (r *expenseRepo) ListExpenses(filter models.ExpenseFilter) ([]*models.Expense, error) {
	var startDate time.Time
	var endDate time.Time
	var category string

	if filter.StartDate != nil {
		startDate = *filter.StartDate
	} else {
		startDate = time.Time{} // Zero time will be handled by the SQL query
	}

	if filter.EndDate != nil {
		endDate = *filter.EndDate
	} else {
		endDate = time.Time{} // Zero time will be handled by the SQL query
	}

	if filter.Category != nil {
		category = *filter.Category
	} else {
		category = ""
	}

	dbExpenses, err := r.queries.ListExpenses(context.Background(), db.ListExpensesParams{
		Column1: startDate,
		Column2: endDate,
		Column3: category,
		Limit:   int32(filter.Limit),
		Offset:  int32(filter.Offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expenses from database: %w", err)
	}

	var expenses []*models.Expense
	for _, dbExpense := range dbExpenses {
		amount, err := decimal.NewFromString(dbExpense.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed to parse expense amount %s for expense %s: %w", dbExpense.Amount, dbExpense.ID.String(), err)
		}

		expense := &models.Expense{
			ID:          dbExpense.ID.String(),
			Category:    dbExpense.Category,
			Description: dbExpense.Description.String,
			Amount:      types.DecimalText(amount),
			Date:        dbExpense.Date,
			CreatedAt:   dbExpense.CreatedAt,
		}

		if dbExpense.UserID.Valid {
			userIDStr := dbExpense.UserID.UUID.String()
			expense.UserID = &userIDStr
		}

		expenses = append(expenses, expense)
	}

	return expenses, nil
}

// UpdateExpense updates an existing expense
func (r *expenseRepo) UpdateExpense(expense *models.Expense) (*models.Expense, error) {
	expenseID, err := uuid.Parse(expense.ID)
	if err != nil {
		return nil, err
	}

	dbExpense, err := r.queries.UpdateExpense(context.Background(), db.UpdateExpenseParams{
		ID:          expenseID,
		Category:    expense.Category,
		Description: sql.NullString{String: expense.Description, Valid: expense.Description != ""},
		Amount:      expense.Amount.String(),
		Date:        expense.Date,
	})
	if err != nil {
		return nil, err
	}

	updatedExpense := &models.Expense{
		ID:          dbExpense.ID.String(),
		Category:    dbExpense.Category,
		Description: dbExpense.Description.String,
		Amount:      types.DecimalText(decimal.RequireFromString(dbExpense.Amount)),
		Date:        dbExpense.Date,
		CreatedAt:   dbExpense.CreatedAt,
	}

	if dbExpense.UserID.Valid {
		userIDStr := dbExpense.UserID.UUID.String()
		updatedExpense.UserID = &userIDStr
	}

	return updatedExpense, nil
}

// DeleteExpense deletes an expense by ID
func (r *expenseRepo) DeleteExpense(id string) error {
	expenseID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = r.queries.DeleteExpense(context.Background(), expenseID)
	if err != nil {
		return err
	}

	return nil
}
// GetExpensesByDateRange retrieves expenses within a date range
func (r *expenseRepo) GetExpensesByDateRange(startDate, endDate time.Time) ([]*models.Expense, error) {
	dbExpenses, err := r.queries.ListExpenses(context.Background(), db.ListExpensesParams{
		Column1: startDate,
		Column2: endDate,
		Column3: "",  // Empty string means no category filter
		Limit:   10000,  // Reasonable limit for expense reports
		Offset:  0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expenses from database: %w", err)
	}

	var expenses []*models.Expense
	for _, dbExpense := range dbExpenses {
		amount, err := decimal.NewFromString(dbExpense.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed to parse expense amount %s for expense %s: %w", dbExpense.Amount, dbExpense.ID.String(), err)
		}

		expense := &models.Expense{
			ID:          dbExpense.ID.String(),
			Category:    dbExpense.Category,
			Description: dbExpense.Description.String,
			Amount:      types.DecimalText(amount),
			Date:        dbExpense.Date,
			CreatedAt:   dbExpense.CreatedAt,
		}

		if dbExpense.UserID.Valid {
			userIDStr := dbExpense.UserID.UUID.String()
			expense.UserID = &userIDStr
		}

		expenses = append(expenses, expense)
	}

	return expenses, nil
}

