package categoryrepo

import (
	"context"

	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
)

func (r *repository) List(ctx context.Context) ([]menu.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		SELECT id, name, created_at, updated_at
		FROM kitchen.categories
		ORDER BY name;
	`

	rows, err := r.pool.GetQuerier(ctx).Query(ctx, sqlQuery)
	if err != nil {
		return nil, errs.Internal("list categories", err)
	}
	defer rows.Close()

	var out []menu.Category
	for rows.Next() {
		var c menu.Category
		c, err = scanCategory(rows)
		if err != nil {
			return nil, errs.Internal("scan categories", err)
		}
		out = append(out, c)
	}
	if err = rows.Err(); err != nil {
		return nil, errs.Internal("iterate categories", err)
	}

	return out, nil
}
