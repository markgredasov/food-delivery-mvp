package catalogservice

import (
	"context"

	"github.com/markgredasov/food-delivery-mvp/internal/model/menu"
)

func (s *service) ListCategories(ctx context.Context) ([]menu.Category, error) {
	return s.category.List(ctx)
}
