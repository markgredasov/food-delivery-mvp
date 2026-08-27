package catalogservice

import (
	"context"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/restaurant"
)

func (s *service) GetRestaurant(ctx context.Context, restaurantID string) (restaurant.Restaurant, error) {
	return s.restaurant.GetByID(ctx, restaurantID)
}
