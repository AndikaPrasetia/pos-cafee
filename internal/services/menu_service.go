package services

import (
	"errors"
	"fmt"

	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/internal/repositories"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/google/uuid"
)

// MenuService handles menu-related business logic
type MenuService struct {
	menuRepo      repositories.MenuRepo
	inventoryRepo repositories.InventoryRepo
}

// NewMenuService creates a new menu service
func NewMenuService(menuRepo repositories.MenuRepo, inventoryRepo repositories.InventoryRepo) *MenuService {
	return &MenuService{
		menuRepo:      menuRepo,
		inventoryRepo: inventoryRepo,
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

	return &types.APIResponse{
		Success: true,
		Data:    createdCategory,
	}, nil
}

// GetCategory retrieves a category by ID
func (s *MenuService) GetCategory(id string) (*types.APIResponse, error) {
	category, err := s.menuRepo.GetCategory(id)
	if err != nil {
		return nil, err
	}

	return &types.APIResponse{
		Success: true,
		Data:    category,
	}, nil
}

// ListCategories retrieves a list of categories
func (s *MenuService) ListCategories(isActive bool, limit, offset int) (*types.APIResponse, error) {
	categories, err := s.menuRepo.ListCategories(isActive, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %v", err)
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

	return &types.APIResponse{
		Success: true,
		Data:    createdItem,
	}, nil
}

// GetMenuItem retrieves a menu item by ID
func (s *MenuService) GetMenuItem(id string) (*types.APIResponse, error) {
	item, err := s.menuRepo.GetMenuItem(id)
	if err != nil {
		return nil, err
	}

	return &types.APIResponse{
		Success: true,
		Data:    item,
	}, nil
}

// ListMenuItems retrieves a list of menu items
func (s *MenuService) ListMenuItems(isAvailable bool, limit, offset int) (*types.APIResponse, error) {
	items, err := s.menuRepo.ListMenuItems(isAvailable, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list menu items: %v", err)
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

	items, err := s.menuRepo.ListMenuItemsByCategory(categoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list menu items by category: %v", err)
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

	return &types.APIResponse{
		Success: true,
		Data:    updatedItem,
	}, nil
}

// DeleteMenuItem deletes (deactivates) a menu item
func (s *MenuService) DeleteMenuItem(id string) (*types.APIResponse, error) {
	err := s.menuRepo.DeleteMenuItem(id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete menu item: %v", err)
	}

	return &types.APIResponse{
		Success: true,
		Message: "Menu item deleted successfully",
	}, nil
}