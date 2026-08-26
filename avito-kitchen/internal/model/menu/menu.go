package menu

import (
	"fmt"
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

// MenuItemStatus are conditions of a MenuItem (available/unavailable).
type MenuItemStatus string

// Statuses of MenuItem.
const (
	MenuItemAvailable   MenuItemStatus = "available"
	MenuItemUnavailable MenuItemStatus = "unavailable"
)

// MenuItem is a single dish/product belonging to a restaurant.
type MenuItem struct {
	ID           uuid.UUID
	RestaurantID uuid.UUID
	Category     *Category
	Name         string
	Description  *string
	Price        money.Money
	Status       MenuItemStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewMenuItem validates and builds new MenuItem.
func NewMenuItem(id, restaurantID uuid.UUID, category *Category, name string,
	description *string, price money.Money, available bool) (MenuItem, error) {
	if name == "" {
		return MenuItem{}, errs.InvalidArgument("menu item name must not be empty")
	}

	var status MenuItemStatus
	if available {
		status = MenuItemAvailable
	} else {
		status = MenuItemUnavailable
	}

	return MenuItem{
		ID:           id,
		RestaurantID: restaurantID,
		Category:     category,
		Name:         name,
		Description:  description,
		Price:        price,
		Status:       status,
	}, nil
}

// EnsureAvailable returns a domain error if the item cannot currently be
// ordered.
func (m MenuItem) EnsureAvailable() error {
	if m.Status == MenuItemUnavailable {
		return errs.Conflict(fmt.Sprintf("menu item \"%s\" is not available", m.Name))
	}
	return nil
}

// BelongsTo reports whether the item belongs to the given restaurant.
func (m MenuItem) BelongsTo(restaurantID uuid.UUID) bool {
	return m.RestaurantID == restaurantID
}

// Menu is a list of products belonging to restaurant.
type Menu struct {
	RestaurantID uuid.UUID
	Items        []MenuItem
}

// NewMenu validates and builds new Menu from MenuItems.
func NewMenu(items []MenuItem) (Menu, error) {
	if len(items) == 0 {
		return Menu{}, errs.InvalidArgument("no items")
	}

	return Menu{
		RestaurantID: items[0].RestaurantID,
		Items:        items,
	}, nil
}
