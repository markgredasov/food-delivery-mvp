package restaurantrepo

import (
	"context"

	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/restaurant"
)

func (r *repository) ListActive(ctx context.Context) ([]restaurant.Restaurant, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		SELECT id, name, description, address, status, service_url, created_at, updated_at
		FROM kitchen.restaurants
		WHERE status = 'active'
		ORDER BY name;
	`

	rows, err := r.pool.GetQuerier(ctx).Query(ctx, sqlQuery)
	if err != nil {
		return nil, errs.Internal("list active restaurants", err)
	}
	defer rows.Close()

	var out []restaurant.Restaurant
	for rows.Next() {
		//nolint:govet // rest is new so no shadows
		rest, err := scanRestaurant(rows)
		if err != nil {
			return nil, errs.Internal("scan restaurant", err)
		}
		out = append(out, rest)
	}
	if err = rows.Err(); err != nil {
		return nil, errs.Internal("iterate restaurants", err)
	}

	return out, nil
}
