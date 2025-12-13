package cache

import (
	"context"
	"fmt"
	"log"
	"strings"
)

// CacheInvalidation provides methods for cache invalidation
type CacheInvalidation struct {
	cache Cache
}

// NewCacheInvalidation creates a new cache invalidation helper
func NewCacheInvalidation(cache Cache) *CacheInvalidation {
	return &CacheInvalidation{
		cache: cache,
	}
}

// InvalidateByPrefix invalidates all cache keys that start with the given prefix
func (ci *CacheInvalidation) InvalidateByPrefix(ctx context.Context, prefix string) error {
	keys, err := ci.cache.Keys(ctx, prefix+"*")
	if err != nil {
		return fmt.Errorf("failed to get keys by prefix %s: %v", prefix, err)
	}

	for _, key := range keys {
		if err := ci.cache.Delete(ctx, key); err != nil {
			log.Printf("Warning: Failed to delete cache key %s: %v", key, err)
		}
	}

	if len(keys) > 0 {
		log.Printf("Invalidated %d cache keys with prefix %s", len(keys), prefix)
	}

	return nil
}

// InvalidateRelatedKeys invalidates cache keys related to a specific entity
func (ci *CacheInvalidation) InvalidateRelatedKeys(ctx context.Context, entityName, entityID string) {
	// Common patterns for entity-related keys
	patterns := []string{
		fmt.Sprintf("%s:%s", entityName, entityID),
		fmt.Sprintf("%s:%s:*", entityName, entityID),
		fmt.Sprintf("%s:*:%s", entityName, entityID),
		fmt.Sprintf("%s:*", entityName),
	}

	for _, pattern := range patterns {
		if err := ci.InvalidateByPrefix(ctx, pattern); err != nil {
			log.Printf("Warning: Failed to invalidate keys for pattern %s: %v", pattern, err)
		}
	}
}

// InvalidateListKeys invalidates cache keys related to list operations
func (ci *CacheInvalidation) InvalidateListKeys(ctx context.Context, entityName string) error {
	return ci.InvalidateByPrefix(ctx, fmt.Sprintf("%s:*", entityName))
}

// InvalidateMenuRelatedKeys invalidates all cache keys related to menu operations
func (ci *CacheInvalidation) InvalidateMenuRelatedKeys(ctx context.Context) {
	// Invalidate all menu items cache
	if err := ci.InvalidateListKeys(ctx, "menu_item"); err != nil {
		log.Printf("Warning: Failed to invalidate menu item cache: %v", err)
	}

	// Invalidate all category cache
	if err := ci.InvalidateListKeys(ctx, "category"); err != nil {
		log.Printf("Warning: Failed to invalidate category cache: %v", err)
	}
}

// InvalidateReportRelatedKeys invalidates all cache keys related to report operations
func (ci *CacheInvalidation) InvalidateReportRelatedKeys(ctx context.Context) {
	// Invalidate daily sales reports
	if err := ci.InvalidateByPrefix(ctx, "daily_sales_report:"); err != nil {
		log.Printf("Warning: Failed to invalidate daily sales reports cache: %v", err)
	}

	// Invalidate top selling items reports
	if err := ci.InvalidateByPrefix(ctx, "top_selling_items:"); err != nil {
		log.Printf("Warning: Failed to invalidate top selling items reports cache: %v", err)
	}
}

// InvalidateUserRelatedKeys invalidates all cache keys related to user operations
func (ci *CacheInvalidation) InvalidateUserRelatedKeys(ctx context.Context) {
	// Invalidate user-related cache
	if err := ci.InvalidateByPrefix(ctx, "user:"); err != nil {
		log.Printf("Warning: Failed to invalidate user cache: %v", err)
	}

	// Invalidate user list cache
	if err := ci.InvalidateByPrefix(ctx, "users:"); err != nil {
		log.Printf("Warning: Failed to invalidate user list cache: %v", err)
	}
}

// InvalidateAllCache clears all cache entries
func (ci *CacheInvalidation) InvalidateAllCache(ctx context.Context) error {
	return ci.cache.FlushDB(ctx)
}

// InvalidateByPattern invalidates cache keys matching a specific pattern
func (ci *CacheInvalidation) InvalidateByPattern(ctx context.Context, pattern string) error {
	return ci.InvalidateByPrefix(ctx, strings.TrimSuffix(pattern, "*"))
}