package repositories

import (
	"database/sql"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/db"
	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock_interfaces.go

// UserRepo defines the interface for user-related database operations
type UserRepo interface {
	GetUser(id string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	UpdateUserPassword(userID, hashedPassword string) error
	UpdateUserStatus(userID string, isActive bool) error
}

// MenuRepo defines the interface for menu-related database operations
type MenuRepo interface {
	GetCategory(id string) (*models.Category, error)
	ListCategories(isActive bool, limit, offset int) ([]*models.Category, error)
	CreateCategory(category *models.Category) (*models.Category, error)
	UpdateCategory(category *models.Category) (*models.Category, error)
	DeleteCategory(id string) error

	GetMenuItem(id string) (*models.MenuItem, error)
	ListMenuItems(isAvailable bool, limit, offset int) ([]*models.MenuItem, error)
	ListMenuItemsByCategory(categoryID string, limit, offset int) ([]*models.MenuItem, error)
	CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error)
	UpdateMenuItem(item *models.MenuItem) (*models.MenuItem, error)
	DeleteMenuItem(id string) error
}

// OrderRepo defines the interface for order-related database operations
type OrderRepo interface {
	GetOrder(id string) (*models.Order, error)
	GetOrderByNumber(orderNumber string) (*models.Order, error)
	ListOrders(filter types.OrderFilter) ([]*models.Order, error)
	CreateOrder(order *models.Order) (*models.Order, error)
	UpdateOrderStatus(orderID string, status string) error
	UpdateOrderPayment(orderID string, paymentMethod, paymentStatus string, completedAt *string) error
	UpdateOrderTotal(orderID string, totalAmount, discountAmount, taxAmount string) error
}

// InventoryRepo defines the interface for inventory-related database operations
type InventoryRepo interface {
	GetInventoryByMenuItem(menuItemID string) (*models.Inventory, error)
	ListInventory(filter models.InventoryFilter) ([]*models.Inventory, error)
	UpdateInventoryStock(menuItemID string, stock int, userID string) error
	CreateInventoryRecord(menuItemID string) error
}

// StockTransactionRepo defines the interface for stock transaction-related database operations
type StockTransactionRepo interface {
	CreateStockTransaction(transaction *models.StockTransaction) (*models.StockTransaction, error)
	ListStockTransactions(filter models.StockTransactionFilter) ([]*models.StockTransaction, error)
}

// ExpenseRepo defines the interface for expense-related database operations
type ExpenseRepo interface {
	CreateExpense(expense *models.Expense) (*models.Expense, error)
	GetExpense(id string) (*models.Expense, error)
	ListExpenses(filter models.ExpenseFilter) ([]*models.Expense, error)
	UpdateExpense(expense *models.Expense) (*models.Expense, error)
	DeleteExpense(id string) error
	GetExpensesByDateRange(startDate, endDate time.Time) ([]*models.Expense, error)
}

// OrderItemRepo defines the interface for order item-related database operations
type OrderItemRepo interface {
	GetOrderItem(id string) (*models.OrderItem, error)
	GetOrderItemsByOrderID(orderID string) ([]*models.OrderItem, error)
	CreateOrderItem(orderItem *models.OrderItem) (*models.OrderItem, error)
	UpdateOrderItem(orderItem *models.OrderItem) (*models.OrderItem, error)
	DeleteOrderItem(id string) error
	GetOrderItemsWithDetails(orderID string) ([]*models.OrderItemWithDetails, error)
}

// Repository holds all repository interfaces
type Repository struct {
	UserRepo             UserRepo
	MenuRepo             MenuRepo
	OrderRepo            OrderRepo
	OrderItemRepo        OrderItemRepo
	InventoryRepo        InventoryRepo
	StockTransactionRepo StockTransactionRepo
	ExpenseRepo          ExpenseRepo
	Queries              *db.Queries
}

// NewRepository creates a new Repository instance with concrete implementations
func NewRepository(dbConn *sql.DB) *Repository {
	queries := db.New(dbConn)

	return &Repository{
		UserRepo:             &userRepo{queries: queries},  // This is defined in user_repository.go
		MenuRepo:             &menuRepo{queries: queries},  // This is defined in menu_repository.go
		OrderRepo:            &orderRepo{queries: queries}, // This is defined in order_repository.go
		OrderItemRepo:        &orderItemRepo{queries: queries}, // This is defined in order_item_repository.go
		InventoryRepo:        &inventoryRepo{queries: queries}, // This is defined in inventory_repository.go
		StockTransactionRepo: &stockTransactionRepo{queries: queries}, // This is defined in stock_transaction_repository.go
		ExpenseRepo:          &expenseRepo{queries: queries}, // This is defined in expense_repository.go
		Queries:              queries,
	}
}

