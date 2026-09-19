package restaurantrepo

import (
	"time"

	"github.com/google/uuid"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/address"
	"github.com/markgredasov/food-delivery-mvp/internal/model/restaurant"
)

type record struct {
	ID          string
	Name        string
	Description *string
	Address     address.Address
	Status      string
	ServiceURL  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func toModel(r record) (restaurant.Restaurant, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return restaurant.Restaurant{}, errs.InvalidArgument(err.Error())
	}

	status, err := restaurant.NewStatus(r.Status)
	if err != nil {
		return restaurant.Restaurant{}, err
	}

	return restaurant.Restaurant{
		ID:          id,
		Name:        r.Name,
		Description: r.Description,
		Address:     r.Address,
		Status:      status,
		ServiceURL:  r.ServiceURL,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}, nil
}
