package orderservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/order"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/restaurant"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/service/webhook"
)

type service struct {
	restaurants RestaurantRepository
	menuItems   MenuRepository
	orders      OrderRepository
	tx          TxManager
	webhook     WebhookSender
}

type RestaurantRepository interface {
	GetByID(ctx context.Context, id string) (restaurant.Restaurant, error)
}

type MenuRepository interface {
	GetManyByIDs(ctx context.Context, ids []uuid.UUID) ([]menu.MenuItem, error)
}

type OrderRepository interface {
	Create(ctx context.Context, o order.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (order.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, from, to order.Status) error
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type WebhookSender interface {
	Send(ctx context.Context, webhookURL string, payload webhook.OrderPayload) error
}

func New(restaurants RestaurantRepository, menu MenuRepository, orders OrderRepository,
	txManager TxManager, webhookSender WebhookSender) *service {
	return &service{
		restaurants: restaurants,
		menuItems:   menu,
		orders:      orders,
		tx:          txManager,
		webhook:     webhookSender,
	}
}
