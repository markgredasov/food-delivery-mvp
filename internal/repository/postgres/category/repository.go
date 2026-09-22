package categoryrepo

import (
	"github.com/jackc/pgx/v5"
	database_postgres "github.com/markgredasov/food-delivery-mvp/internal/database/postgres"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/menu"
)

type Repository struct {
	pool database_postgres.Pool
}

func New(pool database_postgres.Pool) *Repository {
	return &Repository{
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
