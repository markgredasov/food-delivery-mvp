package menurepo

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

func scanMenuItem(row pgx.Row) (menu.MenuItem, error) {
	var r record
	err := row.Scan(&r.ID, &r.RestaurantID, &r.CategoryID, &r.CategoryName,
		&r.Name, &r.Description, &r.Price, &r.Available, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return menu.MenuItem{}, errs.Internal("scan menu item", err)
	}
	return toModel(r)
}

func scanMenuItems(rows pgx.Rows) ([]menu.MenuItem, error) {
	var out []menu.MenuItem
	for rows.Next() {
		item, err := scanMenuItem(rows)
		if err != nil {
			return nil, errs.Internal("scan menu item", err)
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, errs.Internal("iterate menu items", err)
	}

	return out, nil
}
