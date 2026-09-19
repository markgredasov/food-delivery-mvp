package menurepo

import (
	"context"

	"github.com/google/uuid"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/menu"
)

func (r *repository) GetManyByIDs(ctx context.Context, ids []uuid.UUID) ([]menu.MenuItem, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		SELECT m.id, m.restaurant_id, m.category_id, c.name category_name, 
			m.name, m.description, m.price, m.available, 
			m.created_at, m.updated_at
		FROM kitchen.menu_items m
		LEFT JOIN kitchen.categories c 
			ON m.category_id = c.id
		WHERE m.id = ANY($1);
	`

	rows, err := r.pool.GetQuerier(ctx).Query(ctx, sqlQuery, ids)
	if err != nil {
		return nil, errs.Internal("get menu items", err)
	}
	defer rows.Close()

	return scanMenuItems(rows)
}
