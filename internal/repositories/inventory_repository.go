package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AndikaPrasetia/pos-cafee/internal/db"
	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/google/uuid"
)

// inventoryRepo implements the InventoryRepo interface
type inventoryRepo struct {
	queries *db.Queries
}

// GetInventoryByMenuItem retrieves inventory by menu item ID
func (r *inventoryRepo) GetInventoryByMenuItem(menuItemID string) (*models.Inventory, error) {
	menuItemUUID, err := uuid.Parse(menuItemID)
	if err != nil {
		return nil, err
	}

	dbInventory, err := r.queries.GetInventoryByMenuItem(context.Background(), menuItemUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("inventory not found")
		}
		return nil, err
	}

	inventory := &models.Inventory{
		ID:            dbInventory.ID.String(),
		MenuItemID:    dbInventory.MenuItemID.String(),
		CurrentStock:  int(dbInventory.CurrentStock),
		MinimumStock:  int(dbInventory.MinimumStock),
		Unit:          dbInventory.Unit,
		LastUpdatedAt: dbInventory.LastUpdatedAt,
	}

	if dbInventory.LastUpdatedBy.Valid {
		userIDStr := dbInventory.LastUpdatedBy.UUID.String()
		inventory.LastUpdatedBy = &userIDStr
	}

	return inventory, nil
}

// ListInventory retrieves a list of inventory based on filter
func (r *inventoryRepo) ListInventory(filter models.InventoryFilter) ([]*models.Inventory, error) {
	var lowStockFilter bool
	if filter.LowStockOnly {
		lowStockFilter = true
	} else {
		lowStockFilter = false
	}

	dbInventories, err := r.queries.ListInventory(context.Background(), db.ListInventoryParams{
		Column1: lowStockFilter,
		Limit:   int32(filter.Limit),
		Offset:  int32(filter.Offset),
	})
	if err != nil {
		return nil, err
	}

	var inventories []*models.Inventory
	for _, dbInventory := range dbInventories {
		inventory := &models.Inventory{
			ID:             dbInventory.ID.String(),
			MenuItemID:     dbInventory.MenuItemID.String(),
			MenuItemName:   dbInventory.MenuItemName,
			CurrentStock:   int(dbInventory.CurrentStock),
			MinimumStock:   int(dbInventory.MinimumStock),
			Unit:           dbInventory.Unit,
			LastUpdatedAt:  dbInventory.LastUpdatedAt,
			IsLowStock:     dbInventory.CurrentStock <= dbInventory.MinimumStock,
		}

		if dbInventory.LastUpdatedBy.Valid {
			inventory.LastUpdatedByName = &dbInventory.LastUpdatedBy.String
		}
		inventories = append(inventories, inventory)
	}

	return inventories, nil
}

// UpdateInventoryStock updates the stock for a menu item
func (r *inventoryRepo) UpdateInventoryStock(menuItemID string, stock int, userID string) error {
	menuItemUUID, err := uuid.Parse(menuItemID)
	if err != nil {
		return err
	}

	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	userUUID := uuid.NullUUID{
		UUID:  parsedUserUUID,
		Valid: true,
	}

	err = r.queries.UpdateInventoryStock(context.Background(), db.UpdateInventoryStockParams{
		MenuItemID:    menuItemUUID,
		CurrentStock:  int32(stock),
		LastUpdatedBy: userUUID,
	})
	if err != nil {
		return err
	}

	return nil
}

// CreateInventoryRecord creates a new inventory record
func (r *inventoryRepo) CreateInventoryRecord(menuItemID string) error {
	menuItemUUID, err := uuid.Parse(menuItemID)
	if err != nil {
		return err
	}

	err = r.queries.CreateInventoryRecord(context.Background(), menuItemUUID)
	if err != nil {
		return err
	}

	return nil
}