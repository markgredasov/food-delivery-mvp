package menurepo

import (
	"github.com/google/uuid"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/money"
)

type record struct {
	ID           string
	RestaurantID string
	CategoryID   string
	CategoryName string
	Name         string
	Description  *string
	Price        money.Money
	Available    bool
}

func toModel(r record) (menu.MenuItem, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return menu.MenuItem{}, errs.InvalidArgument(err.Error())
	}

	restaurantID, err := uuid.Parse(r.RestaurantID)
	if err != nil {
		return menu.MenuItem{}, errs.InvalidArgument(err.Error())
	}

	categoryID, err := uuid.Parse(r.CategoryID)
	if err != nil {
		return menu.MenuItem{}, errs.InvalidArgument(err.Error())
	}

	category, err := menu.NewCategory(categoryID, r.CategoryName)
	if err != nil {
		return menu.MenuItem{}, err
	}

	return menu.NewMenuItem(id, restaurantID, &category, r.Name, r.Description, r.Price, r.Available)
}
