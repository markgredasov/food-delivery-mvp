package categoryrepo

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

func scanCategory(rows pgx.Rows) (menu.Category, error) {
	var r record
	if err := rows.Scan(&r.ID, &r.Name, &r.CreatedAt, &r.UpdatedAt); err != nil {
		return menu.Category{}, errs.Internal("scan category", err)
	}

	return toModel(r)
}
