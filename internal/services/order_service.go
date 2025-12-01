package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/cache"
	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/internal/repositories"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// OrderService handles order-related business logic
type OrderService struct {
	orderRepo            repositories.OrderRepo
	orderItemRepo        repositories.OrderItemRepo
	menuRepo             repositories.MenuRepo
	inventoryRepo        repositories.InventoryRepo
	stockTransactionRepo repositories.StockTransactionRepo
	cache                cache.Cache
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo repositories.OrderRepo,
	orderItemRepo repositories.OrderItemRepo,
	menuRepo repositories.MenuRepo,
	inventoryRepo repositories.InventoryRepo,
	stockTransactionRepo repositories.StockTransactionRepo,
	cache cache.Cache,
) *OrderService {
	return &OrderService{
		orderRepo:            orderRepo,
		orderItemRepo:        orderItemRepo,
		menuRepo:             menuRepo,
		inventoryRepo:        inventoryRepo,
		stockTransactionRepo: stockTransactionRepo,
		cache:                cache,
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

		// Check inventory stock for the menu item
		inventory, err := s.inventoryRepo.GetInventoryByMenuItem(itemData.MenuItemID)
		if err != nil {
			// If no inventory exists for this item, create a new record
			err = s.inventoryRepo.CreateInventoryRecord(itemData.MenuItemID)
			if err != nil {
				return nil, fmt.Errorf("failed to create inventory record for menu item %s: %v", itemData.MenuItemID, err)
			}

			// Try to fetch again
			inventory, err = s.inventoryRepo.GetInventoryByMenuItem(itemData.MenuItemID)
			if err != nil {
				return nil, fmt.Errorf("failed to get inventory for menu item %s: %v", itemData.MenuItemID, err)
			}
		}

		// Check if sufficient stock is available
		if inventory.CurrentStock < itemData.Quantity {
			return nil, fmt.Errorf("insufficient stock for item %s: only %d available, %d requested", menuItem.Name, inventory.CurrentStock, itemData.Quantity)
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

	// Create order items
	for _, itemWithDetails := range itemsWithDetails {
		orderItem := &models.OrderItem{
			ID:         uuid.New().String(),
			OrderID:    createdOrder.ID,
			MenuItemID: itemWithDetails.MenuItemID,
			Quantity:   itemWithDetails.Quantity,
			UnitPrice:  itemWithDetails.UnitPrice,
			TotalPrice: itemWithDetails.TotalPrice,
		}
		
		_, err := s.orderItemRepo.CreateOrderItem(orderItem)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %v", err)
		}
	}

	// Retrieve order items with details
	orderItemDetails, err := s.orderItemRepo.GetOrderItemsWithDetails(createdOrder.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items with details: %v", err)
	}

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
		Items:          convertOrderItemWithDetailsPtrToSlice(orderItemDetails),
	}

	return &types.APIResponse{
		Success: true,
		Data:    createdOrderWithDetails,
	}, nil
}

// Helper function to convert []*models.OrderItemWithDetails to []models.OrderItemWithDetails
func convertOrderItemWithDetailsPtrToSlice(ptrSlice []*models.OrderItemWithDetails) []models.OrderItemWithDetails {
	if ptrSlice == nil {
		return nil
	}
	
	slice := make([]models.OrderItemWithDetails, len(ptrSlice))
	for i, item := range ptrSlice {
		slice[i] = *item
	}
	return slice
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(id string) (*types.APIResponse, error) {
	order, err := s.orderRepo.GetOrder(id)
	if err != nil {
		return nil, err
	}

	// Fetch order items with details
	orderItemDetails, err := s.orderItemRepo.GetOrderItemsWithDetails(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items with details: %v", err)
	}

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
		Items:          convertOrderItemWithDetailsPtrToSlice(orderItemDetails),
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

	// Calculate the item total
	itemTotal := menuItem.Price.Mul(types.FromDecimal(decimal.NewFromInt(int64(itemData.Quantity))))

	// Create order item
	orderItem := &models.OrderItem{
		ID:         uuid.New().String(),
		OrderID:    orderID,
		MenuItemID: itemData.MenuItemID,
		Quantity:   itemData.Quantity,
		UnitPrice:  menuItem.Price,
		TotalPrice: itemTotal,
	}

	_, err = s.orderItemRepo.CreateOrderItem(orderItem)
	if err != nil {
		return nil, fmt.Errorf("failed to create order item: %v", err)
	}

	// Calculate the new total
	newTotal := order.TotalAmount.Add(itemTotal)

	// Update the order total
	err = s.orderRepo.UpdateOrderTotal(orderID, newTotal.String(), order.DiscountAmount.String(), order.TaxAmount.String())
	if err != nil {
		return nil, fmt.Errorf("failed to update order total: %v", err)
	}

	// Fetch the updated order
	updatedOrder, err := s.orderRepo.GetOrder(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated order: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message":             "Item added to order successfully",
			"updated_order_total": updatedOrder.TotalAmount,
			"added_item":          orderItem,
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

	// Get order items to check inventory
	orderItems, err := s.orderItemRepo.GetOrderItemsByOrderID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %v", err)
	}

	// Check inventory availability for each item before completing the order
	for _, orderItem := range orderItems {
		// Fetch current inventory for the menu item
		inventory, err := s.inventoryRepo.GetInventoryByMenuItem(orderItem.MenuItemID)
		if err != nil {
			// If no inventory exists for this item, create a new record
			err = s.inventoryRepo.CreateInventoryRecord(orderItem.MenuItemID)
			if err != nil {
				return nil, fmt.Errorf("failed to create inventory record for menu item %s: %v", orderItem.MenuItemID, err)
			}

			// Try to fetch again
			inventory, err = s.inventoryRepo.GetInventoryByMenuItem(orderItem.MenuItemID)
			if err != nil {
				return nil, fmt.Errorf("failed to get inventory for menu item %s: %v", orderItem.MenuItemID, err)
			}
		}

		// Check if sufficient stock is available
		if inventory.CurrentStock < orderItem.Quantity {
			return nil, fmt.Errorf("insufficient stock for item %s: only %d available, %d requested", inventory.MenuItemName, inventory.CurrentStock, orderItem.Quantity)
		}
	}

	// Update order with payment information if provided
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

	for _, orderItem := range orderItems {
		// Fetch current inventory for the menu item
		inventory, err := s.inventoryRepo.GetInventoryByMenuItem(orderItem.MenuItemID)
		if err != nil {
			// If no inventory exists for this item, create a new record
			err = s.inventoryRepo.CreateInventoryRecord(orderItem.MenuItemID)
			if err != nil {
				return nil, fmt.Errorf("failed to create inventory record for menu item %s: %v", orderItem.MenuItemID, err)
			}
			
			// Try to fetch again
			inventory, err = s.inventoryRepo.GetInventoryByMenuItem(orderItem.MenuItemID)
			if err != nil {
				return nil, fmt.Errorf("failed to get inventory for menu item %s: %v", orderItem.MenuItemID, err)
			}
		}

		// Update inventory stock (reduce by the quantity ordered)
		newStock := inventory.CurrentStock - orderItem.Quantity
		if newStock < 0 {
			newStock = 0 // Ensure we don't have negative stock
		}

		err = s.inventoryRepo.UpdateInventoryStock(orderItem.MenuItemID, newStock, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to update inventory stock for menu item %s: %v", orderItem.MenuItemID, err)
		}

		// Create a stock transaction record
		stockTransaction := &models.StockTransaction{
			ID:              uuid.New().String(),
			MenuItemID:      orderItem.MenuItemID,
			TransactionType: types.TransactionTypeOut, // Using the defined transaction type
			Quantity:        -orderItem.Quantity, // Negative value for reduction (already int, not int32)
			PreviousStock:   inventory.CurrentStock, // Set the previous stock
			CurrentStock:    newStock, // Set the current stock after reduction
			Reason:          fmt.Sprintf("Order %s completion", orderID),
			// UserID is optional in the model, so setting it to a pointer
			UserID: &userID,
		}

		createdTransaction, err := s.stockTransactionRepo.CreateStockTransaction(stockTransaction)
		if err != nil {
			return nil, fmt.Errorf("failed to create stock transaction for menu item %s: %v", orderItem.MenuItemID, err)
		}
		
		// Verify the transaction was created successfully
		if createdTransaction == nil {
			return nil, fmt.Errorf("stock transaction was not created for menu item %s in order %s", orderItem.MenuItemID, orderID)
		}
	}

	// Fetch the updated order
	updatedOrder, err := s.orderRepo.GetOrder(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated order: %v", err)
	}

	// If the order was completed before cancellation, we should invalidate report caches
	// as the sales data will need to be recalculated
	if order.Status == types.OrderStatusCompleted {
		// Invalidate report caches since sales data has changed
		ctx := context.Background()

		// Delete daily sales report cache for the current date
		today := time.Now().Format("2006-01-02")
		s.cache.Delete(ctx, fmt.Sprintf("daily_sales_report:%s", today))

		// Delete any cached reports that might be affected
		// This would include any cached reports for today or the recent period
		// Using pattern matching to invalidate all daily sales reports
		dailySalesReportKeys, err := s.cache.Keys(ctx, "daily_sales_report:*")
		if err == nil {
			for _, key := range dailySalesReportKeys {
				s.cache.Delete(ctx, key)
			}
		} else {
			fmt.Printf("Warning: Failed to get daily sales report cache keys: %v\n", err)
		}

		// Invalidate top selling items reports
		topSellingReportKeys, err := s.cache.Keys(ctx, "top_selling_items:*")
		if err == nil {
			for _, key := range topSellingReportKeys {
				s.cache.Delete(ctx, key)
			}
		} else {
			fmt.Printf("Warning: Failed to get top selling items report cache keys: %v\n", err)
		}
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

	// If the order was completed before cancellation, we should invalidate report caches
	// as the sales data will need to be recalculated
	if order.Status == types.OrderStatusCompleted {
		// Invalidate report caches since sales data has changed
		ctx := context.Background()

		// Delete daily sales report cache for the current date
		today := time.Now().Format("2006-01-02")
		s.cache.Delete(ctx, fmt.Sprintf("daily_sales_report:%s", today))

		// Delete any cached reports that might be affected
		// This would include any cached reports for today or the recent period
		// Using pattern matching to invalidate all daily sales reports
		dailySalesReportKeys, err := s.cache.Keys(ctx, "daily_sales_report:*")
		if err == nil {
			for _, key := range dailySalesReportKeys {
				s.cache.Delete(ctx, key)
			}
		} else {
			fmt.Printf("Warning: Failed to get daily sales report cache keys: %v\n", err)
		}

		// Invalidate top selling items reports
		topSellingReportKeys, err := s.cache.Keys(ctx, "top_selling_items:*")
		if err == nil {
			for _, key := range topSellingReportKeys {
				s.cache.Delete(ctx, key)
			}
		} else {
			fmt.Printf("Warning: Failed to get top selling items report cache keys: %v\n", err)
		}
	}

	return &types.APIResponse{
		Success: true,
		Data:    updatedOrder,
	}, nil
}