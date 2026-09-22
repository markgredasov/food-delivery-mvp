package catalogservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/markgredasov/food-delivery-mvp/internal/model/restaurant"
)

func (s *Service) GetRestaurant(ctx context.Context, restaurantID uuid.UUID) (restaurant.Restaurant, error) {
	return s.restaurant.GetByID(ctx, restaurantID.String())
}
