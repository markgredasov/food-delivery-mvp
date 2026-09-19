package categoryrepo

import (
	"time"

	"github.com/google/uuid"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/menu"
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
