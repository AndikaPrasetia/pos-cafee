package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/db"
	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/google/uuid"
)

// stockTransactionRepo implements the StockTransactionRepo interface
type stockTransactionRepo struct {
	queries *db.Queries
}

// CreateStockTransaction creates a new stock transaction
func (r *stockTransactionRepo) CreateStockTransaction(transaction *models.StockTransaction) (*models.StockTransaction, error) {
	menuItemUUID, err := uuid.Parse(transaction.MenuItemID)
	if err != nil {
		return nil, err
	}

	var userID uuid.NullUUID
	if transaction.UserID != nil {
		parsedUUID, err := uuid.Parse(*transaction.UserID)
		if err != nil {
			return nil, err
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

	var refUUID uuid.NullUUID
	if transaction.ReferenceID != nil {
		parsedUUID, err := uuid.Parse(*transaction.ReferenceID)
		if err != nil {
			return nil, err
		}
		refUUID = uuid.NullUUID{
			UUID:  parsedUUID,
			Valid: true,
		}
	} else {
		refUUID = uuid.NullUUID{
			UUID:  uuid.UUID{},
			Valid: false,
		}
	}

	var refType sql.NullString
	if transaction.ReferenceType != nil {
		refType = sql.NullString{
			String: *transaction.ReferenceType,
			Valid:  true,
		}
	} else {
		refType = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	dbTransaction, err := r.queries.CreateStockTransaction(context.Background(), db.CreateStockTransactionParams{
		MenuItemID:      menuItemUUID,
		TransactionType: string(transaction.TransactionType),
		Quantity:        int32(transaction.Quantity),
		PreviousStock:   int32(transaction.PreviousStock),
		CurrentStock:    int32(transaction.CurrentStock),
		Reason:          transaction.Reason,
		ReferenceType:   refType,
		ReferenceID:     refUUID,
		UserID:          userID,
	})
	if err != nil {
		return nil, err
	}

	createdTransaction := &models.StockTransaction{
		ID:              dbTransaction.ID.String(),
		MenuItemID:      dbTransaction.MenuItemID.String(),
		TransactionType: types.TransactionType(dbTransaction.TransactionType),
		Quantity:        int(dbTransaction.Quantity),
		PreviousStock:   int(dbTransaction.PreviousStock),
		CurrentStock:    int(dbTransaction.CurrentStock),
		Reason:          dbTransaction.Reason,
		CreatedAt:       dbTransaction.CreatedAt,
	}

	if dbTransaction.ReferenceType.Valid {
		createdTransaction.ReferenceType = &dbTransaction.ReferenceType.String
	}

	if dbTransaction.ReferenceID.Valid {
		refID := dbTransaction.ReferenceID.UUID.String()
		createdTransaction.ReferenceID = &refID
	}

	if dbTransaction.UserID.Valid {
		userID := dbTransaction.UserID.UUID.String()
		createdTransaction.UserID = &userID
	}

	return createdTransaction, nil
}

// ListStockTransactions retrieves a list of stock transactions based on filter
func (r *stockTransactionRepo) ListStockTransactions(filter models.StockTransactionFilter) ([]*models.StockTransaction, error) {
	var menuItemID uuid.UUID
	var startDate time.Time
	var endDate time.Time

	if filter.MenuItemID != nil {
		parsedUUID, err := uuid.Parse(*filter.MenuItemID)
		if err != nil {
			return nil, err
		}
		menuItemID = parsedUUID
	} else {
		menuItemID = uuid.Nil // Use Nil for null case
	}

	if filter.StartDate != nil {
		startDate = *filter.StartDate
	} else {
		startDate = time.Time{}
	}

	if filter.EndDate != nil {
		endDate = *filter.EndDate
	} else {
		endDate = time.Time{}
	}

	dbTransactions, err := r.queries.ListStockTransactions(context.Background(), db.ListStockTransactionsParams{
		Column1: menuItemID,
		Column2: startDate,
		Column3: endDate,
		Limit:   int32(filter.Limit),
		Offset:  int32(filter.Offset),
	})
	if err != nil {
		return nil, err
	}

	var transactions []*models.StockTransaction
	for _, dbTransaction := range dbTransactions {
		transaction := &models.StockTransaction{
			ID:              dbTransaction.ID.String(),
			MenuItemID:      dbTransaction.MenuItemID.String(),
			TransactionType: types.TransactionType(dbTransaction.TransactionType),
			Quantity:        int(dbTransaction.Quantity),
			PreviousStock:   int(dbTransaction.PreviousStock),
			CurrentStock:    int(dbTransaction.CurrentStock),
			Reason:          dbTransaction.Reason,
			CreatedAt:       dbTransaction.CreatedAt,
		}

		if dbTransaction.MenuItemName.Valid {
			transaction.MenuItemName = dbTransaction.MenuItemName.String
		}

		if dbTransaction.ReferenceType.Valid {
			transaction.ReferenceType = &dbTransaction.ReferenceType.String
		}

		if dbTransaction.ReferenceID.Valid {
			refID := dbTransaction.ReferenceID.UUID.String()
			transaction.ReferenceID = &refID
		}

		if dbTransaction.UserID.Valid {
			userID := dbTransaction.UserID.UUID.String()
			transaction.UserID = &userID
		}

		if dbTransaction.UserName.Valid {
			transaction.UserName = &dbTransaction.UserName.String
		}

		transactions = append(transactions, transaction)
	}

	// Ensure we return an empty slice instead of nil if no transactions found
	if transactions == nil {
		transactions = []*models.StockTransaction{}
	}

	return transactions, nil
}