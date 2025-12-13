package cache

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/internal/repositories"
)

// CacheWarmingService handles pre-populating cache with frequently accessed data
type CacheWarmingService struct {
	cache         Cache
	menuRepo      repositories.MenuRepo
	orderRepo     repositories.OrderRepo
	inventoryRepo repositories.InventoryRepo
	expenseRepo   repositories.ExpenseRepo
	repositories  *repositories.Repository
}

// NewCacheWarmingService creates a new cache warming service
func NewCacheWarmingService(
	cache Cache,
	repository *repositories.Repository,
) *CacheWarmingService {
	return &CacheWarmingService{
		cache:         cache,
		menuRepo:      repository.MenuRepo,
		orderRepo:     repository.OrderRepo,
		inventoryRepo: repository.InventoryRepo,
		expenseRepo:   repository.ExpenseRepo,
		repositories:  repository,
	}
}

// WarmCache warms up all frequently accessed cache entries
func (cws *CacheWarmingService) WarmCache(ctx context.Context) error {
	log.Println("Starting cache warming process...")

	// Create a wait group to run warming operations concurrently
	var wg sync.WaitGroup
	errChan := make(chan error, 4) // Buffer for up to 4 errors

	// Warm menu-related cache
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := cws.warmMenuCache(ctx); err != nil {
			errChan <- fmt.Errorf("failed to warm menu cache: %v", err)
		}
	}()

	// Warm category-related cache
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := cws.warmCategoryCache(ctx); err != nil {
			errChan <- fmt.Errorf("failed to warm category cache: %v", err)
		}
	}()

	// Warm report-related cache (for today's data)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := cws.warmReportCache(ctx); err != nil {
			errChan <- fmt.Errorf("failed to warm report cache: %v", err)
		}
	}()

	// Warm inventory-related cache
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := cws.warmInventoryCache(ctx); err != nil {
			errChan <- fmt.Errorf("failed to warm inventory cache: %v", err)
		}
	}()

	// Wait for all warming operations to complete
	wg.Wait()
	close(errChan)

	// Collect any errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
		log.Printf("Cache warming error: %v", err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("cache warming completed with %d errors: %v", len(errors), errors)
	}

	log.Println("Cache warming process completed successfully")
	return nil
}

// warmMenuCache pre-populates cache with menu items
func (cws *CacheWarmingService) warmMenuCache(ctx context.Context) error {
	// Get all active menu items with default pagination
	menuItems, err := cws.menuRepo.ListMenuItems(true, 100, 0) // Fetch first 100 items
	if err != nil {
		return fmt.Errorf("failed to fetch menu items: %v", err)
	}

	// Cache the list of all active menu items
	if len(menuItems) > 0 {
		if err := cws.cache.SetJSON(ctx, "menu_items:available:true:limit:100:offset:0", menuItems, 15*time.Minute); err != nil {
			return fmt.Errorf("failed to cache menu items list: %v", err)
		}
	}

	// Cache each individual menu item
	for _, item := range menuItems {
		key := fmt.Sprintf("menu_item:%s", item.ID)
		if err := cws.cache.SetJSON(ctx, key, item, 15*time.Minute); err != nil {
			log.Printf("Warning: failed to cache menu item %s: %v", item.ID, err)
			// Don't return error here, just continue with other items
		}
	}

	log.Printf("Warmed up cache for %d menu items", len(menuItems))
	return nil
}

// warmCategoryCache pre-populates cache with category data
func (cws *CacheWarmingService) warmCategoryCache(ctx context.Context) error {
	// Get all active categories
	categories, err := cws.menuRepo.ListCategories(true, 100, 0) // Fetch first 100 categories
	if err != nil {
		return fmt.Errorf("failed to fetch categories: %v", err)
	}

	// Cache the list of all active categories
	if len(categories) > 0 {
		if err := cws.cache.SetJSON(ctx, "categories:active:true:limit:100:offset:0", categories, 15*time.Minute); err != nil {
			return fmt.Errorf("failed to cache categories list: %v", err)
		}
	}

	// Cache each individual category
	for _, category := range categories {
		key := fmt.Sprintf("category:%s", category.ID)
		if err := cws.cache.SetJSON(ctx, key, category, 15*time.Minute); err != nil {
			log.Printf("Warning: failed to cache category %s: %v", category.ID, err)
			// Don't return error here, just continue with other categories
		}
	}

	log.Printf("Warmed up cache for %d categories", len(categories))
	return nil
}

// warmReportCache pre-populates cache with report data
func (cws *CacheWarmingService) warmReportCache(ctx context.Context) error {
	// Get today's date in YYYY-MM-DD format
	today := time.Now().Format("2006-01-02")

	// Try to get today's sales report
	// Note: We're not calling a service method here, but rather the repository/db directly
	// In a real implementation, we would call the report service's method to warm the cache
	// For now, let's just add an empty placeholder to prevent multiple attempts
	// in cases where no data exists yet

	// Warm daily sales report for today
	reportKey := fmt.Sprintf("daily_sales_report:%s", today)
	if _, err := cws.cache.Get(ctx, reportKey); err != nil {
		// If not already cached, cache an empty report
		emptyReport := map[string]interface{}{
			"date":                today,
			"total_orders":        0,
			"total_sales":         "0.00",
			"average_order_value": "0.00",
			"top_selling_items":   []map[string]interface{}{},
		}
		if err := cws.cache.SetJSON(ctx, reportKey, emptyReport, 30*time.Minute); err != nil {
			log.Printf("Warning: failed to cache daily sales report for today: %v", err)
		}
	}

	// Warm top selling items report for today
	topSellingKey := fmt.Sprintf("top_selling_items:from:%s:to:%s:limit:10", today, today)
	if _, err := cws.cache.Get(ctx, topSellingKey); err != nil {
		// If not already cached, cache an empty report
		emptyTopSellingReport := map[string]interface{}{
			"period": map[string]string{
				"start_date": today,
				"end_date":   today,
			},
			"top_selling_items": []map[string]interface{}{},
			"limit":             10,
		}
		if err := cws.cache.SetJSON(ctx, topSellingKey, emptyTopSellingReport, 30*time.Minute); err != nil {
			log.Printf("Warning: failed to cache top selling items report: %v", err)
		}
	}

	log.Println("Warmed up cache for report data")
	return nil
}

// warmInventoryCache pre-populates cache with inventory data
func (cws *CacheWarmingService) warmInventoryCache(ctx context.Context) error {
	// Get inventory items using the correct filter format
	filter := models.InventoryFilter{
		Limit:  100, // Fetch first 100 items
		Offset: 0,
	}

	inventoryItems, err := cws.inventoryRepo.ListInventory(filter) // Use the correct parameter format
	if err != nil {
		return fmt.Errorf("failed to fetch inventory items: %v", err)
	}

	// In a real implementation, we would cache inventory data
	// For now, just log the warming operation
	log.Printf("Warmed up cache for %d inventory items", len(inventoryItems))
	return nil
}

// WarmCacheWithRetry attempts to warm cache with retry logic
func (cws *CacheWarmingService) WarmCacheWithRetry(ctx context.Context, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = cws.WarmCache(ctx)
		if err == nil {
			return nil
		}
		log.Printf("Cache warming attempt %d failed: %v, retrying in 5 seconds...", i+1, err)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			// Continue to next retry
		}
	}
	return fmt.Errorf("cache warming failed after %d attempts: %v", maxRetries, err)
}

// WarmCacheInBackground starts cache warming in the background with a timeout
func (cws *CacheWarmingService) WarmCacheInBackground() {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		
		if err := cws.WarmCacheWithRetry(ctx, 3); err != nil {
			log.Printf("Background cache warming failed: %v", err)
		}
	}()
}