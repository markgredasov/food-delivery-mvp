package menurepo

import (
	"time"

	"github.com/google/uuid"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/money"
)

type record struct {
	ID           string
	RestaurantID string
	CategoryID   *string
	CategoryName *string
	Name         string
	Description  *string
	Price        money.Money
	Available    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

	var categoryID uuid.UUID
	if r.CategoryID != nil {
		//nolint:govet // cID is new so no shadows
		cID, err := uuid.Parse(*r.CategoryID)
		if err != nil {
			return menu.MenuItem{}, errs.InvalidArgument(err.Error())
		}
		categoryID = cID
	}

	var category *menu.Category
	if r.CategoryName != nil {
		//nolint:govet // c is new so no shadows
		c, err := menu.NewCategory(categoryID, *r.CategoryName)
		if err != nil {
			return menu.MenuItem{}, err
		}
		category = &c
	}

	return menu.NewMenuItem(id, restaurantID, category, r.Name, r.Description, r.Price, r.Available)
}
