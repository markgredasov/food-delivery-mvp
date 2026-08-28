package orderservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/order"
)

func (s *service) ListRestaurantOrders(ctx context.Context, restaurantID uuid.UUID) ([]order.Order, error) {
	if _, err := s.restaurants.GetByID(ctx, restaurantID.String()); err != nil {
		return nil, err
	}
	return s.orders.ListActiveByRestaurant(ctx, restaurantID)
}
