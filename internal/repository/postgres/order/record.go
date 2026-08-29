package orderrepo

import (
	"time"

	"github.com/google/uuid"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/address"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/money"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/order"
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
