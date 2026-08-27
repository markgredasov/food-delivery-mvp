package catalogservice

import (
	"context"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/restaurant"
)

type service struct {
	restaurant RestaurantRepository
	menu       MenuRepository
}

type RestaurantRepository interface {
	ListActive(ctx context.Context) ([]restaurant.Restaurant, error)
	GetByID(ctx context.Context, restaurantID string) (restaurant.Restaurant, error)
}

type MenuRepository interface {
	ListByRestaurantID(ctx context.Context, restaurantID string) ([]menu.MenuItem, error)
}

func New(restaurant RestaurantRepository, menu MenuRepository) *service {
	return &service{
		restaurant: restaurant,
		menu:       menu,
	}
}
