package catalogservice

import (
	"context"

	"github.com/markgredasov/food-delivery-mvp/internal/model/restaurant"
)

func (s *Service) ListActiveRestaurants(ctx context.Context) ([]restaurant.Restaurant, error) {
	return s.restaurant.ListActive(ctx)
}
