package menurepo

import (
	"context"

	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
)

func (r *repository) ListByRestaurantID(ctx context.Context, restaurantID string) ([]menu.MenuItem, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		SELECT m.id, m.restaurant_id, m.category_id, c.name category_name, 
			m.name, m.description, m.price, m.available, 
			m.created_at, m.updated_at
		FROM kitchen.menu_items m
		LEFT JOIN kitchen.categories c 
			ON m.category_id = c.id
		WHERE restaurant_id = $1;
	`

	rows, err := r.pool.GetQuerier(ctx).Query(ctx, sqlQuery, restaurantID)
	if err != nil {
		return nil, errs.Internal("get menu", err)
	}
	defer rows.Close()

	return scanMenuItems(rows)
}
