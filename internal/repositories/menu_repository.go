package repositories

import (
	"errors"
	"github.com/AndikaPrasetia/pos-cafee/internal/db"
	"github.com/AndikaPrasetia/pos-cafee/internal/models"
)

// menuRepo implements the MenuRepo interface
type menuRepo struct {
	queries *db.Queries
}

// GetCategory retrieves a category by ID
func (r *menuRepo) GetCategory(id string) (*models.Category, error) {
	// TODO: Implement this method
	return nil, NotImplementedError()
}

// ListCategories retrieves a list of categories
func (r *menuRepo) ListCategories(isActive bool, limit, offset int) ([]*models.Category, error) {
	// TODO: Implement this method
	return nil, NotImplementedError()
}

// CreateCategory creates a new category
func (r *menuRepo) CreateCategory(category *models.Category) (*models.Category, error) {
	// TODO: Implement this method
	return nil, NotImplementedError()
}

// UpdateCategory updates an existing category
func (r *menuRepo) UpdateCategory(category *models.Category) (*models.Category, error) {
	// TODO: Implement this method
	return nil, NotImplementedError()
}

// DeleteCategory deletes a category by ID
func (r *menuRepo) DeleteCategory(id string) error {
	// TODO: Implement this method
	return NotImplementedError()
}

// GetMenuItem retrieves a menu item by ID
func (r *menuRepo) GetMenuItem(id string) (*models.MenuItem, error) {
	// TODO: Implement this method
	return nil, NotImplementedError()
}

// ListMenuItems retrieves a list of menu items
func (r *menuRepo) ListMenuItems(isAvailable bool, limit, offset int) ([]*models.MenuItem, error) {
	// TODO: Implement this method
	return nil, NotImplementedError()
}

// ListMenuItemsByCategory retrieves a list of menu items by category
func (r *menuRepo) ListMenuItemsByCategory(categoryID string, limit, offset int) ([]*models.MenuItem, error) {
	// TODO: Implement this method
	return nil, NotImplementedError()
}

// CreateMenuItem creates a new menu item
func (r *menuRepo) CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error) {
	// TODO: Implement this method
	return nil, NotImplementedError()
}

// UpdateMenuItem updates an existing menu item
func (r *menuRepo) UpdateMenuItem(item *models.MenuItem) (*models.MenuItem, error) {
	// TODO: Implement this method
	return nil, NotImplementedError()
}

// DeleteMenuItem deletes a menu item by ID
func (r *menuRepo) DeleteMenuItem(id string) error {
	// TODO: Implement this method
	return NotImplementedError()
}

// NotImplementedError returns a standard error for unimplemented methods
func NotImplementedError() error {
	return errors.New("method not implemented")
}