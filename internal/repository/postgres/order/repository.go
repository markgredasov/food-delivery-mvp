package orderrepo

import (
	"github.com/jackc/pgx/v5"
	database_postgres "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/database/postgres"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/order"
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
