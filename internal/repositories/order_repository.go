package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/db"
	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// orderRepo implements the OrderRepo interface
type orderRepo struct {
	queries *db.Queries
}

// GetOrder retrieves an order by ID
func (r *orderRepo) GetOrder(id string) (*models.Order, error) {
	orderID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	dbOrder, err := r.queries.GetOrder(context.Background(), orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	totalAmount, err := decimal.NewFromString(dbOrder.TotalAmount)
	if err != nil {
		return nil, err
	}
	
	discountAmount, err := decimal.NewFromString(dbOrder.DiscountAmount)
	if err != nil {
		return nil, err
	}
	
	taxAmount, err := decimal.NewFromString(dbOrder.TaxAmount)
	if err != nil {
		return nil, err
	}

	order := &models.Order{
		ID:             dbOrder.ID.String(),
		OrderNumber:    dbOrder.OrderNumber,
		UserID:         dbOrder.UserID.String(),
		Status:         types.OrderStatus(dbOrder.Status),
		TotalAmount:    types.DecimalText(totalAmount),
		DiscountAmount: types.DecimalText(discountAmount),
		TaxAmount:      types.DecimalText(taxAmount),
		PaymentStatus:  types.PaymentStatus(dbOrder.PaymentStatus),
		CreatedAt:      dbOrder.CreatedAt,
		UpdatedAt:      dbOrder.UpdatedAt,
	}

	if dbOrder.PaymentMethod.Valid {
		pm := types.PaymentMethod(dbOrder.PaymentMethod.String)
		order.PaymentMethod = &pm
	}

	if dbOrder.CompletedAt.Valid {
		order.CompletedAt = &dbOrder.CompletedAt.Time
	}

	return order, nil
}

// GetOrderByNumber retrieves an order by order number
func (r *orderRepo) GetOrderByNumber(orderNumber string) (*models.Order, error) {
	dbOrder, err := r.queries.GetOrderByNumber(context.Background(), orderNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	totalAmount, err := decimal.NewFromString(dbOrder.TotalAmount)
	if err != nil {
		return nil, err
	}
	
	discountAmount, err := decimal.NewFromString(dbOrder.DiscountAmount)
	if err != nil {
		return nil, err
	}
	
	taxAmount, err := decimal.NewFromString(dbOrder.TaxAmount)
	if err != nil {
		return nil, err
	}

	order := &models.Order{
		ID:             dbOrder.ID.String(),
		OrderNumber:    dbOrder.OrderNumber,
		UserID:         dbOrder.UserID.String(),
		Status:         types.OrderStatus(dbOrder.Status),
		TotalAmount:    types.DecimalText(totalAmount),
		DiscountAmount: types.DecimalText(discountAmount),
		TaxAmount:      types.DecimalText(taxAmount),
		PaymentStatus:  types.PaymentStatus(dbOrder.PaymentStatus),
		CreatedAt:      dbOrder.CreatedAt,
		UpdatedAt:      dbOrder.UpdatedAt,
	}

	if dbOrder.PaymentMethod.Valid {
		pm := types.PaymentMethod(dbOrder.PaymentMethod.String)
		order.PaymentMethod = &pm
	}

	if dbOrder.CompletedAt.Valid {
		order.CompletedAt = &dbOrder.CompletedAt.Time
	}

	return order, nil
}

// ListOrders retrieves a list of orders based on filter
func (r *orderRepo) ListOrders(filter types.OrderFilter) ([]*models.Order, error) {
	var statusParam sql.NullString
	var userID uuid.UUID
	var startDate, endDate time.Time

	if filter.Status != nil {
		statusParam = sql.NullString{
			String: *filter.Status,
			Valid:  true,
		}
	} else {
		statusParam = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	if filter.UserID != nil {
		parsedUUID, err := uuid.Parse(*filter.UserID)
		if err != nil {
			return nil, err
		}
		userID = parsedUUID
	} else {
		userID = uuid.Nil
	}

	if filter.StartDate != nil {
		startDate = *filter.StartDate
	} else {
		startDate = time.Time{}
	}

	if filter.EndDate != nil {
		endDate = *filter.EndDate
	} else {
		endDate = time.Time{}
	}

	dbOrders, err := r.queries.ListOrders(context.Background(), db.ListOrdersParams{
		Column1: statusParam.String,
		Column2: userID,
		Column3: startDate,
		Column4: endDate,
		Limit:   int32(filter.Limit),
		Offset:  int32(filter.Offset),
	})
	if err != nil {
		return nil, err
	}

	var orders []*models.Order
	for _, dbOrder := range dbOrders {
		totalAmount, err := decimal.NewFromString(dbOrder.TotalAmount)
		if err != nil {
			return nil, err
		}
		
		discountAmount, err := decimal.NewFromString(dbOrder.DiscountAmount)
		if err != nil {
			return nil, err
		}
		
		taxAmount, err := decimal.NewFromString(dbOrder.TaxAmount)
		if err != nil {
			return nil, err
		}

		order := &models.Order{
			ID:             dbOrder.ID.String(),
			OrderNumber:    dbOrder.OrderNumber,
			UserID:         dbOrder.UserID.String(),
			Status:         types.OrderStatus(dbOrder.Status),
			TotalAmount:    types.DecimalText(totalAmount),
			DiscountAmount: types.DecimalText(discountAmount),
			TaxAmount:      types.DecimalText(taxAmount),
			PaymentStatus:  types.PaymentStatus(dbOrder.PaymentStatus),
			CreatedAt:      dbOrder.CreatedAt,
			UpdatedAt:      dbOrder.UpdatedAt,
		}

		if dbOrder.PaymentMethod.Valid {
			pm := types.PaymentMethod(dbOrder.PaymentMethod.String)
			order.PaymentMethod = &pm
		}

		if dbOrder.CompletedAt.Valid {
			order.CompletedAt = &dbOrder.CompletedAt.Time
		}

		orders = append(orders, order)
	}

	return orders, nil
}

// CreateOrder creates a new order
func (r *orderRepo) CreateOrder(order *models.Order) (*models.Order, error) {
	userID, err := uuid.Parse(order.UserID)
	if err != nil {
		return nil, err
	}

	dbOrder, err := r.queries.CreateOrder(context.Background(), db.CreateOrderParams{
		OrderNumber:    order.OrderNumber,
		UserID:         userID,
		TotalAmount:    order.TotalAmount.String(),
		DiscountAmount: order.DiscountAmount.String(),
		TaxAmount:      order.TaxAmount.String(),
	})
	if err != nil {
		return nil, err
	}

	createdOrder := &models.Order{
		ID:             dbOrder.ID.String(),
		OrderNumber:    dbOrder.OrderNumber,
		UserID:         dbOrder.UserID.String(),
		Status:         types.OrderStatus(dbOrder.Status),
		TotalAmount:    types.DecimalText(decimal.RequireFromString(dbOrder.TotalAmount)),
		DiscountAmount: types.DecimalText(decimal.RequireFromString(dbOrder.DiscountAmount)),
		TaxAmount:      types.DecimalText(decimal.RequireFromString(dbOrder.TaxAmount)),
		PaymentStatus:  types.PaymentStatus(dbOrder.PaymentStatus),
		CreatedAt:      dbOrder.CreatedAt,
		UpdatedAt:      dbOrder.UpdatedAt,
	}

	if dbOrder.PaymentMethod.Valid {
		pm := types.PaymentMethod(dbOrder.PaymentMethod.String)
		createdOrder.PaymentMethod = &pm
	}

	if dbOrder.CompletedAt.Valid {
		createdOrder.CompletedAt = &dbOrder.CompletedAt.Time
	}

	return createdOrder, nil
}

// UpdateOrderStatus updates the status of an order
func (r *orderRepo) UpdateOrderStatus(orderID string, status string) error {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return err
	}

	err = r.queries.UpdateOrderStatus(context.Background(), db.UpdateOrderStatusParams{
		ID:     orderUUID,
		Status: status,
	})
	if err != nil {
		return err
	}

	return nil
}

// UpdateOrderPayment updates payment information for an order
func (r *orderRepo) UpdateOrderPayment(orderID string, paymentMethod, paymentStatus string, completedAt *string) error {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return err
	}

	var paymentMethodNull sql.NullString
	if paymentMethod != "" {
		paymentMethodNull = sql.NullString{
			String: paymentMethod,
			Valid:  true,
		}
	} else {
		paymentMethodNull = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	err = r.queries.UpdateOrderPayment(context.Background(), db.UpdateOrderPaymentParams{
		ID:            orderUUID,
		PaymentMethod: paymentMethodNull,
		PaymentStatus: paymentStatus,
	})
	if err != nil {
		return err
	}

	return nil
}

// UpdateOrderTotal updates total amounts for an order
func (r *orderRepo) UpdateOrderTotal(orderID string, totalAmount, discountAmount, taxAmount string) error {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return err
	}

	err = r.queries.UpdateOrderTotal(context.Background(), db.UpdateOrderTotalParams{
		ID:             orderUUID,
		TotalAmount:    totalAmount,
		DiscountAmount: discountAmount,
		TaxAmount:      taxAmount,
	})
	if err != nil {
		return err
	}

	return nil
}