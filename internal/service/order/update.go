package orderservice

import (
	"context"

	"github.com/google/uuid"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/logger"
	"github.com/markgredasov/food-delivery-mvp/internal/model/order"
	"go.uber.org/zap"
)

func (s *Service) AcceptOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID) (order.Order, error) {
	return s.transitionOwnedOrder(ctx, restaurantID, orderID, func(o *order.Order) error {
		return o.Accept()
	})
}

func (s *Service) RejectOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID, reason string) (order.Order, error) {
	log := logger.FromContext(ctx)
	log.Info("rejected order", zap.String("reason", reason), zap.String("orderID", orderID.String()))

	return s.transitionOwnedOrder(ctx, restaurantID, orderID, func(o *order.Order) error {
		return o.Reject()
	})
}

func (s *Service) UpdateStatus(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID, statusString string) (order.Order, error) {
	next, err := order.NewStatus(statusString)
	if err != nil {
		return order.Order{}, err
	}

	return s.transitionOwnedOrder(ctx, restaurantID, orderID, func(o *order.Order) error {
		return o.AdvanceTo(next)
	})
}

func (s *Service) transitionOwnedOrder(ctx context.Context, restaurantID, orderID uuid.UUID, mutate func(o *order.Order) error) (order.Order, error) {
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
