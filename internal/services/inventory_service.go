package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/internal/repositories"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/google/uuid"
)

// InventoryService handles inventory-related business logic
type InventoryService struct {
	inventoryRepo       repositories.InventoryRepo
	stockTransactionRepo repositories.StockTransactionRepo
	menuRepo            repositories.MenuRepo
}

// NewInventoryService creates a new inventory service
func NewInventoryService(
	inventoryRepo repositories.InventoryRepo,
	stockTransactionRepo repositories.StockTransactionRepo,
	menuRepo repositories.MenuRepo,
) *InventoryService {
	return &InventoryService{
		inventoryRepo:       inventoryRepo,
		stockTransactionRepo: stockTransactionRepo,
		menuRepo:            menuRepo,
	}
}

// GetInventoryByMenuItem retrieves inventory information for a specific menu item
func (s *InventoryService) GetInventoryByMenuItem(menuItemID string) (*types.APIResponse, error) {
	// Validate menu item ID
	_, err := uuid.Parse(menuItemID)
	if err != nil {
		return nil, errors.New("invalid menu item ID")
	}

	inventory, err := s.inventoryRepo.GetInventoryByMenuItem(menuItemID)
	if err != nil {
		return nil, fmt.Errorf("inventory not found for item %s: %v", menuItemID, err)
	}

	return &types.APIResponse{
		Success: true,
		Data:    inventory,
	}, nil
}

// ListInventory retrieves a list of inventory items with optional filtering
func (s *InventoryService) ListInventory(filter models.InventoryFilter) (*types.APIResponse, error) {
	inventories, err := s.inventoryRepo.ListInventory(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data:    inventories,
	}, nil
}

// UpdateStock manually adjusts inventory stock levels
func (s *InventoryService) UpdateStock(userID string, updateData *models.InventoryUpdate) (*types.APIResponse, error) {
	// Validate user ID
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Validate menu item ID
	_, err = uuid.Parse(updateData.MenuItemID)
	if err != nil {
		return nil, errors.New("invalid menu item ID")
	}

	// Get current inventory
	currentInventory, err := s.inventoryRepo.GetInventoryByMenuItem(updateData.MenuItemID)
	if err != nil {
		return nil, fmt.Errorf("inventory not found for item %s: %v", updateData.MenuItemID, err)
	}

	// Calculate new stock level
	newStock := currentInventory.CurrentStock + updateData.Quantity

	// Ensure stock doesn't go negative for non-manager users (in a real implementation, this would be handled differently)
	// For now, we'll allow negative stock as the data model allows it
	if newStock < 0 {
		// This might require special authorization in a real system
	}

	// Update the inventory
	err = s.inventoryRepo.UpdateInventoryStock(updateData.MenuItemID, newStock, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update inventory stock: %v", err)
	}

	// Create a stock transaction record
	stockTransaction := &models.StockTransaction{
		ID:              uuid.New().String(),
		MenuItemID:      updateData.MenuItemID,
		TransactionType: getTransactionType(updateData.Quantity),
		Quantity:        updateData.Quantity,
		PreviousStock:   currentInventory.CurrentStock,
		CurrentStock:    newStock,
		Reason:          updateData.Reason,
		UserID:          &userID,
		CreatedAt:       time.Now(),
	}

	// In a complete implementation, we would add the transaction to the database
	_, err = s.stockTransactionRepo.CreateStockTransaction(stockTransaction)
	if err != nil {
		apiResponse := types.APIResponseWithError("Failed to create stock transaction: " + err.Error())
		return &apiResponse, err
	}

	// For now, we'll just return the updated inventory info

	updatedInventory, err := s.inventoryRepo.GetInventoryByMenuItem(updateData.MenuItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated inventory: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data:    updatedInventory,
	}, nil
}

// getTransactionType returns the appropriate transaction type based on the quantity change
func getTransactionType(quantity int) types.TransactionType {
	if quantity > 0 {
		return types.TransactionTypeIn
	} else if quantity < 0 {
		return types.TransactionTypeOut
	}
	return types.TransactionTypeAdjustment
}

// ListStockTransactions retrieves a list of stock transactions with optional filtering
func (s *InventoryService) ListStockTransactions(filter models.StockTransactionFilter) (*types.APIResponse, error) {
	transactions, err := s.stockTransactionRepo.ListStockTransactions(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock transactions: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data:    transactions,
	}, nil
}









// ValidateInventoryForOrder checks if there is sufficient inventory for an order
func (s *InventoryService) ValidateInventoryForOrder(items []models.OrderItemCreate) error {
	for _, item := range items {
		// Get current inventory for the menu item
		inventory, err := s.inventoryRepo.GetInventoryByMenuItem(item.MenuItemID)
		if err != nil {
			return fmt.Errorf("inventory not found for item %s: %v", item.MenuItemID, err)
		}

		// Check if enough stock is available
		if inventory.CurrentStock < item.Quantity {
			menuItem, err := s.menuRepo.GetMenuItem(item.MenuItemID)
			if err != nil {
				return fmt.Errorf("menu item %s not found", item.MenuItemID)
			}
			return fmt.Errorf("insufficient inventory for item %s: required %d, available %d", 
				menuItem.Name, item.Quantity, inventory.CurrentStock)
		}
	}

	return nil
}

// UpdateInventoryAfterOrder updates inventory after an order is completed
func (s *InventoryService) UpdateInventoryAfterOrder(items []models.OrderItemCreate, userID string) error {
	for _, item := range items {
		// Get current inventory for the menu item
		inventory, err := s.inventoryRepo.GetInventoryByMenuItem(item.MenuItemID)
		if err != nil {
			return fmt.Errorf("inventory not found for item %s: %v", item.MenuItemID, err)
		}

		// Calculate new stock
		newStock := inventory.CurrentStock - item.Quantity

		// Update the inventory
		err = s.inventoryRepo.UpdateInventoryStock(item.MenuItemID, newStock, userID)
		if err != nil {
			return fmt.Errorf("failed to update inventory stock for item %s: %v", item.MenuItemID, err)
		}

		// Create a stock transaction record
		stockTransaction := &models.StockTransaction{
			ID:              uuid.New().String(),
			MenuItemID:      item.MenuItemID,
			TransactionType: types.TransactionTypeOut,
			Quantity:        -item.Quantity, // Negative because it's going out
			PreviousStock:   inventory.CurrentStock,
			CurrentStock:    newStock,
			Reason:          "Order fulfillment",
			UserID:          &userID,
			CreatedAt:       time.Now(),
		}

		// In a complete implementation, we would store this transaction
		_, err = s.stockTransactionRepo.CreateStockTransaction(stockTransaction)
		if err != nil {
			// Log the error but don't fail the entire operation
			fmt.Printf("Failed to create stock transaction: %v\n", err)
		}
		// For now, we just log it conceptually
	}

	return nil
}