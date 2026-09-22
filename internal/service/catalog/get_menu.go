package catalogservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/markgredasov/food-delivery-mvp/internal/model/menu"
)

func (s *Service) GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]menu.MenuItem, error) {
	if _, err := s.restaurant.GetByID(ctx, restaurantID.String()); err != nil {
		return nil, err
	}

	return s.menu.ListByRestaurantID(ctx, restaurantID.String())
}
