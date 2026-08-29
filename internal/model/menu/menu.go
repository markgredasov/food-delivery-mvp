package menu

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/money"
)

// Category groups menu items.
type Category struct {
	ID   uuid.UUID
	Name string
}

// NewCategory validates and builds new Category.
func NewCategory(id uuid.UUID, name string) (Category, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return Category{}, errs.InvalidArgument("category name must not be empty")
	}
	return Category{
		ID:   id,
		Name: trimmedName,
	}, nil
}

// MenuItem is a single dish/product belonging to a restaurant.
type MenuItem struct {
	ID           uuid.UUID
	RestaurantID uuid.UUID
	Category     *Category
	Name         string
	Description  *string
	Price        money.Money
	Available    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewMenuItem validates and builds new MenuItem.
func NewMenuItem(id, restaurantID uuid.UUID, category *Category, name string,
	description *string, price money.Money, available bool) (MenuItem, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return MenuItem{}, errs.InvalidArgument("menu item name must not be empty")
	}

	if price.IsZero() {
		return MenuItem{}, errs.InvalidArgument("price must be greater than zero")
	}

	return MenuItem{
		ID:           id,
		RestaurantID: restaurantID,
		Category:     category,
		Name:         trimmedName,
		Description:  description,
		Price:        price,
		Available:    available,
	}, nil
}

// EnsureAvailable returns a domain error if the item cannot currently be
// ordered.
func (m MenuItem) EnsureAvailable() error {
	if !m.Available {
		return errs.Conflict(fmt.Sprintf("menu item \"%s\" is not available", m.Name))
	}
	return nil
}

// BelongsTo reports whether the item belongs to the given restaurant.
func (m MenuItem) BelongsTo(restaurantID uuid.UUID) bool {
	return m.RestaurantID == restaurantID
}
