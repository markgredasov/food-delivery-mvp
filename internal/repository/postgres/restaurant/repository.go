package restaurantrepo

import (
	"github.com/jackc/pgx/v5"
	database_postgres "github.com/markgredasov/food-delivery-mvp/internal/database/postgres"
	"github.com/markgredasov/food-delivery-mvp/internal/model/restaurant"
)

type repository struct {
	pool database_postgres.Pool
}

func New(pool database_postgres.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func scanRestaurant(row pgx.Row) (restaurant.Restaurant, error) {
	var r record
	err := row.Scan(&r.ID, &r.Name, &r.Description, &r.Address, &r.Status, &r.ServiceURL, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return restaurant.Restaurant{}, err
	}
	return toModel(r)
}
