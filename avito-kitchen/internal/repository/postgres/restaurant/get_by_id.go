package restaurantrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/restaurant"
)

func (r *repository) GetByID(ctx context.Context, restaurantID string) (restaurant.Restaurant, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		SELECT id, name, description, address, status, service_url, created_at, updated_at
		FROM kitchen.restaurants
		WHERE id = $1;
	`

	row := r.pool.GetQuerier(ctx).QueryRow(ctx, sqlQuery, restaurantID)
	rest, err := scanRestaurant(row)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return restaurant.Restaurant{}, errs.NotFound("restaurant not found")
		}
		return restaurant.Restaurant{}, errs.Internal("get restaurant", err)
	}

	return rest, nil
}
