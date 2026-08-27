package catalogservice

import (
	"context"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/restaurant"
)

type service struct {
	restaurants RestaurantRepository
}

type RestaurantRepository interface {
	ListActive(ctx context.Context) ([]restaurant.Restaurant, error)
	GetByID(ctx context.Context, restaurantID string) (restaurant.Restaurant, error)
}

func New(restaurants RestaurantRepository) *service {
	return &service{
		restaurants: restaurants,
	}
}
