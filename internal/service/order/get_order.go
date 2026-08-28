package orderservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/order"
)

func (s *service) GetOrder(ctx context.Context, id uuid.UUID) (order.Order, error) {
	return s.orders.GetByID(ctx, id)
}
