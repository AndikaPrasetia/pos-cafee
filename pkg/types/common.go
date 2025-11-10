package types

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

// DecimalText represents a decimal value that implements JSON marshaling
type DecimalText decimal.Decimal

// Scan implements the Scanner interface for database operations
func (d *DecimalText) Scan(value any) error {
	if value == nil {
		*d = DecimalText(decimal.Zero)
		return nil
	}

	switch v := value.(type) {
	case float64:
		*d = DecimalText(decimal.NewFromFloat(v))
	case string:
		dec, err := decimal.NewFromString(v)
		if err != nil {
			return err
		}
		*d = DecimalText(dec)
	case []byte:
		dec, err := decimal.NewFromString(string(v))
		if err != nil {
			return err
		}
		*d = DecimalText(dec)
	default:
		return fmt.Errorf("cannot scan %T into DecimalText", value)
	}

	return nil
}

// Value implements the Valuer interface for database operations
func (d DecimalText) Value() (driver.Value, error) {
	return decimal.Decimal(d).String(), nil
}

// String returns the decimal as a string
func (d DecimalText) String() string {
	return decimal.Decimal(d).String()
}

// Float64 returns the decimal as a float64
func (d DecimalText) Float64() float64 {
	return decimal.Decimal(d).InexactFloat64()
}

// MarshalJSON implements json.Marshaler
func (d DecimalText) MarshalJSON() ([]byte, error) {
	return decimal.Decimal(d).MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler
func (d *DecimalText) UnmarshalJSON(b []byte) error {
	var dec decimal.Decimal
	if err := dec.UnmarshalJSON(b); err != nil {
		return err
	}
	*d = DecimalText(dec)
	return nil
}

// Add returns the sum of this DecimalText and another
func (d DecimalText) Add(other DecimalText) DecimalText {
	return DecimalText(decimal.Decimal(d).Add(decimal.Decimal(other)))
}

// Sub subtracts another DecimalText from this one
func (d DecimalText) Sub(other DecimalText) DecimalText {
	return DecimalText(decimal.Decimal(d).Sub(decimal.Decimal(other)))
}

// Mul multiplies this DecimalText by another
func (d DecimalText) Mul(other DecimalText) DecimalText {
	return DecimalText(decimal.Decimal(d).Mul(decimal.Decimal(other)))
}

// Div divides this DecimalText by another
func (d DecimalText) Div(other DecimalText) DecimalText {
	return DecimalText(decimal.Decimal(d).Div(decimal.Decimal(other)))
}

// Cmp compares this DecimalText with another (-1, 0, or 1)
func (d DecimalText) Cmp(other DecimalText) int {
	return decimal.Decimal(d).Cmp(decimal.Decimal(other))
}

// Equals checks if this DecimalText equals another
func (d DecimalText) Equals(other DecimalText) bool {
	return decimal.Decimal(d).Equal(decimal.Decimal(other))
}

// FromDecimal creates a DecimalText from a decimal.Decimal
func FromDecimal(dec decimal.Decimal) DecimalText {
	return DecimalText(dec)
}

// RegisterValidatorRegistrations registers custom validation functions for this package
func RegisterValidatorRegistrations(validate *validator.Validate) {
	// Register a custom validation for DecimalText type
	validate.RegisterCustomTypeFunc(validateDecimalText, DecimalText{})
	// Also register validation for *DecimalText (pointers)
	validate.RegisterCustomTypeFunc(validateDecimalTextPtr, (*DecimalText)(nil))
}

// validateDecimalText provides validation for DecimalText types
func validateDecimalText(field reflect.Value) interface{} {
	if f, ok := field.Interface().(DecimalText); ok {
		// Return the string representation for validation
		return decimal.Decimal(f).String()
	}
	return nil
}

// validateDecimalTextPtr provides validation for *DecimalText types (pointers)
func validateDecimalTextPtr(field reflect.Value) interface{} {
	if f, ok := field.Interface().(*DecimalText); ok && f != nil {
		// Return the string representation for validation
		return decimal.Decimal(*f).String()
	}
	return nil
}

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusDraft     OrderStatus = "draft"
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// PaymentStatus represents the payment status of an order
type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusFailed  PaymentStatus = "failed"
)

// PaymentMethod represents the payment method used for an order
type PaymentMethod string

const (
	PaymentMethodCash     PaymentMethod = "cash"
	PaymentMethodCard     PaymentMethod = "card"
	PaymentMethodQris     PaymentMethod = "qris"
	PaymentMethodTransfer PaymentMethod = "transfer"
)

// TransactionType represents the type of stock transaction
type TransactionType string

const (
	TransactionTypeIn         TransactionType = "in"
	TransactionTypeOut        TransactionType = "out"
	TransactionTypeAdjustment TransactionType = "adjustment"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	UserRoleCashier UserRole = "cashier"
	UserRoleManager UserRole = "manager"
	UserRoleAdmin   UserRole = "admin"
)

// Pagination represents pagination parameters
type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// APIResponseWithData creates an API response with data
func APIResponseWithData(data any) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
	}
}

// APIResponseWithMessage creates an API response with a message
func APIResponseWithMessage(message string) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
	}
}

// APIResponseWithError creates an API response with an error
func APIResponseWithError(message string) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
	}
}

// OrderFilter represents filter options for listing orders
type OrderFilter struct {
	Status    *string    `json:"status,omitempty"`
	UserID    *string    `json:"user_id,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

