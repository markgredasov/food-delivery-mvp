package orderservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/markgredasov/food-delivery-mvp/internal/model/menu"
	"github.com/markgredasov/food-delivery-mvp/internal/model/order"
	"github.com/markgredasov/food-delivery-mvp/internal/model/restaurant"
	"github.com/markgredasov/food-delivery-mvp/internal/service/webhook"
)

type Service struct {
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
	ListActiveByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]order.Order, error)
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type WebhookSender interface {
	Send(ctx context.Context, webhookURL string, payload webhook.OrderPayload) error
}

func New(restaurants RestaurantRepository, menu MenuRepository, orders OrderRepository,
	txManager TxManager, webhookSender WebhookSender) *Service {
	return &Service{
		restaurants: restaurants,
		menuItems:   menu,
		orders:      orders,
		tx:          txManager,
		webhook:     webhookSender,
	}
}
