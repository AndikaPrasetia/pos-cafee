package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/repositories"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/shopspring/decimal"
)

// ReportService handles financial reporting business logic
type ReportService struct {
	orderRepo       repositories.OrderRepo
	menuRepo        repositories.MenuRepo
	inventoryRepo   repositories.InventoryRepo
	expenseRepo     repositories.ExpenseRepo
}

// NewReportService creates a new report service
func NewReportService(
	orderRepo repositories.OrderRepo,
	menuRepo repositories.MenuRepo,
	inventoryRepo repositories.InventoryRepo,
	expenseRepo repositories.ExpenseRepo,
) *ReportService {
	return &ReportService{
		orderRepo:     orderRepo,
		menuRepo:      menuRepo,
		inventoryRepo: inventoryRepo,
		expenseRepo:   expenseRepo,
	}
}

// GetDailySalesReport generates a daily sales report
func (s *ReportService) GetDailySalesReport(dateStr string) (*types.APIResponse, error) {
	reportDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	// Calculate start and end of the day
	startOfDay := time.Date(reportDate.Year(), reportDate.Month(), reportDate.Day(), 0, 0, 0, 0, reportDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond) // End of the day (23:59:59)

	// Create date range filter for orders
	filter := types.OrderFilter{
		StartDate: &startOfDay,
		EndDate:   &endOfDay,
		Limit:     1000, // Reasonable limit for daily orders
		Offset:    0,
	}

	orders, err := s.orderRepo.ListOrders(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %v", err)
	}

	// Calculate report data
	var totalOrders int
	var totalSales types.DecimalText
	var topSellingItems []map[string]interface{}

	// Group items by menu item ID to count quantities sold
	// itemSales := make(map[string]map[string]interface{})  // Used in a complete implementation

	for _, order := range orders {
		if order.Status == types.OrderStatusCompleted {
			totalOrders++
			totalSales = totalSales.Add(order.TotalAmount)

			// In a complete implementation, we would fetch the order items
			// For now, we'll just simulate the top selling items
		}
	}

	// Calculate average order value
	var averageOrderValue types.DecimalText
	if totalOrders > 0 {
		averageOrderValue = totalSales.Div(types.FromDecimal(decimal.NewFromInt(int64(totalOrders))))
	}

	report := map[string]interface{}{
		"date":                dateStr,
		"total_orders":        totalOrders,
		"total_sales":         totalSales,
		"average_order_value": averageOrderValue,
		"top_selling_items":   topSellingItems,
	}

	return &types.APIResponse{
		Success: true,
		Data:    report,
	}, nil
}

// GetFinancialSummaryReport generates a financial summary report for a date range
func (s *ReportService) GetFinancialSummaryReport(startDateStr, endDateStr string) (*types.APIResponse, error) {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, errors.New("invalid start date format, expected YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, errors.New("invalid end date format, expected YYYY-MM-DD")
	}

	if startDate.After(endDate) {
		return nil, errors.New("start date cannot be after end date")
	}

	// Calculate end of the end date (23:59:59)
	endOfDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())

	// Create date range filter for orders
	filter := types.OrderFilter{
		StartDate: &startDate,
		EndDate:   &endOfDay,
		Limit:     10000, // Reasonable limit for the period
		Offset:    0,
	}

	orders, err := s.orderRepo.ListOrders(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %v", err)
	}

	// Calculate total sales from completed orders
	var totalSales types.DecimalText
	var salesByCategory []map[string]interface{}
	// categorySales := make(map[string]map[string]interface{})  // Used in a complete implementation

	for _, order := range orders {
		if order.Status == types.OrderStatusCompleted {
			totalSales = totalSales.Add(order.TotalAmount)

			// In a complete implementation, we would fetch the order items
			// and calculate sales by category
		}
	}

	// In a complete implementation, we would fetch expenses from the repository
	// For now, we'll simulate this by creating a placeholder
	// Calculate total expenses
	var totalExpenses types.DecimalText = types.DecimalText(decimal.Zero) // Placeholder - would sum actual expenses
	var expenses []map[string]interface{}

	// Calculate profit
	totalProfit := totalSales.Sub(totalExpenses)

	report := map[string]interface{}{
		"period": map[string]string{
			"start_date": startDateStr,
			"end_date":   endDateStr,
		},
		"total_sales":      totalSales,
		"total_expenses":   totalExpenses,
		"total_profit":     totalProfit,
		"sales_by_category": salesByCategory,
		"expenses":         expenses,
	}

	return &types.APIResponse{
		Success: true,
		Data:    report,
	}, nil
}

// GetSalesByCategoryReport generates a sales report grouped by category for a date range
func (s *ReportService) GetSalesByCategoryReport(startDateStr, endDateStr string) (*types.APIResponse, error) {
	// This would be similar to the above but specifically focused on category breakdown
	// For now, we'll return a placeholder implementation
	
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, errors.New("invalid start date format, expected YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, errors.New("invalid end date format, expected YYYY-MM-DD")
	}

	if startDate.After(endDate) {
		return nil, errors.New("start date cannot be after end date")
	}

	// Placeholder implementation - would fetch actual data in a complete implementation
	salesByCategory := []map[string]interface{}{}

	report := map[string]interface{}{
		"period": map[string]string{
			"start_date": startDateStr,
			"end_date":   endDateStr,
		},
		"sales_by_category": salesByCategory,
	}

	return &types.APIResponse{
		Success: true,
		Data:    report,
	}, nil
}

// GetTopSellingItemsReport generates a report of top selling items for a date range
func (s *ReportService) GetTopSellingItemsReport(startDateStr, endDateStr string, limit int) (*types.APIResponse, error) {
	// This would fetch the most sold items by quantity in the given date range
	// For now, we'll return a placeholder implementation
	
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, errors.New("invalid start date format, expected YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, errors.New("invalid end date format, expected YYYY-MM-DD")
	}

	if startDate.After(endDate) {
		return nil, errors.New("start date cannot be after end date")
	}

	if limit <= 0 {
		limit = 10 // Default limit
	}

	// Placeholder implementation - would fetch actual data in a complete implementation
	topSellingItems := []map[string]interface{}{}

	report := map[string]interface{}{
		"period": map[string]string{
			"start_date": startDateStr,
			"end_date":   endDateStr,
		},
		"top_selling_items": topSellingItems,
		"limit":             limit,
	}

	return &types.APIResponse{
		Success: true,
		Data:    report,
	}, nil
}