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

// MockOrderItemRepo is a mock implementation of the OrderItemRepo interface
type MockOrderItemRepo struct {
	mock.Mock
}

func (m *MockOrderItemRepo) GetOrderItem(id string) (*models.OrderItem, error) {
	args := m.Called(id)
	return args.Get(0).(*models.OrderItem), args.Error(1)
}

func (m *MockOrderItemRepo) GetOrderItemsByOrderID(orderID string) ([]*models.OrderItem, error) {
	args := m.Called(orderID)
	return args.Get(0).([]*models.OrderItem), args.Error(1)
}

func (m *MockOrderItemRepo) CreateOrderItem(orderItem *models.OrderItem) (*models.OrderItem, error) {
	args := m.Called(orderItem)
	return args.Get(0).(*models.OrderItem), args.Error(1)
}

func (m *MockOrderItemRepo) UpdateOrderItem(orderItem *models.OrderItem) (*models.OrderItem, error) {
	args := m.Called(orderItem)
	return args.Get(0).(*models.OrderItem), args.Error(1)
}

func (m *MockOrderItemRepo) DeleteOrderItem(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockOrderItemRepo) GetOrderItemsWithDetails(orderID string) ([]*models.OrderItemWithDetails, error) {
	args := m.Called(orderID)
	return args.Get(0).([]*models.OrderItemWithDetails), args.Error(1)
}

func TestOrderServiceWithOrderItemRepo(t *testing.T) {
	// Create mock repositories
	mockOrderRepo := new(MockOrderRepo)
	mockOrderItemRepo := new(MockOrderItemRepo)
	mockMenuRepo := new(MockMenuRepo)
	mockInventoryRepo := new(MockInventoryRepo)
	mockStockTransactionRepo := new(MockStockTransactionRepo)

	// Create service with the mocked repositories
	orderService := NewOrderService(mockOrderRepo, mockOrderItemRepo, mockMenuRepo, mockInventoryRepo, mockStockTransactionRepo)

	userID := "test-user-id"
	orderID := "test-order-id"

	t.Run("CreateOrder with Order Items", func(t *testing.T) {
		// Set up test data
		menuItemID := "test-menu-item-id"
		orderCreateData := &models.OrderCreate{
			Items: []models.OrderItemCreate{
				{
					MenuItemID: menuItemID,
					Quantity:   2,
				},
			},
		}

		// Expected order after creation
		expectedOrder := &models.Order{
			ID:             orderID,
			OrderNumber:    "ORD-TEST-0001",
			UserID:         userID,
			Status:         types.OrderStatusDraft,
			TotalAmount:    types.DecimalText(decimal.NewFromFloat(30.0)),
			DiscountAmount: types.DecimalText(decimal.Zero),
			TaxAmount:      types.DecimalText(decimal.Zero),
			PaymentStatus:  types.PaymentStatusPending,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Expected menu item
		expectedMenuItem := &models.MenuItem{
			ID:          menuItemID,
			Name:        "Test Item",
			Price:       types.DecimalText(decimal.NewFromFloat(15.0)),
			IsAvailable: true,
		}

		// Expected order item
		expectedOrderItem := &models.OrderItem{
			ID:         "test-order-item-id",
			OrderID:    orderID,
			MenuItemID: menuItemID,
			Quantity:   2,
			UnitPrice:  types.DecimalText(decimal.NewFromFloat(15.0)),
			TotalPrice: types.DecimalText(decimal.NewFromFloat(30.0)),
		}

		// Expected order item with details
		expectedOrderItemWithDetails := []*models.OrderItemWithDetails{
			{
				ID:           "test-order-item-id",
				OrderID:      orderID,
				MenuItemID:   menuItemID,
				MenuItemName: "Test Item",
				Quantity:     2,
				UnitPrice:    types.DecimalText(decimal.NewFromFloat(15.0)),
				TotalPrice:   types.DecimalText(decimal.NewFromFloat(30.0)),
			},
		}

		// Set up expectations
		mockMenuRepo.On("GetMenuItem", menuItemID).Return(expectedMenuItem, nil)
		mockOrderRepo.On("CreateOrder", mock.AnythingOfType("*models.Order")).Return(expectedOrder, nil)
		mockOrderItemRepo.On("CreateOrderItem", mock.AnythingOfType("*models.OrderItem")).Return(expectedOrderItem, nil)
		mockOrderItemRepo.On("GetOrderItemsWithDetails", orderID).Return(expectedOrderItemWithDetails, nil)

		// Execute the method
		result, err := orderService.CreateOrder(userID, orderCreateData)

		// Assertions
		assert.NoError(t, err)
		assert.True(t, result.Success)
		assert.NotNil(t, result.Data)

		// Verify the call expectations
		mockMenuRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
		mockOrderItemRepo.AssertExpectations(t)
	})

	t.Run("GetOrder with Order Items", func(t *testing.T) {
		// Expected order
		expectedOrder := &models.Order{
			ID:             orderID,
			OrderNumber:    "ORD-TEST-0001",
			UserID:         userID,
			Status:         types.OrderStatusCompleted,
			TotalAmount:    types.DecimalText(decimal.NewFromFloat(30.0)),
			DiscountAmount: types.DecimalText(decimal.Zero),
			TaxAmount:      types.DecimalText(decimal.Zero),
			PaymentStatus:  types.PaymentStatusPaid,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Expected order items with details
		expectedOrderItemDetails := []*models.OrderItemWithDetails{
			{
				ID:           "test-order-item-id",
				OrderID:      orderID,
				MenuItemID:   "test-menu-item-id",
				MenuItemName: "Test Item",
				Quantity:     2,
				UnitPrice:    types.DecimalText(decimal.NewFromFloat(15.0)),
				TotalPrice:   types.DecimalText(decimal.NewFromFloat(30.0)),
			},
		}

		// Set up expectations
		mockOrderRepo.On("GetOrder", orderID).Return(expectedOrder, nil)
		mockOrderItemRepo.On("GetOrderItemsWithDetails", orderID).Return(expectedOrderItemDetails, nil)

		// Execute the method
		result, err := orderService.GetOrder(orderID)

		// Assertions
		assert.NoError(t, err)
		assert.True(t, result.Success)
		assert.NotNil(t, result.Data)

		// Verify the call expectations
		mockOrderRepo.AssertExpectations(t)
		mockOrderItemRepo.AssertExpectations(t)
	})
}