package menurepo

import (
	"context"

	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/menu"
)

func (r *Repository) ListByRestaurantID(ctx context.Context, restaurantID string) ([]menu.MenuItem, error) {
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
