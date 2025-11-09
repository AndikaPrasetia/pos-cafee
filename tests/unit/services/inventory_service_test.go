package services

import (
	"testing"

	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInventoryService_ValidateInventoryForOrder(t *testing.T) {
	mockInventoryRepo := new(MockInventoryRepo)
	mockMenuRepo := new(MockMenuRepo)
	inventoryService := NewInventoryService(mockInventoryRepo, nil, mockMenuRepo)

	orderItems := []models.OrderItemCreate{
		{
			MenuItemID: "valid-item-id",
			Quantity:   5,
		},
	}

	inventory := &models.Inventory{
		MenuItemID:    "valid-item-id",
		CurrentStock:  10,
		MinimumStock:  2,
		Unit:          "pieces",
		LastUpdatedAt: mock.Anything,
	}

	menuItem := &models.MenuItem{
		ID:    "valid-item-id",
		Name:  "Test Item",
		Price: types.DecimalText(decimal.NewFromFloat(10.50)),
		Cost:  types.DecimalText(decimal.NewFromFloat(5.25)),
	}

	mockInventoryRepo.On("GetInventoryByMenuItem", "valid-item-id").Return(inventory, nil)
	mockMenuRepo.On("GetMenuItem", "valid-item-id").Return(menuItem, nil)

	err := inventoryService.ValidateInventoryForOrder(orderItems)

	assert.NoError(t, err)

	mockInventoryRepo.AssertExpectations(t)
	mockMenuRepo.AssertExpectations(t)
}

func TestInventoryService_ValidateInventoryForOrder_InsufficientStock(t *testing.T) {
	mockInventoryRepo := new(MockInventoryRepo)
	mockMenuRepo := new(MockMenuRepo)
	inventoryService := NewInventoryService(mockInventoryRepo, nil, mockMenuRepo)

	orderItems := []models.OrderItemCreate{
		{
			MenuItemID: "low-stock-item-id",
			Quantity:   15, // Requesting 15 but only 10 available
		},
	}

	inventory := &models.Inventory{
		MenuItemID:    "low-stock-item-id",
		CurrentStock:  10,
		MinimumStock:  2,
		Unit:          "pieces",
		LastUpdatedAt: mock.Anything,
	}

	menuItem := &models.MenuItem{
		ID:    "low-stock-item-id",
		Name:  "Low Stock Item",
		Price: types.DecimalText(decimal.NewFromFloat(10.50)),
		Cost:  types.DecimalText(decimal.NewFromFloat(5.25)),
	}

	mockInventoryRepo.On("GetInventoryByMenuItem", "low-stock-item-id").Return(inventory, nil)
	mockMenuRepo.On("GetMenuItem", "low-stock-item-id").Return(menuItem, nil)

	err := inventoryService.ValidateInventoryForOrder(orderItems)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient inventory")

	mockInventoryRepo.AssertExpectations(t)
	mockMenuRepo.AssertExpectations(t)
}

// Mock implementations for testing
type MockInventoryRepo struct {
	mock.Mock
}

func (m *MockInventoryRepo) GetInventoryByMenuItem(menuItemID string) (*models.Inventory, error) {
	args := m.Called(menuItemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Inventory), args.Error(1)
}

func (m *MockInventoryRepo) ListInventory(filter models.InventoryFilter) ([]*models.Inventory, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Inventory), args.Error(1)
}

func (m *MockInventoryRepo) UpdateInventoryStock(menuItemID string, stock int, userID string) error {
	args := m.Called(menuItemID, stock, userID)
	return args.Error(0)
}

func (m *MockInventoryRepo) CreateInventoryRecord(menuItemID string) error {
	args := m.Called(menuItemID)
	return args.Error(0)
}

type MockMenuRepo struct {
	mock.Mock
}

func (m *MockMenuRepo) GetMenuItem(id string) (*models.MenuItem, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MenuItem), args.Error(1)
}

// Add other required methods for MockMenuRepo to satisfy the interface
func (m *MockMenuRepo) GetCategory(id string) (*models.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockMenuRepo) ListCategories(isActive bool, limit, offset int) ([]*models.Category, error) {
	args := m.Called(isActive, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Category), args.Error(1)
}

func (m *MockMenuRepo) CreateCategory(category *models.Category) (*models.Category, error) {
	args := m.Called(category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockMenuRepo) UpdateCategory(category *models.Category) (*models.Category, error) {
	args := m.Called(category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockMenuRepo) DeleteCategory(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMenuRepo) ListMenuItems(isAvailable bool, limit, offset int) ([]*models.MenuItem, error) {
	args := m.Called(isAvailable, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.MenuItem), args.Error(1)
}

func (m *MockMenuRepo) ListMenuItemsByCategory(categoryID string, limit, offset int) ([]*models.MenuItem, error) {
	args := m.Called(categoryID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.MenuItem), args.Error(1)
}

func (m *MockMenuRepo) CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error) {
	args := m.Called(item)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MenuItem), args.Error(1)
}

func (m *MockMenuRepo) UpdateMenuItem(item *models.MenuItem) (*models.MenuItem, error) {
	args := m.Called(item)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MenuItem), args.Error(1)
}

func (m *MockMenuRepo) DeleteMenuItem(id string) error {
	args := m.Called(id)
	return args.Error(0)
}