package orderservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/markgredasov/food-delivery-mvp/internal/model/order"
)

func (s *service) ListRestaurantOrders(ctx context.Context, restaurantID uuid.UUID) ([]order.Order, error) {
	if _, err := s.restaurants.GetByID(ctx, restaurantID.String()); err != nil {
		return nil, err
	}
	return s.orders.ListActiveByRestaurant(ctx, restaurantID)
}
