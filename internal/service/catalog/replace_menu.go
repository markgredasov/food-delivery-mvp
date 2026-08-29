package catalogservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
)

func (s *service) ReplaceMenu(ctx context.Context, restaurantID uuid.UUID, items []menu.MenuItem) ([]menu.MenuItem, error) {
	if _, err := s.restaurant.GetByID(ctx, restaurantID.String()); err != nil {
		return nil, err
	}

	for i, it := range items {
		if it.ID == uuid.Nil {
			items[i].ID = uuid.New()
		}
	}

	var out []menu.MenuItem
	err := s.tx.WithTx(ctx, func(ctx context.Context) error {
		var err error
		out, err = s.menu.ReplaceMenu(ctx, restaurantID, items)
		return err
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
