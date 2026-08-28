package catalogservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
)

func (s *service) GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]menu.MenuItem, error) {
	if _, err := s.restaurant.GetByID(ctx, restaurantID.String()); err != nil {
		return nil, err
	}

	return s.menu.ListByRestaurantID(ctx, restaurantID.String())
}
