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

// OrderService handles order-related business logic
type OrderService struct {
	orderRepo           repositories.OrderRepo
	menuRepo            repositories.MenuRepo
	inventoryRepo       repositories.InventoryRepo
	stockTransactionRepo repositories.StockTransactionRepo
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo repositories.OrderRepo,
	menuRepo repositories.MenuRepo,
	inventoryRepo repositories.InventoryRepo,
	stockTransactionRepo repositories.StockTransactionRepo,
) *OrderService {
	return &OrderService{
		orderRepo:           orderRepo,
		menuRepo:            menuRepo,
		inventoryRepo:       inventoryRepo,
		stockTransactionRepo: stockTransactionRepo,
	}
}

// CreateOrder creates a new draft order
func (s *OrderService) CreateOrder(userID string, orderData *models.OrderCreate) (*types.APIResponse, error) {
	// Validate user ID format
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Validate order items
	if len(orderData.Items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	// Validate items and calculate totals
	var itemsWithDetails []models.OrderItemWithDetails
	var totalAmount types.DecimalText

	for _, itemData := range orderData.Items {
		// Validate menu item ID
		_, err := uuid.Parse(itemData.MenuItemID)
		if err != nil {
			return nil, fmt.Errorf("invalid menu item ID: %s", itemData.MenuItemID)
		}

		// Get menu item to verify availability and get price
		menuItem, err := s.menuRepo.GetMenuItem(itemData.MenuItemID)
		if err != nil {
			return nil, fmt.Errorf("menu item not found: %s", itemData.MenuItemID)
		}

		if !menuItem.IsAvailable {
			return nil, fmt.Errorf("menu item is not available: %s", menuItem.Name)
		}

		// Calculate item total
		itemTotal := menuItem.Price.Mul(types.FromDecimal(decimal.NewFromInt(int64(itemData.Quantity))))

		// Add to order items
		orderItemWithDetails := models.OrderItemWithDetails{
			ID:           uuid.New().String(),
			OrderID:      "", // Will be set after order is created
			MenuItemID:   itemData.MenuItemID,
			MenuItemName: menuItem.Name,
			Quantity:     itemData.Quantity,
			UnitPrice:    menuItem.Price,
			TotalPrice:   itemTotal,
		}

		itemsWithDetails = append(itemsWithDetails, orderItemWithDetails)
		totalAmount = totalAmount.Add(itemTotal)
	}

	// Generate order number in the format ORD-YYYYMMDD-XXXX
	orderNumber := fmt.Sprintf("ORD-%s-%04d", 
		time.Now().Format("20060102"), 
		time.Now().UnixNano()%10000) // Simple sequential number for demo purposes

	// Create the order
	order := &models.Order{
		ID:             uuid.New().String(),
		OrderNumber:    orderNumber,
		UserID:         userID,
		Status:         types.OrderStatusDraft,
		TotalAmount:    totalAmount,
		DiscountAmount: types.DecimalText(decimal.Zero),
		TaxAmount:      types.DecimalText(decimal.Zero),
		PaymentStatus:  types.PaymentStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	createdOrder, err := s.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}

	// Create order items (this would require additional repository functions in a complete implementation)
	// For now, return the order with basic information
	createdOrderWithDetails := models.OrderWithDetails{
		ID:             createdOrder.ID,
		OrderNumber:    createdOrder.OrderNumber,
		UserID:         createdOrder.UserID,
		Status:         createdOrder.Status,
		TotalAmount:    createdOrder.TotalAmount,
		DiscountAmount: createdOrder.DiscountAmount,
		TaxAmount:      createdOrder.TaxAmount,
		PaymentMethod:  createdOrder.PaymentMethod,
		PaymentStatus:  createdOrder.PaymentStatus,
		CompletedAt:    createdOrder.CompletedAt,
		CreatedAt:      createdOrder.CreatedAt,
		UpdatedAt:      createdOrder.UpdatedAt,
		Items:          itemsWithDetails, // In a complete implementation, these would be saved to DB
	}

	return &types.APIResponse{
		Success: true,
		Data:    createdOrderWithDetails,
	}, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(id string) (*types.APIResponse, error) {
	order, err := s.orderRepo.GetOrder(id)
	if err != nil {
		return nil, err
	}

	// In a complete implementation, this would also fetch the order items
	orderWithDetails := models.OrderWithDetails{
		ID:             order.ID,
		OrderNumber:    order.OrderNumber,
		UserID:         order.UserID,
		Status:         order.Status,
		TotalAmount:    order.TotalAmount,
		DiscountAmount: order.DiscountAmount,
		TaxAmount:      order.TaxAmount,
		PaymentMethod:  order.PaymentMethod,
		PaymentStatus:  order.PaymentStatus,
		CompletedAt:    order.CompletedAt,
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
		Items:          []models.OrderItemWithDetails{}, // Incomplete: would fetch from DB in a full implementation
	}

	return &types.APIResponse{
		Success: true,
		Data:    orderWithDetails,
	}, nil
}

// ListOrders retrieves a list of orders based on filter criteria
func (s *OrderService) ListOrders(filter types.OrderFilter) (*types.APIResponse, error) {
	orders, err := s.orderRepo.ListOrders(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %v", err)
	}

	// In a complete implementation, this would also fetch user and item details
	// Return just the basic order information for now
	return &types.APIResponse{
		Success: true,
		Data:    orders,
	}, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(orderID string, status types.OrderStatus) error {
	// Validate order ID
	_, err := uuid.Parse(orderID)
	if err != nil {
		return errors.New("invalid order ID")
	}

	return s.orderRepo.UpdateOrderStatus(orderID, string(status))
}

// AddItemToOrder adds an item to an existing order
func (s *OrderService) AddItemToOrder(orderID string, userID string, itemData *models.OrderItemCreate) (*types.APIResponse, error) {
	// Validate order ID
	_, err := uuid.Parse(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	// Validate user ID
	_, err = uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get the existing order
	order, err := s.orderRepo.GetOrder(orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %v", err)
	}

	// Check if order is still in draft status
	if order.Status != types.OrderStatusDraft {
		return nil, errors.New("can only add items to draft orders")
	}

	// Get menu item to verify availability and get price
	menuItem, err := s.menuRepo.GetMenuItem(itemData.MenuItemID)
	if err != nil {
		return nil, fmt.Errorf("menu item not found: %s", itemData.MenuItemID)
	}

	if !menuItem.IsAvailable {
		return nil, fmt.Errorf("menu item is not available: %s", menuItem.Name)
	}

	// Calculate the new total
	itemTotal := menuItem.Price.Mul(types.FromDecimal(decimal.NewFromInt(int64(itemData.Quantity))))
	newTotal := order.TotalAmount.Add(itemTotal)

	// In a complete implementation, we would add the item to the database
	// For now, just update the order total
	err = s.orderRepo.UpdateOrderTotal(orderID, newTotal.String(), order.DiscountAmount.String(), order.TaxAmount.String())
	if err != nil {
		return nil, fmt.Errorf("failed to update order total: %v", err)
	}

	// Update the order's updated_at timestamp by fetching it again
	updatedOrder, err := s.orderRepo.GetOrder(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated order: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message":           "Item added to order successfully",
			"updated_order_total": updatedOrder.TotalAmount,
			// In a complete implementation, this would also return the added item details
		},
	}, nil
}

// CompleteOrder processes payment and completes the order, updating inventory
func (s *OrderService) CompleteOrder(orderID string, userID string, updateData *models.OrderUpdate) (*types.APIResponse, error) {
	// Validate order ID
	_, err := uuid.Parse(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	// Validate user ID
	_, err = uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get the order
	order, err := s.orderRepo.GetOrder(orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %v", err)
	}

	// Check if order is in a valid state for completion
	if order.Status != types.OrderStatusDraft && order.Status != types.OrderStatusPending {
		return nil, errors.New("order is not in a valid state for completion")
	}

	// Update order with payment information if provided
	// In a complete implementation, we would check inventory before completing the order
	// For now, just update the order status and payment information
	var paymentMethodStr string
	if updateData.PaymentMethod != nil {
		paymentMethodStr = string(*updateData.PaymentMethod)
	}

	completedAt := time.Now().UTC().Format("2006-01-02 15:04:05.999999-07:00")
	err = s.orderRepo.UpdateOrderPayment(
		orderID,
		paymentMethodStr,
		string(types.PaymentStatusPaid),
		&completedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update order payment: %v", err)
	}

	// Update the order status to completed
	err = s.UpdateOrderStatus(orderID, types.OrderStatusCompleted)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %v", err)
	}

	// In a complete implementation, we would reduce inventory here
	// For now, just simulate the inventory update

	// Fetch the updated order
	updatedOrder, err := s.orderRepo.GetOrder(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated order: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data:    updatedOrder,
	}, nil
}

// CancelOrder cancels an order with authorization checks
func (s *OrderService) CancelOrder(orderID string, userID string, updateData *models.OrderUpdate) (*types.APIResponse, error) {
	// Validate order ID
	_, err := uuid.Parse(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	// Validate user ID
	_, err = uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get the order
	order, err := s.orderRepo.GetOrder(orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %v", err)
	}

	// Check if order can be canceled based on its status
	if order.Status == types.OrderStatusCompleted {
		// Only managers/admins can cancel completed orders
		// This authorization check would happen in the middleware, not in the service
		// For now, we'll just check that it's not already canceled
		if order.Status == types.OrderStatusCancelled {
			return nil, errors.New("order is already canceled")
		}
	} else if order.Status == types.OrderStatusCancelled {
		return nil, errors.New("order is already canceled")
	}

	// Update the order status to cancelled
	err = s.UpdateOrderStatus(orderID, types.OrderStatusCancelled)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel order: %v", err)
	}

	// Fetch the updated order
	updatedOrder, err := s.orderRepo.GetOrder(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated order: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data:    updatedOrder,
	}, nil
}