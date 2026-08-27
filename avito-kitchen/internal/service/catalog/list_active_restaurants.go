package catalogservice

import (
	"context"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/restaurant"
)

func (s *service) ListActiveRestaurants(ctx context.Context) ([]restaurant.Restaurant, error) {
	return s.restaurants.ListActive(ctx)
}
