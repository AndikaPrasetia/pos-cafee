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
)

// MenuService handles menu-related business logic
type MenuService struct {
	menuRepo      repositories.MenuRepo
	inventoryRepo repositories.InventoryRepo
	cache         cache.Cache
}

// NewMenuService creates a new menu service
func NewMenuService(menuRepo repositories.MenuRepo, inventoryRepo repositories.InventoryRepo, cache cache.Cache) *MenuService {
	return &MenuService{
		menuRepo:      menuRepo,
		inventoryRepo: inventoryRepo,
		cache:         cache,
	}
}

// CreateCategory creates a new menu category
func (s *MenuService) CreateCategory(categoryData *models.CategoryCreate) (*types.APIResponse, error) {
	category := &models.Category{
		ID:          uuid.New().String(),
		Name:        categoryData.Name,
		Description: categoryData.Description,
		IsActive:    true, // New categories are active by default
	}

	createdCategory, err := s.menuRepo.CreateCategory(category)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %v", err)
	}

	// Invalidate cached category lists to ensure the new category appears immediately
	ctx := context.Background()
	// Delete all cached ListCategories results (all combinations of active/limit/offset)
	categoryListKeys, err := s.cache.Keys(ctx, "categories:*")
	if err == nil {
		for _, key := range categoryListKeys {
			s.cache.Delete(ctx, key)
		}
	} else {
		fmt.Printf("Warning: Failed to get category list cache keys: %v\n", err)
	}

	return &types.APIResponse{
		Success: true,
		Data:    createdCategory,
	}, nil
}

// GetCategory retrieves a category by ID
func (s *MenuService) GetCategory(id string) (*types.APIResponse, error) {
	// Create cache key
	cacheKey := fmt.Sprintf("category:%s", id)

	// Try to get from cache first
	var category *models.Category
	ctx := context.Background()
	err := s.cache.GetJSON(ctx, cacheKey, &category)
	if err == nil {
		// Cache hit - return cached data
		return &types.APIResponse{
			Success: true,
			Data:    category,
		}, nil
	}

	// Cache miss - get from database
	category, err = s.menuRepo.GetCategory(id)
	if err != nil {
		return nil, err
	}

	// Cache the result for 15 minutes
	cacheErr := s.cache.SetJSON(ctx, cacheKey, category, 15*time.Minute)
	if cacheErr != nil {
		// Log the error but don't fail the request
		fmt.Printf("Warning: Failed to cache category: %v\n", cacheErr)
	}

	return &types.APIResponse{
		Success: true,
		Data:    category,
	}, nil
}

// ListCategories retrieves a list of categories
func (s *MenuService) ListCategories(isActive bool, limit, offset int) (*types.APIResponse, error) {
	// Create cache key
	cacheKey := fmt.Sprintf("categories:active:%t:limit:%d:offset:%d", isActive, limit, offset)

	// Try to get from cache first
	var categories []*models.Category
	ctx := context.Background()
	err := s.cache.GetJSON(ctx, cacheKey, &categories)
	if err == nil {
		// Cache hit - return cached data
		return &types.APIResponse{
			Success: true,
			Data:    categories,
		}, nil
	}

	// Cache miss - get from database
	categories, err = s.menuRepo.ListCategories(isActive, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %v", err)
	}

	// Cache the results for 15 minutes
	cacheErr := s.cache.SetJSON(ctx, cacheKey, categories, 15*time.Minute)
	if cacheErr != nil {
		// Log the error but don't fail the request
		fmt.Printf("Warning: Failed to cache categories: %v\n", cacheErr)
	}

	return &types.APIResponse{
		Success: true,
		Data:    categories,
	}, nil
}

// UpdateCategory updates an existing category
func (s *MenuService) UpdateCategory(id string, updateData *models.CategoryUpdate) (*types.APIResponse, error) {
	category, err := s.menuRepo.GetCategory(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Update fields if provided in updateData
	if updateData.Name != nil {
		category.Name = *updateData.Name
	}
	if updateData.Description != nil {
		category.Description = updateData.Description
	}
	if updateData.IsActive != nil {
		category.IsActive = *updateData.IsActive
	}

	updatedCategory, err := s.menuRepo.UpdateCategory(category)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %v", err)
	}

	// Invalidate cached category lists to ensure consistency
	ctx := context.Background()

	// Delete cached individual category
	s.cache.Delete(ctx, fmt.Sprintf("category:%s", id))

	// Delete all cached ListCategories results (all combinations of active/limit/offset)
	categoryListKeys, err := s.cache.Keys(ctx, "categories:*")
	if err == nil {
		for _, key := range categoryListKeys {
			s.cache.Delete(ctx, key)
		}
	} else {
		fmt.Printf("Warning: Failed to get category list cache keys: %v\n", err)
	}

	// If name changed, also invalidate menu items by category
	if updateData.Name != nil {
		menuItemsByCategoryKeys, err := s.cache.Keys(ctx, fmt.Sprintf("menu_items:category:%s:*", id))
		if err == nil {
			for _, key := range menuItemsByCategoryKeys {
				s.cache.Delete(ctx, key)
			}
		} else {
			fmt.Printf("Warning: Failed to get menu items by category cache keys: %v\n", err)
		}
	}

	return &types.APIResponse{
		Success: true,
		Data:    updatedCategory,
	}, nil
}

// DeleteCategory deletes (deactivates) a category
func (s *MenuService) DeleteCategory(id string) (*types.APIResponse, error) {
	err := s.menuRepo.DeleteCategory(id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete category: %v", err)
	}

	// Invalidate cached category lists to ensure consistency
	ctx := context.Background()

	// Delete cached individual category
	s.cache.Delete(ctx, fmt.Sprintf("category:%s", id))

	// Delete all cached ListCategories results (all combinations of active/limit/offset)
	categoryListKeys, err := s.cache.Keys(ctx, "categories:*")
	if err == nil {
		for _, key := range categoryListKeys {
			s.cache.Delete(ctx, key)
		}
	} else {
		fmt.Printf("Warning: Failed to get category list cache keys: %v\n", err)
	}

	// Also invalidate menu items by this category
	menuItemsByCategoryKeys, err := s.cache.Keys(ctx, fmt.Sprintf("menu_items:category:%s:*", id))
	if err == nil {
		for _, key := range menuItemsByCategoryKeys {
			s.cache.Delete(ctx, key)
		}
	} else {
		fmt.Printf("Warning: Failed to get menu items by category cache keys: %v\n", err)
	}

	return &types.APIResponse{
		Success: true,
		Message: "Category deleted successfully",
	}, nil
}

// CreateMenuItem creates a new menu item
func (s *MenuService) CreateMenuItem(itemData *models.MenuItemCreate) (*types.APIResponse, error) {
	// Validate category ID
	_, err := uuid.Parse(itemData.CategoryID)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}

	item := &models.MenuItem{
		ID:          uuid.New().String(),
		Name:        itemData.Name,
		CategoryID:  itemData.CategoryID,
		Description: itemData.Description,
		Price:       itemData.Price,
		Cost:        itemData.Cost,
		IsAvailable: itemData.IsAvailable,
	}

	createdItem, err := s.menuRepo.CreateMenuItem(item)
	if err != nil {
		return nil, fmt.Errorf("failed to create menu item: %v", err)
	}

	// Create an inventory record for the new menu item
	err = s.inventoryRepo.CreateInventoryRecord(createdItem.ID)
	if err != nil {
		// This is not a critical failure, but should be logged
		fmt.Printf("Warning: Failed to create inventory record for item %s: %v\n", createdItem.ID, err)
	}

	// Invalidate cached menu item lists to ensure the new item appears immediately
	ctx := context.Background()

	// Delete all cached ListMenuItems results (all combinations of available/limit/offset)
	menuItemListKeys, err := s.cache.Keys(ctx, "menu_items:*")
	if err == nil {
		for _, key := range menuItemListKeys {
			s.cache.Delete(ctx, key)
		}
	} else {
		fmt.Printf("Warning: Failed to get menu item list cache keys: %v\n", err)
	}

	// Also invalidate the specific category's cached list
	menuItemsByCategoryKeys, err := s.cache.Keys(ctx, fmt.Sprintf("menu_items:category:%s:*", itemData.CategoryID))
	if err == nil {
		for _, key := range menuItemsByCategoryKeys {
			s.cache.Delete(ctx, key)
		}
	} else {
		fmt.Printf("Warning: Failed to get menu items by category cache keys: %v\n", err)
	}

	return &types.APIResponse{
		Success: true,
		Data:    createdItem,
	}, nil
}

// GetMenuItem retrieves a menu item by ID
func (s *MenuService) GetMenuItem(id string) (*types.APIResponse, error) {
	// Create cache key
	cacheKey := fmt.Sprintf("menu_item:%s", id)

	// Try to get from cache first
	var item *models.MenuItem
	ctx := context.Background()
	err := s.cache.GetJSON(ctx, cacheKey, &item)
	if err == nil {
		// Cache hit - return cached data
		return &types.APIResponse{
			Success: true,
			Data:    item,
		}, nil
	}

	// Cache miss - get from database
	item, err = s.menuRepo.GetMenuItem(id)
	if err != nil {
		return nil, err
	}

	// Cache the result for 15 minutes
	cacheErr := s.cache.SetJSON(ctx, cacheKey, item, 15*time.Minute)
	if cacheErr != nil {
		// Log the error but don't fail the request
		fmt.Printf("Warning: Failed to cache menu item: %v\n", cacheErr)
	}

	return &types.APIResponse{
		Success: true,
		Data:    item,
	}, nil
}

// ListMenuItems retrieves a list of menu items
func (s *MenuService) ListMenuItems(isAvailable bool, limit, offset int) (*types.APIResponse, error) {
	// Create cache key
	cacheKey := fmt.Sprintf("menu_items:available:%t:limit:%d:offset:%d", isAvailable, limit, offset)

	// Try to get from cache first
	var items []*models.MenuItem
	ctx := context.Background()
	err := s.cache.GetJSON(ctx, cacheKey, &items)
	if err == nil {
		// Cache hit - return cached data
		return &types.APIResponse{
			Success: true,
			Data:    items,
		}, nil
	}

	// Cache miss - get from database
	items, err = s.menuRepo.ListMenuItems(isAvailable, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list menu items: %v", err)
	}

	// Cache the results for 15 minutes
	cacheErr := s.cache.SetJSON(ctx, cacheKey, items, 15*time.Minute)
	if cacheErr != nil {
		// Log the error but don't fail the request
		fmt.Printf("Warning: Failed to cache menu items: %v\n", cacheErr)
	}

	return &types.APIResponse{
		Success: true,
		Data:    items,
	}, nil
}

// ListMenuItemsByCategory retrieves a list of menu items in a specific category
func (s *MenuService) ListMenuItemsByCategory(categoryID string, limit, offset int) (*types.APIResponse, error) {
	_, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}

	// Create cache key
	cacheKey := fmt.Sprintf("menu_items:category:%s:limit:%d:offset:%d", categoryID, limit, offset)

	// Try to get from cache first
	var items []*models.MenuItem
	ctx := context.Background()
	err = s.cache.GetJSON(ctx, cacheKey, &items)
	if err == nil {
		// Cache hit - return cached data
		return &types.APIResponse{
			Success: true,
			Data:    items,
		}, nil
	}

	// Cache miss - get from database
	items, err = s.menuRepo.ListMenuItemsByCategory(categoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list menu items by category: %v", err)
	}

	// Cache the results for 15 minutes
	cacheErr := s.cache.SetJSON(ctx, cacheKey, items, 15*time.Minute)
	if cacheErr != nil {
		// Log the error but don't fail the request
		fmt.Printf("Warning: Failed to cache menu items by category: %v\n", cacheErr)
	}

	return &types.APIResponse{
		Success: true,
		Data:    items,
	}, nil
}

// UpdateMenuItem updates an existing menu item
func (s *MenuService) UpdateMenuItem(id string, updateData *models.MenuItemUpdate) (*types.APIResponse, error) {
	item, err := s.menuRepo.GetMenuItem(id)
	if err != nil {
		return nil, errors.New("menu item not found")
	}

	// Store original category ID to invalidate old category cache if needed
	originalCategoryID := item.CategoryID

	// Update fields if provided in updateData
	if updateData.Name != nil {
		item.Name = *updateData.Name
	}
	if updateData.CategoryID != nil {
		_, err := uuid.Parse(*updateData.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category ID")
		}
		item.CategoryID = *updateData.CategoryID
	}
	if updateData.Description != nil {
		item.Description = updateData.Description
	}
	if updateData.Price != nil {
		item.Price = *updateData.Price
	}
	if updateData.Cost != nil {
		item.Cost = *updateData.Cost
	}
	if updateData.IsAvailable != nil {
		item.IsAvailable = *updateData.IsAvailable
	}

	updatedItem, err := s.menuRepo.UpdateMenuItem(item)
	if err != nil {
		return nil, fmt.Errorf("failed to update menu item: %v", err)
	}

	// Invalidate cached menu item lists to ensure consistency
	ctx := context.Background()

	// Delete cached individual menu item
	s.cache.Delete(ctx, fmt.Sprintf("menu_item:%s", id))

	// Delete all cached ListMenuItems results (all combinations of available/limit/offset)
	menuItemListKeys, err := s.cache.Keys(ctx, "menu_items:*")
	if err == nil {
		for _, key := range menuItemListKeys {
			s.cache.Delete(ctx, key)
		}
	} else {
		fmt.Printf("Warning: Failed to get menu item list cache keys: %v\n", err)
	}

	// Invalidate the old category's cached results if category changed
	if updateData.CategoryID != nil && *updateData.CategoryID != originalCategoryID {
		menuItemsByOldCategoryKeys, err := s.cache.Keys(ctx, fmt.Sprintf("menu_items:category:%s:*", originalCategoryID))
		if err == nil {
			for _, key := range menuItemsByOldCategoryKeys {
				s.cache.Delete(ctx, key)
			}
		} else {
			fmt.Printf("Warning: Failed to get menu items by old category cache keys: %v\n", err)
		}
	}

	// Invalidate the new category's cached results if category changed or same category
	newCategoryID := originalCategoryID
	if updateData.CategoryID != nil {
		newCategoryID = *updateData.CategoryID
	}
	menuItemsByNewCategoryKeys, err := s.cache.Keys(ctx, fmt.Sprintf("menu_items:category:%s:*", newCategoryID))
	if err == nil {
		for _, key := range menuItemsByNewCategoryKeys {
			s.cache.Delete(ctx, key)
		}
	} else {
		fmt.Printf("Warning: Failed to get menu items by new category cache keys: %v\n", err)
	}

	return &types.APIResponse{
		Success: true,
		Data:    updatedItem,
	}, nil
}

// DeleteMenuItem deletes (deactivates) a menu item
func (s *MenuService) DeleteMenuItem(id string) (*types.APIResponse, error) {
	// First get the menu item to access its category ID for cache invalidation
	item, err := s.menuRepo.GetMenuItem(id)
	if err != nil {
		return nil, fmt.Errorf("menu item not found: %v", err)
	}

	err = s.menuRepo.DeleteMenuItem(id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete menu item: %v", err)
	}

	// Invalidate cached menu item lists to ensure consistency
	ctx := context.Background()

	// Delete cached individual menu item
	s.cache.Delete(ctx, fmt.Sprintf("menu_item:%s", id))

	// Delete all cached ListMenuItems results (all combinations of available/limit/offset)
	menuItemListKeys, err := s.cache.Keys(ctx, "menu_items:*")
	if err == nil {
		for _, key := range menuItemListKeys {
			s.cache.Delete(ctx, key)
		}
	} else {
		fmt.Printf("Warning: Failed to get menu item list cache keys: %v\n", err)
	}

	// Invalidate the category's cached results that this item belonged to
	menuItemsByCategoryKeys, err := s.cache.Keys(ctx, fmt.Sprintf("menu_items:category:%s:*", item.CategoryID))
	if err == nil {
		for _, key := range menuItemsByCategoryKeys {
			s.cache.Delete(ctx, key)
		}
	} else {
		fmt.Printf("Warning: Failed to get menu items by category cache keys: %v\n", err)
	}

	return &types.APIResponse{
		Success: true,
		Message: "Menu item deleted successfully",
	}, nil
}