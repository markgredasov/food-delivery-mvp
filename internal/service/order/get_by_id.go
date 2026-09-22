package orderservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/markgredasov/food-delivery-mvp/internal/model/order"
)

func (s *Service) GetOrder(ctx context.Context, id uuid.UUID) (order.Order, error) {
	return s.orders.GetByID(ctx, id)
}
