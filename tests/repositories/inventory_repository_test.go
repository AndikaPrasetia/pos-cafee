package repositories

import (
	"testing"

	"github.com/AndikaPrasetia/pos-cafee/internal/db"
	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestInventoryRepository(t *testing.T) {
	// Setup test database connection
	dbConn := setupTestDB(t)
	defer dbConn.Close()
	
	queries := db.New(dbConn)
	inventoryRepo := &inventoryRepo{queries: queries}

	// Create test data
	menuItemID := createTestMenuItem(t, dbConn)
	userID := createTestUser(t, dbConn)

	t.Run("CreateInventoryRecord", func(t *testing.T) {
		err := inventoryRepo.CreateInventoryRecord(menuItemID)
		assert.NoError(t, err)
	})

	t.Run("GetInventoryByMenuItem", func(t *testing.T) {
		// First create an inventory record
		err := inventoryRepo.CreateInventoryRecord(menuItemID)
		assert.NoError(t, err)

		inventory, err := inventoryRepo.GetInventoryByMenuItem(menuItemID)
		assert.NoError(t, err)
		assert.Equal(t, menuItemID, inventory.MenuItemID)
		// Default stock should be 0
		assert.Equal(t, 0, inventory.CurrentStock)
	})

	t.Run("UpdateInventoryStock", func(t *testing.T) {
		// First create an inventory record
		err := inventoryRepo.CreateInventoryRecord(menuItemID)
		assert.NoError(t, err)

		// Update the inventory stock
		newStock := 50
		err = inventoryRepo.UpdateInventoryStock(menuItemID, newStock, userID)
		assert.NoError(t, err)

		// Verify the update
		updatedInventory, err := inventoryRepo.GetInventoryByMenuItem(menuItemID)
		assert.NoError(t, err)
		assert.Equal(t, newStock, updatedInventory.CurrentStock)
	})

	t.Run("ListInventory", func(t *testing.T) {
		// Create multiple inventory records
		menuItemID1 := createTestMenuItem(t, dbConn)
		menuItemID2 := createTestMenuItem(t, dbConn)

		err := inventoryRepo.CreateInventoryRecord(menuItemID1)
		assert.NoError(t, err)

		err = inventoryRepo.CreateInventoryRecord(menuItemID2)
		assert.NoError(t, err)

		// Update stock for the items
		err = inventoryRepo.UpdateInventoryStock(menuItemID1, 25, userID)
		assert.NoError(t, err)

		err = inventoryRepo.UpdateInventoryStock(menuItemID2, 15, userID)
		assert.NoError(t, err)

		filter := models.InventoryFilter{
			Limit: 10,
			Offset: 0,
		}

		inventories, err := inventoryRepo.ListInventory(filter)
		assert.NoError(t, err)

		// Find our test items in the results
		var foundItem1, foundItem2 bool
		for _, inv := range inventories {
			if inv.MenuItemID == menuItemID1 && inv.CurrentStock == 25 {
				foundItem1 = true
			}
			if inv.MenuItemID == menuItemID2 && inv.CurrentStock == 15 {
				foundItem2 = true
			}
		}
		assert.True(t, foundItem1, "First test inventory item not found or has wrong stock")
		assert.True(t, foundItem2, "Second test inventory item not found or has wrong stock")
	})
}