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

// menuRepo implements the MenuRepo interface
type menuRepo struct {
	queries *db.Queries
}

// GetCategory retrieves a category by ID
func (r *menuRepo) GetCategory(id string) (*models.Category, error) {
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	dbCategory, err := r.queries.GetCategory(context.Background(), categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	// Convert database category to model category
	category := &models.Category{
		ID:        dbCategory.ID.String(),
		Name:      dbCategory.Name,
		IsActive:  dbCategory.IsActive,
		CreatedAt: dbCategory.CreatedAt,
		UpdatedAt: dbCategory.UpdatedAt,
	}

	if dbCategory.Description.Valid {
		description := dbCategory.Description.String
		category.Description = &description
	}

	return category, nil
}

// ListCategories retrieves a list of categories
func (r *menuRepo) ListCategories(isActive bool, limit, offset int) ([]*models.Category, error) {
	dbCategories, err := r.queries.ListCategories(context.Background(), db.ListCategoriesParams{
		IsActive: isActive,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, err
	}

	var categories []*models.Category
	for _, dbCategory := range dbCategories {
		category := &models.Category{
			ID:        dbCategory.ID.String(),
			Name:      dbCategory.Name,
			IsActive:  dbCategory.IsActive,
			CreatedAt: dbCategory.CreatedAt,
			UpdatedAt: dbCategory.UpdatedAt,
		}

		if dbCategory.Description.Valid {
			description := dbCategory.Description.String
			category.Description = &description
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// CreateCategory creates a new category
func (r *menuRepo) CreateCategory(category *models.Category) (*models.Category, error) {
	var description sql.NullString
	if category.Description != nil {
		description = sql.NullString{
			String: *category.Description,
			Valid:  true,
		}
	}

	dbCategory, err := r.queries.CreateCategory(context.Background(), db.CreateCategoryParams{
		Name:        category.Name,
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	createdCategory := &models.Category{
		ID:        dbCategory.ID.String(),
		Name:      dbCategory.Name,
		IsActive:  dbCategory.IsActive,
		CreatedAt: dbCategory.CreatedAt,
		UpdatedAt: dbCategory.UpdatedAt,
	}

	if dbCategory.Description.Valid {
		description := dbCategory.Description.String
		createdCategory.Description = &description
	}

	return createdCategory, nil
}

// UpdateCategory updates an existing category
func (r *menuRepo) UpdateCategory(category *models.Category) (*models.Category, error) {
	categoryID, err := uuid.Parse(category.ID)
	if err != nil {
		return nil, err
	}

	var description sql.NullString
	if category.Description != nil {
		description = sql.NullString{
			String: *category.Description,
			Valid:  true,
		}
	} else {
		description = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	dbCategory, err := r.queries.UpdateCategory(context.Background(), db.UpdateCategoryParams{
		ID:          categoryID,
		Name:        category.Name,
		Description: description,
		IsActive:    category.IsActive,
	})
	if err != nil {
		return nil, err
	}

	updatedCategory := &models.Category{
		ID:        dbCategory.ID.String(),
		Name:      dbCategory.Name,
		IsActive:  dbCategory.IsActive,
		CreatedAt: dbCategory.CreatedAt,
		UpdatedAt: dbCategory.UpdatedAt,
	}

	if dbCategory.Description.Valid {
		description := dbCategory.Description.String
		updatedCategory.Description = &description
	}

	return updatedCategory, nil
}

// DeleteCategory deletes a category by ID
func (r *menuRepo) DeleteCategory(id string) error {
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = r.queries.DeleteCategory(context.Background(), categoryID)
	if err != nil {
		return err
	}

	return nil
}

// GetMenuItem retrieves a menu item by ID
func (r *menuRepo) GetMenuItem(id string) (*models.MenuItem, error) {
	itemID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	dbMenuItem, err := r.queries.GetMenuItem(context.Background(), itemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("menu item not found")
		}
		return nil, err
	}

	// Convert database menu item to model menu item
	price, err := decimal.NewFromString(dbMenuItem.Price)
	if err != nil {
		return nil, err
	}

	cost, err := decimal.NewFromString(dbMenuItem.Cost)
	if err != nil {
		return nil, err
	}

	menuItem := &models.MenuItem{
		ID:          dbMenuItem.ID.String(),
		Name:        dbMenuItem.Name,
		CategoryID:  dbMenuItem.CategoryID.String(),
		IsAvailable: dbMenuItem.IsAvailable,
		CreatedAt:   dbMenuItem.CreatedAt,
		UpdatedAt:   dbMenuItem.UpdatedAt,
		Price:       types.DecimalText(price),
		Cost:        types.DecimalText(cost),
	}

	if dbMenuItem.Description.Valid {
		description := dbMenuItem.Description.String
		menuItem.Description = &description
	}

	return menuItem, nil
}

// ListMenuItems retrieves a list of menu items
func (r *menuRepo) ListMenuItems(isAvailable bool, limit, offset int) ([]*models.MenuItem, error) {
	dbMenuItems, err := r.queries.ListMenuItems(context.Background(), db.ListMenuItemsParams{
		IsAvailable: isAvailable,
		Limit:       int32(limit),
		Offset:      int32(offset),
	})
	if err != nil {
		return nil, err
	}

	var menuItems []*models.MenuItem
	for _, dbMenuItem := range dbMenuItems {
		price, err := decimal.NewFromString(dbMenuItem.Price)
		if err != nil {
			return nil, err
		}

		cost, err := decimal.NewFromString(dbMenuItem.Cost)
		if err != nil {
			return nil, err
		}

		menuItem := &models.MenuItem{
			ID:          dbMenuItem.ID.String(),
			Name:        dbMenuItem.Name,
			CategoryID:  dbMenuItem.CategoryID.String(),
			IsAvailable: dbMenuItem.IsAvailable,
			CreatedAt:   dbMenuItem.CreatedAt,
			UpdatedAt:   dbMenuItem.UpdatedAt,
			Price:       types.DecimalText(price),
			Cost:        types.DecimalText(cost),
		}

		if dbMenuItem.Description.Valid {
			description := dbMenuItem.Description.String
			menuItem.Description = &description
		}

		menuItems = append(menuItems, menuItem)
	}

	return menuItems, nil
}

// ListMenuItemsByCategory retrieves a list of menu items by category
func (r *menuRepo) ListMenuItemsByCategory(categoryID string, limit, offset int) ([]*models.MenuItem, error) {
	categoryUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, err
	}

	dbMenuItems, err := r.queries.ListMenuItemsByCategory(context.Background(), db.ListMenuItemsByCategoryParams{
		CategoryID: categoryUUID,
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return nil, err
	}

	var menuItems []*models.MenuItem
	for _, dbMenuItem := range dbMenuItems {
		price, err := decimal.NewFromString(dbMenuItem.Price)
		if err != nil {
			return nil, err
		}

		cost, err := decimal.NewFromString(dbMenuItem.Cost)
		if err != nil {
			return nil, err
		}

		menuItem := &models.MenuItem{
			ID:          dbMenuItem.ID.String(),
			Name:        dbMenuItem.Name,
			CategoryID:  dbMenuItem.CategoryID.String(),
			IsAvailable: dbMenuItem.IsAvailable,
			CreatedAt:   dbMenuItem.CreatedAt,
			UpdatedAt:   dbMenuItem.UpdatedAt,
			Price:       types.DecimalText(price),
			Cost:        types.DecimalText(cost),
		}

		if dbMenuItem.Description.Valid {
			description := dbMenuItem.Description.String
			menuItem.Description = &description
		}

		menuItems = append(menuItems, menuItem)
	}

	return menuItems, nil
}

// CreateMenuItem creates a new menu item
func (r *menuRepo) CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error) {
	categoryID, err := uuid.Parse(item.CategoryID)
	if err != nil {
		return nil, err
	}

	var description sql.NullString
	if item.Description != nil {
		description = sql.NullString{
			String: *item.Description,
			Valid:  true,
		}
	}

	dbMenuItem, err := r.queries.CreateMenuItem(context.Background(), db.CreateMenuItemParams{
		Name:        item.Name,
		CategoryID:  categoryID,
		Description: description,
		Price:       item.Price.String(),
		Cost:        item.Cost.String(),
	})
	if err != nil {
		return nil, err
	}

	createdMenuItem := &models.MenuItem{
		ID:          dbMenuItem.ID.String(),
		Name:        dbMenuItem.Name,
		CategoryID:  dbMenuItem.CategoryID.String(),
		IsAvailable: dbMenuItem.IsAvailable,
		CreatedAt:   dbMenuItem.CreatedAt,
		UpdatedAt:   dbMenuItem.UpdatedAt,
		Price:       types.DecimalText(decimal.RequireFromString(dbMenuItem.Price)),
		Cost:        types.DecimalText(decimal.RequireFromString(dbMenuItem.Cost)),
	}

	if dbMenuItem.Description.Valid {
		description := dbMenuItem.Description.String
		createdMenuItem.Description = &description
	}

	return createdMenuItem, nil
}

// UpdateMenuItem updates an existing menu item
func (r *menuRepo) UpdateMenuItem(item *models.MenuItem) (*models.MenuItem, error) {
	itemID, err := uuid.Parse(item.ID)
	if err != nil {
		return nil, err
	}

	categoryID, err := uuid.Parse(item.CategoryID)
	if err != nil {
		return nil, err
	}

	var description sql.NullString
	if item.Description != nil {
		description = sql.NullString{
			String: *item.Description,
			Valid:  true,
		}
	} else {
		description = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	dbMenuItem, err := r.queries.UpdateMenuItem(context.Background(), db.UpdateMenuItemParams{
		ID:          itemID,
		Name:        item.Name,
		CategoryID:  categoryID,
		Description: description,
		Price:       item.Price.String(),
		Cost:        item.Cost.String(),
		IsAvailable: item.IsAvailable,
	})
	if err != nil {
		return nil, err
	}

	updatedMenuItem := &models.MenuItem{
		ID:          dbMenuItem.ID.String(),
		Name:        dbMenuItem.Name,
		CategoryID:  dbMenuItem.CategoryID.String(),
		IsAvailable: dbMenuItem.IsAvailable,
		CreatedAt:   dbMenuItem.CreatedAt,
		UpdatedAt:   dbMenuItem.UpdatedAt,
		Price:       types.DecimalText(decimal.RequireFromString(dbMenuItem.Price)),
		Cost:        types.DecimalText(decimal.RequireFromString(dbMenuItem.Cost)),
	}

	if dbMenuItem.Description.Valid {
		descriptionStr := dbMenuItem.Description.String
		updatedMenuItem.Description = &descriptionStr
	}

	return updatedMenuItem, nil
}

// DeleteMenuItem deletes a menu item by ID
func (r *menuRepo) DeleteMenuItem(id string) error {
	itemID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = r.queries.DeleteMenuItem(context.Background(), itemID)
	if err != nil {
		return err
	}

	return nil
}

// NotImplementedError returns a standard error for unimplemented methods
func NotImplementedError() error {
	return errors.New("method not implemented")
}
