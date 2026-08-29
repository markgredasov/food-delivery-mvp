package menurepo

import (
	"context"

	"github.com/google/uuid"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
)

func (r *repository) ReplaceMenu(ctx context.Context, restaurantID uuid.UUID, newItems []menu.MenuItem) ([]menu.MenuItem, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	q := r.pool.GetQuerier(ctx)

	sqlQueryInsert := `
		INSERT INTO menu_items (id, restaurant_id, category_id, name, description, price, available)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			category_id = EXCLUDED.category_id,
			name        = EXCLUDED.name,
			description = EXCLUDED.description,
			price       = EXCLUDED.price,
			available   = EXCLUDED.available
		WHERE menu_items.restaurant_id = $2;
	`

	keepIDs := make([]uuid.UUID, 0, len(newItems))
	for _, it := range newItems {
		var categoryID *uuid.UUID
		if it.Category != nil {
			categoryID = &it.Category.ID
		}
		if _, err := q.Exec(ctx, sqlQueryInsert, it.ID, restaurantID, categoryID, it.Name,
			it.Description, it.Price, it.Available); err != nil {
			return nil, errs.Internal("replace menu item", err)
		}
		keepIDs = append(keepIDs, it.ID)
	}

	sqlQueryUpdate := `
		UPDATE menu_items SET available = FALSE WHERE restaurant_id = $1 AND NOT (id = ANY($2));
	`

	if _, err := q.Exec(ctx, sqlQueryUpdate, restaurantID, keepIDs); err != nil {
		return nil, errs.Internal("retire stale menu items", err)
	}

	return r.ListByRestaurantID(ctx, restaurantID.String())
}
