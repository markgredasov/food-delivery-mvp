package catalogservice

import (
	"context"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
)

func (s *service) ListCategories(ctx context.Context) ([]menu.Category, error) {
	return s.category.List(ctx)
}
