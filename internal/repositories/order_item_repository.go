package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AndikaPrasetia/pos-cafee/internal/db"
	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// orderItemRepo implements the OrderItemRepo interface
type orderItemRepo struct {
	queries *db.Queries
}

// GetOrderItem retrieves an order item by ID
func (r *orderItemRepo) GetOrderItem(id string) (*models.OrderItem, error) {
	orderItemID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	dbOrderItem, err := r.queries.GetOrderItem(context.Background(), orderItemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order item not found")
		}
		return nil, err
	}

	unitPrice, err := decimal.NewFromString(dbOrderItem.UnitPrice)
	if err != nil {
		return nil, err
	}

	totalPrice, err := decimal.NewFromString(dbOrderItem.TotalPrice)
	if err != nil {
		return nil, err
	}

	orderItem := &models.OrderItem{
		ID:         dbOrderItem.ID.String(),
		OrderID:    dbOrderItem.OrderID.String(),
		MenuItemID: dbOrderItem.MenuItemID.String(),
		Quantity:   int(dbOrderItem.Quantity),
		UnitPrice:  types.DecimalText(unitPrice),
		TotalPrice: types.DecimalText(totalPrice),
	}

	return orderItem, nil
}

// GetOrderItemsByOrderID retrieves all order items for an order
func (r *orderItemRepo) GetOrderItemsByOrderID(orderID string) ([]*models.OrderItem, error) {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, err
	}

	dbOrderItems, err := r.queries.GetOrderItemsByOrderID(context.Background(), orderUUID)
	if err != nil {
		return nil, err
	}

	var orderItems []*models.OrderItem
	for _, dbOrderItem := range dbOrderItems {
		unitPrice, err := decimal.NewFromString(dbOrderItem.UnitPrice)
		if err != nil {
			return nil, err
		}

		totalPrice, err := decimal.NewFromString(dbOrderItem.TotalPrice)
		if err != nil {
			return nil, err
		}

		orderItem := &models.OrderItem{
			ID:         dbOrderItem.ID.String(),
			OrderID:    dbOrderItem.OrderID.String(),
			MenuItemID: dbOrderItem.MenuItemID.String(),
			Quantity:   int(dbOrderItem.Quantity),
			UnitPrice:  types.DecimalText(unitPrice),
			TotalPrice: types.DecimalText(totalPrice),
		}
		orderItems = append(orderItems, orderItem)
	}

	return orderItems, nil
}

// CreateOrderItem creates a new order item
func (r *orderItemRepo) CreateOrderItem(orderItem *models.OrderItem) (*models.OrderItem, error) {
	orderID, err := uuid.Parse(orderItem.OrderID)
	if err != nil {
		return nil, err
	}

	menuItemID, err := uuid.Parse(orderItem.MenuItemID)
	if err != nil {
		return nil, err
	}

	dbOrderItem, err := r.queries.CreateOrderItem(context.Background(), db.CreateOrderItemParams{
		OrderID:    orderID,
		MenuItemID: menuItemID,
		Quantity:   int32(orderItem.Quantity),
		UnitPrice:  orderItem.UnitPrice.String(),
		TotalPrice: orderItem.TotalPrice.String(),
	})
	if err != nil {
		return nil, err
	}

	createdOrderItem := &models.OrderItem{
		ID:         dbOrderItem.ID.String(),
		OrderID:    dbOrderItem.OrderID.String(),
		MenuItemID: dbOrderItem.MenuItemID.String(),
		Quantity:   int(dbOrderItem.Quantity),
		UnitPrice:  types.DecimalText(decimal.RequireFromString(dbOrderItem.UnitPrice)),
		TotalPrice: types.DecimalText(decimal.RequireFromString(dbOrderItem.TotalPrice)),
	}

	return createdOrderItem, nil
}

// UpdateOrderItem updates an existing order item
func (r *orderItemRepo) UpdateOrderItem(orderItem *models.OrderItem) (*models.OrderItem, error) {
	orderItemID, err := uuid.Parse(orderItem.ID)
	if err != nil {
		return nil, err
	}

	dbOrderItem, err := r.queries.UpdateOrderItem(context.Background(), db.UpdateOrderItemParams{
		ID:         orderItemID,
		Quantity:   int32(orderItem.Quantity),
		UnitPrice:  orderItem.UnitPrice.String(),
		TotalPrice: orderItem.TotalPrice.String(),
	})
	if err != nil {
		return nil, err
	}

	updatedOrderItem := &models.OrderItem{
		ID:         dbOrderItem.ID.String(),
		OrderID:    dbOrderItem.OrderID.String(),
		MenuItemID: dbOrderItem.MenuItemID.String(),
		Quantity:   int(dbOrderItem.Quantity),
		UnitPrice:  types.DecimalText(decimal.RequireFromString(dbOrderItem.UnitPrice)),
		TotalPrice: types.DecimalText(decimal.RequireFromString(dbOrderItem.TotalPrice)),
	}

	return updatedOrderItem, nil
}

// DeleteOrderItem deletes an order item by ID
func (r *orderItemRepo) DeleteOrderItem(id string) error {
	orderItemID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = r.queries.DeleteOrderItem(context.Background(), orderItemID)
	if err != nil {
		return err
	}

	return nil
}

// GetOrderItemsWithDetails retrieves order items with menu item details
func (r *orderItemRepo) GetOrderItemsWithDetails(orderID string) ([]*models.OrderItemWithDetails, error) {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, err
	}

	dbOrderItemDetails, err := r.queries.GetOrderItemsWithDetails(context.Background(), orderUUID)
	if err != nil {
		return nil, err
	}

	var orderItemDetails []*models.OrderItemWithDetails
	for _, dbItem := range dbOrderItemDetails {
		unitPrice, err := decimal.NewFromString(dbItem.UnitPrice)
		if err != nil {
			return nil, err
		}

		totalPrice, err := decimal.NewFromString(dbItem.TotalPrice)
		if err != nil {
			return nil, err
		}

		itemDetail := &models.OrderItemWithDetails{
			ID:           dbItem.ID.String(),
			OrderID:      dbItem.OrderID.String(),
			MenuItemID:   dbItem.MenuItemID.String(),
			MenuItemName: dbItem.MenuItemName,
			Quantity:     int(dbItem.Quantity),
			UnitPrice:    types.DecimalText(unitPrice),
			TotalPrice:   types.DecimalText(totalPrice),
		}
		orderItemDetails = append(orderItemDetails, itemDetail)
	}

	return orderItemDetails, nil
}