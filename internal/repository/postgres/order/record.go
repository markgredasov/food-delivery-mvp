package orderrepo

import (
	"time"

	"github.com/google/uuid"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/address"
	"github.com/markgredasov/food-delivery-mvp/internal/model/money"
	"github.com/markgredasov/food-delivery-mvp/internal/model/order"
)

type record struct {
	ID              string
	RestaurantID    string
	UserID          *string
	DeliveryAddress address.Address
	Status          string
	TotalAmount     money.Money
	Comment         *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func toModel(r record) (order.Order, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return order.Order{}, errs.InvalidArgument(err.Error())
	}

	restaurantID, err := uuid.Parse(r.RestaurantID)
	if err != nil {
		return order.Order{}, errs.InvalidArgument(err.Error())
	}

	var userID *uuid.UUID
	if r.UserID != nil {
		uid, err := uuid.Parse(*r.UserID) //nolint:govet // uid is new so no shadows
		if err != nil {
			return order.Order{}, errs.InvalidArgument(err.Error())
		}
		userID = &uid
	}

	status := order.Status(r.Status)
	if !status.Valid() {
		return order.Order{}, errs.InvalidArgument(r.Status)
	}

	return order.Order{
		ID:              id,
		RestaurantID:    restaurantID,
		UserID:          userID,
		Status:          status,
		DeliveryAddress: r.DeliveryAddress,
		Comment:         r.Comment,
		TotalAmount:     r.TotalAmount,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}, nil
}
