package menurepo

import (
	"github.com/jackc/pgx/v5"
	database_postgres "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/database/postgres"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
)

type repository struct {
	pool database_postgres.Pool
}

func New(pool database_postgres.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func scanMenuItem(row pgx.Row) (menu.MenuItem, error) {
	var r record
	err := row.Scan(&r.ID, &r.RestaurantID, &r.CategoryID, &r.CategoryName,
		&r.Name, &r.Description, &r.Price, &r.Available, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return menu.MenuItem{}, errs.Internal("scan menu item", err)
	}
	return toModel(r)
}
