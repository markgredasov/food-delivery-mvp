package catalogservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/markgredasov/food-delivery-mvp/internal/model/menu"
	"github.com/markgredasov/food-delivery-mvp/internal/model/restaurant"
)

type service struct {
	restaurant RestaurantRepository
	menu       MenuRepository
	category   CategoryRepository
	tx         TxManager
}

type RestaurantRepository interface {
	ListActive(ctx context.Context) ([]restaurant.Restaurant, error)
	GetByID(ctx context.Context, restaurantID string) (restaurant.Restaurant, error)
}

type MenuRepository interface {
	ListByRestaurantID(ctx context.Context, restaurantID string) ([]menu.MenuItem, error)
	ReplaceMenu(ctx context.Context, restaurantID uuid.UUID, items []menu.MenuItem) ([]menu.MenuItem, error)
}

type CategoryRepository interface {
	List(ctx context.Context) ([]menu.Category, error)
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

func New(restaurant RestaurantRepository, menu MenuRepository, category CategoryRepository, tx TxManager) *service {
	return &service{
		restaurant: restaurant,
		menu:       menu,
		category:   category,
		tx:         tx,
	}
}
