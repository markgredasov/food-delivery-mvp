package orderservice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/logger"
	"github.com/markgredasov/food-delivery-mvp/internal/model/menu"
	"github.com/markgredasov/food-delivery-mvp/internal/model/order"
	"github.com/markgredasov/food-delivery-mvp/internal/model/restaurant"
	"github.com/markgredasov/food-delivery-mvp/internal/service/webhook"
	"go.uber.org/zap"
)

func (s *service) CreateOrder(ctx context.Context, in order.Order) (order.Order, error) {
	var created order.Order
	err := s.tx.WithTx(ctx, func(ctx context.Context) error {
		rest, err := s.restaurants.GetByID(ctx, in.RestaurantID.String())
		if err != nil {
			return err
		}
		if err = rest.EnsureAcceptsOrders(); err != nil {
			return err
		}

		items, err := s.resolveOrderItems(ctx, rest, in.Items)
		if err != nil {
			return err
		}

		o, err := order.New(in.ID, rest.ID, in.UserID, in.DeliveryAddress, items, in.Comment)
		if err != nil {
			return err
		}
		if err = s.orders.Create(ctx, o); err != nil {
			return err
		}
		created = o
		return nil
	})
	if err != nil {
		return order.Order{}, err
	}

	s.pushToRestaurant(ctx, created)

	fresh, err := s.orders.GetByID(ctx, created.RestaurantID)
	if err != nil {
		return created, nil //nolint:nilerr // the order was created successfully; a refresh failure is not fatal to the caller
	}
	return fresh, nil
}

func (s *service) resolveOrderItems(ctx context.Context, rest restaurant.Restaurant, requested []order.OrderItem) ([]order.OrderItem, error) {
	ids := make([]uuid.UUID, len(requested))
	for i, r := range requested {
		if r.Quantity <= 0 {
			return nil, errs.InvalidArgument("item quantity must be positive")
		}
		ids[i] = r.MenuItemID
	}

	found, err := s.menuItems.GetManyByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	byID := make(map[uuid.UUID]menu.MenuItem, len(found))
	for _, item := range found {
		byID[item.ID] = item
	}

	items := make([]order.OrderItem, 0, len(requested))
	for _, r := range requested {
		item, ok := byID[r.MenuItemID]
		if !ok {
			return nil, errs.NotFound(fmt.Sprintf("menu item %s not found", r.MenuItemID.String()))
		}
		if !item.BelongsTo(rest.ID) {
			return nil, errs.InvalidArgument(fmt.Sprintf("menu item %s does not belong to this restaurant", r.MenuItemID.String()))
		}
		if err = item.EnsureAvailable(); err != nil {
			return nil, err
		}
		items = append(items, order.OrderItem{
			ID:         uuid.New(),
			Name:       item.Name,
			MenuItemID: item.ID,
			Quantity:   r.Quantity,
			Price:      item.Price,
		})
	}
	return items, nil
}

func (s *service) pushToRestaurant(ctx context.Context, o order.Order) {
	log := logger.FromContext(ctx)

	rest, err := s.restaurants.GetByID(ctx, o.RestaurantID.String())
	if err != nil || !rest.HasServiceURL() {
		return
	}

	payload := webhook.OrderPayload{
		OrderID:         o.ID,
		RestaurantID:    o.RestaurantID,
		DeliveryAddress: o.DeliveryAddress,
		TotalAmount:     o.TotalAmount.String(),
		CreatedAt:       time.Now(),
		Items:           make([]webhook.OrderItem, 0, len(o.Items)),
	}
	for _, it := range o.Items {
		payload.Items = append(payload.Items, webhook.OrderItem{
			MenuItemID: it.MenuItemID,
			Quantity:   it.Quantity,
			Price:      it.Price.String(),
		})
	}

	if err = s.webhook.Send(ctx, rest.ServiceURL, payload); err != nil {
		log.Warn("webhook delivery failed, order stays pending for polling fallback",
			zap.String("order_id", o.ID.String()), zap.Error(err))
		return
	}

	if err = s.orders.UpdateStatus(ctx, o.ID, order.StatusPending, order.StatusSentToRestaurant); err != nil {
		log.Warn("failed to record confirmed after successful webhook delivery",
			zap.String("order_id", o.ID.String()), zap.Error(err))
	}
}
