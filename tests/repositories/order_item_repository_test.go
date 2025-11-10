package repositories

import (
	"context"
	"testing"

	"github.com/AndikaPrasetia/pos-cafee/internal/db"
	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestOrderItemRepository(t *testing.T) {
	// Setup test database connection
	dbConn := setupTestDB(t)
	defer dbConn.Close()
	
	queries := db.New(dbConn)
	orderItemRepo := &orderItemRepo{queries: queries}

	// Create test data
	orderID := createTestOrder(t, dbConn)
	menuItemID := createTestMenuItem(t, dbConn)

	testOrderItem := &models.OrderItem{
		ID:         uuid.New().String(),
		OrderID:    orderID,
		MenuItemID: menuItemID,
		Quantity:   2,
		UnitPrice:  types.DecimalText(decimal.NewFromFloat(10.50)),
		TotalPrice: types.DecimalText(decimal.NewFromFloat(21.00)),
	}

	t.Run("CreateOrderItem", func(t *testing.T) {
		createdOrderItem, err := orderItemRepo.CreateOrderItem(testOrderItem)
		assert.NoError(t, err)
		assert.Equal(t, testOrderItem.OrderID, createdOrderItem.OrderID)
		assert.Equal(t, testOrderItem.MenuItemID, createdOrderItem.MenuItemID)
		assert.Equal(t, testOrderItem.Quantity, createdOrderItem.Quantity)
	})

	t.Run("GetOrderItem", func(t *testing.T) {
		createdOrderItem, err := orderItemRepo.CreateOrderItem(testOrderItem)
		assert.NoError(t, err)

		retrievedOrderItem, err := orderItemRepo.GetOrderItem(createdOrderItem.ID)
		assert.NoError(t, err)
		assert.Equal(t, createdOrderItem.ID, retrievedOrderItem.ID)
		assert.Equal(t, createdOrderItem.OrderID, retrievedOrderItem.OrderID)
		assert.Equal(t, createdOrderItem.MenuItemID, retrievedOrderItem.MenuItemID)
	})

	t.Run("GetOrderItemsByOrderID", func(t *testing.T) {
		orderID := createTestOrder(t, dbConn)
		
		// Create multiple order items for the same order
		orderItem1 := &models.OrderItem{
			ID:         uuid.New().String(),
			OrderID:    orderID,
			MenuItemID: testOrderItem.MenuItemID,
			Quantity:   1,
			UnitPrice:  types.DecimalText(decimal.NewFromFloat(15.00)),
			TotalPrice: types.DecimalText(decimal.NewFromFloat(15.00)),
		}
		
		orderItem2 := &models.OrderItem{
			ID:         uuid.New().String(),
			OrderID:    orderID,
			MenuItemID: testOrderItem.MenuItemID,
			Quantity:   3,
			UnitPrice:  types.DecimalText(decimal.NewFromFloat(5.00)),
			TotalPrice: types.DecimalText(decimal.NewFromFloat(15.00)),
		}
		
		_, err := orderItemRepo.CreateOrderItem(orderItem1)
		assert.NoError(t, err)
		
		_, err = orderItemRepo.CreateOrderItem(orderItem2)
		assert.NoError(t, err)

		orderItems, err := orderItemRepo.GetOrderItemsByOrderID(orderID)
		assert.NoError(t, err)
		assert.Len(t, orderItems, 2)
	})

	t.Run("GetOrderItemsWithDetails", func(t *testing.T) {
		orderID := createTestOrder(t, dbConn)
		
		// Create an order item
		orderItem := &models.OrderItem{
			ID:         uuid.New().String(),
			OrderID:    orderID,
			MenuItemID: testOrderItem.MenuItemID,
			Quantity:   1,
			UnitPrice:  types.DecimalText(decimal.NewFromFloat(12.50)),
			TotalPrice: types.DecimalText(decimal.NewFromFloat(12.50)),
		}
		
		_, err := orderItemRepo.CreateOrderItem(orderItem)
		assert.NoError(t, err)

		orderItemDetails, err := orderItemRepo.GetOrderItemsWithDetails(orderID)
		assert.NoError(t, err)
		assert.Len(t, orderItemDetails, 1)
		assert.Equal(t, orderItem.OrderID, orderItemDetails[0].OrderID)
	})

	t.Run("UpdateOrderItem", func(t *testing.T) {
		createdOrderItem, err := orderItemRepo.CreateOrderItem(testOrderItem)
		assert.NoError(t, err)

		// Update the order item
		updatedQuantity := 5
		updatedUnitPrice := types.DecimalText(decimal.NewFromFloat(20.00))
		updatedTotalPrice := types.DecimalText(decimal.NewFromFloat(100.00))
		
		createdOrderItem.Quantity = updatedQuantity
		createdOrderItem.UnitPrice = updatedUnitPrice
		createdOrderItem.TotalPrice = updatedTotalPrice

		updatedOrderItem, err := orderItemRepo.UpdateOrderItem(createdOrderItem)
		assert.NoError(t, err)
		assert.Equal(t, updatedQuantity, updatedOrderItem.Quantity)
		assert.Equal(t, updatedUnitPrice, updatedOrderItem.UnitPrice)
		assert.Equal(t, updatedTotalPrice, updatedOrderItem.TotalPrice)
	})

	t.Run("DeleteOrderItem", func(t *testing.T) {
		createdOrderItem, err := orderItemRepo.CreateOrderItem(testOrderItem)
		assert.NoError(t, err)

		err = orderItemRepo.DeleteOrderItem(createdOrderItem.ID)
		assert.NoError(t, err)

		// Try to get the deleted order item - should return error
		_, err = orderItemRepo.GetOrderItem(createdOrderItem.ID)
		assert.Error(t, err)
	})
}

// Helper function to create a test order
func createTestOrder(t *testing.T, dbConn *db.DB) string {
	queries := db.New(dbConn)
	
	userID := createTestUser(t, dbConn)
	
	order, err := queries.CreateOrder(context.Background(), db.CreateOrderParams{
		OrderNumber:    "ORD-TEST-0001",
		UserID:         uuid.Must(uuid.Parse(userID)),
		TotalAmount:    "100.00",
		DiscountAmount: "0.00",
		TaxAmount:      "10.00",
	})
	assert.NoError(t, err)
	
	return order.ID.String()
}

// Helper function to create a test menu item
func createTestMenuItem(t *testing.T, dbConn *db.DB) string {
	queries := db.New(dbConn)
	
	categoryID := createTestCategory(t, dbConn)
	
	menuItem, err := queries.CreateMenuItem(context.Background(), db.CreateMenuItemParams{
		Name:        "Test Item",
		CategoryID:  uuid.Must(uuid.Parse(categoryID)),
		Description: "Test description",
		Price:       "15.00",
		Cost:        "10.00",
		IsAvailable: true,
	})
	assert.NoError(t, err)
	
	return menuItem.ID.String()
}

// Helper function to create a test user
func createTestUser(t *testing.T, dbConn *db.DB) string {
	queries := db.New(dbConn)
	
	user, err := queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		Role:      "cashier",
		IsActive:  true,
	})
	assert.NoError(t, err)
	
	return user.ID.String()
}

// Helper function to create a test category
func createTestCategory(t *testing.T, dbConn *db.DB) string {
	queries := db.New(dbConn)
	
	category, err := queries.CreateCategory(context.Background(), db.CreateCategoryParams{
		Name:        "Test Category",
		Description: "Test category description",
		IsActive:    true,
	})
	assert.NoError(t, err)
	
	return category.ID.String()
}

// Helper function to set up test database
func setupTestDB(t *testing.T) *db.DB {
	// This is a placeholder - in a real implementation, 
	// you'd set up a test database connection
	// For now, we'll assume the main database connection works for tests too
	dbConn := connectToTestDB()
	return dbConn
}

// Placeholder function for connecting to test DB
func connectToTestDB() *db.DB {
	// This would connect to a test DB in a real implementation
	// For now we're just providing signature
	return nil
}