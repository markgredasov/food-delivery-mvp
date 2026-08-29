package orderservice

import (
	"context"

	"github.com/google/uuid"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/logger"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/order"
	"go.uber.org/zap"
)

func (s *service) AcceptOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID) (order.Order, error) {
	return s.transitionOwnedOrder(ctx, restaurantID, orderID, func(o *order.Order) error {
		return o.Accept()
	})
}

func (s *service) RejectOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID, reason string) (order.Order, error) {
	log := logger.FromContext(ctx)
	log.Info("rejected order", zap.String("reason", reason), zap.String("orderID", orderID.String()))

	return s.transitionOwnedOrder(ctx, restaurantID, orderID, func(o *order.Order) error {
		return o.Reject()
	})
}

func (s *service) transitionOwnedOrder(ctx context.Context, restaurantID, orderID uuid.UUID, mutate func(o *order.Order) error) (order.Order, error) {
	o, err := s.orders.GetByID(ctx, orderID)
	if err != nil {
		return order.Order{}, err
	}
	if !o.IsOwnedBy(restaurantID) {
		return order.Order{}, errs.Forbidden("order does not belong to this restaurant")
	}

	from := o.Status
	if err = mutate(&o); err != nil {
		return order.Order{}, err
	}
	if err = s.orders.UpdateStatus(ctx, o.ID, from, o.Status); err != nil {
		return order.Order{}, err
	}
	return o, nil
}
