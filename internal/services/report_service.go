package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/cache"
	"github.com/AndikaPrasetia/pos-cafee/internal/db"
	"github.com/AndikaPrasetia/pos-cafee/internal/repositories"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/shopspring/decimal"
)

// ReportService handles financial reporting business logic
type ReportService struct {
	orderRepo     repositories.OrderRepo
	menuRepo      repositories.MenuRepo
	inventoryRepo repositories.InventoryRepo
	expenseRepo   repositories.ExpenseRepo
	queries       *db.Queries
	cache         cache.Cache
}

// NewReportService creates a new report service
func NewReportService(
	orderRepo repositories.OrderRepo,
	menuRepo repositories.MenuRepo,
	inventoryRepo repositories.InventoryRepo,
	expenseRepo repositories.ExpenseRepo,
	queries *db.Queries,
	cache cache.Cache,
) *ReportService {
	return &ReportService{
		orderRepo:     orderRepo,
		menuRepo:      menuRepo,
		inventoryRepo: inventoryRepo,
		expenseRepo:   expenseRepo,
		queries:       queries,
		cache:         cache,
	}
}

// GetDailySalesReport generates a daily sales report
func (s *ReportService) GetDailySalesReport(dateStr string) (*types.APIResponse, error) {
	reportDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	// Create cache key
	cacheKey := fmt.Sprintf("daily_sales_report:%s", dateStr)

	// Try to get from cache first
	var cachedReport map[string]interface{}
	ctx := context.Background()
	err = s.cache.GetJSON(ctx, cacheKey, &cachedReport)
	if err == nil {
		// Cache hit - return cached data
		return &types.APIResponse{
			Success: true,
			Data:    cachedReport,
		}, nil
	}

	// Use the database view to get daily sales data
	reportData, err := s.queries.GetDailySalesReportData(context.Background(), reportDate)
	if err != nil {
		// Return zero values if no data exists for the date
		if err == sql.ErrNoRows {
			report := map[string]interface{}{
				"date":                dateStr,
				"total_orders":        0,
				"total_sales":         types.DecimalText(decimal.Zero),
				"average_order_value": types.DecimalText(decimal.Zero),
				"top_selling_items":   []map[string]interface{}{},
			}

			// Cache the results for 1 hour (reports typically don't change frequently)
			cacheErr := s.cache.SetJSON(ctx, cacheKey, report, time.Hour)
			if cacheErr != nil {
				// Log the error but don't fail the request
				fmt.Printf("Warning: Failed to cache daily sales report: %v\n", cacheErr)
			}

			return &types.APIResponse{
				Success: true,
				Data:    report,
			}, nil
		}
		return nil, fmt.Errorf("failed to fetch daily sales report: %v", err)
	}

	// Calculate average order value
	averageOrderValue := types.DecimalText(decimal.Zero)
	if reportData.TotalOrders > 0 {
		totalSales, err := decimal.NewFromString(reportData.TotalSales)
		if err != nil {
			return nil, fmt.Errorf("failed to parse total sales: %v", err)
		}
		averageOrderValue = types.FromDecimal(totalSales.Div(decimal.NewFromInt(reportData.TotalOrders)))
	}

	// Get top selling items for this date range
	startOfDay := time.Date(reportDate.Year(), reportDate.Month(), reportDate.Day(), 0, 0, 0, 0, reportDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond) // End of the day (23:59:59)

	topSellingItems, err := s.queries.GetTopSellingItemsByDateRange(context.Background(), db.GetTopSellingItemsByDateRangeParams{
		Column1: startOfDay,
		Column2: endOfDay,
		Limit:   10, // Top 10 items
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to fetch top selling items: %v", err)
	}

	// Convert to the expected format
	topItems := make([]map[string]interface{}, 0)
	for _, item := range topSellingItems {
		totalRevenue, err := decimal.NewFromString(item.TotalRevenue)
		if err != nil {
			continue // Skip invalid entries
		}
		topItems = append(topItems, map[string]interface{}{
			"menu_item_name":      item.MenuItemName,
			"total_quantity_sold": item.TotalQuantitySold,
			"total_revenue":       types.FromDecimal(totalRevenue),
		})
	}

	totalSales, err := decimal.NewFromString(reportData.TotalSales)
	if err != nil {
		return nil, fmt.Errorf("failed to parse total sales: %v", err)
	}

	report := map[string]interface{}{
		"date":                dateStr,
		"total_orders":        int(reportData.TotalOrders),
		"total_sales":         types.FromDecimal(totalSales),
		"average_order_value": averageOrderValue,
		"top_selling_items":   topItems,
	}

	// Cache the results for 1 hour (reports typically don't change frequently)
	cacheErr := s.cache.SetJSON(ctx, cacheKey, report, time.Hour)
	if cacheErr != nil {
		// Log the error but don't fail the request
		fmt.Printf("Warning: Failed to cache daily sales report: %v\n", cacheErr)
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

	// Calculate financial summary from orders
	summary, err := s.queries.GetFinancialSummaryByDateRange(context.Background(), db.GetFinancialSummaryByDateRangeParams{
		Column1: startDate,
		Column2: endOfDay,
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to fetch financial summary: %v", err)
	}

	// Convert total sales and other values to decimal
	totalSales, err := decimal.NewFromString(summary.TotalSales)
	if err != nil {
		return nil, fmt.Errorf("failed to parse total sales: %v", err)
	}

	// Calculate expenses from expense repository
	expenses, err := s.expenseRepo.GetExpensesByDateRange(startDate, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expenses: %v", err)
	}

	var totalExpenses decimal.Decimal
	expensesList := make([]map[string]interface{}, 0)
	for _, expense := range expenses {
		amountDecimal := decimal.RequireFromString(expense.Amount.String())
		totalExpenses = totalExpenses.Add(amountDecimal)

		expensesList = append(expensesList, map[string]interface{}{
			"id":          expense.ID,
			"category":    expense.Category,
			"description": expense.Description,
			"amount":      expense.Amount,
			"date":        expense.Date,
			"created_at":  expense.CreatedAt,
		})
	}

	// Calculate profit
	totalProfit :=  totalSales.Sub(totalExpenses)

	// Get sales by  category breakdown
	salesByCategory, err := s.queries.GetSalesByCategoryByDateRange(context.Background(), db.GetSalesByCategoryByDateRangeParams{
		Column1: startDate,
		Column2: endOfDay,
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to fetch sales by category: %v", err)
	}

	salesByCategoryList := make([]map[string]interface{}, 0)
	for _, category := range salesByCategory {
		totalRevenue, err := decimal.NewFromString(category.TotalRevenue)
		if err != nil {
			continue // Skip invalid entries
		}

		salesByCategoryList = append(salesByCategoryList, map[string]interface{}{
			"category_name":   category.CategoryName,
			"items_sold":      int(category.ItemsSold),
			"total_quantity":  int(category.TotalQuantity),
			"total_revenue":   types.FromDecimal(totalRevenue),
		})
	}

	report := map[string]interface{}{
		"period": map[string]string{
			"start_date": startDateStr,
			"end_date":   endDateStr,
		},
		"total_sales":       types.FromDecimal(totalSales),
		"total_expenses":    types.FromDecimal(totalExpenses),
		"total_profit":      types.FromDecimal(totalProfit),
		"sales_by_category": salesByCategoryList,
		"expenses":          expensesList,
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

	// Calculate end of the end date (23:59:59)
	endOfDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())

	// Get sales by category breakdown
	salesByCategory, err := s.queries.GetSalesByCategoryByDateRange(context.Background(), db.GetSalesByCategoryByDateRangeParams{
		Column1: startDate,
		Column2: endOfDay,
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to fetch sales by category: %v", err)
	}

	salesByCategoryList := make([]map[string]interface{}, 0)
	for _, category := range salesByCategory {
		totalRevenue, err := decimal.NewFromString(category.TotalRevenue)
		if err != nil {
			continue // Skip invalid entries
		}

		salesByCategoryList = append(salesByCategoryList, map[string]interface{}{
			"category_name":   category.CategoryName,
			"items_sold":      int(category.ItemsSold),
			"total_quantity":  int(category.TotalQuantity),
			"total_revenue":   types.FromDecimal(totalRevenue),
		})
	}

	report := map[string]interface{}{
		"period": map[string]string{
			"start_date": startDateStr,
			"end_date":   endDateStr,
		},
		"sales_by_category": salesByCategoryList,
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

	// Create cache key
	cacheKey := fmt.Sprintf("top_selling_items:from:%s:to:%s:limit:%d", startDateStr, endDateStr, limit)

	// Try to get from cache first
	var cachedReport map[string]interface{}
	ctx := context.Background()
	err = s.cache.GetJSON(ctx, cacheKey, &cachedReport)
	if err == nil {
		// Cache hit - return cached data
		return &types.APIResponse{
			Success: true,
			Data:    cachedReport,
		}, nil
	}

	// Calculate end of the end date (23:59:59)
	endOfDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())

	// Get top selling items by date range
	topSellingItemsData, err := s.queries.GetTopSellingItemsByDateRange(context.Background(), db.GetTopSellingItemsByDateRangeParams{
		Column1: startDate,
		Column2: endOfDay,
		Limit:   int32(limit),
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to fetch top selling items: %v", err)
	}

	topSellingItems := make([]map[string]interface{}, 0)
	for _, item := range topSellingItemsData {
		totalRevenue, err := decimal.NewFromString(item.TotalRevenue)
		if err != nil {
			continue // Skip invalid entries
		}

		topSellingItems = append(topSellingItems, map[string]interface{}{
			"menu_item_name":      item.MenuItemName,
			"total_quantity_sold": int(item.TotalQuantitySold),
			"total_revenue":       types.FromDecimal(totalRevenue),
		})
	}

	report := map[string]interface{}{
		"period": map[string]string{
			"start_date": startDateStr,
			"end_date":   endDateStr,
		},
		"top_selling_items": topSellingItems,
		"limit":             limit,
	}

	// Cache the results for 30 minutes
	cacheErr := s.cache.SetJSON(ctx, cacheKey, report, 30*time.Minute)
	if cacheErr != nil {
		// Log the error but don't fail the request
		fmt.Printf("Warning: Failed to cache top selling items report: %v\n", cacheErr)
	}

	return &types.APIResponse{
		Success: true,
		Data:    report,
	}, nil
}
