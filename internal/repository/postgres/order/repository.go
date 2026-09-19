package orderrepo

import (
	"github.com/jackc/pgx/v5"
	database_postgres "github.com/markgredasov/food-delivery-mvp/internal/database/postgres"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/order"
)

type repository struct {
	pool database_postgres.Pool
}

func New(pool database_postgres.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func scanOrder(row pgx.Row) (order.Order, error) {
	var r record
	err := row.Scan(&r.ID, &r.RestaurantID, &r.UserID, &r.DeliveryAddress, &r.Status, &r.TotalAmount, &r.Comment, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return order.Order{}, errs.Internal("scan order", err)
	}
	return toModel(r)
}
