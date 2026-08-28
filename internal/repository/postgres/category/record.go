package categoryrepo

import (
	"time"

	"github.com/google/uuid"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
)

type record struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func toModel(r record) (menu.Category, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return menu.Category{}, errs.Internal("parse category id", err)
	}
	return menu.NewCategory(id, r.Name)
}
